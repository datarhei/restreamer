/**
 * @file holds the code for the class Logger
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
'use strict';

const moment = require('moment-timezone');
const LEVEL_ERROR = 1;
const LEVEL_WARN = 2;
const LEVEL_INFO = 3;
const LEVEL_DEBUG = 4;

// set default timezone to use the timezone before the default values are
process.env.TIMEZONE = process.env.TIMEZONE ? process.env.TIMEZONE : 'Europe/Berlin';
if (typeof process.env.LOGGER_LEVEL === 'undefined') {
    process.env.LOGGER_LEVEL = '3';
}

/**
 * Class for logger
 */
class Logger {

    /**
     * check if the logger is muted
     * @returns {boolean}
     */
    static isMuted () {
        return Boolean(process.env.LOGGER_MUTED);
    }

    /**
     * construct a logger object
     * @param {string} context context of the log message (classname.methodname)
     */
    constructor (context) {
        this.context = context;
    }

    /**
     * print a message to stdout
     * @param {string} message
     * @param {string} context
     * @param {string} type
     */
    stdout (message, context, type) {
        var time = moment().tz(process.env.TIMEZONE).format('DD-MM-YYYY HH:mm:ss.SSS');
        var loggerContext = `(${String(context)}`;

        if (Logger.isMuted()) {
            return;
        }

        console.log(`[${time}] [${type}] [${message}] [${loggerContext}]`);
    }

    /**
     * print an info message if LOG_LEVEL >= LEVEL_INFO
     * @param {string} message
     * @param {string} context
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

        if (process.env.LOGGER_LEVEL >= LEVEL_INFO) {
            return this.stdout(message, loggerContext, 'INFO');
        }

        // todo: if alertGui is activated on frontend and websocketcontroller, insert emit here
        if (loggerAlertGui) {
            return;
        }
    }

    /**
     * print a warning message if LOG_LEVEL >= LEVEL_WARN
     * @param {string} message
     * @param {string} context
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

        if (process.env.LOGGER_LEVEL >= LEVEL_WARN) {
            return this.stdout(message, loggerContext, 'WARN');
        }

        // todo: if alertGui is activated on frontend and websocketcontroller, insert emit here
        if (loggerAlertGui) {
            return;
        }
    }

    /**
     * print a debug message if LOG_LEVEL >= LEVEL_DEBUG
     * @param {string} message
     * @param {string} context
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

        if (process.env.LOGGER_LEVEL >= LEVEL_DEBUG) {
            return this.stdout(message, loggerContext, 'DEBUG');
        }

        // todo: if alertGui is activated on frontend and websocketcontroller, insert emit here
        if (loggerAlertGui) {
            return;
        }
    }

    /**
     * print a debug message if LOG_LEVEL >= LEVEL_ERROR
     * sends a string to
     * @param {string} message
     * @param {string} context
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

        if (process.env.LOGGER_LEVEL >= LEVEL_ERROR) {
            return this.stdout(message, loggerContext, 'ERROR');
        }

        // todo: if alertGui is activated on frontend and websocketcontroller, insert emit here
        if (loggerAlertGui) {
            return;
        }
    }
}

// define loglevels in logger class
Logger.LEVEL_ERROR = LEVEL_ERROR;
Logger.LEVEL_WARN = LEVEL_WARN;
Logger.LEVEL_INFO = LEVEL_INFO;
Logger.LEVEL_DEBUG = LEVEL_DEBUG;

module.exports = (context) => {
    return new Logger(context);
};
