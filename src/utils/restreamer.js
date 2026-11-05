import { i18n } from '@lingui/core';
import { t } from '@lingui/macro';
import { v4 as uuidv4 } from 'uuid';
import { jwtDecode } from 'jwt-decode';
import Handlebars from 'handlebars/dist/cjs/handlebars';
import SemverSatisfies from 'semver/functions/satisfies';
import SemverGt from 'semver/functions/gt';
import SemverGte from 'semver/functions/gte';
import SemverMajor from 'semver/functions/major';
import SemverMinor from 'semver/functions/minor';

import * as M from './metadata';
import * as Storage from './storage';
import * as Version from '../version';
import API from './api';
import { anonymize } from './anonymizer';

class Restreamer {
	constructor(address) {
		try {
			new URL(address);
		} catch (e) {
			address = window.location.protocol + '//' + window.location.host;
		}

		// Remove the / at the end
		if (address.slice(-1) === '/') {
			address = address.slice(0, -1);
		}

		this.address = address;
		this.api = new API(this.address);

		this.listeners = [];

		this._reset();
	}

	_reset() {
		this.valid = false;
		this.requiresLogin = true;
		this.connected = false;
		this.refresh = null;
		this.ignoreAPIErrors = false;

		this.about = this._initAbout();

		this.skills = null;
		this.config = null;

		this.cache = {
			assets: new Map(),
		};

		this.channels = new Map();
		this.channel = null;

		this.refreshToken = null;

		this.updates = null;
		this.hasUpdates = false;
		this.hasService = false;

		this._checkForUpdates();
	}

	_initAbout(initialAbout) {
		if (!initialAbout) {
			initialAbout = {};
		}

		const about = {
			id: '',
			name: '',
			created_at: null,
			version: {},
			auths: ['localjwt'],
			...initialAbout,
		};

		about.version = {
			number: '',
			...about.version,
		};

		if (about.created_at !== null) {
			about.created_at = parseRFC3339Date(about.created_at);
		}

		if (!Array.isArray(about.auths)) {
			about.auths = ['localjwt'];
		}

		return about;
	}

	Reset() {
		this._reset();
	}

	ID() {
		return this.about.id;
	}

	CreatedAt() {
		return this.about.created_at;
	}

	Version() {
		return this.about.version;
	}

	App() {
		return this.about.app;
	}

	Name() {
		return this.about.name;
	}

	Address() {
		return this.address;
	}

	SetAddress(address) {
		if (address === this.address) {
			return true;
		}

		try {
			new URL(address);
		} catch (e) {
			return false;
		}

		this.address = address;
		this.api.SetAddress(this.address);
	}

	Auths() {
		return JSON.parse(JSON.stringify(this.about.auths));
	}

	// Events

	AddListener(listener) {
		return this.listeners.push(listener) - 1;
	}

	RemoveListener(id) {
		this.listeners.splice(id, 1);
	}

	_dispatchEvent(severity, type, message) {
		switch (severity) {
			case 'error':
			case 'warning':
			case 'info':
			case 'success':
				break;
			default:
				return;
		}

		for (let l of this.listeners) {
			l({
				severity: severity,
				type: type,
				message: message,
			});
		}
	}

	// API calls

	async _call(fn, ...args) {
		const res = await fn.apply(this.api, args);
		if (res.err !== null && !this.ignoreAPIErrors) {
			if (res.err.code === -1) {
				// Network error
				this._dispatchEvent('error', 'network', res.err.message);
			} else if (res.err.code === 401) {
				if (fn !== this.api.RefreshToken) {
					// Try to refresh access token
					if ((await this.RefreshToken()) === true) {
						// Retry call after successfull token refresh
						const res = await fn.apply(this.api, args);
						if (res.err !== null) {
							if (res.err.status === 'NETWORK_ERROR') {
								// Network error
								this._dispatchEvent('error', 'network', res.err.status, res.err.message);
							} else if (res.err.code === 401) {
								// Auth error
								this._dispatchEvent('error', 'auth', res.err.status, res.err.message);
							}
						}

						return [res.val, res.err];
					}
				}

				if (fn !== this.api.Login && fn !== this.api.LoginWithToken) {
					// Auth error
					this._dispatchEvent('error', 'auth', res.err.message);
				}
			} else {
				// HTTP response error
			}
		}

		return [res.val, res.err];
	}

	IgnoreAPIErrors(toggle) {
		this.ignoreAPIErrors = toggle;
	}

	// Login, Logout and Token handling

	async Validate() {
		let token = this._getRefreshToken();
		if (token !== null) {
			const ok = await this.RefreshToken();
			if (ok === false) {
				this._setRefreshToken(null);
			}
		}

		const about = await this.About();
		if (about === null) {
			return false;
		}

		this.valid = true;
		this.requiresLogin = true;

		if (about.id.length !== 0) {
			if (token === null) {
				this.requiresLogin = false;
			} else {
				this.connected = true;
			}
		}

		this.about = about;

		if (this.IsConnected() === true) {
			await this._init();
		}

		return true;
	}

	async Login(username, password) {
		if (this.requiresLogin === false) {
			await this._init();
			return true;
		}

		const [data, err] = await this._call(this.api.Login, username, password);
		if (err !== null) {
			this._dispatchEvent('error', 'login', i18n._(t`Login failed: ${err.message}`));
			return false;
		}

		this._setAccessToken(data.access_token);
		this._setRefreshToken(data.refresh_token);

		const about = await this.About();
		if (about === null) {
			this._dispatchEvent('error', 'login', i18n._(t`Login failed: Couldn't load API details`));
			return false;
		}

		if (about.id.length === 0) {
			this._dispatchEvent('error', 'login', i18n._(t`Login failed: Couldn't load API details`));
			return false;
		}

		this.about = about;
		this.connected = true;

		await this._init();

		return true;
	}

	async LoginWithToken(token) {
		if (this.requiresLogin === false) {
			await this._init();
			return true;
		}

		const [data, err] = await this._call(this.api.LoginWithToken, token);
		if (err !== null) {
			this._dispatchEvent('error', 'login', i18n._(t`Login failed: ${err.message}`));
			return false;
		}

		this._setAccessToken(data.access_token);
		this._setRefreshToken(data.refresh_token);

		const about = await this.About();
		if (about === null) {
			this._dispatchEvent('error', 'login', i18n._(t`Login failed: Couldn't load API details`));
			return false;
		}

		if (about.id.length === 0) {
			this._dispatchEvent('error', 'login', i18n._(t`Login failed: Couldn't load API details`));
			return false;
		}

		this.about = about;
		this.connected = true;

		await this._init();

		return true;
	}

	Logout() {
		clearTimeout(this.refresh);
		this._setAccessToken(null);
		this._setRefreshToken(null);
	}

	IsConnected() {
		if (this.valid === false) {
			return false;
		}

		if (this.requiresLogin === false) {
			return true;
		}

		if (this.connected === true) {
			return true;
		}

		return false;
	}

	Compatibility() {
		const compatibility = {
			compatible: false,
			core: {
				compatible: false,
				have: '0.0.0',
				want: Version.Core,
			},
			ffmpeg: {
				compatible: false,
				have: '0.0.0',
				want: Version.FFmpeg,
			},
		};

		if (this.IsConnected() === false) {
			return compatibility;
		}

		compatibility.core.have = this.Version().number;
		if (this.skills?.ffmpeg?.version) {
			compatibility.ffmpeg.have = this.skills.ffmpeg.version;
		}

		compatibility.core.compatible = SemverSatisfies(compatibility.core.have, compatibility.core.want);
		compatibility.ffmpeg.compatible = SemverSatisfies(compatibility.ffmpeg.have, compatibility.ffmpeg.want);

		if (compatibility.core.compatible === true && compatibility.ffmpeg.compatible === true) {
			compatibility.compatible = true;
		}

		return compatibility;
	}

	async _init() {
		const compatibility = this.Compatibility();
		if (!compatibility.core.compatible) {
			return;
		}

		await this._initSkills();
		await this._initConfig();
		await this._discoverChannels();
	}

	_setTokenRefresh(expiresIn) {
		clearTimeout(this.refresh);

		if (expiresIn > 60) {
			expiresIn -= 60;
		}

		this.refresh = setTimeout(async () => {
			await this.RefreshToken();
		}, expiresIn * 1000);

		return;
	}

	_setAccessToken(token) {
		if (token === null) {
			this.api.SetToken('');
		} else {
			let claims = null;
			try {
				claims = jwtDecode(token);
				this._setTokenRefresh(claims.exi);
				this.api.SetToken(token);
			} catch (e) {
				this.api.SetToken('');
			}
		}
	}

	_setRefreshToken(token) {
		if (token === null) {
			this.refreshToken = null;
			Storage.Remove('token');
		} else {
			this.refreshToken = token;
			Storage.Set('token', token);
		}
	}

	_getRefreshToken() {
		let token = this.refreshToken;
		if (token === null) {
			token = Storage.Get('token');
		}

		return token;
	}

	async RefreshToken() {
		if (this.requiresLogin === false) {
			return true;
		}

		const token = this._getRefreshToken();
		if (token === null) {
			return false;
		}

		const [data, err] = await this._call(this.api.RefreshToken, token);
		if (err !== null) {
			this._dispatchEvent('error', 'auth', i18n._(t`Failed to refresh token: ${err.message}`));
			return false;
		}

		this._setAccessToken(data.access_token);

		return true;
	}

	// General System Information

	async About() {
		const [val, err] = await this._call(this.api.About);
		if (err !== null) {
			return null;
		}

		if (typeof val !== 'object') {
			return null;
		}

		const about = this._initAbout(val);

		if (about.app !== 'datarhei-core') {
			return null;
		}

		return about;
	}

