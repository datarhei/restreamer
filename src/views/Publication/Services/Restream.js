import React from 'react';

import { Trans } from '@lingui/macro';
import Grid from '@mui/material/Grid';
import Link from '@mui/material/Link';
import MenuItem from '@mui/material/MenuItem';
import TextField from '@mui/material/TextField';

import Logo from './logos/restream.svg';
import FormInlineButton from '../../../misc/FormInlineButton';
import Select from '../../../misc/Select';

const id = 'restream';
const name = 'Restream';
const version = '1.0';
const stream_key_link = 'https://restream.io/settings/streaming-setup?from=datarhei/restreamer';
const description = (
	<Trans>
		Transmit the main source to the Restream RTMP Service. More details about the settings can be found{' '}
		<Link color="secondary" target="_blank" href="https://restream.io">
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
	protocols: ['rtmps'],
	formats: ['flv'],
	codecs: {
		audio: ['aac', 'mp3'],
		video: ['h264'],
	},
};

function ServiceIcon(props) {
	return <img src={Logo} alt="restream.io Logo" {...props} />;
}

function init(settings) {
	const initSettings = {
		region: 'live',
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
			address: 'rtmp://' + settings.region + '.restream.io/live/' + settings.stream_key,
			options: ['-f', 'flv'],
		};

		return output;
	};

	return (
		<Grid container spacing={2}>
			<Grid item xs={12}>
				<Select label={<Trans>Region</Trans>} value={settings.region} onChange={handleChange('region')}>
					<MenuItem value="live">Autodetect</MenuItem>
					<MenuItem value="london">EU-West - London, GB</MenuItem>
					<MenuItem value="amsterdam">EU-West - Amsterdam, NL</MenuItem>
					<MenuItem value="luxembourg">EU-West - Luxembourg</MenuItem>
					<MenuItem value="paris">EU-West - Paris, FR</MenuItem>
					<MenuItem value="milan">EU-West - Milan, IT</MenuItem>
					<MenuItem value="frankfurt">EU-Central - Frankfurt, DE</MenuItem>
					<MenuItem value="falkenstein">EU-East - Falkenstein, DE</MenuItem>
					<MenuItem value="prague">EU-East - Prague, Czech</MenuItem>
					<MenuItem value="madrid">EU-South - Madrid, Spain</MenuItem>
					<MenuItem value="moscow">Russia - Moscow</MenuItem>
					<MenuItem value="istanbul">Turkey - Istanbul</MenuItem>
					<MenuItem value="telaviv">Israel - Tel Aviv</MenuItem>
					<MenuItem value="seattle">US-West - Seattle, WA</MenuItem>
					<MenuItem value="sanjose">US-West - San Jose, CA</MenuItem>
					<MenuItem value="dallas">US-Central - Dallas, TX</MenuItem>
					<MenuItem value="washington">US-East - Washington, DC</MenuItem>
					<MenuItem value="miami.restream.io/live">US-East - Miami, FL</MenuItem>
					<MenuItem value="chicago">US-East - Chicago, IL</MenuItem>
					<MenuItem value="toronto">NA-East - Toronto, Canada</MenuItem>
					<MenuItem value="saopaulo">SA - Saint Paul, Brazil</MenuItem>
					<MenuItem value="bangalore">India - Bangalore</MenuItem>
					<MenuItem value="singapore">Asia - Singapore</MenuItem>
					<MenuItem value="seoul">Asia - Seoul, South Korea</MenuItem>
					<MenuItem value="tokyo">Asia - Tokyo, Japan</MenuItem>
					<MenuItem value="sydney">Australia - Sydney</MenuItem>
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
