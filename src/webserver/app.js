/**
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
'use strict';

// auth stuff
const passport = require('passport');
const flash = require('connect-flash');
const passportConfig = require('./config/passport')

// express
const express = require('express');
const session = require('express-session');
const cookieParser = require('cookie-parser');
const bodyParser = require('body-parser');
const compression = require('compression');
const https = require('https');

// other
const path = require('path');
const Q = require('q');
const crypto = require('crypto');
const exec = require('child_process').exec;

// modules
const packageJson = require('../../package.json');
const logger = require('../classes/Logger')('RestreamerExpressApp');
const indexRouter = require('./controllers/index');

// middlewares
const expressLogger = require('./middlewares/expressLogger');

/**
 * Class for the ReStreamer webserver, powered by express.js
 */
class RestreamerExpressApp {

    /**
     * constructs a new express app with prod or dev config
     */
    constructor () {
        this.app = express();
        this.secretKey = crypto.randomBytes(16).toString('hex');

        if (process.env.NODE_ENV === 'prod') {
            this.initProd();
        }else {
            this.initDev();
        }
    }

    /**
     * use sessions for the express app
     */
    useSessions () {
        this.app.use(session({
            secret: this.secretKey,
            resave: false,

            // session secret
            saveUninitialized: true}));
    }

    /**
     * use passport auth
     */
    useAuth () {
        // add passport auth
        this.app.use(passport.initialize());

        // use connect-flash for flash messages stored in session
        this.app.use(flash());

        // persistent login sessions
        this.app.use(passport.session());

        // add config to passport
        passportConfig(passport);
    }

    /**
     * set the view engine
     * @todo should this be removed? we do not use any view engine by now, only statics
     */
    setViewEngine () {
        this.app.set('views', __dirname + '/views');
        this.app.set('view engine', 'jade');
    }

    /**
     * add automatic parsers for the body
     */
    addParsers () {
        this.app.use(bodyParser.json());
        this.app.use(bodyParser.urlencoded({extended: true}));
        this.app.use(cookieParser());
    }

    /**
     * add content compression on responses
     */
    addCompression () {
        this.app.use(compression());
    }

    /**
     * add express logger
     */
    addExpressLogger () {
        this.app.use('/', expressLogger);
    }

    /**
     * beautify json response
     */
    beautifyJSONResponse () {
        this.app.set('json spaces', 4);
    }

    /**
     * create a promise to check when websockets are ready for bindings
     */
    createPromiseForWebsockets(){
        this.app.set('websocketsReady', Q.defer());
    }

    /**
     * add the restreamer routes
     */
    addRoutes () {
        indexRouter(this.app, passport);
    }

    /**
     * add 404 error handling on pages, that have not been found
     */
    add404ErrorHandling () {
        this.app.use((req, res, next)=>{
            var err = new Error('Not Found ' + req.url);
            err.status = 404;
            next(err);
        });
    }

    /**
     * add ability for internal server errors
     */
    add500ErrorHandling () {
        this.app.use(function (err, req, res, next) {
            res.status(err.status || 500);
            res.send({
                message: err.message,
                error: {}
            });
        });
    }

    /**
     * check for app updates
     */
    checkForRestreamerUpdates () {

        const url = {'host': 'datarhei.org', 'path': '/apps.json'};
        logger.debug('Checking app for updates...');
        https.get(url, (response)=> {
            if (response.statusCode === 200) {
                response.on('data', (body)=>{
                    var updateCheck = JSON.parse(body);
                    var updateAvailable = false;

                    if (updateCheck.restreamer.version === packageJson.version) {
                        updateAvailable = false;
                        logger.debug(`Checking app for updates successful. Update is not available (remote: ${updateCheck.restreamer.version}, local: ${packageJson.version})`);
                    } else {
                        updateAvailable = updateCheck.restreamer.version;
                        logger.debug(`Checking app for updates successful. Update is available (remote: ${updateCheck.restreamer.version}, local: ${packageJson.version})`);
                    }
                    logger.info('Checking app for updates successfull');
                    this.app.set('updateAvailable', updateAvailable);
                });
            } else {
                logger.info('Update check failed', false);
            }
        }).on('error', function () {
            logger.info('Update check failed', false);
        });
    }

    /**
     * get public ip of the app
     */
    getPublicIp () {
        logger.info('Getting public ip...', 'start.publicip');
        exec('public-ip', (err, stdout, stderr)=> {
            if (err) {
                logger.error(err);
            }
            this.app.set('publicIp', stdout.split('\n')[0]);
        });
    }

    /**
     * start the webserver and open the websocket
     * @returns {*|promise}
     */
    startWebserver() {
        var deferred = Q.defer();

        logger.info('Starting webserver...');
        this.app.set('port', process.env.NODEJS_PORT);
        var server = this.app.listen(this.app.get('port'), ()=>{
            this.app.set('io', require('socket.io')(server));
            this.app.set('server', server.address());

            // promise to determine if the webserver has been started to avoid ws binding before
            this.app.get('websocketsReady').resolve(this.app.get('io'));
            logger.info(`Webserver running on port ${process.env.NODEJS_PORT}`);
            deferred.resolve(server.address().port);
        });

        return deferred.promise;
    }

    /**
     * stuff that have always to be added to the webapp
     */
    initAlways () {
        this.useSessions();
        this.useAuth();
        this.setViewEngine();
        this.addParsers();
        this.addCompression();
        this.addExpressLogger();
        this.beautifyJSONResponse();
        this.createPromiseForWebsockets();
        this.addRoutes();

    }

    /**
     * prod config for the express app
     */
    initProd () {
        logger.debug('init webserver with PROD environment');
        this.initAlways();
        this.app.get('/', (req, res)=>{
            res.sendFile(path.join(__dirname, 'public', 'index.html'));
        });
        this.add404ErrorHandling();
        this.add500ErrorHandling();
    }

    /**
     * dev config for the express app
     */
    initDev () {
        logger.debug('init webserver with DEV environment');
        this.initAlways();
        this.app.get('/', (req, res)=>{
            res.sendFile(path.join(__dirname, 'public', 'index-dev.html'));
        });
        this.add404ErrorHandling();
        this.add500ErrorHandling();
    }
}

const restreamerApp = new RestreamerExpressApp();

module.exports = restreamerApp;
