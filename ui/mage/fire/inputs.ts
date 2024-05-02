import * as InputHelpers from '../../core/components/input_helpers';
import { Player } from '../../core/player';
import { Spec } from '../../core/proto/common';
import { TypedEvent } from '../../core/typed_event';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const MageRotationConfig = {
	inputs: [
		// ********************************************************
		//                       FIRE INPUTS
		// ********************************************************
		// InputHelpers.makeRotationEnumInput<Spec.SpecFireMage, PrimaryFireSpell>({
		// 	fieldName: 'primaryFireSpell',
		// 	label: 'Primary Spell',
		// 	values: [
		// 		{ name: 'Fireball', value: PrimaryFireSpell.Fireball },
		// 		{ name: 'Frostfire Bolt', value: PrimaryFireSpell.FrostfireBolt },
		// 		{ name: 'Scorch', value: PrimaryFireSpell.Scorch },
		// 	],
		// 	showWhen: (player: Player<Spec.SpecFireMage>) => player.getTalentTree() == 1,
		// 	changeEmitter: (player: Player<Spec.SpecFireMage>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		// }),

		// InputHelpers.makeRotationBooleanInput<Spec.SpecFireMage>({
		// 	fieldName: 'maintainImprovedScorch',
		// 	label: 'Maintain Imp. Scorch',
		// 	labelTooltip: 'Always use Scorch when below 5 stacks, or < 4s remaining on debuff.',
		// 	showWhen: (player: Player<Spec.SpecFireMage>) => player.getTalents().improvedScorch > 0,
		// 	changeEmitter: (player: Player<Spec.SpecFireMage>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		// }),
	],
};
