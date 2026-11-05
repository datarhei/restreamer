import React from 'react';

import { Trans } from '@lingui/macro';
import Backdrop from '@mui/material/Backdrop';
import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import Divider from '@mui/material/Divider';
import Grid from '@mui/material/Grid';
import Link from '@mui/material/Link';
import Typography from '@mui/material/Typography';
import WarningIcon from '@mui/icons-material/Warning';

import * as M from '../../utils/metadata';
import BoxText from '../../misc/BoxText';
import EncodingSelect from '../../misc/EncodingSelect';
import PaperFooter from '../../misc/PaperFooter';
import ProbeModal from '../../misc/modals/Probe';
import HintModal from '../../misc/modals/Hint';
import SourceSelect from './SourceSelect';
import StreamSelect from './StreamSelect';

import FilterSelect from '../../misc/FilterSelect';

export default function Profile(props) {
	const [$sources, setSources] = React.useState({
		video: M.initSource('video', props.sources[0]),
		audio: M.initSource('audio', props.sources[1]),
	});
	const [$profile, setProfile] = React.useState(M.initProfile(props.profile));
	const [$videoProbe, setVideoProbe] = React.useState({
		probing: false,
		log: [],
		modal: false,
		status: 'none',
	});
	const [$audioProbe, setAudioProbe] = React.useState({
		probing: false,
		log: [],
		modal: false,
		status: 'none',
	});
	const [$skillsRefresh, setSkillsRefresh] = React.useState(false);
	const [$probeModal, setProbeModal] = React.useState({
		open: false,
		data: '',
	});
	const [$hintModal, setHintModal] = React.useState({
		open: false,
		type: '',
		streams: [],
	});
	const [$activeStep, setActiveStep] = React.useState(props.startWith === 'audio' ? 1 : 0);
	const [$ready, setReady] = React.useState(false);

	React.useEffect(() => {
		(async () => {
			await load();
		})();
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

	const load = async () => {
		// Add pseudo sources
		props.skills.sources.noaudio = [];

		let audio = $sources.audio;

		let hasAudio = false;
		for (let i = 0; i < $sources.video.streams.length; i++) {
			if ($sources.video.streams[i].type === 'audio') {
				hasAudio = true;
				break;
			}
		}

		if (hasAudio === true) {
			props.skills.sources.videoaudio = [];
		} else {
			delete props.skills.sources.videoaudio;
		}

		setSources({
			...$sources,
			audio: audio,
		});

		setReady(true);
	};

	const handleNextStep = () => {
		setActiveStep((prevActiveStep) => prevActiveStep + 1);
	};

	const handleBackStep = () => {
		setActiveStep((prevActiveStep) => prevActiveStep - 1);
	};

	const handleProbe = async (type, device, settings, inputs) => {
		if (type === 'video') {
			setVideoProbe({
				...$videoProbe,
				probing: true,
				status: 'none',
			});
		} else {
			setAudioProbe({
				...$audioProbe,
				probing: true,
				status: 'none',
			});
		}

		const res = await props.onProbe(inputs);

		const status = handleProbeStreams(type, device, settings, inputs, res);

		return status === 'success';
	};

	const handleProbeStreams = (type, device, settings, inputs, res) => {
		let status = M.analyzeStreams(type, res.streams);

		if (type === 'video') {
			let audio = $sources.audio;

			const profile = M.preselectProfile('video', res.streams, $profile, props.skills.encoders, audio.type === '');

			// Add pseudo sources
			props.skills.sources.noaudio = [];

			let hasAudio = false;
			for (let i = 0; i < res.streams.length; i++) {
				if (res.streams[i].type === 'audio') {
					hasAudio = true;
					break;
				}
			}

			if (hasAudio === true) {
				props.skills.sources.videoaudio = [];
				if (audio.type === '') {
					audio.type = 'videoaudio';
				}
			} else {
				delete props.skills.sources.videoaudio;
				if (audio.type === '' || audio.type === 'videoaudio') {
					audio.type = 'noaudio';
					profile.audio.source = -1;
					profile.audio.stream = -1;
					profile.custom.selected = false;
					profile.custom.stream = -1;
				}
			}

			audio = M.initSource('audio', audio);

			setProfile({
				...$profile,
				...profile,
			});

			setVideoProbe({
				...$videoProbe,
				probing: false,
				log: res.log,
				status: status,
			});

			setAudioProbe({
				...$audioProbe,
				status: audio.type === 'noaudio' ? 'success' : 'none',
			});

			setSources({
				...$sources,
				audio: audio,
				video: {
					type: device,
					settings: settings,
					inputs: inputs,
					streams: res.streams,
				},
			});
		} else {
			const profile = M.preselectProfile('audio', res.streams, $profile, props.skills.encoders);

			setProfile({
				...$profile,
				...profile,
			});

			setAudioProbe({
				...$audioProbe,
				probing: false,
				log: res.log,
				status: status,
			});

			setSources({
				...$sources,
				audio: {
					type: device,
					settings: settings,
					inputs: inputs,
					streams: res.streams,
				},
			});
		}
	};

	const handleRefresh = async () => {
		setSkillsRefresh(true);
		await props.onRefresh();
		setSkillsRefresh(false);
	};

	const handleStore = async (name, data) => {
		return await props.onStore(name, data);
	};

	const handleEncoding = (type) => (encoder, decoder) => {
		const profile = $profile[type];

		profile.encoder = encoder;
		profile.decoder = decoder;

		setProfile({
			...$profile,
			[type]: profile,
		});
	};

	const handleFilter = (type) => (filter) => {
		const profile = $profile[type];

		profile.filter = filter;

		setProfile({
			...$profile,
			[type]: profile,
		});
	};

	const handleDone = () => {
		const sources = M.cleanupSources($sources);
		const profile = M.cleanupProfile($profile);

		props.onDone(sources, profile);
	};

	const handleAbort = () => {
		props.onAbort();
	};

	const handleProbeLogModal = (type) => (event) => {
		event.preventDefault();

		if (type === 'video') {
			setProbeModal({
				...$probeModal,
				open: true,
				data: $videoProbe.log.join('\n'),
			});
		} else if (type === 'audio') {
			setProbeModal({
				...$probeModal,
				open: true,
				data: $audioProbe.log.join('\n'),
			});
		} else {
			setProbeModal({
				...$probeModal,
				open: false,
				data: '',
			});
		}
	};

	const handleSourceChange = (type, source) => {
		const profile = $profile[type];
		const custom = $profile.custom;

		if (type === 'audio') {
			if (source === 'noaudio') {
				custom.selected = false;
				custom.stream = -1;
				profile.source = -1;
				profile.stream = -1;
			} else if (source === 'videoaudio') {
				custom.selected = false;
				profile.source = 0;

				for (let i = 0; i < $sources.video.streams.length; i++) {
					if ($sources.video.streams[i].type === 'audio') {
						custom.stream = i;
						profile.stream = i;
						break;
					}
				}
			} else {
				custom.selected = true;
				custom.stream = -2;

				profile.source = 1;
				profile.stream = -1;
			}

			let audio = $sources.audio;
			audio.type = source;

			setSources({
				...$sources,
				audio: audio,
			});
		} else {
			let video = $sources.video;
			video.type = source;

			setSources({
				...$sources,
				video: video,
			});
		}

		setProfile({
			...$profile,
			[type]: profile,
			custom: custom,
		});
	};

	const handleSourceSettingsChange = (type, source, settings) => {
		if (type === 'video') {
			setVideoProbe({
				...$videoProbe,
				status: 'none',
			});
		} else {
			setAudioProbe({
				...$audioProbe,
				status: 'none',
			});
		}
	};

	const handleStreamSelect = (type, what) => (stream) => {
		const profile = $profile;

		profile[type].stream = stream;

		if (what === 'custom') {
			profile.custom.stream = stream;
		}

		setProfile({
			...$profile,
			...profile,
		});
	};

	const handleHintModal = (type, streams) => (event) => {
		if (event) {
			event.preventDefault();
		}

		if (!streams) {
			streams = [];
		}

		if (streams.length > 0) {
			return streams;
		}

		if (type === 'video') {
			streams = [
				{
					url: '',
					index: 0,
					stream: 0,
					type: 'video',
					codec: 'h264',
					width: 1920,
					height: 1080,
					pix_fmt: 'yuv420p',
					sampling_hz: 0,
					layout: '',
					channels: 0,
				},
			];
		} else if (type === 'audio') {
			streams = [
				{
					url: '',
					index: 1,
					stream: 0,
					type: 'audio',
					codec: 'aac',
					width: 0,
					height: 0,
					sampling_hz: '44100',
					layout: 'stereo',
					channels: 2,
				},
			];
		}

		if (type === 'video') {
			setHintModal({
				...$hintModal,
				open: true,
				type: type,
				streams: streams,
			});
		} else if (type === 'audio') {
			setHintModal({
				...$hintModal,
				open: true,
				type: type,
				streams: streams,
			});
		} else {
			setHintModal({
				...$hintModal,
				open: false,
				type: '',
				streams: [],
			});
		}
	};

	const handleHintChange = (streams) => {
		setHintModal({
			...$hintModal,
			streams: streams,
		});
	};

	const handleHintCancel = () => {
		setHintModal({
			streams: [],
		});

		handleHintModal('none')(null);
	};

	const handleHintDone = () => {
		const type = $hintModal.type;

		const device = $sources[type].type;
		const settings = $sources[type].settings;
		const inputs = $sources[type].inputs;
		const probe = {
			streams: $hintModal.streams,
			log: [],
		};

		const url = inputs[0].address;

		probe.log.push(`Stream hints for input from '${url}'`);

		for (let s of $hintModal.streams) {
			s.url = url;

			let stream = `Stream #${s.index}:${s.stream}: `;
			if (s.type === 'video') {
				stream += `Video: ${s.codec}, ${s.pix_fmt}, ${s.width}x${s.height}`;
			} else if (s.type === 'audio') {
				stream += `Audio: ${s.codec}, ${s.sampling_hz} Hz, ${s.layout}`;
			}

			probe.log.push(stream);
		}

		handleProbeStreams(type, device, settings, inputs, probe);

		handleHintModal('none')(null);
	};

	if ($ready === false) {
		return null;
	}

	return (
		<React.Fragment>
			{$activeStep === 0 && (
				<React.Fragment>
					<Grid container spacing={2}>
						<Grid item xs={12}>
							<Typography variant="h3">
								<Trans>Video settings</Trans>
							</Typography>
						</Grid>
						<Grid item xs={12}>
							<SourceSelect
								type="video"
								skills={props.skills}
								source={$sources.video}
								config={props.config}
								onProbe={handleProbe}
								onChange={handleSourceSettingsChange}
								onRefresh={handleRefresh}
								onStore={handleStore}
							/>
						</Grid>
						{$videoProbe.status !== 'none' && (
							<React.Fragment>
								{$videoProbe.status === 'error' && (
									<Grid item xs={12} align="center">
										<BoxText color="dark">
											<WarningIcon fontSize="large" color="error" />
											<Typography>
												<Trans>
													Failed to probe the source. Please check the{' '}
													<Link color="textSecondary" href="#!" onClick={handleProbeLogModal('video')}>
														probe details
													</Link>
													.
												</Trans>
											</Typography>
											<Typography>
												<Trans>
													In order to proceed anyways, you can provide{' '}
													<Link color="textSecondary" href="#!" onClick={handleHintModal('video', [])}>
														hints
													</Link>{' '}
													about the available streams.
												</Trans>
											</Typography>
										</BoxText>
									</Grid>
								)}
								{$videoProbe.status === 'nostream' && (
									<Grid item xs={12} align="center">
										<BoxText color="dark">
											<WarningIcon fontSize="large" color="error" />
											<Typography>
												<Trans>
													The source doesn't provide any video streams. Please check the{' '}
													<Link href="#!" onClick={handleProbeLogModal('video')}>
														probe details
													</Link>
													.
												</Trans>
											</Typography>
										</BoxText>
									</Grid>
								)}
								{$videoProbe.status === 'success' && (
									<React.Fragment>
										<Grid item xs={12}>
											<StreamSelect
												type="video"
												streams={$sources.video.streams}
												selected={$profile.video.stream}
												onChange={handleStreamSelect('video')}
											/>
										</Grid>
										<Grid item xs={12} align="right">
											<Typography>
												<Trans>
													<Link href="#!" onClick={handleProbeLogModal('video')}>
														Show probe details
													</Link>
												</Trans>
											</Typography>
										</Grid>
										<Grid item xs={12}>
											<EncodingSelect
												type="video"
												streams={$sources.video.streams}
												profile={$profile.video}
												codecs={['copy', 'h264']}
												skills={props.skills}
												onChange={handleEncoding('video')}
											/>
										</Grid>
										{$profile.video.encoder.coder !== 'none' && $profile.video.encoder.coder !== 'copy' && (
											<Grid item xs={12}>
												<FilterSelect
													type="video"
													profile={$profile.video}
													availableFilters={props.skills.filter}
													onChange={handleFilter('video')}
												/>
											</Grid>
										)}
									</React.Fragment>
								)}
							</React.Fragment>
						)}
						<Grid item xs={12}>
							<Divider />
						</Grid>
					</Grid>
					<PaperFooter
						buttonsLeft={
							<React.Fragment>
								<Button variant="outlined" color="default" onClick={handleAbort}>
									<Trans>Abort</Trans>
								</Button>
								<Button variant="outlined" color="primary" disabled={$videoProbe.status !== 'success'} onClick={handleNextStep}>
									<Trans>Next: Audio</Trans>
								</Button>
							</React.Fragment>
						}
					/>
				</React.Fragment>
			)}
			{$activeStep === 1 && (
				<React.Fragment>
					<Grid container spacing={2}>
						<Grid item xs={12}>
							<Typography variant="h3">
								<Trans>Audio settings</Trans>
							</Typography>
						</Grid>
						<Grid item xs={12}>
							<SourceSelect
								type="audio"
								skills={props.skills}
								source={$sources.audio}
								config={props.config}
								onProbe={handleProbe}
								onSelect={handleSourceChange}
								onChange={handleSourceSettingsChange}
								onRefresh={handleRefresh}
								onStore={handleStore}
							/>
						</Grid>
						{$profile.custom.selected === false && $profile.custom.stream >= 0 && (
							<React.Fragment>
								<Grid item xs={12}>
									<StreamSelect
										type="audio"
										streams={$sources.video.streams}
										selected={$profile.custom.stream}
										onChange={handleStreamSelect('audio', 'custom')}
									/>
								</Grid>
								<Grid item xs={12}>
									<EncodingSelect
										type="audio"
										streams={$sources.video.streams}
										profile={$profile.audio}
										codecs={['copy', 'aac', 'mp3']}
										skills={props.skills}
										onChange={handleEncoding('audio')}
									/>
								</Grid>
								{$profile.audio.encoder.coder !== 'none' && $profile.audio.encoder.coder !== 'copy' && $profile.audio.source !== -1 && (
									<Grid item xs={12}>
										<FilterSelect
											type="audio"
											profile={$profile.audio}
											availableFilters={props.skills.filter}
											onChange={handleFilter('audio')}
										/>
									</Grid>
								)}
							</React.Fragment>
						)}
						{$profile.custom.selected === true && (
							<React.Fragment>
								{$audioProbe.status !== 'none' && (
									<React.Fragment>
										{$audioProbe.status === 'error' && (
											<Grid item xs={12} align="center">
												<BoxText color="dark">
													<WarningIcon fontSize="large" color="error" />
													<Typography>
														<Trans>
															Failed to probe the source. Please check the{' '}
															<Link href="#!" onClick={handleProbeLogModal('audio')}>
																probe details
															</Link>
															.
														</Trans>
													</Typography>
													<Typography>
														<Trans>
															In order to proceed anyways, you can provide{' '}
															<Link color="textSecondary" href="#!" onClick={handleHintModal('audio', [])}>
																hints
															</Link>{' '}
															about the available streams.
														</Trans>
													</Typography>
												</BoxText>
											</Grid>
										)}
										{$audioProbe.status === 'nostream' && (
											<Grid item xs={12} align="center">
												<BoxText color="dark">
													<WarningIcon fontSize="large" color="error" />
													<Typography>
														<Trans>
															The source doesn't provide any audio streams. Please check the{' '}
															<Link href="#!" onClick={handleProbeLogModal('audio')}>
																probe details
															</Link>
															.
														</Trans>
													</Typography>
												</BoxText>
											</Grid>
										)}
										{$audioProbe.status === 'success' && (
											<React.Fragment>
												<Grid item xs={12}>
													<StreamSelect
														type="audio"
														streams={$sources.audio.streams}
														selected={$profile.audio.stream}
														onChange={handleStreamSelect('audio')}
													/>
												</Grid>
												<Grid item xs={12} align="right">
													<Typography>
														<Trans>
															<Link href="#!" onClick={handleProbeLogModal('audio')}>
																Show probe details
															</Link>
														</Trans>
													</Typography>
												</Grid>
												<Grid item xs={12}>
													<EncodingSelect
														type="audio"
														streams={$sources.audio.streams}
														profile={$profile.audio}
														codecs={['copy', 'aac', 'mp3']}
														skills={props.skills}
														onChange={handleEncoding('audio')}
													/>
												</Grid>
												{$profile.audio.encoder.coder !== 'none' && $profile.audio.encoder.coder !== 'copy' && (
													<Grid item xs={12}>
														<FilterSelect
															type="audio"
															profile={$profile.audio}
															availableFilters={props.skills.filter}
															onChange={handleFilter('audio')}
														/>
													</Grid>
												)}
											</React.Fragment>
										)}
									</React.Fragment>
								)}
							</React.Fragment>
						)}
						<Grid item xs={12}>
							<Divider />
						</Grid>
					</Grid>
					<PaperFooter
						buttonsLeft={
							<React.Fragment>
								<Button variant="outlined" onClick={handleBackStep}>
									<Trans>Back</Trans>
								</Button>
								<Button
									variant="outlined"
									color="primary"
									disabled={$profile.custom.selected === true && $audioProbe.status !== 'success'}
									onClick={handleDone}
								>
									<Trans>Finish</Trans>
								</Button>
							</React.Fragment>
						}
					/>
				</React.Fragment>
			)}
			<Backdrop open={$videoProbe.probing || $audioProbe.probing || $skillsRefresh}>
				<CircularProgress color="inherit" />
			</Backdrop>
			<ProbeModal open={$probeModal.open} onClose={handleProbeLogModal('none')} data={$probeModal.data} />
			<HintModal
				open={$hintModal.open}
				onClose={handleHintCancel}
				onChange={handleHintChange}
				onDone={handleHintDone}
				title="Stream hints"
				type={$hintModal.type}
				streams={$hintModal.streams}
			/>
		</React.Fragment>
	);
}

Profile.defaultProps = {
	skills: {},
	sources: [],
	profile: {},
	config: {},
	startWith: '',
	onDone: function (sources, profile) {},
	onAbort: function () {},
	onProbe: function (inputs) {
		return {
			streams: [],
			log: ['onProbe function not provided for this component'],
		};
	},
	onRefresh: function () {},
	onStore: function (name, data) {
		return '';
	},
};
