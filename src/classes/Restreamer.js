/**
 * @file holds the code for the class Restreamer
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
'use strict';

const path = require('path');
const config = require(path.join(global.__base, 'conf', 'live.json'));
const logger = require('./Logger')('Restreamer');
const WebsocketsController = require('./WebsocketsController');
const FfmpegCommand = require('fluent-ffmpeg');
const Q = require('q');
const JsonDB = require('node-json-db');
const exec = require('child_process').exec;
const packageJson = require(path.join(global.__base, 'package.json'));
const https = require('https');

/**
 * class Restreamer creates and manages streams through ffmpeg
 */
class Restreamer {

    /**
     * generate output hls-path from config-file
     * @returns {string}
     */
    static generateOutputHLSPath () {
        var nginx = config.nginx.streaming;
        var token = process.env.RS_TOKEN || config.auth.token;

        if (token != '') {
            return 'rtmp://' + nginx.ip + ':' + nginx.rtmp_port + nginx.rtmp_hls_path + 'live.stream' + '?token=' + token;
        } else {
            return 'rtmp://' + nginx.ip + ':' + nginx.rtmp_port + nginx.rtmp_hls_path + 'live.stream';
        }
    }

    /**
     * generate output snapshot-path from config-file
     * @returns {string}
     */
    static generateSnapshotPath () {
        return path.join(global.__public, 'images', 'live.jpg');
    }

    /**
     * receive snapshot by using first frame of repeated video
     */
    static fetchSnapshot () {
        var command = null;
        if (Restreamer.data.states.repeatToLocalNginx.type === 'connected') {
            command = new FfmpegCommand(Restreamer.generateOutputHLSPath());

            command.output(Restreamer.generateSnapshotPath());
            command.outputOption(config.ffmpeg.options.snapshot);
            command.on('error', (error)=> {
                logger.error('Error on fetching snapshot: ' + error.toString());
            });
            command.on('end', () => {
                logger.info('Updated snapshot');
                WebsocketsController.emit('snapshot', null);
                Q.delay(this.calculateSnapshotRefreshInterval()).then(() => {
                    Restreamer.fetchSnapshot();
                });
            });
            command.exec();
        }
    }

    static calculateSnapshotRefreshInterval () {
        let fallbackRefreshInterval = 60000;
        let snapshotRefreshInterval = process.env.RS_SNAPSHOT_REFRESH_INTERVAL.match(/([0-9]+)([a-z]{1,2})?/);

        let refreshInterval = fallbackRefreshInterval;

        if (typeof snapshotRefreshInterval[2] === 'undefined' || snapshotRefreshInterval[2] === 'ms') {
            refreshInterval = snapshotRefreshInterval[1];
        } else if (snapshotRefreshInterval[2] === 's') {
            refreshInterval = snapshotRefreshInterval[1] * 1000;
        } else if (snapshotRefreshInterval[2] === 'm') {
            refreshInterval = snapshotRefreshInterval[1] * 1000 * 60;
        }

        if (refreshInterval >= 10000) {
            return refreshInterval;
        }

        return fallbackRefreshInterval;
    }

    /**
     * stop stream
     * @param {string} processName
     */
    static stopStream (processName) {
        var processHasBeenSpawned = typeof Restreamer.data.processes[processName].kill === 'function';

        Restreamer.updateState(processName, 'stopped');
        logger.info('stopStream ' + processName);
        clearTimeout(Restreamer.data.retryTimeouts[processName]);
        if (processHasBeenSpawned) {
            Restreamer.data.processes[processName].kill();
            Restreamer.data.processes[processName] = {
                'state': 'not_connected'
            };
        }
    }

