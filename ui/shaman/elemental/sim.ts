import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../../core/components/inputs/other_inputs.js';
import { ReforgeOptimizer } from '../../core/components/suggest_reforges_action';
import * as Mechanics from '../../core/constants/mechanics.js';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui.js';
import { Player } from '../../core/player.js';
import { PlayerClasses } from '../../core/player_classes';
import { APLRotation } from '../../core/proto/apl.js';
import { Debuffs, Faction, IndividualBuffs, ItemSlot, PartyBuffs, PseudoStat, Race, RaidBuffs, Spec, Stat } from '../../core/proto/common.js';
import { Stats, UnitStat } from '../../core/proto_utils/stats.js';
import { TypedEvent } from '../../core/typed_event.js';
import * as ShamanInputs from '../inputs.js';
import * as ElementalInputs from './inputs.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecElementalShaman, {
	cssClass: 'elemental-shaman-sim-ui',
	cssScheme: PlayerClasses.getCssClass(PlayerClasses.Shaman),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],
	warnings: [
		// Warning to use all 4 totems if T6 2pc bonus is active.
		(simUI: IndividualSimUI<Spec.SpecElementalShaman>) => {
			return {
				updateOn: TypedEvent.onAny([simUI.player.rotationChangeEmitter, simUI.player.currentStatsEmitter]),
				getContent: () => {
					const hasT62P = simUI.player.getCurrentStats().sets.includes('Skyshatter Regalia (2pc)');
					const totems = simUI.player.getSpecOptions().classOptions?.totems;
					const hasAll4Totems = totems && totems.earth && totems.air && totems.fire && totems.water;
					if (hasT62P && !hasAll4Totems) {
						return 'T6 2pc bonus is equipped, but inactive because not all 4 totem types are being used.';
					} else {
						return '';
					}
				},
			};
		},
	],

	// All stats for which EP should be calculated.
	epStats: [Stat.StatIntellect, Stat.StatSpirit, Stat.StatSpellPower, Stat.StatHitRating, Stat.StatCritRating, Stat.StatHasteRating, Stat.StatMasteryRating],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatSpellPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: UnitStat.createDisplayStatArray(
		[Stat.StatHealth, Stat.StatMana, Stat.StatStamina, Stat.StatIntellect, Stat.StatSpirit, Stat.StatSpellPower, Stat.StatMasteryRating],
		[PseudoStat.PseudoStatSpellHitPercent, PseudoStat.PseudoStatSpellCritPercent, PseudoStat.PseudoStatSpellHastePercent],
	),
	modifyDisplayStats: (player: Player<Spec.SpecElementalShaman>) => {
		const playerStats = player.getCurrentStats();
		const gearStats = Stats.fromProto(playerStats.gearStats);
		const talentsStats = Stats.fromProto(playerStats.talentsStats);
		const talentsDelta = talentsStats.subtract(gearStats);

		return {
			talents: new Stats().withStat(
				Stat.StatHitRating,
				talentsDelta.getPseudoStat(PseudoStat.PseudoStatSpellHitPercent) * Mechanics.SPELL_HIT_RATING_PER_HIT_PERCENT,
			),
		};
	},

	defaults: {
		// Default equipped gear.
		gear: Presets.P4_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Presets.EP_PRESET_DEFAULT.epWeights,
		// Default stat caps for the Reforge optimizer
		statCaps: (() => {
			return new Stats().withPseudoStat(PseudoStat.PseudoStatSpellHitPercent, 17);
		})(),
		// Default consumes settings.
		consumables: Presets.DefaultConsumables,
		// Default talents.
		talents: Presets.TalentsTotemDuration.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		other: Presets.OtherDefaults,
		// Default raid/party buffs settings.
		raidBuffs: RaidBuffs.create({}),
		partyBuffs: PartyBuffs.create({}),
		individualBuffs: IndividualBuffs.create({}),
		debuffs: Debuffs.create({
			// curseOfElements: true,
			// shadowAndFlame: true,
		}),
	},
	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [ShamanInputs.ShamanShieldInput()],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [ElementalInputs.InThunderstormRange, OtherInputs.InputDelay, OtherInputs.TankAssignment, OtherInputs.DistanceFromTarget],
	},
	itemSwapSlots: [
		ItemSlot.ItemSlotHead,
		ItemSlot.ItemSlotNeck,
		ItemSlot.ItemSlotShoulder,
		ItemSlot.ItemSlotBack,
		ItemSlot.ItemSlotChest,
		ItemSlot.ItemSlotHands,
		ItemSlot.ItemSlotLegs,
		ItemSlot.ItemSlotFinger1,
		ItemSlot.ItemSlotFinger2,
		ItemSlot.ItemSlotTrinket1,
		ItemSlot.ItemSlotTrinket2,
		ItemSlot.ItemSlotMainHand,
	],
	customSections: [ShamanInputs.TotemsSection],
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		epWeights: [Presets.EP_PRESET_DEFAULT],
		// Preset talents that the user can quickly select.
		talents: [Presets.TalentsTotemDuration, Presets.TalentsImprovedShields],
		// Preset rotations that the user can quickly select.
		rotations: [Presets.ROTATION_PRESET_DEFAULT, Presets.ROTATION_PRESET_AOE],
		// Preset gear configurations that the user can quickly select.
		gear: [Presets.PRERAID_PRESET, Presets.P1_PRESET, Presets.P3_PRESET, Presets.P4_PRESET],
		itemSwaps: [Presets.P3_ITEM_SWAP, Presets.P4_ITEM_SWAP],

		builds: [Presets.P3_PRESET_BUILD_DEFAULT, Presets.P3_PRESET_BUILD_CLEAVE],
	},

	autoRotation: (_player: Player<Spec.SpecElementalShaman>): APLRotation => {
		return Presets.ROTATION_PRESET_DEFAULT.rotation.rotation!;
	},

	raidSimPresets: [
		{
			spec: Spec.SpecElementalShaman,
			talents: Presets.TalentsTotemDuration.data,
			specOptions: Presets.DefaultOptions,
			consumables: Presets.DefaultConsumables,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceDraenei,
				[Faction.Horde]: Race.RaceOrc,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.P4_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.P4_PRESET.gear,
				},
			},
			otherDefaults: Presets.OtherDefaults,
		},
	],
});

export class ElementalShamanSimUI extends IndividualSimUI<Spec.SpecElementalShaman> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecElementalShaman>) {
		super(parentElem, player, SPEC_CONFIG);
		player.sim.waitForInit().then(() => {
			new ReforgeOptimizer(this, {
				getEPDefaults: (player: Player<Spec.SpecFuryWarrior>) => {
					const playerWeights = player.getEpWeights();
					const defaultWeights = Presets.EP_PRESET_DEFAULT.epWeights;

					if (playerWeights.equals(defaultWeights)) return defaultWeights;

					return playerWeights;
				},
			});
		});
	}
}
