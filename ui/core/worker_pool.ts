import { REPO_NAME } from './constants/other.js';
import {
	AsyncAPIResult,
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
import { noop } from './utils';

const SIM_WORKER_URL = `/${REPO_NAME}/sim_worker.js`;

export type WorkerProgressCallback = (progressMetrics: ProgressMetrics) => void;
type WorkerOptions = {
	signal?: AbortSignal;
};

export class WorkerPool {
	private workers: Array<SimWorker>;

	constructor(numWorkers: number) {
		this.workers = [];
		for (let i = 0; i < numWorkers; i++) {
			this.workers.push(new SimWorker());
		}
	}

	private getLeastBusyWorker(): SimWorker {
		return this.workers.reduce((curMinWorker, nextWorker) => (curMinWorker.numTasksRunning < nextWorker.numTasksRunning ? curMinWorker : nextWorker));
	}

	async makeApiCall(requestName: string, request: Uint8Array, { signal }: WorkerOptions = {}): Promise<Uint8Array> {
		return await this.getLeastBusyWorker().doApiCall(requestName, request, '', { signal });
	}

	async computeStats(request: ComputeStatsRequest, { signal }: WorkerOptions = {}): Promise<ComputeStatsResult> {
		const result = await this.makeApiCall('computeStats', ComputeStatsRequest.toBinary(request), { signal });
		return ComputeStatsResult.fromBinary(result);
	}

	private getProgressName(id: string) {
		return `${id}progress`;
	}

	async statWeightsAsync(request: StatWeightsRequest, onProgress: WorkerProgressCallback, { signal }: WorkerOptions = {}): Promise<StatWeightsResult> {
		console.log('Stat weights request: ' + StatWeightsRequest.toJsonString(request));
		const worker = this.getLeastBusyWorker();
		const id = worker.makeTaskId();
		// Add handler for the progress events
		worker.addPromiseFunc(this.getProgressName(id), this.newProgressHandler(id, worker, onProgress, { signal }), noop);

		// Now start the async sim
		const resultData = await worker.doApiCall('statWeightsAsync', StatWeightsRequest.toBinary(request), id, { signal });
		const result = ProgressMetrics.fromBinary(resultData);
		console.log('Stat weights result: ' + StatWeightsResult.toJsonString(result.finalWeightResult!));
		return result.finalWeightResult!;
	}

	async bulkSimAsync(request: BulkSimRequest, onProgress: WorkerProgressCallback, { signal }: WorkerOptions = {}): Promise<BulkSimResult> {
		console.log('bulk sim request: ' + BulkSimRequest.toJsonString(request, { enumAsInteger: true }));
		const worker = this.getLeastBusyWorker();
		const id = worker.makeTaskId();
		// Add handler for the progress events
		worker.addPromiseFunc(this.getProgressName(id), this.newProgressHandler(id, worker, onProgress, { signal }), noop);

		// Now start the async sim
		const resultData = await worker.doApiCall('bulkSimAsync', BulkSimRequest.toBinary(request), id, { signal });
		const result = ProgressMetrics.fromBinary(resultData);
		const resultJson = BulkSimResult.toJson(result.finalBulkResult!) as any;
		console.log('bulk sim result: ' + JSON.stringify(resultJson));
		return result.finalBulkResult!;
	}

	async raidSimAsync(request: RaidSimRequest, onProgress: WorkerProgressCallback, { signal }: WorkerOptions = {}): Promise<RaidSimResult> {
		console.log('Raid sim request: ' + RaidSimRequest.toJsonString(request));
		const worker = this.getLeastBusyWorker();
		const id = worker.makeTaskId();
		// Add handler for the progress events
		worker.addPromiseFunc(this.getProgressName(id), this.newProgressHandler(id, worker, onProgress, { signal }), noop);

		// Now start the async sim
		const resultData = await worker.doApiCall('raidSimAsync', RaidSimRequest.toBinary(request), id, { signal });
		const result = ProgressMetrics.fromBinary(resultData);

		// Don't print the logs because it just clogs the console.
		const resultJson = RaidSimResult.toJson(result.finalRaidResult!) as any;
		delete resultJson!['logs'];
		console.log('Raid sim result: ' + JSON.stringify(resultJson));
		return result.finalRaidResult!;
	}

	newProgressHandler(id: string, worker: SimWorker, onProgress: WorkerProgressCallback, { signal }: WorkerOptions = {}) {
		return (progressData: Uint8Array) => {
			const progress = ProgressMetrics.fromBinary(progressData);
			onProgress(progress);
			// If we are done, stop adding the handler.
			if (!!progress.finalRaidResult || !!progress.finalWeightResult) {
				return;
			}

			worker.addPromiseFunc(this.getProgressName(id), this.newProgressHandler(id, worker, onProgress, { signal }), noop);
		};
	}
}

class SimWorker {
	numTasksRunning: number;
	private taskIdsToPromiseFuncs: Record<string, [(result: any) => void, (error: any) => void]>;
	private worker: Worker | undefined;
	private onReady: Promise<void> | undefined;

	constructor() {
		this.numTasksRunning = 0;
		this.taskIdsToPromiseFuncs = {};
		this.initializeWorker();
	}

	initializeWorker() {
		this.worker = new window.Worker(SIM_WORKER_URL);

		let resolveReady: (() => void) | null = null;
		this.onReady = new Promise((_resolve, _reject) => {
			resolveReady = _resolve;
		});

		this.worker.onmessage = event => {
			if (event.data.msg == 'ready') {
				this.worker?.postMessage({ msg: 'setID', id: '1' });
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
				this.cleanupTask(id);
				promiseFuncs[0](event.data.outputData);
			}
		};
	}

	cleanupTask(id: string) {
		delete this.taskIdsToPromiseFuncs[id];
		this.numTasksRunning--;
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

	async doApiCall(requestName: string, request: Uint8Array, id: string, { signal }: WorkerOptions = {}): Promise<Uint8Array> {
		this.numTasksRunning++;
		await this.onReady;

		const taskPromise = new Promise<Uint8Array>((resolve, reject) => {
			if (!id) {
				id = this.makeTaskId();
			}
			this.taskIdsToPromiseFuncs[id] = [resolve, reject];
			signal?.addEventListener('abort', () => {
				// @TODO: Abort current sim logic here
				// Need to send a message to the worker to stop the sim by progress id.
				this.rebuildWorker(id);
			});

			this.worker?.postMessage({
				msg: requestName,
				id: id,
				inputData: request,
			});
		});
		return await taskPromise;
	}

	async rebuildWorker(id: string) {
		this.worker?.terminate();
		this.cleanupTask(id);
		this.initializeWorker();
	}
}
