import React from 'react';

import { Trans } from '@lingui/macro';
import Grid from '@mui/material/Grid';
import Typography from '@mui/material/Typography';

// Import all filters (audio/video)
import * as Filters from './filters';

// Import all encoders (audio/video)
import * as Encoders from './coders/Encoders';

export default function FilterSelect(props) {
	const profile = props.profile;

	// handleFilterChange
	// what: Filter name
	// settings (component settings):  {Key: Value}
	// mapping (FFmpeg -af/-vf args): ['String', ...]
	const handleFilterSettingsChange = (what) => (settings, graph, automatic) => {
		const filter = profile.filter;

		// Store mapping/settings per component
		filter.settings[what] = {
			settings: settings,
			graph: graph,
		};

		// Get the order of the filters
		let filterOrder = [];
		if (props.type === 'video') {
			filterOrder = Filters.Video.Filters();
		} else {
			filterOrder = Filters.Audio.Filters();
		}

		// Create the filter graph in the order as the filters are registered
		const graphs = [];
		for (let f of filterOrder) {
			if (!(f in filter.settings)) {
				continue;
			}

			if (filter.settings[f].graph.length !== 0) {
				graphs.push(filter.settings[f].graph);
			}
		}

		filter.graph = graphs.join(',');

		props.onChange(filter, automatic);
	};

	// Set filterRegistry by type
	let filterRegistry = null;
	if (props.type === 'video') {
		filterRegistry = Filters.Video;
	} else if (props.type === 'audio') {
		filterRegistry = Filters.Audio;
	} else {
		return null;
	}

	// Checks the state of hwaccel (gpu encoding)
	let encoderRegistry = null;
	let hwaccel = false;
	if (props.type === 'video') {
		encoderRegistry = Encoders.Video;
		for (let encoder of encoderRegistry.List()) {
			if (encoder.codec === props.profile.encoder.coder && encoder.hwaccel) {
				hwaccel = true;
			}
		}
	}

	// Creates filter components
	let filterSettings = [];
	if (!hwaccel) {
		for (let c of filterRegistry.List()) {
			// Checks FFmpeg skills (filter is available)
			if (props.availableFilters.includes(c.filter)) {
				const Settings = c.component;

				if (!(c.filter in profile.filter.settings)) {
					profile.filter.settings[c.filter] = c.defaults();
				} else {
					profile.filter.settings[c.filter] = {
						...c.defaults(),
						...profile.filter.settings[c.filter],
					};
				}

				filterSettings.push(
					<Settings key={c.filter} settings={profile.filter.settings[c.filter].settings} onChange={handleFilterSettingsChange(c.filter)} />,
				);
			}
		}
	}

	// No suitable filter found
	if (filterSettings === null && !hwaccel) {
		return (
			<Grid container spacing={2}>
				<Grid item xs={12}>
					<Typography>
						<Trans>No suitable filter found.</Trans>
					</Typography>
				</Grid>
			</Grid>
		);

		// hwaccel requires further settings
	} else if (hwaccel) {
		return false;
	}

	return (
		<Grid container spacing={2}>
			<Grid item xs={12}>
				<Typography>
					<Trans>Select your filter settings (optional):</Trans>
				</Typography>
			</Grid>
			{filterSettings}
		</Grid>
	);
}

FilterSelect.defaultProps = {
	type: '',
	profile: {},
	availableFilters: [],
	onChange: function (filter, automatic) {},
};
