import React from 'react';

import { useLingui } from '@lingui/react';
import { Trans, t } from '@lingui/macro';
import Grid from '@mui/material/Grid';
import MenuItem from '@mui/material/MenuItem';

import Select from '../../misc/Select';

export default function StreamSelect(props) {
	const { i18n } = useLingui();
	let selected = props.selected;

	const handleChange = (event) => {
		const stream = parseInt(event.target.value);

		props.onChange(stream);
	};

	let streamList = [];

	for (let s of props.streams) {
		if (s.type !== props.type) {
			continue;
		}

		if (s.type === 'video') {
			streamList.push(
				<MenuItem value={s.stream} key={s.stream}>
					{i18n._(t`Stream`)} {s.stream}: {s.codec.toUpperCase()} {s.width}x{s.height} {s.pix_fmt}
				</MenuItem>
			);
		} else if (s.type === 'audio') {
			streamList.push(
				<MenuItem value={s.stream} key={s.stream}>
					{i18n._(t`Stream`)} {s.stream}: {s.codec.toUpperCase()} {s.layout} {s.sampling_hz}Hz
				</MenuItem>
			);
		}
	}

	if (props.type === 'video') {
		if (streamList.length === 0) {
			selected = -1;
			streamList.push(
				<MenuItem value="-1" key="unavailable" disabled>
					{i18n._(t`No video stream available`)}
				</MenuItem>
			);
		} else {
			streamList.unshift(
				<MenuItem value="-1" key="none" disabled>
					{i18n._(t`Choose a video stream`)}
				</MenuItem>
			);
		}
	} else if (props.type === 'audio') {
		if (streamList.length === 0 && selected >= 0) {
			selected = -1;
			streamList.push(
				<MenuItem value="-1" key="unavailable" disabled>
					{i18n._(t`No audio stream available`)}
				</MenuItem>
			);
		} else {
			streamList.unshift(
				<MenuItem value="-1" key="none" disabled>
					{i18n._(t`Choose an audio stream`)}
				</MenuItem>
			);
		}
	}

	if (streamList.length === 0) {
		return null;
	}

	return (
		<Grid container spacing={2}>
			<Grid item xs={12}>
				<Select label={<Trans>Stream</Trans>} value={selected} fullWidth onChange={handleChange}>
					{streamList}
				</Select>
			</Grid>
		</Grid>
	);
}

StreamSelect.defaultProps = {
	type: '',
	streams: [],
	selected: -1,
	allowCustom: false,
	allowNone: false,
	onChange: function (stream) {},
};
