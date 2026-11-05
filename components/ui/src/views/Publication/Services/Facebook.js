import React from 'react';

import { faFacebook } from '@fortawesome/free-brands-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Trans } from '@lingui/macro';
import Grid from '@mui/material/Grid';
import TextField from '@mui/material/TextField';

import Checkbox from '../../../misc/Checkbox';
import FormInlineButton from '../../../misc/FormInlineButton';

const id = 'facebook';
const name = 'Facebook Live';
const version = '1.0';
const stream_key_link = 'https://www.facebook.com/live/producer?ref=datarhei/restreamer';
const description = <Trans>Live-Streaming to Facebook Live RTMP service</Trans>;
const image_copyright = <Trans>More about licenses here</Trans>;
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
		audio: ['aac'],
		video: ['h264'],
	},
};

function ServiceIcon(props) {
	return <FontAwesomeIcon icon={faFacebook} style={{ color: '#2D88FF' }} {...props} />;
}

function init(settings) {
	const initSettings = {
		stream_key_primary: '',
		stream_key_backup: '',
		rtmp_primary: true,
		rtmp_backup: false,
		...settings,
	};

	return initSettings;
}

function Service(props) {
	const settings = init(props.settings);

	const handleChange = (what) => (event) => {
		const value = event.target.value;

		if (['rtmp_primary', 'rtmp_backup'].includes(what)) {
			settings[what] = !settings[what];
		} else {
			settings[what] = value;
		}

		const output = createOutput(settings);

		props.onChange(output, settings);
	};

	const createOutput = (settings) => {
		const outputs = [];

		const output_primary = {
			address: 'rtmps://live-api-s.facebook.com:443/rtmp/' + settings.stream_key_primary,
			options: ['-f', 'flv'],
		};

		const output_backup = {
			address: 'rtmps://live-api-s.facebook.com:443/rtmp/' + settings.stream_key_backup,
			options: ['-f', 'flv'],
		};

		if (settings.stream_key_primary.length !== 0) {
			if (settings.rtmp_primary) {
				outputs.push(output_primary);
			}
		}

		if (settings.stream_key_backup.length !== 0) {
			if (settings.rtmp_backup) {
				outputs.push(output_backup);
			}
		}

		return outputs;
	};

	return (
		<Grid container spacing={2}>
			{settings.rtmp_primary === true && (
				<Grid item xs={12} md={9}>
					<TextField
						variant="outlined"
						fullWidth
						label={<Trans>Primary stream key</Trans>}
						value={settings.stream_key_primary}
						onChange={handleChange('stream_key_primary')}
					/>
				</Grid>
			)}
			{settings.rtmp_primary === true && (
				<Grid item xs={12} md={3}>
					<FormInlineButton target="blank" href={stream_key_link} component="a">
						<Trans>GET</Trans>
					</FormInlineButton>
				</Grid>
			)}
			{settings.rtmp_backup === true && (
				<Grid item xs={12} md={9}>
					<TextField
						variant="outlined"
						fullWidth
						label={<Trans>Backup stream key</Trans>}
						value={settings.stream_key_backup}
						onChange={handleChange('stream_key_backup')}
					/>
				</Grid>
			)}
			{settings.rtmp_backup === true && (
				<Grid item xs={12} md={3}>
					<FormInlineButton target="blank" href={stream_key_link} component="a">
						<Trans>GET</Trans>
					</FormInlineButton>
				</Grid>
			)}
			<Grid item xs={12}>
				<Checkbox label={<Trans>Enable primary stream</Trans>} checked={settings.rtmp_primary} onChange={handleChange('rtmp_primary')} />
				<Checkbox label={<Trans>Enable backup stream</Trans>} checked={settings.rtmp_backup} onChange={handleChange('rtmp_backup')} />
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
