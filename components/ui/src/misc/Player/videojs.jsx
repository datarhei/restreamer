import React from 'react';

import Grid from '@mui/material/Grid';

import videojs from 'video.js';
import 'videojs-overlay';
import 'video.js/dist/video-js.min.css';
import './video-js-skin-internal.min.css';
import './video-js-skin-public.min.css';
import 'videojs-overlay/dist/videojs-overlay.css';

export default function VideoJS(props) {
	const videoRef = React.useRef(null);
	const playerRef = React.useRef(null);
	const { options, onReady } = props;

	React.useEffect(() => {
		// make sure Video.js player is only initialized once
		if (!playerRef.current) {
			const videoElement = videoRef.current;
			if (!videoElement) return;

			const player = (playerRef.current = videojs(videoElement, options, () => {
				onReady && onReady(player);
			}));

			// add internal/public skin style
			if (props.type === 'videojs-public') {
				player.addClass('vjs-public');
			} else {
				player.addClass('vjs-internal');
			}
			player.addClass('video-js');
			player.addClass('vjs-16-9');
		} else {
			// you can update player here [update player through props]
			// const player = playerRef.current;
			// player.autoplay(options.autoplay);
			// player.src(options.sources);
		}
	}, [options, videoRef, onReady, props.type]);

	// Dispose the Video.js player when the functional component unmounts
	React.useEffect(() => {
		const player = playerRef.current;

		return () => {
			if (player) {
				player.dispose();
				playerRef.current = null;
			}
		};
	}, [playerRef]);

	return (
		<Grid
			container
			direction="column"
			justifyContent="center"
			alignItems="center"
			spacing={2}
			style={{ position: 'absolute', top: 0, left: 0, bottom: 0, right: 0 }}
		>
			<div data-vjs-player>
				<video ref={videoRef} controls playsInline />
			</div>
		</Grid>
	);
}

VideoJS.defaultProps = {
	type: 'videojs-internal',
};
