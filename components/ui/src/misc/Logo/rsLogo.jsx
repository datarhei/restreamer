import React from 'react';

import makeStyles from '@mui/styles/makeStyles';

import company_logo from './images/rs-logo.svg';

const useStyles = makeStyles((theme) => ({
	Logo: {
		height: 95,
	},
}));

export default function Logo(props) {
	const classes = useStyles();

	let link = 'https://github.com/datarhei/restreamer';

	// eslint-disable-next-line no-useless-escape
	return (
		<a href={link} className={classes.Logo} target="_blank" rel="noopener noreferrer">
			<img src={company_logo} alt="datarhei logo" />
		</a>
	);
}
