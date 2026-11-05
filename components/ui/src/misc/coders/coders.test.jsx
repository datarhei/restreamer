import React from 'react';
import { render, act } from '../../utils/testing';
import '@testing-library/jest-dom';

import * as Decoders from './Decoders';
import * as Encoders from './Encoders';

const audiodecoders = Decoders.Audio.List();
const videodecoders = Decoders.Video.List();
const audioencoders = Encoders.Audio.List();
const videoencoders = Encoders.Video.List();

const testfunc = async (coder) => {
	const defaults = coder.defaults();

	let $settings = {};
	let $mapping = {};

	const handleChange = (settings, mapping) => {
		$settings = settings;
		$mapping = mapping;
	};

	await act(async () => {
		render(<coder.component onChange={handleChange} />);
	});

	expect($settings).toStrictEqual(defaults.settings);
	expect($mapping).toStrictEqual(defaults.mapping);
};

test.each(audiodecoders)('coder:decoder:$type $coder', testfunc);
test.each(videodecoders)('coder:decoder:$type $coder', testfunc);
test.each(audioencoders)('coder:encoder:$type $coder', testfunc);
test.each(videoencoders)('coder:encoder:$type $coder', testfunc);
