import React from 'react';

import { Trans } from '@lingui/macro';
import urlparser from 'url-parse';
import Grid from '@mui/material/Grid';
import Link from '@mui/material/Link';
import TextField from '@mui/material/TextField';
import Typography from '@mui/material/Typography';

import Logo from './logos/datarhei.svg';
import BoxText from '../../../misc/BoxText';

const id = 'datarheicore';
const name = 'datarhei Core';
const version = '2.0';
const stream_key_link_rtmp = 'https://docs.datarhei.com/restreamer/knowledge-base/manual/system-settings/rtmp';
const stream_key_link_srt = 'https://docs.datarhei.com/restreamer/knowledge-base/manual/system-settings/srt';
const description = (
	<Trans>
		Transmit the main source to an datarhei Core Ressource. More details about the settings can be found{' '}
		<Link color="secondary" target="_blank" href="https://docs.datarhei.com/restreamer">
			here
		</Link>
		.
	</Trans>
);
const image_copyright = <Trans>Please contact the operator of the service and check what happens.</Trans>;
const author = {
	creator: {
		name: 'datarhei.com',
		link: 'https://datarhei.com',
	},
	maintainer: {
		name: 'datarhei.com',
		link: 'https://datarhei.com',
	},
};
const category = 'software';
const requires = {
	protocols: ['rtmp', 'rtmps', 'srt'],
	formats: ['flv', 'mpegts'],
	codecs: {
		audio: ['aac'],
		video: ['h264', 'hevc', 'vp9', 'av1'],
	},
};

function ServiceIcon(props) {
	return <img src={Logo} alt="datarhei.com Logo" {...props} />;
}

function init(settings) {
	// v1.0 > v2.0
	if (settings.base_url && !settings.v2_protocol) {
		if (settings.base_url.length !== 0) {
			settings.v2_protocol = settings.protocol;
			settings.v2_host = settings.base_url;
			if (settings.protocol === 'srt') {
				settings.v2_address =
					`srt://${settings.base_url}?mode=caller&transtype=live&streamid=#!:m=publish,r=${settings.stream_name}` +
					(settings.token.length !== 0 ? `,token=${settings.token}` : '') +
					(settings.srt_passphrase.length !== 0 ? `&passphrase=${settings.srt_passphrase}` : '');
				settings.v2_stream_id = `#!:m=publish,r=${settings.stream_name}` + (settings.token.length !== 0 ? `,token=${settings.token}` : '');
				settings.v2_passphrase = settings.srt_passphrase;
			} else {
				settings.v2_address =
					`${settings.protocol}://${settings.base_url}/` +
					(settings.app_path.length !== 0 ? `${settings.app_path}/` : '') +
					settings.stream_name +
					(settings.token.length !== 0 ? `?token=${settings.token}` : '');
				settings.v2_stream_id = `/` + (settings.app_path.length !== 0 ? `${settings.app_path}/` : '') + settings.stream_name;
				settings.v2_token = settings.token;
			}
		}
		settings.base_url = '';
		settings.app_path = '';
		settings.stream_name = '';
		settings.token = '';
		settings.srt_passphrase = '';
	}

	const initSettings = {
		v2_address: '',
		v2_protocol: '',
		v2_host: '',
		v2_stream_id: '',
		v2_token: '',
		v2_passphrase: '',
		...settings,
	};

	return initSettings;
}

