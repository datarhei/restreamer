/**
 * @file config for passport local strategy
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
'use strict';

var LocalStrategy = require('passport-local').Strategy;
var auth = require('../../../conf/live.json').auth;

if (process.env.RESTREAMER_USERNAME) {
    var username = process.env.RESTREAMER_USERNAME;

} else {
    var username = auth.username;

}

if (process.env.RESTREAMER_USERNAME) {
    var password = process.env.RESTREAMER_PASSWORD;

} else {
    var password = auth.password;

}

module.exports = (passport) => {

    // used to serialize the user for the session
    passport.serializeUser(function (user, done) {
        done(null, user);
    });
    passport.deserializeUser(function (user, done) {
        done(null, user);
    });
    passport.use('local-login', new LocalStrategy({
        usernameField: 'user',
        passwordField: 'pass',

        // allows us to pass back the entire request to the callback
        passReqToCallback: true
    },

    // callback with user and pass from our form
    function (req, user, pass, done) {

        // login success
        if (user === username && pass === password) {

            /*
             * WEBSOCKET SECURITY HERE
             */
            done(null, auth);
        } else {
            done(null, false, req.flash('wrong password or wrong user'));
        }
    }));
};
