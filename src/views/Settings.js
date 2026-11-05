import React from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import PropTypes from 'prop-types';

import { useLingui } from '@lingui/react';
import { Trans, t } from '@lingui/macro';
import makeStyles from '@mui/styles/makeStyles';
import Backdrop from '@mui/material/Backdrop';
import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import Divider from '@mui/material/Divider';
import Grid from '@mui/material/Grid';
import Stack from '@mui/material/Stack';
import LinearProgress from '@mui/material/LinearProgress';
import Link from '@mui/material/Link';
import MenuItem from '@mui/material/MenuItem';
import Tab from '@mui/material/Tab';
import Tabs from '@mui/material/Tabs';
import Typography from '@mui/material/Typography';
import WarningIcon from '@mui/icons-material/Warning';

import settingsImage from '../assets/images/settings.png';
import BoxText from '../misc/BoxText';
import BoxTextarea from '../misc/BoxTextarea';
import Checkbox from '../misc/Checkbox';
import Dialog from '../misc/modals/Dialog';
import Env from '../misc/Env';
import H from '../utils/help';
import NotifyContext from '../contexts/Notify';
import Paper from '../misc/Paper';
import PaperHeader from '../misc/PaperHeader';
import PaperContent from '../misc/PaperContent';
import PaperFooter from '../misc/PaperFooter';
import PaperThumb from '../misc/PaperThumb';
import Password from '../misc/Password';
import Select from '../misc/Select';
import TabPanel from '../misc/TabPanel';
import TabsVerticalGrid from '../misc/TabsVerticalGrid';
import Textarea from '../misc/Textarea';
import TextField from '../misc/TextField';
import useInterval from '../hooks/useInterval';

const useStyles = makeStyles((theme) => ({
	inlineEnv: {
		float: 'right',
	},
}));

const configValues = {
	update_check: {
		tab: 'general',
		set: (config, value) => {
			config.update_check = !config.update_check;
		},
		unset: (config) => {
			delete config.update_check;
		},
		validate: (config) => {
			return null;
		},
	},
	id: {
		tab: 'service',
		set: (config, value) => {
			config.id = value;
		},
		unset: (config) => {
			delete config.id;
		},
		validate: (config) => {
			return null;
		},
	},
	name: {
		tab: 'service',
		set: (config, value) => {
			config.name = value;
		},
		unset: (config) => {
			delete config.name;
		},
		validate: (config) => {
			return null;
		},
	},
	'log.level': {
		tab: 'logging',
		set: (config, value) => {
			config.log.level = value;
		},
		unset: (config) => {
			delete config.log.level;
		},
		validate: (config) => {
			return null;
		},
	},
	'log.max_lines': {
		tab: 'logging',
		set: (config, value) => {
			config.log.max_lines = value;
		},
		unset: (config) => {
			delete config.log.max_lines;
		},
		validate: (config) => {
			return null;
		},
	},
	address: {
		tab: 'network',
		set: (config, value) => {
			config.address = value;
		},
		unset: (config) => {
			delete config.address;
		},
		validate: (config) => {
			return null;
		},
	},
	'host.name': {
		tab: 'network',
		set: (config, value) => {
			config.host.name = value;
		},
		unset: (config) => {
			delete config.host.name;
		},
		validate: (config) => {
			if (config.tls.auto === false) {
				return null;
			}

			const names = toArray(config.host.name, ',').filter((name) => {
				if (name.match(/([0-9]+\.)+/) !== null) {
					return true;
				}

				if (name.match(/[:/]/) !== null) {
					return true;
				}

				return false;
			});

			if (names.length !== 0) {
				return 'Only domains names are allowed';
			}

			return null;
		},
	},
	'tls.address': {
		tab: 'network',
		set: (config, value) => {
			config.tls.address = value;
		},
		unset: (config) => {
			delete config.tls.address;
		},
		validate: (config) => {
			return null;
		},
	},
	'tls.auto': {
		tab: 'network',
		set: (config, value) => {
			config.tls.auto = !config.tls.auto;
		},
		unset: (config) => {
			delete config.tls.auto;
		},
		validate: (config) => {
			return null;
		},
	},
	'tls.email': {
		tab: 'network',
		set: (config, value) => {
			config.tls.email = value;
		},
		unset: (config) => {
			delete config.tls.email;
		},
		validate: (config) => {
			return null;
		},
	},
	'api.auth.enable': {
		tab: 'auth',
		set: (config, value) => {
			config.api.auth.enable = !config.api.auth.enable;
		},
		unset: (config) => {
			delete config.api.auth.enable;
		},
		validate: (config) => {
			return null;
		},
	},
	'api.auth.username': {
		tab: 'auth',
		set: (config, value) => {
			config.api.auth.username = value;
		},
		unset: (config) => {
			delete config.api.auth.username;
		},
		validate: (config) => {
			return null;
		},
	},
	'api.auth.password': {
		tab: 'auth',
		set: (config, value) => {
			config.api.auth.password = value;
		},
		unset: (config) => {
			delete config.api.auth.password;
		},
		validate: (config) => {
			return null;
		},
	},
	'rtmp.enable': {
		tab: 'rtmp',
		set: (config, value) => {
			config.rtmp.enable = !config.rtmp.enable;
		},
		unset: (config) => {
			delete config.rtmp.enable;
		},
		validate: (config) => {
			return null;
		},
	},
	'rtmp.enable_tls': {
		tab: 'rtmp',
		set: (config, value) => {
			config.rtmp.enable_tls = !config.rtmp.enable_tls;
		},
		unset: (config) => {
			delete config.rtmp.enable_tls;
		},
		validate: (config) => {
			return null;
		},
	},
	'rtmp.address': {
		tab: 'rtmp',
		set: (config, value) => {
			config.rtmp.address = value;
		},
		unset: (config) => {
			delete config.rtmp.address;
		},
		validate: (config) => {
			return null;
		},
	},
	'rtmp.address_tls': {
		tab: 'rtmp',
		set: (config, value) => {
			config.rtmp.address_tls = value;
		},
		unset: (config) => {
			delete config.rtmp.address_tls;
		},
		validate: (config) => {
			return null;
		},
	},
	'rtmp.app': {
		tab: 'rtmp',
		set: (config, value) => {
			config.rtmp.app = value;
		},
		unset: (config) => {
			delete config.rtmp.app;
		},
		validate: (config) => {
			return null;
		},
	},
	'rtmp.token': {
		tab: 'rtmp',
		set: (config, value) => {
			config.rtmp.token = value;
		},
		unset: (config) => {
			delete config.rtmp.token;
		},
		validate: (config) => {
			return null;
		},
	},
	'ffmpeg.log.max_lines': {
		tab: 'logging',
		set: (config, value) => {
			config.ffmpeg.log.max_lines = value;
		},
		unset: (config) => {
			delete config.ffmpeg.log.max_lines;
		},
		validate: (config) => {
			return null;
		},
	},
	'ffmpeg.log.max_history': {
		tab: 'logging',
		set: (config, value) => {
			config.ffmpeg.log.max_history = value;
		},
		unset: (config) => {
			delete config.ffmpeg.log.max_history;
		},
		validate: (config) => {
			return null;
		},
	},
	'srt.enable': {
		tab: 'srt',
		set: (config, value) => {
			config.srt.enable = !config.srt.enable;
		},
		unset: (config) => {
			delete config.srt.enable;
		},
		validate: (config) => {
			return null;
		},
	},
	'srt.address': {
		tab: 'srt',
		set: (config, value) => {
			config.srt.address = value;
		},
		unset: (config) => {
			delete config.srt.address;
		},
		validate: (config) => {
			return null;
		},
	},
	'srt.passphrase': {
		tab: 'srt',
		set: (config, value) => {
			config.srt.passphrase = value;
		},
		unset: (config) => {
			delete config.srt.passphrase;
		},
		validate: (config) => {
			return null;
		},
	},
	'srt.token': {
		tab: 'srt',
		set: (config, value) => {
			config.srt.token = value;
		},
		unset: (config) => {
			delete config.srt.token;
		},
		validate: (config) => {
			return null;
		},
	},
	'storage.cors.allow_all': {
		tab: 'storage',
		set: (config, value) => {
			config.storage.cors.allow_all = !config.storage.cors.allow_all;
		},
		unset: (config) => {
			return;
		},
		validate: (config) => {
			return null;
		},
	},
	'storage.cors.origins': {
		tab: 'storage',
		set: (config, value) => {
			config.storage.cors.origins = value;
		},
		unset: (config) => {
			delete config.storage.cors.origins;
		},
		validate: (config) => {
			return null;
		},
	},
	'storage.disk.max_size_mbytes': {
		tab: 'storage',
		set: (config, value) => {
			config.storage.disk.max_size_mbytes = value;
		},
		unset: (config) => {
			delete config.storage.disk.max_size_mbytes;
		},
		validate: (config) => {
			return null;
		},
	},
	'storage.disk.cache.enable': {
		tab: 'storage',
		set: (config, value) => {
			config.storage.disk.cache.enable = !config.storage.disk.cache.enable;
		},
		unset: (config) => {
			delete config.storage.disk.cache.enable;
		},
		validate: (config) => {
			return null;
		},
	},
	'storage.disk.cache.max_size_mbytes': {
		tab: 'storage',
		set: (config, value) => {
			config.storage.disk.cache.max_size_mbytes = value;
		},
		unset: (config) => {
			delete config.storage.disk.cache.max_size_mbytes;
		},
		validate: (config) => {
			return null;
		},
	},
	'storage.disk.cache.ttl_seconds': {
		tab: 'storage',
		set: (config, value) => {
			config.storage.disk.cache.ttl_seconds = value;
		},
		unset: (config) => {
			delete config.storage.disk.cache.ttl_seconds;
		},
		validate: (config) => {
			return null;
		},
	},
	'storage.disk.cache.max_file_size_mbytes': {
		tab: 'storage',
		set: (config, value) => {
			config.storage.disk.cache.max_file_size_mbytes = value;
		},
		unset: (config) => {
			delete config.storage.disk.cache.max_file_size_mbytes;
		},
		validate: (config) => {
			return null;
		},
	},
	'storage.disk.cache.types.allow': {
		tab: 'storage',
		set: (config, value) => {
			config.storage.disk.cache.types.allow = value;
		},
		unset: (config) => {
			delete config.storage.disk.cache.types.allow;
		},
		validate: (config) => {
			return null;
		},
	},
	'storage.disk.cache.types.block': {
		tab: 'storage',
		set: (config, value) => {
			config.storage.disk.cache.types.block = value;
		},
		unset: (config) => {
			delete config.storage.disk.cache.types.block;
		},
		validate: (config) => {
			return null;
		},
	},
	'storage.memory.auth.enable': {
		tab: 'storage',
		set: (config, value) => {
			config.storage.memory.auth.enable = !config.storage.memory.auth.enable;
		},
		unset: (config) => {
			delete config.storage.memory.auth.enable;
		},
		validate: (config) => {
			return null;
		},
	},
	'storage.memory.auth.username': {
		tab: 'storage',
		set: (config, value) => {
			config.storage.memory.auth.username = value;
		},
		unset: (config) => {
			delete config.storage.memory.auth.username;
		},
		validate: (config) => {
			return null;
		},
	},
	'storage.memory.auth.password': {
		tab: 'storage',
		set: (config, value) => {
			config.storage.memory.auth.password = value;
		},
		unset: (config) => {
			delete config.storage.memory.auth.password;
		},
		validate: (config) => {
			return null;
		},
	},
	'storage.memory.max_size_mbytes': {
		tab: 'storage',
		set: (config, value) => {
			config.storage.memory.max_size_mbytes = value;
		},
		unset: (config) => {
			delete config.storage.memory.max_size_mbytes;
		},
		validate: (config) => {
			return null;
		},
	},
	'storage.memory.purge': {
		tab: 'storage',
		set: (config, value) => {
			config.storage.memory.purge = !config.storage.memory.purge;
		},
		unset: (config) => {
			delete config.storage.memory.purge;
		},
		validate: (config) => {
			return null;
		},
	},
	'sessions.enable': {
		tab: 'playback',
		set: (config, value) => {
			config.sessions.enable = !config.sessions.enable;
		},
		unset: (config) => {
			delete config.sessions.enable;
		},
		validate: (config) => {
			return null;
		},
	},
	'sessions.ip_ignorelist': {
		tab: 'playback',
		set: (config, value) => {
			config.sessions.ip_ignorelist = value;
		},
		unset: (config) => {
			delete config.sessions.ip_ignorelist;
		},
		validate: (config) => {
			return null;
		},
	},
	'sessions.session_timeout_sec': {
		tab: 'playback',
		set: (config, value) => {
			config.sessions.session_timeout_sec = value;
		},
		unset: (config) => {
			delete config.sessions.session_timeout_sec;
		},
		validate: (config) => {
			return null;
		},
	},
	'sessions.persist': {
		tab: 'playback',
		set: (config, value) => {
			config.sessions.persist = !config.sessions.persist;
		},
		unset: (config) => {
			delete config.sessions.persist;
		},
		validate: (config) => {
			return null;
		},
	},
	'sessions.max_bitrate_mbit': {
		tab: 'network',
		set: (config, value) => {
			config.sessions.max_bitrate_mbit = value;
		},
		unset: (config) => {
			delete config.sessions.max_bitrate_mbit;
		},
		validate: (config) => {
			return null;
		},
	},
	'sessions.max_sessions': {
		tab: 'network',
		set: (config, value) => {
			config.sessions.max_sessions = value;
		},
		unset: (config) => {
			delete config.sessions.max_sessions;
		},
		validate: (config) => {
			return null;
		},
	},
	'service.enable': {
		tab: 'service',
		set: (config, value) => {
			config.service.enable = !config.service.enable;
		},
		unset: (config) => {
			delete config.service.enable;
		},
		validate: (config) => {
			return null;
		},
	},
	'service.token': {
		tab: 'service',
		set: (config, value) => {
			config.service.token = value;
		},
		unset: (config) => {
			delete config.service.token;
		},
		validate: (config) => {
			return null;
		},
	},
};

