// Audio Filter
import * as AResample from './audio/Resample';
import * as Pan from './audio/Pan';
import * as Volume from './audio/Volume';
import * as Loudnorm from './audio/Loudnorm';

// Video Filter
import * as Bwdif from './video/Bwdif';
import * as Framerate from './video/Framerate';
import * as Scale from './video/Scale';
import * as Transpose from './video/Transpose';
import * as HFlip from './video/HFlip';
import * as VFlip from './video/VFlip';

// Register filters type: audio/video
class Registry {
	constructor(type) {
		this.type = type;
		this.services = new Map();
	}

	Register(service) {
		if (service.type !== this.type) {
			return;
		}

		this.services.set(service.filter, service);
	}

	Get(filter) {
		const service = this.services.get(filter);
		if (service) {
			return service;
		}

		return null;
	}

	Filters() {
		return Array.from(this.services.keys());
	}

	List() {
		return Array.from(this.services.values());
	}
}

// Audio Filters
const audioRegistry = new Registry('audio');
audioRegistry.Register(AResample);
audioRegistry.Register(Pan);
audioRegistry.Register(Volume);
audioRegistry.Register(Loudnorm);

// Video Filters
const videoRegistry = new Registry('video');
videoRegistry.Register(Bwdif);
videoRegistry.Register(Framerate);
videoRegistry.Register(Scale);
videoRegistry.Register(Transpose);
videoRegistry.Register(HFlip);
videoRegistry.Register(VFlip);

// Export registrys for ../SelectFilters.js
export { audioRegistry as Audio, videoRegistry as Video };
