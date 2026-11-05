import React from 'react';
import { useNavigate } from 'react-router-dom';

import { useLingui } from '@lingui/react';
import { Trans, t } from '@lingui/macro';
import makeStyles from '@mui/styles/makeStyles';
import Backdrop from '@mui/material/Backdrop';
import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import Divider from '@mui/material/Divider';
import Grid from '@mui/material/Grid';
import MenuItem from '@mui/material/MenuItem';
import Tab from '@mui/material/Tab';
import Tabs from '@mui/material/Tabs';
import TextField from '@mui/material/TextField';
import Typography from '@mui/material/Typography';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';

import * as M from '../utils/metadata';
import playerSiteThumb from '../assets/images/playersite.png';
import Checkbox from '../misc/Checkbox';
import ColorPicker from '../misc/ColorPicker';
import Dialog from '../misc/modals/Dialog';
import Filesize from '../misc/Filesize';
import FormInlineButton from '../misc/FormInlineButton';
import H from '../utils/help';
import MultiSelect from '../misc/MultiSelect';
import MultiSelectOption from '../misc/MultiSelectOption';
import NotifyContext from '../contexts/Notify';
import Paper from '../misc/Paper';
import PaperHeader from '../misc/PaperHeader';
import PaperFooter from '../misc/PaperFooter';
import PaperThumb from '../misc/PaperThumb';
import Select from '../misc/Select';
import TabPanel from '../misc/TabPanel';
import TabsVerticalGrid from '../misc/TabsVerticalGrid';
import UploadButton from '../misc/UploadButton';

const useStyles = makeStyles((theme) => ({
	buttonOpen: {
		float: 'right',
		marginLeft: '.5em',
	},
}));

const imageTypes = [
	{ mimetype: 'image/gif', extension: 'gif', maxSize: 1 * 1024 * 1024 },
	{ mimetype: 'image/png', extension: 'png', maxSize: 1 * 1024 * 1024 },
	{ mimetype: 'image/jpeg', extension: 'jpg', maxSize: 1 * 1024 * 1024 },
	{ mimetype: 'image/svg+xml', extension: 'svg', maxSize: 1 * 1024 * 1024 },
];

const templateTypes = [{ mimetype: 'text/html', extension: 'html', maxSize: 500 * 1024 }];

function CheckboxList(props) {
	const [$checked, setChecked] = React.useState(props.selected);

	const handleToggle = (value) => () => {
		const currentIndex = $checked.indexOf(value);
		const newChecked = [...$checked];

		if (currentIndex === -1) {
			newChecked.push(value);
		} else {
			newChecked.splice(currentIndex, 1);
		}

		setChecked(newChecked);

		props.onChange({
			target: {
				value: [...newChecked],
			},
		});
	};

	return (
		<List sx={{ width: '100%', position: 'relative', overflow: 'auto', maxHeight: 300 }}>
			{props.items.map((value) => {
				const labelId = `checkbox-list-label-${value.id}`;

				return (
					<ListItem key={value.id} disablePadding>
						<ListItemButton onClick={handleToggle(value.id)} dense disabled={props.disabled}>
							<ListItemIcon>
								<Checkbox
									edge="start"
									checked={$checked.indexOf(value.id) !== -1}
									tabIndex={-1}
									disableRipple
									inputProps={{ 'aria-labelledby': labelId }}
								/>
							</ListItemIcon>
							<ListItemText id={labelId} primary={value.name} />
						</ListItemButton>
					</ListItem>
				);
			})}
		</List>
	);
}

CheckboxList.defaultProps = {
	disabled: false,
	selected: [],
	items: [],
	onChange: function () {},
};

