'use strict';

module.exports = function(grunt) {
    // files for this project
    var files = {
        compiledFrontendJS: ['bin/webserver/public/scripts/**/*.js', 'bin/executors/**/gui/js/*.js'],
        es6Src: ['**/*.js'],
        stylesheets: ['static/webserver/public/css/*.css']
    };
    // find all node modules
    var modules = [];
    var pck = require('./package.json')
    if (!!pck.dependencies) {
        modules = Object.keys(pck.dependencies)
            .filter(function (m) {
                return m != 'nodewebkit'
            })
            .map(function (m) {
                return './node_modules/' + m + '/**/*'
            });
    }
    // Project Configuration
    grunt.initConfig({
        watch: {
            scripts: {
                files: ['src/**/*.js'],
                tasks: ['loadConfig','babel', 'minifyFrontendFiles'],
                options: {
                    interrupt: true,
                    livereload: {
                        host: 'localhost',
                        port: 5000
                    }
                }
            },
            statics: {
                files: ['static/**/*.html'],
                tasks: ['loadConfig','shell:copyStatics'],
                options: {
                    interrupt: true,
                    livereload: {
                        host: 'localhost',
                        port: 5000
                    }
                }
            }
        },
        shell: {
            start: {
                command: 'npm start'
            },
            removeOldBinFolder: {
                command: 'rm -Rf bin/'
            },
            createBinFolder: {
                command: 'mkdir bin/'
            },
            copyStatics: {
                command: 'cp -R static/* bin'
            },
            bower: {
                command: 'bower install --allow-root'
            },
            //temp workaround - https://github.com/clappr/clappr/issues/709
            clappr: {
                command: 'git clone https://github.com/clappr/clappr bin/webserver/public/libs/clappr'
            }
        },
        babel: {
            options: {
                sourceMap: true,
                presets: ['es2015']
            },
            all: {
                files: [
                    {
                        expand: true,
                        cwd: 'src/',
                        src:'<%= es6Src %>',
                        dest: 'bin/'
                    }
                ]
            }
        },
        pkg: grunt.file.readJSON('package.json'),
        eslint: {
            all: ['src/**/*.js'],
            options: {
                configFile: '.eslintrc.json'
            }
        },
        csslint: {
            options: {
                csslintrc: '.csslintrc'
                },
                all: {
                src: ['static/webserver/public/css/*.css']
            }
        },
        uglify: {
            production: {
                options: {
                    mangle: true
                },
                files: {
                    'bin/webserver/public/dist/application.min.js': 'bin/webserver/public/dist/application.js'
                }
            }
        },
        cssmin: {
            combine: {
                files: {
                    'bin/webserver/public/css/restreamer.min.css': '<%= stylesheets %>'
                }
            }
        },
        ngAnnotate: {
            production: {
                files: {
                    'bin/webserver/public/dist/application.js': '<%= compiledFrontendJS %>'
                }
            }
            },
        });

    /*
     Load NPM tasks
     */
    require('load-grunt-tasks')(grunt);
    grunt.task.registerTask('loadConfig', 'Task that loads the config into a grunt option.', function() {
        grunt.config.set('es6Src', files.es6Src);
        grunt.config.set('compiledFrontendJS', files.compiledFrontendJS);
        grunt.config.set('stylesheets', files.stylesheets);
    });
    grunt.loadNpmTasks('grunt-shell');
    grunt.loadNpmTasks('gruntify-eslint');
    grunt.loadNpmTasks('grunt-contrib-watch');

    /*
     Helper tasks to keep overview
     */
    // lint
    grunt.registerTask('lint', ['csslint', 'eslint']);
    // clear old bin folder and create new one
    grunt.registerTask('clearOldBuild', ['shell:removeOldBinFolder', 'shell:createBinFolder']);
    // install frontendlibraries (atm through bower)
    grunt.registerTask('installFrontendLibraries', ['shell:bower', 'shell:clappr']);
    // minify the frontend files
    grunt.registerTask('minifyFrontendFiles', ['cssmin', 'ngAnnotate', 'uglify']);

    /*
    Build Tasks
     */
    grunt.registerTask('build', ['loadConfig','clearOldBuild', 'shell:copyStatics', 'babel', 'minifyFrontendFiles', 'installFrontendLibraries']);
    grunt.registerTask('build-code', ['loadConfig', 'shell:copyStatics', 'babel', 'minifyFrontendFiles']);

    /*
    Run Tasks
     */
    // run current build in /bin
    grunt.registerTask('run', ['shell:start']);
    // rebuild and run
    grunt.registerTask('run-clean', ['build', 'run']);
    // update code and run
    grunt.registerTask('run-update-code', ['loadConfig','babel','shell:copyStatics', 'minifyFrontendFiles', 'run'])
};
