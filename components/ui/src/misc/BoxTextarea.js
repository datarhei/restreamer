import React from 'react';

import makeStyles from '@mui/styles/makeStyles';
import Stack from '@mui/material/Stack';

const useStyles = makeStyles((theme) => ({
	box: {
		backgroundColor: theme.palette.background.modalbox,
		borderRadius: 4,
		padding: '0em 1em 0em 1em',
		width: '100%',
	},
}));

export default function Component(props) {
	const classes = useStyles();

	return (
		<Stack direction="column" justifyContent="center" alignItems="center" spacing={1} className={classes.box} {...props}>
			{props.children}
		</Stack>
	);
}

Component.defaultProps = {};
