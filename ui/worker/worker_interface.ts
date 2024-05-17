import type { SimRequest, WorkerReceiveMessage, WorkerSendMessage } from './types';

export type HandlerProgressCb = (outputData: Uint8Array) => void
export type HandlerFunc = (data: Uint8Array, progress: HandlerProgressCb, msg: SimRequest) => Uint8Array | Promise<Uint8Array>
export type HandlerDefs = Record<SimRequest, HandlerFunc>

/**
 * Communication with the UI.
 */
export class WorkerInterface {
	private workerId = "";
	private readonly handlers: HandlerDefs;

	constructor(handlers: HandlerDefs) {
		this.handlers = handlers;

		addEventListener("message", async ({ data }: MessageEvent<WorkerReceiveMessage>) => {
			const { id, msg, inputData } = data;

			if (msg === "setID") {
				this.workerId = id;
				this.postMessage({ msg: "idconfirm" });
				return;
			}

			const handlerFunc = this.handlers?.[msg];

			if (!handlerFunc) {
				console.error(`Request msg: ${msg}, id: ${id}, is not handled!`);
				return;
			}

			const progressCallback: HandlerProgressCb = prog => this.postMessage({
				msg: "progress",
				id: `${id}progress`,
				outputData: prog,
			});

			const outputData = await handlerFunc(inputData, progressCallback, msg);
			this.postMessage({ msg, id, outputData });
		});
	}

	private postMessage(m: WorkerSendMessage) {
		postMessage(m);
	}

	getWorkerId() {
		return this.workerId;
	}

	/** Tell UI that the worker is ready. */
	ready() {
		if (!this.handlers) {
			console.error("WorkerInterface.ready() used but handlers not set!");
			return;
		}
		this.postMessage({msg: "ready"});
	}
}
