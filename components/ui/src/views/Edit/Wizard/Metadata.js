import React from 'react';

import { Trans } from '@lingui/macro';
import Button from '@mui/material/Button';
import Divider from '@mui/material/Divider';
import Grid from '@mui/material/Grid';
import Typography from '@mui/material/Typography';

import MetadataControl from '../../../misc/controls/Metadata';

import Paper from '../../../misc/Paper';
import PaperHeader from '../../../misc/PaperHeader';

export default function Metadata(props) {
	return (
		<Paper xs={12} sm={9} md={6} marginBottom="6em" className="PaperM">
			<PaperHeader spacing={2} variant="h1" title={<Trans>Metadata</Trans>} onAbort={props.onAbort} onHelp={props.onHelp} />
			<Grid container spacing={2}>
				<Grid item xs={12}>
					<Divider />
				</Grid>
				<Grid item xs={12}>
					<Typography>
						<Trans>Briefly describe what the audience will see during the live stream.</Trans>
					</Typography>
				</Grid>
				<Grid item xs={12}>
					<MetadataControl settings={props.metadata} onChange={props.onChange} />
				</Grid>
				<Grid item xs={12}>
					<Divider />
				</Grid>
				<Grid item xs={3}>
					<Button variant="outlined" color="default" fullWidth onClick={props.onBack}>
						<Trans>Back</Trans>
					</Button>
				</Grid>
				<Grid item xs={9}>
					<Button variant="outlined" fullWidth color="primary" onClick={props.onNext}>
						<Trans>Next</Trans>
					</Button>
				</Grid>
			</Grid>
		</Paper>
	);
}

Metadata.defaultProps = {
	onAbort: () => {},
	onHelp: () => {},
	onBack: () => {},
	onNext: () => {},
	onChange: (metadata) => {},
	metadata: {},
};
