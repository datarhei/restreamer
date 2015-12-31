/*
 * @name ReStreamer
 * @namespace https://github.com/datarhei/restreamer
 * @copyright 2015 datarhei.org
 * @license Apache-2.0
 */

'use strict';

/*
requirements
 */
const Logger = require("./classes/Logger");
const logger = new Logger("start");
const EnvVar = require("./classes/EnvVar");
const packageJson = require("../package.json");
const app = require('./webserver/app');
const config = require('../conf/live.json');
const Q = require("q");

if(process.env.CREATE_HEAPDUMPS === 'true'){
    const heapdump = require('heapdump');
    setInterval(function(){
        heapdump.writeSnapshot('heapdump/' + Date.now() + '.heapsnapshot');
    }, 60000);
}

if (typeof process.env.LOGGER_LEVEL === 'undefined'){
    process.env.LOGGER_LEVEL = "3";
}

/*
init simple_streamer with environments
 */
logger.info("     _       _             _           _ ", false);
logger.info("  __| | __ _| |_ __ _ _ __| |___   ___(_)", false);
logger.info(" / _  |/ _  | __/ _  |  __|  _  |/  _ | |", false);
logger.info("| (_| | (_| | || (_| | |  | | | |  __/| |", false);
logger.info("|_____|_____|_||_____|_|  |_| |_|____||_|", false);
logger.info("", false);
logger.info("ReStreamer v" + packageJson.version, false);
logger.info("", false);
logger.info("ENVIRONMENTS", false);
logger.info("More informations in our Docs", false);
logger.info("", false);

// define environment variables
var env_vars = [];
env_vars.push(new EnvVar("NODEJS_PORT", false, 3000, "Webserver port of application"));
env_vars.push(new EnvVar("LOGGER_LEVEL", true, "3", "Logger level to defined, what should be logged"));
env_vars.push(new EnvVar("TIMEZONE", true, "Europe/Berlin", "Set the timezone"));
env_vars.push(new EnvVar("SNAPSHOT_REFRESH_INTERVAL", false, 60000, "Interval to create a new Snapshot"));
env_vars.push(new EnvVar("CREATE_HEAPDUMPS", false, "false", "Create Heapdumps of application"));

// manage all environments
var killProcess = false;
for (let e of env_vars){
    if(typeof process.env[e.name] !== 'undefined'){
        logger.info("ENV \"" + e.name + "=" + process.env[e.name] + "\", " + e.description);
    }else if(e.required === true){
        logger.error("No value set for env " + e.name + ", but it is required");
        killProcess = true;
    }else{
        process.env[e.name] = e.defaultValue;
        logger.info("ENV \"" + e.name + "=" + process.env[e.name] + "\", set to default-value!, " + e.description);
    }
}

// kill process after a short delay to log the error message to console
if(killProcess === true){
    setTimeout(()=> {
        process.exit();
    }, 500);
}
logger.info("", false);

/*
 checking jsondb
 */
const checkJsonDb = function(){
    logger.info("checking jsondb file...");
    var Validator = require('jsonschema').Validator;
    var fs = require("fs");
    var path = require("path");

    var schemadata;
    var dbdata;
    var deferred = Q.defer();

    var readSchema = Q.nfcall(fs.readFile, path.join(__dirname, "../", "conf", "jsondb_v1_schema.json"));
    var readDBFile = Q.nfcall(fs.readFile, path.join(__dirname, "../", "db", "v1.json"));

    readSchema
        .then(function(s){
            schemadata = JSON.parse(s.toString("utf8"));
            return readDBFile;
        })
        .then(function(d){
          dbdata = JSON.parse(d.toString("utf8"));
          var v = new Validator();
          var instance = dbdata;
          var schema = schemadata;
          var validateResult = v.validate(instance, schema);
          if (validateResult.errors.length > 0){
            logger.debug("validation error of v1.db: " + JSON.stringify(validateResult.errors));
            throw new Error(JSON.stringify(validateResult.errors));
          }else{
            logger.debug("v1.db is valid");
            deferred.resolve();
          }
    }).catch(function(error){
        logger.debug("error reading v1.db:" + error.toString());
        var defaultStructure = {
            addresses: {
                srcAddress: "",
                optionalOutputAddress: ""
            },
            states: {
                repeatToLocalNginx: {
                    type: "stopped"
                },
                repeatToOptionalOutput: {
                    type: "stopped"
                }
            },
            userActions:{
                repeatToLocalNginx: "stop",
                repeatToOptionalOutput: "stop"
            },
            progresses: {
                repeatToLocalNginx: {},
                repeatToOptionalOutput: {}
            }
        };
        fs.writeFileSync(path.join(__dirname, "../", "db", "v1.json"), JSON.stringify(defaultStructure));
        deferred.resolve();
    });
    return deferred.promise;
};

/*
 nginx rtmp server
 */
const startNginxRTMPServer = function startNginxRTMPServer(){
    logger.info("starting nginx server....", "start.nginx");
    var deferred = Q.defer();
    const command = config.nginx.exec;
    const  spawn = require('child_process').spawn;
    spawn('sh', ['-c', command], { stdio: 'inherit' });
    deferred.resolve();
    return deferred.promise;
};

/*
 webserver
 */
const startWebserver = function startWebserver(){
    logger.info("starting webserver...");
    var deferred = Q.defer();
    app.set('port', process.env.NODEJS_PORT);
    var server = app.listen(app.get('port'), function () {
        require("./classes/WebsocketController").bindDefaultEvents();
        app.set("io", require('socket.io')(server));
        app.set("server", server.address());
        app.get("websocketsReady").resolve(app.get("io")); //promise to determine if the webserver has been started to avoid ws binding before
        deferred.resolve(server.address().port);
    });
    return deferred.promise;
};

/*
public ip
 */
const getPublicIp = function(){
    logger.info("Getting public ip...", "start.publicip");
    var exec = require("child_process").exec;
    exec("public-ip", (err, stdout, stderr)=>{
        if (err){
            console.log(err);
        }
        app.set("publicIp", stdout.split("\n")[0] );
    });
};

/*
Restore ffmpeg processes after an applicatoin restart or stuff
 */
const restoreFFMPEGProcesses = function(){
    logger.info("Restoring ffmpeg processes...", "start.restore");
    var deferred = Q.defer();
    const ReStreamer = require("./classes/ReStreamer");
    ReStreamer.restoreFFMpegProcesses();
    deferred.resolve();
    return deferred.promise;
};

/*
Init
 */
startNginxRTMPServer()
    .then(function(){
        return checkJsonDb();
    })
    .then(function(){
        return startWebserver();
    })
    .then(function(port){
        logger.info("Webserver running on port " + process.env.NODEJS_PORT, "start.webserver");
    })
    .then(function(){
        getPublicIp();
        return restoreFFMPEGProcesses();
    })
    .catch(function(error){
        logger.error("error starting webserver and nginx for application: " + error);
        setTimeout(()=> {
            process.exit();
        }, 500);
    });
