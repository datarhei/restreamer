import * as AVFoundation from './AVFoundation';
import * as InternalRTMP from './InternalRTMP';
import * as InternalSRT from './InternalSRT';
import * as Network from './Network';
import * as Raspicam from './Raspicam';
import * as V4L from './V4L';

class Registry {
	constructor() {
		this.services = new Map();
	}

	Register(service) {
		this.services.set(service.id, service);
	}

	Get(id) {
		const service = this.services.get(id);
		if (service) {
			return service;
		}

		return null;
	}

	IDs() {
		return Array.from(this.services.keys());
	}

	List() {
		return Array.from(this.services.values());
	}
}

const registry = new Registry();

registry.Register(Network);
registry.Register(InternalRTMP);
registry.Register(InternalSRT);
//registry.Register(InternalHLS);
registry.Register(AVFoundation);
registry.Register(Raspicam);
registry.Register(V4L);

export default registry;
