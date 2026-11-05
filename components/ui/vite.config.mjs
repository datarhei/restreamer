import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import { lingui } from '@lingui/vite-plugin';

export default defineConfig({
	base: './',
	plugins: [react(), lingui({ macroTransform: true })],
	build: {
		outDir: 'build',
		emptyOutDir: true,
	},
	test: {
		environment: 'jsdom',
		globals: true,
		setupFiles: ['./src/setupTests.jsx'],
	},
});
