/**
 * @file holds the code for the class EnvVar
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */
'use strict';

const fs = require('fs');
const Q = require('q');
const path = require('path');
const Validator = require('jsonschema').Validator;

const logger = require('./Logger')('RestreamerData');

const dbPath = path.join(global.__base, 'db');
const dbFile = 'v1.json';

const confPath = path.join(global.__base, 'conf');
const schemaFile = 'jsondb_v1_schema.json';

class RestreamerData {

    static checkJSONDb () {
        var schemadata = {};
        var dbdata = {};
        var deferred = Q.defer();
        var readSchema = Q.nfcall(fs.readFile, path.join(confPath, schemaFile));
        var readDBFile = Q.nfcall(fs.readFile, path.join(dbPath, dbFile));

        logger.info('Checking jsondb file...');
        readSchema
            .then((s) => {
                schemadata = JSON.parse(s.toString('utf8'));
                return readDBFile;
            })
            .then((d) => {
                dbdata = JSON.parse(d.toString('utf8'));
                let v = new Validator();
                let instance = dbdata;
                let schema = schemadata;
                let validateResult = v.validate(instance, schema);

                if (validateResult.errors.length > 0) {
                    logger.debug(`Validation error of v1.db: ${JSON.stringify(validateResult.errors)}`);
                    throw new Error(JSON.stringify(validateResult.errors));
                } else {
                    logger.debug('"v1.db" is valid');

                    // Fill up optional fields if not present
                    if(!('video' in dbdata.options)) {
                        dbdata.options.video = {
                            codec: 'copy',
                            preset: 'ultrafast',
                            bitrate: '4096',
                            fps: '25',
                            profile: 'auto',
                            tune: 'none'
                        };
                    }

                    if(!('audio' in dbdata.options)) {
                        dbdata.options.audio = {
                            codec: 'auto',
                            preset: 'silence',
                            bitrate: '64',
                            channels: 'mono',
                            sampling: '44100'
                        };

                        // Update the defaults according to RS_AUDIO
                        switch(process.env.RS_AUDIO) {
                            case 'auto':
                                dbdata.options.audio.codec = 'auto';
                                break;
                            case 'none':
                                dbdata.options.audio.codec = 'none';
                                break;
                            case 'silence':
                                dbdata.options.audio.codec = 'aac';
                                dbdata.options.audio.preset = 'silence';
                                dbdata.options.audio.bitrate = '8';
                                dbdata.options.audio.channels = 'mono';
                                dbdata.options.audio.sampling = '44100';
                                break;
                            case 'aac':
                                dbdata.options.audio.codec = 'aac';
                                dbdata.options.audio.preset = 'encode';
                                dbdata.options.audio.bitrate = '64';
                                dbdata.options.audio.channels = 'inherit';
                                dbdata.options.audio.sampling = 'inherit';
                                break;
                            case 'mp3':
                                dbdata.options.audio.codec = 'mp3';
                                dbdata.options.audio.preset = 'encode';
                                dbdata.options.audio.bitrate = '64';
                                dbdata.options.audio.channels = 'inherit';
                                dbdata.options.audio.sampling = 'inherit';
                                break;
                            default:
                                break;
                        }
                    }

                    if(!('player' in dbdata.options)) {
                        dbdata.options.player = {
                            autoplay: false,
                            mute: false,
                            statistics: false,
                            color: '#3daa48',
                            logo: {
                                image: '',
                                position: 'bottom-right',
                                link: ''
                            }
                        };
                    }

                    if(!('output' in dbdata.options)) {
                        dbdata.options.output = {
                            type: 'rtmp',
                            rtmp: {},
                            hls: {
                                method: 'POST',
                                time: '2',
                                listSize: '10',
                                timeout: '10'
                            }
                        };
                    }

                    if(parseInt(dbdata.options.output.hls.timeout) > 2147) {
                        dbdata.options.output.hls.timeout = '10';
                    }

                    if (!fs.existsSync(dbPath)) {
                        fs.mkdirSync(dbPath);
                    }
                    fs.writeFileSync(path.join(dbPath, dbFile), JSON.stringify(dbdata));

                    deferred.resolve();
                }
            })
            .catch((error) => {
                var defaultStructure = {
                    addresses: {
                        srcAddress: ''
                    },
                    options: {
                        rtspTcp: true,
                        video: {
                            codec: 'copy',
                            preset: 'ultrafast',
                            bitrate: '4096',
                            fps: '25',
                            profile: 'auto',
                            tune: 'none'
                        },
                        audio: {
                            codec: 'auto',
                            preset: 'silence',
                            bitrate: '64',
                            channels: 'mono',
                            sampling: '44100'
                        },
                        player: {
                            autoplay: false,
                            mute: false,
                            statistics: false,
                            color: '#3daa48',
                            logo: {
                                image: '',
                                position: 'bottom-right',
                                link: ''
                            }
                        },
                        outputs: [{
                            type: 'rtmp',
                            rtmp: {},
                            hls: {
                                method: 'POST',
                                time: '2',
                                listSize: '10',
                                timeout: '10'
                            },
                            outputAddress: '',
                        }]
                    },
                    states: {
                        repeatToLocalNginx: {
                            type: 'stopped'
                        },
                        repeatToOptionalOutput_0: {
                            type: 'stopped'
                        }
                    },
                    userActions: {
                        repeatToLocalNginx: 'stop',
                        repeatToOptionalOutput_0: 'stop',
                    }
                };

                // Set stream source and start streaming on a fresh installation
                if(process.env.RS_INPUTSTREAM != '') {
                    defaultStructure.addresses.srcAddress = process.env.RS_INPUTSTREAM;
                    defaultStructure.states.repeatToLocalNginx.type = 'connected';
                    defaultStructure.userActions.repeatToLocalNginx = 'start';

                    // Set stream destination and start streaming on a fresh installation
                    if(process.env.RS_OUTPUTSTREAM != '') {
                        defaultStructure.addresses.optionalOutputAddress = process.env.RS_OUTPUTSTREAM;
                        defaultStructure.states.repeatToOptionalOutput_0.type = 'connected';
                        defaultStructure.userActions.repeatToOptionalOutput_0 = 'start';
                    }
                }

                // Adjust the structure for multiple outputs
                if (process.env.RS_EXTRA_OUTPUTS > 1) {
                    for (let i = 1; i < process.env.RS_EXTRA_OUTPUTS; i++) {
                        const id = 'repeatToOptionalOutput_' + i;
                        defaultStructure.states[id] = {'type': 'stopped'};
                        defaultStructure.userActions[id] = 'stop';
                        defaultStructure.options.outputs.push({
                            type: 'rtmp',
                            rtmp: {},
                            hls: {
                                method: 'POST',
                                time: '2',
                                listSize: '10',
                                timeout: '10'
                            },
                            outputAddress: ''
                        });
                    }
                }


                logger.debug(`Error reading "v1.db": ${error.toString()}`);
                if (!fs.existsSync(dbPath)) {
                    fs.mkdirSync(dbPath);
                }
                fs.writeFileSync(path.join(dbPath, dbFile), JSON.stringify(defaultStructure));
                deferred.resolve();
            });
        return deferred.promise;
    }
}

module.exports = RestreamerData;
