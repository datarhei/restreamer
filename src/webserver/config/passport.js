/**
 * @file config for passport local strategy
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
'use strict';

var LocalStrategy = require('passport-local').Strategy;
var auth = require(require('path').join(global.__base, 'conf', 'live.json')).auth;

module.exports = (passport) => {
    // used to serialize the user for the session
    passport.serializeUser(function serializeUser (user, done) {
        done(null, user);
    });
    passport.deserializeUser(function deserializeUser (user, done) {
        done(null, user);
    });
    passport.use('local-login', new LocalStrategy(
        {
            'usernameField': 'user',
            'passwordField': 'pass',
            'passReqToCallback': true // allows us to pass back the entire request to the callback
        },
        // callback with user and pass from our form
        function checkLogin (req, user, pass, done) {
            var username = process.env.RS_USERNAME || auth.username;
            var password = process.env.RS_PASSWORD || auth.password;

            // login success
            if (user === username && pass === password) {
                // WEBSOCKET SECURITY HERE
                done(null, auth);
            } else {
                done(null, false);
            }
        }));
};
