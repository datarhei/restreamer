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
    'Main',
    'StreamingInterface',
    'Login']);


app.config(($stateProvider, $urlRouterProvider) => {
    $urlRouterProvider.otherwise('/');
    $stateProvider
        .state('main', {
            'templateUrl': 'main.html',
            'url': '/:error',
            'controller': 'mainController'

        })
        .state('helpSource', {
            'templateUrl': 'help/source.html',
            'url': '/help/source'
        })
        .state('helpOptionalOutput', {
            'templateUrl': 'help/optionalOutput.html',
            'url': '/help/optionalOutput'
        });
});


