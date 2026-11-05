import React from 'react';

import makeStyles from '@mui/styles/makeStyles';
import Box from '@mui/material/Box';
import PropTypes from 'prop-types';

const useStyles = makeStyles((theme) => ({
	root: {
		padding: 0,
	},
	'& .MuiBox-root': {
		padding: 0,
	},
}));

export default function TabPanel(props) {
	const classes = useStyles();
	const { children, value, index, ...other } = props;

	return (
		<div className={classes.root} role="tabpanel" hidden={value !== index} id={`vertical-tabpanel-${index}`} {...other}>
			{value === index && (
				<Box classes={{ root: classes }} p={0}>
					{children}
				</Box>
			)}
		</div>
	);
}

TabPanel.propTypes = {
	children: PropTypes.node,
	index: PropTypes.any.isRequired,
	value: PropTypes.any.isRequired,
};
