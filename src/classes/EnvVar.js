/**
 * @file holds the code for the class EnvVar
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
'use strict';
const logger = require('./Logger')('EnvVar');

/**
 * Class for environment variables with default values
 */
class EnvVar {

    /**
     * constructs an envvar
     * @param {string} name
     * @param {string} required
     * @param {string} defaultValue
     * @param {string} description
     */
    constructor (name, required, defaultValue, description) {
        this.name = name;
        this.required = required;
        this.defaultValue = defaultValue;
        this.description = description;
    }

    /**
     * @param envList
     */
    static checkRequiredAppEnvs (envList) {
        var killProcess = false;

        for (let e of envList) {
            if (typeof process.env[e.name] !== 'undefined') {
                logger.info(`ENV "${e.name}=${process.env[e.name]}"`, e.description);
            } else if (e.required === true) {
                logger.error(`No value set for env "${e.name}", but it is required`);
                killProcess = true;
            } else {
                process.env[e.name] = e.defaultValue;
                logger.info(`ENV "${e.name}=${process.env[e.name]}", set to default-value!`, e.description);
            }
        }

        // kill process after a short delay to log the error message to console
        if (killProcess === true) {
            setTimeout(()=> {
                process.exit();
            }, 500);
        }
    }
}

module.exports = EnvVar;
