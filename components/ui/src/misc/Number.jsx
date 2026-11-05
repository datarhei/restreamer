import React from 'react';

import { i18n } from '@lingui/core';
import { withDefaults } from '../utils/props';

export default function Number(inputProps) {
	const props = withDefaults(inputProps, NumberDefaultProps);
	const options = {
		minimumFractionDigits: props.minDigits,
		maximumFractionDigits: props.digits,
	};

	if (options.minimumFractionDigits > options.maximumFractionDigits) {
		options.maximumFractionDigits = options.minimumFractionDigits;
	}

	return <React.Fragment>{i18n.number(props.value, options)}</React.Fragment>;
}

const NumberDefaultProps = {
	value: 0,
	digits: 0,
	minDigits: 0,
};
