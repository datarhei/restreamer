/**
 * @file holds the Angularjs mainController
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
'use strict';

window.angular.module('Main').controller('mainController',
    [ 'ws', '$scope', '$location', '$rootScope', 'config',
    (ws, $scope, $location, $rootScope, config) => {

    let setup = false;
    let player = null;

    $scope.config = config;

    const initClappr = () => {
        player = new window.Clappr.Player({
            source: 'http://' + window.location.hostname + ':' + window.location.port + '/hls/live.stream.m3u8',
            parentId: '#player',
            baseUrl: '/libs/clappr/dist/',
            poster: 'images/live.jpg',
            mediacontrol: {seekbar: '#3daa48', buttons: '#3daa48'},
            height: '100%',
            width: '100%'
        });
    };

    $scope.optionalOutputInputInvalid = false;
    $scope.nginxRepeatStreamInputInvalid = false;

    $scope.reStreamerData = {
        retryCounter: {
            repeatToLocalNginx: 0,
            repeatToOptionalOutput: 0
        },
        states: {
            repeatToLocalNginx: {
                type: ''
            },
            repeatToOptionalOutput: {
                type: ''
            }
        },
        userActions: {
            repeatToLocalNginx: '',
            repeatToOptionalOutput: ''
        },
        addresses: {
            optionalOutputAddress: '',
            srcAddress: ''
        }
    };

    $rootScope.windowLocationPort = $location.port();

    $scope.optionalOutput = '';

    $scope.showStopButton = function (streamType) {
        return $scope.reStreamerData.userActions[streamType] === 'start';
    };

    $scope.showStartButton = function (streamType) {
        return $scope.reStreamerData.userActions[streamType] === 'stop';
    };

    $scope.nginxRepeatStreamConnecting = function () {
        return $scope.reStreamerData.states.repeatToLocalNginx.type === 'connecting';
    };

    $scope.nginxRepeatStreamConnected = function () {
        return $scope.reStreamerData.states.repeatToLocalNginx.type === 'connected';
    };

    $scope.nginxRepeatStreamError = function () {
        return $scope.reStreamerData.states.repeatToLocalNginx.type === 'error';
    };

    $scope.optionalOutputConnecting = function () {
        return $scope.reStreamerData.states.repeatToOptionalOutput.type === 'connecting';
    };

    $scope.optionalOutputConnected = function () {
        return $scope.reStreamerData.states.repeatToOptionalOutput.type === 'connected';
    };

    $scope.optionalOutputError = function () {
        return $scope.reStreamerData.states.repeatToOptionalOutput.type === 'error';
    };

    $scope.openPlayer = function () {
        if (player === null) {
            initClappr();
        }
        $('#player-modal').modal('show');
        $('#player-modal').on('hide.bs.modal', function (e) {
            player.stop();
            $('#player-modal').off('hide.bs.modal');
            $('#player-modal').modal('hide');
            return e.preventDefault();
        });
    };

    /*
     * Configure Websockets
     */

    ws.emit('getVersion');

    // check states of hls and rtmp stream
    ws.emit('checkStates');

    // check for app updates
    ws.emit('checkForAppUpdates');

    // prohibit double binding of events
    if (!setup) {
        /**
         * test websockets connection (should print below message to browser console if it works)
         */

        ws.on('updateProgress', function (progresses) {
            $scope.reStreamerData.progresses = progresses;
        });
        ws.on('publicIp', function (publicIp) {
            $rootScope.publicIp = publicIp;
        });
        ws.on('updateStreamData', function (reStreamerData) {
            $scope.reStreamerData = reStreamerData;
            if ($scope.showStopButton('repeatToOptionalOutput')) {

                // checkbox
                $scope.activateOptionalOutput = true;
            }
        });
        ws.on('checkForAppUpdatesResult', function (result) {
            $rootScope.checkForAppUpdatesResult = result;
        });
    }

    $scope.startStream = function (streamType) {
        const rtmp_regex = /^(?:rtmp:\/\/|rtsp:\/\/)(?:(?:[^:])+:(?:[^@])+@)?(?:(?:(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?))|(?:(?!\-)(?:[a-zA-Z\d\-]{0,62}[a-zA-Z\d]\.){1,126}(?!\d+)[a-zA-Z\d]{1,63}))(:?:[0-9]{1,4}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])?(?:\/.*)?/;
        var optionalOutput = false;

        if ($scope.activateOptionalOutput === true) {
            optionalOutput = $scope.reStreamerData.addresses.optionalOutputAddress;
        }

        if (streamType === 'repeatToOptionalOutput'){
            $scope.optionalOutputInputInvalid = !rtmp_regex.test(optionalOutput);
            if ($scope.optionalOutputInputInvalid ) {
                return;
            }
        }else {
            $scope.nginxRepeatStreamInputInvalid = !rtmp_regex.test($scope.reStreamerData.addresses.srcAddress);
            if ($scope.nginxRepeatStreamInputInvalid ) {
                return;
            }
        }

        ws.emit('startStream', {
            src: $scope.reStreamerData.addresses.srcAddress,
            streamType: streamType,
            optionalOutput: optionalOutput
        });
    };

    $scope.stopStream = function (streamType) {
        ws.emit('stopStream', streamType);
    };
} ]);
