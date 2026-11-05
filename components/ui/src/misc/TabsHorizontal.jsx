import React from 'react';

import makeStyles from '@mui/styles/makeStyles';
import Box from '@mui/material/Box';
import Tabs from '@mui/material/Tabs';

const useStyles = makeStyles((theme) => ({
	box: {
		borderBottom: `2px solid ${theme.palette.background.light1}`,
		padding: '0em',
		marginBottom: '1em',
		marginTop: '-1.7em',
	},
	tabs: {
		'& .tab': {
			minWidth: '0px',
			margin: 'unset',
			color: theme.palette.text.primary,
			borderRadius: '4px 4px 0px 0px',
		},
		'@media (max-width: 415px)': {
			'& .tab': {
				padding: '10px 6px!important',
			},
		},
	},
}));

export default function Component(props) {
	const classes = useStyles();

	return (
		<Box className={classes.box}>
			<Tabs className={classes.tabs} variant="scrollable" scrollButtons allowScrollButtonsMobile value={props.value} onChange={props.onChange}>
				{props.children}
			</Tabs>
		</Box>
	);
}

Component.defaultProps = {
	value: '',
	children: null,
	onChange: function (event) {},
};
