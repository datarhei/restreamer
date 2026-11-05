import React from 'react';

import FormControl from '@mui/material/FormControl';
import FormHelperText from '@mui/material/FormHelperText';
import IconButton from '@mui/material/IconButton';
import InputAdornment from '@mui/material/InputAdornment';
import InputLabel from '@mui/material/InputLabel';
import OutlinedInput from '@mui/material/OutlinedInput';
import Visibility from '@mui/icons-material/Visibility';
import VisibilityOff from '@mui/icons-material/VisibilityOff';

import Env from './Env';

export default function Password(props) {
	const [$visible, setVisible] = React.useState(props.show);

	const handleClickShowPassword = () => {
		setVisible(!$visible);
	};

	const handleMouseDownPassword = (event) => {
		event.preventDefault();
	};

	let adornment = (
		<InputAdornment position="end">
			<IconButton edge="end" size="large">
				<VisibilityOff />
			</IconButton>
			{props.env ? <Env /> : null}
		</InputAdornment>
	);

	if (props.disabled === false) {
		adornment = (
			<InputAdornment position="end">
				<IconButton onClick={handleClickShowPassword} onMouseDown={handleMouseDownPassword} edge="end" size="large">
					{$visible ? <Visibility /> : <VisibilityOff />}
				</IconButton>
				{props.env ? <Env /> : null}
			</InputAdornment>
		);
	}

	return (
		<FormControl variant="outlined" disabled={props.disabled} fullWidth>
			<InputLabel htmlFor={props.id}>{props.label}</InputLabel>
			<OutlinedInput
				id={props.id}
				type={$visible && !props.disabled ? 'text' : 'password'}
				value={props.value}
				onChange={props.onChange}
				endAdornment={adornment}
				label={props.label}
				autoComplete={props.autoComplete}
				inputProps={props.inputProps}
				error={props.error}
			/>
			{props.helperText && <FormHelperText>{props.helperText}</FormHelperText>}
		</FormControl>
	);
}

Password.defaultProps = {
	id: 'password',
	label: '',
	value: '',
	disabled: false,
	autoComplete: 'current-password',
	env: false,
	show: false,
	helperText: false,
	inputProps: {},
	error: false,
	onChange: function (value) {},
};
