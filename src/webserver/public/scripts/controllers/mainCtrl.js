/**
 * @file holds the Angularjs mainController
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */

window.app.controller('mainCtrl', [ 'ws', '$scope', '$location', '$rootScope', '$translate', function (ws, $scope, $location, $rootScope, $translate) {

    // binding just once
    var setup = false;
    var player = null;

    const initClappr = function () {
        player = new Clappr.Player({
            source: ("https:" === window.location.protocol ? "https:" : "http:") + "//" + window.location.hostname + ':' + window.location.port + '/hls/live.stream.m3u8',
            parentId: '#player',
            baseUrl: '/libs/clappr/dist/',
            poster: 'images/live.jpg',
            mediacontrol: {seekbar: '#3daa48', buttons: '#3daa48'},
            height: '100%',
            width: '100%'
        });
    };

    $translate.use('en_US');

    $rootScope.loggedIn = false;

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

    $scope.nginxRepeatReachedRetryLimit = function () {
        return $scope.reStreamerData.retryCounter.repeatToLocalNginx.current > $scope.reStreamerData.retryCounter.repeatToLocalNginx.max;
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

    $scope.optionalOutputReachedRetryLimit = function () {
        return $scope.reStreamerData.retryCounter.repeatToOptionalOutput.current > $scope.reStreamerData.retryCounter.repeatToOptionalOutput.max;
    };

    $scope.openPlayer = function () {
        if (player === null) {
            initClappr();
        }
        var modal = $('#player-modal');
        modal.modal('show').on('hide.bs.modal', function (e) {
            player.stop();
            modal.off('hide.bs.modal').modal('hide');
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
        ws.on('version', function (version) {
            $rootScope.version = version;
            window.Logger.log('INFO', 'Datarhei ' + version + ' websockets connected');
        });
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
