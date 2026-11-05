import React from 'react';
import { useNavigate } from 'react-router-dom';

import { useLingui } from '@lingui/react';
import { Trans, t } from '@lingui/macro';
import makeStyles from '@mui/styles/makeStyles';
import urlparser from 'url-parse';
import Accordion from '@mui/material/Accordion';
import AccordionDetails from '@mui/material/AccordionDetails';
import AccordionSummary from '@mui/material/AccordionSummary';
import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown';
import Button from '@mui/material/Button';
import Grid from '@mui/material/Grid';
import Icon from '@mui/icons-material/AccountTree';
import MenuItem from '@mui/material/MenuItem';
import RefreshIcon from '@mui/icons-material/Refresh';
import TextField from '@mui/material/TextField';
import Typography from '@mui/material/Typography';
import WarningIcon from '@mui/icons-material/Warning';

import BoxTextarea from '../../../misc/BoxTextarea';
import BoxText from '../../../misc/BoxText';
import Checkbox from '../../../misc/Checkbox';
import FormInlineButton from '../../../misc/FormInlineButton';
import MultiSelect from '../../../misc/MultiSelect';
import MultiSelectOption from '../../../misc/MultiSelectOption';
import Password from '../../../misc/Password';
import Select from '../../../misc/Select';
import Textarea from '../../../misc/Textarea';

const useStyles = makeStyles((theme) => ({
	gridContainer: {
		marginTop: '0.5em',
	},
}));

const initSettings = (initialSettings, config) => {
	if (!initialSettings) {
		initialSettings = {};
	}

	const settings = {
		mode: 'pull',
		address: '',
		username: '',
		password: '',
		push: {},
		rtsp: {},
		http: {},
		general: {},
		...initialSettings,
	};

	settings.push = {
		type: 'rtmp',
		name: config.channelid,
		...settings.push,
	};

	settings.rtsp = {
		udp: false,
		stimeout: 5000000,
		...settings.rtsp,
	};

	settings.http = {
		readNative: true,
		forceFramerate: false,
		framerate: 25,
		userAgent: '',
		referer: '',
		http_proxy: '',
		...settings.http,
	};

	settings.general = {
		analyzeduration: 5_000_000, // microseconds, 5s
		analyzeduration_rtmp: 3_000_000, // 3s
		analyzeduration_http: 20_000_000, // 20s
		probesize: 5_000_000, // bytes
		max_probe_packets: 2500,
		fflags: ['genpts'],
		thread_queue_size: 512,
		copyts: false,
		start_at_zero: false,
		use_wallclock_as_timestamps: false,
		avoid_negative_ts: 'auto',
		...settings.general,
	};

	return settings;
};

const initConfig = (initialConfig) => {
	if (!initialConfig) {
		initialConfig = {};
	}

	const config = {
		rtmp: {},
		srt: {},
		hls: {},
		channelid: 'external',
		...initialConfig,
	};

	config.rtmp = {
		enabled: false,
		secure: false,
		host: 'localhost',
		local: 'localhost',
		app: '',
		token: '',
		...config.rtmp,
	};

	config.srt = {
		enabled: false,
		host: 'localhost',
		local: 'localhost',
		token: '',
		passphrase: '',
		...config.srt,
	};

	config.hls = {
		secure: false,
		host: 'localhost',
		local: 'localhost',
		credentials: '',
		...config.hls,
	};

	return config;
};

const initSkills = (initialSkills) => {
	if (!initialSkills) {
		initialSkills = {};
	}

	const skills = {
		ffmpeg: {},
		formats: {},
		protocols: {},
		...initialSkills,
	};

	skills.ffmpeg = {
		version: '0.0.0',
		version_major: 0,
		version_minor: 0,
		...skills.ffmpeg,
	};

	skills.formats = {
		demuxers: [],
		...skills.formats,
	};

	skills.protocols = {
		input: [],
		...skills.protocols,
	};

	if (skills.formats.demuxers.includes('rtsp')) {
		if (!skills.protocols.input.includes('rtsp')) {
			skills.protocols.input.push('rtsp');
		}

		if (skills.protocols.input.includes('tls')) {
			if (!skills.protocols.input.includes('rtsps')) {
				skills.protocols.input.push('rtsps');
			}
		}
	}

	return skills;
};

