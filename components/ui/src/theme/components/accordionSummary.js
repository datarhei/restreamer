/* eslint-disable import/no-anonymous-default-export */
export default {
	styleOverrides: {
		root: {
			padding: '2px 16px',
			'&.Mui-expanded': {
				minHeight: 48,
			},
			'& .MuiAccordionSummary-expandIconWrapper .MuiSvgIcon-root': {
				color: 'white',
			},
		},
		content: {
			'&.Mui-expanded': {
				margin: '12px 0',
			},
		},
	},
};
