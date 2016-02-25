/**
 * @file holds the code for the class EnvVar
 * @link https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */

/**
 * Class for environment variables with default values
 */
class EnvVar {

    /**
     * constructs an EnvVar
     * @param {string} name
     * @param {boolean} required
     * @param defaultValue
     * @param {string} description
     */
    constructor (name, required, defaultValue, description) {
        this.name = name;
        this.required = required;
        this.defaultValue = defaultValue;
        this.description = description;
    }
}

module.exports = EnvVar;