    /**
     * restore the ffmpeg processes from jsondb (called on app start to restore ffmpeg processes
     * after the application has been killed or stuff
     */
    static restoreFFMpegProcesses () {
        var db = new JsonDB(config.jsondb, true, false);
        var repeatToLocalNginxIsConnected = false;
        var repeatToLocalNginxIsConnecting = false;
        var repeatToOptionalOutputIsConnected = false;
        var repeatToOptionalOutputIsConnecting = false;

        Restreamer.data.addresses = db.getData('/addresses');
        Restreamer.data.states = db.getData('/states');
        Restreamer.data.options = db.getData('/options');
        Restreamer.data.userActions = db.getData('/userActions');

        repeatToLocalNginxIsConnected = Restreamer.data.states.repeatToLocalNginx.type === 'connected';
        repeatToLocalNginxIsConnecting = Restreamer.data.states.repeatToLocalNginx.type === 'connecting';
        repeatToOptionalOutputIsConnected = Restreamer.data.states.repeatToOptionalOutput.type === 'connected';
        repeatToOptionalOutputIsConnecting = Restreamer.data.states.repeatToOptionalOutput.type === 'connecting';

        // check if the srcAddress has been repeated to Local Nginx
        if (Restreamer.data.addresses.srcAddress &&
            (repeatToLocalNginxIsConnected || repeatToLocalNginxIsConnecting)) {
            Restreamer.startStream(
                Restreamer.data.addresses.srcAddress,
                'repeatToLocalNginx'
            );
        }
        if (Restreamer.data.addresses.optionalOutputAddress &&
            (repeatToOptionalOutputIsConnected || repeatToOptionalOutputIsConnecting)) {
            Restreamer.startStream(
                Restreamer.data.addresses.srcAddress,
                'repeatToOptionalOutput',
                Restreamer.data.addresses.optionalOutputAddress
            );
        }
    }

    /**
     * write JSON file for persistence
     */
    static writeToDB () {
        var db = new JsonDB(config.jsondb, true, false);

        db.push('/', Restreamer.dataForJsonDb());
    }

    /**
     * send websocket event to GUI to update the state of the streams
     */
    static updateStreamDataOnGui () {
        WebsocketsController.emit('updateStreamData', Restreamer.extractDataOfStreams());
    }

    /**
     * send websocket event to GUI to update the state of the streams
     */
    static updateProgressOnGui () {
        WebsocketsController.emit('updateProgress', Restreamer.data.progresses);
    }

    /**
     * add output to ffmpeg command
     * @param {FfmpegCommand} ffmpegCommand
     * @param {string} outputAddress
     * @return {Promise}
     */
    static addOutput (ffmpegCommand, outputAddress) {
        ffmpegCommand.output(outputAddress);
        return Restreamer.appendOutputOptionFromConfig(ffmpegCommand);
    }

    static applyOptions (ffmpegCommand, streamType) {
        var ffmpegOptions = config.ffmpeg.options;

        // reduces process data
        ffmpegCommand.inputOptions('-hide_banner', '-stats', '-loglevel', 'quiet');

        // for unclean rtsp-sources
        ffmpegCommand.inputOptions('-err_detect', ffmpegOptions.err_detect);

        // gui option
        if (streamType === 'repeatToLocalNginx') {
            if (Restreamer.data.options.rtspTcp && Restreamer.data.addresses.srcAddress.indexOf('rtsp') === 0) {
                ffmpegCommand.inputOptions('-rtsp_transport', 'tcp');
            }
        }

        // read input at native frame rate (-re)
        ffmpegCommand.native();
    }

    /**
     * append the ffmpeg options of the config file to an output
     * @param {FfmpegCommand} ffmpegCommand
     * @return {Promise}
     */
    static appendOutputOptionFromConfig (ffmpegCommand) {
        var deferred = Q.defer();
        var ffmpegOptions = [];

        ffmpegCommand.ffprobe((err, data) => {
            if (err) {
                return deferred.reject(err);
            }

            if (data.streams.length > 1) {
                ffmpegOptions = config.ffmpeg.options.native_h264;
                logger.debug('Selected ffmpeg.option: native_h264');
            } else {
                ffmpegOptions = config.ffmpeg.options.native_h264_soundless_aac;
                logger.debug('Selected ffmpeg.option: native_h264_soundless_aac');
            }
            for (let option of ffmpegOptions) {
                ffmpegCommand.outputOption(option);
            }
            return deferred.resolve();
        });
        return deferred.promise;
    }

