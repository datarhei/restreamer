/**
 * @file this file is loaded on application start and inits the application
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
'use strict';

const logger = require('./classes/Logger')('start');
const EnvVar = require('./classes/EnvVar');
const packageJson = require('../package.json');
const config = require('../conf/live.json');
const nginxrtmp = require('./classes/Nginxrtmp')(config);
const Q = require('q');
const Restreamer = require('./classes/Restreamer');
const RestreamerData = require('./classes/RestreamerData');
const WebsocketController = require('./classes/WebsocketController');
const restreamerApp = require('./webserver/app');

/*
 show start message
 */
logger.info('     _       _             _           _ ', false);
logger.info('  __| | __ _| |_ __ _ _ __| |___   ___(_)', false);
logger.info(' / _  |/ _  | __/ _  |  __|  _  |/  _ | |', false);
logger.info('| (_| | (_| | || (_| | |  | | | |  __/| |', false);
logger.info('|_____|_____|_||_____|_|  |_| |_|____||_|', false);
logger.info('', false);
logger.info('Restreamer v' + packageJson.version, false);
logger.info('', false);
logger.info('ENVIRONMENTS', false);
logger.info('More informations in our Docs', false);
logger.info('', false);

/*
 * Environments
 */
const envList = [];

envList.push(new EnvVar('NODEJS_PORT', false, 3000, 'Webserver port of application'));
envList.push(new EnvVar('NODE_ENV', false, 'prod', 'Nodejs Environment'));
envList.push(new EnvVar('LOGGER_LEVEL', true, '3', 'Logger level to defined, what should be logged'));
envList.push(new EnvVar('TIMEZONE', true, 'Europe/Berlin', 'Set the timezone'));
envList.push(new EnvVar('SNAPSHOT_REFRESH_INTERVAL', false, 60000, 'Interval to create a new Snapshot'));
envList.push(new EnvVar('CREATE_HEAPDUMPS', false, 'false', 'Create Heapdumps of application'));

// check if all required envs are set and kill app if not
EnvVar.checkRequiredAppEnvs(envList);

// check for app updates
restreamerApp.checkForRestreamerUpdates();
setInterval(restreamerApp.checkForAppUpdates, 24 * 60 * 60 * 2000);

// add default websocket events, @todo this will be removed, when the new websocket workflow is implemented
WebsocketController.bindDefaultEvents();

// start the app
nginxrtmp.init()
    .then(()=>{
        return RestreamerData.checkJSONDb();
    })
    .then(()=>{
        return restreamerApp.startWebserver();
    })
    .then(() => {
        return Q.fcall(Restreamer.restoreFFMpegProcesses);
    })
    .then(()=>{
        restreamerApp.getPublicIp()
    })
    .catch((error)=>{
        console.log(error);
        let errormessage = `Error starting webserver and nginx for application: ${error}`;

        throw new Error(errormessage);
    });
