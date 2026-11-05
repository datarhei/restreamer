/* eslint-disable import/no-anonymous-default-export */
import base from '../base';

export default {
	styleOverrides: {
		root: {
			color: base.palette.text.primary,
			backgroundColor: base.palette.primary.dark,
			'& .MuiOutlinedInput-notchedOutline': {
				borderColor: base.palette.primary.main,
				borderWidth: 2,
			},
			'&:hover .MuiOutlinedInput-notchedOutline': {
				borderColor: base.palette.primary.light,
			},
			'& .Mui-focused': {
				borderColor: base.palette.secondary.main,
			},
			'& .Mui-focused .MuiOutlinedInput-notchedOutline': {
				borderColor: base.palette.secondary.main,
			},
			'&.Mui-disabled': {
				'& fieldset.MuiOutlinedInput-notchedOutline': {
					backgroundColor: base.palette.background.dark2,
					borderColor: `${base.palette.primary.main}!important`,
				},
			},
		},
		input: {
			'&:-webkit-autofill': {
				WebkitTextFillColor: base.palette.text.primary,
				WebkitBoxShadow: `0 0 0 1000px ${base.palette.primary.dark} inset`,
			},
		},
	},
};
