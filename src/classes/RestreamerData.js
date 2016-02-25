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

class RestreamerData {

    static checkJSONDb () {
        var schemadata = {};
        var dbdata = {};
        var deferred = Q.defer();
        var readSchema = Q.nfcall(fs.readFile, path.join(__dirname, '../../', 'conf', 'jsondb_v1_schema.json'));
        var readDBFile = Q.nfcall(fs.readFile, path.join(__dirname, '../../', 'db', 'v1.json'));

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
                    deferred.resolve();
                }
            }).catch((error) => {
                var defaultStructure = {
                    'addresses': {
                        'srcAddress': '',
                        'optionalOutputAddress': ''
                    },
                    'states': {
                        'repeatToLocalNginx': {
                            'type': 'stopped'
                        },
                        'repeatToOptionalOutput': {
                            'type': 'stopped'
                        }
                    },
                    'userActions': {
                        'repeatToLocalNginx': 'stop',
                        'repeatToOptionalOutput': 'stop'
                    }
                };

                logger.debug(`Error reading "v1.db": ${error.toString()}`);
                fs.writeFileSync(path.join(__dirname, '../../', 'db', 'v1.json'), JSON.stringify(defaultStructure));
                deferred.resolve();
            });

        return deferred.promise;
    }
}

module.exports = RestreamerData;
