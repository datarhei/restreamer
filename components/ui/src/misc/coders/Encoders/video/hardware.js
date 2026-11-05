import React from 'react';

import Grid from '@mui/material/Grid';
import MenuItem from '@mui/material/MenuItem';

import { Trans } from '@lingui/macro';

import Select from '../../../Select';
import TextField from '../../../TextField';
import Video from '../../settings/Video';
import Helper from '../../helper';

function createHardwareEncoder(config) {
	function init(initialState) {
		return {
			bitrate: '4096',
			fps: '25',
			gop: '1',
			fps_mode: 'cfr',
			preset: config.defaultPreset || '',
			device: config.defaultDevice,
			...initialState,
		};
	}

	function createMapping(settings, stream, skills) {
		stream = Helper.InitStream(stream);
		skills = Helper.InitSkills(skills);

		const local = [
			'-codec:v',
			config.coder,
			'-b:v',
			`${settings.bitrate}k`,
			'-maxrate:v',
			`${settings.bitrate}k`,
			'-bufsize:v',
			`${settings.bitrate}k`,
			'-r',
			`${settings.fps}`,
			'-bf',
			'0',
		];

		if (settings.gop !== 'auto') {
			local.push('-g', `${Math.round(parseInt(settings.fps) * parseInt(settings.gop)).toFixed(0)}`);
		}

		if (config.presets && settings.preset !== '') {
			local.push('-preset:v', settings.preset);
		}

		if (skills.ffmpeg.version_major >= 5) {
			local.push('-fps_mode', settings.fps_mode);
		}

		local.push(...config.localOptions);

		return {
			global: config.globalOptions(settings),
			local: local,
			filter: config.filters,
		};
	}

	function Preset(props) {
		return (
			<Select label={<Trans>Preset</Trans>} value={props.value} onChange={props.onChange}>
				{config.presets.map((preset) => (
					<MenuItem key={preset} value={preset}>
						{preset}
					</MenuItem>
				))}
			</Select>
		);
	}

	Preset.defaultProps = {
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
			handleChange({
				...settings,
				[what]: event.target.value,
			});
		};

		React.useEffect(() => {
			handleChange(null);
			// eslint-disable-next-line react-hooks/exhaustive-deps
		}, []);

		return (
			<Grid container spacing={2}>
				<Grid item xs={12}>
					<TextField label={<Trans>Hardware device</Trans>} type="text" value={settings.device} onChange={update('device')} />
				</Grid>
				<Grid item xs={12} md={6}>
					<Video.Bitrate value={settings.bitrate} onChange={update('bitrate')} allowCustom />
				</Grid>
				<Grid item xs={12} md={6}>
					<Video.Framerate value={settings.fps} onChange={update('fps')} allowCustom />
				</Grid>
				<Grid item xs={12} md={6}>
					<Video.GOP value={settings.gop} onChange={update('gop')} allowAuto allowCustom />
				</Grid>
				{config.presets && (
					<Grid item xs={12} md={6}>
						<Preset value={settings.preset} onChange={update('preset')} />
					</Grid>
				)}
				{skills.ffmpeg.version_major >= 5 && (
					<Grid item xs={12}>
						<Video.FpsMode value={settings.fps_mode} onChange={update('fps_mode')} />
					</Grid>
				)}
			</Grid>
		);
	}

	Coder.defaultProps = {
		stream: {},
		settings: {},
		skills: {},
		onChange: function (settings, mapping) {},
	};

	function summarize(settings) {
		return `${config.name}, ${settings.bitrate} kbit/s, ${settings.fps} FPS, Device: ${settings.device}`;
	}

	function defaults(stream, skills) {
		const settings = init({});

		return {
			settings: settings,
			mapping: createMapping(settings, stream, skills),
		};
	}

	return {
		coder: config.coder,
		name: config.name,
		codec: config.codec,
		type: 'video',
		hwaccel: true,
		summarize: summarize,
		defaults: defaults,
		component: Coder,
	};
}

export default createHardwareEncoder;
