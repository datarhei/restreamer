import React from 'react';

import { useLingui } from '@lingui/react';
import { Trans, t } from '@lingui/macro';
import Grid from '@mui/material/Grid';
import Typography from '@mui/material/Typography';

import SelectCustom from '../../../misc/SelectCustom';

// Resample Filter
// https://ffmpeg.org/ffmpeg-filters.html#toc-aresample-1

function init(initialState) {
	const state = {
		channels: '2',
		layout: 'stereo',
		sampling: '44100',
		...initialState,
	};

	return state;
}

function createGraph(settings) {
	settings = init(settings);

	const mapping = [];

	const sampling = settings.sampling;
	const layout = settings.layout;

	if (sampling !== 'inherit') {
		mapping.push(`osr=${sampling}`);
	}

	if (layout !== 'inherit') {
		mapping.push(`ochl=${layout}`);
	}

	if (mapping.length === 0) {
		return '';
	}

	return 'aresample=' + mapping.join(':');
}

function Layout(props) {
	const { i18n } = useLingui();
	const options = [
		{ value: 'mono', label: 'mono' },
		{ value: 'stereo', label: 'stereo' },
	];

	if (props.allowAuto === true) {
		options.unshift({ value: 'auto', label: 'auto' });
	}

	if (props.allowInherit === true) {
		options.unshift({ value: 'inherit', label: i18n._(t`Inherit`) });
	}

	if (props.allowCustom === true) {
		options.push({ value: 'custom', label: i18n._(t`Custom ...`) });
	}

	return (
		<React.Fragment>
			<SelectCustom
				options={options}
				label={props.label}
				customLabel={props.customLabel}
				value={props.value}
				onChange={props.onChange}
				variant={props.variant}
				allowCustom={props.allowCustom}
			/>
			<Typography variant="caption">
				<Trans>The layout of the audio stream.</Trans>
			</Typography>
		</React.Fragment>
	);
}

Layout.defaultProps = {
	variant: 'outlined',
	allowAuto: false,
	allowInherit: false,
	allowCustom: false,
	label: <Trans>Layout</Trans>,
	customLabel: <Trans>Custom layout</Trans>,
	onChange: function () {},
};

function Sampling(props) {
	const { i18n } = useLingui();
	const options = [
		{ value: '96000', label: '96000 Hz' },
		{ value: '88200', label: '88200 Hz' },
		{ value: '48000', label: '48000 Hz' },
		{ value: '44100', label: '44100 Hz' },
		{ value: '22050', label: '22050 Hz' },
		{ value: '8000', label: '8000 Hz' },
	];

	if (props.allowAuto === true) {
		options.unshift({ value: 'auto', label: 'auto' });
	}

	if (props.allowInherit === true) {
		options.unshift({ value: 'inherit', label: i18n._(t`Inherit`) });
	}

	if (props.allowCustom === true) {
		options.push({ value: 'custom', label: i18n._(t`Custom ...`) });
	}

	return (
		<React.Fragment>
			<SelectCustom
				options={options}
				label={props.label}
				customLabel={props.customLabel}
				value={props.value}
				onChange={props.onChange}
				variant={props.variant}
				allowCustom={props.allowCustom}
			/>
			<Typography variant="caption">
				<Trans>The sample rate of the audio stream.</Trans>
			</Typography>
		</React.Fragment>
	);
}

Sampling.defaultProps = {
	variant: 'outlined',
	allowAuto: false,
	allowInherit: false,
	allowCustom: false,
	label: <Trans>Sampling</Trans>,
	customLabel: <Trans>Custom sampling (Hz)</Trans>,
	onChange: function () {},
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
		const value = event.target.value;

		const newSettings = {
			...settings,
			[what]: value,
		};

		if (what === 'layout') {
			let channels = 2;

			switch (value) {
				case 'mono':
					channels = 1;
					break;
				case 'stereo':
					channels = 2;
					break;
				default:
					break;
			}

			newSettings.channels = channels;
		}

		handleChange(newSettings);
	};

	React.useEffect(() => {
		handleChange(null);
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

	return (
		<React.Fragment>
			<Grid item xs={12}>
				<Sampling value={settings.sampling} onChange={update('sampling')} allowInherit allowCustom />
			</Grid>
			<Grid item xs={12}>
				<Layout value={settings.layout} onChange={update('layout')} allowInherit />
			</Grid>
		</React.Fragment>
	);
}

Filter.defaultProps = {
	settings: {},
	onChange: function (settings, graph, automatic) {},
};

const filter = 'aresample';
const name = 'Resample';
const type = 'audio';
const hwaccel = false;

function summarize(settings) {
	return `${name} (${settings.layout}, ${settings.sampling}Hz)`;
}

function defaults() {
	const settings = init({});

	return {
		settings: settings,
		graph: createGraph(settings),
	};
}

export { name, filter, type, hwaccel, summarize, defaults, createGraph, Filter as component };
