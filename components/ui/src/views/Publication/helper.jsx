import * as Coders from '../../misc/coders/Encoders';
import * as M from '../../utils/metadata';

/**
 * createSourcesFromStreams creates a sources array with
 * one fake empty input and the given streams.
 * @param {array} streams
 * @returns
 */
export function createSourcesFromStreams(streams) {
	return [
		{
			type: 'network',
			inputs: [
				{
					address: '',
					options: [],
				},
			],
			streams: streams,
		},
	];
}

/**
 * createInputsOutputs creates the inputs and outputs based
 * on the sources, profiles, and outputs of the service.
 * the sources is constructed from the streams of the ingest.
 * the profiles are the encoding settings. the outputs from
 * the service only contain the output format options.
 *
 * @param {*} sources
 * @param {*} profiles encoding settings
 * @param {*} outputs service outputs (format options and address)
 * @returns
 */
export function createInputsOutputs(sources, profiles, outputs, requireVideo = true) {
	const [global, inpts, outpts] = M.createInputsOutputs(sources, profiles, requireVideo);

	if (inpts.length === 0 || outpts.length === 0) {
		return [global, [], []];
	}

	const out = [];

	// iterate through all service outputs and copy its values
	// to a new object. if encoding is enabled, it will be done
	// for each output. this is not optimal. the "tee" output
	// should be considered.
	// https://ffmpeg.org/ffmpeg-formats.html#tee
	// https://stackoverflow.com/questions/67609510/ffmpeg-tee-muxer-failing-with-error-tag-avc1-incompatible-with-output-codec-id
	// and also consider the fifo and onfail option
	// https://ffmpeg.org/ffmpeg-formats.html#fifo
	for (let i = 0; i < outputs.length; i++) {
		let o = {
			address: outputs[i].address,
			options: [...outputs[i].options],
		};

		o.options = [...outpts[0].options, ...o.options];

		out.push(o);
	}

	return [global, inpts, out];
}
/**
 * validateRequirements validates the requirements object the each
 * service provides. Missing fields will be added.
 *
 * @param {*} requires requirement object
 * @returns validated requirement object
 */
export function validateRequirements(requires) {
	if (!requires) {
		requires = {};
	}

	requires = {
		protocols: [],
		formats: [],
		devices: [],
		codecs: {},
		...requires,
	};

	requires.codecs = {
		audio: [],
		video: [],
		...requires.codecs,
	};

	return requires;
}

/**
 *
 * @param {*} requires
 * @param {*} skills
 * @returns {boolean}
 */
export function checkServiceRequirements(requires, skills) {
	if (!skills) {
		return false;
	}

	requires = validateRequirements(requires);

	if (requires.protocols.length !== 0) {
		let hasOne = false;
		for (let protocol of requires.protocols) {
			if (skills.protocols.output.includes(protocol) === true) {
				hasOne = true;
				break;
			}
		}

		if (hasOne === false) {
			return false;
		}
	}

	if (requires.formats.length !== 0) {
		let hasOne = false;
		for (let format of requires.formats) {
			if (skills.formats.muxers.includes(format) === true) {
				hasOne = true;
				break;
			}
		}

		if (hasOne === false) {
			return false;
		}
	}

	if (requires.devices.length !== 0) {
		let hasOne = false;
		for (let device of requires.devices) {
			if (device in skills.sinks) {
				hasOne = true;
				break;
			}
		}

		if (hasOne === false) {
			return false;
		}
	}

	if (requires.codecs.audio.length !== 0) {
		let hasOne = false;
		for (let codec of requires.codecs.audio) {
			if (codec in skills.codecs.audio) {
				hasOne = true;
				break;
			}
		}

		if (hasOne === false) {
			return false;
		}
	}

	if (requires.codecs.video.length !== 0) {
		let hasOne = false;
		for (let codec of requires.codecs.video) {
			if (codec in skills.codecs.video) {
				hasOne = true;
				break;
			}
		}

		if (hasOne === false) {
			return false;
		}
	}

	return true;
}

/**
 *
 * @param {*} requires
 * @param {*} skills
 * @returns
 */
