import React from 'react';

import { Trans } from '@lingui/macro';
import Grid from '@mui/material/Grid';
import TextField from '@mui/material/TextField';
import Typography from '@mui/material/Typography';

import Checkbox from '../Checkbox';

function init(settings) {
	const initSettings = {
		autostart: false,
		reconnect: true,
		delay: 30,
		staleTimeout: 30,
		low_delay: false,
		...settings,
	};

	return initSettings;
}

export default function Control(props) {
	const settings = init(props.settings);

	// Set the defaults
	React.useEffect(() => {
		props.onChange(settings, true);
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

	const handleChange = (what) => (event) => {
		const value = event.target.value;

		if (['autostart', 'reconnect', 'low_delay'].includes(what)) {
			settings[what] = !settings[what];
		} else {
			settings[what] = value;
		}

		props.onChange(settings, false);
	};

	return (
		<Grid container spacing={2}>
			<Grid item xs={12}>
				<Checkbox label={<Trans>Reconnect</Trans>} checked={settings.reconnect} onChange={handleChange('reconnect')} />
				<Checkbox label={<Trans>Low latency (Buffer)</Trans>} checked={settings.low_delay} onChange={handleChange('low_delay')} />
			</Grid>
			<Grid item xs={12} md={6}>
				<TextField
					variant="outlined"
					fullWidth
					type="number"
					label={<Trans>Reconnect delay (seconds)</Trans>}
					disabled={!settings.reconnect}
					value={settings.delay}
					onChange={handleChange('delay')}
				/>
				<Typography variant="caption">
					<Trans>Seconds until a process is restarted.</Trans>
				</Typography>
			</Grid>
			<Grid item xs={12} md={6}>
				<TextField
					variant="outlined"
					fullWidth
					type="number"
					label={<Trans>Stale timeout (seconds)</Trans>}
					value={settings.staleTimeout}
					onChange={handleChange('staleTimeout')}
				/>
				<Typography variant="caption">
					<Trans>Seconds until a staled process is terminated.</Trans>
				</Typography>
			</Grid>
		</Grid>
	);
}

Control.defaulProps = {
	settings: {},
	onChange: function (settings, automatic) {},
};
