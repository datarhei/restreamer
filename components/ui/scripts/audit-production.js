const { spawnSync } = require('node:child_process');

const audit = spawnSync('yarn', ['audit', '--groups', 'dependencies', '--json'], {
	encoding: 'utf8',
	maxBuffer: 16 * 1024 * 1024,
});

if (audit.error) {
	console.error(`Unable to run the production dependency audit: ${audit.error.message}`);
	process.exit(2);
}

const summary = audit.stdout
	.split(/\r?\n/)
	.filter(Boolean)
	.map((line) => {
		try {
			return JSON.parse(line);
		} catch {
			return undefined;
		}
	})
	.findLast((entry) => entry?.type === 'auditSummary');

if (!summary) {
	if (audit.stderr) {
		console.error(audit.stderr.trim());
	}
	console.error('Yarn did not return an audit summary. Refusing to guess that everything is fine.');
	process.exit(2);
}

const vulnerabilities = summary.data.vulnerabilities;
console.log(
	`Production dependency audit: ${vulnerabilities.critical} critical, ` +
		`${vulnerabilities.high} high, ${vulnerabilities.moderate} moderate, ` +
		`${vulnerabilities.low} low`,
);

if (vulnerabilities.critical > 0 || vulnerabilities.high > 0) {
	console.error('Production dependencies contain high or critical vulnerabilities.');
	process.exit(1);
}
