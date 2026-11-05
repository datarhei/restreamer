import { createTheme } from '@mui/material/styles';

import base from './base';
import globals from './global';

import accordion from './components/accordion';
import accordionSummary from './components/accordionSummary';
import alert from './components/alert';
import appBar from './components/appBar';
import backdrop from './components/backdrop';
import box from './components/box';
import button from './components/button';
import checkobox from './components/checkbox';
import ctypography from './components/typography';
import dialog from './components/dialog';
import divider from './components/divider';
import fab from './components/fab';
import formControlLabel from './components/formControlLabel';
import formHelperText from './components/formHelperText';
import formLabel from './components/formLabel';
import iconButton from './components/iconButton';
import inputLabel from './components/inputLabel';
import link from './components/link';
import listItem from './components/listItem';
import listSubheader from './components/listSubheader';
import menu from './components/menu';
import menuItem from './components/menuItem';
import modal from './components/modal';
import outlinedInput from './components/outlinedInput';
import paper from './components/paper';
import popover from './components/popover';
import select from './components/select';
import snackbar from './components/snackbar';
import tab from './components/tab';
import tabPanel from './components/tabPanel';
import tabs from './components/tabs';
import tabScrollButton from './components/tabScrollButton';
import toggleButton from './components/toggleButton';
import toggleButtonGroup from './components/toggleButtonGroup';
import tooltip from './components/tooltip';

// https://mui.com/customization/default-theme/

const theme = createTheme({
	...base,

	components: {
		MuiCssBaseline: {
			styleOverrides: {
				...globals,
			},
		},
		MuiAccordion: { ...accordion },
		MuiAccordionSummary: { ...accordionSummary },
		MuiAlert: { ...alert },
		MuiAppBar: { ...appBar },
		MuiBackdrop: { ...backdrop },
		MuiBox: { ...box },
		MuiButton: { ...button },
		MuiCheckbox: { ...checkobox },
		MuiDialog: { ...dialog },
		MuiDivider: { ...divider },
		MuiFab: { ...fab },
		MuiFormControlLabel: { ...formControlLabel },
		MuiFormHelperText: { ...formHelperText },
		MuiFormLabel: { ...formLabel },
		MuiIconButton: { ...iconButton },
		MuiInputLabel: { ...inputLabel },
		MuiLink: { ...link },
		MuiListItem: { ...listItem },
		MuiListSubheader: { ...listSubheader },
		MuiMenu: { ...menu },
		MuiMenuItem: { ...menuItem },
		MuiModal: { ...modal },
		MuiOutlinedInput: { ...outlinedInput },
		MuiPaper: { ...paper },
		MuiPopover: { ...popover },
		MuiSelect: { ...select },
		MuiSnackbar: { ...snackbar },
		MuiTab: { ...tab },
		MuiTabPanel: { ...tabPanel },
		MuiTabs: { ...tabs },
		MuiTabScrollButton: { ...tabScrollButton },
		MuiToggleButton: { ...toggleButton },
		MuiToggleButtonGroup: { ...toggleButtonGroup },
		MuiTooltip: { ...tooltip },
		MuiTypography: { ...ctypography },
	},
});

export default theme;
