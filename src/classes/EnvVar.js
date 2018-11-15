/**
 * @file holds the code for the class EnvVar
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
'use strict';

const logBlacklist = ['RS_PASSWORD'];

/**
 * Class for environment variables with default values
 */
class EnvVar {
    constructor() {
        this.reset();
    }

    log(message, level) {
        this.messages.push({
            level: level,
            message: message
        });
    }

    init(config) {
        // Cycle through all defined environment variables
        for(let envVar of config.envVars) {
            // Check if the environment variable is set. If not, cycle through the aliases.
            if(!(envVar.name in process.env)) {
                for(let i in envVar.alias) {
                    let alias = envVar.alias[i];
                    // If the alias exists, copy it to the actual name and delete it.
                    if(alias in process.env) {
                        this.log('The use of ' + alias + ' is deprecated. Please use ' + envVar.name + ' instead', 'warn');
                        process.env[envVar.name] = process.env[alias];
                        delete process.env[alias];
                    }
                }
            }

            // Check if the environment variable is set and display it, if it is not set
            // apply the default value. In case the environment variable is required and
            // not set, stop the process.
            if(envVar.name in process.env) {
                // Adjust the given value to the required type
                switch(envVar.type) {
                    case 'int':
                        process.env[envVar.name] = parseInt(process.env[envVar.name], 10);
                        break;
                    case 'bool':
                        process.env[envVar.name] = process.env[envVar.name] == 'true';
                        break;
                    default: // keep strings
                        break;
                }

                // Cover blacklisted values
                let value = process.env[envVar.name];
                if(logBlacklist.indexOf(envVar.name) != -1) {
                    value = '******';
                }

                this.log(envVar.name + ' = ' + value + ' - ' + envVar.description, 'info');
            }
            else {
                if(envVar.required == true) {
                    this.log(envVar.name + ' not set, but required', 'error');
                    this.errors = true;
                }
                else {
                    this.log(envVar.name + ' = ' + envVar.defaultValue + ' (using default) - ' + envVar.description, 'info');
                    process.env[envVar.name] = envVar.defaultValue;
                }
            }
        }
    }

    list(logger) {
        for(let i = 0; i < this.messages.length; i++) {
            let m = this.messages[i];
            switch(m.level) {
                case 'info':
                    logger.info(m.message, 'ENV');
                    break;
                case 'warn':
                    logger.warn(m.message, 'ENV');
                    break;
                case 'error':
                    logger.error(m.message, 'ENV');
                    break;
                default:
                    break;
            }
        }

        this.messages = [];
    }

    hasErrors() {
        return this.errors;
    }

    reset() {
        this.messages = [];
        this.errors = false;
    }
}

module.exports = new EnvVar;