const createInputs = (settings, config, skills) => {
	config = initConfig(config);
	settings = initSettings(settings, config);
	skills = initSkills(skills);

	const input = {
		address: '',
		options: [],
	};

	if (settings.mode === 'push') {
		let name = settings.push.name;
		if (settings.push.type === 'hls') {
			if (name === 'none') {
				name = config.channelid;
			}
			input.address = getLocalHLS(name);
		} else if (settings.push.type === 'rtmp') {
			if (name === config.channelid) {
				name += '.stream';
			}
			input.address = getLocalRTMP(name);
		} else if (settings.push.type === 'srt') {
			if (name === config.channelid) {
				name += '.stream';
			}
			input.address = getLocalSRT(name);
		} else {
			input.address = '';
		}
	} else {
		input.address = settings.address;
	}

	// registrate protocol by address
	const protocol = getProtocolClass(input.address);

	// general settings (pull/push)
	if (settings.general.fflags.length !== 0) {
		input.options.push('-fflags', '+' + settings.general.fflags.join('+'));
	}
	input.options.push('-thread_queue_size', settings.general.thread_queue_size);
	input.options.push('-probesize', settings.general.probesize);

	if (settings.general.max_probe_packets !== 2500) {
		input.options.push('-max_probe_packets', settings.general.max_probe_packets);
	}
	if (settings.general.copyts) {
		input.options.push('-copyts');
	}
	if (settings.general.start_at_zero) {
		input.options.push('-start_at_zero');
	}
	if (settings.general.use_wallclock_as_timestamps) {
		input.options.push('-use_wallclock_as_timestamps', '1');
	}
	if (skills.ffmpeg.version_major >= 5 && settings.general.avoid_negative_ts !== 'auto') {
		input.options.push('-avoid_negative_ts', settings.general.avoid_negative_ts);
	}

	// general settings > analyzeduration by protocol
	//
	// old settings:
	// analyzeduration: 20s for http and 3s for rtmp streams
	if (settings.mode === 'push') {
		if (settings.push.type === 'hls') {
			input.options.push('-analyzeduration', settings.general.analyzeduration_http);
		} else if (settings.push.type === 'rtmp') {
			input.options.push('-analyzeduration', settings.general.analyzeduration_rtmp);

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
					input.options.push('-rtmp_enhanced_codecs', codecs.join(','));
				}
			}
		} else if (settings.push.type === 'srt') {
			input.options.push('-analyzeduration', settings.general.analyzeduration);
		}
	} else {
		input.address = addUsernamePassword(input.address, settings.username, settings.password);

		if (protocol === 'http') {
			input.options.push('-analyzeduration', settings.general.analyzeduration_http);

			if (settings.http.readNative === true) {
				input.options.push('-re');
			}

			if (settings.http.forceFramerate === true) {
				input.options.push('-r', settings.http.framerate);
			}

			if (settings.http.userAgent.length !== 0) {
				input.options.push('-user_agent', settings.http.userAgent);
			}

			if (settings.http.referer.length !== 0) {
				input.options.push('-referer', settings.http.referer);
			}

			if (settings.http.http_proxy.length !== 0) {
				input.options.push('-http_proxy', settings.http.http_proxy);
			}
		} else if (protocol === 'rtmp') {
			input.options.push('-analyzeduration', settings.general.analyzeduration_rtmp);

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
					input.options.push('-rtmp_enhanced_codecs', codecs.join(','));
				}
			}
		} else {
			input.options.push('-analyzeduration', settings.general.analyzeduration);

			if (protocol === 'rtsp') {
				if (skills.ffmpeg.version_major === 4) {
					input.options.push('-stimeout', settings.rtsp.stimeout);
				} else {
					input.options.push('-timeout', settings.rtsp.stimeout);
				}

				if (settings.rtsp.udp === true) {
					input.options.push('-rtsp_transport', 'udp');
				} else {
					input.options.push('-rtsp_transport', 'tcp');
				}
			}
		}
	}

	return [input];
};

