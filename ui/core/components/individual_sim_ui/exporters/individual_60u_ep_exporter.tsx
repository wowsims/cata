import { IndividualSimUI } from '../../../individual_sim_ui';
import { PseudoStat, Spec, Stat } from '../../../proto/common';
import { UnitStat } from '../../../proto_utils/stats';
import { IndividualExporter } from './individual_exporter';

export class Individual60UEPExporter<SpecType extends Spec> extends IndividualExporter<SpecType> {
	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, simUI, { title: 'Sixty Upgrades Cataclysm EP Export', allowDownload: true });
	}

	getData(): string {
		const player = this.simUI.player;
		const epValues = player.getEpWeights();
		const allUnitStats = UnitStat.getAll();

		const namesToWeights: Record<string, number> = {};
		allUnitStats.forEach(stat => {
			const statName = Individual60UEPExporter.getName(stat);
			const weight = epValues.getUnitStat(stat);
			if (weight == 0 || statName == '') {
				return;
			}

			// Need to add together stats with the same name (e.g. hit/crit/haste).
			if (namesToWeights[statName]) {
				namesToWeights[statName] += weight;
			} else {
				namesToWeights[statName] = weight;
			}
		});

		return (
			`https://sixtyupgrades.com/mop/ep/import?name=${encodeURIComponent(`${player.getPlayerSpec().friendlyName} WoWSims Weights`)}` +
			Object.keys(namesToWeights)
				.map(statName => `&${statName}=${namesToWeights[statName].toFixed(3)}`)
				.join('')
		);
	}

	static getName(stat: UnitStat): string {
		if (stat.isStat()) {
			return Individual60UEPExporter.statNames[stat.getStat()];
		} else {
			return Individual60UEPExporter.pseudoStatNames[stat.getPseudoStat()] || '';
		}
	}

	static statNames: Record<Stat, string> = {
		[Stat.StatStrength]: 'strength',
		[Stat.StatAgility]: 'agility',
		[Stat.StatStamina]: 'stamina',
		[Stat.StatIntellect]: 'intellect',
		[Stat.StatSpirit]: 'spirit',
		[Stat.StatSpellPower]: 'spellDamage',
		[Stat.StatMP5]: 'mp5',
		[Stat.StatHitRating]: 'hitRating',
		[Stat.StatCritRating]: 'critRating',
		[Stat.StatHasteRating]: 'hasteRating',
		[Stat.StatAttackPower]: 'attackPower',
		[Stat.StatMasteryRating]: 'masteryRating',
		[Stat.StatExpertiseRating]: 'expertiseRating',
		// TODO: Change PVP Resilience and Power once 60U exists for MoP
		[Stat.StatPvpResilienceRating]: 'pvpResilienceRating',
		[Stat.StatPvpPowerRating]: 'pvpPowerRating',
		[Stat.StatMana]: 'mana',
		[Stat.StatArmor]: 'armor',
		[Stat.StatRangedAttackPower]: 'attackPower',
		[Stat.StatDodgeRating]: 'dodgeRating',
		[Stat.StatParryRating]: 'parryRating',
		[Stat.StatHealth]: 'health',
		[Stat.StatBonusArmor]: 'armorBonus',
	};
	static pseudoStatNames: Partial<Record<PseudoStat, string>> = {
		[PseudoStat.PseudoStatMainHandDps]: 'dps',
		[PseudoStat.PseudoStatRangedDps]: 'rangedDps',
	};
}
