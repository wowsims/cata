import { StatMods } from '../core/components/character_stats';
import * as Mechanics from '../core/constants/mechanics';
import { Player } from '../core/player';
import { ItemSlot, PseudoStat, Race, RangedWeaponType, Spec, Stat } from '../core/proto/common';
import { Stats } from '../core/proto_utils/stats';

export const sharedHunterDisplayStatsModifiers = (
	player: Player<Spec.SpecBeastMasteryHunter> | Player<Spec.SpecMarksmanshipHunter> | Player<Spec.SpecSurvivalHunter>,
): StatMods => {
	return {};
};
