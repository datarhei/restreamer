/**
 * @file controller for routing from /
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
/* eslint no-unused-vars: 0 */
'use strict';

const path = require('path');
var auth = require(require('path').join(global.__base, 'conf', 'live.json')).auth;

module.exports = (app) => {
    /* Handle Login POST */
    app.post('/login', (req, res, next) => {
        var username = process.env.RS_USERNAME || auth.username;
        var password = process.env.RS_PASSWORD || auth.password;
        var success = false;
        var message = '';
        if (req.body.user === username && req.body.pass === password) {
            req.session.authenticated = true;
            success = true;
        } else {
            message = 'login_invalid';
            req.session.destroy();
            success = false;
        }
        res.json({
            'success': success,
            'message': message
        });
    });
    app.get('/authenticated', (req, res) => {
        res.json(req.session.authenticated === true);
    });
    app.get('/logout', (req, res) => {
        req.session.destroy();
        res.end();
    });
    /* Handle NGINX-RTMP token */
    app.get('/token', (req, res) => {
        var token = process.env.RS_TOKEN || auth.token;
        if (token != '') {
            if (req.query.token == token) {
                res.writeHead(200, {
                    'Content-Type': 'text/plain'
                });
                res.end('Authorized');
            } else {
                res.writeHead(401, {
                    'Content-Type': 'text/plain'
                });
                res.end('Unauthorized');
            }
        } else {
            res.writeHead(200, {
                'Content-Type': 'text/plain'
            });
            res.end('Authorized');
        }
    });
};



