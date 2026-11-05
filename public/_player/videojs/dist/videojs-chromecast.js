(function(){function r(e,n,t){function o(i,f){if(!n[i]){if(!e[i]){var c="function"==typeof require&&require;if(!f&&c)return c(i,!0);if(u)return u(i,!0);var a=new Error("Cannot find module '"+i+"'");throw a.code="MODULE_NOT_FOUND",a}var p=n[i]={exports:{}};e[i][0].call(p.exports,function(r){var n=e[i][1][r];return o(n||r)},p,p.exports,r,e,n,t)}return n[i].exports}for(var u="function"==typeof require&&require,i=0;i<t.length;i++)o(t[i]);return o}return r})()({1:[function(require,module,exports){
/**
 * @license
 * Copyright (c) 2014 The Polymer Project Authors. All rights reserved.
 * This code may only be used under the BSD style license found at http://polymer.github.io/LICENSE.txt
 * The complete set of authors may be found at http://polymer.github.io/AUTHORS.txt
 * The complete set of contributors may be found at http://polymer.github.io/CONTRIBUTORS.txt
 * Code distributed by Google as part of the polymer project is also
 * subject to an additional IP rights grant found at http://polymer.github.io/PATENTS.txt
 */
// @version 0.7.24
(function() {
  window.WebComponents = window.WebComponents || {
    flags: {}
  };
  var file = "webcomponents-lite.js";
  var script = document.querySelector('script[src*="' + file + '"]');
  var flags = {};
  if (!flags.noOpts) {
    location.search.slice(1).split("&").forEach(function(option) {
      var parts = option.split("=");
      var match;
      if (parts[0] && (match = parts[0].match(/wc-(.+)/))) {
        flags[match[1]] = parts[1] || true;
      }
    });
    if (script) {
      for (var i = 0, a; a = script.attributes[i]; i++) {
        if (a.name !== "src") {
          flags[a.name] = a.value || true;
        }
      }
    }
    if (flags.log && flags.log.split) {
      var parts = flags.log.split(",");
      flags.log = {};
      parts.forEach(function(f) {
        flags.log[f] = true;
      });
    } else {
      flags.log = {};
    }
  }
  if (flags.register) {
    window.CustomElements = window.CustomElements || {
      flags: {}
    };
    window.CustomElements.flags.register = flags.register;
  }
  WebComponents.flags = flags;
})();

(function(scope) {
  "use strict";
  var hasWorkingUrl = false;
  if (!scope.forceJURL) {
    try {
      var u = new URL("b", "http://a");
      u.pathname = "c%20d";
      hasWorkingUrl = u.href === "http://a/c%20d";
    } catch (e) {}
  }
  if (hasWorkingUrl) return;
  var relative = Object.create(null);
  relative["ftp"] = 21;
  relative["file"] = 0;
  relative["gopher"] = 70;
  relative["http"] = 80;
  relative["https"] = 443;
  relative["ws"] = 80;
  relative["wss"] = 443;
  var relativePathDotMapping = Object.create(null);
  relativePathDotMapping["%2e"] = ".";
  relativePathDotMapping[".%2e"] = "..";
  relativePathDotMapping["%2e."] = "..";
  relativePathDotMapping["%2e%2e"] = "..";
  function isRelativeScheme(scheme) {
    return relative[scheme] !== undefined;
  }
  function invalid() {
    clear.call(this);
    this._isInvalid = true;
  }
  function IDNAToASCII(h) {
    if ("" == h) {
      invalid.call(this);
    }
    return h.toLowerCase();
  }
  function percentEscape(c) {
    var unicode = c.charCodeAt(0);
    if (unicode > 32 && unicode < 127 && [ 34, 35, 60, 62, 63, 96 ].indexOf(unicode) == -1) {
      return c;
    }
    return encodeURIComponent(c);
  }
  function percentEscapeQuery(c) {
    var unicode = c.charCodeAt(0);
    if (unicode > 32 && unicode < 127 && [ 34, 35, 60, 62, 96 ].indexOf(unicode) == -1) {
      return c;
    }
    return encodeURIComponent(c);
  }
  var EOF = undefined, ALPHA = /[a-zA-Z]/, ALPHANUMERIC = /[a-zA-Z0-9\+\-\.]/;
  function parse(input, stateOverride, base) {
    function err(message) {
      errors.push(message);
    }
    var state = stateOverride || "scheme start", cursor = 0, buffer = "", seenAt = false, seenBracket = false, errors = [];
    loop: while ((input[cursor - 1] != EOF || cursor == 0) && !this._isInvalid) {
      var c = input[cursor];
      switch (state) {
       case "scheme start":
        if (c && ALPHA.test(c)) {
          buffer += c.toLowerCase();
          state = "scheme";
        } else if (!stateOverride) {
          buffer = "";
          state = "no scheme";
          continue;
        } else {
          err("Invalid scheme.");
          break loop;
        }
        break;

       case "scheme":
        if (c && ALPHANUMERIC.test(c)) {
          buffer += c.toLowerCase();
        } else if (":" == c) {
          this._scheme = buffer;
          buffer = "";
          if (stateOverride) {
            break loop;
          }
          if (isRelativeScheme(this._scheme)) {
            this._isRelative = true;
          }
          if ("file" == this._scheme) {
            state = "relative";
          } else if (this._isRelative && base && base._scheme == this._scheme) {
            state = "relative or authority";
          } else if (this._isRelative) {
            state = "authority first slash";
          } else {
            state = "scheme data";
          }
        } else if (!stateOverride) {
          buffer = "";
          cursor = 0;
          state = "no scheme";
          continue;
        } else if (EOF == c) {
          break loop;
        } else {
          err("Code point not allowed in scheme: " + c);
          break loop;
        }
        break;

       case "scheme data":
        if ("?" == c) {
          this._query = "?";
          state = "query";
        } else if ("#" == c) {
          this._fragment = "#";
          state = "fragment";
        } else {
          if (EOF != c && "\t" != c && "\n" != c && "\r" != c) {
            this._schemeData += percentEscape(c);
          }
        }
        break;

       case "no scheme":
        if (!base || !isRelativeScheme(base._scheme)) {
          err("Missing scheme.");
          invalid.call(this);
        } else {
          state = "relative";
          continue;
        }
        break;

       case "relative or authority":
        if ("/" == c && "/" == input[cursor + 1]) {
          state = "authority ignore slashes";
        } else {
          err("Expected /, got: " + c);
          state = "relative";
          continue;
        }
        break;

       case "relative":
        this._isRelative = true;
        if ("file" != this._scheme) this._scheme = base._scheme;
        if (EOF == c) {
          this._host = base._host;
          this._port = base._port;
          this._path = base._path.slice();
          this._query = base._query;
          this._username = base._username;
          this._password = base._password;
          break loop;
        } else if ("/" == c || "\\" == c) {
          if ("\\" == c) err("\\ is an invalid code point.");
          state = "relative slash";
        } else if ("?" == c) {
          this._host = base._host;
          this._port = base._port;
          this._path = base._path.slice();
          this._query = "?";
          this._username = base._username;
          this._password = base._password;
          state = "query";
        } else if ("#" == c) {
          this._host = base._host;
          this._port = base._port;
          this._path = base._path.slice();
          this._query = base._query;
          this._fragment = "#";
          this._username = base._username;
          this._password = base._password;
          state = "fragment";
        } else {
          var nextC = input[cursor + 1];
          var nextNextC = input[cursor + 2];
          if ("file" != this._scheme || !ALPHA.test(c) || nextC != ":" && nextC != "|" || EOF != nextNextC && "/" != nextNextC && "\\" != nextNextC && "?" != nextNextC && "#" != nextNextC) {
            this._host = base._host;
            this._port = base._port;
            this._username = base._username;
            this._password = base._password;
            this._path = base._path.slice();
            this._path.pop();
          }
          state = "relative path";
          continue;
        }
        break;

       case "relative slash":
        if ("/" == c || "\\" == c) {
          if ("\\" == c) {
            err("\\ is an invalid code point.");
          }
          if ("file" == this._scheme) {
            state = "file host";
          } else {
            state = "authority ignore slashes";
          }
        } else {
          if ("file" != this._scheme) {
            this._host = base._host;
            this._port = base._port;
            this._username = base._username;
            this._password = base._password;
          }
          state = "relative path";
          continue;
        }
        break;

       case "authority first slash":
        if ("/" == c) {
          state = "authority second slash";
        } else {
          err("Expected '/', got: " + c);
          state = "authority ignore slashes";
          continue;
        }
        break;

       case "authority second slash":
        state = "authority ignore slashes";
        if ("/" != c) {
          err("Expected '/', got: " + c);
          continue;
        }
        break;

       case "authority ignore slashes":
        if ("/" != c && "\\" != c) {
          state = "authority";
          continue;
        } else {
          err("Expected authority, got: " + c);
        }
        break;

       case "authority":
        if ("@" == c) {
          if (seenAt) {
            err("@ already seen.");
            buffer += "%40";
          }
          seenAt = true;
          for (var i = 0; i < buffer.length; i++) {
            var cp = buffer[i];
            if ("\t" == cp || "\n" == cp || "\r" == cp) {
              err("Invalid whitespace in authority.");
              continue;
            }
            if (":" == cp && null === this._password) {
              this._password = "";
              continue;
            }
            var tempC = percentEscape(cp);
            null !== this._password ? this._password += tempC : this._username += tempC;
          }
          buffer = "";
        } else if (EOF == c || "/" == c || "\\" == c || "?" == c || "#" == c) {
          cursor -= buffer.length;
          buffer = "";
          state = "host";
          continue;
        } else {
          buffer += c;
        }
        break;

       case "file host":
        if (EOF == c || "/" == c || "\\" == c || "?" == c || "#" == c) {
          if (buffer.length == 2 && ALPHA.test(buffer[0]) && (buffer[1] == ":" || buffer[1] == "|")) {
            state = "relative path";
          } else if (buffer.length == 0) {
            state = "relative path start";
          } else {
            this._host = IDNAToASCII.call(this, buffer);
            buffer = "";
            state = "relative path start";
          }
          continue;
        } else if ("\t" == c || "\n" == c || "\r" == c) {
          err("Invalid whitespace in file host.");
        } else {
          buffer += c;
        }
        break;

       case "host":
       case "hostname":
        if (":" == c && !seenBracket) {
          this._host = IDNAToASCII.call(this, buffer);
          buffer = "";
          state = "port";
          if ("hostname" == stateOverride) {
            break loop;
          }
        } else if (EOF == c || "/" == c || "\\" == c || "?" == c || "#" == c) {
          this._host = IDNAToASCII.call(this, buffer);
          buffer = "";
          state = "relative path start";
          if (stateOverride) {
            break loop;
          }
          continue;
        } else if ("\t" != c && "\n" != c && "\r" != c) {
          if ("[" == c) {
            seenBracket = true;
          } else if ("]" == c) {
            seenBracket = false;
          }
          buffer += c;
        } else {
          err("Invalid code point in host/hostname: " + c);
        }
        break;

       case "port":
        if (/[0-9]/.test(c)) {
          buffer += c;
        } else if (EOF == c || "/" == c || "\\" == c || "?" == c || "#" == c || stateOverride) {
          if ("" != buffer) {
            var temp = parseInt(buffer, 10);
            if (temp != relative[this._scheme]) {
              this._port = temp + "";
            }
            buffer = "";
          }
          if (stateOverride) {
            break loop;
          }
          state = "relative path start";
          continue;
        } else if ("\t" == c || "\n" == c || "\r" == c) {
          err("Invalid code point in port: " + c);
        } else {
          invalid.call(this);
        }
        break;

       case "relative path start":
        if ("\\" == c) err("'\\' not allowed in path.");
        state = "relative path";
        if ("/" != c && "\\" != c) {
          continue;
        }
        break;

       case "relative path":
        if (EOF == c || "/" == c || "\\" == c || !stateOverride && ("?" == c || "#" == c)) {
          if ("\\" == c) {
            err("\\ not allowed in relative path.");
          }
          var tmp;
          if (tmp = relativePathDotMapping[buffer.toLowerCase()]) {
            buffer = tmp;
          }
          if (".." == buffer) {
            this._path.pop();
            if ("/" != c && "\\" != c) {
              this._path.push("");
            }
          } else if ("." == buffer && "/" != c && "\\" != c) {
            this._path.push("");
          } else if ("." != buffer) {
            if ("file" == this._scheme && this._path.length == 0 && buffer.length == 2 && ALPHA.test(buffer[0]) && buffer[1] == "|") {
              buffer = buffer[0] + ":";
            }
            this._path.push(buffer);
          }
          buffer = "";
          if ("?" == c) {
            this._query = "?";
            state = "query";
          } else if ("#" == c) {
            this._fragment = "#";
            state = "fragment";
          }
        } else if ("\t" != c && "\n" != c && "\r" != c) {
          buffer += percentEscape(c);
        }
        break;

       case "query":
        if (!stateOverride && "#" == c) {
          this._fragment = "#";
          state = "fragment";
        } else if (EOF != c && "\t" != c && "\n" != c && "\r" != c) {
          this._query += percentEscapeQuery(c);
        }
        break;

       case "fragment":
        if (EOF != c && "\t" != c && "\n" != c && "\r" != c) {
          this._fragment += c;
        }
        break;
      }
      cursor++;
    }
  }
  function clear() {
    this._scheme = "";
    this._schemeData = "";
    this._username = "";
    this._password = null;
    this._host = "";
    this._port = "";
    this._path = [];
    this._query = "";
    this._fragment = "";
    this._isInvalid = false;
    this._isRelative = false;
  }
  function jURL(url, base) {
    if (base !== undefined && !(base instanceof jURL)) base = new jURL(String(base));
    this._url = url;
    clear.call(this);
    var input = url.replace(/^[ \t\r\n\f]+|[ \t\r\n\f]+$/g, "");
    parse.call(this, input, null, base);
  }
  jURL.prototype = {
    toString: function() {
      return this.href;
    },
    get href() {
      if (this._isInvalid) return this._url;
      var authority = "";
      if ("" != this._username || null != this._password) {
        authority = this._username + (null != this._password ? ":" + this._password : "") + "@";
      }
      return this.protocol + (this._isRelative ? "//" + authority + this.host : "") + this.pathname + this._query + this._fragment;
    },
    set href(href) {
      clear.call(this);
      parse.call(this, href);
    },
    get protocol() {
      return this._scheme + ":";
    },
    set protocol(protocol) {
      if (this._isInvalid) return;
      parse.call(this, protocol + ":", "scheme start");
    },
    get host() {
      return this._isInvalid ? "" : this._port ? this._host + ":" + this._port : this._host;
    },
    set host(host) {
      if (this._isInvalid || !this._isRelative) return;
      parse.call(this, host, "host");
    },
    get hostname() {
      return this._host;
    },
    set hostname(hostname) {
      if (this._isInvalid || !this._isRelative) return;
      parse.call(this, hostname, "hostname");
    },
    get port() {
      return this._port;
    },
    set port(port) {
      if (this._isInvalid || !this._isRelative) return;
      parse.call(this, port, "port");
    },
    get pathname() {
      return this._isInvalid ? "" : this._isRelative ? "/" + this._path.join("/") : this._schemeData;
    },
    set pathname(pathname) {
      if (this._isInvalid || !this._isRelative) return;
      this._path = [];
      parse.call(this, pathname, "relative path start");
    },
    get search() {
      return this._isInvalid || !this._query || "?" == this._query ? "" : this._query;
    },
    set search(search) {
      if (this._isInvalid || !this._isRelative) return;
      this._query = "?";
      if ("?" == search[0]) search = search.slice(1);
      parse.call(this, search, "query");
    },
    get hash() {
      return this._isInvalid || !this._fragment || "#" == this._fragment ? "" : this._fragment;
    },
    set hash(hash) {
      if (this._isInvalid) return;
      this._fragment = "#";
      if ("#" == hash[0]) hash = hash.slice(1);
      parse.call(this, hash, "fragment");
    },
    get origin() {
      var host;
      if (this._isInvalid || !this._scheme) {
        return "";
      }
      switch (this._scheme) {
       case "data":
       case "file":
       case "javascript":
       case "mailto":
        return "null";
      }
      host = this.host;
      if (!host) {
        return "";
      }
      return this._scheme + "://" + host;
    }
  };
  var OriginalURL = scope.URL;
  if (OriginalURL) {
    jURL.createObjectURL = function(blob) {
      return OriginalURL.createObjectURL.apply(OriginalURL, arguments);
    };
    jURL.revokeObjectURL = function(url) {
      OriginalURL.revokeObjectURL(url);
    };
  }
  scope.URL = jURL;
})(self);

if (typeof WeakMap === "undefined") {
  (function() {
    var defineProperty = Object.defineProperty;
    var counter = Date.now() % 1e9;
    var WeakMap = function() {
      this.name = "__st" + (Math.random() * 1e9 >>> 0) + (counter++ + "__");
    };
    WeakMap.prototype = {
      set: function(key, value) {
        var entry = key[this.name];
        if (entry && entry[0] === key) entry[1] = value; else defineProperty(key, this.name, {
          value: [ key, value ],
          writable: true
        });
        return this;
      },
      get: function(key) {
        var entry;
        return (entry = key[this.name]) && entry[0] === key ? entry[1] : undefined;
      },
      "delete": function(key) {
        var entry = key[this.name];
        if (!entry || entry[0] !== key) return false;
        entry[0] = entry[1] = undefined;
        return true;
      },
      has: function(key) {
        var entry = key[this.name];
        if (!entry) return false;
        return entry[0] === key;
      }
    };
    window.WeakMap = WeakMap;
  })();
}

(function(global) {
  if (global.JsMutationObserver) {
    return;
  }
  var registrationsTable = new WeakMap();
  var setImmediate;
  if (/Trident|Edge/.test(navigator.userAgent)) {
    setImmediate = setTimeout;
  } else if (window.setImmediate) {
    setImmediate = window.setImmediate;
  } else {
    var setImmediateQueue = [];
    var sentinel = String(Math.random());
    window.addEventListener("message", function(e) {
      if (e.data === sentinel) {
        var queue = setImmediateQueue;
        setImmediateQueue = [];
        queue.forEach(function(func) {
          func();
        });
      }
    });
    setImmediate = function(func) {
      setImmediateQueue.push(func);
      window.postMessage(sentinel, "*");
    };
  }
  var isScheduled = false;
  var scheduledObservers = [];
  function scheduleCallback(observer) {
    scheduledObservers.push(observer);
    if (!isScheduled) {
      isScheduled = true;
      setImmediate(dispatchCallbacks);
    }
  }
  function wrapIfNeeded(node) {
    return window.ShadowDOMPolyfill && window.ShadowDOMPolyfill.wrapIfNeeded(node) || node;
  }
  function dispatchCallbacks() {
    isScheduled = false;
    var observers = scheduledObservers;
    scheduledObservers = [];
    observers.sort(function(o1, o2) {
      return o1.uid_ - o2.uid_;
    });
    var anyNonEmpty = false;
    observers.forEach(function(observer) {
      var queue = observer.takeRecords();
      removeTransientObserversFor(observer);
      if (queue.length) {
        observer.callback_(queue, observer);
        anyNonEmpty = true;
      }
    });
    if (anyNonEmpty) dispatchCallbacks();
  }
  function removeTransientObserversFor(observer) {
    observer.nodes_.forEach(function(node) {
      var registrations = registrationsTable.get(node);
      if (!registrations) return;
      registrations.forEach(function(registration) {
        if (registration.observer === observer) registration.removeTransientObservers();
      });
    });
  }
  function forEachAncestorAndObserverEnqueueRecord(target, callback) {
    for (var node = target; node; node = node.parentNode) {
      var registrations = registrationsTable.get(node);
      if (registrations) {
        for (var j = 0; j < registrations.length; j++) {
          var registration = registrations[j];
          var options = registration.options;
          if (node !== target && !options.subtree) continue;
          var record = callback(options);
          if (record) registration.enqueue(record);
        }
      }
    }
  }
  var uidCounter = 0;
  function JsMutationObserver(callback) {
    this.callback_ = callback;
    this.nodes_ = [];
    this.records_ = [];
    this.uid_ = ++uidCounter;
  }
  JsMutationObserver.prototype = {
    observe: function(target, options) {
      target = wrapIfNeeded(target);
      if (!options.childList && !options.attributes && !options.characterData || options.attributeOldValue && !options.attributes || options.attributeFilter && options.attributeFilter.length && !options.attributes || options.characterDataOldValue && !options.characterData) {
        throw new SyntaxError();
      }
      var registrations = registrationsTable.get(target);
      if (!registrations) registrationsTable.set(target, registrations = []);
      var registration;
      for (var i = 0; i < registrations.length; i++) {
        if (registrations[i].observer === this) {
          registration = registrations[i];
          registration.removeListeners();
          registration.options = options;
          break;
        }
      }
      if (!registration) {
        registration = new Registration(this, target, options);
        registrations.push(registration);
        this.nodes_.push(target);
      }
      registration.addListeners();
    },
    disconnect: function() {
      this.nodes_.forEach(function(node) {
        var registrations = registrationsTable.get(node);
        for (var i = 0; i < registrations.length; i++) {
          var registration = registrations[i];
          if (registration.observer === this) {
            registration.removeListeners();
            registrations.splice(i, 1);
            break;
          }
        }
      }, this);
      this.records_ = [];
    },
    takeRecords: function() {
      var copyOfRecords = this.records_;
      this.records_ = [];
      return copyOfRecords;
    }
  };
  function MutationRecord(type, target) {
    this.type = type;
    this.target = target;
    this.addedNodes = [];
    this.removedNodes = [];
    this.previousSibling = null;
    this.nextSibling = null;
    this.attributeName = null;
    this.attributeNamespace = null;
    this.oldValue = null;
  }
  function copyMutationRecord(original) {
    var record = new MutationRecord(original.type, original.target);
    record.addedNodes = original.addedNodes.slice();
    record.removedNodes = original.removedNodes.slice();
    record.previousSibling = original.previousSibling;
    record.nextSibling = original.nextSibling;
    record.attributeName = original.attributeName;
    record.attributeNamespace = original.attributeNamespace;
    record.oldValue = original.oldValue;
    return record;
  }
  var currentRecord, recordWithOldValue;
  function getRecord(type, target) {
    return currentRecord = new MutationRecord(type, target);
  }
  function getRecordWithOldValue(oldValue) {
    if (recordWithOldValue) return recordWithOldValue;
    recordWithOldValue = copyMutationRecord(currentRecord);
    recordWithOldValue.oldValue = oldValue;
    return recordWithOldValue;
  }
  function clearRecords() {
    currentRecord = recordWithOldValue = undefined;
  }
  function recordRepresentsCurrentMutation(record) {
    return record === recordWithOldValue || record === currentRecord;
  }
  function selectRecord(lastRecord, newRecord) {
    if (lastRecord === newRecord) return lastRecord;
    if (recordWithOldValue && recordRepresentsCurrentMutation(lastRecord)) return recordWithOldValue;
    return null;
  }
  function Registration(observer, target, options) {
    this.observer = observer;
    this.target = target;
    this.options = options;
    this.transientObservedNodes = [];
  }
  Registration.prototype = {
    enqueue: function(record) {
      var records = this.observer.records_;
      var length = records.length;
      if (records.length > 0) {
        var lastRecord = records[length - 1];
        var recordToReplaceLast = selectRecord(lastRecord, record);
        if (recordToReplaceLast) {
          records[length - 1] = recordToReplaceLast;
          return;
        }
      } else {
        scheduleCallback(this.observer);
      }
      records[length] = record;
    },
    addListeners: function() {
      this.addListeners_(this.target);
    },
    addListeners_: function(node) {
      var options = this.options;
      if (options.attributes) node.addEventListener("DOMAttrModified", this, true);
      if (options.characterData) node.addEventListener("DOMCharacterDataModified", this, true);
      if (options.childList) node.addEventListener("DOMNodeInserted", this, true);
      if (options.childList || options.subtree) node.addEventListener("DOMNodeRemoved", this, true);
    },
    removeListeners: function() {
      this.removeListeners_(this.target);
    },
    removeListeners_: function(node) {
      var options = this.options;
      if (options.attributes) node.removeEventListener("DOMAttrModified", this, true);
      if (options.characterData) node.removeEventListener("DOMCharacterDataModified", this, true);
      if (options.childList) node.removeEventListener("DOMNodeInserted", this, true);
      if (options.childList || options.subtree) node.removeEventListener("DOMNodeRemoved", this, true);
    },
    addTransientObserver: function(node) {
      if (node === this.target) return;
      this.addListeners_(node);
      this.transientObservedNodes.push(node);
      var registrations = registrationsTable.get(node);
      if (!registrations) registrationsTable.set(node, registrations = []);
      registrations.push(this);
    },
    removeTransientObservers: function() {
      var transientObservedNodes = this.transientObservedNodes;
      this.transientObservedNodes = [];
      transientObservedNodes.forEach(function(node) {
        this.removeListeners_(node);
        var registrations = registrationsTable.get(node);
        for (var i = 0; i < registrations.length; i++) {
          if (registrations[i] === this) {
            registrations.splice(i, 1);
            break;
          }
        }
      }, this);
    },
    handleEvent: function(e) {
      e.stopImmediatePropagation();
      switch (e.type) {
       case "DOMAttrModified":
        var name = e.attrName;
        var namespace = e.relatedNode.namespaceURI;
        var target = e.target;
        var record = new getRecord("attributes", target);
        record.attributeName = name;
        record.attributeNamespace = namespace;
        var oldValue = e.attrChange === MutationEvent.ADDITION ? null : e.prevValue;
        forEachAncestorAndObserverEnqueueRecord(target, function(options) {
          if (!options.attributes) return;
          if (options.attributeFilter && options.attributeFilter.length && options.attributeFilter.indexOf(name) === -1 && options.attributeFilter.indexOf(namespace) === -1) {
            return;
          }
          if (options.attributeOldValue) return getRecordWithOldValue(oldValue);
          return record;
        });
        break;

       case "DOMCharacterDataModified":
        var target = e.target;
        var record = getRecord("characterData", target);
        var oldValue = e.prevValue;
        forEachAncestorAndObserverEnqueueRecord(target, function(options) {
          if (!options.characterData) return;
          if (options.characterDataOldValue) return getRecordWithOldValue(oldValue);
          return record;
        });
        break;

       case "DOMNodeRemoved":
        this.addTransientObserver(e.target);

       case "DOMNodeInserted":
        var changedNode = e.target;
        var addedNodes, removedNodes;
        if (e.type === "DOMNodeInserted") {
          addedNodes = [ changedNode ];
          removedNodes = [];
        } else {
          addedNodes = [];
          removedNodes = [ changedNode ];
        }
        var previousSibling = changedNode.previousSibling;
        var nextSibling = changedNode.nextSibling;
        var record = getRecord("childList", e.target.parentNode);
        record.addedNodes = addedNodes;
        record.removedNodes = removedNodes;
        record.previousSibling = previousSibling;
        record.nextSibling = nextSibling;
        forEachAncestorAndObserverEnqueueRecord(e.relatedNode, function(options) {
          if (!options.childList) return;
          return record;
        });
      }
      clearRecords();
    }
  };
  global.JsMutationObserver = JsMutationObserver;
  if (!global.MutationObserver) {
    global.MutationObserver = JsMutationObserver;
    JsMutationObserver._isPolyfilled = true;
  }
})(self);

(function() {
  var needsTemplate = typeof HTMLTemplateElement === "undefined";
  if (/Trident/.test(navigator.userAgent)) {
    (function() {
      var importNode = document.importNode;
      document.importNode = function() {
        var n = importNode.apply(document, arguments);
        if (n.nodeType === Node.DOCUMENT_FRAGMENT_NODE) {
          var f = document.createDocumentFragment();
          f.appendChild(n);
          return f;
        } else {
          return n;
        }
      };
    })();
  }
  var needsCloning = function() {
    if (!needsTemplate) {
      var t = document.createElement("template");
      var t2 = document.createElement("template");
      t2.content.appendChild(document.createElement("div"));
      t.content.appendChild(t2);
      var clone = t.cloneNode(true);
      return clone.content.childNodes.length === 0 || clone.content.firstChild.content.childNodes.length === 0;
    }
  }();
  var TEMPLATE_TAG = "template";
  var TemplateImpl = function() {};
  if (needsTemplate) {
    var contentDoc = document.implementation.createHTMLDocument("template");
    var canDecorate = true;
    var templateStyle = document.createElement("style");
    templateStyle.textContent = TEMPLATE_TAG + "{display:none;}";
    var head = document.head;
    head.insertBefore(templateStyle, head.firstElementChild);
    TemplateImpl.prototype = Object.create(HTMLElement.prototype);
    TemplateImpl.decorate = function(template) {
      if (template.content) {
        return;
      }
      template.content = contentDoc.createDocumentFragment();
      var child;
      while (child = template.firstChild) {
        template.content.appendChild(child);
      }
      template.cloneNode = function(deep) {
        return TemplateImpl.cloneNode(this, deep);
      };
      if (canDecorate) {
        try {
          Object.defineProperty(template, "innerHTML", {
            get: function() {
              var o = "";
              for (var e = this.content.firstChild; e; e = e.nextSibling) {
                o += e.outerHTML || escapeData(e.data);
              }
              return o;
            },
            set: function(text) {
              contentDoc.body.innerHTML = text;
              TemplateImpl.bootstrap(contentDoc);
              while (this.content.firstChild) {
                this.content.removeChild(this.content.firstChild);
              }
              while (contentDoc.body.firstChild) {
                this.content.appendChild(contentDoc.body.firstChild);
              }
            },
            configurable: true
          });
        } catch (err) {
          canDecorate = false;
        }
      }
      TemplateImpl.bootstrap(template.content);
    };
    TemplateImpl.bootstrap = function(doc) {
      var templates = doc.querySelectorAll(TEMPLATE_TAG);
      for (var i = 0, l = templates.length, t; i < l && (t = templates[i]); i++) {
        TemplateImpl.decorate(t);
      }
    };
    document.addEventListener("DOMContentLoaded", function() {
      TemplateImpl.bootstrap(document);
    });
    var createElement = document.createElement;
    document.createElement = function() {
      "use strict";
      var el = createElement.apply(document, arguments);
      if (el.localName === "template") {
        TemplateImpl.decorate(el);
      }
      return el;
    };
    var escapeDataRegExp = /[&\u00A0<>]/g;
    function escapeReplace(c) {
      switch (c) {
       case "&":
        return "&amp;";

       case "<":
        return "&lt;";

       case ">":
        return "&gt;";

       case " ":
        return "&nbsp;";
      }
    }
    function escapeData(s) {
      return s.replace(escapeDataRegExp, escapeReplace);
    }
  }
  if (needsTemplate || needsCloning) {
    var nativeCloneNode = Node.prototype.cloneNode;
    TemplateImpl.cloneNode = function(template, deep) {
      var clone = nativeCloneNode.call(template, false);
      if (this.decorate) {
        this.decorate(clone);
      }
      if (deep) {
        clone.content.appendChild(nativeCloneNode.call(template.content, true));
        this.fixClonedDom(clone.content, template.content);
      }
      return clone;
    };
    TemplateImpl.fixClonedDom = function(clone, source) {
      if (!source.querySelectorAll) return;
      var s$ = source.querySelectorAll(TEMPLATE_TAG);
      var t$ = clone.querySelectorAll(TEMPLATE_TAG);
      for (var i = 0, l = t$.length, t, s; i < l; i++) {
        s = s$[i];
        t = t$[i];
        if (this.decorate) {
          this.decorate(s);
        }
        t.parentNode.replaceChild(s.cloneNode(true), t);
      }
    };
    var originalImportNode = document.importNode;
    Node.prototype.cloneNode = function(deep) {
      var dom = nativeCloneNode.call(this, deep);
      if (deep) {
        TemplateImpl.fixClonedDom(dom, this);
      }
      return dom;
    };
    document.importNode = function(element, deep) {
      if (element.localName === TEMPLATE_TAG) {
        return TemplateImpl.cloneNode(element, deep);
      } else {
        var dom = originalImportNode.call(document, element, deep);
        if (deep) {
          TemplateImpl.fixClonedDom(dom, element);
        }
        return dom;
      }
    };
    if (needsCloning) {
      HTMLTemplateElement.prototype.cloneNode = function(deep) {
        return TemplateImpl.cloneNode(this, deep);
      };
    }
  }
  if (needsTemplate) {
    window.HTMLTemplateElement = TemplateImpl;
  }
})();

(function(scope) {
  "use strict";
  if (!(window.performance && window.performance.now)) {
    var start = Date.now();
    window.performance = {
      now: function() {
        return Date.now() - start;
      }
    };
  }
  if (!window.requestAnimationFrame) {
    window.requestAnimationFrame = function() {
      var nativeRaf = window.webkitRequestAnimationFrame || window.mozRequestAnimationFrame;
      return nativeRaf ? function(callback) {
        return nativeRaf(function() {
          callback(performance.now());
        });
      } : function(callback) {
        return window.setTimeout(callback, 1e3 / 60);
      };
    }();
  }
  if (!window.cancelAnimationFrame) {
    window.cancelAnimationFrame = function() {
      return window.webkitCancelAnimationFrame || window.mozCancelAnimationFrame || function(id) {
        clearTimeout(id);
      };
    }();
  }
  var workingDefaultPrevented = function() {
    var e = document.createEvent("Event");
    e.initEvent("foo", true, true);
    e.preventDefault();
    return e.defaultPrevented;
  }();
  if (!workingDefaultPrevented) {
    var origPreventDefault = Event.prototype.preventDefault;
    Event.prototype.preventDefault = function() {
      if (!this.cancelable) {
        return;
      }
      origPreventDefault.call(this);
      Object.defineProperty(this, "defaultPrevented", {
        get: function() {
          return true;
        },
        configurable: true
      });
    };
  }
  var isIE = /Trident/.test(navigator.userAgent);
  if (!window.CustomEvent || isIE && typeof window.CustomEvent !== "function") {
    window.CustomEvent = function(inType, params) {
      params = params || {};
      var e = document.createEvent("CustomEvent");
      e.initCustomEvent(inType, Boolean(params.bubbles), Boolean(params.cancelable), params.detail);
      return e;
    };
    window.CustomEvent.prototype = window.Event.prototype;
  }
  if (!window.Event || isIE && typeof window.Event !== "function") {
    var origEvent = window.Event;
    window.Event = function(inType, params) {
      params = params || {};
      var e = document.createEvent("Event");
      e.initEvent(inType, Boolean(params.bubbles), Boolean(params.cancelable));
      return e;
    };
    window.Event.prototype = origEvent.prototype;
  }
})(window.WebComponents);

window.HTMLImports = window.HTMLImports || {
  flags: {}
};

(function(scope) {
  var IMPORT_LINK_TYPE = "import";
  var useNative = Boolean(IMPORT_LINK_TYPE in document.createElement("link"));
  var hasShadowDOMPolyfill = Boolean(window.ShadowDOMPolyfill);
  var wrap = function(node) {
    return hasShadowDOMPolyfill ? window.ShadowDOMPolyfill.wrapIfNeeded(node) : node;
  };
  var rootDocument = wrap(document);
  var currentScriptDescriptor = {
    get: function() {
      var script = window.HTMLImports.currentScript || document.currentScript || (document.readyState !== "complete" ? document.scripts[document.scripts.length - 1] : null);
      return wrap(script);
    },
    configurable: true
  };
  Object.defineProperty(document, "_currentScript", currentScriptDescriptor);
  Object.defineProperty(rootDocument, "_currentScript", currentScriptDescriptor);
  var isIE = /Trident/.test(navigator.userAgent);
  function whenReady(callback, doc) {
    doc = doc || rootDocument;
    whenDocumentReady(function() {
      watchImportsLoad(callback, doc);
    }, doc);
  }
  var requiredReadyState = isIE ? "complete" : "interactive";
  var READY_EVENT = "readystatechange";
  function isDocumentReady(doc) {
    return doc.readyState === "complete" || doc.readyState === requiredReadyState;
  }
  function whenDocumentReady(callback, doc) {
    if (!isDocumentReady(doc)) {
      var checkReady = function() {
        if (doc.readyState === "complete" || doc.readyState === requiredReadyState) {
          doc.removeEventListener(READY_EVENT, checkReady);
          whenDocumentReady(callback, doc);
        }
      };
      doc.addEventListener(READY_EVENT, checkReady);
    } else if (callback) {
      callback();
    }
  }
  function markTargetLoaded(event) {
    event.target.__loaded = true;
  }
  function watchImportsLoad(callback, doc) {
    var imports = doc.querySelectorAll("link[rel=import]");
    var parsedCount = 0, importCount = imports.length, newImports = [], errorImports = [];
    function checkDone() {
      if (parsedCount == importCount && callback) {
        callback({
          allImports: imports,
          loadedImports: newImports,
          errorImports: errorImports
        });
      }
    }
    function loadedImport(e) {
      markTargetLoaded(e);
      newImports.push(this);
      parsedCount++;
      checkDone();
    }
    function errorLoadingImport(e) {
      errorImports.push(this);
      parsedCount++;
      checkDone();
    }
    if (importCount) {
      for (var i = 0, imp; i < importCount && (imp = imports[i]); i++) {
        if (isImportLoaded(imp)) {
          newImports.push(this);
          parsedCount++;
          checkDone();
        } else {
          imp.addEventListener("load", loadedImport);
          imp.addEventListener("error", errorLoadingImport);
        }
      }
    } else {
      checkDone();
    }
  }
  function isImportLoaded(link) {
    return useNative ? link.__loaded || link.import && link.import.readyState !== "loading" : link.__importParsed;
  }
  if (useNative) {
    new MutationObserver(function(mxns) {
      for (var i = 0, l = mxns.length, m; i < l && (m = mxns[i]); i++) {
        if (m.addedNodes) {
          handleImports(m.addedNodes);
        }
      }
    }).observe(document.head, {
      childList: true
    });
    function handleImports(nodes) {
      for (var i = 0, l = nodes.length, n; i < l && (n = nodes[i]); i++) {
        if (isImport(n)) {
          handleImport(n);
        }
      }
    }
    function isImport(element) {
      return element.localName === "link" && element.rel === "import";
    }
    function handleImport(element) {
      var loaded = element.import;
      if (loaded) {
        markTargetLoaded({
          target: element
        });
      } else {
        element.addEventListener("load", markTargetLoaded);
        element.addEventListener("error", markTargetLoaded);
      }
    }
    (function() {
      if (document.readyState === "loading") {
        var imports = document.querySelectorAll("link[rel=import]");
        for (var i = 0, l = imports.length, imp; i < l && (imp = imports[i]); i++) {
          handleImport(imp);
        }
      }
    })();
  }
  whenReady(function(detail) {
    window.HTMLImports.ready = true;
    window.HTMLImports.readyTime = new Date().getTime();
    var evt = rootDocument.createEvent("CustomEvent");
    evt.initCustomEvent("HTMLImportsLoaded", true, true, detail);
    rootDocument.dispatchEvent(evt);
  });
  scope.IMPORT_LINK_TYPE = IMPORT_LINK_TYPE;
  scope.useNative = useNative;
  scope.rootDocument = rootDocument;
  scope.whenReady = whenReady;
  scope.isIE = isIE;
})(window.HTMLImports);

(function(scope) {
  var modules = [];
  var addModule = function(module) {
    modules.push(module);
  };
  var initializeModules = function() {
    modules.forEach(function(module) {
      module(scope);
    });
  };
  scope.addModule = addModule;
  scope.initializeModules = initializeModules;
})(window.HTMLImports);

window.HTMLImports.addModule(function(scope) {
  var CSS_URL_REGEXP = /(url\()([^)]*)(\))/g;
  var CSS_IMPORT_REGEXP = /(@import[\s]+(?!url\())([^;]*)(;)/g;
  var path = {
    resolveUrlsInStyle: function(style, linkUrl) {
      var doc = style.ownerDocument;
      var resolver = doc.createElement("a");
      style.textContent = this.resolveUrlsInCssText(style.textContent, linkUrl, resolver);
      return style;
    },
    resolveUrlsInCssText: function(cssText, linkUrl, urlObj) {
      var r = this.replaceUrls(cssText, urlObj, linkUrl, CSS_URL_REGEXP);
      r = this.replaceUrls(r, urlObj, linkUrl, CSS_IMPORT_REGEXP);
      return r;
    },
    replaceUrls: function(text, urlObj, linkUrl, regexp) {
      return text.replace(regexp, function(m, pre, url, post) {
        var urlPath = url.replace(/["']/g, "");
        if (linkUrl) {
          urlPath = new URL(urlPath, linkUrl).href;
        }
        urlObj.href = urlPath;
        urlPath = urlObj.href;
        return pre + "'" + urlPath + "'" + post;
      });
    }
  };
  scope.path = path;
});

window.HTMLImports.addModule(function(scope) {
  var xhr = {
    async: true,
    ok: function(request) {
      return request.status >= 200 && request.status < 300 || request.status === 304 || request.status === 0;
    },
    load: function(url, next, nextContext) {
      var request = new XMLHttpRequest();
      if (scope.flags.debug || scope.flags.bust) {
        url += "?" + Math.random();
      }
      request.open("GET", url, xhr.async);
      request.addEventListener("readystatechange", function(e) {
        if (request.readyState === 4) {
          var redirectedUrl = null;
          try {
            var locationHeader = request.getResponseHeader("Location");
            if (locationHeader) {
              redirectedUrl = locationHeader.substr(0, 1) === "/" ? location.origin + locationHeader : locationHeader;
            }
          } catch (e) {
            console.error(e.message);
          }
          next.call(nextContext, !xhr.ok(request) && request, request.response || request.responseText, redirectedUrl);
        }
      });
      request.send();
      return request;
    },
    loadDocument: function(url, next, nextContext) {
      this.load(url, next, nextContext).responseType = "document";
    }
  };
  scope.xhr = xhr;
});

window.HTMLImports.addModule(function(scope) {
  var xhr = scope.xhr;
  var flags = scope.flags;
  var Loader = function(onLoad, onComplete) {
    this.cache = {};
    this.onload = onLoad;
    this.oncomplete = onComplete;
    this.inflight = 0;
    this.pending = {};
  };
  Loader.prototype = {
    addNodes: function(nodes) {
      this.inflight += nodes.length;
      for (var i = 0, l = nodes.length, n; i < l && (n = nodes[i]); i++) {
        this.require(n);
      }
      this.checkDone();
    },
    addNode: function(node) {
      this.inflight++;
      this.require(node);
      this.checkDone();
    },
    require: function(elt) {
      var url = elt.src || elt.href;
      elt.__nodeUrl = url;
      if (!this.dedupe(url, elt)) {
        this.fetch(url, elt);
      }
    },
    dedupe: function(url, elt) {
      if (this.pending[url]) {
        this.pending[url].push(elt);
        return true;
      }
      var resource;
      if (this.cache[url]) {
        this.onload(url, elt, this.cache[url]);
        this.tail();
        return true;
      }
      this.pending[url] = [ elt ];
      return false;
    },
    fetch: function(url, elt) {
      flags.load && console.log("fetch", url, elt);
      if (!url) {
        setTimeout(function() {
          this.receive(url, elt, {
            error: "href must be specified"
          }, null);
        }.bind(this), 0);
      } else if (url.match(/^data:/)) {
        var pieces = url.split(",");
        var header = pieces[0];
        var body = pieces[1];
        if (header.indexOf(";base64") > -1) {
          body = atob(body);
        } else {
          body = decodeURIComponent(body);
        }
        setTimeout(function() {
          this.receive(url, elt, null, body);
        }.bind(this), 0);
      } else {
        var receiveXhr = function(err, resource, redirectedUrl) {
          this.receive(url, elt, err, resource, redirectedUrl);
        }.bind(this);
        xhr.load(url, receiveXhr);
      }
    },
    receive: function(url, elt, err, resource, redirectedUrl) {
      this.cache[url] = resource;
      var $p = this.pending[url];
      for (var i = 0, l = $p.length, p; i < l && (p = $p[i]); i++) {
        this.onload(url, p, resource, err, redirectedUrl);
        this.tail();
      }
      this.pending[url] = null;
    },
    tail: function() {
      --this.inflight;
      this.checkDone();
    },
    checkDone: function() {
      if (!this.inflight) {
        this.oncomplete();
      }
    }
  };
  scope.Loader = Loader;
});

window.HTMLImports.addModule(function(scope) {
  var Observer = function(addCallback) {
    this.addCallback = addCallback;
    this.mo = new MutationObserver(this.handler.bind(this));
  };
  Observer.prototype = {
    handler: function(mutations) {
      for (var i = 0, l = mutations.length, m; i < l && (m = mutations[i]); i++) {
        if (m.type === "childList" && m.addedNodes.length) {
          this.addedNodes(m.addedNodes);
        }
      }
    },
    addedNodes: function(nodes) {
      if (this.addCallback) {
        this.addCallback(nodes);
      }
      for (var i = 0, l = nodes.length, n, loading; i < l && (n = nodes[i]); i++) {
        if (n.children && n.children.length) {
          this.addedNodes(n.children);
        }
      }
    },
    observe: function(root) {
      this.mo.observe(root, {
        childList: true,
        subtree: true
      });
    }
  };
  scope.Observer = Observer;
});

window.HTMLImports.addModule(function(scope) {
  var path = scope.path;
  var rootDocument = scope.rootDocument;
  var flags = scope.flags;
  var isIE = scope.isIE;
  var IMPORT_LINK_TYPE = scope.IMPORT_LINK_TYPE;
  var IMPORT_SELECTOR = "link[rel=" + IMPORT_LINK_TYPE + "]";
  var importParser = {
    documentSelectors: IMPORT_SELECTOR,
    importsSelectors: [ IMPORT_SELECTOR, "link[rel=stylesheet]:not([type])", "style:not([type])", "script:not([type])", 'script[type="application/javascript"]', 'script[type="text/javascript"]' ].join(","),
    map: {
      link: "parseLink",
      script: "parseScript",
      style: "parseStyle"
    },
    dynamicElements: [],
    parseNext: function() {
      var next = this.nextToParse();
      if (next) {
        this.parse(next);
      }
    },
    parse: function(elt) {
      if (this.isParsed(elt)) {
        flags.parse && console.log("[%s] is already parsed", elt.localName);
        return;
      }
      var fn = this[this.map[elt.localName]];
      if (fn) {
        this.markParsing(elt);
        fn.call(this, elt);
      }
    },
    parseDynamic: function(elt, quiet) {
      this.dynamicElements.push(elt);
      if (!quiet) {
        this.parseNext();
      }
    },
    markParsing: function(elt) {
      flags.parse && console.log("parsing", elt);
      this.parsingElement = elt;
    },
    markParsingComplete: function(elt) {
      elt.__importParsed = true;
      this.markDynamicParsingComplete(elt);
      if (elt.__importElement) {
        elt.__importElement.__importParsed = true;
        this.markDynamicParsingComplete(elt.__importElement);
      }
      this.parsingElement = null;
      flags.parse && console.log("completed", elt);
    },
    markDynamicParsingComplete: function(elt) {
      var i = this.dynamicElements.indexOf(elt);
      if (i >= 0) {
        this.dynamicElements.splice(i, 1);
      }
    },
    parseImport: function(elt) {
      elt.import = elt.__doc;
      if (window.HTMLImports.__importsParsingHook) {
        window.HTMLImports.__importsParsingHook(elt);
      }
      if (elt.import) {
        elt.import.__importParsed = true;
      }
      this.markParsingComplete(elt);
      if (elt.__resource && !elt.__error) {
        elt.dispatchEvent(new CustomEvent("load", {
          bubbles: false
        }));
      } else {
        elt.dispatchEvent(new CustomEvent("error", {
          bubbles: false
        }));
      }
      if (elt.__pending) {
        var fn;
        while (elt.__pending.length) {
          fn = elt.__pending.shift();
          if (fn) {
            fn({
              target: elt
            });
          }
        }
      }
      this.parseNext();
    },
    parseLink: function(linkElt) {
      if (nodeIsImport(linkElt)) {
        this.parseImport(linkElt);
      } else {
        linkElt.href = linkElt.href;
        this.parseGeneric(linkElt);
      }
    },
    parseStyle: function(elt) {
      var src = elt;
      elt = cloneStyle(elt);
      src.__appliedElement = elt;
      elt.__importElement = src;
      this.parseGeneric(elt);
    },
    parseGeneric: function(elt) {
      this.trackElement(elt);
      this.addElementToDocument(elt);
    },
    rootImportForElement: function(elt) {
      var n = elt;
      while (n.ownerDocument.__importLink) {
        n = n.ownerDocument.__importLink;
      }
      return n;
    },
    addElementToDocument: function(elt) {
      var port = this.rootImportForElement(elt.__importElement || elt);
      port.parentNode.insertBefore(elt, port);
    },
    trackElement: function(elt, callback) {
      var self = this;
      var done = function(e) {
        elt.removeEventListener("load", done);
        elt.removeEventListener("error", done);
        if (callback) {
          callback(e);
        }
        self.markParsingComplete(elt);
        self.parseNext();
      };
      elt.addEventListener("load", done);
      elt.addEventListener("error", done);
      if (isIE && elt.localName === "style") {
        var fakeLoad = false;
        if (elt.textContent.indexOf("@import") == -1) {
          fakeLoad = true;
        } else if (elt.sheet) {
          fakeLoad = true;
          var csr = elt.sheet.cssRules;
          var len = csr ? csr.length : 0;
          for (var i = 0, r; i < len && (r = csr[i]); i++) {
            if (r.type === CSSRule.IMPORT_RULE) {
              fakeLoad = fakeLoad && Boolean(r.styleSheet);
            }
          }
        }
        if (fakeLoad) {
          setTimeout(function() {
            elt.dispatchEvent(new CustomEvent("load", {
              bubbles: false
            }));
          });
        }
      }
    },
    parseScript: function(scriptElt) {
      var script = document.createElement("script");
      script.__importElement = scriptElt;
      script.src = scriptElt.src ? scriptElt.src : generateScriptDataUrl(scriptElt);
      scope.currentScript = scriptElt;
      this.trackElement(script, function(e) {
        if (script.parentNode) {
          script.parentNode.removeChild(script);
        }
        scope.currentScript = null;
      });
      this.addElementToDocument(script);
    },
    nextToParse: function() {
      this._mayParse = [];
      return !this.parsingElement && (this.nextToParseInDoc(rootDocument) || this.nextToParseDynamic());
    },
    nextToParseInDoc: function(doc, link) {
      if (doc && this._mayParse.indexOf(doc) < 0) {
        this._mayParse.push(doc);
        var nodes = doc.querySelectorAll(this.parseSelectorsForNode(doc));
        for (var i = 0, l = nodes.length, n; i < l && (n = nodes[i]); i++) {
          if (!this.isParsed(n)) {
            if (this.hasResource(n)) {
              return nodeIsImport(n) ? this.nextToParseInDoc(n.__doc, n) : n;
            } else {
              return;
            }
          }
        }
      }
      return link;
    },
    nextToParseDynamic: function() {
      return this.dynamicElements[0];
    },
    parseSelectorsForNode: function(node) {
      var doc = node.ownerDocument || node;
      return doc === rootDocument ? this.documentSelectors : this.importsSelectors;
    },
    isParsed: function(node) {
      return node.__importParsed;
    },
    needsDynamicParsing: function(elt) {
      return this.dynamicElements.indexOf(elt) >= 0;
    },
    hasResource: function(node) {
      if (nodeIsImport(node) && node.__doc === undefined) {
        return false;
      }
      return true;
    }
  };
  function nodeIsImport(elt) {
    return elt.localName === "link" && elt.rel === IMPORT_LINK_TYPE;
  }
  function generateScriptDataUrl(script) {
    var scriptContent = generateScriptContent(script);
    return "data:text/javascript;charset=utf-8," + encodeURIComponent(scriptContent);
  }
  function generateScriptContent(script) {
    return script.textContent + generateSourceMapHint(script);
  }
  function generateSourceMapHint(script) {
    var owner = script.ownerDocument;
    owner.__importedScripts = owner.__importedScripts || 0;
    var moniker = script.ownerDocument.baseURI;
    var num = owner.__importedScripts ? "-" + owner.__importedScripts : "";
    owner.__importedScripts++;
    return "\n//# sourceURL=" + moniker + num + ".js\n";
  }
  function cloneStyle(style) {
    var clone = style.ownerDocument.createElement("style");
    clone.textContent = style.textContent;
    path.resolveUrlsInStyle(clone);
    return clone;
  }
  scope.parser = importParser;
  scope.IMPORT_SELECTOR = IMPORT_SELECTOR;
});

window.HTMLImports.addModule(function(scope) {
  var flags = scope.flags;
  var IMPORT_LINK_TYPE = scope.IMPORT_LINK_TYPE;
  var IMPORT_SELECTOR = scope.IMPORT_SELECTOR;
  var rootDocument = scope.rootDocument;
  var Loader = scope.Loader;
  var Observer = scope.Observer;
  var parser = scope.parser;
  var importer = {
    documents: {},
    documentPreloadSelectors: IMPORT_SELECTOR,
    importsPreloadSelectors: [ IMPORT_SELECTOR ].join(","),
    loadNode: function(node) {
      importLoader.addNode(node);
    },
    loadSubtree: function(parent) {
      var nodes = this.marshalNodes(parent);
      importLoader.addNodes(nodes);
    },
    marshalNodes: function(parent) {
      return parent.querySelectorAll(this.loadSelectorsForNode(parent));
    },
    loadSelectorsForNode: function(node) {
      var doc = node.ownerDocument || node;
      return doc === rootDocument ? this.documentPreloadSelectors : this.importsPreloadSelectors;
    },
    loaded: function(url, elt, resource, err, redirectedUrl) {
      flags.load && console.log("loaded", url, elt);
      elt.__resource = resource;
      elt.__error = err;
      if (isImportLink(elt)) {
        var doc = this.documents[url];
        if (doc === undefined) {
          doc = err ? null : makeDocument(resource, redirectedUrl || url);
          if (doc) {
            doc.__importLink = elt;
            this.bootDocument(doc);
          }
          this.documents[url] = doc;
        }
        elt.__doc = doc;
      }
      parser.parseNext();
    },
    bootDocument: function(doc) {
      this.loadSubtree(doc);
      this.observer.observe(doc);
      parser.parseNext();
    },
    loadedAll: function() {
      parser.parseNext();
    }
  };
  var importLoader = new Loader(importer.loaded.bind(importer), importer.loadedAll.bind(importer));
  importer.observer = new Observer();
  function isImportLink(elt) {
    return isLinkRel(elt, IMPORT_LINK_TYPE);
  }
  function isLinkRel(elt, rel) {
    return elt.localName === "link" && elt.getAttribute("rel") === rel;
  }
  function hasBaseURIAccessor(doc) {
    return !!Object.getOwnPropertyDescriptor(doc, "baseURI");
  }
  function makeDocument(resource, url) {
    var doc = document.implementation.createHTMLDocument(IMPORT_LINK_TYPE);
    doc._URL = url;
    var base = doc.createElement("base");
    base.setAttribute("href", url);
    if (!doc.baseURI && !hasBaseURIAccessor(doc)) {
      Object.defineProperty(doc, "baseURI", {
        value: url
      });
    }
    var meta = doc.createElement("meta");
    meta.setAttribute("charset", "utf-8");
    doc.head.appendChild(meta);
    doc.head.appendChild(base);
    doc.body.innerHTML = resource;
    if (window.HTMLTemplateElement && HTMLTemplateElement.bootstrap) {
      HTMLTemplateElement.bootstrap(doc);
    }
    return doc;
  }
  if (!document.baseURI) {
    var baseURIDescriptor = {
      get: function() {
        var base = document.querySelector("base");
        return base ? base.href : window.location.href;
      },
      configurable: true
    };
    Object.defineProperty(document, "baseURI", baseURIDescriptor);
    Object.defineProperty(rootDocument, "baseURI", baseURIDescriptor);
  }
  scope.importer = importer;
  scope.importLoader = importLoader;
});

window.HTMLImports.addModule(function(scope) {
  var parser = scope.parser;
  var importer = scope.importer;
  var dynamic = {
    added: function(nodes) {
      var owner, parsed, loading;
      for (var i = 0, l = nodes.length, n; i < l && (n = nodes[i]); i++) {
        if (!owner) {
          owner = n.ownerDocument;
          parsed = parser.isParsed(owner);
        }
        loading = this.shouldLoadNode(n);
        if (loading) {
          importer.loadNode(n);
        }
        if (this.shouldParseNode(n) && parsed) {
          parser.parseDynamic(n, loading);
        }
      }
    },
    shouldLoadNode: function(node) {
      return node.nodeType === 1 && matches.call(node, importer.loadSelectorsForNode(node));
    },
    shouldParseNode: function(node) {
      return node.nodeType === 1 && matches.call(node, parser.parseSelectorsForNode(node));
    }
  };
  importer.observer.addCallback = dynamic.added.bind(dynamic);
  var matches = HTMLElement.prototype.matches || HTMLElement.prototype.matchesSelector || HTMLElement.prototype.webkitMatchesSelector || HTMLElement.prototype.mozMatchesSelector || HTMLElement.prototype.msMatchesSelector;
});

(function(scope) {
  var initializeModules = scope.initializeModules;
  var isIE = scope.isIE;
  if (scope.useNative) {
    return;
  }
  initializeModules();
  var rootDocument = scope.rootDocument;
  function bootstrap() {
    window.HTMLImports.importer.bootDocument(rootDocument);
  }
  if (document.readyState === "complete" || document.readyState === "interactive" && !window.attachEvent) {
    bootstrap();
  } else {
    document.addEventListener("DOMContentLoaded", bootstrap);
  }
})(window.HTMLImports);

window.CustomElements = window.CustomElements || {
  flags: {}
};

(function(scope) {
  var flags = scope.flags;
  var modules = [];
  var addModule = function(module) {
    modules.push(module);
  };
  var initializeModules = function() {
    modules.forEach(function(module) {
      module(scope);
    });
  };
  scope.addModule = addModule;
  scope.initializeModules = initializeModules;
  scope.hasNative = Boolean(document.registerElement);
  scope.isIE = /Trident/.test(navigator.userAgent);
  scope.useNative = !flags.register && scope.hasNative && !window.ShadowDOMPolyfill && (!window.HTMLImports || window.HTMLImports.useNative);
})(window.CustomElements);

window.CustomElements.addModule(function(scope) {
  var IMPORT_LINK_TYPE = window.HTMLImports ? window.HTMLImports.IMPORT_LINK_TYPE : "none";
  function forSubtree(node, cb) {
    findAllElements(node, function(e) {
      if (cb(e)) {
        return true;
      }
      forRoots(e, cb);
    });
    forRoots(node, cb);
  }
  function findAllElements(node, find, data) {
    var e = node.firstElementChild;
    if (!e) {
      e = node.firstChild;
      while (e && e.nodeType !== Node.ELEMENT_NODE) {
        e = e.nextSibling;
      }
    }
    while (e) {
      if (find(e, data) !== true) {
        findAllElements(e, find, data);
      }
      e = e.nextElementSibling;
    }
    return null;
  }
  function forRoots(node, cb) {
    var root = node.shadowRoot;
    while (root) {
      forSubtree(root, cb);
      root = root.olderShadowRoot;
    }
  }
  function forDocumentTree(doc, cb) {
    _forDocumentTree(doc, cb, []);
  }
  function _forDocumentTree(doc, cb, processingDocuments) {
    doc = window.wrap(doc);
    if (processingDocuments.indexOf(doc) >= 0) {
      return;
    }
    processingDocuments.push(doc);
    var imports = doc.querySelectorAll("link[rel=" + IMPORT_LINK_TYPE + "]");
    for (var i = 0, l = imports.length, n; i < l && (n = imports[i]); i++) {
      if (n.import) {
        _forDocumentTree(n.import, cb, processingDocuments);
      }
    }
    cb(doc);
  }
  scope.forDocumentTree = forDocumentTree;
  scope.forSubtree = forSubtree;
});

window.CustomElements.addModule(function(scope) {
  var flags = scope.flags;
  var forSubtree = scope.forSubtree;
  var forDocumentTree = scope.forDocumentTree;
  function addedNode(node, isAttached) {
    return added(node, isAttached) || addedSubtree(node, isAttached);
  }
  function added(node, isAttached) {
    if (scope.upgrade(node, isAttached)) {
      return true;
    }
    if (isAttached) {
      attached(node);
    }
  }
  function addedSubtree(node, isAttached) {
    forSubtree(node, function(e) {
      if (added(e, isAttached)) {
        return true;
      }
    });
  }
  var hasThrottledAttached = window.MutationObserver._isPolyfilled && flags["throttle-attached"];
  scope.hasPolyfillMutations = hasThrottledAttached;
  scope.hasThrottledAttached = hasThrottledAttached;
  var isPendingMutations = false;
  var pendingMutations = [];
  function deferMutation(fn) {
    pendingMutations.push(fn);
    if (!isPendingMutations) {
      isPendingMutations = true;
      setTimeout(takeMutations);
    }
  }
  function takeMutations() {
    isPendingMutations = false;
    var $p = pendingMutations;
    for (var i = 0, l = $p.length, p; i < l && (p = $p[i]); i++) {
      p();
    }
    pendingMutations = [];
  }
  function attached(element) {
    if (hasThrottledAttached) {
      deferMutation(function() {
        _attached(element);
      });
    } else {
      _attached(element);
    }
  }
  function _attached(element) {
    if (element.__upgraded__ && !element.__attached) {
      element.__attached = true;
      if (element.attachedCallback) {
        element.attachedCallback();
      }
    }
  }
  function detachedNode(node) {
    detached(node);
    forSubtree(node, function(e) {
      detached(e);
    });
  }
  function detached(element) {
    if (hasThrottledAttached) {
      deferMutation(function() {
        _detached(element);
      });
    } else {
      _detached(element);
    }
  }
  function _detached(element) {
    if (element.__upgraded__ && element.__attached) {
      element.__attached = false;
      if (element.detachedCallback) {
        element.detachedCallback();
      }
    }
  }
  function inDocument(element) {
    var p = element;
    var doc = window.wrap(document);
    while (p) {
      if (p == doc) {
        return true;
      }
      p = p.parentNode || p.nodeType === Node.DOCUMENT_FRAGMENT_NODE && p.host;
    }
  }
  function watchShadow(node) {
    if (node.shadowRoot && !node.shadowRoot.__watched) {
      flags.dom && console.log("watching shadow-root for: ", node.localName);
      var root = node.shadowRoot;
      while (root) {
        observe(root);
        root = root.olderShadowRoot;
      }
    }
  }
  function handler(root, mutations) {
    if (flags.dom) {
      var mx = mutations[0];
      if (mx && mx.type === "childList" && mx.addedNodes) {
        if (mx.addedNodes) {
          var d = mx.addedNodes[0];
          while (d && d !== document && !d.host) {
            d = d.parentNode;
          }
          var u = d && (d.URL || d._URL || d.host && d.host.localName) || "";
          u = u.split("/?").shift().split("/").pop();
        }
      }
      console.group("mutations (%d) [%s]", mutations.length, u || "");
    }
    var isAttached = inDocument(root);
    mutations.forEach(function(mx) {
      if (mx.type === "childList") {
        forEach(mx.addedNodes, function(n) {
          if (!n.localName) {
            return;
          }
          addedNode(n, isAttached);
        });
        forEach(mx.removedNodes, function(n) {
          if (!n.localName) {
            return;
          }
          detachedNode(n);
        });
      }
    });
    flags.dom && console.groupEnd();
  }
  function takeRecords(node) {
    node = window.wrap(node);
    if (!node) {
      node = window.wrap(document);
    }
    while (node.parentNode) {
      node = node.parentNode;
    }
    var observer = node.__observer;
    if (observer) {
      handler(node, observer.takeRecords());
      takeMutations();
    }
  }
  var forEach = Array.prototype.forEach.call.bind(Array.prototype.forEach);
  function observe(inRoot) {
    if (inRoot.__observer) {
      return;
    }
    var observer = new MutationObserver(handler.bind(this, inRoot));
    observer.observe(inRoot, {
      childList: true,
      subtree: true
    });
    inRoot.__observer = observer;
  }
  function upgradeDocument(doc) {
    doc = window.wrap(doc);
    flags.dom && console.group("upgradeDocument: ", doc.baseURI.split("/").pop());
    var isMainDocument = doc === window.wrap(document);
    addedNode(doc, isMainDocument);
    observe(doc);
    flags.dom && console.groupEnd();
  }
  function upgradeDocumentTree(doc) {
    forDocumentTree(doc, upgradeDocument);
  }
  var originalCreateShadowRoot = Element.prototype.createShadowRoot;
  if (originalCreateShadowRoot) {
    Element.prototype.createShadowRoot = function() {
      var root = originalCreateShadowRoot.call(this);
      window.CustomElements.watchShadow(this);
      return root;
    };
  }
  scope.watchShadow = watchShadow;
  scope.upgradeDocumentTree = upgradeDocumentTree;
  scope.upgradeDocument = upgradeDocument;
  scope.upgradeSubtree = addedSubtree;
  scope.upgradeAll = addedNode;
  scope.attached = attached;
  scope.takeRecords = takeRecords;
});

window.CustomElements.addModule(function(scope) {
  var flags = scope.flags;
  function upgrade(node, isAttached) {
    if (node.localName === "template") {
      if (window.HTMLTemplateElement && HTMLTemplateElement.decorate) {
        HTMLTemplateElement.decorate(node);
      }
    }
    if (!node.__upgraded__ && node.nodeType === Node.ELEMENT_NODE) {
      var is = node.getAttribute("is");
      var definition = scope.getRegisteredDefinition(node.localName) || scope.getRegisteredDefinition(is);
      if (definition) {
        if (is && definition.tag == node.localName || !is && !definition.extends) {
          return upgradeWithDefinition(node, definition, isAttached);
        }
      }
    }
  }
  function upgradeWithDefinition(element, definition, isAttached) {
    flags.upgrade && console.group("upgrade:", element.localName);
    if (definition.is) {
      element.setAttribute("is", definition.is);
    }
    implementPrototype(element, definition);
    element.__upgraded__ = true;
    created(element);
    if (isAttached) {
      scope.attached(element);
    }
    scope.upgradeSubtree(element, isAttached);
    flags.upgrade && console.groupEnd();
    return element;
  }
  function implementPrototype(element, definition) {
    if (Object.__proto__) {
      element.__proto__ = definition.prototype;
    } else {
      customMixin(element, definition.prototype, definition.native);
      element.__proto__ = definition.prototype;
    }
  }
  function customMixin(inTarget, inSrc, inNative) {
    var used = {};
    var p = inSrc;
    while (p !== inNative && p !== HTMLElement.prototype) {
      var keys = Object.getOwnPropertyNames(p);
      for (var i = 0, k; k = keys[i]; i++) {
        if (!used[k]) {
          Object.defineProperty(inTarget, k, Object.getOwnPropertyDescriptor(p, k));
          used[k] = 1;
        }
      }
      p = Object.getPrototypeOf(p);
    }
  }
  function created(element) {
    if (element.createdCallback) {
      element.createdCallback();
    }
  }
  scope.upgrade = upgrade;
  scope.upgradeWithDefinition = upgradeWithDefinition;
  scope.implementPrototype = implementPrototype;
});

window.CustomElements.addModule(function(scope) {
  var isIE = scope.isIE;
  var upgradeDocumentTree = scope.upgradeDocumentTree;
  var upgradeAll = scope.upgradeAll;
  var upgradeWithDefinition = scope.upgradeWithDefinition;
  var implementPrototype = scope.implementPrototype;
  var useNative = scope.useNative;
  function register(name, options) {
    var definition = options || {};
    if (!name) {
      throw new Error("document.registerElement: first argument `name` must not be empty");
    }
    if (name.indexOf("-") < 0) {
      throw new Error("document.registerElement: first argument ('name') must contain a dash ('-'). Argument provided was '" + String(name) + "'.");
    }
    if (isReservedTag(name)) {
      throw new Error("Failed to execute 'registerElement' on 'Document': Registration failed for type '" + String(name) + "'. The type name is invalid.");
    }
    if (getRegisteredDefinition(name)) {
      throw new Error("DuplicateDefinitionError: a type with name '" + String(name) + "' is already registered");
    }
    if (!definition.prototype) {
      definition.prototype = Object.create(HTMLElement.prototype);
    }
    definition.__name = name.toLowerCase();
    if (definition.extends) {
      definition.extends = definition.extends.toLowerCase();
    }
    definition.lifecycle = definition.lifecycle || {};
    definition.ancestry = ancestry(definition.extends);
    resolveTagName(definition);
    resolvePrototypeChain(definition);
    overrideAttributeApi(definition.prototype);
    registerDefinition(definition.__name, definition);
    definition.ctor = generateConstructor(definition);
    definition.ctor.prototype = definition.prototype;
    definition.prototype.constructor = definition.ctor;
    if (scope.ready) {
      upgradeDocumentTree(document);
    }
    return definition.ctor;
  }
  function overrideAttributeApi(prototype) {
    if (prototype.setAttribute._polyfilled) {
      return;
    }
    var setAttribute = prototype.setAttribute;
    prototype.setAttribute = function(name, value) {
      changeAttribute.call(this, name, value, setAttribute);
    };
    var removeAttribute = prototype.removeAttribute;
    prototype.removeAttribute = function(name) {
      changeAttribute.call(this, name, null, removeAttribute);
    };
    prototype.setAttribute._polyfilled = true;
  }
  function changeAttribute(name, value, operation) {
    name = name.toLowerCase();
    var oldValue = this.getAttribute(name);
    operation.apply(this, arguments);
    var newValue = this.getAttribute(name);
    if (this.attributeChangedCallback && newValue !== oldValue) {
      this.attributeChangedCallback(name, oldValue, newValue);
    }
  }
  function isReservedTag(name) {
    for (var i = 0; i < reservedTagList.length; i++) {
      if (name === reservedTagList[i]) {
        return true;
      }
    }
  }
  var reservedTagList = [ "annotation-xml", "color-profile", "font-face", "font-face-src", "font-face-uri", "font-face-format", "font-face-name", "missing-glyph" ];
  function ancestry(extnds) {
    var extendee = getRegisteredDefinition(extnds);
    if (extendee) {
      return ancestry(extendee.extends).concat([ extendee ]);
    }
    return [];
  }
  function resolveTagName(definition) {
    var baseTag = definition.extends;
    for (var i = 0, a; a = definition.ancestry[i]; i++) {
      baseTag = a.is && a.tag;
    }
    definition.tag = baseTag || definition.__name;
    if (baseTag) {
      definition.is = definition.__name;
    }
  }
  function resolvePrototypeChain(definition) {
    if (!Object.__proto__) {
      var nativePrototype = HTMLElement.prototype;
      if (definition.is) {
        var inst = document.createElement(definition.tag);
        nativePrototype = Object.getPrototypeOf(inst);
      }
      var proto = definition.prototype, ancestor;
      var foundPrototype = false;
      while (proto) {
        if (proto == nativePrototype) {
          foundPrototype = true;
        }
        ancestor = Object.getPrototypeOf(proto);
        if (ancestor) {
          proto.__proto__ = ancestor;
        }
        proto = ancestor;
      }
      if (!foundPrototype) {
        console.warn(definition.tag + " prototype not found in prototype chain for " + definition.is);
      }
      definition.native = nativePrototype;
    }
  }
  function instantiate(definition) {
    return upgradeWithDefinition(domCreateElement(definition.tag), definition);
  }
  var registry = {};
  function getRegisteredDefinition(name) {
    if (name) {
      return registry[name.toLowerCase()];
    }
  }
  function registerDefinition(name, definition) {
    registry[name] = definition;
  }
  function generateConstructor(definition) {
    return function() {
      return instantiate(definition);
    };
  }
  var HTML_NAMESPACE = "http://www.w3.org/1999/xhtml";
  function createElementNS(namespace, tag, typeExtension) {
    if (namespace === HTML_NAMESPACE) {
      return createElement(tag, typeExtension);
    } else {
      return domCreateElementNS(namespace, tag);
    }
  }
  function createElement(tag, typeExtension) {
    if (tag) {
      tag = tag.toLowerCase();
    }
    if (typeExtension) {
      typeExtension = typeExtension.toLowerCase();
    }
    var definition = getRegisteredDefinition(typeExtension || tag);
    if (definition) {
      if (tag == definition.tag && typeExtension == definition.is) {
        return new definition.ctor();
      }
      if (!typeExtension && !definition.is) {
        return new definition.ctor();
      }
    }
    var element;
    if (typeExtension) {
      element = createElement(tag);
      element.setAttribute("is", typeExtension);
      return element;
    }
    element = domCreateElement(tag);
    if (tag.indexOf("-") >= 0) {
      implementPrototype(element, HTMLElement);
    }
    return element;
  }
  var domCreateElement = document.createElement.bind(document);
  var domCreateElementNS = document.createElementNS.bind(document);
  var isInstance;
  if (!Object.__proto__ && !useNative) {
    isInstance = function(obj, ctor) {
      if (obj instanceof ctor) {
        return true;
      }
      var p = obj;
      while (p) {
        if (p === ctor.prototype) {
          return true;
        }
        p = p.__proto__;
      }
      return false;
    };
  } else {
    isInstance = function(obj, base) {
      return obj instanceof base;
    };
  }
  function wrapDomMethodToForceUpgrade(obj, methodName) {
    var orig = obj[methodName];
    obj[methodName] = function() {
      var n = orig.apply(this, arguments);
      upgradeAll(n);
      return n;
    };
  }
  wrapDomMethodToForceUpgrade(Node.prototype, "cloneNode");
  wrapDomMethodToForceUpgrade(document, "importNode");
  document.registerElement = register;
  document.createElement = createElement;
  document.createElementNS = createElementNS;
  scope.registry = registry;
  scope.instanceof = isInstance;
  scope.reservedTagList = reservedTagList;
  scope.getRegisteredDefinition = getRegisteredDefinition;
  document.register = document.registerElement;
});

(function(scope) {
  var useNative = scope.useNative;
  var initializeModules = scope.initializeModules;
  var isIE = scope.isIE;
  if (useNative) {
    var nop = function() {};
    scope.watchShadow = nop;
    scope.upgrade = nop;
    scope.upgradeAll = nop;
    scope.upgradeDocumentTree = nop;
    scope.upgradeSubtree = nop;
    scope.takeRecords = nop;
    scope.instanceof = function(obj, base) {
      return obj instanceof base;
    };
  } else {
    initializeModules();
  }
  var upgradeDocumentTree = scope.upgradeDocumentTree;
  var upgradeDocument = scope.upgradeDocument;
  if (!window.wrap) {
    if (window.ShadowDOMPolyfill) {
      window.wrap = window.ShadowDOMPolyfill.wrapIfNeeded;
      window.unwrap = window.ShadowDOMPolyfill.unwrapIfNeeded;
    } else {
      window.wrap = window.unwrap = function(node) {
        return node;
      };
    }
  }
  if (window.HTMLImports) {
    window.HTMLImports.__importsParsingHook = function(elt) {
      if (elt.import) {
        upgradeDocument(wrap(elt.import));
      }
    };
  }
  function bootstrap() {
    upgradeDocumentTree(window.wrap(document));
    window.CustomElements.ready = true;
    var requestAnimationFrame = window.requestAnimationFrame || function(f) {
      setTimeout(f, 16);
    };
    requestAnimationFrame(function() {
      setTimeout(function() {
        window.CustomElements.readyTime = Date.now();
        if (window.HTMLImports) {
          window.CustomElements.elapsed = window.CustomElements.readyTime - window.HTMLImports.readyTime;
        }
        document.dispatchEvent(new CustomEvent("WebComponentsReady", {
          bubbles: true
        }));
      });
    });
  }
  if (document.readyState === "complete" || scope.flags.eager) {
    bootstrap();
  } else if (document.readyState === "interactive" && !window.attachEvent && (!window.HTMLImports || window.HTMLImports.ready)) {
    bootstrap();
  } else {
    var loadEvent = window.HTMLImports && !window.HTMLImports.ready ? "HTMLImportsLoaded" : "DOMContentLoaded";
    window.addEventListener(loadEvent, bootstrap);
  }
})(window.CustomElements);

(function(scope) {
  var style = document.createElement("style");
  style.textContent = "" + "body {" + "transition: opacity ease-in 0.2s;" + " } \n" + "body[unresolved] {" + "opacity: 0; display: block; overflow: hidden; position: relative;" + " } \n";
  var head = document.querySelector("head");
  head.insertBefore(style, head.firstChild);
})(window.WebComponents);
},{}],2:[function(require,module,exports){
"use strict";

function _defineProperty(e, r, t) { return (r = _toPropertyKey(r)) in e ? Object.defineProperty(e, r, { value: t, enumerable: !0, configurable: !0, writable: !0 }) : e[r] = t, e; }
function _toPropertyKey(t) { var i = _toPrimitive(t, "string"); return "symbol" == typeof i ? i : i + ""; }
function _toPrimitive(t, r) { if ("object" != typeof t || !t) return t; var e = t[Symbol.toPrimitive]; if (void 0 !== e) { var i = e.call(t, r || "default"); if ("object" != typeof i) return i; throw new TypeError("@@toPrimitive must return a primitive value."); } return ("string" === r ? String : Number)(t); }
/** @lends ChromecastSessionManager.prototype **/
class ChromecastSessionManager {
  /**
   * Stores the state of the current Chromecast session and its associated objects such
   * as the
   * [RemotePlayerController](https://developers.google.com/cast/docs/reference/chrome/cast.framework.RemotePlayerController),
   * and the
   * [RemotePlayer](https://developers.google.com/cast/docs/reference/chrome/cast.framework.RemotePlayer).
   *
   * WARNING: Do not instantiate this class until the
   * [CastContext](https://developers.google.com/cast/docs/reference/chrome/cast.framework.CastContext)
   * has been configured.
   *
   * For an undocumented (and thus unknown) reason, RemotePlayer and
   * RemotePlayerController instances created before the cast context has been configured
   * or after requesting a session or loading media will not stay in sync with media
   * items that are loaded later.
   *
   * For example, the first item that you cast will work as expected: events on
   * RemotePlayerController will fire and the state (currentTime, duration, etc) of the
   * RemotePlayer instance will update as the media item plays. However, if a new media
   * item is loaded via a `loadMedia` request, the media item will play, but the
   * remotePlayer will be in a "media unloaded" state where the duration is 0, the
   * currentTime does not update, and no change events are fired (except, strangely,
   * displayStatus updates).
   *
   * @param player {object} Video.js Player
   * @constructs ChromecastSessionManager
   */
  constructor(player) {
    this.player = player;
    this._sessionListener = this._onSessionStateChange.bind(this);
    this._castListener = this._onCastStateChange.bind(this);
    this._addCastContextEventListeners();

    // Remove global event listeners when this player instance is destroyed to prevent
    // memory leaks.
    this.player.on('dispose', this._removeCastContextEventListeners.bind(this));
    this._notifyPlayerOfDevicesAvailabilityChange(this.getCastContext().getCastState());
    this.remotePlayer = new cast.framework.RemotePlayer();
    this.remotePlayerController = new cast.framework.RemotePlayerController(this.remotePlayer);
  }
  /**
   * Add event listeners for events triggered on the current CastContext.
   *
   * @private
   */
  _addCastContextEventListeners() {
    var sessionStateChangedEvt = cast.framework.CastContextEventType.SESSION_STATE_CHANGED,
      castStateChangedEvt = cast.framework.CastContextEventType.CAST_STATE_CHANGED;
    this.getCastContext().addEventListener(sessionStateChangedEvt, this._sessionListener);
    this.getCastContext().addEventListener(castStateChangedEvt, this._castListener);
  }

  /**
   * Remove event listeners that were added in {@link
   * ChromecastSessionManager#_addCastContextEventListeners}.
   *
   * @private
   */
  _removeCastContextEventListeners() {
    var sessionStateChangedEvt = cast.framework.CastContextEventType.SESSION_STATE_CHANGED,
      castStateChangedEvt = cast.framework.CastContextEventType.CAST_STATE_CHANGED;
    this.getCastContext().removeEventListener(sessionStateChangedEvt, this._sessionListener);
    this.getCastContext().removeEventListener(castStateChangedEvt, this._castListener);
  }

  /**
   * Handle the CastContext's SessionState change event.
   *
   * @private
   */
  _onSessionStateChange(event) {
    if (event.sessionState === cast.framework.SessionState.SESSION_ENDED) {
      this.player.trigger('chromecastDisconnected');
      this._reloadTech();
    }
  }

  /**
   * Handle the CastContext's CastState change event.
   *
   * @private
   */
  _onCastStateChange(event) {
    this._notifyPlayerOfDevicesAvailabilityChange(event.castState);
  }

  /**
   * Triggers player events that notifies listeners that Chromecast devices are
   * either available or unavailable.
   *
   * @private
   */
  _notifyPlayerOfDevicesAvailabilityChange(castState) {
    if (this.hasAvailableDevices(castState)) {
      this.player.trigger('chromecastDevicesAvailable');
    } else {
      this.player.trigger('chromecastDevicesUnavailable');
    }
  }

  /**
   * Returns whether or not there are Chromecast devices available to cast to.
   *
   * @see https://developers.google.com/cast/docs/reference/chrome/cast.framework#.CastState
   * @param {String} castState
   * @return {boolean} true if there are Chromecast devices available to cast to.
   */
  hasAvailableDevices(castState) {
    castState = castState || this.getCastContext().getCastState();
    return castState === cast.framework.CastState.NOT_CONNECTED || castState === cast.framework.CastState.CONNECTING || castState === cast.framework.CastState.CONNECTED;
  }

  /**
   * Opens the Chromecast casting menu by requesting a CastSession. Does nothing if the
   * Video.js player does not have a source.
   */
  openCastMenu() {
    var onSessionSuccess;
    if (!this.player.currentSource()) {
      // Do not cast if there is no media item loaded in the player
      return;
    }
    onSessionSuccess = function () {
      ChromecastSessionManager.hasConnected = true;
      this.player.trigger('chromecastConnected');
      this._reloadTech();
    }.bind(this);

    // It is the `requestSession` function call that actually causes the cast menu to
    // open.
    // The second parameter to `.then` is an error handler. We use a noop function here
    // because we handle errors in the ChromecastTech class and we do not want an
    // error to bubble up to the console. This error handler is also triggered when
    // the user closes out of the chromecast selector pop-up without choosing a
    // casting destination.
    this.getCastContext().requestSession().then(onSessionSuccess, function () {/* noop */});
  }

  /**
   * Reloads the Video.js player's Tech. This causes the player to re-evaluate which
   * Tech should be used for the current source by iterating over available Tech and
   * calling `Tech.isSupported` and `Tech.canPlaySource`. Video.js uses the first
   * Tech that returns true from both of those functions. This is what allows us to
   * switch back and forth between the Chromecast Tech and other available Tech when a
   * CastSession is connected or disconnected.
   *
   * @private
   */
  _reloadTech() {
    var player = this.player,
      currentTime = player.currentTime(),
      wasPaused = player.paused(),
      sources = player.currentSources();

    // Reload the current source(s) to re-lookup and use the currently available Tech.
    // The chromecast Tech gets used if `ChromecastSessionManager.isChromecastConnected`
    // is true (effectively, if a chromecast session is currently in progress),
    // otherwise Video.js continues to search through the Tech list for other eligible
    // Tech to use, such as the HTML5 player.
    player.src(sources);
    player.ready(function () {
      if (wasPaused) {
        player.pause();
      } else {
        player.play();
      }
      player.currentTime(currentTime || 0);
    });
  }

  /**
   * @see https://developers.google.com/cast/docs/reference/chrome/cast.framework.CastContext
   * @returns {object} the current CastContext, if one exists
   */

  getCastContext() {
    return cast.framework.CastContext.getInstance();
  }

  /**
   * @see https://developers.google.com/cast/docs/reference/chrome/cast.framework.RemotePlayer
   * @returns {object} the current RemotePlayer, if one exists
   */
  getRemotePlayer() {
    return this.remotePlayer;
  }

  /**
   * @see https://developers.google.com/cast/docs/reference/chrome/cast.framework.RemotePlayerController
   * @returns {object} the current RemotePlayerController, if one exists
   */
  getRemotePlayerController() {
    return this.remotePlayerController;
  }

  /**
   * Returns whether or not the current Chromecast API is available (that is,
   * `window.chrome`, `window.chrome.cast`, and `window.cast` exist).
   *
   * @static
   * @returns {boolean} true if the Chromecast API is available
   */
  static isChromecastAPIAvailable() {
    return window.chrome && window.chrome.cast && window.cast;
  }

  /**
   * Returns whether or not there is a current CastSession and it is connected.
   *
   * @static
   * @returns {boolean} true if the current CastSession exists and is connected
   */
  static isChromecastConnected() {
    // We must also check the `hasConnected` flag because
    // `getCastContext().getCastState()` returns `CONNECTED` even when the current
    // casting session was initiated by another tab in the browser or by another process
    return ChromecastSessionManager.isChromecastAPIAvailable() && cast.framework.CastContext.getInstance().getCastState() === cast.framework.CastState.CONNECTED && ChromecastSessionManager.hasConnected;
  }
}
_defineProperty(ChromecastSessionManager, "hasConnected", false);
module.exports = ChromecastSessionManager;

},{}],3:[function(require,module,exports){
"use strict";

module.exports = function (videojs) {
  /**
   * Registers the ChromecastButton Component with Video.js. Calls
   * {@link http://docs.videojs.com/Component.html#.registerComponent}, which will add a
   * component called `chromecastButton` to the list of globally registered Video.js
   * components. The `chromecastButton` is added to the player's control bar UI
   * automatically once {@link module:enableChromecast} has been called. If you would
   * like to specify the order of the buttons that appear in the control bar, including
   * this button, you can do so in the options that you pass to the `videojs` function
   * when creating a player:
   *
   * ```
   * videojs('playerID', {
   *    controlBar: {
   *       children: [
   *          'playToggle',
   *          'progressControl',
   *          'volumePanel',
   *          'fullscreenToggle',
   *          'chromecastButton',
   *       ],
   *    }
   * });
   * ```
   *
   * @param videojs {object} A reference to {@link http://docs.videojs.com/module-videojs.html|Video.js}
   * @see http://docs.videojs.com/module-videojs.html#~registerPlugin
   */

  /**
   * The Video.js Button class is the base class for UI button components.
   *
   * @external Button
   * @see {@link http://docs.videojs.com/Button.html|Button}
   */
  const ButtonComponent = videojs.getComponent('Button');

  /**
   * The ChromecastButton module contains both the ChromecastButton class definition and
   * the function used to register the button as a Video.js Component.
   * @module ChromecastButton
   */

  /** @lends ChromecastButton.prototype **/
  class ChromecastButton extends ButtonComponent {
    /**
     * This class is a button component designed to be displayed in the
     * player UI's control bar. It opens the Chromecast menu when clicked.
     *
     * @constructs
     * @extends external:Button
     * @param player {Player} the video.js player instance
     */
    constructor(player, options) {
      super(player, options);
      player.on('chromecastConnected', this._onChromecastConnected.bind(this));
      player.on('chromecastDisconnected', this._onChromecastDisconnected.bind(this));
      player.on('chromecastDevicesAvailable', this._onChromecastDevicesAvailable.bind(this));
      player.on('chromecastDevicesUnavailable', this._onChromecastDevicesUnavailable.bind(this));

      // Use the initial state of `hasAvailableDevices` to call the corresponding event
      // handlers because the corresponding events may have already been emitted before
      // binding the listeners above.
      if (player.chromecastSessionManager && player.chromecastSessionManager.hasAvailableDevices()) {
        this._onChromecastDevicesAvailable();
      } else {
        this._onChromecastDevicesUnavailable();
      }
      if (options.addCastLabelToButton) {
        this.el().classList.add('vjs-chromecast-button-lg');
        this._labelEl = document.createElement('span');
        this._labelEl.classList.add('vjs-chromecast-button-label');
        this._updateCastLabelText();
        this.el().appendChild(this._labelEl);
      } else {
        this.controlText('Open Chromecast menu');
      }
    }

    /**
     * Overrides Button#buildCSSClass to return the classes used on the button element.
     *
     * @param el {DOMElement}
     * @see {@link http://docs.videojs.com/Button.html#buildCSSClass|Button#buildCSSClass}
     */
    buildCSSClass() {
      return 'vjs-chromecast-button ' + (this._isChromecastConnected ? 'vjs-chromecast-casting-state ' : '') + (this.options_.addCastLabelToButton ? 'vjs-chromecast-button-lg ' : '') + ButtonComponent.prototype.buildCSSClass();
    }

    /**
     * Overrides Button#handleClick to handle button click events. Chromecast
     * functionality is handled outside of this class, which should be limited
     * to UI related logic.  This function simply triggers an event on the player.
     *
     * @fires ChromecastButton#chromecastRequested
     * @param el {DOMElement}
     * @see {@link http://docs.videojs.com/Button.html#handleClick|Button#handleClick}
     */
    handleClick() {
      this.player().trigger('chromecastRequested');
    }

    /**
     * Handles `chromecastConnected` player events.
     *
     * @private
     */
    _onChromecastConnected() {
      this._isChromecastConnected = true;
      this._reloadCSSClasses();
      this._updateCastLabelText();
    }

    /**
     * Handles `chromecastDisconnected` player events.
     *
     * @private
     */
    _onChromecastDisconnected() {
      this._isChromecastConnected = false;
      this._reloadCSSClasses();
      this._updateCastLabelText();
    }

    /**
     * Handles `chromecastDevicesAvailable` player events.
     *
     * @private
     */
    _onChromecastDevicesAvailable() {
      this.show();
    }

    /**
     * Handles `chromecastDevicesUnavailable` player events.
     *
     * @private
     */
    _onChromecastDevicesUnavailable() {
      this.hide();
    }

    /**
     * Re-calculates which CSS classes the button needs and sets them on the buttons'
     * DOMElement.
     *
     * @private
     */
    _reloadCSSClasses() {
      if (!this.el_) {
        return;
      }
      this.el_.className = this.buildCSSClass();
    }

    /**
     * Updates the optional cast label text based on whether the chromecast is connected
     * or disconnected.
     *
     * @private
     */
    _updateCastLabelText() {
      if (!this._labelEl) {
        return;
      }
      this._labelEl.textContent = this._isChromecastConnected ? this.localize('Disconnect Cast') : this.localize('Cast');
    }
  }
  videojs.registerComponent('chromecastButton', ChromecastButton);
};

},{}],4:[function(require,module,exports){
"use strict";

/**
 * @module enableChromecast
 */

var ChromecastSessionManager = require('./chromecast/ChromecastSessionManager'),
  CHECK_AVAILABILITY_INTERVAL = 1000,
  // milliseconds
  CHECK_AVAILABILITY_TIMEOUT = 30 * 1000; // milliseconds

/**
 * Configures the Chromecast
 * [casting context](https://developers.google.com/cast/docs/reference/chrome/cast.framework.CastContext),
 * which is required before casting.
 *
 * @private
 * @param options {object} the plugin options
 */
function configureCastContext(options) {
  var context = cast.framework.CastContext.getInstance();
  context.setOptions({
    receiverApplicationId: options.receiverAppID || chrome.cast.media.DEFAULT_MEDIA_RECEIVER_APP_ID,
    // Setting autoJoinPolicy to ORIGIN_SCOPED prevents this plugin from automatically
    // trying to connect to a preexisting Chromecast session, if one exists. The user
    // must end any existing session before trying to cast from this player instance.
    autoJoinPolicy: chrome.cast.AutoJoinPolicy.ORIGIN_SCOPED
  });
}

/**
 * Handles the `chromecastRequested` event. Delegates to a `chromecastSessionManager`
 * instance.
 *
 * @private
 * @param player {object} a Video.js player instance
 */
function onChromecastRequested(player) {
  player.chromecastSessionManager.openCastMenu();
}

/**
 * Adds the Chromecast button to the player's control bar, if one does not already exist,
 * then starts listening for the `chromecastRequested` event.
 *
 * @private
 * @param player {object} a Video.js player instance
 * @param options {object} the plugin options
 */
function setUpChromecastButton(player, options) {
  var indexOpt;

  // Ensure Chromecast button exists
  if (options.addButtonToControlBar && !player.controlBar.getChild('chromecastButton')) {
    // Figure out Chromecast button's index
    indexOpt = player.controlBar.children().length;
    if (typeof options.buttonPositionIndex !== 'undefined') {
      indexOpt = options.buttonPositionIndex >= 0 ? options.buttonPositionIndex : player.controlBar.children().length + options.buttonPositionIndex;
    }
    player.controlBar.addChild('chromecastButton', options, indexOpt);
  }
  // Respond to requests for casting. The ChromecastButton component triggers this event
  // when the user clicks the Chromecast button.
  player.on('chromecastRequested', onChromecastRequested.bind(null, player));
}

/**
 * Creates a {@link ChromecastSessionManager} and assigns it to the player.
 *
 * @private
 * @param player {object} a Video.js player instance
 */
function createSessionManager(player) {
  if (!player.chromecastSessionManager) {
    player.chromecastSessionManager = new ChromecastSessionManager(player);
  }
}

/**
 * Sets up and configures the casting context and Chromecast button.
 *
 * @private
 * @param options {object} the plugin options
 */
function enableChromecast(player, options) {
  configureCastContext(options);
  createSessionManager(player);
  setUpChromecastButton(player, options);
}

/**
 * Waits for the Chromecast APIs to become available, then configures the casting context
 * and configures the Chromecast button. The Chromecast APIs are loaded asynchronously,
 * so we must wait until they are available before initializing the casting context and
 * Chromecast button.
 *
 * @private
 * @param player {object} a Video.js player instance
 * @param options {object} the plugin options
 */
function waitUntilChromecastAPIsAreAvailable(player, options) {
  var maxTries = CHECK_AVAILABILITY_TIMEOUT / CHECK_AVAILABILITY_INTERVAL,
    tries = 1,
    intervalID;

  // The Chromecast APIs are loaded asynchronously, so they may not be loaded and
  // initialized at this point. The Chromecast APIs do provide a callback function that
  // is called after the framework has loaded, but it requires you to define the callback
  // function **before** loading the APIs. That would require us to expose some callback
  // function to `window` here, and would require users of this plugin to define a
  // Chromecast API callback on `window` that calls our callback function in their HTML
  // file. To avoid all of this, we simply check to see if the Chromecast API is
  // available periodically, and stop after a timeout threshold has passed.
  //
  // See https://developers.google.com/cast/docs/chrome_sender_integrate#initialization
  intervalID = setInterval(function () {
    if (tries > maxTries) {
      clearInterval(intervalID);
      return;
    }
    if (ChromecastSessionManager.isChromecastAPIAvailable()) {
      clearInterval(intervalID);
      enableChromecast(player, options);
    }
    tries = tries + 1;
  }, CHECK_AVAILABILITY_INTERVAL);
}

/**
 * Registers the Chromecast plugin with Video.js. Calls
 * [videojs#registerPlugin](http://docs.videojs.com/module-videojs.html#~registerPlugin),
 * which will add a plugin function called `chromecast` to any instance of a Video.js
 * player that is created after calling this function. Call `player.chromecast(options)`,
 * passing in configuration options, to enable the Chromecast plugin on your Player
 * instance.
 *
 * Currently, there are only two configuration options:
 *
 *    * **`receiverAppID`** - the string ID of a [Chromecast receiver
 *      app](https://developers.google.com/cast/docs/receiver_apps) to use. Defaults to
 *      the [default Media Receiver
 *      ID](https://developers.google.com/cast/docs/receiver_apps#default).
 *    * **`addButtonToControlBar`** - flag that tells the plugin
 *      whether or not it should automatically add the Chromecast button the the Video.js
 *      player's control bar component. Defaults to `true`.
 *
 * Other configuration options are set through the player's Chromecast Tech configuration:
 *
 * ```
 * var playerOptions, player, pluginOptions;
 *
 * playerOptions = {
 *    chromecast: {
 *       requestTitleFn: function(source) {
 *          return titles[source.url];
 *       },
 *       requestSubtitleFn: function(source) {
 *          return subtitles[source.url];
 *       },
 *       requestCustomDataFn: function(source) {
 *          return customData[source.url];
 *       }
 *    }
 * };
 *
 * pluginOptions = {
 *    receiverAppID: '1234',
 *    addButtonToControlBar: false,
 * };
 *
 * player = videojs(document.getElementById('myVideoElement'), playerOptions);
 * player.chromecast(pluginOptions); // initializes the Chromecast plugin
 * ```
 *
 * @param {object} videojs
 * @see http://docs.videojs.com/module-videojs.html#~registerPlugin
 */
module.exports = function (videojs) {
  videojs.registerPlugin('chromecast', function (options) {
    var pluginOptions = Object.assign({
      addButtonToControlBar: true
    }, options || {});

    // `this` is an instance of a Video.js Player.
    // Wait until the player is "ready" so that the player's control bar component has
    // been created.
    this.ready(function () {
      if (!this.controlBar) {
        return;
      }
      if (ChromecastSessionManager.isChromecastAPIAvailable()) {
        enableChromecast(this, pluginOptions);
      } else {
        waitUntilChromecastAPIsAreAvailable(this, pluginOptions);
      }
    }.bind(this));
  });
};

},{"./chromecast/ChromecastSessionManager":2}],5:[function(require,module,exports){
"use strict";

/* eslint-disable global-require */
var preloadWebComponents = require('./preloadWebComponents'),
  createChromecastButton = require('./components/ChromecastButton'),
  createChromecastTech = require('./tech/ChromecastTech'),
  enableChromecast = require('./enableChromecast');

/**
 * @module index
 */

/**
 * Registers the Chromecast plugin and ChromecastButton Component with Video.js. See
 * {@link module:ChromecastButton} and {@link module:enableChromecast} for more details
 * about how the plugin and button are registered and configured.
 *
 * @param videojs {object} the videojs library. If `undefined`, this plugin
 * will look to `window.videojs`.
 * @param userOpts {object} the options to use for configuration
 * @see module:enableChromecast
 * @see module:ChromecastButton
 */
module.exports = function (videojs, userOpts) {
  var options = Object.assign({
    preloadWebComponents: false
  }, userOpts);
  if (options.preloadWebComponents) {
    preloadWebComponents();
  }
  videojs = videojs || window.videojs;
  createChromecastButton(videojs);
  createChromecastTech(videojs);
  enableChromecast(videojs);
};

},{"./components/ChromecastButton":3,"./enableChromecast":4,"./preloadWebComponents":6,"./tech/ChromecastTech":8}],6:[function(require,module,exports){
"use strict";

function doesUserAgentContainString(str) {
  return typeof window.navigator.userAgent === 'string' && window.navigator.userAgent.indexOf(str) >= 0;
}

// For information as to why this is needed, please see:
// https://github.com/silvermine/videojs-chromecast/issues/17
// https://github.com/silvermine/videojs-chromecast/issues/22

module.exports = function () {
  var needsWebComponents = !document.registerElement,
    iosChrome = doesUserAgentContainString('CriOS'),
    androidChrome;
  androidChrome = doesUserAgentContainString('Android') && doesUserAgentContainString('Chrome/') && window.navigator.presentation;

  // These checks are based on the checks found in `cast_sender.js` which
  // determine if `cast_framework.js` needs to be loaded
  if ((androidChrome || iosChrome) && needsWebComponents) {
    // This is requiring webcomponents.js@0.7.24 because that's what was used
    // by the Chromecast framework at the time this was added.
    // We are using webcomponents-lite.js because it doesn't interfere with jQuery as
    // badly (e.g. it doesn't interfere with jQuery's fix for consistently bubbling
    // events, see #21). While the "lite" version does not include the shadow DOM
    // polyfills that the Chromecast framework may need for the <google-cast-button>
    // component to work properly, this plugin does not use the <google-cast-button>
    // component.
    require('webcomponents.js/webcomponents-lite.js'); // eslint-disable-line global-require
  }
};

},{"webcomponents.js/webcomponents-lite.js":1}],7:[function(require,module,exports){
"use strict";

// This file is used to create a standalone javascript file for use in a script tag. The
// file that is output assumes that Video.js is available at `window.videojs`.

require('./index')(undefined, window.SILVERMINE_VIDEOJS_CHROMECAST_CONFIG);

},{"./index":5}],8:[function(require,module,exports){
"use strict";

var ChromecastSessionManager = require('../chromecast/ChromecastSessionManager'),
  ChromecastTechUI = require('./ChromecastTechUI');

/**
 * Registers the ChromecastTech Tech with Video.js. Calls {@link
 * http://docs.videojs.com/Tech.html#.registerTech}, which will add a Tech called
 * `chromecast` to the list of globally registered Video.js Tech implementations.
 *
 * [Video.js Tech](http://docs.videojs.com/Tech.html) are initialized and used
 * automatically by Video.js Player instances. Whenever a new source is set on the player,
 * the player iterates through the list of available Tech to determine which to use to
 * play the source.
 *
 * @param videojs {object} A reference to
 * {@link http://docs.videojs.com/module-videojs.html|Video.js}
 * @see http://docs.videojs.com/Tech.html#.registerTech
 */
module.exports = function (videojs) {
  var Tech = videojs.getComponent('Tech'),
    SESSION_TIMEOUT = 10 * 1000; // milliseconds

  /**
   * @module ChomecastTech
   */

  /**
   * The Video.js Tech class is the base class for classes that provide media playback
   * technology implementations to Video.js such as HTML5, Flash and HLS.
   *
   * @external Tech
   * @see {@link http://docs.videojs.com/Tech.html|Tech}
   */

  /** @lends ChromecastTech.prototype */
  class ChromecastTech extends Tech {
    /**
     * Implements Video.js playback {@link http://docs.videojs.com/tutorial-tech_.html|Tech}
     * for {@link https://developers.google.com/cast/|Google's Chromecast}.
     *
     * @constructs ChromecastTech
     * @extends external:Tech
     * @param options {object} The options to use for configuration
     * @see {@link https://developers.google.com/cast/|Google Cast}
     */
    constructor(options) {
      super(options);
      this.featuresVolumeControl = true;
      this.featuresPlaybackRate = false;
      this.movingMediaElementInDOM = false;
      this.featuresFullscreenResize = true;
      this.featuresTimeupdateEvents = true;
      this.featuresProgressEvents = false;
      // Text tracks are not supported in this version
      this.featuresNativeTextTracks = false;
      this.featuresNativeAudioTracks = false;
      this.featuresNativeVideoTracks = false;

      // Give ChromecastTech class instances a reference to videojs
      this.videojs = videojs;
      this._eventListeners = [];
      this.videojsPlayer = this.videojs(options.playerId);
      this._chromecastSessionManager = this.videojsPlayer.chromecastSessionManager;
      this._ui.updatePoster(this.videojsPlayer.poster());
      this._remotePlayer = this._chromecastSessionManager.getRemotePlayer();
      this._remotePlayerController = this._chromecastSessionManager.getRemotePlayerController();
      this._listenToPlayerControllerEvents();
      this.on('dispose', this._removeAllEventListeners.bind(this));
      this._hasPlayedAnyItem = false;
      this._requestTitle = options.requestTitleFn || function () {/* noop */};
      this._requestSubtitle = options.requestSubtitleFn || function () {/* noop */};
      this._requestCustomData = options.requestCustomDataFn || function () {/* noop */};
      this._modifyLoadRequestFn = options.modifyLoadRequestFn || function (request) {
        return request;
      };
      // See `currentTime` function
      this._initialStartTime = options.startTime || 0;
      this._playSource(options.source, this._initialStartTime);
      this.ready(function () {
        this.setMuted(options.muted);
      }.bind(this));
    }

    /**
     * Creates a DOMElement that Video.js displays in its player UI while this Tech is
     * active.
     *
     * @returns {DOMElement}
     * @see {@link http://docs.videojs.com/Tech.html#createEl}
     */
    createEl() {
      // We have to initialize the UI here, because the super.constructor
      // calls `createEl`, which references `this._ui`.
      this._ui = this._ui || new ChromecastTechUI();
      return this._ui.getDOMElement();
    }

    /**
     * Resumes playback if a media item is paused or restarts an item from
     * its beginning if the item has played and ended.
     *
     * @see {@link http://docs.videojs.com/Player.html#play}
     */
    play() {
      if (!this.paused()) {
        return;
      }
      if (this.ended() && !this._isMediaLoading) {
        // Restart the current item from the beginning
        this._playSource({
          src: this.videojsPlayer.src()
        }, 0);
      } else {
        this._remotePlayerController.playOrPause();
      }
    }

    /**
     * Pauses playback if the player is not already paused and if the current media item
     * has not ended yet.
     *
     * @see {@link http://docs.videojs.com/Player.html#pause}
     */
    pause() {
      if (!this.paused() && this._remotePlayer.canPause) {
        this._remotePlayerController.playOrPause();
      }
    }

    /**
     * Returns whether or not the player is "paused". Video.js'
     * definition of "paused" is "playback paused" OR "not playing".
     *
     * @returns {boolean} true if playback is paused
     * @see {@link http://docs.videojs.com/Player.html#paused}
     */
    paused() {
      return this._remotePlayer.isPaused || this.ended() || this._remotePlayer.playerState === null;
    }

    /**
     * Stores the given source and begins playback, starting at the beginning
     * of the media item.
     *
     * @param source {object} the source to store and play
     * @see {@link http://docs.videojs.com/Player.html#src}
     */
    setSource(source) {
      if (this._currentSource && this._currentSource.src === source.src && this._currentSource.type === source.type) {
        // Skip setting the source if the `source` argument is the
        // same as what's already been set. This `setSource` function
        // calls `this._playSource` which sends a "load media" request
        // to the Chromecast PlayerController. Because this function
        // may be called multiple times in rapid succession with the same `source`
        // argument, we need to de-duplicate calls with the same `source` argument to
        // prevent overwhelming the Chromecast PlayerController with expensive "load
        // media" requests, which it itself does not de-duplicate.
        return;
      }
      // We cannot use `this.videojsPlayer.currentSource()` because the
      // value returned by that function is not the same as what's returned
      // by the Video.js Player's middleware after they are run. Also, simply
      // using `this.videojsPlayer.src()` does not include mimetype information
      // which we pass to the Chromecast player.
      this._currentSource = source;
      this._playSource(source, 0);
    }

    /**
     * Plays the given source, beginning at an optional starting time.
     *
     * @private
     * @param source {object} the source to play
     * @param [startTime] The time to start playback at, in seconds
     * @see {@link http://docs.videojs.com/Player.html#src}
     */
    _playSource(source, startTime) {
      var castSession = this._getCastSession(),
        mediaInfo = new chrome.cast.media.MediaInfo(source.src, source.type),
        title = this._requestTitle(source),
        subtitle = this._requestSubtitle(source),
        poster = this.poster(),
        customData = this._requestCustomData(source),
        request;
      this.trigger('waiting');
      this._clearSessionTimeout();
      mediaInfo.metadata = new chrome.cast.media.GenericMediaMetadata();
      mediaInfo.metadata.metadataType = chrome.cast.media.MetadataType.GENERIC;
      mediaInfo.metadata.title = title;
      mediaInfo.metadata.subtitle = subtitle;
      mediaInfo.streamType = this.videojsPlayer.liveTracker && this.videojsPlayer.liveTracker.isLive() ? chrome.cast.media.StreamType.LIVE : chrome.cast.media.StreamType.BUFFERED;
      if (poster) {
        mediaInfo.metadata.images = [{
          url: poster
        }];
      }
      if (customData) {
        mediaInfo.customData = customData;
      }
      this._ui.updateTitle(title);
      this._ui.updateSubtitle(subtitle);
      request = new chrome.cast.media.LoadRequest(mediaInfo);
      request.autoplay = true;
      request.currentTime = startTime;
      request = this._modifyLoadRequestFn(request);
      this._isMediaLoading = true;
      this._hasPlayedCurrentItem = false;
      castSession.loadMedia(request).then(function () {
        this._clearSessionTimeout();
        if (!this._hasPlayedAnyItem) {
          // `triggerReady` is required here to notify the Video.js
          // player that the Tech has been initialized and is ready.
          this.triggerReady();
        }
        this.trigger('loadstart');
        this.trigger('loadeddata');
        this.trigger('play');
        this.trigger('playing');
        this._hasPlayedAnyItem = true;
        this._isMediaLoading = false;
        this._getMediaSession().addUpdateListener(this._onMediaSessionStatusChanged.bind(this));
      }.bind(this), this._triggerErrorEvent.bind(this));
    }

    /**
     * Manually updates the current time. The playback position will jump to
     * the given time and continue playing if the item was playing when `setCurrentTime`
     * was called, or remain paused if the item was paused.
     *
     * @param time {number} the playback time position to jump to
     * @see {@link http://docs.videojs.com/Tech.html#setCurrentTime}
     */
    setCurrentTime(time) {
      var duration = this.duration();
      if (time > duration || !this._remotePlayer.canSeek) {
        return;
      }
      // Seeking to any place within (approximately) 1 second of the end of the item
      // causes the Video.js player to get stuck in a BUFFERING state. To work around
      // this, we only allow seeking to within 1 second of the end of an item.
      this._remotePlayer.currentTime = Math.min(duration - 1, time);
      this._remotePlayerController.seek();
      this._triggerTimeUpdateEvent();
    }

    /**
     * Returns the current playback time position.
     *
     * @returns {number} the current playback time position
     * @see {@link http://docs.videojs.com/Player.html#currentTime}
     */
    currentTime() {
      // There is a brief period of time when Video.js has switched to the chromecast
      // Tech, but chromecast has not yet loaded its first media item. During
      // that time, Video.js calls this `currentTime` function to update
      // its player UI. In that period, `this._remotePlayer.currentTime`
      // will be 0 because the media has not loaded yet. To prevent the
      // UI from using a 0 second currentTime, we use the currentTime passed
      // in to the first media item that was provided to the Tech until
      // chromecast plays its first item.
      if (!this._hasPlayedAnyItem) {
        return this._initialStartTime;
      }
      return this._remotePlayer.currentTime;
    }

    /**
     * Returns the duration of the current media item, or `0` if the source
     * is not set or if the duration of the item is not available from the
     * Chromecast API yet.
     *
     * @returns {number} the duration of the current media item
     * @see {@link http://docs.videojs.com/Player.html#duration}
     */
    duration() {
      // There is a brief period of time when Video.js has switched to the chromecast
      // Tech, but chromecast has not yet loaded its first media item.
      // During that time, Video.js calls this `duration` function to update its player
      // UI. In that period, `this._remotePlayer.duration` will be 0 because the media
      // has not loaded yet. To prevent the UI from using a 0 second duration, we
      // use the duration passed in to the first media item that was provided to
      // the Tech until chromecast plays its first item.
      if (!this._hasPlayedAnyItem) {
        return this.videojsPlayer.duration();
      }
      return this._remotePlayer.duration;
    }

    /**
     * Returns whether or not the current media item has finished playing.
     * Returns `false` if a media item has not been loaded, has not been played,
     * or has not yet finished playing.
     *
     * @returns {boolean} true if the current media item has finished playing
     * @see {@link http://docs.videojs.com/Player.html#ended}
     */
    ended() {
      var mediaSession = this._getMediaSession();
      if (!mediaSession && this._hasMediaSessionEnded) {
        return true;
      }
      return mediaSession ? mediaSession.idleReason === chrome.cast.media.IdleReason.FINISHED : false;
    }

    /**
     * Returns the current volume level setting as a decimal number between `0` and `1`.
     *
     * @returns {number} the current volume level
     * @see {@link http://docs.videojs.com/Player.html#volume}
     */
    volume() {
      return this._remotePlayer.volumeLevel;
    }

    /**
     * Sets the current volume level. Volume level is a decimal number
     * between `0` and `1`, where `0` is muted and `1` is the loudest volume level.
     *
     * @param volumeLevel {number}
     * @returns {number} the current volume level
     * @see {@link http://docs.videojs.com/Player.html#volume}
     */
    setVolume(volumeLevel) {
      this._remotePlayer.volumeLevel = volumeLevel;
      this._remotePlayerController.setVolumeLevel();
      // This event is triggered by the listener on
      // `RemotePlayerEventType.VOLUME_LEVEL_CHANGED`, but waiting for
      // that event to fire in response to calls to `setVolume` introduces
      // noticeable lag in the updating of the player UI's volume slider bar,
      // which makes user interaction with the volume slider choppy.
      this._triggerVolumeChangeEvent();
    }

    /**
     * Returns whether or not the player is currently muted.
     *
     * @returns {boolean} true if the player is currently muted
     * @see {@link http://docs.videojs.com/Player.html#muted}
     */
    muted() {
      return this._remotePlayer.isMuted;
    }

    /**
     * Mutes or un-mutes the player. Does nothing if the player is currently
     * muted and the `isMuted` parameter is true or if the player is not muted and
     * `isMuted` is false.
     *
     * @param isMuted {boolean} whether or not the player should be muted
     * @see {@link http://docs.videojs.com/Html5.html#setMuted} for an example
     */
    setMuted(isMuted) {
      if (this._remotePlayer.isMuted && !isMuted || !this._remotePlayer.isMuted && isMuted) {
        this._remotePlayerController.muteOrUnmute();
      }
    }

    /**
     * Gets the URL to the current poster image.
     *
     * @returns {string} URL to the current poster image or `undefined` if none exists
     * @see {@link http://docs.videojs.com/Player.html#poster}
     */
    poster() {
      return this._ui.getPoster();
    }

    /**
     * Sets the URL to the current poster image. The poster image shown
     * in the Chromecast Tech UI view is updated with this new URL.
     *
     * @param poster {string} the URL to the new poster image
     * @see {@link http://docs.videojs.com/Tech.html#setPoster}
     */
    setPoster(poster) {
      this._ui.updatePoster(poster);
    }

    /**
     * This function is "required" when implementing {@link external:Tech}
     * and is supposed to return a mock
     * {@link https://developer.mozilla.org/en-US/docs/Web/API/TimeRanges|TimeRanges}
     * object that represents the portions of the current media item that have been
     * buffered. However, the Chromecast API does not currently provide a way
     * to determine how much the media item has buffered, so we always
     * return `undefined`.
     *
     * Returning `undefined` is safe: the player will simply not display
     * the buffer amount indicator in the scrubber UI.
     *
     * @returns {undefined} always returns `undefined`
     * @see {@link http://docs.videojs.com/Player.html#buffered}
     */
    buffered() {
      return undefined;
    }

    /**
     * This function is "required" when implementing {@link external:Tech}
     * and is supposed to return a mock
     * {@link https://developer.mozilla.org/en-US/docs/Web/API/TimeRanges|TimeRanges}
     * object that represents the portions of the current media item that has playable
     * content. However, the Chromecast API does not currently provide a
     * way to determine how much the media item has playable content, so
     * we'll just assume the entire video is an available seek target.
     *
     * The risk here lies with live streaming, where there may exist a sliding window of
     * playable content and seeking is only possible within the last X number of
     * minutes, rather than for the entire video.
     *
     * Unfortunately we have no way of detecting when this is the case. Returning
     * anything other than the full range of the video means that we lose the ability
     * to seek during VOD.
     *
     * @returns {TimeRanges} always returns a `TimeRanges` object with one
     * `TimeRange` that starts at `0` and ends at the `duration` of the
     * current media item
     * @see {@link http://docs.videojs.com/Player.html#seekable}
     */
    seekable() {
      // TODO Investigate if there's a way to detect
      // if the source is live, so that we can
      // possibly adjust the seekable `TimeRanges` accordingly.
      return this.videojs.createTimeRange(0, this.duration());
    }

    /**
     * Returns whether the native media controls should be shown (`true`) or hidden
     * (`false`). Not applicable to this Tech.
     *
     * @returns {boolean} always returns `false`
     * @see {@link http://docs.videojs.com/Html5.html#controls} for an example
     */
    controls() {
      return false;
    }

    /**
     * Returns whether or not the browser should show the player
     * "inline" (non-fullscreen) by default. This function always
     * returns true to tell the browser that non-fullscreen playback is preferred.
     *
     * @returns {boolean} always returns `true`
     * @see {@link http://docs.videojs.com/Html5.html#playsinline} for an example
     */
    playsinline() {
      return true;
    }

    /**
     * Returns whether or not fullscreen is supported by this Tech.
     * Always returns `true` because fullscreen is always supported.
     *
     * @returns {boolean} always returns `true`
     * @see {@link http://docs.videojs.com/Html5.html#supportsFullScreen} for an example
     */
    supportsFullScreen() {
      return true;
    }

    /**
     * Sets a flag that determines whether or not the media should automatically begin
     * playing on page load. This is not supported because a Chromecast session must be
     * initiated by casting via the casting menu and cannot autoplay.
     *
     * @see {@link http://docs.videojs.com/Html5.html#setAutoplay} for an example
     */
    setAutoplay() {
      // Not supported
    }

    /**
     * @returns {number} the chromecast player's playback rate, if available. Otherwise,
     * the return value defaults to `1`.
     */
    playbackRate() {
      var mediaSession = this._getMediaSession();
      return mediaSession ? mediaSession.playbackRate : 1;
    }

    /**
     * Does nothing. Changing the playback rate is not supported.
     */
    setPlaybackRate() {
      // Not supported
    }

    /**
     * Does nothing. Satisfies calls to the missing preload method.
     */
    preload() {
      // Not supported
    }

    /**
     * Causes the Tech to begin loading the current source. `load`
     * is not supported in this ChromecastTech because setting the
     * source on the `Chromecast` automatically causes it to begin loading.
     */
    load() {
      // Not supported
    }

    /**
     * Gets the Chromecast equivalent of HTML5 Media Element's `readyState`.
     *
     * @see https://developer.mozilla.org/en-US/docs/Web/API/HTMLMediaElement/readyState
     */
    readyState() {
      if (this._remotePlayer.playerState === 'IDLE' || this._remotePlayer.playerState === 'BUFFERING') {
        return 0; // HAVE_NOTHING
      }
      return 4;
    }

    /**
     * Wires up event listeners for
     * [RemotePlayerController](https://developers.google.com/cast/docs/reference/chrome/cast.framework.RemotePlayerController)
     * events.
     *
     * @private
     */
    _listenToPlayerControllerEvents() {
      var eventTypes = cast.framework.RemotePlayerEventType;
      this._addEventListener(this._remotePlayerController, eventTypes.PLAYER_STATE_CHANGED, this._onPlayerStateChanged, this);
      this._addEventListener(this._remotePlayerController, eventTypes.VOLUME_LEVEL_CHANGED, this._triggerVolumeChangeEvent, this);
      this._addEventListener(this._remotePlayerController, eventTypes.IS_MUTED_CHANGED, this._triggerVolumeChangeEvent, this);
      this._addEventListener(this._remotePlayerController, eventTypes.CURRENT_TIME_CHANGED, this._triggerTimeUpdateEvent, this);
      this._addEventListener(this._remotePlayerController, eventTypes.DURATION_CHANGED, this._triggerDurationChangeEvent, this);
    }

    /**
     * Registers an event listener on the given target object.
     * Because many objects in the Chromecast API are either singletons
     * or must be shared between instances of `ChromecastTech` for the
     * lifetime of the player, we must unbind the listeners when this Tech
     * instance is destroyed to prevent memory leaks. To do that, we need to keep
     * a reference to listeners that are added to global objects so that we can
     * use those references to remove the listener when this Tech is destroyed.
     *
     * @param target {object} the object to register the event listener on
     * @param type {string} the name of the event
     * @param callback {Function} the listener's callback function that
     * executes when the event is emitted
     * @param context {object} the `this` context to use when executing the `callback`
     * @private
     */
    _addEventListener(target, type, callback, context) {
      var listener;
      listener = {
        target: target,
        type: type,
        callback: callback,
        context: context,
        listener: callback.bind(context)
      };
      target.addEventListener(type, listener.listener);
      this._eventListeners.push(listener);
    }

    /**
     * Removes all event listeners that were registered with global objects during the
     * lifetime of this Tech. See {@link _addEventListener} for more information
     * about why this is necessary.
     *
     * @private
     */
    _removeAllEventListeners() {
      while (this._eventListeners.length > 0) {
        this._removeEventListener(this._eventListeners[0]);
      }
      this._eventListeners = [];
    }

    /**
     * Removes a single event listener that was registered with global objects
     * during the lifetime of this Tech. See {@link _addEventListener} for
     * more information about why this is necessary.
     *
     * @private
     */
    _removeEventListener(listener) {
      var index = -1,
        pass = false,
        i;
      listener.target.removeEventListener(listener.type, listener.listener);
      for (i = 0; i < this._eventListeners.length; i++) {
        pass = this._eventListeners[i].target === listener.target && this._eventListeners[i].type === listener.type && this._eventListeners[i].callback === listener.callback && this._eventListeners[i].context === listener.context;
        if (pass) {
          index = i;
          break;
        }
      }
      if (index !== -1) {
        this._eventListeners.splice(index, 1);
      }
    }

    /**
     * Handles Chromecast player state change events. The player may "change state" when
     * paused, played, buffering, etc.
     *
     * @private
     */
    _onPlayerStateChanged() {
      var states = chrome.cast.media.PlayerState,
        playerState = this._remotePlayer.playerState;
      if (playerState === states.PLAYING) {
        this._hasPlayedCurrentItem = true;
        this.trigger('play');
        this.trigger('playing');
      } else if (playerState === states.PAUSED) {
        this.trigger('pause');
      } else if (playerState === states.IDLE && this.ended() || playerState === null && this._hasPlayedCurrentItem) {
        this._hasPlayedCurrentItem = false;
        this._closeSessionOnTimeout();
        this.trigger('ended');
        this._triggerTimeUpdateEvent();
      } else if (playerState === states.BUFFERING) {
        this.trigger('waiting');
      }
    }

    /**
     * Handles Chromecast MediaSession state change events. The only property sent
     * to this event is whether the session is alive. This is useful for determining
     * if an item has ended as the MediaSession will fire this event with `false` then
     * be immediately destroyed. This means that we cannot trust `idleReason` to show
     * whether an item has ended since we may no longer have access to the MediaSession.
     *
     * @private
     */
    _onMediaSessionStatusChanged(isAlive) {
      this._hasMediaSessionEnded = !!isAlive;
    }

    /**
     * Ends the session after a certain number of seconds of inactivity.
     *
     * If the Chromecast player is in the "IDLE" state after an item has ended, and no
     * further items are queued up to play, the session is considered inactive. Once a
     * period of time (currently 10 seconds) has elapsed with no activity, we manually
     * end the session to prevent long periods of a blank Chromecast screen that is
     * shown at the end of item playback.
     *
     * @private
     */
    _closeSessionOnTimeout() {
      // Ensure that there's never more than one session timeout active
      this._clearSessionTimeout();
      this._sessionTimeoutID = setTimeout(function () {
        var castSession = this._getCastSession();
        if (castSession) {
          castSession.endSession(true);
        }
        this._clearSessionTimeout();
      }.bind(this), SESSION_TIMEOUT);
    }

    /**
     * Stops the timeout that is waiting during a period of inactivity in order to close
     * the session.
     *
     * @private
     * @see _closeSessionOnTimeout
     */
    _clearSessionTimeout() {
      if (this._sessionTimeoutID) {
        clearTimeout(this._sessionTimeoutID);
        this._sessionTimeoutID = false;
      }
    }

    /**
     * @private
     * @return {object} the current CastContext, if one exists
     */
    _getCastContext() {
      return this._chromecastSessionManager.getCastContext();
    }

    /**
     * @private
     * @return {object} the current CastSession, if one exists
     */
    _getCastSession() {
      return this._getCastContext().getCurrentSession();
    }

    /**
     * @private
     * @return {object} the current MediaSession, if one exists
     * @see https://developers.google.com/cast/docs/reference/chrome/chrome.cast.media.Media
     */
    _getMediaSession() {
      var castSession = this._getCastSession();
      return castSession ? castSession.getMediaSession() : null;
    }

    /**
     * Triggers a 'volumechange' event
     * @private
     * @see http://docs.videojs.com/Player.html#event:volumechange
     */
    _triggerVolumeChangeEvent() {
      this.trigger('volumechange');
    }

    /**
     * Triggers a 'timeupdate' event
     * @private
     * @see http://docs.videojs.com/Player.html#event:timeupdate
     */
    _triggerTimeUpdateEvent() {
      this.trigger('timeupdate');
    }

    /**
     * Triggers a 'durationchange' event
     * @private
     * @see http://docs.videojs.com/Player.html#event:durationchange
     */
    _triggerDurationChangeEvent() {
      this.trigger('durationchange');
    }

    /**
     * Triggers an 'error' event
     * @private
     * @see http://docs.videojs.com/Player.html#event:error
     */
    _triggerErrorEvent() {
      this.trigger('error');
    }
  }

  // Required for Video.js Tech implementations.
  // TODO Consider a more comprehensive check based on mimetype.
  ChromecastTech.canPlaySource = () => {
    return ChromecastSessionManager.isChromecastConnected();
  };
  ChromecastTech.isSupported = () => {
    return ChromecastSessionManager.isChromecastConnected();
  };
  videojs.registerTech('chromecast', ChromecastTech);
};

},{"../chromecast/ChromecastSessionManager":2,"./ChromecastTechUI":9}],9:[function(require,module,exports){
"use strict";

/**
 * This class represents the UI that is shown in the player while the Chromecast Tech is
 * active. The UI has a single root DOM element that displays the poster image of the
 * current item and title and subtitle. This class receives updates to the poster, title
 * and subtitle when the media item that the player is playing changes.
 *
 * @class ChromecastTechUI
 */
class ChromecastTechUI {
  constructor() {
    this._el = this._createDOMElement();
  }

  /**
   * Creates and returns a single DOMElement that contains the UI. This implementation
   * of the Chromecast Tech's UI displays a poster image, a title and a subtitle.
   *
   * @private
   * @returns {DOMElement}
   */
  _createDOMElement() {
    var el = this._createElement('div', 'vjs-tech vjs-tech-chromecast'),
      posterContainerEl = this._createElement('div', 'vjs-tech-chromecast-poster'),
      posterImageEl = this._createElement('img', 'vjs-tech-chromecast-poster-img'),
      titleEl = this._createElement('div', 'vjs-tech-chromecast-title'),
      subtitleEl = this._createElement('div', 'vjs-tech-chromecast-subtitle'),
      titleContainer = this._createElement('div', 'vjs-tech-chromecast-title-container');
    posterContainerEl.appendChild(posterImageEl);
    titleContainer.appendChild(titleEl);
    titleContainer.appendChild(subtitleEl);
    el.appendChild(titleContainer);
    el.appendChild(posterContainerEl);
    return el;
  }

  /**
   * A helper method for creating DOMElements of the given type and with the given class
   * name(s).
   *
   * @param type {string} the kind of DOMElement to create (ex: 'div')
   * @param className {string} the class name(s) to give to the DOMElement. May also be
   * a space-delimited list of class names.
   * @returns {DOMElement}
   */
  _createElement(type, className) {
    var el = document.createElement(type);
    el.className = className;
    return el;
  }

  /**
   * Gets the root DOMElement to be shown in the player's UI.
   *
   * @returns {DOMElement}
   */
  getDOMElement() {
    return this._el;
  }

  /**
   * Finds the poster's DOMElement in the root UI element.
   *
   * @private
   * @returns {DOMElement}
   */
  _findPosterEl() {
    return this._el.querySelector('.vjs-tech-chromecast-poster');
  }

  /**
   * Finds the poster's <img> DOMElement in the root UI element.
   *
   * @private
   * @returns {DOMElement}
   */
  _findPosterImageEl() {
    return this._el.querySelector('.vjs-tech-chromecast-poster-img');
  }

  /**
   * Finds the title's DOMElement in the root UI element.
   *
   * @private
   * @returns {DOMElement}
   */
  _findTitleEl() {
    return this._el.querySelector('.vjs-tech-chromecast-title');
  }

  /**
   * Finds the subtitle's DOMElement in the root UI element.
   *
   * @private
   * @returns {DOMElement}
   */
  _findSubtitleEl() {
    return this._el.querySelector('.vjs-tech-chromecast-subtitle');
  }

  /**
   * Sets the current poster image URL and updates the poster image DOMElement with the
   * new poster image URL.
   *
   * @param poster {string} a URL for a poster image
   */
  updatePoster(poster) {
    var posterImageEl = this._findPosterImageEl();
    this._poster = poster ? poster : null;
    if (poster) {
      posterImageEl.setAttribute('src', poster);
      posterImageEl.classList.remove('vjs-tech-chromecast-poster-img-empty');
    } else {
      posterImageEl.removeAttribute('src');
      posterImageEl.classList.add('vjs-tech-chromecast-poster-img-empty');
    }
  }

  /**
   * Gets the current poster image URL.
   *
   * @returns {string} the URL for th current poster image
   */
  getPoster() {
    return this._poster;
  }

  /**
   * Sets the current title and updates the title's DOMElement with the new text.
   *
   * @param title {string} a title to show
   */
  updateTitle(title) {
    var titleEl = this._findTitleEl();
    this._title = title;
    if (title) {
      titleEl.innerHTML = title;
      titleEl.classList.remove('vjs-tech-chromecast-title-empty');
    } else {
      titleEl.classList.add('vjs-tech-chromecast-title-empty');
    }
  }

  /**
   * Sets the current subtitle and updates the subtitle's DOMElement with the new text.
   *
   * @param subtitle {string} a subtitle to show
   */
  updateSubtitle(subtitle) {
    var subtitleEl = this._findSubtitleEl();
    this._subtitle = subtitle;
    if (subtitle) {
      subtitleEl.innerHTML = subtitle;
      subtitleEl.classList.remove('vjs-tech-chromecast-subtitle-empty');
    } else {
      subtitleEl.classList.add('vjs-tech-chromecast-subtitle-empty');
    }
  }
}
module.exports = ChromecastTechUI;

},{}]},{},[7]);
