import React from 'react';

import { Trans } from '@lingui/macro';
import Accordion from '@mui/material/Accordion';
import AccordionDetails from '@mui/material/AccordionDetails';
import AccordionSummary from '@mui/material/AccordionSummary';
import Divider from '@mui/material/Divider';
import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown';
import Grid from '@mui/material/Grid';
import Link from '@mui/material/Link';
import TextField from '@mui/material/TextField';
import Typography from '@mui/material/Typography';

import Logo from './logos/icecast.svg';

const id = 'icecast';
const name = 'Icecast';
const version = '1.0';
const stream_key_link = '';
const description = (
	<Trans>
		Transmit the audio channel of the main source to an Icecast Server. More details about the settings can be found{' '}
		<Link color="secondary" target="_blank" href="https://ffmpeg.org/ffmpeg-all.html#Icecast">
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
	protocols: ['icecast'],
	formats: ['adts', 'opus', 'mp3', 'ogg', 'webm'],
	codecs: {
		audio: ['aac', 'opus', 'vorbis', 'mp3'],
		video: ['none', 'vp9'],
	},
};

function ServiceIcon(props) {
	return <img src={Logo} alt="Icecat Logo" {...props} />;
}

function init(settings, metadata) {
	const initSettings = {
		protocol: 'icecast://',
		address: '',
		options: {},
		...settings,
	};

	initSettings.options = {
		ice_genre: '',
		ice_name: metadata.name,
		ice_description: metadata.description,
		ice_url: '',
		ice_public: true,
		user_agent: '',
		legacy_icecast: false,
		tls: false,
		...initSettings.options,
	};

	return initSettings;
}

function createOutputs(settings, skills, metadata, streams) {
	settings = init(settings, metadata);
	let hasVideo = false;
	let audioCodec = '';

	for (let i = 0; i < streams.length; i++) {
		if (streams[i].type === 'video') {
			hasVideo = true;
		} else if (streams[i].type === 'audio') {
			audioCodec = streams[i].codec;
		}
	}

	const options = [];

	for (let key in settings.options) {
		if (settings.options[key].length === 0) {
			continue;
		}

		options.push('-' + key, String(settings.options[key]));
	}

	// https://gist.github.com/keiya/c8a5cbd4fe2594ddbb3390d9cf7dcac9#file-readme-md
	// https://wiki.xiph.org/Icecast_Server/Streaming_WebM_to_Icecast_with_FFmpeg

	if (hasVideo === true) {
		options.push('-f', 'webm', '-cluster_size_limit', '2', '-cluster_time_limit', '5100', '-content_type', 'video/webm');
	} else {
		switch (audioCodec) {
			case 'aac':
				options.push('-content_type', 'audio/aac', '-f', 'adts');
				break;
			case 'vorbis':
				options.push('-content_type', 'audio/ogg', '-f', 'ogg');
				break;
			case 'opus':
				options.push('-content_type', 'audio/ogg', '-f', 'opus');
				break;
			case 'mp3':
				options.push('-content_type', 'audio/mpeg', '-f', 'mp3');
				break;
			default:
				break;
		}
	}

	const output = {
		address: settings.protocol + settings.address,
		options: options,
	};

	return [output];
}

function Service(props) {
	const settings = init(props.settings, props.metadata);

	const handleChange = (what) => (event) => {
		const value = event.target.value;

		if (what in settings.options) {
			if (['ice_public', 'legacy_icecast', 'tls'].includes(what)) {
				settings.options[what] = !settings.options[what];
			} else {
				settings.options[what] = value;
			}
		} else {
			settings[what] = value;
		}

		const outputs = createOutputs(settings, props.skills, props.metadata, props.streams);

		props.onChange(outputs, settings);
	};

	return (
		<Grid container spacing={2}>
			<Grid item xs={12} md={3}>
				<TextField variant="outlined" fullWidth label={<Trans>Protocol</Trans>} value={settings.protocol} readOnly disabled />
			</Grid>
			<Grid item xs={12} md={9}>
				<TextField
					variant="outlined"
					fullWidth
					label={<Trans>Address</Trans>}
					value={settings.address}
					placeholder="username:password@server:port/mountpoint"
					onChange={handleChange('address')}
				/>
			</Grid>
			<Grid item xs={12}>
				<Accordion className="accordion">
					<AccordionSummary className="accordion-summary" elevation={0} expandIcon={<ArrowDropDownIcon />}>
						<Typography>
							<Trans>Advanced settings</Trans>
						</Typography>
					</AccordionSummary>
					<AccordionDetails>
						<Grid container spacing={2}>
							<Grid item xs={12}>
								<Typography variant="h3">
									<Trans>General</Trans>
								</Typography>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									type="text"
									fullWidth
									label="ice_name"
									value={settings.options.ice_name}
									onChange={handleChange('ice_name')}
									disabled
								/>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									type="text"
									fullWidth
									label="ice_genre"
									value={settings.options.ice_genre}
									onChange={handleChange('ice_genre')}
								/>
							</Grid>
							<Grid item xs={12}>
								<TextField
									multiline
									variant="outlined"
									fullWidth
									rows={4}
									label="ice_description"
									value={settings.options.ice_description}
									onChange={handleChange('ice_description')}
									disabled
								/>
							</Grid>
							<Grid item xs={12}>
								<TextField
									variant="outlined"
									type="text"
									fullWidth
									label="ice_url"
									value={settings.options.ice_url}
									onChange={handleChange('ice_url')}
								/>
							</Grid>
							<Grid item xs={12}>
								<TextField
									variant="outlined"
									type="text"
									fullWidth
									label="user_agent"
									value={settings.options.user_agent}
									onChange={handleChange('user_agent')}
								/>
							</Grid>
							<Grid item xs={12}>
								<Divider />
							</Grid>
							<Grid item xs={12}>
								<Typography>
									<Trans>Documentation</Trans>{' '}
									<Link color="secondary" target="_blank" href="http://ffmpeg.org/ffmpeg-all.html#Icecast">
										http://ffmpeg.org/ffmpeg-all.html#Icecast
									</Link>
								</Typography>
							</Grid>
						</Grid>
					</AccordionDetails>
				</Accordion>
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

export {
	id,
	name,
	version,
	stream_key_link,
	description,
	image_copyright,
	author,
	category,
	requires,
	ServiceIcon as icon,
	Service as component,
	createOutputs,
};
