import React from 'react';

import { Trans } from '@lingui/macro';
import makeStyles from '@mui/styles/makeStyles';
import Grid from '@mui/material/Grid';
import Tab from '@mui/material/Tab';
import TextField from '@mui/material/TextField';

import TabPanel from '../TabPanel';
import TabsHorizontal from '../TabsHorizontal';

const useStyles = makeStyles((theme) => ({
	root: {
		'& .MuiOutlinedInput-input': {
			whiteSpace: 'pre-line',
		},
	},
}));

function init(settings) {
	const initSettings = {
		name: 'Livestream',
		description: 'Live from earth. Powered by datarhei Restreamer.',
		author: {},
		...settings,
	};

	initSettings.author = {
		name: '',
		description: '',
		...initSettings.author,
	};

	return initSettings;
}

export default function Control(props) {
	const classes = useStyles();
	const [$tab, setTab] = React.useState('content');
	const settings = init(props.settings);

	// Set the defaults
	React.useEffect(() => {
		props.onChange(settings, true);
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

	const handleChange = (what) => (event) => {
		const value = event.target.value;

		if (what === 'author.name') {
			settings.author.name = value;
		} else if (what === 'author.description') {
			settings.author.description = value;
		} else {
			settings[what] = value;
		}

		props.onChange(settings, false);
	};

	const handleChangeTab = (event, value) => {
		setTab(value);
	};

	return (
		<Grid container spacing={0}>
			<Grid item xs={12}>
				<TabsHorizontal value={$tab} onChange={handleChangeTab}>
					<Tab className="tab" label={<Trans>Content</Trans>} value="content" />
					<Tab className="tab" label={<Trans>Author</Trans>} value="author" />
				</TabsHorizontal>
				<TabPanel value={$tab} index="content">
					<Grid container spacing={2}>
						<Grid item xs={12}>
							<TextField variant="outlined" fullWidth label={<Trans>Name</Trans>} value={settings.name} onChange={handleChange('name')} />
						</Grid>
						<Grid item xs={12}>
							<TextField
								className={classes.root}
								variant="outlined"
								fullWidth
								multiline
								rows={10}
								label={<Trans>Description</Trans>}
								value={settings.description}
								onChange={handleChange('description')}
							/>
						</Grid>
					</Grid>
				</TabPanel>
				<TabPanel value={$tab} index="author">
					<Grid container spacing={2}>
						<Grid item xs={12}>
							<TextField
								variant="outlined"
								fullWidth
								label={<Trans>Name</Trans>}
								value={settings.author.name}
								onChange={handleChange('author.name')}
							/>
						</Grid>
						<Grid item xs={12}>
							<TextField
								className={classes.root}
								variant="outlined"
								fullWidth
								multiline
								rows={10}
								label={<Trans>Description</Trans>}
								value={settings.author.description}
								onChange={handleChange('author.description')}
							/>
						</Grid>
					</Grid>
				</TabPanel>
			</Grid>
		</Grid>
	);
}

Control.defaultProps = {
	settings: {},
	onChange: function (metadata) {},
};
