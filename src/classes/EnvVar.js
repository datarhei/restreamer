/**
 * @file holds the code for the class EnvVar
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
'use strict';

const logger = require('./Logger')('ENV');
const logBlacklist = ['RS_PASSWORD'];

/**
 * Class for environment variables with default values
 */
class EnvVar {
    static init(config) {
        var killProcess = false;

        // Cycle through all defined environment variables
        for(let envVar of config.envVars) {
            // Check if the environment variable is set. If not, cycle through the aliases.
            if(!(envVar.name in process.env)) {
                for(let i in envVar.alias) {
                    let alias = envVar.alias[i];
                    // If the alias exists, copy it to the actual name and delete it.
                    if(alias in process.env) {
                        logger.warn('The use of ' + alias + ' is deprecated. Please use ' + envVar.name + ' instead')
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

                logger.info(envVar.name + ' = ' + value + ' - ' + envVar.description);
            }
            else {
                if(envVar.required == true) {
                    logger.error(envVar.name + ' not set but required');
                    killProcess = true;
                }
                else {
                    logger.info(envVar.name + ' = ' + envVar.defaultValue + ' (using default) - ' + envVar.description);
                    process.env[envVar.name] = envVar.defaultValue;
                }
            }
        }

        if(killProcess == true) {
            setTimeout(()=> {
                process.exit();
            }, 500);
        }
    }
}

module.exports = EnvVar;
