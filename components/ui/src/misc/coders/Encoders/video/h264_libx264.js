import React from 'react';

import Grid from '@mui/material/Grid';
import MenuItem from '@mui/material/MenuItem';

import { Trans } from '@lingui/macro';

import Select from '../../../Select';
import Video from '../../settings/Video';
import Helper from '../../helper';

function init(initialState) {
	const state = {
		preset: 'ultrafast',
		bitrate: '4096',
		fps: '25',
		gop: '2',
		profile: 'auto',
		tune: 'zerolatency',
		fps_mode: 'auto',
		...initialState,
	};

	return state;
}

function createMapping(settings, stream, skills) {
	stream = Helper.InitStream(stream);
	skills = Helper.InitSkills(skills);

	const local = [
		'-codec:v',
		'libx264',
		'-preset:v',
		`${settings.preset}`,
		'-b:v',
		`${settings.bitrate}k`,
		'-maxrate:v',
		`${settings.bitrate}k`,
		'-bufsize:v',
		`${settings.bitrate}k`,
		'-r',
		`${settings.fps}`,
		'-sc_threshold',
		'0',
		'-pix_fmt',
		'yuv420p',
	];

	if (settings.gop !== 'auto') {
		local.push(
			'-g',
			`${Math.round(parseInt(settings.fps) * parseInt(settings.gop)).toFixed(0)}`,
			'-keyint_min',
			`${Math.round(parseInt(settings.fps) * parseInt(settings.gop)).toFixed(0)}`,
		);
	}

	if (skills.ffmpeg.version_major >= 5) {
		local.push('-fps_mode', `${settings.fps_mode}`);
	}

	if (settings.profile !== 'auto') {
		local.push('-profile:v', `${settings.profile}`);
	}

	if (settings.tune !== 'none') {
		local.push('-tune:v', `${settings.tune}`);
	}

	const mapping = {
		global: [],
		local: local,
		filter: [],
	};

	return mapping;
}

function Preset(props) {
	return (
		<Select label={<Trans>Preset</Trans>} value={props.value} onChange={props.onChange}>
			<MenuItem value="ultrafast">ultrafast</MenuItem>
			<MenuItem value="superfast">superfast</MenuItem>
			<MenuItem value="veryfast">veryfast</MenuItem>
			<MenuItem value="faster">faster</MenuItem>
			<MenuItem value="fast">fast</MenuItem>
			<MenuItem value="medium">medium</MenuItem>
			<MenuItem value="slow">slow</MenuItem>
			<MenuItem value="slower">slower</MenuItem>
			<MenuItem value="veryslow">veryslow</MenuItem>
		</Select>
	);
}

Preset.defaultProps = {
	value: '',
	onChange: function (event) {},
};

function Tune(props) {
	return (
		<Select label={<Trans>Tune</Trans>} value={props.value} onChange={props.onChange}>
			<MenuItem value="none">none</MenuItem>
			<MenuItem value="animation">animation</MenuItem>
			<MenuItem value="fastdecode">fastdecode</MenuItem>
			<MenuItem value="film">film</MenuItem>
			<MenuItem value="grain">grain</MenuItem>
			<MenuItem value="stillimage">stillimage</MenuItem>
			<MenuItem value="zerolatency">zerolatency</MenuItem>
		</Select>
	);
}

Tune.defaultProps = {
	value: '',
	onChange: function (event) {},
};

function Coder(props) {
	const settings = init(props.settings);
	const stream = Helper.InitStream(props.stream);
	const skills = Helper.InitSkills(props.skills);

	const handleChange = (newSettings) => {
		let automatic = false;
		if (!newSettings) {
			newSettings = settings;
			automatic = true;
		}

		props.onChange(newSettings, createMapping(newSettings, stream, skills), automatic);
	};

	const update = (what) => (event) => {
		const newSettings = {
			...settings,
			[what]: event.target.value,
		};

		handleChange(newSettings);
	};

	React.useEffect(() => {
		handleChange(null);
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

	return (
		<Grid container spacing={2}>
			<Grid item xs={12}>
				<Video.Bitrate value={settings.bitrate} onChange={update('bitrate')} allowCustom />
			</Grid>
			<Grid item xs={12} md={6}>
				<Video.Framerate value={settings.fps} onChange={update('fps')} allowCustom />
			</Grid>
			<Grid item xs={12} md={6}>
				<Video.GOP value={settings.gop} onChange={update('gop')} allowAuto allowCustom />
			</Grid>
			{skills.ffmpeg.version_major >= 5 && (
				<Grid item xs={12}>
					<Video.FpsMode value={settings.fps_mode} onChange={update('fps_mode')} />
				</Grid>
			)}
			<Grid item xs={6}>
				<Preset value={settings.preset} onChange={update('preset')} />
			</Grid>
			<Grid item xs={6}>
				<Video.Profile value={settings.profile} onChange={update('profile')} />
			</Grid>
			<Grid item xs={6}>
				<Tune value={settings.tune} onChange={update('tune')} />
			</Grid>
		</Grid>
	);
}

Coder.defaultProps = {
	stream: {},
	settings: {},
	skills: {},
	onChange: function (settings, mapping) {},
};

const coder = 'libx264';
const name = 'H.264 (libx264)';
const codec = 'h264';
const type = 'video';
const hwaccel = false;

function summarize(settings) {
	return `${name}, ${settings.bitrate} kbit/s, ${settings.fps} FPS, Preset: ${settings.preset}, Profile: ${settings.profile}`;
}

function defaults(stream, skills) {
	const settings = init({});

	return {
		settings: settings,
		mapping: createMapping(settings, stream, skills),
	};
}

export { coder, name, codec, type, hwaccel, summarize, defaults, Coder as component };
