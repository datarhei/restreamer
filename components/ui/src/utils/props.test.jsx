import { withDefaults } from './props';

test('withDefaults applies only missing and undefined values', () => {
	const props = { disabled: false, empty: '', missing: undefined, nullable: null, zero: 0 };
	const defaults = {
		disabled: true,
		empty: 'default',
		missing: 'default',
		newValue: 'default',
		nullable: 'default',
		zero: 42,
	};

	expect(withDefaults(props, defaults)).toEqual({
		disabled: false,
		empty: '',
		missing: 'default',
		newValue: 'default',
		nullable: null,
		zero: 0,
	});
	expect(props.missing).toBeUndefined();
	expect(defaults.newValue).toBe('default');
});
