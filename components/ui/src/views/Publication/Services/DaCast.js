import React from 'react';

import { Trans } from '@lingui/macro';
import Grid from '@mui/material/Grid';
import Link from '@mui/material/Link';
import MenuItem from '@mui/material/MenuItem';
import TextField from '@mui/material/TextField';

import Checkbox from '../../../misc/Checkbox';
import Logo from './logos/dacast.svg';
import Select from '../../../misc/Select';

const id = 'dacast';
const name = 'DaCast';
const version = '1.0';
const stream_key_link = '';
const description = (
	<Trans>
		Transmit the main source to an DaCast RTMP Service. More about the setup{' '}
		<Link color="secondary" target="_blank" href="https://www.dacast.com/support/knowledgebase/intro-to-live-streaming-a-step-by-step-walkthrough/">
			here
		</Link>
		.
	</Trans>
);
const image_copyright = <Trans>Please contact the operator of the service and check what happens.</Trans>;
const author = {
	creator: {
		name: 'datarhei',
		link: 'https://github.com/datarhei',
	},
	maintainer: {
		name: 'datarhei',
		link: 'https://github.com/datarhei',
	},
};
const category = 'platform';
const requires = {
	protocols: ['rtmp'],
	formats: ['flv'],
	codecs: {
		audio: ['aac', 'mp3'],
		video: ['h264'],
	},
};

function ServiceIcon(props) {
	return <img src={Logo} alt="DaCast Logo" {...props} />;
}

function init(settings) {
	const initSettings = {
		protocol: 'rtmp://',
		server: '',
		backup_server: '',
		enable_backup_server: false,
		username: '',
		password: '',
		stream_key: '',
		...settings,
	};

	return initSettings;
}

function Service(props) {
	const settings = init(props.settings);

	const handleChange = (what) => (event) => {
		const value = event.target.value;

		if (what === 'enable_backup_server') {
			settings[what] = !settings[what];
		} else {
			settings[what] = value;
		}

		const output = createOutput(settings);

		props.onChange(output, settings);
	};

	const createOutput = (settings) => {
		const outputs = [];

		const options = ['-f', 'flv'];

		const output_0 = {
			address: settings.protocol + settings.username + ':' + settings.password + '@' + settings.server + '/' + settings.stream_key,
			options: options,
		};

		outputs.push(output_0);

		if (settings.enable_backup_server) {
			const output_1 = {
				address: settings.protocol + settings.username + ':' + settings.password + '@' + settings.server + '/' + settings.stream_key,
				options: options,
			};

			outputs.push(output_1);
		}

		return outputs;
	};

	return (
		<Grid container spacing={2}>
			<Grid item xs={12} md={3}>
				<Select type="select" label={<Trans>Protocol</Trans>} value={settings.protocol} onChange={handleChange('protocol')}>
					<MenuItem value="rtmp://">rtmp://</MenuItem>
					<MenuItem value="rtmps://">rtmps://</MenuItem>
				</Select>
			</Grid>
			<Grid item xs={12} md={9}>
				<TextField
					variant="outlined"
					fullWidth
					type="url"
					label={<Trans>Server</Trans>}
					value={settings.server}
					onChange={handleChange('server')}
					required
				/>
			</Grid>
			<Grid item xs={12} md={6}>
				<TextField
					variant="outlined"
					fullWidth
					type="url"
					label={<Trans>Username</Trans>}
					value={settings.username}
					onChange={handleChange('username')}
					required
				/>
			</Grid>
			<Grid item xs={12} md={6}>
				<TextField
					variant="outlined"
					fullWidth
					type="url"
					label={<Trans>Password</Trans>}
					value={settings.password}
					onChange={handleChange('password')}
					required
				/>
			</Grid>
			<Grid item xs={12} md={12}>
				<TextField
					variant="outlined"
					fullWidth
					type="url"
					label={<Trans>Stream key</Trans>}
					value={settings.stream_key}
					onChange={handleChange('stream_key')}
					required
				/>
			</Grid>
			<Grid item xs={12}>
				<Checkbox label="Enable Backup Server" checked={settings.enable_backup_server} onChange={handleChange('enable_backup_server')} />
			</Grid>
			{settings.enable_backup_server && (
				<Grid item xs={12} md={3}>
					<TextField
						variant="outlined"
						fullWidth
						type="url"
						label={<Trans>Protocol</Trans>}
						value={settings.protocol}
						onChange={handleChange('protocol')}
						readOnly
						disabled
					/>
				</Grid>
			)}
			{settings.enable_backup_server && (
				<Grid item xs={12} md={9}>
					<TextField
						variant="outlined"
						fullWidth
						type="url"
						label={<Trans>Backup server</Trans>}
						value={settings.backup_server}
						onChange={handleChange('backup_server')}
					/>
				</Grid>
			)}
		</Grid>
	);
}

Service.defaultProps = {
	settings: {},
	skills: {},
	metadata: {},
	streams: [],
	onChange: function (output, settings) {},
};

export { id, name, version, stream_key_link, description, image_copyright, author, category, requires, ServiceIcon as icon, Service as component };
