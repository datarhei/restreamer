import React from 'react';

import { Trans } from '@lingui/macro';
import Grid from '@mui/material/Grid';
import MenuItem from '@mui/material/MenuItem';
import TextField from '@mui/material/TextField';
import Typography from '@mui/material/Typography';

import Checkbox from '../Checkbox';
import Select from '../Select';

function init(settings) {
	const initSettings = {
		lhls: false,
		segmentDuration: 2,
		listSize: 6,
		cleanup: true,
		version: 3,
		storage: 'memfs',
		master_playlist: true,
		...settings,
	};

	return initSettings;
}

export default function Control(props) {
	const settings = init(props.settings);

	// Set the defaults
	React.useEffect(() => {
		props.onChange(settings, true);
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

	const handleChange = (what) => (event) => {
		const value = event.target.value;

		if (['lhls', 'cleanup', 'master_playlist'].includes(what)) {
			settings[what] = !settings[what];
		} else {
			settings[what] = value;
		}

		props.onChange(settings, false);
	};

	return (
		<Grid container spacing={2}>
			{/*
				lhls problems:
				- naming of the manifests and media files does not work
				- segmentation of audio and video creates a playback error
				- hls inputs requires encoding
				- have to wait for the final integration in hls.js

				<Grid item xs={12}>
					<Checkbox label={<Trans>Low latency Mode (LHLS/CMAF)</Trans>} checked={settings.lhls} onChange={handleChange('lhls')} />
					<Typography variant="caption">
						<Trans>Experimental function. Older browsers and clients may not be supported.</Trans>
					</Typography>
				</Grid>
			*/}
			<Grid item xs={12}>
				<Select label={<Trans>Storage</Trans>} value={settings.storage} onChange={handleChange('storage')}>
					<MenuItem value="memfs">
						<Trans>In-memory</Trans>
					</MenuItem>
					<MenuItem value="diskfs">
						<Trans>Disk</Trans>
					</MenuItem>
				</Select>
				<Typography variant="caption">
					<Trans>Where to store the HLS playlist and segments. In-Memory is recommended.</Trans>
				</Typography>
			</Grid>
			<Grid item xs={12}>
				<Select label={<Trans>EXT-X-VERSION</Trans>} value={settings.version} onChange={handleChange('version')}>
					<MenuItem value={3}>3</MenuItem>
					<MenuItem value={6}>
						<Trans>6 (+ guaranteed to start with a Key frame)</Trans>
					</MenuItem>
					<MenuItem value={7}>
						<Trans>7 (+ fragmented MP4 format)</Trans>
					</MenuItem>
				</Select>
				<Typography variant="caption">
					<Trans>Playlist version (M3U8). Version 3 has the best browser/client compatibility.</Trans>
				</Typography>
			</Grid>
			<Grid item xs={12} md={6}>
				<TextField
					variant="outlined"
					fullWidth
					label={<Trans>Segment length (seconds)</Trans>}
					value={settings.segmentDuration}
					onChange={handleChange('segmentDuration')}
				/>
				<Typography variant="caption">
					<Trans>Segment will be cut on the following keyframe after this time has passed. 2 is recommended.</Trans>
				</Typography>
			</Grid>
			<Grid item xs={12} md={6}>
				<TextField
					variant="outlined"
					fullWidth
					label={<Trans>List size (segments)</Trans>}
					value={settings.listSize}
					onChange={handleChange('listSize')}
				/>
				<Typography variant="caption">
					<Trans>The maximum number of playlist segments. 0 will contain all the segments. 6 is recommended.</Trans>
				</Typography>
			</Grid>
			<Grid item xs={12}>
				<Checkbox
					label={<Trans>Master playlist (increases browser/client compatibility)</Trans>}
					checked={settings.master_playlist}
					onChange={handleChange('master_playlist')}
				/>
			</Grid>
			<Grid item xs={12}>
				<Checkbox label={<Trans>Automatic cleanup of all media data</Trans>} checked={settings.cleanup} onChange={handleChange('cleanup')} />
			</Grid>
		</Grid>
	);
}

Control.defaulProps = {
	settings: {},
	onChange: function (settings, automatic) {},
};
