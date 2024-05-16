import { buildProgressId, workerPostMessage } from './shared';
import { SimRequest, SimRequestSync, WorkerReceiveMessage } from './types';

// Wasm binary calls this function when its done loading.
// eslint-disable-next-line @typescript-eslint/no-unused-vars
globalThis.wasmready = function () {
	workerPostMessage({
		msg: 'ready',
	});
};

const go = new Go();
let inst: WebAssembly.Instance | null = null;

WebAssembly.instantiateStreaming(fetch('lib.wasm'), go.importObject).then(async result => {
	inst = result.instance;
	// console.log("loading wasm...")
	await go.run(inst);
});

// eslint-disable-next-line @typescript-eslint/no-unused-vars
let workerID = '';

addEventListener(
	'message',
	async ({ data }: MessageEvent<WorkerReceiveMessage>) => {
		const { id, msg, inputData } = data;

		let handled = false;

		const requests: [SimRequest, SimRequestSync][] = [
			[
				SimRequest.bulkSimAsync,
				data =>
					bulkSimAsync(data, result => {
						workerPostMessage({
							id: buildProgressId(id),
							msg: 'progress',
							outputData: result,
						});
					}),
			],
			[SimRequest.computeStats, computeStats],
			[SimRequest.computeStatsJson, computeStatsJson],
			[SimRequest.raidSim, raidSim],
			[SimRequest.raidSimJson, raidSimJson],
			[
				SimRequest.raidSimAsync,
				data =>
					raidSimAsync(data, result => {
						workerPostMessage({
							id: buildProgressId(id),
							msg: 'progress',
							outputData: result,
						});
					}),
			],
			[
				SimRequest.statWeightsAsync,
				data =>
					statWeightsAsync(data, result => {
						workerPostMessage({
							id: buildProgressId(id),
							msg: 'progress',
							outputData: result,
						});
					}),
			],
		];
		requests.forEach(([funcName, func]) => {
			if (msg === funcName) {
				const outputData = func(inputData);

				workerPostMessage({
					id: id,
					msg: funcName,
					outputData: outputData,
				});
				handled = true;
			}
		});

		if (handled) {
			return;
		}

		if (msg === 'setID') {
			workerID = id;
			workerPostMessage({ msg: 'idconfirm' });
		}
	},
	false,
);

export {};
