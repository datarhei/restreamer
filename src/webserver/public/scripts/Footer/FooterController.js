/**
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
'use strict';

window.angular.module('Footer').controller('footerController',
    ['ws', '$scope', '$http', '$rootScope', 'config', (ws, $scope, $http, $rootScope, config) => {
        $http.get('v1/version').then((response) => {
            $scope.version = response.data.version;
            $rootScope.checkForAppUpdatesResult = response.data.update;
            $scope.config = config;
        });
        $scope.logout = () => {
            $http.get('logout').then(() => {
                $rootScope.loggedIn = false;
            });
        };
    }]
);