const addUsernamePassword = (address, username, password) => {
	if (username.length === 0 && password.length === 0) {
		return address;
	}

	if (isAuthProtocol(address) === false) {
		return address;
	}

	const url = urlparser(address, {});
	if (username.length !== 0) {
		url.set('username', username);
	}
	if (password.length !== 0) {
		url.set('password', password);
	}

	return url.toString();
};

const getProtocol = (url) => {
	const re = new RegExp('^([a-z][a-z0-9.+-:]*)://', 'i');
	const matches = url.match(re);

	if (matches === null || matches.length === 0) {
		return '';
	}

	return matches[1];
};

const getProtocolClass = (url) => {
	if (typeof url !== 'string') {
		url = '';
	}

	const protocol = getProtocol(url);
	if (protocol.length === 0) {
		return '';
	}

	if (/rtmp(e|s|t)?/.test(protocol) === true) {
		return 'rtmp';
	} else if (/https?/.test(protocol) === true) {
		return 'http';
	} else if (/mms(t|h)/.test(protocol) === true) {
		return 'mms';
	} else if (/rtsps?/.test(protocol) === true) {
		return 'rtsp';
	}

	return protocol;
};

const isAuthProtocol = (url) => {
	const protocolClass = getProtocolClass(url);

	// eslint-disable-next-line default-case
	switch (protocolClass) {
		case 'amqp':
		case 'ftp':
		case 'http':
		case 'icecast':
		case 'mms':
		case 'rtmp':
		case 'sftp':
		case 'rtsp':
			return true;
	}

	return false;
};

const isSupportedProtocol = (url, supportedProtocols) => {
	const protocol = getProtocol(url);
	if (protocol.length === 0) {
		return false;
	}

	if (!supportedProtocols.includes(protocol)) {
		return false;
	}

	return true;
};

const getHLSAddress = (host, credentials, name, secure) => {
	// Test for IPv6 addresses and put brackets around
	let url = 'http' + (secure ? 's' : '') + '://' + (credentials.length !== 0 ? credentials + '@' : '') + host + '/memfs/' + name + '.m3u8';

	return url;
};

const getRTMPAddress = (host, app, name, token, secure) => {
	let url = 'rtmp' + (secure ? 's' : '') + '://' + host + app + '/' + name + '.stream';

	if (token.length !== 0) {
		url += '/' + encodeURIComponent(token);
	}

	return url;
};

const getSRTAddress = (host, name, token, passphrase, publish) => {
	let url =
		'srt' +
		'://' +
		host +
		'?mode=caller&transtype=live&streamid=' +
		name +
		'.stream,mode:' +
		(publish ? 'publish' : 'request') +
		(token.length !== 0 ? ',token:' + encodeURIComponent(token) : '');

	if (passphrase.length !== 0) {
		url += '&passphrase=' + encodeURIComponent(passphrase);
	}

	return url;
};

const getHLS = (config, name) => {
	const url = getHLSAddress(config.hls.host, config.hls.credentials, config.channelid, config.hls.secure);

	return url;
};

const getRTMP = (config) => {
	const url = getRTMPAddress(config.rtmp.host, config.rtmp.app, config.channelid, config.rtmp.token, config.rtmp.secure);

	return url;
};

const getSRT = (config) => {
	const url = getSRTAddress(config.srt.host, config.channelid, config.srt.token, config.srt.passphrase, true);

	return url;
};

const getLocalHLS = (config, name) => {
	let url = getHLSAddress(config.hls.local, '', config.channelid, false);

	return url;
};

const getLocalRTMP = (name) => {
	return '{rtmp,name=' + name + '}';
};

const getLocalSRT = (name) => {
	return '{srt,name=' + name + ',mode=request}';
};

const isValidURL = (address) => {
	const protocol = getProtocolClass(address);
	if (protocol.length === 0) {
		return false;
	}

	return true;
};

