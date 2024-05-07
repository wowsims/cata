import { REPO_NAME } from './constants/other.js';
import {
	BulkSimRequest,
	BulkSimResult,
	ComputeStatsRequest,
	ComputeStatsResult,
	ProgressMetrics,
	RaidSimRequest,
	RaidSimResult,
	StatWeightsRequest,
	StatWeightsResult,
} from './proto/api.js';
import SimProgress from './proto_utils/sim_progress';
import { noop } from './utils';

const SIM_WORKER_URL = `/${REPO_NAME}/sim_worker.js`;

export type WorkerProgressCallback = (progressMetrics: ProgressMetrics) => void;

export class WorkerPool {
	private concurrency = 1;
	private workers: SimWorker[];

	constructor(numWorkers: number) {
		this.concurrency = numWorkers;
		this.workers = [];
		for (let i = 0; i < numWorkers; i++) {
			this.workers.push(new SimWorker());
		}
	}

	getWorkers() {
		return this.workers;
	}

	private getLeastBusyWorker(): SimWorker {
		return this.workers.reduce((curMinWorker, nextWorker) => (curMinWorker.numTasksRunning < nextWorker.numTasksRunning ? curMinWorker : nextWorker));
	}

	async makeApiCall(requestName: string, request: Uint8Array): Promise<Uint8Array> {
		return this.getLeastBusyWorker().doApiCall(requestName, request, '');
	}

	async computeStats(request: ComputeStatsRequest): Promise<ComputeStatsResult> {
		const result = await this.makeApiCall('computeStats', ComputeStatsRequest.toBinary(request));
		return ComputeStatsResult.fromBinary(result);
	}

	async statWeightsAsync(request: StatWeightsRequest, onProgress: WorkerProgressCallback): Promise<StatWeightsResult> {
		console.log('Stat weights request: ' + StatWeightsRequest.toJsonString(request));
		const worker = this.getLeastBusyWorker();
		const id = worker.makeTaskId();
		// Add handler for the progress events
		worker.addPromiseFunc(id + 'progress', this.newProgressHandler(id, worker, onProgress), noop);

		// Now start the async sim
		const resultData = await worker.doApiCall('statWeightsAsync', StatWeightsRequest.toBinary(request), id);
		const result = ProgressMetrics.fromBinary(resultData);
		console.log('Stat weights result: ' + StatWeightsResult.toJsonString(result.finalWeightResult!));
		return result.finalWeightResult!;
	}

	async bulkSimAsync(request: BulkSimRequest, onProgress: WorkerProgressCallback): Promise<BulkSimResult> {
		console.log('bulk sim request: ' + BulkSimRequest.toJsonString(request, { enumAsInteger: true }));
		const worker = this.getLeastBusyWorker();
		const id = worker.makeTaskId();
		// Add handler for the progress events
		worker.addPromiseFunc(id + 'progress', this.newProgressHandler(id, worker, onProgress), noop);

		// Now start the async sim
		const resultData = await worker.doApiCall('bulkSimAsync', BulkSimRequest.toBinary(request), id);
		const result = ProgressMetrics.fromBinary(resultData);
		const resultJson = BulkSimResult.toJson(result.finalBulkResult!) as any;
		console.log('bulk sim result: ' + JSON.stringify(resultJson));
		return result.finalBulkResult!;
	}

	async raidSimAsyncConcurrent(requests: RaidSimRequest[], onProgress: WorkerProgressCallback): Promise<RaidSimResult> {
		performance.mark('raidSimAsync:start');
		console.log('Raid sim request: ' + RaidSimRequest.toJsonString(requests[0]));
		const iterations = requests[0].simOptions!.iterations || 0;
		const iterationsTotal = requests.reduce<number>((total, r) => (total += r.simOptions!.iterations), 0);
		const concurrency = Math.min(iterations, this.concurrency);
		const simData = new SimProgress({
			concurrency,
			iterationsTotal,
		});

		const handleConcurrentProgress = (workerIndex: number, progress: ProgressMetrics, simData: SimProgress) => {
			simData.updateProgress(workerIndex, progress);
			const { presimRunning, totalSims, completedSims } = progress;
			onProgress({
				totalIterations: simData.data.iterationsTotal,
				completedIterations: simData.getIterationsDone(),
				totalSims,
				completedSims,
				dps: simData.getDpsAvg(),
				hps: simData.getHpsAvg(),
				presimRunning,
			});
		};

		await Promise.all(
			requests.map(async (request, index) => {
				const worker = this.getLeastBusyWorker();
				const id = worker.makeTaskId();

				console.log(`Started worker ${index + 1} with ${request.simOptions?.iterations} iterations and seed ${request.simOptions?.randomSeed}`);
				// Add handler for the progress events
				worker.addPromiseFunc(
					id + 'progress',
					this.newProgressHandler(id, worker, progress => handleConcurrentProgress(index, progress, simData)),
					noop,
				);

				// Now start the async sim
				const resultData = await worker.doApiCall('raidSimAsync', RaidSimRequest.toBinary(request), id);
				const result = ProgressMetrics.fromBinary(resultData);

				return result.finalRaidResult!;
			}),
		);

		const hasMissingFinalResult = simData.data.finalResults.some(result => !result);

		if (hasMissingFinalResult) {
			throw new Error('Missing one or more final sim result(s)!');
		}
		const results = simData.getCombinedFinalResult();
		const { logs: _, ...resultJson } = results;

		performance.mark('raidSimAsync:end');
		console.log(`Raid sim result: ${JSON.stringify(resultJson)}`);
		console.log(
			`Finished ${simData.getIterationsDone()} in ${(
				performance.measure('raidSimAsync', 'raidSimAsync:start', 'raidSimAsync:end').duration / 1000
			).toFixed(2)}s`,
		);
		return results;
	}

