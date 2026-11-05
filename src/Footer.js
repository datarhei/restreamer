import React from 'react';

import { isMobile } from 'react-device-detect';
import { Trans } from '@lingui/macro';
import makeStyles from '@mui/styles/makeStyles';
//import Button from '@mui/material/Button';
import Grid from '@mui/material/Grid';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemText from '@mui/material/ListItemText';
import ListSubheader from '@mui/material/ListSubheader';
import Popover from '@mui/material/Popover';
import Stack from '@mui/material/Stack';
import Typography from '@mui/material/Typography';
import WarningIcon from '@mui/icons-material/Warning';

import useInterval from './hooks/useInterval';
import Duration from './misc/Duration';
import Logo from './misc/Logo';
import Number from './misc/Number';

const useStyles = makeStyles((theme) => ({
	footer: {
		zIndex: '2',
		position: 'fixed',
		bottom: 0,
		width: '100%',
		height: 60,
		background: `-webkit-linear-gradient(left, ${theme.palette.background.footer1} 0, ${theme.palette.background.footer2} 100%)`,
		color: theme.palette.text.secondary,
		'& .footerLeft': {
			textOverflow: 'ellipsis',
			overflow: 'hidden !important',
			whiteSpace: 'nowrap',
			marginLeft: 40,
		},
		'& .footerRight': {
			marginLeft: 20,
			marginRight: 40,
		},
		'& .footerVersion': {
			textOverflow: 'ellipsis',
			overflow: 'hidden !important',
			whiteSpace: 'nowrap',
		},
		'@media (max-width: 599px)': {
			'& .footerLeft': {
				marginRight: 20,
			},
			'& .footerRight': {
				marginLeft: 20,
			},
			'& .footerVersion': {
				display: 'none',
			},
		},
	},
	warningIcon: {
		fontSize: '1.1rem',
		marginTop: -1,
		marginRight: 5,
	},
	subheader: {
		color: `${theme.palette.service.main}`,
		textTransform: 'uppercase',
		fontWeight: 'bold',
	},
}));

