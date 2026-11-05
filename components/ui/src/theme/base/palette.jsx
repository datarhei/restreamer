/* eslint-disable import/no-anonymous-default-export */
export default {
	common: {
		// MuiButton, ChannelList.js, index.js, Player.js
		black: '#000',
		// Header.js, ChannelList.js, Edit->index.js
		white: '#fff',
	},
	background: {
		// MuiAccordion
		accordion: 'rgba(77, 77, 77, .4)',
		// MuiAccordion
		accordion_border: 'rgba(128, 128, 128, .6)',
		// MuiMenuItem, Header.js (MenuIcon, MenuItem), Progress.js (Progress boxes)
		box_default: 'rgb(77, 77, 77)',
		// MuiListItem, Progress.js (Progress boxes - danger)
		box_danger: 'rgb(193, 39, 45)',
		// MuiButton, global.js (Page background)
		button_disabled: '#333333',
		// MuiButton, MuiMenu, MuiToggleButton, MultiSelectOption.js (MenuItem)
		dark1: 'rgba(0, 0, 0, .1)',
		// MuiOutlinedInput, BoxText.js (color=dark)
		dark2: 'rgba(0, 0, 0, .25)',
		// Footer.js, Textarea.js, global.js (Scrollbar)
		footer1: 'rgba(66, 61, 63, .9)',
		// Footer.js, Textarea.js
		footer2: 'rgba(39, 36, 37, .9)',
		// MuiAppBar, BoxText.js (color=light = default), TabsHorizontal.js, Paper.js
		light1: 'rgba(255, 255, 255, .1)',
		// MuiDialog, MuiPaper, MuiTooltip, Header.js, ModalContent.js, Dialog.js
		modal: 'rgb(71, 71, 71)',
		// BoxTextarea.js, Progress.js, Process.js, Textarea.js, gobal.js (textarea)
		modalbox: 'rgb(91, 91, 91)',
		// MuiFab, MuiListItem, MuiPaper, Headerjs,
		paper: 'rgba(56, 56, 56, .95)',
	},
	default: {
		main: '#686868',
	},
	primary: {
		// MuiOutlinedInput
		light: '#919090',
		// MuiAlert, MuiMenu, MuiOutlinedInput, MuiSelect
		main: '#686868',
		// MuiMenu, MuiOutlinedInput
		dark: '#4D4D4D',
		// MuiButton
		contrastText: '#fff',
	},
	secondary: {
		// MuiAlert, MuiButton, MuiFab, MuiOutlinedInput, MuiSelect, MuiToggleButton
		main: '#39B54A',
		// MuiButton, MuiFab
		contrastText: '#fff',
	},
	service: {
		// MuiButton, Login.js
		main: '#EAEA05',
		// Paper.js
		contrastText: 'rgba(43,41,42,.95)',
	},
	error: {
		// MuiAlert, MuiButton, BoxText.js (color=danger), Env.js
		main: '#C1272D',
		// MuiButton
		contrastText: '#fff',
	},
	warning: {
		// MuiAlert, Main->index.js
		main: '#E28014',
	},
	selected: {
		// MuiFab
		main: '#fff',
	},
	text: {
		// MuiBackdrop, MuiButton, MuiDialog, MuiFormLabel, MuiIconButton,
		// MuiInputLabel, MuiMenu, MuiOutlinedInput, MuiPaper, MuiSelect, MuiTab,
		// MuiToggleButton, MuiTypography,
		// Header.js, Checkbox.js, Env.js, LanguageSelect.js, ModalContent.js,
		// TabsHorizontal.js, Dialog.js, global.js (Textarea)
		primary: '#fff',
		// MuiSelect, MuiTab, Footer.js, Header.js
		secondary: 'rgba(255, 255, 255, 0.7)',
		// MuiCheckbox, Checkbox.js
		disabled: 'rgba(255, 255, 255, 0.38)',
		// MuiTypography
		hint: 'rgba(255, 255, 255, 0.7)',
	},
};
