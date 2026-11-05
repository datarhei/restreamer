import React from 'react';

import { Trans } from '@lingui/react/macro';
import makeStyles from '@mui/styles/makeStyles';
import withStyles from '@mui/styles/withStyles';
import Chip from '@mui/material/Chip';
import Tooltip from '@mui/material/Tooltip';
import { withDefaults } from '../utils/props';

const useStyles = makeStyles((theme) => ({
	root: {
		color: theme.palette.text.primary,
		backgroundColor: theme.palette.error.main,
		borderRadius: 4,
		fontWeight: 'bold',
		height: 20,
	},
	labelSmall: {
		paddingLeft: 6,
		paddingRight: 6,
	},
}));

const HtmlTooltip = withStyles((theme) => ({
	tooltip: {
		backgroundColor: theme.palette.error.main,
		color: theme.palette.text.primary,
		maxWidth: 100,
		fontSize: '.8rem',
	},
	arrow: {
		color: theme.palette.error.main,
	},
}))(Tooltip);

export default function Component(inputProps) {
	const props = withDefaults(inputProps, ComponentDefaultProps);
	const classes = useStyles();
	return (
		<HtmlTooltip
			title={
				<React.Fragment>
					<Trans>An environment variable sets this value.</Trans>
				</React.Fragment>
			}
			placement="right"
			arrow
		>
			<Chip size="small" label="ENV" className={classes.root} {...props} />
		</HtmlTooltip>
	);
}

const ComponentDefaultProps = {};
