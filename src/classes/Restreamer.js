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
const publicIp = require('public-ip');
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
        let nginx = config.nginx.streaming;
        let token = process.env.RS_TOKEN || config.auth.token;

        let hlspath = 'rtmp://' + nginx.ip + ':' + nginx.rtmp_port + nginx.rtmp_hls_path + 'live.stream';

        if(token != '') {
            hlspath += '?token=' + token;
        }

        return hlspath;
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
        if(Restreamer.data.states.repeatToLocalNginx.type != 'connected') {
            return;
        }

        let command = new FfmpegCommand(Restreamer.generateOutputHLSPath());

        command.output(Restreamer.generateSnapshotPath());
        command.outputOption(config.ffmpeg.options.snapshot);
        command.on('error', (error) => {
            logger.error('Error fetching snapshot: ' + error.toString().trim());
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
     * @param {string} streamType
     */
    static stopStream(streamType) {
        Restreamer.updateState(streamType, 'stopped');
        logger.info('Stopping ' + streamType);

        if(Restreamer.data.processes[streamType] !== null) {
            Restreamer.data.processes[streamType].kill('SIGKILL');
            Restreamer.data.processes[streamType] = null;
        }
    }

    /**
     * restore the ffmpeg processes from jsondb (called on app start to restore ffmpeg processes
     * after the application has been killed or stuff
     */
    static restoreFFMpegProcesses () {
        var db = new JsonDB(config.jsondb, true, false);

        Restreamer.data.addresses = db.getData('/addresses');
        Restreamer.data.states = db.getData('/states');
        Restreamer.data.options = db.getData('/options');
        Restreamer.data.userActions = db.getData('/userActions');

        let state = '';

        state = Restreamer.getState('repeatToLocalNginx');
        let repeatToLocalNginxReconnecting = (state == 'connected' || state == 'connecting');

        state = Restreamer.getState('repeatToOptionalOutput');
        let repeatToOptionalOutputReconnecting = (state == 'connected' || state == 'connecting');

        // check if a stream was repeated locally
        if(Restreamer.data.addresses.srcAddress && repeatToLocalNginxReconnecting) {
            Restreamer.startStream(
                Restreamer.data.addresses.srcAddress,
                'repeatToLocalNginx',
                null,
                true
            );
        }
        else {
            Restreamer.updateState('repeatToLocalNginx', 'disconnected')
        }

        // check if the stream was repeated to an output address
        if(Restreamer.data.addresses.optionalOutputAddress && repeatToOptionalOutputReconnecting) {
            Restreamer.startStream(
                Restreamer.data.addresses.srcAddress,
                'repeatToOptionalOutput',
                Restreamer.data.addresses.optionalOutputAddress,
                true
            );
        }
        else {
            Restreamer.updateState('repeatToOptionalOutput', 'disconnected')
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
     * append the ffmpeg options of the config file to an output
     * @param {FfmpegCommand} ffmpegCommand
     * @return {Promise}
     */
    static appendOutputOptionFromConfig(ffmpegCommand, streamType) {
        var deferred = Q.defer();
        var ffmpegOptions = [];

        let state = Restreamer.getState(streamType);
        if(state != 'connecting') {
            logger.debug('Skipping "startStream" because state is not "connecting". Current state is "' + state + '".', streamType);
            return null;
        }

        ffmpegCommand.ffprobe((err, data) => {
            if(err) {
                let lines = err.toString().split(/\r\n|\r|\n/);
                lines = lines.filter(function (line) {
                    return line.length > 0;
                });

                return deferred.reject(lines[lines.length - 1]);
            }

            let video = null;
            let audio = null;

            for(let s of data.streams) {
                if(s.codec_type == 'video' && video === null) {
                    video = s;
                }
                else if(s.codec_type == 'audio' && audio === null) {
                    audio = s;
                }

                if(video !== null && audio !== null) {
                    break;
                }
            }

            if(video === null) {
                return deferred.reject("no video stream detected");
            }

            if(video.codec_name != 'h264') {
                return deferred.reject("video stream must be h264, found " + video.codec_name);
            }

            let option = "native_h264_no_audio";

            if(streamType == 'repeatToLocalNginx') {
                if(audio !== null) {
                    switch(audio.codec_name) {  // consider all allowed audio codecs for FLV
                        case 'mp3':
                        case 'pcm_alaw':
                        case 'pcm_mulaw':
                            option = "native_h264_native_audio"; break;
                        case 'aac':
                            option = "native_h264_native_aac"; break;
                        default:
                            option = "native_h264_transcode_aac"; break;
                    }
                }

                if(process.env.RS_AUDIO == "none") {
                    option = "native_h264_no_audio";
                }
                else if(process.env.RS_AUDIO == "silence") {
                    option = "native_h264_silence_aac";
                }
                else if(process.env.RS_AUDIO == "aac") {
                    if(audio !== null && audio.codec_name != 'aac') {
                        option = "native_h264_transcode_aac";
                    }
                }
                else if(process.env.RS_AUDIO == "mp3") {
                    if(audio !== null && audio.codec_name != 'mp3') {
                        option = "native_h264_transcode_mp3";
                    }
                }
            }
            else {
                option = "native_h264_native_audio";
            }

            logger.debug('Selected ffmpeg.option: ' + option);
            ffmpegOptions = config.ffmpeg.options[option];

            for(let o of ffmpegOptions) {
                ffmpegCommand.outputOptions(o);
            }

            return deferred.resolve();
        });

        return deferred.promise;
    }

    /**
     * update the state of the stream
     * @param {string} streamType
     * @param {string} state
     * @param {string} message
     * @return {string} name of the new state
     */
    static updateState(streamType, state, message) {
        let previousState = Restreamer.setState(streamType, state, message);

        if(previousState == state) {
            return state;
        }

        logger.debug('Update state from "' + previousState + '" to "' + state + '"', streamType);

        if(streamType === 'repeatToLocalNginx' && state === 'connected') {
            Restreamer.fetchSnapshot();
        }

        Restreamer.writeToDB();
        Restreamer.updateStreamDataOnGui();

        return state;
    }

    static setState(streamType, state, message) {
        let previousState = Restreamer.data.states[streamType].type;

        Restreamer.data.states[streamType] = {
            'type': state,
            'message': message
        };

        return previousState;
    }

    /**
     * get the state of the stream
     * @param {string} streamType
     * @return {string} name of the new state
     */
    static getState(streamType) {
        return Restreamer.data.states[streamType].type;
    }

    /**
     * get the current user action
     * @param {string} streamType
     * @return {string} name of the user action
     */
    static getUserAction(streamType) {
        return Restreamer.data.userActions[streamType];
    }

    /**
     * set the current user action
     * @param {string} streamType
     * @param {string} action
     * @return {string} name of the previous user action
     */
    static setUserAction(streamType, action) {
        let previousAction = Restreamer.data.userActions[streamType];
        Restreamer.data.userActions[streamType] = action;
        return previousAction;
    }

    /**
     * update the last submitted user action (like click on stop stream, click on start stream)
     * @param {string} streamType
     * @param {string} action user action
     * @return {string} name of the new user action
     */
    static updateUserAction(streamType, action) {
        let previousAction = Restreamer.setUserAction(streamType, action);

        if(previousAction == action) {
            return action;
        }

        logger.debug('Set user action from "' + previousAction + '" to "' + action +'"', streamType);

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
     * @param {string} optionalOutput address of the optional output
     * @param {string} force force the start of the stream
     */
    static startStream(src, streamType, optionalOutput, force) {
        // remove any running timeouts
        Restreamer.setTimeout(streamType, 'retry', null);
        Restreamer.setTimeout(streamType, 'stale', null);

        force = force || false;
        if(force != true) {
            // check if there's currently no other stream connected or connecting
            let state = Restreamer.getState(streamType);
            if(state == 'connected' || state == 'connecting') {
                logger.debug('Skipping "startStream" because state is "connected" or "connecting". Current state is "' + state + '".', streamType);
                return;
            }
        } 

        // check if the user has clicked 'stop' meanwhile, so the startStream process has to be skipped
        if(Restreamer.getUserAction(streamType) == 'stop') {
            Restreamer.stopStream(streamType);
            logger.debug('Skipping "startStream" because "stop" has been clicked', streamType);
            return;
        }

        // update the state and src address on the frontend
        Restreamer.data.addresses.srcAddress = src;
        Restreamer.updateState(streamType, 'connecting');

        /** @var {FfmpegCommand} instance of FFmpeg command of module fluent-ffmpeg */
        var command = null;

        /** @var {Promise} Promise to make sure, the add output process has been finished */
        var addOutputPromise = null;

        logger.info('Start stream', streamType);

        let rtmpPath = Restreamer.generateOutputHLSPath();

        if(streamType == 'repeatToLocalNginx') {
            // repeat to local nginx server
            command = new FfmpegCommand(src, {
                stdoutLines: 1
            });

            // add outputs to the ffmpeg stream
            command.output(rtmpPath);
            addOutputPromise = Restreamer.appendOutputOptionFromConfig(command, streamType)
        }
        else {
            // repeat to optional output
            command = new FfmpegCommand(rtmpPath, {
                stdoutLines: 1
            });

            Restreamer.data.addresses.optionalOutputAddress = optionalOutput;

            // add outputs to the ffmpeg stream
            command.output(optionalOutput);
            addOutputPromise = Restreamer.appendOutputOptionFromConfig(command, streamType)
        }

        if(addOutputPromise === null) {
            logger.debug('Skipping "startStream" because promise is null', streamType);
            return;
        }

        let retry = () => {
            logger.info('Retry to connect to "' + src + '" in ' + config.ffmpeg.monitor.restart_wait + 'ms', streamType);
            Restreamer.setTimeout(streamType, 'retry', () => {
                logger.info('Retry to connect to "' + src + '"', streamType);

                if(Restreamer.data.userActions[streamType] == 'stop') {
                    logger.debug('Skipping retry because "stop" has been clicked', streamType);
                    Restreamer.updateState(streamType, 'disconnected');
                    return;
                }

                Restreamer.startStream(src, streamType, optionalOutput);
            }, config.ffmpeg.monitor.restart_wait);
        }

        // after adding outputs, define events on the new FFmpeg stream
        addOutputPromise.then(() => {
            let ffmpegOptions = config.ffmpeg.options.input;

            for(let o of ffmpegOptions) {
                command.inputOptions(o);
            }

            // GUI option
            if(streamType == 'repeatToLocalNginx') {
                if(Restreamer.data.options.rtspTcp && Restreamer.data.addresses.srcAddress.indexOf('rtsp') == 0) {
                    command.inputOptions('-rtsp_transport', 'tcp');
                }
            }

            command
                .on('start', (commandLine) => {
                    Restreamer.data.processes[streamType] = command;

                    if(Restreamer.data.userActions[streamType] == 'stop') {
                        Restreamer.stopStream(streamType);
                        logger.debug('Skipping on "start" event of FFmpeg command because "stop" has been clicked', streamType);
                        return;
                    }

                    logger.debug('FFmpeg spawned: ' + commandLine);
                })

                // stream ended
                .on('end', () => {
                    Restreamer.data.processes[streamType] = null;

                    Restreamer.setTimeout(streamType, 'retry', null);
                    Restreamer.setTimeout(streamType, 'stale', null);

                    Restreamer.updateState(streamType, 'disconnected');

                    if(Restreamer.data.userActions[streamType] == 'stop') {
                        logger.debug('Skipping retry because "stop" has been clicked', streamType);
                        return;
                    }

                    logger.info(streamType + ': ended normally');

                    retry();
                })

                // stream error handler
                .on('error', (error, stdout, stderr) => {
                    Restreamer.data.processes[streamType] = null;

                    Restreamer.setTimeout(streamType, 'retry', null);
                    Restreamer.setTimeout(streamType, 'stale', null);

                    if(Restreamer.data.userActions[streamType] == 'stop') {
                        Restreamer.updateState(streamType, 'disconnected');
                        logger.debug('Skipping retry since "stop" has been clicked', streamType);
                        return;
                    }

                    Restreamer.updateState(streamType, 'error', error.toString());

                    retry();
                })

                // progress handler
                .on('progress', (progress) => {
                    if(Restreamer.data.states[streamType].type == 'connecting') {
                        Restreamer.updateState(streamType, 'connected');
                    }

                    Restreamer.data.progresses[streamType] = progress;
                    Restreamer.updateProgressOnGui();

                    // add a stale timeout
                    Restreamer.setTimeout(streamType, 'stale', () => {
                        logger.info('Stale connection', streamType);
                        Restreamer.stopStream(streamType);
                    }, config.ffmpeg.monitor.stale_wait);
                });

            command.exec();
        }).catch((error) => {
            logger.debug('Failed to probe stream: ' + error.toString(), streamType);

            if(Restreamer.data.userActions[streamType] == 'stop') {
                Restreamer.updateState(streamType, 'disconnected');
                logger.debug('Skipping retry since "stop" has been clicked', streamType);
                return;
            }

            Restreamer.updateState(streamType, 'error', error.toString());

            retry();
        });
    }

    /**
     * set a timeout
     * @param {string} streamType Either 'repeatToLocalNginx' or 'repeatToOptionalOutput'
     * @param {string} target Kind of timeout, either 'retry' or 'stale'
     * @param {function} func Callback function
     * @param {int} delay Delay for the timeout
     * @return {void}
     */
    static setTimeout(streamType, target, func, delay) {
        clearTimeout(Restreamer.data.timeouts[target][streamType]);

        if(typeof func == 'function') {
            Restreamer.data.timeouts[target][streamType] = setTimeout(func, delay);
        }
    }

    /**
     * bind websocket events on application start
     */
    static bindWebsocketEvents () {
        WebsocketsController.setConnectCallback((socket) => {
            socket.emit('publicIp', Restreamer.data.publicIp);
            socket.on('startStream', (options) => {
                logger.debug('Got "startStream" event for ' + options.streamType);
                Restreamer.updateUserAction(options.streamType, 'start');
                Restreamer.updateOptions(options.options);
                Restreamer.startStream(options.src, options.streamType, options.optionalOutput);
            });
            socket.on('stopStream', (streamType) => {
                logger.debug('Got "stopStream" event for ' + streamType);
                Restreamer.updateUserAction(streamType, 'stop');
                Restreamer.stopStream(streamType);
            });
            socket.on('checkStates', () => {
                logger.debug('Got "checkStates" event');
                Restreamer.updateStreamDataOnGui()
            });
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
            'states': Restreamer.data.states
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
        logger.info('Getting public IP ...', 'start.publicip');
        publicIp.v4().then(ip => {
            Restreamer.data.publicIp = ip;
            logger.info('Found public IP: ' + ip, 'start.publicip');
        }).catch(err => {
            logger.warn('Failed to get public IP', 'start.publicip');
            Restreamer.data.publicIp = '127.0.0.1';
        });
    }
}

/*
 define data structure of Restreamer Data
 */
Restreamer.data = {
    'timeouts': {
        'retry': {
            'repeatToLocalNginx': null,
            'repeatToOptionalOutput': null
        },
        'stale': {
            'repeatToLocalNginx': null,
            'repeatToOptionalOutput': null
        }
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
        'repeatToLocalNginx': null,
        'repeatToOptionalOutput': null
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