function AdvancedSettings(props) {
	const settings = props.settings;
	let protocolClass = getProtocolClass(settings.address);
	if (settings.mode === 'push') {
		switch (settings.push.type) {
			case 'rtmp':
				protocolClass = 'rtmp';
				break;
			case 'srt':
				protocolClass = 'srt';
				break;
			default:
		}
	}

	return (
		<Grid item xs={12}>
			<Accordion className="accordion">
				<AccordionSummary elevation={0} expandIcon={<ArrowDropDownIcon />}>
					<Typography>
						<Trans>Advanced settings</Trans>
					</Typography>
				</AccordionSummary>
				<AccordionDetails>
					<Grid container spacing={2}>
						{protocolClass === 'rtsp' && (
							<React.Fragment>
								<Grid item xs={12}>
									<Typography variant="h3">
										<Trans>RTSP</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Checkbox label={<Trans>UDP transport</Trans>} checked={settings.rtsp.udp} onChange={props.onChange('rtsp', 'udp')} />
								</Grid>
								<Grid item xs={12}>
									<TextField
										variant="outlined"
										type="number"
										min="0"
										step="1"
										fullWidth
										label={<Trans>Socket timeout (microseconds)</Trans>}
										value={settings.rtsp.stimeout}
										onChange={props.onChange('rtsp', 'stimeout')}
									/>
								</Grid>
							</React.Fragment>
						)}
						{protocolClass === 'http' && (
							<React.Fragment>
								<Grid item xs={12}>
									<Typography variant="h3">
										<Trans>HTTP and HTTPS</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Checkbox
										label={<Trans>Read input at native speed</Trans>}
										checked={settings.http.readNative}
										onChange={props.onChange('http', 'readNative')}
									/>
									<Checkbox
										label={<Trans>Force input framerate</Trans>}
										checked={settings.http.forceFramerate}
										onChange={props.onChange('http', 'forceFramerate')}
									/>
								</Grid>
								{settings.http.forceFramerate === true && (
									<Grid item xs={12}>
										<TextField
											variant="outlined"
											type="number"
											min="0"
											step="1"
											fullWidth
											label={<Trans>Framerate</Trans>}
											value={settings.http.framerate}
											onChange={props.onChange('http', 'framerate')}
										/>
									</Grid>
								)}
								<Grid item xs={12}>
									<TextField
										variant="outlined"
										fullWidth
										label="User-Agent"
										value={settings.http.userAgent}
										onChange={props.onChange('http', 'userAgent')}
									/>
								</Grid>
								<Grid item xs={12}>
									<TextField
										variant="outlined"
										fullWidth
										label="Referrer"
										value={settings.http.referer}
										onChange={props.onChange('http', 'referer')}
									/>
								</Grid>
								<Grid item xs={12}>
									<TextField
										variant="outlined"
										fullWidth
										label="HTTP proxy"
										value={settings.http.http_proxy}
										onChange={props.onChange('http', 'http_proxy')}
										placeholder="https://123.123.123.123:443"
									/>
								</Grid>
							</React.Fragment>
						)}
						<Grid item xs={12}>
							<Typography variant="h3">
								<Trans>General</Trans>
							</Typography>
						</Grid>
						<Grid item xs={12}>
							<TextField
								variant="outlined"
								type="number"
								min="0"
								step="1"
								fullWidth
								label="thread_queue_size"
								value={settings.general.thread_queue_size}
								onChange={props.onChange('general', 'thread_queue_size')}
							/>
						</Grid>
						<Grid item xs={12}>
							<TextField
								variant="outlined"
								type="number"
								min="32"
								step="1"
								fullWidth
								label="probesize (bytes)"
								value={settings.general.probesize}
								onChange={props.onChange('general', 'probesize')}
							/>
							<Typography variant="caption">
								<Trans>
									Mininum {32}, default {5000000}
								</Trans>
							</Typography>
						</Grid>
						<Grid item xs={12}>
							<TextField
								variant="outlined"
								type="number"
								min="1"
								step="1"
								fullWidth
								label="max_probe_packets"
								value={settings.general.max_probe_packets}
								onChange={props.onChange('general', 'max_probe_packets')}
							/>
							<Typography variant="caption">
								<Trans>Default {2500}</Trans>
							</Typography>
						</Grid>
						<Grid item xs={12}>
							<TextField
								variant="outlined"
								type="number"
								min="1"
								step="1"
								fullWidth
								label="analyzeduration (microseconds)"
								value={
									settings.mode === 'push'
										? settings.push.type === 'hls'
											? settings.general.analyzeduration_http
											: settings.push.type === 'rtmp'
												? settings.general.analyzeduration_rtmp
												: settings.general.analyzeduration
										: protocolClass === 'http'
											? settings.general.analyzeduration_http
											: protocolClass === 'rtmp'
												? settings.general.analyzeduration_rtmp
												: settings.general.analyzeduration
								}
								onChange={props.onChange(
									'general',
									settings.mode === 'push'
										? settings.push.type === 'hls'
											? 'analyzeduration_http'
											: settings.push.type === 'rtmp'
												? 'analyzeduration_rtmp'
												: 'analyzeduration'
										: protocolClass === 'http'
											? 'analyzeduration_http'
											: protocolClass === 'rtmp'
												? 'analyzeduration_rtmp'
												: 'analyzeduration',
								)}
							/>
							<Typography variant="caption">
								<Trans>
									Default {5000000} ({5} seconds)
								</Trans>
							</Typography>
						</Grid>
						<Grid item xs={12}>
							<MultiSelect type="select" label="flags" value={settings.general.fflags} onChange={props.onChange('general', 'fflags')}>
								<MultiSelectOption value="discardcorrupt" name="discardcorrupt" />
								<MultiSelectOption value="fastseek" name="fastseek" />
								<MultiSelectOption value="genpts" name="genpts" />
								<MultiSelectOption value="igndts" name="igndts" />
								<MultiSelectOption value="ignidx" name="ignidx" />
								<MultiSelectOption value="nobuffer" name="nobuffer" />
								<MultiSelectOption value="nofillin" name="nofillin" />
								<MultiSelectOption value="noparse" name="noparse" />
								<MultiSelectOption value="sortdts" name="sortdts" />
							</MultiSelect>
						</Grid>
						<Grid item xs={12}>
							<Checkbox label={<Trans>copyts</Trans>} checked={settings.general.copyts} onChange={props.onChange('general', 'copyts')} />
							<Checkbox
								label={<Trans>start_at_zero</Trans>}
								checked={settings.general.start_at_zero}
								onChange={props.onChange('general', 'start_at_zero')}
							/>
							<Checkbox
								label={<Trans>use_wallclock_as_timestamps</Trans>}
								checked={settings.general.use_wallclock_as_timestamps}
								onChange={props.onChange('general', 'use_wallclock_as_timestamps')}
							/>
						</Grid>
						<Grid item xs={12}>
							<Select
								type="select"
								label={<Trans>avoid_negative_ts</Trans>}
								value={settings.general.avoid_negative_ts}
								onChange={props.onChange('general', 'avoid_negative_ts')}
							>
								<MenuItem value="make_non_negative">make_non_negative</MenuItem>
								<MenuItem value="make_zero">make_zero</MenuItem>
								<MenuItem value="auto">auto (default)</MenuItem>
								<MenuItem value="disabled">disabled</MenuItem>
							</Select>
						</Grid>
					</Grid>
				</AccordionDetails>
			</Accordion>
		</Grid>
	);
}

