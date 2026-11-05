import React from 'react';

import TextareaModal from './Textarea';
import { withDefaults } from '../../utils/props';

const Component = function (inputProps) {
	const props = withDefaults(inputProps, ComponentDefaultProps);
	return (
		<TextareaModal open={props.open} title={props.title} onClose={props.onClose} onHelp={props.onHelp} rows={18} value={props.data} readOnly allowCopy />
	);
};

export default Component;

const ComponentDefaultProps = {
	open: false,
	data: '',
	title: '',
	onClose: null,
	onHelp: null,
};
