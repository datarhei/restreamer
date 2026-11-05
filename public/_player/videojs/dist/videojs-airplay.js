(function(){function r(e,n,t){function o(i,f){if(!n[i]){if(!e[i]){var c="function"==typeof require&&require;if(!f&&c)return c(i,!0);if(u)return u(i,!0);var a=new Error("Cannot find module '"+i+"'");throw a.code="MODULE_NOT_FOUND",a}var p=n[i]={exports:{}};e[i][0].call(p.exports,function(r){var n=e[i][1][r];return o(n||r)},p,p.exports,r,e,n,t)}return n[i].exports}for(var u="function"==typeof require&&require,i=0;i<t.length;i++)o(t[i]);return o}return r})()({1:[function(require,module,exports){
"use strict";

var hasAirPlayAPISupport = require('../lib/hasAirPlayAPISupport');

/**
 * Registers the AirPlayButton Component with Video.js. Calls
 * {@link http://docs.videojs.com/Component.html#.registerComponent}, which will add a
 * component called `airPlayButton` to the list of globally registered Video.js
 * components. The `airPlayButton` is added to the player's control bar UI automatically
 * once {@link module:enableAirPlay} has been called. If you would like to specify the
 * order of the buttons that appear in the control bar, including this button, you can do
 * so in the options that you pass to the `videojs` function when creating a player:
 *
 * ```
 * videojs('playerID', {
 *    controlBar: {
 *       children: [
 *          'playToggle',
 *          'progressControl',
 *          'volumePanel',
 *          'fullscreenToggle',
 *          'airPlayButton',
 *       ],
 *    }
 * });
 * ```
 *
 * @param videojs {object} A reference to {@link http://docs.videojs.com/module-videojs.html|Video.js}
 * @see http://docs.videojs.com/module-videojs.html#~registerPlugin
 */
module.exports = function (videojs) {
  /**
   * The AirPlayButton module contains both the AirPlayButton class definition and the
   * function used to register the button as a Video.js Component.
   *
   * @module AirPlayButton
   */

  const ButtonComponent = videojs.getComponent('Button');

  /**
  * The Video.js Button class is the base class for UI button components.
  *
  * @external Button
  * @see {@link http://docs.videojs.com/Button.html|Button}
  */

  /** @lends AirPlayButton.prototype */
  class AirPlayButton extends ButtonComponent {
    /**
     * This class is a button component designed to be displayed in the
     * player UI's control bar. It displays an Apple AirPlay selection
     * list when clicked.
     *
     * @constructs
     * @extends external:Button
     */
    constructor(player, options) {
      super(player, options);
      if (!hasAirPlayAPISupport()) {
        this.hide();
      }
      this._reactToAirPlayAvailableEvents();
      if (options.addAirPlayLabelToButton) {
        this.el().classList.add('vjs-airplay-button-lg');
        this._labelEl = document.createElement('span');
        this._labelEl.classList.add('vjs-airplay-button-label');
        this._labelEl.textContent = this.localize('AirPlay');
        this.el().appendChild(this._labelEl);
      } else {
        this.controlText('Start AirPlay');
      }
    }

    /**
     * Overrides Button#buildCSSClass to return the classes used on the button element.
     *
     * @param {DOMElement} el
     * @see {@link http://docs.videojs.com/Button.html#buildCSSClass|Button#buildCSSClass}
     */
    buildCSSClass() {
      return 'vjs-airplay-button ' + super.buildCSSClass();
    }

    /**
     * Overrides Button#handleClick to handle button click events. AirPlay
     * functionality is handled outside of this class, which should be limited
     * to UI related logic. This function simply triggers an event on the player.
     *
     * @fires AirPlayButton#airPlayRequested
     * @param {DOMElement} el
     * @see {@link http://docs.videojs.com/Button.html#handleClick|Button#handleClick}
     */
    handleClick() {
      this.player().trigger('airPlayRequested');
    }

    /**
     * Gets the underlying DOMElement used by the player.
     *
     * @private
     * @returns {DOMElement} either an <audio> or <video> tag, depending on the type of
     * player
     */
    _getMediaEl() {
      var playerEl = this.player().el();
      return playerEl.querySelector('video, audio');
    }

    /**
     * Binds a listener to the `webkitplaybacktargetavailabilitychanged` event, if it is
     * supported, that will show or hide this button Component based on the availability
     * of the AirPlay function.
     *
     * @private
     */
    _reactToAirPlayAvailableEvents() {
      var mediaEl = this._getMediaEl(),
        self = this;
      if (!mediaEl || !hasAirPlayAPISupport()) {
        return;
      }
      function onTargetAvailabilityChanged(event) {
        if (event.availability === 'available') {
          self.show();
        } else {
          self.hide();
        }
      }
      mediaEl.addEventListener('webkitplaybacktargetavailabilitychanged', onTargetAvailabilityChanged);
      this.on('dispose', function () {
        mediaEl.removeEventListener('webkitplaybacktargetavailabilitychanged', onTargetAvailabilityChanged);
      });
    }
  }
  videojs.registerComponent('airPlayButton', AirPlayButton);
};

},{"../lib/hasAirPlayAPISupport":4}],2:[function(require,module,exports){
"use strict";

/**
 * @module enableAirPlay
 */

var hasAirPlayAPISupport = require('./lib/hasAirPlayAPISupport');

/**
 * @private
 * @param {object} the Video.js Player instance
 * @returns {AirPlayButton} or `undefined` if it does not exist
 */
function getExistingAirPlayButton(player) {
  return player.controlBar.getChild('airPlayButton');
}

/**
 * Adds the AirPlayButton Component to the player's ControlBar component, if the
 * AirPlayButton does not already exist in the ControlBar.
 * @private
 * @param player {object} the Video.js Player instance
 * @param options {object}
 */
function ensureAirPlayButtonExists(player, options) {
  var existingAirPlayButton = getExistingAirPlayButton(player),
    indexOpt;
  if (options.addButtonToControlBar && !existingAirPlayButton) {
    // Figure out AirPlay button's index
    indexOpt = player.controlBar.children().length;
    if (typeof options.buttonPositionIndex !== 'undefined') {
      indexOpt = options.buttonPositionIndex >= 0 ? options.buttonPositionIndex : player.controlBar.children().length + options.buttonPositionIndex;
    }
    player.controlBar.addChild('airPlayButton', options, indexOpt);
  }
}

/**
 * Handles requests for AirPlay triggered by the AirPlayButton Component.
 *
 * @private
 * @param player {object} the Video.js Player instance
 */
function onAirPlayRequested(player) {
  var mediaEl = player.el().querySelector('video, audio');
  if (mediaEl && mediaEl.webkitShowPlaybackTargetPicker) {
    mediaEl.webkitShowPlaybackTargetPicker();
  }
}

/**
 * Adds an event listener for the `airPlayRequested` event triggered by the AirPlayButton
 * Component.
 *
 * @private
 * @param player {object} the Video.js Player instance
 */
function listenForAirPlayEvents(player) {
  // Respond to requests for AirPlay. The AirPlayButton component triggers this event
  // when the user clicks the AirPlay button.
  player.on('airPlayRequested', onAirPlayRequested.bind(null, player));
}

/**
 * Sets up the AirPlay plugin.
 *
 * @private
 * @param player {object} the Video.js player
 * @param options {object} the plugin options
 */
function enableAirPlay(player, options) {
  if (!player.controlBar) {
    return;
  }
  if (hasAirPlayAPISupport()) {
    listenForAirPlayEvents(player);
    ensureAirPlayButtonExists(player, options);
  }
}

/**
 * Registers the AirPlay plugin with Video.js. Calls
 * {@link http://docs.videojs.com/module-videojs.html#~registerPlugin|videojs#registerPlugin},
 * which will add a plugin function called `airPlay` to any instance of a Video.js player
 * that is created after calling this function. Call `player.airPlay(options)`, passing in
 * configuration options, to enable the AirPlay plugin on your Player instance.
 *
 * Currently, the only configuration option is:
 *
 *    * **buttonText** - the text to display inside of the button component. By default,
 *    this text is hidden and is used for accessibility purposes.
 *
 * @param {object} videojs
 * @see http://docs.videojs.com/module-videojs.html#~registerPlugin
 */
module.exports = function (videojs) {
  videojs.registerPlugin('airPlay', function (options) {
    var pluginOptions = Object.assign({
      addButtonToControlBar: true
    }, options || {});

    // `this` is an instance of a Video.js Player.
    // Wait until the player is "ready" so that the player's control bar component has
    // been created.
    this.ready(enableAirPlay.bind(this, this, pluginOptions));
  });
};

},{"./lib/hasAirPlayAPISupport":4}],3:[function(require,module,exports){
"use strict";

var createAirPlayButton = require('./components/AirPlayButton'),
  createAirPlayPlugin = require('./enableAirPlay');

/**
 * @module index
 */

/**
 * Registers the AirPlay plugin and AirPlayButton Component with Video.js. See
 * {@link module:AirPlayButton} and {@link module:enableAirPlay} for more details about
 * how the plugin and button are registered and configured.
 *
 * @param {object} videojs
 * @see module:enableAirPlay
 * @see module:AirPlayButton
 */
module.exports = function (videojs) {
  videojs = videojs || window.videojs;
  createAirPlayButton(videojs);
  createAirPlayPlugin(videojs);
};

},{"./components/AirPlayButton":1,"./enableAirPlay":2}],4:[function(require,module,exports){
"use strict";

/**
 * @module hasAirPlayAPISupport
 */

/**
 * Returns whether or not the current browser environment supports AirPlay.
 *
 * @private
 * @returns {boolean} true if AirPlay support is available
 */
module.exports = function () {
  return !!window.WebKitPlaybackTargetAvailabilityEvent;
};

},{}],5:[function(require,module,exports){
"use strict";

/**
 * This module is used as an entry point for the build system to bundle this plugin into a
 * single javascript file that can be loaded by a script tag on a web page. The javascript
 * file that is built assumes that `videojs` is available globally at `window.videojs`, so
 * Video.js must be loaded **before** this plugin is loaded.
 *
 * Run `npm install` and then `grunt build` to build the plugin's bundled javascript
 * file, as well as the CSS and image assets into the project's `./dist/` folder.
 *
 * @module standalone
 */

require('./index')();

},{"./index":3}]},{},[5]);
