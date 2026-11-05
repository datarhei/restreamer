import React from 'react';

import { Trans } from '@lingui/macro';
import Grid from '@mui/material/Grid';
import Link from '@mui/material/Link';
import TextField from '@mui/material/TextField';

import Logo from './logos/akamai.svg';

import Checkbox from '../../../misc/Checkbox';

const id = 'akamai';
const name = 'Akamai';
const version = '1.0';
const stream_key_link = '';
const description = (
	<Trans>
		Transmit the main source to the Akamai (MSL) Media Services Live. More details about the MSL Encoder settings can be found on{' '}
		<Link
			color="secondary"
			target="_blank"
			href="https://learn.akamai.com/en-us/webhelp/media-services-live/media-services-live-encoder-compatibility-testing-and-qualification-guide-v4.0/GUID-620AF7FB-F1D7-4CC5-AC43-F4D56DDC02DB.html"
		>
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
	return <img src={Logo} alt="Akamai Logo" {...props} />;
}

function init(settings) {
	const initSettings = {
		protocol: 'rtmp://',
		stream_url: '',
		backup_url: '',
		address_extensions: {},
		enable_backup_stream: false,
		...settings,
	};

	initSettings.address_extensions = {
		flashver: 'FMLE/3.020(compatible;20FMSc/1.0)',
		live: 'true',
		pubUser: '',
		pubPasswd: '',
		playpath: '',
		...initSettings.address_extensions,
	};

	return initSettings;
}

function Service(props) {
	const settings = init(props.settings);

	const handleChange = (what) => (event) => {
		const value = event.target.value;

		if (what in settings.address_extensions) {
			settings.address_extensions[what] = value;
		} else {
			if (what === 'enable_backup_stream') {
				settings[what] = !settings[what];
			} else {
				settings[what] = value;
			}
		}

		const output = createOutput(settings);

		props.onChange([output], settings);
	};

	const createOutput = (settings) => {
		const address_extensions = [];

		for (let key in settings.address_extensions) {
			if (settings.address_extensions[key].length === 0) {
				continue;
			}
			address_extensions.push(key + '=' + settings.address_extensions[key]);
		}

		const options = ['-f', 'flv'];

		const output = {
			address: settings.protocol + settings.stream_url + ' ' + address_extensions.join(' '),
			options: options,
		};

		return output;
	};

	return (
		<Grid container spacing={2}>
			<Grid item xs={12} md={6}>
				<TextField
					variant="outlined"
					fullWidth
					type="url"
					label={<Trans>Login</Trans>}
					value={settings.address_extensions.pubUser}
					onChange={handleChange('pubUser')}
				/>
			</Grid>
			<Grid item xs={12} md={6}>
				<TextField
					variant="outlined"
					fullWidth
					type="url"
					label={<Trans>Password</Trans>}
					value={settings.address_extensions.pubPasswd}
					onChange={handleChange('pubPasswd')}
				/>
			</Grid>
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
			<Grid item xs={12} md={9}>
				<TextField
					variant="outlined"
					fullWidth
					type="url"
					label={<Trans>Stream URL</Trans>}
					placeholder="p.ab12345.i.akamaientrypoint.net/EntryPoint"
					value={settings.stream_url}
					onChange={handleChange('stream_url')}
				/>
			</Grid>
			<Grid item xs={12} md={12}>
				<TextField
					variant="outlined"
					fullWidth
					type="url"
					label={<Trans>Stream names</Trans>}
					placeholder="abc_%i@12345"
					value={settings.address_extensions.playpath}
					onChange={handleChange('playpath')}
				/>
			</Grid>
			<Grid item xs={12}>
				<Checkbox label="Enable backup URL" checked={settings.enable_backup_stream} onChange={handleChange('enable_backup_stream')} />
			</Grid>
			{settings.enable_backup_stream && (
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
			{settings.enable_backup_stream && (
				<Grid item xs={12} md={9}>
					<TextField
						variant="outlined"
						fullWidth
						type="url"
						label={<Trans>Backup URL</Trans>}
						placeholder="b.ab12345.i.akamaientrypoint.net/EntryPoint"
						value={settings.backup_url}
						onChange={handleChange('backup_url')}
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
