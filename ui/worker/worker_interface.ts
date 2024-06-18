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
				console.error(`Request msg: ${msg}, id: ${id}, is not handled!`);
				return;
			}

			const progressCallback: HandlerProgressCallback = prog => {
				this.postMessage({
					msg: 'progress',
					id: `${id}progress`,
					outputData: prog,
				});
			};

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

	/**
	 * Tell UI that the worker is ready.
	 * @param isWasm true if worker is using wasm.
	 */
	ready(isWasm: boolean) {
		this.postMessage({ msg: 'ready', outputData: new Uint8Array([+isWasm]) });
	}
}
