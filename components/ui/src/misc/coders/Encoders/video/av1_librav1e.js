import React from 'react';

import Grid from '@mui/material/Grid';
import MenuItem from '@mui/material/MenuItem';
import TextField from '@mui/material/TextField';
import Typography from '@mui/material/Typography';

import { Trans } from '@lingui/macro';

import Select from '../../../Select';
import Video from '../../settings/Video';
import Helper from '../../helper';

function init(initialState) {
	const state = {
		bitrate: '4096',
		fps: '25',
		fps_mode: 'auto',
		gop: '2',
		qp: '-1',
		speed: '-1',
		tiles: '0',
		tile_rows: '0',
		tile_columns: '0',
		params: '',
		...initialState,
	};

	return state;
}

function createMapping(settings, stream, skills) {
	stream = Helper.InitStream(stream);
	skills = Helper.InitSkills(skills);

	const local = [
		'-codec:v',
		'librav1e',
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
		'-qp',
		`${settings.qp}`,
		'-speed',
		`${settings.speed}`,
		'-tiles',
		`${settings.tiles}`,
		'-tile-rows',
		`${settings.tile_rows}`,
		'-tile-columns',
		`${settings.tile_columns}`,
	];

	if (settings.params.length !== 0) {
		local.push('-rav1e-params', `${settings.params}`);
	}

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

	const mapping = {
		global: [],
		local: local,
		filter: [],
	};

	return mapping;
}

function Speed(props) {
	return (
		<React.Fragment>
			<Select label={<Trans>Speed Preset</Trans>} value={props.value} onChange={props.onChange}>
				<MenuItem value="-1">auto (-1)</MenuItem>
				<MenuItem value="0">slowest (0)</MenuItem>
				<MenuItem value="1">1</MenuItem>
				<MenuItem value="2">2</MenuItem>
				<MenuItem value="3">3</MenuItem>
				<MenuItem value="4">4</MenuItem>
				<MenuItem value="5">5</MenuItem>
				<MenuItem value="6">6</MenuItem>
				<MenuItem value="7">7</MenuItem>
				<MenuItem value="8">8</MenuItem>
				<MenuItem value="9">9</MenuItem>
				<MenuItem value="10">fastest (10)</MenuItem>
			</Select>
			<Typography variant="caption">
				<Trans>What speed preset to use.</Trans>
			</Typography>
		</React.Fragment>
	);
}

Speed.defaultProps = {
	value: '-1',
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
				<Speed value={settings.speed} onChange={update('speed')} />
			</Grid>
			<Grid item xs={6}>
				<TextField
					variant="outlined"
					fullWidth
					type="number"
					inputProps={{ min: -1, max: 255 }}
					label={<Trans>QP</Trans>}
					value={settings.qp}
					onChange={update('qp')}
				/>
				<Typography variant="caption">
					<Trans>Constant Quantizer Mode (-1 to 255).</Trans>
				</Typography>
			</Grid>
			<Grid item xs={4}>
				<TextField
					variant="outlined"
					fullWidth
					type="number"
					inputProps={{ min: -1 }}
					label={<Trans>Tiles</Trans>}
					value={settings.tiles}
					onChange={update('tiles')}
				/>
				<Typography variant="caption">
					<Trans>Number of tiles encode with.</Trans>
				</Typography>
			</Grid>
			<Grid item xs={4}>
				<TextField
					variant="outlined"
					fullWidth
					type="number"
					inputProps={{ min: -1 }}
					label={<Trans>Tile Rows</Trans>}
					value={settings.tile_rows}
					onChange={update('tile_rows')}
				/>
				<Typography variant="caption">
					<Trans>Number of tiles rows to encode with.</Trans>
				</Typography>
			</Grid>
			<Grid item xs={4}>
				<TextField
					variant="outlined"
					fullWidth
					type="number"
					inputProps={{ min: -1 }}
					label={<Trans>Tile Columns</Trans>}
					value={settings.tile_columns}
					onChange={update('tile_columns')}
				/>
				<Typography variant="caption">
					<Trans>Number of tiles columns to encode with.</Trans>
				</Typography>
			</Grid>
			<Grid item xs={12}>
				<TextField variant="outlined" fullWidth label={<Trans>rav1e Parameters</Trans>} value={settings.params} onChange={update('params')} />
				<Typography variant="caption">
					<Trans>Set the rav1e configuration using a :-separated list of key=value parameters.</Trans>
				</Typography>
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

const coder = 'librav1e';
const name = 'AV1 (librav1e)';
const codec = 'av1';
const type = 'video';
const hwaccel = false;

function summarize(settings) {
	return `${name}, ${settings.bitrate} kbit/s, ${settings.fps} FPS, Speed: ${settings.speed}, QP: ${settings.qp}`;
}

function defaults(stream, skills) {
	const settings = init({});

	return {
		settings: settings,
		mapping: createMapping(settings, stream, skills),
	};
}

export { coder, name, codec, type, hwaccel, summarize, defaults, Coder as component };
