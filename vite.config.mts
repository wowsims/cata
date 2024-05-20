/** @type {import('vite').UserConfig} */
import glob from 'glob';
import path from 'path';
import copy from 'rollup-plugin-copy';
import { ConfigEnv, defineConfig, UserConfigExport } from 'vite';
import { checker } from 'vite-plugin-checker';

import { serveExternalAssets, serveFile } from './node_modules/@wowsims/assets/helpers/assets';

export const BASE_PATH = path.resolve(__dirname, 'ui');
export const OUT_DIR = path.join(__dirname, 'dist', 'cata');

const workerMappings = {
	'/cata/sim_worker.js': '/cata/local_worker.js',
	'/cata/net_worker.js': '/cata/net_worker.js',
	'/cata/lib.wasm': '/cata/lib.wasm',
};

const replacePaths = [
	{
		replacePath: '/cata/assets',
		sourcePath: path.resolve(__dirname, './assets'),
	},
];

export const getBaseConfig = ({ command, mode }: ConfigEnv) =>
	({
		base: '/cata/',
		root: path.join(__dirname, 'ui'),
		publicDir: false,
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
			serveExternalAssets({
				assets: replacePaths,
				additionalMiddlewareHook: (req, res) => {
					const url = req.url!;
					if (Object.keys(workerMappings).includes(url)) {
						const targetPath = workerMappings[url as keyof typeof workerMappings];
						const assetsPath = path.resolve(__dirname, './dist/cata');
						const requestedPath = path.join(assetsPath, targetPath.replace('/cata/', ''));
						serveFile(res, requestedPath);
						return true;
					}
					return false;
				},
				transform(code) {
					// Replace all wowsims repo assets to local assets
					if (code.match(/\@wowsims\/assets/g)) {
						code = code.replace(/\@wowsims\/assets/g, '/cata/assets/@wowsims');
					}
					return code;
				},
			}),
			copy({
				copyOnce: mode === 'development',
				recursive: true,
				targets: [
					{
						src: `node_modules/@wowsims/assets/public/**/*`,
						dest: 'assets/@wowsims',
					},
				],
				hook: 'buildEnd',
			}),
			checker({
				root: path.resolve(__dirname, 'ui'),
				typescript: true,
				enableBuild: true,
			}),
		],
		esbuild: {
			jsxFactory: 'element',
			jsxFragment: 'fragment',
			jsxInject: "import { element, fragment } from 'tsx-vanilla';",
		},
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
		},
		server: {
			origin: 'http://localhost:3000',
			// Adding custom middleware to serve 'dist' directory in development
		},
	};
});
