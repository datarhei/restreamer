import React from 'react';

import Grid from '@mui/material/Grid';

const Component = function (props) {
	return (
		<Grid container justifyContent="center" spacing={props.spacing} align={props.textAlign}>
			<Grid item xs={12}>
				{props.children}
			</Grid>
		</Grid>
	);
};

export default Component;

Component.defaultProps = {
	spacing: 3,
	textAlign: 'left',
};