    /**
     * update the state of the stream
     * @param {string} processName
     * @param {string} state
     * @param {string=} message
     * @return {string} name of the new state
     */
    static updateState (processName, state, message) {
        logger.debug(`Update state of "${processName}" from state "${Restreamer.data.states[processName].type}" to state "${state}"`);
        Restreamer.data.states[processName] = {
            'type': state,
            'message': message
        };
        if (processName === 'repeatToLocalNginx' && state === 'connected') {
            Restreamer.fetchSnapshot();
        }
        Restreamer.writeToDB();
        Restreamer.updateStreamDataOnGui();
        return state;
    }

    /**
     * update the last submitted user action (like click on stop stream, click on start stream)
     * @param {string} processName
     * @param {string} action user action
     * @return {string} name of the new user action
     */
    static updateUserAction (processName, action) {
        logger.debug(`Set user action of "${processName}" from "${Restreamer.data.userActions[processName]}" to "${action}"`);
        Restreamer.data.userActions[processName] = action;
        Restreamer.writeToDB();
        Restreamer.updateStreamDataOnGui();
        return action;
    }

    /**
     * update options
     * @param {Object} options
     */
    static updateOptions (options) {
        Restreamer.data.options = options;
        Restreamer.writeToDB();
        Restreamer.updateStreamDataOnGui();
    }

