import * as AudioCopy from './audio/copy';
import * as AudioNone from './audio/none';
import * as AAC from './audio/aac';
import * as AACAudioToolbox from './audio/aac_audiotoolbox';
import * as Libopus from './audio/opus_libopus';
import * as Libvorbis from './audio/vorbis_libvorbis';
import * as MP3 from './audio/mp3_libmp3lame';
import * as Opus from './audio/opus';
import * as Vorbis from './audio/vorbis';

import * as X264 from './video/h264_libx264';
import * as X265 from './video/hevc_libx265';
import * as H264VideoToolbox from './video/h264_videotoolbox';
import * as H264NVENC from './video/h264_nvenc';
import * as H264OMX from './video/h264_omx';
import * as H264V4L2M2M from './video/h264_v4l2m2m';
import * as H264VAAPI from './video/h264_vaapi';
import * as HEVCVAAPI from './video/hevc_vaapi';
import * as HEVCVideoToolbox from './video/hevc_videotoolbox';
import * as VP9VAAPI from './video/vp9_vaapi';
import * as VideoCopy from './video/copy';
import * as VideoNone from './video/none';
import * as VideoRaw from './video/rawvideo';
import * as VP9 from './video/vp9_libvpx';
import * as AV1Rav1e from './video/av1_librav1e';

class Registry {
	constructor(type) {
		this.type = type;
		this.services = new Map();
	}

	Register(service) {
		if (service.type !== this.type) {
			return;
		}

		this.services.set(service.coder, service);
	}

	Get(coder) {
		const service = this.services.get(coder);
		if (service) {
			return service;
		}

		return null;
	}

	// Get the first coder for a codec that is in a
	// list of available coders.
	GetCoderForCodec(codec, availableCoders) {
		for (let coder of this.services.values()) {
			if (coder.codec !== codec) {
				continue;
			}

			if (!availableCoders.includes(coder.coder)) {
				continue;
			}

			return coder;
		}

		return null;
	}

	// Get a list of coders for a codec that is in a list of
	// availabled coders. The option for hwAcceleration can be
	// 'any', 'no', or 'yes'.
	GetCodersForCodec(codec, availableCoders, hwAcceleration) {
		const coders = [];

		for (let coder of this.services.values()) {
			if (coder.codec !== codec) {
				continue;
			}

			if (!availableCoders.includes(coder.coder)) {
				continue;
			}

			if (hwAcceleration === 'any') {
				coders.push(coder);
				continue;
			}

			if (hwAcceleration === 'no' && coder.hwaccel === false) {
				coders.push(coder);
				continue;
			}

			if (hwAcceleration === 'yes' && coder.hwaccel === true) {
				coders.push(coder);
				continue;
			}
		}

		return coders;
	}

	Coders() {
		return Array.from(this.services.keys());
	}

	List() {
		return Array.from(this.services.values());
	}
}

const audioRegistry = new Registry('audio');

audioRegistry.Register(AudioCopy);
audioRegistry.Register(AudioNone);
audioRegistry.Register(AAC);
audioRegistry.Register(AACAudioToolbox);
audioRegistry.Register(MP3);
audioRegistry.Register(Opus);
audioRegistry.Register(Libopus);
audioRegistry.Register(Vorbis);
audioRegistry.Register(Libvorbis);

const videoRegistry = new Registry('video');

videoRegistry.Register(VideoCopy);
videoRegistry.Register(VideoNone);
videoRegistry.Register(X264);
videoRegistry.Register(H264VideoToolbox);
videoRegistry.Register(H264NVENC);
videoRegistry.Register(H264OMX);
videoRegistry.Register(H264V4L2M2M);
videoRegistry.Register(H264VAAPI);
videoRegistry.Register(X265);
videoRegistry.Register(HEVCVAAPI);
videoRegistry.Register(HEVCVideoToolbox);
videoRegistry.Register(VP9VAAPI);
videoRegistry.Register(VP9);
videoRegistry.Register(AV1Rav1e);
videoRegistry.Register(VideoRaw);

export { audioRegistry as Audio, videoRegistry as Video };
