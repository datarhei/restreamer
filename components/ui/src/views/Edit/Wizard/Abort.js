import React from 'react';

import { Trans } from '@lingui/macro';
import Grid from '@mui/material/Grid';
import Typography from '@mui/material/Typography';
import Button from '@mui/material/Button';

import Paper from '../../../misc/Paper';
import PaperHeader from '../../../misc/PaperHeader';

export default function Abort(props) {
	return (
		<Paper xs={12} sm={8} md={6} marginBottom="6em" className="PaperM">
			<PaperHeader spacing={3} variant="h1" title={<Trans>Abort</Trans>} onHelp={props.onHelp} />
			<Grid container spacing={3}>
				{props.nchannels <= 1 ? (
					<React.Fragment>
						<Grid item xs={12}>
							<Typography>
								<Trans>You can't abort the wizard because at least one input must be defined.</Trans>
							</Typography>
						</Grid>
						<Grid item xs={12}>
							<Button variant="outlined" color="primary" fullWidth onClick={props.onBack}>
								<Trans>Back</Trans>
							</Button>
						</Grid>
					</React.Fragment>
				) : (
					<React.Fragment>
						<Grid item xs={12}>
							<Typography>
								<Trans>Are you sure you want to abort the wizard?</Trans>
							</Typography>
						</Grid>
						<Grid item xs={6}>
							<Button variant="outlined" color="default" fullWidth onClick={props.onNext}>
								<Trans>Yes</Trans>
							</Button>
						</Grid>
						<Grid item xs={6}>
							<Button variant="outlined" color="primary" fullWidth onClick={props.onBack}>
								<Trans>No</Trans>
							</Button>
						</Grid>
					</React.Fragment>
				)}
			</Grid>
		</Paper>
	);
}

Abort.defaultProps = {
	onHelp: () => {},
	onBack: () => {},
	onNext: () => {},
	nchannels: 0,
};
