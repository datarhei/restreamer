const fs = require('node:fs');
const path = require('node:path');

const buildDirectory = path.resolve(__dirname, '..', 'build');
const assetDirectory = path.join(buildDirectory, 'assets');
const maximumInitialBytes = 1_300_000;
const maximumChunkBytes = 750_000;
const minimumDynamicCatalogs = 14;
const routeChunks = ['Add-', 'ChannelSelect-', 'Edit-', 'Main-', 'Player-', 'Playersite-', 'Settings-', 'Wizard-'];

const html = fs.readFileSync(path.join(buildDirectory, 'index.html'), 'utf8');
const initialAssets = [
	...new Set(
		[...html.matchAll(/(?:src|href)="\.\/(assets\/[^"#]+)"/g)].map((match) => match[1]),
	),
];

let initialBytes = 0;
for (const asset of initialAssets) {
	const target = path.join(buildDirectory, asset);
	if (!fs.existsSync(target)) {
		throw new Error(`Initial HTML references missing asset: ${asset}`);
	}
	initialBytes += fs.statSync(target).size;
}

if (initialBytes > maximumInitialBytes) {
	throw new Error(`Initial assets total ${initialBytes} bytes; limit is ${maximumInitialBytes}.`);
}

const scripts = fs.readdirSync(assetDirectory).filter((file) => file.endsWith('.js'));
const oversized = scripts
	.map((file) => ({ file, size: fs.statSync(path.join(assetDirectory, file)).size }))
	.filter(({ size }) => size > maximumChunkBytes);
if (oversized.length > 0) {
	throw new Error(
		`JavaScript chunks exceed ${maximumChunkBytes} bytes: ${oversized.map(({ file, size }) => `${file} (${size})`).join(', ')}`,
	);
}

const catalogs = scripts.filter((file) => file.startsWith('messages-'));
if (catalogs.length < minimumDynamicCatalogs) {
	throw new Error(`Expected at least ${minimumDynamicCatalogs} dynamic locale chunks, found ${catalogs.length}.`);
}

const missingRoutes = routeChunks.filter((prefix) => !scripts.some((file) => file.startsWith(prefix)));
if (missingRoutes.length > 0) {
	throw new Error(`Missing route chunks: ${missingRoutes.join(', ')}`);
}

const largestChunk = scripts
	.map((file) => ({ file, size: fs.statSync(path.join(assetDirectory, file)).size }))
	.sort((left, right) => right.size - left.size)[0];

console.log(
	`Build assets validated: ${initialAssets.length} initial files / ${initialBytes} bytes, ` +
		`${catalogs.length} dynamic catalogs, largest chunk ${largestChunk.file} / ${largestChunk.size} bytes.`,
);
