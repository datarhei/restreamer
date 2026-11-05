import React from 'react';

import makeStyles from '@mui/styles/makeStyles';
import Button from '@mui/material/Button';

const useStyles = makeStyles((theme) => ({
	button: {
		fontSize: '.9rem!important',
		height: '56px!important',
	},
}));

export default function Component(props) {
	const classes = useStyles();

	return (
		<Button variant="outlined" size="large" fullWidth color="primary" className={classes.button} {...props}>
			{props.children}
		</Button>
	);
}

Component.defaultProps = {};
