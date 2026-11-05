import React from 'react';

import { useLingui } from '@lingui/react';
import { Trans, t } from '@lingui/macro';
import MenuItem from '@mui/material/MenuItem';

import Select from '../../../misc/Select';
import SelectCustom from '../../../misc/SelectCustom';

function Bitrate(props) {
	const { i18n } = useLingui();
	const bitrates = [
		{ value: '32768', label: '32768 kbit/s' },
		{ value: '24576', label: '24576 kbit/s' },
		{ value: '20480', label: '20480 kbit/s' },
		{ value: '16384', label: '16384 kbit/s' },
		{ value: '12288', label: '12288 kbit/s' },
		{ value: '8192', label: '8192 kbit/s' },
		{ value: '4096', label: '4096 kbit/s' },
		{ value: '2048', label: '2048 kbit/s' },
		{ value: '1024', label: '1024 kbit/s' },
		{ value: '512', label: '512 kbit/s' },
		{ value: '256', label: '256 kbit/s' },
	];

	if (props.allowAuto === true) {
		bitrates.unshift({ value: 'auto', label: 'auto' });
	}

	if (props.allowCustom === true) {
		bitrates.push({ value: 'custom', label: i18n._(t`Custom ...`) });
	}

	return (
		<SelectCustom
			options={bitrates}
			label={props.label}
			customLabel={props.customLabel}
			value={props.value}
			onChange={props.onChange}
			variant={props.variant}
			allowCustom={props.allowCustom}
		/>
	);
}

Bitrate.defaultProps = {
	allowAuto: false,
	allowCustom: false,
	variant: 'outlined',
	label: <Trans>Bitrate</Trans>,
	customLabel: <Trans>Custom bitrate (kbit/s)</Trans>,
	onChange: function (event) {},
};

function GOP(props) {
	const { i18n } = useLingui();
	const bitrates = [
		{ value: '1', label: '1' },
		{ value: '2', label: '2' },
		{ value: '4', label: '4' },
		{ value: '10', label: '10' },
	];

	if (props.allowAuto === true) {
		bitrates.unshift({ value: 'auto', label: 'auto' });
	}

	if (props.allowCustom === true) {
		bitrates.push({ value: 'custom', label: i18n._(t`Custom ...`) });
	}

	return (
		<SelectCustom
			options={bitrates}
			label={props.label}
			customLabel={props.customLabel}
			value={props.value}
			onChange={props.onChange}
			variant={props.variant}
			allowCustom={props.allowCustom}
		/>
	);
}

GOP.defaultProps = {
	allowAuto: false,
	allowCustom: false,
	variant: 'outlined',
	label: <Trans>Keyframe interval (seconds)</Trans>,
	customLabel: <Trans>Custom keyframe interval</Trans>,
	onChange: function (event) {},
};

function Framerate(props) {
	const { i18n } = useLingui();
	const sizes = [
		{ value: '60', label: '60' },
		{ value: '59.94', label: '59.94' },
		{ value: '50', label: '50' },
		{ value: '30', label: '30' },
		{ value: '29.97', label: '29.97 (NTSC)' },
		{ value: '25', label: '25 (PAL)' },
		{ value: '24', label: '24 (Film)' },
		{ value: '23.97', label: '23.97 (NTSC Film)' },
		{ value: '15', label: '15' },
		{ value: '10', label: '10' },
	];

	if (props.allowAuto === true) {
		sizes.unshift({ value: 'auto', label: 'auto' });
	}

	if (props.allowCustom === true) {
		sizes.push({ value: 'custom', label: i18n._(t`Custom ...`) });
	}

	return (
		<SelectCustom
			options={sizes}
			label={props.label}
			customLabel={props.customLabel}
			value={props.value}
			onChange={props.onChange}
			variant={props.variant}
			allowCustom={props.allowCustom}
		/>
	);
}

