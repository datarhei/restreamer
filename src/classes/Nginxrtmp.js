/**
 * @file holds the code for the class Nginx
 * @link https://github.com/datarhei/restreamer
 * @copyright 2016 datarhei.org
 * @license Apache-2.0
 */
'use strict';

const Q = require('q');
const config = require("../../conf/live.json");
const spawn = require('child_process').spawn;
const logger = require('./Logger')('Nginxrtmp');
const rp = require('request-promise');
const RUNNING = "running";
const NOT_RUNNING = "not_running";

/**
 * class to watch and control the NGINX RTMP server process
 */
class Nginxrtmp {

    /**
     * constructs the Nginx rtmp with injection of config to use
     * @param config
     */
    constructor (config) {
        this.config = config;
        this.process = null;
        this.logger = logger;
    }

    /**
     * init the nginx rtmp server
     * @returns {Promise.<boolean>}
     */
    async init () {
        const currentState = await this.getState();
        if (currentState === NOT_RUNNING){
            await this.start();
        }else{
            return true;
        }
    }

    /**
     * start the nginx rtmp server
     * @returns {Promise.<boolean>}
     */
    async start () {
        this.logger.info('Starting NGINX...', 'start.nginx');
        let state = NOT_RUNNING;
        let tries = 0;
        let timeout = 250;
        this.process = spawn('sh', [
            '-c',
            this.config.nginx.exec
        ]);
        this.bindNginxProcessEvents(this.process);

        // todo: code a watcher instead of fix 1s delay for state detection
        while(state === NOT_RUNNING){
            state = await this.getState(timeout);
            tries++;
        }
        if (state === RUNNING){
            this.logger.info(`NGINX successfully started in ${tries * timeout}ms`);
        }else{
            this.logger.info('NGINX could not be started');
        }
        return true;
    }

    /**
     * bind events on outputs of the nginx process
     * @param {process} process
     */
    bindNginxProcessEvents (process) {
        process.stdout.on('data', (data) => {
            this.logger.info(`The NGINX RTMP process created an output: ${data}`);
        });

        process.stderr.on('data', (data) => {
            this.logger.error(`The NGINX RTMP process created an error output: ${data}`);
        });

        process.stderr.on('close', (code) => {
            this.logger.error(`The NGINX RTMP process exited with code: ${code}`);
            this.start();
        });
    }

    /**
     * get current state of nginx rtmp server
     * @returns {*|promise}
     */
    async getState (delayMs) {
        const url = "http://" + config.nginx.streaming.ip + ":" +  config.nginx.streaming.http_port  +  config.nginx.streaming.http_health_path;
        const delay = typeof delayMs === 'undefined' ? 0 : delayMs;
        // delay the state detection if waiting for process is needed
        try{
            await Q.delay(delay);
            const response = (await rp("http://" + config.nginx.streaming.ip + ":" +  config.nginx.streaming.http_port  +  config.nginx.streaming.http_health_path));
            return response === 'pong' ? RUNNING: NOT_RUNNING;
        }catch(error){
            return NOT_RUNNING;
        }
    }
}

module.exports = (config)=> {
    return new Nginxrtmp(config);
};
