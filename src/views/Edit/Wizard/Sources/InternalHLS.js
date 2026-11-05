import React from 'react';

import { Trans } from '@lingui/macro';
import Grid from '@mui/material/Grid';
import Icon from '@mui/icons-material/KeyboardTab';
import Typography from '@mui/material/Typography';

import * as S from '../../Sources/Network';
import BoxTextarea from '../../../../misc/BoxTextarea';
import Textarea from '../../../../misc/Textarea';

const initSettings = (initialSettings, config) => {
	const settings = {
		...S.func.initSettings(initialSettings, config),
		mode: 'push',
	};

	settings.push.type = 'hls';

	return settings;
};

function Source(props) {
	const config = S.func.initConfig(props.config);
	const settings = initSettings(props.settings, config);
	const skills = S.func.initSkills(props.skills);

	const handleChange = (newSettings) => {
		newSettings = newSettings || settings;

		const inputs = S.func.createInputs(newSettings, config, skills);
		newSettings.address = inputs[0].address;

		props.onChange(S.id, newSettings, inputs, newSettings.address.length !== 0);
	};

	React.useEffect(() => {
		handleChange();
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

	const HLS = S.func.getHLS(config, settings.push.name);

	return (
		<React.Fragment>
			<Grid item xs={12}>
				<Typography>
					<Trans>Send stream to this address:</Trans>
				</Typography>
			</Grid>
			<Grid item xs={12}>
				<BoxTextarea>
					<Textarea rows={1} value={HLS} readOnly allowCopy />
				</BoxTextarea>
			</Grid>
		</React.Fragment>
	);
}

Source.defaultProps = {
	settings: {},
	config: null,
	skills: null,
	onChange: function (type, settings, inputs, ready) {},
};

function SourceIcon(props) {
	return <Icon style={{ color: '#FFF' }} {...props} />;
}

const id = 'hls';
const type = 'network';
const name = <Trans>HLS server</Trans>;
const capabilities = ['audio', 'video'];

export { id, type, name, capabilities, SourceIcon as icon, Source as component };
