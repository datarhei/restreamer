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

import Select from '../../../misc/Select';
import MultiSelect from '../../../misc/MultiSelect';
import MultiSelectOption from '../../../misc/MultiSelectOption';
import Password from '../../../misc/Password';

const id = 'rtsp';
const name = 'RTSP';
const version = '1.0';
const stream_key_link = '';
const description = (
	<Trans>
		Transmit the main source to an RTSP Server. More details about the settings can be found{' '}
		<Link color="secondary" target="_blank" href="https://ffmpeg.org/ffmpeg-all.html#rtsp">
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
	protocols: ['rtp'],
	formats: ['rtsp'],
	codecs: {
		audio: ['aac', 'mp3'],
		video: ['h264', 'hevc', 'vp9'],
	},
};

function ServiceIcon(props) {
	return <FontAwesomeIcon icon={faTools} style={{ color: '#39B54A' }} {...props} />;
}

function init(settings) {
	const initSettings = {
		protocol: 'rtsp://',
		address: '',
		username: '',
		password: '',
		options: {},
		...settings,
	};

	initSettings.options = {
		rtsp_transport: '',
		rtsp_flags: '',
		allowed_media_types: [],
		reorder_queue_size: '',
		user_agent: '',
		muxdelay: '0.1',
		...initSettings.options,
	};

	return initSettings;
}

function Service(props) {
	const settings = init(props.settings);

	const handleChange = (what) => (event) => {
		const value = event.target.value;

		if (what in settings.options) {
			settings.options[what] = value;
		} else {
			settings[what] = value;
		}

		const output = createOutput(settings);

		props.onChange([output], settings);
	};

	const createOutput = (settings) => {
		const options = ['-f', 'rtsp'];

		for (let key in settings.options) {
			if (settings.options[key].length !== 0) {
				if (key !== 'user_agent') {
					options.push('-' + key, settings.options[key]);
				} else {
					options.push('-user-agent', settings.options.user_agent);
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
				<TextField
					variant="outlined"
					fullWidth
					type="url"
					label={<Trans>Protocol</Trans>}
					value={settings.protocol}
					onChange={handleChange('protocol')}
				/>
			</Grid>
			<Grid item xs={12} md={9}>
				<TextField variant="outlined" fullWidth type="url" label={<Trans>Address</Trans>} value={settings.address} onChange={handleChange('address')} />
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
								<Select type="select" label="rtsp_transport" value={settings.options.rtsp_transport} onChange={handleChange('rtsp_transport')}>
									<MenuItem value="">default</MenuItem>
									<MenuItem value="udp">udp</MenuItem>
									<MenuItem value="tcp">tcp</MenuItem>
									<MenuItem value="udp_multicast">udp_multicast</MenuItem>
									<MenuItem value="http">http</MenuItem>
								</Select>
							</Grid>
							<Grid item xs={12}>
								<Select type="select" label="rtsp_flags" value={settings.options.rtsp_flags} onChange={handleChange('rtsp_flags')}>
									<MenuItem value="">default</MenuItem>
									<MenuItem value="prefer_tcp">prefer_tcp</MenuItem>
								</Select>
							</Grid>
							<Grid item xs={12}>
								<MultiSelect
									type="select"
									label="allowed_media_types"
									value={settings.options.allowed_media_types}
									onChange={handleChange('allowed_media_types')}
								>
									<MultiSelectOption value="" name="all media types" />
									<MultiSelectOption value="video" name="video" />
									<MultiSelectOption value="audio" name="audio" />
									<MultiSelectOption value="data" name="data" />
								</MultiSelect>
							</Grid>
							<Grid item xs={12}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="reorder_queue_size"
									value={settings.reorder_queue_size}
									onChange={handleChange('reorder_queue_size')}
								/>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="user-agent"
									value={settings.user_agent}
									onChange={handleChange('user_agent')}
								/>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="max_delay"
									value={settings.max_delay}
									onChange={handleChange('max_delay')}
								/>
							</Grid>
							<Grid item xs={12}>
								<Divider />
							</Grid>
							<Grid item xs={12}>
								<Typography>
									<Trans>Documentation</Trans>{' '}
									<Link color="secondary" target="_blank" href="https://ffmpeg.org/ffmpeg-all.html#rtsp">
										https://ffmpeg.org/ffmpeg-all.html#rtsp
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
