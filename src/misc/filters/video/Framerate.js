import React from 'react';

import { useLingui } from '@lingui/react';
import { Trans, t } from '@lingui/macro';
import Grid from '@mui/material/Grid';

import Checkbox from '../../Checkbox';
import SelectCustom from '../../SelectCustom';

// Framerate Filter
// http://ffmpeg.org/ffmpeg-all.html#framerate

function init(initialState) {
	const state = {
		enabled: false,
		fps: '30',
		...initialState,
	};

	return state;
}

function createGraph(settings) {
	settings = init(settings);

	const mapping = [];

	if (settings.enabled) {
		mapping.push(`framerate=fps=${settings.fps}`);
	}

	return mapping.join(',');
}

function Framerate(props) {
	const { i18n } = useLingui();
	const sizes = [
		{ value: '60', label: '60' },
		{ value: '59.94', label: '59.94' },
		{ value: '50', label: '50' },
		{ value: '30', label: '30' },
		{ value: '29.97', label: '29.97 (NTSC)' },
		{ value: '25', label: '25 (PAL)' },
		{ value: '24', label: '24 (Film)' },
		{ value: '23.97', label: '23.97 (NTSC Film)' },
		{ value: '15', label: '15' },
		{ value: '10', label: '10' },
		{ value: 'custom', label: i18n._(t`Custom ...`) },
	];

	return (
		<SelectCustom
			options={sizes}
			label={props.label}
			customLabel={props.customLabel}
			value={props.value}
			onChange={props.onChange}
			variant={props.variant}
			allowCustom={props.allowCustom}
		/>
	);
}

Framerate.defaultProps = {
	label: <Trans>Framerate</Trans>,
	customLabel: <Trans>Custom framerate</Trans>,
	value: '',
	variant: 'outlined',
	allowCustom: true,
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
			<Grid item xs={12}>
				<Checkbox label={<Trans>Framerate conversion (frame interpolation)</Trans>} checked={settings.enabled} onChange={update('enabled')} />
			</Grid>
			{settings.enabled && (
				<React.Fragment>
					<Grid item xs={12}>
						<Framerate value={settings.fps} onChange={update('fps')}></Framerate>
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

const filter = 'fps';
const name = 'Frame Interpolation';
const type = 'video';
const hwaccel = false;

function summarize(settings) {
	return `${name} (${settings.fps}fps)`;
}

function defaults() {
	const settings = init({});

	return {
		settings: settings,
		graph: createGraph(settings),
	};
}

export { name, filter, type, hwaccel, summarize, defaults, createGraph, Filter as component };
