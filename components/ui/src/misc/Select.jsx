import React from 'react';

import InputLabel from '@mui/material/InputLabel';
import FormControl from '@mui/material/FormControl';
import Select from '@mui/material/Select';
import { withDefaults } from '../utils/props';

export default function Component(inputProps) {
	const props = withDefaults(inputProps, ComponentDefaultProps);
	return (
		<FormControl variant={props.variant} fullWidth>
			<InputLabel>{props.label}</InputLabel>
			<Select value={props.value} onChange={props.onChange} disabled={props.disabled} label={props.label} MenuProps={{ disableScrollLock: true }}>
				{props.children}
			</Select>
		</FormControl>
	);
}

const ComponentDefaultProps = {
	variant: 'outlined',
	label: '',
	value: '',
	disabled: false,
	onChange: function (event) {},
};
