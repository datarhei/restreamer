/* eslint-disable import/no-anonymous-default-export */
import base from '../base';

export default {
	styleOverrides: {
		list: {
			outline: '0',
			color: base.palette.text.primary,
			backgroundColor: base.palette.primary.main,
			ontSize: '1rem',
			fontWeight: 400,
			lineHeight: '2rem',
			'&option': {
				paddingLeft: '10px!important',
				paddingRight: '10px!important',
				'&:hover': {
					backgroundColor: base.palette.primary.dark,
				},
				'&.Mui-focused': {
					backgroundColor: base.palette.primary.dark,
				},
				'&[aria-selected="true"]': {
					fontWeight: 'bold',
					backgroundColor: base.palette.background.dark1,
				},
			},
		},
		paper: {
			backgroundColor: base.palette.primary.main,
			padding: '.3em 0em .3em 0em!important',
		},
	},
};
