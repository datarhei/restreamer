import React from 'react';

import { Trans } from '@lingui/macro';
import Grid from '@mui/material/Grid';
import Link from '@mui/material/Link';
import MenuItem from '@mui/material/MenuItem';
import TextField from '@mui/material/TextField';

import Logo from './logos/picarto-tv.svg';
import Select from '../../../misc/Select';

const id = 'picartotv';
const name = 'Picarto TV';
const version = '1.0';
const stream_key_link = 'https://picarto.tv';
const description = (
	<Trans>
		Transmit the main source to the Picarto TV RTMP Service. More details about the settings can be found{' '}
		<Link color="secondary" target="_blank" href="https://help.picarto.tv/help/">
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
	return <img src={Logo} alt="Picarto TV Logo" {...props} />;
}

function init(settings) {
	const initSettings = {
		region: 'live.us',
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
		const output = {
			address: 'rtmp://' + settings.region + '.picarto.tv/golive/' + settings.stream_key,
			options: ['-f', 'flv'],
		};

		return output;
	};

	return (
		<Grid container spacing={2}>
			<Grid item xs={12}>
				<Select label={<Trans>Region</Trans>} value={settings.region} onChange={handleChange('region')}>
					<MenuItem value="live.us">Autodetect</MenuItem>
					<MenuItem value="live.us-losangeles">Los Angeles, USA</MenuItem>
					<MenuItem value="live.us-dalla">Dallas, USA</MenuItem>
					<MenuItem value="live.us-miami">Miami, USA</MenuItem>
					<MenuItem value="live.us-newyork">New York, USA</MenuItem>
					<MenuItem value="live.eu-west1">Europe</MenuItem>
				</Select>
			</Grid>
			<Grid item xs={12} md={9}>
				<TextField variant="outlined" fullWidth label={<Trans>Stream key</Trans>} value={settings.stream_key} onChange={handleChange('stream_key')} />
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
