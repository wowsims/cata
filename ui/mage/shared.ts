import { StatMods } from '../core/components/character_stats';
import { Player } from '../core/player';
import { PseudoStat, Spec } from '../core/proto/common';
import { Stats } from '../core/proto_utils/stats';

export const sharedMageDisplayStatsModifiers = (player: Player<Spec.SpecArcaneMage> | Player<Spec.SpecFireMage> | Player<Spec.SpecFrostMage>): StatMods => {
	const stats = new Stats().withPseudoStat(PseudoStat.PseudoStatSpellCritPercent, player.getTalents().piercingIce * 1);
	return {
		talents: stats,
	};
};
