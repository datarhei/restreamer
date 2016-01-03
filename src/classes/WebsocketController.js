/*
 * @name Restreamer
 * @namespace https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */

'use strict';

const async = require("async");
const Logger = require("../classes/Logger");
const logger = new Logger("WebsocketsController");

class WebsocketsController {

    static emitToNamespace(namespace, event, data) {
        var app = require("../webserver/app");
        app.get("websocketsReady").promise.then(function (io) {
            logger.debug("websocket got event " + event + " to namespace " + namespace + "","Websockets");
            io.of(namespace).emit(event, data);
        });
    }

    static addOnConnectionEventToNamespace(namespace, callback) {
        var app = require("../webserver/app");
        app.get("websocketsReady").promise.then(function (io) {
            var nsp = io.of(namespace);
            nsp.on("connection", function (socket) {
                callback(socket);
            });
        });
    }

    static bindDefaultEvents() {
        /*
         update Endpoint Class
         */
        WebsocketsController.addOnConnectionEventToNamespace("/", function (socket) {
            socket.on("testConnection", (options)=> {
                var packageJson = require("../../package.json");
                socket.emit("connection", packageJson.version);
            });
            var app = require("../webserver/app");
            socket.emit("publicIp", app.get("publicIp"));
        });
        require("./Restreamer").bindWebsocketEvents();
    }
}

module.exports = WebsocketsController;