export default function Playersite(props) {
	const classes = useStyles();
	const navigate = useNavigate();
	const { i18n } = useLingui();
	const address = props.restreamer.Address() + '/';
	const playersiteUrl = props.restreamer.GetPlayersiteUrl();
	const notify = React.useContext(NotifyContext);
	const [$ready, setReady] = React.useState(false);
	const [$ingest, setIngest] = React.useState(false);
	const [$available, setAvailable] = React.useState(false);
	const [$data, setData] = React.useState(M.getDefaultMetadata());
	const [$settings, setSettings] = React.useState({});
	const [$channels, setChannels] = React.useState([]);
	const [$templates, setTemplates] = React.useState([]);
	const [$tab, setTab] = React.useState('general');
	const [$saving, setSaving] = React.useState(false);
	const [$error, setError] = React.useState({
		open: false,
		title: '',
		message: '',
	});

	React.useEffect(() => {
		(async () => {
			await mount();
		})();
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

	const mount = async () => {
		const data = await props.restreamer.GetMetadata();

		setData(data);
		setSettings(props.restreamer.InitPlayersiteSettings(data.playersite));
		setChannels(props.restreamer.ListChannels());
		setIngest(props.restreamer.HasIngest());

		setTemplates(await props.restreamer.ListPlayersiteTemplates());
		setAvailable(await props.restreamer.HasPlayersite());

		setReady(true);
	};

	const handleChange = (what) => (event) => {
		const value = event.target.value;
		const settings = $settings;

		if (['playersite', 'header', 'share', 'support', 'chromecast', 'airplay'].includes(what)) {
			settings[what] = !settings[what];
		} else {
			settings[what] = value;
		}

		setSettings({
			...$settings,
			...settings,
		});
	};

	const handleBackgroundImageUpload = async (data, extension) => {
		const path = await props.restreamer.UploadPlayersiteBackgroundImage(data, extension);

		handleChange('bgimage_url')({
			target: {
				value: path,
			},
		});

		setSaving(false);
	};

	const handleTemplateUpload = async (data, extension) => {
		const name = await props.restreamer.UploadPlayersiteTemplate(data, $settings.templatename);

		setTemplates(await props.restreamer.ListPlayersiteTemplates());
		setSettings({
			...$settings,
			template: name,
			templatename: '',
		});

		setSaving(false);
	};

	const handleTemplateDelete = async () => {
		setSaving(true);

		await props.restreamer.DeletePlayersiteTemplate($settings.template);
		setSettings({
			...$settings,
			template: '!default',
		});
		setTemplates(await props.restreamer.ListPlayersiteTemplates());

		setSaving(false);
	};

	const handleUploadStart = () => {
		setSaving(true);
	};

	const handleUploadError = (title) => (err) => {
		let message = null;

		switch (err.type) {
			case 'nofiles':
				message = <Trans>Please select a file to upload.</Trans>;
				break;
			case 'mimetype':
				message = (
					<Trans>
						The selected file type ({err.actual}) is not allowed. Allowed file types are {err.allowed.join(', ')}
					</Trans>
				);
				break;
			case 'size':
				message = (
					<Trans>
						The selected file is too big (<Filesize bytes={err.actual} />
						). Only <Filesize bytes={err.allowed} /> are allowed.
					</Trans>
				);
				break;
			case 'read':
				message = <Trans>There was an error during upload: {err.message}</Trans>;
				break;
			default:
				message = <Trans>Unknown upload error</Trans>;
		}

		setSaving(false);

		showUploadError(title, message);
	};

	const showUploadError = (title, message) => {
		setError({
			...$error,
			open: true,
			title: title,
			message: message,
		});
	};

	const hideUploadError = () => {
		setError({
			...$error,
			open: false,
		});
	};

	const handleDone = async () => {
		setSaving(true);

		const data = {
			...$data,
			playersite: $settings,
		};

		let res = await props.restreamer.SetMetadata(data);
		if (res === false) {
			notify.Dispatch('error', 'save:playersite', i18n._(t`Failed to store player size setting.`));
			setSaving(false);
			return;
		}

		res = await props.restreamer.UpdatePlayersite();
		if (res === false) {
			notify.Dispatch('error', 'save:playersite', i18n._(t`Failed to create publication website files.`));
			setSaving(false);
			return;
		}

		setAvailable(await props.restreamer.HasPlayersite());

		setSaving(false);

		notify.Dispatch('success', 'save:playersite', i18n._(t`Publication website settings saved`));
	};

	const handleChangeTab = (event, value) => {
		setTab(value);
	};

	const handleAbort = () => {
		navigate('/');
	};

	const handleHelp = () => {
		H('playersite-' + $tab);
	};

	if ($ready === false) {
		return null;
	}

	if ($ingest === false) {
		navigate('/');
		return null;
	}

	let main_channelid = props.restreamer.GetCurrentChannelID();
	let channel = props.restreamer.GetChannel($settings.channelid);
	if (channel !== null) {
		main_channelid = channel.channelid;
	}

	return (
		<React.Fragment>
			<Paper xs={12} md={10}>
				<PaperHeader title={<Trans>EDIT: Publication Website</Trans>} onAbort={handleAbort} onHelp={handleHelp} />
				<Grid container spacing={2}>
					<TabsVerticalGrid>
						<Tabs orientation="vertical" variant="scrollable" value={$tab} onChange={handleChangeTab} className="tabs">
							<Tab className="tab" label={<Trans>General</Trans>} value="general" />
							<Tab className="tab" label={<Trans>Template</Trans>} value="template" disabled={!$settings.playersite} />
							<Tab className="tab" label={<Trans>Design</Trans>} value="design" disabled={!$settings.playersite} />
							<Tab className="tab" label={<Trans>Notes</Trans>} value="notes" disabled={!$settings.playersite} />
							<Tab className="tab" label={<Trans>Code injection</Trans>} value="code_injection" disabled={!$settings.playersite} />
						</Tabs>
						<TabPanel value={$tab} index="general">
							<Grid container spacing={2}>
								<Grid item xs={12}>
									<PaperThumb image={playerSiteThumb} title="Playersite" />
								</Grid>
								<Grid item xs={12}>
									<Typography variant="h1">
										<Trans>General</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Typography variant="body1">
										<Trans>
											In addition to the player, the Restreamer offers a complete landingpage, which you can use to present your live
											stream easily and quickly.
										</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Checkbox label={<Trans>Playersite</Trans>} checked={$settings.playersite} onChange={handleChange('playersite')} />
								</Grid>
								<Grid item xs={12}>
									<TextField
										variant="outlined"
										fullWidth
										label={<Trans>Sitename</Trans>}
										disabled={!$settings.playersite}
										value={$settings.title}
										onChange={handleChange('title')}
									/>
								</Grid>
								<Grid item xs={12}>
									<Select
										label={<Trans>Main channel</Trans>}
										value={$settings.channelid}
										disabled={!$settings.playersite}
										onChange={handleChange('channelid')}
									>
										<MenuItem value="current">
											<Trans>Selected channel</Trans>
										</MenuItem>
										{$channels
											.sort((a, b) => {
												const aname = a.name.toUpperCase();
												const bname = b.name.toUpperCase();
												return aname < bname ? -1 : aname > bname ? 1 : 0;
											})
											.map((c) => {
												return (
													<MenuItem key={c.channelid} value={c.channelid}>
														{c.name}
													</MenuItem>
												);
											})}
									</Select>
									<Typography variant="caption">
										<Trans>Main page channel (index.html).</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<MultiSelect
										disabled={!$settings.playersite}
										type="select"
										label={<Trans>Additional channels</Trans>}
										value={$settings.channel_list.filter((c) => c !== main_channelid)}
										renderValue={(selected) => {
											return selected
												.map((id) => {
													let channel = props.restreamer.GetChannel(id);
													if (channel === null) {
														return '';
													}
													return channel.name;
												})
												.filter((name) => name.length !== 0)
												.join(', ');
										}}
										onChange={handleChange('channel_list')}
									>
										{$channels
											.sort((a, b) => {
												const aname = a.name.toUpperCase();
												const bname = b.name.toUpperCase();
												return aname < bname ? -1 : aname > bname ? 1 : 0;
											})
											.map((c) => {
												return { id: c.channelid, name: c.name };
											})
											.filter((c) => c.id !== main_channelid)
											.map((c) => {
												return <MultiSelectOption key={c.id} value={c.id} name={c.name} />;
											})}
									</MultiSelect>
									<Typography variant="caption">
										<Trans>Additional channels to display on the playersite.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Checkbox
										label={<Trans>Share button</Trans>}
										checked={$settings.share}
										disabled={!$settings.playersite}
										onChange={handleChange('share')}
									/>
								</Grid>
								<Grid item xs={12}>
									<Checkbox
										label={<Trans>Chromecast</Trans>}
										checked={$settings.chromecast}
										disabled={!$settings.playersite}
										onChange={handleChange('chromecast')}
									/>
								</Grid>
								<Grid item xs={12}>
									<Checkbox
										label={<Trans>AirPlay</Trans>}
										checked={$settings.airplay}
										disabled={!$settings.playersite}
										onChange={handleChange('airplay')}
									/>
								</Grid>
								<Grid item xs={12}>
									<Checkbox
										label={<Trans>Support datarhei Restreamer</Trans>}
										checked={$settings.support}
										disabled={!$settings.playersite}
										onChange={handleChange('support')}
									/>
									<Typography variant="caption">
										<Trans>Shows a reference to the project.</Trans>
									</Typography>
								</Grid>
							</Grid>
						</TabPanel>
						<TabPanel value={$tab} index="template">
							<Grid container spacing={2}>
								<Grid item xs={12}>
									<Typography variant="h1">
										<Trans>Template</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Typography variant="body1">
										<Trans>Available</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12} md={9}>
									<Select label={<Trans>Selection</Trans>} value={$settings.template} onChange={handleChange('template')}>
										<MenuItem value="!default">
											<Trans>Default</Trans>
										</MenuItem>
										{$templates.sort().map((a) => {
											return (
												<MenuItem key={a} value={a}>
													{a}
												</MenuItem>
											);
										})}
									</Select>
									<Typography variant="caption">
										<Trans>
											Template to be used for creating the publication website. The delete button removes the selection from the system.
										</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12} md={3}>
									<FormInlineButton
										variant="outlined"
										color="secondary"
										disabled={$settings.template === '!default'}
										onClick={handleTemplateDelete}
									>
										<Trans>Delete</Trans>
									</FormInlineButton>
								</Grid>
								<Grid item xs={12}>
									<Divider />
								</Grid>
								<Grid item xs={12}>
									<Typography variant="body1">
										<Trans>Upload</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12} md={9}>
									<TextField
										variant="outlined"
										fullWidth
										label={<Trans>Name</Trans>}
										value={$settings.templatename}
										onChange={handleChange('templatename')}
									/>
									<Typography variant="caption">
										<Trans>Name for the template. If the name already exists, it will be overwritten.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12} md={3}>
									<UploadButton
										variant="outlined"
										color="primary"
										disabled={$settings.templatename.length === 0}
										label={<Trans>Upload</Trans>}
										acceptTypes={templateTypes}
										onStart={handleUploadStart}
										onError={handleUploadError(<Trans>Uploading the file failed</Trans>)}
										onUpload={handleTemplateUpload}
									/>
								</Grid>
							</Grid>
						</TabPanel>
						<TabPanel value={$tab} index="design">
							<Grid container spacing={2}>
								<Grid item xs={12}>
									<Typography variant="h1">
										<Trans>Design</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Typography variant="body1">
										<Trans>Adjust publication site colors and background as you like.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Typography variant="h3">
										<Trans>Text colors</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<ColorPicker
										variant="outlined"
										fullWidth
										label={<Trans>Default</Trans>}
										value={$settings.textcolor_default}
										onChange={handleChange('textcolor_default')}
									/>
								</Grid>
								<Grid item xs={12}>
									<ColorPicker
										variant="outlined"
										fullWidth
										label={<Trans>Headline</Trans>}
										value={$settings.textcolor_title}
										onChange={handleChange('textcolor_title')}
									/>
								</Grid>
								<Grid item xs={6}>
									<ColorPicker
										variant="outlined"
										fullWidth
										label={<Trans>Link</Trans>}
										value={$settings.textcolor_link}
										onChange={handleChange('textcolor_link')}
									/>
								</Grid>
								<Grid item xs={6}>
									<ColorPicker
										variant="outlined"
										fullWidth
										label={<Trans>Link, mouseover</Trans>}
										value={$settings.textcolor_link_hover}
										onChange={handleChange('textcolor_link_hover')}
									/>
								</Grid>
								<Grid item xs={12}>
									<Divider />
								</Grid>
								<Grid item xs={12}>
									<Typography variant="h3">
										<Trans>Background colors</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<ColorPicker
										variant="outlined"
										fullWidth
										label={<Trans>Default</Trans>}
										value={$settings.bgcolor_default}
										onChange={handleChange('bgcolor_default')}
									/>
								</Grid>
								<Grid item xs={12}>
									<ColorPicker
										variant="outlined"
										fullWidth
										label={<Trans>Header</Trans>}
										value={$settings.bgcolor_header}
										onChange={handleChange('bgcolor_header')}
									/>
								</Grid>
								<Grid item xs={12}>
									<ColorPicker
										variant="outlined"
										fullWidth
										label={<Trans>Selected</Trans>}
										value={$settings.bgcolor_selected}
										onChange={handleChange('bgcolor_selected')}
									/>
								</Grid>
								<Grid item xs={12}>
									<ColorPicker
										variant="outlined"
										fullWidth
										label={<Trans>Unselected</Trans>}
										value={$settings.bgcolor_unselected}
										onChange={handleChange('bgcolor_unselected')}
									/>
								</Grid>
								<Grid item xs={12}>
									<ColorPicker
										variant="outlined"
										fullWidth
										label={<Trans>Linecolor</Trans>}
										value={$settings.hrcolor}
										onChange={handleChange('hrcolor')}
									/>
								</Grid>
								<Grid item xs={12}>
									<Divider />
								</Grid>
								<Grid item xs={12}>
									<Typography variant="h3">
										<Trans>Background image</Trans>
									</Typography>
								</Grid>

								<Grid item xs={12} md={8}>
									<TextField
										variant="outlined"
										fullWidth
										label={<Trans>Image URL</Trans>}
										value={$settings.bgimage_url}
										onChange={handleChange('bgimage_url')}
									/>
									<Typography variant="caption">
										<Trans>Address for the background image.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12} md={4}>
									<UploadButton
										label={<Trans>Upload</Trans>}
										acceptTypes={imageTypes}
										onStart={handleUploadStart}
										onError={handleUploadError(<Trans>Uploading the file failed</Trans>)}
										onUpload={handleBackgroundImageUpload}
									/>
								</Grid>
							</Grid>
						</TabPanel>
						<TabPanel value={$tab} index="notes">
							<Grid container spacing={2}>
								<Grid item xs={12}>
									<Typography variant="h2">
										<Trans>Notes</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Typography variant="h3">
										<Trans>Imprint</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<TextField
										className={classes.root}
										variant="outlined"
										fullWidth
										multiline
										rows={10}
										label={<Trans>Content</Trans>}
										value={$settings.imprint}
										onChange={handleChange('imprint')}
									/>
								</Grid>
								<Grid item xs={12}>
									<Typography variant="h3">
										<Trans>Terms</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<TextField
										className={classes.root}
										variant="outlined"
										fullWidth
										multiline
										rows={10}
										label={<Trans>Content</Trans>}
										value={$settings.terms}
										onChange={handleChange('terms')}
									/>
								</Grid>
							</Grid>
						</TabPanel>
						<TabPanel value={$tab} index="code_injection">
							<Grid container spacing={2}>
								<Grid item xs={12}>
									<Typography variant="h2">
										<Trans>Custom code injection</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Typography variant="body1">
										<Trans>Add external widgets and styles to the publication site. You can find some examples on the help page.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Typography variant="h3">
										<Trans>Extend header</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<TextField
										className={classes.root}
										variant="outlined"
										fullWidth
										multiline
										rows={10}
										label={<Trans>Inject 1</Trans>}
										value={$settings.inject1}
										onChange={handleChange('inject1')}
									/>
									<Typography variant="caption">
										<Trans>For Stylesheets.</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Typography variant="h3">
										<Trans>Extend channel list</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<TextField
										className={classes.root}
										variant="outlined"
										fullWidth
										multiline
										rows={10}
										label={<Trans>Inject 2</Trans>}
										value={$settings.inject2}
										onChange={handleChange('inject2')}
									/>
									<Typography variant="caption">
										<Trans>Expands the area above the channel list (live chat).</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Typography variant="h3">
										<Trans>Extend content</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<TextField
										className={classes.root}
										variant="outlined"
										fullWidth
										multiline
										rows={10}
										label={<Trans>Inject 3</Trans>}
										value={$settings.inject3}
										onChange={handleChange('inject3')}
									/>
									<Typography variant="caption">
										<Trans>Expands the area under the channel description (comment boxes).</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<Typography variant="h3">
										<Trans>Extend footer</Trans>
									</Typography>
								</Grid>
								<Grid item xs={12}>
									<TextField
										className={classes.root}
										variant="outlined"
										fullWidth
										multiline
										rows={10}
										label={<Trans>Inject 4</Trans>}
										value={$settings.inject4}
										onChange={handleChange('inject4')}
									/>
									<Typography variant="caption">
										<Trans>For Javascripts.</Trans>
									</Typography>
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
							<Button variant="outlined" color="primary" onClick={handleDone}>
								<Trans>Save</Trans>
							</Button>
							<Button
								variant="outlined"
								color="primary"
								className={classes.buttonOpen}
								href={address + playersiteUrl}
								target="blank"
								disabled={!$available}
							>
								<Trans>Open</Trans>
							</Button>
						</React.Fragment>
					}
				/>
			</Paper>
			<Backdrop open={$saving}>
				<CircularProgress color="inherit" />
			</Backdrop>
			<Dialog
				open={$error.open}
				title={$error.title}
				onClose={hideUploadError}
				buttonsRight={
					<Button variant="outlined" color="primary" onClick={hideUploadError}>
						<Trans>OK</Trans>
					</Button>
				}
			>
				<Typography variant="body1">{$error.message}</Typography>
			</Dialog>
		</React.Fragment>
	);
}

Playersite.defaultProps = {
	restreamer: null,
};
