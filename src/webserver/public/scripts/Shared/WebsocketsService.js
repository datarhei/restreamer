/**
 * @file Service to handle websocket connections and events
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */

/* eslint no-undef: 0*/
'use strict';

const WebsocketsService = function websockeService ($rootScope, loggerService) {
    this.$rootScope = $rootScope;
    this.loggerService = loggerService;
    this.socket = io.connect();
    this.loggerService.websocketsNamespace('websockets connected');


    /**
     * emit an event to socket
     * @param event
     * @param data
     * @returns {WebsocketsService}
     */
    this.emit = (event, data) => {
        this.loggerService.websocketsOut(`emit event "${event}"`);
        this.socket.emit(event, data);
        return this;
    };

    /**
     * react on an event to socket with callback
     * @param event
     * @param {function} callback
     * @returns {WebsocketsService}
     */
    this.on = (event, callback) => {
        var self = this;
        this.loggerService.websocketsIn(`got event "${event}"`);
        this.socket.on(event, function woEvent () {
            var args = arguments;
            self.$rootScope.$apply(function weApply () {
                callback.apply(null, args);
            });
        });
        return this;
    };

    /**
     * disable an event on socket
     * @param event
     * @param callback
     */
    this.off = (event, callback) => {
        this.socket.removeListener(event, callback);
    };
};

// connect service to angular.js
window.angular.module('app').factory('ws', ['$rootScope', 'loggerService', ($rootScope, loggerService) => {
    return new WebsocketsService($rootScope, loggerService);
}]);
