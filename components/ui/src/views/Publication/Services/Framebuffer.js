import React from 'react';

import { faImages as icon } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Trans } from '@lingui/macro';
import Grid from '@mui/material/Grid';
import MenuItem from '@mui/material/MenuItem';

import Select from '../../../misc/Select';
import SelectCustom from '../../../misc/SelectCustom';

const id = 'framebuffer';
const name = 'Framebuffer';
const version = '1.0';
const stream_key_link = '';
const description = <Trans>Send video to Framebuffer</Trans>;
const image_copyright = '';

const author = {
	creator: {
		name: 'datarhei',
		link: 'https://github.com/datarhei',
	},
	maintainer: {
		name: 'datarhei',
		link: 'https://github.com/datarhei',
	},
};

const category = 'universal';
const requires = {
	codecs: {
		audio: ['none'],
		video: ['rawvideo'],
	},
	formats: ['fbdev'],
	devices: ['fbdev'],
};

function ServiceIcon(props) {
	return <FontAwesomeIcon icon={icon} style={{ color: 'rgba(57, 181, 74, 1)' }} {...props} />;
}

function init(settings) {
	const initSettings = {
		device: 'none',
		pix_fmt: 'bgra',
		...settings,
	};

	return initSettings;
}

function Service(props) {
	const settings = init(props.settings);

	const handleChange = (what) => (event) => {
		const value = event.target.value;

		settings[what] = value;

		const outputs = createOutput(settings);

		props.onChange(outputs, settings);
	};

	const createOutput = (settings) => {
		const outputs = [];

		outputs.push({
			address: settings.device,
			options: ['-pix_fmt', settings.pix_fmt, '-f', 'fbdev'],
		});

		return outputs;
	};

	const filteredDevices = props.skills.devices.fbdev.filter((device) => device.extra !== '');
	const options = filteredDevices.map((device) => {
		return {
			value: device.id,
			label: device.name + ' (' + device.extra + ')',
		};
	});

	options.unshift({
		value: 'none',
		label: `Choose an output device ...`,
		disabled: true,
	});

	const videoDevices = (
		<SelectCustom options={options} label={<Trans>Device</Trans>} value={settings.device} onChange={handleChange('device')} variant="outlined" />
	);

	return (
		<Grid container spacing={2}>
			<Grid item xs={12}>
				{videoDevices}
			</Grid>
			<Grid item xs={12}>
				<Select label={<Trans>Pixel format</Trans>} value={settings.pix_fmt} onChange={handleChange('pix_fmt')}>
					<MenuItem value="rgba">RGBA</MenuItem>
					<MenuItem value="bgra">BGRA</MenuItem>
					<MenuItem value="argb">ARGB</MenuItem>
					<MenuItem value="abgr">ABGR</MenuItem>
					<MenuItem value="rgb24">RGB24</MenuItem>
					<MenuItem value="bgr24">BGR24</MenuItem>
					<MenuItem value="rgb565le">RGB565 Little Endian</MenuItem>
					<MenuItem value="rgb565be">RGB565 Big Endian</MenuItem>
				</Select>
			</Grid>
		</Grid>
	);
}

Service.defaultProps = {
	settings: {},
	skills: {},
	metadata: {},
	streams: [],
	onChange: function (output, settings) {},
};

export { id, name, version, stream_key_link, description, image_copyright, author, category, requires, ServiceIcon as icon, Service as component };
