/**
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
'use strict';

window.angular.module('Footer').controller('footerController', ['ws', '$scope', 'config', (ws, $scope, config) => {
    ws.emit('getVersion');
    ws.on('version', (version) => {
        $scope.version = version;
        $scope.config = config;
    });
}]);
