import * as H264QSV from './h264_qsv';
import * as HEVCQSV from './hevc_qsv';
import * as H264Vulkan from './h264_vulkan';
import * as HEVCVulkan from './hevc_vulkan';

const skills = { ffmpeg: { version_major: 8 } };

test.each([
	[H264QSV, 'h264_qsv'],
	[HEVCQSV, 'hevc_qsv'],
	[H264Vulkan, 'h264_vulkan'],
	[HEVCVulkan, 'hevc_vulkan'],
])('$coder creates a low-latency hardware mapping', (service, expectedCoder) => {
	const { settings, mapping } = service.defaults({}, skills);

	expect(mapping.local).toEqual(expect.arrayContaining(['-codec:v', expectedCoder, '-bf', '0', '-fps_mode', 'cfr']));
	expect(settings.gop).toBe('1');
	expect(mapping.filter).toEqual(['format=nv12', 'hwupload=extra_hw_frames=16']);
});

test('QSV binds to the selected DRM render node', () => {
	const { mapping } = H264QSV.defaults({}, skills);

	expect(mapping.global).toContainEqual(['-init_hw_device', 'qsv=restreamer_qsv:hw,child_device=/dev/dri/renderD128']);
	expect(mapping.local).toEqual(expect.arrayContaining(['-async_depth', '1']));
});

test('Vulkan binds to the selected Vulkan device', () => {
	const { mapping } = H264Vulkan.defaults({}, skills);

	expect(mapping.global).toContainEqual(['-init_hw_device', 'vulkan=restreamer_vk:0']);
});
