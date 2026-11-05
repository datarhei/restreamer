import React from 'react';
import { render, fireEvent } from '../../../utils/testing';
import '@testing-library/jest-dom';

import * as Network from './Network';

const $skills_ffmpeg6 = {
	ffmpeg: {
		version: '6.1.1',
		version_major: 6,
		version_minor: 1,
	},
	formats: {
		demuxers: ['rtsp'],
	},
	protocols: {
		input: ['http', 'https', 'rtmp', 'rtmps', 'srt'],
	},
	codecs: {
		audio: {},
		video: {
			av1: ['librav1e'],
			hevc: ['hevc'],
			vp9: ['libvpx-vp9'],
		},
	},
};

const $skills_ffmpeg5 = {
	ffmpeg: {
		version: '5.1.2',
		version_major: 5,
		version_minor: 1,
	},
	formats: {
		demuxers: ['rtsp'],
	},
	protocols: {
		input: ['http', 'https', 'rtmp', 'rtmps', 'srt'],
	},
};

const $skills_ffmpeg4 = {
	ffmpeg: {
		version: '4.4.1',
		version_major: 4,
		version_minor: 4,
	},
	formats: {
		demuxers: ['rtsp'],
	},
	protocols: {
		input: ['http', 'https', 'rtmp', 'rtmps', 'srt'],
	},
};

const $config = {
	channelid: 'external',
	rtmp: {
		enabled: true,
		app: '/live',
		token: 'foobar',
	},
	srt: {
		enabled: true,
		token: 'foobar',
		passphrase: 'bazfoobazfoo',
	},
};

test('source:network pull', async () => {
	let $settings = {
		mode: 'pull',
	};
	const handleChange = (settings) => {
		$settings = settings;
	};

	const Source = Network.component;
	let { getByLabelText, queryByText, rerender } = render(<Source onChange={handleChange} knownDevices={[]} />);

	expect(queryByText(`This protocol is unknown or not supported by the available FFmpeg binary.`)).toBe(null);

	const input = getByLabelText('Address');
	fireEvent.change(input, { target: { value: 'rtsp://127.0.0.1/live/stream' } });

	expect($settings.mode).toBe('pull');
	expect($settings.address).toBe('rtsp://127.0.0.1/live/stream');

	rerender(<Source settings={$settings} onChange={handleChange} knownDevices={[]} />);

	expect(queryByText(`This protocol is unknown or not supported by the available FFmpeg binary.`)).toBeInTheDocument();

	rerender(<Source settings={$settings} skills={$skills_ffmpeg4} onChange={handleChange} knownDevices={[]} />);

	expect(queryByText(`This protocol is unknown or not supported by the available FFmpeg binary.`)).toBe(null);

	rerender(<Source settings={$settings} skills={$skills_ffmpeg5} onChange={handleChange} knownDevices={[]} />);

	expect(queryByText(`This protocol is unknown or not supported by the available FFmpeg binary.`)).toBe(null);

	rerender(<Source settings={$settings} skills={$skills_ffmpeg6} onChange={handleChange} knownDevices={[]} />);

	expect(queryByText(`This protocol is unknown or not supported by the available FFmpeg binary.`)).toBe(null);
});

const pullmatrix = {
	settings: {
		mode: 'pull',
		address: '',
		username: 'admin',
		password: 'foobar',
		rtsp: {
			udp: false,
			stimeout: 5000000,
		},
		http: {
			readNative: true,
			forceFramerate: true,
			framerate: 25,
			userAgent: 'foobaz/1',
			referer: 'http://example.com',
			http_proxy: '',
		},
		general: {
			analyzeduration: 5000000,
			analyzeduration_rtmp: 3000000,
			analyzeduration_http: 20000000,
			probesize: 5000000,
			max_probe_packets: 2500,
			fflags: ['genpts'],
			thread_queue_size: 512,
			copyts: false,
			start_at_zero: false,
			use_wallclock_as_timestamps: false,
			avoid_negative_ts: 'auto',
		},
	},
	tests: [],
};

