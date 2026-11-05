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
		profile: 'auto',
		entropy: 'default',
		...initialState,
	};

	return state;
}

function createMapping(settings, stream, skills) {
	stream = Helper.InitStream(stream);
	skills = Helper.InitSkills(skills);

	const local = [
		'-codec:v',
		'h264_videotoolbox',
		'-b:v',
		`${settings.bitrate}k`,
		'-maxrate:v',
		`${settings.bitrate}k`,
		'-bufsize:v',
		`${settings.bitrate}k`,
		'-r',
		`${settings.fps}`,
		'-pix_fmt',
		'yuv420p',
		'-realtime',
		'true',
	];

	if (settings.gop !== 'auto') {
		local.push('-g', `${Math.round(parseInt(settings.fps) * parseInt(settings.gop)).toFixed(0)}`);
	}

	if (settings.profile !== 'auto') {
		local.push('-profile:v', `${settings.profile}`);
	}

	if (settings.entropy !== 'default') {
		local.push('-coder:v', `${settings.entropy}`);
	}

	const mapping = {
		global: [],
		local: local,
		filter: [],
	};

	return mapping;
}

function Entropy(props) {
	return (
		<Select label={<Trans>Entropy coder</Trans>} value={props.value} onChange={props.onChange}>
			<MenuItem value="default">default</MenuItem>
			<MenuItem value="cavlc">CAVLC</MenuItem>
			<MenuItem value="cabac">CABAC</MenuItem>
		</Select>
	);
}

Entropy.defaultProps = {
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
			<Grid item xs={12}>
				<Video.Profile value={settings.profile} onChange={update('profile')} />
			</Grid>
			<Grid item xs={12}>
				<Entropy value={settings.entropy} onChange={update('entropy')} />
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

const coder = 'h264_videotoolbox';
const name = 'H.264 (VideoToolbox)';
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
