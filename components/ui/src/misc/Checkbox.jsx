import React from 'react';

import Checkbox from '@mui/material/Checkbox';
import FormControlLabel from '@mui/material/FormControlLabel';
import makeStyles from '@mui/styles/makeStyles';
import { withDefaults } from '../utils/props';

const useStyles = makeStyles((theme) => ({
	root: {
		color: theme.palette.text.pirmary,
		'&.Mui-checked': {
			color: theme.palette.text.primary,
		},
		'&.Mui-disabled': {
			color: theme.palette.text.disabled,
		},
	},
	checked: {},
	disabled: {},
}));

export default function Component(inputProps) {
	const props = withDefaults(inputProps, ComponentDefaultProps);
	const classes = useStyles();

	return (
		<FormControlLabel
			className={classes.root}
			control={<Checkbox className={classes.root} checked={props.checked} onChange={props.onChange} />}
			label={props.label}
			disabled={props.disabled}
		/>
	);
}

const ComponentDefaultProps = {
	label: '',
	checked: false,
	disabled: false,
	onChange: function (event) {},
};
