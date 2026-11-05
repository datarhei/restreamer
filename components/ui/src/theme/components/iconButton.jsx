/* eslint-disable import/no-anonymous-default-export */
import base from '../base';

export default {
	styleOverrides: {
		root: {
			color: base.palette.text.primary,
		},
		label: {
			color: base.palette.text.primary,
			'&>img': {
				maxHeight: '19px!important',
				maxWidth: '19px!important',
				marginTop: '-5px!important',
			},
		},
		sizeSmall: {
			padding: '0 2px 0 2px',
		},
	},
};
