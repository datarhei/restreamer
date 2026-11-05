export function withDefaults(props, defaults) {
	const resolved = { ...props };

	for (const [name, value] of Object.entries(defaults)) {
		if (resolved[name] === undefined) {
			resolved[name] = value;
		}
	}

	return resolved;
}
