/**
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
'use strict';

window.angular.module('Login').controller('loginController',
    ['$scope', '$http', '$rootScope', function loginController ($scope, $http, $rootScope) {
        $scope.submit = function submit () {
            $http.post('login', {'user': $scope.user, 'pass': $scope.pass}).then((response) => {
                $scope.message = response.data.message;
                $rootScope.loggedIn = response.data.success;
            });
        };
    }]
);
