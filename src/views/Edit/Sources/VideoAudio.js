import React from 'react';

import { Trans } from '@lingui/macro';
import Icon from '@mui/icons-material/Movie';

// This is a pseudo audio source for selecting the audio streams from the video source

const initSettings = (initialSettings) => {
	if (!initialSettings) {
		initialSettings = {};
	}

	const settings = {
		...initialSettings,
	};

	return settings;
};

const createInputs = (settings) => {
	const input = {
		address: '',
		options: [],
	};

	return [input];
};

function Source(props) {
	return null;
}

Source.defaultProps = {
	knownDevices: [],
	settings: {},
	onChange: function (settings) {},
	onProbe: function (settings, inputs) {},
};

function SourceIcon(props) {
	return <Icon {...props} />;
}

const id = 'videoaudio';
const name = <Trans>Video source</Trans>;
const capabilities = ['audio'];
const ffversion = '^4.1.0 || ^5.0.0 || ^6.1.0';

const func = {
	initSettings,
	createInputs,
};

export { id, name, capabilities, ffversion, SourceIcon as icon, Source as component, func };
