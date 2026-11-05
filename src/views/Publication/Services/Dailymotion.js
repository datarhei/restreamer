import React from 'react';

import { Trans } from '@lingui/macro';
import Grid from '@mui/material/Grid';
import Link from '@mui/material/Link';
import TextField from '@mui/material/TextField';

import Logo from './logos/dailymotion.svg';

const id = 'dailymotion';
const name = 'Dailymotion';
const version = '1.0';
const stream_key_link = '';
const description = (
	<Trans>
		Transmit your Livestream to an Dailymotion RTMP service.{' '}
		<Link color="secondary" target="_blank" href="https://faq.dailymotion.com/hc/en-us/articles/115009103088-Create-and-configure-a-live-stream">
			Here{' '}
		</Link>
		you can find more details about the settings.
	</Trans>
);
const image_copyright = <Trans>Please get in touch with the operator of the service and check what happens.</Trans>;
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
	protocols: ['rtmp', 'rtmps'],
	formats: ['flv'],
	codecs: {
		audio: ['aac'],
		video: ['h264'],
	},
};

function ServiceIcon(props) {
	return <img src={Logo} alt="Dailymotion Logo" {...props} />;
}

function init(settings) {
	const initSettings = {
		server_url: '',
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
			address: settings.server_url,
			options: ['-rtmp_playpath', settings.stream_key, '-f', 'flv'],
		};

		return output;
	};

	return (
		<Grid container spacing={2}>
			<Grid item xs={12}>
				<TextField
					variant="outlined"
					fullWidth
					label={<Trans>Stream URL</Trans>}
					value={settings.server_url}
					onChange={handleChange('server_url')}
					error={settings.server_url.includes('rtmp://') || settings.server_url.includes('rtmps://') ? false : true}
					helperText={settings.server_url.includes('rtmp://') || settings.server_url.includes('rtmps://') ? false : 'Please enter a valid RTMP/S URL.'}
				/>
			</Grid>
			<Grid item xs={12} md={9}>
				<TextField
					variant="outlined"
					fullWidth
					placeholder="abc123"
					label={<Trans>Stream key</Trans>}
					value={settings.stream_key}
					onChange={handleChange('stream_key')}
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
