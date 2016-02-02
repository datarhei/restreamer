/**
 * @file holds the code for the class WebsocketsController
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */

const logger = require('../classes/Logger')('WebsocketsController');

/**
 * static class websocket controller, that helps communicating through websockets to different namespaces and ensures
 * that websocket events are binded, if the websocket server has been initialized (through promise made on app start)
 * @todo since currently the restream is a singlepage application, there is no need to use different namespaces
 */
class WebsocketsController {

    /**
     *
     * @param {string} namespace namespace to emit the event to
     * @param {string} event name of the event
     * @param {object} data data to emit to the client event listener
     */
    static emitToNamespace (namespace, event, data) {
        var app = require('../webserver/app');

        app.get('websocketsReady').promise.then(function (io) {
            logger.debug('websocket got event ' + event + ' to namespace ' + namespace + '', 'Websockets');
            io.of(namespace).emit(event, data);
        });
    }

    /**
     * add event, that is emmi
     * @param {string} namespace
     * @param {function} callback
     */
    static addOnConnectionEventToNamespace (namespace, callback) {
        var app = require('../webserver/app');

        app.get('websocketsReady').promise.then(function (io) {
            var nsp = io.of(namespace);

            nsp.on('connection', function (socket) {
                callback(socket);
            });
        });
    }

    /**
     * bind default events of all classes that are using websocketevents
     */
    static bindDefaultEvents () {
        WebsocketsController.addOnConnectionEventToNamespace('/', function (socket) {
            socket.on('getVersion', (options)=> {
                var packageJson = require('../../package.json');

                socket.emit('version', packageJson.version);
            });
            var app = require('../webserver/app');

            socket.emit('publicIp', app.get('publicIp'));
        });
        require('./Restreamer').bindWebsocketEvents();
    }
}

module.exports = WebsocketsController;