	async _initSkills() {
		const skills = {
			ffmpeg: {
				version: '',
				version_major: 0,
				version_minor: 0,
			},
			codecs: {
				audio: {
					none: ['none'],
				},
				video: {
					none: ['none'],
				},
			},
			encoders: {
				audio: ['copy', 'none'],
				video: ['copy', 'none'],
			},
			decoders: {
				audio: ['default'],
				video: ['default'],
			},
			formats: {
				demuxers: [],
				muxers: [],
			},
			protocols: {
				input: [],
				output: [],
			},
			filter: [],
			sources: {
				network: [],
				virtualaudio: [],
				virtualvideo: [],
				videoloop: [],
				audioloop: [],
			},
			sinks: {},
		};

		let [val, err] = await this._call(this.api.Skills);
		if (err !== null) {
			this.skills = skills;
			return;
		}

		val = {
			ffmpeg: {},
			codecs: {},
			hwaccels: [],
			formats: {},
			protocols: {},
			devices: {},
			filter: [],
			...val,
		};

		skills.ffmpeg = {
			version: '0.0.0',
			...val.ffmpeg,
		};

		skills.ffmpeg.version_major = SemverMajor(skills.ffmpeg.version);
		skills.ffmpeg.version_minor = SemverMinor(skills.ffmpeg.version);

		val.codecs = {
			audio: {},
			video: {},
			...val.codecs,
		};

		for (let codec of val.codecs.audio) {
			if (codec.encoders !== null) {
				skills.encoders.audio.push(...codec.encoders);

				skills.codecs.audio[codec.id] = [...codec.encoders];
			}

			if (codec.decoders !== null) {
				skills.decoders.audio.push(...codec.decoders);
			}
		}

		for (let codec of val.codecs.video) {
			if (codec.encoders !== null) {
				skills.encoders.video.push(...codec.encoders);

				skills.codecs.video[codec.id] = [...codec.encoders];
			}

			if (codec.decoders !== null) {
				skills.decoders.video.push(...codec.decoders);
			}
		}

		for (let hwaccel of val.hwaccels) {
			skills.decoders.video.push(hwaccel.id);
		}

		for (let filter of val.filter) {
			skills.filter.push(filter.id);
		}

		val.formats = {
			demuxers: [],
			muxers: [],
			...val.formats,
		};

		for (let format of val.formats.demuxers) {
			skills.formats.demuxers.push(format.id);
		}

		for (let format of val.formats.muxers) {
			skills.formats.muxers.push(format.id);
		}

		val.protocols = {
			input: [],
			output: [],
			...val.protocols,
		};

		for (let protocol of val.protocols.input) {
			skills.protocols.input.push(protocol.id);
		}

		for (let protocol of val.protocols.output) {
			skills.protocols.output.push(protocol.id);
		}

		val.devices = {
			demuxers: [],
			muxers: [],
			...val.devices,
		};

		for (let device of val.devices.demuxers) {
			if (!['avfoundation', 'video4linux2', 'alsa', 'fbdev'].includes(device.id)) {
				continue;
			}

			// It's OK to have an empty list of devices because a device might get
			// plugged meanwhile and a refresh is required.
			skills.sources[device.id] = [];

			// Split out a Raspberry Pi camera and create a dedicated source
			if (device.id === 'video4linux2') {
				for (let d of device.devices) {
					if (d.extra.match(/bcm2835[-_]v4l2/) !== null) {
						if (!('raspicam' in skills.sources)) {
							skills.sources['raspicam'] = [];
						}
						skills.sources['raspicam'].push({ ...d });
					} else {
						skills.sources[device.id].push({ ...d });
					}
				}
			} else {
				for (let d of device.devices) {
					skills.sources[device.id].push({ ...d });
				}
			}
		}

		for (let device of val.devices.muxers) {
			if (['fbdev'].includes(device.id)) {
				if (device.devices.length === 0) {
					continue;
				}

				skills.sinks[device.id] = [];

				for (let d of device.devices) {
					skills.sinks[device.id].push({ ...d });
				}
			}
		}

		let channels = (await this.ListRTMPChannels()).map((name) => {
			return {
				media: 'rtmp',
				id: name,
				name: name,
			};
		});

		skills.sources['network'].push(...channels);

		channels = (await this.ListSRTChannels()).map((name) => {
			return {
				media: 'srt',
				id: name,
				name: name,
			};
		});

		skills.sources['network'].push(...channels);

		this.skills = skills;
	}

	Skills() {
		return JSON.parse(JSON.stringify(this.skills));
	}

	async RefreshSkills() {
		const [, err] = await this._call(this.api.SkillsReload);
		if (err !== null) {
			return;
		}

		await this._initSkills();

		return;
	}

	async Config() {
		const [config, err] = await this._call(this.api.Config);
		if (err !== null) {
			return null;
		}

		config.created_at = parseRFC3339Date(config.created_at);
		config.loaded_at = parseRFC3339Date(config.loaded_at);
		config.updated_at = parseRFC3339Date(config.updated_at);

		return config;
	}

	async ConfigSet(config) {
		const res = await this._call(this.api.ConfigSet, config);

		return res;
	}

	async _initConfig() {
		const config = {
			source: {
				network: {
					rtmp: {
						enabled: false,
						secure: false,
						host: '',
						local: 'localhost',
						app: '',
						token: '',
					},
					srt: {
						enabled: false,
						host: '',
						local: 'localhost',
						token: '',
						passphrase: '',
					},
					hls: {
						secure: false,
						host: '',
						local: 'localhost',
						credentials: '',
					},
				},
			},
			http: {
				secure: false,
				host: '',
			},
			memfs: {
				auth: {
					enable: false,
					username: '',
					password: '',
				},
			},
			api: {
				auth: {
					enable: false,
					username: '',
					password: '',
				},
			},
			hostname: '',
			overrides: [],
		};

		const val = await this.Config();
		if (val === null) {
			this.config = config;
			return;
		}

		const isIP = (host) => {
			if (host === 'localhost') {
				return true;
			}

			// IPv4
			if (host.match(/^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$/)) {
				return true;
			}

			// IPv6
			if (host.indexOf(':') !== -1) {
				return true;
			}

			return false;
		};

		const splitHostPort = (address) => {
			let host = '';
			let port = '';
			let hostport = address.split(/:([0-9]+)$/);

			if (hostport.length === 3) {
				host = hostport[0];
				port = hostport[1];
			} else if (hostport.length === 1) {
				port = hostport[0];
			}

			return [host, port];
		};

		const address = new URL(this.Address());

		// The hostname we're using to connect to the API is the one we're
		// going to use to display in the UI.
		let hostname = address.hostname;

		// However, if the provided hostname is an IP and we have at least one
		// name provided, we'll use the first one of the provided names.
		if (isIP(hostname) && val.config.host.name.length !== 0) {
			hostname = val.config.host.name[0];
		}

		config.hostname = hostname;

		// If we're connecting to the API with TLS or if the API is TLS-enabled
		// we upgrade to TLS.
		let protocol = 'http:';
		if (address.protocol === 'https:' || val.config.tls.enable) {
			protocol = 'https:';
		}

		config.http.secure = protocol === 'https:';

		// Set defaults for the port if it's not set.
		let port = address.port;
		if (port.length === 0) {
			port = config.http.secure ? '443' : '80';
		}

		// Set the HTTP host and only add the port if it is not the default one.
		config.http.host = config.hostname;
		if ((config.http.secure && port !== '443') || (!config.http.secure && port !== '80')) {
			config.http.host += ':' + port;
		}

		// HLS

		config.source.network.hls.secure = config.http.secure;
		config.source.network.hls.host = config.http.host;

		// This is used for FFmpeg to access the HLS stream. This will happen always via HTTP.
		// If the HTTP server is bound to a specific address, we'll use this one, localhost otherwise.
		let [http_host, http_port] = splitHostPort(val.config.address);
		config.source.network.hls.local = http_host.length !== 0 ? http_host : 'localhost';
		if (http_port !== '80') {
			config.source.network.hls.local += ':' + http_port;
		}

		// RTMP

		config.source.network.rtmp.enabled = val.config.rtmp.enable;
		config.source.network.rtmp.secure = val.config.rtmp.enable_tls;
		config.source.network.rtmp.token = encodeURIComponent(val.config.rtmp.token);

		// Sanity check on the RTMP app
		let app = val.config.rtmp.app;
		const re = new RegExp('/+', 'g');
		app = app.replace(re, '/');
		if (app !== '/') {
			if (app[app.length - 1] === '/') {
				app = app.substring(0, app.length - 2);
			}

			if (app[0] !== '/') {
				app = '/' + app;
			}
		} else {
			app = '';
		}

		config.source.network.rtmp.app = app;
		config.source.network.rtmp.host = config.hostname;

		// This is used for FFmpeg to access the RTMP stream. If the RTMP server is bound to a
		// specific address, we'll use this one, localhost otherwise.
		const [rtmp_host, rtmp_port] = splitHostPort(val.config.rtmp.address);
		config.source.network.rtmp.local = rtmp_host.length !== 0 ? rtmp_host : 'localhost';
		if (rtmp_port !== '1935') {
			config.source.network.rtmp.host += ':' + rtmp_port;
			config.source.network.rtmp.local += ':' + rtmp_port;
		}

		if (config.source.network.rtmp.secure === true) {
			const [, rtmp_port] = splitHostPort(val.config.rtmp.address_tls);
			if (rtmp_port !== '1935') {
				config.source.network.rtmp.host = config.hostname + ':' + rtmp_port;
			}
		}

		// SRT

		config.source.network.srt.enabled = val.config.srt.enable;
		config.source.network.srt.passphrase = val.config.srt.passphrase;
		config.source.network.srt.token = val.config.srt.token;

		config.source.network.srt.host = config.hostname;

		let [srt_host, srt_port] = splitHostPort(val.config.srt.address);
		config.source.network.srt.local = srt_host.length !== 0 ? srt_host : 'localhost';
		config.source.network.srt.host += ':' + srt_port;
		config.source.network.srt.local += ':' + srt_port;

		// Memfs

		config.memfs.auth.enable = val.config.storage.memory.auth.enable;
		config.memfs.auth.username = val.config.storage.memory.auth.username;
		config.memfs.auth.password = val.config.storage.memory.auth.password;

		if (config.memfs.auth.enable === true) {
			config.source.network.hls.credentials = encodeURIComponent(config.memfs.auth.username) + ':' + encodeURIComponent(config.memfs.auth.password);
		}

		// API Auth

		config.api.auth.enable = val.config.api.auth.enable;
		config.api.auth.username = val.config.api.auth.username;
		config.api.auth.password = val.config.api.auth.password;

		// Environment Config Overrides

		config.overrides = val.overrides;

		this.config = config;
	}

	ConfigActive() {
		const config = JSON.parse(JSON.stringify(this.config));

		config.source.network.channelid = this.channel.channelid;

		return config;
	}

	async ConfigReload() {
		const [, err] = await this._call(this.api.ConfigReload);
		if (err !== null) {
			return false;
		}

		return true;
	}

	ConfigOverrides(name) {
		if (!this.config) {
			return false;
		}

		return this.config.overrides.includes(name);
	}

	ConfigValue(name) {
		if (!this.config) {
			return null;
		}

		const elms = name.split('.');

		let config = this.config;
		for (let e of elms) {
			if (!(e in config)) {
				return null;
			}

			config = config[e];
		}

		return config;
	}

	// Get system metadata
	async GetMetadata(defaults = true) {
		let metadata = await this._getMetadata();

		if (defaults === false) {
			return metadata;
		}

		return M.initMetadata(metadata);
	}

	// Set system metadata
	async SetMetadata(metadata) {
		return await this._setMetadata(metadata);
	}

	// Get the system log
	async Log() {
		const [val, err] = await this._call(this.api.Log);
		if (err !== null) {
			return [];
		}

		return val;
	}

	// Get system resources
	async Resources() {
		return await this._getResources();
	}

	// Get the public HTTP address
	GetPublicHTTPAddress() {
		const config = this.ConfigActive();
		const address = (config.http.secure === true ? 'https://' : 'http://') + config.http.host;

		return address;
	}

