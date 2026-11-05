import React from 'react';

import { useLingui } from '@lingui/react';
import { Trans, t } from '@lingui/macro';
import Grid from '@mui/material/Grid';
import Icon from '@mui/icons-material/DeviceHub';
import MenuItem from '@mui/material/MenuItem';
import TextField from '@mui/material/TextField';
import Typography from '@mui/material/Typography';

import FormInlineButton from '../../../misc/FormInlineButton';
import Select from '../../../misc/Select';
import Video from '../../../misc/coders/settings/Video';

const initSettings = (initialSettings) => {
	if (!initialSettings) {
		initialSettings = {};
	}

	const settings = {
		source: 'none',
		fps: '25',
		size: '1920x1080',
		rule: 'S23/B3',
		scale: '1920x1080',
		ratio: '2/3',
		mold: '10',
		death_color: '#383838',
		life_color: '#FCEE21',
		flags: 'fast_bilinear',
		...initialSettings,
	};

	return settings;
};

const createInputs = (settings) => {
	const inputs = [];

	let address = '';

	switch (settings.source) {
		case 'testsrc':
			address = `testsrc=rate=${settings.fps}:size=${settings.size}`;
			break;
		case 'testsrc2':
			address = `testsrc2=rate=${settings.fps}:size=${settings.size}`;
			break;
		case 'pal75bars':
			address = `pal75bars=rate=${settings.fps}:size=${settings.size}`;
			break;
		case 'pal100bars':
			address = `pal100bars=rate=${settings.fps}:size=${settings.size}`;
			break;
		case 'rgbtestsrc':
			address = `rgbtestsrc=rate=${settings.fps}:size=${settings.size}`;
			break;
		case 'smptebars':
			address = `smptebars=rate=${settings.fps}:size=${settings.size}`;
			break;
		case 'smptehdbars':
			address = `smptehdbars=rate=${settings.fps}:size=${settings.size}`;
			break;
		case 'yuvtestsrc':
			address = `yuvtestsrc=rate=${settings.fps}:size=${settings.size}`;
			break;
		case 'life':
			address = `life=rate=${settings.fps}:size=${settings.size}:rule=${settings.rule}:ratio=${settings.ratio}:mold=${settings.mold}:death_color=${settings.death_color}:life_color=${settings.life_color}`;
			if (settings.scale !== settings.size) {
				address += `,scale=${settings.scale}:flags=${settings.flags}`;
			}
			break;
		default:
			break;
	}

	if (address.length !== 0) {
		inputs.push({
			address: address,
			options: ['-f', 'lavfi', '-re'],
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
					<Trans>Select video source:</Trans>
				</Typography>
			</Grid>
			<Grid item xs={12}>
				<Select label={<Trans>Source</Trans>} value={settings.source} onChange={handleChange('source')}>
					<MenuItem value="none" disabled>
						{i18n._(t`Select source ...`)}
					</MenuItem>
					<MenuItem value="testsrc">{i18n._(t`Test pattern`)}</MenuItem>
					<MenuItem value="testsrc2">{i18n._(t`Test pattern (extended)`)}</MenuItem>
					<MenuItem value="rgbtestsrc">{i18n._(t`RGB test pattern`)}</MenuItem>
					<MenuItem value="yuvtestsrc">{i18n._(t`YUV test pattern`)}</MenuItem>
					<MenuItem value="pal75bars">EBU PAL 75%</MenuItem>
					<MenuItem value="pal100bars">EBU PAL 100%</MenuItem>
					<MenuItem value="smptebars">SMPTE EG 1-1990</MenuItem>
					<MenuItem value="smptehdbars">SMPTE RP 219-2002</MenuItem>
					<MenuItem value="life">Game Of Life</MenuItem>
				</Select>
			</Grid>
			{['testsrc', 'testsrc2', 'pal75bars', 'pal100bars', 'rgbtestsrc', 'smptebars', 'smptehdbars', 'yuvtestsrc'].includes(settings.source) && (
				<React.Fragment>
					<Grid item xs={12}>
						<Video.Framerate value={settings.fps} onChange={handleChange('fps')} allowCustom />
					</Grid>
					<Grid item xs={12}>
						<Video.Size value={settings.size} onChange={handleChange('size')} allowCustom />
					</Grid>
				</React.Fragment>
			)}
			{settings.source === 'life' && (
				<React.Fragment>
					<Grid item xs={12}>
						<Video.Framerate value={settings.fps} onChange={handleChange('fps')} allowCustom />
					</Grid>
					<Grid item xs={12}>
						<Video.Size value={settings.size} onChange={handleChange('size')} allowCustom />
					</Grid>
					<Grid item xs={12}>
						<Video.Size
							value={settings.scale}
							label={<Trans>Scale</Trans>}
							customLabel={<Trans>Custom scale</Trans>}
							onChange={handleChange('scale')}
							allowCustom={true}
						/>
					</Grid>
					<Grid item xs={12}>
						<TextField variant="outlined" fullWidth label={<Trans>Rule</Trans>} value={settings.rule} onChange={handleChange('rule')} />
					</Grid>
					<Grid item xs={12}>
						<TextField variant="outlined" fullWidth label={<Trans>Mold</Trans>} value={settings.mold} onChange={handleChange('mold')} />
					</Grid>
					<Grid item xs={12}>
						<TextField variant="outlined" fullWidth label={<Trans>Ratio</Trans>} value={settings.ratio} onChange={handleChange('ratio')} />
					</Grid>
					<Grid item xs={12}>
						<TextField
							variant="outlined"
							fullWidth
							label={<Trans>Death color</Trans>}
							value={settings.death_color}
							onChange={handleChange('death_color')}
						/>
					</Grid>
					<Grid item xs={12}>
						<TextField
							variant="outlined"
							fullWidth
							label={<Trans>Life color</Trans>}
							value={settings.life_color}
							onChange={handleChange('life_color')}
						/>
					</Grid>
					<Grid item xs={12}>
						<TextField variant="outlined" fullWidth label={<Trans>Flags</Trans>} value={settings.flags} onChange={handleChange('flags')} />
					</Grid>
				</React.Fragment>
			)}
			<Grid item xs={12}>
				<FormInlineButton disabled={settings.source === 'none'} onClick={handleProbe}>
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

const id = 'virtualvideo';
const name = <Trans>Virtual source</Trans>;
const capabilities = ['video'];
const ffversion = '^4.1.0 || ^5.0.0 || ^6.1.0';

const func = {
	initSettings,
	createInputs,
};

export { id, name, capabilities, ffversion, SourceIcon as icon, Source as component, func };
