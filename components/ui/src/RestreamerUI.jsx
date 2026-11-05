import React from 'react';

import makeStyles from '@mui/styles/makeStyles';
import Alert from '@mui/material/Alert';
import Backdrop from '@mui/material/Backdrop';
import CircularProgress from '@mui/material/CircularProgress';
import Grid from '@mui/material/Grid';
import Snackbar from '@mui/material/Snackbar';

import { NotifyProvider } from './contexts/Notify';
import * as auth0 from './utils/auth0';
import useInterval from './hooks/useInterval';
import ChannelList from './misc/ChannelList';
import Footer from './Footer';
import I18n from './I18n';
import Header from './Header';
import * as M from './utils/metadata';
import Restreamer from './utils/restreamer';
import Router from './Router';
import Views from './views';
import { UI as Version } from './version';
import Changelog from './misc/Changelog';

import SemverGt from 'semver/functions/gt';
import SemverValid from 'semver/functions/valid';

const useStyles = makeStyles((theme) => ({
	MainHeader: {
		height: '132px',
	},
	// todo: one layer
	MainContent: {
		height: '100%',
		'& .MainContent-container': {
			minHeight: 'calc(100vh - 230px)',
		},
		'& .MainContent-item': {
			maxWidth: '980px',
		},
	},
}));

