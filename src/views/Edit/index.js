import React from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import PropTypes from 'prop-types';

import { useLingui } from '@lingui/react';
import { Trans, t } from '@lingui/macro';
import makeStyles from '@mui/styles/makeStyles';
import AutoFixHighIcon from '@mui/icons-material/AutoFixHigh';
import Backdrop from '@mui/material/Backdrop';
import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import Divider from '@mui/material/Divider';
import EditIcon from '@mui/icons-material/Edit';
import Grid from '@mui/material/Grid';
import Tab from '@mui/material/Tab';
import Tabs from '@mui/material/Tabs';
import Typography from '@mui/material/Typography';

import * as M from '../../utils/metadata';
import sourceThumb from '../../assets/images/livesource.png';
import Dialog from '../../misc/modals/Dialog';
import H from '../../utils/help';
import HLSControl from '../../misc/controls/HLS';
import LicenseControl from '../../misc/controls/License';
import LimitsControl from '../../misc/controls/Limits';
import MetadataControl from '../../misc/controls/Metadata';
import NotifyContext from '../../contexts/Notify';
import Paper from '../../misc/Paper';
import PaperHeader from '../../misc/PaperHeader';
import PaperFooter from '../../misc/PaperFooter';
import PaperThumb from '../../misc/PaperThumb';
import ProcessControl from '../../misc/controls/Process';
import Profile from './Profile';
import ProfileSummary from './ProfileSummary';
import RTMPControl from '../../misc/controls/RTMP';
import SnapshotControl from '../../misc/controls/Snapshot';
import SRTControl from '../../misc/controls/SRT';
import TabPanel from '../../misc/TabPanel';
import TabsVerticalGrid from '../../misc/TabsVerticalGrid';

const useStyles = makeStyles((theme) => ({
	wizardButtonElement: {
		display: 'flex',
		alignItems: 'left',
	},
	wizardButton: {
		marginLeft: '1em',
		padding: ' 0em 2em 0em 2em',
	},
	link: {
		color: theme.palette.common.white,
	},
	inlineIcon: {
		marginBottom: '-.2rem',
	},
}));

