import React from 'react';

import { Trans } from '@lingui/react/macro';
import CircularProgress from '@mui/material/CircularProgress';
import Grid from '@mui/material/Grid';
import Typography from '@mui/material/Typography';

import Paper from '../../../misc/Paper';
import PaperHeader from '../../../misc/PaperHeader';
import { withDefaults } from '../../../utils/props';

export default function Probe(inputProps) {
	const props = withDefaults(inputProps, ProbeDefaultProps);
	return (
		<Paper xs={12} md={5} marginBottom="6em" className="PaperM">
			<PaperHeader spacing={2} variant="h1" onAbort={props.onAbort} />
			<Grid container justifyContent="center" spacing={2} align="center">
				<Grid item xs={12}>
					<CircularProgress color="inherit" />
				</Grid>
				<Grid item xs={12}>
					<Typography textAlign="center">
						<Trans>Please wait. Probe stream data ...</Trans>
					</Typography>
				</Grid>
			</Grid>
		</Paper>
	);
}

const ProbeDefaultProps = {
	onAbort: () => {},
};
