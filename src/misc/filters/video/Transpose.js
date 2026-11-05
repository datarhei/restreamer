import React from 'react';

import { Trans } from '@lingui/macro';
import Grid from '@mui/material/Grid';
import MenuItem from '@mui/material/MenuItem';

import Select from '../../Select';

// Transpose Filter
// http://ffmpeg.org/ffmpeg-all.html#transpose-1

function init(initialState) {
	const state = {
		value: 'none',
		...initialState,
	};

	return state;
}

function createGraph(settings) {
	settings = init(settings);

	const mapping = [];

	switch (settings.value) {
		case '90':
			mapping.push('transpose=dir=clock:passthrough=none');
			break;
		case '180':
			mapping.push('transpose=dir=clock:passthrough=none', 'transpose=dir=clock:passthrough=none');
			break;
		case '270':
			mapping.push('transpose=dir=cclock:passthrough=none');
			break;
		default:
			break;
	}

	return mapping.join(',');
}

// filter
function Rotate(props) {
	return (
		<Select label={<Trans>Rotate</Trans>} value={props.value} onChange={props.onChange}>
			<MenuItem value="none">None</MenuItem>
			<MenuItem value="90">90°</MenuItem>
			<MenuItem value="180">180°</MenuItem>
			<MenuItem value="270">270°</MenuItem>
		</Select>
	);
}

Rotate.defaultProps = {
	value: '',
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
		<Grid item xs={12}>
			<Rotate value={settings.value} onChange={update('value')} allowCustom />
		</Grid>
	);
}

Filter.defaultProps = {
	settings: {},
	onChange: function (settings, mapping) {},
};

const filter = 'transpose';
const name = 'Transpose';
const type = 'video';
const hwaccel = false;

function summarize(settings) {
	return `${name} (${settings.value}° clockwise)`;
}

function defaults() {
	const settings = init({});

	return {
		settings: settings,
		graph: createGraph(settings),
	};
}

export { name, filter, type, hwaccel, summarize, defaults, createGraph, Filter as component };
