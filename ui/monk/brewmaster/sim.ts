import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../../core/components/inputs/other_inputs';
import { ReforgeOptimizer } from '../../core/components/suggest_reforges_action';
import * as Mechanics from '../../core/constants/mechanics.js';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui';
import { Player } from '../../core/player';
import { PlayerClasses } from '../../core/player_classes';
import { APLRotation } from '../../core/proto/apl';
import { Debuffs, Faction, IndividualBuffs, PartyBuffs, PseudoStat, Race, RaidBuffs, Spec, Stat } from '../../core/proto/common';
import { StatCapType } from '../../core/proto/ui';
import { StatCap, Stats, UnitStat } from '../../core/proto_utils/stats';
import { Sim } from '../../core/sim';
import { TypedEvent } from '../../core/typed_event';
import * as Presets from './presets';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecBrewmasterMonk, {
	cssClass: 'brewmaster-monk-sim-ui',
	cssScheme: PlayerClasses.getCssClass(PlayerClasses.Monk),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],

	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatAgility,
		Stat.StatAttackPower,
		Stat.StatHitRating,
		Stat.StatExpertiseRating,
		Stat.StatCritRating,
		Stat.StatHasteRating,
		Stat.StatDodgeRating,
		Stat.StatParryRating,
		Stat.StatMasteryRating,
	],
	epPseudoStats: [PseudoStat.PseudoStatMainHandDps, PseudoStat.PseudoStatOffHandDps],
	// Reference stat against which to calculate EP.
	epReferenceStat: Stat.StatAttackPower,
	consumableStats: [
		Stat.StatAgility,
		Stat.StatBonusArmor,
		Stat.StatStamina,
		Stat.StatArmor,
		Stat.StatAttackPower,
		Stat.StatDodgeRating,
		Stat.StatParryRating,
		Stat.StatHitRating,
		Stat.StatHasteRating,
		Stat.StatCritRating,
		Stat.StatExpertiseRating,
		Stat.StatMasteryRating,
	],
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: UnitStat.createDisplayStatArray(
		[Stat.StatHealth, Stat.StatStamina, Stat.StatAgility, Stat.StatStrength, Stat.StatAttackPower, Stat.StatMasteryRating, Stat.StatExpertiseRating],
		[
			PseudoStat.PseudoStatPhysicalHitPercent,
			PseudoStat.PseudoStatSpellHitPercent,
			PseudoStat.PseudoStatPhysicalCritPercent,
			PseudoStat.PseudoStatSpellCritPercent,
			PseudoStat.PseudoStatMeleeHastePercent,
			PseudoStat.PseudoStatDodgePercent,
			PseudoStat.PseudoStatParryPercent,
		],
	),

	defaults: {
		// Default equipped gear.
		gear: Presets.P1_BIS_BALANCED_DW_GEAR_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Presets.PREPATCH_EP_PRESET.epWeights,
		// Stat caps for reforge optimizer
		statCaps: (() => {
			const hitCap = new Stats().withPseudoStat(PseudoStat.PseudoStatPhysicalHitPercent, 7.5);
			return hitCap;
		})(),
		softCapBreakpoints: (() => {
			const expertiseCap = StatCap.fromStat(Stat.StatExpertiseRating, {
				breakpoints: [7.5 * 4 * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION, 15 * 4 * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION],
				capType: StatCapType.TypeSoftCap,
				// These are set by the active EP weight in the updateSoftCaps callback
				postCapEPs: [6.04, 0],
			});

			return [expertiseCap];
		})(),
		other: Presets.OtherDefaults,
		// Default consumes settings.
		consumables: Presets.DefaultConsumables,
		// Default talents.
		talents: Presets.DefaultTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		// Default raid/party buffs settings.
		raidBuffs: RaidBuffs.create({
			legacyOfTheEmperor: true,
			legacyOfTheWhiteTiger: true,
			darkIntent: true,
			trueshotAura: true,
			unleashedRage: true,
			moonkinAura: true,
			blessingOfMight: true,
			bloodlust: true,
			skullBannerCount: 2,
			stormlashTotemCount: 4,
		}),
		partyBuffs: PartyBuffs.create({}),
		individualBuffs: IndividualBuffs.create({}),
		debuffs: Debuffs.create({
			curseOfElements: true,
			physicalVulnerability: true,
			weakenedArmor: true,
		}),
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [BuffDebuffInputs.CritBuff, BuffDebuffInputs.MajorArmorDebuff],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			OtherInputs.InputDelay,
			OtherInputs.TankAssignment,
			OtherInputs.HpPercentForDefensives,
			OtherInputs.IncomingHps,
			OtherInputs.HealingCadence,
			OtherInputs.HealingCadenceVariation,
			OtherInputs.AbsorbFrac,
			OtherInputs.BurstWindow,
			OtherInputs.InFrontOfTarget,
		],
	},
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		epWeights: [Presets.PREPATCH_EP_PRESET],
		// Preset talents that the user can quickly select.
		talents: [Presets.DefaultTalents, Presets.DungeonTalents],
		// Preset rotations that the user can quickly select.
		rotations: [Presets.ROTATION_PRESET],
		// Preset gear configurations that the user can quickly select.
		gear: [
			Presets.P1_PREBIS_RICH_GEAR_PRESET,
			Presets.P1_PREBIS_POOR_GEAR_PRESET,
			Presets.P1_BIS_BALANCED_DW_GEAR_PRESET,
			Presets.P1_BIS_BALANCED_2H_GEAR_PRESET,
			Presets.P1_BIS_OFFENSIVE_DW_GEAR_PRESET,
			Presets.P1_BIS_OFFENSIVE_2H_GEAR_PRESET,
		],
	},

	autoRotation: (_: Player<Spec.SpecBrewmasterMonk>): APLRotation => {
		return Presets.ROTATION_PRESET.rotation.rotation!;
	},

	raidSimPresets: [
		{
			spec: Spec.SpecBrewmasterMonk,
			talents: Presets.DefaultTalents.data,
			specOptions: Presets.DefaultOptions,
			consumables: Presets.DefaultConsumables,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceAlliancePandaren,
				[Faction.Horde]: Race.RaceHordePandaren,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.P1_BIS_BALANCED_DW_GEAR_PRESET.gear,
					2: Presets.P1_BIS_BALANCED_DW_GEAR_PRESET.gear,
					3: Presets.P1_BIS_BALANCED_DW_GEAR_PRESET.gear,
					4: Presets.P1_BIS_BALANCED_DW_GEAR_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.P1_BIS_BALANCED_DW_GEAR_PRESET.gear,
					2: Presets.P1_BIS_BALANCED_DW_GEAR_PRESET.gear,
					3: Presets.P1_BIS_BALANCED_DW_GEAR_PRESET.gear,
					4: Presets.P1_BIS_BALANCED_DW_GEAR_PRESET.gear,
				},
			},
			otherDefaults: Presets.OtherDefaults,
		},
	],
});

const getActiveEPWeight = (player: Player<Spec.SpecBrewmasterMonk>, sim: Sim): Stats => {
	if (sim.getUseCustomEPValues()) {
		return player.getEpWeights();
	} else {
		return Presets.PREPATCH_EP_PRESET.epWeights;
	}
};

export class BrewmasterMonkSimUI extends IndividualSimUI<Spec.SpecBrewmasterMonk> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecBrewmasterMonk>) {
		super(parentElem, player, SPEC_CONFIG);

		const setTalentBasedSettings = () => {
			const talents = player.getTalents();
			// Zen sphere can be on 2 targets, so we set the target dummies to 1 if it is talented.
			player.getRaid()?.setTargetDummies(TypedEvent.nextEventID(), talents.zenSphere ? 2 : 0);
		};

		setTalentBasedSettings();
		player.talentsChangeEmitter.on(() => {
			setTalentBasedSettings();
		});

		player.sim.waitForInit().then(() => {
			new ReforgeOptimizer(this, {
				getEPDefaults: (player: Player<Spec.SpecBrewmasterMonk>) => {
					return getActiveEPWeight(player, this.sim);
				},
			});
		});
	}
}
