import React from 'react';
import { useNavigate } from 'react-router-dom';

import { Trans, t } from '@lingui/macro';
import { useLingui } from '@lingui/react';
import Button from '@mui/material/Button';
import Grid from '@mui/material/Grid';
import Icon from '@mui/icons-material/KeyboardTab';
import MenuItem from '@mui/material/MenuItem';
import RefreshIcon from '@mui/icons-material/Refresh';
import Typography from '@mui/material/Typography';

import * as S from '../../Sources/Network';
import BoxTextarea from '../../../../misc/BoxTextarea';
import Select from '../../../../misc/Select';
import Textarea from '../../../../misc/Textarea';

const initSettings = (initialSettings, config) => {
	const settings = {
		...S.func.initSettings(initialSettings, config),
		mode: 'push',
	};

	settings.push.type = 'rtmp';

	return settings;
};

function Source(props) {
	const { i18n } = useLingui();
	const navigate = useNavigate();
	const config = S.func.initConfig(props.config);
	const settings = initSettings(props.settings, config);
	const skills = S.func.initSkills(props.skills);

	const handleChange = (newSettings) => {
		newSettings = newSettings || settings;

		const inputs = S.func.createInputs(newSettings, config, skills);
		newSettings.address = inputs[0].address;

		props.onChange(S.id, newSettings, inputs, config.rtmp.enabled);
	};

	const handleRefresh = () => {
		props.onRefresh();
	};

	const update = (what) => (event) => {
		const value = event.target.value;
		const newSettings = {
			...settings,
		};

		if (what in newSettings.push) {
			newSettings.push[what] = value;
		}

		handleChange(newSettings);
	};

	React.useEffect(() => {
		handleChange();
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

	let form = null;

	if (config.rtmp.enabled === false) {
		form = (
			<React.Fragment>
				<Grid item xs={12}>
					<Typography>
						<Trans>RTMP server is not enabled</Trans>
					</Typography>
				</Grid>
				<Grid item xs={12}>
					<Button variant="outlined" size="large" fullWidth color="primary" onClick={() => navigate('/settings/rtmp')}>
						<Trans>Enable RTMP server ...</Trans>
					</Button>
				</Grid>
			</React.Fragment>
		);
	} else {
		const filteredDevices = props.knownDevices.filter((device) => device.media === 'rtmp');
		const options = filteredDevices.map((device) => {
			return (
				<MenuItem key={device.id} value={device.id}>
					{device.name}
				</MenuItem>
			);
		});

		options.unshift(
			<MenuItem key="none" value="none" disabled>
				{i18n._(t`Choose an input stream ...`)}
			</MenuItem>,
		);

		options.push(
			<MenuItem key={config.channelid} value={config.channelid}>
				{i18n._(t`Send stream to address ...`)}
			</MenuItem>,
		);

		const RTMP = S.func.getRTMP(config, settings.push.name);

		form = (
			<React.Fragment>
				<Grid item xs={12}>
					<Select type="select" label={<Trans>Input stream</Trans>} value={settings.push.name} onChange={update('name')}>
						{options}
					</Select>
					<Button size="small" startIcon={<RefreshIcon />} onClick={handleRefresh} sx={{ float: 'right' }}>
						<Trans>Refresh</Trans>
					</Button>
				</Grid>
				{settings.push.name === config.channelid && (
					<React.Fragment>
						<Grid item xs={12}>
							<Typography>
								<Trans>Address:</Trans>
							</Typography>
						</Grid>
						<Grid item xs={12}>
							<BoxTextarea>
								<Textarea rows={1} value={RTMP} readOnly allowCopy />
							</BoxTextarea>
						</Grid>
					</React.Fragment>
				)}
			</React.Fragment>
		);
	}

	return form;
}

Source.defaultProps = {
	knownDevices: [],
	settings: {},
	config: null,
	skills: null,
	onChange: function (type, settings, inputs, ready) {},
	onRefresh: function () {},
};

function SourceIcon(props) {
	return <Icon style={{ color: '#FFF' }} {...props} />;
}

const id = 'rtmp';
const type = 'network';
const name = <Trans>RTMP server</Trans>;
const capabilities = ['audio', 'video'];

export { id, type, name, capabilities, SourceIcon as icon, Source as component };
