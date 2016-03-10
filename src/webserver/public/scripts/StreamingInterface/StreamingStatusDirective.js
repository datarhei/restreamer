/**
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
'use strict';

window.angular.module('StreamingInterface').directive('streamingStatus', () => {
    return {
        'scope': {
            'data': '=',
            'name': '@name'
        },
        'restrict': 'E',
        'replace': true,
        'templateUrl': '/scripts/StreamingInterface/_streamingStatus.html',
        'controller': 'streamingStatusController'
    };
});
