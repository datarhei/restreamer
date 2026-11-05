import React from 'react';

import { Trans } from '@lingui/react/macro';
import Grid from '@mui/material/Grid';
import MenuItem from '@mui/material/MenuItem';
import Typography from '@mui/material/Typography';

import Select from '../../Select';
import { withDefaults } from '../../../utils/props';

// Pan Filter
// https://ffmpeg.org/ffmpeg-filters.html#pan-1

function init(initialState) {
	const state = {
		value: 'inherit',
		...initialState,
	};

	return state;
}

function createGraph(settings) {
	settings = init(settings);

	const mapping = [];

	switch (settings.value) {
		case 'mute_left':
			mapping.push('pan=stereo|c1=c1');
			break;
		case 'mute_right':
			mapping.push('pan=stereo|c0=c0');
			break;
		default:
			break;
	}

	return mapping;
}

// filter
function Pan(inputProps) {
	const props = withDefaults(inputProps, PanDefaultProps);
	return (
		<React.Fragment>
			<Select label={<Trans>Pan</Trans>} value={props.value} onChange={props.onChange}>
				<MenuItem value="inherit">
					<Trans>Inherit</Trans>
				</MenuItem>
				<MenuItem value="mute_left">
					<Trans>Mute left</Trans>
				</MenuItem>
				<MenuItem value="mute_right">
					<Trans>Mute right</Trans>
				</MenuItem>
			</Select>
			<Typography variant="caption">
				<Trans>Mute a channel.</Trans>
			</Typography>
		</React.Fragment>
	);
}

const PanDefaultProps = {
	value: '',
	onChange: function (event) {},
};

function Filter(inputProps) {
	const props = withDefaults(inputProps, FilterDefaultProps);
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
		};

		newSettings[what] = event.target.value;

		handleChange(newSettings);
	};

	React.useEffect(() => {
		handleChange(null);
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

	return (
		<React.Fragment>
			<Grid item xs={12}>
				<Pan value={settings.value} onChange={update('value')} allowCustom />
			</Grid>
		</React.Fragment>
	);
}

const FilterDefaultProps = {
	settings: {},
	onChange: function (settings, graph, automatic) {},
};

const filter = 'pan';
const name = 'Pan';
const type = 'audio';
const hwaccel = false;

function summarize(settings) {
	return `${name} (${settings.value.replace(/_/i, ' ')})`;
}

function defaults() {
	const settings = init({});

	return {
		settings: settings,
		graph: createGraph(settings),
	};
}

export { name, filter, type, hwaccel, summarize, defaults, createGraph, Filter as component };
