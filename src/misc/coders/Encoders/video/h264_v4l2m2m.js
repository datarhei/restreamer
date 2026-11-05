import React from 'react';

import Grid from '@mui/material/Grid';
import Typography from '@mui/material/Typography';

import { Trans } from '@lingui/macro';

import BoxText from '../../../BoxText';
import TextField from '../../../TextField';
import Video from '../../settings/Video';
import Helper from '../../helper';

function init(initialState) {
	const state = {
		bitrate: '4096',
		fps: '25',
		gop: '4',
		profile: 'auto',
		force_key_frames: 'expr:gte(t,n_forced*1)',
		num_capture_buffers: '256',
		num_output_buffers: '512',
		fps_mode: 'cfr',

		...initialState,
	};

	return state;
}

// https://forums.raspberrypi.com/viewtopic.php?t=294161
// ffmpeg -y -nostdin -f v4l2 -threads auto -input_format yuyv422 -fflags +genpts -flags +global_header \
// -i /dev/video0 -s 1280x720 -r 25 \
// -vcodec h264_v4l2m2m \
// -num_output_buffers 32 -num_capture_buffers 16 \
// -keyint_min 25 -force_key_frames "expr:gte(t,n_forced*1)" \
// -g 50 -b:v 6M -pix_fmt nv12 \
// -f mp4 -f segment -segment_format_options movflags=+faststart -segment_time 1 -reset_timestamps 1 -segment_time 1 -segment_time_delta 1 -segment_format mp4 -muxdelay 0 -muxpreload 0 /home/pi/Videos/%d-output.mp4

/**
v4l2-ctl --list-ctrls-menu -d 11

Codec Controls

             video_bitrate_mode 0x009909ce (menu)   : min=0 max=1 default=0 value=0 (Variable Bitrate) flags=update
				0: Variable Bitrate
				1: Constant Bitrate
                  video_bitrate 0x009909cf (int)    : min=25000 max=25000000 step=25000 default=10000000 value=10000000
           sequence_header_mode 0x009909d8 (menu)   : min=0 max=1 default=1 value=1 (Joined With 1st Frame)
				0: Separate Buffer
				1: Joined With 1st Frame
         repeat_sequence_header 0x009909e2 (bool)   : default=0 value=0
                force_key_frame 0x009909e5 (button) : value=0 flags=write-only, execute-on-write
          h264_minimum_qp_value 0x00990a61 (int)    : min=0 max=51 step=1 default=20 value=20
          h264_maximum_qp_value 0x00990a62 (int)    : min=0 max=51 step=1 default=51 value=51
            h264_i_frame_period 0x00990a66 (int)    : min=0 max=2147483647 step=1 default=60 value=60
                     h264_level 0x00990a67 (menu)   : min=0 max=13 default=11 value=11 (4)
				0: 1
				1: 1b
				2: 1.1
				3: 1.2
				4: 1.3
				5: 2
				6: 2.1
				7: 2.2
				8: 3
				9: 3.1
				10: 3.2
				11: 4
				12: 4.1
				13: 4.2
                   h264_profile 0x00990a6b (menu)   : min=0 max=4 default=4 value=4 (High)
				0: Baseline
				1: Constrained Baseline
				2: Main
				4: High
 */

function createMapping(settings, stream, skills) {
	stream = Helper.InitStream(stream);
	skills = Helper.InitSkills(skills);

	const local = [
		'-codec:v',
		'h264_v4l2m2m',
		'-num_capture_buffers',
		`${settings.num_capture_buffers}`,
		'-num_output_buffers',
		`${settings.num_output_buffers}`,
		'-b:v',
		`${settings.bitrate}k`,
		'-maxrate',
		`${settings.bitrate}k`,
		'-bufsize',
		`${settings.bitrate}k`,
		'-r',
		`${settings.fps}`,
		'-pix_fmt',
		'yuv420p',
		'-sc_threshold',
		'0',
		'-copyts',
		'-avoid_negative_ts',
		'disabled',
		'-max_muxing_queue_size',
		'2048',
	];

	if (settings.force_key_frames !== '') {
		local.push('-force_key_frames', `${settings.force_key_frames}`);
	}

	if (settings.gop !== 'auto') {
		local.push('-g', `${Math.round(parseInt(settings.fps) * parseInt(settings.gop)).toFixed(0)}`);
		local.push('-keyint_min', `${parseInt(settings.fps)}`);
	}

	if (skills.ffmpeg.version_major >= 5) {
		local.push('-fps_mode', `${settings.fps_mode}`);
	}

	if (settings.profile !== 'auto') {
		local.push('-profile:v', `${settings.profile}`);
	}

	const mapping = {
		global: [],
		local: local,
		filter: [],
	};

	return mapping;
}

function Coder(props) {
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

	const update = (what) => (event) => {
		const newSettings = {
			...settings,
			[what]: event.target.value,
		};

		handleChange(newSettings);
	};

	React.useEffect(() => {
		handleChange(null);
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

	return (
		<Grid container spacing={2}>
			<Grid item xs={12}>
				<BoxText color="danger">
					<Trans>V4L2_M2M is experimental.</Trans>
					<br />
					<Trans>We recommend OpenMAX IL for Raspberry PI (3/4) with a 32-bit operating system.</Trans>
				</BoxText>
			</Grid>
			<Grid item xs={12}>
				<Video.Bitrate value={settings.bitrate} onChange={update('bitrate')} allowCustom />
			</Grid>
			<Grid item xs={12}>
				<Video.Framerate value={settings.fps} onChange={update('fps')} allowCustom />
			</Grid>
			<Grid item xs={12}>
				<Video.GOP value={settings.gop} onChange={update('gop')} allowAuto allowCustom />
				<Typography variant="caption">
					<Trans>To stabilize the system, increase the HLS segment length for the keyframe interval by 2-3 * (Processing and Control).</Trans>
				</Typography>
			</Grid>
			{skills.ffmpeg.version_major >= 5 && (
				<Grid item xs={12}>
					<Video.FpsMode value={settings.fps_mode} onChange={update('fps_mode')} />
				</Grid>
			)}
			<Grid item xs={12}>
				<TextField label={<Trans>Force key frames</Trans>} type="text" value={settings.force_key_frames} onChange={update('force_key_frames')} />
			</Grid>
			<Grid item xs={6}>
				<TextField label={<Trans>Capture buffer</Trans>} type="number" value={settings.num_capture_buffers} onChange={update('num_capture_buffers')} />
			</Grid>
			<Grid item xs={6}>
				<TextField label={<Trans>Output buffer</Trans>} type="number" value={settings.num_output_buffers} onChange={update('num_output_buffers')} />
			</Grid>
		</Grid>
	);
}

Coder.defaultProps = {
	stream: {},
	settings: {},
	skills: {},
	onChange: function (settings, mapping) {},
};

const coder = 'h264_v4l2m2m';
const name = 'H.264 (V4L2 Memory to Memory)';
const codec = 'h264';
const type = 'video';
const hwaccel = true;

function summarize(settings) {
	return `${name}, ${settings.bitrate} kbit/s, ${settings.fps} FPS, Profile: ${settings.profile}`;
}

function defaults(stream, skills) {
	const settings = init({});

	return {
		settings: settings,
		mapping: createMapping(settings, stream, skills),
	};
}

export { coder, name, codec, type, hwaccel, summarize, defaults, Coder as component };
