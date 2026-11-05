import React from 'react';

import { Trans } from '@lingui/macro';
import Backdrop from '@mui/material/Backdrop';
import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import Grid from '@mui/material/Grid';
import Link from '@mui/material/Link';
import MenuItem from '@mui/material/MenuItem';
import TextField from '@mui/material/TextField';
import ToggleButton from '@mui/material/ToggleButton';
import ToggleButtonGroup from '@mui/material/ToggleButtonGroup';
import Typography from '@mui/material/Typography';

import * as auth0 from '../utils/auth0';
import Paper from '../misc/Paper';
import Password from '../misc/Password';
import Select from '../misc/Select';

function hasAuthType(auths, type) {
	for (let i = 0; i < auths.length; i++) {
		if (auths[i].type === type) {
			return true;
		}
	}

	return false;
}

function initAuths(auths) {
	auths = auths
		.map((a) => {
			const params = a.split(' ');
			if (params[0] === 'localjwt') {
				return {
					type: 'local',
				};
			} else if (params[0] === 'auth0') {
				const auth0 = {
					type: 'auth0',
					id: '',
					domain: '',
					audience: '',
					client_id: '',
				};

				for (let i = 1; i < params.length; i++) {
					const kv = params[i].split('=');
					if (kv.length !== 2) {
						continue;
					}

					switch (kv[0]) {
						case 'domain':
							auth0.domain = kv[1];
							break;
						case 'audience':
							auth0.audience = kv[1];
							break;
						case 'clientid':
							auth0.client_id = kv[1];
							break;
						default:
							break;
					}
				}

				if (auth0.domain.length === 0 || auth0.audience.length === 0 || auth0.client_id.length === 0) {
					return null;
				}

				auth0.id = auth0.domain + auth0.audience + auth0.client_id;

				return auth0;
			}

			return null;
		})
		.filter((a) => a !== null)
		.sort((a, b) => {
			if (a.type !== 'auth0' || b.type !== 'auth0' || a.type !== b.type) {
				if (a.type < b.type) {
					return -1;
				}
				if (a.type > b.type) {
					return 1;
				}

				return 0;
			}

			const domaina = a.domain.toUpperCase();
			const domainb = b.domain.toUpperCase();
			if (domaina < domainb) {
				return -1;
			}
			if (domaina > domainb) {
				return 1;
			}

			return 0;
		});

	return auths;
}

