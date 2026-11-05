import React from 'react';

import { i18n } from '@lingui/core';

export default function Number(props) {
	const options = {
		minimumFractionDigits: props.minDigits,
		maximumFractionDigits: props.digits,
	};

	if (options.minimumFractionDigits > options.maximumFractionDigits) {
		options.maximumFractionDigits = options.minimumFractionDigits;
	}

	return <React.Fragment>{i18n.number(props.value, options)}</React.Fragment>;
}

Number.defaultProps = {
	value: 0,
	digits: 0,
	minDigits: 0,
};
