import { Debuffs,  PseudoStat, RaidBuffs } from '../core/proto/common';
import { UnitStat, UnitStatPresets } from '../core/proto_utils/stats';

export const LIVING_BOMB_BREAKPOINTS: UnitStatPresets = {
	unitStat: UnitStat.fromPseudoStat(PseudoStat.PseudoStatSpellHastePercent),
	presets: new Map([
		['5-tick - Living Bomb', 12.507036],
		['6-tick - Living Bomb', 37.520061],
		['7-tick - Living Bomb', 62.469546],
		['8-tick - Living Bomb', 87.441436],
		['9-tick - Living Bomb', 112.539866],
		['10-tick - Living Bomb', 137.435713],
		['11-tick - Living Bomb', 162.58208],
		['12-tick - Living Bomb', 187.494038],
	]),
};

export const NETHER_TEMPEST_BREAKPOINTS: UnitStatPresets = {
	unitStat: UnitStat.fromPseudoStat(PseudoStat.PseudoStatSpellHastePercent),
	presets: new Map([
		['13-tick - Nether Tempest', 4.220959],
		['14-tick - Nether Tempest', 12.549253],
		['15-tick - Nether Tempest', 20.845936],
		['16-tick - Nether Tempest', 29.115575],
		['17-tick - Nether Tempest', 37.457064],
		['18-tick - Nether Tempest', 45.878942],
		['19-tick - Nether Tempest', 54.202028],
		['20-tick - Nether Tempest', 62.469563],
		['21-tick - Nether Tempest', 70.794222],
		['22-tick - Nether Tempest', 79.051062],
		['23-tick - Nether Tempest', 87.44146],
		['24-tick - Nether Tempest', 95.886424],
		['25-tick - Nether Tempest', 104.290134],
		['26-tick - Nether Tempest', 112.539896],
		['27-tick - Nether Tempest', 120.994524],
		['28-tick - Nether Tempest', 129.095127],
		['29-tick - Nether Tempest', 137.24798],
		['30-tick - Nether Tempest', 146.002521],
	]),
};

export const MAGE_BREAKPOINTS: UnitStatPresets = {
	unitStat: UnitStat.fromPseudoStat(PseudoStat.PseudoStatSpellHastePercent),
	presets: new Map([...LIVING_BOMB_BREAKPOINTS.presets, ...NETHER_TEMPEST_BREAKPOINTS.presets].sort((a, b) => a[1] - b[1])),
};

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	blessingOfKings: true,
	mindQuickening: true,
	leaderOfThePack: true,
	blessingOfMight: true,
	unholyAura: true,
	bloodlust: true,
	skullBannerCount: 2,
	stormlashTotemCount: 4,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfElements: true,
});
