import React from 'react';

import { Trans } from '@lingui/macro';
import Grid from '@mui/material/Grid';
import Link from '@mui/material/Link';
import TextField from '@mui/material/TextField';

import Logo from './logos/livespotting.svg';

const id = 'livespotting';
const name = 'livespotting';
const version = '1.0';
const stream_key_link = '';
const description = (
	<Trans>
		Transmit the main source to an livespotting.com Ressource. More details about the settings can be found{' '}
		<Link color="secondary" target="_blank" href="https://livespotting.com/products-vpu.html">
			here
		</Link>
		.
	</Trans>
);
const image_copyright = <Trans>Please contact the operator of the service and check what happens.</Trans>;
const author = {
	creator: {
		name: 'livespotting.com',
		link: 'https://livespotting.com',
	},
	maintainer: {
		name: 'livespotting.com',
		link: 'https://livespotting.com',
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
	return <img src={Logo} alt="livespotting.com Logo" {...props} />;
}

function init(settings) {
	const initSettings = {
		protocol: 'rtmp://',
		base_url: '-vpu.livespotting.com:1935/live/',
		vpu_id: '',
		livesource_id: '',
		rtmp_token: '',
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
		const options = ['-f', 'flv'];

		const output = {
			address: settings.protocol + settings.vpu_id + settings.base_url + settings.livesource_id + '?token=' + settings.rtmp_token,
			options: options,
		};

		return output;
	};

	return (
		<Grid container spacing={2}>
			<Grid item xs={12} md={6}>
				<TextField variant="outlined" fullWidth label={<Trans>VPU ID</Trans>} value={settings.vpu_id} onChange={handleChange('vpu_id')} />
			</Grid>
			<Grid item xs={12} md={6}>
				<TextField
					variant="outlined"
					fullWidth
					label={<Trans>Livesource ID</Trans>}
					value={settings.livesource_id}
					onChange={handleChange('livesource_id')}
				/>
			</Grid>
			<Grid item xs={12} md={12}>
				<TextField
					variant="outlined"
					fullWidth
					label={<Trans>Security token</Trans>}
					value={settings.rtmp_token}
					onChange={handleChange('rtmp_token')}
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
