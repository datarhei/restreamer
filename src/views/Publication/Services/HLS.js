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
import MultiSelect from '../../../misc/MultiSelect';
import MultiSelectOption from '../../../misc/MultiSelectOption';
import Password from '../../../misc/Password';

const id = 'hls';
const name = 'HLS';
const version = '1.0';
const stream_key_link = '';
const description = (
	<Trans>
		Transmit the main source as HTTP-Live-Streaming (HLS) to an HTTP/S Server. More details about the settings can be found{' '}
		<Link color="secondary" target="_blank" href="http://ffmpeg.org/ffmpeg-all.html#hls">
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
	formats: ['hls'],
	codecs: {
		audio: ['aac', 'mp3'],
		video: ['h264'],
	},
};

function ServiceIcon(props) {
	return <FontAwesomeIcon icon={faTools} style={{ color: '#39B54A' }} {...props} />;
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
		hls_init_time: '0',
		hls_time: '2',
		hls_list_size: '5',
		hls_delete_threshold: '1',
		hls_ts_options: '',
		hls_start_number_source: 'generic',
		start_number: '0',
		hls_allow_cache: false,
		hls_base_url: '',
		hls_segment_filename: '',
		strftime_mkdir: '',
		hls_key_info_file: '',
		hls_enc: false,
		hls_enc_key: '',
		hls_enc_key_url: '',
		hls_enc_iv: '',
		hls_segment_type: 'mpegts',
		hls_fmp4_init_filename: 'init.mp4',
		hls_fmp4_init_resend: '0',
		hls_flags: [],
		hls_playlist_type: '',
		method: 'PUT',
		http_user_agent: '',
		var_stream_map: '',
		master_pl_name: '',
		master_pl_publish_rate: '',
		http_persistent: false,
		timeout: '',
		ignore_io_errors: false,
		headers: '',
		...initSettings.options,
	};

	return initSettings;
}

