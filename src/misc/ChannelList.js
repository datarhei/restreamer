import React from 'react';

import { Trans } from '@lingui/macro';
import { styled } from '@mui/material/styles';
import { useTheme } from '@mui/styles';
import makeStyles from '@mui/styles/makeStyles';
import useMediaQuery from '@mui/material/useMediaQuery';
import AddIcon from '@mui/icons-material/Add';
import Button from '@mui/material/Button';
import ButtonBase from '@mui/material/ButtonBase';
import CloseIcon from '@mui/icons-material/Close';
import DoNotDisturbAltIcon from '@mui/icons-material/DoNotDisturbAlt';
import FullscreenIcon from '@mui/icons-material/Fullscreen';
import FullscreenExitIcon from '@mui/icons-material/FullscreenExit';
import Grid from '@mui/material/Grid';
import IconButton from '@mui/material/IconButton';
import LensIcon from '@mui/icons-material/Lens';
import NavigateBeforeIcon from '@mui/icons-material/NavigateBefore';
import NavigateNextIcon from '@mui/icons-material/NavigateNext';
import Stack from '@mui/material/Stack';
import SwipeableDrawer from '@mui/material/SwipeableDrawer';
import TextField from '@mui/material/TextField';
import Typography from '@mui/material/Typography';

import Dialog from './modals/Dialog';

const useStyles = makeStyles((theme) => ({
	drawerPaper: {
		marginBottom: '60px',
		boxShadow: '0px -20px 10px -14px rgb(0 0 0 / 25%), 0px -20px 10px -10px rgb(0 0 0 / 10%)',
		paddingLeft: 50,
		paddingRight: 50,
	},
	drawerButtons: {
		marginRight: '-15px!important',
	},
	channel: {
		marginBottom: '1em',
	},
	channelBar: {
		maxWidth: '1200px',
	},
	channelItems: {
		width: '100%',
		maxWidth: '980px',
	},
	iconButton: {
		'& .MuiSvgIcon-root': {
			fontSize: '2.5em',
		},
	},
	iconButtonLeft: {
		marginLeft: '-1.5em!important',
		paddingBottom: '1em',
	},
	iconButtonRight: {
		marginRight: '-1.5em!important',
		paddingBottom: '1em',
	},
	imageTitle: {
		textAlign: 'initial',
		padding: '.5em 0em 0em .1em',
	},
}));

const ImageButton = styled(ButtonBase)(({ theme }) => ({
	'&:hover, &.Mui-focusVisible': {
		zIndex: 1,
		'& .MuiImageBackdrop-root': {
			opacity: 0.2,
		},
	},
	'&:disabled': {
		'& .MuiImageBackdrop-root': {
			opacity: 0.2,
		},
		'& .MuiTypography-root': {
			color: 'white',
		},
	},
}));

const Image = styled('span')(({ theme }) => ({
	position: 'relative',
}));

const ImageSrc = styled('span')(({ theme }) => ({
	position: 'absolute',
	left: 0,
	right: 0,
	top: 0,
	bottom: 0,
	backgroundSize: 'cover',
	backgroundPosition: 'center 40%',
	borderRadius: 4,
	border: `2px solid ${theme.palette.primary.dark}`,
}));

const ImageAlt = styled('span')(({ theme }) => ({
	position: 'absolute',
	left: 0,
	right: 0,
	top: 0,
	bottom: 0,
	display: 'flex',
	alignItems: 'center',
	justifyContent: 'center',
	color: theme.palette.common.white,
}));

const ImageBackdrop = styled('span')(({ theme }) => ({
	position: 'absolute',
	left: 0,
	right: 0,
	top: 0,
	bottom: 0,
	backgroundColor: theme.palette.common.black,
	opacity: 0.5,
	transition: theme.transitions.create('opacity'),
	borderRadius: 4,
	border: `2px solid ${theme.palette.primary.dark}`,
}));

