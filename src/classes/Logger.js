/**
 * @file holds the code for the class Logger
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
'use strict';

const moment = require('moment-timezone');
const printf = require('printf');
const fs = require('fs');

const LEVEL_MUTE = 0;
const LEVEL_ERROR = 1;
const LEVEL_WARN = 2;
const LEVEL_INFO = 3;
const LEVEL_DEBUG = 4;

// set default timezone to use the timezone before the default values are
// @todo: it is really ugly and wrong to log with hardcoded timezone before environment is read
process.env.RS_TIMEZONE = process.env.RS_TIMEZONE || 'Europe/Berlin';
process.env.RS_LOGLEVEL = process.env.RS_LOGLEVEL || 3;

/**
 * Class for logger
 */
class Logger {

    /**
     * check if the logger is muted
     * @returns {boolean}
     */
    static isMuted () {
        return process.env.RS_LOGLEVEL === LEVEL_MUTE;
    }

    /**
     * construct a logger object
     * @param {string} context context of the log message (classname.methodname)
     */
    constructor (context) {
        process.env.RS_LOGLEVEL = process.env.RS_LOGLEVEL || LEVEL_INFO;
        this.context = context;

        this.debuglog = null;

        if(process.env.RS_DEBUG == 'true') {
            let identifier = process.pid + '-' + process.platform + '-' + process.arch;
            try {
                this.debuglog = fs.openSync('/restreamer/src/webserver/public/debug/Restreamer-' + identifier + '.txt', 'a');
            } catch(err) {
                this.debuglog = null;
                this.stdout('Error opening debug file ' + identifier + ': ' + err, false, 'INFO');
            } finally {
                this.stdout('Enabled logging to ' + identifier, context, 'INFO');
            }
        }
    }

    logline(message, context, type) {
        let time = moment().tz(process.env.RS_TIMEZONE).format('DD-MM-YYYY HH:mm:ss.SSS');

        let logline = '';
        if(context) {
            logline = printf('[%s] [%-5s] [%22s] %s', time, type, context, message);
        } else {
            logline = printf('[%s] [%-5s] %s', time, type, message);
        }

        return logline;
    }

    /**
     * print a message to stdout
     * @param {string} message
     * @param {string} context
     * @param {string} type
     */
    stdout (message, context, type) {
        if(Logger.isMuted()) {
            return;
        }

        let logline = this.logline(message, context, type);

        process.stdout.write(logline + '\n');
    }

    /**
     * print a message to a file
     * @param {string} message
     * @param {string} context
     * @param {string} type
     */
    file (message, context, type) {
        let logline = this.logline(message, context, type);

        if(this.debuglog !== null) {
            fs.appendFile(this.debuglog, logline + '\n', 'utf8', (err) => {
                // ignore errors
                if(err) {
                    return;
                }

                fs.fsync(this.debuglog, (err) => {
                    return;
                });

                return;
            });
        }
    }

    /**
     * print an info message if LOG_LEVEL >= LEVEL_INFO
     * @param {string} message
     * @param {string=} context
     * @param {boolean=} alertGui
     */
    info (message, context, alertGui) {
        var loggerContext = context;
        var loggerAlertGui = alertGui;

        if (typeof context === 'undefined') {
            loggerContext = this.context;
        }

        if (typeof alertGui === 'undefined') {
            loggerAlertGui = false;
        }

        if(process.env.RS_DEBUG == 'true') {
            this.file(message, loggerContext, 'INFO');
        }

        if (process.env.RS_LOGLEVEL >= LEVEL_INFO) {
            return this.stdout(message, loggerContext, 'INFO');
        }

        // todo: if alertGui is activated on frontend and websockets controller, insert emit here
        if (loggerAlertGui) {
            return;
        }
    }

    /**
     * print a warning message if LOG_LEVEL >= LEVEL_WARN
     * @param {string} message
     * @param {string=} context
     * @param {boolean=} alertGui
     */
    warn (message, context, alertGui) {
        var loggerContext = context;
        var loggerAlertGui = alertGui;

        if (typeof context === 'undefined') {
            loggerContext = this.context;
        }

        if (typeof alertGui === 'undefined') {
            loggerAlertGui = false;
        }

        if(process.env.RS_DEBUG == 'true') {
            this.file(message, loggerContext, 'WARN');
        }

        if (process.env.RS_LOGLEVEL >= LEVEL_WARN) {
            return this.stdout(message, loggerContext, 'WARN');
        }

        // todo: if alertGui is activated on frontend and websockets controller, insert emit here
        if (loggerAlertGui) {
            return;
        }
    }

    /**
     * print a debug message if LOG_LEVEL >= LEVEL_DEBUG
     * @param {string} message
     * @param {string=} context
     * @param {boolean=} alertGui
     */
    debug (message, context, alertGui) {
        var loggerContext = context;
        var loggerAlertGui = alertGui;

        if (typeof context === 'undefined') {
            loggerContext = this.context;
        }

        if (typeof alertGui === 'undefined') {
            loggerAlertGui = false;
        }

        if(process.env.RS_DEBUG == 'true') {
            this.file(message, loggerContext, 'DEBUG');
        }

        if (process.env.RS_LOGLEVEL >= LEVEL_DEBUG) {
            return this.stdout(message, loggerContext, 'DEBUG');
        }

        // todo: if alertGui is activated on frontend and websockets controller, insert emit here
        if (loggerAlertGui) {
            return;
        }
    }

    /**
     * print a debug message if LOG_LEVEL >= LEVEL_ERROR
     * sends a string to
     * @param {string} message
     * @param {string=} context
     * @param {boolean=} alertGui
     */
    error (message, context, alertGui) {
        var loggerContext = context;
        var loggerAlertGui = alertGui;

        if (typeof context === 'undefined') {
            loggerContext = this.context;
        }

        if (typeof alertGui === 'undefined') {
            loggerAlertGui = false;
        }

        if(process.env.RS_DEBUG == 'true') {
            this.file(message, loggerContext, 'ERROR');
        }

        if (process.env.RS_LOGLEVEL >= LEVEL_ERROR) {
            return this.stdout(message, loggerContext, 'ERROR');
        }

        // todo: if alertGui is activated on frontend and websockets controller, insert emit here
        if (loggerAlertGui) {
            return;
        }
    }
}

// define log levels in logger class
Logger.LEVEL_ERROR = LEVEL_ERROR;
Logger.LEVEL_WARN = LEVEL_WARN;
Logger.LEVEL_INFO = LEVEL_INFO;
Logger.LEVEL_DEBUG = LEVEL_DEBUG;

module.exports = (context) => {
    return new Logger(context);
};
