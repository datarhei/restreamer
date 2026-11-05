import React from 'react';

// https://github.com/omgovich/react-colorful
import { RgbaStringColorPicker } from 'react-colorful';
import TextField from '@mui/material/TextField';

// https://stackoverflow.com/questions/5623838/rgb-to-hex-and-hex-to-rgb
function hexToRGBA(value) {
	// Expand shorthand form (e.g. "03F") to full form (e.g. "0033FF")
	var shorthandRegex = /^#([a-f\d])([a-f\d])([a-f\d])$/i;
	value = value.replace(shorthandRegex, function (m, r, g, b) {
		return '#' + r + r + g + g + b + b;
	});

	var result = /^#([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(value);
	return result ? `rgba(${parseInt(result[1], 16)},${parseInt(result[2], 16)},${parseInt(result[3], 16)},1)` : value;
}

export default function ColorPicker(props) {
	const [$open, setOpen] = React.useState(false);

	const handleOpen = () => {
		setOpen(true);
	};

	const handleClose = () => {
		setOpen(false);
	};

	const handleChange = (color) => {
		props.onChange({
			target: {
				value: color,
			},
		});
	};

	const value = hexToRGBA(props.value);

	return (
		<React.Fragment>
			<TextField variant={props.variant} fullWidth={props.fullWidth} label={props.label} value={value} onClick={handleOpen} onChange={props.onChange} />
			{$open ? (
				<div style={{ position: 'absolute', zIndex: '2' }}>
					<div
						style={{
							position: 'fixed',
							top: '0px',
							right: '0px',
							bottom: '0px',
							left: '0px',
						}}
						onClick={handleClose}
					/>
					<RgbaStringColorPicker color={value} onChange={handleChange} />
				</div>
			) : null}
		</React.Fragment>
	);
}

ColorPicker.defaultProps = {
	variant: 'default',
	label: '',
	fullWidth: false,
	value: 'rgba(255, 255, 255, 1)',
	onChange: () => {},
};
