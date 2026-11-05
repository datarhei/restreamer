import React from 'react';
import PropTypes from 'prop-types';

import { Trans } from '@lingui/macro';

import makeStyles from '@mui/styles/makeStyles';
import Divider from '@mui/material/Divider';
import Grid from '@mui/material/Grid';
import Link from '@mui/material/Link';
import Stack from '@mui/material/Stack';
import Typography from '@mui/material/Typography';

const useStyles = makeStyles((theme) => ({
	serviceIcon: {
		fontSize: '4rem!important',
		maxHeight: 64,
		marginTop: '-0.065em',
	},
	serviceName: {
		marginTop: '-.2rem',
	},
}));

export default function TabContent(props) {
	const classes = useStyles();

	return (
		<Grid container spacing={2}>
			<Grid item xs={12}>
				<Stack direction="row" justifyContent="flex-start" alignItems="center" spacing={2}>
					<props.service.icon className={classes.serviceIcon} />
					<Stack direction="column" justifyContent="center" alignItems="flex-start" spacing={0}>
						<Typography variant="h1" className={classes.serviceName}>
							{props.service.name}
						</Typography>
						<Typography>v{props.service.version}</Typography>
					</Stack>
				</Stack>
			</Grid>
			<Grid item xs={12}>
				<Divider />
			</Grid>
			{props.children}
			<Grid item xs={12}>
				<Divider />
			</Grid>
			<Grid item xs={12}>
				<Typography>
					<Trans>Maintainer:</Trans>{' '}
					<Link color="secondary" target="_blank" href={props.service.author.maintainer.link}>
						{props.service.author.maintainer.name}
					</Link>
				</Typography>
			</Grid>
		</Grid>
	);
}

TabContent.defaultProps = {
	service: null,
};

TabContent.propTypes = {
	service: PropTypes.object.isRequired,
};
