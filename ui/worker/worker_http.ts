import { sleep } from '../core/utils';
import { type SimRequest } from './types';
import { WorkerInterface } from './worker_interface';

const defaultRequestOptions = {
	method: 'POST',
	headers: {
		'Content-Type': 'application/x-protobuf',
	},
};

export function setupHttpWorker(baseURL: string) {

	function makeHttpApiRequest(endPoint: string, inputData: Uint8Array) {
		return fetch(`${baseURL}/${endPoint}`, {
			...defaultRequestOptions,
			body: inputData,
		});
	}

	async function syncHandler(inputData: Uint8Array, _: any, msg: SimRequest) {
		const response = await makeHttpApiRequest(msg, inputData);
		const ab = await response.arrayBuffer();
		return new Uint8Array(ab);
	}

	async function asyncHandler(inputData: Uint8Array, progress: (outputData: Uint8Array) => void, msg: SimRequest) {
		const asyncApiResult = await syncHandler(inputData, null, msg);
		let outputData = new Uint8Array();
		while (true) {
			const progressResponse = await makeHttpApiRequest("asyncProgress", asyncApiResult);

			// If no new data available, stop querying.
			if (progressResponse.status === 204) {
				break;
			}

			const ab = await progressResponse.arrayBuffer();
			outputData = new Uint8Array(ab);
			progress(outputData);
			await sleep(500);
		}
		return outputData
	}

	new WorkerInterface({
		"bulkSimAsync": asyncHandler,
		"computeStats": syncHandler,
		"computeStatsJson": syncHandler,
		"raidSim": syncHandler,
		"raidSimJson": syncHandler,
		"raidSimAsync": asyncHandler,
		"statWeights": syncHandler,
		"statWeightsAsync": asyncHandler,
	}).ready();
}