Framerate.defaultProps = {
	allowAuto: false,
	allowCustom: false,
	variant: 'outlined',
	label: <Trans>Framerate</Trans>,
	customLabel: <Trans>Custom framerate</Trans>,
	onChange: function (event) {},
};

function Profile(props) {
	return (
		<Select label={<Trans>Profile</Trans>} value={props.value} onChange={props.onChange}>
			<MenuItem value="auto">auto</MenuItem>
			<MenuItem value="baseline">baseline</MenuItem>
			<MenuItem value="main">main</MenuItem>
			<MenuItem value="high">high</MenuItem>
		</Select>
	);
}

Profile.defaultProps = {
	value: '',
	onChange: function (event) {},
};

// https://en.wikipedia.org/wiki/Graphics_display_resolution
function Size(props) {
	const { i18n } = useLingui();
	const sizes = [
		{ value: '7680x4320', label: '7680x4320 (8K UHD)' },
		{ value: '5120x2880', label: '5120x2880 (5K)' },
		{ value: '4096x2160', label: '4096x2160 (DCI 4K)' },
		{ value: '3840x2160', label: '3840x2160 (4K UHD)' },
		{ value: '3200x1800', label: '3200x1800 (QHD+)' },
		{ value: '2560x1600', label: '2560x1600 (WQXGA)' },
		{ value: '2560x1440', label: '2560x1440 (QHD, WQHD)' },
		{ value: '2048x1080', label: '2048x1080 (DCI 2K)' },
		{ value: '1920x1080', label: '1920x1080 (Full HD)' },
		{ value: '1600x900', label: '1600x900 (HD+)' },
		{ value: '1280x1080', label: '1280x1080' },
		{ value: '1280x720', label: '1280x720 (HD)' },
		{ value: '960x540', label: '960x540 (qHD)' },
		{ value: '640x360', label: '640x360 (nHD)' },
	];

	if (props.allowAuto === true) {
		sizes.unshift({ value: 'auto', label: 'auto' });
	}

	if (props.allowCustom === true) {
		sizes.push({ value: 'custom', label: i18n._(t`Custom ...`) });
	}

	return (
		<SelectCustom
			options={sizes}
			label={props.label}
			customLabel={props.customLabel}
			value={props.value}
			onChange={props.onChange}
			variant={props.variant}
			allowCustom={props.allowCustom}
		/>
	);
}

Size.defaultProps = {
	allowAuto: false,
	allowCustom: false,
	variant: 'outlined',
	label: <Trans>Size</Trans>,
	customLabel: <Trans>Custom size</Trans>,
	onChange: function (event) {},
};

function Height(props) {
	const { i18n } = useLingui();
	const height = [
		{ value: '4320', label: '4320' },
		{ value: '2880', label: '2880' },
		{ value: '2160', label: '2160' },
		{ value: '1800', label: '1800' },
		{ value: '1600', label: '1600' },
		{ value: '1440', label: '1440' },
		{ value: '1080', label: '1080' },
		{ value: '900', label: '900' },
		{ value: '720', label: '720' },
		{ value: '540', label: '540' },
		{ value: '360', label: '360' },
	];

	if (props.allowCustom === true) {
		height.push({ value: 'custom', label: i18n._(t`Custom ...`) });
	}

	return (
		<SelectCustom
			options={height}
			label={props.label}
			customLabel={props.customLabel}
			value={props.value}
			onChange={props.onChange}
			variant={props.variant}
			allowCustom={props.allowCustom}
		/>
	);
}

Height.defaultProps = {
	allowNone: false,
	allowCustom: false,
	variant: 'outlined',
	label: <Trans>Height</Trans>,
	customLabel: <Trans>Custom size</Trans>,
	onChange: function (event) {},
};

