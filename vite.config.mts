/** @type {import('vite').UserConfig} */

import fs from 'fs';
import glob from 'glob';
import { IncomingMessage, ServerResponse } from 'http';
import path from 'path';
import { ConfigEnv, defineConfig, PluginOption, UserConfigExport } from 'vite';
import { checker } from 'vite-plugin-checker';

export const BASE_PATH = path.resolve(__dirname, 'ui');
export const OUT_DIR = path.join(__dirname, 'dist', 'cata');

function serveExternalAssets() {
	const workerMappings = {
		'/cata/sim_worker.js': '/cata/local_worker.js',
		'/cata/net_worker.js': '/cata/net_worker.js',
		'/cata/lib.wasm': '/cata/lib.wasm',
	};

	return {
		name: 'serve-external-assets',
		configureServer(server) {
			server.middlewares.use((req, res, next) => {
				const url = req.url!;

				if (Object.keys(workerMappings).includes(url)) {
					const targetPath = workerMappings[url as keyof typeof workerMappings];
					const assetsPath = path.resolve(__dirname, './dist/cata');
					const requestedPath = path.join(assetsPath, targetPath.replace('/cata/', ''));
					serveFile(res, requestedPath);
					return;
				}

				if (url.includes('/cata/assets')) {
					const assetsPath = path.resolve(__dirname, './assets');
					const assetRelativePath = url.split('/cata/assets')[1];
					const requestedPath = path.join(assetsPath, assetRelativePath);

					serveFile(res, requestedPath);
					return;
				} else {
					next();
				}
			});
		},
	} satisfies PluginOption;
}

function serveFile(res: ServerResponse<IncomingMessage>, filePath: string) {
	if (fs.existsSync(filePath)) {
		const contentType = determineContentType(filePath);
		res.writeHead(200, { 'Content-Type': contentType });
		fs.createReadStream(filePath).pipe(res);
	} else {
		console.log('Not found on filesystem: ', filePath);
		res.writeHead(404, { 'Content-Type': 'text/plain' });
		res.end('Not Found');
	}
}

function determineContentType(filePath: string) {
	const extension = path.extname(filePath).toLowerCase();
	switch (extension) {
		case '.jpg':
		case '.jpeg':
			return 'image/jpeg';
		case '.png':
			return 'image/png';
		case '.gif':
			return 'image/gif';
		case '.css':
			return 'text/css';
		case '.js':
			return 'text/javascript';
		case '.woff':
		case '.woff2':
			return 'font/woff2';
		case '.json':
			return 'application/json';
		case '.wasm':
			return 'application/wasm'; // Adding MIME type for WebAssembly files
		// Add more cases as needed
		default:
			return 'application/octet-stream';
	}
}

export const getBaseConfig = ({ command, mode }: ConfigEnv) =>
	({
		base: '/cata/',
		root: path.join(__dirname, 'ui'),
		build: {
			outDir: OUT_DIR,
			minify: mode === 'development' ? false : 'terser',
			sourcemap: command === 'serve' ? 'inline' : false,
			target: ['es2020'],
		},
	}) satisfies Partial<UserConfigExport>;

export default defineConfig(({ command, mode }) => {
	const baseConfig = getBaseConfig({ command, mode });
	return {
		...baseConfig,
		plugins: [
			serveExternalAssets(),
			checker({
				root: path.resolve(__dirname, 'ui'),
				typescript: true,
				enableBuild: true,
			}),
		],
		build: {
			...baseConfig.build,
			rollupOptions: {
				input: {
					...glob.sync(path.resolve(BASE_PATH, '**/index.html').replace(/\\/g, '/')).reduce<Record<string, string>>((acc, cur) => {
						const name = path.relative(__dirname, cur);
						acc[name] = cur;
						return acc;
					}, {}),
					// Add shared.scss as a separate entry if needed or handle it separately
				},
				output: {
					assetFileNames: () => 'bundle/[name]-[hash].style.css',
					entryFileNames: () => 'bundle/[name]-[hash].entry.js',
					chunkFileNames: () => 'bundle/[name]-[hash].chunk.js',
				},
			},
			server: {
				origin: 'http://localhost:3000',
				// Adding custom middleware to serve 'dist' directory in development
			},
		},
	};
});
