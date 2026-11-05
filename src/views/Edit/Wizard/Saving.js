import React from 'react';

import { Trans } from '@lingui/macro';
import CircularProgress from '@mui/material/CircularProgress';
import Grid from '@mui/material/Grid';
import Typography from '@mui/material/Typography';

import Paper from '../../../misc/Paper';
import PaperHeader from '../../../misc/PaperHeader';

export default function Saving(props) {
	return (
		<Paper xs={12} md={5} marginBottom="6em" className="PaperM">
			<PaperHeader spacing={2} variant="h1" onAbort={props.onAbort} />
			<Grid container justifyContent="center" spacing={2} align="center">
				<Grid item xs={12}>
					<CircularProgress color="inherit" />
				</Grid>
				<Grid item xs={12}>
					<Typography>
						<Trans>Please wait. Setting up the stream ...</Trans>
					</Typography>
				</Grid>
			</Grid>
		</Paper>
	);
}

Saving.defaultProps = {
	onAbort: () => {},
};
