/**
 * @file controller for routing from /
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
'use strict';

const express = require('express');
const path = require('path');
const Restreamer = require('../../classes/Restreamer');

module.exports = (app, passport) => {

    // static paths
    app.get('/main.html', (req, res)=>{
        if (req.isAuthenticated()) {
            res.sendFile(path.join(__dirname, '../', 'public', 'main.html'));
        } else {
            res.sendFile(path.join(__dirname, '../', 'public', 'login.html'));
        }
    });

    /* Handle Login POST */
    app.post('/login',
        passport.authenticate('local-login', {
            successRedirect: '/',
            failureRedirect: '/',
            failureFlash : true
        })
    );
    app.get('/logout', (req, res)=>{
        req.logout();
        res.redirect('/');
    });

    // small get.api
    app.get('/v1/states', (req, res)=>{
        var states = Restreamer.data.states;

        res.json({
            repeat_to_local_nginx: states.repeatToLocalNginx,
            repeat_to_optional_output: states.repeatToOptionalOutput
        });
    });
    app.get('/v1/progresses', (req, res)=>{
        var progresses = Restreamer.data.progresses;

        res.json({
            repeat_to_local_nginx: {
                frames: progresses.repeatToLocalNginx.frames,
                current_fps: progresses.repeatToLocalNginx.currentFps,
                current_kbps: progresses.repeatToLocalNginx.currentKbps,
                target_size: progresses.repeatToLocalNginx.targetSize,
                timemark: progresses.repeatToLocalNginx.timemark
            },
            repeat_to_optional_output: {
                frames: progresses.repeatToOptionalOutput.frames,
                current_fps: progresses.repeatToOptionalOutput.currentFps,
                current_kbps: progresses.repeatToOptionalOutput.currentKbps,
                target_size: progresses.repeatToOptionalOutput.targetSize,
                timemark: progresses.repeatToOptionalOutput.timemark
            }
        });
    });
};
