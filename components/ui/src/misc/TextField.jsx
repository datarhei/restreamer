import React from 'react';

import { v4 as uuidv4 } from 'uuid';
import FormControl from '@mui/material/FormControl';
import FormHelperText from '@mui/material/FormHelperText';
import InputAdornment from '@mui/material/InputAdornment';
import InputLabel from '@mui/material/InputLabel';
import OutlinedInput from '@mui/material/OutlinedInput';

import Env from './Env';

export default function Component(props) {
	const id = props.id === null ? uuidv4() : props.id;
	let adornment = null;

	if (props.env) {
		adornment = (
			<InputAdornment position="end">
				<Env />
			</InputAdornment>
		);
	}

	return (
		<FormControl variant="outlined" disabled={props.disabled} fullWidth>
			<InputLabel htmlFor={id}>{props.label}</InputLabel>
			<OutlinedInput
				id={id}
				value={props.value}
				onChange={props.onChange}
				endAdornment={adornment}
				label={props.label}
				multiline={props.multiline}
				rows={props.rows}
				type={props.type}
			/>
			{props.helperText && <FormHelperText>{props.helperText}</FormHelperText>}
		</FormControl>
	);
}

Component.defaultProps = {
	id: null,
	label: '',
	value: '',
	disabled: false,
	multiline: false,
	rows: 1,
	env: false,
	type: 'text',
	helperText: null,
	onChange: function (value) {},
};
