import React from 'react';
import { useNavigate, useParams } from 'react-router-dom';

import { useLingui } from '@lingui/react';
import { Trans, t } from '@lingui/macro';
import makeStyles from '@mui/styles/makeStyles';
import Backdrop from '@mui/material/Backdrop';
import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import Grid from '@mui/material/Grid';
import MenuItem from '@mui/material/MenuItem';
import Tab from '@mui/material/Tab';
import TextField from '@mui/material/TextField';
import Typography from '@mui/material/Typography';

import Checkbox from '../../misc/Checkbox';
import Dialog from '../../misc/modals/Dialog';
import Filesize from '../../misc/Filesize';
import H from '../../utils/help';
import NotifyContext from '../../contexts/Notify';
import Paper from '../../misc/Paper';
import PaperHeader from '../../misc/PaperHeader';
import PaperFooter from '../../misc/PaperFooter';
import Player from '../../misc/Player';
import Select from '../../misc/Select';
import UploadButton from '../../misc/UploadButton';
import TabPanel from '../../misc/TabPanel';
import TabsHorizontal from '../../misc/TabsHorizontal';
import TextFieldCopy from '../../misc/TextFieldCopy';

const useStyles = makeStyles((theme) => ({
	gridContainer: {
		paddingTop: '1em',
	},
	playerL1: {
		padding: '4px 1px 4px 9px',
	},
	playerL2: {
		position: 'relative',
		width: '100%',
		paddingTop: '56.25%',
	},
	playerL3: {
		position: 'absolute',
		top: 0,
		left: 0,
		bottom: 0,
		right: 0,
		backgroundColor: theme.palette.common.black,
	},
}));

const logoImageTypes = [
	{ mimetype: 'image/gif', extension: 'gif', maxSize: 1 * 1024 * 1024 },
	{ mimetype: 'image/png', extension: 'png', maxSize: 1 * 1024 * 1024 },
	{ mimetype: 'image/jpeg', extension: 'jpg', maxSize: 1 * 1024 * 1024 },
	{ mimetype: 'image/svg+xml', extension: 'svg', maxSize: 1 * 1024 * 1024 },
];

const posterImageTypes = [
	{ mimetype: 'image/gif', extension: 'gif', maxSize: 1 * 1024 * 1024 },
	{ mimetype: 'image/png', extension: 'png', maxSize: 1 * 1024 * 1024 },
	{ mimetype: 'image/jpeg', extension: 'jpg', maxSize: 1 * 1024 * 1024 },
];

