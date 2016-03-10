/**
 * @file holds the code for the class Nginx
 * @link https://github.com/datarhei/restreamer
 * @copyright 2016 datarhei.org
 * @license Apache-2.0
 */
'use strict';

const Q = require('q');
const psFind = require('ps-find');
const spawn = require('child_process').spawn;
const logger = require('./Logger')('Nginxrtmp');

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
     * generate output hls-path from config-file
     * @returns {string}
     */
    init () {
        return this.getState()
                   .then((state) => {
                       switch (state) {
                           case 'not_running':
                               this.start();
                               break;
                           case 'running':
                               this.logger.info('NGINX already started');
                               break;
                           default:
                               throw new Error('NGINX state could not be detected');
                       }
                   });
    }

    /**
     * start nginx rtmp server
     */
    start () {
        this.logger.info('Starting NGINX...', 'start.nginx');
        this.process = spawn('sh', [
            '-c',
            this.config.nginx.exec
        ]);
        this.bindNginxProcessEvents(this.process);

        // todo: code a watcher instead of fix 1s delay for state detection
        this.getState(1000)
            .then((state)=> {
                switch (state) {

                    case 'not_running':
                        this.logger.info('NGINX could not be started');
                        break;

                    case 'running':
                        this.logger.info('NGINX successfully started');
                        break;

                    default:
                        throw new Error('NGINX state could not be detected');
                }
            });
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
    getState (delayMs) {
        let delay = typeof delayMs === 'undefined' ? 0 : delayMs;
        let nginxProcessString = 'nginx: master process ' + this.config.nginx.exec;
        let state = 'not_running';
        let deferred = Q.defer();

        // delay the state detection if waiting for process is needed
        Q.delay(delay)
         .then(()=> {
             return Q.nfcall(psFind.find, nginxProcessString);
         })
         .then(()=> {
             state = 'running';
         })
         .catch(()=> { // ps-find throws exception in case of 'not found' so we have to handle that
             state = 'not_running';
         })
         .finally(()=> {
             deferred.resolve(state);
         });
        return deferred.promise;
    }
}

module.exports = (config)=> {
    return new Nginxrtmp(config);
};
