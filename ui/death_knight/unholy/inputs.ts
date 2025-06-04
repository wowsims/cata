
// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

// export const SelfUnholyFrenzy = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecUnholyDeathKnight>({
// 	fieldName: 'unholyFrenzyTarget',
// 	label: 'Self Unholy Frenzy',
// 	labelTooltip: 'Cast Unholy Frenzy on yourself.',
// 	extraCssClasses: ['within-raid-sim-hide'],
// 	getValue: (player: Player<Spec.SpecUnholyDeathKnight>) => player.getSpecOptions().unholyFrenzyTarget?.type == UnitType.Player,
// 	setValue: (eventID: EventID, player: Player<Spec.SpecUnholyDeathKnight>, newValue: boolean) => {
// 		const newOptions = player.getSpecOptions();
// 		newOptions.unholyFrenzyTarget = UnitReference.create({
// 			type: newValue ? UnitType.Player : UnitType.Unknown,
// 			index: 0,
// 		});
// 		player.setSpecOptions(eventID, newOptions);
// 	},
// 	showWhen: (player: Player<Spec.SpecUnholyDeathKnight>) => player.getTalents().hysteria,
// 	changeEmitter: (player: Player<Spec.SpecUnholyDeathKnight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
// });