export default function Edit(props) {
	const classes = useStyles();
	const navigate = useNavigate();
	const { channelid: _channelid } = useParams();
	const { i18n } = useLingui();
	const address = props.restreamer.Address();
	const timeout = React.useRef();
	const notify = React.useContext(NotifyContext);
	const [$player] = React.useState('videojs-public');
	const [$ready, setReady] = React.useState(false);
	const [$state, setState] = React.useState('disconnected');
	const [$metadata, setMetadata] = React.useState({});
	const [$settings, setSettings] = React.useState({});
	const [$tab, setTab] = React.useState('embed');
	const [$revision, setRevision] = React.useState(0);
	const [$saving, setSaving] = React.useState(false);
	const [$error, setError] = React.useState({
		open: false,
		title: '',
		message: '',
	});
	const [$invalid, setInvalid] = React.useState('');

	React.useEffect(() => {
		(async () => {
			await mount();
		})();
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

	React.useEffect(() => {
		if ($invalid.length !== 0) {
			navigate($invalid, { replace: true });
		}
	}, [navigate, $invalid]);

	const mount = async () => {
		const channelid = props.restreamer.SelectChannel(_channelid);
		if (channelid === '' || channelid !== _channelid) {
			setInvalid('/');
			return;
		}

		const proc = await props.restreamer.GetIngest(channelid, ['state', 'metadata']);
		if (proc === null) {
			notify.Dispatch('warning', 'notfound:ingest', i18n._(t`Main channel not found`));
			setInvalid(`/${_channelid}/`);
			return;
		}

		setMetadata(proc.metadata);
		setState(proc.progress.state);
		setSettings(props.restreamer.InitPlayerSettings(proc.metadata.player));

		setReady(true);
	};

	const handleChange =
		(what, section = '') =>
		(event) => {
			const value = event.target.value;
			const settings = $settings;

			if (section === '') {
				if (['autoplay', 'mute', 'statistics', 'chromecast', 'airplay'].includes(what)) {
					settings[what] = !settings[what];
				} else {
					settings[what] = value;
				}
			} else if (section === 'color') {
				settings.color[what] = value;
			} else if (section === 'ga') {
				settings.ga[what] = value;
			} else if (section === 'logo') {
				settings.logo[what] = value;
			}

			if (timeout.current !== null) {
				clearTimeout(timeout.current);
				timeout.current = null;
			}

			timeout.current = setTimeout(() => {
				timeout.current = null;
				setRevision($revision + 1);
			}, 500);

			setSettings({
				...$settings,
				...settings,
			});
		};

	const handleLogoUpload = async (data, extension) => {
		const path = await props.restreamer.UploadLogo(_channelid, data, extension);

		handleChange(
			'image',
			'logo',
		)({
			target: {
				value: path,
			},
		});

		setSaving(false);
	};

	const handlePosterUpload = async (data, extension) => {
		const path = await props.restreamer.UploadPoster(_channelid, data, extension);

		handleChange('poster')({
			target: {
				value: path,
			},
		});

		setSaving(false);
	};

	const handleUploadStart = () => {
		setSaving(true);
	};

	const handleUploadError = (title) => (err) => {
		let message = null;

		switch (err.type) {
			case 'nofiles':
				message = <Trans>Please select a file to upload.</Trans>;
				break;
			case 'mimetype':
				message = (
					<Trans>
						The selected file type ({err.actual}) is not allowed. Allowed file types are {err.allowed.join(', ')}
					</Trans>
				);
				break;
			case 'size':
				message = (
					<Trans>
						The selected file is too big (<Filesize bytes={err.actual} />
						). Only <Filesize bytes={err.allowed} /> are allowed.
					</Trans>
				);
				break;
			case 'read':
				message = <Trans>There was an error during upload: {err.message}</Trans>;
				break;
			default:
				message = <Trans>Unknown upload error</Trans>;
		}

		setSaving(false);

		showUploadError(title, message);
	};

	const showUploadError = (title, message) => {
		setError({
			...$error,
			open: true,
			title: title,
			message: message,
		});
	};

	const hideUploadError = () => {
		setError({
			...$error,
			open: false,
		});
	};

	const handleLogoReset = (event) => {
		// For the cleanup of the core, we need to check the following:
		// 1. is the image on the core or external?
		// 2. is the image used somewhere else?
		// 3. OK via dialog

		handleChange(
			'image',
			'logo',
		)({
			target: {
				value: '',
			},
		});

		handleChange(
			'position',
			'logo',
		)({
			target: {
				value: 'top-left',
			},
		});

		handleChange(
			'link',
			'logo',
		)({
			target: {
				value: '',
			},
		});
	};

	const handlePosterReset = (event) => {
		// For the cleanup of the core, we need to check the following:
		// 1. is the image on the core or external?
		// 2. is the image used somewhere else?
		// 3. OK via dialog

		handleChange('poster')({
			target: {
				value: '',
			},
		});
	};

	const handleDone = async () => {
		setSaving(true);

		const metadata = {
			...$metadata,
			player: $settings,
		};

		await props.restreamer.SetIngestMetadata(_channelid, metadata);
		await props.restreamer.UpdatePlayer(_channelid);

		setSaving(false);

		notify.Dispatch('success', 'save:player', i18n._(t`Player settings saved`));
	};

	const handleChangeTab = (event, value) => {
		setTab(value);
	};

	const handleAbort = () => {
		navigate(`/${_channelid}/`);
	};

	const handleHelp = () => {
		H('player-' + $tab);
	};

	const prepareUrl = (url) => {
		if (url.length === 0) {
			return '';
		}

		if (url.match(/^https?:\/\//) === null) {
			url = address + url;
		}

		try {
			let u = new URL(url);
			u.searchParams.set('_rscache', Math.random());
			return u.href;
		} catch (e) {
			return url + '?' + Math.random();
		}
	};

	if ($ready === false) {
		return null;
	}

	const storage = $metadata.control.hls.storage;
	const manifest = props.restreamer.GetChannelAddress('hls+' + storage, _channelid);
	const poster = $settings.poster ? prepareUrl($settings.poster) : props.restreamer.GetChannelAddress('snapshot+' + storage, _channelid);
	const playerAddress = props.restreamer.GetPublicAddress('player', _channelid);
	const iframeCode = props.restreamer.GetPublicIframeCode(_channelid);
	const logo = { ...$settings.logo, image: prepareUrl($settings.logo.image) };

	return (
		<React.Fragment>
			<Paper xs={12} md={10}>
				<PaperHeader title={<Trans>EDIT: Player</Trans>} onAbort={handleAbort} onHelp={handleHelp} />
				<Grid container spacing={1} className={classes.gridContainer}>
					<Grid item xs={12}>
						<Grid container spacing={0} className={classes.playerL1}>
							<Grid item xs={12} className={classes.playerL2}>
								{$state !== 'connected' ? (
									<Grid container direction="column" className={classes.playerL3} justifyContent="center" alignItems="center" spacing={1}>
										<Grid item>
											<Typography variant="h2">
												<Trans>No video</Trans>
											</Typography>
										</Grid>
									</Grid>
								) : (
									<Player
										key={$revision}
										type={$player}
										source={manifest}
										autoplay={$settings.autoplay}
										mute={$settings.mute}
										poster={poster}
										logo={logo}
										colors={$settings.color}
										statistics={$settings.statistics}
										controls
									/>
								)}
							</Grid>
						</Grid>
					</Grid>
				</Grid>
				<Grid container spacing={0}>
					<Grid item xs={12}>
						<TabsHorizontal value={$tab} onChange={handleChangeTab}>
							<Tab className="tab" label={<Trans>Embed</Trans>} value="embed" />
							<Tab className="tab" label={<Trans>Logo</Trans>} value="logo" />
							<Tab className="tab" label={<Trans>Poster</Trans>} value="poster" />
							<Tab className="tab" label={<Trans>Playback</Trans>} value="playback" />
						</TabsHorizontal>
						<TabPanel value={$tab} index="embed">
							<Grid container spacing={2}>
								<Grid item xs={12}>
									<TextFieldCopy label={<Trans>Player URL</Trans>} value={playerAddress} />
								</Grid>
								<Grid item xs={12}>
									<TextFieldCopy label={<Trans>iframe code</Trans>} value={iframeCode} />
								</Grid>
							</Grid>
						</TabPanel>
						<TabPanel value={$tab} index="colors">
							<Grid container spacing={2}>
								<Grid item xs={6}>
									<TextField
										variant="outlined"
										fullWidth
										id="color-1"
										label={<Trans>Seekbar color</Trans>}
										value={$settings.color.seekbar}
										onChange={handleChange('seekbar', 'color')}
									/>
								</Grid>
								<Grid item xs={6}>
									<TextField
										variant="outlined"
										fullWidth
										id="color-2"
										label={<Trans>Button color</Trans>}
										value={$settings.color.buttons}
										onChange={handleChange('buttons', 'color')}
									/>
								</Grid>
							</Grid>
						</TabPanel>
						<TabPanel value={$tab} index="logo">
							<Grid container spacing={2}>
								<Grid item xs={12} md={9}>
									<TextField
										variant="outlined"
										fullWidth
										id="logo-url"
										label={<Trans>Image URL</Trans>}
										value={$settings.logo.image}
										onChange={handleChange('image', 'logo')}
									/>
								</Grid>
								<Grid item xs={12} md={3}>
									<UploadButton
										label={<Trans>Upload</Trans>}
										acceptTypes={logoImageTypes}
										onStart={handleUploadStart}
										onError={handleUploadError(<Trans>Uploading the logo failed</Trans>)}
										onUpload={handleLogoUpload}
									/>
								</Grid>
								<Grid item xs={12} md={4}>
									<Select
										variant="outlined"
										fullWidth
										label={<Trans>Position</Trans>}
										value={$settings.logo.position}
										onChange={handleChange('position', 'logo')}
									>
										<MenuItem value="top-left">Top-Left</MenuItem>
										<MenuItem value="top-right">Top-Right</MenuItem>
										<MenuItem value="bottom-left">Bottom-Left</MenuItem>
										<MenuItem value="bottom-right">Bottom-Right</MenuItem>
									</Select>
								</Grid>
								<Grid item xs={12} md={8}>
									<TextField
										variant="outlined"
										fullWidth
										id="logo-link"
										label={<Trans>Link</Trans>}
										value={$settings.logo.link}
										onChange={handleChange('link', 'logo')}
									/>
								</Grid>
							</Grid>
						</TabPanel>
						<TabPanel value={$tab} index="poster">
							<Grid container spacing={2}>
								<Grid item xs={12} md={9}>
									<TextField
										variant="outlined"
										fullWidth
										id="poster-url"
										label={<Trans>Poster image URL</Trans>}
										value={$settings.poster}
										onChange={handleChange('poster')}
									/>
								</Grid>
								<Grid item xs={12} md={3}>
									<UploadButton
										label={<Trans>Upload</Trans>}
										acceptTypes={posterImageTypes}
										onStart={handleUploadStart}
										onError={handleUploadError(<Trans>Uploading the poster failed</Trans>)}
										onUpload={handlePosterUpload}
									/>
								</Grid>
							</Grid>
						</TabPanel>
						<TabPanel value={$tab} index="statistic">
							<Grid container spacing={2}>
								<Grid item xs={12} md={6}>
									<TextField
										variant="outlined"
										fullWidth
										id="ga-id"
										label={<Trans>Google Analytics ID</Trans>}
										value={$settings.gaAccount}
										onChange={handleChange('account', 'ga')}
									/>
								</Grid>
								<Grid item xs={12} md={6}>
									<TextField
										variant="outlined"
										fullWidth
										id="ga-name"
										label={<Trans>Google Analytics Tracker Name</Trans>}
										value={$settings.gaName}
										onChange={handleChange('name', 'ga')}
									/>
								</Grid>
								<Grid item xs={12}>
									<Checkbox
										label={<Trans>Enable nerd statistics</Trans>}
										checked={$settings.statistics}
										onChange={handleChange('statistics')}
									/>
								</Grid>
							</Grid>
						</TabPanel>
						<TabPanel value={$tab} index="playback">
							<Grid container spacing={2}>
								<Grid item xs={12}>
									<Checkbox label={<Trans>Autoplay</Trans>} checked={$settings.autoplay} onChange={handleChange('autoplay')} />
									<Checkbox label={<Trans>Mute</Trans>} checked={$settings.mute} onChange={handleChange('mute')} />
									<Checkbox label={<Trans>Chromecast</Trans>} checked={$settings.chromecast} onChange={handleChange('chromecast')} />
									<Checkbox label={<Trans>AirPlay</Trans>} checked={$settings.airplay} onChange={handleChange('airplay')} />
								</Grid>
							</Grid>
						</TabPanel>
					</Grid>
				</Grid>
				<PaperFooter
					buttonsLeft={
						<Button variant="outlined" color="default" onClick={handleAbort}>
							<Trans>Close</Trans>
						</Button>
					}
					buttonsRight={
						<React.Fragment>
							<Button variant="outlined" color="primary" onClick={handleDone}>
								<Trans>Save</Trans>
							</Button>
							{$settings.logo.image && $tab === 'logo' && (
								<Button variant="outlined" color="secondary" onClick={handleLogoReset}>
									<Trans>Reset logo</Trans>
								</Button>
							)}
							{$settings.poster && $tab === 'poster' && (
								<Button variant="outlined" color="secondary" onClick={handlePosterReset}>
									<Trans>Reset poster</Trans>
								</Button>
							)}
						</React.Fragment>
					}
				/>
			</Paper>
			<Backdrop open={$saving}>
				<CircularProgress color="inherit" />
			</Backdrop>
			<Dialog
				open={$error.open}
				title={$error.title}
				onClose={hideUploadError}
				buttonsRight={
					<Button variant="outlined" color="primary" onClick={hideUploadError}>
						<Trans>OK</Trans>
					</Button>
				}
			>
				<Typography variant="body1">{$error.message}</Typography>
			</Dialog>
		</React.Fragment>
	);
}

Edit.defaultProps = {
	restreamer: null,
};
