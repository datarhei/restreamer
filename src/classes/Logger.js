/**
 * @file holds the code for the class Logger
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
'use strict';

const moment = require('moment-timezone');
const printf = require('printf');

const LEVEL_MUTE = 0;
const LEVEL_ERROR = 1;
const LEVEL_WARN = 2;
const LEVEL_INFO = 3;
const LEVEL_DEBUG = 4;

// set default timezone to use the timezone before the default values are
// @todo: it is really ugly and wrong to log with hardcoded timezone before environment is read
process.env.RS_TIMEZONE = process.env.RS_TIMEZONE || 'Europe/Berlin';
process.env.RS_LOGGER_LEVEL = process.env.RS_LOGGER_LEVEL || 3;

/**
 * Class for logger
 */
class Logger {

    /**
     * check if the logger is muted
     * @returns {boolean}
     */
    static isMuted () {
        return process.env.RS_LOGGER_LEVEL === LEVEL_MUTE;
    }

    /**
     * construct a logger object
     * @param {string} context context of the log message (classname.methodname)
     */
    constructor (context) {
        process.env.RS_LOGGER_LEVEL = process.env.RS_LOGGER_LEVEL || LEVEL_INFO;
        this.context = context;
    }

    /**
     * print a message to stdout
     * @param {string} message
     * @param {string} context
     * @param {string} type
     */
    stdout (message, context, type) {
        var time = moment().tz(process.env.RS_TIMEZONE).format('DD-MM-YYYY HH:mm:ss.SSS');

        if(Logger.isMuted()) {
            return;
        }

        let logline = '';
        if(context) {
            logline = printf('[%s] [%-5s] [%22s] %s', time, type, context, message);
        } else {
            logline = printf('[%s] [%-5s] %s', time, type, message);
        }

        process.stdout.write(logline + '\n');
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

        if (process.env.RS_LOGGER_LEVEL >= LEVEL_INFO) {
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

        if (process.env.RS_LOGGER_LEVEL >= LEVEL_WARN) {
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

        if (process.env.RS_LOGGER_LEVEL >= LEVEL_DEBUG) {
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

        if (process.env.RS_LOGGER_LEVEL >= LEVEL_ERROR) {
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
