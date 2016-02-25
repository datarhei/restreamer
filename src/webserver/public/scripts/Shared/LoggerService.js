/**
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
'use strict';

//import {app} from '../App';

/*
styles of the logging outputs
 */
const INFO = 'color: #0000FF; font-weight: bold';
const DEBUG = 'color: #AABBCC; font-weight: bold';
const ERROR = 'color: #FF0011d; font-weight: bold';
const WEBSOCKETS_IN = 'color: #00BFFF; font-weight: bold';
const WEBSOCKETS_OUT = 'color: #00BF00; font-weight: bold';
const WEBSOCKETS_NAMESPACE = 'color: #00BF00; font-weight: bold';

class LoggerService {

    /*
     * no further dependencies needed
     */
    constructor () {}

    /*
     * log an info message
     * @param {string} message
     */
    info (message) {
        this.log(INFO, message, 'INFO');
    }

    /*
     * log an info message
     * @param {string} message
     */
    debug (message) {
        this.log(DEBUG, message, 'DEBUG');
    }

    /*
     * log an info message
     * @param {string} message
     */
    error (message) {
        this.log(ERROR, message, 'ERROR');
    }

    /*
     * log an info message
     * @param {string} message
     */
    websockets_in (message) {
        this.log(WEBSOCKETS_IN, message, 'WS_IN');
    }

    /*
     * log an info message
     * @param {string} message
     */
    websockets_out (message) {
        this.log(WEBSOCKETS_OUT, message, 'WS_OUT');
    }

    /*
     * log an info message
     * @param {string} message
     */
    websockets_namespace (message) {
        this.log(WEBSOCKETS_NAMESPACE, message, 'WS_CONNECT');
    }

    /*
     * log a message with style
     * @param {string} style
     * @param {string} message
     */
    log (style, message, type) {
        console.log('%c ' + '[' + type + ']' + message, style);
    }
}

/*
 * configure loggerService as angulerjs Service
 */
window.angular.module('app').service('loggerService', () => {
    return new LoggerService();
});
