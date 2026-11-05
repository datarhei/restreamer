import React from 'react';

import makeStyles from '@mui/styles/makeStyles';
import Grid from '@mui/material/Grid';
import Modal from '@mui/material/Modal';
import Stack from '@mui/material/Stack';

import ModalContent from '../ModalContent';
import Textarea from '../Textarea';

const useStyles = makeStyles((theme) => ({
	box: {
		backgroundColor: theme.palette.background.modalbox,
		borderRadius: 4,
		padding: '0em 1em 1em 1em',
	},
}));

const Component = function (props) {
	const classes = useStyles();

	const { open, title, onClose, onHelp, ...other } = props;

	return (
		<Modal open={props.open} onClose={props.onClose} className="modal">
			<ModalContent title={props.title} onClose={props.onClose} onHelp={props.onHelp} style={{ overflow: 'hidden' }}>
				<Grid container spacing={1}>
					<Grid item xs={12}>
						<Stack direction="column" justifyContent="center" alignItems="center" spacing={1} className={classes.box}>
							<Textarea {...other} />
						</Stack>
					</Grid>
				</Grid>
			</ModalContent>
		</Modal>
	);
};

export default Component;

Component.defaultProps = {
	open: false,
	title: '',
	onClose: null,
	onHelp: null,
};
