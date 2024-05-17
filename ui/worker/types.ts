/**
 * API endpoints and exposed wasm function names. Also used as request identifier.
 */
export enum SimRequest {
	bulkSimAsync = 'bulkSimAsync',
	computeStats = 'computeStats',
	computeStatsJson = 'computeStatsJson',
	raidSim = 'raidSim',
	raidSimJson = 'raidSimJson',
	raidSimAsync = 'raidSimAsync',
	statWeights = 'statWeights',
	statWeightsAsync = 'statWeightsAsync',
}

/**
 * What the Worker receives from the UI
 */
export type WorkerReceiveMessageType = keyof typeof SimRequest | 'setID';

export interface WorkerReceiveMessageBodyBase {
	id: string;
	msg: WorkerReceiveMessageType;
	inputData?: Uint8Array;
}

export interface WorkerReceiveMessageSetId extends WorkerReceiveMessageBodyBase {
	msg: 'setID';
}

export interface WorkerReceiveMessageSimRequest extends Required<WorkerReceiveMessageBodyBase> {
	msg: SimRequest;
}

export type WorkerReceiveMessage = WorkerReceiveMessageSetId | WorkerReceiveMessageSimRequest;

/**
 * What the Worker sends to the UI
 */
export type WorkerSendMessageType = 'ready' | 'idconfirm' | 'progress' | keyof typeof SimRequest;

export interface WorkerSendMessageBodyBase {
	id?: string;
	msg: WorkerSendMessageType;
	outputData?: Uint8Array;
}

export interface WorkerSendMessageIdConfirm extends WorkerSendMessageBodyBase {
	msg: 'idconfirm';
}

export interface WorkerSendMessageReady extends WorkerSendMessageBodyBase {
	msg: 'ready';
}

export interface WorkerSendMessageProgress extends Required<WorkerSendMessageBodyBase> {
	msg: 'progress';
}

export interface WorkerSendMessageSimRequest extends Required<WorkerSendMessageBodyBase>, Required<Omit<WorkerReceiveMessageSimRequest, 'inputData'>> {
	msg: SimRequest;
}

export type WorkerSendMessage = WorkerSendMessageReady | WorkerSendMessageIdConfirm | WorkerSendMessageProgress | WorkerSendMessageSimRequest;
