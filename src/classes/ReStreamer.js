/*
 * @name ReStreamer
 * @namespace https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */

const config = require("../../conf/live.json");
const Logger = require("./Logger");
const logger = new Logger("ReStreamer");
const WebsocketsController = require("./WebsocketController");
const FfmpegCommand = require('fluent-ffmpeg');
const Q = require("q");

/**
 * class ReStreamer creates and manages streams through ffmpeg
 */
class ReStreamer{

    /**
     * generate output rtmp-path from config-file
     * @returns {string}
     */
    static generateOutputRTMPPath (){
        var nginx =  config.nginx.streaming;
        return "rtmp://" + nginx.ip + ":" + nginx.rtmp_port + nginx.rtmp_path + "live.stream";
    }

    /**
     * generate output hls-path from config-file
     * @returns {string}
     */
    static generateOutputHLSPath (){
        var nginx =  config.nginx.streaming;
        return "rtmp://" + nginx.ip + ":" + nginx.rtmp_port + nginx.hls_path + "live.stream";
    }

    /**
     * generate output snapshot-path from config-file
     * @returns {string}
     */
    static generateSnapshotPath (){
        const path = require("path");
        return path.join(__dirname, "..", "webserver", "public", "images", "live.jpg");
    }

    /**
     * receive snapshot by using first frame of repeated video
     */
    static fetchSnapshot(firstSnapshot){
        if (ReStreamer.data.states.repeatToLocalNginx.type === 'connected' || firstSnapshot){
            logger.debug("updating snapshot");
            var command = new FfmpegCommand(ReStreamer.generateOutputRTMPPath());
            command.output(ReStreamer.generateSnapshotPath());
            command.outputOption(config.ffmpeg.options.snapshot);
            command.on('error', (error)=> {
                logger.error("Error on fetching snapshot: " + error.toString());
            });
            command.on("end", () =>{
                logger.info("updated snapshot");
                Q.delay(process.env.SNAPSHOT_REFRESH_INTERVAL).then(function(){
                    ReStreamer.fetchSnapshot(false);
                });
            });
            command.exec();
        }
    }

    /**
     * stop stream
     */
    static stopStream(processName){
        logger.info("stopStream " + processName);
        ReStreamer.updateState(processName, "stopped");
        var processHasBeenSpawned = typeof ReStreamer.data.processes[processName].kill === 'function';
        if(processHasBeenSpawned){
            ReStreamer.data.processes[processName].kill();
            ReStreamer.data.processes[processName] = {
                state: "not_connected"
            };
        }
    }

    /**
     * restore the ffmpeg processes from jsondb (called on app start to restore ffmpeg processes
     * after the applicatoin has been killed or stuff
     */
    static restoreFFMpegProcesses(){
        var JsonDB = require("node-json-db");
        var db = new JsonDB(config.jsondb, true, false);
        try {
            ReStreamer.data.addresses = db.getData("/addresses");
            ReStreamer.data.states = db.getData("/states");
            ReStreamer.data.userActions = db.getData("/userActions");
            //check if the srcAddress has been repeated to Local Nginx
            var repeatToLocalNginxIsConnected = ReStreamer.data.states.repeatToLocalNginx.type === 'connected';
            var repeatToLocalNginxIsConnecting = ReStreamer.data.states.repeatToLocalNginx.type === 'connecting';
            var repeatToOptionalOutputIsConnected = ReStreamer.data.states.repeatToOptionalOutput.type === 'connected';
            var repeatToOptionalOutputIsConnecting = ReStreamer.data.states.repeatToOptionalOutput.type === 'connecting';
            if (ReStreamer.data.addresses.srcAddress && (!!repeatToLocalNginxIsConnected || !!repeatToLocalNginxIsConnecting)) {
                ReStreamer.startStream(ReStreamer.data.addresses.srcAddress, "repeatToLocalNginx");
            }
            if (ReStreamer.data.addresses.optionalOutputAddress && (!!repeatToOptionalOutputIsConnected || !!repeatToOptionalOutputIsConnecting)) {
                ReStreamer.startStream(ReStreamer.data.addresses.srcAddress, "repeatToOptionalOutput", ReStreamer.data.addresses.optionalOutputAddress);
            }
        }
        catch(error){
            logger.error("error restoring ffmpeg process: " + error);
        }
    }

    /**
     * write JSON file for persistency
     */
    static writeToDB(){
        var JsonDB = require("node-json-db");
        var db = new JsonDB(config.jsondb, true, false);
        db.push("/", ReStreamer.extractDataOfStreams());
    }

    /**
     * send websocket event to GUI to update the state of the streams
     */
    static updateStreamDataOnGui(){
        WebsocketsController.emitToNamespace("/", "updateStreamData", ReStreamer.extractDataOfStreams());
    }

    /**
     * send websocket event to GUI to update the state of the streams
     */
    static updateProgressOnGui(){
        WebsocketsController.emitToNamespace("/", "updateProgress", ReStreamer.data.progresses);
    }

