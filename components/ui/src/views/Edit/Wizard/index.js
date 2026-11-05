import React from 'react';
import { useNavigate, useParams } from 'react-router-dom';

import { useLingui } from '@lingui/react';
import { t } from '@lingui/macro';
import Backdrop from '@mui/material/Backdrop';
import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import Grid from '@mui/material/Grid';
import MenuItem from '@mui/material/MenuItem';
import Typography from '@mui/material/Typography';

import * as Decoders from '../../../misc/coders/Decoders';
import * as Encoders from '../../../misc/coders/Encoders';
import * as M from '../../../utils/metadata';
import FullSources from '../Sources';
import H from '../../../utils/help';
import NotifyContext from '../../../contexts/Notify';
import Sources from './Sources';

import Source from './Source';
import Video from './Video';
import VideoProfile from './VideoProfile';
import Audio from './Audio';
import Abort from './Abort';
import Error from './Error';
import Saving from './Saving';
import Probe from './Probe';
import License from './License';
import Metadata from './Metadata';

export default function Wizard(props) {
	const { i18n } = useLingui();
	const navigate = useNavigate();
	const { channelid: _channelid } = useParams();
	const notify = React.useContext(NotifyContext);
	const [$data, setData] = React.useState(M.getDefaultIngestMetadata());
	const [$sources, setSources] = React.useState({
		video: M.initSource('video', null),
		audio: M.initSource('audio', null),
	});
	const [$profile, setProfile] = React.useState(M.initProfile());
	const [$skills, setSkills] = React.useState({});
	const [$config, setConfig] = React.useState({});
	const [$step, setStep] = React.useState('TYPE');
	const [$sourceid, setSourceid] = React.useState('');
	const [$probe, setProbe] = React.useState({
		probing: false,
		status: 'none',
	});
	const [$skillsRefresh, setSkillsRefresh] = React.useState(false);
	const [$abort, setAbort] = React.useState({
		step: 'TYPE',
	});
	const [$ready, setReady] = React.useState(false);
	const [$invalid, setInvalid] = React.useState(false);

	React.useEffect(() => {
		(async () => {
			await load();
		})();
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

	React.useEffect(() => {
		if ($invalid === true) {
			navigate('/', { replace: true });
		}
	}, [navigate, $invalid]);

	const load = async () => {
		const channelid = props.restreamer.SelectChannel(_channelid);
		if (channelid === '' || channelid !== _channelid) {
			setInvalid(true);
			return;
		}

		const skills = await props.restreamer.Skills();
		setSkills(skills);

		const config = await props.restreamer.ConfigActive();
		setConfig(config);

		setData({
			...$data,
			meta: {
				name: props.restreamer.GetChannel(_channelid).name,
			},
		});

		setReady(true);
	};

	const refreshSkills = async () => {
		await props.restreamer.RefreshSkills();

		const skills = await props.restreamer.Skills();
		setSkills(skills);
	};

	const probe = async (type, source) => {
		setProbe({
			...$probe,
			probing: true,
			status: 'none',
		});

		let [res, err] = await props.restreamer.Probe(_channelid, source.inputs);
		if (err !== null) {
			res = {
				streams: [],
				log: [err.message],
			};
		}

		let status = M.analyzeStreams(type, res.streams);

		if (status === 'success') {
			if (type === 'video') {
				const profile = M.preselectProfile('video', res.streams, $profile, $skills.encoders);

				setProfile({
					...$profile,
					...profile,
				});

				if (profile.video.encoder.coder === 'none') {
					status = 'nocoder';
				}
			} else {
				const profile = M.preselectProfile('audio', res.streams, $profile, $skills.encoders);

				setProfile({
					...$profile,
					...profile,
				});

				if (profile.audio.encoder.coder === 'none') {
					status = 'nocoder';
				}
			}
		}

		setProbe({
			...$probe,
			probing: false,
			status: status,
		});

		setSources({
			...$sources,
			[type]: {
				type: source.type,
				settings: source.settings,
				inputs: source.inputs,
				streams: res.streams,
			},
		});

		return status === 'success';
	};

	const handleDone = async () => {
		const data = $data;

		data.sources = M.cleanupSources($sources);
		data.profiles = [M.cleanupProfile($profile)];

		const sources = data.sources;
		const profiles = data.profiles;
		const control = data.control;

		const [global, inputs, outputs] = M.createInputsOutputs(sources, profiles, true);

		if (inputs.length === 0 || outputs.length === 0) {
			notify.Dispatch('error', 'save:ingest', i18n._(t`The input profile is not complete. Please define a video and audio source.`));
			return false;
		}

		data.streams = M.createOutputStreams(sources, profiles);

		const [, err] = await props.restreamer.UpsertIngest(_channelid, global, inputs, outputs, control);
		if (err !== null) {
			notify.Dispatch('error', 'save:ingest', err.message);
			return false;
		}

		// Save the metadata
		await props.restreamer.SetIngestMetadata(_channelid, data);

		// Create update the ingest snapshot process
		await props.restreamer.UpsertIngestSnapshot(_channelid, control);

		// Create/update the player
		await props.restreamer.UpdatePlayer(_channelid);

		// Create/update the playersite
		await props.restreamer.UpdatePlayersite();

		notify.Dispatch('success', 'save:ingest', i18n._(t`Main channel saved`));

		navigate(`/${_channelid}`);

		return true;
	};

	const handleAbort = () => {
		setAbort({
			...$abort,
			step: $step,
		});

		setStep('ABORT');
	};

	const handleAdvanced = () => {
		navigate(`/${_channelid}/edit`);
	};

	const handleHelp = (what) => () => {
		H('wizard-' + what);
	};

	if ($ready === false) {
		return null;
	}

	let handleNext = null;
	let handleBack = null;

	if ($step === 'TYPE') {
		handleNext = (sourceid) => () => {
			setSourceid(sourceid);
			setStep('VIDEO SETTINGS');
		};

		handleBack = () => {};

		let knownSources = [];
		for (let s in $skills.sources) {
			if (s === 'network') {
				knownSources.push('network');
				if ($skills.protocols.input.includes('rtmp')) {
					knownSources.push('rtmp');
				}
				if ($skills.protocols.input.includes('http')) {
					knownSources.push('hls');
				}
				if ($skills.protocols.input.includes('srt')) {
					knownSources.push('srt');
				}
			} else if (s === 'video4linux2') {
				knownSources.push('video4linux2');
			} else if (s === 'raspicam') {
				knownSources.push('raspicam');
			} else if (s === 'avfoundation') {
				knownSources.push('avfoundation');
			}
		}

		let availableSources = [];

		for (let s of Sources.List()) {
			if (knownSources.indexOf(s.id) === -1) {
				continue;
			}

			const Icon = s.icon;

			availableSources.push(
				<Grid item xs={6} align="center" key={s.id}>
					<Button variant="big" onClick={handleNext(s.id)}>
						<div>
							<Icon fontSize="large" />
							<Typography>{s.name}</Typography>
						</div>
					</Button>
				</Grid>
			);
		}

		// STEP 1 - Source Type Selection
		return <Source onAbort={handleAbort} onHelp={handleHelp('video-setup')} onAdvanced={handleAdvanced} sources={availableSources} />;
	} else if ($step === 'VIDEO SETTINGS') {
		handleNext = async () => {
			// probing ...
			setStep('VIDEO PROBE');

			const source = $sources.video;

			const status = await probe('video', source);
			if (status === true) {
				setStep('VIDEO RESULT');
				return;
			}

			setStep('VIDEO SETTINGS');
		};

		handleBack = () => {
			setProbe({
				...$probe,
				status: 'none',
			});

			setSources({
				...$sources,
				video: M.initSource('video', null),
			});

			setStep('TYPE');
		};

		const handleChange = (type, settings, inputs, ready) => {
			const source = $sources.video;

			source.type = type;
			source.settings = settings;
			source.inputs = inputs;
			source.ready = ready;

			setSources({
				...$sources,
				video: source,
			});
		};

		const handleRefresh = async () => {
			setSkillsRefresh(true);
			await refreshSkills();
			setSkillsRefresh(false);
		};

		const s = Sources.Get($sourceid);
		if (s === null) {
			setStep('TYPE');
			return null;
		}

		const Component = s.component;

		// STEP 2 - Source Settings
		return (
			<Video
				onAbort={handleAbort}
				onHelp={handleHelp('video-settings')}
				onBack={handleBack}
				onNext={handleNext}
				status={$probe.status}
				sourceid={$sourceid}
				ready={$sources.video.ready}
			>
				<Component
					knownDevices={$skills.sources[s.type]}
					config={$config.source[s.type]}
					settings={$sources.video.settings}
					skills={$skills}
					onChange={handleChange}
					onRefresh={handleRefresh}
				/>
				<Backdrop open={$skillsRefresh}>
					<CircularProgress color="inherit" />
				</Backdrop>
			</Video>
		);
	}
	// STEP 3 - Source Probe
	else if ($step === 'VIDEO PROBE') {
		return <Probe onAbort={handleAbort} />;
	} else if ($step === 'VIDEO RESULT') {
		handleNext = () => {
			const streams = $sources.video.streams;
			const videoprofile = $profile.video;

			const encoder = Encoders.Video.Get(videoprofile.encoder.coder);
			let defaults = encoder.defaults({}, $skills);

			videoprofile.encoder.settings = defaults.settings;
			videoprofile.encoder.mapping = defaults.mapping;

			const decoder = Decoders.Video.Get(videoprofile.decoder.coder);
			defaults = decoder.defaults({}, $skills);

			videoprofile.decoder.settings = defaults.settings;
			videoprofile.decoder.mapping = defaults.mapping;

			const audioprofile = $profile.audio;
			audioprofile.source = -1;

			// set default for first video audio track
			for (let s of streams) {
				if (s.type !== 'audio') {
					continue;
				}

				audioprofile.source = 0;
				audioprofile.stream = s.stream;

				break;
			}

			// set default for silence audio track if the video doesn't have an audio track
			if (audioprofile.source === -1) {
				audioprofile.source = 1;

				const fullSource = FullSources.Get('virtualaudio');
				const source = $sources.audio;

				source.type = fullSource.id;
				source.settings = fullSource.func.initSettings({
					source: 'silence',
					layout: 'stereo',
					sampling: '44100',
				});
				source.inputs = fullSource.func.createInputs(source.settings);

				setSources({
					...$sources,
					audio: M.initSource('audio', source),
				});
			}

			setProfile({
				...$profile,
				audio: audioprofile,
				video: profile,
			});

			setStep('AUDIO SETTINGS');
		};

		handleBack = () => {
			setStep('VIDEO SETTINGS');
		};

		const handleStreamChange = (event) => {
			const value = event.target.value;

			const profile = $profile.video;
			profile.stream = parseInt(value);

			setProfile({
				...$profile,
				video: profile,
			});
		};

		const handleEncoderChange = (event) => {
			const value = event.target.value;

			const profile = $profile.video;
			profile.encoder.coder = value;

			setProfile({
				...$profile,
				video: profile,
			});
		};

		const handleDecoderChange = (event) => {
			const value = event.target.value;

			const profile = $profile.video;
			profile.decoder.coder = value;

			setProfile({
				...$profile,
				video: profile,
			});
		};

		const isCompatible = (stream) => {
			if (stream.codec === 'h264') {
				return true;
			}

			return false;
		};

		const streams = $sources.video.streams;
		const profile = $profile.video;

		const compatible = isCompatible(streams[profile.stream]);
		const decoders = Decoders.Video.GetCodersForCodec(streams[profile.stream].codec, $skills.decoders.video, 'any');

		let decodersList = [];

		for (let c of decoders) {
			decodersList.push(
				<MenuItem value={c.coder} key={c.coder}>
					{c.name}
				</MenuItem>
			);
		}

		const encoders = Encoders.Video.GetCodersForCodec('h264', $skills.encoders.video, 'any');

		let encodersList = [];

		for (let c of encoders) {
			encodersList.push(
				<MenuItem value={c.coder} key={c.coder}>
					{c.name}
				</MenuItem>
			);
		}

		let streamList = [];

		for (let s of streams) {
			if (s.type !== 'video') {
				continue;
			}

			streamList.push(
				<MenuItem value={s.stream} key={s.stream}>
					{s.width}x{s.height}, {s.codec.toUpperCase()}
				</MenuItem>
			);
		}

		// STEP 4 - Video Profile Selection
		return (
			<VideoProfile
				onAbort={handleAbort}
				onHelp={handleHelp('video-result')}
				onBack={handleBack}
				onNext={handleNext}
				compatible={compatible}
				stream={$profile.video.stream}
				streamList={streamList}
				onStreamChange={handleStreamChange}
				decoder={$profile.video.decoder.coder}
				decodersList={decodersList}
				onDecoderChange={handleDecoderChange}
				encoder={$profile.video.encoder.coder}
				encodersList={encodersList}
				onEncoderChange={handleEncoderChange}
			/>
		);
	} else if ($step === 'AUDIO SETTINGS') {
		handleNext = async () => {
			if ($profile.audio.source === 1) {
				const source = $sources.audio;

				// probing ...

				setSources({
					...$sources,
					audio: source,
				});

				setStep('AUDIO PROBE');

				const status = await probe('audio', source);
				if (status === true) {
					setStep('AUDIO RESULT');
					return;
				}

				setStep('AUDIO SETTINGS');
			} else {
				setStep('AUDIO RESULT');
			}
		};

		handleBack = () => {
			setStep('VIDEO RESULT');
		};

		const handleAudioStreamChange = (event) => {
			const value = event.target.value;

			const profile = $profile.audio;
			profile.stream = parseInt(value);

			setProfile({
				...$profile,
				audio: profile,
			});
		};

		const handleAudioDeviceChange = (event) => {
			const value = event.target.value;

			const source = $sources.audio;
			source.settings.address = value;
			source.settings.device = value.replace(/^hw:/, '');

			setSources({
				...$sources,
				audio: source,
			});
		};

		const handleStream = (event) => {
			const value = event.target.value;

			const profile = $profile.audio;
			let source = null;

			if (value === 'video') {
				profile.source = 0;

				source = null;
			} else if (value === 'alsa') {
				profile.source = 1;

				// The first ALSA device is selected by default
				let address = '';
				if ($skills.sources['alsa'].length !== 0) {
					address = $skills.sources['alsa'][0].id;
				}

				source = M.initSource('audio', null);

				const fullSource = FullSources.Get('alsa');

				source.type = fullSource.id;
				source.settings = fullSource.func.initSettings({
					address: address,
					device: address.replace(/^hw:/, ''),
				});
				source.inputs = fullSource.func.createInputs(source.settings);
			} else if (value === 'silence') {
				profile.source = 1;

				source = M.initSource('audio', null);

				const fullSource = FullSources.Get('virtualaudio');

				source.type = fullSource.id;
				source.settings = fullSource.func.initSettings({
					source: 'silence',
					layout: 'stereo',
					sampling: 44100,
				});
				source.inputs = fullSource.func.createInputs(source.settings);
			} else {
				profile.source = -1;

				source = null;
			}

			setSources({
				...$sources,
				audio: M.initSource('audio', source),
			});

			setProfile({
				...$profile,
				audio: profile,
			});
		};

		const profile = $profile.audio;
		const source = $sources.audio;

		let streamList = [];
		const streams = $sources.video.streams;

		for (let s of streams) {
			if (s.type !== 'audio') {
				continue;
			}

			streamList.push(
				<MenuItem value={s.stream} key={s.stream}>
					{s.codec.toUpperCase()} {s.layout} {s.sampling_hz}Hz
				</MenuItem>
			);
		}

		let deviceList = [];

		if ('alsa' in $skills.sources && ['video4linux2', 'raspicam'].includes($sources.video.type)) {
			for (let device of $skills.sources['alsa']) {
				if (!source.settings.address) {
					source.settings.address = device.id;
				}

				deviceList.push(
					<MenuItem key={device.id} value={device.id}>
						{device.name} ({device.id})
					</MenuItem>
				);
			}
		}

		let radioValue = 'video';

		if (profile.source === 0) {
			// video input has an audio track
			// options:
			// 1. use audio track (and encode it if it is not compatible)
			// 2. add silence
			// 3. don't use any audio
		} else if (profile.source === 1) {
			if ($sources.audio.type === 'alsa') {
				// ALSA audio device
				radioValue = 'alsa';
			} else {
				// silence
				radioValue = 'silence';
			}
		} else {
			radioValue = 'none';
			// video input doesn't have an audio track

			// options:
			// 1. add silence
			// 2. don't use any audio
		}

		return (
			<Audio
				onAbort={handleAbort}
				onHelp={handleHelp('audio-settings')}
				onBack={handleBack}
				onNext={handleNext}
				status={$probe.status}
				source={radioValue}
				onSource={handleStream}
				streamList={streamList}
				deviceList={deviceList}
				stream={profile.stream}
				onAudioStreamChange={handleAudioStreamChange}
				address={source.settings.address}
				onAudioDeviceChange={handleAudioDeviceChange}
			/>
		);
	} else if ($step === 'AUDIO PROBE') {
		return <Probe onAbort={handleAbort} />;
	} else if ($step === 'AUDIO RESULT') {
		handleNext = () => {
			let stream = null;
			const profile = $profile.audio;

			if (profile.source === 0) {
				stream = $sources.video.streams[profile.stream];
			} else if (profile.source === 1) {
				stream = $sources.audio.streams[profile.stream];
			} else {
				profile.source = -1;
				profile.stream = -1;
			}

			if (stream !== null) {
				const compatible = isCompatible(stream);

				if (compatible === true) {
					profile.coder = 'copy';
				} else {
					profile.coder = 'aac';
				}

				const encoder = Encoders.Audio.Get(profile.coder);
				const defaults = encoder.defaults(stream, $skills);

				profile.encoder.settings = defaults.settings;
				profile.encoder.mapping = defaults.mapping;
			}

			setProfile({
				...$profile,
				audio: profile,
			});

			setStep('META');
		};

		const isCompatible = (stream) => {
			if (stream.codec === 'aac' || stream.codec === 'mp3') {
				return true;
			}

			return false;
		};

		handleNext();

		return null;
	} else if ($step === 'META') {
		handleNext = () => {
			setStep('LICENSE');
		};

		handleBack = () => {
			setStep('AUDIO SETTINGS');
		};

		const handleMetadataChange = (settings) => {
			setData({
				...$data,
				meta: settings,
			});
		};

		return (
			<Metadata
				onAbort={handleAbort}
				onHelp={handleHelp('audio-result')}
				onBack={handleBack}
				onNext={handleNext}
				onChange={handleMetadataChange}
				metadata={$data.meta}
			/>
		);
	} else if ($step === 'LICENSE') {
		handleNext = async () => {
			setStep('SAVING');

			const res = await handleDone();

			if (res === false) {
				setStep('DONE');
			} else {
				setStep('ERROR');
			}
		};

		handleBack = () => {
			setStep('META');
		};

		const handleLicenseChange = (license) => {
			setData({
				...$data,
				license: license,
			});
		};

		return (
			<License
				onAbort={handleAbort}
				onHelp={handleHelp('license')}
				onBack={handleBack}
				onNext={handleNext}
				onChange={handleLicenseChange}
				license={$data.license}
			/>
		);
	} else if ($step === 'SAVING') {
		return <Saving onAbort={handleAbort} />;
	} else if ($step === 'DONE') {
		return null;
	} else if ($step === 'ERROR') {
		handleNext = () => {
			setStep('TYPE');
		};

		return <Error onAbort={handleAbort} onHelp={handleHelp('error')} />;
	} else if ($step === 'ABORT') {
		const nchannels = props.restreamer.ListChannels().length;

		handleBack = () => {
			setStep($abort.step);
		};

		handleNext = () => {
			props.restreamer.DeleteChannel(_channelid);

			// Select a channel to jump back to
			const channels = props.restreamer.ListChannels();
			props.restreamer.SelectChannel(channels[0].channelid);

			navigate(`/`);
		};

		return <Abort onHelp={handleHelp('abort')} onBack={handleBack} onNext={handleNext} nchannels={nchannels} />;
	}

	return null;
}

Wizard.defaultProps = {
	restreamer: null,
};
