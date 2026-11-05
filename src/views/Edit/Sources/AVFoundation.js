import React from 'react';

import { useLingui } from '@lingui/react';
import { Trans, t } from '@lingui/macro';
import makeStyles from '@mui/styles/makeStyles';
import Button from '@mui/material/Button';
import Grid from '@mui/material/Grid';
import Icon from '@mui/icons-material/Apple';
import RefreshIcon from '@mui/icons-material/Refresh';
import Typography from '@mui/material/Typography';

import Checkbox from '../../../misc/Checkbox';
import FormInlineButton from '../../../misc/FormInlineButton';
import SelectCustom from '../../../misc/SelectCustom';
import Video from '../../../misc/coders/settings/Video';

const useStyles = makeStyles((theme) => ({
	gridContainer: {
		marginTop: '0.5em',
	},
}));

const initSettings = (initialSettings) => {
	if (!initialSettings) {
		initialSettings = {};
	}

	const settings = {
		aindex: 'default',
		vindex: 'default',
		format: 'nv12',
		framerate: '25',
		size: 'auto',
		cursor: false,
		clicks: false,
		...initialSettings,
	};

	return settings;
};

const createInputs = (settings) => {
	const vindex = settings.vindex;
	const aindex = settings.aindex;

	const address = `${vindex}:${aindex}`;
	const input = {
		address: address,
		options: [],
	};

	input.options.push('-f', 'avfoundation');

	if (parseInt(settings.framerate) !== 0) {
		input.options.push('-framerate', '' + settings.framerate);
	}

	if (settings.format.length !== 0) {
		input.options.push('-pixel_format', settings.format);
	}

	if (settings.size.length !== 0 && settings.size !== 'auto') {
		input.options.push('-video_size', settings.size);
	}

	if (settings.cursor) {
		input.options.push('-capture_cursor');
	}

	if (settings.clicks) {
		input.options.push('-capture_mouse_clicks');
	}

	return [input];
};

function Source(props) {
	const classes = useStyles();
	const { i18n } = useLingui();
	const settings = initSettings(props.settings);

	const handleChange = (what) => (event) => {
		let data = {};

		if (['cursor', 'clicks'].includes(what)) {
			data[what] = !settings[what];
		} else {
			data[what] = event.target.value;
		}

		props.onChange({
			...settings,
			...data,
		});
	};

	const handleRefresh = () => {
		props.onRefresh();
	};

	const handleProbe = () => {
		props.onProbe(settings, createInputs(settings));
	};

	let filteredDevices = props.knownDevices.filter((device) => device.media === 'video');
	let options = filteredDevices.map((device) => {
		return {
			value: device.id,
			label: device.name,
		};
	});

	options.unshift({
		value: 'default',
		label: i18n._(t`Default`),
	});

	options.push({
		value: 'custom',
		label: i18n._(t`Custom ...`),
	});

	const videoDevices = (
		<SelectCustom
			options={options}
			label={<Trans>Video device</Trans>}
			customLabel={<Trans>Custom video index</Trans>}
			value={settings.vindex}
			onChange={handleChange('vindex')}
			variant="outlined"
			allowCustom
		/>
	);

	filteredDevices = props.knownDevices.filter((device) => device.media === 'audio');
	options = filteredDevices.map((device) => {
		return {
			value: device.id,
			label: device.name,
		};
	});

	options.unshift({
		value: 'none',
		label: i18n._(t`None`),
	});
	options.unshift({
		value: 'default',
		label: i18n._(t`Default`),
	});
	options.push({
		value: 'custom',
		label: i18n._(t`Custom ...`),
	});

	const audioDevices = (
		<SelectCustom
			options={options}
			label={<Trans>Audio device</Trans>}
			customLabel={<Trans>Custom audio index</Trans>}
			value={settings.aindex}
			onChange={handleChange('aindex')}
			variant="outlined"
			allowCustom
		/>
	);

	return (
		<Grid container alignItems="flex-start" spacing={2} className={classes.gridContainer}>
			<Grid item xs={12}>
				<Typography>
					<Trans>Select a device:</Trans>
				</Typography>
			</Grid>
			<Grid item xs={12}>
				{videoDevices}
			</Grid>
			<Grid item xs={12}>
				{audioDevices}
				<Button size="small" startIcon={<RefreshIcon />} onClick={handleRefresh} sx={{ float: 'right' }}>
					<Trans>Refresh</Trans>
				</Button>
			</Grid>
			<Grid item xs={12}>
				<Video.Format value={settings.format} onChange={handleChange('format')} allowCustom />
			</Grid>
			<Grid item xs={12}>
				<Video.Framerate value={settings.framerate} onChange={handleChange('framerate')} allowCustom />
			</Grid>
			<Grid item xs={12}>
				<Video.Size value={settings.size} onChange={handleChange('size')} allowAuto allowCustom />
			</Grid>
			<Grid item xs={12}>
				<Checkbox label={<Trans>Capture cursor</Trans>} checked={settings.cursor} onChange={handleChange('cursor')} />
			</Grid>
			<Grid item xs={12}>
				<Checkbox label={<Trans>Capture clicks</Trans>} checked={settings.clicks} onChange={handleChange('clicks')} />
			</Grid>
			<Grid item xs={12}>
				<FormInlineButton onClick={handleProbe}>
					<Trans>Probe</Trans>
				</FormInlineButton>
			</Grid>
		</Grid>
	);
}

Source.defaultProps = {
	knownDevices: [],
	settings: {},
	onChange: function (settings) {},
	onProbe: function (settings, inputs) {},
	onRefresh: function () {},
};

function SourceIcon(props) {
	return <Icon style={{ color: '#FFF' }} {...props} />;
}

const id = 'avfoundation';
const name = <Trans>AVFoundation</Trans>;
const capabilities = ['audio', 'video'];
const ffversion = '^4.1.0 || ^5.0.0 || ^6.1.0';

const func = {
	initSettings,
	createInputs,
};

export { id, name, capabilities, ffversion, SourceIcon as icon, Source as component, func };
