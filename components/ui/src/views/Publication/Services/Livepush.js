import React from 'react';

import { Trans } from '@lingui/macro';
import Grid from '@mui/material/Grid';
import Link from '@mui/material/Link';
import MenuItem from '@mui/material/MenuItem';
import TextField from '@mui/material/TextField';

import Logo from './logos/livepush.svg';
import Select from '../../../misc/Select';

const id = 'livepush';
const name = 'Livepush';
const version = '1.0';
const stream_key_link = 'https://livepush.io';
const description = (
	<Trans>
		Transmit the main source to the Livepush RTMP Service. More details about the settings can be found{' '}
		<Link color="secondary" target="_blank" href="https://docs.livepush.io/en/">
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
	return <img src={Logo} alt="Livepush Logo" {...props} />;
}

function init(settings) {
	const initSettings = {
		region: 'dc-global',
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
			address: 'rtmp://' + settings.region + '.livepush.io/live/' + settings.stream_key,
			options: ['-f', 'flv'],
		};

		return output;
	};

	return (
		<Grid container spacing={2}>
			<Grid item xs={12}>
				<Select label={<Trans>Region</Trans>} value={settings.region} onChange={handleChange('region')}>
					<MenuItem value="dc-global">Global</MenuItem>
					<MenuItem value="us-central-ch">Chicago, US</MenuItem>
					<MenuItem value="us-east-ny">New York, US</MenuItem>
					<MenuItem value="us-west-la">Los Angeles, US</MenuItem>
					<MenuItem value="us-south-mia">Miami, US</MenuItem>
					<MenuItem value="us-central-dal">Dallas, US</MenuItem>
					<MenuItem value="ca-central-mon">Montreal, CA</MenuItem>
					<MenuItem value="ca-south-tor">Toronto, CA</MenuItem>
					<MenuItem value="au-east-syd">Sydney, AU</MenuItem>
					<MenuItem value="uk-central-ldn">London, UK</MenuItem>
					<MenuItem value="it-north-mln">Milan, IT</MenuItem>
					<MenuItem value="fr-central-par">Paris, FR</MenuItem>
					<MenuItem value="as-southeast-sg">Singapore</MenuItem>
					<MenuItem value="in-south-blr">Bangalore, IN</MenuItem>
					<MenuItem value="tr-north-ist">Turkiye</MenuItem>
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
