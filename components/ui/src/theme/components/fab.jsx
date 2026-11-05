/* eslint-disable import/no-anonymous-default-export */
import base from '../base';

export default {
	styleOverrides: {
		primary: {
			color: base.palette.selected.main,
			backgroundColor: base.palette.background.paper,
			'&:hover': {
				color: base.palette.secondary.contrastText,
				backgroundColor: base.palette.secondary.main,
			},
		},
	},
};
