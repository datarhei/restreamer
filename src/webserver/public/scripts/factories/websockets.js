/*
 * @name ReStreamer
 * @namespace https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */

/* jslint browser: true */

'use strict';

datarheiApp.factory('ws', [ '$rootScope', 'alertService', function ($rootScope, alertService) {
    var ws = function(namespace){
        var wsHandler = {};

        window.Logger.log ("WEBSOCKETS_NAMESPACE", "connect to namespace " + namespace);
        var socket;
        if (namespace === "/"){
            socket = io.connect();
        }else{
            socket = io.connect(namespace);
        }
        wsHandler.emit = function (event, data, callback) {
            window.Logger.log("WEBSOCKETS_OUT", "emit event '" + event + "'");
            socket.emit(event, data);
            return wsHandler;
        };
        wsHandler.on = function (event, callback) {
            socket.on(event, function () {
                window.Logger.log("WEBSOCKETS_IN", "got event '" + event + "'");
                var args = arguments;
                $rootScope.$apply(function () {
                    callback.apply(null, args);
                });
            });
            return wsHandler;
        };
        wsHandler.off = function (event, callback) {
            socket.removeListener(event, callback);
        };
        wsHandler
            .on("alert", function(alert) {
                return alertService.add(alert);
            });
        return wsHandler;
    };
    return ws;
}]);