    /**
     * add output to ffmpeg command
     * @param {FfmpegCommand} ffmpegCommand
     * @param {string} outputAddress
     * @return {Promise}
     */
    static addOutput(ffmpegCommand, outputAddress){
        ffmpegCommand.output(outputAddress);
        return ReStreamer.appendOutputOptionFromConfig(ffmpegCommand);
    }

    /**
     * append the ffmpeg options of the config file to an output
     * @param {FfmpegCommand} ffmpegCommand
     * @return {Promise}
     */
    static appendOutputOptionFromConfig(ffmpegCommand){
        var deferred = Q.defer();
        ffmpegCommand.ffprobe(function(err, data) {
            if (err){
                return deferred.reject(err);
            }else{
                var ffmpegOptions;
                if(data.streams.length > 1) {
                    ffmpegOptions = config.ffmpeg.options.native_h264;
                    logger.debug("Selected ffmpeg.option: native_h264");
                }else{
                    ffmpegOptions = config.ffmpeg.options.native_h264_soundless_aac;
                    logger.debug("Selected ffmpeg.option: native_h264_soundless_aac");
                }
                for (let option of ffmpegOptions){
                    ffmpegCommand.outputOption(option);
                }
                return deferred.resolve();
            }
        });
        return deferred.promise;
    }

    /**
     *
     * @param {string} processName
     * @param {string} state
     * @param {string} message
     */
    static updateState(processName, state, message){
        logger.debug("Update state of '" + processName + "' from state '" + ReStreamer.data.states[processName].type + "' to state '" + state + "'");
        ReStreamer.data.states[processName] = {
            type: state,
            message: message
        };
        ReStreamer.writeToDB();
        ReStreamer.updateStreamDataOnGui();
    }

    /**
     *
     * @param {string} processName
     * @param {string} action useraction
     */
    static setUserAction(processName, action){
        logger.debug("Set useraction of '" + processName + "' from '" + ReStreamer.data.userActions[processName] + "' to '" + action + "'");
        ReStreamer.data.userActions[processName] = action;
        ReStreamer.writeToDB();
        ReStreamer.updateStreamDataOnGui();
    }

    /**
     * repeat stream from src to local nginx server (hls and rtmp)
     * @param {string} src
     * @param {string} streamType repeatToOptionalOutput or repeatToLocalNginx
     */
    static startStream(src, streamType, optionalOutput, retryCounter){
        logger.info("start stream " + streamType);
        ReStreamer.data.retryCounter[streamType] = typeof retryCounter === 'undefined' ? 0 : retryCounter;
        if (ReStreamer.data.userActions[streamType] === 'stop'){
            logger.debug("skipping 'startStream' since 'stopped' has been clicked");
            return;
        }

        var command;

        //booleans for easier to read code-structure
        var repeatToLocalNginx = streamType === 'repeatToLocalNginx';
        var repeatToOptionalOutput = streamType === 'repeatToOptionalOutput';

        var addOutputsPromise;
        ReStreamer.updateState(streamType, "connecting");
        //repeat to local nginx server
        if (repeatToLocalNginx){
            command = new FfmpegCommand(src, {
                outputLineLimit: 50
            });
            addOutputsPromise = ReStreamer.addOutput(command, ReStreamer.generateOutputRTMPPath())
            .then(function(){
                return ReStreamer.addOutput(command, ReStreamer.generateOutputHLSPath());
            })
            .catch(function(error){
                logger.error("error adding one or more outputs: " + error.toString);
            });
        //repeat to optional output
        }else if (repeatToOptionalOutput){
            command = new FfmpegCommand(ReStreamer.generateOutputRTMPPath(src, {
                outputLineLimit: 50
            }));
            ReStreamer.data.addresses.optionalOutputAddress = optionalOutput;
            addOutputsPromise = ReStreamer.addOutput(command, optionalOutput);
        }

        addOutputsPromise.then(function(){
            command
                /*
                 STREAMING STARTED
                 */
                .on('start', function(commandLine){
                    if (ReStreamer.data.userActions[streamType] === 'stop'){
                        logger.debug("skipping on 'start' event of ffmpeg command since 'stopped' has been clicked");
                        return;
                    }else{
                        logger.debug("FFMPEG spawned: " + commandLine);
                        ReStreamer.data.processes[streamType] = command;
                        ReStreamer.data.addresses.srcAddress = src;
                        //fetch snapshot only, if repeated to local nginx
                        if (repeatToLocalNginx){
                            ReStreamer.fetchSnapshot(true);
                        }
                    }
                })

                /*
                 ERROR HANDLER
                 */
                .on('error', (error)=>{
                    if (error.toString().indexOf("SIGKILL") > -1) {
                        ReStreamer.updateState(streamType, "disconnected");
                        logger.info("FFMPEG streaming stopped for " + streamType);
                    }else{
                        ReStreamer.data.retryCounter[streamType]++;
                        if (ReStreamer.data.retryCounter[streamType] <= config.ffmpeg.monitor.retries) {
                            if (ReStreamer.data.userActions[streamType] === 'stop'){
                                logger.debug("skipping retry since 'stopped' has been clicked");
                                return;
                            }else{
                                ReStreamer.updateState(streamType, "error", error.toString());
                                logger.info("Retrying ffmpeg connection to " + src + " after " + config.ffmpeg.monitor.restart_wait + "ms");
                                Q.delay(config.ffmpeg.monitor.restart_wait).then(function () {
                                    logger.info("Retry ffmpeg connection to " + src + " retrycounter:" + ReStreamer.data.retryCounter[streamType]);
                                    ReStreamer.startStream(src, streamType, optionalOutput, ReStreamer.data.retryCounter[streamType]);
                                });
                            }
                        }
                    }
                })

                /*
                 STREAMING ENDED
                 */
                .on('end',()=>{
                    ReStreamer.updateState(streamType, "disconnected");
                    ReStreamer.data.retryCounter[streamType]++;
                    if (ReStreamer.data.retryCounter[streamType] <= config.ffmpeg.monitor.retries) {
                        if (ReStreamer.data.userActions[streamType] === 'stop'){
                            logger.debug("skipping retry since 'stopped' has been clicked");
                            return;
                        }else{
                            logger.info("Retrying ffmpeg connection to " + src + " after " + config.ffmpeg.monitor.restart_wait + "ms");
                            Q.delay(config.ffmpeg.monitor.restart_wait).then(function () {
                                logger.info("Retry ffmpeg connection to " + src + " retrycounter:" + ReStreamer.data.retryCounter[streamType]);
                                ReStreamer.startStream(src, streamType, optionalOutput, ReStreamer.data.retryCounter[streamType]);
                            });
                        }
                    }
                });

            /*
             PROGRESS
             */
            var progressMethod = function (progress){
                if (ReStreamer.data.states[streamType].type === 'connecting'){
                    ReStreamer.updateState(streamType, "connected");
                }
                ReStreamer.data.progresses[streamType] = progress;
                ReStreamer.updateProgressOnGui();
                command.removeAllListeners("progress");
            };
            command.on('progress', progressMethod);
            setInterval(()=>{
                if(command.listeners("progress").length === 0) {
                    command.on('progress', progressMethod);
                }
            }, 1000);
            command.exec();
        }).catch(function(error){
            logger.error("Error starting ffmpeg command: " + error.toString());
        });
    }

