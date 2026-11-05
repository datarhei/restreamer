import React from 'react';

import { useLingui } from '@lingui/react';
import { Trans, t } from '@lingui/macro';
import Grid from '@mui/material/Grid';
import Icon from '@mui/icons-material/DeviceHub';
import MenuItem from '@mui/material/MenuItem';
import TextField from '@mui/material/TextField';
import Typography from '@mui/material/Typography';

import Audio from '../../../misc/coders/settings/Audio';
import FormInlineButton from '../../../misc/FormInlineButton';
import Select from '../../../misc/Select';

const initSettings = (initialSettings) => {
	if (!initialSettings) {
		initialSettings = {};
	}

	const settings = {
		source: 'silence',
		layout: 'stereo',
		sampling: '44100',
		color: 'white',
		frequency: 440,
		beepfactor: 4,
		amplitude: 1.0,
		...initialSettings,
	};

	return settings;
};

const createInputs = (settings) => {
	const inputs = [];

	let address = '';

	switch (settings.source) {
		case 'silence':
			address = `anullsrc=r=${settings.sampling}:cl=${settings.layout}`;
			break;
		case 'noise':
			address = `anoisesrc=r=${settings.sampling}:color=${settings.color}:amplitude=${settings.amplitude}`;
			break;
		case 'sine':
			address = `sine=r=${settings.sampling}:frequency=${settings.frequency}:beep_factor=${settings.beepfactor}`;
			break;
		default:
			break;
	}

	if (address.length !== 0) {
		inputs.push({
			address: address,
			options: ['-f', 'lavfi'],
		});
	}

	return inputs;
};

function Source(props) {
	const { i18n } = useLingui();
	const settings = initSettings(props.settings);

	const handleChange = (what) => (event) => {
		const value = event.target.value;

		props.onChange({
			...settings,
			[what]: value,
		});
	};

	const handleProbe = () => {
		props.onProbe(settings, createInputs(settings));
	};

	return (
		<Grid container alignItems="flex-start" spacing={2} style={{ marginTop: '0.5em' }}>
			<Grid item xs={12}>
				<Typography>
					<Trans>Select audio source:</Trans>
				</Typography>
			</Grid>
			<Grid item xs={12}>
				<Select label={<Trans>Source</Trans>} value={settings.source} onChange={handleChange('source')}>
					<MenuItem value="silence">{i18n._(t`Silence`)}</MenuItem>
					<MenuItem value="noise">{i18n._(t`Noise`)}</MenuItem>
					<MenuItem value="sine">{i18n._(t`Sine`)}</MenuItem>
				</Select>
			</Grid>
			{settings.source === 'silence' && (
				<React.Fragment>
					<Grid item xs={12}>
						<Audio.Sampling value={settings.sampling} onChange={handleChange('sampling')} allowCustom />
					</Grid>
					<Grid item xs={12}>
						<Audio.Layout value={settings.layout} onChange={handleChange('layout')} />
					</Grid>
				</React.Fragment>
			)}
			{settings.source === 'noise' && (
				<React.Fragment>
					<Grid item xs={12}>
						<Audio.Sampling value={settings.sampling} onChange={handleChange('sampling')} allowCustom />
					</Grid>
					<Grid item xs={6}>
						<Select label={<Trans>Color</Trans>} value={settings.color} onChange={handleChange('color')}>
							<MenuItem value="white">{i18n._(t`white`)}</MenuItem>
							<MenuItem value="pink">{i18n._(t`pink`)}</MenuItem>
							<MenuItem value="brown">{i18n._(t`brown`)}</MenuItem>
							<MenuItem value="blue">{i18n._(t`blue`)}</MenuItem>
							<MenuItem value="violet">{i18n._(t`violet`)}</MenuItem>
							<MenuItem value="velvet">{i18n._(t`velvet`)}</MenuItem>
						</Select>
						<Typography variant="caption">
							<Trans>The noise color</Trans>
						</Typography>
					</Grid>
					<Grid item xs={6}>
						<TextField
							variant="outlined"
							fullWidth
							label={<Trans>Amplitude</Trans>}
							value={settings.amplitude}
							onChange={handleChange('amplitude')}
						/>
						<Typography variant="caption">
							<Trans>The amplitude (0.0 - 1.0) of the generated audio stream</Trans>
						</Typography>
					</Grid>
				</React.Fragment>
			)}
			{settings.source === 'sine' && (
				<React.Fragment>
					<Grid item xs={12}>
						<Audio.Sampling value={settings.sampling} onChange={handleChange('sampling')} allowCustom />
					</Grid>
					<Grid item xs={6}>
						<TextField
							variant="outlined"
							fullWidth
							label={<Trans>Frequency (Hz)</Trans>}
							value={settings.frequency}
							onChange={handleChange('frequency')}
						/>
						<Typography variant="caption">
							<Trans>The carrier frequency</Trans>
						</Typography>
					</Grid>
					<Grid item xs={6}>
						<TextField
							variant="outlined"
							fullWidth
							label={<Trans>Beep factor</Trans>}
							value={settings.beepfactor}
							onChange={handleChange('beepfactor')}
						/>
						<Typography variant="caption">
							<Trans>Enable a periodic beep every second with this value times the carrier frequency</Trans>
						</Typography>
					</Grid>
				</React.Fragment>
			)}
			<Grid item xs={12}>
				<FormInlineButton onClick={handleProbe}>
					<Trans>Probe</Trans>
				</FormInlineButton>
			</Grid>
		</Grid>
	);
}

Source.defaultProps = {
	settings: {},
	onChange: function (settings) {},
	onProbe: function (settings, inputs) {},
};

function SourceIcon(props) {
	return <Icon style={{ color: '#FFF' }} {...props} />;
}

const id = 'virtualaudio';
const name = <Trans>Virtual source</Trans>;
const capabilities = ['audio'];
const ffversion = '^4.1.0 || ^5.0.0 || ^6.1.0';

const func = {
	initSettings,
	createInputs,
};

export { id, name, capabilities, ffversion, SourceIcon as icon, Source as component, func };