    /**
     *
     * @param {string} src src-address of the ffmpeg stream
     * @param {string} streamType repeatToOptionalOutput or repeatToLocalNginx
     * @param {string=} optionalOutput address of the optional output
     * @param {string=} retryCounter current value of the retry counter (startStream retries automatically if anything fails)
     */
    static startStream (src, streamType, optionalOutput, retryCounter) {
        /** @var {FfmpegCommand} instance of FFmpeg command of module fluent-ffmpeg */
        var command = null;

        /** @var {Promise} Promise to make sure, the add output process has been finished */
        var addOutputPromise = null;

        logger.info(`Start stream "${streamType}"`);

        // update the retry counter for the streamType
        Restreamer.data.retryCounter[streamType].current = typeof retryCounter === 'undefined' ? 0 : retryCounter;

        /** @var {Boolean} */
        const repeatToLocalNginx = streamType === 'repeatToLocalNginx';

        /** @var {Boolean} */
        const repeatToOptionalOutput = streamType === 'repeatToOptionalOutput';

        // check if the user has clicked 'stop' meanwhile, so the startStream process has to be skipped
        if (Restreamer.data.userActions[streamType] === 'stop') {
            logger.debug('Skipping "startStream" since "stopped" has been clicked');
            return;
        }

        // update the state and src address on the frontend
        Restreamer.data.addresses.srcAddress = src;
        Restreamer.updateState(streamType, 'connecting');


        // repeat to local nginx server
        if (repeatToLocalNginx) {
            command = new FfmpegCommand(src, {
                'outputLineLimit': 1
            });

            // add outputs to the ffmpeg stream
            addOutputPromise = Restreamer.addOutput(command, Restreamer.generateOutputHLSPath()).catch((error) => {
                logger.error(`Error adding one or more outputs: ${error.toString}`);
            });

            // repeat to optional output
        } else if (repeatToOptionalOutput) {
            command = new FfmpegCommand(Restreamer.generateOutputHLSPath(), {
                'outputLineLimit': 1
            });
            Restreamer.data.addresses.optionalOutputAddress = optionalOutput;
            addOutputPromise = Restreamer.addOutput(command, optionalOutput);
        }

        // after adding outputs, define events on the new FFmpeg stream
        addOutputPromise.then(() => {
            var progressMethod = (progress) => {
                if (Restreamer.data.states[streamType].type === 'connecting') {
                    Restreamer.data.retryCounter[streamType].current = 1;
                    Restreamer.updateState(streamType, 'connected');
                }
                Restreamer.data.progresses[streamType] = progress;
                Restreamer.updateProgressOnGui();
                command.removeAllListeners('progress');
            };
            Restreamer.applyOptions(command, streamType);
            command
                .on('start', (commandLine) => {
                    if (Restreamer.data.userActions[streamType] === 'stop') {
                        logger.debug('Skipping on "start" event of FFmpeg command since "stopped" has been clicked');
                        return;
                    }
                    logger.debug(`FFmpeg spawned: ${commandLine}`);
                    Restreamer.data.processes[streamType] = command;
                })

                // stream ended
                .on('end', ()=> {
                    Restreamer.updateState(streamType, 'disconnected');
                    Restreamer.data.retryCounter[streamType].current++;
                    if (Restreamer.data.retryCounter[streamType].current <= config.ffmpeg.monitor.retries) {
                        if (Restreamer.data.userActions[streamType] === 'stop') {
                            logger.debug('Skipping retry since "stopped" has been clicked');
                            return;
                        }
                        logger.info(`Retrying FFmpeg connection to "${src}" after "${config.ffmpeg.monitor.restart_wait}" ms`);
                        Restreamer.data.retryTimeouts[streamType] = setTimeout(() => {
                            logger.info(`Retry FFmpeg connection to "${src}" retry counter: ${Restreamer.data.retryCounter[streamType].current}`);
                            Restreamer.startStream(src, streamType, optionalOutput, Restreamer.data.retryCounter[streamType].current);
                        }, config.ffmpeg.monitor.restart_wait);
                    }
                })

                // stream error handler
                .on('error', (error) => {
                    if (error.toString().indexOf('SIGKILL') > -1) {
                        Restreamer.updateState(streamType, 'disconnected');
                        logger.info(`FFmpeg streaming stopped for "${streamType}"`);
                    } else {
                        Restreamer.data.retryCounter[streamType].current++;
                        if (Restreamer.data.retryCounter[streamType].current <= config.ffmpeg.monitor.retries) {
                            if (Restreamer.data.userActions[streamType] === 'stop') {
                                logger.debug('Skipping retry since "stopped" has been clicked');
                                return;
                            }
                            Restreamer.updateState(streamType, 'error', error.toString());
                            logger.error(`Error on stream ${streamType}: ${error.toString()}`);
                            logger.info(`Retrying FFmpeg connection to "${src}" after "${config.ffmpeg.monitor.restart_wait}" ms`);
                            Restreamer.data.retryTimeouts[streamType] = setTimeout(() => {
                                logger.info(`Retry FFmpeg connection to "${src}" retry counter: ${Restreamer.data.retryCounter[streamType].current}`);
                                Restreamer.startStream(src, streamType, optionalOutput, Restreamer.data.retryCounter[streamType].current);
                            }, config.ffmpeg.monitor.restart_wait);
                        } else {
                            Restreamer.updateState(streamType, 'error', error.toString());
                        }
                    }
                });

            command.on('progress', progressMethod);
            setInterval(()=> {
                if (command.listeners('progress').length === 0) {
                    command.on('progress', progressMethod);
                }
            }, 1000);
            command.exec();
        }).catch((error) => {
            logger.error(`Error starting FFmpeg command: ${error.toString()}`);
        });
    }

    /**
     * bind websocket events on application start
     */
    static bindWebsocketEvents () {
        WebsocketsController.setConnectCallback((socket) => {
            socket.emit('publicIp', Restreamer.data.publicIp);
            socket.on('startStream', (options)=> {
                Restreamer.updateUserAction(options.streamType, 'start');
                Restreamer.updateOptions(options.options);
                Restreamer.startStream(options.src, options.streamType, options.optionalOutput);
            });
            socket.on('stopStream', (streamType)=> {
                Restreamer.updateUserAction(streamType, 'stop');
                Restreamer.stopStream(streamType);
            });
            socket.on('checkStates', Restreamer.updateStreamDataOnGui);
        });
    }

