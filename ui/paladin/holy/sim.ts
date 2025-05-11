import * as OtherInputs from '../../core/components/inputs/other_inputs.js';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui.js';
import { Player } from '../../core/player.js';
import { PlayerClasses } from '../../core/player_classes';
import { APLRotation } from '../../core/proto/apl.js';
import { Faction, IndividualBuffs, PartyBuffs, PseudoStat, Race, Spec, Stat } from '../../core/proto/common.js';
import { Stats, UnitStat } from '../../core/proto_utils/stats.js';
import * as HolyInputs from '../../paladin/holy/inputs.js';
import * as PaladinInputs from '../inputs.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecHolyPaladin, {
	cssClass: 'holy-paladin-sim-ui',
	cssScheme: PlayerClasses.getCssClass(PlayerClasses.Paladin),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],

	// All stats for which EP should be calculated.
	epStats: [Stat.StatIntellect, Stat.StatSpirit, Stat.StatSpellPower, Stat.StatHasteRating, Stat.StatCritRating, Stat.StatMasteryRating],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatSpellPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: UnitStat.createDisplayStatArray(
		[Stat.StatHealth, Stat.StatMana, Stat.StatIntellect, Stat.StatSpirit, Stat.StatSpellPower, Stat.StatMasteryRating, Stat.StatArmor, Stat.StatStamina],
		[PseudoStat.PseudoStatSpellHastePercent, PseudoStat.PseudoStatSpellCritPercent, PseudoStat.PseudoStatSpellHitPercent],
	),
	defaults: {
		// Default equipped gear.
		gear: Presets.P1_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Presets.P1_EP_PRESET.epWeights,
		// Default consumes settings.
		consumables: Presets.DefaultConsumables,
		// Default talents.
		talents: Presets.StandardTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		other: Presets.OtherDefaults,
		// Default raid/party buffs settings.
		raidBuffs: Presets.DefaultRaidBuffs,
		partyBuffs: PartyBuffs.create({}),

		individualBuffs: IndividualBuffs.create({}),
		debuffs: Presets.DefaultDebuffs,
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [PaladinInputs.AuraSelection()],
	// Inputs to include in the 'Rotation' section on the settings tab.
	rotationInputs: HolyInputs.PaladinRotationConfig,
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [OtherInputs.InputDelay, OtherInputs.TankAssignment],
	},
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		epWeights: [Presets.P1_EP_PRESET],
		// Preset talents that the user can quickly select.
		talents: [Presets.StandardTalents],
		rotations: [],
		// Preset gear configurations that the user can quickly select.
		gear: [Presets.PRERAID_PRESET, Presets.P1_PRESET],
	},

	autoRotation: (_player: Player<Spec.SpecHolyPaladin>): APLRotation => {
		return APLRotation.create();
	},

	raidSimPresets: [
		{
			spec: Spec.SpecHolyPaladin,
			talents: Presets.StandardTalents.data,
			specOptions: Presets.DefaultOptions,
			consumables: Presets.DefaultConsumables,
			otherDefaults: Presets.OtherDefaults,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceHuman,
				[Faction.Horde]: Race.RaceBloodElf,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.P1_PRESET.gear,
					2: Presets.PRERAID_PRESET.gear,
					// 3: Presets.P3_PRESET.gear,
					// 4: Presets.P4_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.P1_PRESET.gear,
					2: Presets.PRERAID_PRESET.gear,
					// 3: Presets.P3_PRESET.gear,
					// 4: Presets.P4_PRESET.gear,
				},
			},
		},
	],
});

export class HolyPaladinSimUI extends IndividualSimUI<Spec.SpecHolyPaladin> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecHolyPaladin>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