function Pull(props) {
	const classes = useStyles();
	const settings = props.settings;
	const authProtocol = isAuthProtocol(settings.address);
	const validURL = isValidURL(settings.address);
	const supportedProtocol = isSupportedProtocol(settings.address, props.skills.protocols.input);

	return (
		<Grid container alignItems="flex-start" spacing={2} className={classes.gridContainer}>
			<Grid item xs={12}>
				<Typography>
					<Trans>Enter the address of your network source:</Trans>
				</Typography>
			</Grid>
			<Grid item xs={12}>
				<TextField
					variant="outlined"
					fullWidth
					label={<Trans>Address</Trans>}
					placeholder="rtsp://ip:port/path"
					value={settings.address}
					onChange={props.onChange('', 'address')}
				/>
				<Typography variant="caption">
					<Trans>Supports HTTP (HLS, DASH), RTP, RTSP, RTMP, SRT and more.</Trans>
				</Typography>
			</Grid>
			{validURL === true && (
				<React.Fragment>
					{!supportedProtocol ? (
						<Grid item xs={12} align="center">
							<BoxText color="dark">
								<WarningIcon fontSize="large" color="error" />
								<Typography>
									<Trans>This protocol is unknown or not supported by the available FFmpeg binary.</Trans>
								</Typography>
							</BoxText>
						</Grid>
					) : (
						<React.Fragment>
							{authProtocol && (
								<React.Fragment>
									<Grid item md={6} xs={12}>
										<TextField
											variant="outlined"
											fullWidth
											label={<Trans>Username</Trans>}
											value={settings.username}
											onChange={props.onChange('', 'username')}
										/>
										<Typography variant="caption">
											<Trans>Username for the device.</Trans>
										</Typography>
									</Grid>
									<Grid item md={6} xs={12}>
										<Password
											variant="outlined"
											fullWidth
											label={<Trans>Password</Trans>}
											value={settings.password}
											onChange={props.onChange('', 'password')}
										/>
										<Typography variant="caption">
											<Trans>Password for the device.</Trans>
										</Typography>
									</Grid>
								</React.Fragment>
							)}
							<AdvancedSettings {...props}></AdvancedSettings>
						</React.Fragment>
					)}
				</React.Fragment>
			)}
			<Grid item xs={12}>
				<FormInlineButton disabled={!validURL || !supportedProtocol} onClick={props.onProbe}>
					<Trans>Probe</Trans>
				</FormInlineButton>
			</Grid>
		</Grid>
	);
}

