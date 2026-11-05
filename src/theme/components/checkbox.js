/* eslint-disable import/no-anonymous-default-export */
import base from '../base';

export default {
	styleOverrides: {
		colorSecondary: {
			'&.Mui-disabled': {
				'&.MuiIconButton-label': {
					color: base.palette.text.disabled,
				},
			},
		},
	},
};
