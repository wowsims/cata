import { StatMods } from '../core/components/character_stats';
import * as Mechanics from '../core/constants/mechanics';
import { Player } from '../core/player';
import { Spec, Stat } from '../core/proto/common';
import { Stats } from '../core/proto_utils/stats';

export const sharedMageDisplayStatsModifiers = (player: Player<Spec.SpecArcaneMage> | Player<Spec.SpecFireMage> | Player<Spec.SpecFrostMage>): StatMods => {
	let stats = new Stats();
	stats = stats.addStat(Stat.StatSpellCrit, player.getTalents().piercingIce * 1 * Mechanics.SPELL_CRIT_RATING_PER_CRIT_CHANCE);
	return {
		talents: stats,
	};
};
