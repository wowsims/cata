import { Encounter } from '../encounter.js';
import { IndividualSimUI } from '../individual_sim_ui.js';
import { InputType, MobType, SpellSchool, Stat, Target, Target as TargetProto, TargetInput } from '../proto/common.js';
import { getStatName } from '../proto_utils/names.js';
import { Stats } from '../proto_utils/stats.js';
import { Raid } from '../raid.js';
import { SimUI } from '../sim_ui.js';
import { EventID, TypedEvent } from '../typed_event.js';
import { randomUUID } from '../utils';
import { BaseModal } from './base_modal.js';
import { Component } from './component.js';
import { Input } from './input.js';
import { BooleanPicker } from './pickers/boolean_picker.js';
import { EnumPicker } from './pickers/enum_picker.js';
import { ListItemPickerConfig, ListPicker } from './pickers/list_picker.jsx';
import { NumberPicker } from './pickers/number_picker.js';

export interface EncounterPickerConfig {
	showExecuteProportion: boolean;
}

export class EncounterPicker extends Component {
	constructor(parent: HTMLElement, modEncounter: Encounter, config: EncounterPickerConfig, simUI: SimUI) {
		super(parent, 'encounter-picker-root');

		addEncounterFieldPickers(this.rootElem, modEncounter, config.showExecuteProportion);

		// Need to wait so that the encounter and target presets will be loaded.
		modEncounter.sim.waitForInit().then(() => {
			const presetTargets = modEncounter.sim.db.getAllPresetTargets();

			// new EnumPicker<Encounter>(this.rootElem, modEncounter, {
			// 	extraCssClasses: ['damage-metrics', 'npc-picker'],
			// 	label: 'NPC',
			// 	labelTooltip: 'Selects a preset NPC configuration.',
			// 	values: [{ name: 'Custom', value: -1 }].concat(
			// 		presetTargets.map((pe, i) => {
			// 			return {
			// 				name: pe.path,
			// 				value: i,
			// 			};
			// 		}),
			// 	),
			// 	changedEvent: (encounter: Encounter) => encounter.changeEmitter,
			// 	getValue: (encounter: Encounter) => presetTargets.findIndex(pe => equalTargetsIgnoreInputs(encounter.primaryTarget, pe.target)),
			// 	setValue: (eventID: EventID, encounter: Encounter, newValue: number) => {
			// 		if (newValue != -1) {
			// 			encounter.applyPresetTarget(eventID, presetTargets[newValue], 0);
			// 		}
			// 	},
			// });

			const presetEncounters = modEncounter.sim.db.getAllPresetEncounters();
			new EnumPicker<Encounter>(this.rootElem, modEncounter, {
				id: 'encounter-preset-encouter',
				label: 'Encounter',
				//extraCssClasses: ['encounter-picker', 'mb-0', 'pe-2', 'order-first'],
				extraCssClasses: ['damage-metrics', 'npc-picker'],
				values: [{ name: 'Custom', value: -1 }].concat(
					presetEncounters.map((pe, i) => {
						return {
							name: pe.path,
							value: i,
						};
					}),
				),
				changedEvent: (encounter: Encounter) => encounter.changeEmitter,
				getValue: (encounter: Encounter) => presetEncounters.findIndex(pe => encounter.matchesPreset(pe)),
				setValue: (eventID: EventID, encounter: Encounter, newValue: number) => {
					if (newValue != -1) {
						encounter.applyPreset(eventID, presetEncounters[newValue]);
					}
				},
			});

			//new EnumPicker<Encounter>(this.rootElem, modEncounter, {
			//	label: 'Target Level',
			//	values: [
			//		{ name: '83', value: 83 },
			//		{ name: '82', value: 82 },
			//		{ name: '81', value: 81 },
			//		{ name: '80', value: 80 },
			//	],
			//	changedEvent: (encounter: Encounter) => encounter.changeEmitter,
			//	getValue: (encounter: Encounter) => encounter.primaryTarget.getLevel(),
			//	setValue: (eventID: EventID, encounter: Encounter, newValue: number) => {
			//		encounter.primaryTarget.setLevel(eventID, newValue);
			//	},
			//});

			//new EnumPicker(this.rootElem, modEncounter, {
			//	label: 'Mob Type',
			//	values: mobTypeEnumValues,
			//	changedEvent: (encounter: Encounter) => encounter.changeEmitter,
			//	getValue: (encounter: Encounter) => encounter.primaryTarget.getMobType(),
			//	setValue: (eventID: EventID, encounter: Encounter, newValue: number) => {
			//		encounter.primaryTarget.setMobType(eventID, newValue);
			//	},
			//});

			// Leaving this commented in case we want it later. But it takes up a lot of
			// screen space and none of these fields get changed much.
			//if (config.simpleTargetStats) {
			//	config.simpleTargetStats.forEach(stat => {
			//		new NumberPicker(this.rootElem, modEncounter, {
			//			label: statNames[stat],
			//			changedEvent: (encounter: Encounter) => encounter.changeEmitter,
			//			getValue: (encounter: Encounter) => encounter.primaryTarget.getStats().getStat(stat),
			//			setValue: (eventID: EventID, encounter: Encounter, newValue: number) => {
			//				encounter.primaryTarget.setStats(eventID, encounter.primaryTarget.getStats().withStat(stat, newValue));
			//			},
			//		});
			//	});
			//}

			if (simUI.isIndividualSim() && (simUI as IndividualSimUI<any>).player.canEnableTargetDummies()) {
				const player = (simUI as IndividualSimUI<any>).player;
				new NumberPicker(this.rootElem, simUI.sim.raid, {
					id: 'encounter-num-allies',
					label: 'Num Allies',
					labelTooltip: 'Number of allied players in the raid.',
					changedEvent: (raid: Raid) => TypedEvent.onAny([raid.targetDummiesChangeEmitter, player.itemSwapSettings.changeEmitter]),
					getValue: (raid: Raid) => raid.getTargetDummies(),
					setValue: (eventID: EventID, raid: Raid, newValue: number) => {
						raid.setTargetDummies(eventID, newValue);
					},
					showWhen: (raid: Raid) => {
						const shouldEnable = player.shouldEnableTargetDummies();

						if (!shouldEnable) {
							raid.setTargetDummies(TypedEvent.nextEventID(), 0);
						}

						return shouldEnable;
					},
				});
			}

			if (simUI.isIndividualSim() && (simUI as IndividualSimUI<any>).player.getPlayerSpec().isTankSpec) {
				new NumberPicker(this.rootElem, modEncounter, {
					id: 'encounter-min-base-damage',
					label: 'Min Base Damage',
					labelTooltip: 'Base damage for auto attacks, i.e. lowest roll with 0 AP against a 0-armor Player.',
					changedEvent: (encounter: Encounter) => encounter.changeEmitter,
					getValue: (encounter: Encounter) => encounter.primaryTarget.minBaseDamage,
					setValue: (eventID: EventID, encounter: Encounter, newValue: number) => {
						encounter.primaryTarget.minBaseDamage = newValue;
						encounter.targetsChangeEmitter.emit(eventID);
					},
				});
			}

			// Transfer Target Inputs from target Id if they dont match (possible when custom AI is selected)
			const targetIndex = presetTargets.findIndex(pe => modEncounter.primaryTarget.id == pe.target?.id);
			const targetInputs = presetTargets[targetIndex]?.target?.targetInputs || [];
			if (
				targetInputs.length != modEncounter.primaryTarget.targetInputs.length ||
				modEncounter.primaryTarget.targetInputs.some((ti, i) => ti.label != targetInputs[i].label)
			) {
				modEncounter.primaryTarget.targetInputs = targetInputs;
				modEncounter.targetsChangeEmitter.emit(TypedEvent.nextEventID());
			}

			makeTargetInputsPicker(this.rootElem, modEncounter, 0);

			const advancedModal = new AdvancedEncounterModal(simUI.rootElem, simUI, modEncounter);
			const advancedButton = document.createElement('button');
			advancedButton.classList.add('advanced-button', 'btn', 'btn-primary');
			advancedButton.textContent = 'Advanced';
			advancedButton.addEventListener('click', () => advancedModal.open());
			this.rootElem.appendChild(advancedButton);
		});
	}
}

