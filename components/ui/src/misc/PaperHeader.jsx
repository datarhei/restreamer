import React from 'react';

import makeStyles from '@mui/styles/makeStyles';
import AddIcon from '@mui/icons-material/Add';
import CloseIcon from '@mui/icons-material/Close';
import EditIcon from '@mui/icons-material/Edit';
import Grid from '@mui/material/Grid';
import HelpIcon from '@mui/icons-material/HelpOutline';
import IconButton from '@mui/material/IconButton';
import Typography from '@mui/material/Typography';

const useStyles = makeStyles((theme) => ({
	root: {
		marginBottom: '.3em',
		'& button': {
			float: 'right',
			marginLeft: '.5em',
		},
	},
}));

const Component = function (props) {
	const classes = useStyles();

	return (
		<Grid container spacing={props.spacing} padding={props.padding}>
			<Grid item xs={12} className={classes.root}>
				{typeof props.onAbort === 'function' && (
					<IconButton color="inherit" size="small" onClick={props.onAbort}>
						<CloseIcon />
					</IconButton>
				)}
				{typeof props.onEdit === 'function' && (
					<IconButton color="inherit" size="small" onClick={props.onEdit}>
						<EditIcon />
					</IconButton>
				)}
				{typeof props.onAdd === 'function' && (
					<IconButton color="inherit" size="small" onClick={props.onAdd}>
						<AddIcon />
					</IconButton>
				)}
				{typeof props.onHelp === 'function' && (
					<IconButton color="inherit" size="small" onClick={props.onHelp}>
						<HelpIcon />
					</IconButton>
				)}
				<Typography variant={props.variant}>{props.title}</Typography>
			</Grid>
		</Grid>
	);
};

export default Component;

Component.defaultProps = {
	spacing: 0,
	padding: null,
	title: '',
	variant: 'pagetitle',
	onAbort: null,
	onHelp: null,
	onEdit: null,
	onAdd: null,
};