    /**
     * binded on app-start to bind websocketevents of ReStreamer
     */
    static bindWebsocketEvents(){
        WebsocketsController.addOnConnectionEventToNamespace("/", function(socket){
            socket.on("startStream", (options)=> {
                ReStreamer.setUserAction(options.streamType, "start");
                ReStreamer.startStream(options.src, options.streamType, options.optionalOutput);
            });
            socket.on("stopStream", (streamType)=>{
                ReStreamer.setUserAction(streamType, "stop");
                ReStreamer.stopStream(streamType);
            });
            socket.on("checkForAppUpdates", ()=>{
                const app = require("../webserver/app");
                socket.emit("checkForAppUpdatesResult", app.get("updateAvailable"));
            });
            socket.on("checkStates", ReStreamer.updateStreamDataOnGui);
        });
    }

    /**
     * creates an object of available streaming data like states and addresses to be able
     * to write it to filesystem
     * @returns {object}
     */
    static extractDataOfStreams(){
        var sData = {};
        sData.progresses = ReStreamer.data.progresses;
        sData.userActions = ReStreamer.data.userActions;
        sData.addresses = ReStreamer.data.addresses;
        sData.states = ReStreamer.data.states;
        return sData;
    }
}

/*
define data structure of ReStreamer Data
 */
ReStreamer.data = {
    retryCounter: {
        repeatToLocalNginx: 0,
        repeatToOptionalOutput: 0
    },
    states: {
        repeatToLocalNginx: {
            type: "disconnected",
            message: ""
        },
        repeatToOptionalOutput: {
            type: "disconnected",
            message: ""
        }
    },
     userActions: {
        repeatToLocalNginx: "start",
        repeatToOptionalOutput: "start"
    },
    processes: {
        //overwritten with ffmpeg process if stream has been started
        repeatToLocalNginx: {
            state: "not_connected"
        },
        //overwritten with ffmpeg process if stream has been started
        repeatToOptionalOutput: {
            state: "not_connected"
        }
    },
    progresses: {
        //overwritten with ffmpeg process if stream has been started
        repeatToLocalNginx:{},
        //overwritten with ffmpeg process if stream has been started
        repeatToOptionalOutput: {}
    },
    addresses: {
        srcAddress: "",
        optionalOutputAddress: ""
    }
};

module.exports = ReStreamer;
