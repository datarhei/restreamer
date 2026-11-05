/*! @name videojs-license @version 0.1.0 @license MIT */
(function (global, factory) {
  typeof exports === 'object' && typeof module !== 'undefined' ? module.exports = factory(require('video.js')) :
  typeof define === 'function' && define.amd ? define(['video.js'], factory) :
  (global = typeof globalThis !== 'undefined' ? globalThis : global || self, global.videojsLicense = factory(global.videojs));
})(this, (function (videojs) { 'use strict';

  function _interopDefaultLegacy (e) { return e && typeof e === 'object' && 'default' in e ? e : { 'default': e }; }

  var videojs__default = /*#__PURE__*/_interopDefaultLegacy(videojs);

  var version = "0.1.0";

  const Plugin = videojs__default["default"].getPlugin('plugin');
  const Component = videojs__default["default"].getComponent('Component');
  const Button = videojs__default["default"].getComponent('MenuButton');

  // Default options for the plugin.
  const defaults = {
    license: 'none',
    title: '',
    author: '',
    languages: {
      license: 'License',
      loading: 'Loading'
    }
  };

  /**
   * An advanced Video.js plugin. For more information on the API
   *
   * See: https://blog.videojs.com/feature-spotlight-advanced-plugins/
   */
  class License extends Plugin {
    /**
     * Create a License plugin instance.
     *
     * @param  {Player} player
     *         A Video.js Player instance.
     *
     * @param  {Object} [options]
     *         An optional options object.
     *
     *         While not a core part of the Video.js plugin architecture, a
     *         second argument of options is a convenient way to accept inputs
     *         from your plugin's caller.
     */
    constructor(player, options) {
      // the parent class will add player under this.player
      super(player);
      this.playerId = this.player.id();
      this.options = videojs__default["default"].mergeOptions(defaults, options);
      if (options.license === 'none') {
        return;
      }
      player.addClass('vjs-license');
      this.buildUI();
      if (videojs__default["default"].browser.IS_IOS || videojs__default["default"].browser.IS_ANDROID) {
        this.mobileBuildUI();
      }

      // close the menu if open on userinactive
      this.player.on('userinactive', () => {
        document.getElementById(this.playerId).querySelectorAll('.vjs-menu').forEach(element => {
          element.classList.remove('vjs-lock-open');
        });
      });

      // close the menu if anywhere in the player is clicked
      this.player.on('click', evt => {
        if (evt.target.tagName === 'VIDEO') {
          document.getElementById(this.playerId).querySelectorAll('.vjs-menu').forEach(element => {
            element.classList.remove('vjs-lock-open');
          });
        }
      });
      this.player.on('loadstart', _event => {
        this.removeElementsByClass('vjs-license-clear');
        if (videojs__default["default"].browser.IS_IOS || videojs__default["default"].browser.IS_ANDROID) {
          this.mobileBuildTopLevelMenu();
        } else {
          this.buildTopLevelMenu();
        }
      });
    }

    /**
     * Add the menu ui button to the controlbar
     */
    buildUI() {
      const playerId = this.playerId;
      const that = this;

      /**
       * LicenseMenuButton
       */
      class LicenseMenuButton extends Button {
        /**
         * Contructor
         *
         * @param {*} player videojs player instance
         * @param {*} options videojs player options
         */
        constructor(player, options) {
          super(player, options);
          this.addClass('vjs-license');
          this.controlText(that.options.languages.loading);
          player.one('canplaythrough', _event => {
            this.controlText(that.options.languages.settings);
          });
          this.menu.contentEl_.id = playerId + '-vjs-license-default';
        }

        /**
         * Handle click
         */
        handleClick() {
          if (videojs__default["default"].browser.IS_IOS || videojs__default["default"].browser.IS_ANDROID) {
            this.player.getChild('licenseMenuMobileModal').el().style.display = 'block';
          } else {
            this.el().classList.toggle('vjs-toogle-btn');
            this.menu.el().classList.toggle('vjs-lock-open');
          }
        }
      }
      videojs__default["default"].registerComponent('licenseMenuButton', LicenseMenuButton);
      this.player.getChild('controlBar').addChild('licenseMenuButton');
      if (this.player.getChild('controlBar').getChild('fullscreenToggle')) {
        this.player.getChild('controlBar').el().insertBefore(this.player.getChild('controlBar').getChild('licenseMenuButton').el(), this.player.getChild('controlBar').getChild('fullscreenToggle').el());
      }
    }

    /**
     *
     * This is just build the top level menu no sub menus
     */
    buildTopLevelMenu() {
      const settingsButton = this.player.getChild('controlBar').getChild('licenseMenuButton');

      // settingsButton.addClass('vjs-license-is-loaded');

      const main = settingsButton.menu.contentEl_;

      // Empty the main menu div to repopulate
      main.innerHTML = '';
      main.classList.add('vjs-license-top-level');

      // Start building new list items
      const menuTitle = document.createElement('li');
      menuTitle.className = 'vjs-license-top-level-header';
      const menuTitleInner = document.createElement('span');
      menuTitleInner.innerHTML = 'About';
      menuTitleInner.className = 'vjs-license-top-level-header-titel';
      menuTitle.appendChild(menuTitleInner);
      main.appendChild(menuTitle);
      const itemTitel = document.createElement('li');
      itemTitel.innerHTML = this.buildItemTitel();
      itemTitel.className = 'vjs-license-top-level-item';
      main.appendChild(itemTitel);
      if (this.options.author) {
        const itemAuthor = document.createElement('li');
        itemAuthor.innerHTML = this.buildItemAuthor();
        itemAuthor.className = 'vjs-license-top-level-item';
        main.appendChild(itemAuthor);
      }
      const itemLicense = document.createElement('li');
      itemLicense.innerHTML = this.buildItemLicense();
      itemLicense.className = 'vjs-license-top-level-item';
      main.appendChild(itemLicense);
    }

    /**
     * Add the menu ui button to the controlbar
     */
    mobileBuildUI() {
      /**
       * bla
       */
      class LicenseMenuMobileModal extends Component {
        /**
         * Contructor
         *
         * @param {*} player videojs player instance
         * @param {*} options videojs player options
         */
        constructor(player, options) {
          super(player, options);
        }

        /**
         * Creates an HTML element
         *
         * @return {Object} HTML element
         */
        createEl() {
          return videojs__default["default"].createEl('div', {
            className: 'vjs-license-mobile'
          });
        }
      }
      videojs__default["default"].registerComponent('licenseMenuMobileModal', LicenseMenuMobileModal);
      videojs__default["default"].dom.prependTo(this.player.addChild('licenseMenuMobileModal').el(), document.body);
    }

    /**
     * Add the menu ui button to the controlbar
     */
    mobileBuildTopLevelMenu() {
      const settingsButton = this.player.getChild('licenseMenuMobileModal');
      const menuTopLevel = document.createElement('ul');
      menuTopLevel.className = 'vjs-license-mob-top-level vjs-setting-menu-clear';
      settingsButton.el().appendChild(menuTopLevel);

      // Empty the main menu div to repopulate
      const menuTitle = document.createElement('li');
      menuTitle.innerHTML = 'About';
      menuTitle.className = 'vjs-setting-menu-mobile-top-header';
      menuTopLevel.appendChild(menuTitle);
      const itemTitel = document.createElement('li');
      itemTitel.innerHTML = this.buildItemTitel();
      itemTitel.className = 'vjs-license-top-level-item';
      if (this.options.author) {
        const itemAuthor = document.createElement('li');
        itemAuthor.innerHTML = this.buildItemAuthor();
        itemAuthor.className = 'vjs-license-top-level-item';
      }
      const itemLicense = document.createElement('li');
      itemLicense.innerHTML = this.buildItemLicense();
      itemLicense.className = 'vjs-license-top-level-item';
      const menuClose = document.createElement('li');
      menuClose.innerHTML = 'Close';
      menuClose.className = 'setting-menu-footer-default';
      menuClose.onclick = e => {
        this.player.getChild('settingsMenuMobileModal').el().style.display = 'none';
      };
      menuTopLevel.appendChild(menuClose);
    }

    /**
     * Add the menu ui button to the controlbar
     *
     * @return {string} Returns license text
     */
    buildItemTitel() {
      let titel = '';
      if (this.options.title) {
        titel = `${this.options.title}`;
      }
      return 'Title: ' + titel;
    }

    /**
     * Add the menu ui button to the controlbar
     *
     * @return {string} Returns license text
     */
    buildItemAuthor() {
      let author = '';
      if (this.options.author) {
        author = ` by ${this.options.author}`;
      }
      return 'Author: ' + author;
    }

    /**
     * Add the menu ui button to the controlbar
     *
     * @return {string} Returns license text
     */
    buildItemLicense() {
      let license = '';
      const reVersion = new RegExp('[0-9]+.[0-9]+$');
      let version = '4.0';
      const matches = this.options.license.match(reVersion);
      if (matches !== null) {
        version = matches[0];
      }
      const which = this.options.license.replace(reVersion, '').trim();
      let deed = null;
      switch (which) {
        case 'CC0':
          deed = 'https://creativecommons.org/licenses/zero/1.0/';
          break;
        case 'CC BY':
          deed = `https://creativecommons.org/licenses/by/${version}/`;
          break;
        case 'CC BY-SA':
          deed = `https://creativecommons.org/licenses/by-sa/${version}/`;
          break;
        case 'CC BY-NC':
          deed = `https://creativecommons.org/licenses/by-nc/${version}/`;
          break;
        case 'CC BY-NC-SA':
          deed = `https://creativecommons.org/licenses/by-nc-sa/${version}/`;
          break;
        case 'CC BY-ND':
          deed = `https://creativecommons.org/licenses/by-nd/${version}/`;
          break;
        case 'CC BY-NC-ND':
          deed = `https://creativecommons.org/licenses/by-nc-nd/${version}/`;
          break;
      }
      if (deed) {
        license = `<a href='${deed}' onclick='window.open('${deed}')' target='_blank' rel='noopener'>${this.options.license}</a>`;
      } else {
        license = this.options.license;
      }
      return 'License: ' + license;
    }

    /**
     *
     * Helper class to clear menu items before rebuild
     *
     * @param {*} className Name of a class
     */
    removeElementsByClass(className) {
      // Need to prevent the menu from not showing sometimes
      document.querySelectorAll('.vjs-sm-top-level').forEach(element => {
        element.classList.remove('vjs-hidden');
      });
      const elements = document.getElementsByClassName(className);
      while (elements.length > 0) {
        elements[0].parentNode.removeChild(elements[0]);
      }
    }
  }

  // Define default values for the plugin's `state` object here.
  License.defaultState = {};

  // Include the version number.
  License.VERSION = version;

  // Register the plugin with video.js.
  videojs__default["default"].registerPlugin('license', License);

  return License;

}));