function Push(props) {
	const classes = useStyles();
	const settings = props.settings;

	//const supportsHLS = isSupportedProtocol('http://', props.skills.protocols.input);
	const supportsRTMP = isSupportedProtocol('rtmp://', props.skills.protocols.input);
	const supportsSRT = isSupportedProtocol('srt://', props.skills.protocols.input);

	if (!supportsRTMP && !supportsSRT) {
		return (
			<Grid container alignItems="flex-start" spacing={2} className={classes.gridContainer}>
				<Grid item xs={12} align="center">
					<BoxText color="dark">
						<WarningIcon fontSize="large" color="error" />
						<Typography>
							<Trans>The available FFmpeg binary doesn't support any of the required protocols.</Trans>
						</Typography>
					</BoxText>
				</Grid>
			</Grid>
		);
	}

	return (
		<React.Fragment>
			<Grid container alignItems="flex-start" spacing={2} className={classes.gridContainer}>
				<Grid item xs={12}>
					<Select type="select" label={<Trans>Protocol</Trans>} value={settings.push.type} onChange={props.onChange('push', 'type')}>
						<MenuItem value="rtmp" disabled={!supportsRTMP}>
							RTMP
						</MenuItem>
						<MenuItem value="srt" disabled={!supportsSRT}>
							SRT
						</MenuItem>
					</Select>
				</Grid>
			</Grid>
			{settings.push.type === 'rtmp' && <PushRTMP {...props} />}
			{settings.push.type === 'hls' && <PushHLS {...props} />}
			{settings.push.type === 'srt' && <PushSRT {...props} />}
		</React.Fragment>
	);
}

Push.defaultProps = {
	knownDevices: [],
	settings: {},
	config: {},
	skills: null,
	onChange: function (settings) {},
	onProbe: function (settings, inputs) {},
	onRefresh: function () {},
};

function PushHLS(props) {
	const classes = useStyles();
	const config = props.config;

	const HLS = getHLS(config);

	return (
		<Grid container alignItems="flex-start" spacing={2} className={classes.gridContainer}>
			<Grid item xs={12}>
				<Typography>
					<Trans>Send stream to this address:</Trans>
				</Typography>
			</Grid>
			<Grid item xs={12}>
				<BoxTextarea>
					<Textarea rows={1} value={HLS} readOnly allowCopy />
				</BoxTextarea>
			</Grid>
			<AdvancedSettings {...props} />
			<Grid item xs={12}>
				<FormInlineButton onClick={props.onProbe}>
					<Trans>Probe</Trans>
				</FormInlineButton>
			</Grid>
		</Grid>
	);
}