	PrefixPublicHTTPAddress(path) {
		const address = this.GetPublicHTTPAddress();

		if (path.match(/^https?:\/\//) !== null) {
			return path;
		}

		if (path.match(/^\//) === null) {
			path = '/' + path;
		}

		return address + path;
	}

	// Get all RTMP/SRT/SNAPSHOT+MEMFS/HLS+MEMFS addresses
	GetPublicAddress(what, channelid) {
		const config = this.ConfigActive();
		const host = config.hostname;

		let address = '';

		function getPort(servicePort) {
			let port = servicePort.split(/:([0-9]+)$/)[1];
			if (port && !port.includes(':')) {
				port = `:${port}`;
			}
			if (port) {
				return port;
			} else {
				return '';
			}
		}

		if (what === 'rtmp') {
			// rtmp/s
			const cfg = config.source.network.rtmp;
			const port = getPort(cfg.host);

			address = 'rtmp';
			if (cfg.secure) {
				address += 's';
			}

			address +=
				`://${host}${port}` +
				(cfg.app.length !== 0 ? cfg.app : '') +
				'/' +
				channelid +
				'.stream' +
				(cfg.token.length !== 0 ? `?token=${cfg.token}` : '');
		} else if (what === 'srt') {
			// srt
			const cfg = config.source.network.srt;
			const port = getPort(cfg.host);

			address =
				`srt://${host}${port}/?mode=caller&transtype=live&streamid=${channelid},mode:request` +
				(cfg.token.length !== 0 ? `,token:${cfg.token}` : '') +
				(cfg.passphrase.length !== 0 ? `&passphrase=${cfg.passphrase}` : '');
		} else if (what === 'snapshot+memfs') {
			// snapshot+memfs
			const port = getPort(config.source.network.hls.host);

			address = (config.http.secure === true ? 'https://' : 'http://') + `${host}${port}/` + this.GetChannelPosterPath(channelid, 'memfs');
		} else if (what === 'snapshot+diskfs') {
			// snapshot+diskfs
			const port = getPort(config.source.network.hls.host);

			address = (config.http.secure === true ? 'https://' : 'http://') + `${host}${port}/` + this.GetChannelPosterPath(channelid, 'diskfs');
		} else if (what === 'hls+memfs') {
			// hls+memfs
			const port = getPort(config.source.network.hls.host);

			address = (config.http.secure === true ? 'https://' : 'http://') + `${host}${port}/` + this.GetChannelManifestPath(channelid, 'memfs');
		} else if (what === 'hls+diskfs') {
			// hls+diskfs
			const port = getPort(config.source.network.hls.host);

			address = (config.http.secure === true ? 'https://' : 'http://') + `${host}${port}/` + this.GetChannelManifestPath(channelid, 'diskfs');
		} else if (what === 'player') {
			// player
			address = (config.http.secure === true ? 'https://' : 'http://') + `${config.http.host}/` + this.GetChannelPlayerPath(channelid);
		}

		return address;
	}

	// Get the iframe codes for the player
	GetPublicIframeCode(channelid) {
		const channel = this.GetChannel(channelid);
		if (channel === null) {
			return '';
		}

		const address = this.GetPublicHTTPAddress();

		return `<iframe src="${address}/${channel.channelid}.html" width="640" height="360" frameborder="no" scrolling="no" allowfullscreen="true"></iframe>`;
	}

	// Channels

	async _discoverChannels() {
		const channels = new Map();

		const reIngest = new RegExp('^restreamer-ui:ingest:([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})$');
		const reEgress = new RegExp('^restreamer-ui:egress:([0-9a-z]+):([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})$');

		const processes = await this._listProcesses(['metadata']);
		const egresses = new Map();

		let hasImported = false;

		for (let p of processes) {
			let matches = reIngest.exec(p.id);

			if (matches === null) {
				matches = reEgress.exec(p.id);
				if (matches === null) {
					continue;
				}

				p.metadata = M.initEgressMetadata(p.metadata);

				const service = matches[1];
				const index = matches[2];
				const channelid = p.reference;

				let egressList = [];

				if (egresses.has(channelid)) {
					egressList = egresses.get(channelid);
				}

				egressList.push({
					id: p.id,
					service: service,
					index: index,
					channelid: p.reference,
					name: p.metadata.name,
				});

				egresses.set(channelid, egressList);

				continue;
			}

			p.metadata = M.initIngestMetadata(p.metadata);

			const channelid = matches[1];
			if (channelid !== p.reference) {
				continue;
			}

			if (!channels.has(channelid)) {
				channels.set(channelid, {
					id: p.id,
					channelid: channelid,
					name: p.metadata.meta.name,
					egresses: new Map(),
					available: true,
				});

				if (p.metadata.imported && p.metadata.imported === true) {
					hasImported = true;
				}
			} else {
				// a channel ID shouldn't exist more than once
				continue;
			}
		}

		for (let [channelid, channel] of channels) {
			if (!egresses.has(channelid)) {
				continue;
			}

			const egressList = egresses.get(channelid);
			for (let egress of egressList) {
				channel.egresses.set(egress.id, egress);
			}

			channels.set(channelid, channel);
		}

		this.channels = channels;

		if (this.channels.size === 0) {
			this.CreateChannel('Livestream');
		}

		let channelid = Storage.Get(this.ID() + ':channel');
		if (!this.channels.has(channelid)) {
			channelid = null;
		}

		if (channelid === null) {
			// Set the first detected channel as default selected channel
			for (let [id] of this.channels) {
				channelid = id;
				break;
			}
		}

		this.SelectChannel(channelid);

		if (hasImported === true) {
			await this.UpdatePlayer(channelid);

			const metadata = await this.GetIngestMetadata(channelid);
			delete metadata.imported;
			await this.SetIngestMetadata(channelid, metadata);
		}
	}

	CreateChannel(name) {
		const channelid = uuidv4();
		this.channels.set(channelid, {
			id: `restreamer-ui:ingest:${channelid}`,
			channelid: channelid,
			name: name,
			egresses: new Map(),
			available: false,
		});

		return channelid;
	}

	async DeleteChannel(channelid) {
		const channel = this.GetChannel(channelid);
		if (!channel) {
			return false;
		}

		await this.StopAllEgresses(channel.channelid);
		await this.DeleteIngest(channel.channelid);
		await this.DeleteIngestSnapshot(channel.channelid);

		for (let egressid of channel.egresses) {
			await this.DeleteEgress(channel.channelid, egressid);
		}

		this.channels.delete(channel.channelid);

		if (this.channels.size === 0) {
			this.CreateChannel();
		}

		// select one of the remaining channels
		for (let [channelid] of this.channels) {
			this.SelectChannel(channelid);
			break;
		}

		return true;
	}

	SelectChannel(channelid) {
		if (!this.channels.has(channelid)) {
			return '';
		}

		Storage.Set(this.ID() + ':channel', channelid);

		this.channel = this.channels.get(channelid);

		return channelid;
	}

	ListChannels() {
		const channels = [];

		for (let channel of this.channels.values()) {
			channels.push({
				id: channel.id,
				channelid: channel.channelid,
				name: channel.name,
				available: channel.available,
				thumbnail: this.GetChannelAddress('snapshot+memfs', channel.channelid),
				egresses: Array.from(channel.egresses.keys()),
			});
		}

		return channels;
	}

	GetChannel(channelid) {
		const channel = this.channels.get(channelid);
		if (!channel) {
			return null;
		}

		return {
			id: channel.id,
			channelid: channel.channelid,
			name: channel.name,
			available: channel.available,
			thumbnail: this.GetChannelAddress('snapshot+memfs', channel.channelid),
			egresses: Array.from(channel.egresses.keys()),
		};
	}

	SetChannel(channelid, channel) {
		let c = this.channels.get(channelid);
		if (!c) {
			return false;
		}

		this.channels.set(channelid, {
			...c,
			...channel,
			egresses: c.egresses,
		});

		return true;
	}

	GetChannelEgress(channelid, id) {
		let channel = this.channels.get(channelid);
		if (!channel) {
			return null;
		}

		const egress = channel.egresses.get(id);
		if (!egress) {
			return null;
		}

		return {
			id: egress.id,
			service: egress.service,
			index: egress.index,
			channelid: egress.channelid,
			name: egress.name,
		};
	}

	SetChannelEgress(channelid, id, data) {
		let channel = this.channels.get(channelid);
		if (!channel) {
			return false;
		}

		channel.egresses.set(id, data);
	}

	DeleteChannelEgress(channelid, id) {
		let channel = this.channels.get(channelid);
		if (!channel) {
			return false;
		}

		channel.egresses.delete(id);
	}

	GetCurrentChannelID() {
		if (this.channel === null) {
			return '';
		}

		return this.channel.channelid;
	}

	// Get the path for the HLS manifest
	GetChannelManifestPath(channelid, storage) {
		if (!storage) {
			storage = 'memfs';
		}

		let url = `${channelid}.m3u8`;
		if (storage === 'memfs') {
			url = 'memfs/' + url;
		}

		return url;
	}

	// Get the path for the poster image
	GetChannelPosterPath(channelid, storage) {
		return `memfs/${channelid}.jpg`;
	}

	// Get the path for the player
	GetChannelPlayerPath(channelid) {
		return `${channelid}.html`;
	}

	GetChannelAddress(what, channelid) {
		const address = this.Address();

		if (what === 'hls+memfs') {
			return `${address}/${this.GetChannelManifestPath(channelid, 'memfs')}`;
		} else if (what === 'hls+diskfs') {
			return `${address}/${this.GetChannelManifestPath(channelid, 'diskfs')}`;
		} else if (what === 'snapshot+memfs') {
			return `${address}/${this.GetChannelPosterPath(channelid, 'memfs')}`;
		} else if (what === 'snapshot+diskfs') {
			return `${address}/${this.GetChannelPosterPath(channelid, 'diskfs')}`;
		} else if (what === 'player') {
			return `${address}/${this.GetChannelPlayerPath(channelid)}`;
		}
	}

	// Sessions

	async CurrentSessions(protocols) {
		const sessions = {
			sessions: 0,
			bitrate_kbit: 0,
		};

		const [val, err] = await this._call(this.api.ActiveSessions, protocols);
		if (err !== null) {
			return sessions;
		}

		for (let p of protocols) {
			if (!(p in val)) {
				continue;
			}

			for (let s of val[p]) {
				if (!s.reference.startsWith(this.channel.channelid)) {
					continue;
				}

				// Don't count viewers where nothing is transmitted.
				if (s.bandwidth_tx_kbit === 0) {
					continue;
				}

				// Don't count ffmpeg as session, only the bandwidth, otherwise
				// the ingest and every running publication service would be
				// considered as a viewer.
				if (p !== 'ffmpeg') {
					sessions.sessions++;
				}

				sessions.bitrate_kbit += s.bandwidth_tx_kbit;
			}
		}

		return sessions;
	}

	// Ingest

	// Check whether there's an ingest defined or not
	HasIngest() {
		if (!this.channel) {
			return false;
		}

		return this.channel.available;
	}

	// Get process information for ingest
	async GetIngest(channelid, filter = []) {
		const channel = this.GetChannel(channelid);
		if (channel === null) {
			return null;
		}

		return await this._getProcess(channel.id, filter);
	}

	// Get the ingest metadata
	async GetIngestMetadata(channelid) {
		const channel = this.GetChannel(channelid);
		if (channel === null) {
			return M.initIngestMetadata({});
		}

		let metadata = await this._getProcessMetadata(channel.id);

		metadata = M.initIngestMetadata(metadata);
		if (metadata.meta.name.length === 0) {
			metadata.meta.name = this.channel.name;
		}

		return metadata;
	}

	// Set the ingest metadata
	async SetIngestMetadata(channelid, metadata) {
		const channel = this.GetChannel(channelid);
		if (channel === null) {
			return false;
		}

		this.SetChannel(channelid, {
			name: metadata.meta.name ? metadata.meta.name : channel.name,
		});

		return await this._setProcessMetadata(channel.id, metadata);
	}

	// Get the ingest progress
	async GetIngestProgress(channelid) {
		const channel = this.GetChannel(channelid);
		if (channel === null) {
			return this._getProgressFromState(null);
		}

		const state = await this._getProcessState(channel.id);

		return this._getProgressFromState(state);
	}

	// Get the ingest log
	async GetIngestLog(channelid) {
		const channel = this.GetChannel(channelid);
		if (channel === null) {
			return null;
		}

		return await this._getProcessLog(channel.id);
	}

	// Get the ingest debug log
	async GetIngestDebug(channelid) {
		const channel = this.GetChannel(channelid);
		if (channel === null) {
			return {};
		}

		return await this.GetDebug(channel.id);
	}

	// Start the ingest process
	async StartIngest(channelid) {
		const channel = this.GetChannel(channelid);
		if (channel === null) {
			return false;
		}

		return await this._startProcess(channel.id);
	}

	// Start the ingest snapshot process
	async StartIngestSnapshot(channelid) {
		const channel = this.GetChannel(channelid);
		if (channel === null) {
			return false;
		}

		return await this._startProcess(channel.id + '_snapshot');
	}

	// Stop the ingest process
	async StopIngest(channelid) {
		const channel = this.GetChannel(channelid);
		if (channel === null) {
			return false;
		}

		return await this._stopProcess(channel.id);
	}

	// Stop the ingest snapshot process
	async StopIngestSnapshot(channelid) {
		const channel = this.GetChannel(channelid);
		if (channel === null) {
			return false;
		}

		return await this._stopProcess(channel.id + '_snapshot');
	}

	// Delete the ingest process
	async DeleteIngest(channelid) {
		const channel = this.GetChannel(channelid);
		if (channel === null) {
			return false;
		}

		return await this._deleteProcess(channel.id);
	}

	// Delete the ingest snaphot process
	async DeleteIngestSnapshot(channelid) {
		const channel = this.GetChannel(channelid);
		if (channel === null) {
			return false;
		}

		return await this._deleteProcess(channel.id + '_snapshot');
	}

	// Update/Insert the ingest process
	async UpsertIngest(channelid, global, inputs, outputs, control) {
		const channel = this.GetChannel(channelid);
		if (!channel) {
			return [null, { message: 'Unknown channel ID' }];
		}

		const proc = {
			type: 'ffmpeg',
			id: channel.id,
			reference: channel.channelid,
			input: [],
			output: [],
			options: ['-loglevel', 'level+info', '-err_detect', 'ignore_err', ...global],
			autostart: control.process.autostart,
			reconnect: control.process.reconnect,
			reconnect_delay_seconds: parseInt(control.process.delay),
			stale_timeout_seconds: parseInt(control.process.staleTimeout),
			limits: {
				cpu_usage: parseInt(control.limits.cpu_usage),
				memory_mbytes: parseInt(control.limits.memory_mbytes),
				waitfor_seconds: parseInt(control.limits.waitfor_seconds),
			},
		};

		for (let i in inputs) {
			const input = inputs[i];

			proc.input.push({
				id: 'input_' + i,
				address: input.address,
				options: input.options.map((o) => '' + o),
			});
		}

		// 1. Add output address and options

		// 1.1 Fetch core config
		const core_config = this.ConfigActive();

		// 1.2 Fetch rtmp settings
		const rtmp_config = core_config.source.network.rtmp;
		let rtmp_enabled = false;
		if (control.rtmp && control.rtmp.enable && rtmp_config.enabled) {
			rtmp_enabled = true;
		}
		if (
			proc.input[0].address.includes('rtmp://localhost') &&
			proc.input[0].address.includes(channel.channelid) &&
			!proc.input[0].address.includes('ingest')
		) {
			rtmp_enabled = false;
			control.rtmp.enable = true;
		}

		// 1.3 Fetch srt settings
		const srt_config = core_config.source.network.srt;
		let srt_enabled = false;
		if (control.srt.enable && srt_config.enabled) {
			srt_enabled = true;
		}
		if (
			proc.input[0].address.includes('srt://localhost') &&
			proc.input[0].address.includes(channel.channelid) &&
			!proc.input[0].address.includes('ingest')
		) {
			srt_enabled = false;
			control.srt.enable = true;
		}

		// 1.4 'tee_muxer' is required for the delivery of one output to multiple endpoints without processing the input for each output
		// http://ffmpeg.org/ffmpeg-all.html#tee-1
		let tee_muxer = false;
		if (rtmp_enabled || srt_enabled) {
			tee_muxer = true;
		}

		// 1.5 Set hls filename vars
		const hlsStorage = control.hls.storage;
		let segmentPlaylistPath = `${channel.channelid}` + (control.hls.master_playlist ? `_{outputid}` : '');
		let segmentFilePath = `${channel.channelid}` + (control.hls.master_playlist ? `_{outputid}_%04d` : '_%04d');
		if (hlsStorage === 'diskfs') {
			// diskfs (path structure)
			segmentFilePath = `${channel.channelid}` + (control.hls.master_playlist ? `/{outputid}/%Y%m%d/%s` : '/%Y%m%d/%s');
		}

		// 1.6 Set hls filenames
		const hls_master_playlist = `${channel.channelid}.m3u8`;
		const hls_fmp4_init_filename = `${channel.channelid}.mp4`;
		const hls_segment_playlist = `{${hlsStorage}}/${segmentPlaylistPath}.m3u8`;
		const hls_segment_filename =
			`{${hlsStorage}` + (tee_muxer ? '^:' : '') + `}/${segmentFilePath}.` + (!control.hls.lhls && control.hls.version === 7 ? 'mp4' : 'ts');

		// 2. Add output address

		const output = {
			id: 'output_0',
			address: hls_segment_playlist,
			options: ['-dn', '-sn', ...outputs[0].options.map((o) => '' + o)],
			cleanup: [],
		};

		// 3. Extend output options

		// 3.1 Injects a metadata link as title
		const metadata = `${this.GetPublicHTTPAddress()}/${channel.channelid}/oembed.json`;
		const metadata_options = ['-metadata', `title=${metadata}`, '-metadata', 'service_provider=datarhei-Restreamer'];
		output.options.push(...metadata_options);

		// Manifest versions
		// https://developer.apple.com/documentation/http_live_streaming/about_the_ext-x-version_tag
		// https://ffmpeg.org/ffmpeg-all.html#Options-53

		// 3.2 Returns the raw l/hls parameters for an EXT-X-VERSION

		// fix Malformed AAC bitstream detected for hls version 7
		let hls_aac_adtstoasc = false;

		const getHLSParams = (lhls, version) => {
			if (lhls) {
				// lhls
				return [
					['f', 'dash'],
					['strict', 'experimental'],
					['hls_playlist', '1'],
					['init_seg_name', `init-${channel.channelid}.$ext$`],
					['media_seg_name', `chunk-${channel.channelid}-$Number%05d$.$ext$`],
					['master_m3u8_publish_rate', '1'],
					['adaptation_sets', 'id=0,streams=v id=1,streams=a'],
					['lhls', '1'],
					['streaming', '1'],
					['seg_duration', '' + parseInt(control.hls.segmentDuration)],
					['frag_duration', '0.5'],
					['use_template', '1'],
					['remove_at_exit', '0'],
					['window_size', '' + parseInt(control.hls.listSize)],
					['http_persistent', '0'],
				];
			} else {
				// hls
				switch (version) {
					case 6:
						return [
							['f', 'hls'],
							['start_number', '0'],
							['hls_time', '' + parseInt(control.hls.segmentDuration)],
							['hls_list_size', '' + parseInt(control.hls.listSize)],
							['hls_flags', 'append_list+delete_segments+program_date_time+independent_segments+temp_file'],
							['hls_delete_threshold', '4'],
							['hls_segment_filename', hls_segment_filename],
						];
					case 7:
						// fix Malformed AAC bitstream detected for hls version 7
						if (output.options.includes('-codec:a') && output.options.includes('copy')) {
							if (!tee_muxer) {
								output.options.push('-bsf:a', 'aac_adtstoasc');
							}
							hls_aac_adtstoasc = true;
						}
						return [
							['f', 'hls'],
							['start_number', '0'],
							['hls_time', '' + parseInt(control.hls.segmentDuration)],
							['hls_list_size', '' + parseInt(control.hls.listSize)],
							['hls_flags', 'append_list+delete_segments+program_date_time+independent_segments+temp_file'],
							['hls_delete_threshold', '4'],
							['hls_segment_type', 'fmp4'],
							['hls_fmp4_init_filename', hls_fmp4_init_filename],
							['hls_fmp4_init_resend', '1'],
							['hls_segment_filename', hls_segment_filename],
						];
					// case 3
					default:
						return [
							['f', 'hls'],
							['start_number', '0'],
							['hls_time', '' + parseInt(control.hls.segmentDuration)],
							['hls_list_size', '' + parseInt(control.hls.listSize)],
							['hls_flags', 'append_list+delete_segments+program_date_time+temp_file'],
							['hls_delete_threshold', '4'],
							['hls_segment_filename', hls_segment_filename],
						];
				}
			}
		};
		const hls_params_raw = getHLSParams(control.hls.lhls, control.hls.version);

		// 3.3 Use strftime for DiskFS
		if (control.hls.storage && control.hls.storage === 'diskfs') {
			hls_params_raw.push(['strftime', '1'], ['strftime_mkdir', '1']);
		}

		// 3.4 Push master playlist params
		if (control.hls.master_playlist) {
			hls_params_raw.push(['master_pl_name', hls_master_playlist], ['master_pl_publish_rate', `${control.hls.segmentDuration}`]);
		}

		// 3.5 Use HTTP method
		if (control.hls.storage && control.hls.storage !== 'diskfs') {
			hls_params_raw.push(['method', 'PUT']);
		}

		// 3.6 Overwrite output files
		proc.options.push('-y');

		// 3.7 Returns the l/hls parameters with or without tee_muxer
		if (tee_muxer) {
			// f=hls:start_number=0...
			const hls_params = hls_params_raw
				.filter((o) => {
					// unsupported in tee_muxer
					return !(o[0] === 'segment_format_options' || o[0] === 'max_muxing_queue_size');
				})
				.map((o) => o[0] + '=' + o[1])
				.join(':');

			// set flags
			if (control.process.low_delay) {
				output.options.push('-flags', '+low_delay+global_header');
			} else {
				output.options.push('-flags', '+global_header');
			}

			output.options.push('-tag:v', '7', '-tag:a', '10', '-f', 'tee');
			// ['f=hls:start_number=0...]address.m3u8
			// use tee_muxer formatting
			output.address =
				`[${hls_aac_adtstoasc ? 'bsfs/a=aac_adtstoasc:' : ''}${hls_params}]${hls_segment_playlist}` +
				(rtmp_enabled ? `|[f=flv]{rtmp,name=${channel.channelid}.stream}` : '') +
				(srt_enabled ? `|[f=mpegts]{srt,name=${channel.channelid},mode=publish}` : '');
		} else {
			// ['-f', 'hls', '-start_number', '0', ...]
			// adding the '-' in front of the first option, then flatten everything
			const hls_params = hls_params_raw.map((o) => ['-' + o[0], o[1]]).reduce((acc, val) => acc.concat(val), []);

			// set flags
			if (control.process.low_delay) {
				output.options.push('-flags', '+low_delay');
			}

			output.options.push(...hls_params);
		}

		// 4. Add output cleanup jobs

		// 4.1 Set hls cleanup filename vars
		let cleanupSegmentFilePath = `${channel.channelid}` + (control.hls.master_playlist ? `_{outputid}_**` : '_**');
		if (hlsStorage === 'diskfs') {
			// diskfs (path structure)
			cleanupSegmentFilePath = `${channel.channelid}` + (control.hls.master_playlist ? `/{outputid}/**` : '/**');
		}

		// 4.2 Set hls cleanup filenames
		const cleanup_global = `${hlsStorage}:/${channel.channelid}**`;
		const cleanup_hls_master_playlist = `${hlsStorage}:/${channel.channelid}.m3u8`;
		const cleanup_hls_fmp4_init_filename = `${hlsStorage}:/${channel.channelid}.mp4`;
		const cleanup_hls_segment_playlist = `${hlsStorage}:/${segmentPlaylistPath}.m3u8`;
		const cleanup_hls_segment_filename = `${hlsStorage}:/${cleanupSegmentFilePath}.` + (!control.hls.lhls && control.hls.version === 7 ? 'mp4' : 'ts');

		// 4.3 Cleanup id* (process is deleted) + continuously hls_segment_playlist and hls_segment_filename
		output.cleanup.push(
			{
				pattern: cleanup_global,
				purge_on_delete: true,
			},
			{
				pattern: cleanup_hls_segment_playlist,
				max_file_age_seconds: control.hls.cleanup ? parseInt(control.hls.segmentDuration) * (parseInt(control.hls.listSize) + 6) : 0,
				purge_on_delete: true,
			},
			{
				pattern: cleanup_hls_segment_filename,
				max_files: parseInt(control.hls.listSize) + 6,
				max_file_age_seconds: control.hls.cleanup ? parseInt(control.hls.segmentDuration) * (parseInt(control.hls.listSize) + 6) : 0,
				purge_on_delete: true,
			},
		);

		// 4.4 Cleanup hls_master_playlist
		if (control.hls.master_playlist) {
			output.cleanup.push({
				pattern: cleanup_hls_master_playlist,
				max_file_age_seconds: control.hls.cleanup ? parseInt(control.hls.segmentDuration) * (parseInt(control.hls.listSize) + 6) : 0,
				purge_on_delete: true,
			});
		}

		// 4.5 Cleanup hls_fmp4_init_filename
		if (!control.hls.lhls && control.hls.version === 7) {
			output.cleanup.push({
				pattern: cleanup_hls_fmp4_init_filename,
				max_file_age_seconds: control.hls.cleanup ? parseInt(control.hls.segmentDuration) * (parseInt(control.hls.listSize) + 6) : 0,
				purge_on_delete: true,
			});
		}

		// 5. Push output
		proc.output.push(output);

		const [val, err] = await this._upsertProcess(channel.id, proc);
		if (err !== null) {
			return [val, err];
		}

		this.SetChannel(channelid, {
			...channel,
			available: true,
		});

		return [val, null];
	}

	// Upsert the ingest snapshot process
	async UpsertIngestSnapshot(channelid, control) {
		const channel = this.GetChannel(channelid);
		if (channel === null) {
			return [null, { message: 'Unknown channel ID' }];
		}

		// Set hls storage endpoint
		const hlsStorage = control.hls.storage;

		const snapshot = {
			type: 'ffmpeg',
			id: channel.id + '_snapshot',
			reference: channel.channelid,
			input: [
				{
					id: 'input_0',
					address: `{${hlsStorage}}/${channel.channelid}.m3u8`,
					options: [],
				},
			],
			output: [
				{
					id: 'output_0',
					address: `{memfs}/${channel.channelid}.jpg`,
					options: ['-vframes', '1', '-f', 'image2', '-update', '1'],
					cleanup: [
						{
							pattern: `memfs:/${channel.channelid}.jpg`,
							purge_on_delete: true,
						},
					],
				},
			],
			options: ['-err_detect', 'ignore_err'],
			autostart: control.process.autostart,
			reconnect: true,
			reconnect_delay_seconds: parseInt(control.snapshot.interval),
			stale_timeout_seconds: 30,
		};

		const [val, err] = await this._upsertProcess(channel.id + '_snapshot', snapshot);
		if (err !== null) {
			return [val, err];
		}

		return [val, null];
	}

	// Check whether the manifest of the ingest process is available
	async HasIngestFiles(channelid) {
		const channel = this.GetChannel(channelid);
		if (channel === null) {
			return false;
		}

		const [, err] = await this._call(this.api.MemFSHasFile, `/${channel.channelid}.m3u8`);
		if (err !== null) {
			return false;
		}

		return true;
	}

	// Probe an external stream
	async Probe(channelid, inputs) {
		const channel = this.GetChannel(channelid);
		if (channel === null) {
			return [null, { message: 'Unknown channel ID' }];
		}

		const id = `${channel.id}_probe`;

		const config = {
			type: 'ffmpeg',
			id: id,
			reference: channel.channelid,
			input: [],
			output: [
				{
					id: 'output_0',
					address: '-',
					options: ['-dn', '-sn', '-codec', 'copy', '-f', 'null'],
				},
			],
			options: [],
			autostart: false,
			reconnect: false,
		};

		for (let i in inputs) {
			const input = inputs[i];

			config.input.push({
				id: 'input_' + i,
				address: input.address,
				options: input.options.map((o) => '' + o),
			});
		}

		await this._deleteProcess(id);

		let [val, err] = await this._call(this.api.ProcessAdd, config);
		if (err !== null) {
			return [val, err];
		}

		[val, err] = await this._call(this.api.ProcessProbe, id);
		await this._deleteProcess(id);

		return [val, err];
	}

	// Probe the ingest stream
	async ProbeIngest(channelid) {
		const channel = this.GetChannel(channelid);
		if (channel === null) {
			return [null, { message: 'Unknown channel ID' }];
		}

		const id = `${channel.id}_probe`;

		const config = {
			type: 'ffmpeg',
			id: id,
			reference: channel.channelid,
			input: [
				{
					id: 'input_0',
					address: `#${channel.id}:output=output_0`,
					options: [],
				},
			],
			output: [
				{
					id: 'output_0',
					address: '-',
					options: ['-dn', '-sn', '-codec', 'copy', '-f', 'null'],
				},
			],
			options: [],
			autostart: false,
			reconnect: false,
		};

		await this._deleteProcess(id);

		let [val, err] = await this._call(this.api.ProcessAdd, config);
		if (err !== null) {
			return [val, err];
		}

		[val, err] = await this._call(this.api.ProcessProbe, id);
		await this._deleteProcess(id);

		return [val, err];
	}

	// Selfhosted Player

	// Set defaults for the settings of the selfhosted player
	InitPlayerSettings(initSettings) {
		if (!initSettings) {
			initSettings = {};
		}

		const settings = {
			autoplay: false,
			mute: false,
			statistics: false,
			chromecast: false,
			airplay: false,
			color: {},
			ga: {},
			logo: {},
			poster: '',
			...initSettings,
		};

		settings.color = {
			seekbar: '#ffffff',
			buttons: '#ffffff',
			...settings.color,
		};

		settings.ga = {
			account: '',
			name: '',
			...settings.ga,
		};

		settings.logo = {
			image: '',
			position: 'top-left',
			link: '',
			...settings.logo,
		};

		return settings;
	}

	// Update the player the selfthosted player
	async UpdatePlayer(channelid) {
		const channel = this.GetChannel(channelid);
		if (channel === null) {
			return false;
		}

		let metadata = await this.GetIngestMetadata(channelid);

		// update the player files
		const playerType = 'videojs';
		if ((await this._updatePlayerEssentials(playerType)) === false) {
			return false;
		}

		metadata.player = this.InitPlayerSettings(metadata.player);

		const templateData = {
			channelid: channelid,
			name: metadata.meta.name,
			description: metadata.meta.description,
			author_name: metadata.meta.author.name,
			author_url: this.GetPublicAddress('player', channelid),
			license: metadata.license,
			iframecode: this.GetPublicIframeCode(channelid),
			poster: this.GetChannelPosterPath(channelid, metadata.control.hls.storage),
			poster_url: this.GetPublicAddress('snapshot+memfs', channelid),
			width: 640,
			height: 360,
			chromecast: metadata.player.chromecast,
			airplay: metadata.player.airplay,
		};

		if (metadata.player.poster.length !== 0) {
			templateData.poster = metadata.player.poster.replace(/^\/+/, '');
			templateData.poster_url = this.PrefixPublicHTTPAddress(metadata.player.poster);
		}

		// upload player.html
		let player = await this._getLocalAssetAsString(`/_player/${playerType}/player.html`);
		player = Handlebars.compile(player)(templateData);
		await this._uploadAssetData(`/${channelid}.html`, player);

		// upload oembed.json
		let embed = await this._getLocalAssetAsString('/_player/oembed.json.in');
		embed = Handlebars.compile(embed)({
			...templateData,
			name: JSON.stringify(templateData.name),
			description: JSON.stringify(templateData.description),
			author_name: JSON.stringify(templateData.author_name),
			author_url: JSON.stringify(templateData.author_url),
			license: JSON.stringify(templateData.license),
			iframecode: JSON.stringify(templateData.iframecode),
			poster: JSON.stringify(templateData.poster),
			poster_url: JSON.stringify(templateData.poster_url),
		});
		await this._uploadAssetData(`/channels/${channelid}/oembed.json`, embed);

		// upload oembed.xml
		embed = await this._getLocalAssetAsString('/_player/oembed.xml.in');
		embed = Handlebars.compile(embed)(templateData);
		await this._uploadAssetData(`/channels/${channelid}/oembed.xml`, embed);

		await this.UpdatePlayerConfig(channelid, metadata);

		return true;
	}

	async UpdatePlayerConfig(channelid, metadata) {
		metadata.player = this.InitPlayerSettings(metadata.player);

		const playerConfig = {
			...metadata.player,
			source: this.GetChannelManifestPath(channelid, metadata.control.hls.storage),
			poster: this.GetChannelPosterPath(channelid, metadata.control.hls.storage),
			license: {
				license: metadata.license,
				title: metadata.meta.name,
				author: metadata.meta.author.name,
			},
		};

		if (metadata.player.logo.image.length !== 0) {
			playerConfig.logo.image = metadata.player.logo.image.replace(/^\/+/, '');
		}

		if (metadata.player.poster.length !== 0) {
			playerConfig.poster = metadata.player.poster.replace(/^\/+/, '');
		}

		await this._uploadAssetData(`/channels/${channelid}/config.js`, 'var playerConfig = ' + JSON.stringify(playerConfig));
	}

	// Upload channel specific channel data
	async UploadData(channelid, name, data) {
		if (channelid.length === 0) {
			channelid = this.GetCurrentChannelID();
		}

		const channel = this.GetChannel(channelid);
		if (channel === null) {
			return '';
		}

		// sanitize name
		name = name.replace(/[^0-9a-z.]/gi, '');

		const path = `/channels/${channel.channelid}/${name}`;

		await this._uploadAssetData(path, data);

		return path;
	}

	// Upload a logo for the selfhosted player
	async UploadLogo(channelid, data, extension) {
		return this.UploadData(channelid, 'logo.' + extension, data);
	}

	// Upload a poster image for the selfhosted player
	async UploadPoster(channelid, data, extension) {
		return this.UploadData(channelid, 'poster.' + extension, data);
	}

	// Playersite

	// Set defaults for the settings of the playersite
	InitPlayersiteSettings(initSettings) {
		if (!initSettings) {
			initSettings = {};
		}

		const settings = {
			player: 'videojs',
			playersite: false,
			channelid: 'current',
			channel_list: this.ListChannels().map((c) => c.channelid),
			title: 'restreamer',
			share: true,
			support: true,
			chromecast: false,
			airplay: false,
			template: '!default',
			templatename: '',
			textcolor_title: 'rgba(255,255,255,1)',
			textcolor_default: 'rgba(230,230,230,1)',
			textcolor_link: 'rgba(230,230,230,1)',
			textcolor_link_hover: 'rgba(255,255,255,1)',
			bgcolor_default: 'rgba(56,56,56,1)',
			bgcolor_selected: 'rgba(0,0,0,.15)',
			bgcolor_unselected: 'rgba(255,255,255,.05)',
			bgcolor_header: 'rgba(44,44,44,1)',
			hrcolor: 'rgba(102,102,102,1)',
			bgimage_url: '',
			inject1: '',
			inject2: '',
			inject3: '',
			inject4: '',
			imprint: '',
			terms: '',

			...initSettings,
		};

		return settings;
	}

	// Get the URL for the playersite
	GetPlayersiteUrl() {
		return `index.html`;
	}

	// Is a playersite already available
	async HasPlayersite() {
		return await this._hasAsset('/index.html');
	}

	// Update the playersite
	async UpdatePlayersite() {
		let metadata = await this.GetMetadata();

		const settings = this.InitPlayersiteSettings(metadata.playersite);
		settings.player = 'videojs';

		if (settings.playersite === false) {
			await this._removePlayersiteEssentials();
			return true;
		}

		// update the player files
		await this._updatePlayerEssentials(settings.player);

		// update the playersite files
		await this._updatePlayersiteEssentials();

		// use preferred channel
		let channel = this.GetChannel(settings.channelid);
		if (channel === null) {
			channel = this.GetChannel(this.GetCurrentChannelID());
		}

		// filter channels based on the main channel and additional channels
		const channels = this.ListChannels().filter((c) => settings.channel_list.indexOf(c.channelid) !== -1 || c.channelid === channel.channelid);

		// Handlebars function ifEquals
		Handlebars.registerHelper('ifEquals', function (arg1, arg2, options) {
			return arg1 === arg2 ? options.fn(this) : options.inverse(this);
		});

		Handlebars.registerHelper('ifnoteq', function (arg1, arg2, options) {
			if (arg1 !== arg2) {
				return options.fn(this);
			}
			return options.inverse(this);
		});

		for (const item of channels) {
			const ingestMetadata = await this.GetIngestMetadata(item.channelid);
			const templateData = {
				player: settings.player,
				playersite: settings.playersite,
				title: settings.title,
				share: settings.share,
				support: settings.support,
				chromecast: settings.chromecast,
				airplay: settings.airplay,
				url: this.GetPlayersiteUrl(),
				textcolor_title: settings.textcolor_title,
				textcolor_default: settings.textcolor_default,
				textcolor_link: settings.textcolor_link,
				textcolor_link_hover: settings.textcolor_link_hover,
				bgcolor_header: settings.bgcolor_header,
				bgcolor_selected: settings.bgcolor_selected,
				bgcolor_unselected: settings.bgcolor_unselected,
				hrcolor: settings.hrcolor,
				bgcolor_default: settings.bgcolor_default,
				bgimage_url: settings.bgimage_url,
				imprint_html: settings.imprint.replace(/(?:\r\n|\r|\n)/g, '<br />'),
				terms_html: settings.terms.replace(/(?:\r\n|\r|\n)/g, '<br />'),
				inject1: settings.inject1,
				inject2: settings.inject2,
				inject3: settings.inject3,
				inject4: settings.inject4,
				channels: channels,
				channel_id: item.channelid,
				channel_name: ingestMetadata.meta.name,
				channel_description: ingestMetadata.meta.description,
				channel_description_html: ingestMetadata.meta.description.replace(/(?:\r\n|\r|\n)/g, '<br />'),
				channel_creator_name: ingestMetadata.meta.author.name,
				channel_creator_description: ingestMetadata.meta.author.description,
				channel_creator_description_html: ingestMetadata.meta.author.description.replace(/(?:\r\n|\r|\n)/g, '<br />'),
				channel_license: ingestMetadata.license,
				channel_poster: this.GetChannelPosterPath(item.channelid, ingestMetadata.control.hls.storage),
				channel_width: 640,
				channel_height: 360,
			};

			// upload playersite.html
			let playersite = '';
			if (settings.template !== '!default') {
				playersite = await this.GetPlayersiteTemplate(settings.template);
				if (playersite.length === 0) {
					settings.template = '!default';
				}
			}

			if (settings.template === '!default') {
				playersite = await this._getLocalAssetAsString('/_playersite/index.html');
			}

			if (item.channelid === channel.channelid) {
				playersite = Handlebars.compile(playersite)(templateData);
				await this._uploadAssetData('/index.html', playersite);
			}

			playersite = Handlebars.compile(playersite)(templateData);
			await this._uploadAssetData('/playersite_' + item.channelid + '.html', playersite);

			// Upload player config for each channel
			await this.UpdatePlayerConfig(item.channelid, ingestMetadata);
		}

		// Upload player implementation
		await this._uploadAssetData('/playersite/player.js', await this._getLocalAssetAsString(`/_playersite/${settings.player}.js`));

		return true;
	}

	// Upload the background image for the playersite
	async UploadPlayersiteBackgroundImage(data, extension) {
		// sanitize extension
		extension = extension.replace(/[^0-9a-z]/gi, '');

		const path = `/playersite/bg.${extension}`;

		await this._uploadAssetData(path, data);

		return path;
	}

	// Upload a playersite template file
	async UploadPlayersiteTemplate(data, name) {
		// sanitize name
		name = name.replace(/[^0-9a-z]/gi, '');

		const path = `/playersite/templates/${name}.html`;

		await this._uploadAssetData(path, data);

		return name;
	}

	// Delete a playersite template file
	async DeletePlayersiteTemplate(name) {
		// sanitize name
		name = name.replace(/[^0-9a-z]/gi, '');

		const path = `/playersite/templates/${name}.html`;

		await this._deleteAsset(path);

		return true;
	}

	async GetPlayersiteTemplate(name) {
		// sanitize name
		name = name.replace(/[^0-9a-z]/gi, '');

		const path = `/playersite/templates/${name}.html`;

		const data = await this._getAssetAsString(path);

		return data;
	}

	async ListPlayersiteTemplates() {
		let templates = await this._listAssets('/playersite/templates/*');

		templates = templates.map((t) => {
			const components = t.split('/');
			const name = components[components.length - 1].split('.', 1)[0];
			return name;
		});

		return templates;
	}

	// Egress

	GetEgressId(service, id) {
		return `restreamer-ui:egress:${service}:${id}`;
	}

	// Get process information for egress
	async GetEgress(channelid, id, filter = []) {
		const channel = this.GetChannel(channelid);
		if (!channel) {
			return;
		}

		if (!channel.egresses.includes(id)) {
			return null;
		}

		return await this._getProcess(id, filter);
	}

	// Get metadata for egress
	async GetEgressMetadata(channelid, id) {
		let metadata = null;

		const channel = this.GetChannel(channelid);
		if (!channel) {
			return M.initEgressMetadata(metadata);
		}

		if (channel.egresses.includes(id)) {
			metadata = await this._getProcessMetadata(id);
		}

		return M.initEgressMetadata(metadata);
	}

	// Set metadata for egress
	async SetEgressMetadata(channelid, id, metadata) {
		const channel = this.GetChannel(channelid);
		if (!channel) {
			return null;
		}

		if (!channel.egresses.includes(id)) {
			return null;
		}

		const egress = this.channel.egresses.get(id);

		egress.name = metadata.name ? metadata.name : '';

		return await this._setProcessMetadata(id, metadata);
	}

	// Start egress process
	async StartEgress(channelid, id) {
		const channel = this.GetChannel(channelid);
		if (!channel) {
			return null;
		}

		if (!channel.egresses.includes(id)) {
			return null;
		}

		return await this._startProcess(id);
	}

	// Stop egress process
	async StopEgress(channelid, id) {
		const channel = this.GetChannel(channelid);
		if (!channel) {
			return null;
		}

		if (!channel.egresses.includes(id)) {
			return null;
		}

		return await this._stopProcess(id);
	}

	// Stop all egress processes
	async StopAllEgresses(channelid) {
		const channel = this.GetChannel(channelid);
		if (!channel) {
			return;
		}

		for (let egressid of channel.egresses) {
			await this._stopProcess(egressid);
		}

		return;
	}

	// Delete egress process
	async DeleteEgress(channelid, id) {
		const channel = this.GetChannel(channelid);
		if (!channel) {
			return null;
		}

		if (!channel.egresses.includes(id)) {
			return null;
		}

		const res = await this._deleteProcess(id);

		if (res === true) {
			this.channel.egresses.delete(id);
		}

		return res;
	}

	// Get the egress log
	async GetEgressLog(channelid, id) {
		const channel = this.GetChannel(channelid);
		if (!channel) {
			return null;
		}

		if (!channel.egresses.includes(id)) {
			return null;
		}

		return await this._getProcessLog(id);
	}

	// Get the egress debug log
	async GetEgressDebug(channelid, id) {
		const channel = this.GetChannel(channelid);
		if (!channel) {
			return null;
		}

		if (!channel.egresses.includes(id)) {
			return null;
		}

		return await this.GetDebug(id);
	}

	// Update an egress process
	async UpdateEgress(channelid, id, global, inputs, outputs, control) {
		const channel = this.GetChannel(channelid);
		if (!channel) {
			return null;
		}

		if (!channel.egresses.includes(id)) {
			return null;
		}

		const egress = this.GetChannelEgress(channelid, id);
		if (!egress) {
			return null;
		}

		if (!Array.isArray(outputs)) {
			outputs = [outputs];
		}

		// from the inputs only the first is used and only its options are considered.

		let address = '';
		let options = [];
		if (control.source.source === 'hls+memfs') {
			address = `{memfs}/${channel.channelid}.m3u8`;
			options.push('-re');
		} else if (control.source.source === 'hls+diskfs') {
			address = `{diskfs}/${channel.channelid}.m3u8`;
			options.push('-re');
		} else if (control.source.source === 'rtmp') {
			address = `{rtmp,name=${channel.channelid}.stream}`;
			const skills = this.Skills();
			if (skills.ffmpeg.version_major >= 6) {
				const codecs = [];
				if (skills.codecs.video.hevc?.length > 0) {
					codecs.push('hvc1');
				}
				if (skills.codecs.video.av1?.length > 0) {
					codecs.push('av01');
				}
				if (skills.codecs.video.vp9?.length > 0) {
					codecs.push('vp09');
				}

				if (codecs.length !== 0) {
					options.push('-rtmp_enhanced_codecs', codecs.join(','));
				}
			}
		} else if (control.source.source === 'srt') {
			address = `{srt,name=${channel.channelid},mode=request}`;
		}

		const config = {
			type: 'ffmpeg',
			id: egress.id,
			reference: egress.channelid,
			input: [
				{
					id: 'input_0',
					address: address,
					options: [...options, ...inputs[0].options],
				},
			],
			output: [],
			options: ['-loglevel', 'level+info', '-err_detect', 'ignore_err', ...global],
			autostart: control.process.autostart,
			reconnect: control.process.reconnect,
			reconnect_delay_seconds: parseInt(control.process.delay),
			stale_timeout_seconds: parseInt(control.process.staleTimeout),
			limits: {
				cpu_usage: parseInt(control.limits.cpu_usage),
				memory_mbytes: parseInt(control.limits.memory_mbytes),
				waitfor_seconds: parseInt(control.limits.waitfor_seconds),
			},
		};

		for (let i in outputs) {
			const output = outputs[i];

			if (!Array.isArray(output.options)) {
				output.options = [];
			}

			// set flags
			let options = [];
			if (control.process.low_delay) {
				options.push('-flags', '+low_delay');
			}

			options.push(...output.options.map((o) => '' + o));

			config.output.push({
				id: 'output_' + i,
				address: output.address,
				options: options,
			});
		}

		let [val, err] = await this._upsertProcess(egress.id, config);
		return [val, err];
	}

	// Create an egress process
	async CreateEgress(channelid, service, global, inputs, outputs, control) {
		const channel = this.GetChannel(channelid);
		if (!channel) {
			return ['', { message: 'Unknown channel ID' }];
		}

		const index = uuidv4();
		const egress = {
			id: `restreamer-ui:egress:${service}:${index}`,
			index: index,
			service: service,
			channelid: channel.channelid,
			name: '',
		};

		this.SetChannelEgress(channelid, egress.id, egress);

		const [, err] = await this.UpdateEgress(channelid, egress.id, global, inputs, outputs, control);
		if (err !== null) {
			this.DeleteChannelEgress(channelid, egress.id);
		}

		return [egress.id, err];
	}

	// Ingest + Egresses

	async ListIngestEgresses(channelid, services = []) {
		const channel = this.GetChannel(channelid);
		if (!channel) {
			return [];
		}

		const re = new RegExp('^restreamer-ui:egress:');

		let list = await this._listProcesses(['state'], channel.channelid);

		list = list.filter((p) => {
			if (p.id === channel.id) {
				p.index = '';
				p.service = 'player';
				p.name = channel.name;

				return true;
			}

			const matches = re.exec(p.id);
			if (matches !== null) {
				const egress = this.GetChannelEgress(channelid, p.id);
				if (!egress) {
					return false;
				}

				p.service = egress.service;
				p.index = egress.index;
				p.name = egress.name;

				return true;
			}

			return false;
		});

		list.sort((a, b) => {
			if (a.service === 'player') {
				return -1;
			} else if (b.service === 'player') {
				return 1;
			}

			let astring = a.name;
			let bstring = b.name;

			if (a.name === b.name) {
				astring = a.index;
				bstring = b.index;
			}

			astring = astring.toUpperCase();
			bstring = bstring.toUpperCase();

			if (astring < bstring) {
				return -1;
			}

			if (astring > bstring) {
				return 1;
			}

			return 0;
		});

		return list;
	}

	async ListProcesses(filter = [], ids = []) {
		return await this._listProcesses(filter, '', ids);
	}

	async GetDebug(processid) {
		const about = await this._getAboutDebug();
		const skills = await this.Skills();
		const config = await this._getConfigDebug();
		const proc = await this._getProcessDebug(processid);

		const data = {
			about: about,
			ffmpeg: skills.ffmpeg,
			config: config,
			process: proc,
		};

		return data;
	}

	// RTMP

	async ListRTMPChannels() {
		return await this._listRTMPChannels();
	}

	// SRT

	async ListSRTChannels() {
		return await this._listSRTChannels();
	}

	// Expert Mode

	IsExpert() {
		return Storage.Get('expert') === 'true';
	}

	SetExpert(value) {
		Storage.Set('expert', !!value);
	}

	// Check for Updates

	CheckForUpdates() {
		return Storage.Get('updates') === 'true';
	}

	SetCheckForUpdates(value) {
		Storage.Set('updates', !!value);

		this._checkForUpdates();
	}

	HasUpdates() {
		if (!this.CheckForUpdates()) {
			return false;
		}

		return this.hasUpdates;
	}

	HasService() {
		return this.hasService;
	}

	async _checkForUpdates() {
		if (Storage.Get('updates') !== 'false') {
			Storage.Set('updates', true);
		}

		clearTimeout(this.updates);

		if (!this.CheckForUpdates()) {
			return;
		}

		(async () => {
			let response = null;

			try {
				response = await fetch('https://service.datarhei.com/api/v1/app_version', {
					method: 'PUT',
					headers: {
						'Content-Type': 'application/json',
					},
					body: JSON.stringify({
						app_version: Version.UI,
					}),
				});
			} catch (err) {
				return;
			}

			const contentType = response.headers.get('Content-Type');
			let isJSON = false;

			if (contentType != null) {
				isJSON = contentType.indexOf('application/json') !== -1;
			}

			if (isJSON === false) {
				return;
			}

			if (response.ok === false) {
				return;
			}

			const value = {
				latest_version: Version.UI,
				...(await response.json()),
			};

			const findVersion = (name) => {
				const matches = name.match(/v(\d+\.\d+\.\d+)\s*$/);
				if (matches === null) {
					return '0.0.0';
				}

				return matches[1];
			};

			const currentVersion = findVersion(Version.UI);
			const announcedVersion = findVersion(value.latest_version);

			if (currentVersion !== '0.0.0') {
				if (SemverGt(announcedVersion, currentVersion)) {
					this.hasUpdates = true;
				} else {
					this.hasUpdates = false;
				}
			}

			const serviceVersion = findVersion(value.service_version);
			if (SemverGte(serviceVersion, '1.0.0')) {
				this.hasService = true;
			} else {
				this.hasService = false;
			}
		})();

		this.updates = setTimeout(
			() => {
				this._checkForUpdates();
			},
			1000 * 60 * 60,
		);
	}

	// Private system related function

	async _setMetadata(data) {
		const [, err] = await this._call(this.api.SetMetadata, 'restreamer-ui', data);
		if (err !== null) {
			return false;
		}

		return true;
	}

	async _getMetadata() {
		const [val, err] = await this._call(this.api.GetMetadata, 'restreamer-ui');
		if (err !== null) {
			return null;
		}

		return val;
	}

	// Private process related functions

	async _listProcesses(filter = [], reference = '', ids = []) {
		const [val, err] = await this._call(this.api.Processes, reference, ids, filter);
		if (err !== null) {
			return [];
		}

		for (let i = 0; i < val.length; i++) {
			val[i] = this._sanitizeProcess(val[i]);
		}

		return val;
	}

	async _getProcess(id, filter = []) {
		const [val, err] = await this._call(this.api.Process, id, filter);
		if (err !== null) {
			return null;
		}

		return this._sanitizeProcess(val);
	}

	_sanitizeProcess(proc) {
		if (!proc.id) {
			proc.id = '';
		}

		if (!proc.config) {
			proc.config = null;
		}

		if (!proc.state) {
			proc.state = null;
		}

		proc.progress = this._getProgressFromState(proc.state);

		if (!proc.report) {
			proc.report = null;
		}

		if (!proc.metadata) {
			proc.metadata = {};
		}

		if (proc.metadata['restreamer-ui']) {
			proc.metadata = proc.metadata['restreamer-ui'];
		} else {
			proc.metadata = {};
		}

		return proc;
	}

	async _getProcessConfig(id) {
		const [val, err] = await this._call(this.api.ProcessConfig, id);
		if (err !== null) {
			return null;
		}

		return val;
	}

	async _getProcessState(id) {
		const [val, err] = await this._call(this.api.ProcessState, id);
		if (err !== null) {
			return null;
		}

		return val;
	}

	async _getProcessLog(id) {
		const [val, err] = await this._call(this.api.ProcessReport, id);
		if (err !== null) {
			return null;
		}

		return val;
	}

	async _getProcessDebug(id) {
		const [p, err] = await this._call(this.api.Process, id, ['config', 'state', 'report']);
		if (err !== null) {
			return null;
		}

		if (p.config) {
			p.config.options = p.config.options.map(anonymize);

			for (let i in p.config.input) {
				p.config.input[i].address = anonymize(p.config.input[i].address);
				p.config.input[i].options = p.config.input[i].options.map(anonymize);
			}

			for (let i in p.config.output) {
				p.config.output[i].address = anonymize(p.config.output[i].address);
				p.config.output[i].options = p.config.output[i].options.map(anonymize);
			}
		}

		if (p.state) {
			for (let i in p.state.progress.inputs) {
				p.state.progress.inputs[i].address = anonymize(p.state.progress.inputs[i].address);
			}

			for (let i in p.state.progress.outputs) {
				p.state.progress.outputs[i].address = anonymize(p.state.progress.outputs[i].address);
			}

			if (!p.state.command) {
				p.state.command = [];
			}

			p.state.command = p.state.command.map(anonymize);
			p.state.last_logline = anonymize(p.state.last_logline);
		}

		if (p.report) {
			if (!p.report.prelude) {
				p.report.prelude = [];
			}

			p.report.prelude = p.report.prelude.map(anonymize);
			p.report.log = p.report.log.map((l) => [l[0], anonymize(l[1])]);

			for (let i in p.report.history) {
				p.report.history[i].prelude = p.report.history[i].prelude.map(anonymize);
				p.report.history[i].log = p.report.history[i].log.map((l) => [l[0], anonymize(l[1])]);
			}
		}

		return p;
	}

	async _startProcess(id) {
		const [, err] = await this._call(this.api.ProcessCommand, id, 'start');
		if (err !== null) {
			return false;
		}

		return true;
	}

	async _stopProcess(id) {
		const [, err] = await this._call(this.api.ProcessCommand, id, 'stop');
		if (err !== null) {
			return false;
		}

		return true;
	}

	async _upsertProcess(id, config) {
		const [val, err] = await this._call(this.api.ProcessUpdate, id, config);
		if (err !== null) {
			if (err.code === 404) {
				return await this._call(this.api.ProcessAdd, config);
			}
		}

		return [val, err];
	}

	async _deleteProcess(id) {
		const [, err] = await this._call(this.api.ProcessDelete, id);
		if (err !== null) {
			if (err.code === 404) {
				return true;
			}

			return false;
		}

		return true;
	}

	async _setProcessMetadata(id, data) {
		const [, err] = await this._call(this.api.ProcessSetMetadata, id, 'restreamer-ui', data);
		if (err !== null) {
			return false;
		}

		return true;
	}

	async _getProcessMetadata(id) {
		const [val, err] = await this._call(this.api.ProcessGetMetadata, id, 'restreamer-ui');
		if (err !== null) {
			return null;
		}

		return val;
	}

	// Assets

	async _updatePlayerEssentials(player) {
		// get the list of supplemental files for the player
		const data = await this._getLocalAssetAsString(`/_player/${player}/files.txt`);
		if (data === null) {
			return false;
		}

		const files = data.split(/\n/);

		// upload player files
		for (let file of files) {
			if (file.length === 0) {
				continue;
			}

			await this._uploadLocalAsset(`/_player/${player}/${file}`, `/player/${player}/${file}`);
		}

		await this._updatePublicEssentials();

		return true;
	}

	async _updatePlayersiteEssentials() {
		// upload playersite background
		await this._uploadLocalAsset('/_playersite/bg.jpg', '/playersite/default_bg.jpg');

		await this._updatePublicEssentials();

		return true;
	}

	async _removePlayersiteEssentials() {
		await this._deleteAsset('/index.html');
		await this._deleteAsset('/playersite/default_bg.jpg');
		await this._deleteAsset('/playersite/bg.jpg');
		await this._deleteAsset('/playersite/bg.png');
	}

	async _updatePublicEssentials() {
		// upload robots.txt
		await this._uploadLocalAsset('/robots.txt', '/robots.txt');

		// upload playersite favicons
		await this._uploadLocalAsset('/favicon.ico', '/favicon.ico');
		await this._uploadLocalAsset('/logo192.png', '/logo192.png');
		await this._uploadLocalAsset('/logo512.png', '/logo512.png');
	}

	async _removePublicEssentials() {
		await this._deleteAsset('/robots.txt');
		await this._deleteAsset('/favicon.ico');
		await this._deleteAsset('/logo192.png');
		await this._deleteAsset('/logo512.png');
	}

	async _getLocalAssetAsString(localPath) {
		let data = await this._getLocalAsset(localPath);
		if (data === null) {
			return null;
		}

		const text = await data.text();

		return text;
	}

	async _getLocalAsset(localPath) {
		let data = this.cache.assets.get(localPath);
		if (data === undefined) {
			let response = null;

			try {
				response = await fetch(process.env.PUBLIC_URL + localPath, {
					method: 'GET',
				});
			} catch (err) {
				return null;
			}

			data = await response.blob();

			if (response.ok === false) {
				return null;
			}

			this.cache.assets.set(localPath, data);
		}

		return data;
	}

	async _uploadLocalAsset(localPath, remotePath) {
		const data = await this._getLocalAsset(localPath);
		if (data === null) {
			return false;
		}

		await this._uploadAssetData(remotePath, data);

		return true;
	}

	async _uploadAssetData(remotePath, data) {
		await this._call(this.api.DataPutFile, remotePath, data);

		return true;
	}

	async _deleteAsset(remotePath) {
		await this._call(this.api.DataDeleteFile, remotePath);

		return true;
	}

	async _hasAsset(remotePath) {
		const [, err] = await this._call(this.api.DataHasFile, remotePath);
		if (err !== null) {
			return false;
		}

		return true;
	}

	async _getAssetAsString(remotePath) {
		const [val, err] = await this._call(this.api.DataGetFile, remotePath);
		if (err !== null) {
			return '';
		}

		return val;
	}

	async _listAssets(remotePathPattern) {
		const [val, err] = await this._call(this.api.DataListFiles, remotePathPattern);
		if (err !== null) {
			return [];
		}

		return val.map((f) => f.name);
	}

	async _listRTMPChannels() {
		const [val, err] = await this._call(this.api.RTMPChannels);
		if (err !== null) {
			return [];
		}

		return val;
	}

	async _listSRTChannels() {
		const [val, err] = await this._call(this.api.SRTChannels);
		if (err !== null) {
			return [];
		}

		return val;
	}

	async _getAboutDebug() {
		const about = await this.About();

		about.auths = about.auths.map((a) => (a.startsWith('auth0 ') ? 'auth0' : a));

		return about;
	}

	async _getConfigDebug() {
		const [data, err] = await this._call(this.api.Config);
		if (err !== null) {
			return null;
		}

		const config = data.config;

		config.host.name = config.host.name.map((e) => '[anonymized]');

		config.api.auth.username = '[anonymized]';
		config.api.auth.password = '[anonymized]';
		config.api.auth.jwt.secret = '[anonymized]';

		config.api.auth.auth0.tenants = config.host.name.map((e) => '[anonymized]');

		config.api.access.http.allow = config.api.access.http.allow.map((e) => '[anonymized]');
		config.api.access.http.block = config.api.access.http.block.map((e) => '[anonymized]');
		config.api.access.https.allow = config.api.access.https.allow.map((e) => '[anonymized]');
		config.api.access.https.block = config.api.access.https.block.map((e) => '[anonymized]');

		config.storage.memory.auth.username = '[anonymized]';
		config.storage.memory.auth.password = '[anonymized]';

		config.storage.s3 = config.storage.s3.map((e) => {
			return {
				...e,
				auth: {
					...e.auth,
					username: '[anonymized]',
					password: '[anonymized]',
				},
				endpoint: '[anonymized]',
				access_key_id: '[anonymized]',
				secret_access_key: '[anonymized]',
			};
		});

		if (config.storage.cors.origins.length !== 1 || config.storage.cors.origins[0] !== '*') {
			config.storage.cors.origins = config.storage.cors.origins.map((e) => '[anonymized]');
		}

		config.rtmp.app = '[anonymized]';
		config.rtmp.token = '[anonymized]';

		config.ffmpeg.binary = anonymize(config.ffmpeg.binary);

		config.service.token = '[anonymized]';

		config.sessions.ip_ignorelist = config.sessions.ip_ignorelist.map((e) => '[anonymized]');

		return config;
	}

	_getProgressFromState(state) {
		const progress = {
			valid: false,
			order: 'stop',
			state: 'disconnected',
			error: '',
			reconnect: -1,
			bitrate: 0,
			fps: 0,
			time: 0,
			speed: 0,
			q: -1,
			frames: 0,
			drop: 0,
			dup: 0,
			command: [],
			cpu: 0,
			memory: 0,
		};

		if (state === null) {
			return progress;
		}

		progress.valid = true;
		progress.order = state.order;
		progress.command = state.command.slice();

		const fps = state.progress.fps || 0;

		if (state.exec === 'starting') {
			progress.state = 'connecting';
		} else if (state.exec === 'running') {
			if (state.runtime_seconds >= 10 && fps !== 0) {
				progress.state = 'connected';
			} else {
				progress.state = 'connecting';
			}
		} else if (state.exec === 'finishing') {
			progress.state = 'disconnecting';
		} else if (state.exec === 'finished') {
			progress.state = 'disconnected';
			progress.reconnect = state.reconnect_seconds === undefined ? -1 : state.reconnect_seconds;
		} else if (state.exec === 'killed' || state.exec === 'failed') {
			progress.state = 'error';
			progress.error = state.last_logline || '';
			progress.reconnect = state.reconnect_seconds === undefined ? -1 : state.reconnect_seconds;
		}

		if (progress.state === 'connected') {
			progress.bitrate = state.progress.bitrate_kbit || 0;
			progress.fps = state.progress.fps || 0;
			progress.time = state.runtime_seconds || 0;
			progress.speed = state.progress.speed || 0;
			progress.q = state.progress.quality === undefined ? -1 : state.progress.quality;
			progress.frames = state.progress.frames || 0;
			progress.drop = state.progress.drop || 0;
			progress.dup = state.progress.dup || 0;
			progress.cpu = state.cpu_usage || 0;
			progress.memory = state.memory_bytes || 0;
		}

		return progress;
	}

	async _getResources() {
		const query = {
			metrics: [
				{ name: 'uptime_uptime' },
				{ name: 'cpu_ncpu' },
				{ name: 'cpu_idle' },
				{ name: 'mem_total' },
				{ name: 'mem_free' },
				{ name: 'filesystem_limit' },
				{ name: 'filesystem_usage' },
				{ name: 'session_limit' },
				{ name: 'session_active' },
				{ name: 'session_txbitrate' },
				{ name: 'session_maxtxbitrate' },
			],
		};
		const [data, err] = await this._call(this.api.Metrics, query);
		if (err !== null) {
			return null;
		}

		const getMetrics = (metrics, metric, labels) => {
			loop: for (const m of metrics) {
				if (m.name !== metric) {
					continue;
				}

				if (labels !== null) {
					for (const key in labels) {
						const value = labels[key];

						if (!(key in m.labels)) {
							continue loop;
						}

						if (m.labels[key] !== value) {
							continue loop;
						}
					}
				}

				return m;
			}

			return null;
		};

		const getValue = (metrics, metric, labels) => {
			const m = getMetrics(metrics, metric, labels);
			if (m === null) {
				return 0;
			}

			return m.values[0][1];
		};

		const metrics = data.metrics;

		const resources = {
			uptime_seconds: getValue(metrics, 'uptime_uptime'),
			system: {
				cpu_used: 100 - getValue(metrics, 'cpu_idle'),
				cpu_ncores: getValue(metrics, 'cpu_ncpu'),
				mem_used_bytes: getValue(metrics, 'mem_total') - getValue(metrics, 'mem_free'),
				mem_total_bytes: getValue(metrics, 'mem_total'),
			},
			core: {
				session_used: getValue(metrics, 'session_active', { collector: 'hls' }),
				session_limit: getValue(metrics, 'session_limit', { collector: 'hls' }),
				net_used_kbit: getValue(metrics, 'session_txbitrate', { collector: 'hls' }) / 1024,
				net_limit_kbit: getValue(metrics, 'session_maxtxbitrate', { collector: 'hls' }) / 1024,
				memfs_used_bytes: getValue(metrics, 'filesystem_usage', { name: 'memfs' }),
				memfs_limit_bytes: getValue(metrics, 'filesystem_limit', { name: 'memfs' }),
				disk_used_bytes: getValue(metrics, 'filesystem_usage', { name: 'diskfs' }),
				disk_limit_bytes: getValue(metrics, 'filesystem_limit', { name: 'diskfs' }),
			},
		};

		return resources;
	}
}

const dateRegex = /^(\d{4})-(\d{2})-(\d{2})[T ](\d{2}):(\d{2}):(\d{2})(?:\.?(\d+))?(?:(?:([+-]\d{2}):?(\d{2}))|Z)?$/;

function parseRFC3339Date(d) {
	const m = dateRegex.exec(d);
	if (m === null) {
		return null;
	}

	// Milliseconds are optional.
	if (m[7] === undefined) {
		m[7] = 0;
	} else {
		m[7] = parseInt((1.0 / parseFloat(m[7])) * 100);
	}

	// If timezone is undefined, it must be Z or nothing (otherwise the group would have captured).
	if (m[8] === undefined && m[9] === undefined) {
		// Use UTC.
		m[8] = 0;
		m[9] = 0;
	}

	var year = +m[1];
	var month = +m[2];
	var day = +m[3];
	var hour = +m[4];
	var minute = +m[5];
	var second = +m[6];
	var msec = +m[7];
	var tzHour = +m[8];
	var tzMin = +m[9];
	var tzOffset = tzHour * 60 + tzMin;

	return new Date(Date.UTC(year, month - 1, day, hour, minute - tzOffset, second, msec));
}

export default Restreamer;
