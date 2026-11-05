import React from 'react';

import { Trans } from '@lingui/macro';
import Button from '@mui/material/Button';
import Divider from '@mui/material/Divider';
import Grid from '@mui/material/Grid';
import Typography from '@mui/material/Typography';

import LicenseControl from '../../../misc/controls/License';

import Paper from '../../../misc/Paper';
import PaperHeader from '../../../misc/PaperHeader';

export default function License(props) {
	return (
		<Paper xs={12} sm={9} md={6} marginBottom="6em" className="PaperM">
			<PaperHeader spacing={2} variant="h1" title={<Trans>License</Trans>} onAbort={props.onAbort} onHelp={props.onHelp} />
			<Grid container spacing={2}>
				<Grid item xs={12}>
					<Divider />
				</Grid>
				<Grid item xs={12}>
					<Typography>
						<Trans>
							Use your copyright and choose the correct image license. Whether free for all or highly restricted. Briefly discuss what others are
							allowed to do with your image.
						</Trans>
					</Typography>
				</Grid>
				<Grid item xs={12}>
					<LicenseControl license={props.license} onChange={props.onChange} />
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
					<Button variant="outlined" color="primary" fullWidth onClick={props.onNext}>
						<Trans>Save</Trans>
					</Button>
				</Grid>
			</Grid>
		</Paper>
	);
}

License.defaultProps = {
	onAbort: () => {},
	onHelp: () => {},
	onBack: () => {},
	onNext: () => {},
	onChange: (license) => {},
	license: '',
};
