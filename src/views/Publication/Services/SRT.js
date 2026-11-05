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

import Checkbox from '../../../misc/Checkbox';
import Select from '../../../misc/Select';

const id = 'srt';
const name = 'SRT';
const version = '1.0';
const stream_key_link = '';
const description = (
	<Trans>
		Transmit the main source to an SRT Server. More details about the settings can be found{' '}
		<Link color="secondary" target="_blank" href="https://ffmpeg.org/ffmpeg-all.html#srt">
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
	protocols: ['srt'],
	formats: ['mpegts'],
	codecs: {
		audio: ['aac', 'mp3'],
		video: ['h264', 'hevc', 'vp9', 'av1'],
	},
};

function ServiceIcon(props) {
	return <FontAwesomeIcon icon={faTools} style={{ color: '#39B54A' }} {...props} />;
}

const defaults = {
	connect_timeout: '3000',
	ffs: '25600',
	inputbw: '0',
	iptos: '0xB8',
	ipttl: '64',
	latency: '120000',
	listen_timeout: '0',
	maxbw: '0',
	mode: 'caller',
	mss: '1500',
	nakreport: true,
	oheadbw: '25',
	payload_size: '-1',
	send_buffer_size: '0',
	timeout: '3000000',
	tlpktdrop: true,
	sndbuf: '0',
	rcvbuf: '0',
	lossmaxttl: '0',
	smoother: 'live',
	streamid: '',
	transtype: '',
	passphrase: '',
	enforced_encryption: false,
	kmrefreshrate: '-1',
	kmpreannounce: '-1',
	pbkeylen: '16',
};

function init(settings) {
	const initSettings = {
		protocol: 'srt://',
		address: '',
		params: {},
		...settings,
	};

	initSettings.params = {
		...defaults,
		transtype: 'live',
		...initSettings.params,
	};

	return initSettings;
}

