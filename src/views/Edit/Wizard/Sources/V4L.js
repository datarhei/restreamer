import React from 'react';

import { useLingui } from '@lingui/react';
import { faUsb } from '@fortawesome/free-brands-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Trans, t } from '@lingui/macro';
import Button from '@mui/material/Button';
import Grid from '@mui/material/Grid';
import MenuItem from '@mui/material/MenuItem';
import RefreshIcon from '@mui/icons-material/Refresh';
import Typography from '@mui/material/Typography';

import * as S from '../../Sources/V4L';
import Select from '../../../../misc/Select';

function initSettings(initialSettings, knownDevices) {
	const settings = {
		...S.func.initSettings(initialSettings),
		format: 'nv12',
		framerate: '25',
		size: '1280x720',
	};

	const devices = initDevices(knownDevices);

	if (devices.length !== 0) {
		let found = false;
		for (let device of devices) {
			if (settings.device === device.id) {
				found = true;
				break;
			}
		}

		if (found === false) {
			settings.device = '';
		}
	}

	if (devices.length !== 0) {
		const device = devices[0];

		if (settings.device.length === 0) {
			settings.device = device.id;
		}
	}

	return settings;
}

function initDevices(initialDevices) {
	const devices = initialDevices.filter((device) => device.media === 'video' && device.extra.match(/bcm2835[-_]v4l2/) === null);

	return devices;
}

function Source(props) {
	const { i18n } = useLingui();
	const settings = initSettings(props.settings, props.knownDevices);
	const devices = initDevices(props.knownDevices);

	const handleChange = (newSettings) => {
		newSettings = newSettings || settings;

		props.onChange(S.id, newSettings, S.func.createInputs(newSettings), devices.length !== 0 ? true : false);
	};

	const handleRefresh = () => {
		props.onRefresh();
	};

	const update = (what) => (event) => {
		const value = event.target.value;
		const newSettings = settings;

		if (what in newSettings) {
			newSettings[what] = value;
		}

		handleChange(newSettings);
	};

	React.useEffect(() => {
		handleChange();
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

	const options = devices.map((device) => {
		return (
			<MenuItem key={device.id} value={device.id}>
				{device.name} ({device.id})
			</MenuItem>
		);
	});

	if (options.length === 0) {
		options.push(
			<MenuItem key="none" value="none" disabled={true}>
				{i18n._(t`No input device available`)}
			</MenuItem>,
		);
	}

	const videoDevices = (
		<Select label={<Trans>Video device</Trans>} value={settings.device} onChange={update('device')}>
			{options}
		</Select>
	);

	return (
		<React.Fragment>
			<Grid item xs={12}>
				<Typography>
					<Trans>Select a device:</Trans>
				</Typography>
			</Grid>
			<Grid item xs={12}>
				{videoDevices}
				<Button size="small" startIcon={<RefreshIcon />} onClick={handleRefresh} sx={{ float: 'right' }}>
					<Trans>Refresh</Trans>
				</Button>
			</Grid>
		</React.Fragment>
	);
}

Source.defaultProps = {
	knownDevices: [],
	settings: {},
	onChange: function (type, settings, inputs, ready) {},
	onRefresh: function () {},
};

function SourceIcon(props) {
	return <FontAwesomeIcon icon={faUsb} style={{ color: '#FFF' }} {...props} />;
}

const id = 'video4linux2';
const type = 'video4linux2';
const name = <Trans>Hardware device</Trans>;
const capabilities = ['video'];

export { id, type, name, capabilities, SourceIcon as icon, Source as component };
