/**
 * @file holds the code for the class Restreamer
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */

const config = require('../../conf/live.json');
const logger = require('./Logger')('Restreamer');
const WebsocketsController = require('./WebsocketsController');
const FfmpegCommand = require('fluent-ffmpeg');
const Q = require('q');

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

        return 'rtmp://' + nginx.ip + ':' + nginx.rtmp_port + nginx.rtmp_hls_path + 'live.stream';
    }

    /**
     * generate output snapshot-path from config-file
     * @returns {string}
     */
    static generateSnapshotPath () {
        const path = require('path');

        return path.join(__dirname, '..', 'webserver', 'public', 'images', 'live.jpg');
    }

    /**
     * receive snapshot by using first frame of repeated video
     * @param {boolean} firstSnapshot
     */
    static fetchSnapshot (firstSnapshot) {
        if (Restreamer.data.states.repeatToLocalNginx.type === 'connected' || firstSnapshot) {
            var command = new FfmpegCommand(Restreamer.generateOutputHLSPath());

            command.output(Restreamer.generateSnapshotPath());
            command.outputOption(config.ffmpeg.options.snapshot);
            command.on('error', (error)=> {
                logger.error('Error on fetching snapshot: ' + error.toString());
            });
            command.on('end', () =>{
                logger.info('updated snapshot');
                Q.delay(process.env.SNAPSHOT_REFRESH_INTERVAL).then(function () {
                    Restreamer.fetchSnapshot(false);
                });
            });
            command.exec();
        }
    }

    /**
     * stop stream
     * @param {string} processName
     */
    static stopStream (processName) {
        var processHasBeenSpawned = typeof Restreamer.data.processes[processName].kill === 'function';

        Restreamer.updateState(processName, 'stopped');
        logger.info('stopStream ' + processName);
        if (processHasBeenSpawned) {
            Restreamer.data.processes[processName].kill();
            Restreamer.data.processes[processName] = {
                state: 'not_connected'
            };
        }
    }

    /**
     * restore the ffmpeg processes from jsondb (called on app start to restore ffmpeg processes
     * after the application has been killed or stuff
     */
    static restoreFFMpegProcesses () {
        var JsonDB = require('node-json-db');
        var db = new JsonDB(config.jsondb, true, false);

        try {
            Restreamer.data.addresses = db.getData('/addresses');
            Restreamer.data.states = db.getData('/states');
            Restreamer.data.userActions = db.getData('/userActions');

            // check if the srcAddress has been repeated to Local Nginx
            var repeatToLocalNginxIsConnected = Restreamer.data.states.repeatToLocalNginx.type === 'connected';
            var repeatToLocalNginxIsConnecting = Restreamer.data.states.repeatToLocalNginx.type === 'connecting';
            var repeatToOptionalOutputIsConnected = Restreamer.data.states.repeatToOptionalOutput.type === 'connected';
            var repeatToOptionalOutputIsConnecting = Restreamer.data.states.repeatToOptionalOutput.type === 'connecting';

            if (Restreamer.data.addresses.srcAddress && (!!repeatToLocalNginxIsConnected || !!repeatToLocalNginxIsConnecting)) {
                Restreamer.startStream(Restreamer.data.addresses.srcAddress, 'repeatToLocalNginx');
            }
            if (Restreamer.data.addresses.optionalOutputAddress && (!!repeatToOptionalOutputIsConnected || !!repeatToOptionalOutputIsConnecting)) {
                Restreamer.startStream(Restreamer.data.addresses.srcAddress, 'repeatToOptionalOutput', Restreamer.data.addresses.optionalOutputAddress);
            }
        }
        catch(error) {
            logger.error('error restoring ffmpeg process: ' + error);
        }
    }

    /**
     * write JSON file for persistence
     */
    static writeToDB () {
        var JsonDB = require('node-json-db');
        var db = new JsonDB(config.jsondb, true, false);

        db.push('/', Restreamer.dataForJsonDb());
    }

    /**
     * send websocket event to GUI to update the state of the streams
     */
    static updateStreamDataOnGui () {
        WebsocketsController.emitToNamespace('/', 'updateStreamData', Restreamer.extractDataOfStreams());
    }

    /**
     * send websocket event to GUI to update the state of the streams
     */
    static updateProgressOnGui () {
        WebsocketsController.emitToNamespace('/', 'updateProgress', Restreamer.data.progresses);
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

    /**
     * append the ffmpeg options of the config file to an output
     * @param {FfmpegCommand} ffmpegCommand
     * @return {promise}
     */
    static appendOutputOptionFromConfig (ffmpegCommand) {
        var deferred = Q.defer();

        ffmpegCommand.ffprobe(function (err, data) {
            if (err) {
                return deferred.reject(err);
            } else {
                var ffmpegOptions;

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
            }
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
            type: state,
            message: message
        };
        Restreamer.writeToDB();
        Restreamer.updateStreamDataOnGui();
        return state;
    }

    /**
     * update the last submitted useraction (like click on stop stream, click on start stream)
     * @param {string} processName
     * @param {string} action useraction
     * @return {string} name of the new user action
     */
    static updateUserAction (processName, action) {
        logger.debug(`Set useraction of "${processName}" from "${Restreamer.data.userActions[processName]}" to "${action}"`);
        Restreamer.data.userActions[processName] = action;
        Restreamer.writeToDB();
        Restreamer.updateStreamDataOnGui();
        return action;
    }

    /**
     *
     * @param {string} src src-address of the ffmpeg stream
     * @param {string} streamType repeatToOptionalOutput or repeatToLocalNginx
     * @param {string=} optionalOutput address of the optional output
     * @param {number=} retryCounter current value of the retry counter (startStream retries automatically if anything fails)
     */
    static startStream (src, streamType, optionalOutput, retryCounter) {
        logger.info(`Start stream "${streamType}"`);

        // update the retry counter for the streamType
        Restreamer.data.retryCounter[streamType].current = typeof retryCounter === 'undefined' ? 0 : retryCounter;

        /** @var {FfmpegCommand} instance of FFmpeg command of module fluent-ffmpeg*/
        var command;

        /** @var {Promise} Promise to make sure, the add output process has been finished */
        var addOutputPromise;

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
                outputLineLimit: 1
            });

            // add outputs to the ffmpeg stream
            addOutputPromise = Restreamer.addOutput(command, Restreamer.generateOutputHLSPath())
            .catch(function (error) {
                logger.error(`Error adding one or more outputs: ${error.toString}`);
            });

        // repeat to optional output
        } else if (repeatToOptionalOutput) {
            command = new FfmpegCommand(Restreamer.generateOutputHLSPath(), {
                outputLineLimit: 1
            });
            Restreamer.data.addresses.optionalOutputAddress = optionalOutput;
            addOutputPromise = Restreamer.addOutput(command, optionalOutput);
        }

        // after adding outputs, define events on the new FFmpeg stream
        addOutputPromise.then(function () {
            command

                /*
                 * STREAMING STARTED
                 */
                .on('start', function (commandLine) {
                    if (Restreamer.data.userActions[streamType] === 'stop') {
                        logger.debug('Skipping on "start" event of FFmpeg command since "stopped" has been clicked');
                        return;
                    } else {
                        logger.debug(`FFmpeg spawned: ${commandLine}`);
                        Restreamer.data.processes[streamType] = command;

                        // fetch snapshot only, if repeated to local nginx
                        if (repeatToLocalNginx) {
                            Restreamer.fetchSnapshot(true);
                        }
                    }
                })

                /*
                 * ERROR HANDLER
                 */
                .on('error', (error)=>{
                    if (error.toString().indexOf('SIGKILL') > -1) {
                        Restreamer.updateState(streamType, 'disconnected');
                        logger.info(`FFmpeg streaming stopped for "${streamType}"`);
                    } else {
                        Restreamer.data.retryCounter[streamType].current++;
                        if (Restreamer.data.retryCounter[streamType].current <= config.ffmpeg.monitor.retries) {
                            if (Restreamer.data.userActions[streamType] === 'stop') {
                                logger.debug('Skipping retry since "stopped" has been clicked');
                                return;
                            } else {
                                Restreamer.updateState(streamType, 'error', error.toString());
                                logger.info(`Retrying FFmpeg connection to "${src}" after "${config.ffmpeg.monitor.restart_wait} ms`);
                                Q.delay(config.ffmpeg.monitor.restart_wait).then(function () {
                                    logger.info(`Retry FFmpeg connection to "${src}" retrycounter: ${Restreamer.data.retryCounter[streamType].current}`);
                                    Restreamer.startStream(src, streamType, optionalOutput, Restreamer.data.retryCounter[streamType].current);
                                });
                            }
                        } else {
                            Restreamer.updateState(streamType, 'error', error.toString());
                        }
                    }
                })

                /*
                 STREAMING ENDED
                 */
                .on('end', ()=>{
                    Restreamer.updateState(streamType, 'disconnected');
                    Restreamer.data.retryCounter[streamType].current++;
                    if (Restreamer.data.retryCounter[streamType].current <= config.ffmpeg.monitor.retries) {
                        if (Restreamer.data.userActions[streamType] === 'stop') {
                            logger.debug('Skipping retry since "stopped" has been clicked');
                            return;
                        } else {
                            logger.info(`Retrying FFmpeg connection to "${src}" after "${config.ffmpeg.monitor.restart_wait}" ms`);
                            Q.delay(config.ffmpeg.monitor.restart_wait).then(function () {
                                logger.info(`Retry FFmpeg connection to "${src}" retrycounter: ${Restreamer.data.retryCounter[streamType].current}`);
                                Restreamer.startStream(src, streamType, optionalOutput, Restreamer.data.retryCounter[streamType].current);
                            });
                        }
                    }
                });

            /*
             * PROGRESS
             */

            var progressMethod = function (progress) {
                if (Restreamer.data.states[streamType].type === 'connecting') {
                    Restreamer.updateState(streamType, 'connected');
                }
                Restreamer.data.progresses[streamType] = progress;
                Restreamer.updateProgressOnGui();
                command.removeAllListeners('progress');
            };

            command.on('progress', progressMethod);
            setInterval(()=>{
                if (command.listeners('progress').length === 0) {
                    command.on('progress', progressMethod);
                }
            }, 1000);
            command.exec();
        }).catch(function (error) {
            logger.error(`Error starting FFmpeg command: ${error.toString()}`);
        });
    }

    /**
     * binded on app-start to bind websocketevents of Restreamer
     */
    static bindWebsocketEvents () {
        WebsocketsController.addOnConnectionEventToNamespace('/', function (socket) {
            socket.on('startStream', (options)=> {
                Restreamer.updateUserAction(options.streamType, 'start');
                Restreamer.startStream(options.src, options.streamType, options.optionalOutput);
            });
            socket.on('stopStream', (streamType)=>{
                Restreamer.updateUserAction(streamType, 'stop');
                Restreamer.stopStream(streamType);
            });
            socket.on('checkForAppUpdates', ()=>{
                const app = require('../webserver/app');

                socket.emit('checkForAppUpdatesResult', app.get('updateAvailable'));
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
        var sData = {};
        sData.userActions = Restreamer.data.userActions;
        sData.addresses = Restreamer.data.addresses;
        sData.states = Restreamer.data.states;
        sData.retryCounter = Restreamer.data.retryCounter;
        return sData;
    }

    /**
     * create with only the data, that is needed by the jsonDb
     * @return {object}
     */
    static dataForJsonDb () {
        var dbData = {};
        dbData.addresses = Restreamer.data.addresses;
        dbData.userActions = Restreamer.data.userActions;
        dbData.states = Restreamer.data.states;
        return dbData;
    }
}

/*
define data structure of Restreamer Data
 */
Restreamer.data = {
    retryCounter: {
        repeatToLocalNginx: {
            current: 0,
            max: config.ffmpeg.monitor.retries
        },
        repeatToOptionalOutput: {
            current: 0,
            max: config.ffmpeg.monitor.retries
        }
    },
    states: {
        repeatToLocalNginx: {
            type: 'disconnected',
            message: ''
        },
        repeatToOptionalOutput: {
            type: 'disconnected',
            message: ''
        }
    },
    userActions: {
        repeatToLocalNginx: 'start',
        repeatToOptionalOutput: 'start'
    },
    processes: {
        // overwritten with ffmpeg process if stream has been started
        repeatToLocalNginx: {
            state: 'not_connected'
        },
        // overwritten with ffmpeg process if stream has been started
        repeatToOptionalOutput: {
            state: 'not_connected'
        }
    },
    progresses: {
        // overwritten with ffmpeg process if stream has been started
        repeatToLocalNginx: {},
        // overwritten with ffmpeg process if stream has been started
        repeatToOptionalOutput: {}
    },
    addresses: {
        srcAddress: '',
        optionalOutputAddress: ''
    }
};

module.exports = Restreamer;