function PushRTMP(props) {
	const { i18n } = useLingui();
	const classes = useStyles();
	const navigate = useNavigate();
	const config = props.config;

	let form = null;

	if (config.rtmp.enabled === false) {
		form = (
			<Grid container alignItems="flex-start" spacing={2} className={classes.gridContainer}>
				<Grid item xs={12}>
					<Typography>
						<Trans>RTMP server is not enabled</Trans>
					</Typography>
				</Grid>
				<Grid item xs={12}>
					<Button variant="outlined" size="large" fullWidth color="primary" onClick={() => navigate('/settings/rtmp')}>
						<Trans>Enable RTMP server ...</Trans>
					</Button>
				</Grid>
			</Grid>
		);
	} else {
		const RTMP = getRTMP(config);

		const filteredDevices = props.knownDevices.filter((device) => device.media === 'rtmp');
		const options = filteredDevices.map((device) => {
			return (
				<MenuItem key={device.id} value={device.id}>
					{device.name}
				</MenuItem>
			);
		});

		options.unshift(
			<MenuItem key="none" value="none" disabled>
				{i18n._(t`Choose an input stream ...`)}
			</MenuItem>,
		);

		options.push(
			<MenuItem key={config.channelid} value={config.channelid}>
				{i18n._(t`Send stream to address ...`)}
			</MenuItem>,
		);

		form = (
			<Grid container alignItems="flex-start" spacing={2} className={classes.gridContainer}>
				<Grid item xs={12}>
					<Select type="select" label={<Trans>Input stream</Trans>} value={props.settings.push.name} onChange={props.onChange('push', 'name')}>
						{options}
					</Select>
					<Button size="small" startIcon={<RefreshIcon />} onClick={props.onRefresh} sx={{ float: 'right' }}>
						<Trans>Refresh</Trans>
					</Button>
				</Grid>
				{props.settings.push.name === config.channelid && (
					<React.Fragment>
						<Grid item xs={12}>
							<Typography>
								<Trans>Address:</Trans>
							</Typography>
						</Grid>
						<Grid item xs={12}>
							<BoxTextarea>
								<Textarea rows={1} value={RTMP} readOnly allowCopy />
							</BoxTextarea>
						</Grid>
					</React.Fragment>
				)}
				<AdvancedSettings {...props} />
				<Grid item xs={12}>
					<FormInlineButton onClick={props.onProbe} disabled={props.settings.push.name === 'none'}>
						<Trans>Probe</Trans>
					</FormInlineButton>
				</Grid>
			</Grid>
		);
	}

	return form;
}

PushRTMP.defaultProps = {
	knownDevices: [],
	settings: {},
	config: {},
	skills: null,
	onChange: function (settings) {},
	onProbe: function (settings, inputs) {},
	onRefresh: function () {},
};

function PushSRT(props) {
	const { i18n } = useLingui();
	const classes = useStyles();
	const navigate = useNavigate();
	const config = props.config;

	let form = null;

	if (config.srt.enabled === false) {
		form = (
			<Grid container alignItems="flex-start" spacing={2} className={classes.gridContainer}>
				<Grid item xs={12}>
					<Typography>
						<Trans>SRT server is not enabled</Trans>
					</Typography>
				</Grid>
				<Grid item xs={12}>
					<Button variant="outlined" size="large" fullWidth color="primary" onClick={() => navigate('/settings/srt')}>
						<Trans>Enable SRT server ...</Trans>
					</Button>
				</Grid>
			</Grid>
		);
	} else {
		const SRT = getSRT(config);

		const filteredDevices = props.knownDevices.filter((device) => device.media === 'srt');
		const options = filteredDevices.map((device) => {
			return (
				<MenuItem key={device.id} value={device.id}>
					{device.name}
				</MenuItem>
			);
		});

		options.unshift(
			<MenuItem key="none" value="none" disabled>
				{i18n._(t`Choose an input stream ...`)}
			</MenuItem>,
		);

		options.push(
			<MenuItem key={config.channelid} value={config.channelid}>
				{i18n._(t`Send stream to address ...`)}
			</MenuItem>,
		);

		form = (
			<Grid container alignItems="flex-start" spacing={2} className={classes.gridContainer}>
				<Grid item xs={12}>
					<Select type="select" label={<Trans>Input stream</Trans>} value={props.settings.push.name} onChange={props.onChange('push', 'name')}>
						{options}
					</Select>
					<Button size="small" startIcon={<RefreshIcon />} onClick={props.onRefresh} sx={{ float: 'right' }}>
						<Trans>Refresh</Trans>
					</Button>
				</Grid>
				{props.settings.push.name === config.channelid && (
					<React.Fragment>
						<Grid item xs={12}>
							<Typography>
								<Trans>Address:</Trans>
							</Typography>
						</Grid>
						<Grid item xs={12}>
							<BoxTextarea>
								<Textarea rows={1} value={SRT} readOnly allowCopy />
							</BoxTextarea>
						</Grid>
					</React.Fragment>
				)}
				<AdvancedSettings {...props} />
				<Grid item xs={12}>
					<FormInlineButton onClick={props.onProbe} disabled={props.settings.push.name === 'none'}>
						<Trans>Probe</Trans>
					</FormInlineButton>
				</Grid>
			</Grid>
		);
	}

	return form;
}

