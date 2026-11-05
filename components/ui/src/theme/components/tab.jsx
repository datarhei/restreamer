/* eslint-disable import/no-anonymous-default-export */
import base from '../base';

export default {
	styleOverrides: {
		root: {
			minWidth: '30%',
			textTransform: 'initial',
			margin: '0em 1.4em .2em 0em',
			color: base.palette.text.secondary,
			textAlign: 'left',
			minHeight: '30px',
			borderRadius: '5px',
			alignItems: 'flex-start',
			'&.Mui-selected': {
				color: base.palette.text.primary,
				backgroundColor: `${base.palette.primary.dark}`,
			},
		},
		wrapper: {
			alignItems: 'flex-start',
			padding: 0,
		},
	},
};
