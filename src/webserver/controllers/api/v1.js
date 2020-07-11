/**
 * @file controller for routing from /v1
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
'use strict';

const express = require('express');
const router = new express.Router();
const version = require(require('path').join(global.__base, 'package.json')).version;

// TODO: solve the circular dependency problem and place Restreamer require here

router.get('/version', (req, res) => {
    res.json({
        'version': version,
        'update': require.main.require('./classes/Restreamer').data.updateAvailable
    });
});
router.get('/ip', (req, res) => {
    res.end(require.main.require('./classes/Restreamer').data.publicIp);
});
router.get('/states', (req, res) => {
    const states = require.main.require('./classes/Restreamer').data.states;

    const response = {
        'repeat_to_local_nginx': {
            type: states.repeatToLocalNginx.type,
            message: states.repeatToLocalNginx.message.replace(/\?token=[^\s]+/, '?token=***'),
        },
        'repeat_to_optional_output': {
            type: states.repeatToOptionalOutput_0.type,
            message: states.repeatToOptionalOutput_0.message.replace(/\?token=[^\s]+/, '?token=***'),
        },
    };

    res.json(response);
});
router.get('/progresses', (req, res) => {
    const progresses = require.main.require('./classes/Restreamer').data.progresses;

    res.json({
        'repeat_to_local_nginx': {
            'frames': progresses.repeatToLocalNginx.frames,
            'current_fps': progresses.repeatToLocalNginx.currentFps,
            'current_kbps': progresses.repeatToLocalNginx.currentKbps,
            'target_size': progresses.repeatToLocalNginx.targetSize,
            'timemark': progresses.repeatToLocalNginx.timemark
        },
        'repeat_to_optional_output': {
            'frames': progresses.repeatToOptionalOutput_0.frames,
            'current_fps': progresses.repeatToOptionalOutput_0.currentFps,
            'current_kbps': progresses.repeatToOptionalOutput_0.currentKbps,
            'target_size': progresses.repeatToOptionalOutput_0.targetSize,
            'timemark': progresses.repeatToOptionalOutput_0.timemark
        }
    });
});

module.exports = router;
