/*
 * @name Restreamer
 * @namespace https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */

window.Logger = {

    level: {
        INFO: 'color: #0000FF; font-weight: bold',
        DEBUG: 'color: #AABBCC; font-weight: bold',
        ERROR: 'color: #FF0011d; font-weight: bold',
        WEBSOCKETS_IN: 'color: #00BFFF; font-weight: bold',
        WEBSOCKETS_OUT: 'color: #00BF00; font-weight: bold',
        WEBSOCKETS_NAMESPACE: 'color: #00BF00; font-weight: bold'
    },
    log: function(level, message){
        console.log("%c " +"[" + level + "]  " + message, this.level[level]);
    }
};
