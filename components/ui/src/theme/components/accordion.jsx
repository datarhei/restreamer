/* eslint-disable import/no-anonymous-default-export */
import base from '../base';

export default {
	styleOverrides: {
		root: {
			padding: '0',
		},
		rounded: {
			boxShadow: 'unset',
			backgroundColor: base.palette.background.accordion,
			border: `2px ${base.palette.background.accordion_border} solid`,
		},
	},
};
