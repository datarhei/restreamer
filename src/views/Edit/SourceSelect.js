import React from 'react';
import SemverSatisfies from 'semver/functions/satisfies';

import Grid from '@mui/material/Grid';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';

import { Trans } from '@lingui/macro';

import Sources from './Sources';

function initConfig(initialConfig) {
	let config = {};

	for (let s of Sources.List()) {
		config[s.id] = {};
	}

	config = {
		...config,
		...initialConfig,
	};

	return config;
}

function init(source) {
	const settings = {};

	for (let id of Sources.IDs()) {
		settings[id] = {};
	}

	settings[source.type] = source.settings;

	return settings;
}

function reducer(settings, data) {
	const newSettings = {
		...settings,
		...data,
	};

	return newSettings;
}

export default function SourceSelect(props) {
	// $source holds the currently selected device. It is initialized with the
	// last stored source.
	const [$source, setSource] = React.useState(props.source.type);

	// $settings is for storing the settings of the different devices, such that if
	// the user switches between them, they can be restored. It takes the last
	// stored source settings as initial value.
	const [$settings, setSettings] = React.useReducer(reducer, props.source, init);

	const config = initConfig(props.config);

	const handleSource = (source) => {
		props.onChange(props.type);
		setSource(source);

		props.onSelect(props.type, source);
	};

	const handleRefresh = async () => {
		await props.onRefresh();
	};

	const handleStore = async (name, data) => {
		return await props.onStore(name, data);
	};

	const handleProbe = async (settings, inputs) => {
		await props.onProbe(props.type, $source, settings, inputs);
	};

	const handleChange = (source) => (settings) => {
		setSettings({
			...$settings,
			[source]: settings,
		});

		props.onChange(props.type, source, settings);
	};

	let sourceControl = null;

	const s = Sources.Get($source);
	if (s !== null) {
		const Component = s.component;

		if (SemverSatisfies(props.skills.ffmpeg.version, s.ffversion)) {
			sourceControl = (
				<Component
					knownDevices={props.skills.sources[$source]}
					skills={props.skills}
					config={config[$source]}
					settings={$settings[$source]}
					onChange={handleChange($source)}
					onProbe={handleProbe}
					onRefresh={handleRefresh}
					onStore={handleStore}
				/>
			);
		}
	}

	return (
		<Grid container spacing={1}>
			<Grid item xs={12}>
				<Select
					type={props.type}
					selected={$source}
					ffversion={props.skills.ffmpeg.version}
					availableSources={props.skills.sources}
					onSelect={handleSource}
				/>
			</Grid>
			<Grid item xs={12}>
				{sourceControl}
			</Grid>
		</Grid>
	);
}

SourceSelect.defaultProps = {
	type: '',
	skills: {},
	source: {},
	config: {},
	onProbe: function (type, device, settings, inputs) {},
	onSelect: function (type, device) {},
	onChange: function (type, device, settings) {},
	onRefresh: function () {},
	onStore: function (name, data) {},
};

function Select(props) {
	const handleSource = (source) => () => {
		props.onSelect(source);
	};

	let availableSources = [];

	for (let s of Sources.List()) {
		if (!(s.id in props.availableSources)) {
			continue;
		}

		if (!s.capabilities.includes(props.type)) {
			continue;
		}

		if (!SemverSatisfies(props.ffversion, s.ffversion)) {
			continue;
		}

		const variant = s.id === props.selected ? 'bigSelected' : 'big';
		const Icon = s.icon;

		availableSources.push(
			<Grid item xs={6} md={4} align="center" key={s.id}>
				<Button variant={variant} onClick={handleSource(s.id)}>
					<div>
						<Icon />
						<Typography>{s.name}</Typography>
					</div>
				</Button>
			</Grid>,
		);
	}

	if (availableSources.length === 0) {
		return (
			<Grid container spacing={1}>
				<Grid item xs={12}>
					<Typography variant="body1">
						<Trans>No sources available</Trans>
					</Typography>
				</Grid>
			</Grid>
		);
	}

	return (
		<Grid container spacing={1}>
			{availableSources}
		</Grid>
	);
}

Select.defaultProps = {
	type: '',
	selected: '',
	ffversion: '0.0.0',
	availableSources: {},
	onSelect: function (source) {},
};
