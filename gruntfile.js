'use strict';

module.exports = function(grunt) {

    // Project Configuration
    grunt.initConfig({

        /*
         Config for shell commands
         */
        shell: {
            bower: {
                command: 'bower install --allow-root'
            },
            npminstall: {
                command: 'npm install'
            },
            eslint: {
                command: 'eslint src/*'
            },
            //temp workaround - https://github.com/clappr/clappr/issues/709
            clappr: {
                command: 'curl -LOks https://github.com/clappr/clappr/archive/master.tar.gz && tar xzvf master.tar.gz && rm master.tar.gz && mv clappr-master src/webserver/public/libs/clappr'
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
        }
    })


    /*
     Load NPM tasks
     */
    require('load-grunt-tasks')(grunt);


    grunt.loadNpmTasks('grunt-shell');

    // lint
    grunt.registerTask('lint', ['shell:eslint']);
    // install frontendlibraries (atm through bower)
    grunt.registerTask('installFrontendLibraries', ['shell:bower', 'shell:clappr']);

    grunt.registerTask('build', ['shell:npminstall', 'installFrontendLibraries']);
};