    /**
     * creates an object of available streaming data like states and addresses to be able
     * to write it to filesystem
     * @returns {object}
     */
    static extractDataOfStreams () {
        return {
            'addresses': Restreamer.data.addresses,
            'options': Restreamer.data.options,
            'userActions': Restreamer.data.userActions,
            'states': Restreamer.data.states,
            'retryCounter': Restreamer.data.retryCounter
        };
    }

    /**
     * create with only the data, that is needed by the jsonDb
     * @return {object}
     */
    static dataForJsonDb () {
        return {
            'addresses': Restreamer.data.addresses,
            'options': Restreamer.data.options,
            'userActions': Restreamer.data.userActions,
            'states': Restreamer.data.states
        };
    }

    /**
     * check for updates
     */
    static checkForUpdates () {
        const url = {'host': 'datarhei.org', 'path': '/apps.json'};
        logger.debug('Checking for updates...');
        https.get(url, (response) => {
            if (response.statusCode === 200) {
                response.on('data', (body) => {
                    var updateCheck = JSON.parse(body);
                    var updateAvailable = require('semver').lt(packageJson.version, updateCheck.restreamer.version);
                    logger.info(`Update checking succeeded. ${updateAvailable ? 'Update' : 'No updates'} available`, 'checkForUpdates');
                    logger.debug(`local: ${packageJson.version}; remote: ${updateCheck.restreamer.version}`, 'checkForUpdates');
                    Restreamer.data.updateAvailable = updateAvailable;
                    WebsocketsController.emit('update', updateAvailable);
                });
            } else {
                logger.warn(`Got ${String(response.statusCode)} status while trying to fetch update info`, 'checkForUpdates');
            }
        }).on('error', () => {
            logger.warn('Failed fetching update info', 'checkForUpdates');
        });
        setTimeout(Restreamer.checkForUpdates, 12 * 3600 * 1000);
    }

    /**
     * get public ip
     */
    static getPublicIp () {
        logger.info('Getting public ip...', 'start.publicip');
        exec('public-ip', (err, stdout, stderr) => {
            if (err) {
                logger.error(err);
            }
            Restreamer.data.publicIp = stdout.split('\n')[0];
        });
    }
}

/*
 define data structure of Restreamer Data
 */
Restreamer.data = {
    'retryCounter': {
        'repeatToLocalNginx': {
            'current': 0,
            'max': config.ffmpeg.monitor.retries
        },
        'repeatToOptionalOutput': {
            'current': 0,
            'max': config.ffmpeg.monitor.retries
        }
    },
    'retryTimeouts': {
        'repeatToLocalNginx': null,
        'repeatToOptionalOutput': null
    },
    'options': {
        'rtspTcp': false
    },
    'states': {
        'repeatToLocalNginx': {
            'type': 'disconnected',
            'message': ''
        },
        'repeatToOptionalOutput': {
            'type': 'disconnected',
            'message': ''
        }
    },
    'userActions': {
        'repeatToLocalNginx': 'start',
        'repeatToOptionalOutput': 'start'
    },
    'processes': {

        // overwritten with ffmpeg process if stream has been started
        'repeatToLocalNginx': {
            'state': 'not_connected'
        },

        // overwritten with ffmpeg process if stream has been started
        'repeatToOptionalOutput': {
            'state': 'not_connected'
        }
    },
    'progresses': {

        // overwritten with ffmpeg process if stream has been started
        'repeatToLocalNginx': {},

        // overwritten with ffmpeg process if stream has been started
        'repeatToOptionalOutput': {}
    },
    'addresses': {
        'srcAddress': '',
        'optionalOutputAddress': ''
    },
    'updateAvailable': false,
    'publicIp': '127.0.0.1'
};

module.exports = Restreamer;
