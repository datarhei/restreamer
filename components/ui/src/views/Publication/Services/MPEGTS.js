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
import MultiSelect from '../../../misc/MultiSelect';
import MultiSelectOption from '../../../misc/MultiSelectOption';

const id = 'mpegts';
const name = 'MPEG-TS';
const version = '1.0';
const stream_key_link = '';
const description = (
	<Trans>
		Transmit the main source to an MPEG-TS Service. More details about the settings can be found here{' '}
		<Link color="secondary" target="_blank" href="https://ffmpeg.org/ffmpeg-all.html#mpegts-1">
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
	protocols: ['udp', 'tcp'],
	formats: ['mpegts'],
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
		protocol: 'udp://',
		address: '',
		options: {},
		...settings,
	};

	initSettings.options = {
		service_provider: 'FFmpeg',
		service_name: 'Service01',
		mpegts_transport_stream_id: '0x0001',
		mpegts_original_network_id: '0x0001',
		mpegts_service_id: '0x0001',
		mpegts_service_type: 'digital_tv',
		mpegts_service_type_hex_value: '',
		mpegts_pmt_start_pid: '0x1000',
		mpegts_start_pid: '0x0100',
		mpegts_m2ts_mode: false,
		muxrate: '',
		pes_payload_size: '2930',
		mpegts_flags: [],
		mpegts_copyts: false,
		omit_video_pes_length: true,
		pcr_period: '-1',
		pat_period: '0.1',
		sdt_period: '0.5',
		tables_version: '0',
		...initSettings.options,
	};

	return initSettings;
}