class AdvancedEncounterModal extends BaseModal {
	private readonly encounter: Encounter;

	constructor(parent: HTMLElement, simUI: SimUI, encounter: Encounter) {
		super(parent, 'advanced-encounter-picker-modal', { disposeOnClose: false });

		this.encounter = encounter;

		this.addHeader();
		this.body.innerHTML = `
			<div class="encounter-header"></div>
			<div class="encounter-targets"></div>
		`;

		const header = this.rootElem.getElementsByClassName('encounter-header')[0] as HTMLElement;
		const targetsElem = this.rootElem.getElementsByClassName('encounter-targets')[0] as HTMLElement;

		addEncounterFieldPickers(header, this.encounter, true);
		if (!simUI.isIndividualSim()) {
			new BooleanPicker<Encounter>(header, encounter, {
				id: 'aem-use-health',
				label: 'Use Health',
				labelTooltip: 'Uses a damage limit in place of a duration limit. Damage limit is equal to sum of all targets health.',
				inline: true,
				changedEvent: (encounter: Encounter) => encounter.changeEmitter,
				getValue: (encounter: Encounter) => encounter.getUseHealth(),
				setValue: (eventID: EventID, encounter: Encounter, newValue: boolean) => {
					encounter.setUseHealth(eventID, newValue);
				},
			});
		}
		new ListPicker<Encounter, TargetProto>(targetsElem, this.encounter, {
			extraCssClasses: ['targets-picker', 'mb-0'],
			itemLabel: 'Target',
			changedEvent: (encounter: Encounter) => encounter.targetsChangeEmitter,
			getValue: (encounter: Encounter) => encounter.targets,
			setValue: (eventID: EventID, encounter: Encounter, newValue: Array<TargetProto>) => {
				encounter.targets = newValue;
				encounter.targetsChangeEmitter.emit(eventID);
			},
			newItem: () => Encounter.defaultTargetProto(),
			copyItem: (oldItem: TargetProto) => TargetProto.clone(oldItem),
			newItemPicker: (
				parent: HTMLElement,
				listPicker: ListPicker<Encounter, TargetProto>,
				index: number,
				config: ListItemPickerConfig<Encounter, TargetProto>,
			) => new TargetPicker(parent, encounter, index, config),
			minimumItems: 1,
		});
	}