const RETRIES = 60;

function ErrorTab(props) {
	let { label, errors, ...other } = props;

	if (errors === true) {
		label = (
			<Stack direction="row" spacing={1} alignItems="center">
				<WarningIcon color="error" />
				<Typography>{label}</Typography>
			</Stack>
		);
	}

	return <Tab label={label} {...other} />;
}

function ErrorBox(props) {
	const messages = props.messages.filter((m) => props.configvalue === '' || m.configvalue === props.configvalue);

	if (messages.length === 0) {
		return null;
	}

	return (
		<BoxText color="danger">
			<Typography variant="body2" gutterBottom>
				{messages.map((e, i) => (
					<span key={i}>
						{e.error}
						<br />
					</span>
				))}
			</Typography>
		</BoxText>
	);
}

ErrorBox.defaultProps = {
	configvalue: '',
	messages: [],
};

const toArray = (val, separator) => {
	return val
		.split(separator)
		.map((l) => l.trim())
		.filter((l) => l.length !== 0);
};

const toInt = (val) => {
	if (typeof val === 'string') {
		return val.length === 0 ? 0 : parseInt(val);
	}

	return val;
};

export default function Settings(props) {
	const classes = useStyles();
	const { i18n } = useLingui();
	const navigate = useNavigate();
	const { tab: _tab } = useParams();
	const notify = React.useContext(NotifyContext);
	const [$config, setConfig] = React.useState({
		ready: false,
		modified: false,
		data: null,
		overrides: [],
		created_at: null,
		loaded_at: null,
		updated_at: null,
		outdated: false,
		core: '',
		service: false,
	});
	const [$expert, setExpert] = React.useState(props.restreamer.IsExpert());
	const [$updates, setUpdates] = React.useState({
		has: props.restreamer.HasUpdates(),
		want: props.restreamer.CheckForUpdates(),
	});
	const [$tab, setTab] = React.useState(_tab ? _tab : 'general');
	const [$tabs, setTabs] = React.useState({
		general: { errors: false, messages: [] }, // messages is an array of objects: {configvalue: '', error: ''}
		network: { errors: false, messages: [] },
		auth: { errors: false, messages: [] },
		playback: { errors: false, messages: [] },
		storage: { errors: false, messages: [] },
		rtmp: { errors: false, messages: [] },
		srt: { errors: false, messages: [] },
		logging: { errors: false, messages: [] },
		service: { errors: false, messages: [] },
	});
	const [$logdata, setLogdata] = React.useState('');
	const logTimer = React.useRef();
	const [$reloadKey, setReloadKey] = React.useState('');
	const [$dialogs, setDialogs] = React.useState({
		restart: false,
		saved: false,
	});
	const [$saving, setSaving] = React.useState(false);
	const [$restart, setRestart] = React.useState({
		restarting: false,
		timeout: false,
		changes: {
			tls: false,
			port: false,
		},
	});

	React.useEffect(() => {
		(async () => {
			await load();
		})();

		return () => {
			clearInterval(logTimer.current);
		};
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

	useInterval(() => {
		setUpdates({
			...$updates,
			has: props.restreamer.HasUpdates(),
		});
	}, 1000 * 2);

	const load = async () => {
		setReloadKey(props.restreamer.CreatedAt().toISOString());

		const data = await props.restreamer.Config();

		let config = null;
		let overrides = [];
		let outdated = false;

		if (data !== null) {
			config = JSON.parse(JSON.stringify(data.config));

			config.host.name = config.host.name.join(', ');

			config.storage.cors.allow_all = false;
			if (config.storage.cors.origins.length === 1 && config.storage.cors.origins[0] === '*') {
				config.storage.cors.allow_all = true;
			}
			config.storage.cors.origins = config.storage.cors.origins.join('\n');
			config.storage.disk.cache.types.allow = config.storage.disk.cache.types.allow.join('\n');
			config.storage.disk.cache.types.block = config.storage.disk.cache.types.block.join('\n');

			config.sessions.ip_ignorelist = config.sessions.ip_ignorelist.join('\n');

			if (config.tls.enable === false) {
				config.tls.auto = false;
			}

			config.address = config.address.split(':').join('');
			config.tls.address = config.tls.address.split(':').join('');
			config.rtmp.address = config.rtmp.address.split(':').join('');
			config.rtmp.address_tls = config.rtmp.address_tls.split(':').join('');
			config.srt.address = config.srt.address.split(':').join('');

			if (config.tls.auto === true) {
				config.tls.enable = true;
			}

			overrides = data.overrides;

			if (data.loaded_at < data.updated_at) {
				outdated = true;
			}
		}

		setConfig({
			...$config,
			ready: true,
			modified: false,
			data: config,
			overrides: overrides,
			outdated: outdated,
			core: '',
			service: props.restreamer.HasService(),
		});
	};

	const handleExpertMode = () => {
		props.restreamer.SetExpert(!$expert);
		setExpert(!$expert);
	};

	const handleCheckForUpdates = () => {
		props.restreamer.SetCheckForUpdates(!$updates.want);
		setUpdates({
			...$updates,
			want: !$updates.want,
		});
	};

	const handleChange = (what) => (event) => {
		const value = event.target.value;
		const config = $config.data;

		if (!(what in configValues)) {
			console.warn(`config value "${what}" is not a valid configValue`);
			return;
		}

		configValues[what].set(config, value);

		const tabs = $tabs;

		// Create tab if it doesn't exist
		const tab = configValues[what].tab;
		if (!(tab in tabs)) {
			tabs[tab] = {
				errors: false,
				messages: [],
			};
		}

		// Remove all error for a config value
		tabs[tab].messages = tabs[tab].messages.filter((m) => m.configvalue !== what);
		if (tabs[tab].messages.length === 0) {
			tabs[tab].errors = false;
		}

		// Validate current config value and add error if validation fails
		const err = configValues[what].validate(config);
		if (err !== null) {
			tabs[tab].errors = true;

			tabs[tab].messages.push({
				configvalue: what,
				error: err,
			});
		}

		setTabs({
			...$tabs,
			...tabs,
		});

		setConfig({
			...$config,
			data: config,
			modified: true,
		});
	};

	const handleCoreConfig = (event) => {
		const value = event.target.value;

		setConfig({
			...$config,
			core: value,
			modified: true,
		});
	};

	const handleChangeTab = async (event, value) => {
		if (value === 'logging') {
			await updateLogdata();

			logTimer.current = setInterval(async () => {
				await updateLogdata();
			}, 1000);
		} else {
			clearInterval(logTimer.current);
		}

		setTab(value);
	};

	const updateLogdata = async () => {
		const logdata = await props.restreamer.Log();
		setLogdata(logdata.join('\n'));
	};

	const handleSave = async () => {
		setSaving(true);

		let config = null;

		if ($config.core.length === 0) {
			// Poor man's deep clone, but that's OK because we have only basic data types
			config = JSON.parse(JSON.stringify($config.data));

			config.log.max_lines = toInt(config.log.max_lines);

			config.host.name = toArray(config.host.name, ',');

			config.storage.cors.origins = toArray(config.storage.cors.origins, '\n');
			if (config.storage.cors.allow_all === true) {
				config.storage.cors.origins = ['*'];
			}
			delete config.storage.cors.allow_all;

			config.storage.disk.cache.types.allow = toArray(config.storage.disk.cache.types.allow, '\n');
			config.storage.disk.cache.types.block = toArray(config.storage.disk.cache.types.block, '\n');

			config.sessions.ip_ignorelist = toArray(config.sessions.ip_ignorelist, '\n');

			config.ffmpeg.log.max_lines = toInt(config.ffmpeg.log.max_lines);
			config.ffmpeg.log.max_history = toInt(config.ffmpeg.log.max_history);
			config.storage.disk.max_size_mbytes = toInt(config.storage.disk.max_size_mbytes);
			config.storage.disk.cache.max_size_mbytes = toInt(config.storage.disk.cache.max_size_mbytes);
			config.storage.disk.cache.ttl_seconds = toInt(config.storage.disk.cache.ttl_seconds);
			config.storage.disk.cache.max_file_size_mbytes = toInt(config.storage.disk.cache.max_file_size_mbytes);
			config.storage.memory.max_size_mbytes = toInt(config.storage.memory.max_size_mbytes);
			config.sessions.session_timeout_sec = toInt(config.sessions.session_timeout_sec);
			config.sessions.max_bitrate_mbit = toInt(config.sessions.max_bitrate_mbit);
			config.sessions.max_sessions = toInt(config.sessions.max_sessions);

			config.address = ':' + config.address;
			config.tls.address = ':' + config.tls.address;
			config.rtmp.address = ':' + config.rtmp.address;
			config.rtmp.address_tls = ':' + config.rtmp.address_tls;
			config.rtmp.app = !config.rtmp.app.startsWith('/') ? '/' + config.rtmp.app : config.rtmp.app;
			config.srt.address = ':' + config.srt.address;

			if (config.tls.auto === true) {
				config.tls.enable = true;
			} else {
				config.tls.enable = false;
			}

			for (let what of $config.overrides) {
				if (!(what in configValues)) {
					continue;
				}

				configValues[what].unset(config);
			}
		} else {
			try {
				config = JSON.parse($config.core);
			} catch (e) {
				const tabs = {};
				tabs['service'] = {
					errors: true,
					messages: [
						{
							configvalue: 'coreconfig',
							error: e.message,
						},
					],
				};

				setTabs({
					...$tabs,
					...tabs,
				});

				setSaving(false);

				return;
			}
		}

		let outdated = $config.outdated;

		const [, err] = await props.restreamer.ConfigSet(config);
		if (err !== null) {
			if (err.code === 404) {
				notify.Dispatch('error', 'save:settings', i18n._(t`API endpoint not found. Settings not saved.`));
			} else if (err.code === 409) {
				notify.Dispatch('error', 'save:settings', i18n._(t`There were some errors in the settings. Settings not saved.`));
			} else {
				notify.Dispatch('error', 'save:settings', i18n._(t`There was a problem storing the settings. Settings not saved.`));
			}

			const tabs = {};
			const ucfirst = (string) => {
				return string.charAt(0).toUpperCase() + string.slice(1);
			};

			if (err.code === 409) {
				for (let errorfield in err.message) {
					if (!(errorfield in configValues)) {
						continue;
					}

					const tab = configValues[errorfield].tab;

					if (!(tab in tabs)) {
						tabs[tab] = {
							errors: true,
							messages: [],
						};
					}

					for (let m of err.message[errorfield]) {
						tabs[tab].messages.push({
							configvalue: errorfield,
							error: ucfirst(m),
						});
					}
				}
			} else if (err.code === 400) {
				tabs['service'] = {
					errors: true,
					messages: [
						{
							configvalue: 'coreconfig',
							error: err.message,
						},
					],
				};
			}

			setTabs({
				...$tabs,
				...tabs,
			});
		} else {
			notify.Dispatch('success', 'save:settings', i18n._(t`Settings saved. All changes will be applied after restarting the application.`));

			const tabs = {};

			for (let t in $tabs) {
				tabs[t] = { errors: false, messages: [] };
			}

			setTabs({
				...$tabs,
				...tabs,
			});

			outdated = true;

			handleSavedDialog();
		}

		setConfig({
			...$config,
			modified: false,
			outdated: outdated,
			core: err === null ? '' : $config.core,
		});

		setSaving(false);
	};

	const handleSavedDialog = () => {
		setDialogs({
			...$dialogs,
			saved: !$dialogs.saved,
		});
	};

	const handleReloadDialog = () => {
		setDialogs({
			...$dialogs,
			restart: !$dialogs.restart,
		});
	};

	const handleReload = async () => {
		setDialogs({
			...$dialogs,
			saved: false,
			restart: false,
		});

		setRestart({
			...$restart,
			restarting: true,
			timeout: false,
		});

		const res = await props.restreamer.ConfigReload();
		if (res === false) {
			notify.Dispatch('error', 'restart', i18n._(t`Restarting the application failed.`));

			setRestart({
				...$restart,
				restarting: false,
			});

			return;
		}

		const waitFor = (ms) => {
			return new Promise((resolve) => {
				setTimeout(resolve, ms);
			});
		};

		props.restreamer.IgnoreAPIErrors(true);

		let restarted = false;

		for (let retries = 0; retries <= RETRIES; retries++) {
			await waitFor(1000);

			let currentKey = $reloadKey;

			const about = await props.restreamer.About();
			if (about === null) {
				// API is not yet available
				continue;
			}

			if (about.id.length === 0) {
				// The API requires a login and after the restart the token
				// is not valid anymore. This means the restart happened.
				restarted = true;
				break;
			}

			let createdAt = about.created_at;
			if (createdAt !== null) {
				currentKey = createdAt.toISOString();
			}

			if ($reloadKey !== currentKey) {
				// The API tells us that it changed from before. This means
				// the restart happened.
				restarted = true;
				break;
			}
		}

		if (restarted === false) {
			setRestart({
				...$restart,
				restarting: true,
				timeout: true,
			});

			//notify.Dispatch('error', 'restart', i18n._(t`Restarting the application failed.`));

			return false;
		}

		await props.restreamer.Validate();
		await props.restreamer.Login($config.data.api.auth.username, $config.data.api.auth.password);

		window.location.reload();

		return true;
	};

	const handleAbort = () => {
		navigate(-1);
	};

	const handleHelp = (topic) => () => {
		H('settings-' + topic);
	};

	const env = (what) => {
		if ($config.overrides.indexOf(what) !== -1) {
			return true;
		}

		return false;
	};

	if ($config.ready === false) {
		return null;
	}

	if ($config.data === null) {
		return (
			<Paper xs={8} sm={6} md={6}>
				<PaperHeader title={<Trans>Error</Trans>} />
				<PaperContent>
					<Trans>Unable to load the config.</Trans>
				</PaperContent>
				<PaperFooter
					buttonsRight={
						<Button variant="outlined" color="primary" onClick={() => window.location.reload()}>
							<Trans>Retry</Trans>
						</Button>
					}
				/>
			</Paper>
		);
	}

	const config = $config.data;

	let title = <Trans>Settings</Trans>;
	if ($expert === true) {
		title = <Trans>Settings (expert mode)</Trans>;
	}

	return (
		<React.Fragment>
			<Paper xs={12} md={10}>
				<PaperHeader title={title} onAbort={handleAbort} onHelp={handleHelp($tab)} />
				<Grid container spacing={1}>
					<TabsVerticalGrid>
						<Tabs orientation="vertical" variant="scrollable" value={$tab} onChange={handleChangeTab}>
							<ErrorTab className="tab" label={<Trans>General</Trans>} value="general" errors={$tabs.general.errors} />
							{$config.service === true && (
								<ErrorTab className="tab" label={<Trans>Service</Trans>} value="service" errors={$tabs.service.errors} />
							)}
							<ErrorTab className="tab" label={<Trans>Network</Trans>} value="network" errors={$tabs.network.errors} />
							<ErrorTab className="tab" label={<Trans>Authorization</Trans>} value="auth" errors={$tabs.auth.errors} />
							{$expert === true && <ErrorTab className="tab" label={<Trans>Playback</Trans>} value="playback" errors={$tabs.playback.errors} />}
							{$expert === true && <ErrorTab className="tab" label={<Trans>Storage</Trans>} value="storage" errors={$tabs.storage.errors} />}
							<ErrorTab className="tab" label={<Trans>RTMP</Trans>} value="rtmp" errors={$tabs.rtmp.errors} />
							<ErrorTab className="tab" label={<Trans>SRT</Trans>} value="srt" errors={$tabs.srt.errors} />
							{$expert === true && <ErrorTab className="tab" label={<Trans>Logging</Trans>} value="logging" errors={$tabs.logging.errors} />}
						</Tabs>
						<TabPanel value={$tab} index="general" className="panel">
							<Grid container spacing={2}>
								<Grid item xs={12}>
									<Grid item xs={12}>
										<PaperThumb image={settingsImage} title="Welcome to Restreamer v2" height="200px" />
									</Grid>
								</Grid>
								<Grid item xs={12}>
									<Typography variant="h2">
										<Trans>General</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Typography variant="body1">
										<Trans>All important system settings.</Trans>
									</Typography>
								</Grid>
								{$tabs.general.errors && (
									<Grid item xs={12}>
										<ErrorBox messages={$tabs.general.messages} />
									</Grid>
								)}
								<Grid item xs={12}>
									<Checkbox label={<Trans>Check for updates</Trans>} checked={$updates.want} onChange={handleCheckForUpdates} />
									{$updates.has === true && (
										<BoxText color="success">
											<Typography variant="inherit" color="inherit" width="100%">
												<Link
													color="inherit"
													style={{ textDecoration: 'underline', cursor: 'pointer' }}
													onClick={handleHelp('update-link')}
													target="_blank"
												>
													<Trans>There are updates available. Here you get more information.</Trans>
												</Link>
											</Typography>
										</BoxText>
									)}
								</Grid>
								<Grid item xs={12}>
									<Checkbox
										label={<Trans>Send anonymous metrics (helps us for future development)</Trans>}
										checked={config.update_check}
										disabled={env('update_check')}
										onChange={handleChange('update_check')}
									/>
									{env('update_check') && <Env />}
									<ErrorBox configvalue="update_check" messages={$tabs.general.messages} />
								</Grid>
								<Grid item xs={12}>
									<Checkbox label={<Trans>Expert mode</Trans>} checked={$expert} onChange={handleExpertMode} />
								</Grid>
								<Grid item xs={12}>
									<Divider />
								</Grid>
								<Grid item xs={12}>
									<Button variant="outlined" color="secondary" onClick={handleReloadDialog}>
										<Trans>Restart</Trans>
									</Button>
									{$config.outdated === true && (
										<Typography variant="caption">
											<Trans>The application is using an older version of the settings.</Trans>
										</Typography>
									)}
								</Grid>
							</Grid>
						</TabPanel>
						<TabPanel value={$tab} index="service" className="panel">
							<Grid container spacing={2}>
								<Grid item xs={12}>
									<Typography variant="h2">
										<Trans>Service</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Typography variant="body1">
										<Trans>Setting for connection to the service.</Trans>{' '}
										<Link color="secondary" href="https://service.datarhei.com" target="_blank">
											<Trans>More about the service</Trans>
										</Link>
										.
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Checkbox
										label={<Trans>Restreamer Service</Trans>}
										checked={config.service.enable}
										disabled={env('service.enable')}
										onChange={handleChange('service.enable')}
									/>{' '}
									{env('service.enable') && <Env />}
									<ErrorBox configvalue="service.enable" messages={$tabs.service.messages} />
								</Grid>
								<Grid item xs={12}>
									<TextField label={<Trans>Node ID</Trans>} env={env('id')} disabled value={config.id} onChange={handleChange('id')} />
									<Typography variant="caption">
										<Trans>Unique ident on the service.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<TextField
										label={<Trans>Name</Trans>}
										env={env('name')}
										disabled={env('name') || !config.service.enable}
										value={config.name}
										onChange={handleChange('name')}
									/>
									<ErrorBox configvalue="name" messages={$tabs.service.messages} />
									<Typography variant="caption">
										<Trans>Human readable name on the service.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<TextField
										label={<Trans>Token</Trans>}
										env={env('service.token')}
										disabled={env('service.token') || !config.service.enable}
										value={config.service.token}
										onChange={handleChange('service.token')}
									/>
									<ErrorBox configvalue="service.token" messages={$tabs.service.messages} />
									<Typography variant="caption">
										<Trans>Service token for monitoring.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<TextField
										className={classes.root}
										variant="outlined"
										fullWidth
										multiline
										rows={10}
										disabled={!config.service.enable}
										label={<Trans>Config</Trans>}
										value={$config.core}
										onChange={handleCoreConfig}
									/>
									<ErrorBox configvalue="coreconfig" messages={$tabs.service.messages} />
									<Typography variant="caption">
										<Trans>Custom JSON config for datarhei Core.</Trans>
									</Typography>
								</Grid>
							</Grid>
						</TabPanel>
						<TabPanel value={$tab} index="network" className="panel">
							<Grid container spacing={2}>
								<Grid item xs={12}>
									<Typography variant="h2">
										<Trans>Network</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Typography variant="h3">
										<Trans>Address</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<TextField
										label={<Trans>Public domain/s</Trans>}
										env={env('host.name')}
										disabled={env('host.name')}
										value={config.host.name}
										onChange={handleChange('host.name')}
									/>
									<ErrorBox configvalue="host.name" messages={$tabs.network.messages} />
									<Typography variant="caption">
										<Trans>
											The public reachable domain name of the host this Restreamer is running on. Separate multiple domain names by a
											comma.
										</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12} md={6}>
									<TextField
										type="number"
										label={<Trans>HTTP port</Trans>}
										env={env('address')}
										disabled={env('address')}
										value={config.address}
										onChange={handleChange('address')}
									/>
									<ErrorBox configvalue="address" messages={$tabs.network.messages} />
									<Typography variant="caption">
										<Trans>Address to listen on for HTTP requests.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12} md={6}>
									<TextField
										type="number"
										label={<Trans>HTTPS port</Trans>}
										env={env('tls.address')}
										disabled={env('tls.address') || !config.tls.auto}
										value={config.tls.address}
										onChange={handleChange('tls.address')}
									/>
									<ErrorBox configvalue="tls.address" messages={$tabs.network.messages} />
									<Typography variant="caption">
										<Trans>Address to listen on for HTTPS requests.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Divider />
								</Grid>
								<Grid item xs={12}>
									<Typography variant="h3">
										<Trans>HTTPS (SSL/TLS)</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Checkbox
										label={<Trans>Let's Encrypt certification</Trans>}
										checked={config.tls.auto}
										disabled={env('tls.auto') || config.host.name.length === 0}
										onChange={handleChange('tls.auto')}
									/>{' '}
									{env('tls.auto') && <Env />}
									<ErrorBox configvalue="tls.auto" messages={$tabs.network.messages} />
									<Typography variant="caption">
										<Trans>Let's Encrypt requires one or more public domain names and an accessible port 80/TCP.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<TextField
										label={<Trans>Email</Trans>}
										env={env('tls.email')}
										disabled={env('tls.email') || !config.tls.auto}
										value={config.tls.email}
										onChange={handleChange('tls.email')}
									/>
									<ErrorBox configvalue="tls.email" messages={$tabs.network.messages} />
									<Typography variant="caption">
										<Trans>
											Please enter your email address to signify agreement with the Let's Encrypt CA's terms of service and to be notified
											in case of issues.
										</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Divider />
								</Grid>
								<Grid item xs={12}>
									<Typography variant="h3">
										<Trans>Bandwidth control</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12} md={12}>
									<TextField
										type="number"
										label={<Trans>Maximum viewers</Trans>}
										env={env('sessions.max_sessions')}
										disabled={env('sessions.max_sessions')}
										value={config.sessions.max_sessions}
										onChange={handleChange('sessions.max_sessions')}
									/>
									<ErrorBox configvalue="sessions.max_sessions" messages={$tabs.network.messages} />
									<Typography variant="caption">
										<Trans>
											Sets a viewer limit for HLS sessions. If the limit is exceeded, HLS viewers receive the HTTP status code 509
											(Bandwidth Limit Exceeded). 0 is unlimited.
										</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12} md={12}>
									<TextField
										type="number"
										label={<Trans>Maximum bandwidth Mbit/s</Trans>}
										env={env('sessions.max_bitrate_mbit')}
										disabled={env('sessions.max_bitrate_mbit')}
										value={config.sessions.max_bitrate_mbit}
										onChange={handleChange('sessions.max_bitrate_mbit')}
									/>
									<ErrorBox configvalue="sessions.max_bitrate_mbit" messages={$tabs.network.messages} />
									<Typography variant="caption">
										<Trans>
											Sets a bandwidth limit in Mbit per second for outgoing HLS data transfer. All services, such as RTMP and outgoing
											processes, are included in the calculation. If the bandwidth is exceeded, HLS viewers receive the HTTP status code
											509 (Bandwidth Limit Exceeded). 0 is unlimited.
										</Trans>
									</Typography>
								</Grid>
							</Grid>
						</TabPanel>
						<TabPanel value={$tab} index="auth" className="panel">
							<Grid container spacing={2}>
								<Grid item xs={12}>
									<Typography variant="h2">
										<Trans>Authorization</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Typography variant="h3">
										<Trans>Basic</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Checkbox
										label={<Trans>Login/JWT authorization</Trans>}
										checked={config.api.auth.enable}
										// prob: interface enforces auth.
										// disabled={env('api.auth.enable')}
										disabled
										onChange={handleChange('api.auth.enable')}
									/>{' '}
									{env('api.auth.enable') && <Env />}
									<ErrorBox configvalue="api.auth.enable" messages={$tabs.auth.messages} />
									<Typography variant="caption">
										<Trans>Enabling authorization is strongly advised. Otherwise, anybody can access this instance.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<TextField
										env={env('api.auth.username')}
										disabled={env('api.auth.username') || !config.api.auth.enable}
										label={<Trans>Username</Trans>}
										value={config.api.auth.username}
										onChange={handleChange('api.auth.username')}
									/>
									<ErrorBox configvalue="api.auth.username" messages={$tabs.auth.messages} />
									<Typography variant="caption">
										<Trans>Username for authorization.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Password
										env={env('api.auth.password')}
										disabled={env('api.auth.password') || !config.api.auth.enable}
										label={<Trans>Password</Trans>}
										value={config.api.auth.password}
										onChange={handleChange('api.auth.password')}
									/>
									<ErrorBox configvalue="api.auth.password" messages={$tabs.auth.messages} />
									<Typography variant="caption">
										<Trans>Password for authorization.</Trans>
									</Typography>
								</Grid>
							</Grid>
						</TabPanel>
						<TabPanel value={$tab} index="playback" className="panel">
							<Grid container spacing={2}>
								<Grid item xs={12}>
									<Typography variant="h2">
										<Trans>Playback</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Typography variant="h3">
										<Trans>Security</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Checkbox
										label={<Trans>Allow all referrer</Trans>}
										checked={config.storage.cors.allow_all}
										disabled={env('storage.cors.origins')}
										onChange={handleChange('storage.cors.allow_all')}
									/>{' '}
									{env('storage.cors.allow_all') && <Env />}
								</Grid>
								<Grid item xs={12}>
									<TextField
										multiline
										rows={5}
										label="CORS: Access-Control-Allow-Origin"
										value={config.storage.cors.origins}
										env={env('storage.cors.origins')}
										disabled={env('storage.cors.origins') || config.storage.cors.allow_all}
										onChange={handleChange('storage.cors.origins')}
									/>
									<ErrorBox configvalue="storage.cors.origins" messages={$tabs.playback.messages} />
									<Typography variant="caption">
										<Trans>One referrer per line, e.g. http://www.example.com</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Divider />
								</Grid>
								<Grid item xs={12}>
									<Typography variant="h3">
										<Trans>Statistics</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Checkbox
										label={<Trans>HLS statistic for the In-memory storage</Trans>}
										checked={config.sessions.enable}
										disabled={env('sessions.enable')}
										onChange={handleChange('sessions.enable')}
									/>{' '}
									{env('sessions.enable') && <Env />}
									<ErrorBox configvalue="sessions.enable" messages={$tabs.playback.messages} />
									<Typography variant="caption">
										<Trans>Allow counting how many viewers the stream has.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<TextField
										multiline
										rows={5}
										label={<Trans>Ignore IP ranges</Trans>}
										env={env('sessions.ip_ignorelist')}
										disabled={env('sessions.ip_ignorelist') || !config.sessions.enable}
										value={config.sessions.ip_ignorelist}
										onChange={handleChange('sessions.ip_ignorelist')}
									/>
									<ErrorBox configvalue="sessions.ip_ignorelist" messages={$tabs.playback.messages} />
									<Typography variant="caption">
										<Trans>
											List of IP ranges in CIDR notation, e.g., 127.0.0.1/32, that the statistics will not record—one IP range per line.
											Leave empty to record all sessions.
										</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<TextField
										type="number"
										label={<Trans>Maximum viewer idle time (Seconds)</Trans>}
										env={env('sessions.session_timeout_sec')}
										disabled={env('sessions.session_timeout_sec') || !config.sessions.enable}
										value={config.sessions.session_timeout_sec}
										onChange={handleChange('sessions.session_timeout_sec')}
									/>
									<ErrorBox configvalue="sessions.session_timeout_sec" messages={$tabs.playback.messages} />
									<Typography variant="caption">
										<Trans>Time until an inactive viewer connection is treated as closed.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Checkbox
										label={<Trans>Persist viewer statistics</Trans>}
										checked={config.sessions.persist}
										disabled={env('sessions.persist') || !config.sessions.enable}
										onChange={handleChange('sessions.persist')}
									/>{' '}
									{env('sessions.persist') && <Env />}
									<ErrorBox configvalue="sessions.persist" messages={$tabs.playback.messages} />
									<Typography variant="caption">
										<Trans>Stores the viewer statistics to the disk.</Trans>
									</Typography>
								</Grid>
							</Grid>
						</TabPanel>
						<TabPanel value={$tab} index="storage" className="panel">
							<Grid container spacing={2}>
								<Grid item xs={12}>
									<Typography variant="h2">
										<Trans>Storage</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Typography variant="h3">
										<Trans>In-memory</Trans>
									</Typography>
									<Typography variant="body1">
										<Trans>Settings for /memfs path.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Checkbox
										label={<Trans>Write protection</Trans>}
										checked={config.storage.memory.auth.enable}
										env={env('storage.memory.auth.enable')}
										disabled={env('storage.memory.auth.enable')}
										onChange={handleChange('storage.memory.auth.enable')}
									/>
									<ErrorBox configvalue="storage.memory.auth.enable" messages={$tabs.storage.messages} />
									<Typography variant="caption">
										<Trans>Enabling basic auth is strongly advised. Otherwise, anybody could write data to /memfs.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12} md={6}>
									<TextField
										label={<Trans>Username</Trans>}
										env={env('storage.memory.auth.username')}
										disabled={env('storage.memory.auth.username') || !config.storage.memory.auth.enable}
										value={config.storage.memory.auth.username}
										onChange={handleChange('storage.memory.auth.username')}
									/>
									<ErrorBox configvalue="storage.memory.auth.username" messages={$tabs.storage.messages} />
									<Typography variant="caption">
										<Trans>Username for authorization.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12} md={6}>
									<Password
										label={<Trans>Password</Trans>}
										env={env('storage.memory.auth.password')}
										disabled={env('storage.memory.auth.password') || !config.storage.memory.auth.enable}
										value={config.storage.memory.auth.password}
										onChange={handleChange('storage.memory.auth.password')}
									/>
									<ErrorBox configvalue="storage.memory.auth.password" messages={$tabs.storage.messages} />
									<Typography variant="caption">
										<Trans>Password for authorization.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<TextField
										type="number"
										label={<Trans>Maximum size (Megabytes)</Trans>}
										env={env('storage.memory.max_size_mbytes')}
										disabled={env('storage.memory.max_size_mbytes')}
										value={config.storage.memory.max_size_mbytes}
										onChange={handleChange('storage.memory.max_size_mbytes')}
									/>
									<ErrorBox configvalue="storage.memory.max_size_mbytes" messages={$tabs.storage.messages} />
									<Typography variant="caption">
										<Trans>Maximum allowed megabytes of RAM for /memfs, 0 for unlimited.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Checkbox
										label={<Trans>Remove the oldest entries if the /memfs is full</Trans>}
										checked={config.storage.memory.purge}
										disabled={env('storage.memory.purge') || toInt(config.storage.memory.max_size_mbytes) <= 0}
										onChange={handleChange('storage.memory.purge')}
									/>{' '}
									{env('storage.memory.purge') && <Env />}
									<ErrorBox configvalue="storage.memory.purge" messages={$tabs.storage.messages} />
								</Grid>
								<Grid item xs={12}>
									<Divider />
								</Grid>
								<Grid item xs={12}>
									<Typography variant="h3">
										<Trans>Disk</Trans>
									</Typography>
									<Typography variant="body1">
										<Trans>Settings for /data path. The access is protected by </Trans>{' '}
										<Link
											color="secondary"
											href="#/settings/auth"
											onClick={() => {
												setTab('auth');
											}}
										>
											login (JWT)
										</Link>
										.
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<TextField
										type="number"
										label={<Trans>Maximum size (Megabytes)</Trans>}
										value={config.storage.disk.max_size_mbytes}
										env={env('storage.disk.max_size_mbytes')}
										disabled={env('storage.disk.max_size_mbytes')}
										onChange={handleChange('storage.disk.max_size_mbytes')}
									/>
									<ErrorBox configvalue="storage.disk.max_size_mbytes" messages={$tabs.storage.messages} />
									<Typography variant="caption">
										<Trans>Maximum allowed megabytes to consume from hard disk. 0 for unlimited.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Divider />
								</Grid>
								<Grid item xs={12}>
									<Typography variant="h3">
										<Trans>Disk cache</Trans>
									</Typography>
									<Typography variant="body1">
										<Trans>Cache for files on /data.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Checkbox
										label={<Trans>Disk cache</Trans>}
										checked={config.storage.disk.cache.enable}
										disabled={env('storage.disk.cache.enable')}
										onChange={handleChange('storage.disk.cache.enable')}
									/>{' '}
									{env('storage.disk.cache.enable') && <Env />}
									<ErrorBox configvalue="storage.disk.cache.enable" messages={$tabs.storage.messages} />
								</Grid>
								<Grid item xs={12}>
									<TextField
										type="number"
										label={<Trans>Maximum size (Megabytes)</Trans>}
										env={env('storage.disk.cache.max_size_mbytes')}
										disabled={env('storage.disk.cache.max_size_mbytes') || !config.storage.disk.cache.enable}
										value={config.storage.disk.cache.max_size_mbytes}
										onChange={handleChange('storage.disk.cache.max_size_mbytes')}
										error={
											parseInt(config.storage.disk.cache.max_file_size_mbytes) > parseInt(config.storage.disk.cache.max_size_mbytes) &&
											parseInt(config.storage.disk.cache.max_size_mbytes) !== 0
										}
										helperText={
											parseInt(config.storage.disk.cache.max_file_size_mbytes) > parseInt(config.storage.disk.cache.max_size_mbytes) &&
											parseInt(config.storage.disk.cache.max_size_mbytes) !== 0 ? (
												<Trans>Must be larger than maximum file size in cache.</Trans>
											) : (
												''
											)
										}
									/>
									<ErrorBox configvalue="storage.disk.cache.max_size_mbytes" messages={$tabs.storage.messages} />
									<Typography variant="caption">
										<Trans>Maximum allowed cache size, 0 for unlimited.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12} md={6}>
									<TextField
										type="number"
										label={<Trans>Cache time (Seconds)</Trans>}
										env={env('storage.disk.cache.ttl_seconds')}
										disabled={env('storage.disk.cache.ttl_seconds') || !config.storage.disk.cache.enable}
										value={config.storage.disk.cache.ttl_seconds}
										onChange={handleChange('storage.disk.cache.ttl_seconds')}
									/>
									<ErrorBox configvalue="storage.disk.cache.ttl_seconds" messages={$tabs.storage.messages} />
									<Typography variant="caption">
										<Trans>Seconds to keep files in cache.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12} md={6}>
									<TextField
										type="number"
										label={<Trans>Maximum file size (Megabytes)</Trans>}
										env={env('storage.disk.cache.max_file_size_mbytes')}
										disabled={env('storage.disk.cache.max_file_size_mbytes') || !config.storage.disk.cache.enable}
										value={config.storage.disk.cache.max_file_size_mbytes}
										onChange={handleChange('storage.disk.cache.max_file_size_mbytes')}
										error={
											parseInt(config.storage.disk.cache.max_file_size_mbytes) > parseInt(config.storage.disk.cache.max_size_mbytes) &&
											parseInt(config.storage.disk.cache.max_size_mbytes) !== 0
										}
										helperText={
											parseInt(config.storage.disk.cache.max_file_size_mbytes) > parseInt(config.storage.disk.cache.max_size_mbytes) &&
											parseInt(config.storage.disk.cache.max_size_mbytes) !== 0 ? (
												<Trans>Must be smaller than maximum cache size.</Trans>
											) : (
												''
											)
										}
									/>
									<ErrorBox configvalue="storage.disk.cache.max_file_size_mbytes" messages={$tabs.storage.messages} />
									<Typography variant="caption">
										<Trans>Maximum file size to put in cache.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={6}>
									<TextField
										multiline
										rows={5}
										label={<Trans>Cache types</Trans>}
										env={env('storage.disk.cache.types.allow')}
										disabled={env('storage.disk.cache.types.allow') || !config.storage.disk.cache.enable}
										value={config.storage.disk.cache.types.allow}
										onChange={handleChange('storage.disk.cache.types.allow')}
									/>
									<ErrorBox configvalue="storage.disk.cache.types.allow" messages={$tabs.storage.messages} />
									<Typography variant="caption">
										<Trans>List of file extensions to cache (e.g. ".html"), one per line. Leave empty to cache all file types.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={6}>
									<TextField
										multiline
										rows={5}
										label={<Trans>Block cache types</Trans>}
										env={env('storage.disk.cache.types.block')}
										disabled={env('storage.disk.cache.types.block') || !config.storage.disk.cache.enable}
										value={config.storage.disk.cache.types.block}
										onChange={handleChange('storage.disk.cache.types.block')}
									/>
									<ErrorBox configvalue="storage.disk.cache.types.block" messages={$tabs.storage.messages} />
									<Typography variant="caption">
										<Trans>List of file extensions not to cache (e.g. ".m3u8"), one per line. Leave empty for none.</Trans>
									</Typography>
								</Grid>
							</Grid>
						</TabPanel>
						<TabPanel value={$tab} index="rtmp" className="panel">
							<Grid container spacing={2}>
								<Grid item xs={12}>
									<Typography variant="h2">
										<Trans>RTMP</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Checkbox
										label={<Trans>RTMP server</Trans>}
										checked={config.rtmp.enable}
										disabled={env('rtmp.enable') || config.rtmp.enable_tls}
										onChange={handleChange('rtmp.enable')}
									/>{' '}
									{env('rtmp.enable') && <Env style={{ marginRight: '2em' }} />}
									<ErrorBox configvalue="rtmp.enable" messages={$tabs.rtmp.messages} />
									<Checkbox
										label={<Trans>RTMPS server</Trans>}
										checked={config.rtmp.enable_tls}
										disabled={env('rtmp.enable_tls')}
										onChange={handleChange('rtmp.enable_tls')}
									/>{' '}
									{env('rtmp.enable_tls') && <Env />}
									<ErrorBox configvalue="rtmp.enable_tls" messages={$tabs.rtmp.messages} />
									{config.rtmp.enable_tls && !config.tls.auto && (
										<Typography variant="caption">
											<Trans>Requires activation</Trans>{' '}
											<Link
												color="secondary"
												href="#/settings/network"
												onClick={() => {
													setTab('network');
												}}
											>
												TLS/HTTPS
											</Link>
											.
										</Typography>
									)}
								</Grid>
								<Grid item xs={12}>
									<Divider />
								</Grid>
								<Grid item xs={6} md={3}>
									<TextField
										type="number"
										label={<Trans>RTMP Port</Trans>}
										env={env('rtmp.address')}
										disabled={env('rtmp.address') || (!config.rtmp.enable && !config.rtmp.enable_tls)}
										value={config.rtmp.address}
										onChange={handleChange('rtmp.address')}
									/>
									<ErrorBox configvalue="rtmp.address" messages={$tabs.rtmp.messages} />
									<Typography variant="caption">
										<Trans>RTMP server listen address.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={6} md={3}>
									<TextField
										type="number"
										label={<Trans>RTMPS Port</Trans>}
										env={env('rtmp.address_tls')}
										disabled={env('rtmp.address_tls') || !config.rtmp.enable_tls || !config.tls.auto}
										value={config.rtmp.address_tls}
										onChange={handleChange('rtmp.address_tls')}
									/>
									<ErrorBox configvalue="rtmp.address_tls" messages={$tabs.rtmp.messages} />
									<Typography variant="caption">
										<Trans>RTMPS server listen address.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12} md={6}>
									<TextField
										label={<Trans>App</Trans>}
										env={env('rtmp.app')}
										disabled={env('rtmp.app') || !config.rtmp.enable}
										value={config.rtmp.app}
										onChange={handleChange('rtmp.app')}
									/>
									<ErrorBox configvalue="rtmp.app" messages={$tabs.rtmp.messages} />
									<Typography variant="caption">
										<Trans>RTMP app for publishing.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Password
										label={<Trans>Token</Trans>}
										env={env('rtmp.token')}
										disabled={env('rtmp.token') || !config.rtmp.enable}
										value={config.rtmp.token}
										onChange={handleChange('rtmp.token')}
									/>
									<ErrorBox configvalue="rtmp.token" messages={$tabs.rtmp.messages} />
									<Typography variant="caption">
										<Trans>RTMP token for publishing and playing. The token is the value of the URL query parameter 'token.'</Trans>
									</Typography>
								</Grid>
							</Grid>
						</TabPanel>
						<TabPanel value={$tab} index="srt" className="panel">
							<Grid container spacing={2}>
								<Grid item xs={12}>
									<Typography variant="h2">
										<Trans>SRT</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Checkbox
										label={<Trans>SRT server</Trans>}
										checked={config.srt.enable}
										disabled={env('srt.enable')}
										onChange={handleChange('srt.enable')}
									/>{' '}
									{env('srt.enable') && <Env style={{ marginRight: '2em' }} />}
									<ErrorBox configvalue="srt.enable" messages={$tabs.srt.messages} />
								</Grid>
								<Grid item xs={12}>
									<Divider />
								</Grid>
								<Grid item xs={6} md={4}>
									<TextField
										type="number"
										label={<Trans>Port</Trans>}
										env={env('srt.address')}
										disabled={env('srt.address') || !config.srt.enable}
										value={config.srt.address}
										onChange={handleChange('srt.address')}
									/>
									<ErrorBox configvalue="srt.address" messages={$tabs.srt.messages} />
									<Typography variant="caption">
										<Trans>SRT server listen address.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={6} md={8}>
									<Password
										label={<Trans>Token</Trans>}
										env={env('srt.token')}
										disabled={env('srt.token') || !config.srt.enable}
										value={config.srt.token}
										onChange={handleChange('srt.token')}
									/>
									<ErrorBox configvalue="srt.token" messages={$tabs.srt.messages} />
									<Typography variant="caption">
										<Trans>SRT token for publishing and playing. The token is the value of the streamid parameter 'token.'</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Password
										label={<Trans>Passphrase</Trans>}
										env={env('srt.passphrase')}
										disabled={env('srt.passphrase') || !config.srt.enable}
										value={config.srt.passphrase}
										onChange={handleChange('srt.passphrase')}
										inputProps={{ maxLength: 79 }}
										error={config.srt.passphrase && config.srt.passphrase.length < 10}
										helperText={
											config.srt.passphrase && config.srt.passphrase.length < 10 ? (
												<Trans>Passphrase must be between 10 and 79 characters long</Trans>
											) : (
												false
											)
										}
									/>
									<ErrorBox configvalue="srt.passphrase" messages={$tabs.srt.messages} />
									<Typography variant="caption">
										<Trans>Passphrase for SRT encryption.</Trans>
									</Typography>
								</Grid>
							</Grid>
						</TabPanel>
						<TabPanel value={$tab} index="logging" className="panel">
							<Grid container spacing={2}>
								<Grid item xs={12}>
									<Typography variant="h2">
										<Trans>Logging</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Typography variant="h3">
										<Trans>System</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<BoxTextarea style={{ paddingTop: '.2em', paddingBottom: '1em' }}>
										<Textarea
											rows={17}
											value={$logdata}
											scrollTo="bottom"
											readOnly
											allowModal
											allowCopy
											allowDownload={true}
											downloadName="core.log"
										/>
									</BoxTextarea>
								</Grid>
								<Grid item xs={6}>
									<Select
										label={<Trans>Log level</Trans>}
										value={config.log.level}
										disabled={env('log.level')}
										onChange={handleChange('log.level')}
									>
										<MenuItem value="silent">
											Silent
											{env('log.level') && (
												<div className={classes.inlineEnv}>
													<Env />
												</div>
											)}
										</MenuItem>
										<MenuItem value="error">
											Error
											{env('log.level') && (
												<div className={classes.inlineEnv}>
													<Env />
												</div>
											)}
										</MenuItem>
										<MenuItem value="warn">
											Warn
											{env('log.level') && (
												<div className={classes.inlineEnv}>
													<Env />
												</div>
											)}
										</MenuItem>
										<MenuItem value="info">
											Info
											{env('log.level') && (
												<div className={classes.inlineEnv}>
													<Env />
												</div>
											)}
										</MenuItem>
										<MenuItem value="debug">
											Debug
											{env('loglevel') && (
												<div className={classes.inlineEnv}>
													<Env />
												</div>
											)}
										</MenuItem>
									</Select>
									<ErrorBox configvalue="log.level" messages={$tabs.logging.messages} />
									<Typography variant="caption">
										<Trans>Level of system protocol.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={6}>
									<TextField
										type="number"
										label={<Trans>Maximum log lines</Trans>}
										value={config.log.max_lines}
										env={env('log.max_lines')}
										disabled={env('log.max_lines')}
										onChange={handleChange('log.max_lines')}
									/>
									<ErrorBox configvalue="log.max_lines" messages={$tabs.logging.messages} />
									<Typography variant="caption">
										<Trans>Number of log lines to keep.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Typography variant="h3">
										<Trans>FFmpeg</Trans>
									</Typography>
								</Grid>
								<Grid item xs={6}>
									<TextField
										type="number"
										label={<Trans>Maximum log lines</Trans>}
										value={config.ffmpeg.log.max_lines}
										env={env('ffmpeg.log.max_lines')}
										disabled={env('ffmpeg.log.max_lines')}
										onChange={handleChange('ffmpeg.log.max_lines')}
									/>
									<ErrorBox configvalue="ffmpeg.log.max_lines" messages={$tabs.logging.messages} />
									<Typography variant="caption">
										<Trans>Number of log lines to keep.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={6}>
									<TextField
										type="number"
										label={<Trans>Maximum log histroy</Trans>}
										value={config.ffmpeg.log.max_history}
										env={env('ffmpeg.log.max_history')}
										disabled={env('ffmpeg.log.max_history')}
										onChange={handleChange('ffmpeg.log.max_history')}
									/>
									<ErrorBox configvalue="ffmpeg.log.max_history" messages={$tabs.logging.messages} />
									<Typography variant="caption">
										<Trans>Number of logs to keep for each process.</Trans>
									</Typography>
								</Grid>
							</Grid>
						</TabPanel>
					</TabsVerticalGrid>
				</Grid>
				<PaperFooter
					buttonsLeft={
						<Button variant="outlined" color="default" onClick={handleAbort}>
							<Trans>Close</Trans>
						</Button>
					}
					buttonsRight={
						<Button variant="outlined" color="primary" disabled={!$config.modified} onClick={handleSave}>
							<Trans>Save</Trans>
						</Button>
					}
				/>
			</Paper>
			<Backdrop open={$saving}>
				<CircularProgress color="inherit" />
			</Backdrop>
			<Backdrop open={$restart.restarting}>
				<Paper xs={4} sm={4} md={4}>
					<PaperHeader title={<Trans>Restarting</Trans>} />
					<PaperContent>
						{$restart.timeout === false ? (
							<React.Fragment>
								<Typography variant="body1">
									<Trans>Restarting Restreamer Core ...</Trans>
								</Typography>
								<LinearProgress sx={{ mt: '1em' }} />
							</React.Fragment>
						) : (
							<React.Fragment>
								<Typography variant="body1">
									<Trans>Reconnecting to Restreamer Core failed for the last {RETRIES} seconds.</Trans>
								</Typography>
								<Typography variant="body1" sx={{ mt: '1em' }}>
									<Trans>
										If you enabled Let's Encrypt TLS it might take some time to acquire the certificates. Make sure that Restreamer Core is
										reachable via port 80 from the internet. Please check the console log of Restreamer Core.
									</Trans>
								</Typography>
								<Typography variant="body1" sx={{ mt: '1em' }}>
									<Trans>
										If you changed the ports, it might be that Restreamer Core restarted already, but it is now available on a different
										port.
									</Trans>
								</Typography>
							</React.Fragment>
						)}
					</PaperContent>
					<PaperFooter
						buttonsRight={
							<Button variant="outlined" color="primary" disabled={!$restart.timeout} onClick={() => window.location.reload()}>
								<Trans>Reload</Trans>
							</Button>
						}
					/>
				</Paper>
			</Backdrop>
			<Dialog
				open={$dialogs.saved}
				title={<Trans>Restart required</Trans>}
				onClose={handleSavedDialog}
				buttonsLeft={
					<Button variant="outlined" color="default" onClick={handleSavedDialog}>
						<Trans>Abort</Trans>
					</Button>
				}
				buttonsRight={
					<Button variant="outlined" color="secondary" onClick={handleReload}>
						<Trans>Restart</Trans>
					</Button>
				}
			>
				<Typography variant="body1">
					<Trans>
						You have changed the configuration. In order for the changes to take effect, you have to restart the application. Do you want to restart
						now?
					</Trans>
				</Typography>
			</Dialog>
			<Dialog
				open={$dialogs.restart}
				title={<Trans>Restart</Trans>}
				onClose={handleReloadDialog}
				buttonsLeft={
					<Button variant="outlined" color="default" onClick={handleReloadDialog}>
						<Trans>Abort</Trans>
					</Button>
				}
				buttonsRight={
					<Button variant="outlined" color="secondary" onClick={handleReload}>
						<Trans>Restart</Trans>
					</Button>
				}
			>
				<Typography variant="body1">
					<Trans>Do you really want to restart the application now?</Trans>
				</Typography>
			</Dialog>
		</React.Fragment>
	);
}

Settings.defaultProps = {
	restreamer: null,
};

Settings.propTypes = {
	restreamer: PropTypes.object.isRequired,
};
