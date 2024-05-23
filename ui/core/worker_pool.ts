import { SimRequest, WorkerReceiveMessage, WorkerSendMessage } from '../worker/types';
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
	StatWeightsRequest,
	StatWeightsResult,
} from './proto/api.js';
import { noop } from './utils';

const SIM_WORKER_URL = `/${REPO_NAME}/sim_worker.js`;
export type WorkerProgressCallback = (progressMetrics: ProgressMetrics) => void;

export class WorkerPool {
	private readonly workers: Array<SimWorker>;
	private readonly workersDisabled: Array<SimWorker>;

	constructor(numWorkers: number) {
		this.workers = [];
		this.workersDisabled = [];
		this.setNumWorkers(numWorkers);
	}

	setNumWorkers(numWorkers: number) {
		if (numWorkers < this.workers.length) {
			for (let i = this.workers.length - 1; i >= numWorkers; i--) {
				this.workers[i].disable();
				this.workersDisabled[i] = this.workers[i];
			}
			this.workers.length = numWorkers;
			return;
		}

		for (let i = 0; i < numWorkers; i++) {
			if (!this.workers[i]) {
				if (this.workersDisabled[i]) {
					this.workers[i] = this.workersDisabled[i];
					delete this.workersDisabled[i];
					this.workers[i].enable();
				} else {
					this.workers[i] = new SimWorker(i);
				}
			}
		}
	}

	getNumWorkers() {
		return this.workers.length;
	}

	private getLeastBusyWorker(): SimWorker {
		return this.workers.reduce((curMinWorker, nextWorker) => (curMinWorker.getSimTaskWorkAmount() < nextWorker.getSimTaskWorkAmount() ? curMinWorker : nextWorker));
	}

	async makeApiCall(requestName: SimRequest, request: Uint8Array): Promise<Uint8Array> {
		return await this.getLeastBusyWorker().doApiCall(requestName, request, '');
	}

	async computeStats(request: ComputeStatsRequest): Promise<ComputeStatsResult> {
		const result = await this.makeApiCall(SimRequest.computeStats, ComputeStatsRequest.toBinary(request));
		return ComputeStatsResult.fromBinary(result);
	}

	private getProgressName(id: string) {
		return `${id}progress`;
	}

	async statWeightsAsync(request: StatWeightsRequest, onProgress: WorkerProgressCallback): Promise<StatWeightsResult> {
		const worker = this.getLeastBusyWorker();
		worker.log('Stat weights request: ' + StatWeightsRequest.toJsonString(request));
		const id = worker.makeTaskId();

		try {
			worker.addSimTaskRunning(id, request.simOptions ? request.simOptions.iterations * request.statsToWeigh.length : 30000);

			// Add handler for the progress events
			worker.addPromiseFunc(this.getProgressName(id), this.newProgressHandler(id, worker, onProgress), noop);

			// Now start the async sim
			const resultData = await worker.doApiCall(SimRequest.statWeightsAsync, StatWeightsRequest.toBinary(request), id);
			const result = ProgressMetrics.fromBinary(resultData);
			worker.log('Stat weights result: ' + StatWeightsResult.toJsonString(result.finalWeightResult!));
			return result.finalWeightResult!;
		} finally {
			worker.updateSimTask(id, 0);
		}
	}

	async bulkSimAsync(request: BulkSimRequest, onProgress: WorkerProgressCallback): Promise<BulkSimResult> {
		const worker = this.getLeastBusyWorker();
		worker.log('bulk sim request: ' + BulkSimRequest.toJsonString(request, { enumAsInteger: true }));
		const id = worker.makeTaskId();

		try {
			worker.addSimTaskRunning(id, request.baseSettings?.simOptions?.iterations ?? 30000);

			// Add handler for the progress events
			worker.addPromiseFunc(this.getProgressName(id), this.newProgressHandler(id, worker, onProgress), noop);

			// Now start the async sim
			const resultData = await worker.doApiCall(SimRequest.bulkSimAsync, BulkSimRequest.toBinary(request), id);
			const result = ProgressMetrics.fromBinary(resultData);
			const resultJson = BulkSimResult.toJson(result.finalBulkResult!) as any;
			worker.log('bulk sim result: ' + JSON.stringify(resultJson));
			return result.finalBulkResult!;
		} finally {
			worker.updateSimTask(id, 0);
		}
	}

	async raidSimAsync(request: RaidSimRequest, onProgress: WorkerProgressCallback): Promise<RaidSimResult> {
		const worker = this.getLeastBusyWorker();
		worker.log('Raid sim request: ' + RaidSimRequest.toJsonString(request));
		const id = worker.makeTaskId();

		try {
			worker.addSimTaskRunning(id, request.simOptions?.iterations ?? 3000);

			// Add handler for the progress events
			worker.addPromiseFunc(this.getProgressName(id), this.newProgressHandler(id, worker, onProgress), noop);

			// Now start the async sim
			const resultData = await worker.doApiCall(SimRequest.raidSimAsync, RaidSimRequest.toBinary(request), id);
			const result = ProgressMetrics.fromBinary(resultData);

			// Don't print the logs because it just clogs the console.
			const resultJson = RaidSimResult.toJson(result.finalRaidResult!) as any;
			delete resultJson!['logs'];
			worker.log('Raid sim result: ' + JSON.stringify(resultJson));
			return result.finalRaidResult!;
		} finally {
			worker.updateSimTask(id, 0);
		}
	}

	async raidSimRequestSplit(request: RaidSimRequestSplitRequest): Promise<RaidSimRequestSplitResult> {
		const result = await this.makeApiCall(SimRequest.raidSimRequestSplit, RaidSimRequestSplitRequest.toBinary(request));
		return RaidSimRequestSplitResult.fromBinary(result);
	}

