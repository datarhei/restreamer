import { createTheme } from '@mui/material/styles';

import breakpoints from './breakpoints';
import palette from './palette';
import typography from './typography';
import shape from './shape';

const theme = createTheme({
	breakpoints: { ...breakpoints },
	palette: { ...palette },
	typography: { ...typography },
	shape: { ...shape },
});

export default theme;
