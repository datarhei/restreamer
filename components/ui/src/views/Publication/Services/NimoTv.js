import React from 'react';

import { Trans } from '@lingui/macro';
import Grid from '@mui/material/Grid';
import Link from '@mui/material/Link';
import MenuItem from '@mui/material/MenuItem';
import TextField from '@mui/material/TextField';

import FormInlineButton from '../../../misc/FormInlineButton';
import Logo from './logos/nimo-tv.svg';
import Select from '../../../misc/Select';

const id = 'nimotv';
const name = 'Nimo TV';
const version = '1.0';
const stream_key_link = 'https://dashboard.nimo.tv/i/stream-url';
const description = (
	<Trans>
		Transmit the main source to the Nimo TV RTMP Service. More details about the settings can be found{' '}
		<Link color="secondary" target="_blank" href="https://help.nimo.tv/help/">
			here
		</Link>
		.
	</Trans>
);
const image_copyright = <Trans>Please contact the operator of the service and check what happens.</Trans>;
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
const category = 'platform';
const requires = {
	protocols: ['rtmp'],
	formats: ['flv'],
	codecs: {
		audio: ['aac', 'mp3'],
		video: ['h264'],
	},
};

function ServiceIcon(props) {
	return <img src={Logo} alt="Nimo TV Logo" {...props} />;
}

function init(settings) {
	const initSettings = {
		lines: '1',
		stream_key: '',
		...settings,
	};

	return initSettings;
}

function Service(props) {
	const settings = init(props.settings);

	const handleChange = (what) => (event) => {
		const value = event.target.value;

		settings[what] = value;

		const output = createOutput(settings);

		props.onChange([output], settings);
	};

	const createOutput = (settings) => {
		if (!settings.stream_key) {
			return null;
		}
		let address = 'rtmp://txpush.rtmp.nimo.tv/live/';
		if (settings.lines === '2') {
			address = 'rtmp://162.62.215.16/txpush.rtmp.nimo.tv/live/';
		}
		const output = {
			address: address + settings.stream_key,
			options: ['-f', 'flv'],
		};

		return output;
	};

	return (
		<Grid container spacing={2}>
			<Grid item xs={12}>
				<Select label={<Trans>Lines</Trans>} value={settings.lines} onChange={handleChange('lines')}>
					<MenuItem value="1">Line 1</MenuItem>
					<MenuItem value="2">Line 2</MenuItem>
				</Select>
			</Grid>
			<Grid item xs={12} md={9}>
				<TextField variant="outlined" fullWidth label={<Trans>Stream key</Trans>} value={settings.stream_key} onChange={handleChange('stream_key')} />
			</Grid>
			<Grid item xs={12} md={3}>
				<FormInlineButton target="blank" href={stream_key_link} component="a">
					<Trans>GET</Trans>
				</FormInlineButton>
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
