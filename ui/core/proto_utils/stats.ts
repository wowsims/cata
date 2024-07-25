import * as Mechanics from '../constants/mechanics.js';
import { Class, PseudoStat, Stat, UnitStats } from '../proto/common.js';
import { getEnumValues } from '../utils.js';
import { getClassStatName, pseudoStatNames } from './names.js';
import { migrateOldProto, ProtoConversionMap } from './utils.js';
import { CURRENT_API_VERSION } from '../constants/other.js';

const STATS_LEN = getEnumValues(Stat).length;
const PSEUDOSTATS_LEN = getEnumValues(PseudoStat).length;

export class UnitStat {
	private readonly stat: Stat | null;
	private readonly pseudoStat: PseudoStat | null;

	private constructor(stat: Stat | null, pseudoStat: PseudoStat | null) {
		this.stat = stat;
		this.pseudoStat = pseudoStat;
	}

	isStat(): boolean {
		return this.stat != null;
	}
	isPseudoStat(): boolean {
		return this.pseudoStat != null;
	}

	getStat(): Stat {
		if (!this.isStat()) {
			throw new Error('Not a stat!');
		}
		return this.stat!;
	}
	getPseudoStat(): PseudoStat {
		if (!this.isPseudoStat()) {
			throw new Error('Not a pseudo stat!');
		}
		return this.pseudoStat!;
	}

	equals(other: UnitStat): boolean {
		return this.stat == other.stat && this.pseudoStat == other.pseudoStat;
	}

	getName(clazz: Class): string {
		if (this.isStat()) {
			return getClassStatName(this.stat!, clazz);
		} else {
			return pseudoStatNames.get(this.pseudoStat!)!;
		}
	}

	getProtoValue(proto: UnitStats): number {
		if (this.isStat()) {
			return proto.stats[this.stat!];
		} else {
			return proto.pseudoStats[this.pseudoStat!];
		}
	}

	setProtoValue(proto: UnitStats, val: number) {
		if (this.isStat()) {
			proto.stats[this.stat!] = val;
		} else {
			proto.pseudoStats[this.pseudoStat!] = val;
		}
	}

	static fromStat(stat: Stat): UnitStat {
		return new UnitStat(stat, null);
	}
	static fromPseudoStat(pseudoStat: PseudoStat): UnitStat {
		return new UnitStat(null, pseudoStat);
	}

	static getAll(): Array<UnitStat> {
		const allStats = getEnumValues(Stat) as Array<Stat>;
		const allPseudoStats = getEnumValues(PseudoStat) as Array<PseudoStat>;
		return [allStats.map(stat => UnitStat.fromStat(stat)), allPseudoStats.map(stat => UnitStat.fromPseudoStat(stat))].flat();
	}
}

/**
 * Represents values for all character stats (stam, agi, spell power, hit raiting, etc).
 *
 * This is an immutable type.
 */
export class Stats {
	private readonly stats: Array<number>;
	private readonly pseudoStats: Array<number>;

	constructor(stats?: Array<number>, pseudoStats?: Array<number>) {
		this.stats = Stats.initStatsArray(STATS_LEN, stats);
		this.pseudoStats = Stats.initStatsArray(PSEUDOSTATS_LEN, pseudoStats);
	}

	private static initStatsArray(expectedLen: number, newStats?: Array<number>): Array<number> {
		let stats = newStats?.slice(0, expectedLen) || [];

		if (stats.length < expectedLen) {
			stats = stats.concat(new Array(expectedLen - (newStats?.length || 0)).fill(0));
		}

		for (let i = 0; i < expectedLen; i++) {
			if (stats[i] == null) stats[i] = 0;
		}
		return stats;
	}

	equals(other: Stats): boolean {
		return (
			this.stats.every((newStat, statIdx) => newStat == other.getStat(statIdx)) &&
			this.pseudoStats.every((newStat, statIdx) => newStat == other.getPseudoStat(statIdx))
		);
	}

