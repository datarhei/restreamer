import React from 'react';

import { I18nProvider } from '@lingui/react';
import { i18n } from '@lingui/core';

import { messages as EN } from './locales/en/messages.js';
import * as Storage from './utils/storage';

const catalogLoaders = {
	da: () => import('./locales/da/messages.js'),
	de: () => import('./locales/de/messages.js'),
	el: () => import('./locales/el/messages.js'),
	es: () => import('./locales/es/messages.js'),
	fr: () => import('./locales/fr/messages.js'),
	it: () => import('./locales/it/messages.js'),
	ko: () => import('./locales/ko/messages.js'),
	pl: () => import('./locales/pl/messages.js'),
	'pt-br': () => import('./locales/pt-br/messages.js'),
	ru: () => import('./locales/ru/messages.js'),
	sl: () => import('./locales/sl/messages.js'),
	tr: () => import('./locales/tr/messages.js'),
	uk: () => import('./locales/uk/messages.js'),
	'zh-hans': () => import('./locales/zh-hans/messages.js'),
};

const supportedLanguages = ['en', ...Object.keys(catalogLoaders)];

i18n.loadAndActivate({ locale: 'en', messages: EN });

const aliases = {
	pt: 'pt-br',
	'zh-cn': 'zh-hans',
};

const getAlias = (lang) => {
	if (lang in aliases) {
		return aliases[lang];
	}

	return lang;
};

const getLanguage = (defaultLanguage, supportedLanguages) => {
	let lang = getAlias(Storage.Get('language'));
	if (supportedLanguages.indexOf(lang) === -1) {
		lang = getAlias(getBrowserLanguage(defaultLanguage));

		if (supportedLanguages.indexOf(lang) === -1) {
			lang = defaultLanguage;
		}
	}

	Storage.Set('language', lang);

	return lang;
};

const getBrowserLanguage = (defaultLanguage) => {
	let lang = window.navigator.language;

	const match = lang.match(/^[a-z]+(-[a-z]+)?/i);
	if (!match) {
		return defaultLanguage;
	}

	return match[0].toLowerCase();
};

let activationSequence = 0;

export async function activateLanguage(language) {
	const sequence = ++activationSequence;
	let locale = supportedLanguages.includes(language) ? language : 'en';
	let messages = EN;

	if (locale !== 'en') {
		try {
			({ messages } = await catalogLoaders[locale]());
		} catch {
			locale = 'en';
			messages = EN;
		}
	}

	if (sequence !== activationSequence) {
		return i18n.locale;
	}

	i18n.loadAndActivate({ locale, messages });
	return locale;
}

const initialLanguage = getLanguage('en', supportedLanguages);
const initialActivation = initialLanguage === 'en' ? Promise.resolve('en') : activateLanguage(initialLanguage);

export default function Provider(props) {
	const [ready, setReady] = React.useState(initialLanguage === 'en');

	React.useEffect(() => {
		if (ready) {
			return undefined;
		}

		let mounted = true;
		initialActivation.then((locale) => {
			if (mounted) {
				Storage.Set('language', locale);
				setReady(true);
			}
		});

		return () => {
			mounted = false;
		};
	}, [ready]);

	if (!ready) {
		return null;
	}

	return <I18nProvider i18n={i18n}>{props.children}</I18nProvider>;
}
