/**
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */

'use strict';

window.angular.module('app').config(($stateProvider, $urlRouterProvider) => {
    $urlRouterProvider.otherwise('/');
    $stateProvider

        .state('main', {
            templateUrl: 'main.html',
            url: '/',
            controller: 'mainController'}
        )

        .state('helpSource', {
            templateUrl: 'help/source.html',
            url: '/help/source'}
        )

        .state('helpOptionalOutput', {
            templateUrl: 'help/optionalOutput.html',
            url: '/help/optionalOutput'}
        );
});