	getStat(stat: Stat): number {
		return this.stats[stat];
	}
	getPseudoStat(stat: PseudoStat): number {
		return this.pseudoStats[stat];
	}
	getUnitStat(stat: UnitStat): number {
		if (stat.isStat()) {
			return this.stats[stat.getStat()];
		} else {
			return this.pseudoStats[stat.getPseudoStat()];
		}
	}

	withStat(stat: Stat, value: number): Stats {
		const newStats = this.stats.slice();
		newStats[stat] = value;
		return new Stats(newStats, this.pseudoStats);
	}
	withPseudoStat(stat: PseudoStat, value: number): Stats {
		const newStats = this.pseudoStats.slice();
		newStats[stat] = value;
		return new Stats(this.stats, newStats);
	}
	withUnitStat(stat: UnitStat, value: number): Stats {
		if (stat.isStat()) {
			return this.withStat(stat.getStat(), value);
		} else {
			return this.withPseudoStat(stat.getPseudoStat(), value);
		}
	}

	addStat(stat: Stat, value: number): Stats {
		return this.withStat(stat, this.getStat(stat) + value);
	}

	add(other: Stats): Stats {
		return new Stats(
			this.stats.map((value, stat) => value + other.stats[stat]),
			this.pseudoStats.map((value, stat) => value + other.pseudoStats[stat]),
		);
	}

	subtract(other: Stats): Stats {
		return new Stats(
			this.stats.map((value, stat) => value - other.stats[stat]),
			this.pseudoStats.map((value, stat) => value - other.pseudoStats[stat]),
		);
	}

	scale(scalar: number): Stats {
		return new Stats(
			this.stats.map((value, _stat) => value * scalar),
			this.pseudoStats.map((value, _stat) => value * scalar),
		);
	}

	computeEP(epWeights: Stats): number {
		let total = 0;
		this.stats.forEach((stat, idx) => {
			total += stat * epWeights.stats[idx];
		});
		this.pseudoStats.forEach((stat, idx) => {
			total += stat * epWeights.pseudoStats[idx];
		});
		return total;
	}

	belowCaps(statCaps: Stats): boolean {
		for (const [idx, stat] of this.stats.entries()) {
			if (statCaps.stats[idx] > 0 && stat > statCaps.stats[idx]) {
				return false;
			}
		}

		return true;
	}

	getHasteMultipliers(playerClass: Class): number[] {
		const baseMeleeHasteMultiplier = 1 + this.getStat(Stat.StatMeleeHaste) / (Mechanics.HASTE_RATING_PER_HASTE_PERCENT * 100);
		const meleeHasteBuffsMultiplier =
			playerClass == Class.ClassHunter
				? this.getPseudoStat(PseudoStat.PseudoStatRangedSpeedMultiplier)
				: this.getPseudoStat(PseudoStat.PseudoStatMeleeSpeedMultiplier);
		const baseSpellHasteMultiplier = 1 + this.getStat(Stat.StatSpellHaste) / (Mechanics.HASTE_RATING_PER_HASTE_PERCENT * 100);
		const spellHasteBuffsMultiplier = this.getPseudoStat(PseudoStat.PseudoStatCastSpeedMultiplier);
		return [baseMeleeHasteMultiplier, meleeHasteBuffsMultiplier, baseSpellHasteMultiplier, spellHasteBuffsMultiplier];
	}

	// Apply any multiplicative Haste buffs stored via PseudoStats to the Stats entries for MeleeHaste and SpellHaste
	withHasteMultipliers(playerClass: Class): Stats {
		const [baseMeleeMulti, meleeBuffsMulti, baseSpellMulti, spellBuffsMulti] = this.getHasteMultipliers(playerClass);
		const newStats = this.stats.slice();
		newStats[Stat.StatMeleeHaste] = (baseMeleeMulti * meleeBuffsMulti - 1) * 100 * Mechanics.HASTE_RATING_PER_HASTE_PERCENT;
		newStats[Stat.StatSpellHaste] = (baseSpellMulti * spellBuffsMulti - 1) * 100 * Mechanics.HASTE_RATING_PER_HASTE_PERCENT;
		return new Stats(newStats, this.pseudoStats);
	}

