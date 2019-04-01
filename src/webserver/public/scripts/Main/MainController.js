/**
 * @file holds the AngularJS mainController
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
'use strict';

window.angular.module('Main').controller('mainController',
    ['ws', '$scope', '$location', '$rootScope', '$stateParams', 'config', function mainController (ws, $scope, $location, $rootScope, $stateParams, config) {
        let setup = false;
        let player = null;
        let posterPlugin = null;

        $scope.config = config;

        const updateSnapshot = () => {
            if (posterPlugin !== null) {
                posterPlugin.options.poster = 'images/live.jpg?t=' + String(new Date().getTime());
                if (!player.isPlaying()) {
                    posterPlugin.render();
                }
            }
        };

        const initClappr = () => {
            player = new window.Clappr.Player({
                'source': 'hls/live.stream.m3u8',
                'parentId': '#player',
                'baseUrl': 'libs/scripts/',
                'poster': 'images/live.jpg?t=' + String(new Date().getTime()),
                'mediacontrol': {'seekbar': '#3daa48', 'buttons': '#3daa48'},
                'height': '100%',
                'width': '100%'
            });
            posterPlugin = player.core.mediaControl.container.getPlugin('poster');
            player.on(window.Clappr.Events.PLAYER_STOP, () => {
                posterPlugin.render();
            });
        };

        $scope.optionalOutputInputInvalid = false;
        $scope.nginxRepeatStreamInputInvalid = false;

        $scope.reStreamerData = {
            'options': {
                'rtspTcp': false
            },
            'states': {
                'repeatToLocalNginx': {
                    'type': '',
                    'message': ''
                },
                'repeatToOptionalOutput': {
                    'type': '',
                    'message': ''
                }
            },
            'userActions': {
                'repeatToLocalNginx': '',
                'repeatToOptionalOutput': ''
            },
            'progresses': {
                'repeatToLocalNginx': '',
                'repeatToOptionalOutput': ''
            },
            'addresses': {
                'optionalOutputAddress': '',
                'srcAddress': ''
            }
        };

        $rootScope.windowLocationPort = $location.port();

        $scope.optionalOutput = '';

        $scope.showStartButton = (streamType) => {
            return ($scope.reStreamerData.states[streamType].type == 'disconnected');
        };

        $scope.showStopButton = (streamType) => {
            let state = $scope.reStreamerData.states[streamType].type;

            return (state == 'connected' || state == 'connecting' || state == 'error');
        };

        $scope.disableInput = (streamType) => {
            return ($scope.reStreamerData.states[streamType].type != 'disconnected');
        };

        $scope.openPlayer = () => {
            if (player === null) {
                initClappr();
            }
            $('#player-modal').modal('show').on('hide.bs.modal', function closeModal (e) {
                player.stop();
                $(this).off('hide.bs.modal');
                $(this).modal('hide');
                return e.preventDefault();
            });
        };

        /**
         * Configure Websockets
         */

        ws.emit('checkStates'); // check states of hls and rtmp stream

        // prohibit double binding of events
        if (!setup) {
            /*
             * test websockets connection (should print below message to browser console if it works)
             */
            ws.on('updateProgress', (progresses) => {
                $scope.reStreamerData.progresses = progresses;
            });
            ws.on('publicIp', (publicIp) => {
                $rootScope.publicIp = publicIp;
            });
            ws.on('updateStreamData', (reStreamerData) => {
                $scope.reStreamerData = reStreamerData;
                if ($scope.showStopButton('repeatToOptionalOutput')) {
                    // checkbox
                    $scope.activateOptionalOutput = true;
                }
            });
            ws.on('snapshot', updateSnapshot);
        }

        $scope.startStream = (streamType) => {
            const inputRegex = /^(rtmp(s|t)?|rtsp|https?):\/\//;
            const outputRegex = /^rtmp(s|t)?/;

            var optionalOutput = '';
            if($scope.activateOptionalOutput === true) {
                optionalOutput = $scope.reStreamerData.addresses.optionalOutputAddress;
            }

            if(streamType == 'repeatToOptionalOutput') {
                $scope.optionalOutputInputInvalid = !outputRegex.test(optionalOutput);
                if($scope.optionalOutputInputInvalid) {
                    return;
                }
            }
            else {
                $scope.nginxRepeatStreamInputInvalid = !inputRegex.test($scope.reStreamerData.addresses.srcAddress);
                if($scope.nginxRepeatStreamInputInvalid) {
                    return;
                }
            }

            ws.emit('startStream', {
                'src': $scope.reStreamerData.addresses.srcAddress,
                'options': $scope.reStreamerData.options,
                'streamType': streamType,
                'optionalOutput': optionalOutput
            });
        };

        $scope.stopStream = (streamType) => {
            ws.emit('stopStream', streamType);
        };
    }]
);
