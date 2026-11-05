/* eslint-disable import/no-anonymous-default-export */
import base from '../base';

export default {
	styleOverrides: {
		root: {
			'&:hover': {
				backgroundColor: base.palette.background.box_default,
			},
			'&.Mui-selected': {
				backgroundColor: base.palette.background.box_default,
				'&:hover': {
					backgroundColor: base.palette.background.box_default,
				},
			},
		},
	},
};
