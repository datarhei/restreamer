/**
 * @file controller for routing from /
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
/* eslint no-unused-vars: 0 */
'use strict';

const path = require('path');

module.exports = (app, passport) => {
    // static paths
    app.get('/favicon.ico', (req, res) => {
        res.sendFile(path.join(global.__public, 'images', 'favicon.ico'));
    });

    app.get('/main.html', (req, res) => {
        if (req.isAuthenticated()) {
            res.sendFile(path.join(global.__public, 'main.html'));
        } else {
            res.sendFile(path.join(global.__public, 'login.html'));
        }
    });

    /* Handle Login POST */
    app.post('/login',
        passport.authenticate('local-login', {
            'successRedirect': '/',
            'failureRedirect': '/#/login_invalid'
        })
    );
    app.get('/logout', (req, res) => {
        req.logout();
        res.redirect('/');
    });
};
