const isAvailable = () => {
	// https://developer.mozilla.org/en-US/docs/Web/API/Web_Storage_API/Using_the_Web_Storage_API
	let storage = null;

	try {
		storage = window.localStorage;
		const x = '__storage_test__';
		storage.setItem(x, x);
		storage.removeItem(x);
		return true;
	} catch (e) {
		return (
			e instanceof DOMException &&
			// everything except Firefox
			(e.code === 22 ||
				// Firefox
				e.code === 1014 ||
				// test name field too, because code might not be present
				// everything except Firefox
				e.name === 'QuotaExceededError' ||
				// Firefox
				e.name === 'NS_ERROR_DOM_QUOTA_REACHED') &&
			// acknowledge QuotaExceededError only if there's something already stored
			storage &&
			storage.length !== 0
		);
	}
};

const Set = (key, value) => {
	if (isAvailable() === false) {
		return;
	}

	window.localStorage.setItem('@@restreamer-ui@@' + key, value);
};

const Get = (key) => {
	if (isAvailable() === false) {
		return null;
	}

	return window.localStorage.getItem('@@restreamer-ui@@' + key);
};

const Remove = (key) => {
	if (isAvailable() === false) {
		return;
	}

	window.localStorage.removeItem('@@restreamer-ui@@' + key);
};

export { Set, Get, Remove };
