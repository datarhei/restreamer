/**
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
'use strict';

window.angular.module('Header').controller('headerController', ['$scope', '$translate', 'loggerService', ($scope, $translate, loggerService) => {
    $scope.currentLocale = $translate.preferredLanguage();
    $scope.switchLanguage = (locale) => {
        $scope.currentLocale = locale;
        $translate.use(locale).then(
            () => {
                loggerService.info('Switched language to ' + locale);
            },
            (error) => {
                loggerService.error('INFO', 'Switching language to ' + locale + ' failed: ' + error);
            });
    };
    $scope.langIs = function langIs (locale) {
        return locale === $scope.currentLocale;
    };
}]);
