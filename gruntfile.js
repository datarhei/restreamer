'use strict';

module.exports = function(grunt) {
    // files for this project
    var files = {
        compiledFrontendJS: ['bin/webserver/public/scripts/**/*.js', 'bin/executors/**/gui/js/*.js'],
        es6Src: ['**/*.js'],
        stylesheets: ['static/webserver/public/css/*.css']
    };

    // Project Configuration
    grunt.initConfig({

        /*
            Watcher config with livereload
        */
        watch: {
            scripts: {
                files: ['src/**/*.js'],
                tasks: ['compile-code'],
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
                tasks: ['copy-statics'],
                options: {
                    interrupt: true,
                    livereload: {
                        host: 'localhost',
                        port: 5000
                    }
                }
            }
        },

        /*
            Config for shell commands
        */
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
            copyTemplates: {
                command: 'cp -R src/webserver/views/. bin/webserver/views'
            },
            bower: {
                command: 'bower install --allow-root'
            },
            eslint: {
                command: 'eslint src/*'
            },
            //temp workaround - https://github.com/clappr/clappr/issues/709
            clappr: {
                command: 'curl -LOks https://github.com/clappr/clappr/archive/master.tar.gz && tar xzvf master.tar.gz && rm master.tar.gz && mv clappr-master bin/webserver/public/libs/clappr'
            }
        },

        /*
            Config for Babel compiling
        */
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

        /*
            Config for eslinter
        */
        eslint: {
            all: ['src/**/*.js'],
            options: {
                configFile: '.eslintrc.json'
            }
        },

        /*
            config for css linter
         */
        csslint: {
            options: {
                csslintrc: '.csslintrc'
                },
                all: {
                src: ['static/webserver/public/css/*.css']
            }
        },

        /*
            uglify and minify frontend javascript
         */
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

        /*
            minify css files
         */
        cssmin: {
            combine: {
                files: {
                    'bin/webserver/public/css/restreamer.min.css': '<%= stylesheets %>'
                }
            }
        },

        /*
            produces one file from all fontend javascript bewaring DI naming of angular
         */
        ngAnnotate: {
            production: {
                files: {
                    'bin/webserver/public/dist/application.js': '<%= compiledFrontendJS %>'
                }
            }
            }
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
    grunt.loadNpmTasks('grunt-contrib-watch');

    /*
     Helper tasks to keep overview
     */
    // lint
    grunt.registerTask('lint', ['csslint', 'shell:eslint']);
    // clear old bin folder and create new one
    grunt.registerTask('clearOldBuild', ['shell:removeOldBinFolder', 'shell:createBinFolder']);
    // install frontendlibraries (atm through bower)
    grunt.registerTask('installFrontendLibraries', ['shell:bower', 'shell:clappr']);
    // minify the frontend files
    grunt.registerTask('minifyFrontendFiles', ['cssmin', 'ngAnnotate', 'uglify']);

    /*
    Build Tasks
     */
    grunt.registerTask('build', ['loadConfig','clearOldBuild', 'shell:copyStatics', 'shell:copyTemplates', 'babel', 'minifyFrontendFiles', 'installFrontendLibraries']);
    grunt.registerTask('compile-code', ['loadConfig','babel', 'shell:copyStatics', 'shell:copyTemplates', 'minifyFrontendFiles']);
    grunt.registerTask('copy-statics', ['loadConfig', 'shell:copyStatics']);
    grunt.registerTask('copy-templates', ['loadConfig', 'shell:copyTemplates']);

    /*
    Run Tasks
     */
    // run current build in /bin
    grunt.registerTask('run', ['shell:start']);
    // rebuild and run
    grunt.registerTask('run-clean', ['build', 'run']);
    // update code and run
    grunt.registerTask('run-update-code', ['loadConfig','babel','shell:copyStatics', 'shell:copyTemplates', 'minifyFrontendFiles', 'run'])
};
