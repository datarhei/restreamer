import React from 'react';

import { faImages as icon } from '@fortawesome/free-solid-svg-icons';
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

const id = 'image2';
const name = 'Image2';
const version = '1.0';
const stream_key_link = '';
const description = (
	<Trans>
		Transmit snapshots of the main source to an HTTP/S Server. More details about the settings can be found{' '}
		<Link color="secondary" target="_blank" href="https://ffmpeg.org/ffmpeg-all.html#image2">
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
	formats: ['image2'],
	codecs: {
		audio: ['none'],
		video: ['h264'],
	},
};

function ServiceIcon(props) {
	return <FontAwesomeIcon icon={icon} style={{ color: '#39B54A' }} {...props} />;
}

function init(settings) {
	const initSettings = {
		protocol: 'https://',
		address: '',
		options: {},
		...settings,
	};

	initSettings.options = {
		frame_pts: false,
		start_number: '1',
		update: true,
		strftime: false,
		method: 'PUT',
		...initSettings.options,
	};

	return initSettings;
}

function Service(props) {
	const settings = init(props.settings);

	const handleChange = (what) => (event) => {
		const value = event.target.value;

		if (what in settings.options) {
			if (['frame_pts', 'update', 'strftime'].includes(what)) {
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
		const options = ['-f', 'image2'];

		for (let key in settings.options) {
			if (settings.options[key].length === 0) {
				continue;
			}
			options.push('-' + key, String(settings.options[key]));
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
					<MenuItem value="http://">http://</MenuItem>
					<MenuItem value="https://">https://</MenuItem>
				</Select>
			</Grid>
			<Grid item xs={12} md={9}>
				<TextField variant="outlined" fullWidth label={<Trans>Address</Trans>} value={settings.address} onChange={handleChange('address')} />
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
									type="number"
									min="1"
									step="1"
									label="start_number"
									value={settings.options.start_number}
									onChange={handleChange('start_number')}
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
								<Checkbox label="frame_pts" checked={settings.options.frame_pts} onChange={handleChange('frame_pts')} />
								<Checkbox label="update" checked={settings.options.update} onChange={handleChange('update')} />
								<Checkbox label="strftime" checked={settings.options.strftime} onChange={handleChange('strftime')} />
							</Grid>
							<Grid item xs={12}>
								<Divider />
							</Grid>
							<Grid item xs={12}>
								<Typography>
									<Trans>Documentation</Trans>{' '}
									<Link color="secondary" target="_blank" href="https://ffmpeg.org/ffmpeg-all.html#image2">
										https://ffmpeg.org/ffmpeg-all.html#image2
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
