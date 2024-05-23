import { AbortResponse } from "./proto/api";
import { WorkerPool } from "./worker_pool";

export const enum RequestType {
	RaidSim = 0x1,
	StatWeights = 0x2,
	BulkSim = 0x4,
	All = 0xF,
}

class TriggerSignal {
	private readonly callbacks: (() => void)[] = [];
	private triggered = false;

	trigger() {
		if (this.triggered) return;
		this.triggered = true;
		for (const cb of this.callbacks) {
			cb();
		}
	}

	isTriggered() {
		return this.triggered;
	}

	onTrigger(cb: () => void) {
		this.callbacks.push(cb);
		if (this.triggered) cb();
	}
}

export type SimSignals = {
	abort: TriggerSignal;
}

function newSignals(): SimSignals {
	return {
		abort: new TriggerSignal(),
	}
}

export class SimManager {
	private readonly workerPool: WorkerPool;
	private readonly running: Map<string, {type: RequestType, signals?: SimSignals}>;

	constructor(wp: WorkerPool) {
		this.workerPool = wp;
		this.running = new Map<string, {type: RequestType, signals?: SimSignals}>();
	}

	/**
	 * Add running sim. Makes it available for manager functions.
	 * @param id The unique id used for the request
	 * @param isManagedInJs Set true if immediate request is managed by JS and not in wasm or net workers.
	 * @returns Signal object to be used in managing functions if isManagedInJs is true.
	 */
	registerRunning(id: string, type: RequestType, isManagedInJs: false): void
	registerRunning(id: string, type: RequestType, isManagedInJs: true): SimSignals
	registerRunning(id: string, type: RequestType, isManagedInJs: boolean): SimSignals | void {
		if (this.running.has(id)) throw new Error("Tried to add already existing id!");

		if (isManagedInJs) {
			const sigObj = newSignals();
			this.running.set(id, {type, signals: sigObj});
			return sigObj;
		}

		this.running.set(id, {type});
	}

	unregisterRunning(id: string) {
		this.running.delete(id);
	}

	/**
	 * Trigger abort for all registered request ids.
	 * @param type Limit to specific types of requests or all requests.
	 */
	async abortAll(type: RequestType) {
		for (const [id, cfg] of this.running) {
			if ((cfg.type & type) == 0) continue;
			if (cfg.signals) {
				cfg.signals.abort.trigger();
			} else {
				await this.workerPool.abortById(id);
			}
		}
	}

	/** Trigger abort for specific id, if it is registered. */
	async abortId(id: string): Promise<AbortResponse> {
		const cfg = this.running.get(id);
		if (cfg) {
			if (cfg.signals) {
				cfg.signals.abort.trigger();
				return AbortResponse.create({id, wasTriggered: true});
			} else {
				return this.workerPool.abortById(id);
			}
		}
		return AbortResponse.create({id, wasTriggered: false});
	}
}
