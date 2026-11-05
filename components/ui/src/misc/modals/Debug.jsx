import React from 'react';

import TextareaModal from './Textarea';

const Component = function (props) {
	return (
		<TextareaModal
			open={props.open}
			title={props.title}
			onClose={props.onClose}
			onHelp={props.onHelp}
			rows={20}
			value={props.data}
			readOnly
			allowCopy
			allowDownload
			downloadName="report.txt"
			marginBottom="1em"
		/>
	);
};

export default Component;

Component.defaultProps = {
	open: false,
	data: '',
	title: '',
	onClose: null,
	onHelp: null,
};