	private addHeader() {
		const presetEncounters = this.encounter.sim.db.getAllPresetEncounters();

		new EnumPicker<Encounter>(this.header as HTMLElement, this.encounter, {
			id: 'aem-encounter-picker',
			label: 'Encounter',
			extraCssClasses: ['encounter-picker', 'mb-0', 'pe-2', 'order-first'],
			values: [{ name: 'Custom', value: -1 }].concat(
				presetEncounters.map((pe, i) => {
					return {
						name: pe.path,
						value: i,
					};
				}),
			),
			changedEvent: (encounter: Encounter) => encounter.changeEmitter,
			getValue: (encounter: Encounter) => presetEncounters.findIndex(pe => encounter.matchesPreset(pe)),
			setValue: (eventID: EventID, encounter: Encounter, newValue: number) => {
				if (newValue != -1) {
					encounter.applyPreset(eventID, presetEncounters[newValue]);
				}
			},
		});
	}
}

class TargetPicker extends Input<Encounter, TargetProto> {
	private readonly encounter: Encounter;
	private readonly targetIndex: number;
	private readonly aiPicker: Input<null, number>;
	private readonly levelPicker: Input<null, number>;
	private readonly mobTypePicker: Input<null, number>;
	private readonly tankIndexPicker: Input<null, number>;
	private readonly statPickers: Array<Input<null, number>>;
	private readonly swingSpeedPicker: Input<null, number>;
	private readonly minBaseDamagePicker: Input<null, number>;
	private readonly dualWieldPicker: Input<null, boolean>;
	private readonly dwMissPenaltyPicker: Input<null, boolean>;
	private readonly parryHastePicker: Input<null, boolean>;
	private readonly spellSchoolPicker: Input<null, number>;
	private readonly damageSpreadPicker: Input<null, number>;
	private readonly targetInputPickers: ListPicker<Encounter, TargetInput>;

	private getTarget(): TargetProto {
		return this.encounter.targets[this.targetIndex] || Target.create();
	}

