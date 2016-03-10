/**
 * @file holds the code for the class WebsocketsController
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
'use strict';

const logger = require.main.require('./classes/Logger')('WebsocketsController');
const app = require.main.require('./webserver/app').app;
const packageJson = require(require('path').join(global.__base, 'package.json'));

/**
 * static class websocket controller, that helps communicating through websockets to different namespaces and ensures
 * that websocket events are bound, if the websocket server has been initialized (through promise made on app start)
 * @todo since currently the Restreamer is a single page application, there is no need to use different namespaces
 */
class WebsocketsController {

    /**
     *
     * @param {string} namespace namespace to emit the event to
     * @param {string} event name of the event
     * @param {object} data data to emit to the client event listener
     */
    static emitToNamespace (namespace, event, data) {
        app.get('websocketsReady').promise.then((io) => {
            logger.debug(`websocket got event ${event} to namespace ${namespace}`, 'Websockets');
            io.of(namespace).emit(event, data);
        });
    }

    /**
     * add event, that is emmi
     * @param {string} namespace
     * @param {function} callback
     */
    static addOnConnectionEventToNamespace (namespace, callback) {
        app.get('websocketsReady').promise.then((io) => {
            var nsp = io.of(namespace);

            nsp.on('connection', (socket) => {
                callback(socket);
            });
        });
    }

    /**
     * bind default events of all classes that are using websockets events
     */
    static bindDefaultEvents () {
        WebsocketsController.addOnConnectionEventToNamespace('/', (socket) => {
            socket.on('getVersion', () => {
                socket.emit('version', packageJson.version);
            });
            socket.emit('publicIp', app.get('publicIp'));
        });
        require('./Restreamer').bindWebsocketEvents();
    }
}

module.exports = WebsocketsController;
