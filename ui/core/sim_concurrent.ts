import { ProgressMetrics, RaidSimRequest, RaidSimRequestSplitRequest, RaidSimRequestSplitResult, RaidSimResult, RaidSimResultCombinationRequest, RaidSimResultCombinationResult } from "./proto/api";
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
		for (const done of this.dpsValues){
			total += done;
		}
		return total / this.concurrency;
	}

	getHpsAvg(): number {
		let total = 0;
		for (const done of this.hpsValues){
			total += done;
		}
		return total / this.concurrency;
	}

	updateProgress(idx: number, msg: ProgressMetrics): boolean {
		this.iterationsDone[idx] = msg.completedIterations;
		this.dpsValues[idx] = msg.dps;
		this.hpsValues[idx] = msg.hps;

		if (msg.finalRaidResult) {
			this.finalResults[idx] = msg.finalRaidResult;
			return true;
		}

		return false;
	}

	makeProgressMetrics(): ProgressMetrics {
		return ProgressMetrics.create({
			totalIterations:     this.iterationsTotal,
			completedIterations: this.getIterationsDone(),
			dps:                 this.getDpsAvg(),
			hps:                 this.getHpsAvg(),
		});
	}
}

function makeAndSendErrorResult(err: string, onProgress: WorkerProgressCallback): RaidSimResult {
	const errRes = RaidSimResult.create({errorResult: err});
	onProgress(ProgressMetrics.create({finalRaidResult: errRes}));
	console.error(err);
	return errRes;
}

function splitRequest(request: RaidSimRequest, wp: WorkerPool): Promise<RaidSimRequestSplitResult> {
	return wp.raidSimRequestSplit(RaidSimRequestSplitRequest.create({
		splitCount: wp.getNumWorkers(),
		request: request,
	}));
}

function combineResults(results: RaidSimResult[], wp: WorkerPool): Promise<RaidSimResultCombinationResult> {
	return wp.raidSimResultCombination(RaidSimResultCombinationRequest.create({
		results: results,
	}));
}

interface SimRunResult {
	errorResult?: RaidSimResult;
	results: RaidSimResult[];
	progressMetricsFinal: ProgressMetrics;
}

async function runSims(requests: RaidSimRequest[], totalIterations: number, wp: WorkerPool, onProgress: WorkerProgressCallback): Promise<SimRunResult> {
	let resolve: ((r: SimRunResult) => void) | undefined = undefined;

	const csp = new ConcurrentSimProgress(requests.length, totalIterations);
	let progressCounter = 0;
	let running = requests.length;

	const progressHandler = async (idx: number, pm: ProgressMetrics) => {
		if (!resolve) return;

		if (csp.updateProgress(idx, pm)) {
			if (pm.finalRaidResult) {
				let errRes: RaidSimResult | undefined;
				if (pm.finalRaidResult.errorResult) {
					console.error(`Worker ${idx} had an error!`);
					errRes = pm.finalRaidResult;
					// This sucks, but it's better than having long running workers forever.
					if (requests[0].simOptions!.iterations > 1000) {
						console.log("Terminating all workers to get going again.");
						const num = wp.getNumWorkers()
						wp.setNumWorkers(0);
						wp.setNumWorkers(num);
					}
				}

				running--;
				if (errRes || running == 0) {
					resolve({
						errorResult: errRes,
						results: csp.finalResults,
						progressMetricsFinal: csp.makeProgressMetrics(),
					});
					resolve = undefined;
				}
			}
		}

		progressCounter++;
		if (progressCounter % csp.concurrency == 0) {
			onProgress(csp.makeProgressMetrics());
		}
	}

	for (let i = 0; i < requests.length; i++) {
		wp.raidSimAsync(requests[i], pm => { progressHandler(i, pm) });
	}

	return new Promise(res => {
		resolve = res;
	});
}

export async function runConcurrentSim(request: RaidSimRequest, workerPool: WorkerPool, onProgress: WorkerProgressCallback): Promise<RaidSimResult> {
	console.log(`Sending requests split for ${workerPool.getNumWorkers()} splits.`);
	const splitRes = await splitRequest(request, workerPool);
	if (splitRes.errorResult) return makeAndSendErrorResult(splitRes.errorResult, onProgress);
	console.log(`Got ${splitRes.splitsDone} splits.`);

	console.log(`Running ${request.simOptions?.iterations} iterations on ${splitRes.splitsDone} concurrent sims...`);
	const simRes = await runSims(splitRes.requests, request.simOptions!.iterations, workerPool, onProgress);
	if (simRes.errorResult) {
		console.error(simRes.errorResult.errorResult);
		return simRes.errorResult;
	}
	console.log(`All ${splitRes.splitsDone} sims finished successfully.`);
	onProgress(simRes.progressMetricsFinal);

	console.log(`Combining ${simRes.results.length} results.`);
	const combinationResult = await combineResults(simRes.results, workerPool);
	if (combinationResult.errorResult) makeAndSendErrorResult(combinationResult.errorResult, onProgress);
	if (!combinationResult.combinedResult) return makeAndSendErrorResult("Could not get combined result!", onProgress);
	console.log("Combination successfull!");

	simRes.progressMetricsFinal.finalRaidResult = combinationResult.combinedResult;
	onProgress(simRes.progressMetricsFinal);
	return combinationResult.combinedResult;
}
