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
import TextField from '@mui/material/TextField';
import Typography from '@mui/material/Typography';

import Checkbox from '../../../misc/Checkbox';

const id = 'udp';
const name = 'UDP';
const version = '1.0';
const stream_key_link = '';
const description = (
	<Trans>
		Transmit the main source to an UDP Server. More details about the settings can be found{' '}
		<Link color="secondary" target="_blank" href="https://ffmpeg.org/ffmpeg-all.html#udp">
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
	protocols: ['udp'],
	formats: ['mpegts'],
	codecs: {
		audio: ['aac', 'mp3'],
		video: ['h264', 'hevc', 'vp9', 'av1'],
	},
};

function ServiceIcon(props) {
	return <FontAwesomeIcon icon={faTools} style={{ color: '#39B54A' }} {...props} />;
}

// init merges the given settings with the default settings
function init(settings) {
	const initSettings = {
		protocol: 'udp://',
		address: '',
		params: {},
		...settings,
	};

	initSettings.params = {
		buffer_size: '32',
		bitrate: '',
		burst_bits: '',
		localport: '',
		localaddr: '',
		pkt_size: '',
		ttl: '',
		fifo_size: '',
		overrun_nonfatal: false,
		broadcast: false,
		...initSettings.params,
	};

	return initSettings;
}

function Service(props) {
	const settings = init(props.settings);

	const handleChange = (what) => (event) => {
		const value = event.target.value;

		if (what in settings.params) {
			if (['overrun_nonfatal', 'broadcast'].includes(what)) {
				settings.params[what] = !settings.params[what];
			} else {
				settings.params[what] = value;
			}
		} else {
			settings[what] = value;
		}

		const output = createOutput(settings);

		props.onChange([output], settings);
	};

	const createOutput = (settings) => {
		const params = [];

		for (let key in settings.params) {
			if (settings.params[key].length === 0) {
				continue;
			}
			params.push(key + '=' + settings.params[key]);
		}

		const output = {
			address: settings.protocol + settings.address + '?' + params.join('&'),
			options: ['-f', 'mpegts'],
		};

		return output;
	};

	return (
		<Grid container spacing={2}>
			<Grid item xs={12} md={3}>
				<TextField variant="outlined" fullWidth label={<Trans>Protocol</Trans>} value={settings.protocol} readOnly disabled />
			</Grid>
			<Grid item xs={12} md={9}>
				<TextField variant="outlined" type="url" fullWidth label={<Trans>Address</Trans>} value={settings.address} onChange={handleChange('address')} />
			</Grid>
			<Grid item xs={12}>
				<Accordion className="accordion">
					<AccordionSummary expandIcon={<ArrowDropDownIcon />} className="accordion-summary">
						<Typography>
							<Trans>Advanced settings</Trans>
						</Typography>
					</AccordionSummary>
					<AccordionDetails>
						<Grid container spacing={2}>
							<Grid item xs={12} md={4}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="0"
									step="1"
									label="buffer_size"
									value={settings.params.buffer_size}
									onChange={handleChange('buffer_size')}
								/>
							</Grid>
							<Grid item xs={12} md={4}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="0"
									step="1"
									label="bitrate"
									value={settings.params.bitrate}
									onChange={handleChange('bitrate')}
								/>
							</Grid>
							<Grid item xs={12} md={4}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="0"
									step="1"
									label="burst_bits"
									value={settings.params.burst_bits}
									onChange={handleChange('burst_bits')}
								/>
							</Grid>
							<Grid item xs={12} md={8}>
								<TextField
									variant="outlined"
									fullWidth
									type="url"
									label="localaddr"
									value={settings.params.localaddr}
									onChange={handleChange('localaddr')}
								/>
							</Grid>
							<Grid item xs={12} md={4}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="0"
									step="1"
									label="localport"
									value={settings.params.localport}
									onChange={handleChange('localport')}
								/>
							</Grid>
							<Grid item xs={12} md={4}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="0"
									step="1"
									label="pkt_size"
									value={settings.params.pkt_size}
									onChange={handleChange('pkt_size')}
								/>
							</Grid>
							<Grid item xs={12} md={4}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="0"
									step="1"
									label="ttl"
									value={settings.params.ttl}
									onChange={handleChange('ttl')}
								/>
							</Grid>
							<Grid item xs={12} md={4}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="0"
									step="1"
									label="fifo_size"
									value={settings.params.fifo_size}
									onChange={handleChange('fifo_size')}
								/>
							</Grid>
							<Grid item xs={12}>
								<Checkbox label="overrun_nonfatal" checked={settings.params.overrun_nonfatal} onChange={handleChange('overrun_nonfatal')} />
								<Checkbox label="broadcast" checked={settings.params.broadcast} onChange={handleChange('broadcast')} />
							</Grid>
							<Grid item xs={12}>
								<Divider />
							</Grid>
							<Grid item xs={12}>
								<Typography>
									<Trans>Documentation</Trans>{' '}
									<Link color="secondary" target="_blank" href="https://ffmpeg.org/ffmpeg-all.html#udp">
										https://ffmpeg.org/ffmpeg-all.html#udp
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