export default function RestreamerUI(props) {
	const classes = useStyles();

	const [$state, setState] = React.useState({
		initialized: false,
		valid: false,
		connected: false,
		compatibility: { compatible: false },
		ingest: false,
		password: false,
		updates: false,
		service: false,
	});
	const [$ready, setReady] = React.useState(false);
	const [$snack, setSnack] = React.useState({
		open: false,
		message: '',
		severity: 'info',
	});
	const [$channelList, setChannelList] = React.useState({
		open: false,
		channelid: '',
		channels: [],
	});
	const [$metadata, setMetadata] = React.useState({});
	const [$changelog, setChangelog] = React.useState({
		open: false,
		current: '',
		previous: '',
	});

	const restreamer = React.useRef(null);

	React.useEffect(() => {
		(async () => {
			await handleMount();
		})();

		return () => {};

		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

	useInterval(() => {
		setState({
			...$state,
			updates: restreamer.current.HasUpdates(),
		});
	}, 1000 * 60);

	const notify = (severity, type, message) => {
		setSnack({
			...$snack,
			open: true,
			message: message,
			severity: severity,
		});

		if (severity === 'success') {
			if (type === 'save:ingest') {
				setState({
					...$state,
					ingest: true,
				});
			}
		} else if (severity === 'error') {
			if (type === 'network') {
				setState({
					...$state,
					initialized: true,
					valid: false,
				});
			} else if (type === 'auth') {
				(async () => {
					await handleLogout();
				})();
			}
		}
	};

	const handleMount = async () => {
		restreamer.current = new Restreamer(props.address);
		restreamer.current.AddListener((event) => {
			notify(event.severity, event.type, event.message);
		});

		// Try if there's still an auth0 session
		if (auth0.init() === true) {
			if (await auth0.isAuthenticated()) {
				const token = await auth0.getToken();
				await restreamer.current.LoginWithToken(token);
			} else {
				const result = await auth0.handleRedirectCallback();
				if (result.initialized === true) {
					if (result.error === true) {
						notify('error', 'auth0', 'Auth0: ' + result.description);
					}
				}
			}
		}

		const valid = await restreamer.current.Validate();

		await checkChangelog();

		setState({
			...$state,
			initialized: true,
			valid: valid,
			connected: restreamer.current.IsConnected(),
			compatibility: restreamer.current.Compatibility(),
			ingest: restreamer.current.HasIngest(),
			password: restreamer.current.Auths().length === 0 && !restreamer.current.ConfigOverrides('api.auth.enable'),
			updates: restreamer.current.HasUpdates(),
			service: restreamer.current.HasService(),
		});

		setReady(true);
	};

	const checkChangelog = async () => {
		let showChangelog = true;

		if (restreamer.current.IsConnected() === true) {
			let metadata = await restreamer.current.GetMetadata(false);
			const channels = await restreamer.current.ListChannels();
			let current = Version.replace('restreamer-', '');
			let previous = '';

			if (SemverValid(current) === null) {
				showChangelog = false;
			}

			if (metadata === null) {
				if (channels.length === 1) {
					const progress = await restreamer.current.GetIngestProgress(channels[0].channelid);
					if (progress.valid === false) {
						// assume fresh installation
						metadata = M.initMetadata(metadata);
						await restreamer.current.SetMetadata({
							...metadata,
							bundle: {
								...metadata.bundle,
								version: current,
							},
						});
						return false;
					}
				}
			}

			metadata = M.initMetadata(metadata);

			if ('version' in metadata.bundle) {
				if (SemverValid(metadata.bundle.version) !== null) {
					previous = metadata.bundle.version;
				}
			}

			if (showChangelog === true) {
				if (SemverValid(previous) === null) {
					previous = '';
				}

				if (previous.length !== 0) {
					if (!SemverGt(current, previous)) {
						showChangelog = false;
					}
				}
			}

			setMetadata({
				...$metadata,
				...metadata,
			});

			setChangelog({
				...$changelog,
				open: showChangelog,
				current: current,
				previous: previous,
			});
		}

		return showChangelog;
	};

	const handleCloseChangelog = async () => {
		await restreamer.current.SetMetadata({
			...$metadata,
			bundle: {
				...$metadata.bundle,
				version: $changelog.current,
			},
		});

		setChangelog({
			...$changelog,
			open: false,
		});
	};

	const handleLogin = async (username, password) => {
		const connected = await restreamer.current.Login(username, password);

		await checkChangelog();

		setState({
			...$state,
			connected: connected,
			compatibility: restreamer.current.Compatibility(),
			ingest: restreamer.current.HasIngest(),
		});

		return connected;
	};

	const handleAuth0 = async () => {
		const token = await auth0.getToken();
		const connected = await restreamer.current.LoginWithToken(token);

		setState({
			...$state,
			connected: connected,
			compatibility: restreamer.current.Compatibility(),
		});
	};

	const handleLogout = async () => {
		setState({
			...$state,
			initialized: false,
			connected: false,
		});

		restreamer.current.Logout();

		if (await auth0.isAuthenticated()) {
			await auth0.logout();
		}

		restreamer.current.Reset();

		const valid = await restreamer.current.Validate();

		setState({
			...$state,
			initialized: true,
			valid: valid,
			connected: restreamer.current.IsConnected(),
			compatibility: restreamer.current.Compatibility(),
			ingest: restreamer.current.HasIngest(),
		});
	};

	const handlePasswordReset = async (username, loginUsername, password, loginPassword) => {
		const data = {
			api: {
				auth: {
					enable: true,
				},
			},
			version: 3,
		};

		if (username.length !== 0) {
			data.api.auth.username = username;
		}

		if (password.length !== 0) {
			data.api.auth.password = password;
		}

		const [, err] = await restreamer.current.ConfigSet(data);
		if (err !== null) {
			notify('error', 'save:settings', `There was an error resetting the password.`);
			return 'ERROR';
		}

		const res = await restreamer.current.ConfigReload();
		if (res === false) {
			notify('error', 'restart', `Restarting the application failed.`);
			return 'ERROR';
		}

		restreamer.current.IgnoreAPIErrors(true);

		const waitFor = (ms) => {
			return new Promise((resolve) => {
				setTimeout(resolve, ms);
			});
		};

		let restarted = false;
		const key = restreamer.current.CreatedAt().toISOString();

		for (let retries = 0; retries <= 60; retries++) {
			await waitFor(1000);

			const about = await restreamer.current.About();
			if (about === null) {
				// Restarted API not yet available
				continue;
			}

			if (about.created_at?.toISOString() === key) {
				// API did not yet restart
				continue;
			}

			restarted = true;
			break;
		}

		if (restarted === true) {
			// After the restart the API requires a login and this means the restart happened
			await restreamer.current.Validate();
			await restreamer.current.Login(loginUsername, loginPassword);

			window.location.reload();
		} else {
			return 'TIMEOUT';
		}

		return 'OK';
	};

	const handlePlayersite = () => {
		document.location.hash = '#/playersite';
	};

	const handleSettings = () => {
		document.location.hash = '#/settings';
	};

	const handleChannelList = () => {
		const channelid = restreamer.current.GetCurrentChannelID();
		const channels = restreamer.current.ListChannels();

		setChannelList({
			...$channelList,
			open: true,
			channelid: channelid,
			channels: channels,
		});
	};

	const handleSelectChannel = (channelid) => {
		restreamer.current.SelectChannel(channelid);
		handleChannelList();

		document.location.hash = `#/${channelid}`;
	};

	const handleCloseChannelList = () => {
		setChannelList({
			...$channelList,
			open: false,
		});
	};

	const handleAddChannel = (name) => {
		const channelid = restreamer.current.CreateChannel(name);
		restreamer.current.SelectChannel(channelid);

		setChannelList({
			...$channelList,
			open: false,
		});

		document.location.hash = `#/${channelid}/edit/wizard`;
	};

	const handleStateChannel = async (channelids) => {
		const processes = await restreamer.current.ListProcesses(['state'], channelids);
		const states = {};

		for (let p of processes) {
			states[p.id] = p.progress.state;
		}

		return states;
	};

	const handleCloseSnack = () => {
		setSnack({
			...$snack,
			open: false,
		});
	};

	const handleResources = async () => {
		return await restreamer.current.Resources();
	};

	if ($ready === false) {
		return (
			<Backdrop open={true}>
				<CircularProgress color="inherit" />
			</Backdrop>
		);
	}

	let version = {};
	let app = '';
	let name = '';
	if ($state.initialized === true) {
		version = restreamer.current.Version();
		app = restreamer.current.App();
		name = restreamer.current.Name();
	}

	let resources = () => {
		return null;
	};

	let view = null;
	if ($state.initialized === false) {
		view = <Views.Initializing />;
	} else {
		if ($state.valid === false) {
			view = <Views.Invalid address={restreamer.current.Address()} />;
		} else if ($state.connected === false) {
			view = (
				<Views.Login
					onLogin={handleLogin}
					auths={restreamer.current.Auths()}
					hasService={$state.service}
					address={restreamer.current.Address()}
					onAuth0={handleAuth0}
				/>
			);
		} else if ($state.compatibility.compatible === false) {
			if ($state.compatibility.core.compatible === false) {
				view = <Views.Incompatible type="core" have={$state.compatibility.core.have} want={$state.compatibility.core.want} />;
			} else if ($state.compatibility.ffmpeg.compatible === false) {
				view = <Views.Incompatible type="ffmpeg" have={$state.compatibility.ffmpeg.have} want={$state.compatibility.ffmpeg.want} />;
			}
		} else if ($state.password === true) {
			view = (
				<Views.Password
					onReset={handlePasswordReset}
					username={restreamer.current.ConfigValue('api.auth.username')}
					usernameOverride={restreamer.current.ConfigOverrides('api.auth.username')}
					password={restreamer.current.ConfigValue('api.auth.password')}
					passwordOverride={restreamer.current.ConfigOverrides('api.auth.password')}
				/>
			);
		} else {
			view = <Router restreamer={restreamer.current} />;
			resources = handleResources;
		}
	}

	const expand = $state.connected && $state.compatibility.compatible && !$state.password;

	return (
		<I18n>
			<NotifyProvider value={{ Dispatch: notify }}>
				<Grid container direction="column" justifyContent="flex-start" alignItems="stretch" spacing={0}>
					<Grid className={classes.MainHeader}>
						<Header
							expand={expand}
							showPlayersite={$state.ingest}
							showSettings={$state.compatibility.compatible}
							hasUpdates={$state.updates}
							hasService={$state.service}
							onChannel={handleChannelList}
							onPlayersite={handlePlayersite}
							onSettings={handleSettings}
							onLogout={handleLogout}
						/>
					</Grid>
					<Grid item className={classes.MainContent}>
						<Grid container className="MainContent-container" justifyContent="center" alignItems="center" spacing={0}>
							<Grid item sm={1}></Grid>
							<Grid item xs={12} sm={10} className="MainContent-item">
								{view}
							</Grid>
							<Grid item sm={1}></Grid>
						</Grid>
					</Grid>
				</Grid>
				<Footer expand={$state.connected} app={app} version={version} name={name} resources={resources} />
				<Snackbar
					anchorOrigin={{
						vertical: 'top',
						horizontal: 'right',
					}}
					open={$snack.open}
					autoHideDuration={6000}
					onClose={handleCloseSnack}
				>
					<Alert variant="filled" elevation={6} onClose={handleCloseSnack} severity={$snack.severity}>
						{$snack.message}
					</Alert>
				</Snackbar>
				{expand && (
					<ChannelList
						open={$channelList.open}
						channels={$channelList.channels}
						channelid={$channelList.channelid}
						onClose={handleCloseChannelList}
						onClick={handleSelectChannel}
						onAdd={handleAddChannel}
						onState={handleStateChannel}
					/>
				)}
				{expand && $changelog.open && (
					<Changelog open={$changelog.open} onClose={handleCloseChangelog} current={$changelog.current} previous={$changelog.previous} />
				)}
			</NotifyProvider>
		</I18n>
	);
}

RestreamerUI.defaultProps = {
	address: '',
};
