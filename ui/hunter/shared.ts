import { StatMods } from '../core/components/character_stats';
import * as Mechanics from '../core/constants/mechanics';
import { Player } from '../core/player';
import { ItemSlot, PseudoStat, Race, RangedWeaponType, Spec, Stat } from '../core/proto/common';
import { Stats } from '../core/proto_utils/stats';

export const sharedHunterDisplayStatsModifiers = (
	player: Player<Spec.SpecBeastMasteryHunter> | Player<Spec.SpecMarksmanshipHunter> | Player<Spec.SpecSurvivalHunter>,
): StatMods => {
	let stats = new Stats();

	// TODO: Update for MOP Scopes
	const rangedWeapon = player.getEquippedItem(ItemSlot.ItemSlotMainHand);
	if (rangedWeapon?.enchant?.effectId == 3608) {
		stats = stats.addStat(Stat.StatCritRating, 40);
	}
	if (rangedWeapon?.enchant?.effectId == 4176) {
		stats = stats.addStat(Stat.StatHitRating, 88);
	}
	if (player.getRace() == Race.RaceDwarf) {
		stats = stats.withStat(Stat.StatExpertiseRating, Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION * 4);
	}
	if (player.getRace() == Race.RaceTroll) {
		stats = stats.withStat(Stat.StatExpertiseRating, Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION * 4);
	}
	return {
		talents: stats,
	};
};
