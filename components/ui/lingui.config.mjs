import { formatter } from '@lingui/format-po';

export default {
	catalogs: [
		{
			path: 'src/locales/{locale}/messages',
			include: ['src/'],
			exclude: ['**/node_modules/**'],
		},
	],
	format: formatter({ lineNumbers: false }),
	sourceLocale: 'en',
	locales: ['en', 'da', 'de', 'el', 'es', 'fr', 'it', 'ko', 'pl', 'pt-br', 'ru', 'sl', 'tr', 'uk', 'zh-hans'],
};
