import React from 'react';
import urlparser from 'url-parse';

import { faTools } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Trans } from '@lingui/macro';
import Accordion from '@mui/material/Accordion';
import AccordionDetails from '@mui/material/AccordionDetails';
import AccordionSummary from '@mui/material/AccordionSummary';
import Divider from '@mui/material/Divider';
import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown';
import Grid from '@mui/material/Grid';
import Link from '@mui/material/Link';
import MenuItem from '@mui/material/MenuItem';
import TextField from '@mui/material/TextField';
import Typography from '@mui/material/Typography';

import Checkbox from '../../../misc/Checkbox';
import Select from '../../../misc/Select';
import Password from '../../../misc/Password';

const id = 'dash';
const name = 'DASH';
const version = '1.0';
const stream_key_link = '';
const description = (
	<Trans>
		Transmit the main source as MPEG-DASH to an HTTP/S Server. More details about the settings can be found{' '}
		<Link color="secondary" target="_blank" href="http://ffmpeg.org/ffmpeg-all.html#dash">
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
const category = 'universal';
const requires = {
	protocols: ['http', 'https'],
	formats: ['dash'],
	codecs: {
		audio: ['none'],
		video: ['h264'],
	},
};

function ServiceIcon(props) {
	return <FontAwesomeIcon icon={faTools} style={{ color: 'rgba(57, 181, 74, 1)' }} {...props} />;
}

function init(settings) {
	const initSettings = {
		protocol: 'https://',
		address: '',
		username: '',
		password: '',
		options: {},
		...settings,
	};

	initSettings.options = {
		seg_duration: '2',
		frag_duration: '',
		frag_type: '',
		window_size: '6',
		extra_window_size: '',
		remove_at_exit: false,
		use_template: true,
		use_timeline: true,
		single_file: false,
		single_file_name: '',
		init_seg_name: 'init-stream$RepresentationID$.$ext$',
		media_seg_name: 'chunk-stream$RepresentationID$-$Number%05d$.$ext$',
		utc_timing_url: '',
		method: 'PUT',
		http_user_agent: '',
		http_persistent: false,
		hls_playlist: false,
		streaming: true,
		adaptation_sets: 'id=0,streams=v id=1,streams=a',
		timeout: '',
		index_correction: true,
		format_options: '',
		global_sidx: '',
		dash_segment_type: 'auto',
		ignore_io_errors: false,
		lhls: false,
		ldash: false,
		write_prft: '',
		mpd_profile: '',
		http_opts: '',
		target_latency: '',
		min_playback_rate: '',
		max_playback_rate: '',
		update_period: '',
		...initSettings.options,
	};

	return initSettings;
}

function Service(props) {
	const settings = init(props.settings);

	const handleChange = (what) => (event) => {
		const value = event.target.value;

		if (what in settings.options) {
			if (
				[
					'remove_at_exit',
					'use_template',
					'use_timeline',
					'single_file',
					'http_persistent',
					'hls_playlist',
					'streaming',
					'index_correction',
					'ignore_io_errors',
					'lhls',
					'ldash',
				].includes(what)
			) {
				settings.options[what] = !settings.options[what];
			} else {
				settings.options[what] = value;
			}
		} else {
			settings[what] = value;
		}

		const output = createOutput(settings);

		props.onChange([output], settings);
	};

	const createOutput = (settings) => {
		const options = ['-strict', 'experimental', '-f', 'dash'];

		for (let key in settings.options) {
			if (settings.options[key].length === 0) {
				continue;
			}
			if (typeof settings.options[key] === 'boolean') {
				options.push('-' + key, Number(settings.options[key]));
			} else {
				options.push('-' + key, String(settings.options[key]));
			}
		}

		let address = settings.protocol + settings.address;
		if (settings.username.length !== 0 || settings.password.length !== 0) {
			const url = urlparser(address);

			if (settings.username.length !== 0) {
				url.set('username', encodeURIComponent(settings.username));
			}

			if (settings.password.length !== 0) {
				url.set('password', encodeURIComponent(settings.password));
			}

			address = url.toString();
		}

		const output = {
			address: address,
			options: options,
		};

		return output;
	};

	return (
		<Grid container spacing={2}>
			<Grid item xs={12} md={3}>
				<Select type="select" label={<Trans>Protocol</Trans>} value={settings.protocol} onChange={handleChange('protocol')}>
					<MenuItem value="http://">http://</MenuItem>
					<MenuItem value="https://">https://</MenuItem>
				</Select>
			</Grid>
			<Grid item xs={12} md={9}>
				<TextField variant="outlined" fullWidth label={<Trans>Address</Trans>} value={settings.address} onChange={handleChange('address')} />
			</Grid>
			<Grid item xs={6}>
				<TextField variant="outlined" fullWidth label={<Trans>Username</Trans>} value={settings.username} onChange={handleChange('username')} />
			</Grid>
			<Grid item xs={6}>
				<Password variant="outlined" fullWidth label={<Trans>Password</Trans>} value={settings.password} onChange={handleChange('password')} />
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
							<Grid item xs={12}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="http_user_agent"
									value={settings.options.http_user_agent}
									onChange={handleChange('http_user_agent')}
								/>
							</Grid>
							<Grid item xs={12}>
								<Select type="select" label="method" value={settings.options.method} onChange={handleChange('method')}>
									<MenuItem value="PUT">PUT</MenuItem>
									<MenuItem value="POST">POST</MenuItem>
								</Select>
							</Grid>
							<Grid item xs={12}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="http_opts"
									value={settings.options.http_opts}
									onChange={handleChange('http_opts')}
								/>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="0"
									step="1"
									label="timeout"
									value={settings.options.timeout}
									onChange={handleChange('timeout')}
								/>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="0"
									step="1"
									label="target_latency"
									value={settings.options.target_latency}
									onChange={handleChange('target_latency')}
								/>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="0"
									step="1"
									label="min_playback_rate"
									value={settings.options.min_playback_rate}
									onChange={handleChange('min_playback_rate')}
								/>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="0"
									step="1"
									label="max_playback_rate"
									value={settings.options.max_playback_rate}
									onChange={handleChange('max_playback_rate')}
								/>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="0"
									step="1"
									label="update_period"
									value={settings.options.update_period}
									onChange={handleChange('update_period')}
								/>
							</Grid>
							<Grid item xs={12}>
								<Checkbox label="http_persistent" checked={settings.options.http_persistent} onChange={handleChange('http_persistent')} />
								<Checkbox label="streaming" checked={settings.options.streaming} onChange={handleChange('streaming')} />
								<Checkbox label="ignore_io_errors" checked={settings.options.ignore_io_errors} onChange={handleChange('ignore_io_errors')} />
								<Checkbox label="lhls" checked={settings.options.lhls} onChange={handleChange('lhls')} />
								<Checkbox label="ldash" checked={settings.options.ldash} onChange={handleChange('ldash')} />
							</Grid>
							<Grid item xs={12}>
								<Typography variant="h3">
									<Trans>Chunk</Trans>
								</Typography>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="0"
									step="1"
									label="seg_duration"
									value={settings.options.seg_duration}
									onChange={handleChange('seg_duration')}
								/>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="0"
									step="1"
									label="frag_duration"
									value={settings.options.frag_duration}
									onChange={handleChange('frag_duration')}
								/>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="frag_type"
									value={settings.options.frag_type}
									onChange={handleChange('frag_type')}
								/>
							</Grid>
							<Grid item xs={12}>
								<Typography variant="h3">
									<Trans>Playlist</Trans>
								</Typography>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="0"
									step="1"
									label="window_size"
									value={settings.options.window_size}
									onChange={handleChange('window_size')}
								/>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="0"
									step="1"
									label="extra_window_size"
									value={settings.options.extra_window_size}
									onChange={handleChange('extra_window_size')}
								/>
							</Grid>
							<Grid item xs={12}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="adaptation_sets"
									value={settings.options.adaptation_sets}
									onChange={handleChange('adaptation_sets')}
								/>
							</Grid>
							<Grid item xs={12}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="format_options"
									value={settings.options.format_options}
									onChange={handleChange('format_options')}
								/>
							</Grid>
							{!settings.options.streaming && (
								<Grid item xs={12}>
									<TextField
										variant="outlined"
										fullWidth
										type="text"
										label="global_sidx"
										value={settings.options.global_sidx}
										onChange={handleChange('global_sidx')}
									/>
								</Grid>
							)}
							<Grid item xs={12}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="mpd_profile"
									value={settings.options.mpd_profile}
									onChange={handleChange('mpd_profile')}
								/>
							</Grid>
							<Grid item xs={12}>
								{!settings.options.streaming && (
									<Checkbox label="single_file" checked={settings.options.single_file} onChange={handleChange('single_file')} />
								)}
								<Checkbox label="remove_at_exit" checked={settings.options.remove_at_exit} onChange={handleChange('remove_at_exit')} />
								<Checkbox label="use_template" checked={settings.options.use_template} onChange={handleChange('use_template')} />
								<Checkbox label="use_timeline" checked={settings.options.use_timeline} onChange={handleChange('use_timeline')} />
								<Checkbox label="index_correction" checked={settings.options.index_correction} onChange={handleChange('index_correction')} />
								<Checkbox label="hls_playlist" checked={settings.options.hls_playlist} onChange={handleChange('hls_playlist')} />
							</Grid>
							<Grid item xs={12}>
								<Typography variant="h3">
									<Trans>Segmentation</Trans>
								</Typography>
							</Grid>
							<Grid item xs={12}>
								<Select
									type="select"
									label="dash_segment_type"
									value={settings.options.dash_segment_type}
									onChange={handleChange('dash_segment_type')}
								>
									<MenuItem value="auto">auto</MenuItem>
									<MenuItem value="mp4">mp4</MenuItem>
									<MenuItem value="webm">webm</MenuItem>
								</Select>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="init_seg_name"
									value={settings.options.init_seg_name}
									onChange={handleChange('init_seg_name')}
								/>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="media_seg_name"
									value={settings.options.media_seg_name}
									onChange={handleChange('media_seg_name')}
								/>
							</Grid>
							{settings.options.single_file && (
								<Grid item xs={12} md={6}>
									<TextField
										variant="outlined"
										fullWidth
										type="text"
										label="single_file_name"
										value={settings.options.single_file_name}
										onChange={handleChange('single_file_name')}
									/>
								</Grid>
							)}
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="utc_timing_url"
									value={settings.options.utc_timing_url}
									onChange={handleChange('utc_timing_url')}
								/>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="write_prft"
									value={settings.options.write_prft}
									onChange={handleChange('write_prft')}
								/>
							</Grid>
							<Grid item xs={12}>
								<Divider />
							</Grid>
							<Grid item xs={12}>
								<Typography>
									<Trans>Documentation</Trans>{' '}
									<Link color="secondary" target="_blank" href="https://ffmpeg.org/ffmpeg-all.html#dash">
										https://ffmpeg.org/ffmpeg-all.html#dash
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

export { id, name, version, stream_key_link, description, image_copyright, author, category, requires, ServiceIcon as icon, Service as component };
