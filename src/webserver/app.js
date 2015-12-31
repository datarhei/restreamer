/*
 * @name ReStreamer
 * @namespace https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */

'use strict';

/*
    libraries
 */
//auth stuff
const passport = require('passport');
const flash    = require('connect-flash');

//express
const express = require('express');
const session = require('express-session');
const i18n = require('i18n');
const favicon = require('serve-favicon');
const cookieParser = require('cookie-parser');
const bodyParser = require('body-parser');
const compression = require("compression");

//other
const path = require('path');
const Q = require("q");
const request = require("request");

/*
    modules
 */
const packageJson = require("../../package.json");
const Logger = require("../classes/Logger");
const logger = new Logger("Webserver");
//middlewares
const expressLogger = require("./middlewares/expressLogger");

//create express app
const app = express();

app.use(session({
    secret: 'keyboard cat',
    resave: false,
    saveUninitialized: true})); // session secret

//create promise for "websockets ready"
app.set("websocketsReady", Q.defer());

//add passport auth
app.use(passport.initialize());
app.use(flash()); // use connect-flash for flash messages stored in session
app.use(passport.session()); // persistent login sessions
require("./config/passport")(passport);

//configure i18n
i18n.configure({
    locales:['en','de'],
    directory: __dirname + '/locales',
    defaultLocale: 'en',
    cookie: 'locale',
    updateFiles: false
});

//configure express app
app.use(favicon(__dirname + '/public/images/favicon.ico'));
app.use(i18n.init);
app.use(bodyParser.json());
app.set('views', __dirname + '/views');
app.set('view engine', 'jade');
app.use(bodyParser.urlencoded({ extended: true }));
app.use(cookieParser());
app.use(compression());
app.set('json spaces', 4);

require("./controllers/index")(app, passport);
app.use('/', expressLogger);

// catch 404 and forward to error handler
app.use(function(req, res, next) {
    var err = new Error('Not Found ' + req.url);
    err.status = 404;
    next(err);
});

// no stacktraces leaked to user
app.use(function(err, req, res, next) {
    res.status(err.status || 500);
    res.send({
        message: err.message,
        error: {}
    });
});

var checkForAppUpdates = function(){
    logger.debug("checking app for updates...");
    const url = 'http://datarhei.org/apps.json';
    request(url, function (error, response, body) {
        if (!error && response.statusCode === 200) {
            var updateCheck = JSON.parse(body);
            var updateAvailable = false;
            if (updateCheck.restreamer.version === packageJson.version){
                updateAvailable = false;
                logger.debug("checking app for updates successful. update is not available (remote: " + updateCheck.restreamer.version + ", local: " + packageJson.version + ")");
            }else{
                updateAvailable = updateCheck.restreamer.version;
                logger.debug("checking app for updates successful. update is available (remote: " + updateCheck.restreamer.version + ", local: " + packageJson.version + ")");
            }
            logger.info("checking app for updates successful.");
            app.set("updateAvailable", updateAvailable);
        } else {
            logger.info("Update check failed", false);
        }
    });
};

//start interval to check for updates
checkForAppUpdates();
setInterval(checkForAppUpdates, 24 * 60 * 60 * 1000);

module.exports = app;