	constructor(parent: HTMLElement, encounter: Encounter, targetIndex: number, config: ListItemPickerConfig<Encounter, TargetProto>) {
		super(parent, 'target-picker-root', encounter, config);
		this.encounter = encounter;
		this.targetIndex = targetIndex;

		this.rootElem.innerHTML = `
			<div class="picker-group target-picker-section target-picker-section1"></div>
			<div class="picker-group target-picker-section target-picker-section2"></div>
			<div class="picker-group target-picker-section target-picker-section3 threat-metrics"></div>
		`;

		const section1 = this.rootElem.querySelector<HTMLElement>('.target-picker-section1')!;
		const section2 = this.rootElem.querySelector<HTMLElement>('.target-picker-section2')!;
		const section3 = this.rootElem.querySelector<HTMLElement>('.target-picker-section3')!;

		const presetTargets = encounter.sim.db.getAllPresetTargets();
		new EnumPicker<null>(section1, null, {
			id: 'target-picker-npc',
			extraCssClasses: ['npc-picker'],
			label: 'NPC',
			labelTooltip: 'Selects a preset NPC configuration.',
			values: [{ name: 'Custom', value: -1 }].concat(
				presetTargets.map((pe, i) => {
					return {
						name: pe.path,
						value: i,
					};
				}),
			),
			changedEvent: () => encounter.targetsChangeEmitter,
			getValue: () => presetTargets.findIndex(pe => equalTargetsIgnoreInputs(this.getTarget(), pe.target)),
			setValue: (eventID: EventID, _: null, newValue: number) => {
				if (newValue != -1) {
					encounter.applyPresetTarget(eventID, presetTargets[newValue], this.targetIndex);
					encounter.targetsChangeEmitter.emit(eventID);
				}
			},
		});

		this.aiPicker = new EnumPicker<null>(section1, null, {
			id: 'target-picker-ai',
			extraCssClasses: ['ai-picker'],
			label: 'AI',
			labelTooltip: `
				<p>Determines the target\'s ability rotation.</p>
				<p>Note that most rotations are not yet implemented.</p>
			`,
			values: [{ name: 'None', value: 0 }].concat(
				presetTargets.map(pe => {
					return {
						name: pe.path,
						value: pe.target!.id,
					};
				}),
			),
			changedEvent: () => encounter.targetsChangeEmitter,
			getValue: () => this.getTarget().id,
			setValue: (eventID: EventID, _: null, newValue: number) => {
				const target = this.getTarget();
				target.id = newValue;

				// Transfer Target Inputs from the AI of the selected target
				target.targetInputs = (presetTargets.find(pe => target.id == pe.target?.id)?.target?.targetInputs || []).map(ti => TargetInput.clone(ti));

				encounter.targetsChangeEmitter.emit(eventID);
			},
		});

		this.levelPicker = new EnumPicker<null>(section1, null, {
			id: 'target-picker-level',
			label: 'Level',
			values: [
				{ name: '93', value: 93 },
				{ name: '92', value: 92 },
				{ name: '91', value: 91 },
				{ name: '90', value: 90 },
				{ name: '88', value: 88 },

			],
			changedEvent: () => encounter.targetsChangeEmitter,
			getValue: () => this.getTarget().level,
			setValue: (eventID: EventID, _: null, newValue: number) => {
				this.getTarget().level = newValue;
				encounter.targetsChangeEmitter.emit(eventID);
			},
		});
		this.mobTypePicker = new EnumPicker(section1, null, {
			id: 'target-picker-mob-type',
			label: 'Mob Type',
			values: mobTypeEnumValues,
			changedEvent: () => encounter.targetsChangeEmitter,
			getValue: () => this.getTarget().mobType,
			setValue: (eventID: EventID, _: null, newValue: number) => {
				this.getTarget().mobType = newValue;
				encounter.targetsChangeEmitter.emit(eventID);
			},
		});
		this.tankIndexPicker = new EnumPicker<null>(section1, null, {
			id: 'target-picker-tanked-by',
			extraCssClasses: ['threat-metrics'],
			label: 'Tanked By',
			labelTooltip:
				'Determines which player in the raid this enemy will attack. If no player is assigned to the specified tank slot, this enemy will not attack.',
			values: [
				{ name: 'None', value: -1 },
				{ name: 'Main Tank', value: 0 },
				{ name: 'Tank 2', value: 1 },
				{ name: 'Tank 3', value: 2 },
				{ name: 'Tank 4', value: 3 },
			],
			changedEvent: () => encounter.targetsChangeEmitter,
			getValue: () => this.getTarget().tankIndex,
			setValue: (eventID: EventID, _: null, newValue: number) => {
				this.getTarget().tankIndex = newValue;
				encounter.targetsChangeEmitter.emit(eventID);
			},
		});

		this.targetInputPickers = makeTargetInputsPicker(section1, encounter, this.targetIndex);

		this.statPickers = ALL_TARGET_STATS.map(statData => {
			const stat = statData.stat;
			return new NumberPicker(section2, null, {
				id: `target-${this.targetIndex}-picker-stats-${statData.stat}`,
				inline: true,
				extraCssClasses: statData.extraCssClasses,
				label: getStatName(stat),
				labelTooltip: statData.tooltip,
				changedEvent: () => encounter.targetsChangeEmitter,
				getValue: () => this.getTarget().stats[stat],
				setValue: (eventID: EventID, _: null, newValue: number) => {
					this.getTarget().stats[stat] = newValue;
					encounter.targetsChangeEmitter.emit(eventID);
				},
			});
		});

		this.swingSpeedPicker = new NumberPicker(section3, null, {
			id: `target-${this.targetIndex}-picker-swing-speed`,
			label: 'Swing Speed',
			labelTooltip: 'Time in seconds between auto attacks. Set to 0 to disable auto attacks.',
			float: true,
			changedEvent: () => encounter.targetsChangeEmitter,
			getValue: () => this.getTarget().swingSpeed,
			setValue: (eventID: EventID, _: null, newValue: number) => {
				this.getTarget().swingSpeed = newValue;
				encounter.targetsChangeEmitter.emit(eventID);
			},
		});
		this.minBaseDamagePicker = new NumberPicker(section3, null, {
			id: `target-${this.targetIndex}-picker-min-base-damage`,
			label: 'Min Base Damage',
			labelTooltip: 'Base damage for auto attacks, i.e. lowest roll with 0 AP against a 0-armor Player.',
			changedEvent: () => encounter.targetsChangeEmitter,
			getValue: () => this.getTarget().minBaseDamage,
			setValue: (eventID: EventID, _: null, newValue: number) => {
				this.getTarget().minBaseDamage = newValue;
				encounter.targetsChangeEmitter.emit(eventID);
			},
		});
		this.damageSpreadPicker = new NumberPicker(section3, null, {
			id: `target-${this.targetIndex}-picker-damage-spread`,
			label: 'Damage Spread',
			labelTooltip: 'Fractional spread between the minimum and maximum auto-attack damage from this enemy at 0 Attack Power.',
			float: true,
			changedEvent: () => encounter.targetsChangeEmitter,
			getValue: () => this.getTarget().damageSpread,
			setValue: (eventID: EventID, _: null, newValue: number) => {
				this.getTarget().damageSpread = newValue;
				encounter.targetsChangeEmitter.emit(eventID);
			},
		});
		this.dualWieldPicker = new BooleanPicker(section3, null, {
			id: `target-${this.targetIndex}-picker-dual-wield`,
			label: 'Dual Wield',
			labelTooltip: 'Uses 2 separate weapons to attack.',
			inline: true,
			reverse: true,
			changedEvent: () => encounter.targetsChangeEmitter,
			getValue: () => this.getTarget().dualWield,
			setValue: (eventID: EventID, _: null, newValue: boolean) => {
				this.getTarget().dualWield = newValue;
				encounter.targetsChangeEmitter.emit(eventID);
			},
		});
		this.dwMissPenaltyPicker = new BooleanPicker(section3, null, {
			id: `target-${this.targetIndex}-picker-dw-miss-penalty`,
			label: 'DW Miss Penalty',
			labelTooltip:
				'Enables the Dual Wield Miss Penalty (+19% chance to miss) if dual wielding. Bosses in Hyjal/BT/SWP usually have this disabled to stop tanks from avoidance stacking.',
			inline: true,
			reverse: true,
			changedEvent: () => encounter.targetsChangeEmitter,
			getValue: () => this.getTarget().dualWieldPenalty,
			setValue: (eventID: EventID, _: null, newValue: boolean) => {
				this.getTarget().dualWieldPenalty = newValue;
				encounter.targetsChangeEmitter.emit(eventID);
			},
			enableWhen: () => this.getTarget().dualWield,
		});
		this.parryHastePicker = new BooleanPicker(section3, null, {
			id: `target-${this.targetIndex}-picker-parry-haste`,
			label: 'Parry Haste',
			labelTooltip: 'Whether this enemy will gain parry haste when parrying attacks.',
			inline: true,
			reverse: true,
			changedEvent: () => encounter.targetsChangeEmitter,
			getValue: () => this.getTarget().parryHaste,
			setValue: (eventID: EventID, _: null, newValue: boolean) => {
				this.getTarget().parryHaste = newValue;
				encounter.targetsChangeEmitter.emit(eventID);
			},
		});
		this.spellSchoolPicker = new EnumPicker<null>(section3, null, {
			id: `target-${this.targetIndex}-picker-spell-school`,
			label: 'Spell School',
			labelTooltip: 'Type of damage caused by auto attacks. This is usually Physical, but some enemies have elemental attacks.',
			values: [
				{ name: 'Physical', value: SpellSchool.SpellSchoolPhysical },
				{ name: 'Arcane', value: SpellSchool.SpellSchoolArcane },
				{ name: 'Fire', value: SpellSchool.SpellSchoolFire },
				{ name: 'Frost', value: SpellSchool.SpellSchoolFrost },
				{ name: 'Holy', value: SpellSchool.SpellSchoolHoly },
				{ name: 'Nature', value: SpellSchool.SpellSchoolNature },
				{ name: 'Shadow', value: SpellSchool.SpellSchoolShadow },
			],
			changedEvent: () => encounter.targetsChangeEmitter,
			getValue: () => this.getTarget().spellSchool,
			setValue: (eventID: EventID, _: null, newValue: number) => {
				this.getTarget().spellSchool = newValue;
				encounter.targetsChangeEmitter.emit(eventID);
			},
		});

		this.init();
	}

