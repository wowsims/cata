import { SimRequest } from '../worker/types';
import {
	ErrorOutcome,
	ErrorOutcomeType,
	ProgressMetrics,
	RaidSimRequest,
	RaidSimRequestSplitRequest,
	RaidSimResult,
	RaidSimResultCombinationRequest,
	StatWeightsCalcRequest,
	StatWeightsRequest,
	StatWeightsResult,
	StatWeightsStatResultData,
} from './proto/api';
import { SimSignals } from './sim_signal_manager';
import { generateRequestId, WorkerPool, WorkerProgressCallback } from './worker_pool';

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

const runSims = (
	requests: RaidSimRequest[],
	totalIterations: number,
	wp: WorkerPool,
	onProgress: WorkerProgressCallback,
	signals: SimSignals,
): Promise<SimRunResult> => {
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
				let errorResult: RaidSimResult | undefined;

				if (pm.finalRaidResult.error) {
					if (pm.finalRaidResult.error.type == ErrorOutcomeType.ErrorOutcomeError) {
						console.error(`Worker ${idx} had an error!`);
					}
					errorResult = pm.finalRaidResult;
					signals.abort.trigger();
				}

				if (errorResult || running == 0) {
					running = 0;
					const finalProgressMetrics = csp.makeProgressMetrics();
					finalProgressMetrics.finalRaidResult = errorResult;
					onProgress(finalProgressMetrics);
					resolve({
						errorResult: errorResult,
						results: csp.finalResults,
						progressMetricsFinal: finalProgressMetrics,
					});
					return;
				}
			}
		};

		for (let i = 0; i < requests.length; i++) {
			wp.raidSimAsync(requests[i], pm => progressHandler(i, pm), signals);
		}
	});
};

const makeAndSendRaidSimError = (err: string | ErrorOutcome, onProgress: WorkerProgressCallback): RaidSimResult => {
	const errRes = RaidSimResult.create();
	if (typeof err === 'string') {
		console.error(err);
		errRes.error = ErrorOutcome.create({ message: err });
	} else {
		if (err.message) console.error(err.message);
		errRes.error = err;
	}
	onProgress(ProgressMetrics.create({ finalRaidResult: errRes }));
	return errRes;
};

export const runConcurrentSim = async (
	request: RaidSimRequest,
	workerPool: WorkerPool,
	onProgress: WorkerProgressCallback,
	signals: SimSignals,
): Promise<RaidSimResult> => {
	console.log(`Sending requests split for ${workerPool.getNumWorkers()} splits.`);

	const splitResult = await workerPool.raidSimRequestSplit(
		RaidSimRequestSplitRequest.create({
			splitCount: workerPool.getNumWorkers(),
			request: request,
		}),
	);

	if (splitResult.errorResult) {
		return makeAndSendRaidSimError(splitResult.errorResult, onProgress);
	}

	if (signals.abort.isTriggered()) {
		return makeAndSendRaidSimError(ErrorOutcome.create({ type: ErrorOutcomeType.ErrorOutcomeAborted }), onProgress);
	}

	console.log(`Running ${request.simOptions!.iterations} iterations on ${splitResult.splitsDone} concurrent sims...`);

	const simRes = await runSims(splitResult.requests, request.simOptions!.iterations, workerPool, onProgress, signals);

	if (simRes.errorResult && simRes.errorResult.error) {
		return makeAndSendRaidSimError(simRes.errorResult.error, onProgress);
	}

	if (signals.abort.isTriggered()) {
		return makeAndSendRaidSimError(ErrorOutcome.create({ type: ErrorOutcomeType.ErrorOutcomeAborted }), onProgress);
	}

	console.log(`All ${splitResult.splitsDone} sims finished successfully. Combining ${simRes.results.length} results.`);

	const combiResult = await workerPool.raidSimResultCombination(
		RaidSimResultCombinationRequest.create({
			results: simRes.results,
		}),
	);

	if (combiResult.error) {
		return makeAndSendRaidSimError(combiResult.error, onProgress);
	}

	simRes.progressMetricsFinal.finalRaidResult = combiResult;
	onProgress(simRes.progressMetricsFinal);

	return combiResult;
};

