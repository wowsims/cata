import { sleep } from '../core/utils';
import { buildProgressId, defaultRequestOptions, workerPostMessage } from './shared';
import { SimRequest, type WorkerReceiveMessage } from './types';

// eslint-disable-next-line @typescript-eslint/no-unused-vars
let workerID = '';

addEventListener(
	'message',
	async ({ data }: MessageEvent<WorkerReceiveMessage>) => {
		const { msg, id, inputData } = data;

		if (msg === 'setID') {
			workerID = id;
			workerPostMessage({ msg: 'idconfirm' });
			return;
		}

		const url = new URL(msg, 'http://localhost:3333');
		const response = await fetch(url, {
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
					const progressResponse = await fetch('http://localhost:3333/asyncProgress', {
						...defaultRequestOptions,
						body: content,
					});

					// If no new data available, stop querying.
					if (progressResponse.status === 204) {
						break;
					}

					outputData = await progressResponse.arrayBuffer();
					workerPostMessage({
						id: buildProgressId(id),
						msg,
						outputData: new Uint8Array(outputData),
					});
					await sleep(500);
				}
				break;
			default:
				outputData = content;
				break;
		}

		workerPostMessage({
			id,
			msg,
			outputData: new Uint8Array(outputData!),
		});
	},
	false,
);

// Let UI know worker is ready.
workerPostMessage({
	msg: 'ready',
});

export {};
