import { WorkerSendMessage } from './types';

export const defaultRequestOptions = {
	method: 'POST',
	headers: {
		'Content-Type': 'application/x-protobuf',
	},
};

export const workerPostMessage = (data: WorkerSendMessage) => postMessage(data);

export const buildProgressId = (id: string) => `${id}progress`;
