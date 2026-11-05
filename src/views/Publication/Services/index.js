import * as Akamai from './Akamai';
import * as Azure from './Azure';
import * as Brightcove from './Brightcove';
import * as CDN77 from './CDN77';
import * as Core from './Core';
import * as DaCast from './DaCast';
import * as Dailymotion from './Dailymotion';
import * as DASH from './DASH';
import * as DLive from './DLive';
import * as Facebook from './Facebook';
import * as Framebuffer from './Framebuffer';
import * as HLS from './HLS';
import * as Icecast from './Icecast';
import * as Image2 from './Image2';
import * as Instagram from './Instagram';
import * as Kick from './Kick';
import * as Linkedin from './Linkedin';
import * as Livepush from './Livepush';
import * as Livespotting from './Livespotting';
import * as MediaNetwork from './MediaNetwork';
import * as MPEGTS from './MPEGTS';
import * as NimoTv from './NimoTv';
import * as Owncast from './Owncast';
import * as PeerTube from './PeerTube';
import * as PicartoTv from './PicartoTv';
import * as Red5 from './Red5';
import * as Restream from './Restream';
import * as RTMP from './RTMP';
import * as RTSP from './RTSP';
import * as Rumble from './Rumble';
import * as SRT from './SRT';
import * as Trovo from './Trovo';
import * as Telegram from './Telegram';
import * as Twitch from './Twitch';
import * as Twitter from './Twitter';
import * as UDP from './UDP';
import * as Vimeo from './Vimeo';
import * as WOWZA from './WOWZA';
import * as Wettercom from './Wettercom';
import * as Youtube from './Youtube';

class Registry {
	constructor() {
		this.services = new Map();
	}

	Register(service) {
		if (service.id.match(/[^0-9a-z]/)) {
			console.warn(`the service.id "${service.id}" is invalid. only [0-9a-z] is allowed.`);
			return;
		}

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

// The order the services are registered is relevant, i.e. on
// the "add service" screen they appear in this order.
registry.Register(Facebook);
registry.Register(Youtube);
registry.Register(Twitter);
registry.Register(Twitch);
registry.Register(Kick);
registry.Register(Instagram);
registry.Register(Vimeo);
registry.Register(Restream);
registry.Register(Telegram);
registry.Register(Linkedin);
registry.Register(DLive);
registry.Register(Trovo);
registry.Register(PicartoTv);
registry.Register(NimoTv);
registry.Register(PeerTube);
registry.Register(MediaNetwork);
registry.Register(Livespotting);
registry.Register(Wettercom);
registry.Register(Azure);
registry.Register(Brightcove);
registry.Register(Akamai);
registry.Register(DaCast);
registry.Register(Dailymotion);
registry.Register(Livepush);
registry.Register(Rumble);
registry.Register(CDN77);
registry.Register(Core);
registry.Register(Owncast);
registry.Register(WOWZA);
registry.Register(Red5);
registry.Register(Icecast);
registry.Register(Image2);
registry.Register(RTSP);
registry.Register(RTMP);
registry.Register(HLS);
registry.Register(DASH);
registry.Register(SRT);
registry.Register(UDP);
registry.Register(MPEGTS);
registry.Register(Framebuffer);

export default registry;
