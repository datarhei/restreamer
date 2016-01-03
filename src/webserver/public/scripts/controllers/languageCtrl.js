/*
 * @name Restreamer
 * @namespace https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */

window.datarheiApp.controller('languageCtrl',['$scope','$translate', function($scope, $translate) {
    $scope.currentLocale = $translate.preferredLanguage();
    $scope.switchLanguage = function(locale){
        $scope.currentLocale = locale;
        $translate.use(locale).then(function(){
            window.Logger.log("INFO", "Switched language to " + locale);
        }, function(error){
            window.Logger.error("INFO", "Switching language to " + locale + " failed: " + error);
        });
    };
    $scope.langIs = function(locale){
        return locale === $scope.currentLocale;
    };
}]);
