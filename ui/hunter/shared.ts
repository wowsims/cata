import { StatMods } from '../core/components/character_stats';
import * as Mechanics from '../core/constants/mechanics';
import { Player } from '../core/player';
import { ItemSlot, Race, RangedWeaponType, Spec, Stat } from '../core/proto/common';
import { Stats } from '../core/proto_utils/stats';

export const sharedHunterDisplayStatsModifiers = (
	player: Player<Spec.SpecBeastMasteryHunter> | Player<Spec.SpecMarksmanshipHunter> | Player<Spec.SpecSurvivalHunter>,
): StatMods => {
	let stats = new Stats();

	const rangedWeapon = player.getEquippedItem(ItemSlot.ItemSlotRanged);
	if (rangedWeapon?.enchant?.effectId == 3608) {
		stats = stats.addStat(Stat.StatMeleeCrit, 40);
	}
	if (rangedWeapon?.enchant?.effectId == 4176) {
		stats = stats.addStat(Stat.StatMeleeHit, 88);
	}
	if (player.getRace() == Race.RaceDwarf && rangedWeapon?.item.rangedWeaponType == RangedWeaponType.RangedWeaponTypeGun) {
		stats = stats.addStat(Stat.StatMeleeCrit, 1 * Mechanics.MELEE_CRIT_RATING_PER_CRIT_CHANCE);
	}
	if (player.getRace() == Race.RaceTroll && rangedWeapon?.item.rangedWeaponType == RangedWeaponType.RangedWeaponTypeBow) {
		stats = stats.addStat(Stat.StatMeleeCrit, 1 * Mechanics.MELEE_CRIT_RATING_PER_CRIT_CHANCE);
	}
	return {
		talents: stats,
	};
};
