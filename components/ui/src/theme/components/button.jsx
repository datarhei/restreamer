/* eslint-disable import/no-anonymous-default-export */
import base from '../base';

const root = {
	textAlign: 'center',
	textTransform: 'uppercase',
	userSelect: 'none',
};

const outlined = {
	base: {
		color: base.palette.text.primary,
		backgroundColor: base.palette.background.dark1,
		border: `2px solid ${base.palette.primary.main}`,
		'&:hover': {
			color: base.palette.primary.contrastText,
			backgroundColor: base.palette.background.dark1,
			border: `2px solid ${base.palette.primary.light}`,
		},
		'&:disabled': {
			color: `${base.palette.text.disabled}`,
			border: `2px solid ${base.palette.primary.light}`,
		},
		'&.Mui-disabled': {
			backgroundColor: `${base.palette.background.button_disabled}`,
			border: `2px solid ${base.palette.text.disabled}`,
			color: `${base.palette.text.disabled}`,
		},
	},
	primary: {
		color: base.palette.text.primary,
		backgroundColor: base.palette.background.dark1,
		border: `2px solid ${base.palette.secondary.main}`,
		'&:hover': {
			color: base.palette.secondary.contrastText,
			backgroundColor: base.palette.secondary.main,
			border: `2px solid ${base.palette.secondary.main}`,
		},
	},
	// color secondary: danger
	secondary: {
		color: base.palette.text.primary,
		backgroundColor: base.palette.background.dark1,
		border: `2px solid ${base.palette.error.main}`,
		'&:hover': {
			color: base.palette.error.contrastText,
			backgroundColor: base.palette.error.main,
			border: `2px solid ${base.palette.error.main}`,
		},
	},
};

export default {
	variants: [
		{
			props: { variant: 'big' },
			style: {
				height: 130,
				width: '100%',
				textTransform: 'initial!important',
				color: base.palette.text.primary,
				backgroundColor: base.palette.background.dark1,
				border: `2px solid ${base.palette.primary.main}`,
				'&:hover': {
					color: base.palette.primary.contrastText,
					backgroundColor: base.palette.background.dark1,
					border: `2px solid ${base.palette.secondary.main}!important`,
				},
				'&:active': {
					border: `2px solid ${base.palette.secondary.main}!important`,
				},
				'&.Mui-disabled': {
					border: `2px solid ${base.palette.text.disabled}!important`,
					color: `${base.palette.text.primary}!important`,
					opacity: 0.4,
				},
				'& svg': {
					fontSize: '2.5em',
					maxHeight: '40px',
					maxWidth: '40px',
				},
				'& img': {
					fontSize: '2.5em',
					maxHeight: '40px',
					maxWidth: '40px',
				},
			},
		},
		{
			props: { variant: 'bigSelected' },
			style: {
				height: 130,
				width: '100%',
				textTransform: 'initial!important',
				color: base.palette.text.primary,
				backgroundColor: base.palette.background.dark1,
				border: `2px solid ${base.palette.secondary.main}!important`,
				'&:hover': {
					color: base.palette.primary.contrastText,
					backgroundColor: base.palette.background.dark1,
					border: `2px solid ${base.palette.secondary.main}!important`,
				},
				'& svg': {
					fontSize: '2.5em',
					maxHeight: '40px',
					maxWidth: '40px',
				},
				'& img': {
					fontSize: '2.5em',
					maxHeight: '40px',
					maxWidth: '40px',
				},
			},
		},
		{
			props: { variant: 'service' },
			style: {
				color: base.palette.text.primary,
				backgroundColor: base.palette.background.dark1,
				border: `2px solid ${base.palette.service.main}!important`,
				'&:hover': {
					color: base.palette.common.black,
					backgroundColor: base.palette.service.main,
					border: `2px solid ${base.palette.service.main}!important`,
				},
			},
		},
	],
	styleOverrides: {
		root: { ...root },
		outlined: { ...outlined.base },
		outlinedPrimary: { ...outlined.primary },
		outlinedSecondary: { ...outlined.secondary },
	},
};
