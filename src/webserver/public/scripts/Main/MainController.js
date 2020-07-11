/**
 * @file holds the AngularJS mainController
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
'use strict';

window.angular.module('Main').controller('mainController',
    ['ws', '$scope', '$location', '$rootScope', '$stateParams', 'config', 'loggerService', function mainController (ws, $scope, $location, $rootScope, $stateParams, config, logger) {
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
            const plugins = [];
            if($scope.reStreamerData.options.player.statistics == true) {
                plugins.push(ClapprNerdStats);
                plugins.push(ClapprStats);
            }

            const config = {
                source: 'hls/live.stream.m3u8',
                parentId: '#player',
                baseUrl: 'libs/scripts/',
                poster: 'images/live.jpg?t=' + String(new Date().getTime()),
                mediacontrol: {
                    seekbar: $scope.reStreamerData.options.player.color,
                    buttons: $scope.reStreamerData.options.player.color
                },
                height: '100%',
                width: '100%',
                disableCanAutoPlay: true,
                autoPlay: $scope.reStreamerData.options.player.autoplay,
                mute: $scope.reStreamerData.options.player.mute,
                plugins: plugins,
                clapprStats: {
                    runEach: 1000,
                    onReport: (metrics) => {},
                },
                clapprNerdStats: {
                    shortcut: ['command+shift+s', 'ctrl+shift+s'],
                    iconPosition: 'top-right'
                }
            };

            if($scope.reStreamerData.options.player.logo.image.length != 0) {
                config.watermark = $scope.reStreamerData.options.player.logo.image;
                config.position = $scope.reStreamerData.options.player.logo.position;

                if($scope.reStreamerData.options.player.logo.link.length != 0) {
                    config.watermarkLink = $scope.reStreamerData.options.player.logo.link;
                }
            }

            $('#player').empty();

            player = new window.Clappr.Player(config);
            posterPlugin = player.core.mediaControl.container.getPlugin('poster');
            player.on(window.Clappr.Events.PLAYER_STOP, () => {
                posterPlugin.render();
            });
        };

        $scope.optionalOutputInputInvalid = false;
        $scope.nginxRepeatStreamInputInvalid = false;

        $scope.reStreamerData = {
            options: {
                rtspTcp: false,
                video: {
                    codec: 'copy',
                    preset: 'ultrafast',
                    bitrate: 4096,
                    profile: 'auto',
                    tune: 'none'
                },
                audio: {
                    codec: 'copy',
                    preset: 'silence',
                    bitrate: 64,
                    channels: 'mono',
                    sampling: 41000
                },
                player: {
                    autoplay: false,
                    mute: false,
                    statistics: false,
                    color: '#3daa48',
                    logo: {
                        image: '',
                        position: 'bottom-right',
                        link: ''
                    }
                },
                outputs: [{
                    label: '',
                    type: 'rtmp',
                    rtmp: {},
                    hls: {
                        method: 'POST'
                    },
                    outputAddress: ''
                }],
            },
            states: {
                repeatToLocalNginx: {
                    type: '',
                    message: ''
                },
                repeatToOptionalOutput_0: {
                    type: '',
                    message: ''
                }
            },
            userActions: {
                repeatToLocalNginx: '',
                repeatToOptionalOutput_0: ''
            },
            progresses: {
                repeatToLocalNginx: '',
                repeatToOptionalOutput_0: ''
            },
            addresses: {
                srcAddress: ''
            },
        };

        $rootScope.windowProtocol = window.location.protocol;
        $rootScope.windowLocationPort = window.location.port ? `:${window.location.port}` : '';
        $rootScope.windowLocationPath = window.location.pathname;

        $scope.optionalOutput = '';
        $scope.hasMultipleOutputs = $scope.reStreamerData.options.outputs.length > 1;

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
            initClappr();

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
                $scope.hasMultipleOutputs = $scope.reStreamerData.options.outputs.length > 1;
                if ($scope.showStopButton('repeatToOptionalOutput_0')) {
                    // checkbox
                    $scope.activateOptionalOutput = true;
                }
            });
            ws.on('snapshot', updateSnapshot);
        }

        $scope.startStream = (streamType) => {
            const inputRegex = /^(rtmp(s|t)?|rtsp|https?):\/\//;
            const outputRegexRTMP = /^rtmp(s|t)?:\/\//;
            const outputRegexHLS = /^https?:\/\/.*\.m3u8/;
            const outputRegexOptionalOutput = /(repeatToOptionalOutput)_(\d)$/;

            var optionalOutput = '';
            const optionalMatches = streamType.match(outputRegexOptionalOutput);
            let streamTypeMatch = null;
            let streamIndexMatch = null;
            if (optionalMatches && optionalMatches.length > 2) {
                streamTypeMatch = optionalMatches[1];
                streamIndexMatch = optionalMatches[2];
            }

            if ($scope.activateOptionalOutput === true && streamTypeMatch === 'repeatToOptionalOutput') {
                $scope.optionalOutputInputInvalid = true;

                let currentOption = false;

                if (streamIndexMatch >= 0 || streamIndexMatch < 5) {
                    currentOption = $scope.reStreamerData.options.outputs[streamIndexMatch];
                }

                if (!currentOption) {
                    return;
                }

                optionalOutput = currentOption.outputAddress;

                if (currentOption.type == 'rtmp') {
                    $scope.optionalOutputInputInvalid = !outputRegexRTMP.test(optionalOutput);
                }
                else if(currentOption.type == 'hls') {
                    $scope.optionalOutputInputInvalid = !outputRegexHLS.test(optionalOutput);
                }

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
                src: $scope.reStreamerData.addresses.srcAddress,
                options: $scope.reStreamerData.options,
                streamType: streamType,
                optionalOutput: optionalOutput
            });
        };

        $scope.stopStream = (streamType) => {
            ws.emit('stopStream', streamType);
        };

        $scope.playerOptions = () => {
            ws.emit('playerOptions', $scope.reStreamerData.options.player);
        };
    }]
);
