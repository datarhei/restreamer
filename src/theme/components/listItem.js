/* eslint-disable import/no-anonymous-default-export */
import base from '../base';

export default {
	styleOverrides: {
		root: {
			backgroundColor: base.palette.background.paper,
			borderBottom: `1px solid ${base.palette.background.darker2}`,
			'&.Mui-selected': {
				backgroundColor: base.palette.background.box_danger,
			},
		},
		button: {
			'&:hover': {
				backgroundColor: `${base.palette.background.darker1} !important`,
			},
		},
		'&.Mui-selected': {
			backgroundColor: base.palette.background.box_danger,
		},
	},
};
