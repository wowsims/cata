import fs from 'node:fs/promises';

import { exec as syncExec } from 'child_process';
import { promisify } from 'util';
const execAsync = promisify(syncExec);
import minimist from 'minimist';
import path from 'path';
import { build } from 'vite';

import { BASE_PATH, getBaseConfig } from './vite.config.mjs';

const WORKER_BASE_PATH = path.resolve(BASE_PATH, 'worker');

const workers = {
	local_worker: path.resolve(WORKER_BASE_PATH, 'local_worker.ts'),
	net_worker: path.resolve(WORKER_BASE_PATH, 'net_worker.ts'),
	sim_worker: path.resolve(WORKER_BASE_PATH, 'sim_worker.ts'),
};

const args = minimist(process.argv.slice(2), { boolean: ['watch'] });

const buildWorkers = async () => {
	const { stdout } = await execAsync('go env GOROOT');
	const GO_ROOT = stdout.replace('\n', '');
	const wasmExecutablePath = path.join(GO_ROOT, '/misc/wasm/wasm_exec.js');
	const wasmFile = await fs.readFile(wasmExecutablePath, 'utf8');

	Object.entries(workers).forEach(async ([name, sourcePath]) => {
		const baseConfig = getBaseConfig({ command: 'build', mode: 'production' });
		await build({
			...baseConfig,
			configFile: false,
			plugins: [
				{
					name: 'add-wasm-exec-file',
					transform: async (code, id) => {
						if (id.includes('sim_worker.ts')) {
							code = wasmFile + code;
						}
						return code;
					},
				},
			],
			build: {
				...baseConfig.build,
				watch: args.watch === true ? {} : undefined,
				minify: false,
				emptyOutDir: false,
				lib: {
					entry: { [name]: sourcePath },
					name: `${name}.js`,
					formats: ['cjs'],
				},
			},
		});
	});
};

buildWorkers();