	async raidSimResultCombination(request: RaidSimResultCombinationRequest): Promise<RaidSimResult> {
		const result = await this.makeApiCall(SimRequest.raidSimResultCombination, RaidSimResultCombinationRequest.toBinary(request));
		return RaidSimResult.fromBinary(result);
	}

	async isWasm() {
		return await this.workers[0].isWasmWorker();
	}

	newProgressHandler(id: string, worker: SimWorker, onProgress: WorkerProgressCallback): (progressData: any) => void {
		return (progressData: any) => {
			const progress = ProgressMetrics.fromBinary(progressData);
			onProgress(progress);
			worker.updateSimTask(id, progress.totalIterations - progress.completedIterations);
			// If we are done, stop adding the handler.
			if (progress.finalRaidResult != null || progress.finalWeightResult != null) {
				return;
			}

			worker.addPromiseFunc(this.getProgressName(id), this.newProgressHandler(id, worker, onProgress), noop);
		};
	}
}

class SimWorker {
	readonly workerId: number;
	private readonly simTasksRunning: Record<string, {workLeft: number}>;
	private taskIdsToPromiseFuncs: Record<string, [(result: any) => void, (error: any) => void]>;
	private worker: Worker | undefined;
	private onReady: Promise<void> | undefined;
	private resolveReady: (() => void) | undefined;
	private wasmWorker: boolean;
	private shouldDestroy: boolean;

	constructor(id: number) {
		this.workerId = id;
		this.simTasksRunning = {};
		this.taskIdsToPromiseFuncs = {};
		this.wasmWorker = false;
		this.shouldDestroy = false;
		this.onReady = new Promise((_resolve, _reject) => {
			this.resolveReady = _resolve;
		});
		this.setupWorker();
		this.log('Created.');
	}

	private setupWorker() {
		this.setTaskActive('setup', true); // Make it prefer ready workers.

		this.onReady = new Promise((_resolve, _reject) => {
			this.resolveReady = _resolve;
		});

		this.worker = new window.Worker(SIM_WORKER_URL);

		this.worker.addEventListener('message', ({ data }: MessageEvent<WorkerSendMessage>) => {
			const { id, msg, outputData } = data;
			switch (msg) {
				case 'ready':
					this.wasmWorker = !!outputData && !!outputData[0];
					this.postMessage({ msg: 'setID', id: this.workerId.toString() });
					this.resolveReady!();
					this.setTaskActive('setup', false);
					this.log(`Ready, isWasm: ${this.wasmWorker}`);
					break;
				case 'idConfirm':
					break;
				default:
					const promiseFuncs = this.taskIdsToPromiseFuncs[id];
					if (!promiseFuncs) {
						console.warn(`Unrecognized result id ${id} for msg ${msg}`);
						return;
					}
					if (!id.includes('progress')) this.setTaskActive(id, false);
					delete this.taskIdsToPromiseFuncs[id];
					promiseFuncs[0](outputData);
			}
		});
	}

	/** Add sim work amount (iterations) used for load balancing. */
	addSimTaskRunning(id: string, workLeft: number) {
		this.simTasksRunning[id] = {workLeft};
		this.log(`Added work ${id}, current work amount: ${this.getSimTaskWorkAmount()}`);
	}

	/** Update sim work amount (iterations left) used for load balancing. */
	updateSimTask(id: string, workLeft: number) {
		if (workLeft <= 0) {
			delete this.simTasksRunning[id];
			this.log(`Work ${id} done, current work amount: ${this.getSimTaskWorkAmount()}`);
			if (this.shouldDestroy && this.getSimTaskWorkAmount() == 0) {
				this.disable(true);
			}
			return;
		}
		this.simTasksRunning[id].workLeft = workLeft;
	}

	/** Get total iterative work left on this worker. */
	getSimTaskWorkAmount() {
		let work = 0;
		for (const t of Object.values(this.simTasksRunning)) {
			work += t.workLeft;
		}
		return work;
	}

	private setTaskActive(id: string, active: boolean) {
		if (active) {
			this.addSimTaskRunning(id + 'task', 1); // Add 1 work to track pending tasks
		} else {
			this.updateSimTask(id + 'task', 0);
		}
	}

	async isWasmWorker() {
		if (!this.onReady || this.shouldDestroy) throw new Error('Disabled worker was used!');
		await this.onReady;
		return this.wasmWorker;
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

	async doApiCall(requestName: SimRequest, request: Uint8Array, id: string): Promise<Uint8Array> {
		if (!this.onReady || this.shouldDestroy) throw new Error('Disabled worker was used!');
		if (!id) id = this.makeTaskId();
		this.setTaskActive(id, true);
		await this.onReady;

		const taskPromise = new Promise<Uint8Array>((resolve, reject) => {
			this.taskIdsToPromiseFuncs[id] = [resolve, reject];

			this.postMessage({
				msg: requestName,
				id,
				inputData: request,
			});
		});
		return await taskPromise;
	}

	postMessage(message: WorkerReceiveMessage) {
		if (!this.worker) throw new Error(`Worker ${this.workerId} postMessage while disabled!`);
		this.worker.postMessage(message);
	}

	disable(force = false) {
		this.shouldDestroy = true;
		if (!this.worker || (!force && this.getSimTaskWorkAmount())) return;
		this.worker.terminate();
		delete this.worker;
		delete this.onReady;
		delete this.resolveReady;
		this.log('Disabled.');
	}

	enable() {
		this.shouldDestroy = false;
		if (this.worker) return;
		this.setupWorker();
		this.log('Enabled.')
	}

	log(s: string) {
		console.log(`Worker ${this.workerId}: ${s}`);
	}
}
