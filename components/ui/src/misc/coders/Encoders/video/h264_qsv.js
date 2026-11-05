import createHardwareEncoder from './hardware';

const service = createHardwareEncoder({
	coder: 'h264_qsv',
	name: 'H.264 (Intel Quick Sync)',
	codec: 'h264',
	defaultDevice: '/dev/dri/renderD128',
	defaultPreset: 'veryfast',
	presets: ['veryfast', 'faster', 'fast', 'medium', 'slow', 'slower', 'veryslow'],
	globalOptions: (settings) => [
		['-init_hw_device', `qsv=restreamer_qsv:hw,child_device=${settings.device}`],
		['-filter_hw_device', 'restreamer_qsv'],
	],
	filters: ['format=nv12', 'hwupload=extra_hw_frames=16'],
	localOptions: ['-async_depth', '1', '-forced_idr', '1'],
});

const { coder, name, codec, type, hwaccel, summarize, defaults, component } = service;
export { coder, name, codec, type, hwaccel, summarize, defaults, component };