	// Assumes that Haste multipliers have already been applied to both Stats arrays
	computeStatCapsDelta(statCaps: Stats, playerClass: Class): Stats {
		const [_finalMeleeHasteMulti, meleeHasteBuffsMulti, _finalSpellHasteMulti, spellHasteBuffsMulti] = this.getHasteMultipliers(playerClass);
		return new Stats(
			this.stats.map((value, stat) => {
				if (statCaps.stats[stat] > 0) {
					let statDelta = statCaps.stats[stat] - value;

					if (stat == Stat.StatMeleeHaste) {
						statDelta /= meleeHasteBuffsMulti;
					} else if (stat == Stat.StatSpellHaste) {
						statDelta /= spellHasteBuffsMulti;
					}

					return statDelta;
				}

				return 0;
			}),
			this.pseudoStats,
		);
	}

	asArray(): Array<number> {
		return this.stats.slice();
	}

	toJson(): object {
		return UnitStats.toJson(this.toProto()) as object;
	}

	toProto(): UnitStats {
		return UnitStats.create({
			stats: this.stats.slice(),
			pseudoStats: this.pseudoStats.slice(),
			apiVersion: CURRENT_API_VERSION,
		});
	}

	static fromJson(obj: any): Stats {
		return Stats.fromProto(UnitStats.fromJson(obj));
	}

	static fromMap(statsMap: Partial<Record<Stat, number>>, pseudoStatsMap?: Partial<Record<PseudoStat, number>>): Stats {
		const statsArr = new Array(STATS_LEN).fill(0);
		Object.entries(statsMap).forEach(entry => {
			const [statStr, value] = entry;
			statsArr[Number(statStr)] = value;
		});

		const pseudoStatsArr = new Array(PSEUDOSTATS_LEN).fill(0);
		if (pseudoStatsMap) {
			Object.entries(pseudoStatsMap).forEach(entry => {
				const [pseudoStatstr, value] = entry;
				pseudoStatsArr[Number(pseudoStatstr)] = value;
			});
		}

		return new Stats(statsArr, pseudoStatsArr);
	}

	static fromProto(unitStats?: UnitStats): Stats {
		if (unitStats) {
			// Fix out of-date protos before importing
			if (unitStats.apiVersion < CURRENT_API_VERSION) {
				Stats.updateProtoVersion(unitStats);
			}

			return new Stats(unitStats.stats, unitStats.pseudoStats);
		} else {
			return new Stats();
		}
	}

	static updateProtoVersion(proto: UnitStats) {
		// First migrate the stats array.
		proto.stats = Stats.migrateStatsArray(proto.stats, proto.apiVersion);

		// Any other required data migration code (such as for the
		// pseudoStats array) should go here.

		// Flag the version as up-to-date once all migrations are done.
		proto.apiVersion = CURRENT_API_VERSION;
	}

	// Takes in a stats array that was generated from an out-of-date proto version, and
	// converts it to an array that is consistent with the current proto version.
	static migrateStatsArray(oldStats: Array<number>, oldApiVersion: number): Array<number> {
		const conversionMap: ProtoConversionMap<Array<number>> = new Map([
			[1, (oldArray: Array<number>) => {
				// Revision 1 simply re-orders the stats for clarity
				const newIndices = [0, 1, 2, 3, 4, 17, 18, 6, 8, 10, 19, 15, 5, 7, 9, 11, 29, 26, 16, 30, 12, 13, 20, 28, 21, 22, 23, 24, 25, 27, 14];
				const newArray: Array<number> = new Array(oldArray.length);
				oldArray.forEach((value, idx) => {
					newArray[newIndices[idx]] = value;
				});
				return newArray;
			}],
		]);
		return migrateOldProto<Array<number>>(oldStats, oldApiVersion, conversionMap);
	}
}