function Service(props) {
	const settings = init(props.settings);

	const handleChange = (what) => (event) => {
		const value = event.target.value;

		settings[what] = value;

		const url = urlparser(value);

		if (url.protocol === 'rtmp:' || url.protocol === 'rtmps:' || url.protocol === 'srt:') {
			settings.v2_protocol = url.protocol.split(':')[0];
			settings.v2_host = url.host;
			if (url.protocol === 'srt:') {
				if (url.hash) {
					if (url.hash.includes('passphrase=')) {
						settings.v2_stream_id = url.hash.split('&')[0];
						settings.v2_passphrase = url.hash.split('&passphrase=')[1].split('&')[0];
					} else {
						settings.v2_stream_id = url.hash.split('&')[0];
						settings.v2_passphrase = '';
					}
				} else if (url.query) {
					if (url.query.includes('passphrase=')) {
						settings.v2_stream_id = url.query.split('streamid=')[1].split('&')[0];
						settings.v2_passphrase = url.query.split('&passphrase=')[1].split('&')[0];
					} else {
						settings.v2_stream_id = url.query.split('streamid=')[1];
						settings.v2_passphrase = '';
					}
				} else {
					settings.v2_stream_id = '';
					settings.v2_passphrase = '';
				}
			} else {
				settings.v2_stream_id = url.pathname;
				if (url.query) {
					if (url.query.includes('token=')) {
						settings.v2_token = url.query.split('token=')[1].split('&')[0];
					} else if (url.query.includes('token:')) {
						settings.v2_token = url.query.split('token:')[1].split('&')[0];
					} else {
						settings.v2_token = '';
					}
				} else {
					settings.v2_token = '';
				}
			}
		} else {
			settings.v2_protocol = '';
			settings.v2_host = '';
			settings.v2_port = '';
			settings.v2_streamid = '';
			settings.v2_token = '';
			settings.v2_passphrase = '';
		}

		const output = createOutput(settings);

		props.onChange([output], settings);
	};

	const createOutput = (settings) => {
		const output = {
			address: null,
			options: null,
		};

		output.address = settings.v2_address;
		if (settings.v2_address.includes('srt://')) {
			output.options = ['-bsf:v', 'dump_extra', '-f', 'mpegts'];
		} else {
			output.options = ['-f', 'flv'];

			if (props.skills.ffmpeg.version_major >= 6) {
				const codecs = [];
				if (props.skills.codecs.video.includes('hevc')) {
					codecs.push('hvc1');
				}
				if (props.skills.codecs.video.includes('av1')) {
					codecs.push('av01');
				}
				if (props.skills.codecs.video.includes('vp9')) {
					codecs.push('vp09');
				}

				if (codecs.length !== 0) {
					output.options.push('-rtmp_enhanced_codecs', codecs.join(','));
				}
			}
		}

		return output;
	};

	return (
		<Grid container spacing={2}>
			<Grid item xs={12}>
				<TextField
					variant="outlined"
					fullWidth
					label={<Trans>Target address</Trans>}
					value={settings.v2_address}
					onChange={handleChange('v2_address')}
				/>
			</Grid>
			<Grid item xs={12} align="left">
				<BoxText alignItems="left" justifyContent="left">
					{!settings.v2_address && (
						<Typography>
							<strong>
								<Trans>Restreamer instructions</Trans>:
							</strong>
							<br />
							<Trans>1. Switch to the interface of the target Restreamer.</Trans>
							<br />
							<Trans>2. Create a new channel and select RTMP or SRT server.</Trans>
							<br />
							<Trans>3. Copy the URL and paste it in the "Target address" field.</Trans>
						</Typography>
					)}
					{settings.v2_address && (
						<Typography>
							<strong>
								<Trans>Protocol</Trans>:
							</strong>{' '}
							{settings.v2_protocol}
							<br />
							<strong>
								<Trans>Address</Trans>:
							</strong>{' '}
							{settings.v2_host}
							<br />
							<strong>Stream ID:</strong> {settings.v2_stream_id}
							<br />
							{settings.v2_token && (
								<React.Fragment>
									<strong>Token:</strong> {settings.v2_token}
								</React.Fragment>
							)}
							{settings.v2_passphrase && (
								<React.Fragment>
									<strong>
										<Trans>Passphrase</Trans>:
									</strong>{' '}
									{settings.v2_passphrase}
								</React.Fragment>
							)}
						</Typography>
					)}
				</BoxText>
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
	stream_key_link_rtmp,
	stream_key_link_srt,
	description,
	image_copyright,
	author,
	category,
	requires,
	ServiceIcon as icon,
	Service as component,
};
