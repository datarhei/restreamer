import React from 'react';
import { useNavigate, useParams } from 'react-router-dom';

import { Trans } from '@lingui/macro';
import Button from '@mui/material/Button';
import Grid from '@mui/material/Grid';
import Link from '@mui/material/Link';
import Typography from '@mui/material/Typography';

import welcomeImage from '../assets/images/welcome.png';
import Paper from '../misc/Paper';
import PaperThumb from '../misc/PaperThumb';

export default function Welcome(props) {
	const navigate = useNavigate();
	const { channelid: _channelid } = useParams();

	return (
		<Paper xs={12} md={6} className="PaperM">
			<Grid container justifyContent="center" spacing={2}>
				<Grid item xs={12}>
					<PaperThumb image={welcomeImage} title="Welcome to Restreamer v2" height="200px" />
				</Grid>
				<Grid item xs={12}></Grid>
				<Grid item xs={12}>
					<Typography align="center">
						<Trans>
							Welcome to Restreamer v2, the solution for fast and easy video publishing. Free for private and commercial use. Further help in the{' '}
							<Link color="secondary" target="_blank" href="https://docs.datarhei.com">
								docs
							</Link>
							.
						</Trans>
					</Typography>
				</Grid>
				<Grid item xs={12}></Grid>
				<Grid item xs={12}>
					<Button fullWidth variant="outlined" color="primary" onClick={() => navigate(`/${_channelid}/edit/wizard`)}>
						<Trans>Next: Video setup</Trans>
					</Button>
				</Grid>
			</Grid>
		</Paper>
	);
}
