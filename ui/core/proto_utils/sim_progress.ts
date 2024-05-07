import { ActionMetrics, AuraMetrics, DistributionMetrics, PartyMetrics, ProgressMetrics, RaidSimResult, ResourceMetrics, UnitMetrics } from '../proto/api';
import { lazyDeepClone } from '../utils';
import { ActionId } from './action_id';

type ConcurrentSimData = {
	concurrency: number;
	iterationsTotal: number;
	iterationsDone: number[];

	dpsValues: number[];
	hpsValues: number[];

	finalResults: RaidSimResult[];
};

const baseProgressData: ConcurrentSimData = {
	concurrency: 0,
	iterationsTotal: 0,
	iterationsDone: [],

	dpsValues: [],
	hpsValues: [],

	finalResults: [],
};

const baseDistMetric: DistributionMetrics = {
	avg: 0,
	stdev: 0,
	max: 0,
	maxSeed: BigInt(0),
	min: 0,
	minSeed: BigInt(0),
	hist: {},
	allValues: [],
};

const baseRaidSimResult: RaidSimResult = {
	raidMetrics: {
		dps: baseDistMetric,
		hps: baseDistMetric,
		parties: [],
	},
	encounterMetrics: {
		targets: [],
	},
	logs: '',
	firstIterationDuration: 0,
	avgIterationDuration: 0,
	errorResult: '',
};

const sumArrayNumberTotals = (arr: number[]): number => arr.reduce((total, a) => total + a, 0);
const bigIntMax = (...args: bigint[]) => args.reduce((m, e) => (e > m ? e : m));
const bigIntMin = (...args: bigint[]) => args.reduce((m, e) => (e < m ? e : m));

declare global {
	interface BigInt {
		/** Convert to BigInt to string form in JSON.stringify */
		toJSON: () => string;
	}
}

BigInt.prototype.toJSON = function () {
	return this.toString();
};

class SimProgress {
	data: ConcurrentSimData = lazyDeepClone(baseProgressData);
	result: RaidSimResult = lazyDeepClone(baseRaidSimResult);
	constructor(options: Partial<ConcurrentSimData>) {
		this.data = { ...this.data, ...options };

		this.data.iterationsDone = new Array(this.data.concurrency).fill(null);
		this.data.dpsValues = new Array(this.data.concurrency).fill(null);
		this.data.hpsValues = new Array(this.data.concurrency).fill(null);
		this.data.finalResults = new Array(this.data.concurrency).fill(null);
	}

	setBaseResult(result: RaidSimResult) {
		this.result = lazyDeepClone(result);

		result.raidMetrics?.parties.forEach((party, index) => {
			this.result.raidMetrics!.parties[index] = party;
		});

		result.encounterMetrics?.targets.forEach((target, index) => {
			this.result.encounterMetrics!.targets[index] = target;
		});
	}

	addResult(result: RaidSimResult, isLast: boolean, weight: number) {
		this.combineDistMetrics(this.result.raidMetrics!.dps!, result.raidMetrics!.dps!, isLast, weight);
		this.combineDistMetrics(this.result.raidMetrics!.hps!, result.raidMetrics!.hps!, isLast, weight);
		result.raidMetrics?.parties.forEach((party, partyIndex) => {
			if (!this.result.raidMetrics!.parties[partyIndex]) this.result.raidMetrics!.parties[partyIndex] = this.newPartyMetrics(party);
			const baseParty = this.result.raidMetrics!.parties[partyIndex];

			this.combineDistMetrics(baseParty.dps!, party.dps!, isLast, weight);
			this.combineDistMetrics(baseParty.hps!, party.hps!, isLast, weight);

			party.players.forEach((player, playerIndex) => {
				this.combineUnitMetrics(baseParty.players[playerIndex], player, isLast, weight);
			});
		});

		result.encounterMetrics!.targets.forEach((target, targetIndex) => {
			if (!this.result.encounterMetrics!.targets[targetIndex]) this.result.encounterMetrics!.targets[targetIndex] = this.newUnitMetrics(target);
			this.combineUnitMetrics(this.result.encounterMetrics!.targets[targetIndex], target, isLast, weight);
		});

		this.result.avgIterationDuration += result.avgIterationDuration * weight;
	}

	newDistMetrics(): DistributionMetrics {
		const copy = lazyDeepClone(baseDistMetric);
		return { ...copy, maxSeed: BigInt(copy.maxSeed), minSeed: BigInt(copy.minSeed) };
	}

