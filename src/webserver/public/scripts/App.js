/**
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
'use strict';

var app = window.angular.module('app', [
    'ui.router',
    'ui.bootstrap',
    'pascalprecht.translate',
    'Footer',
    'Header',
    'Login',
    'Main',
    'StreamingInterface'
]);

app.config(($stateProvider) => {
    $stateProvider
        .state('login', {
            'controller': 'loginController',
            'templateUrl': 'views/login.html'
        })
        .state('logged-in', {
            'controller': 'mainController',
            'templateUrl': 'views/main.html'
        });
});

app.controller('appController',
    ['$rootScope', '$state', '$http', ($rootScope, $state, $http) => {
        $http.get('/authenticated').then((response) => {
            $rootScope.loggedIn = response.data;
        });
        $rootScope.$watch('loggedIn', (value) => {
            $state.go(value ? 'logged-in' : 'login');
        });
    }]
);
