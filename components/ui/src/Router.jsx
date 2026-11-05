import React from 'react';
import { Route, Navigate, Routes, HashRouter as DOMRouter } from 'react-router-dom';

import { withDefaults } from './utils/props';

const AddService = React.lazy(() => import('./views/Publication/Add'));
const ChannelSelect = React.lazy(() => import('./views/ChannelSelect'));
const Edit = React.lazy(() => import('./views/Edit'));
const EditPlayer = React.lazy(() => import('./views/Publication/Player'));
const EditService = React.lazy(() => import('./views/Publication/Edit'));
const Main = React.lazy(() => import('./views/Main'));
const Playersite = React.lazy(() => import('./views/Playersite'));
const Settings = React.lazy(() => import('./views/Settings'));
const Wizard = React.lazy(() => import('./views/Edit/Wizard'));

function LoadingRoute() {
	return (
		<div role="status" aria-live="polite">
			Loading…
		</div>
	);
}

export default function Router(inputProps) {
	const props = withDefaults(inputProps, RouterDefaultProps);
	if (props.restreamer === null) {
		return null;
	}

	const channelid = props.restreamer.GetCurrentChannelID();

	return (
		<DOMRouter>
			<React.Suspense fallback={<LoadingRoute />}>
				<Routes>
					<Route path="/" element={<ChannelSelect channelid={channelid} />} />
					<Route path="/playersite" element={<Playersite restreamer={props.restreamer} />} />
					<Route path="/settings" element={<Settings restreamer={props.restreamer} />} />
					<Route path="/settings/:tab" element={<Settings restreamer={props.restreamer} />} />
					<Route path="/:channelid" element={<Main key={channelid} restreamer={props.restreamer} />} />
					<Route path="/:channelid/edit" element={<Edit key={channelid} restreamer={props.restreamer} />} />
					<Route path="/:channelid/edit/wizard" element={<Wizard key={channelid} restreamer={props.restreamer} />} />
					<Route path="/:channelid/edit/:tab" element={<Edit key={channelid} restreamer={props.restreamer} />} />
					<Route path="/:channelid/publication" element={<AddService key={channelid} restreamer={props.restreamer} />} />
					<Route path="/:channelid/publication/player" element={<EditPlayer key={channelid} restreamer={props.restreamer} />} />
					<Route path="/:channelid/publication/:service/:index" element={<EditService key={channelid} restreamer={props.restreamer} />} />
					<Route path="*" element={<Navigate to="/" replace />} />
				</Routes>
			</React.Suspense>
		</DOMRouter>
	);
}

const RouterDefaultProps = {
	restreamer: null,
};
