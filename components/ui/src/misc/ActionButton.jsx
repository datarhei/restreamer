import React from 'react';

import { Trans } from '@lingui/react/macro';
import Button from '@mui/material/Button';
import { withDefaults } from '../utils/props';

export default function ActionButton(inputProps) {
	const props = withDefaults(inputProps, ActionButtonDefaultProps);
	let button = null;

	if (props.state === 'connecting') {
		button = (
			<Button variant="outlined" fullWidth disabled>
				<Trans>Connecting ...</Trans>
			</Button>
		);
	} else if (props.state === 'disconnecting') {
		button = (
			<Button variant="outlined" fullWidth disabled>
				<Trans>Disconnecting ...</Trans>
			</Button>
		);
	} else if (props.state === 'connected') {
		button = (
			<Button variant="outlined" fullWidth color="secondary" disabled={props.disabled} onClick={props.onDisconnect}>
				<Trans>Disconnect</Trans>
			</Button>
		);
	} else if (props.state === 'disconnected') {
		if (props.reconnect < 0) {
			if (props.order === 'start') {
				button = (
					<Button variant="outlined" fullWidth color="primary" disabled={props.disabled} onClick={props.onReconnect}>
						<Trans>Reconnect</Trans>
					</Button>
				);
			} else {
				button = (
					<Button variant="outlined" fullWidth color="primary" disabled={props.disabled} onClick={props.onConnect}>
						<Trans>Connect</Trans>
					</Button>
				);
			}
		} else {
			button = (
				<Button variant="outlined" fullWidth color="secondary" disabled={props.disabled} onClick={props.onDisconnect}>
					<Trans>Disconnect</Trans>
				</Button>
			);
		}
	} else if (props.state === 'error') {
		if (props.reconnect < 0) {
			button = (
				<Button variant="outlined" fullWidth color="primary" disabled={props.disabled} onClick={props.onReconnect}>
					<Trans>Reconnect</Trans>
				</Button>
			);
		} else {
			button = (
				<Button variant="outlined" fullWidth color="secondary" disabled={props.disabled} onClick={props.onDisconnect}>
					<Trans>Disconnect</Trans>
				</Button>
			);
		}
	}

	return button;
}

const ActionButtonDefaultProps = {
	order: 'stop',
	state: 'disconnected',
	reconnect: -1,
	disabled: false,
	onDisconnect: function () {},
	onConnect: function () {},
	onReconnect: function () {},
};