	getInputElem(): HTMLElement | null {
		return null;
	}
	getInputValue(): TargetProto {
		return TargetProto.create({
			id: this.aiPicker.getInputValue(),
			level: this.levelPicker.getInputValue(),
			mobType: this.mobTypePicker.getInputValue(),
			tankIndex: this.tankIndexPicker.getInputValue(),
			swingSpeed: this.swingSpeedPicker.getInputValue(),
			minBaseDamage: this.minBaseDamagePicker.getInputValue(),
			dualWield: this.dualWieldPicker.getInputValue(),
			dualWieldPenalty: this.dwMissPenaltyPicker.getInputValue(),
			parryHaste: this.parryHastePicker.getInputValue(),
			spellSchool: this.spellSchoolPicker.getInputValue(),
			damageSpread: this.damageSpreadPicker.getInputValue(),
			stats: this.statPickers
				.map(picker => picker.getInputValue())
				.map((statValue, i) => new Stats().withStat(ALL_TARGET_STATS[i].stat, statValue))
				.reduce((totalStats, curStats) => totalStats.add(curStats))
				.asProtoArray(),
			targetInputs: this.targetInputPickers.getInputValue(),
		});
	}
	setInputValue(newValue: TargetProto) {
		if (!newValue) {
			return;
		}
		this.aiPicker.setInputValue(newValue.id);
		this.levelPicker.setInputValue(newValue.level);
		this.mobTypePicker.setInputValue(newValue.mobType);
		this.tankIndexPicker.setInputValue(newValue.tankIndex);
		this.swingSpeedPicker.setInputValue(newValue.swingSpeed);
		this.minBaseDamagePicker.setInputValue(newValue.minBaseDamage);
		this.dualWieldPicker.setInputValue(newValue.dualWield);
		this.dwMissPenaltyPicker.setInputValue(newValue.dualWieldPenalty);
		this.parryHastePicker.setInputValue(newValue.parryHaste);
		this.spellSchoolPicker.setInputValue(newValue.spellSchool);
		this.damageSpreadPicker.setInputValue(newValue.damageSpread);
		ALL_TARGET_STATS.forEach((statData, i) => this.statPickers[i].setInputValue(newValue.stats[statData.stat]));
		this.targetInputPickers.setInputValue(newValue.targetInputs);
	}
}