export default function Login(props) {
	const [$auths] = React.useState(initAuths(props.auths));
	const [$login, setLogin] = React.useState({
		username: '',
		password: '',
		showPassword: false,
	});
	const [$auth0, setAuth0] = React.useState('none');
	const [$isAuthenticated, setIsAuthenticated] = React.useState(false);
	const [$canUseAuth0] = React.useState(auth0.canUseAuth0());
	const [$loginTarget, setLoginTarget] = React.useState('none');
	const [$loginCheck, setLoginCheck] = React.useState(false);

	React.useEffect(() => {
		(async () => {
			if (hasAuthType($auths, 'auth0') && $canUseAuth0) {
				setLoginTarget('auth0');
			} else if (hasAuthType($auths, 'local')) {
				setLoginTarget('local');
			}

			const auth0auths = $auths.filter((a) => a.type === 'auth0');
			if (auth0auths.length === 1) {
				setAuth0(auth0auths[0].id);
			}

			if ($canUseAuth0 && auth0.init() === true) {
				const isAuthenticated = await auth0.isAuthenticated();
				setIsAuthenticated(isAuthenticated);
			}
		})();
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

	const handleLogin = async (event) => {
		event.preventDefault();

		setLoginCheck(true);

		const res = await props.onLogin($login.username, $login.password);
		if (res === false) {
			setLoginCheck(false);
		}
	};

	const handleAuth0Login = async () => {
		setLoginCheck(true);

		let config = null;
		for (let i = 0; i < $auths.length; i++) {
			if ($auths[i].type !== 'auth0') {
				continue;
			}

			if ($auths[i].id === $auth0) {
				config = {
					domain: $auths[i].domain,
					audience: $auths[i].audience,
					client_id: $auths[i].client_id,
					redirect_uri: window.location.origin,
				};

				break;
			}
		}

		if (config === null) {
			setLoginCheck(false);
			return;
		}

		auth0.setConfig(config);

		auth0.init();
		const isAuthenticated = await auth0.isAuthenticated();

		if (isAuthenticated === false) {
			const redirect = await auth0.login({
				address: props.address,
				hash: window.location.hash.substring(1),
			});

			if (redirect === true) {
				return;
			}
		} else {
			setIsAuthenticated(true);
		}

		setLoginCheck(false);
	};

	const handleAuth0 = async () => {
		setLoginCheck(true);
		await props.onAuth0();
	};

	const handleChange = (what) => (event) => {
		const value = event.target.value;

		setLogin({
			...$login,
			[what]: value,
		});
	};

	const handleAuth0Change = (event) => {
		const value = event.target.value;

		setAuth0(value);
	};

	const handleLoginTargetChange = (event, value) => {
		if (!value) {
			return;
		}

		setLoginTarget(value);
	};

	return (
		<Paper xs={11} sm={10} md={6}>
			<Grid container spacing={3} padding={{ xs: 1, sm: 5 }}>
				{$loginTarget !== 'none' && (
					<Grid item xs={12} align="center">
						{$auths.length > 1 && (
							<ToggleButtonGroup value={$loginTarget} exclusive onChange={handleLoginTargetChange}>
								{hasAuthType($auths, 'auth0') && (
									<ToggleButton value="auth0">
										<Trans>Auth0</Trans>
									</ToggleButton>
								)}
								{hasAuthType($auths, 'local') && (
									<ToggleButton value="local">
										<Trans>Basic</Trans>
									</ToggleButton>
								)}
								{props.hasService === true && !hasAuthType($auths, 'auth0') && (
									<ToggleButton value="service">
										<Trans>Service</Trans>
									</ToggleButton>
								)}
							</ToggleButtonGroup>
						)}
					</Grid>
				)}
				{$loginTarget === 'none' && (
					<React.Fragment>
						<Grid item xs={12} align="center">
							<Typography>
								<Trans>There's no login method available.</Trans>
							</Typography>
						</Grid>
					</React.Fragment>
				)}
				{$loginTarget === 'local' && hasAuthType($auths, 'local') && (
					<React.Fragment>
						<Grid item xs={12}>
							<form noValidate>
								<Grid container spacing={3}>
									<Grid item xs={12}>
										<TextField
											variant="outlined"
											fullWidth
											id="username"
											label={<Trans>Username</Trans>}
											value={$login.username}
											onChange={handleChange('username')}
											autoComplete="username"
										/>
									</Grid>
									<Grid item xs={12}>
										<Password
											value={$login.password}
											id="password"
											label={<Trans>Password</Trans>}
											onChange={handleChange('password')}
											autoComplete="current-password"
										/>
									</Grid>
									<Grid item xs={12}>
										<Button
											variant="outlined"
											color="primary"
											fullWidth
											size="large"
											type="submit"
											name="login_local"
											onClick={handleLogin}
										>
											<Trans>Login</Trans>
										</Button>
									</Grid>
								</Grid>
							</form>
						</Grid>
					</React.Fragment>
				)}
				{props.hasService === true && $loginTarget === 'service' && !hasAuthType($auths, 'auth0') && (
					<React.Fragment>
						<Grid item xs={12} align="center">
							<Grid container spacing={0}>
								<Grid item xs={12}>
									<Typography variant="h3" style={{ marginBottom: 10 }}>
										<Trans>
											Plan: <strong>Starter</strong>
										</Trans>
									</Typography>
									<Typography>
										<Trans>Social-login (OAuth2, 2FA)</Trans>
									</Typography>
									<Typography>
										<Trans>Hosted Restreamer interface</Trans>
									</Typography>
									<Typography>
										<Trans>Advanced monitoring</Trans>
									</Typography>
									<Typography>
										<Trans>Alerting by email</Trans>
									</Typography>
								</Grid>
							</Grid>
						</Grid>
						<Grid item xs={12}>
							<Button variant="service" color="primary" fullWidth size="large" href="https://service.datarhei.com" target="blank">
								<Trans>Sign up (free)</Trans>
							</Button>
						</Grid>
					</React.Fragment>
				)}
				{$loginTarget === 'auth0' && hasAuthType($auths, 'auth0') && (
					<React.Fragment>
						<Grid item xs={12}>
							<Grid container spacing={3}>
								<Grid item xs={12} align="center">
									<Typography>
										<Trans>
											Use Auth0 for your running Restreamer Core. More{' '}
											<Link color="secondary" target="_blank" href="https://github.com/datarhei">
												details
											</Link>
											.
										</Trans>
									</Typography>
								</Grid>
								{$canUseAuth0 === false ? (
									<React.Fragment>
										<Grid item xs={12} align="center">
											<Typography>
												<Trans>Auth0 is currently not available because this interface is loaded from an insecure origin.</Trans>
											</Typography>
										</Grid>
									</React.Fragment>
								) : (
									<React.Fragment>
										{$isAuthenticated === false ? (
											<React.Fragment>
												<Grid item xs={12}>
													<Select label={<Trans>Auth0 Tenant</Trans>} value={$auth0} onChange={handleAuth0Change}>
														<MenuItem value="none" disabled>
															<Trans>Choose tenant ...</Trans>
														</MenuItem>
														{$auths
															.filter((a) => a.type === 'auth0')
															.map((a) => {
																return (
																	<MenuItem key={a.id} value={a.id}>
																		{a.domain}
																	</MenuItem>
																);
															})}
													</Select>
												</Grid>
												<Grid item xs={12}>
													<Button
														variant="outlined"
														color="primary"
														fullWidth
														size="large"
														type="submit"
														name="login_auth0"
														onClick={handleAuth0Login}
														disabled={$auth0 === 'none'}
													>
														<Trans>Login</Trans>
													</Button>
												</Grid>
											</React.Fragment>
										) : (
											<Grid item xs={12}>
												<Button
													variant="outlined"
													color="primary"
													fullWidth
													size="large"
													type="submit"
													name="login_auth0"
													onClick={handleAuth0}
												>
													<Trans>Login</Trans>
												</Button>
											</Grid>
										)}
									</React.Fragment>
								)}
							</Grid>
						</Grid>
					</React.Fragment>
				)}
			</Grid>
			<Backdrop open={$loginCheck}>
				<CircularProgress color="inherit" />
			</Backdrop>
		</Paper>
	);
}

Login.defaultProps = {
	address: '',
	auths: [],
	hasService: false,
	onLogin: function (username, password) {},
	onAuth0: function () {},
};
