import React from 'react';

import { Trans } from '@lingui/macro';
import Grid from '@mui/material/Grid';
import Typography from '@mui/material/Typography';

import Duration from '../../misc/Duration';
import Number from '../../misc/Number';
import Palette from '../../theme/base/palette';

export default function Progress(props) {
	const uptime = props.progress.time;
	const bitrate = props.progress.bitrate;
	const fps = props.progress.fps;
	const speed = props.progress.speed;

	const valueStyle = {
		fontWeight: 'bold',
		marginBottom: -5,
	};

	const divStyle = {
		backgroundColor: Palette.background.box_default,
		textAlign: 'center',
		borderRadius: 4,
		padding: 8,
		marginTop: '.5em',
		marginBottom: '.2em',
	};

	const uptimeStyle = {
		...divStyle,
	};

	const bitrateStyle = {
		...divStyle,
	};

	const fpsStyle = {
		...divStyle,
	};

	if (fps && uptime > 30 && (fps < 10 || speed < 0.99 || speed > 1.01)) {
		fpsStyle.backgroundColor = Palette.background.box_danger;
	} else {
		fpsStyle.backgroundColor = Palette.background.box_default;
	}

	return (
		<Grid container spacing={1}>
			<Grid item xs={4}>
				<div style={uptimeStyle}>
					<Typography variant="body1" style={valueStyle}>
						<Duration seconds={uptime} />
					</Typography>
					<Typography variant="body2">
						<Trans>Uptime</Trans>
					</Typography>
				</div>
			</Grid>
			<Grid item xs={4}>
				<div style={bitrateStyle}>
					<Typography variant="body1" style={valueStyle}>
						<Number value={bitrate} digits={0} />
					</Typography>
					<Typography variant="body2">
						<Trans>kbit/s</Trans>
					</Typography>
				</div>
			</Grid>
			<Grid item xs={4}>
				<div style={fpsStyle}>
					<Typography variant="body1" style={valueStyle}>
						<Number value={fps} digits={0} />
					</Typography>
					<Typography variant="body2">
						<Trans>FPS</Trans>
					</Typography>
				</div>
			</Grid>
		</Grid>
	);
}

Progress.defaultProps = {
	progress: {},
};
