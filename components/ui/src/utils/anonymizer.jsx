import urlparser from 'url-parse';

const anonymize_url = (url) => {
	let u = urlparser(url, true);

	if (u.hostname !== 'localhost') {
		return `${u.protocol}//[anonymized]`;
	}

	if (u.auth.length !== 0) {
		u.username = '***';
		u.password = '***';
	}

	if (u.protocol === 'rtmp:') {
		if ('token' in u.query) {
			u.query.token = '***';
		}
	} else if (u.protocol === 'srt') {
		if ('streamid' in u.query) {
			u.query.streamid = '***';
		}

		if ('passphrase' in u.query) {
			u.query.passphrase = '***';
		}
	}

	return u.toString();
};

const anonymize = (text) => {
	const regex = /(?:([a-z0-9\\]+):)?\/[A-Za-z0-9-._~!$&'()*+,;=:@?/{}%\\]*/gm;

	return text.replaceAll(regex, (match, scheme) => {
		if (scheme) {
			match = match.replace(scheme, scheme.replaceAll('\\', ''));
			return anonymize_url(match);
		}

		const pathElm = match.split('/').filter((p) => p.length !== 0);
		if (pathElm.length < 2) {
			return match;
		}

		return `/[anonymized]/${pathElm.pop()}`;
	});
};

export { anonymize, anonymize_url };
