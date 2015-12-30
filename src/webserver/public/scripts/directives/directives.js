/*
 * @name ReStreamer
 * @namespace https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */

/* jslint browser: true */

'use strict';

window.datarheiApp
.directive('ngModalNewinput', function() {
    return {
        templateUrl: 'ngDirectives/modals/newinput.html'
    };
});
const RTMP_REGEXP = /^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$/;
window.datarheiApp
    .directive('rtmpaddress', function($q, $timeout) {
    return {
        require: 'ngModel',
        link: function(scope, elm, attrs, ctrl) {
            ctrl.$validators.integer = function(modelValue, viewValue) {
                if (ctrl.$isEmpty(modelValue)) {
                    // consider empty models to be valid
                    return true;
                }
                if (RTMP_REGEXP.test(viewValue)) {
                    return true;
                }
                return false;
            };
        }
    };
});
