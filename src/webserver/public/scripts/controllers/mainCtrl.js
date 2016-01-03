/*
 * @name Restreamer
 * @namespace https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */

/**
 Angularjs Controller mainController

 controlls the restreamer workflow
 */

window.datarheiApp.controller('mainCtrl',['ws', '$scope', '$location', '$rootScope', '$translate', function(ws, $scope, $location, $rootScope, $translate) {
    //binding just once
    var setup = false;
    $translate.use("en_US");

    $scope.reStreamerData = {
      states: {
          repeatToLocalNginx: {
              type: ""
          },
          repeatToOptionalOutput: {
              type: ""
          }
      },
      userActions: {
          repeatToLocalNginx: "",
          repeatToOptionalOutput: ""
      },
      addresses: {
          optionalOutputAddress: "",
          srcAddress: ""
      }
    };

    $rootScope.windowLocationPort = $location.port();

    $scope.optionalOutput = "";

    $scope.showStopButton = function(streamType){
        return $scope.reStreamerData.userActions[streamType] === "start";
    };

    $scope.showStartButton = function(streamType){
        return $scope.reStreamerData.userActions[streamType] === "stop";
    };

    $scope.nginxRepeatStreamConnecting = function(){
        return $scope.reStreamerData.states.repeatToLocalNginx.type  === "connecting";
    };

    $scope.nginxRepeatStreamConnected = function(){
        return $scope.reStreamerData.states.repeatToLocalNginx.type === "connected";
    };

    $scope.nginxRepeatStreamError = function(){
        return  $scope.reStreamerData.states.repeatToLocalNginx.type === "error";
    };

    $scope.optionalOutputConnecting = function(){
        return $scope.reStreamerData.states.repeatToOptionalOutput.type  === "connecting";
    };

    $scope.optionalOutputConnected = function(){
        return $scope.reStreamerData.states.repeatToOptionalOutput.type === "connected";
    };

    $scope.optionalOutputError = function(){
        return  $scope.reStreamerData.states.repeatToOptionalOutput.type === "error";
    };

    /*
     Configure Websockets
     */

    //connect to namespace /
    ws = ws("/");

    ws.emit('testConnection');

    //check states of hls and rtmp stream
    ws.emit('checkStates');

    ws.emit('checkForAppUpdates');

    //prohibit double binding of events
    if (!setup) {
        /**
         * test websockets connection (should print below message to browser console if it works)
         */
        ws.on("connection", function (version) {
            $rootScope.version = version;
            window.Logger.log("INFO", "Datarhei " + version + " websockets connected");
            var player = new Clappr.Player({
                source: "http://" + window.location.hostname + ":" + window.location.port + "/hls/live.stream.m3u8",
                parentId: "#player",
                baseUrl: '/libs/clappr/dist/',
                poster: "images/live.jpg",
                mediacontrol: {seekbar: "#3daa48", buttons: "#3daa48"},
                height: "100%",
                width: "100%"
            });
        });

        ws.on("updateProgress", function(progresses){
          $scope.reStreamerData.progresses = progresses;
        });
        ws.on("publicIp", function(publicIp){
            $rootScope.publicIp = publicIp;
        });
        ws.on("updateStreamData", function(reStreamerData) {
            $scope.reStreamerData = reStreamerData;
            if ($scope.showStopButton('repeatToOptionalOutput')) {
                $scope.activateOptionalOutput = true; //checkbox;
            }
        });
        ws.on("checkForAppUpdatesResult", function(result){
            $rootScope.checkForAppUpdatesResult = result;
        });
    }

    $scope.startStream = function(streamType){
        var optionalOutput = false;
        if($scope.activateOptionalOutput === true){
            optionalOutput = $scope.reStreamerData.addresses.optionalOutputAddress;
        }
        ws.emit("startStream", {
           src: $scope.reStreamerData.addresses.srcAddress,
           streamType: streamType,
           optionalOutput: optionalOutput
        });
    };

    $scope.stopStream = function(streamType){
        ws.emit("stopStream", streamType);
    };
}]);
