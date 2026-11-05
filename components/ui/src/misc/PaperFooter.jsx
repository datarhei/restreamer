import React from 'react';

import makeStyles from '@mui/styles/makeStyles';
import Grid from '@mui/material/Grid';

const useStyles = makeStyles((theme) => ({
	root: {
		marginBottom: '.3em',
		marginTop: '1.2em',
		minHeight: '38px',
		'& button': {
			marginRight: '.5em',
		},
		'& div button': {
			float: 'right',
			marginRight: '0',
			marginLeft: '.5em',
		},
	},
}));

const Component = function (props) {
	const classes = useStyles();

	return (
		<Grid container spacing={3}>
			<Grid item xs={12} className={classes.root}>
				<div>{props.buttonsRight}</div>
				{props.buttonsLeft}
			</Grid>
		</Grid>
	);
};

export default Component;

Component.defaultProps = {
	buttonsLeft: null,
	buttonsRight: null,
};
