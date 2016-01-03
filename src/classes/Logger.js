/*
 * @name Restreamer
 * @namespace https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */

'use strict';

/*
Requirements
 */
const moment = require("moment-timezone");
require("colors");

/*
global vars
 */
const LEVEL_ERROR = 1;
const LEVEL_WARN = 2;
const LEVEL_INFO = 3;
const LEVEL_DEBUG = 4;

//set default timezone
process.env.TIMEZONE = process.env.TIMEZONE ? process.env.TIMEZONE : "Europe/Berlin";

/**
 * Class representing a logger
 */
class Logger {

    /**
     * check if the logger is muted
     * @returns {boolean}
     */
    static isMuted() {
        return !!process.env.LOGGER_MUTED;
    }

    /**
     * construct a logger object
     * @param {string} context context of the log message (classname.methodname)
     */
    constructor(context) {
        this.context = context;
    }

    /**
     * print a message to stdout
     * @param {string} message
     * @param {string} context
     * @param {string} type
     * @param {string} color
     */
    stdout(message, context, type, color) {
        if (Logger.isMuted()) {
            return;
        }
        if (context === false){
            context = "";
        }else{
            context = "(" + context + ")";
        }
        var str = "[" + (moment().tz(process.env.TIMEZONE).format('DD-MM-YYYY HH:mm:ss.SSS')) + "] [" + type + "] " + message + " " + context;
        if (process.env.LOGGER_LEVEL === '4') {
            console.log(str[color]);
        } else {
            console.log(str);
        }
    }

    /**
     * print an info message if LOG_LEVEL >= LEVEL_INFO
     * @param {string} message
     * @param {string} context
     */
    info(message, context, alertGui) {
        if (typeof context === 'undefined') {
            context = this.context;
        }
        if (typeof alertGui === 'undefined') {
            alertGui = false;
        }
        if (process.env.LOGGER_LEVEL >= LEVEL_INFO) {
            return this.stdout(message, context, "INFO", "blue");
        }
        if (alertGui){
            //todo: if alertGui is activated on frontend and websocketcontroller, insert emit here
        }
    }

    /**
     * print a warning message if LOG_LEVEL >= LEVEL_WARN
     * @param {string} message
     * @param {string} context
     */
    warn(message, context, alertGui) {
        if (typeof context === 'undefined') {
            context = this.context;
        }
        if (typeof alertGui === 'undefined') {
            alertGui = false;
        }
        if (process.env.LOGGER_LEVEL >= LEVEL_WARN) {
            return this.stdout(message, context, "WARN", "orange");
        }
       if (alertGui){
            //todo: if alertGui is activated on frontend and websocketcontroller, insert emit here
        }
    }

    /**
     * print a debug message if LOG_LEVEL >= LEVEL_DEBUG
     * @param {string} message
     * @param {string} context
     */
    debug(message, context, alertGui) {
        if (typeof context === 'undefined') {
            context = this.context;
        }
        if (typeof alertGui === 'undefined') {
            alertGui = false;
        }
        if (process.env.LOGGER_LEVEL >= LEVEL_DEBUG) {
            return this.stdout(message, context, "DEBUG", "yellow");
        }
        if (alertGui){
            //todo: if alertGui is activated on frontend and websocketcontroller, insert emit here
        }
    }

    /**
     * print a debug message if LOG_LEVEL >= LEVEL_ERROR
     * sends a string to
     * @param {string} message
     * @param {string} context
     */
    error(message, context, alertGui) {
        if (typeof context === 'undefined') {
            context = this.context;
        }
        if (typeof alertGui === 'undefined') {
            alertGui = false;
        }
        if (process.env.LOGGER_LEVEL >= LEVEL_ERROR) {
            return this.stdout(message, context, "ERROR", "red");
        }
        if (alertGui){
            //todo: if alertGui is activated on frontend and websocketcontroller, insert emit here
        }
    }
}

Logger.LEVEL_ERROR = LEVEL_ERROR;
Logger.LEVEL_WARN = LEVEL_WARN;
Logger.LEVEL_INFO = LEVEL_INFO;
Logger.LEVEL_DEBUG = LEVEL_DEBUG;

module.exports = Logger;