	newUnitMetrics(baseUnit: UnitMetrics) {
		const newUnitMetrics: UnitMetrics = {
			name: baseUnit.name,
			unitIndex: baseUnit.unitIndex,
			dps: this.newDistMetrics(),
			dpasp: this.newDistMetrics(),
			threat: this.newDistMetrics(),
			dtps: this.newDistMetrics(),
			tmi: this.newDistMetrics(),
			hps: this.newDistMetrics(),
			tto: this.newDistMetrics(),
			secondsOomAvg: 0,
			chanceOfDeath: 0,
			actions: [],
			auras: [],
			resources: [],
			pets: [],
		};

		baseUnit.pets.forEach((pet, index) => {
			newUnitMetrics.pets[index] = this.newUnitMetrics(pet);
		});

		return newUnitMetrics;
	}

	newPartyMetrics(baseParty: PartyMetrics) {
		const newPartyMetrics: PartyMetrics = {
			dps: this.newDistMetrics(),
			hps: this.newDistMetrics(),
			players: [],
		};

		baseParty.players.forEach((player, index) => {
			newPartyMetrics.players[index] = this.newUnitMetrics(player);
		});

		return newPartyMetrics;
	}

	combineDistMetrics(base: DistributionMetrics, add: DistributionMetrics, isLast: boolean, weight: number) {
		base.avg += add.avg * weight;
		base.stdev += add.stdev * add.stdev * weight;
		if (isLast) {
			base.stdev = Math.sqrt(base.stdev);
		}

		base.max = Math.max(base.max, add.max);
		base.maxSeed = bigIntMax(base.maxSeed, add.maxSeed);

		base.min = base.min === 0 ? add.min : Math.min(base.min, add.min);
		base.minSeed = base.minSeed === BigInt(0) ? add.minSeed : bigIntMin(base.minSeed, add.minSeed);

		Object.entries(add.hist).forEach(([key, value]) => {
			const baseHist = base.hist[Number(key)];
			if (typeof baseHist !== 'number') {
				base.hist[Number(key)] = 0;
			}
			if (!isNaN(value)) base.hist[Number(key)] += value;
		});

		base.allValues.push(...add.allValues);
	}

	combineUnitMetrics(base: UnitMetrics, add: UnitMetrics, isLast: boolean, weight: number) {
		if (base.name !== add.name) {
			throw new Error('Names do not match?!');
		}

		if (base.unitIndex !== add.unitIndex) {
			throw new Error('UnitIndices do not match?!');
		}

		this.combineDistMetrics(base.dps!, add.dps!, isLast, weight);
		this.combineDistMetrics(base.dpasp!, add.dpasp!, isLast, weight);
		this.combineDistMetrics(base.threat!, add.threat!, isLast, weight);
		this.combineDistMetrics(base.dtps!, add.dtps!, isLast, weight);
		this.combineDistMetrics(base.tmi!, add.tmi!, isLast, weight);
		this.combineDistMetrics(base.hps!, add.hps!, isLast, weight);
		this.combineDistMetrics(base.tto!, add.tto!, isLast, weight);

		base.secondsOomAvg += add.secondsOomAvg * weight;
		base.chanceOfDeath += add.chanceOfDeath * weight;

		add.actions.forEach(addAction => {
			this.addActionMetrics(base, addAction);
		});

		add.auras.forEach(addAura => {
			this.addAuraMetrics(base, addAura, isLast, weight);
		});

		add.resources.forEach(addResource => {
			this.addResourceMetrics(base, addResource);
		});

		add.pets.forEach((addPet, index) => {
			this.combineUnitMetrics(base.pets[index], addPet, isLast, weight);
		});
	}