export function conflateServiceSkills(requires, skills) {
	if (!skills) {
		return null;
	}

	requires = validateRequirements(requires);

	const serviceSkills = {
		ffmpeg: {
			...skills.ffmpeg,
		},
		protocols: [],
		formats: [],
		devices: {},
		codecs: {
			audio: ['copy'],
			video: ['copy'],
		},
	};

	for (let protocol of requires.protocols) {
		if (skills.protocols.output.includes(protocol) === true) {
			serviceSkills.protocols.push(protocol);
		}
	}

	for (let format of requires.formats) {
		if (skills.formats.muxers.includes(format) === true) {
			serviceSkills.formats.push(format);
		}
	}

	for (let device of requires.devices) {
		if (device in skills.sinks) {
			serviceSkills.devices[device] = [...skills.sinks[device]];
		}
	}

	for (let codec of requires.codecs.audio) {
		if (codec in skills.codecs.audio) {
			serviceSkills.codecs.audio.push(codec);
		}
	}

	for (let codec of requires.codecs.video) {
		if (codec in skills.codecs.video) {
			serviceSkills.codecs.video.push(codec);
		}
	}

	return serviceSkills;
}

/**
 * Preselects an encoder to a profile.
 * @param {*} profile A profile
 * @param {*} type Either 'audio' or 'video'
 * @param {*} streams List of available streams
 * @param {*} codecs List of target codecs
 * @param {*} skills FFmpeg skills
 * @returns {boolean} Whether the provided profile is valid
 */
export function preselectProfile(profile, type, streams, codecs, skills) {
	const encoders = skills.encoders[type];

	/**
	 * Checks if the given profile makes sense, i.e. matches to the available
	 * streams and codecs.
	 * @param {*} profile A profile
	 * @param {*} type Either 'audio' or 'video'
	 * @param {*} streams List of available streams
	 * @param {*} codecs List of target codecs
	 * @returns {boolean} Whether the provided profile is valid
	 */
	const isPlausible = (profile, type, streams, codecs) => {
		if (profile.stream < 0) {
			return false;
		}

		if (profile.stream >= streams.length) {
			return false;
		}

		if (streams[profile.stream].type !== type) {
			return false;
		}

		if (!codecs.includes(streams[profile.stream].codec)) {
			if (profile.encoder.coder === 'copy') {
				return false;
			}
		} else {
			if (profile.encoder.coder === 'copy') {
				return true;
			}
		}

		let coder = null;
		if (type === 'audio') {
			coder = Coders.Audio.Get(profile.encoder.coder);
		} else if (type === 'video') {
			coder = Coders.Video.Get(profile.encoder.coder);
		}

		if (coder === null) {
			return false;
		}

		if (!codecs.includes(coder.codec)) {
			return false;
		}

		return true;
	};

	if (isPlausible(profile, type, streams, codecs)) {
		return profile;
	}

	profile.source = 0;
	profile.stream = -1;
	profile.encoder.coder = 'none';

	for (let i = 0; i < streams.length; i++) {
		if (streams[i].type !== type) {
			continue;
		}

		// Select the first stream with matching type.
		profile.stream = i;

		if (!codecs.includes(streams[i].codec)) {
			// The codec doesn't match. Select the first available coder for one of the target codecs.
			for (let codec of codecs) {
				if (codec === 'copy') {
					continue;
				}

				let coder = null;

				if (type === 'audio') {
					coder = Coders.Audio.GetCoderForCodec(codec, encoders);
				} else if (type === 'video') {
					coder = Coders.Video.GetCoderForCodec(codec, encoders);
				}

				if (coder === null) {
					profile.encoder.coder = 'none';
				} else {
					const defaults = coder.defaults(streams[i], skills);
					profile.encoder.coder = coder.coder;
					profile.encoder.settings = defaults.settings;
					profile.encoder.mapping = defaults.mapping;
					break;
				}
			}
		} else {
			// The codec matches. Select the copy coder.
			profile.encoder.coder = 'copy';

			let coder = type === 'audio' ? Coders.Audio.Get('copy') : Coders.Video.Get('copy');

			const defaults = coder.defaults(streams[i], skills);

			profile.encoder.settings = defaults.settings;
			profile.encoder.mapping = defaults.mapping;
		}

		break;
	}

	return profile;
}