class TargetInputPicker extends Input<Encounter, TargetInput> {
	private readonly encounter: Encounter;
	private readonly targetIndex: number;
	private readonly targetInputIndex: number;

	private boolPicker: Input<null, boolean> | null;
	private numberPicker: Input<null, number> | null;
	private enumPicker: EnumPicker<null> | null;

	private getTargetInput(): TargetInput {
		return this.encounter.targets[this.targetIndex].targetInputs[this.targetInputIndex] || TargetInput.create();
	}

	private clearPickers() {
		if (this.boolPicker) {
			this.boolPicker.rootElem.remove();
			this.boolPicker = null;
		}
		if (this.numberPicker) {
			this.numberPicker.rootElem.remove();
			this.numberPicker = null;
		}
		if (this.enumPicker) {
			this.enumPicker.rootElem.remove();
			this.enumPicker = null;
		}
	}

	constructor(
		parent: HTMLElement,
		encounter: Encounter,
		targetIndex: number,
		targetInputIndex: number,
		config: ListItemPickerConfig<Encounter, TargetInput>,
	) {
		super(parent, 'target-input-picker-root', encounter, config);
		this.encounter = encounter;
		this.targetIndex = targetIndex;
		this.targetInputIndex = targetInputIndex;

		this.boolPicker = null;
		this.numberPicker = null;
		this.enumPicker = null;
		this.init();
	}

