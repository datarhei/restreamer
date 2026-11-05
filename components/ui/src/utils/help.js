import { i18n } from '@lingui/core';

const topics = {
	'edit-general': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/edit-livestream/general',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/streameinstellungen/allgemeine-streameinstellungen',
	},
	'edit-control': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/edit-livestream/processing-and-control',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/streameinstellungen/verarbeitung-und-steuerung',
	},
	'edit-meta': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/edit-livestream/meta-information',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/streameinstellungen/metainformationen',
	},
	'edit-license': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/edit-livestream/license',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/streameinstellungen/lizenz',
	},
	login: {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/login',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/dashboard',
	},
	main: {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/main-screen',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/hauptbildschirm',
	},
	'player-embed': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/video-player-settings',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/videoplayer-einstellungen',
	},
	'player-colors': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/video-player-settings',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/videoplayer-einstellungen',
	},
	'player-logo': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/video-player-settings',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/videoplayer-einstellungen',
	},
	'player-statistic': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/video-player-settings',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/videoplayer-einstellungen',
	},
	'player-playback': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/video-player-settings',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/videoplayer-einstellungen',
	},
	'player-poster': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/video-player-settings',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/videoplayer-einstellungen',
	},
	'playersite-general': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/publication-website/general-settings',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/publication-webseite/allgemeine-einstellungen',
	},
	'playersite-template': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/publication-website/template',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/publication-webseite/template',
	},
	'playersite-design': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/publication-website/design',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/publication-webseite/design',
	},
	'playersite-notes': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/publication-website/notes',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/publication-webseite/notizen',
	},
	'playersite-code_injection': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/publication-website/code-injection',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/publication-webseite/code-injection',
	},
	'process-details': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/process-details',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/prozess-report',
	},
	'process-report': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/process-report',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/prozess-details',
	},
	publication: {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/publications',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/publication-services',
	},
	'publication-add': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/publications',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/publication-services',
	},
	'publication-general': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/publications',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/publication-services',
	},
	'publication-process': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/publications',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/publication-services',
	},
	'publication-encoding': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/publications',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/publication-services',
	},
	'settings-general': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/system-settings/general-system-settings',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/restreamer-einstellungen/allgemeine-systemeinstellungen',
	},
	'settings-update-link': {
		en: 'https://docs.datarhei.com/restreamer/installing/update',
		de: 'https://docs.datarhei.com/restreamer/v/de/installation/aktualisieren',
	},
	'settings-service': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/system-settings/service',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/restreamer-einstellungen/service',
	},
	'settings-network': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/system-settings/network',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/restreamer-einstellungen/netzwerk',
	},
	'settings-auth': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/system-settings/authorization',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/restreamer-einstellungen/autorisierung',
	},
	'settings-playback': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/system-settings/playback',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/restreamer-einstellungen/wiedergabe',
	},
	'settings-storage': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/system-settings/disk-space',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/restreamer-einstellungen/speicherplatz',
	},
	'settings-rtmp': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/system-settings/rtmp',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/restreamer-einstellungen/rtmp',
	},
	'settings-srt': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/system-settings/srt',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/restreamer-einstellungen/srt',
	},
	'settings-logging': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/system-settings/logging',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/restreamer-einstellungen/protokollierung-and-fehlersuche',
	},
	'wizard-video-setup': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/wizard',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/assistent',
	},
	'wizard-video-settings': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/wizard',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/assistent',
	},
	'wizard-video-result': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/wizard',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/assistent',
	},
	'wizard-audio-settings': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/wizard',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/assistent',
	},
	'wizard-audio-result': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/wizard',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/assistent',
	},
	'wizard-license': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/wizard',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/assistent',
	},
	'wizard-error': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/wizard',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/assistent',
	},
	'wizard-abort': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/manual/wizard',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/user-guides/assistent',
	},
	'encoder-h264_omx': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/troubleshooting/encoding-compatibility-list#h264_omx',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/fehlersuche/kompatibilitaetsliste-fuer-encoder#h264_omx',
	},
	'encoder-h264_v4l2m2m': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/troubleshooting/encoding-compatibility-list#h264_v4l2m2m',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/fehlersuche/kompatibilitaetsliste-fuer-encoder#h264_v4l2m2m',
	},
	'encoder-h264_nvenc': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/troubleshooting/encoding-compatibility-list#h264_nvenc',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/fehlersuche/kompatibilitaetsliste-fuer-encoder#h264_nvenc',
	},
	'encoder-h264_vaapi': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/troubleshooting/encoding-compatibility-list#h264_vaapi',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/fehlersuche/kompatibilitaetsliste-fuer-encoder#h264_vaapi',
	},
	'encoder-h264_videotoolbox': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/troubleshooting/encoding-compatibility-list#h264_videotoolbox',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/fehlersuche/kompatibilitaetsliste-fuer-encoder#h264_videotoolbox',
	},
	'encoder-libx264': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/troubleshooting/encoding-compatibility-list#libx264',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/fehlersuche/kompatibilitaetsliste-fuer-encoder#libx264',
	},
	'encoder-libx265': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/troubleshooting/encoding-compatibility-list#libx265',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/fehlersuche/kompatibilitaetsliste-fuer-encoder#libx265',
	},
	'encoder-libvpx-vp9': {
		en: 'https://docs.datarhei.com/restreamer/knowledge-base/troubleshooting/encoding-compatibility-list#libvpx-vp9',
		de: 'https://docs.datarhei.com/restreamer/v/de/wissensdatenbank/fehlersuche/kompatibilitaetsliste-fuer-encoder#libvpx-vp9',
	},
};

function getTopicURL(topic, locale) {
	if (!(topic in topics)) {
		console.warn(`help topic "${topic}" not found`);
		// If topic doesn't exist, return default URL
		return 'https://docs.datarhei.com/restreamer';
	}

	if (!(locale in topics[topic])) {
		// If locale doesn't exist, return default locale (en) URL
		// This requires that "en" exists for all topics.
		return topics[topic].en;
	}

	return topics[topic][locale];
}

export default function Help(topic) {
	const url = getTopicURL(topic, i18n.locale);

	window.open(url);
}
