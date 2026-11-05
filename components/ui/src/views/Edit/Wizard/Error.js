import React from 'react';

import { Trans } from '@lingui/macro';
import Grid from '@mui/material/Grid';
import Typography from '@mui/material/Typography';
import Button from '@mui/material/Button';
import WarningIcon from '@mui/icons-material/Warning';

import BoxText from '../../../misc/BoxText';
import Paper from '../../../misc/Paper';
import PaperHeader from '../../../misc/PaperHeader';

export default function Error(props) {
	return (
		<Paper xs={12} sm={8} md={6} marginBottom="6em" className="PaperM">
			<PaperHeader spacing={3} variant="h1" title={<Trans>Error</Trans>} onAbort={props.onAbort} onHelp={props.onHelp} />
			<Grid container spacing={3}>
				<BoxText color="dark">
					<WarningIcon fontSize="large" color="error" />
					<Typography textAlign="center">
						<Trans>There was an error setting up the stream.</Trans>
					</Typography>
				</BoxText>
				<Grid item xs={12}>
					<Button variant="outlined" fullWidth color="primary" onClick={props.onNext}>
						<Trans>Retry</Trans>
					</Button>
				</Grid>
			</Grid>
		</Paper>
	);
}

Error.defaultProps = {
	onAbort: () => {},
	onHelp: () => {},
	onNext: () => {},
};
