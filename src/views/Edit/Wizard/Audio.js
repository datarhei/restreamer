import React from 'react';

import { Trans, t } from '@lingui/macro';
import Button from '@mui/material/Button';
import Divider from '@mui/material/Divider';
import Grid from '@mui/material/Grid';
import Typography from '@mui/material/Typography';
import { useLingui } from '@lingui/react';
import FormControlLabel from '@mui/material/FormControlLabel';
import Radio from '@mui/material/Radio';
import RadioGroup from '@mui/material/RadioGroup';
import WarningIcon from '@mui/icons-material/Warning';

import BoxText from '../../../misc/BoxText';

import Paper from '../../../misc/Paper';
import PaperHeader from '../../../misc/PaperHeader';
import Select from '../../../misc/Select';

export default function Audio(props) {
	const { i18n } = useLingui();

	return (
		<Paper xs={12} sm={9} md={6} marginBottom="6em" className="PaperM">
			<PaperHeader spacing={2} variant="h1" title={<Trans>Audio setup</Trans>} onAbort={props.onAbort} onHelp={props.onHelp} />
			<Grid container spacing={2}>
				<Grid item xs={12}>
					<Divider />
				</Grid>
				<Grid item xs={12}>
					{props.status === 'error' && (
						<BoxText color="dark">
							<WarningIcon fontSize="large" color="error" />
							<Typography textAlign="center">
								<Trans>Failed to verify the source. Please check the address.</Trans>
							</Typography>
						</BoxText>
					)}
					{props.status === 'nostream' && (
						<BoxText color="dark">
							<WarningIcon fontSize="large" color="error" />
							<Typography textAlign="center">
								<Trans>The source doesn't provide any audio streams.</Trans>
							</Typography>
						</BoxText>
					)}
					{props.status === 'nocoder' && (
						<BoxText color="dark">
							<WarningIcon fontSize="large" color="error" />
							<Typography textAlign="center">
								<Trans>The source doesn't provide any compatible audio streams.</Trans>
							</Typography>
						</BoxText>
					)}
				</Grid>
				<Grid item xs={12}>
					<RadioGroup row value={props.source} onChange={props.onSource}>
						<Grid container spacing={2}>
							{props.streamList.length === 0 && (
								<Grid item xs={12}>
									<Typography>
										<Trans>
											The video source doesn't provide any compatible audio stream. <strong>Silence audio</strong> is recommended.
											Services e.g. YouTube, Facebook &amp; Co. require an audio channel.
										</Trans>
									</Typography>
								</Grid>
							)}
							{props.streamList.length !== 0 && (
								<React.Fragment>
									<Grid item xs={12}>
										<FormControlLabel value="video" control={<Radio />} label={i18n._(t`Audio from device`)} />
									</Grid>
									<Grid item xs={12}>
										<Select label={<Trans>Stream</Trans>} value={props.stream} onChange={props.onAudioStreamChange}>
											{props.streamList}
										</Select>
									</Grid>
								</React.Fragment>
							)}
							{props.deviceList.length !== 0 && (
								<React.Fragment>
									<Grid item xs={12}>
										<FormControlLabel value="alsa" control={<Radio />} label={i18n._(t`Audio from device`)} />
									</Grid>
									<Grid item xs={12}>
										<Select label={<Trans>Device</Trans>} value={props.address} onChange={props.onAudioDeviceChange}>
											{props.deviceList}
										</Select>
									</Grid>
								</React.Fragment>
							)}
							<Grid item xs={12}>
								<div>
									<FormControlLabel value="silence" control={<Radio />} label={i18n._(t`Silence Audio`)} />
								</div>
								<div>
									<FormControlLabel value="none" control={<Radio />} label={i18n._(t`No audio`)} />
								</div>
							</Grid>
						</Grid>
					</RadioGroup>
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

Audio.defaultProps = {
	onAbort: () => {},
	onHelp: () => {},
	onBack: () => {},
	onNext: () => {},
	onSource: () => {},
	source: '',
	onAudioStreamChange: () => {},
	onAudioDeviceChange: () => {},
	streamList: [],
	deviceList: [],
	status: '',
	stream: 0,
	address: {},
};
