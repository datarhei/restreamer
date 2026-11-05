import React from 'react';

import { Trans } from '@lingui/macro';
import Grid from '@mui/material/Grid';
import TextField from '@mui/material/TextField';
import Typography from '@mui/material/Typography';

import Checkbox from '../Checkbox';

function init(settings) {
	const initSettings = {
		enable: true,
		interval: 60,
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

		if (['enable'].includes(what)) {
			settings[what] = !settings[what];
		} else {
			settings[what] = value;
		}

		props.onChange(settings, false);
	};

	return (
		<Grid container spacing={2}>
			<Grid item xs={12}>
				<Checkbox label={<Trans>Enable snapshots</Trans>} checked={settings.enable} onChange={handleChange('enable')} />
			</Grid>
			<Grid item xs={12} md={6}>
				<TextField
					variant="outlined"
					fullWidth
					label={<Trans>Interval (seconds)</Trans>}
					disabled={!settings.enable}
					value={settings.interval}
					onChange={handleChange('interval')}
				/>
				<Typography variant="caption">
					<Trans>Seconds until the snapshot/thumbnail of the video source is updated.</Trans>
				</Typography>
			</Grid>
		</Grid>
	);
}

Control.defaulProps = {
	settings: {},
	onChange: function (settings, automatic) {},
};
