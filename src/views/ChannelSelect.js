import React from 'react';
import { useNavigate } from 'react-router-dom';

export default function ChannelSelector(props) {
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

ChannelSelector.defaultProps = {
	channelid: '',
};
