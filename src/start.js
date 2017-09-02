/**
 * @file this file is loaded on application start and initializes the application
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
'use strict';
const path = require('path');

global.__src = __dirname;
global.__base = path.join(__dirname, '..');
global.__public = path.join(__dirname, 'webserver', 'public');

const logger = require('./classes/Logger')('start');
const EnvVar = require('./classes/EnvVar');
const packageJson = require(path.join('..', 'package.json'));
const config = require(path.join(global.__base, 'conf', 'live.json'));
const nginxrtmp = require('./classes/Nginxrtmp')(config);
const Q = require('q');
const Restreamer = require('./classes/Restreamer');
const RestreamerData = require('./classes/RestreamerData');
const restreamerApp = require('./webserver/app');

// show start message
logger.info('     _       _             _           _ ', false);
logger.info('  __| | __ _| |_ __ _ _ __| |___   ___(_)', false);
logger.info(' / _  |/ _  | __/ _  |  __|  _  |/  _ | |', false);
logger.info('| (_| | (_| | || (_| | |  | | | |  __/| |', false);
logger.info('|_____|_____|_||_____|_|  |_| |_|____||_|', false);
logger.info('', false);
logger.info('Restreamer v' + packageJson.version, false);
logger.info('', false);
logger.info('ENVIRONMENTS', false);
logger.info('More information in our Docs', false);
logger.info('', false);

// setup environment vars
EnvVar.init(config);

// ffmpeg reporting
if (process.env.RS_FFMPEG_DEBUG === 'true') {
    process.env['FFREPORT'] = 'file=/restreamer/src/webserver/public/debug/%p-%t.log';
}

// start the app
nginxrtmp.init()
    .then(() => {
        return RestreamerData.checkJSONDb();
    })
    .then(() => {
        Restreamer.checkForUpdates();
        Restreamer.getPublicIp();
        Restreamer.bindWebsocketEvents();
        return restreamerApp.startWebserver();
    })
    .then(() => {
        return Q.fcall(Restreamer.restoreFFMpegProcesses);
    })
    .catch((error)=> {
        let errorMessage = `Error starting webserver and nginx for application: ${error}`;
        throw new Error(errorMessage);
    });