export const statToPercentageOrPoints = (stat: Stat, value: number, stats: Stats) => {
	let statInPercentage: number | null = null;
	switch (stat) {
		case Stat.StatMeleeHit:
			statInPercentage = value / Mechanics.MELEE_HIT_RATING_PER_HIT_CHANCE;
			break;
		case Stat.StatSpellHit:
			statInPercentage = value / Mechanics.SPELL_HIT_RATING_PER_HIT_CHANCE;
			break;
		case Stat.StatMeleeCrit:
		case Stat.StatSpellCrit:
			statInPercentage = value / Mechanics.SPELL_CRIT_RATING_PER_CRIT_CHANCE;
			break;
		case Stat.StatMeleeHaste:
			statInPercentage = value / Mechanics.HASTE_RATING_PER_HASTE_PERCENT;
			break;
		case Stat.StatSpellHaste:
			statInPercentage = value / Mechanics.HASTE_RATING_PER_HASTE_PERCENT;
			break;
		case Stat.StatExpertise:
			statInPercentage = value / Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION / 4;
			break;
		case Stat.StatBlock:
			statInPercentage = value / Mechanics.BLOCK_RATING_PER_BLOCK_CHANCE + 5.0;
			break;
		case Stat.StatDodge:
			statInPercentage = stats.getPseudoStat(PseudoStat.PseudoStatDodge) * 100;
			break;
		case Stat.StatParry:
			statInPercentage = stats.getPseudoStat(PseudoStat.PseudoStatParry) * 100;
			break;
		case Stat.StatResilience:
			statInPercentage = value / Mechanics.RESILIENCE_RATING_PER_CRIT_REDUCTION_CHANCE;
			break;
		case Stat.StatMastery:
			statInPercentage = value / Mechanics.MASTERY_RATING_PER_MASTERY_POINT;
			break;
		default:
			statInPercentage = value;
			break;
	}
	return statInPercentage;
};

export const statPercentageOrPointsToNumber = (stat: Stat, value: number, stats: Stats) => {
	let statInPoints: number | null = null;
	switch (stat) {
		case Stat.StatMeleeHit:
			statInPoints = value * Mechanics.MELEE_HIT_RATING_PER_HIT_CHANCE;
			break;
		case Stat.StatSpellHit:
			statInPoints = value * Mechanics.SPELL_HIT_RATING_PER_HIT_CHANCE;
			break;
		case Stat.StatMeleeCrit:
		case Stat.StatSpellCrit:
			statInPoints = value * Mechanics.SPELL_CRIT_RATING_PER_CRIT_CHANCE;
			break;
		case Stat.StatMeleeHaste:
			statInPoints = value * Mechanics.HASTE_RATING_PER_HASTE_PERCENT;
			break;
		case Stat.StatSpellHaste:
			statInPoints = value * Mechanics.HASTE_RATING_PER_HASTE_PERCENT;
			break;
		case Stat.StatExpertise:
			statInPoints = value * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION * 4;
			break;
		case Stat.StatBlock:
			statInPoints = value * Mechanics.BLOCK_RATING_PER_BLOCK_CHANCE - 5.0;
			break;
		case Stat.StatDodge:
			statInPoints = stats.getPseudoStat(PseudoStat.PseudoStatDodge) / 100;
			break;
		case Stat.StatParry:
			statInPoints = stats.getPseudoStat(PseudoStat.PseudoStatParry) / 100;
			break;
		case Stat.StatResilience:
			statInPoints = value * Mechanics.RESILIENCE_RATING_PER_CRIT_REDUCTION_CHANCE;
			break;
		case Stat.StatMastery:
			statInPoints = value * Mechanics.MASTERY_RATING_PER_MASTERY_POINT;
			break;
		default:
			statInPoints = value;
			break;
	}
	return statInPoints;
};
