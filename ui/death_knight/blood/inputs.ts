import * as InputHelpers from '../../core/components/input_helpers';
import { Player } from '../../core/player';
import { Spec } from '../../core/proto/common';
import {
	BloodDeathKnight_Rotation_BloodSpell as BloodSpell,
	BloodDeathKnight_Rotation_BloodTapPrio as BloodTapPrio,
	BloodDeathKnight_Rotation_Opener as Opener,
	BloodDeathKnight_Rotation_OptimizationSetting as OptimizationSetting,
	BloodDeathKnight_Rotation_Presence as Presence,
	DeathKnightMajorGlyph,
} from '../../core/proto/death_knight';
import { TypedEvent } from '../../core/typed_event';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const DrwPestiApply = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecBloodDeathKnight>({
	fieldName: 'drwPestiApply',
	label: 'DRW Pestilence Add',
	labelTooltip:
		'There is currently an interaction with DRW and pestilence where you can use pestilence to force DRW to apply diseases if they are already applied by the DK. It only works with Glyph of Disease and if there is an off target. This toggle forces the sim to assume there is an off target.',
	showWhen: (player: Player<Spec.SpecBloodDeathKnight>) =>
		player.getTalentTree() == 0 &&
		(player.getGlyphs().major1 == DeathKnightMajorGlyph.GlyphOfDisease ||
			player.getGlyphs().major2 == DeathKnightMajorGlyph.GlyphOfDisease ||
			player.getGlyphs().major3 == DeathKnightMajorGlyph.GlyphOfDisease),
	changeEmitter: (player: Player<Spec.SpecBloodDeathKnight>) =>
		TypedEvent.onAny([player.specOptionsChangeEmitter, player.rotationChangeEmitter, player.talentsChangeEmitter]),
});

export const DefensiveCdDelay = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecBloodDeathKnight>({
	fieldName: 'defensiveDelay',
	label: 'Defensives Delay',
	labelTooltip: 'Minimum delay between using more defensive cooldowns.',
});

export const BloodDeathKnightRotationConfig = {
	inputs: [
		InputHelpers.makeRotationEnumInput<Spec.SpecBloodDeathKnight, Presence>({
			fieldName: 'presence',
			label: 'Presence',
			labelTooltip: 'Presence to be in during the encounter.',
			values: [
				{ name: 'Blood', value: Presence.Blood },
				{ name: 'Frost', value: Presence.Frost },
				{ name: 'Unholy', value: Presence.Unholy },
			],
			changeEmitter: (player: Player<Spec.SpecBloodDeathKnight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecBloodDeathKnight, Opener>({
			fieldName: 'opener',
			label: 'Opener',
			labelTooltip:
				'Chose what opener to perform:<br>\
				<b>Regular</b>: Regular opener.<br>\
				<b>Threat</b>: Full IT spam for max threat.',
			values: [
				{ name: 'Regular', value: Opener.Regular },
				{ name: 'Threat', value: Opener.Threat },
			],
			changeEmitter: (player: Player<Spec.SpecBloodDeathKnight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecBloodDeathKnight, OptimizationSetting>({
			fieldName: 'optimizationSetting',
			label: 'Optimization Setting',
			labelTooltip:
				'Chose what metric to optimize:<br>\
				<b>Hps</b>: Prioritizes holding runes for healing after damage taken.<br>\
				<b>Tps</b>: Prioritizes spending runes for icy touch spam.',
			values: [
				{ name: 'Hps', value: OptimizationSetting.Hps },
				{ name: 'Tps', value: OptimizationSetting.Tps },
			],
			changeEmitter: (player: Player<Spec.SpecBloodDeathKnight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecBloodDeathKnight, BloodSpell>({
			fieldName: 'bloodSpell',
			label: 'Blood Spell',
			labelTooltip: 'Chose what blood rune spender to use.',
			values: [
				{ name: 'Blood Strike', value: BloodSpell.BloodStrike },
				{ name: 'Blood Boil', value: BloodSpell.BloodBoil },
				{ name: 'Heart Strike', value: BloodSpell.HeartStrike },
			],
			changeEmitter: (player: Player<Spec.SpecBloodDeathKnight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecBloodDeathKnight, BloodTapPrio>({
			fieldName: 'bloodTapPrio',
			label: 'Blood Tap',
			labelTooltip:
				'Chose how to use Blood Tap:<br>\
				<b>Use as Defensive Cooldown</b>: Use as defined in Cooldowns (Requires T10 4pc).<br>\
				<b>Offensive</b>: Use Blood Tap for extra Icy Touches.',
			values: [
				{ name: 'Use as Defensive Cooldown', value: BloodTapPrio.Defensive },
				{ name: 'Offensive', value: BloodTapPrio.Offensive },
			],
			changeEmitter: (player: Player<Spec.SpecBloodDeathKnight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
	],
};
