import React from 'react';

import { Trans } from '@lingui/macro';
import Button from '@mui/material/Button';
import Divider from '@mui/material/Divider';
import Grid from '@mui/material/Grid';
import Typography from '@mui/material/Typography';

import Paper from '../../../misc/Paper';
import PaperHeader from '../../../misc/PaperHeader';
import Select from '../../../misc/Select';

export default function VideoProfile(props) {
	return (
		<Paper xs={12} sm={9} md={6} marginBottom="6em" className="PaperM">
			<PaperHeader spacing={2} variant="h1" title={<Trans>Video setup</Trans>} onAbort={props.onAbort} onHelp={props.onHelp} />
			<Grid container spacing={2}>
				<Grid item xs={12}>
					<Divider />
				</Grid>
				<Grid item xs={12}>
					<Typography>
						<Trans>The video source is compatible. Select the desired resolution:</Trans>
					</Typography>
				</Grid>
				<Grid item xs={12}>
					<Select label={<Trans>Profile</Trans>} value={props.stream} onChange={props.onStreamChange}>
						{props.streamList}
					</Select>
				</Grid>
				{props.compatible === false && (
					<React.Fragment>
						{props.encodersList.length === 0 ? (
							<Grid item xs={12}>
								<Typography>
									<Trans>Your stream needs to be encoded, but there's no suitable encoder available.</Trans>
								</Typography>
							</Grid>
						) : (
							<React.Fragment>
								<Grid item xs={12}>
									<Typography>
										<Trans>Your stream needs to be encoded. Choose the desired encoder:</Trans>
									</Typography>
								</Grid>
								{props.decodersList.length >= 2 && (
									<Grid item xs={12}>
										<Select label={<Trans>Decoder</Trans>} value={props.decoder} onChange={props.onDecoderChange}>
											{props.decodersList}
										</Select>
									</Grid>
								)}
								<Grid item xs={12}>
									<Select label={<Trans>Encoder</Trans>} value={props.encoder} onChange={props.onEncoderChange}>
										{props.encodersList}
									</Select>
								</Grid>
							</React.Fragment>
						)}
					</React.Fragment>
				)}
				<Grid item xs={12}>
					<Divider />
				</Grid>
				<Grid item xs={3}>
					<Button variant="outlined" color="default" fullWidth onClick={props.onBack}>
						<Trans>Back</Trans>
					</Button>
				</Grid>
				<Grid item xs={9}>
					<Button
						variant="outlined"
						fullWidth
						color="primary"
						disabled={props.compatible === false && props.encodersList.length === 0}
						onClick={props.onNext}
					>
						<Trans>Next</Trans>
					</Button>
				</Grid>
			</Grid>
		</Paper>
	);
}

VideoProfile.defaultProps = {
	onAbort: () => {},
	onHelp: () => {},
	onBack: () => {},
	onNext: () => {},
	compatible: false,
	onStreamChange: () => {},
	streamList: [],
	stream: 0,
	onDecoderChange: () => {},
	decodersList: [],
	decoder: '',
	onEncoderChange: () => {},
	encodersList: [],
	encoder: '',
};
