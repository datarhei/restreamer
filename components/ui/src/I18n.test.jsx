import { i18n } from '@lingui/core';

import { activateLanguage } from './I18n';

afterEach(async () => {
	await activateLanguage('en');
});

test('loads and activates a translated catalog on demand', async () => {
	expect(await activateLanguage('de')).toBe('de');
	expect(i18n.locale).toBe('de');
});

test('falls back to English for an unsupported language', async () => {
	expect(await activateLanguage('not-a-language')).toBe('en');
	expect(i18n.locale).toBe('en');
});
