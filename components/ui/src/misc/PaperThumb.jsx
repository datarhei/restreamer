import React from 'react';

import makeStyles from '@mui/styles/makeStyles';
import CardMedia from '@mui/material/CardMedia';
import { withDefaults } from '../utils/props';

const useStyles = makeStyles((theme) => ({
	media: {
		paddingTop: '39.25%',
		borderRadius: 4,
	},
}));

export default function Component(inputProps) {
	const props = withDefaults(inputProps, ComponentDefaultProps);
	const classes = useStyles();

	return <CardMedia className={classes.media} style={{ height: props.height }} image={props.image} title={props.title} />;
}

const ComponentDefaultProps = {
	image: '',
	title: '',
	height: '0px',
};
