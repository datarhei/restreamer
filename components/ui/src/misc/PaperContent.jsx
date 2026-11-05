import React from 'react';

import Grid from '@mui/material/Grid';
import { withDefaults } from '../utils/props';

const Component = function (inputProps) {
	const props = withDefaults(inputProps, ComponentDefaultProps);
	return (
		<Grid container justifyContent="center" spacing={props.spacing} align={props.textAlign}>
			<Grid item xs={12}>
				{props.children}
			</Grid>
		</Grid>
	);
};

export default Component;

const ComponentDefaultProps = {
	spacing: 3,
	textAlign: 'left',
};
