import * as AudioDefault from './audio/default';

import * as NVDEC from './video/nvdec';
import * as H264MMAL from './video/h264_mmal';
import * as H264CUVID from './video/h264_cuvid';
import * as HEVCCUVID from './video/hevc_cuvid';
import * as MJPEGCUVID from './video/mjpeg_cuvid';
import * as MPEG1CUVID from './video/mpeg1_cuvid';
import * as MPEG2CUVID from './video/mpeg2_cuvid';
import * as MPEG2MMAL from './video/mpeg2_mmal';
import * as MPEG4CUVID from './video/mpeg4_cuvid';
import * as MPEG4MMAL from './video/mpeg4_mmal';
import * as VC1CUVID from './video/vc1_cuvid';
import * as VC1MMAL from './video/vc1_mmal';
import * as VideoDefault from './video/default';
import * as VideoToolbox from './video/videotoolbox';
import * as VP8CUVID from './video/vp8_cuvid';
import * as VP9CUVID from './video/vp9_cuvid';
import * as AV1CUVID from './video/av1_cuvid';

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
			if (!coder.codecs.includes(codec)) {
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
			// An empty codecs list is a catch-all
			if (coder.codecs.length !== 0) {
				if (!coder.codecs.includes(codec)) {
					continue;
				}
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

audioRegistry.Register(AudioDefault);

const videoRegistry = new Registry('video');

videoRegistry.Register(VideoDefault);
videoRegistry.Register(VideoToolbox);
videoRegistry.Register(AV1CUVID);
videoRegistry.Register(NVDEC);
videoRegistry.Register(H264MMAL);
videoRegistry.Register(H264CUVID);
videoRegistry.Register(HEVCCUVID);
videoRegistry.Register(MJPEGCUVID);
videoRegistry.Register(MPEG1CUVID);
videoRegistry.Register(MPEG2CUVID);
videoRegistry.Register(MPEG2MMAL);
videoRegistry.Register(MPEG4CUVID);
videoRegistry.Register(MPEG4MMAL);
videoRegistry.Register(VC1CUVID);
videoRegistry.Register(VC1MMAL);
videoRegistry.Register(VP8CUVID);
videoRegistry.Register(VP9CUVID);

export { audioRegistry as Audio, videoRegistry as Video };
