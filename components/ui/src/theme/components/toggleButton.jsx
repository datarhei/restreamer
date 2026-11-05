/* eslint-disable import/no-anonymous-default-export */
import base from '../base';

export default {
	styleOverrides: {
		root: {
			color: base.palette.text.primary,
			border: `2px solid ${base.palette.text.disabled}`,
			backgroundColor: base.palette.background.dark1,
			'&:hover': {
				backgroundColor: base.palette.secondary.main,
				border: `2px solid ${base.palette.secondary.main}`,
			},
			'&.Mui-selected': {
				color: base.palette.text.primary,
				backgroundColor: base.palette.secondary.main,
				border: `2px solid ${base.palette.secondary.main}`,
				'&:hover': {
					backgroundColor: `${base.palette.secondary.main}!important`,
				},
			},
		},
	},
};
