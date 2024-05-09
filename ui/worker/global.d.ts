import type { SimRequestAsync, SimRequestSync } from './types';

declare global {
	let workerID: string;
	const bulkSimAsync: SimRequestAsync;
	const computeStats: SimRequestSync;
	const computeStatsJson: SimRequestSync;
	const raidSim: SimRequestSync;
	const raidSimJson: SimRequestSync;
	const raidSimAsync: SimRequestAsync;
	const statWeights: SimRequestSync;
	const statWeightsAsync: SimRequestAsync;
}
