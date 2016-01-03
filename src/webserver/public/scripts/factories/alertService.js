/*
 * @name Restreamer
 * @namespace https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */

window.datarheiApp.factory('alertService', function($rootScope) {
    var alertService = {};
    $rootScope.alerts = [];
    alertService.add = function(alertObject) {
        var alert, found, _i, _len, _ref;
        found = false;
        for (var i = 0; i < $rootScope.alerts.length; i++) {
            alert = $rootScope.alerts[i];
            if (alert.message === alertObject.message)
            {
                found = true;
                alert.time = alertObject.time;
            }
        }
        if (!found) {
            alertObject.close = function() {
                return alertService.closeAlert(this);
            };
            $rootScope.alerts.push(alertObject);
            window.setTimeout(function() {
                alertService.closeAlertIdx($rootScope.alerts.indexOf(alertObject));
                return $rootScope.$apply();
            }, 10000);
            console.log($rootScope.alerts);
        }
    };
    alertService.closeAlert = function(alert) {
        return this.closeAlertIdx($rootScope.alerts.indexOf(alert));
    };
    alertService.closeAlertIdx = function(index) {
        return $rootScope.alerts.splice(index, 1);
    };
    return alertService;
});
