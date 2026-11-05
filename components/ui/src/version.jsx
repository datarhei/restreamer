import pkg from '../package.json';

const Core = '^16.11.0';
const FFmpeg = '^5.1.0 || ^6.1.0';
const UI = pkg.bundle ? pkg.bundle : pkg.name + ' v' + pkg.version;
const Version = pkg.version;

export { Core, FFmpeg, UI, Version };
