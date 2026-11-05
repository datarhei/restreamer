import React from 'react';

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

const id = 'rtmp';
const name = 'RTMP';
const version = '1.0';
const stream_key_link = '';
const description = (
	<Trans>
		Transmit the main source to an RTMP(e|s|t|te|ts) Server. More details about the settings can be found{' '}
		<Link color="secondary" target="_blank" href="https://ffmpeg.org/ffmpeg-all.html#rtmp">
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
	protocols: ['rtmp', 'rtmps', 'rtmpt', 'rtmpts'],
	formats: ['flv'],
	codecs: {
		audio: ['aac', 'mp3'],
		video: ['h264', 'hevc', 'vp9', 'av1'],
	},
};

function ServiceIcon(props) {
	return <FontAwesomeIcon icon={faTools} style={{ color: '#39B54A' }} {...props} />;
}

function init(settings) {
	const initSettings = {
		protocol: 'rtmp://',
		address: '',
		options: {},
		...settings,
	};

	initSettings.options = {
		rtmp_app: '',
		rtmp_playpath: '',
		rtmp_pageurl: '',
		rtmp_tcurl: '',
		rtmp_conn: '',
		rtmp_flashver: 'FMLE/3.0',
		rtmp_flush_interval: '10',
		rtmp_swfhash: '',
		rtmp_swfsize: '',
		...initSettings.options,
	};

	return initSettings;
}

function Service(props) {
	const settings = init(props.settings);

	const handleChange = (what) => (event) => {
		let value = event.target.value;

		if (what in settings.options) {
			settings.options[what] = value;
		} else {
			if (what === 'address') {
				const matches = value.match(/(rtmp.):\/\//);
				if (matches !== null) {
					if (props.skills.protocols.includes(matches[1])) {
						settings['protocol'] = matches[0];
						value = value.replace(matches[0], '');
					}
				}
			}

			settings[what] = value;
		}

		const output = createOutput(settings);

		props.onChange([output], settings);
	};

	const createOutput = (settings) => {
		const options = ['-f', 'flv'];

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
				options.push('-rtmp_enhanced_codecs', codecs.join(','));
			}
		}

		for (let key in settings.options) {
			if (settings.options[key].length !== 0) {
				if (key !== 'rtmp_flush_interval') {
					options.push('-' + key, settings.options[key]);
				} else if (settings.protocol.includes('rtmpt')) {
					options.push('-' + key, settings.options[key]);
				}
			}
		}

		const output = {
			address: settings.protocol + settings.address,
			options: options,
		};

		return output;
	};

	return (
		<Grid container spacing={2}>
			<Grid item xs={12} md={3}>
				<Select type="select" label={<Trans>Protocol</Trans>} value={settings.protocol} onChange={handleChange('protocol')}>
					<MenuItem value="rtmp://">rtmp://</MenuItem>
					<MenuItem value="rtmpe://">rtmpe://</MenuItem>
					<MenuItem value="rtmps://">rtmps://</MenuItem>
					<MenuItem value="rtmpt://">rtmpt://</MenuItem>
					<MenuItem value="rtmpte://">rtmpte://</MenuItem>
					<MenuItem value="rtmpts://">rtmpts://</MenuItem>
				</Select>
			</Grid>
			<Grid item xs={12} md={9}>
				<TextField variant="outlined" fullWidth type="url" label={<Trans>Address</Trans>} value={settings.address} onChange={handleChange('address')} />
			</Grid>
			<Grid item xs={12}>
				<TextField
					variant="outlined"
					fullWidth
					label={<Trans>Stream Key</Trans>}
					value={settings.options.rtmp_playpath}
					onChange={handleChange('rtmp_playpath')}
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
							<Grid item xs={12}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="rtmp_app"
									value={settings.options.rtmp_app}
									onChange={handleChange('rtmp_app')}
								/>
							</Grid>
							<Grid item xs={12}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="rtmp_playpath"
									value={settings.options.rtmp_playpath}
									onChange={handleChange('rtmp_playpath')}
								/>
							</Grid>
							<Grid item xs={12}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="rtmp_pageurl"
									value={settings.options.rtmp_pageurl}
									onChange={handleChange('rtmp_pageurl')}
								/>
							</Grid>
							<Grid item xs={12}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="rtmp_tcurl"
									value={settings.options.rtmp_tcurl}
									onChange={handleChange('rtmp_tcurl')}
								/>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="rtmp_conn"
									value={settings.options.rtmp_conn}
									onChange={handleChange('rtmp_conn')}
								/>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="rtmp_flashver"
									value={settings.options.rtmp_flashver}
									onChange={handleChange('rtmp_flashver')}
								/>
							</Grid>
							{settings.protocol.includes('rtmpt') && (
								<Grid item xs={12} md={6}>
									<TextField
										variant="outlined"
										fullWidth
										type="number"
										label="rtmp_flush_interval"
										value={settings.options.rtmp_flush_interval}
										onChange={handleChange('rtmp_flush_interval')}
									/>
								</Grid>
							)}
							<Grid item xs={12}>
								<Typography variant="h3">
									<Trans>Encryption</Trans>
								</Typography>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="rtmp_swfhash"
									value={settings.options.rtmp_swfhash}
									onChange={handleChange('rtmp_swfhash')}
								/>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									label="rtmp_swfsize"
									value={settings.options.rtmp_swfsize}
									onChange={handleChange('rtmp_swfsize')}
								/>
							</Grid>
							<Grid item xs={12}>
								<Divider />
							</Grid>
							<Grid item xs={12}>
								<Typography>
									<Trans>Documentation</Trans>{' '}
									<Link color="secondary" target="_blank" href="https://ffmpeg.org/ffmpeg-all.html#rtmp">
										https://ffmpeg.org/ffmpeg-all.html#rtmp
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
