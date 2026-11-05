const { spawnSync } = require('node:child_process');

const audit = spawnSync('yarn', ['audit', '--groups', 'dependencies', '--json'], {
	encoding: 'utf8',
	maxBuffer: 16 * 1024 * 1024,
});

if (audit.error) {
	console.error(`Unable to run the production dependency audit: ${audit.error.message}`);
	process.exit(2);
}

const entries = audit.stdout
	.split(/\r?\n/)
	.filter(Boolean)
	.map((line) => {
		try {
			return JSON.parse(line);
		} catch {
			return undefined;
		}
	});

const summary = entries.findLast((entry) => entry?.type === 'auditSummary');

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

const advisories = new Map();
for (const entry of entries) {
	const advisory = entry?.type === 'auditAdvisory' ? entry.data?.advisory : undefined;
	if (!advisory) {
		continue;
	}

	const id = advisory.github_advisory_id ?? `${advisory.module_name}:${advisory.title}`;
	advisories.set(id, advisory);
}

for (const advisory of [...advisories.values()].sort((left, right) =>
	left.module_name.localeCompare(right.module_name),
)) {
	console.log(
		`- ${advisory.severity}: ${advisory.module_name}: ${advisory.title} ` +
			`(${advisory.github_advisory_id ?? 'no GHSA'}; patched ${advisory.patched_versions})`,
	);
}

if (vulnerabilities.critical > 0 || vulnerabilities.high > 0 || vulnerabilities.moderate > 0) {
	console.error('Production dependencies contain moderate, high, or critical vulnerabilities.');
	process.exit(1);
}
