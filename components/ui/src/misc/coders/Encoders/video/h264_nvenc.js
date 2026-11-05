import React from 'react';

import Grid from '@mui/material/Grid';
import MenuItem from '@mui/material/MenuItem';

import { Trans } from '@lingui/macro';

import Select from '../../../Select';
import Video from '../../settings/Video';
import Helper from '../../helper';

function init(initialState) {
	const state = {
		bitrate: '4096',
		fps: '25',
		gop: '2',
		preset: 'medium',
		profile: 'auto',
		level: 'auto',
		rc: 'auto',
		...initialState,
	};

	return state;
}

function createMapping(settings, stream, skills) {
	stream = Helper.InitStream(stream);
	skills = Helper.InitSkills(skills);

	const local = [
		'-codec:v',
		'h264_nvenc',
		'-preset:v',
		`${settings.preset}`,
		'-b:v',
		`${settings.bitrate}k`,
		'-maxrate',
		`${settings.bitrate}k`,
		'-bufsize',
		`${settings.bitrate}k`,
		'-r',
		`${settings.fps}`,
		'-pix_fmt',
		'yuv420p',
	];

	if (settings.gop !== 'auto') {
		local.push('-g', `${Math.round(parseInt(settings.fps) * parseInt(settings.gop)).toFixed(0)}`);
	}

	if (settings.profile !== 'auto') {
		local.push('-profile:v', `${settings.profile}`);
	}

	if (settings.level !== 'auto') {
		local.push('-level:v', `${settings.level}`);
	}

	if (settings.rc !== 'auto') {
		local.push('-rc:v', `${settings.rc}`);
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
			<MenuItem value="default">default</MenuItem>
			<MenuItem value="slow">slow</MenuItem>
			<MenuItem value="medium">medium</MenuItem>
			<MenuItem value="fast">fast</MenuItem>
			<MenuItem value="hp">hp</MenuItem>
			<MenuItem value="hq">hq</MenuItem>
			<MenuItem value="bd">db</MenuItem>
			<MenuItem value="ll">ll</MenuItem>
			<MenuItem value="llhq">llhq</MenuItem>
			<MenuItem value="llhp">llhp</MenuItem>
			<MenuItem value="lossless">lossless</MenuItem>
			<MenuItem value="losslesshp">losslesshp</MenuItem>
			<MenuItem value="losslesshq">losslesshq</MenuItem>
		</Select>
	);
}

Preset.defaultProps = {
	value: '',
	onChange: function (event) {},
};

function Profile(props) {
	return (
		<Select label={<Trans>Profile</Trans>} value={props.value} onChange={props.onChange}>
			<MenuItem value="auto">auto</MenuItem>
			<MenuItem value="baseline">baseline</MenuItem>
			<MenuItem value="main">main</MenuItem>
			<MenuItem value="high">high</MenuItem>
			<MenuItem value="high444p">high444p</MenuItem>
		</Select>
	);
}

Profile.defaultProps = {
	value: '',
	onChange: function (event) {},
};

function Level(props) {
	return (
		<Select label={<Trans>Level</Trans>} value={props.value} onChange={props.onChange}>
			<MenuItem value="auto">auto</MenuItem>
			<MenuItem value="1">1</MenuItem>
			<MenuItem value="1.0">1.0</MenuItem>
			<MenuItem value="1b">1b</MenuItem>
			<MenuItem value="1.0b">1.0b</MenuItem>
			<MenuItem value="1.1">1.1</MenuItem>
			<MenuItem value="1.2">1.2</MenuItem>
			<MenuItem value="1.3">1.3</MenuItem>
			<MenuItem value="2">2</MenuItem>
			<MenuItem value="2.0">2.0</MenuItem>
			<MenuItem value="2.1">2.1</MenuItem>
			<MenuItem value="2.2">2.2</MenuItem>
			<MenuItem value="3">3</MenuItem>
			<MenuItem value="3.0">3.0</MenuItem>
			<MenuItem value="3.1">3.1</MenuItem>
			<MenuItem value="3.2">3.2</MenuItem>
			<MenuItem value="4">4</MenuItem>
			<MenuItem value="4.0">4.0</MenuItem>
			<MenuItem value="4.1">4.1</MenuItem>
			<MenuItem value="4.2">4.2</MenuItem>
			<MenuItem value="5">5</MenuItem>
			<MenuItem value="5.0">5.0</MenuItem>
			<MenuItem value="5.1">5.1</MenuItem>
		</Select>
	);
}

Level.defaultProps = {
	value: '',
	onChange: function (event) {},
};

function RateControl(props) {
	return (
		<Select label={<Trans>Rate control</Trans>} value={props.value} onChange={props.onChange}>
			<MenuItem value="auto">auto</MenuItem>
			<MenuItem value="constqp">constqp</MenuItem>
			<MenuItem value="vbr">vbr</MenuItem>
			<MenuItem value="cbr">cbr</MenuItem>
			<MenuItem value="cbr_ld_hq">cbr_ld_hq</MenuItem>
			<MenuItem value="cbr_hq">cbr_hq</MenuItem>
			<MenuItem value="vbr_hq">vbr_hq</MenuItem>
		</Select>
	);
}

RateControl.defaultProps = {
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
			<Grid item xs={12}>
				<Video.Framerate value={settings.fps} onChange={update('fps')} allowCustom />
			</Grid>
			<Grid item xs={12}>
				<Video.GOP value={settings.gop} onChange={update('gop')} allowAuto allowCustom />
			</Grid>
			<Grid item xs={6}>
				<Preset value={settings.preset} onChange={update('preset')} />
			</Grid>
			<Grid item xs={6}>
				<Profile value={settings.profile} onChange={update('profile')} />
			</Grid>
			<Grid item xs={6}>
				<Level value={settings.level} onChange={update('level')} />
			</Grid>
			<Grid item xs={6}>
				<RateControl value={settings.rc} onChange={update('rc')} />
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

const coder = 'h264_nvenc';
const name = 'H.264 (NVENC)';
const codec = 'h264';
const type = 'video';
const hwaccel = true;

function summarize(settings) {
	return `${name}, ${settings.bitrate} kbit/s, ${settings.fps} FPS, Profile: ${settings.profile}`;
}

function defaults(stream, skills) {
	const settings = init({});

	return {
		settings: settings,
		mapping: createMapping(settings, stream, skills),
	};
}

export { coder, name, codec, type, hwaccel, summarize, defaults, Coder as component };