	getInputElem(): HTMLElement | null {
		return this.rootElem;
	}
	getInputValue(): TargetInput {
		return TargetInput.create({
			boolValue: this.boolPicker ? this.boolPicker.getInputValue() : undefined,
			numberValue: this.numberPicker ? this.numberPicker.getInputValue() : undefined,
			enumValue: this.enumPicker ? this.enumPicker.getInputValue() : undefined,
		});
	}
	setInputValue(newValue: TargetInput) {
		if (!newValue) {
			return;
		}

		if (newValue.inputType == InputType.Number) {
			if (this.numberPicker && this.numberPicker.inputConfig.label === newValue.label) {
				return;
			}

			this.clearPickers();
			this.numberPicker = new NumberPicker(this.rootElem, null, {
				id: randomUUID(),
				float: true,
				label: newValue.label,
				labelTooltip: newValue.tooltip,
				changedEvent: () => this.encounter.targetsChangeEmitter,
				getValue: () => this.getTargetInput().numberValue,
				setValue: (eventID: EventID, _: null, newValue: number) => {
					this.getTargetInput().numberValue = newValue;
					this.encounter.targetsChangeEmitter.emit(eventID);
				},
			});
		} else if (newValue.inputType == InputType.Bool) {
			if (this.boolPicker && this.boolPicker.inputConfig.label === newValue.label) {
				return;
			}

			this.clearPickers();
			this.boolPicker = new BooleanPicker(this.rootElem, null, {
				id: randomUUID(),
				label: newValue.label,
				labelTooltip: newValue.tooltip,
				extraCssClasses: ['input-inline'],
				changedEvent: () => this.encounter.targetsChangeEmitter,
				getValue: () => this.getTargetInput().boolValue,
				setValue: (eventID: EventID, _: null, newValue: boolean) => {
					this.getTargetInput().boolValue = newValue;
					this.encounter.targetsChangeEmitter.emit(eventID);
				},
			});
		} else if (newValue.inputType == InputType.Enum) {
			this.clearPickers();
			this.enumPicker = new EnumPicker<null>(this.rootElem, null, {
				id: randomUUID(),
				label: newValue.label,
				values: newValue.enumOptions.map((option, index) => {
					return { value: index, name: option };
				}),
				changedEvent: () => this.encounter.targetsChangeEmitter,
				getValue: () => this.getTargetInput().enumValue,
				setValue: (eventID: EventID, _: null, newValue: number) => {
					this.getTargetInput().enumValue = newValue;
					this.encounter.targetsChangeEmitter.emit(eventID);
				},
			});
		}
	}
}

function addEncounterFieldPickers(rootElem: HTMLElement, encounter: Encounter, showExecuteProportion: boolean) {
	const durationGroup = Input.newGroupContainer();
	rootElem.appendChild(durationGroup);

	new NumberPicker(durationGroup, encounter, {
		id: 'encounter-duration',
		label: 'Duration',
		labelTooltip: 'The fight length for each sim iteration, in seconds.',
		changedEvent: (encounter: Encounter) => encounter.changeEmitter,
		getValue: (encounter: Encounter) => encounter.getDuration(),
		setValue: (eventID: EventID, encounter: Encounter, newValue: number) => {
			encounter.setDuration(eventID, newValue);
		},
		enableWhen: _obj => {
			return !encounter.getUseHealth();
		},
	});
	new NumberPicker(durationGroup, encounter, {
		id: 'encounter-duration-variation',
		label: 'Duration +/-',
		labelTooltip:
			'Adds a random amount of time, in seconds, between [value, -1 * value] to each sim iteration. For example, setting Duration to 180 and Duration +/- to 10 will result in random durations between 170s and 190s.',
		changedEvent: (encounter: Encounter) => encounter.changeEmitter,
		getValue: (encounter: Encounter) => encounter.getDurationVariation(),
		setValue: (eventID: EventID, encounter: Encounter, newValue: number) => {
			encounter.setDurationVariation(eventID, newValue);
		},
		enableWhen: _obj => {
			return !encounter.getUseHealth();
		},
	});

	if (showExecuteProportion) {
		const executeGroup = Input.newGroupContainer();
		executeGroup.classList.add('execute-group');
		rootElem.appendChild(executeGroup);

		new NumberPicker(executeGroup, encounter, {
			id: 'encounter-execute-proportion',
			label: 'Execute Duration 20 (%)',
			labelTooltip:
				'Percentage of the total encounter duration, for which the targets will be considered to be in execute range (< 20% HP) for the purpose of effects like Warrior Execute or Mage Molten Fury.',
			changedEvent: (encounter: Encounter) => encounter.changeEmitter,
			getValue: (encounter: Encounter) => encounter.getExecuteProportion20() * 100,
			setValue: (eventID: EventID, encounter: Encounter, newValue: number) => {
				encounter.setExecuteProportion20(eventID, newValue / 100);
			},
			enableWhen: _obj => {
				return !encounter.getUseHealth();
			},
		});
		new NumberPicker(executeGroup, encounter, {
			id: 'encounter-execute-proportion-25',
			label: 'Execute Duration 25 (%)',
			labelTooltip:
				"Percentage of the total encounter duration, for which the targets will be considered to be in execute range (< 25% HP) for the purpose of effects like Warlock's Drain Soul.",
			changedEvent: (encounter: Encounter) => encounter.changeEmitter,
			getValue: (encounter: Encounter) => encounter.getExecuteProportion25() * 100,
			setValue: (eventID: EventID, encounter: Encounter, newValue: number) => {
				encounter.setExecuteProportion25(eventID, newValue / 100);
			},
			enableWhen: _obj => {
				return !encounter.getUseHealth();
			},
		});
		new NumberPicker(executeGroup, encounter, {
			id: 'encounter-execute-proportion-35',
			label: 'Execute Duration 35 (%)',
			labelTooltip:
				'Percentage of the total encounter duration, for which the targets will be considered to be in execute range (< 35% HP) for the purpose of effects like Death Knight Soul Reaper.',
			changedEvent: (encounter: Encounter) => encounter.changeEmitter,
			getValue: (encounter: Encounter) => encounter.getExecuteProportion35() * 100,
			setValue: (eventID: EventID, encounter: Encounter, newValue: number) => {
				encounter.setExecuteProportion35(eventID, newValue / 100);
			},
			enableWhen: _obj => {
				return !encounter.getUseHealth();
			},
		});
		new NumberPicker(executeGroup, encounter, {
			id: 'encounter-execute-proportion-45',
			label: 'Execute Duration 45 (%)',
			labelTooltip:
				'Percentage of the total encounter duration, for which the targets will be considered to be in execute range (< 45% HP) for the purpose of effects like Death Knight Soul Reaper with T15 DPS 4pc.',
			changedEvent: (encounter: Encounter) => encounter.changeEmitter,
			getValue: (encounter: Encounter) => encounter.getExecuteProportion45() * 100,
			setValue: (eventID: EventID, encounter: Encounter, newValue: number) => {
				encounter.setExecuteProportion45(eventID, newValue / 100);
			},
			enableWhen: _obj => {
				return !encounter.getUseHealth();
			},
		});
		new NumberPicker(executeGroup, encounter, {
			id: 'encounter-execute-proportion-90',
			label: 'Duration spent below high-HP regime (%)',
			labelTooltip:
				'Percentage of the total encounter duration, for which the targets are considered out of range for effects like Hunter Careful Aim (<90% HP) or Druid Predatory Strikes (<80% HP).',
			changedEvent: (encounter: Encounter) => encounter.changeEmitter,
			getValue: (encounter: Encounter) => encounter.getExecuteProportion90() * 100,
			setValue: (eventID: EventID, encounter: Encounter, newValue: number) => {
				encounter.setExecuteProportion90(eventID, newValue / 100);
			},
			enableWhen: _obj => {
				return !encounter.getUseHealth();
			},
		});
	}
}

