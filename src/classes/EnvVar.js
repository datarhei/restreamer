/**
 * @file holds the code for the class EnvVar
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
'use strict';
const logger = require('./Logger')('EnvVar');
const logBlacklist = ['RS_PASSWORD'];

/**
 * Class for environment variables with default values
 */
class EnvVar {
    static init (config) {
        var killProcess = false;

        for (let envVar of config.envVars) {
            if (typeof process.env[envVar.alias] !== 'undefined') {
                process.env[envVar.name] = process.env[envVar.alias];
                delete process.env[envVar.alias];
            }
            if (typeof process.env[envVar.name] !== 'undefined') {
                logger.info(`ENV "${envVar.name} = ${(logBlacklist.indexOf(envVar.name) === -1 ? process.env[envVar.name] : '[hidden]')}"`, envVar.description);
            } else if (envVar.required === true) {
                logger.error(`No value set for env "${envVar.name}", but it is required`);
                killProcess = true;
            } else {
                process.env[envVar.name] = envVar.defaultValue;
                logger.info(`ENV "${envVar.name} = ${(logBlacklist.indexOf(envVar.name) === -1 ? process.env[envVar.name] : '[hidden]')}", set to default value`, envVar.description);
            }

            if (typeof process.env[envVar.name] !== 'undefined') {
                switch (envVar.type) {
                    case 'int':
                        process.env[envVar.name] = parseInt(process.env[envVar.name], 10);
                        break;
                    case 'bool':
                        process.env[envVar.name] = process.env[envVar.name] === 'true';
                        break;
                    default: // keep strings
                        break;
                }
            }
        }

        if (killProcess === true) {
            setTimeout(()=> {
                process.exit();
            }, 500);
        }
    }
}

module.exports = EnvVar;
