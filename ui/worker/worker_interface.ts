import type { SimRequest, WorkerReceiveMessage, WorkerSendMessage } from './types';

export type HandlerProgressCallback = (outputData: Uint8Array) => void;
export type HandlerFunction = (data: Uint8Array, progress: HandlerProgressCallback, msg: SimRequest) => Uint8Array | Promise<Uint8Array>;
export type Handlers = Record<SimRequest, HandlerFunction>;

/**
 * Communication with the UI.
 */
export class WorkerInterface {
	private _workerId = '';
	private readonly handlers: Handlers;

	constructor(handlers: Handlers) {
		this.handlers = handlers;

		addEventListener('message', async ({ data }: MessageEvent<WorkerReceiveMessage>) => {
			const { id, msg, inputData } = data;

			if (msg === 'setID') {
				this._workerId = id;
				this.postMessage({ msg: 'idConfirm' });
				return;
			}

			const handlerFunc = this.handlers?.[msg];

			if (!handlerFunc) {
				console.error(`Request msg: ${msg}, id: ${this.workerId}, is not handled!`);
				return;
			}

			const progressCallback: HandlerProgressCallback = prog =>
				this.postMessage({
					msg: 'progress',
					id: `${this.workerId}progress`,
					outputData: prog,
				});

			const outputData = await handlerFunc(inputData, progressCallback, msg);
			this.postMessage({ msg, id, outputData });
		});
	}

	private postMessage(m: WorkerSendMessage) {
		postMessage(m);
	}

	get workerId() {
		return this._workerId;
	}

	/** Tell UI that the worker is ready. */
	ready() {
		if (!this.handlers) {
			console.error('WorkerInterface.ready() used but handlers not set!');
			return;
		}
		this.postMessage({ msg: 'ready' });
	}
}