function makeTargetInputsPicker(parent: HTMLElement, encounter: Encounter, targetIndex: number) {
	return new ListPicker<Encounter, TargetInput>(parent, encounter, {
		allowedActions: [],
		itemLabel: 'Target Input',
		extraCssClasses: ['mt-2'],
		isCompact: true,
		changedEvent: (encounter: Encounter) => encounter.targetsChangeEmitter,
		getValue: (encounter: Encounter) => encounter.targets[targetIndex].targetInputs,
		setValue: (eventID: EventID, encounter: Encounter, newValue: Array<TargetInput>) => {
			encounter.targets[targetIndex].targetInputs = newValue;
			encounter.targetsChangeEmitter.emit(eventID);
		},
		newItem: () => TargetInput.create(),
		copyItem: (oldItem: TargetInput) => TargetInput.clone(oldItem),
		newItemPicker: (
			parent: HTMLElement,
			listPicker: ListPicker<Encounter, TargetInput>,
			index: number,
			config: ListItemPickerConfig<Encounter, TargetInput>,
		) => new TargetInputPicker(parent, encounter, targetIndex, index, config),
	});
}

function equalTargetsIgnoreInputs(target1: TargetProto | undefined, target2: TargetProto | undefined): boolean {
	if (!!target1 != !!target2) {
		return false;
	}
	if (!target1) {
		return true;
	}
	const modTarget2 = TargetProto.clone(target2!);
	modTarget2.targetInputs = target1.targetInputs;
	return TargetProto.equals(target1, modTarget2);
}

const ALL_TARGET_STATS: Array<{ stat: Stat; tooltip: string; extraCssClasses: Array<string> }> = [
	{ stat: Stat.StatHealth, tooltip: '', extraCssClasses: [] },
	{ stat: Stat.StatArmor, tooltip: '', extraCssClasses: [] },
	{ stat: Stat.StatAttackPower, tooltip: '', extraCssClasses: ['threat-metrics'] },
];

const mobTypeEnumValues = [
	{ name: 'None', value: MobType.MobTypeUnknown },
	{ name: 'Beast', value: MobType.MobTypeBeast },
	{ name: 'Demon', value: MobType.MobTypeDemon },
	{ name: 'Dragonkin', value: MobType.MobTypeDragonkin },
	{ name: 'Elemental', value: MobType.MobTypeElemental },
	{ name: 'Giant', value: MobType.MobTypeGiant },
	{ name: 'Humanoid', value: MobType.MobTypeHumanoid },
	{ name: 'Mechanical', value: MobType.MobTypeMechanical },
	{ name: 'Undead', value: MobType.MobTypeUndead },
];
