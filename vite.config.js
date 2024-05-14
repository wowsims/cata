import fs from 'fs';
import glob from 'glob';
import path from 'path';
import { defineConfig } from 'vite';

function serveExternalAssets() {
	return {
		name: 'serve-external-assets',
		configureServer(server) {
			server.middlewares.use((req, res, next) => {
				const workerMappings = {
					'/cata/sim_worker.js': '/cata/local_worker.js',
					'/cata/net_worker.js': '/cata/net_worker.js',
					'/cata/lib.wasm': '/cata/lib.wasm',
				};

				if (Object.keys(workerMappings).includes(req.url)) {
					const targetPath = workerMappings[req.url];
					const assetsPath = path.resolve(__dirname, './dist/cata');
					const requestedPath = path.join(assetsPath, targetPath.replace('/cata/', ''));
					serveFile(res, requestedPath);
					return;
				}

				if (req.url.includes('/cata/assets')) {
					const assetsPath = path.resolve(__dirname, './assets');
					const assetRelativePath = req.url.split('/cata/assets')[1];
					const requestedPath = path.join(assetsPath, assetRelativePath);

					serveFile(res, requestedPath);
					return;
				} else {
					next();
				}
			});
		},
	};
}

function serveFile(res, filePath) {
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

function determineContentType(filePath) {
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

export default defineConfig(({ command, mode }) => ({
	plugins: [serveExternalAssets()],
	base: '/cata/',
	root: path.join(__dirname, 'ui'),
	build: {
		outDir: path.join(__dirname, 'dist', 'cata'),
		minify: mode === 'development' ? false : 'terser',
		sourcemap: command === 'serve' ? 'inline' : 'false',
		target: ['es2020'],
		rollupOptions: {
			input: {
				...glob.sync(path.resolve(__dirname, 'ui', '**/index.html').replace(/\\/g, '/')).reduce((acc, cur) => {
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
}));
