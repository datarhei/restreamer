import React from 'react';

import { faYoutube } from '@fortawesome/free-brands-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Trans } from '@lingui/macro';
import { v4 as uuidv4 } from 'uuid';
import Grid from '@mui/material/Grid';
import Link from '@mui/material/Link';
import MenuItem from '@mui/material/MenuItem';
import TextField from '@mui/material/TextField';

import Checkbox from '../../../misc/Checkbox';
import FormInlineButton from '../../../misc/FormInlineButton';
import Select from '../../../misc/Select';

const id = 'youtube';
const name = 'YouTube Live';
const version = '1.1';
const stream_key_link = 'https://www.youtube.com/live_dashboard';
const description = (
	<Trans>
		Transmits your video as an RTMP stream with the required key generated in YouTube Studio. You can find more information on setting up a live stream at
		YouTube's{' '}
		<Link color="secondary" target="_blank" href="https://creatoracademy.youtube.com/">
			Creator Academy
		</Link>
		.
	</Trans>
);
const image_copyright = (
	<Trans>
		More about YouTube's copyright{' '}
		<Link color="secondary" target="_blank" href="https://www.youtube.com/howyoutubeworks/policies/copyright/#support-and-troubleshooting">
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
		video: ['h264', 'hevc', 'av1'],
	},
};

function ServiceIcon(props) {
	return <FontAwesomeIcon icon={faYoutube} style={{ color: '#FF0000' }} {...props} />;
}

function init(settings) {
	const initSettings = {
		mode: 'rtmps',
		stream_key: '',
		primary: true,
		backup: false,
		...settings,
	};

	return initSettings;
}

function Service(props) {
	const settings = init(props.settings);

	const handleChange = (what) => (event) => {
		const value = event.target.value;

		if (['primary', 'backup'].includes(what)) {
			settings[what] = !settings[what];
		} else {
			settings[what] = value;
		}

		const outputs = createOutput(settings);

		props.onChange(outputs, settings);
	};

	const createOutput = (settings) => {
		const outputs = [];

		if (settings.stream_key.length === 0) {
			return outputs;
		}

		if (settings.mode === 'rtmps') {
			let options = ['-f', 'flv'];

			if (props.skills.ffmpeg.version_major >= 6) {
				const codecs = [];
				if (props.skills.codecs.video.includes('hevc')) {
					codecs.push('hvc1');
				}
				if (props.skills.codecs.video.includes('av1')) {
					codecs.push('av01');
				}

				if (codecs.length !== 0) {
					options.push('-rtmp_enhanced_codecs', codecs.join(','));
				}
			}

			// https://developers.google.com/youtube/v3/live/guides/rtmps-ingestion
			if (settings.primary === true) {
				outputs.push({
					address: 'rtmps://a.rtmp.youtube.com/live2/' + settings.stream_key,
					options: options.slice(),
				});
			}

			if (settings.backup === true) {
				outputs.push({
					address: 'rtmps://b.rtmp.youtube.com/live2?backup=1/' + settings.stream_key,
					options: options.slice(),
				});
			}
		} else if (settings.mode === 'hls') {
			// https://developers.google.com/youtube/v3/live/guides/hls-ingestion
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

			if (settings.primary === true) {
				const base = `https://a.upload.youtube.com/http_upload_hls?cid=${settings.stream_key}&copy=0&file=`;
				outputs.push({
					address: base + name + '.m3u8',
					options: [...options, '-hls_segment_filename', base + name + '_%d.ts'],
				});
			}

			if (settings.backup === true) {
				const base = `https://b.upload.youtube.com/http_upload_hls?cid=${settings.stream_key}&copy=1&file=`;
				outputs.push({
					address: base + name + '.m3u8',
					options: [...options, '-hls_segment_filename', base + name + '_%d.ts'],
				});
			}
		}

		return outputs;
	};

	const allowRTMPS = props.skills.protocols.includes('rtmps') && props.skills.formats.includes('flv');
	const allowHLS = props.skills.protocols.includes('https') && props.skills.formats.includes('hls');

	return (
		<Grid container spacing={2}>
			<Grid item xs={12}>
				<Select label={<Trans>Delivering mode</Trans>} value={settings.mode} onChange={handleChange('mode')}>
					{allowRTMPS === true && <MenuItem value="rtmps">RTMP</MenuItem>}
					{allowHLS === true && <MenuItem value="hls">HLS</MenuItem>}
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
			<Grid item xs={12}>
				<Checkbox label={<Trans>Primary stream</Trans>} checked={settings.primary} onChange={handleChange('primary')} />
				<Checkbox label={<Trans>Backup stream</Trans>} checked={settings.backup} onChange={handleChange('backup')} />
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
