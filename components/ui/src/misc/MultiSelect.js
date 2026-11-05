import React from 'react';

import FormControl from '@mui/material/FormControl';
import InputLabel from '@mui/material/InputLabel';
import OutlinedInput from '@mui/material/OutlinedInput';
import Select from '@mui/material/Select';

const MenuProps = {
	PaperProps: {
		style: {
			width: 250,
		},
	},
};

export default function Component(props) {
	return (
		<FormControl variant={props.variant} disabled={props.disabled} fullWidth>
			<InputLabel>{props.label}</InputLabel>
			<Select multiple value={props.value} onChange={props.onChange} input={<OutlinedInput />} renderValue={props.renderValue} MenuProps={MenuProps}>
				{props.children}
			</Select>
		</FormControl>
	);
}

Component.defaultProps = {
	variant: 'outlined',
	label: '',
	value: [],
	disabled: false,
	renderValue: (selected) => selected.join(', '),
	onChange: function (event) {},
};