function Service(props) {
	const settings = init(props.settings);

	const handleChange = (what) => (event) => {
		const value = event.target.value;

		if (what in settings.params) {
			if (['enforced_encryption', 'tlpktdrop', 'nakreport'].includes(what)) {
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

			if (settings.params[key] === defaults[key]) {
				continue;
			}

			if (typeof settings.params[key] === 'boolean') {
				params.push(key + '=' + (settings.params[key] ? '1' : '0'));
			} else {
				params.push(key + '=' + settings.params[key]);
			}
		}

		const options = ['-bsf:v', 'dump_extra', '-f', 'mpegts'];

		const output = {
			address: settings.protocol + settings.address + '?' + params.join('&'),
			options: options,
		};

		return output;
	};

	return (
		<Grid container spacing={2}>
			<Grid item xs={12} md={3}>
				<TextField variant="outlined" fullWidth label={<Trans>Protocol</Trans>} value={settings.protocol} readOnly disabled />
			</Grid>
			<Grid item xs={12} md={9}>
				<TextField variant="outlined" fullWidth label={<Trans>Address</Trans>} value={settings.address} onChange={handleChange('address')} />
			</Grid>
			<Grid item xs={12}>
				<TextField variant="outlined" fullWidth label="streamid" value={settings.params.streamid} onChange={handleChange('streamid')} />
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
							<Grid item xs={12}>
								<Typography variant="h3">
									<Trans>General</Trans>
								</Typography>
							</Grid>
							<Grid item xs={12} md={4}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="0"
									step="1"
									label="connect_timeout"
									value={settings.params.connect_timeout}
									onChange={handleChange('connect_timeout')}
								/>
							</Grid>
							<Grid item xs={12} md={4}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="0"
									step="1"
									label="ffs"
									value={settings.params.ffs}
									onChange={handleChange('ffs')}
								/>
							</Grid>
							<Grid item xs={12} md={4}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="0"
									step="1"
									label="inputbw"
									value={settings.params.inputbw}
									onChange={handleChange('inputbw')}
								/>
							</Grid>
							<Grid item xs={12} md={4}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="iptos"
									value={settings.params.iptos}
									onChange={handleChange('iptos')}
								/>
							</Grid>
							<Grid item xs={12} md={4}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="0"
									step="1"
									label="ipttl"
									value={settings.params.ipttl}
									onChange={handleChange('ipttl')}
								/>
							</Grid>
							<Grid item xs={12} md={4}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="0"
									step="1000"
									label="latency"
									value={settings.params.latency}
									onChange={handleChange('latency')}
								/>
							</Grid>
							<Grid item xs={12} md={4}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="-1"
									step="1"
									label="maxbw"
									value={settings.params.maxbw}
									onChange={handleChange('maxbw')}
								/>
							</Grid>
							<Grid item xs={12} md={4}>
								<Select label="mode" value={settings.params.mode} onChange={handleChange('mode')}>
									<MenuItem value="caller">caller</MenuItem>
									<MenuItem value="listener">listener</MenuItem>
									<MenuItem value="rendezvous">rendezvous</MenuItem>
								</Select>
							</Grid>
							<Grid item xs={12} md={4}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="0"
									step="1"
									label="mss"
									value={settings.params.mss}
									onChange={handleChange('mss')}
								/>
							</Grid>
							<Grid item xs={12} md={4}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="0"
									max="100"
									step="1"
									label="oheadbw"
									value={settings.params.oheadbw}
									onChange={handleChange('oheadbw')}
								/>
							</Grid>
							<Grid item xs={12} md={4}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="-1"
									step="1"
									label="payload_size"
									value={settings.params.payload_size}
									onChange={handleChange('payload_size')}
								/>
							</Grid>
							<Grid item xs={12} md={4}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="0"
									step="1"
									label="send_buffer_size"
									value={settings.params.send_buffer_size}
									onChange={handleChange('send_buffer_size')}
								/>
							</Grid>
							<Grid item xs={12} md={4}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="0"
									step="1"
									label="sndbuf"
									value={settings.params.sndbuf}
									onChange={handleChange('sndbuf')}
								/>
							</Grid>
							<Grid item xs={12} md={4}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="0"
									step="1"
									label="rcvbuf"
									value={settings.params.rcvbuf}
									onChange={handleChange('rcvbuf')}
								/>
							</Grid>
							<Grid item xs={12} md={4}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="0"
									step="1"
									label="lossmaxttl"
									value={settings.params.lossmaxttl}
									onChange={handleChange('lossmaxttl')}
								/>
							</Grid>
							<Grid item xs={12} md={4}>
								<Select label="smoother" value={settings.params.smoother} onChange={handleChange('smoother')}>
									<MenuItem value="live">live</MenuItem>
									<MenuItem value="file">file</MenuItem>
								</Select>
							</Grid>
							<Grid item xs={12} md={4}>
								<Select label="transtype" value={settings.params.transtype} onChange={handleChange('transtype')}>
									<MenuItem value="live">live</MenuItem>
									<MenuItem value="file">file</MenuItem>
								</Select>
							</Grid>
							<Grid item xs={12}>
								<Checkbox label="nakreport" checked={settings.params.nakreport} onChange={handleChange('nakreport')} />
								<Checkbox label="tlpktdrop" checked={settings.params.tlpktdrop} onChange={handleChange('tlpktdrop')} />
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
									label="passphrase"
									value={settings.params.passphrase}
									onChange={handleChange('passphrase')}
								/>
							</Grid>
							<Grid item xs={12} md={4}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="-1"
									step="1"
									label="kmrefreshrate"
									value={settings.params.kmrefreshrate}
									onChange={handleChange('kmrefreshrate')}
								/>
							</Grid>
							<Grid item xs={12} md={4}>
								<TextField
									variant="outlined"
									fullWidth
									type="number"
									min="-1"
									step="1"
									label="kmpreannounce"
									value={settings.params.kmpreannounce}
									onChange={handleChange('kmpreannounce')}
								/>
							</Grid>
							<Grid item xs={12} md={4}>
								<Select label="pbkeylen" value={settings.params.pbkeylen} onChange={handleChange('pbkeylen')}>
									<MenuItem value="16">16</MenuItem>
									<MenuItem value="24">24</MenuItem>
									<MenuItem value="32">32</MenuItem>
								</Select>
							</Grid>
							<Grid item xs={12}>
								<Checkbox
									label="enforced_encryption"
									checked={settings.params.enforced_encryption}
									onChange={handleChange('enforced_encryption')}
								/>
							</Grid>
							<Grid item xs={12}>
								<Divider />
							</Grid>
							<Grid item xs={12}>
								<Typography>
									<Trans>Documentation</Trans>{' '}
									<Link color="secondary" target="_blank" href="https://ffmpeg.org/ffmpeg-all.html#srt">
										https://ffmpeg.org/ffmpeg-all.html#srt
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
