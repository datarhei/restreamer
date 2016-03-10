/**
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */

/* eslint no-console: 0*/
'use strict';

// styles of the logging output
const INFO = 'color: #0000FF; font-weight: bold';
const DEBUG = 'color: #AABBCC; font-weight: bold';
const ERROR = 'color: #FF0011d; font-weight: bold';
const WEBSOCKETS_IN = 'color: #00BFFF; font-weight: bold';
const WEBSOCKETS_OUT = 'color: #00BF00; font-weight: bold';
const WEBSOCKETS_NAMESPACE = 'color: #00BF00; font-weight: bold';

const LoggerService = function loggerService () {
    /**
     * log an info message
     * @param {string} message
     */
    this.info = (message) => {
        this.log(INFO, message, 'INFO');
    };

    /**
     * log an debug message
     * @param {string} message
     */
    this.debug = (message) => {
        this.log(DEBUG, message, 'DEBUG');
    };

    /**
     * log an error message
     * @param {string} message
     */
    this.error = (message) => {
        this.log(ERROR, message, 'ERROR');
    };

    /**
     * log an websocket in message
     * @param {string} message
     */
    this.websocketsIn = (message) => {
        this.log(WEBSOCKETS_IN, message, 'WS_IN');
    };

    /**
     * log an websocket out message
     * @param {string} message
     */
    this.websocketsOut = (message) => {
        this.log(WEBSOCKETS_OUT, message, 'WS_OUT');
    };

    /**
     * log an websocket namespace message
     * @param {string} message
     */
    this.websocketsNamespace = (message) => {
        this.log(WEBSOCKETS_NAMESPACE, message, 'WS_CONNECT');
    };

    /**
     * log a message with style
     * @param {string} style
     * @param {string} message
     * @param {string} type
     */
    this.log = (style, message, type) => {
        console.log('%c [' + type + ']' + message, style);
    };
};

// configure loggerService as AngularJS Service
window.angular.module('app').factory('loggerService', () => {
    return new LoggerService();
});
