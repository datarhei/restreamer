import React from 'react';

import { useLingui } from '@lingui/react';
import makeStyles from '@mui/styles/makeStyles';
import MenuItem from '@mui/material/MenuItem';
import Select from '@mui/material/Select';

const useStyles = makeStyles((theme) => ({
	root: {
		fontSize: '1rem',
		fontWeight: '400',
		backgroundColor: 'none',
		color: theme.palette.text.primary,
		borderButtom: 'none',
		'&:hover': {
			color: theme.palette.text.primary,
		},
		'& fieldset': {
			display: 'none',
		},
		'& .MuiSelect-select': {
			padding: '0px 0px',
			marginRight: '10px',
		},
		'&:before': {
			borderBottom: 'none',
		},
	},
}));

export default function LanguageSelect(props) {
	const classes = useStyles();
	const { i18n } = useLingui();

	const handleChange = (event) => {
		const language = event.target.value;

		i18n.activate(language);

		props.onChange(language);
	};

	return (
		<Select className={classes.root} variant="standard" displayEmpty value={i18n.locale} onChange={handleChange}>
			<MenuItem value="en">English</MenuItem>
			<MenuItem value="da">Dansk</MenuItem>
			<MenuItem value="de">Deutsch</MenuItem>
			<MenuItem value="el">Ελληνικά</MenuItem>
			<MenuItem value="es">Español</MenuItem>
			<MenuItem value="fr">Français</MenuItem>
			<MenuItem value="it">Italiano</MenuItem>
			<MenuItem value="ko">한국어</MenuItem>
			<MenuItem value="pl">Polski</MenuItem>
			<MenuItem value="pt-br">Português</MenuItem>
			<MenuItem value="ru">Русский</MenuItem>
			<MenuItem value="sl">Slovenščina</MenuItem>
			<MenuItem value="tr">Türkçe</MenuItem>
			<MenuItem value="uk">Українська мова</MenuItem>
			<MenuItem value="zh-hans">中文(简体)</MenuItem>
		</Select>
	);
}

LanguageSelect.defaultProps = {
	onChange: function (lang) {},
};
