/* eslint-disable import/no-anonymous-default-export */
import base from '../base';

export default {
	styleOverrides: {
		root: {
			color: base.palette.text.primary,
			backgroundColor: base.palette.background.paper,
			padding: '.85em 1.25em .85em 1.25em',
		},
	},
	variants: [
		{
			props: { variant: 'modal' },
			style: {
				padding: '1em 1.5em 1.3em 1.5em',
				maxHeight: '95%',
				overflow: 'scroll',
				backgroundColor: base.palette.background.modal,
				color: base.palette.text.primary,
			},
		},
		{
			props: { variant: 'select' },
			style: {
				padding: 0,
				overflow: 'scroll',
				backgroundColor: base.palette.background.modal,
				color: base.palette.text.primary,
			},
		},
	],
};
