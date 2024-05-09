import { sleep } from '../core/utils';
import { buildProgressId, defaultRequestOptions } from './shared';
import { SimRequest, WorkerReceiveMessage } from './types';

// eslint-disable-next-line @typescript-eslint/no-unused-vars
let workerID = '';

addEventListener(
	'message',
	async ({ data }: MessageEvent<WorkerReceiveMessage>) => {
		const { msg, id, inputData } = data;

		if (msg == 'setID') {
			workerID = id;
			postMessage({ msg: 'idconfirm' });
			return;
		}

		const response = await fetch(`/${msg}`, {
			...defaultRequestOptions,
			body: inputData,
		});

		const content = await response.arrayBuffer();
		let outputData: ArrayBuffer | null = null;
		switch (msg) {
			case SimRequest.raidSimAsync:
			case SimRequest.statWeightsAsync:
			case SimRequest.bulkSimAsync:
				while (true) {
					console.log({ msg, content });
					const progressResponse = await fetch('/asyncProgress', {
						...defaultRequestOptions,
						body: content,
					});

					// If no new data available, stop querying.
					if (progressResponse.status === 204) {
						break;
					}

					outputData = await progressResponse.arrayBuffer();
					postMessage({
						msg: msg,
						outputData: new Uint8Array(outputData),
						id: buildProgressId(id),
					});
					await sleep(500);
				}
				break;
			default:
				outputData = content;
				break;
		}

		postMessage({
			msg: msg,
			outputData: new Uint8Array(outputData!),
			id: id,
		});
	},
	false,
);

// Let UI know worker is ready.
postMessage({
	msg: 'ready',
});

export {};
