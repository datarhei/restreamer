/* eslint-disable import/no-anonymous-default-export */
import palette from './palette';

import '@fontsource/dosis/300.css';
import '@fontsource/dosis/400.css';
import '@fontsource/dosis/500.css';
import '@fontsource/dosis/700.css';
import '@fontsource/roboto/300.css';
import '@fontsource/roboto/400.css';
import '@fontsource/roboto/500.css';
import '@fontsource/roboto/700.css';

const base = {
	htmlFontSize: 16,
	fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
	fontSize: 12,
	fontWeightLight: 300,
	fontWeightRegular: 400,
	fontWeightMedium: 500,
	fontWeightBold: 700,
	h1: {
		fontFamily: '"Dosis", "Roboto", "Helvetica", "Arial", sans-serif',
		fontSize: '2rem',
		marginBottom: '.25rem',
	},
	h2: {
		fontFamily: '"Dosis", "Roboto", "Helvetica", "Arial", sans-serif',
		fontSize: '1.8rem',
		marginTop: '.25rem',
		marginBottom: '.25rem',
		wordBreak: 'break-word',
	},
	h3: {
		fontFamily: '"Dosis", "Roboto", "Helvetica", "Arial", sans-serif',
		fontSize: '1.35rem',
		marginTop: '.5em',
		wordBreak: 'break-word',
	},
	h4: {
		fontFamily: '"Dosis", "Roboto", "Helvetica", "Arial", sans-serif',
		fontSize: '1.2rem',
		marginTop: '.4em',
	},
	h5: {
		fontFamily: '"Dosis", "Roboto", "Helvetica", "Arial", sans-serif',
		fontSize: '1.1rem',
		marginBottom: '.5rem',
	},
	h6: {
		fontFamily: '"Dosis", "Roboto", "Helvetica", "Arial", sans-serif',
		fontSize: '1rem',
		fontWeight: 500,
	},
	subtitle1: {
		fontSize: '1.2rem',
		marginBottom: '.5rem',
	},
	subtitle2: {
		fontSize: '1.1rem',
		fontWeight: 500,
	},
	body1: {
		fontSize: '1rem',
	},
	body2: {
		fontSize: '.9rem',
		fontWeight: 500,
	},
	body3: {
		fontSize: '1rem',
		fontWeight: 500,
		marginBottom: '.5rem',
	},
	button: {
		fontSize: '.9rem',
	},
	pagetitle: {
		fontSize: '1rem',
		textTransform: 'uppercase',
		fontWeight: 500,
	},
	caption: {
		color: palette.text.hint,
		marginTop: '0.5em',
		fontSize: '.75rem',
		display: 'block',
		fontStyle: 'italic',
	},
};

export default {
	...base,
};
