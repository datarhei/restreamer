/* eslint-disable import/no-anonymous-default-export */
import base from '../base';

export default {
	styleOverrides: {
		root: {
			color: base.palette.text.primary,
			borderColor: base.palette.primary.main,
			borderWidth: '2px',
			'&:focus': {
				borderColor: base.palette.secondary.main,
			},
			'&.Mui-selected': {
				backgroundColor: '#fff',
			},
		},
		icon: {
			color: base.palette.text.secondary,
			'&.Mui-disabled': {
				color: base.palette.text.secondary,
			},
		},
	},
};