function Resources(props) {
	const classes = useStyles();
	const [$popover, setPopover] = React.useState(null);
	const [$resources, setResources] = React.useState(null);

	const handlePopoverOpen = (event) => {
		setPopover(event.currentTarget);
	};

	const handlePopoverClose = () => {
		setPopover(null);
	};

	const open = Boolean($popover);

	useInterval(async () => {
		await update();
	}, 2000);

	React.useEffect(() => {
		(async () => {
			await update();
		})();
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

	const update = async () => {
		const resources = await props.resources();
		if (resources === null) {
			return;
		}

		resources.system.mem_used = (resources.system.mem_used_bytes / resources.system.mem_total_bytes) * 100;

		resources.core.disk_used = -1;
		if (resources.core.disk_limit_bytes > 0) {
			resources.core.disk_used = (resources.core.disk_used_bytes / resources.core.disk_limit_bytes) * 100;
		}

		resources.core.memfs_used = -1;
		if (resources.core.memfs_limit_bytes > 0) {
			resources.core.memfs_used = (resources.core.memfs_used_bytes / resources.core.memfs_limit_bytes) * 100;
		}

		resources.core.net_used = -1;
		if (resources.core.net_limit_kbit > 0) {
			resources.core.net_used = (resources.core.net_used_kbit / resources.core.net_limit_kbit) * 100;
		}

		resources.core.sessions = -1;
		if (resources.core.session_limit > 0) {
			resources.core.sessions = (resources.core.session_used / resources.core.session_limit) * 100;
		}

		setResources(resources);
	};

	if ($resources === null) {
		return null;
	}

	const system = $resources.system;
	const core = $resources.core;

	return (
		<Stack className="footerRight" direction="row" alignItems="center" spacing={0}>
			{(system.cpu_used >= 75 || system.mem_used >= 75 || core.memfs_used >= 75 || core.disk_used >= 75 || core.net_used >= 75) && (
				<WarningIcon className={classes.warningIcon} color="service" />
			)}
			<Typography variant="button" noWrap aria-owns={open ? 'mouse-over-popover' : undefined} aria-haspopup="true" onMouseOver={handlePopoverOpen}>
				{system.cpu_used.toFixed(0)}% CPU / {system.mem_used.toFixed(0)}% Mem
			</Typography>
			<Popover
				id="mouse-over-popover"
				open={open}
				onClose={handlePopoverClose}
				anchorEl={$popover}
				disableRestoreFocus
				disableScrollLock
				anchorOrigin={{
					vertical: 'top',
					horizontal: 'left',
				}}
				transformOrigin={{
					vertical: 'top',
					horizontal: 'left',
				}}
				PaperProps={{
					onMouseLeave: isMobile ? null : handlePopoverClose,
				}}
			>
				<List
					size="small"
					subheader={
						<ListSubheader className={classes.subheader}>
							<Trans>Uptime</Trans>
						</ListSubheader>
					}
				>
					<ListItem divider>
						<ListItemText
							primary={
								<Typography variant="body3">
									<Duration seconds={$resources.uptime_seconds} />
								</Typography>
							}
							secondary={``}
						/>
					</ListItem>
				</List>
				<List
					size="small"
					subheader={
						<ListSubheader className={classes.subheader}>
							<Trans>System</Trans>
						</ListSubheader>
					}
				>
					<ListItem divider selected={system.cpu_used >= 75}>
						<ListItemText
							primary={<Typography variant="body3">{system.cpu_used.toFixed(0)}% CPU</Typography>}
							secondary={
								<Typography variant="body2">
									{system.cpu_ncores} <Trans>Cores</Trans>
								</Typography>
							}
						/>
					</ListItem>
					<ListItem divider selected={system.mem_used >= 75}>
						<ListItemText
							primary={
								<Typography variant="body3">
									{system.mem_used.toFixed(0)}% <Trans>Memory</Trans>
								</Typography>
							}
							secondary={
								<Typography variant="body2">
									<Number value={system.mem_used_bytes / 1024 / 1024} /> / <Number value={system.mem_total_bytes / 1024 / 1024} />{' '}
									<Trans>MB</Trans>
								</Typography>
							}
						/>
					</ListItem>
				</List>
				<List
					size="small"
					subheader={
						<ListSubheader className={classes.subheader}>
							<Trans>Application</Trans>
						</ListSubheader>
					}
				>
					<ListItem divider>
						{core.sessions >= 0 ? (
							<ListItemText
								primary={
									<Typography variant="body3">
										{core.sessions.toFixed(0)}% <Trans>Viewer</Trans>
									</Typography>
								}
								secondary={
									<Typography variant="body2">
										<Number value={core.session_used} /> / <Number value={core.session_limit} /> <Trans>Viewer</Trans>
									</Typography>
								}
							/>
						) : (
							<ListItemText
								primary={
									<Typography variant="body3">
										<Trans>Sessions</Trans>
									</Typography>
								}
								secondary={
									<Typography variant="body2">
										<Number value={core.session_used} /> <Trans>Viewer</Trans>
									</Typography>
								}
							/>
						)}
					</ListItem>
					<ListItem divider>
						{core.net_used >= 0 ? (
							<ListItemText
								primary={
									<Typography variant="body3">
										{core.net_used.toFixed(0)}% <Trans>Bandwidth</Trans>
									</Typography>
								}
								secondary={
									<Typography variant="body2">
										<Number value={core.net_used_kbit / 1024} /> / <Number value={core.net_limit_kbit / 1024} /> Mbit/s
									</Typography>
								}
							/>
						) : (
							<ListItemText
								primary={
									<Typography variant="body3">
										<Trans>Bandwidth</Trans>
									</Typography>
								}
								secondary={
									<Typography variant="body2">
										<Number value={core.net_used_kbit / 1024} /> Mbit/s
									</Typography>
								}
							/>
						)}
					</ListItem>
					<ListItem divider>
						{core.memfs_used >= 0 ? (
							<ListItemText
								primary={
									<Typography variant="body3">
										{core.memfs_used.toFixed(0)}% <Trans>In-memory storage</Trans>
									</Typography>
								}
								secondary={
									<Typography variant="body2">
										<Number value={core.memfs_used_bytes / 1024 / 1024} /> / <Number value={core.memfs_limit_bytes / 1024 / 1024} /> MB
									</Typography>
								}
							/>
						) : (
							<ListItemText
								primary={
									<Typography variant="body3">
										<Trans>In-memory storage</Trans>
									</Typography>
								}
								secondary={
									<Typography variant="body2">
										<Number value={core.memfs_used_bytes / 1024 / 1024} /> MB
									</Typography>
								}
							/>
						)}
					</ListItem>
					{/* <ListItem divider> */}
					<ListItem>
						{core.disk_used >= 0 ? (
							<ListItemText
								primary={
									<Typography variant="body3">
										{core.disk_used.toFixed(0)}% <Trans>Disk storage</Trans>
									</Typography>
								}
								secondary={
									<Typography variant="body2">
										<Number value={core.disk_used_bytes / 1024 / 1024} /> / <Number value={core.disk_limit_bytes / 1024 / 1024} /> MB
									</Typography>
								}
							/>
						) : (
							<ListItemText
								primary={
									<Typography variant="body3">
										<Trans>Disk storage</Trans>
									</Typography>
								}
								secondary={
									<Typography variant="body2">
										<Number value={core.disk_used_bytes / 1024 / 1024} /> MB
									</Typography>
								}
							/>
						)}
					</ListItem>
					{/* <ListItem divider>
						<Button variant="service" color="primary" fullWidth size="large" component="a" href="https://service.datarhei.com" target="blank">
							<Trans>More details</Trans>
						</Button>
					</ListItem> */}
				</List>
			</Popover>
		</Stack>
	);
}

Resources.defaultProps = {
	resources: () => {
		return null;
	},
};

const initVersion = (initialVersion) => {
	if (!initialVersion) {
		initialVersion = {};
	}

	const version = {
		number: 0,
		arch: 'unknown',
		...initialVersion,
	};

	return version;
};

export default function Footer(props) {
	const classes = useStyles();
	const version = initVersion(props.version);

	if (props.expand === true) {
		return (
			<Grid container className={classes.footer} spacing={0} direction="row" alignItems="center">
				<Grid item xs={12}>
					<Stack direction="row" justifyContent="space-between" alignItems="center" spacing={0}>
						<Stack className="footerLeft" direction="row" alignItems="center" spacing={0}>
							<Logo className={classes.logo} />
							<Typography className="footerVersion">
								{props.app} v{version.number} ({version.arch}) {props.name ? '- ' + props.name : ''}
							</Typography>
						</Stack>
						<Resources resources={props.resources} />
					</Stack>
				</Grid>
			</Grid>
		);
	} else {
		return (
			<Grid container className={classes.footer} spacing={0} direction="row" alignItems="center">
				<Grid item xs={12}>
					<Stack direction="row" justifyContent="space-between" alignItems="center" spacing={0}>
						<Stack className="footerLeft" direction="row" alignItems="center" spacing={0}>
							<Logo className={classes.logo} />
						</Stack>
					</Stack>
				</Grid>
			</Grid>
		);
	}
}

Footer.defaultProps = {
	expand: false,
	app: '',
	name: '',
	version: initVersion(),
	resources: () => {
		return null;
	},
};
