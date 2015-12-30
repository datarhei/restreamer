/*
 * @name ReStreamer
 * @namespace https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */

'use strict';

var datarheiApp = angular.module('datarheiApp', ['ui.router','ui.bootstrap', 'pascalprecht.translate']);

datarheiApp.config(function ($translateProvider) {
    // deutsche Sprache
    $translateProvider.useStaticFilesLoader({
        prefix: 'locales/lang-',
        suffix: '.json'
    });

    $translateProvider.useSanitizeValueStrategy('escape');
    $translateProvider.preferredLanguage("en_US");
});

datarheiApp.config(function($stateProvider, $urlRouterProvider){
    // For any unmatched url, redirect to /
    $urlRouterProvider.otherwise("/");

    $stateProvider
    .state('main',{
        templateUrl: 'main.html',
        url:'/',
        controller: 'mainCtrl'}
    )
    .state('helpSource',{
        templateUrl: 'help/source.html',
        url:'/help/source'}
    )
    .state('helpOptionalOutput',{
        templateUrl: 'help/optionalOutput.html',
        url:'/help/optionalOutput'}
    );
});

datarheiApp.filter('inArray', function($filter){
    return function(list, arrayFilter, element){
        if(arrayFilter){
            return $filter("filter")(list, function(listItem){
                return arrayFilter.indexOf(listItem[element]) != -1;
            });
        }
    };
});