pullmatrix.tests = [
	{
		name: 'RTSP',
		settings: { ...pullmatrix.settings, address: 'rtsp://127.0.0.1/live/stream' },
		skills: $skills_ffmpeg4,
		input: {
			address: 'rtsp://admin:foobar@127.0.0.1/live/stream',
			options: [
				'-fflags',
				'+genpts',
				'-thread_queue_size',
				512,
				'-probesize',
				5000000,
				'-analyzeduration',
				5000000,
				'-stimeout',
				5000000,
				'-rtsp_transport',
				'tcp',
			],
		},
	},
	{
		name: 'RTMP',
		settings: { ...pullmatrix.settings, address: 'rtmp://127.0.0.1/live/stream' },
		skills: $skills_ffmpeg4,
		input: {
			address: 'rtmp://admin:foobar@127.0.0.1/live/stream',
			options: ['-fflags', '+genpts', '-thread_queue_size', 512, '-probesize', 5000000, '-analyzeduration', 3000000],
		},
	},
	{
		name: 'HTTP',
		settings: { ...pullmatrix.settings, address: 'http://127.0.0.1/live/stream.m3u8' },
		skills: $skills_ffmpeg4,
		input: {
			address: 'http://admin:foobar@127.0.0.1/live/stream.m3u8',
			options: [
				'-fflags',
				'+genpts',
				'-thread_queue_size',
				512,
				'-probesize',
				5000000,
				'-analyzeduration',
				20000000,
				'-re',
				'-r',
				25,
				'-user_agent',
				'foobaz/1',
				'-referer',
				'http://example.com',
			],
		},
	},
	{
		name: 'SRT',
		settings: { ...pullmatrix.settings, address: 'srt://127.0.0.1?mode=caller&streamid=foobar' },
		skills: $skills_ffmpeg4,
		input: {
			address: 'srt://127.0.0.1?mode=caller&streamid=foobar',
			options: ['-fflags', '+genpts', '-thread_queue_size', 512, '-probesize', 5000000, '-analyzeduration', 5000000],
		},
	},
	{
		name: 'RTSP',
		settings: { ...pullmatrix.settings, address: 'rtsp://127.0.0.1/live/stream' },
		skills: $skills_ffmpeg5,
		input: {
			address: 'rtsp://admin:foobar@127.0.0.1/live/stream',
			options: [
				'-fflags',
				'+genpts',
				'-thread_queue_size',
				512,
				'-probesize',
				5000000,
				'-analyzeduration',
				5000000,
				'-timeout',
				5000000,
				'-rtsp_transport',
				'tcp',
			],
		},
	},
	{
		name: 'RTMP',
		settings: { ...pullmatrix.settings, address: 'rtmp://127.0.0.1/live/stream' },
		skills: $skills_ffmpeg5,
		input: {
			address: 'rtmp://admin:foobar@127.0.0.1/live/stream',
			options: ['-fflags', '+genpts', '-thread_queue_size', 512, '-probesize', 5000000, '-analyzeduration', 3000000],
		},
	},
	{
		name: 'HTTP',
		settings: { ...pullmatrix.settings, address: 'http://127.0.0.1/live/stream.m3u8' },
		skills: $skills_ffmpeg5,
		input: {
			address: 'http://admin:foobar@127.0.0.1/live/stream.m3u8',
			options: [
				'-fflags',
				'+genpts',
				'-thread_queue_size',
				512,
				'-probesize',
				5000000,
				'-analyzeduration',
				20000000,
				'-re',
				'-r',
				25,
				'-user_agent',
				'foobaz/1',
				'-referer',
				'http://example.com',
			],
		},
	},
	{
		name: 'SRT',
		settings: { ...pullmatrix.settings, address: 'srt://127.0.0.1?mode=caller&streamid=foobar' },
		skills: $skills_ffmpeg5,
		input: {
			address: 'srt://127.0.0.1?mode=caller&streamid=foobar',
			options: ['-fflags', '+genpts', '-thread_queue_size', 512, '-probesize', 5000000, '-analyzeduration', 5000000],
		},
	},
	{
		name: 'RTSP',
		settings: { ...pullmatrix.settings, address: 'rtsp://127.0.0.1/live/stream' },
		skills: $skills_ffmpeg6,
		input: {
			address: 'rtsp://admin:foobar@127.0.0.1/live/stream',
			options: [
				'-fflags',
				'+genpts',
				'-thread_queue_size',
				512,
				'-probesize',
				5000000,
				'-analyzeduration',
				5000000,
				'-timeout',
				5000000,
				'-rtsp_transport',
				'tcp',
			],
		},
	},
	{
		name: 'RTMP',
		settings: { ...pullmatrix.settings, address: 'rtmp://127.0.0.1/live/stream' },
		skills: $skills_ffmpeg6,
		input: {
			address: 'rtmp://admin:foobar@127.0.0.1/live/stream',
			options: [
				'-fflags',
				'+genpts',
				'-thread_queue_size',
				512,
				'-probesize',
				5000000,
				'-analyzeduration',
				3000000,
				'-rtmp_enhanced_codecs',
				'hvc1,av01,vp09',
			],
		},
	},
	{
		name: 'HTTP',
		settings: { ...pullmatrix.settings, address: 'http://127.0.0.1/live/stream.m3u8' },
		skills: $skills_ffmpeg6,
		input: {
			address: 'http://admin:foobar@127.0.0.1/live/stream.m3u8',
			options: [
				'-fflags',
				'+genpts',
				'-thread_queue_size',
				512,
				'-probesize',
				5000000,
				'-analyzeduration',
				20000000,
				'-re',
				'-r',
				25,
				'-user_agent',
				'foobaz/1',
				'-referer',
				'http://example.com',
			],
		},
	},
	{
		name: 'SRT',
		settings: { ...pullmatrix.settings, address: 'srt://127.0.0.1?mode=caller&streamid=foobar' },
		skills: $skills_ffmpeg6,
		input: {
			address: 'srt://127.0.0.1?mode=caller&streamid=foobar',
			options: ['-fflags', '+genpts', '-thread_queue_size', 512, '-probesize', 5000000, '-analyzeduration', 5000000],
		},
	},
];

