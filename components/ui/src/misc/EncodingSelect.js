import React from 'react';

import { useLingui } from '@lingui/react';
import { Trans, t } from '@lingui/macro';
import Grid from '@mui/material/Grid';
import MenuItem from '@mui/material/MenuItem';
import Typography from '@mui/material/Typography';
import Link from '@mui/material/Link';

import * as Encoders from './coders/Encoders';
import * as Decoders from './coders/Decoders';
import Select from './Select';
import H from '../utils/help';

export default function EncodingSelect(props) {
	const { i18n } = useLingui();

	const profile = props.profile;
	let availableEncoders = [];
	let availableDecoders = [];

	if (props.type === 'video') {
		availableEncoders = props.skills.encoders.video;
		availableDecoders = props.skills.decoders.video;
	} else if (props.type === 'audio') {
		availableEncoders = props.skills.encoders.audio;
	}

	const handleDecoderChange = (event) => {
		const decoder = profile.decoder;
		const stream = props.streams[profile.stream];
		decoder.coder = event.target.value;

		// If the coder changes, use the coder's default settings
		let c = null;
		if (props.type === 'audio') {
			c = Decoders.Audio.Get(decoder.coder);
		} else if (props.type === 'video') {
			c = Decoders.Video.Get(decoder.coder);
		}

		if (c !== null) {
			const defaults = c.defaults(stream, props.skills);
			decoder.settings = defaults.settings;
			decoder.mapping = defaults.mapping;
		}

		props.onChange(profile.encoder, decoder, false);
	};

	const handleDecoderSettingsChange = (settings, mapping, automatic) => {
		const decoder = profile.decoder;

		decoder.settings = settings;
		decoder.mapping = mapping;

		props.onChange(profile.encoder, decoder, automatic);
	};

	const handleEncoderChange = (event) => {
		const encoder = profile.encoder;
		const stream = props.streams[profile.stream];
		encoder.coder = event.target.value;

		// If the coder changes, use the coder's default settings
		let c = null;
		if (props.type === 'audio') {
			c = Encoders.Audio.Get(encoder.coder);
		} else if (props.type === 'video') {
			c = Encoders.Video.Get(encoder.coder);
		}

		if (c !== null) {
			const defaults = c.defaults(stream, props.skills);
			encoder.settings = defaults.settings;
			encoder.mapping = defaults.mapping;
		}

		props.onChange(encoder, profile.decoder, false);
	};

	const handleEncoderSettingsChange = (settings, mapping, automatic) => {
		const encoder = profile.encoder;

		encoder.settings = settings;
		encoder.mapping = mapping;

		props.onChange(encoder, profile.decoder, automatic);
	};

	const handleEncoderHelp = (topic) => (event) => {
		event.preventDefault();
		H('encoder-' + topic);
	};

	let stream = null;

	if (profile.stream >= 0) {
		if (profile.stream < props.streams.length) {
			stream = props.streams[profile.stream];
		}
	}

	if (stream === null) {
		return null;
	}

	if (stream.type !== props.type) {
		return null;
	}

	let allowCopy = props.codecs.includes(stream.codec);
	let encoderRegistry = null;
	let decoderRegistry = null;

	if (props.type === 'video') {
		encoderRegistry = Encoders.Video;
		decoderRegistry = Decoders.Video;
	} else if (props.type === 'audio') {
		encoderRegistry = Encoders.Audio;
		decoderRegistry = Decoders.Audio;
	} else {
		return null;
	}

	let encoderSettings = null;
	let encoderSettingsHelp = null;

	let coder = encoderRegistry.Get(profile.encoder.coder);
	if (coder !== null && availableEncoders.includes(coder.coder)) {
		const Settings = coder.component;

		encoderSettings = <Settings stream={stream} settings={profile.encoder.settings} skills={props.skills} onChange={handleEncoderSettingsChange} />;

		if (props.type === 'video' && !['copy', 'none', 'rawvideo'].includes(coder.coder)) {
			encoderSettingsHelp = handleEncoderHelp(coder.coder);
		}
	}

	let encoderList = [];

	for (let c of encoderRegistry.List()) {
		// Does ffmpeg support the coder?
		if (!availableEncoders.includes(c.coder)) {
			continue;
		}

		// Is the encoder in the list of codec we allow as target?
		if (!props.codecs.includes(c.codec)) {
			continue;
		}

		if (c.coder === 'copy') {
			if (allowCopy === true) {
				encoderList.push(
					<MenuItem value={c.coder} key={c.coder}>
						{c.name}
					</MenuItem>,
				);
			}
		} else {
			encoderList.push(
				<MenuItem value={c.coder} key={c.coder}>
					{c.name}
				</MenuItem>,
			);
		}
	}

	if (encoderSettings === null || encoderList.length === 0) {
		return (
			<Grid container spacing={2}>
				<Grid item xs={12}>
					<Typography>
						<Trans>No suitable encoder found.</Trans>
					</Typography>
				</Grid>
			</Grid>
		);
	}

	let decoderSettings = null;
	let decoderList = [];

	if (coder.coder !== 'copy' && coder.coder !== 'none') {
		let c = decoderRegistry.Get(profile.decoder.coder);
		if (c !== null && availableDecoders.includes(c.coder)) {
			const Settings = c.component;

			decoderSettings = <Settings stream={stream} settings={profile.decoder.settings} skills={props.skills} onChange={handleDecoderSettingsChange} />;
		}

		// List all decoders for the codec of the stream
		for (let c of decoderRegistry.GetCodersForCodec(stream.codec, availableDecoders, 'any')) {
			decoderList.push(
				<MenuItem value={c.coder} key={c.coder}>
					{c.name}
				</MenuItem>,
			);
		}
	}

	// TODO: in case there's no decoder for a codec it should be mentioned.

	return (
		<Grid container spacing={2}>
			<Grid item xs={12}>
				<Typography>
					<Trans>Select your encoding setting:</Trans>
				</Typography>
			</Grid>
			<Grid item xs={12}>
				<Select label={<Trans>Codec</Trans>} value={profile.encoder.coder} onChange={handleEncoderChange}>
					<MenuItem value="none" disabled>
						{i18n._(t`Choose codec ...`)}
					</MenuItem>
					{encoderList}
				</Select>
			</Grid>
			{decoderList.length >= 2 && (
				<React.Fragment>
					<Grid item xs={12}>
						<Select label={<Trans>Decoder</Trans>} value={profile.decoder.coder} onChange={handleDecoderChange}>
							{decoderList}
						</Select>
					</Grid>
					<Grid item xs={12}>
						{decoderSettings}
					</Grid>
				</React.Fragment>
			)}
			<Grid item xs={12}>
				{encoderSettings}
			</Grid>
			{encoderSettingsHelp !== null && (
				<Grid item xs={12}>
					<Trans>
						<Link color="secondary" href="#" onClick={encoderSettingsHelp}>
							Compatibility list
						</Link>
					</Trans>
				</Grid>
			)}
		</Grid>
	);
}

EncodingSelect.defaultProps = {
	type: '',
	streams: [],
	profile: {},
	codecs: [],
	skills: {},
	onChange: function (encoder, decoder, automatic) {},
};
