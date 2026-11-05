import React from 'react';

import makeStyles from '@mui/styles/makeStyles';
import MenuItem from '@mui/material/MenuItem';
import { withDefaults } from '../utils/props';

const useStyles = makeStyles((theme) => ({
	root: {
		fontWeight: 'bold',
		backgroundColor: theme.palette.background.dark1,
	},
}));

const Component = React.forwardRef((inputProps, ref) => {
	const props = withDefaults(inputProps, ComponentDefaultProps);
	const classes = useStyles();

	const { name, value, selected, ...other } = props;

	return (
		<MenuItem value={props.value} className={props.selected ? classes.root : ''} ref={ref} {...other}>
			{props.name}
		</MenuItem>
	);
});

export default Component;

const ComponentDefaultProps = {
	name: '',
	value: '',
	selected: false,
};
