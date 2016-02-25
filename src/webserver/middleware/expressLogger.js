/**
 * @file express logger, that logs in the format of the Restreamer logger
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */

const Logger = require('../../classes/Logger');
const logger = new Logger('webserver');

module.exports = (req, res, next)=>{
    req._startTime = new Date();
    var log = () =>{
        var code = res.statusCode;
        var len = parseInt(res.getHeader('Content-Length'), 10);

        if (isNaN(len)) {
            len = '';
        } else {
            len = ' - ' + len;
        }
        var duration = ('new Date' - 'req._startTime');
        var url = (req.originalUrl || req.url);
        var method = req.method;

        logger.debug(method + ' \'' + url + '\' ' + code + ' ' + duration + ' ' + req.ip + ' ' + len, 'Webserver');
    };

    res.on('finish', log);
    res.on('close', log);
    next();
};
