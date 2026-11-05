import React from 'react';
import { useNavigate } from 'react-router-dom';
import { withDefaults } from '../utils/props';

export default function ChannelSelector(inputProps) {
	const props = withDefaults(inputProps, ChannelSelectorDefaultProps);
	const navigate = useNavigate();
	const [$channelid, setChannelid] = React.useState('');

	React.useEffect(() => {
		onMount();
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

	React.useEffect(() => {
		navigate(`/${$channelid}`, { replace: true });
	}, [navigate, $channelid]);

	const onMount = () => {
		setChannelid(props.channelid);
	};

	return null;
}

const ChannelSelectorDefaultProps = {
	channelid: '',
};
