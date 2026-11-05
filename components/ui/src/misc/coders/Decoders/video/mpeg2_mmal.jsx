import React from 'react';

import Helper from '../../helper';
import { withDefaults } from '../../../../utils/props';

function init(initialState) {
	const state = {
		...initialState,
	};

	return state;
}

function createMapping(settings, stream, skills) {
	stream = Helper.InitStream(stream);
	skills = Helper.InitSkills(skills);

	const mapping = {
		global: [],
		local: ['-c:v', 'mpeg2_mmal'],
		filter: [],
	};

	return mapping;
}

function Coder(inputProps) {
	const props = withDefaults(inputProps, CoderDefaultProps);
	const settings = init(props.settings);
	const stream = Helper.InitStream(props.stream);
	const skills = Helper.InitSkills(props.skills);

	const handleChange = (newSettings) => {
		let automatic = false;
		if (!newSettings) {
			newSettings = settings;
			automatic = true;
		}

		props.onChange(newSettings, createMapping(newSettings, stream, skills), automatic);
	};

	React.useEffect(() => {
		handleChange(null);
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

	return null;
}

const CoderDefaultProps = {
	stream: {},
	settings: {},
	skills: {},
	onChange: function (settings, mapping) {},
};

const coder = 'mpeg2_mmal';
const name = 'MPEG2 (MMAL)';
const codecs = ['mpeg2'];
const type = 'video';
const hwaccel = true;

function defaults(stream, skills) {
	const settings = init({});

	return {
		settings: settings,
		mapping: createMapping(settings, stream, skills),
	};
}

export { coder, name, codecs, type, hwaccel, defaults, Coder as component };
