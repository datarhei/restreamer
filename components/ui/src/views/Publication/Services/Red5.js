import React from 'react';

import { Trans } from '@lingui/macro';
import Grid from '@mui/material/Grid';
import Link from '@mui/material/Link';
import MenuItem from '@mui/material/MenuItem';
import TextField from '@mui/material/TextField';

import Logo from './logos/red5.svg';
import Select from '../../../misc/Select';

const id = 'red5';
const name = 'Red5';
const version = '1.0';
const stream_key_link = '';
const description = (
	<Trans>
		Transmit the main source to an Red5/Pro Server. More details about the settings can be found{' '}
		<Link color="secondary" target="_blank" href="https://www.red5pro.com/docs/server/ffmpeg/publishing-live-streams/">
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
const category = 'software';
const requires = {
	protocols: ['rtmp', 'rtp'],
	formats: ['flv', 'rtsp'],
	codecs: {
		audio: ['aac', 'mp3'],
		video: ['h264'],
	},
};

function ServiceIcon(props) {
	return <img src={Logo} alt="Red5 Logo" {...props} />;
}

function init(settings) {
	const initSettings = {
		protocol: 'rtmp://',
		username: '',
		password: '',
		address: '',
		application: 'live',
		stream_name: 'myStream',
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
		const options = [];

		if (settings.protocol === 'rtsp') {
			options.push('-f', 'rtsp');
		} else {
			options.push('-f', 'flv');
		}

		let protocol = '';

		if (settings.username.length !== 0 && settings.password.length !== 0) {
			protocol = settings.protocol + settings.username + ':' + settings.password + '@';
		} else {
			protocol = settings.protocol;
		}

		const output = {
			address: protocol + settings.address + '/' + settings.application + '/' + settings.stream_name,
			options: options,
		};

		return output;
	};

	return (
		<Grid container spacing={2}>
			<Grid item xs={12} md={3}>
				<Select type="select" label={<Trans>Protocol</Trans>} value={settings.protocol} onChange={handleChange('protocol')}>
					<MenuItem value="rtsp://">rtsp://</MenuItem>
					<MenuItem value="rtmp://">rtmp://</MenuItem>
				</Select>
			</Grid>
			<Grid item xs={12} md={9}>
				<TextField
					variant="outlined"
					fullWidth
					type="url"
					label={<Trans>Address</Trans>}
					placeholder="red5-address:port"
					value={settings.address}
					onChange={handleChange('address')}
				/>
			</Grid>
			<Grid item xs={12} md={6}>
				<TextField
					variant="outlined"
					fullWidth
					type="text"
					label={<Trans>Application</Trans>}
					value={settings.application}
					onChange={handleChange('application')}
				/>
			</Grid>
			<Grid item xs={12} md={6}>
				<TextField
					variant="outlined"
					fullWidth
					type="text"
					label={<Trans>Stream name</Trans>}
					value={settings.stream_name}
					onChange={handleChange('stream_name')}
				/>
			</Grid>
			<Grid item xs={12} md={6}>
				<TextField
					variant="outlined"
					fullWidth
					type="text"
					label={<Trans>Username</Trans>}
					placeholder="username"
					value={settings.username}
					onChange={handleChange('username')}
				/>
			</Grid>
			<Grid item xs={12} md={6}>
				<TextField
					variant="outlined"
					fullWidth
					type="text"
					label={<Trans>Password</Trans>}
					placeholder="username"
					value={settings.password}
					onChange={handleChange('password')}
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
