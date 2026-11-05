/* eslint-disable import/no-anonymous-default-export */
import base from '../base';

export default {
	styleOverrides: {
		tooltip: {
			fontSize: '1rem',
			fontWeight: 'normal',
			boxShadow: base.shadows[2],
			backgroundColor: base.palette.background.modal,
		},
		arrow: {
			color: base.palette.background.modal,
		},
	},
};
