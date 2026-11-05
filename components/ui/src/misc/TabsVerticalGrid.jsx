import React from 'react';

import makeStyles from '@mui/styles/makeStyles';
import Grid from '@mui/material/Grid';
import { withDefaults } from '../utils/props';

const useStyles = makeStyles((theme) => ({
	grid: {
		flexGrow: '1',
		display: 'flex',
		height: '100%',
	},
}));

export default function Component(inputProps) {
	const props = withDefaults(inputProps, ComponentDefaultProps);
	const classes = useStyles();

	return (
		<Grid item xs={12} className={classes.grid}>
			{props.children}
		</Grid>
	);
}

const ComponentDefaultProps = {
	children: null,
};
