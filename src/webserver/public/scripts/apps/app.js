/**
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */

var app = angular.module('app', [ 'ui.router', 'ui.bootstrap', 'pascalprecht.translate' ]);

app.config(function ($translateProvider) {

    // lang
    $translateProvider.useStaticFilesLoader({
        prefix: 'locales/lang-',
        suffix: '.json'
    });

    $translateProvider.useSanitizeValueStrategy('escape');
    $translateProvider.preferredLanguage('en_US');
});

app.config(function ($stateProvider, $urlRouterProvider) {

    // For any unmatched url, redirect to /
    $urlRouterProvider.otherwise('/');

    $stateProvider
    .state('main', {
        templateUrl: 'main.html',
        url: '/',
        controller: 'mainCtrl'}
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

app.filter('inArray', function ($filter) {
    return function (list, arrayFilter, element) {
        if (arrayFilter) {
            return $filter('filter')(list, function (listItem) {
                return arrayFilter.indexOf(listItem[element]) !== -1;
            });
        }
    };
});
