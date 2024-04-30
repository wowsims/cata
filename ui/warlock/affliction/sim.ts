import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../../core/components/other_inputs';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui';
import { Player } from '../../core/player';
import { PlayerClasses } from '../../core/player_classes';
import { APLRotation } from '../../core/proto/apl';
import { Faction, ItemSlot, PartyBuffs, Race, Spec, Stat } from '../../core/proto/common';
import { Stats } from '../../core/proto_utils/stats';
import * as WarlockInputs from '../inputs';
import * as Presets from './presets';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecAfflictionWarlock, {
	cssClass: 'affliction-warlock-sim-ui',
	cssScheme: PlayerClasses.getCssClass(PlayerClasses.Warlock),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: ['Drain Soul is currently disabled for APL rotations'],

	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatIntellect,
		Stat.StatSpirit,
		Stat.StatSpellPower,
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatSpellHaste,
		Stat.StatStamina,
		Stat.StatMastery,
	],
	// Reference stat against which to calculate EP. DPS classes use either spell power or attack power.
	epReferenceStat: Stat.StatSpellPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		Stat.StatHealth,
		Stat.StatIntellect,
		Stat.StatSpirit,
		Stat.StatSpellPower,
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatSpellHaste,
		Stat.StatMP5,
		Stat.StatStamina,
		Stat.StatMastery,
	],

	defaults: {
		// Default equipped gear.
		gear: Presets.P3_AFFLICTION_HORDE_PRESET.gear,

		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap({
			[Stat.StatIntellect]: 0.18,
			[Stat.StatSpirit]: 0.54,
			[Stat.StatSpellPower]: 1,
			[Stat.StatSpellHit]: 0.93,
			[Stat.StatSpellCrit]: 0.53,
			[Stat.StatSpellHaste]: 0.81,
			[Stat.StatStamina]: 0.01,
		}),
		// Default consumes settings.
		consumes: Presets.DefaultConsumes,

		// Default talents.
		talents: Presets.AfflictionTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.AfflictionOptions,

		// Default buffs and debuffs settings.
		raidBuffs: Presets.DefaultRaidBuffs,

		partyBuffs: PartyBuffs.create({}),

		individualBuffs: Presets.DefaultIndividualBuffs,

		debuffs: Presets.DefaultDebuffs,

		other: Presets.OtherDefaults,
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [
		// WarlockInputs.PetInput(),
		// WarlockInputs.ArmorInput(),
		// WarlockInputs.WeaponImbueInput()
	],

	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
		BuffDebuffInputs.ReplenishmentBuff,
		BuffDebuffInputs.MajorArmorDebuff,
		BuffDebuffInputs.PhysicalDamageDebuff,
		BuffDebuffInputs.MeleeHasteBuff,
		BuffDebuffInputs.CritBuff,
		BuffDebuffInputs.MP5Buff,
		BuffDebuffInputs.AttackPowerPercentBuff,
		BuffDebuffInputs.StrengthAndAgilityBuff,
		BuffDebuffInputs.StaminaBuff,
	],
	excludeBuffDebuffInputs: [],
	petConsumeInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			WarlockInputs.DetonateSeed(),
			OtherInputs.InputDelay,
			OtherInputs.DistanceFromTarget,
			OtherInputs.TankAssignment,
			OtherInputs.ChannelClipDelay,
		],
	},
	itemSwapSlots: [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand, ItemSlot.ItemSlotRanged],
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		// Preset talents that the user can quickly select.
		talents: [Presets.AfflictionTalents],
		// Preset rotations that the user can quickly select.
		rotations: [Presets.APL_Affliction_Default],

		// Preset gear configurations that the user can quickly select.
		gear: [
			Presets.PRERAID_AFFLICTION_PRESET,
			Presets.P1_AFFLICTION_PRESET,
			Presets.P2_AFFLICTION_PRESET,
			Presets.P3_AFFLICTION_ALLIANCE_PRESET,
			Presets.P3_AFFLICTION_HORDE_PRESET,
			Presets.P4_AFFLICTION_PRESET,
		],
	},

	autoRotation: (_player: Player<Spec.SpecAfflictionWarlock>): APLRotation => {
		return Presets.APL_Affliction_Default.rotation.rotation!;
	},

	raidSimPresets: [
		{
			spec: Spec.SpecAfflictionWarlock,
			talents: Presets.AfflictionTalents.data,
			specOptions: Presets.AfflictionOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceHuman,
				[Faction.Horde]: Race.RaceOrc,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.P1_AFFLICTION_PRESET.gear,
					2: Presets.P2_AFFLICTION_PRESET.gear,
					3: Presets.P3_AFFLICTION_ALLIANCE_PRESET.gear,
					4: Presets.P4_AFFLICTION_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.P1_AFFLICTION_PRESET.gear,
					2: Presets.P2_AFFLICTION_PRESET.gear,
					3: Presets.P3_AFFLICTION_HORDE_PRESET.gear,
					4: Presets.P4_AFFLICTION_PRESET.gear,
				},
			},
			otherDefaults: Presets.OtherDefaults,
		},
	],
});

export class AfflictionWarlockSimUI extends IndividualSimUI<Spec.SpecAfflictionWarlock> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecAfflictionWarlock>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