function ChannelButton(props, largeChannelList) {
	const classes = useStyles();
	const theme = useTheme();

	let color = theme.palette.primary.main;
	switch (props.state) {
		case 'disconnected':
			color = theme.palette.primary.main;
			break;
		case 'connected':
			color = theme.palette.secondary.main;
			break;
		case 'disconnecting':
		case 'connecting':
			color = theme.palette.warning.main;
			break;
		case 'error':
			color = theme.palette.error.main;
			break;
		default:
			color = theme.palette.primary.main;
			break;
	}

	let color_active = theme.palette.primary.main;
	switch (props.disabled) {
		case true:
			color_active = theme.palette.primary.light;
			break;
		default:
			color_active = theme.palette.primary.dark;
			break;
	}

	return (
		<Grid item xs={12} sm={6} md={4} lg={3} style={{ paddingBottom: largeChannelList ? '10px' : 'auto' }}>
			<ImageButton focusRipple disabled={props.disabled} onClick={props.onClick} style={{ width: props.width }}>
				<Stack direction="column" spacing={0.5}>
					<Image
						style={{
							width: props.width,
							height: parseInt((props.width / 16) * 9),
						}}
					>
						<ImageAlt>
							<DoNotDisturbAltIcon fontSize="large" />
						</ImageAlt>
						<ImageSrc style={{ backgroundImage: `url(${props.url})`, borderColor: color_active }} />
						<ImageBackdrop className="MuiImageBackdrop-root" style={{ borderColor: color_active }} />
					</Image>
					<Stack direction="row" alignItems="flex-start" justifyContent="space-between" className={classes.imageTitle}>
						<Typography variant="body2" color="inherit">
							{props.title}
						</Typography>
						<Typography variant="body2" color="inherit">
							<LensIcon fontSize="small" style={{ color: color }} />
						</Typography>
					</Stack>
				</Stack>
			</ImageButton>
		</Grid>
	);
}

ChannelButton.defaultProps = {
	url: '',
	width: 200,
	title: '',
	state: '',
	disabled: false,
	onClick: () => {},
};

const calculateColumnsPerRow = (breakpointSmall, breakpointMedium, breakpointLarge) => {
	if (breakpointLarge) {
		return 4;
	} else if (breakpointMedium) {
		return 3;
	} else if (breakpointSmall) {
		return 2;
	}
	return 1;
};

const calculateRowsToFit = (windowHeight, thumbnailHeight, otherUIHeight) => {
	return Math.floor((windowHeight - otherUIHeight) / thumbnailHeight);
};

