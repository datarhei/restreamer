import React from 'react';

import makeStyles from '@mui/styles/makeStyles';
import CloseIcon from '@mui/icons-material/Close';
import Grid from '@mui/material/Grid';
import HelpIcon from '@mui/icons-material/HelpOutline';
import IconButton from '@mui/material/IconButton';
import Modal from '@mui/material/Modal';
import Paper from '@mui/material/Paper';
import Typography from '@mui/material/Typography';

const useStyles = makeStyles((theme) => ({
	modalHeader: {
		marginBottom: '.7em',
		'& button': {
			float: 'right',
			marginLeft: '.5em',
		},
	},
	modalFooter: {
		marginTop: '1.2em',
		minHeight: '38px',
		'& button': {
			marginRight: '.5em',
		},
		'& div button': {
			float: 'right',
			marginRight: '0',
			marginLeft: '.5em',
		},
	},
	modalPaper: {
		padding: '1em 1.5em 1.3em 1.5em',
		width: '95%',
		maxWidth: 540,
		maxHeight: '95%',
		overflow: 'scroll',
		backgroundColor: theme.palette.background.modal,
		color: theme.palette.text.primary,
	},
}));

// todo: use MuiDialog

const Component = React.forwardRef((props, ref) => {
	const classes = useStyles();

	const paperStyle = {};

	if (props.maxWidth > 0) {
		paperStyle.maxWidth = props.maxWidth + 'px';
	}

	return (
		<Modal open={props.open} onClose={props.onClose} className="modal" disableScrollLock>
			<Paper className={classes.modalPaper} elevation={0} ref={ref} tabIndex={-1} style={paperStyle}>
				<Grid container spacing={0}>
					<Grid item xs={12} className={classes.modalHeader}>
						<Typography variant="button">{props.title}</Typography>
						{typeof props.onClose === 'function' && (
							<IconButton color="inherit" size="small" onClick={props.onClose}>
								<CloseIcon fontSize="small" />
							</IconButton>
						)}
						{typeof props.onHelp === 'function' && (
							<IconButton color="inherit" size="small" onClick={props.onHelp}>
								<HelpIcon fontSize="small" />
							</IconButton>
						)}
					</Grid>
				</Grid>
				<Grid item xs={12}>
					{props.children}
				</Grid>
				<Grid container spacing={0}>
					<Grid item xs={12} className={classes.modalFooter}>
						<div>{props.buttonsRight}</div>
						{props.buttonsLeft}
					</Grid>
				</Grid>
			</Paper>
		</Modal>
	);
});

export default Component;

Component.defaultProps = {
	open: false,
	title: '',
	onClose: null,
	onHelp: null,
	buttonsRight: null,
	buttonsLefts: null,
	maxWidth: -1,
};
