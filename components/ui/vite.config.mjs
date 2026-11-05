import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import { lingui } from '@lingui/vite-plugin';

export default defineConfig({
	base: './',
	plugins: [react(), lingui({ macroTransform: true })],
	build: {
		outDir: 'build',
		emptyOutDir: true,
		// The largest chunk is the lazily loaded video.js player (about 708 kB raw,
		// 208 kB gzip). Route and locale chunks remain independently loadable.
		chunkSizeWarningLimit: 750,
	},
	test: {
		environment: 'jsdom',
		globals: true,
		setupFiles: ['./src/setupTests.jsx'],
	},
});
