import React from 'react';

import { Trans } from '@lingui/macro';
import Grid from '@mui/material/Grid';
import MenuItem from '@mui/material/MenuItem';
import TextField from '@mui/material/TextField';

import Select from '../../Select';

// Volume Filter
// http://ffmpeg.org/ffmpeg-all.html#volume

function init(initialState) {
	const state = {
		level: 'inherit',
		db: 0,
		...initialState,
	};

	return state;
}

function createGraph(settings) {
	settings = init(settings);

	const mapping = [];

	switch (settings.level) {
		case 'inherit':
			break;
		case 'custom':
			mapping.push(`volume=volume=${settings.db}dB`);
			break;
		default:
			mapping.push(`volume=volume=${parseInt(settings.level) / 100}`);
			break;
	}

	return mapping.join(',');
}

function VolumeLevel(props) {
	return (
		<Select label={<Trans>Volume</Trans>} value={props.value} onChange={props.onChange}>
			<MenuItem value="inherit">
				<Trans>Inherit</Trans>
			</MenuItem>
			<MenuItem value="10">10%</MenuItem>
			<MenuItem value="20">20%</MenuItem>
			<MenuItem value="30">30%</MenuItem>
			<MenuItem value="40">40%</MenuItem>
			<MenuItem value="50">50%</MenuItem>
			<MenuItem value="60">60%</MenuItem>
			<MenuItem value="70">70%</MenuItem>
			<MenuItem value="80">80%</MenuItem>
			<MenuItem value="90">90%</MenuItem>
			<MenuItem value="custom">
				<Trans>Custom ...</Trans>
			</MenuItem>
		</Select>
	);
}

VolumeLevel.defaultProps = {
	value: '',
	onChange: function (event) {},
};

function VolumeDB(props) {
	return (
		<TextField
			variant="outlined"
			fullWidth
			label={<Trans>Decibels (dB)</Trans>}
			type="number"
			value={props.value}
			disabled={props.disabled}
			onChange={props.onChange}
		/>
	);
}

VolumeDB.defaultProps = {
	value: '',
	disabled: false,
	onChange: function (event) {},
};

function Filter(props) {
	const settings = init(props.settings);

	const handleChange = (newSettings) => {
		let automatic = false;
		if (!newSettings) {
			newSettings = settings;
			automatic = true;
		}

		props.onChange(newSettings, createGraph(newSettings), automatic);
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
		<React.Fragment>
			<Grid item xs={6}>
				<VolumeLevel value={settings.level} onChange={update('level')} />
			</Grid>
			<Grid item xs={6}>
				<VolumeDB value={settings.db} onChange={update('db')} disabled={settings.level !== 'custom'} />
			</Grid>
		</React.Fragment>
	);
}

Filter.defaultProps = {
	settings: {},
	onChange: function (settings, graph, automatic) {},
};

const filter = 'volume';
const name = 'Volume';
const type = 'audio';
const hwaccel = false;

function summarize(settings) {
	let summary = `${name} (`;

	if (settings.level === 'custom') {
		summary += `${settings.db}dB`;
	} else {
		summary += `${settings.level}%`;
	}

	summary += ')';

	return summary;
}

function defaults() {
	const settings = init({});

	return {
		settings: settings,
		graph: createGraph(settings),
	};
}

export { name, filter, type, hwaccel, summarize, defaults, createGraph, Filter as component };
