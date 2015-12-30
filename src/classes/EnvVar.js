/*
 * @name ReStreamer
 * @namespace https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */

'use strict';
/**
 * Class representing a environment variable with default values
 */
class EnvVar {
    constructor(name, required, defaultValue, description) {
        this.name = name;
        this.required = required;
        this.defaultValue = defaultValue;
        this.description = description;
    }
}
module.exports = EnvVar;
