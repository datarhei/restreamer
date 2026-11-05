export default async function CopyToClipboard(content) {
	let success = false;

	if (!navigator.clipboard) {
		success = writeTextDeprecated(content);
	} else {
		success = await writeText(navigator.clipboard.writeText(content));
	}

	return success;
}

const writeText = (promise) => {
	return promise
		.then(() => true)
		.catch((err) => {
			console.warn(err);
			return false;
		});
};

const writeTextDeprecated = (value) => {
	// Strangely, this doesn't seem to work if the text is longer than just one row
	var element = document.createElement('textarea');
	element.value = value;
	element.style.position = 'absolute';
	element.style.width = '10px';
	element.style.height = '10px';
	element.style.top = '-1000px';
	element.style.left = '-1000px';

	document.body.appendChild(element);

	element.select();

	try {
		document.execCommand('copy');
	} catch (err) {
		console.warn(err);
		return false;
	}

	document.body.removeChild(element);

	return true;
};
