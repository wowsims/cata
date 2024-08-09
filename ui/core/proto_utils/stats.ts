import * as Mechanics from '../constants/mechanics.js';
import { CURRENT_API_VERSION } from '../constants/other.js';
import { Class, PseudoStat, Stat, UnitStats } from '../proto/common.js';
import { StatCapConfig, StatCapType, UIStat as UnitStatProto } from '../proto/ui.js';
import { getEnumValues } from '../utils.js';
import { getStatName, getClassPseudoStatName } from './names.js';
import { migrateOldProto, ProtoConversionMap } from './utils.js';

const STATS_LEN = getEnumValues(Stat).length;
const PSEUDOSTATS_LEN = getEnumValues(PseudoStat).length;

export class UnitStat {
	private readonly stat: Stat | null;
	private readonly pseudoStat: PseudoStat | null;

	// Used to link a "child" PseudoStat like PhysicalHitPercent to a
	// "parent" Stat like HitRating, so that both values can be displayed
	// together in the character sheet.
	private readonly rootStat: Stat | null;

	private constructor(stat: Stat | null, pseudoStat: PseudoStat | null, rootStat: Stat | null) {
		this.stat = stat;
		this.pseudoStat = pseudoStat;
		this.rootStat = rootStat;
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

	equalsStat(other: Stat): boolean {
		return this.isStat() && (this.stat == other);
	}

	equalsPseudoStat(other: PseudoStat): boolean {
		return this.isPseudoStat() && (this.pseudoStat == other);
	}

	linkedToStat(other: Stat): boolean {
		return (this.stat == other) || (this.rootStat == other);
	}

	getFullName(playerClass: Class): string {
		if (this.isStat()) {
			return getStatName(this.stat!);
		} else {
			return getClassPseudoStatName(this.getPseudoStat(), playerClass);
		}
	}

	getShortName(playerClass: Class): string {
		const fullName = this.getFullName(playerClass);
		return fullName.replace(' Rating', '').replace(' Percent', '');
	}

	// Convert a UnitStat value from its Rating representation to a percentage representation
	// (0-100). If a percentage representation does not make sense for the stat in question
	// (Strength for example), then null is returned. Mastery is special cased to return
	// Mastery points rather than %.
	getPercentOrPointsValue(ratingValue: number): number | null {
		if (this.linkedToStat(Stat.StatCritRating)) {
			return ratingValue / Mechanics.CRIT_RATING_PER_CRIT_PERCENT;
		} else if (this.linkedToStat(Stat.StatHasteRating)) {
			return ratingValue / Mechanics.HASTE_RATING_PER_HASTE_PERCENT;
		} else if (this.equalsStat(Stat.StatExpertiseRating)) {
			return ratingValue / Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION / 4;
		} else if (this.linkedToStat(Stat.StatDodgeRating)) {
			return ratingValue / Mechanics.DODGE_RATING_PER_DODGE_PERCENT;
		} else if (this.linkedToStat(Stat.StatParryRating)) {
			return ratingValue / Mechanics.PARRY_RATING_PER_PARRY_PERCENT;
		} else if (this.equalsStat(Stat.StatMasteryRating)) {
			return ratingValue / Mechanics.MASTERY_RATING_PER_MASTERY_POINT;
		} else if (this.equalsPseudoStat(PseudoStat.PseudoStatPhysicalHitPercent)) {
			return ratingValue / Mechanics.PHYSICAL_HIT_RATING_PER_HIT_PERCENT;
		} else if (this.equalsPseudoStat(PseudoStat.PseudoStatSpellHitPercent)) {
			return ratingValue / Mechanics.SPELL_HIT_RATING_PER_HIT_PERCENT;
		} else {
			return null;
		}
	}

	// Convert a UnitStat value from its percentage representation (0-100) to the equivalent amount of
	// Rating. If a Rating representation does not make sense for the stat in question (Block in Cata
	// for example), then null is returned. Mastery is special cased to assume a Mastery points input
	// rather than a percentage.
	getRatingValue(percentOrPointsValue: number): number | null {
		if (this.linkedToStat(Stat.StatCritRating)) {
			return percentOrPointsValue * Mechanics.CRIT_RATING_PER_CRIT_PERCENT;
		} else if (this.linkedToStat(Stat.StatHasteRating)) {
			return percentOrPointsValue * Mechanics.HASTE_RATING_PER_HASTE_PERCENT;
		} else if (this.equalsStat(Stat.StatExpertiseRating)) {
			return percentOrPointsValue * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION * 4;
		} else if (this.linkedToStat(Stat.StatDodgeRating)) {
			return percentOrPointsValue * Mechanics.DODGE_RATING_PER_DODGE_PERCENT;
		} else if (this.linkedToStat(Stat.StatParryRating)) {
			return percentOrPointsValue * Mechanics.PARRY_RATING_PER_PARRY_PERCENT;
		} else if (this.equalsStat(Stat.StatMasteryRating)) {
			return percentOrPointsValue * Mechanics.MASTERY_RATING_PER_MASTERY_POINT;
		} else if (this.equalsPseudoStat(PseudoStat.PseudoStatPhysicalHitPercent)) {
			return percentOrPointsValue * Mechanics.PHYSICAL_HIT_RATING_PER_HIT_PERCENT;
		} else if (this.equalsPseudoStat(PseudoStat.PseudoStatSpellHitPercent)) {
			return percentOrPointsValue * Mechanics.SPELL_HIT_RATING_PER_HIT_PERCENT;
		} else {
			return null;
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

	toProto(): UnitStatProto {
		const protoMessage = UnitStatProto.create({});

		if (this.isStat()) {
			protoMessage.unitStat = {
				oneofKind: 'stat',
				stat: this.stat!,
			};
		} else if (this.isPseudoStat()) {
			protoMessage.unitStat = {
				oneofKind: 'pseudoStat',
				pseudoStat: this.pseudoStat!,
			};
		} else {
			throw new Error('Neither a Stat nor a PseudoStat!');
		}

		return protoMessage;
	}

	static fromProto(protoMessage: UnitStatProto): UnitStat {
		if (protoMessage.unitStat.oneofKind == 'stat') {
			return UnitStat.fromStat(protoMessage.unitStat.stat);
		} else if (protoMessage.unitStat.oneofKind == 'pseudoStat') {
			return UnitStat.fromPseudoStat(protoMessage.unitStat.pseudoStat);
		} else {
			return new UnitStat(null, null, null);
		}
	}

	static fromStat(stat: Stat): UnitStat {
		return new UnitStat(stat, null, null);
	}
	static fromPseudoStat(pseudoStat: PseudoStat): UnitStat {
		return new UnitStat(null, pseudoStat, UnitStat.getRootStat(pseudoStat));
	}

	static getAll(): Array<UnitStat> {
		const allStats = getEnumValues(Stat) as Array<Stat>;
		const allPseudoStats = getEnumValues(PseudoStat) as Array<PseudoStat>;
		return [allStats.map(stat => UnitStat.fromStat(stat)), allPseudoStats.map(stat => UnitStat.fromPseudoStat(stat))].flat();
	}

	// Returns the "parent" Stat (such as HitRating) associated with a
	// "child" PseudoStat (such as PhysicalHitPercent), or null if there is
	// no such root Stat.
	static getRootStat(pseudoStat: PseudoStat): Stat | null {
		const pseudoStatName = PseudoStat[pseudoStat];

		if (pseudoStatName.includes('Dodge')) {
			return Stat.StatDodgeRating;
		} else if (pseudoStatName.includes('Parry')) {
			return Stat.StatParryRating;
		} else if (pseudoStatName.includes('Haste')) {
			return Stat.StatHasteRating;
		} else if (pseudoStatName.includes('Hit')) {
			return Stat.StatHitRating;
		} else if (pseudoStatName.includes('Crit')) {
			return Stat.StatCritRating;
		} else {
			return null;
		}
	}

	// Returns the other school variant of a school-specific PseudoStat, or
	// null if not applicable.
	static getSiblingPseudoStat(pseudoStat: PseudoStat): PseudoStat | null {
		switch (pseudoStat) {
			case PseudoStat.PseudoStatPhysicalHitPercent:
				return PseudoStat.PseudoStatSpellHitPercent;
			case PseudoStat.PseudoStatSpellHitPercent:
				return PseudoStat.PseudoStatPhysicalHitPercent;
			case PseudoStat.PseudoStatPhysicalCritPercent:
				return PseudoStat.PseudoStatSpellCritPercent;
			case PseudoStat.PseudoStatSpellCritPercent:
				return PseudoStat.PseudoStatPhysicalHitPercent;
			default:
				return null;
		}
	}

	static createDisplayStatArray(statList: Stat[], pseudoStatList: PseudoStat[]): UnitStat[] {
		const displayStats: UnitStat[] = [];

		statList.forEach(stat => {
			displayStats.push(UnitStat.fromStat(stat));
		});

		pseudoStatList.forEach(pseudoStat => {
			displayStats.push(UnitStat.fromPseudoStat(pseudoStat));
		});

		return displayStats;
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

	computeGapToCap(unitStat: UnitStat, cap: number): number {
		let statDelta = cap - this.getUnitStat(unitStat);

		if (unitStat.equalsPseudoStat(PseudoStat.PseudoStatMeleeHastePercent)) {
			statDelta /= this.getPseudoStat(PseudoStat.PseudoStatMeleeSpeedMultiplier);
		} else if (unitStat.equalsPseudoStat(PseudoStat.PseudoStatRangedHastePercent)) {
			statDelta /= this.getPseudoStat(PseudoStat.PseudoStatRangedSpeedMultiplier);
		} else if (unitStat.equalsPseudoStat(PseudoStat.PseudoStatSpellHastePercent)) {
			statDelta /= this.getPseudoStat(PseudoStat.PseudoStatCastSpeedMultiplier);
		}

		return statDelta;
	}

	computeStatCapsDelta(statCaps: Stats): Stats {
		return new Stats(
			statCaps.stats.map((value, key) => {
				return (value > 0) ? this.computeGapToCap(UnitStat.fromStat(key), value) : 0;
			}),
			statCaps.pseudoStats.map((value, key) => {
				return (value > 0) ? this.computeGapToCap(UnitStat.fromPseudoStat(key), value) : 0;
			}),
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
	static migrateStatsArray(oldStats: number[], oldApiVersion: number, fallbackStats?: number[]): number[] {
		const conversionMap: ProtoConversionMap<number[]> = new Map([
			[
				1,
				(oldArray: number[]) => {
					// Revision 1 simply re-orders the stats for clarity
					const newIndices = [0, 1, 2, 3, 4, 17, 18, 6, 8, 10, 19, 15, 5, 7, 9, 11, 29, 26, 16, 30, 12, 13, 20, 28, 21, 22, 23, 24, 25, 27, 14];
					const newArray: number[] = new Array(oldArray.length);
					oldArray.forEach((value, idx) => {
						newArray[newIndices[idx]] = value;
					});
					return newArray;
				},
			],
		]);
		const migratedProto = migrateOldProto<number[]>(oldStats, oldApiVersion, conversionMap);

		// If there is a fallback array, use it if the lengths don't match
		if (fallbackStats && migratedProto.length !== fallbackStats.length) return fallbackStats;

		return migratedProto;
	}
}

export interface BreakpointConfig {
	breakpoints: number[];
	capType: StatCapType;
	postCapEPs: number[];
}

// Represents a StatCapConfig proto message as a proper class for UI convenience
export class StatCap {
	readonly unitStat: UnitStat;
	readonly capType: StatCapType;
	breakpoints: number[] = [];
	postCapEPs: number[] = [];

	private constructor(unitStat: UnitStat, breakpoints: number[], capType: StatCapType, postCapEPs: number[]) {
		// Check for valid inputs
		if (capType == StatCapType.TypeSoftCap) {
			if (breakpoints.length != postCapEPs.length) {
				throw new Error('Breakpoint and EP counts do not match!');
			}
		} else if (capType == StatCapType.TypeThreshold) {
			if (postCapEPs.length != 1) {
				throw new Error('Exactly 1 post-cap EP value must be specified for Threshold cap types!');
			}
		} else {
			throw new Error('Only SoftCap and Threshold cap types are supported currently!');
		}

		this.unitStat = unitStat;
		this.capType = capType;
		this.breakpoints = breakpoints;
		this.postCapEPs = postCapEPs;
	}

	static fromStat(stat: Stat, config: BreakpointConfig): StatCap {
		return new StatCap(UnitStat.fromStat(stat), config.breakpoints, config.capType, config.postCapEPs);
	}

	static fromPseudoStat(pseudoStat: PseudoStat, config: BreakpointConfig): StatCap {
		return new StatCap(UnitStat.fromPseudoStat(pseudoStat), config.breakpoints, config.capType, config.postCapEPs);
	}

	static fromProto(message: StatCapConfig[]): StatCap[] {
		const statCapObjects: StatCap[] = [];

		message.filter(config => config.unitStat).forEach(config => {
			statCapObjects.push(new StatCap(UnitStat.fromProto(config.unitStat!), config.breakpoints, config.capType, config.postCapEPs));
		});

		return statCapObjects;
	}
}

export function convertHastePresetBreakpointsToPercent(ratingPresets: Map<string, number>): Map<string, number> {
	const convertedPresets = new Map<string, number>();

	for (const [presetName, ratingValue] of ratingPresets.entries()) {
		convertedPresets.set(presetName, ratingValue / Mechanics.HASTE_RATING_PER_HASTE_PERCENT);
	}

	return convertedPresets;
}
