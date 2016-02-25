/**
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
'use strict';

window.angular.module('Header', []);

window.angular.module('Header').config(($translateProvider) => {
    $translateProvider.useStaticFilesLoader({
        'prefix': 'locales/lang-',
        'suffix': '.json'
    });

    $translateProvider.useSanitizeValueStrategy('escape');
    $translateProvider.preferredLanguage('en_US');
});
