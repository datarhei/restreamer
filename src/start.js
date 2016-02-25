/**
 * @file this file is loaded on application start and initializes the application
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
const logger = require('./classes/Logger')('start');
const EnvVar = require('./classes/EnvVar');
const packageJson = require('../package.json');
const app = require('./webserver/app');
const config = require('../conf/live.json');
const nginxrtmp = require('./classes/Nginxrtmp')(config);
const Q = require('q');
const Validator = require('jsonschema').Validator;
const fs = require('fs');
const path = require('path');

/*
init simple_streamer with environments
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
logger.info('More information in our Docs', false);
logger.info('', false);

// define environment variables
var env_vars = [];

env_vars.push(new EnvVar('NODEJS_PORT', false, 3000, 'Webserver port of application'));
env_vars.push(new EnvVar('LOGGER_LEVEL', true, '3', 'Logger level to defined, what should be logged'));
env_vars.push(new EnvVar('TIMEZONE', true, 'Europe/Berlin', 'Set the timezone'));
env_vars.push(new EnvVar('SNAPSHOT_REFRESH_INTERVAL', false, 60000, 'Interval to create a new Snapshot'));
env_vars.push(new EnvVar('CREATE_HEAPDUMPS', false, 'false', 'Create Heapdumps of application'));

// manage all environments
var killProcess = false;

for (let e of env_vars) {
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
logger.info('', false);

/**
 * check if the data from jsondb is valid against the Restreamer jsondb schema
 * @returns {promise}
 */
const checkJsonDb = function () {
    logger.info('Checking jsondb file...');

    var schemadata;
    var dbdata;
    var deferred = Q.defer();

    var readSchema = Q.nfcall(fs.readFile, path.join(__dirname, '../', 'conf', 'jsondb_v1_schema.json'));
    var readDBFile = Q.nfcall(fs.readFile, path.join(__dirname, '../', 'db', 'v1.json'));

    readSchema
    .then((s)=> {
        schemadata = JSON.parse(s.toString('utf8'));
        return readDBFile;
    })
    .then((d)=> {
        dbdata = JSON.parse(d.toString('utf8'));
        let v = new Validator();
        let instance = dbdata;
        let schema = schemadata;
        let validateResult = v.validate(instance, schema);

        if (validateResult.errors.length > 0) {
            logger.debug(`Validation error of v1.db: ${JSON.stringify(validateResult.errors)}`);
            throw new Error(JSON.stringify(validateResult.errors));
        } else {
            logger.debug('"v1.db" is valid');
            deferred.resolve();
        }
    }).catch(function (error) {
        logger.debug(`Error reading "v1.db": ${error.toString()}`);
        var defaultStructure = {
            addresses: {
                srcAddress: '',
                optionalOutputAddress: ''
            },
            states: {
                repeatToLocalNginx: {
                    type: 'stopped'
                },
                repeatToOptionalOutput: {
                    type: 'stopped'
                }
            },
            userActions: {
                repeatToLocalNginx: 'stop',
                repeatToOptionalOutput: 'stop'
            }
        };

        fs.writeFileSync(path.join(__dirname, '../', 'db', 'v1.json'), JSON.stringify(defaultStructure));
        deferred.resolve();
    });
    return deferred.promise;
};


/**
 * start the express webserver
 * @returns {promise}
 */
const startWebserver = function startWebserver () {
    var deferred = Q.defer();

    logger.info('Starting webserver...');
    app.set('port', process.env.NODEJS_PORT);
    var server = app.listen(app.get('port'), function () {
        require('./classes/WebsocketsController').bindDefaultEvents();
        app.set('io', require('socket.io')(server));
        app.set('server', server.address());

        // promise to determine if the webserver has been started to avoid ws binding before
        app.get('websocketsReady').resolve(app.get('io'));
        logger.info(`Webserver running on port ${process.env.NODEJS_PORT}`, 'start.webserver');
        deferred.resolve(server.address().port);
    });

    return deferred.promise;
};

/**
 * get the public ip of the server, on that the Restreamer is running
 */
const getPublicIp = function () {
    logger.info('Getting public ip...', 'start.publicip');
    var exec = require('child_process').exec;

    exec('public-ip', (err, stdout, stderr)=> {
        if (err) {
            logger.error(err);
        }
        app.set('publicIp', stdout.split('\n')[0]);
    });
};

/**
 * Restores FFmpeg processes, that have been spawned i.E. on the last app-start (stored on jsondb)
 * @returns {promise}
 */
const restoreFFMPEGProcesses = function () {
    var deferred = Q.defer();
    const Restreamer = require('./classes/Restreamer');

    logger.info('Restoring FFmpeg processes...', 'start.restore');
    Restreamer.restoreFFMpegProcesses();
    deferred.resolve();
    return deferred.promise;
};

/**
 * let's do it
 */

nginxrtmp.init()
    .then(checkJsonDb)
    .then(startWebserver)
    .then(checkJsonDb)
    .then(restoreFFMPEGProcesses)
    .then(getPublicIp)
    .catch(function (error) {
        logger.error(`Error starting webserver and nginx for application: ${error}`);
        setTimeout(()=> {
            process.exit();
        }, 500);
    });
