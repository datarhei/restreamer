/**
 * @file holds the code for the class WebsocketsController
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
'use strict';

const logger = require.main.require('./classes/Logger')('WebsocketsController');
const app = require.main.require('./webserver/app').app;

/**
 * static class websocket controller, that helps communicating through websockets to different namespaces and ensures
 * that websocket events are bound, if the websocket server has been initialized (through promise made on app start)
 */
class WebsocketsController {

    /**
     * emit an event to WS
     * @param {string} event name of the event
     * @param {object} data data to emit to the client event listener
     */
    static emit (event, data) {
        app.get('websocketsReady').promise.then((io) => {
            logger.debug(`Emitting ${event}`);
            io.sockets.emit(event, data);
        });
    }

    /**
     * add callback on WS connection
     * @param {function} callback
     */
    static setConnectCallback (callback) {
        app.get('websocketsReady').promise.then((io) => {
            io.on('connection', (socket) => {
                callback(socket);
            });
        });
    }
}

module.exports = WebsocketsController;
