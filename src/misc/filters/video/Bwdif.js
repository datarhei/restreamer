import React from 'react';

import { Trans } from '@lingui/macro';
import Grid from '@mui/material/Grid';
import MenuItem from '@mui/material/MenuItem';

import Checkbox from '../../Checkbox';
import Select from '../../Select';

// Deinterlace the input video ("bwdif" stands for "Bob Weaver Deinterlacing Filter").
// http://ffmpeg.org/ffmpeg-all.html#bwdif

function init(initialState) {
	const state = {
		enabled: false,
		mode: '1',
		parity: '-1',
		deint: '0',
		...initialState,
	};

	return state;
}

function createGraph(settings) {
	settings = init(settings);

	const mapping = [];

	if (settings.enabled) {
		mapping.push(`bwdif=mode=${settings.mode}:parity=${settings.parity}:deint=${settings.deint}`);
	}

	return mapping.join(',');
}

function Mode(props) {
	return (
		<Select label={<Trans>Deinterlace mode</Trans>} value={props.value} onChange={props.onChange}>
			<MenuItem value="0">
				<Trans>Each frames</Trans>
			</MenuItem>
			<MenuItem value="1">
				<Trans>Each field</Trans>
			</MenuItem>
		</Select>
	);
}

Mode.defaultProps = {
	value: '',
	onChange: function (event) {},
};

function Parity(props) {
	return (
		<Select label={<Trans>Deinterlace parity</Trans>} value={props.value} onChange={props.onChange}>
			<MenuItem value="0">
				<Trans>Top field</Trans>
			</MenuItem>
			<MenuItem value="1">
				<Trans>Bottom field</Trans>
			</MenuItem>
			<MenuItem value="-1">
				<Trans>Auto</Trans>
			</MenuItem>
		</Select>
	);
}

Parity.defaultProps = {
	value: '',
	onChange: function (event) {},
};

function Deint(props) {
	return (
		<Select label={<Trans>Deinterlace deint</Trans>} value={props.value} onChange={props.onChange}>
			<MenuItem value="0">
				<Trans>All frames</Trans>
			</MenuItem>
			<MenuItem value="1">
				<Trans>Marked frames</Trans>
			</MenuItem>
		</Select>
	);
}

Deint.defaultProps = {
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
		};
		if (['enabled'].includes(what)) {
			newSettings[what] = !settings.enabled;
		} else {
			newSettings[what] = event.target.value;
		}

		handleChange(newSettings);
	};

	React.useEffect(() => {
		handleChange(null);
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

	return (
		<React.Fragment>
			<Grid item>
				<Checkbox label={<Trans>Deinterlace (bwdif)</Trans>} checked={settings.enabled} onChange={update('enabled')} />
			</Grid>
			{settings.enabled && (
				<React.Fragment>
					<Grid item xs={12}>
						<Mode value={settings.mode} onChange={update('mode')}></Mode>
					</Grid>
					<Grid item xs={6}>
						<Parity value={settings.parity} onChange={update('parity')}></Parity>
					</Grid>
					<Grid item xs={6}>
						<Deint value={settings.deint} onChange={update('deint')}></Deint>
					</Grid>
				</React.Fragment>
			)}
		</React.Fragment>
	);
}

Filter.defaultProps = {
	settings: {},
	onChange: function (settings, mapping) {},
};

const filter = 'bwdif';
const name = 'Deinterlacing Filter';
const type = 'video';
const hwaccel = false;

function summarize(settings) {
	return `${name}`;
}

function defaults() {
	const settings = init({});

	return {
		settings: settings,
		graph: createGraph(settings),
	};
}

export { name, filter, type, hwaccel, summarize, defaults, createGraph, Filter as component };