function Service(props) {
	const settings = init(props.settings);

	const handleChange = (what) => (event) => {
		const value = event.target.value;

		if (what in settings.options) {
			if (['mpegts_m2ts_mode', 'mpegts_copyts', 'omit_video_pes_length'].includes(what)) {
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
		const options = ['-f', 'mpegts'];

		for (let key in settings.options) {
			if (settings.options[key].length === 0) {
				continue;
			}
			if (['service_provider', 'service_name'].includes(key)) {
				options.push('-metadata', key + '="' + String(settings.options[key]) + '"');
			} else if (key === 'mpegts_service_type' && settings.options.mpegts_service_type === 'hex_value') {
				if (settings.options.mpegts_service_type_hex_value.length !== 0) {
					options.push('-mpegts_service_type', String(settings.options.mpegts_service_type_hex_value));
				}
			} else if (key === 'mpegts_service_type_hex_value') {
				continue;
			} else if (['mpegts_m2ts_mode', 'mpegts_copyts', 'omit_video_pes_length'].includes(key)) {
				if (settings.options[key]) {
					options.push('-' + key, '1');
				} else {
					options.push('-' + key, '-1');
				}
			} else if (key === 'mpegts_flags') {
				let flags = '';

				settings.options[key].forEach((flag) => {
					flags += '+' + flag;
				});

				options.push('-mpegts_flags', String(flags));
			} else {
				options.push('-' + key, String(settings.options[key]));
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
			<Grid item xs={12}>
				<Accordion className="accordion">
					<AccordionSummary className="accordion-summary" elevation={0} expandIcon={<ArrowDropDownIcon />}>
						<Typography>
							<Trans>Advanced settings</Trans>
						</Typography>
					</AccordionSummary>
					<AccordionDetails>
						<Grid container spacing={2}>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="service_provider"
									value={settings.options.service_provider}
									onChange={handleChange('service_provider')}
								/>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="service_name"
									value={settings.options.service_name}
									onChange={handleChange('service_name')}
								/>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="mpegts_transport_stream_id"
									value={settings.options.mpegts_transport_stream_id}
									onChange={handleChange('mpegts_transport_stream_id')}
								/>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="mpegts_original_network_id"
									value={settings.options.mpegts_original_network_id}
									onChange={handleChange('mpegts_original_network_id')}
								/>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="mpegts_service_id"
									value={settings.options.mpegts_service_id}
									onChange={handleChange('mpegts_service_id')}
								/>
							</Grid>
							<Grid item xs={12} md={6}>
								<Select
									type="select"
									label="mpegts_service_type"
									value={settings.options.mpegts_service_type}
									onChange={handleChange('mpegts_service_type')}
								>
									<MenuItem value="hex_value">hex_value</MenuItem>
									<MenuItem value="digital_tv">digital_tv</MenuItem>
									<MenuItem value="digital_radio">digital_radio</MenuItem>
									<MenuItem value="teletext">teletext</MenuItem>
									<MenuItem value="advanced_codec_digital_radio">advanced_codec_digital_radio</MenuItem>
									<MenuItem value="mpeg2_digital_hdtv">mpeg2_digital_hdtv</MenuItem>
									<MenuItem value="advanced_codec_digital_sdtv">advanced_codec_digital_sdtv</MenuItem>
									<MenuItem value="advanced_codec_digital_hdtv">advanced_codec_digital_hdtv</MenuItem>
								</Select>
							</Grid>
							{settings.options.mpegts_service_type === 'hex_value' && (
								<Grid item xs={12}>
									<TextField
										variant="outlined"
										fullWidth
										type="text"
										label="mpegts_service_type: hex_value"
										value={settings.options.mpegts_service_type_hex_value}
										onChange={handleChange('mpegts_service_type_hex_value')}
									/>
								</Grid>
							)}
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="mpegts_pmt_start_pid"
									value={settings.options.mpegts_pmt_start_pid}
									onChange={handleChange('mpegts_pmt_start_pid')}
								/>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="mpegts_start_pid"
									value={settings.options.mpegts_start_pid}
									onChange={handleChange('mpegts_start_pid')}
								/>
							</Grid>
							<Grid item xs={12}>
								<Checkbox label="mpegts_m2ts_mode" checked={settings.options.mpegts_m2ts_mode} onChange={handleChange('mpegts_m2ts_mode')} />
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="muxrate"
									value={settings.options.muxrate}
									onChange={handleChange('muxrate')}
								/>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="pes_payload_size"
									value={settings.options.pes_payload_size}
									onChange={handleChange('pes_payload_size')}
								/>
							</Grid>
							<Grid item xs={12}>
								<MultiSelect type="select" label="mpegts_flags" value={settings.options.mpegts_flags} onChange={handleChange('mpegts_flags')}>
									<MultiSelectOption value="resend_headers" name="resend_headers" />
									<MultiSelectOption value="latm" name="latm" />
									<MultiSelectOption value="pat_pmt_at_frames" name="pat_pmt_at_frames" />
									<MultiSelectOption value="system_b" name="system_b" />
									<MultiSelectOption value="initial_discontinuity" name="initial_discontinuity" />
								</MultiSelect>
							</Grid>
							<Grid item xs={12}>
								<Checkbox label="mpegts_copyts" checked={settings.options.mpegts_copyts} onChange={handleChange('mpegts_copyts')} />
								<Checkbox
									label="omit_video_pes_length"
									checked={settings.options.omit_video_pes_length}
									onChange={handleChange('omit_video_pes_length')}
								/>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="pcr_period"
									value={settings.pcr_period}
									onChange={handleChange('pcr_period')}
								/>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="pat_period"
									value={settings.pat_period}
									onChange={handleChange('pat_period')}
								/>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="sdt_period"
									value={settings.sdt_period}
									onChange={handleChange('sdt_period')}
								/>
							</Grid>
							<Grid item xs={12} md={6}>
								<TextField
									variant="outlined"
									fullWidth
									type="text"
									label="tables_version"
									value={settings.tables_version}
									onChange={handleChange('tables_version')}
								/>
							</Grid>
							<Grid item xs={12}>
								<Divider />
							</Grid>
							<Grid item xs={12}>
								<Typography>
									<Trans>Documentation</Trans>{' '}
									<Link color="secondary" target="_blank" href="https://ffmpeg.org/ffmpeg-all.html#mpegts-1">
										https://ffmpeg.org/ffmpeg-all.html#mpegts-1
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
