import React from 'react';

import { Trans } from '@lingui/macro';
import Grid from '@mui/material/Grid';
import Typography from '@mui/material/Typography';

import ActionButton from '../../misc/ActionButton';
import BoxText from '../../misc/BoxText';
import Duration from '../../misc/Duration';
import Progress from '../../misc/Progress';

function init(progress) {
	const initProgress = {
		state: 'disconnected',
		reconnect: -1,
		error: '',
		time: 0,
		fps: 0,
		bitrate: 0,
		speed: 0,
		...progress,
	};

	return initProgress;
}

export default function Process(props) {
	const progress = init(props.progress);

	const handleAction = (action) => () => {
		props.onAction(action);
	};

	return (
		<Grid container spacing={1}>
			<Grid item xs={12}>
				{progress.state === 'connecting' && (
					<BoxText>
						<Typography variant="body2" gutterBottom>
							<Trans>Connecting ...</Trans>
						</Typography>
					</BoxText>
				)}
				{progress.state === 'disconnecting' && (
					<BoxText>
						<Typography variant="body2" gutterBottom>
							<Trans>Disconnecting ...</Trans>
						</Typography>
					</BoxText>
				)}
				{progress.state === 'connected' && (
					<BoxText>
						<Typography variant="body2" gutterBottom>
							<Trans>
								Connected since <Duration seconds={progress.time} />
							</Trans>
						</Typography>
					</BoxText>
				)}
				{progress.state === 'error' && (
					<BoxText color="danger">
						<Typography variant="body2">
							{progress.error === '' && <Trans>Error</Trans>}
							{progress.error !== '' && <Trans>Error: {progress.error}</Trans>}
						</Typography>
						{progress.reconnect !== -1 && (
							<Typography variant="body2" style={{ marginTop: 0, marginBottom: '.2em' }}>
								<Trans>Reconnecting in {progress.reconnect}s</Trans>
							</Typography>
						)}
						{progress.reconnect === -1 && (
							<Typography variant="body2" style={{ marginTop: 0, marginBottom: '.2em' }}>
								<Trans>You have to reconnect manually</Trans>
							</Typography>
						)}
					</BoxText>
				)}
				{progress.state === 'disconnected' && progress.reconnect !== -1 && (
					<BoxText>
						<Typography variant="body2" style={{ marginTop: 0, marginBottom: '.2em' }}>
							<Trans>Reconnecting in {progress.reconnect}s</Trans>
						</Typography>
					</BoxText>
				)}
			</Grid>
			<Grid item xs={12}>
				<ActionButton
					order={progress.order}
					state={progress.state}
					reconnect={progress.reconnect}
					onDisconnect={handleAction('disconnect')}
					onConnect={handleAction('connect')}
					onReconnect={handleAction('reconnect')}
				/>
			</Grid>
		</Grid>
	);
}

Progress.defaultProps = {
	progress: {},
	onAction: function (action) {},
};
