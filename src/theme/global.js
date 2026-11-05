/* eslint-disable import/no-anonymous-default-export */
import base from './base';
import universe from '../assets/images/universe-4609408.jpg';

export default {
	html: {
		width: '100%',
		height: '100%',
		fontSize: '16px/1.5',
	},
	body: {
		background: `${base.palette.background.button_disabled} url(${universe}) no-repeat fixed left top`,
		backgroundSize: 'cover',
		overflowX: 'hidden',
		overflowY: 'scroll',
	},
	code: {
		fontFamily: 'soure-code-pro, monospace',
	},
	textarea: {
		width: '100%',
		backgroundColor: base.palette.background.modalbox,
		fontFamily:
			'Consolas, "Andale Mono WT", "Andale Mono", "Lucida Console", "Lucida Sans Typewriter", "DejaVu Sans Mono", "Bitstream Vera Sans Mono", "Liberation Mono", "Nimbus Mono L", Monaco, "Courier New", Courier, monospace',
		fontSize: 14,
		whiteSpace: 'pre',
		overflow: 'auto',
		color: base.palette.text.primary,
		border: '0',
		resize: 'none',
		wordWrap: 'inherit',
		wordBreak: 'inherit',
		overflowWrap: 'normal',
	},
	'*:focus': {
		outline: 'none',
	},
	'::-webkit-scrollbar': {
		width: 6,
		height: 6,
	},
	'::-webkit-scrollbar:vertical': {
		width: 6,
	},
	'::-webkit-scrollbar:horizontal': {
		height: 6,
	},
	'::-webkit-scrollbar-thumb': {
		borderRadius: 4,
		color: base.palette.background.footer1,
		backgroundColor: base.palette.background.footer1,
	},
	'::-webkit-scrollbar-track': {
		backgroundColor: 'transparent',
		borderRadius: 8,
	},
	'::-webkit-scrollbar-corner': {
		background: 'transparent',
	},
	// firefox
	'@-moz-document url-prefix(http://),url-prefix(https://)': {
		scrollbar: {
			MozAppearance: 'none !important',
			background: `${base.palette.background.footer1}!important`,
		},
		'thumb,scrollbarbutton': {
			MozAppearance: 'none !important',
			backgroundColor: 'transparent !important',
		},
		'thumb:hover,scrollbarbutton:hover': {
			MozAppearance: 'none !important',
			backgroundColor: 'transparent !important',
		},
		scrollbarbutton: {
			display: 'none !important',
		},
		'scrollbar[orient="vertical"]': {
			minWidth: '6px !important',
		},
	},
	// Extend build in Player
	// -> removes black bottom line
	'[data-player]': {
		margin: '2px -2px 2px -2px',
	},
};
