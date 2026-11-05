import React from 'react';

import { Trans } from '@lingui/macro';
import Grid from '@mui/material/Grid';
import TextField from '@mui/material/TextField';
import Typography from '@mui/material/Typography';

function init(settings) {
	const initSettings = {
		cpu_usage: 0,
		memory_mbytes: 0,
		waitfor_seconds: 5,
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

		settings[what] = value;

		props.onChange(settings, false);
	};

	return (
		<Grid container spacing={2}>
			<Grid item xs={12} md={4}>
				<TextField
					variant="outlined"
					fullWidth
					type="number"
					inputProps={{ min: 0, max: 100 }}
					label={<Trans>CPU Limit (percent)</Trans>}
					value={settings.cpu_usage}
					onChange={handleChange('cpu_usage')}
				/>
				<Typography variant="caption">
					<Trans>CPU usage limit in percent (0-100%), 0 for unlimited.</Trans>
				</Typography>
			</Grid>
			<Grid item xs={12} md={4}>
				<TextField
					variant="outlined"
					fullWidth
					type="number"
					inputProps={{ min: 0 }}
					label={<Trans>Memory Limit (megabytes)</Trans>}
					value={settings.memory_mbytes}
					onChange={handleChange('memory_mbytes')}
				/>
				<Typography variant="caption">
					<Trans>Memory usage limit in megabytes, 0 for unlimited.</Trans>
				</Typography>
			</Grid>
			<Grid item xs={12} md={4}>
				<TextField
					variant="outlined"
					fullWidth
					type="number"
					inputProps={{ min: 0 }}
					label={<Trans>Threshold (seconds)</Trans>}
					value={settings.waitfor_seconds}
					onChange={handleChange('waitfor_seconds')}
				/>
				<Typography variant="caption">
					<Trans>Number of seconds the limits are allowed to be exceeded.</Trans>
				</Typography>
			</Grid>
		</Grid>
	);
}

Control.defaulProps = {
	settings: {},
	onChange: function (settings, automatic) {},
};
