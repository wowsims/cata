import * as Mechanics from '../constants/mechanics.js';
import { CURRENT_API_VERSION } from '../constants/other.js';
import { Class, PseudoStat, Stat, UnitStats } from '../proto/common.js';
import { StatCapConfig, StatCapType, UIStat as UnitStatProto } from '../proto/ui.js';
import { getEnumValues } from '../utils.js';
import { getClassPseudoStatName, getStatName } from './names.js';
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
	hasRootStat(): boolean {
		return this.rootStat != null;
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
	getRootStat(): Stat {
		if (!this.hasRootStat()) {
			throw new Error('No root stat for this PseudoStat');
		}
		return this.rootStat!;
	}

	equals(other: UnitStat): boolean {
		return this.stat == other.stat && this.pseudoStat == other.pseudoStat;
	}

	equalsStat(other: Stat): boolean {
		return this.isStat() && this.stat == other;
	}

	equalsPseudoStat(other: PseudoStat): boolean {
		return this.isPseudoStat() && this.pseudoStat == other;
	}

	linkedToStat(other: Stat): boolean {
		return this.stat == other || this.rootStat == other;
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

	getKey(): string {
		return this.isStat() ? Stat[this.stat!] : PseudoStat[this.getPseudoStat()];
	}

	// Convert a UnitStat value from its Rating representation to a percentage representation
	// (0-100). If a percentage representation does not make sense for the stat in question
	// (Strength for example), then null is returned. Mastery is special cased to return
	// Mastery points rather than %.
	convertRatingToPercent(ratingValue: number): number | null {
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
	convertPercentToRating(percentOrPointsValue: number): number | null {
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

	convertDefaultUnitsToRating(value: number): number | null {
		if (this.isStat()) {
			// Proper Stats are either already in Rating units, or
			// do not have a Rating representation.
			if (Stat[this.stat!].includes('Rating')) {
				return value;
			} else {
				return null;
			}
		} else {
			// PseudoStats are either in percent units, or do not
			// have a Rating representation.
			return this.convertPercentToRating(value);
		}
	}

	convertDefaultUnitsToPercent(value: number): number | null {
		if (this.isStat()) {
			// Proper stats are either in Rating units, or do not
			// have a percent representation.
			return this.convertRatingToPercent(value);
		} else {
			// PseudoStats are either already in percent units, or do
			// not have a percent representation.
			if (PseudoStat[this.getPseudoStat()].includes('Percent')) {
				return value;
			} else {
				return null;
			}
		}
	}

	convertRatingToDefaultUnits(ratingValue: number): number | null {
		if (this.isStat()) {
			return Stat[this.stat!].includes('Rating') ? ratingValue : null;
		} else {
			return this.convertRatingToPercent(ratingValue);
		}
	}

	convertPercentToDefaultUnits(percentValue: number): number | null {
		if (this.isStat()) {
			return this.convertPercentToRating(percentValue);
		} else {
			return PseudoStat[this.getPseudoStat()].includes('Percent') ? percentValue : null;
		}
	}

	convertEpToRatingScale(epValue: number): number {
		if (this.isPseudoStat() && PseudoStat[this.pseudoStat!].includes('Percent')) {
			return this.convertRatingToPercent(epValue)!;
		} else {
			return epValue;
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

	toJson(): object {
		return UnitStatProto.toJson(this.toProto()) as object;
	}

	static fromJson(obj: any): UnitStat {
		return UnitStat.fromProto(UnitStatProto.fromJson(obj));
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
		return new UnitStat(stat, null, stat);
	}
	static fromPseudoStat(pseudoStat: PseudoStat): UnitStat {
		return new UnitStat(null, pseudoStat, UnitStat.getRootStat(pseudoStat));
	}

	static getAllStats(): Array<UnitStat> {
		const allStats = getEnumValues(Stat) as Array<Stat>;
		return allStats.map(stat => UnitStat.fromStat(stat));
	}
	static getAllPseudoStats(): Array<UnitStat> {
		const allPseudoStats = getEnumValues(PseudoStat) as Array<PseudoStat>;
		return allPseudoStats.map(pseudoStat => UnitStat.fromPseudoStat(pseudoStat));
	}
	static getAll(): Array<UnitStat> {
		return [UnitStat.getAllStats(), UnitStat.getAllPseudoStats()].flat();
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

	// Inverse of the above
	static getChildren(parentStat: Stat): PseudoStat[] {
		switch (parentStat) {
			case Stat.StatHitRating:
				return [PseudoStat.PseudoStatPhysicalHitPercent, PseudoStat.PseudoStatSpellHitPercent];
			case Stat.StatCritRating:
				return [PseudoStat.PseudoStatPhysicalCritPercent, PseudoStat.PseudoStatSpellCritPercent];
			case Stat.StatHasteRating:
				return [PseudoStat.PseudoStatMeleeHastePercent, PseudoStat.PseudoStatRangedHastePercent, PseudoStat.PseudoStatSpellHastePercent];
			default:
				return [];
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
		return displayStatOrder.filter(
			displayStat =>
				(displayStat.isStat() && statList.includes(displayStat.stat!)) ||
				(displayStat.isPseudoStat() && pseudoStatList.includes(displayStat.pseudoStat!)),
		);
	}
}

export const displayStatOrder: Array<UnitStat> = [
	UnitStat.fromStat(Stat.StatHealth),
	UnitStat.fromStat(Stat.StatMana),
	UnitStat.fromStat(Stat.StatArmor),
	UnitStat.fromStat(Stat.StatBonusArmor),
	UnitStat.fromStat(Stat.StatStamina),
	UnitStat.fromStat(Stat.StatStrength),
	UnitStat.fromStat(Stat.StatAgility),
	UnitStat.fromStat(Stat.StatIntellect),
	UnitStat.fromStat(Stat.StatSpirit),
	UnitStat.fromStat(Stat.StatSpellPower),
	UnitStat.fromPseudoStat(PseudoStat.PseudoStatSpellHitPercent),
	UnitStat.fromPseudoStat(PseudoStat.PseudoStatSpellCritPercent),
	UnitStat.fromPseudoStat(PseudoStat.PseudoStatSpellHastePercent),
	UnitStat.fromStat(Stat.StatMP5),
	UnitStat.fromStat(Stat.StatAttackPower),
	UnitStat.fromStat(Stat.StatRangedAttackPower),
	UnitStat.fromPseudoStat(PseudoStat.PseudoStatPhysicalHitPercent),
	UnitStat.fromPseudoStat(PseudoStat.PseudoStatPhysicalCritPercent),
	UnitStat.fromPseudoStat(PseudoStat.PseudoStatMeleeHastePercent),
	UnitStat.fromPseudoStat(PseudoStat.PseudoStatRangedHastePercent),
	UnitStat.fromStat(Stat.StatExpertiseRating),
	UnitStat.fromStat(Stat.StatMasteryRating),
	UnitStat.fromPseudoStat(PseudoStat.PseudoStatBlockPercent),
	UnitStat.fromPseudoStat(PseudoStat.PseudoStatDodgePercent),
	UnitStat.fromPseudoStat(PseudoStat.PseudoStatParryPercent),
];

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
	addPseudoStat(pseudoStat: PseudoStat, value: number): Stats {
		return this.withPseudoStat(pseudoStat, this.getPseudoStat(pseudoStat) + value);
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
				return value > 0 ? this.computeGapToCap(UnitStat.fromStat(key), value) : 0;
			}),
			statCaps.pseudoStats.map((value, key) => {
				return value > 0 ? this.computeGapToCap(UnitStat.fromPseudoStat(key), value) : 0;
			}),
		);
	}

	asProtoArray(): Array<number> {
		return this.stats.slice();
	}

	asUnitStatArray(): [UnitStat, number][] {
		const statValues = this.stats.map((value, key) => [UnitStat.fromStat(key), value] as [UnitStat, number]);
		const pseudoStatValues = this.pseudoStats.map((value, key) => [UnitStat.fromPseudoStat(key), value] as [UnitStat, number]);
		return statValues.concat(pseudoStatValues);
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
			Stats.updateProtoVersion(unitStats);

			return new Stats(unitStats.stats, unitStats.pseudoStats);
		} else {
			return new Stats();
		}
	}

	static updateProtoVersion(proto: UnitStats) {
		if (!(proto.apiVersion < CURRENT_API_VERSION)) {
			return;
		}
	}

	// Takes in a stats array that was generated from an out-of-date proto version, and converts it to an array that is consistent with the current proto version.
	static migrateStatsArray(oldStats: number[], oldApiVersion: number, fallbackStats?: number[], targetApiVersion?: number): number[] {
		const conversionMap: ProtoConversionMap<number[]> = new Map([
			[
				1,
				(oldArray: number[]) => {
					// Revision 1 simply re-orders the stats for clarity.
					const newIndices = [0, 1, 2, 3, 4, 17, 18, 6, 8, 10, 19, 15, 5, 7, 9, 11, 29, 26, 16, 30, 12, 13, 20, 28, 21, 22, 23, 24, 25, 27, 14];
					const newArray: number[] = new Array(oldArray.length);
					oldArray.forEach((value, idx) => {
						newArray[newIndices[idx]] = value;
					});
					return newArray;
				},
			],
			[
				2,
				(oldArray: number[]) => {
					// Revision 2 collapses school-specific Hit/Crit/Haste into generic Rating stats.
					const newArray: number[] = new Array(oldArray.length - 4).fill(0);
					newArray[5] = oldArray[5] + oldArray[6]; // MeleeHit + SpellHit --> HitRating
					newArray[6] = oldArray[7] + oldArray[8]; // MeleeCrit + SpellCrit --> CritRating
					newArray[7] = oldArray[9] + oldArray[10]; // MeleeHaste + SpellHaste --> HasteRating
					newArray[26] = oldArray[18]; // MP5 was moved

					// Other entries are simply shifted over
					// and copied.
					for (let idx = 0; idx < 5; idx++) {
						newArray[idx] = oldArray[idx];
					}
					for (let idx = 8; idx < 15; idx++) {
						newArray[idx] = oldArray[idx + 3];
					}
					for (let idx = 15; idx < 26; idx++) {
						newArray[idx] = oldArray[idx + 4];
					}

					return newArray;
				},
			],
		]);
		const migratedProto = migrateOldProto<number[]>(oldStats, oldApiVersion, conversionMap, targetApiVersion);

		// If there is a fallback array, use it if the lengths don't match
		if (fallbackStats && migratedProto.length !== fallbackStats.length) return fallbackStats;

		return migratedProto;
	}
}

// Used for spec specific stat presets to be used
// as a easy to access reference
export interface UnitStatPresets {
	unitStat: UnitStat;
	// Name of the preset and the value in percentage
	presets: Map<string, number>;
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

	constructor(unitStat: UnitStat, breakpoints: number[], capType: StatCapType, postCapEPs: number[]) {
		// Check for valid inputs
		if (capType == StatCapType.TypeSoftCap && breakpoints.length != postCapEPs.length) {
			throw new Error('Breakpoint and EP counts do not match!');
		}
		if (capType != StatCapType.TypeSoftCap && capType != StatCapType.TypeThreshold) {
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

		message
			.filter(config => config.unitStat)
			.forEach(config => {
				statCapObjects.push(new StatCap(UnitStat.fromProto(config.unitStat!), config.breakpoints, config.capType, config.postCapEPs));
			});

		return statCapObjects;
	}

	static cloneSoftCaps(softCaps: StatCap[]): StatCap[] {
		const clonedSoftCaps: StatCap[] = [];

		softCaps.forEach(config => {
			clonedSoftCaps.push(new StatCap(config.unitStat, config.breakpoints.slice(), config.capType, config.postCapEPs.slice()));
		});

		return clonedSoftCaps;
	}
}

export function convertHastePresetBreakpointsToPercent(ratingPresets: Map<string, number>): Map<string, number> {
	const convertedPresets = new Map<string, number>();

	for (const [presetName, ratingValue] of ratingPresets.entries()) {
		convertedPresets.set(presetName, ratingValue / Mechanics.HASTE_RATING_PER_HASTE_PERCENT);
	}

	return convertedPresets;
}

// Helper utility to determine whether a particular PseudoStat has been configured as either a hard cap or
// soft cap.
export function pseudoStatIsCapped(pseudoStat: PseudoStat, hardCaps: Stats, softCaps: StatCap[]): boolean {
	if (hardCaps.getPseudoStat(pseudoStat) != 0) {
		return true;
	}

	for (const config of softCaps) {
		if (config.unitStat.equalsPseudoStat(pseudoStat)) {
			return true;
		}
	}

	return false;
}
