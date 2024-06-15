import { WorkerInterface } from './worker_interface';

type SimRequestAsync = (data: Uint8Array, progress: (result: Uint8Array) => void) => Uint8Array;
type SimRequestSync = (data: Uint8Array) => Uint8Array;

// Functions provided or used by the wasm lib.
declare global {
	function wasmready(): void;
	const bulkSimAsync: SimRequestAsync;
	const computeStats: SimRequestSync;
	const computeStatsJson: SimRequestSync;
	const raidSim: SimRequestSync;
	const raidSimJson: SimRequestSync;
	const raidSimAsync: SimRequestAsync;
	const statWeights: SimRequestSync;
	const statWeightsAsync: SimRequestAsync;
	const statWeightRequests: SimRequestSync;
	const statWeightCompute: SimRequestSync;
	const raidSimResultCombination: SimRequestSync;
	const raidSimRequestSplit: SimRequestSync;
	const abortById: SimRequestSync;
}

// Wasm binary calls this function when its done loading.
// eslint-disable-next-line @typescript-eslint/no-unused-vars
globalThis.wasmready = function () {
	new WorkerInterface({
		bulkSimAsync: (data, progress) => bulkSimAsync(data, progress),
		computeStats: computeStats,
		computeStatsJson: computeStatsJson,
		raidSim: raidSim,
		raidSimJson: raidSimJson,
		raidSimAsync: (data, progress) => raidSimAsync(data, progress),
		statWeights: statWeights,
		statWeightsAsync: (data, progress) => statWeightsAsync(data, progress),
		statWeightRequests: statWeightRequests,
		statWeightCompute: statWeightCompute,
		raidSimRequestSplit: raidSimRequestSplit,
		raidSimResultCombination: raidSimResultCombination,
		abortById: abortById,
	}).ready(true);
};

const go = new Go();
let inst: WebAssembly.Instance | null = null;

WebAssembly.instantiateStreaming(fetch('lib.wasm'), go.importObject).then(async result => {
	inst = result.instance;
	// console.log("loading wasm...")
	await go.run(inst);
});

export {};