export default function ChannelList(props) {
	const classes = useStyles();
	const theme = useTheme();
	const breakpointSmall = useMediaQuery(theme.breakpoints.up('sm'));
	const breakpointMedium = useMediaQuery(theme.breakpoints.up('md'));
	const breakpointLarge = useMediaQuery(theme.breakpoints.up('lg'));
	const [$pos, setPos] = React.useState(-1);
	const [$nChannels, setNChannels] = React.useState(breakpointSmall ? (breakpointMedium ? (breakpointLarge ? 4 : 3) : 2) : 1);
	const [$channels, setChannels] = React.useState([]);
	const [$addChannel, setAddChannel] = React.useState({
		open: false,
		name: '',
	});
	const [windowWidth, setWindowWidth] = React.useState(window.innerWidth);
	const [windowHeight, setWindowHeight] = React.useState(window.innerHeight);

	const { channels: allChannels, channelid, onClick, onClose, onState } = props;

	const [$largeChannelList, setLargeChannelList] = React.useState(false);

	React.useEffect(() => {
		(async () => {
			await onMount();
		})();
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

	React.useEffect(() => {
		setNChannels(breakpointSmall ? (breakpointMedium ? (breakpointLarge ? 4 : 3) : 2) : 1);
	}, [breakpointSmall, breakpointMedium, breakpointLarge]);

	React.useEffect(() => {
		(async () => {
			if (allChannels.length === 0) {
				return;
			}

			let channels = allChannels
				.sort((a, b) => {
					const aname = a.name.toUpperCase();
					const bname = b.name.toUpperCase();
					return aname < bname ? -1 : aname > bname ? 1 : 0;
				})
				.slice($pos, $pos + $nChannels);

			const states = await onState(channels.map((channel) => channel.id));

			channels = channels.map((channel) => {
				return (
					<ChannelButton
						key={channel.channelid}
						url={channel.thumbnail}
						width={200}
						title={channel.name}
						state={states[channel.id]}
						disabled={channelid === channel.channelid}
						onClick={() => {
							onClick(channel.channelid);
							if ($largeChannelList) {
								onClose();
							}
						}}
						largeChannelList={$largeChannelList}
					/>
				);
			});

			setChannels(channels);
		})();
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [$pos, allChannels, $nChannels, channelid, onClick, onState]);

	const onMount = async () => {
		setPos(0);
	};

	React.useEffect(() => {
		const handleResize = () => {
			setWindowWidth(window.innerWidth);
			setWindowHeight(window.innerHeight);
		};

		window.addEventListener('resize', handleResize);
		return () => window.removeEventListener('resize', handleResize);
	}, []);

	React.useEffect(() => {
		const columns = calculateColumnsPerRow(breakpointSmall, breakpointMedium, breakpointLarge);
		const rows = $largeChannelList ? calculateRowsToFit(windowHeight, 200, 60) : 1;
		setNChannels(rows * columns);
	}, [breakpointSmall, breakpointMedium, breakpointLarge, windowHeight, windowWidth, $largeChannelList]);

	const handleAddChannelDialog = () => {
		setAddChannel({
			...$addChannel,
			open: !$addChannel.open,
			name: '',
		});
	};

	const handleAddChannelChange = (event) => {
		setAddChannel({
			...$addChannel,
			name: event.target.value,
		});
	};

	if (props.open === false) {
		return null;
	}

	if ($pos < 0) {
		return null;
	}

	const handleLargeChannelList = () => {
		setLargeChannelList(!$largeChannelList);
		setPos(0);
	};

	return (
		<React.Fragment>
			<SwipeableDrawer
				anchor="bottom"
				open={props.open}
				onOpen={() => {}}
				onClose={props.onClose}
				sx={{
					marginButtom: 60,
					'& .MuiDrawer-paper': {
						top: $largeChannelList ? '0px!important' : 'auto!important',
						height: $largeChannelList ? '100vh' : 'auto',
					},
				}}
				BackdropProps={{ invisible: true }}
				classes={{
					paper: classes.drawerPaper,
				}}
				disableScrollLock
			>
				<React.Fragment>
					<Grid container spacing={2} justifyContent="center" className={classes.channel}>
						<Grid item xs={12}>
							<Stack direction="row" spacing={1} justifyContent="space-between">
								<Stack direction="row" spacing={1} justifyContent="flex-start">
									<Typography variant="h2">Channels</Typography>
								</Stack>
								<Stack direction="row" spacing={1} justifyContent="flex-end" className={classes.drawerButtons}>
									<IconButton color="inherit" size="large" onClick={handleAddChannelDialog}>
										<AddIcon />
									</IconButton>
									<IconButton color="inherit" size="large" onClick={handleLargeChannelList}>
										{$largeChannelList ? <FullscreenExitIcon /> : <FullscreenIcon />}
									</IconButton>
									<IconButton color="inherit" size="large" onClick={props.onClose}>
										<CloseIcon />
									</IconButton>
								</Stack>
							</Stack>
						</Grid>
						<Grid item xs={12} textAlign="center">
							<Stack direction="row" spacing={0} justifyContent="space-between">
								<Stack direction="row" spacing={0} alignItems="center" className={classes.iconButtonLeft}>
									<IconButton
										onClick={() => {
											setPos($pos - 1);
										}}
										disabled={$pos === 0}
										className={classes.iconButton}
									>
										<NavigateBeforeIcon />
									</IconButton>
								</Stack>
								<Stack direction="row" spacing={0} className={classes.channelItems}>
									<Grid container spacing={0} justifyContent={$largeChannelList ? 'flex-start' : 'center'}>
										{$channels}
									</Grid>
								</Stack>
								<Stack direction="row" spacing={0} alignItems="center" className={classes.iconButtonRight}>
									<IconButton
										onClick={() => {
											setPos($pos + 1);
										}}
										disabled={$pos + $nChannels >= allChannels.length}
										className={classes.iconButton}
									>
										<NavigateNextIcon />
									</IconButton>
								</Stack>
							</Stack>
						</Grid>
					</Grid>
				</React.Fragment>
			</SwipeableDrawer>
			<Dialog
				open={$addChannel.open}
				onClose={handleAddChannelDialog}
				title={<Trans>Add new channel</Trans>}
				buttonsLeft={
					<Button variant="outlined" color="default" onClick={handleAddChannelDialog}>
						<Trans>Abort</Trans>
					</Button>
				}
				buttonsRight={
					<Button
						variant="outlined"
						color="primary"
						disabled={$addChannel.name.length === 0}
						onClick={() => {
							handleAddChannelDialog();
							props.onAdd($addChannel.name);
						}}
					>
						<Trans>Add</Trans>
					</Button>
				}
			>
				<Grid container spacing={2}>
					<Grid item xs={12}>
						<Typography>
							<Trans>Enter a name for the new channel.</Trans>
						</Typography>
					</Grid>
					<Grid item xs={12}>
						<TextField variant="outlined" fullWidth label={<Trans>Name</Trans>} value={$addChannel.name} onChange={handleAddChannelChange} />
					</Grid>
				</Grid>
			</Dialog>
		</React.Fragment>
	);
}

ChannelList.defaultProps = {
	open: false,
	channelid: '',
	channels: [],
	onClose: () => {},
	onClick: (channelid) => {},
	onAdd: (name) => {},
	onState: (channelids) => {
		const states = {};
		for (let channelid of channelids) {
			states[channelid] = '';
		}

		return states;
	},
};
