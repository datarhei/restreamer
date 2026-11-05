import React from 'react';

import { Trans } from '@lingui/macro';
import Grid from '@mui/material/Grid';
import Link from '@mui/material/Link';
import TextField from '@mui/material/TextField';

import Logo from './logos/brightcove.svg';

const id = 'brightcove';
const name = 'Brightcove Live';
const version = '1.0';
const stream_key_link = '';
const description = (
	<Trans>
		Transmit the main source to a Brightcove Live Service. More details about the settings can be found{' '}
		<Link color="secondary" target="_blank" href="https://live.support.brightcove.com/get-started/brightcove-live-training-videos.html">
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
	return <img src={Logo} alt="brightcove live Logo" {...props} />;
}

function init(settings) {
	const initSettings = {
		protocol: 'rtmp://',
		address: '',
		stream_name: '',
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
			address: settings.protocol + settings.address + settings.stream_name,
			options: ['-f', 'flv'],
		};

		return output;
	};

	return (
		<Grid container spacing={2}>
			<Grid item xs={12} md={3}>
				<TextField
					variant="outlined"
					fullWidth
					label={<Trans>Protocol</Trans>}
					value={settings.protocol}
					onChange={handleChange('protocol')}
					readOnly
					disabled
				/>
			</Grid>
			<Grid item xs={12} md={9}>
				<TextField variant="outlined" fullWidth label={<Trans>Address</Trans>} value={settings.address} onChange={handleChange('address')} />
			</Grid>
			<Grid item xs={12}>
				<TextField
					variant="outlined"
					fullWidth
					label={<Trans>Stream name</Trans>}
					value={settings.stream_name}
					onChange={handleChange('stream_name')}
				/>
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
