import React from 'react';

// Adapted from https://stackoverflow.com/questions/10420352/converting-file-size-in-bytes-to-human-readable-string
export default function Filesize(props) {
	let bytes = props.bytes;
	const thresh = props.si ? 1000 : 1024;

	if (Math.abs(bytes) < thresh) {
		return bytes + ' B';
	}

	const units = props.si ? ['kB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'] : ['KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB'];
	let u = -1;
	const r = 10 ** props.digits;

	do {
		bytes /= thresh;
		++u;
	} while (Math.round(Math.abs(bytes) * r) / r >= thresh && u < units.length - 1);

	return <React.Fragment>{bytes.toFixed(props.digits) + ' ' + units[u]}</React.Fragment>;
}

Filesize.defaultProps = {
	bytes: 0,
	si: false,
	digits: 1,
};