test.each(pullmatrix.tests)('source:network pull $name input with ffmpeg $skills.ffmpeg.version', async (data) => {
	let $inputs = [];
	const handleProbe = (_, inputs) => {
		$inputs = inputs;
	};

	const Source = Network.component;

	let { getByText, getByRole } = render(<Source settings={data.settings} skills={data.skills} onProbe={handleProbe} knownDevices={[]} />);

	expect(getByText('Probe')).toBeInTheDocument();

	const button = getByRole('button', { name: 'Probe' });
	fireEvent.click(button, { bubbles: true });

	expect($inputs.length).toBe(1);
	expect($inputs[0]).toStrictEqual(data.input);
});

test('source:network push', async () => {
	let $settings = {
		mode: 'push',
		push: {
			type: 'rtmp',
			name: 'external',
		},
	};
	const handleChange = (settings) => {
		$settings = settings;
	};

	const Source = Network.component;
	let { queryByText, rerender } = render(<Source settings={$settings} onChange={handleChange} knownDevices={[]} />);

	expect($settings.mode).toBe('push');

	expect(queryByText(`The available FFmpeg binary doesn't support any of the required protocols.`)).toBeInTheDocument();

	rerender(<Source settings={$settings} skills={$skills_ffmpeg4} onChange={handleChange} knownDevices={[]} />);

	expect(queryByText(`The available FFmpeg binary doesn't support any of the required protocols.`)).toBe(null);

	rerender(<Source settings={$settings} skills={$skills_ffmpeg5} onChange={handleChange} knownDevices={[]} />);

	expect(queryByText(`The available FFmpeg binary doesn't support any of the required protocols.`)).toBe(null);

	rerender(<Source settings={$settings} skills={$skills_ffmpeg6} onChange={handleChange} knownDevices={[]} />);

	expect(queryByText(`The available FFmpeg binary doesn't support any of the required protocols.`)).toBe(null);
});

test('source:network push RTMP', async () => {
	let $settings = {
		mode: 'push',
		push: {
			type: 'rtmp',
			name: 'external',
		},
	};
	const handleChange = (settings) => {
		$settings = settings;
	};

	const Source = Network.component;
	let { getByText, queryByText, rerender } = render(<Source settings={$settings} skills={$skills_ffmpeg5} onChange={handleChange} knownDevices={[]} />);

	expect($settings.mode).toBe('push');
	expect($settings.push.type).toBe('rtmp');

	expect(queryByText(`Enable RTMP server ...`)).toBeInTheDocument();

	rerender(<Source settings={$settings} config={$config} skills={$skills_ffmpeg5} onChange={handleChange} knownDevices={[]} />);

	expect(getByText('Probe')).toBeInTheDocument();
});