function Service(props) {
	const settings = init(props.settings);

	const handleChange = (what) => (event) => {
		let value = event.target.value;

		if (what in settings.options) {
			if (['hls_allow_cache', 'http_persistent', 'ignore_io_errors', 'hls_enc'].includes(what)) {
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
		const options = ['-f', 'hls'];

		for (let key in settings.options) {
			if (settings.options[key].length === 0) {
				continue;
			}

			if (typeof settings.options[key] === 'boolean') {
				options.push('-' + key, settings.options[key] === true ? '1' : '0');
			} else {
				if (key === 'hls_flags') {
					options.push('-' + key, settings.options[key].join('+'));
				} else {
					options.push('-' + key, '' + settings.options[key]);
				}
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
									label="headers"
									value={settings.options.headers}
									onChange={handleChange('headers')}
								/>
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
									<MenuItem value="">default</MenuItem>
									<MenuItem value="PUT">PUT</MenuItem>
									<MenuItem value="POST">POST</MenuItem>
								</Select>
							</Grid>
							<Grid item xs={12}>
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
							<Grid item xs={12}>
								<Checkbox label="hls_allow_cache" checked={settings.options.hls_allow_cache} onChange={handleChange('hls_allow_cache')} />
								<Checkbox label="http_persistent" checked={settings.options.http_persistent} onChange={handleChange('http_persistent')} />
								<Checkbox label="ignore_io_errors" checked={settings.options.ignore_io_errors} onChange={handleChange('ignore_io_errors')} />
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
									label="hls_init_time"
									value={settings.options.hls_init_time}
									onChange={handleChange('hls_init_time')}
								/>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="0"
									step="1"
									label="hls_time"
									value={settings.options.hls_time}
									onChange={handleChange('hls_time')}
								/>
							</Grid>
							<Grid item xs={12}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="hls_ts_options"
									value={settings.options.hls_ts_options}
									onChange={handleChange('hls_ts_options')}
								/>
							</Grid>
							<Grid item xs={12}>
								<Typography variant="h3">
									<Trans>Playlist</Trans>
								</Typography>
							</Grid>
							<Grid item xs={12}>
								<Select
									type="select"
									label="hls_playlist_type"
									value={settings.options.hls_playlist_type}
									onChange={handleChange('hls_playlist_type')}
								>
									<MenuItem value="">none</MenuItem>
									<MenuItem value="event">event</MenuItem>
									<MenuItem value="vod">vod</MenuItem>
								</Select>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="0"
									step="1"
									label="hls_list_size"
									value={settings.options.hls_list_size}
									onChange={handleChange('hls_list_size')}
								/>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="0"
									step="1"
									label="hls_delete_threshold"
									value={settings.options.hls_delete_threshold}
									onChange={handleChange('hls_delete_threshold')}
								/>
							</Grid>
							<Grid item xs={12}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="hls_base_url"
									value={settings.options.hls_base_url}
									onChange={handleChange('hls_base_url')}
								/>
							</Grid>
							<Grid item xs={12}>
								<MultiSelect type="select" label="hls_flags" value={settings.options.hls_flags} onChange={handleChange('hls_flags')}>
									<MultiSelectOption value="single_file" name="single_file" />
									<MultiSelectOption value="delete_segments" name="delete_segments" />
									<MultiSelectOption value="append_list" name="append_list" />
									<MultiSelectOption value="round_durations" name="round_durations" />
									<MultiSelectOption value="discont_start" name="discont_start" />
									<MultiSelectOption value="omit_endlist" name="omit_endlist" />
									<MultiSelectOption value="periodic_rekey" name="periodic_rekey" />
									<MultiSelectOption value="independent_segments" name="independent_segments" />
									<MultiSelectOption value="iframes_only" name="iframes_only" />
									<MultiSelectOption value="split_by_time" name="split_by_time" />
									<MultiSelectOption value="program_date_time" name="program_date_time" />
									<MultiSelectOption value="second_level_segment_index" name="second_level_segment_index" />
									<MultiSelectOption value="second_level_segment_size" name="second_level_segment_size" />
									<MultiSelectOption value="second_level_segment_duration" name="second_level_segment_duration" />
									<MultiSelectOption value="temp_file" name="temp_file" />
								</MultiSelect>
							</Grid>
							<Grid item xs={12}>
								<Typography variant="h3">
									<Trans>Segmentation</Trans>
								</Typography>
							</Grid>
							<Grid item xs={12}>
								<Select
									type="select"
									label="hls_segment_type"
									value={settings.options.hls_segment_type}
									onChange={handleChange('hls_segment_type')}
								>
									<MenuItem value="mpegts">mpegts</MenuItem>
									<MenuItem value="fmp4">fmp4</MenuItem>
								</Select>
							</Grid>
							<Grid item xs={12}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="hls_segment_filename"
									value={settings.options.hls_segment_filename}
									onChange={handleChange('hls_segment_filename')}
								/>
							</Grid>
							{settings.options.hls_segment_type === 'fmp4' && (
								<Grid item xs={6}>
									<TextField
										variant="outlined"
										fullWidth
										type="text"
										label="hls_fmp4_init_filename"
										value={settings.options.hls_fmp4_init_filename}
										onChange={handleChange('hls_fmp4_init_filename')}
									/>
								</Grid>
							)}
							{settings.options.hls_segment_type === 'fmp4' && (
								<Grid item xs={6}>
									<TextField
										variant="outlined"
										fullWidth
										type="number"
										min="0"
										step="1"
										label="hls_fmp4_init_resend"
										value={settings.options.hls_fmp4_init_resend}
										onChange={handleChange('hls_fmp4_init_resend')}
									/>
								</Grid>
							)}
							<Grid item xs={12} md={6}>
								<Select
									type="select"
									label="hls_start_number_source"
									value={settings.options.hls_start_number_source}
									onChange={handleChange('hls_start_number_source')}
								>
									<MenuItem value="generic">generic</MenuItem>
									<MenuItem value="epoch">epoch</MenuItem>
									<MenuItem value="epoch_us">epoch_us</MenuItem>
									<MenuItem value="datetime">datetime</MenuItem>
								</Select>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="0"
									step="1"
									label="start_number"
									value={settings.options.start_number}
									onChange={handleChange('start_number')}
								/>
							</Grid>
							<Grid item xs={12}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="strftime_mkdir"
									value={settings.options.strftime_mkdir}
									onChange={handleChange('strftime_mkdir')}
								/>
							</Grid>
							<Grid item xs={12}>
								<Typography variant="h3">
									<Trans>Encryption</Trans>
								</Typography>
							</Grid>
							<Grid item xs={12}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="hls_key_info_file"
									value={settings.options.hls_key_info_file}
									onChange={handleChange('hls_key_info_file')}
								/>
							</Grid>
							<Grid item xs={12}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="hls_enc_key"
									value={settings.options.hls_enc_key}
									onChange={handleChange('hls_enc_key')}
								/>
							</Grid>
							<Grid item xs={12}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="hls_enc_key_url"
									value={settings.options.hls_enc_key_url}
									onChange={handleChange('hls_enc_key_url')}
								/>
							</Grid>
							<Grid item xs={12}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="hls_enc_iv"
									value={settings.options.hls_enc_iv}
									onChange={handleChange('hls_enc_iv')}
								/>
							</Grid>
							<Grid item xs={12}>
								<Checkbox label="hls_enc" checked={settings.options.hls_enc} onChange={handleChange('hls_enc')} />
							</Grid>
							<Grid item xs={12}>
								<Divider />
							</Grid>
							<Grid item xs={12}>
								<Typography>
									<Trans>Documentation</Trans>{' '}
									<Link color="secondary" target="_blank" href="https://ffmpeg.org/ffmpeg-all.html#hls">
										https://ffmpeg.org/ffmpeg-all.html#hls
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