export default function Edit(props) {
	const classes = useStyles();
	const { i18n } = useLingui();
	const navigate = useNavigate();
	const { channelid: _channelid, tab: _tab } = useParams();
	const notify = React.useContext(NotifyContext);
	const [$tab, setTab] = React.useState(_tab ? _tab : 'general');
	const [$state, setState] = React.useState({
		editing: false,
		edit: '',
		complete: false,
		saving: false,
	});
	const [$data, setData] = React.useState(M.getDefaultIngestMetadata());
	const [$skills, setSkills] = React.useState({});
	const [$config, setConfig] = React.useState({});
	const [$process, setProcess] = React.useState({});
	const [$ready, setReady] = React.useState(false);
	const [$deleteDialog, setDeleteDialog] = React.useState(false);
	const [$editDialog, setEditDialog] = React.useState({
		open: false,
		target: '',
		what: '',
	});
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

		const proc = await props.restreamer.GetIngestProgress(_channelid);
		setProcess(proc);

		let metadata = await props.restreamer.GetIngestMetadata(_channelid);
		setData({
			...$data,
			...metadata,
		});

		const skills = await props.restreamer.Skills();
		setSkills(skills);

		const config = await props.restreamer.ConfigActive();
		setConfig(config);

		const complete = M.validateProfile(metadata.sources, metadata.profiles[0]);

		const state = {
			complete: complete,
		};

		if (metadata.sources.length === 0) {
			state.editing = true;
			state.edit = 'video';
		}

		setState({
			...$state,
			...state,
		});

		setReady(true);
	};

	const handleChangeTab = (event, value) => {
		setTab(value);
	};

	const handleSourceEditDialog = (target) => (what) => {
		if ($process.order === 'start') {
			setEditDialog({
				...$editDialog,
				open: true,
				target: target,
				what: what,
			});

			return;
		}

		if (target === 'wizard') {
			handleWizard();
		} else {
			handleSourceEdit(what);
		}
	};

	const handleSourceEditDialogAbort = () => {
		setEditDialog({
			...$editDialog,
			open: false,
			target: '',
			what: '',
		});
	};

	const handleSourceEditDialogDone = async () => {
		let stopped = false;

		stopped = await props.restreamer.StopIngest(_channelid);
		stopped = stopped ? await props.restreamer.StopIngestSnapshot(_channelid) : false;

		const target = $editDialog.target;
		const what = $editDialog.what;

		setEditDialog({
			...$editDialog,
			open: false,
			target: '',
			what: '',
		});

		if (stopped === false) {
			notify.Dispatch('error', 'edit:ingest', t`Failed to stop process`);
			return;
		}

		if (target === 'wizard') {
			handleWizard();
		} else {
			handleSourceEdit(what);
		}
	};

	const handleSourceEdit = (what) => {
		setState({
			...$state,
			editing: true,
			edit: what,
		});
	};

	const handleSkillsRefresh = async () => {
		await props.restreamer.RefreshSkills();

		const skills = await props.restreamer.Skills();
		setSkills(skills);
	};

	const handleSourceStore = async (name, data) => {
		return await props.restreamer.UploadData('', name, data);
	};

	const handleSourceProbe = async (inputs) => {
		let [res, err] = await props.restreamer.Probe(_channelid, inputs);
		if (err !== null) {
			res = {
				streams: [],
				log: [err.message],
			};
		}

		return res;
	};

	const handleSourceDone = (sources, profile) => {
		const complete = M.validateProfile(sources, profile);

		let streams = [];
		if (complete === true) {
			streams = M.createOutputStreams(sources, [profile]);
		}

		setData({
			...$data,
			sources: sources,
			profiles: [profile],
			streams: streams,
		});

		setState({
			...$state,
			editing: false,
			complete: complete,
		});
	};

	const handleSourceAbort = () => {
		setState({
			...$state,
			editing: false,
		});
	};

	const handleWizard = () => {
		navigate(`/${_channelid}/edit/wizard`);
	};

	const handleControlChange = (what) => (settings) => {
		const control = {
			...$data.control,
			[what]: settings,
		};

		setData({
			...$data,
			control: control,
		});
	};

	const handleMetadataChange = (settings) => {
		setData({
			...$data,
			meta: settings,
		});
	};

	const handleLicenseChange = (license) => {
		setData({
			...$data,
			license: license,
		});
	};

	const handleDone = async () => {
		setState({
			...$state,
			saving: true,
		});

		const save = async () => {
			const sources = $data.sources;
			const profiles = $data.profiles;
			const control = $data.control;

			const [global, inputs, outputs] = M.createInputsOutputs(sources, profiles, true);

			if (inputs.length === 0 || outputs.length === 0) {
				notify.Dispatch('error', 'save:ingest', i18n._(t`The input profile is not complete. Please define a video and audio source.`));
				return false;
			}

			// Create/update the ingest
			let [, err] = await props.restreamer.UpsertIngest(_channelid, global, inputs, outputs, control);
			if (err !== null) {
				notify.Dispatch('error', 'save:ingest', i18n._(t`Failed to update ingest process (${err.message})`));
				return false;
			}

			// Save the metadata
			let res = await props.restreamer.SetIngestMetadata(_channelid, $data);
			if (res === false) {
				notify.Dispatch('warning', 'save:ingest', i18n._(t`Failed to save ingest metadata`));
			}

			// Create/update the ingest snapshot process
			[, err] = await props.restreamer.UpsertIngestSnapshot(_channelid, control);
			if (err !== null) {
				notify.Dispatch('error', 'save:ingest', i18n._(t`Failed to update ingest snapshot process (${err.message})`));
			}

			// Create/update the player
			res = await props.restreamer.UpdatePlayer(_channelid);
			if (res === false) {
				notify.Dispatch('warning', 'save:ingest', i18n._(t`Failed to update the player`));
			}

			// Create/update the playersite
			res = await props.restreamer.UpdatePlayersite();
			if (res === false) {
				notify.Dispatch('warning', 'save:ingest', i18n._(t`Failed to update the playersite`));
			}

			return true;
		};

		const res = await save();

		setState({
			...$state,
			saving: false,
		});

		if (res === false) {
			return;
		}

		notify.Dispatch('success', 'save:ingest', i18n._(t`Channel "${$data.meta.name}" saved`));

		navigate(`/${_channelid}/`);
	};

	const handleAbort = () => {
		navigate(`/${_channelid}/`);
	};

	const handleChannelDeleteDialog = () => {
		setDeleteDialog(!$deleteDialog);
	};

	const handleChannelDelete = async () => {
		setState({
			...$state,
			saving: true,
		});

		const res = await props.restreamer.DeleteChannel(_channelid);
		if (res === false) {
			setState({
				...$state,
				saving: false,
			});
			notify.Dispatch('warning', 'delete:ingest', i18n._(t`The channel "${$data.meta.name}" could not be deleted`));
			return;
		}

		setState({
			...$state,
			saving: false,
		});

		notify.Dispatch('success', 'delete:ingest', i18n._(t`The channel "${$data.meta.name}" has been deleted`));

		navigate('/');
	};

	const handleHelp = () => {
		H('edit-' + $tab);
	};

	if ($ready === false) {
		return null;
	}

	let title = <Trans>Main Source</Trans>;
	if ($data.meta.name.length !== '') {
		title = $data.meta.name;
	}

	return (
		<React.Fragment>
			<Paper xs={12} md={10}>
				<PaperHeader
					title={
						<React.Fragment>
							<Trans>Edit</Trans>: {title}
						</React.Fragment>
					}
					onAbort={handleAbort}
					onHelp={handleHelp}
				/>
				<Grid container spacing={1}>
					<TabsVerticalGrid>
						<Tabs orientation="vertical" variant="scrollable" value={$tab} onChange={handleChangeTab}>
							<Tab className="tab" label={<Trans>General</Trans>} value="general" />
							<Tab className="tab" label={<Trans>Processing &amp; Control</Trans>} value="control" />
							<Tab className="tab" label={<Trans>Meta information</Trans>} value="meta" />
							<Tab className="tab" label={<Trans>License</Trans>} value="license" />
						</Tabs>
						<TabPanel value={$tab} index="general" className="panel">
							<Grid container spacing={2}>
								<Grid item xs={12}>
									<PaperThumb image={sourceThumb} title="General" />
								</Grid>
								<Grid item xs={12}>
									<Typography variant="h2">
										<Trans>General</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Typography variant="body1">
										<Trans>
											Edit the audio and video sources for the live stream. Add a description, and set your desired content license.
										</Trans>
									</Typography>
								</Grid>
								{$state.editing === false && (
									<Grid item xs={12}>
										<div className={classes.wizardButtonElement}>
											<Typography>
												<Trans>
													Use the wizard (<AutoFixHighIcon className={classes.inlineIcon} />) for a quick and easy setup, or edit (
													<EditIcon className={classes.inlineIcon} />) the sources directly in custom mode.
												</Trans>
											</Typography>
										</div>
									</Grid>
								)}
								<Grid item xs={12}>
									<Divider />
								</Grid>
							</Grid>
							{$state.editing === false ? (
								<ProfileSummary
									sources={$data.sources}
									profile={$data.profiles[0]}
									onWizard={handleSourceEditDialog('wizard')}
									onEdit={handleSourceEditDialog('extended')}
								/>
							) : (
								<Profile
									skills={$skills}
									sources={$data.sources}
									profile={$data.profiles[0]}
									config={$config.source}
									startWith={$state.edit}
									onProbe={handleSourceProbe}
									onRefresh={handleSkillsRefresh}
									onDone={handleSourceDone}
									onAbort={handleSourceAbort}
									onStore={handleSourceStore}
								/>
							)}
						</TabPanel>
						<TabPanel value={$tab} index="control" className="panel">
							<Grid container spacing={2}>
								<Grid item xs={12}>
									<Typography variant="h2">
										<Trans>Processing &amp; Control</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Divider />
								</Grid>
								<Grid item xs={12}>
									<Typography variant="h3">
										<Trans>HLS output</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<HLSControl settings={$data.control.hls} onChange={handleControlChange('hls')} />
								</Grid>
								<Grid item xs={12}>
									<Divider />
								</Grid>
								<Grid item xs={12}>
									<Grid container spacing={2}>
										<Grid item xs={12} md={6}>
											<Grid container spacing={2}>
												<Grid item xs={12}>
													<Typography variant="h3">
														<Trans>RTMP output</Trans>
													</Typography>
												</Grid>
												<Grid item xs={12}>
													<RTMPControl
														settings={$data.control.rtmp}
														enabled={$config.source.network.rtmp.enabled}
														onChange={handleControlChange('rtmp')}
													/>
												</Grid>
											</Grid>
										</Grid>
										<Grid item xs={12} display={{ xs: 'block', md: 'none' }}>
											<Divider />
										</Grid>
										<Grid item xs={12} md={6}>
											<Grid container spacing={2}>
												<Grid item xs={12}>
													<Typography variant="h3">
														<Trans>SRT output</Trans>
													</Typography>
												</Grid>
												<Grid item xs={12}>
													<SRTControl
														settings={$data.control.srt}
														enabled={$config.source.network.srt.enabled}
														onChange={handleControlChange('srt')}
													/>
												</Grid>
											</Grid>
										</Grid>
									</Grid>
								</Grid>
								<Grid item xs={12}>
									<Divider />
								</Grid>
								<Grid item xs={12}>
									<Typography variant="h3">
										<Trans>Snapshot</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<SnapshotControl settings={$data.control.snapshot} onChange={handleControlChange('snapshot')} />
								</Grid>
								<Grid item xs={12}>
									<Divider />
								</Grid>
								<Grid item xs={12}>
									<Typography variant="h3">
										<Trans>Process</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<ProcessControl settings={$data.control.process} onChange={handleControlChange('process')} />
								</Grid>
								<Grid item xs={12}>
									<Divider />
								</Grid>
								<Grid item xs={12}>
									<Typography variant="h3">
										<Trans>Limits</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<LimitsControl settings={$data.control.limits} onChange={handleControlChange('limits')} />
								</Grid>
							</Grid>
						</TabPanel>
						<TabPanel value={$tab} index="meta" className="panel">
							<Grid container spacing={2}>
								<Grid item xs={12}>
									<Typography variant="h2">
										<Trans>Metadata</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Typography>
										<Trans>Briefly describe what the audience will see during the live stream.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Divider />
								</Grid>
								<Grid item xs={12}>
									<MetadataControl settings={$data.meta} onChange={handleMetadataChange} />
								</Grid>
							</Grid>
						</TabPanel>
						<TabPanel value={$tab} index="license" className="panel">
							<Grid container spacing={2}>
								<Grid item xs={12}>
									<Typography variant="h2">
										<Trans>License</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Typography>
										<Trans>
											Use your copyright and choose the right image licence. Whether free for all or highly restricted. Briefly discuss
											what others are allowed to do with your image.
										</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Divider />
								</Grid>
								<Grid item xs={12}>
									<LicenseControl license={$data.license} onChange={handleLicenseChange} />
								</Grid>
							</Grid>
						</TabPanel>
					</TabsVerticalGrid>
				</Grid>
				<PaperFooter
					buttonsLeft={
						<Button variant="outlined" color="default" onClick={handleAbort}>
							<Trans>Close</Trans>
						</Button>
					}
					buttonsRight={
						<React.Fragment>
							<Button
								variant="outlined"
								color="primary"
								disabled={$state.editing === true || $state.complete === false || $state.saving === true}
								onClick={handleDone}
							>
								<Trans>Save</Trans>
							</Button>
							<Button variant="outlined" color="secondary" onClick={handleChannelDeleteDialog}>
								<Trans>Delete</Trans>
							</Button>
						</React.Fragment>
					}
				/>
			</Paper>
			<Dialog
				open={$editDialog.open}
				onClose={handleSourceEditDialogAbort}
				title={<Trans>Do you want to disconnect "{$data.meta.name}"?</Trans>}
				buttonsLeft={
					<Button variant="outlined" color="default" onClick={handleSourceEditDialogAbort}>
						<Trans>Abort</Trans>
					</Button>
				}
				buttonsRight={
					<Button variant="outlined" color="secondary" onClick={handleSourceEditDialogDone}>
						<Trans>Disconnect &amp; Continue</Trans>
					</Button>
				}
			>
				<Typography>
					<Trans>This source cannot be edited while it is in use. To continue, you have to disconnect the source.</Trans>
				</Typography>
			</Dialog>
			<Dialog
				open={$deleteDialog}
				onClose={handleChannelDeleteDialog}
				title={<Trans>Do you want to delete "{$data.meta.name}"?</Trans>}
				buttonsLeft={
					<Button variant="outlined" color="default" onClick={handleChannelDeleteDialog}>
						<Trans>Abort</Trans>
					</Button>
				}
				buttonsRight={
					<Button variant="outlined" color="secondary" onClick={handleChannelDelete}>
						<Trans>Delete</Trans>
					</Button>
				}
			>
				<Typography>
					<Trans>The deletion of this channel can not be recovered. All publications of this channel will be removed.</Trans>
				</Typography>
			</Dialog>
			<Backdrop open={$state.saving}>
				<CircularProgress color="inherit" />
			</Backdrop>
		</React.Fragment>
	);
}

Edit.defaultProps = {
	restreamer: null,
};

Edit.propTypes = {
	restreamer: PropTypes.object.isRequired,
};