test('source:network push SRT', async () => {
	let $settings = {
		mode: 'push',
		push: {
			type: 'srt',
			name: 'external',
		},
	};
	const handleChange = (settings) => {
		$settings = settings;
	};

	const Source = Network.component;
	let { getByText, queryByText, rerender } = render(<Source settings={$settings} skills={$skills_ffmpeg5} onChange={handleChange} knownDevices={[]} />);

	expect($settings.mode).toBe('push');
	expect($settings.push.type).toBe('srt');

	expect(queryByText(`Enable SRT server ...`)).toBeInTheDocument();

	rerender(<Source settings={$settings} config={$config} skills={$skills_ffmpeg5} onChange={handleChange} knownDevices={[]} />);

	expect(getByText('Probe')).toBeInTheDocument();
});

const pushmatrix = {
	settings: {
		mode: 'push',
		push: {
			type: '',
			name: 'external',
		},
	},
	tests: [],
};

pushmatrix.tests = [
	{
		name: 'RTMP',
		settings: { ...pushmatrix.settings, push: { ...pushmatrix.push, type: 'rtmp' } },
		skills: $skills_ffmpeg4,
		config: $config,
		input: {
			address: '{rtmp,name=external.stream}',
			options: ['-fflags', '+genpts', '-thread_queue_size', 512, '-probesize', 5000000, '-analyzeduration', 3000000],
		},
	},
	{
		name: 'SRT',
		settings: { ...pushmatrix.settings, push: { ...pushmatrix.push, type: 'srt' } },
		skills: $skills_ffmpeg4,
		config: $config,
		input: {
			address: '{srt,name=external.stream,mode=request}',
			options: ['-fflags', '+genpts', '-thread_queue_size', 512, '-probesize', 5000000, '-analyzeduration', 5000000],
		},
	},
	{
		name: 'RTMP',
		settings: { ...pushmatrix.settings, push: { ...pushmatrix.push, type: 'rtmp' } },
		skills: $skills_ffmpeg5,
		config: $config,
		input: {
			address: '{rtmp,name=external.stream}',
			options: ['-fflags', '+genpts', '-thread_queue_size', 512, '-probesize', 5000000, '-analyzeduration', 3000000],
		},
	},
	{
		name: 'SRT',
		settings: { ...pushmatrix.settings, push: { ...pushmatrix.push, type: 'srt' } },
		skills: $skills_ffmpeg5,
		config: $config,
		input: {
			address: '{srt,name=external.stream,mode=request}',
			options: ['-fflags', '+genpts', '-thread_queue_size', 512, '-probesize', 5000000, '-analyzeduration', 5000000],
		},
	},
	{
		name: 'RTMP',
		settings: { ...pushmatrix.settings, push: { ...pushmatrix.push, type: 'rtmp' } },
		skills: $skills_ffmpeg6,
		config: $config,
		input: {
			address: '{rtmp,name=external.stream}',
			options: [
				'-fflags',
				'+genpts',
				'-thread_queue_size',
				512,
				'-probesize',
				5000000,
				'-analyzeduration',
				3000000,
				'-rtmp_enhanced_codecs',
				'hvc1,av01,vp09',
			],
		},
	},
	{
		name: 'SRT',
		settings: { ...pushmatrix.settings, push: { ...pushmatrix.push, type: 'srt' } },
		skills: $skills_ffmpeg6,
		config: $config,
		input: {
			address: '{srt,name=external.stream,mode=request}',
			options: ['-fflags', '+genpts', '-thread_queue_size', 512, '-probesize', 5000000, '-analyzeduration', 5000000],
		},
	},
];

test.each(pushmatrix.tests)('source:network push $name input with ffmpeg $skills.ffmpeg.version', async (data) => {
	let $inputs = [];
	const handleProbe = (_, inputs) => {
		$inputs = inputs;
	};

	const Source = Network.component;
	let { getByText, getByRole } = render(
		<Source settings={data.settings} config={data.config} skills={data.skills} onProbe={handleProbe} knownDevices={[]} />,
	);

	expect(getByText('Probe')).toBeInTheDocument();

	const button = getByRole('button', { name: 'Probe' });
	fireEvent.click(button, { bubbles: true });

	expect($inputs.length).toBe(1);
	expect($inputs[0]).toStrictEqual(data.input);
});
