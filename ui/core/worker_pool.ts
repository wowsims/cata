import { REPO_NAME } from './constants/other.js';
import {
	BulkSimRequest,
	BulkSimResult,
	ComputeStatsRequest,
	ComputeStatsResult,
	ProgressMetrics,
	RaidSimRequest,
	RaidSimRequestSplitRequest,
	RaidSimRequestSplitResult,
	RaidSimResult,
	RaidSimResultCombinationRequest,
	RaidSimResultCombinationResult,
	StatWeightsRequest,
	StatWeightsResult,
} from './proto/api.js';
import { noop } from './utils';

const SIM_WORKER_URL = `/${REPO_NAME}/sim_worker.js`;

export type WorkerProgressCallback = (progressMetrics: ProgressMetrics) => void;

export class WorkerPool {
	private workers: Array<SimWorker>;

	constructor(numWorkers: number) {
		this.workers = [];
		this.setNumWorkers(numWorkers);
	}

	setNumWorkers(numWorkers: number) {
		if (numWorkers < this.workers.length) {
			for (let i = this.workers.length - 1; i >= numWorkers; i--) {
				this.workers[i].destroy();
			}
			this.workers.length = numWorkers;
			return;
		}

		for (let i = 0; i < numWorkers; i++) {
			if (!this.workers[i]) {
				this.workers[i] = new SimWorker(i);
			}
		}
	}

	getNumWorkers() {
		return this.workers.length;
	}

	private getLeastBusyWorker(): SimWorker {
		return this.workers.reduce((curMinWorker, nextWorker) => (curMinWorker.numTasksRunning < nextWorker.numTasksRunning ? curMinWorker : nextWorker));
	}

	async makeApiCall(requestName: string, request: Uint8Array): Promise<Uint8Array> {
		return await this.getLeastBusyWorker().doApiCall(requestName, request, '');
	}

	async computeStats(request: ComputeStatsRequest): Promise<ComputeStatsResult> {
		const result = await this.makeApiCall('computeStats', ComputeStatsRequest.toBinary(request));
		return ComputeStatsResult.fromBinary(result);
	}

	private getProgressName(id: string) {
		return `${id}progress`;
	}

	async statWeightsAsync(request: StatWeightsRequest, onProgress: WorkerProgressCallback): Promise<StatWeightsResult> {
		console.log('Stat weights request: ' + StatWeightsRequest.toJsonString(request));
		const worker = this.getLeastBusyWorker();
		const id = worker.makeTaskId();
		// Add handler for the progress events
		worker.addPromiseFunc(this.getProgressName(id), this.newProgressHandler(id, worker, onProgress), noop);

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
		worker.addPromiseFunc(this.getProgressName(id), this.newProgressHandler(id, worker, onProgress), noop);

		// Now start the async sim
		const resultData = await worker.doApiCall('bulkSimAsync', BulkSimRequest.toBinary(request), id);
		const result = ProgressMetrics.fromBinary(resultData);
		const resultJson = BulkSimResult.toJson(result.finalBulkResult!) as any;
		console.log('bulk sim result: ' + JSON.stringify(resultJson));
		return result.finalBulkResult!;
	}

	async raidSimAsync(request: RaidSimRequest, onProgress: WorkerProgressCallback): Promise<RaidSimResult> {
		console.log('Raid sim request: ' + RaidSimRequest.toJsonString(request));
		const worker = this.getLeastBusyWorker();
		const id = worker.makeTaskId();
		// Add handler for the progress events
		worker.addPromiseFunc(this.getProgressName(id), this.newProgressHandler(id, worker, onProgress), noop);

		console.log(`Running raid sim on worker ${worker.id}`);

		// Now start the async sim
		const resultData = await worker.doApiCall('raidSimAsync', RaidSimRequest.toBinary(request), id);
		const result = ProgressMetrics.fromBinary(resultData);

		// Don't print the logs because it just clogs the console.
		const resultJson = RaidSimResult.toJson(result.finalRaidResult!) as any;
		delete resultJson!['logs'];
		console.log('Raid sim result: ' + JSON.stringify(resultJson));
		return result.finalRaidResult!;
	}

	async raidSimRequestSplit(request: RaidSimRequestSplitRequest): Promise<RaidSimRequestSplitResult> {
		const result = await this.makeApiCall('raidSimRequestSplit', RaidSimRequestSplitRequest.toBinary(request));
		return RaidSimRequestSplitResult.fromBinary(result);
	}

	async raidSimResultCombination(request: RaidSimResultCombinationRequest): Promise<RaidSimResultCombinationResult> {
		const result = await this.makeApiCall('raidSimResultCombination', RaidSimResultCombinationRequest.toBinary(request));
		return RaidSimResultCombinationResult.fromBinary(result);
	}

	newProgressHandler(id: string, worker: SimWorker, onProgress: WorkerProgressCallback): (progressData: any) => void {
		return (progressData: any) => {
			const progress = ProgressMetrics.fromBinary(progressData);
			onProgress(progress);
			// If we are done, stop adding the handler.
			if (progress.finalRaidResult != null || progress.finalWeightResult != null) {
				return;
			}

			worker.addPromiseFunc(this.getProgressName(id), this.newProgressHandler(id, worker, onProgress), noop);
		};
	}
}

class SimWorker {
	readonly id: number;
	numTasksRunning: number;
	private taskIdsToPromiseFuncs: Record<string, [(result: any) => void, (error: any) => void]>;
	private eventIdsToPromiseFuncs: Record<string, [(result: any) => void, (error: any) => void]>;
	private worker: Worker;
	private onReady: Promise<void>;

	constructor(id: number) {
		this.id = id;
		this.numTasksRunning = 0;
		this.taskIdsToPromiseFuncs = {};
		this.eventIdsToPromiseFuncs = {};
		this.worker = new window.Worker(SIM_WORKER_URL);

		let resolveReady: (() => void) | null = null;
		this.onReady = new Promise((_resolve, _reject) => {
			resolveReady = _resolve;
		});

		this.worker.onmessage = event => {
			if (event.data.msg == 'ready') {
				this.worker.postMessage({ msg: 'setID', id: '1' });
				resolveReady!();
			} else if (event.data.msg == 'idconfirm') {
				// Do nothing
			} else {
				const id = event.data.id;
				let promiseFuncs: [(result: any) => void, (error: any) => void] | undefined;

				if (this.taskIdsToPromiseFuncs[id]) {
					promiseFuncs = this.taskIdsToPromiseFuncs[id];
					delete this.taskIdsToPromiseFuncs[id];
					this.numTasksRunning--;
					if (this.numTasksRunning < 0) {
						alert(`Worker ${this.id} API response ${id} caused numTasksRunning to become negative!`);
					}
				} else if (this.eventIdsToPromiseFuncs[id]) {
					promiseFuncs = this.eventIdsToPromiseFuncs[id];
					delete this.eventIdsToPromiseFuncs[id];
				}

				if (!promiseFuncs) {
					console.warn('Unrecognized result id: ' + id);
					return;
				}

				promiseFuncs[0](event.data.outputData);
			}
		};
	}

	addPromiseFunc(id: string, callback: (result: any) => void, onError: (error: any) => void) {
		this.eventIdsToPromiseFuncs[id] = [callback, onError];
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

	destroy() {
		this.worker.terminate();
	}
}
