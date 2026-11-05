import React from 'react';

import { Trans } from '@lingui/macro';
import makeStyles from '@mui/styles/makeStyles';
import Button from '@mui/material/Button';
import CloseIcon from '@mui/icons-material/Close';
import Grid from '@mui/material/Grid';
import HelpIcon from '@mui/icons-material/HelpOutline';
import IconButton from '@mui/material/IconButton';
import Paper from '@mui/material/Paper';
import Stack from '@mui/material/Stack';
import Typography from '@mui/material/Typography';

const useStyles = makeStyles((theme) => ({
	modalHeader: {
		marginBottom: '.7em',
		'& button': {
			float: 'right',
			marginLeft: '.5em',
			paddingTop: '.25em',
			marginRight: '-.7em',
		},
	},
	modalFooter: {
		marginTop: '1.2em',
		minHeight: '38px',
		'& button': {
			marginRight: '.5em',
		},
		'& .right': {
			float: 'right',
			marginRight: '0',
			marginLeft: '.5em',
		},
	},
	modalPaper: {
		padding: '1em 1.5em 1.3em 1.5em',
		width: '95%',
		maxWidth: 980,
		maxHeight: '95%',
		overflow: 'scroll',
		backgroundColor: theme.palette.background.modal,
		color: theme.palette.text.primary,
	},
}));

const Component = React.forwardRef((props, ref) => {
	const classes = useStyles();

	const { title, onClose, onHelp, ...other } = props;

	return (
		<Paper className={classes.modalPaper} elevation={0} tabIndex={-1} ref={ref} {...other}>
			<Grid container spacing={0}>
				<Grid item xs={12} className={classes.modalHeader}>
					<Stack direction="row" justifyContent="space-between" alignItems="center" spacing={2}>
						<Typography variant="button">{props.title}</Typography>
						<Stack direction="row" justifyContent="flex-end" alignItems="center" spacing={2}>
							{typeof props.onHelp === 'function' && (
								<IconButton color="inherit" size="small" onClick={props.onHelp}>
									<HelpIcon fontSize="small" />
								</IconButton>
							)}
							{typeof props.onClose === 'function' && (
								<IconButton color="inherit" size="small" onClick={props.onClose}>
									<CloseIcon fontSize="small" />
								</IconButton>
							)}
						</Stack>
					</Stack>
				</Grid>
			</Grid>
			{props.children}
			<Grid container spacing={0}>
				<Grid item xs={12} className={classes.modalFooter}>
					<Button variant="outlined" color="default" onClick={props.onClose}>
						<Trans>Close</Trans>
					</Button>
				</Grid>
			</Grid>
		</Paper>
	);
});

export default Component;

Component.defaultProps = {
	title: '',
	onClose: null,
	onHelp: null,
};