	async raidSimAsync(request: RaidSimRequest, onProgress: WorkerProgressCallback): Promise<RaidSimResult> {
		console.log('Raid sim request: ' + RaidSimRequest.toJsonString(request));
		const worker = this.getLeastBusyWorker();
		const id = worker.makeTaskId();
		// Add handler for the progress events
		worker.addPromiseFunc(id + 'progress', this.newProgressHandler(id, worker, onProgress), noop);

		// Now start the async sim
		const resultData = await worker.doApiCall('raidSimAsync', RaidSimRequest.toBinary(request), id);
		const result = ProgressMetrics.fromBinary(resultData);

		// Don't print the logs because it just clogs the console.
		const results = RaidSimResult.toJson(result.finalRaidResult!) as any;
		const { logs: _, ...resultJson } = results;
		console.log(`Raid sim result: ${JSON.stringify(resultJson)}`);
		return result.finalRaidResult!;
	}

	newProgressHandler(id: string, worker: SimWorker, onProgress: WorkerProgressCallback): (progressData: any) => void {
		return (progressData: any) => {
			const progress = ProgressMetrics.fromBinary(progressData);
			onProgress(progress);
			// If we are done, stop adding the handler.
			if (progress.finalRaidResult != null || progress.finalWeightResult != null) {
				return;
			}

			worker.addPromiseFunc(id + 'progress', this.newProgressHandler(id, worker, onProgress), noop);
		};
	}
}

class SimWorker {
	numTasksRunning: number;
	private taskIdsToPromiseFuncs: Record<string, [(result: any) => void, (error: any) => void]>;
	private worker: Worker;
	public type: 'net' | 'local' | null = null;
	private onReady: Promise<void>;

	constructor() {
		this.numTasksRunning = 0;
		this.taskIdsToPromiseFuncs = {};

		this.worker = new window.Worker(SIM_WORKER_URL);

		let resolveReady: (() => void) | null = null;
		this.onReady = new Promise((_resolve, _reject) => {
			resolveReady = _resolve;
		});

		this.worker.onmessage = event => {
			if (event.data.msg == 'ready') {
				this.type = event.data.workerType;
				this.worker.postMessage({ msg: 'setID', id: '1' });
				resolveReady!();
			} else if (event.data.msg == 'idconfirm') {
				// Do nothing
			} else {
				const id = event.data.id;
				if (!this.taskIdsToPromiseFuncs[id]) {
					console.warn('Unrecognized result id: ' + id);
					return;
				}

				const promiseFuncs = this.taskIdsToPromiseFuncs[id];
				delete this.taskIdsToPromiseFuncs[id];
				this.numTasksRunning--;

				promiseFuncs[0](event.data.outputData);
			}
		};
	}

	addPromiseFunc(id: string, callback: (result: any) => void, onError: (error: any) => void) {
		this.taskIdsToPromiseFuncs[id] = [callback, onError];
	}

	makeTaskId(): string {
		let id = '';
		const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
		for (let i = 0; i < 16; i++) {
			id += characters.charAt(Math.floor(Math.random() * characters.length));
		}
		return id;
	}

	async doApiCall(requestName: string, request: Uint8Array, id: string): Promise<Uint8Array> {
		this.numTasksRunning++;
		await this.onReady;

		const taskPromise = new Promise<Uint8Array>((resolve, reject) => {
			if (!id) {
				id = this.makeTaskId();
			}
			this.taskIdsToPromiseFuncs[id] = [resolve, reject];

			this.worker.postMessage({
				msg: requestName,
				id: id,
				inputData: request,
			});
		});
		return await taskPromise;
	}
}
