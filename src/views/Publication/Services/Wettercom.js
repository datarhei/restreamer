import React from 'react';

import { Trans } from '@lingui/macro';
import Grid from '@mui/material/Grid';
import Link from '@mui/material/Link';
import TextField from '@mui/material/TextField';

import FormInlineButton from '../../../misc/FormInlineButton';
import Logo from './logos/wettercom.png';

const id = 'wettercom';
const name = 'wetter.com';
const version = '1.0';
const stream_key_link = 'https://www.wetter.com/hd-live-webcams/informationen/';
const description = (
	<Trans>
		Transmit your Livestream to wetter.com. Contact the publisher {' '}
		<Link color="secondary" target="_blank" href="https://www.wetter.com/hd-live-webcams/informationen/">
			here{' '}
		</Link>
		.
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
const category = 'software';
const requires = {
	protocols: ['rtmp'],
	formats: ['flv'],
	codecs: {
		audio: ['aac'],
		video: ['h264'],
	},
};

function ServiceIcon(props) {
	return <img src={Logo} alt="wetter.com Logo" {...props} />;
}

function init(settings) {
	const initSettings = {
		protocol: 'rtmp://',
		server_url: 'wettercom-c1.livespotting.com/c61cba4f-68a3-488e-b2e4-ffb8e182794f',
		stream_id: '',
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
			address: settings.protocol + settings.server_url + '/' + settings.stream_id + '.stream/' + settings.stream_id + '_token:' + settings.stream_key,
			options: ['-f', 'flv'],
		};

		return output;
	};

	return (
		<Grid container spacing={2}>
			<Grid item xs={12} md={12}>
				<TextField
					variant="outlined"
					fullWidth
					placeholder=""
					label={<Trans>Stream ID</Trans>}
					value={settings.stream_id}
					onChange={handleChange('stream_id')}
				/>
			</Grid>
			<Grid item xs={12} md={9}>
				<TextField
					variant="outlined"
					fullWidth
					placeholder=""
					label={<Trans>Stream key</Trans>}
					value={settings.stream_key}
					onChange={handleChange('stream_key')}
				/>
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
