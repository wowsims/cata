import { SimRequest, WorkerReceiveMessage, WorkerSendMessage } from '../worker/types';
import { REPO_NAME } from './constants/other.js';
import {
	AbortRequest,
	AbortResponse,
	BulkSimCombosRequest,
	BulkSimCombosResult,
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
	StatWeightRequestsData,
	StatWeightsCalcRequest,
	StatWeightsRequest,
	StatWeightsResult,
} from './proto/api.js';
import { SimSignals } from './sim_signal_manager';
import { isDevMode, noop } from './utils';

const SIM_WORKER_URL = `/${REPO_NAME}/sim_worker.js`;
export type WorkerProgressCallback = (progressMetrics: ProgressMetrics) => void;

/**
 * Create random id for requests.
 * @param type The request type to prepend.
 * @returns Random id in the format type-randomhex
 */
export const generateRequestId = (type: SimRequest) => {
	const chars = Array.from(Array(4)).map(() => Math.floor(Math.random() * 0x10000).toString(16));
	return type + '-' + chars.join('');
};

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
		return this.workers.reduce((curMinWorker, nextWorker) =>
			curMinWorker.getSimTaskWorkAmount() < nextWorker.getSimTaskWorkAmount() ? curMinWorker : nextWorker,
		);
	}

	async makeApiCall(requestName: SimRequest, request: Uint8Array): Promise<Uint8Array> {
		return await this.getLeastBusyWorker().doApiCall(requestName, request, generateRequestId(requestName));
	}

	async computeStats(request: ComputeStatsRequest): Promise<ComputeStatsResult> {
		const result = await this.makeApiCall(SimRequest.computeStats, ComputeStatsRequest.toBinary(request));
		return ComputeStatsResult.fromBinary(result);
	}

	private getProgressName(id: string) {
		return `${id}progress`;
	}

	async statWeightsAsync(request: StatWeightsRequest, onProgress: WorkerProgressCallback, signals: SimSignals): Promise<StatWeightsResult> {
		const worker = this.getLeastBusyWorker();
		worker.log('Stat weights request: ' + StatWeightsRequest.toJsonString(request));
		const id = generateRequestId(SimRequest.statWeightsAsync);

		signals.abort.onTrigger(async () => {
			await worker.sendAbortById(id);
		});

		const iterations = request.simOptions ? request.simOptions.iterations * request.statsToWeigh.length : 30000;
		const result = await this.doAsyncRequest(SimRequest.statWeightsAsync, StatWeightsRequest.toBinary(request), id, worker, onProgress, iterations);

		worker.log('Stat weights result: ' + StatWeightsResult.toJsonString(result.finalWeightResult!));
		return result.finalWeightResult!;
	}

	async statWeightRequests(request: StatWeightsRequest): Promise<StatWeightRequestsData> {
		const result = await this.makeApiCall(SimRequest.statWeightRequests, StatWeightsRequest.toBinary(request));
		return StatWeightRequestsData.fromBinary(result);
	}

	async statWeightCompute(request: StatWeightsCalcRequest): Promise<StatWeightsResult> {
		const result = await this.makeApiCall(SimRequest.statWeightCompute, StatWeightsCalcRequest.toBinary(request));
		return StatWeightsResult.fromBinary(result);
	}

	async bulkSimAsync(request: BulkSimRequest, onProgress: WorkerProgressCallback, signals: SimSignals): Promise<BulkSimResult> {
		const worker = this.getLeastBusyWorker();
		worker.log('bulk sim request: ' + BulkSimRequest.toJsonString(request, { enumAsInteger: true }));
		const id = generateRequestId(SimRequest.bulkSimAsync);

		signals.abort.onTrigger(async () => {
			await worker.sendAbortById(id);
		});

		const iterations = request.baseSettings?.simOptions?.iterations ?? 30000;
		const result = await this.doAsyncRequest(SimRequest.bulkSimAsync, BulkSimRequest.toBinary(request), id, worker, onProgress, iterations);

		const resultJson = BulkSimResult.toJson(result.finalBulkResult!) as any;
		worker.log('bulk sim result: ' + JSON.stringify(resultJson));
		return result.finalBulkResult!;
	}

	// Calculate combos and return counts
	async bulkSimCombosAsync(request: BulkSimCombosRequest): Promise<BulkSimCombosResult> {
		const worker = this.getLeastBusyWorker();
		worker.log('bulk sim combinations request: ' + BulkSimCombosRequest.toJsonString(request, { enumAsInteger: true }));
		const id = generateRequestId(SimRequest.bulkSimCombos);

		// Now start the async sim
		const resultData = await worker.doApiCall(SimRequest.bulkSimCombos, BulkSimCombosRequest.toBinary(request), id);
		return BulkSimCombosResult.fromBinary(resultData);
	}

	async raidSimAsync(request: RaidSimRequest, onProgress: WorkerProgressCallback, signals: SimSignals): Promise<RaidSimResult> {
		const worker = this.getLeastBusyWorker();
		worker.log('Raid sim request: ' + RaidSimRequest.toJsonString(request));
		const id = request.requestId;

		signals.abort.onTrigger(async () => {
			await worker.sendAbortById(id);
		});

		const iterations = request.simOptions?.iterations ?? 3000;
		const result = await this.doAsyncRequest(SimRequest.raidSimAsync, RaidSimRequest.toBinary(request), id, worker, onProgress, iterations);

		// Don't print the logs because it just clogs the console.
		const resultJson = RaidSimResult.toJson(result.finalRaidResult!) as any;
		delete resultJson!['logs'];
		worker.log('Raid sim result: ' + JSON.stringify(resultJson));
		return result.finalRaidResult!;
	}

	async raidSimRequestSplit(request: RaidSimRequestSplitRequest): Promise<RaidSimRequestSplitResult> {
		const result = await this.makeApiCall(SimRequest.raidSimRequestSplit, RaidSimRequestSplitRequest.toBinary(request));
		return RaidSimRequestSplitResult.fromBinary(result);
	}

	async raidSimResultCombination(request: RaidSimResultCombinationRequest): Promise<RaidSimResult> {
		const result = await this.makeApiCall(SimRequest.raidSimResultCombination, RaidSimResultCombinationRequest.toBinary(request));
		return RaidSimResult.fromBinary(result);
	}

	/**
	 * Check if workers are net workers or wasm workers.
	 * @returns True if workers are running wasm.
	 */
	isWasm() {
		return this.workers[0].isWasmWorker();
	}

	/**
	 * Start an async request, handling progress and returning the final ProgressMetrics.
	 * @param requestName
	 * @param request
	 * @param id The task id used.
	 * @param worker
	 * @param onProgress
	 * @param totalIterations Used for initial work amount tracking for worker task.
	 * @returns The final ProgressMetrics.
	 */
	private async doAsyncRequest(
		requestName: SimRequest.raidSimAsync | SimRequest.bulkSimAsync | SimRequest.statWeightsAsync,
		request: Uint8Array,
		id: string,
		worker: SimWorker,
		onProgress: WorkerProgressCallback,
		totalIterations: number,
	): Promise<ProgressMetrics> {
		try {
			worker.addSimTaskRunning(id, totalIterations);
			worker.doApiCall(requestName, request, id);
			const finalProgress: Promise<ProgressMetrics> = new Promise(resolve => {
				// Add handler for the progress events
				worker.addPromiseFunc(
					this.getProgressName(id),
					this.newProgressHandler(id, worker, onProgress, pm => resolve(pm)),
					noop,
				);
			});
			return await finalProgress;
		} finally {
			worker.updateSimTask(id, 0);
		}
	}

	private newProgressHandler(
		id: string,
		worker: SimWorker,
		onProgress: WorkerProgressCallback,
		onFinal: (pm: ProgressMetrics) => void,
	): (progressData: Uint8Array) => void {
		return (progressData: any) => {
			const progress = ProgressMetrics.fromBinary(progressData);
			onProgress(progress);
			worker.updateSimTask(id, Math.max(1, progress.totalIterations - progress.completedIterations));
			// If we are done, stop adding the handler.
			if (progress.finalRaidResult != null || progress.finalWeightResult != null || progress.finalBulkResult != null) {
				onFinal(progress);
				return;
			}

			worker.addPromiseFunc(this.getProgressName(id), this.newProgressHandler(id, worker, onProgress, onFinal), noop);
		};
	}
}

