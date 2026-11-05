import React from 'react';

import { I18nProvider } from '@lingui/react';
import { i18n } from '@lingui/core';

import { messages as EN } from './locales/en/messages.js';
import { messages as DA } from './locales/da/messages.js';
import { messages as DE } from './locales/de/messages.js';
import { messages as EL } from './locales/el/messages.js';
import { messages as ES } from './locales/es/messages.js';
import { messages as FR } from './locales/fr/messages.js';
import { messages as IT } from './locales/it/messages.js';
import { messages as KO } from './locales/ko/messages.js';
import { messages as PL } from './locales/pl/messages.js';
import { messages as PT } from './locales/pt-br/messages.js';
import { messages as RU } from './locales/ru/messages.js';
import { messages as SL } from './locales/sl/messages.js';
import { messages as TR } from './locales/tr/messages.js';
import { messages as UK } from './locales/uk/messages.js';
import { messages as ZH } from './locales/zh-hans/messages.js';
import * as Storage from './utils/storage';

i18n.load({
	en: EN,
	da: DA,
	de: DE,
	el: EL,
	es: ES,
	fr: FR,
	it: IT,
	ko: KO,
	pl: PL,
	'pt-br': PT,
	ru: RU,
	sl: SL,
	tr: TR,
	uk: UK,
	'zh-hans': ZH,
});

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

i18n.activate(getLanguage('en', ['en', 'da', 'de', 'el', 'es', 'fr', 'it', 'ko', 'pl', 'pt-br', 'ru', 'sl', 'tr', 'uk', 'zh-hans']));

export default function Provider(props) {
	return <I18nProvider i18n={i18n}>{props.children}</I18nProvider>;
}
