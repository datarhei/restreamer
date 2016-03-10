/**
 * @file holds the Angularjs mainController
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
'use strict';

window.angular.module('Main').controller('mainController',
    ['ws', '$scope', '$location', '$rootScope', '$stateParams', 'config',
        function mainController (ws, $scope, $location, $rootScope, $stateParams, config) {
            let setup = false;
            let player = null;

            if ($stateParams.error === 'login_invalid') {
                $scope.login_invalid = 'login_invalid';
            }

            $scope.config = config;

            const initClappr = () => {
                player = new window.Clappr.Player({
                    'source': (window.location.protocol === 'https:' ? 'https:' : 'http:') +
                    '//' + window.location.hostname + ':' + window.location.port + '/hls/live.stream.m3u8',
                    'parentId': '#player',
                    'baseUrl': '/libs/clappr/dist/',
                    'poster': 'images/live.jpg',
                    'mediacontrol': {'seekbar': '#3daa48', 'buttons': '#3daa48'},
                    'height': '100%',
                    'width': '100%'
                });
            };

            $rootScope.loggedIn = false;

            $scope.optionalOutputInputInvalid = false;
            $scope.nginxRepeatStreamInputInvalid = false;

            $scope.reStreamerData = {
                'retryCounter': {
                    'repeatToLocalNginx': 0,
                    'repeatToOptionalOutput': 0
                },
                'options': {
                    'rtspTcp': false
                },
                'states': {
                    'repeatToLocalNginx': {
                        'type': ''
                    },
                    'repeatToOptionalOutput': {
                        'type': ''
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

            $scope.showStopButton = (streamType) => {
                return $scope.reStreamerData.userActions[streamType] === 'start';
            };

            $scope.showStartButton = (streamType) => {
                return $scope.reStreamerData.userActions[streamType] === 'stop';
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

            /*
             * Configure Websockets
             */

            // check states of hls and rtmp stream
            ws.emit('checkStates');

            // check for app updates
            ws.emit('checkForAppUpdates');

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
                ws.on('checkForAppUpdatesResult', (result) => {
                    $rootScope.checkForAppUpdatesResult = result;
                });
            }

            $scope.startStream = (streamType) => {
                const rtmpRegex = /^(?:rtmp:\/\/|rtsp:\/\/)(?:(?:[^:])+:(?:[^@])+@)?(?:(?:(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?))|(?:(?!\-)(?:[a-zA-Z\d\-]{0,62}[a-zA-Z\d]\.){1,126}(?!\d+)[a-zA-Z\d]{1,63}))(:?:[0-9]{1,4}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])?(?:\/.*)?/;
                var optionalOutput = '';

                if ($scope.activateOptionalOutput === true) {
                    optionalOutput = $scope.reStreamerData.addresses.optionalOutputAddress;
                }

                if (streamType === 'repeatToOptionalOutput') {
                    $scope.optionalOutputInputInvalid = !rtmpRegex.test(optionalOutput);
                    if ($scope.optionalOutputInputInvalid) {
                        return;
                    }
                } else {
                    $scope.nginxRepeatStreamInputInvalid = !rtmpRegex.test($scope.reStreamerData.addresses.srcAddress);
                    if ($scope.nginxRepeatStreamInputInvalid) {
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
        }]);
