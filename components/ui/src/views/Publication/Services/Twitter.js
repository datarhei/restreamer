import React from 'react';

import { faTwitter } from '@fortawesome/free-brands-svg-icons';
import { v4 as uuidv4 } from 'uuid';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Trans } from '@lingui/macro';
import Grid from '@mui/material/Grid';
import Link from '@mui/material/Link';
import MenuItem from '@mui/material/MenuItem';
import TextField from '@mui/material/TextField';

import FormInlineButton from '../../../misc/FormInlineButton';
import Select from '../../../misc/Select';

const id = 'twitter';
const name = 'Twitter';
const version = '1.0';
const stream_key_link = 'https://studio.twitter.com/producer/sources';
const description = (
	<Trans>
		Transmits your video stream with the required key, which was generated in Twitter Producer. You can find more information on seting up a live stream at
		Twitter's{' '}
		<Link color="secondary" target="_blank" href="https://help.twitter.com/en/using-twitter/how-to-use-live-producer">
			Producer
		</Link>
		.
	</Trans>
);
const image_copyright = (
	<Trans>
		More about Twitter's copyright{' '}
		<Link color="secondary" target="_blank" href="https://help.twitter.com/en/rules-and-policies/copyright-policy">
			here
		</Link>
		.
	</Trans>
);

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
	protocols: ['rtmps', 'https'],
	formats: ['flv', 'hls'],
	codecs: {
		audio: ['aac', 'mp3'],
		video: ['h264'],
	},
};

function ServiceIcon(props) {
	return <FontAwesomeIcon icon={faTwitter} style={{ color: 'rgb(29, 161, 242)' }} {...props} />;
}

function init(settings) {
	const initSettings = {
		mode: 'rtmps',
		stream_key: '',
		region: 'de',
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

		if (settings.stream_key.length === 0) {
			return outputs;
		}

		if (settings.mode === 'rtmps') {
			// https://help.twitter.com/en/using-twitter/how-to-use-live-producer#RTMP
			outputs.push({
				address: 'rtmps://' + settings.region + '.pscp.tv:443/x/' + settings.stream_key,
				options: ['-f', 'flv'],
			});
		} else if (settings.mode === 'hls') {
			// https://help.twitter.com/en/using-twitter/how-to-use-live-producer#HLS
			const name = uuidv4();
			const options = [
				'-f',
				'hls',
				'-start_number',
				'0',
				'-hls_time',
				'2',
				'-hls_delete_threshold',
				'3',
				'-hls_list_size',
				'5',
				'-hls_flags',
				'append_list',
				'-hls_segment_type',
				'mpegts',
				'-http_persistent',
				'1',
				'-y',
				'-method',
				'PUT',
			];

			const base = 'https://' + settings.region + '.pscp.tv:443/x/' + settings.stream_key;
			outputs.push({
				address: base + name + '.m3u8',
				options: [...options, '-hls_segment_filename', base + name + '_%d.ts'],
			});
		}

		return outputs;
	};

	const allowRTMPS = props.skills.protocols.includes('rtmps') && props.skills.formats.includes('flv');

	return (
		<Grid container spacing={2}>
			<Grid item xs={12}>
				<Select label={<Trans>Delivering mode</Trans>} value={settings.mode} onChange={handleChange('mode')}>
					{allowRTMPS === true && <MenuItem value="rtmps">RTMPS</MenuItem>}
					{/* 2021-05-20 Endpoints currently not visible
					{allowHLS === true && <MenuItem value="hls">HLS</MenuItem>}
					*/}
				</Select>
			</Grid>
			<Grid item xs={12}>
				{settings.mode === 'rtmps' && (
					<Select label={<Trans>Region</Trans>} value={settings.region} onChange={handleChange('region')}>
						<MenuItem value="jp">Asia - Tokyo, Japan</MenuItem>
						<MenuItem value="kr">Asia - Seoul, South Korea</MenuItem>
						<MenuItem value="in">Asia - Mumbai, India</MenuItem>
						<MenuItem value="sg">Asia - Singapore, Singapore</MenuItem>
						<MenuItem value="au">Asia - Sydney, Australia</MenuItem>
						<MenuItem value="de">EU - Frankfurt, Germany</MenuItem>
						<MenuItem value="ie">EU - Dublin, Ireland</MenuItem>
						<MenuItem value="fr">EU - Paris, France</MenuItem>
						<MenuItem value="br">South America - São Paulo, Brazil</MenuItem>
						<MenuItem value="va">US East - Northern Virginia, USA</MenuItem>
						<MenuItem value="ca">US West - Oregon, USA</MenuItem>
						<MenuItem value="or">US West - Northern California, USA</MenuItem>
					</Select>
				)}
				{settings.mode === 'hls' && (
					<Select label={<Trans>Region</Trans>} value={settings.region} onChange={handleChange('region')}>
						<MenuItem value="jp">Asia Pacific - Tokyo, Japan</MenuItem>
						<MenuItem value="de">EU - Frankfurt, Germany</MenuItem>
						<MenuItem value="va">US East - Northern Virginia, USA</MenuItem>
						<MenuItem value="or">US West - Oregon, USA</MenuItem>
					</Select>
				)}
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
