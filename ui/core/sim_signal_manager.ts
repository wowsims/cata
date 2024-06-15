export const enum RequestTypes {
	RaidSim = 0x1,
	StatWeights = 0x2,
	BulkSim = 0x4,
	All = 0xf,
}

class TriggerSignal {
	private readonly callbacks: (() => Promise<void>)[] = [];
	private triggered = false;

	async trigger() {
		if (this.triggered) return;
		this.triggered = true;
		return await Promise.all(this.callbacks.map(cb => cb()));
	}

	isTriggered() {
		return this.triggered;
	}

	onTrigger(cb: () => Promise<void>) {
		this.callbacks.push(cb);
		if (this.triggered) cb();
	}
}

export type SimSignals = {
	readonly id: string;
	readonly abort: TriggerSignal;
};

function newSignals(id: string): SimSignals {
	return {
		id,
		abort: new TriggerSignal(),
	};
}

export class SimSignalManager {
	private readonly running: Map<string, { type: RequestTypes; signals: SimSignals }>;

	constructor() {
		this.running = new Map<string, { type: RequestTypes; signals: SimSignals }>();
	}

	/**
	 * Create signals and register with id to make it available for manager functions.
	 * @param id The unique id used for the request.
	 * @returns The SimSignals object.
	 */
	registerRunning(id: string, type: RequestTypes): SimSignals {
		if (this.running.has(id)) throw new Error('Tried to add already existing id!');
		const signals = newSignals(id);
		this.running.set(id, { type, signals: signals });
		return signals;
	}

	/**
	 * Remove signals from internal registry.
	 * @param id The id they were created with.
	 */
	unregisterRunning(id: string) {
		this.running.delete(id);
	}

	/**
	 * Trigger abort for all registered request ids.
	 * @param typeMask Limit to specific types of requests or all requests.
	 */
	async abortType(typeMask: RequestTypes) {
		for (const info of this.running.values()) {
			if (!(info.type & typeMask)) continue;
			await info.signals.abort.trigger();
		}
	}

	/**
	 * Trigger abort for a specific request.
	 * @param requestId The request id of the request to abort.
	 * @returns True if signals for that id existed.
	 */
	async abortId(requestId: string): Promise<boolean> {
		const cfg = this.running.get(requestId);
		if (cfg) {
			await cfg.signals.abort.trigger();
			return true;
		}
		return false;
	}
}