const makeAndSendWeightsError = (err: string | ErrorOutcome, onProgress: WorkerProgressCallback): StatWeightsResult => {
	const errRes = RaidSimResult.create();
	if (typeof err === 'string') {
		console.error(err);
		errRes.error = ErrorOutcome.create({ message: err });
	} else {
		if (err.message) console.error(err.message);
		errRes.error = err;
	}
	onProgress(ProgressMetrics.create({ finalWeightResult: errRes }));
	return errRes;
};

export const runConcurrentStatWeights = async (
	request: StatWeightsRequest,
	workerPool: WorkerPool,
	onProgress: WorkerProgressCallback,
	signals: SimSignals,
): Promise<StatWeightsResult> => {
	console.log('Getting stat weight sim requests.');

	const id = generateRequestId(SimRequest.statWeightsAsync);

	const manualResponse = await workerPool.statWeightRequests(request);
	manualResponse.baseRequest!.requestId = id;

	if (signals.abort.isTriggered()) {
		return makeAndSendWeightsError(ErrorOutcome.create({ type: ErrorOutcomeType.ErrorOutcomeAborted }), onProgress);
	}

	let iterationsTotal = manualResponse.baseRequest!.simOptions!.iterations;
	let iterationsDone = 0;
	let simsTotal = 1;
	let simsDone = 0;

	for (const statReqData of manualResponse.statSimRequests) {
		statReqData.requestLow!.requestId = id;
		statReqData.requestHigh!.requestId = id;
		iterationsTotal += statReqData.requestLow!.simOptions!.iterations + statReqData.requestHigh!.simOptions!.iterations;
		simsTotal += 2;
	}

	console.log(`Need to run a total of ${simsTotal} sims and ${iterationsTotal} iterations.`);

	let lastIterations = 0;
	const progressHandler = (pm: ProgressMetrics) => {
		iterationsDone += pm.completedIterations - lastIterations;
		lastIterations = pm.completedIterations;

		onProgress(
			ProgressMetrics.create({
				totalIterations: iterationsTotal,
				completedIterations: iterationsDone,
				totalSims: simsTotal,
				completedSims: simsDone,
			}),
		);

		if (pm.finalRaidResult) simsDone++;
	};

	const baseLine = await runConcurrentSim(manualResponse.baseRequest!, workerPool, progressHandler, signals);
	if (baseLine.error) return makeAndSendWeightsError(baseLine.error, onProgress);

	const calcRequest = StatWeightsCalcRequest.create({
		baseResult: baseLine,
		epReferenceStat: manualResponse.epReferenceStat,
		statSimResults: [],
	});

	for (const statReqData of manualResponse.statSimRequests) {
		if (signals.abort.isTriggered()) return makeAndSendWeightsError(ErrorOutcome.create({ type: ErrorOutcomeType.ErrorOutcomeAborted }), onProgress);

		lastIterations = 0;
		const lowRes = await runConcurrentSim(statReqData.requestLow!, workerPool, progressHandler, signals);
		if (lowRes.error) return makeAndSendWeightsError(lowRes.error, onProgress);

		lastIterations = 0;
		const highRes = await runConcurrentSim(statReqData.requestHigh!, workerPool, progressHandler, signals);
		if (highRes.error) return makeAndSendWeightsError(highRes.error, onProgress);

		calcRequest.statSimResults.push(
			StatWeightsStatResultData.create({
				statData: statReqData.statData,
				resultLow: lowRes,
				resultHigh: highRes,
			}),
		);
	}

	console.log(`All ${simsTotal} sims finished successfully. Computing weights.`);

	const weightResult = await workerPool.statWeightCompute(calcRequest);
	if (weightResult.error) return makeAndSendWeightsError(weightResult.error, onProgress);
	onProgress(ProgressMetrics.create({ finalWeightResult: weightResult }));
	return weightResult;
};
