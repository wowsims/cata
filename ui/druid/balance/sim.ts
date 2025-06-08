import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../../core/components/inputs/other_inputs';
import { ReforgeOptimizer } from '../../core/components/suggest_reforges_action';
import * as Mechanics from '../../core/constants/mechanics';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui';
import { Player } from '../../core/player';
import { PlayerClasses } from '../../core/player_classes';
import { APLRotation, APLRotation_Type } from '../../core/proto/apl';
import { Faction, ItemSlot, PseudoStat, Race, Spec, Stat } from '../../core/proto/common';
import { Stats, UnitStat } from '../../core/proto_utils/stats';
import * as DruidInputs from '../inputs';
import * as BalanceInputs from './inputs';
import * as Presets from './presets';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecBalanceDruid, {
	cssClass: 'balance-druid-sim-ui',
	cssScheme: PlayerClasses.getCssClass(PlayerClasses.Druid),
	// List any known bugs / issues here, and they'll be shown on the site.
	knownIssues: [],

	// All stats for which EP should be calculated.
	epStats: [Stat.StatIntellect, Stat.StatSpirit, Stat.StatSpellPower, Stat.StatHitRating, Stat.StatCritRating, Stat.StatHasteRating, Stat.StatMasteryRating],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatSpellPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: UnitStat.createDisplayStatArray(
		[Stat.StatHealth, Stat.StatMana, Stat.StatStamina, Stat.StatIntellect, Stat.StatSpirit, Stat.StatSpellPower, Stat.StatMasteryRating],
		[PseudoStat.PseudoStatSpellHitPercent, PseudoStat.PseudoStatSpellCritPercent, PseudoStat.PseudoStatSpellHastePercent],
	),

	modifyDisplayStats: (player: Player<Spec.SpecBalanceDruid>) => {
		const playerStats = player.getCurrentStats();
		const gearStats = Stats.fromProto(playerStats.gearStats);
		const talentsStats = Stats.fromProto(playerStats.talentsStats);
		const talentsDelta = talentsStats.subtract(gearStats);
		const talentsMod = new Stats().withStat(
			Stat.StatHitRating,
			talentsDelta.getPseudoStat(PseudoStat.PseudoStatSpellHitPercent) * Mechanics.SPELL_HIT_RATING_PER_HIT_PERCENT,
		);

		return {
			talents: talentsMod,
		};
	},

	defaults: {
		// Default equipped gear.
		gear: Presets.T13PresetGear.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Presets.StandardEPWeights.epWeights,
		// Default stat caps for the Reforge optimizer
		statCaps: (() => {
			return new Stats().withPseudoStat(PseudoStat.PseudoStatSpellHitPercent, 17);
		})(),
		// Default consumes settings.

		consumables: Presets.DefaultConsumables,
		// Default talents.
		talents: Presets.StandardTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		// Default raid/party buffs settings.
		raidBuffs: Presets.DefaultRaidBuffs,
		partyBuffs: Presets.DefaultPartyBuffs,
		individualBuffs: Presets.DefaultIndividualBuffs,
		debuffs: Presets.DefaultDebuffs,
		other: Presets.OtherDefaults,
		rotationType: APLRotation_Type.TypeAuto,
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [DruidInputs.SelfInnervate()],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
		BuffDebuffInputs.ReplenishmentBuff,
		BuffDebuffInputs.CritBuff,
		BuffDebuffInputs.MP5Buff,
		BuffDebuffInputs.AttackPowerPercentBuff,
		BuffDebuffInputs.StaminaBuff,
		BuffDebuffInputs.ManaBuff,
	],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [BalanceInputs.OkfUptime, OtherInputs.TankAssignment, OtherInputs.InputDelay, OtherInputs.DistanceFromTarget, OtherInputs.DarkIntentUptime],
	},
	itemSwapSlots: [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotHands, ItemSlot.ItemSlotTrinket1, ItemSlot.ItemSlotTrinket2],
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		epWeights: [Presets.StandardEPWeights],
		// Preset talents that the user can quickly select.
		talents: [Presets.StandardTalents],
		rotations: [Presets.T11PresetRotation, Presets.T12PresetRotation, Presets.T13PresetRotation],
		// Preset gear configurations that the user can quickly select.
		gear: [Presets.PreraidPresetGear, Presets.T11PresetGear, Presets.T12PresetGear, Presets.T13PresetGear],
		itemSwaps: [Presets.T13PresetItemSwapGear],
		builds: [Presets.PresetBuildPreraid, Presets.PresetBuildT11, Presets.PresetBuildT12, Presets.PresetBuildT13],
	},

	autoRotation: (_player: Player<Spec.SpecBalanceDruid>): APLRotation => {
		return Presets.T13PresetRotation.rotation.rotation!;
	},

	raidSimPresets: [
		{
			spec: Spec.SpecBalanceDruid,
			talents: Presets.StandardTalents.data,
			specOptions: Presets.DefaultOptions,

			consumables: Presets.DefaultConsumables,
			otherDefaults: Presets.OtherDefaults,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceWorgen,
				[Faction.Horde]: Race.RaceTroll,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.T13PresetGear.gear,
				},
				[Faction.Horde]: {
					1: Presets.T13PresetGear.gear,
				},
			},
		},
	],
});

export class BalanceDruidSimUI extends IndividualSimUI<Spec.SpecBalanceDruid> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecBalanceDruid>) {
		super(parentElem, player, SPEC_CONFIG);
		player.sim.waitForInit().then(() => {
			new ReforgeOptimizer(this);
		});
	}
}
