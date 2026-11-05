import React from 'react';

import { Trans } from '@lingui/macro';
import makeStyles from '@mui/styles/makeStyles';
import Grid from '@mui/material/Grid';
import IconButton from '@mui/material/IconButton';
import Typography from '@mui/material/Typography';

import AutoFixHighIcon from '@mui/icons-material/AutoFixHigh';
import EditIcon from '@mui/icons-material/Edit';

import * as M from '../../utils/metadata';
import Summary from './Summary';

function IconWizard(props) {
	return <AutoFixHighIcon {...props} />;
}

function IconEdit(props) {
	return <EditIcon {...props} />;
}

const useStyles = makeStyles((theme) => ({
	title: {
		float: 'left',
	},
	edit: {
		float: 'right',
		marginTop: '.5em',
	},
}));

export default function ProfileSummary(props) {
	const classes = useStyles();

	const profile = M.initProfile(props.profile);

	const handleEdit = (what) => () => {
		props.onEdit(what);
	};

	const handleWizard = () => () => {
		props.onWizard();
	};

	return (
		<Grid container spacing={2}>
			<Grid
				item
				xs={12}
				sx={{
					'& button': {
						marginTop: '1em',
						float: 'right',
						marginLeft: '1em',
					},
				}}
			>
				<IconButton size="small" color="inherit" onClick={handleEdit('video')}>
					<IconEdit />
				</IconButton>
				<IconButton size="small" color="inherit" onClick={handleWizard()}>
					<IconWizard />
				</IconButton>
				<Typography variant="h3">
					<Trans>Video settings</Trans>
				</Typography>
			</Grid>
			<Grid item xs={12}>
				<Summary type="video" sources={props.sources} profile={profile.video} />
			</Grid>
			<Grid
				item
				xs={12}
				sx={{
					'& button': {
						marginTop: '1em',
						float: 'right',
						marginLeft: '1em',
					},
				}}
			>
				<IconButton size="small" color="inherit" onClick={handleEdit('audio')}>
					<IconEdit />
				</IconButton>
				<IconButton size="small" color="inherit" onClick={handleWizard()}>
					<IconWizard />
				</IconButton>
				<Typography variant="h3" className={classes.title}>
					<Trans>Audio settings</Trans>
				</Typography>
			</Grid>
			<Grid item xs={12}>
				<Summary type="audio" sources={props.sources} profile={profile.audio} />
			</Grid>
		</Grid>
	);
}

ProfileSummary.defaultProps = {
	sources: [],
	profile: null,
	onEdit: function (type) {},
	onWizard: function () {},
};