PushSRT.defaultProps = {
	knownDevices: [],
	settings: {},
	config: {},
	skills: null,
	onChange: function (settings) {},
	onProbe: function (settings, inputs) {},
	onRefresh: function () {},
};

function Source(props) {
	const classes = useStyles();
	const { i18n } = useLingui();
	const config = initConfig(props.config);
	const settings = initSettings(props.settings, config);
	const skills = initSkills(props.skills);

	const handleChange = (section, what) => (event) => {
		const value = event.target.value;

		if (section === 'http') {
			if (['readNative', 'forceFramerate'].includes(what)) {
				settings.http[what] = !settings.http[what];
			} else {
				settings.http[what] = value;
			}
		} else if (section === 'rtsp') {
			if (['udp'].includes(what)) {
				settings.rtsp[what] = !settings.rtsp[what];
			} else {
				settings.rtsp[what] = value;
			}
		} else if (section === 'general') {
			if (['copyts', 'start_at_zero', 'use_wallclock_as_timestamps'].includes(what)) {
				settings.general[what] = !settings.general[what];
			} else {
				settings.general[what] = value;
			}
		} else if (section === 'push') {
			settings.push[what] = value;
			if (what === 'type') {
				settings.push.name = config.channelid;
			}
		} else {
			settings[what] = value;
		}

		props.onChange({
			...settings,
		});
	};

	const handleProbe = () => {
		props.onProbe(settings, createInputs(settings, config, skills));
	};

	const handleRefresh = () => {
		props.onRefresh();
	};

	return (
		<React.Fragment>
			<Grid container alignItems="flex-start" spacing={2} className={classes.gridContainer}>
				<Grid item xs={12}>
					<Select type="select" label={<Trans>Pull or recieve the data:</Trans>} value={settings.mode} onChange={handleChange('', 'mode')}>
						<MenuItem value="pull">{i18n._(t`Pull Mode`)}</MenuItem>
						<MenuItem value="push">{i18n._(t`Receive Mode`)}</MenuItem>
					</Select>
				</Grid>
			</Grid>
			{settings.mode === 'pull' ? (
				<Pull settings={settings} config={config} skills={skills} onChange={handleChange} onProbe={handleProbe} />
			) : (
				<Push
					settings={settings}
					config={config}
					skills={skills}
					knownDevices={props.knownDevices}
					onChange={handleChange}
					onProbe={handleProbe}
					onRefresh={handleRefresh}
				/>
			)}
		</React.Fragment>
	);
}

Source.defaultProps = {
	knownDevices: [],
	settings: {},
	config: {},
	skills: null,
	onChange: function (settings) {},
	onProbe: function (settings, inputs) {},
};

function SourceIcon(props) {
	return <Icon style={{ color: '#FFF' }} {...props} />;
}

const id = 'network';
const name = <Trans>Network source</Trans>;
const capabilities = ['audio', 'video'];
const ffversion = '^4.1.0 || ^5.0.0 || ^6.1.0';

const func = {
	initSettings,
	initConfig,
	initSkills,
	createInputs,
	getProtocolClass,
	getHLS,
	getRTMP,
	getSRT,
	isValidURL,
	isAuthProtocol,
	isSupportedProtocol,
};

export { id, name, capabilities, ffversion, SourceIcon as icon, Source as component, func };