class SimWorker {
	readonly workerId: number;
	private readonly simTasksRunning: Record<string, { workLeft: number }>;
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
		this.simTasksRunning[id] = { workLeft };
		this.log(`Added work ${id}, current work amount: ${this.getSimTaskWorkAmount()}`);
	}

	/**
	 * Update sim work amount (iterations left) used for load balancing.
	 * @param id
	 * @param workLeft Set to 0 to remove the task.
	 */
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

	/** Check if worker has a running task with id. */
	hasTaskId(id: string) {
		return !!this.simTasksRunning[id];
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

	addPromiseFunc(id: string, callback: (result: Uint8Array) => void, onError: (error: any) => void) {
		this.taskIdsToPromiseFuncs[id] = [callback, onError];
	}

	async doApiCall(requestName: SimRequest, request: Uint8Array, id: string): Promise<Uint8Array> {
		if (!this.onReady || this.shouldDestroy) throw new Error('Disabled worker was used!');
		if (!id) throw new Error('ApiCall with empty id!');
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
		this.log('Enabled.');
	}

	log(s: string) {
		if (isDevMode()) console.log(`Worker ${this.workerId}: ${s}`);
	}

	async sendAbortById(requestId: string) {
		const abortReqBinary = AbortRequest.toBinary(AbortRequest.create({ requestId }));
		const rid = generateRequestId(SimRequest.abortById);
		const result = await this.doApiCall(SimRequest.abortById, abortReqBinary, rid);
		return AbortResponse.fromBinary(result);
	}
}
