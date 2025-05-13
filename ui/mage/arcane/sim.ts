import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../../core/components/inputs/other_inputs';
import { ReforgeOptimizer } from '../../core/components/suggest_reforges_action';
import * as Mechanics from '../../core/constants/mechanics';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui';
import { Player } from '../../core/player';
import { PlayerClasses } from '../../core/player_classes';
import { Mage } from '../../core/player_classes/mage';
import { APLRotation } from '../../core/proto/apl';
import { Faction, IndividualBuffs, ItemSlot, PartyBuffs, PseudoStat, Race, Spec, Stat } from '../../core/proto/common';
import { StatCapType } from '../../core/proto/ui';
import { StatCap, Stats, UnitStat } from '../../core/proto_utils/stats';
import * as ArcaneInputs from './inputs';
import * as Presets from './presets';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecArcaneMage, {
	cssClass: 'arcane-mage-sim-ui',
	cssScheme: PlayerClasses.getCssClass(Mage),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],

	// All stats for which EP should be calculated.
	epStats: [Stat.StatIntellect, Stat.StatSpellPower, Stat.StatHitRating, Stat.StatCritRating, Stat.StatHasteRating, Stat.StatMasteryRating], // Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatSpellPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: UnitStat.createDisplayStatArray(
		[Stat.StatHealth, Stat.StatMana, Stat.StatStamina, Stat.StatIntellect, Stat.StatSpellPower, Stat.StatMasteryRating],
		[PseudoStat.PseudoStatSpellHitPercent, PseudoStat.PseudoStatSpellCritPercent, PseudoStat.PseudoStatSpellHastePercent],
	),

	defaults: {
		// Default equipped gear.
		gear: Presets.ARCANE_P4_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Presets.P1_EP_PRESET.epWeights,
		// Default stat caps for the Reforge Optimizer
		statCaps: (() => {
			return new Stats().withPseudoStat(PseudoStat.PseudoStatSpellHitPercent, 17);
		})(),
		// Default soft caps for the Reforge optimizer
		softCapBreakpoints: (() => {
			// Sources:
			// https://www.icy-veins.com/cataclysm-classic/arcane-mage-pve-stat-priority
			// https://www.wowhead.com/mop-classic/guide/classes/mage/arcane/dps-stat-priority-attributes-pve
			const breakpoints = [23.14];
			const hasteSoftCapConfig = StatCap.fromPseudoStat(PseudoStat.PseudoStatSpellHastePercent, {
				breakpoints,
				capType: StatCapType.TypeSoftCap,
				postCapEPs: [0.48 * Mechanics.HASTE_RATING_PER_HASTE_PERCENT],
			});

			return [hasteSoftCapConfig];
		})(),
		// Default consumes settings.
		consumables: Presets.DefaultConsumables,
		// Default talents.
		talents: Presets.ArcaneTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultArcaneOptions,
		other: Presets.OtherDefaults,
		// Default raid/party buffs settings.
		raidBuffs: Presets.DefaultRaidBuffs,

		partyBuffs: PartyBuffs.create({
			manaTideTotems: 1,
		}),
		individualBuffs: IndividualBuffs.create({}),
		debuffs: Presets.DefaultDebuffs,
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [],
	// Inputs to include in the 'Rotation' section on the settings tab.
	rotationInputs: ArcaneInputs.MageRotationConfig,
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [ArcaneInputs.FocusMagicUptime, OtherInputs.InputDelay, OtherInputs.DistanceFromTarget, OtherInputs.TankAssignment],
	},
	itemSwapSlots: [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotHands, ItemSlot.ItemSlotTrinket1, ItemSlot.ItemSlotTrinket2],
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: true,
	},

	presets: {
		epWeights: [Presets.P1_EP_PRESET],
		// Preset rotations that the user can quickly select.
		rotations: [Presets.ARCANE_ROTATION_PRESET_DEFAULT],
		// Preset talents that the user can quickly select.
		talents: [Presets.ArcaneTalents],
		// Preset gear configurations that the user can quickly select.
		gear: [Presets.ARCANE_P1_PRESET, Presets.ARCANE_P3_PREBIS_PRESET, Presets.ARCANE_P3_PRESET, Presets.ARCANE_P4_PRESET],
	},

	autoRotation: (player: Player<Spec.SpecArcaneMage>): APLRotation => {
		/* 		const numTargets = player.sim.encounter.targets.length;
		if (numTargets > 3) {
			return Presets.ARCANE_ROTATION_PRESET_AOE.rotation.rotation!;
		} else {
			return Presets.ARCANE_ROTATION_PRESET_DEFAULT.rotation.rotation!;
		} */
		return Presets.ARCANE_ROTATION_PRESET_DEFAULT.rotation.rotation!;
	},

	/* simpleRotation: (player: Player<Spec.SpecArcaneMage>, simple: ArcaneMage_Rotation, cooldowns: Cooldowns): APLRotation => {
		const [prepullActions, actions] = AplUtils.standardCooldownDefaults(cooldowns);

		const prepullMirrorImage = APLPrepullAction.fromJsonString(
			`{"action":{"castSpell":{"spellId":{"spellId":55342}}},"doAtValue":{"const":{"val":"-2s"}}}`,
		);

		const berserking = APLAction.fromJsonString(
			`{"condition":{"not":{"val":{"auraIsActive":{"auraId":{"spellId":12472}}}}},"castSpell":{"spellId":{"spellId":26297}}}`,
		);
		const hyperspeedAcceleration = APLAction.fromJsonString(
			`{"condition":{"not":{"val":{"auraIsActive":{"auraId":{"spellId":12472}}}}},"castSpell":{"spellId":{"spellId":54758}}}`,
		);
		const combatPot = APLAction.fromJsonString(
			`{"condition":{"not":{"val":{"auraIsActive":{"auraId":{"spellId":12472}}}}},"castSpell":{"spellId":{"otherId":"OtherActionPotion"}}}`,
		);
		const evocation = APLAction.fromJsonString(
			`{"condition":{"cmp":{"op":"OpLe","lhs":{"currentManaPercent":{}},"rhs":{"const":{"val":"25%"}}}},"castSpell":{"spellId":{"spellId":12051}}}`,
		);

		const arcaneBlastBelowStacks = APLAction.fromJsonString(
			`{"condition":{"or":{"vals":[{"cmp":{"op":"OpLt","lhs":{"auraNumStacks":{"auraId":{"spellId":36032}}},"rhs":{"const":{"val":"4"}}}},{"and":{"vals":[{"cmp":{"op":"OpLt","lhs":{"auraNumStacks":{"auraId":{"spellId":36032}}},"rhs":{"const":{"val":"3"}}}},{"cmp":{"op":"OpLt","lhs":{"currentManaPercent":{}},"rhs":{"const":{"val":"${(
				simple.only3ArcaneBlastStacksBelowManaPercent * 100
			).toFixed(0)}%"}}}}]}}]}},"castSpell":{"spellId":{"spellId":42897}}}`,
		);
		const arcaneMissilesWithMissileBarrageBelowMana = APLAction.fromJsonString(
			`{"condition":{"and":{"vals":[{"auraIsActiveWithReactionTime":{"auraId":{"spellId":44401}}},{"cmp":{"op":"OpLt","lhs":{"currentManaPercent":{}},"rhs":{"const":{"val":"${(
				simple.missileBarrageBelowManaPercent * 100
			).toFixed(0)}%"}}}}]}},"castSpell":{"spellId":{"spellId":42846}}}`,
		);
		const arcaneMisslesWithMissileBarrage = APLAction.fromJsonString(
			`{"condition":{"auraIsActiveWithReactionTime":{"auraId":{"spellId":44401}}},"castSpell":{"spellId":{"spellId":42846}}}`,
		);
		const arcaneBlastAboveMana = APLAction.fromJsonString(
			`{"condition":{"cmp":{"op":"OpGt","lhs":{"currentManaPercent":{}},"rhs":{"const":{"val":"${(
				simple.blastWithoutMissileBarrageAboveManaPercent * 100
			).toFixed(0)}%"}}}},"castSpell":{"spellId":{"spellId":42897}}}`,
		);
		const arcaneMissiles = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":42846}}}`);
		const arcaneBarrage = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":44781}}}`);

		prepullActions.push(prepullMirrorImage);

		actions.push(
			...([
				berserking,
				hyperspeedAcceleration,
				combatPot,
				simple.missileBarrageBelowManaPercent > 0 ? arcaneMissilesWithMissileBarrageBelowMana : null,
				arcaneBlastBelowStacks,
				arcaneMisslesWithMissileBarrage,
				evocation,
				arcaneBlastAboveMana,
				simple.useArcaneBarrage ? arcaneBarrage : null,
				arcaneMissiles,
			].filter(a => a) as Array<APLAction>),
		);

		return APLRotation.create({
			prepullActions: prepullActions,
			priorityList: actions.map(action =>
				APLListItem.create({
					action: action,
				}),
			),
		});
	}, */

	raidSimPresets: [
		{
			spec: Spec.SpecArcaneMage,
			talents: Presets.ArcaneTalents.data,
			specOptions: Presets.DefaultArcaneOptions,
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
					1: Presets.ARCANE_P4_PRESET.gear,
					2: Presets.ARCANE_P3_PREBIS_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.ARCANE_P4_PRESET.gear,
					2: Presets.ARCANE_P3_PREBIS_PRESET.gear,
				},
			},
		},
	],
});

export class ArcaneMageSimUI extends IndividualSimUI<Spec.SpecArcaneMage> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecArcaneMage>) {
		super(parentElem, player, SPEC_CONFIG);

		player.sim.waitForInit().then(() => {
			new ReforgeOptimizer(this, {
				updateSoftCaps: softCaps => {
					const gear = player.getGear();
					const hasT114P = gear.getItemSetCount("Firelord's Vestments") >= 4;

					if (hasT114P) {
						const softCapToModify = softCaps.find(sc => sc.unitStat.equalsPseudoStat(PseudoStat.PseudoStatSpellHastePercent));
						if (softCapToModify) {
							softCapToModify.breakpoints[0] = 17.48;
						}
					}

					return softCaps;
				},
			});
		});
	}
}
