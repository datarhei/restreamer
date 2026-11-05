import React from 'react';

import makeStyles from '@mui/styles/makeStyles';
import Grid from '@mui/material/Grid';

const useStyles = makeStyles((theme) => ({
	grid: {
		flexGrow: '1',
		display: 'flex',
		height: '100%',
	},
}));

export default function Component(props) {
	const classes = useStyles();

	return (
		<Grid item xs={12} className={classes.grid}>
			{props.children}
		</Grid>
	);
}

Component.defaultProps = {
	children: null,
};
