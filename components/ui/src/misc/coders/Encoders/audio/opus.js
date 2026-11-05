import React from 'react';

import Grid from '@mui/material/Grid';
import Typography from '@mui/material/Typography';

import { useLingui } from '@lingui/react';
import { Trans, t } from '@lingui/macro';

import Audio from '../../settings/Audio';
import SelectCustom from '../../../SelectCustom';
import Helper from '../../helper';

function init(initialState) {
	const state = {
		bitrate: '64',
		delay: 'auto',
		...initialState,
	};

	return state;
}

function createMapping(settings, stream, skills) {
	stream = Helper.InitStream(stream);
	skills = Helper.InitSkills(skills);

	let sampling = settings.sampling;
	let layout = settings.layout;

	if (sampling === 'inherit') {
		sampling = stream.sampling_hz;
	}

	if (layout === 'inherit') {
		layout = stream.layout;
	}

	const local = ['-codec:a', 'opus', '-b:a', `${settings.bitrate}k`, '-vbr', 'on', '-shortest'];

	if (settings.delay !== 'auto') {
		local.push('opus_delay', settings.delay);
	}

	const mapping = {
		global: [['-vsync', 'drop']],
		local: local,
		filter: [],
	};

	return mapping;
}

function Delay(props) {
	const { i18n } = useLingui();
	const options = [
		{ value: '20', label: '20ms' },
		{ value: '30', label: '30ms' },
		{ value: '50', label: '50ms' },
	];

	if (props.allowAuto === true) {
		options.unshift({ value: 'auto', label: 'auto' });
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
				<Trans>Maximum delay in milliseconds.</Trans>
			</Typography>
		</React.Fragment>
	);
}

Delay.defaultProps = {
	variant: 'outlined',
	allowAuto: false,
	allowCustom: false,
	label: <Trans>Delay</Trans>,
	customLabel: <Trans>Custom delay</Trans>,
	onChange: function () {},
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
		const value = event.target.value;

		const newSettings = {
			...settings,
			[what]: value,
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
				<Audio.Bitrate value={settings.bitrate} onChange={update('bitrate')} allowCustom />
			</Grid>
			<Grid item xs={12}>
				<Delay value={settings.delay} onChange={update('delay')} allowAuto allowCustom />
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

const coder = 'opus';
const name = 'Opus';
const codec = 'opus';
const type = 'audio';
const hwaccel = false;

function summarize(settings) {
	return `${name}, ${settings.bitrate} kbit/s`;
}

function defaults(stream, skills) {
	const settings = init({});

	return {
		settings: settings,
		mapping: createMapping(settings, stream, skills),
	};
}

export { coder, name, codec, type, hwaccel, summarize, defaults, Coder as component };
