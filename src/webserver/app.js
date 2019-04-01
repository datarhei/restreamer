/**
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
'use strict';

// express
const express = require('express');
const session = require('express-session');
const cookie = require('cookie');
const cookieParser = require('cookie-parser');
const bodyParser = require('body-parser');
const compression = require('compression');

// other
const path = require('path');
const Q = require('q');
const crypto = require('crypto');

// modules
const logger = require.main.require('./classes/Logger')('Webserver');
const indexRouter = require('./controllers/index');
const apiV1 = require('./controllers/api/v1');

// middleware
const expressLogger = require('./middleware/expressLogger');

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
        this.sessionKey = 'restreamer-session';
        this.sessionStore = new session.MemoryStore();

        if (process.env.RS_NODEJS_ENV === 'dev') {
            this.initDev();
        } else {
            this.initProd();
        }
    }

    /**
     * use sessions for the express app
     */
    useSessions () {
        this.app.use(session({
            'resave': true,
            'saveUninitialized': false,
            'key': this.sessionKey,
            'secret': this.secretKey,
            'unset': 'destroy',
            'store': this.sessionStore
        }));
    }

    /**
     * add automatic parsers for the body
     */
    addParsers () {
        this.app.use(bodyParser.json());
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
    createPromiseForWebsockets () {
        this.app.set('websocketsReady', Q.defer());
    }

    /**
     * add the restreamer routes
     */
    addRoutes () {
        indexRouter(this.app);
        this.app.use('/v1', apiV1);
    }

    /**
     * add 404 error handling on pages, that have not been found
     */
    add404ErrorHandling () {
        this.app.use((req, res, next) => {
            var err = new Error('Not Found ' + req.url);
            err.status = 404;
            next(err);
        });
    }

    /**
     * add ability for internal server errors
     */
    add500ErrorHandling () {
        this.app.use((err, req, res, next) => {
            logger.error(err);
            res.status(err.status || 500);
            res.send({
                'message': err.message,
                'error': {}
            });
        });
    }

    /**
     * enable websocket session validation
     */
    secureSockets () {
        this.app.get('io').set('authorization', (handshakeData, accept) => {
            if (handshakeData.headers.cookie) {
                this.sessionStore.get(cookieParser.signedCookie(
                    cookie.parse(handshakeData.headers.cookie)[this.sessionKey], this.secretKey
                ), (err, s) => {
                    if (!err && s && s.authenticated) {
                        return accept(null, true);
                    }
                });
            } else {
                return accept(null, false);
            }
        });
    }

    /**
     * start the webserver and open the websocket
     * @returns {*|promise}
     */
    startWebserver () {
        var deferred = Q.defer();
        var server = null;

        logger.info('Starting ...');
        this.app.set('port', process.env.RS_NODEJS_PORT);
        server = this.app.listen(this.app.get('port'), ()=> {
            this.app.set('io', require('socket.io')(server, {path: '/socket.io'}));
            this.secureSockets();
            this.app.set('server', server.address());

            // promise to avoid ws binding before the webserver has been started
            this.app.get('websocketsReady').resolve(this.app.get('io'));
            logger.info('Running on port ' + process.env.RS_NODEJS_PORT);
            deferred.resolve(server.address().port);
        });

        return deferred.promise;
    }

    /**
     * stuff that have always to be added to the webapp
     */
    initAlways () {
        this.useSessions();
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
        logger.debug('Init webserver with PROD environment');
        this.initAlways();
        this.app.get('/', (req, res)=> {
            res.sendFile(path.join(global.__public, 'index.prod.html'));
        });
        this.add404ErrorHandling();
        this.add500ErrorHandling();
    }

    /**
     * dev config for the express app
     */
    initDev () {
        logger.debug('Init webserver with DEV environment');
        this.initAlways();
        this.app.get('/', (req, res)=> {
            res.sendFile(path.join(global.__public, 'index.dev.html'));
        });
        this.add404ErrorHandling();
        this.add500ErrorHandling();
    }
}

const restreamerApp = new RestreamerExpressApp();

module.exports = restreamerApp;
