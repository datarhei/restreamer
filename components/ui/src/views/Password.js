import React from 'react';

import { Trans } from '@lingui/macro';
import Backdrop from '@mui/material/Backdrop';
import Button from '@mui/material/Button';
import Divider from '@mui/material/Divider';
import Grid from '@mui/material/Grid';
import LinearProgress from '@mui/material/LinearProgress';
import Typography from '@mui/material/Typography';

import Paper from '../misc/Paper';
import Password from '../misc/Password';
import PaperHeader from '../misc/PaperHeader';
import PaperContent from '../misc/PaperContent';
import PaperFooter from '../misc/PaperFooter';
import TextField from '../misc/TextField';

const generatePassword = (length) => {
	const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';

	let password = '';
	for (let i = 0; i < length; i++) {
		password += characters.charAt(Math.floor(Math.random() * characters.length));
	}

	return password;
};

export default function ResetPassword(props) {
	const [$login, setLogin] = React.useState({
		username: props.username.length === 0 ? 'admin' : props.username,
		password: props.password.length === 0 ? generatePassword(6) + '-' + generatePassword(6) + '-' + generatePassword(6) : props.password,
		passwordConfirm: props.password.length !== 0 ? props.password : '',
		showPassword: props.password.length === 0 ? true : false,
	});
	const [$restart, setRestart] = React.useState({
		restarting: false,
		timeout: false,
	});

	const handleReset = async (event) => {
		event.preventDefault();

		setRestart({
			...$restart,
			restarting: true,
			timeout: false,
		});

		// If the username and/or password are set by an environment variable (override == true), then don't
		// store that password to the config file. By setting them as empty string, the currently stored
		// values won't be changed.

		let username = $login.username;
		if (props.usernameOverride) {
			username = '';
		}

		let password = $login.password;
		if (props.passwordOverride) {
			password = '';
		}

		const res = await props.onReset(username, $login.username, password, $login.password);
		switch (res) {
			case 'ERROR':
				setRestart({
					...$restart,
					restarting: false,
				});
				break;
			case 'TIMEOUT':
				setRestart({
					...$restart,
					restarting: true,
					timeout: true,
				});
				break;
			default:
		}
	};

	const handleChange = (what) => (event) => {
		let value = event.target.value;

		setLogin({
			...$login,
			[what]: value,
		});
	};

	const invalid = $login.username.length === 0 || $login.password.length === 0 || $login.password !== $login.passwordConfirm;

	return (
		<Paper xs={12} sm={10} md={6} className="PaperL">
			<Grid container spacing={3}>
				<Grid item xs={12} align="center">
					<Typography variant="h1">
						<Trans>User registration</Trans>
					</Typography>
				</Grid>
				<Grid item xs={12}>
					<form noValidate>
						<Grid container spacing={3}>
							<Grid item xs={12}>
								<TextField
									variant="outlined"
									fullWidth
									id="username"
									label={<Trans>Enter username</Trans>}
									value={$login.username}
									onChange={handleChange('username')}
									autoComplete="username"
									disabled={props.usernameOverride}
									env={props.usernameOverride}
								/>
							</Grid>
							<Grid item xs={12}>
								<Password
									value={$login.password}
									id="password"
									label={<Trans>Enter password</Trans>}
									onChange={handleChange('password')}
									show={$login.showPassword}
									autoComplete="current-password"
									disabled={props.passwordOverride}
									env={props.passwordOverride}
								/>
							</Grid>
							<Grid item xs={12}>
								<Password
									value={$login.passwordConfirm}
									id="password_confirm"
									label={<Trans>Confirm password</Trans>}
									onChange={handleChange('passwordConfirm')}
									show={$login.showPassword}
									disabled={props.passwordOverride}
									env={props.passwordOverride}
								/>
							</Grid>
							<Grid item xs={12}>
								<Divider />
							</Grid>
							<Grid item xs={12}>
								<Button variant="outlined" color="primary" fullWidth size="large" disabled={invalid} type="submit" onClick={handleReset}>
									<Trans>Register user</Trans>
								</Button>
							</Grid>
						</Grid>
					</form>
				</Grid>
			</Grid>
			<Backdrop open={$restart.restarting}>
				<Paper xs={12} sm={4} md={3}>
					<PaperHeader title={<Trans>Restarting</Trans>} />
					<PaperContent>
						{$restart.timeout === false ? (
							<React.Fragment>
								<Typography variant="body1">
									<Trans>Restarting Restreamer Core ...</Trans>
								</Typography>
								<LinearProgress sx={{ mt: '1em' }} />
							</React.Fragment>
						) : (
							<React.Fragment>
								<Typography variant="body1">
									<Trans>Reconnecting to Restreamer Core failed.</Trans>
								</Typography>
								<Typography variant="body1" sx={{ mt: '1em' }}>
									<Trans>This is not necessarily an error. However, it may take a bit longer for Restreamer Core to restart..</Trans>
								</Typography>
							</React.Fragment>
						)}
					</PaperContent>
					<PaperFooter
						buttonsRight={
							<Button variant="outlined" color="primary" disabled={!$restart.timeout} onClick={() => window.location.reload()}>
								<Trans>Reload</Trans>
							</Button>
						}
					/>
				</Paper>
			</Backdrop>
		</Paper>
	);
}

ResetPassword.defaultProps = {
	onReset: function (username, password) {},
	username: '',
	usernameOverride: false,
	password: '',
	passwordOverride: false,
};
