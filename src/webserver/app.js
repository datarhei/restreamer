/**
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */

/*
 * libraries
 */
// auth stuff
const passport = require('passport');
const flash = require('connect-flash');

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

/*
 * modules
 */
const packageJson = require('../../package.json');
const logger = require('../classes/Logger')('Webserver');

// middlewares
const expressLogger = require('./middlewares/expressLogger');

// create express app
const app = express();

// generate random key
const secretKey = crypto.randomBytes(16).toString('hex');

app.use(session({
    secret: secretKey,
    resave: false,

    // session secret
    saveUninitialized: true}));

// create promise for 'websockets ready'
app.set('websocketsReady', Q.defer());

// add passport auth
app.use(passport.initialize());

// use connect-flash for flash messages stored in session
app.use(flash());

// persistent login sessions
app.use(passport.session());
require('./config/passport')(passport);

// configure express app
app.use(bodyParser.json());
app.set('views', __dirname + '/views');
app.set('view engine', 'jade');
app.use(bodyParser.urlencoded({extended: true}));
app.use(cookieParser());
app.use(compression());
app.set('json spaces', 4);

require('./controllers/index')(app, passport);
app.use('/', expressLogger);

// catch 404 and forward to error handler
app.use(function (req, res, next) {
    var err = new Error('Not Found ' + req.url);

    err.status = 404;
    next(err);
});

// no stacktraces leaked to user
app.use(function (err, req, res, next) {
    res.status(err.status || 500);
    res.send({
        message: err.message,
        error: {}
    });
});

var checkForAppUpdates = function () {
    const url = {'host': 'datarhei.org', 'path': '/apps.json'};

    logger.debug('Checking app for updates...');
    https.get(url, function (response) {
        if (response.statusCode === 200) {
            response.on('data', function (body) {
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
                app.set('updateAvailable', updateAvailable);
            });
        } else {
            logger.info('Update check failed', false);
        }
    }).on('error', function () {
        logger.info('Update check failed', false);
    });
};

// start interval to check for updates
checkForAppUpdates();
setInterval(checkForAppUpdates, 24 * 60 * 60 * 1000);

module.exports = app;
