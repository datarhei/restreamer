import React from 'react';

export default function Duration(props) {
	const fullSeconds = parseInt(Math.floor(props.seconds));
	const s = fullSeconds % 60;
	const m = parseInt(fullSeconds / 60) % 60;
	const h = parseInt(fullSeconds / (60 * 60)) % 24;
	const d = parseInt(fullSeconds / (60 * 60 * 24));

	let duration = '.' + ((props.seconds - fullSeconds) * 100).toFixed(0);

	if (s < 10) {
		duration = ':0' + s + duration;
	} else {
		duration = ':' + s + duration;
	}

	if (m < 10) {
		duration = ':0' + m + duration;
	} else {
		duration = ':' + m + duration;
	}

	if (h < 10) {
		duration = '0' + h + duration;
	} else {
		duration = '' + h + duration;
	}

	if (d !== 0) {
		duration = d + ':' + duration;
	}

	return <React.Fragment>{duration}</React.Fragment>;
}

Duration.defaultProps = {
	seconds: 0,
};
