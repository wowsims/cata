import { ProgressMetrics, RaidSimRequest, RaidSimRequestSplitRequest, RaidSimResult, RaidSimResultCombinationRequest } from "./proto/api";
import { SimSignals } from "./sim_manager";
import { WorkerPool, WorkerProgressCallback } from "./worker_pool";

class ConcurrentSimProgress {
	readonly concurrency: number;
	readonly iterationsTotal: number;
	private readonly iterationsDone: number[];
	private readonly dpsValues: number[];
	private readonly hpsValues: number[];
	readonly finalResults: RaidSimResult[];

	constructor(concurrency: number, totalIterations: number) {
		this.concurrency = concurrency;
		this.iterationsTotal = totalIterations;
		this.iterationsDone = Array(this.concurrency).fill(0);
		this.dpsValues = Array(this.concurrency).fill(0);
		this.hpsValues = Array(this.concurrency).fill(0);
		this.finalResults = Array(this.concurrency);
	}

	getIterationsDone(): number {
		let total = 0;
		for (const done of this.iterationsDone) {
			total += done;
		}
		return total;
	}

	getDpsAvg(): number {
		let total = 0;
		for (const done of this.dpsValues) {
			total += done;
		}
		return total / this.concurrency;
	}

	getHpsAvg(): number {
		let total = 0;
		for (const done of this.hpsValues) {
			total += done;
		}
		return total / this.concurrency;
	}

	updateProgress(idx: number, msg: ProgressMetrics) {
		this.iterationsDone[idx] = msg.completedIterations;
		this.dpsValues[idx] = msg.dps;
		this.hpsValues[idx] = msg.hps;

		if (msg.finalRaidResult) {
			this.finalResults[idx] = msg.finalRaidResult;
		}
	}

	makeProgressMetrics(): ProgressMetrics {
		return ProgressMetrics.create({
			totalIterations: this.iterationsTotal,
			completedIterations: this.getIterationsDone(),
			dps: this.getDpsAvg(),
			hps: this.getHpsAvg(),
		});
	}
}

interface SimRunResult {
	errorResult?: RaidSimResult;
	results: RaidSimResult[];
	progressMetricsFinal: ProgressMetrics;
}

function runSims(requests: RaidSimRequest[], totalIterations: number, wp: WorkerPool, onProgress: WorkerProgressCallback, signals: SimSignals): Promise<SimRunResult> {
	return new Promise(resolve => {
		const csp = new ConcurrentSimProgress(requests.length, totalIterations);
		let progressCounter = 0;
		let running = requests.length;

		const progressHandler = (idx: number, pm: ProgressMetrics) => {
			if (!running) return;

			csp.updateProgress(idx, pm);

			progressCounter++;
			if (progressCounter % running == 0) {
				onProgress(csp.makeProgressMetrics());
			}

			if (pm.finalRaidResult) {
				running--;
				let errRes: RaidSimResult | undefined;

				if (pm.finalRaidResult.errorResult) {
					console.error(`Worker ${idx} had an error!`);
					errRes = pm.finalRaidResult;
					for (const req of requests) {
						wp.abortById(req.id);
					}
				}

				if (errRes || running == 0) {
					running = 0;
					resolve({
						errorResult: errRes,
						results: csp.finalResults,
						progressMetricsFinal: csp.makeProgressMetrics(),
					});
					return;
				}
			}
		}

		for (let i = 0; i < requests.length; i++) {
			wp.raidSimAsync(requests[i], pm => progressHandler(i, pm));
		}

		signals.abort.onTrigger(() => {
			for (const req of requests) {
				wp.abortById(req.id);
			}
		});
	});
}

function makeAndSendErrorResult(err: string, onProgress: WorkerProgressCallback): RaidSimResult {
	const errRes = RaidSimResult.create({ errorResult: err });
	onProgress(ProgressMetrics.create({ finalRaidResult: errRes }));
	console.error(err);
	return errRes;
}

export async function runConcurrentSim(request: RaidSimRequest, workerPool: WorkerPool, onProgress: WorkerProgressCallback, signals: SimSignals): Promise<RaidSimResult> {
	console.log(`Sending requests split for ${workerPool.getNumWorkers()} splits.`);

	const splitResult = await workerPool.raidSimRequestSplit(RaidSimRequestSplitRequest.create({
		splitCount: workerPool.getNumWorkers(),
		request: request,
	}));

	if (splitResult.errorResult) {
		return makeAndSendErrorResult(splitResult.errorResult, onProgress);
	}

	if (signals.abort.isTriggered()) {
		return makeAndSendErrorResult("aborted", onProgress);
	}

	console.log(`Running ${request.simOptions!.iterations} iterations on ${splitResult.splitsDone} concurrent sims...`);

	const simRes = await runSims(splitResult.requests, request.simOptions!.iterations, workerPool, onProgress, signals);

	if (simRes.errorResult) {
		console.error(simRes.errorResult.errorResult);
		return simRes.errorResult;
	}

	if (signals.abort.isTriggered()) {
		return makeAndSendErrorResult("aborted", onProgress);
	}

	console.log(`All ${splitResult.splitsDone} sims finished successfully. Combining ${simRes.results.length} results.`);

	const combiResult = await workerPool.raidSimResultCombination(RaidSimResultCombinationRequest.create({
		results: simRes.results,
	}));

	if (combiResult.errorResult) {
		return makeAndSendErrorResult(combiResult.errorResult, onProgress);
	}

	simRes.progressMetricsFinal.finalRaidResult = combiResult;
	onProgress(simRes.progressMetricsFinal);

	return combiResult;
}
