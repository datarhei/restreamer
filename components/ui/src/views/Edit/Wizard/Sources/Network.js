import React from 'react';

import { Trans } from '@lingui/macro';
import Grid from '@mui/material/Grid';
import Icon from '@mui/icons-material/SettingsEthernet';
import TextField from '@mui/material/TextField';
import Typography from '@mui/material/Typography';

import * as S from '../../Sources/Network';
import Checkbox from '../../../../misc/Checkbox';
import Password from '../../../../misc/Password';

const initSettings = (initialSettings, config) => {
	const settings = {
		...S.func.initSettings(initialSettings, config),
		mode: 'pull',
	};

	return settings;
};

function Source(props) {
	const config = S.func.initConfig(props.config);
	const settings = initSettings(props.settings, config);
	const skills = S.func.initSkills(props.skills);

	const handleChange = (newSettings) => {
		newSettings = newSettings || settings;

		props.onChange(S.id, newSettings, S.func.createInputs(newSettings, config, skills), S.func.isValidURL(newSettings.address));
	};

	const update = (protocol, what) => (event) => {
		const value = event.target.value;
		const newSettings = settings;

		if (protocol === 'rtsp') {
			if (['udp'].includes(what)) {
				newSettings.rtsp[what] = !newSettings.rtsp[what];
			}
		} else {
			newSettings[what] = value;
		}

		handleChange(newSettings);
	};

	React.useEffect(() => {
		handleChange();
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

	const protocol = S.func.getProtocolClass(settings.address);

	return (
		<React.Fragment>
			<Grid item xs={12}>
				<Typography>
					<Trans>Enter the address of your network source:</Trans>
				</Typography>
			</Grid>
			<Grid item xs={12}>
				<TextField
					variant="outlined"
					fullWidth
					label={<Trans>Address</Trans>}
					placeholder="rtsp://ip:port/path"
					value={settings.address}
					onChange={update('', 'address')}
				/>
				<Typography variant="caption">
					<Trans>Supports HTTP (HLS, DASH), RTP, RTSP, RTMP, SRT and more.</Trans>
				</Typography>
			</Grid>
			{protocol === 'rtsp' && (
				<Grid item xs={12}>
					<Checkbox label={<Trans>UDP transport</Trans>} checked={settings.rtsp.udp} onChange={update('rtsp', 'udp')} />
				</Grid>
			)}
			<Grid item md={6} xs={12}>
				<TextField
					variant="outlined"
					fullWidth
					label={<Trans>Username</Trans>}
					value={settings.username}
					onChange={update('', 'username')}
					disabled={protocol === 'srt'}
				/>
				<Typography variant="caption">
					<Trans>Username for the device.</Trans>
				</Typography>
			</Grid>
			<Grid item md={6} xs={12}>
				<Password
					variant="outlined"
					fullWidth
					label={<Trans>Password</Trans>}
					value={settings.password}
					onChange={update('', 'password')}
					disabled={protocol === 'srt'}
				/>
				<Typography variant="caption">
					<Trans>Password for the device.</Trans>
				</Typography>
			</Grid>
		</React.Fragment>
	);
}

Source.defaultProps = {
	settings: {},
	config: null,
	skills: null,
	onChange: function (type, settings, inputs, ready) {},
};

function SourceIcon(props) {
	return <Icon style={{ color: '#FFF' }} {...props} />;
}

const id = 'network';
const type = 'network';
const name = <Trans>Network source</Trans>;
const capabilities = ['audio', 'video'];

export { id, type, name, capabilities, SourceIcon as icon, Source as component };
