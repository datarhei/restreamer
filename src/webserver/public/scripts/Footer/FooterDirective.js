/**
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
'use strict';

//import {footerModule} from 'FooterModule';

window.angular.module('Footer').directive('footer', () => {
    return {
        'restrict': 'A',
        'replace': true,
        'templateUrl': '/scripts/Footer/_footer.html',
        'controller': 'footerController'
    };
});
