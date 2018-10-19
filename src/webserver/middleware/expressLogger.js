/**
 * @file express logger, that logs in the format of the Restreamer logger
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
'use strict';
/* eslint vars-on-top: 0 */

const logger = require.main.require('./classes/Logger')('Webserver');

module.exports = (req, res, next) => {
    req._startTime = new Date();
    var log = () => {
        let code = res.statusCode;
        let len = parseInt(res.getHeader('Content-Length'), 10);
        let duration = new Date() - req._startTime;
        let url = (req.originalUrl || req.url);
        let method = req.method;

        if (isNaN(len)) {
            len = '-';
        }

        logger.debug(method + ' "' + url + '" ' + code + ' ' + duration + 'ms ' + req.ip + ' ' + len);
    };

    res.on('finish', log);
    res.on('close', log);
    next();
};