	addActionMetrics(unit: UnitMetrics, add: ActionMetrics) {
		let actionMetrics: ActionMetrics | null = null;

		const addKey = ActionId.fromProto(add.id!).toString();
		for (const baseAction of unit.actions) {
			if (ActionId.fromProto(baseAction.id!).toString() === addKey) {
				actionMetrics = baseAction;
				break;
			}
		}

		if (!actionMetrics) {
			actionMetrics = {
				id: add.id,
				isMelee: add.isMelee,
				targets: [],
			};
			add.targets.forEach((addTgt, index) => {
				actionMetrics!.targets[index] = {
					unitIndex: addTgt.unitIndex,
					blocks: 0,
					casts: 0,
					crits: 0,
					damage: 0,
					dodges: 0,
					glances: 0,
					hits: 0,
					misses: 0,
					parries: 0,
					shielding: 0,
					threat: 0,
					healing: 0,
					castTimeMs: 0,
				};
			});
			unit.actions.push(actionMetrics);
		}

		actionMetrics.targets.forEach((baseTgt, index) => {
			const addTgt = add.targets[index];
			if (baseTgt.unitIndex != addTgt.unitIndex) {
				throw new Error("UnitIndex doesn't match?!");
			}
			baseTgt.casts += addTgt.casts;
			baseTgt.hits += addTgt.hits;
			baseTgt.crits += addTgt.crits;
			baseTgt.misses += addTgt.misses;
			baseTgt.dodges += addTgt.dodges;
			baseTgt.parries += addTgt.parries;
			baseTgt.blocks += addTgt.blocks;
			baseTgt.glances += addTgt.glances;
			baseTgt.damage += addTgt.damage;
			baseTgt.threat += addTgt.threat;
			baseTgt.healing += addTgt.healing;
			baseTgt.shielding += addTgt.shielding;
			baseTgt.castTimeMs += addTgt.castTimeMs;
		});
	}

	addAuraMetrics(unit: UnitMetrics, add: AuraMetrics, isLast: boolean, weight: number) {
		let auraMetrics: AuraMetrics | null = null;

		const addKey = ActionId.fromProto(add.id!).toString();
		for (const baseAura of unit.auras) {
			if (ActionId.fromProto(baseAura.id!).toString() === addKey) {
				auraMetrics = baseAura;
				break;
			}
		}

		if (!auraMetrics) {
			auraMetrics = {
				id: add.id,
				uptimeSecondsAvg: 0,
				procsAvg: 0,
				uptimeSecondsStdev: 0,
			};
			unit.auras.push(auraMetrics);
		}

		auraMetrics.uptimeSecondsAvg += add.uptimeSecondsAvg * weight;
		auraMetrics.procsAvg += add.procsAvg * weight;
		auraMetrics.uptimeSecondsStdev += add.uptimeSecondsStdev * add.uptimeSecondsStdev * weight;
		if (isLast) {
			auraMetrics.uptimeSecondsStdev = Math.sqrt(auraMetrics.uptimeSecondsStdev);
		}
	}

	addResourceMetrics(unit: UnitMetrics, add: ResourceMetrics) {
		let resourceMetrics: ResourceMetrics | null = null;

		const addKey = ActionId.fromProto(add.id!).toString();
		for (const baseResource of unit.resources) {
			if (ActionId.fromProto(baseResource.id!).toString() === addKey) {
				resourceMetrics = baseResource;
				break;
			}
		}

		if (!resourceMetrics) {
			resourceMetrics = {
				id: add.id,
				type: add.type,
				actualGain: 0,
				events: 0,
				gain: 0,
			};
			unit.resources.push(resourceMetrics);
		}

		resourceMetrics.events += add.events;
		resourceMetrics.gain += add.gain;
		resourceMetrics.actualGain += add.actualGain;
	}

	updateProgress(index: number, metrics: ProgressMetrics) {
		if (metrics.presimRunning) return false;

		this.data.iterationsDone[index] = metrics.completedIterations;
		this.data.dpsValues[index] = metrics.dps;
		this.data.hpsValues[index] = metrics.hps;

		if (!!metrics.finalRaidResult) {
			this.data.finalResults[index] = metrics.finalRaidResult;
			return true;
		}

		return false;
	}

	getCombinedFinalResult() {
		if (this.data.concurrency == 1) {
			return this.data.finalResults[0];
		}
		this.setBaseResult(baseRaidSimResult);
		this.result.firstIterationDuration = this.data.finalResults[0].firstIterationDuration;
		this.result.logs = this.data.finalResults[0].logs;

		this.data.finalResults.forEach((result, index) => {
			const resultWeight = this.data.iterationsDone[index] / this.data.iterationsTotal;
			this.addResult(result, index === this.data.finalResults.length - 1, resultWeight);
		});
		return this.result;
	}

	getIterationsDone() {
		return sumArrayNumberTotals(this.data.iterationsDone);
	}

	getDpsAvg(): number {
		const count = this.data.dpsValues.filter(val => !!val).length;
		return sumArrayNumberTotals(this.data.dpsValues) / count;
	}

	getHpsAvg(): number {
		const count = this.data.hpsValues.filter(val => !!val).length;
		return sumArrayNumberTotals(this.data.hpsValues) / count;
	}
}

export default SimProgress;