function Width(props) {
	const { i18n } = useLingui();
	const width = [
		{ value: '7680', label: '7680' },
		{ value: '5120', label: '5120' },
		{ value: '4096', label: '4096' },
		{ value: '3840', label: '3840' },
		{ value: '3200', label: '3200' },
		{ value: '2560', label: '2560' },
		{ value: '2048', label: '2048' },
		{ value: '1920', label: '1920' },
		{ value: '1600', label: '1600' },
		{ value: '1280', label: '1280' },
		{ value: '960', label: '960' },
		{ value: '640', label: '640' },
	];

	if (props.allowCustom === true) {
		width.push({ value: 'custom', label: i18n._(t`Custom ...`) });
	}

	return (
		<SelectCustom
			options={width}
			label={props.label}
			customLabel={props.customLabel}
			value={props.value}
			onChange={props.onChange}
			variant={props.variant}
			allowCustom={props.allowCustom}
		/>
	);
}

Width.defaultProps = {
	allowNone: false,
	allowCustom: false,
	variant: 'outlined',
	label: <Trans>Width</Trans>,
	customLabel: <Trans>Custom size</Trans>,
	onChange: function (event) {},
};

function Format(props) {
	const { i18n } = useLingui();
	const sizes = [
		{ value: 'yuv420p', label: 'yuv420p' },
		{ value: 'nv12', label: 'nv12' },
		{ value: 'h264', label: 'H.264' },
		{ value: 'mjpeg', label: 'MJPEG' },
	];

	if (props.allowAuto === true) {
		sizes.unshift({ value: 'auto', label: 'auto' });
	}

	if (props.allowCustom === true) {
		sizes.push({ value: 'custom', label: i18n._(t`Custom ...`) });
	}

	return (
		<SelectCustom
			options={sizes}
			label={props.label}
			customLabel={props.customLabel}
			value={props.value}
			onChange={props.onChange}
			variant={props.variant}
			allowCustom={props.allowCustom}
		/>
	);
}

Format.defaultProps = {
	allowAuto: false,
	allowCustom: false,
	variant: 'outlined',
	label: <Trans>Format</Trans>,
	customLabel: <Trans>Custom format</Trans>,
	onChange: function (event) {},
};

function PixFormat(props) {
	const { i18n } = useLingui();
	const sizes = [
		{ value: 'yuv420p', label: 'yuv420p' },
		{ value: 'nv12', label: 'nv12' },
	];

	if (props.allowAuto === true) {
		sizes.unshift({ value: 'auto', label: 'auto' });
	}

	if (props.allowCustom === true) {
		sizes.push({ value: 'custom', label: i18n._(t`Custom ...`) });
	}

	return (
		<SelectCustom
			options={sizes}
			label={props.label}
			customLabel={props.customLabel}
			value={props.value}
			onChange={props.onChange}
			variant={props.variant}
			allowCustom={props.allowCustom}
		/>
	);
}

PixFormat.defaultProps = {
	allowAuto: false,
	allowCustom: false,
	variant: 'outlined',
	label: <Trans>Pixel Format</Trans>,
	customLabel: <Trans>Custom format</Trans>,
	onChange: function (event) {},
};

function FpsMode(props) {
	return (
		<Select label={<Trans>Framerate mode</Trans>} value={props.value} onChange={props.onChange}>
			<MenuItem value="passthrough">
				<Trans>Frame is passed through (Passthrough)</Trans>
			</MenuItem>
			<MenuItem value="cfr">
				<Trans>Constant frame rate (CFR)</Trans>
			</MenuItem>
			<MenuItem value="vfr">
				<Trans>Variable frame rate (VFR)</Trans>
			</MenuItem>
			<MenuItem value="auto">
				<Trans>Choose between CFR and VFR (Auto)</Trans>
			</MenuItem>
		</Select>
	);
}

FpsMode.defaultProps = {
	value: '',
	onChange: function (event) {},
};

export default {
	Bitrate,
	GOP,
	Framerate,
	Profile,
	Size,
	Width,
	Height,
	Format,
	PixFormat,
	FpsMode,
};
