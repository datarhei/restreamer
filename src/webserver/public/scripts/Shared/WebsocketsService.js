/**
 * @file Service to handle websocket connections and events
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */

/* jslint browser: true */
'use strict';

//import {app} from '../App';

class WebsocketsService {

    /*
     * construct the Websockets Service
     * @param $rootScope
     * @param loggerService
     */
    constructor ($rootScope, loggerService) {
        this.$rootScope = $rootScope;
        this.loggerService = loggerService;
        this.socket = io.connect();
        this.loggerService.websockets_namespace('websockets connected');
    }

    /*
     * emit an event to socket
     * @param event
     * @param data
     * @returns {WebsocketsService}
     */
    emit (event, data) {
        this.loggerService.websockets_out(`emit event "${event}"`);
        this.socket.emit(event, data);
        return this;
    }

    /*
     * react on an event to socket with callback
     * @param event
     * @param {function} callback
     * @returns {WebsocketsService}
     */
    on (event, callback) {
        this.loggerService.websockets_in(`got event "${event}"`);
        var self = this;

        this.socket.on(event, function () {
            var args = arguments;

            self.$rootScope.$apply(function () {
                callback.apply(null, args);
            });
        });
        return this;
    }

    /*
     * disable an event on socket
     * @param event
     * @param callback
     */
    off (event, callback) {
        this.socket.removeListener(event, callback);
    }
}

/*
 * connect service to angular.js
 */

window.angular.module('app').factory('ws', [ '$rootScope', 'loggerService', ($rootScope, loggerService) => {
    return new WebsocketsService($rootScope, loggerService);
} ]);
