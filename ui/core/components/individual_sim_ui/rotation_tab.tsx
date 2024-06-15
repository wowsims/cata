import { ref } from 'tsx-vanilla';

import * as Tooltips from '../../constants/tooltips';
import { IndividualSimUI, InputSection } from '../../individual_sim_ui';
import { Player } from '../../player';
import { APLRotation, APLRotation_Type as APLRotationType } from '../../proto/apl';
import { SavedRotation } from '../../proto/ui';
import { EventID, TypedEvent } from '../../typed_event';
import { ContentBlock } from '../content_block';
import * as IconInputs from '../icon_inputs';
import { Input } from '../input';
import { BooleanPicker } from '../pickers/boolean_picker';
import { EnumPicker } from '../pickers/enum_picker';
import { NumberPicker } from '../pickers/number_picker';
import { SavedDataManager } from '../saved_data_manager';
import { SimTab } from '../sim_tab';
import { APLRotationPicker } from './apl_rotation_picker';
import { CooldownsPicker } from './cooldowns_picker';

export class RotationTab extends SimTab {
	protected simUI: IndividualSimUI<any>;

	readonly leftPanel: HTMLElement;
	readonly rightPanel: HTMLElement;

	constructor(parentElem: HTMLElement, simUI: IndividualSimUI<any>) {
		super(parentElem, simUI, { identifier: 'rotation-tab', title: 'Rotation' });
		this.simUI = simUI;

		this.leftPanel = (<div className="rotation-tab-left tab-panel-left" />) as HTMLElement;

		this.rightPanel = (<div className="rotation-tab-right tab-panel-right" />) as HTMLElement;

		this.contentContainer.appendChild(this.leftPanel);
		this.contentContainer.appendChild(this.rightPanel);

		this.buildTabContent();

		this.updateSections();
		this.simUI.player.rotationChangeEmitter.on(() => this.updateSections());
	}

	protected buildTabContent() {
		this.buildHeader();

		this.buildAutoContent();
		this.buildAplContent();
		this.buildSimpleContent();

		this.buildSavedDataPickers();
	}

	private updateSections() {
		this.rootElem.classList.remove('rotation-type-auto', 'rotation-type-simple', 'rotation-type-apl', 'rotation-type-legacy');

		const rotationType = this.simUI.player.getRotationType();
		let rotationClass = '';
		switch (rotationType) {
			case APLRotationType.TypeAuto:
				rotationClass = 'rotation-type-auto';
				break;
			case APLRotationType.TypeSimple:
				rotationClass = 'rotation-type-simple';
				break;
			case APLRotationType.TypeAPL:
				rotationClass = 'rotation-type-apl';
				break;
		}

		this.rootElem.classList.add(rotationClass);
	}

	private buildHeader() {
		const headerRef = ref<HTMLDivElement>();
		const resetButtonRef = ref<HTMLButtonElement>();
		const rotationTypeSelectRef = ref<HTMLDivElement>();
		this.leftPanel.appendChild(
			<div ref={headerRef} className="rotation-tab-header d-flex justify-content-between align-items-baseline">
				<div ref={rotationTypeSelectRef} />
				<button ref={resetButtonRef} className="btn btn-sm btn-link btn-reset summary-table-reset-button">
					<i className="fas fa-times me-1"></i>
					Reset APL
				</button>
			</div>,
		);

		resetButtonRef.value!.addEventListener('click', () => {
			this.simUI.applyEmptyAplRotation(TypedEvent.nextEventID());
		});

		this.simUI.player.rotationChangeEmitter.on(() => {
			const type = this.simUI.player.getRotationType();
			resetButtonRef.value?.classList[type === APLRotationType.TypeAPL ? 'remove' : 'add']('hide');
		});

		new EnumPicker(rotationTypeSelectRef.value!, this.simUI.player, {
			extraCssClasses: ['w-auto'],
			id: 'rotation-tab-rotation-type',
			label: 'Rotation Type',
			labelTooltip: 'Which set of options to use for specifying the rotation.',
			inline: true,
			values: this.simUI.player.hasSimpleRotationGenerator()
				? [
						{ value: APLRotationType.TypeAuto, name: 'Auto' },
						{ value: APLRotationType.TypeSimple, name: 'Simple' },
						{ value: APLRotationType.TypeAPL, name: 'APL' },
				  ]
				: [
						{ value: APLRotationType.TypeAuto, name: 'Auto' },
						{ value: APLRotationType.TypeAPL, name: 'APL' },
				  ],
			changedEvent: (player: Player<any>) => player.rotationChangeEmitter,
			getValue: (player: Player<any>) => player.getRotationType(),
			setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
				player.aplRotation.type = newValue;
				player.rotationChangeEmitter.emit(eventID);
			},
		});
	}

	private buildAutoContent() {
		this.leftPanel.appendChild(<div className="rotation-tab-auto" />);
	}

	private buildAplContent() {
		const contentRef = ref<HTMLDivElement>();
		this.leftPanel.appendChild(<div ref={contentRef} className="rotation-tab-apl" />);

		new APLRotationPicker(contentRef.value!, this.simUI, this.simUI.player);
	}

	private buildSimpleContent() {
		if (!this.simUI.player.hasSimpleRotationGenerator() || !this.simUI.individualConfig.rotationInputs) {
			return;
		}
		const cssClass = 'rotation-tab-simple';

		const contentBlock = new ContentBlock(this.leftPanel, 'rotation-settings', {
			header: { title: 'Rotation' },
		});
		contentBlock.rootElem.classList.add(cssClass);

		const rotationIconGroup = Input.newGroupContainer();
		rotationIconGroup.classList.add('rotation-icon-group', 'icon-group');
		contentBlock.bodyElement.appendChild(rotationIconGroup);

		if (this.simUI.individualConfig.rotationIconInputs?.length) {
			this.configureIconSection(
				rotationIconGroup,
				this.simUI.individualConfig.rotationIconInputs.map(iconInput => IconInputs.buildIconInput(rotationIconGroup, this.simUI.player, iconInput)),
				true,
			);
		}

		this.configureInputSection(contentBlock.bodyElement, this.simUI.individualConfig.rotationInputs);
		const cooldownsContentBlock = new ContentBlock(this.leftPanel, 'cooldown-settings', {
			header: { title: 'Cooldowns', tooltip: Tooltips.COOLDOWNS_SECTION },
		});
		cooldownsContentBlock.rootElem.classList.add(cssClass);

		new CooldownsPicker(cooldownsContentBlock.bodyElement, this.simUI.player);
	}

	private configureInputSection(sectionElem: HTMLElement, sectionConfig: InputSection) {
		sectionConfig.inputs.forEach(inputConfig => {
			inputConfig.extraCssClasses = [...(inputConfig.extraCssClasses || []), 'input-inline'];
			if (inputConfig.type == 'number') {
				new NumberPicker(sectionElem, this.simUI.player, { ...inputConfig, inline: true });
			} else if (inputConfig.type == 'boolean') {
				new BooleanPicker(sectionElem, this.simUI.player, { ...inputConfig, inline: true });
			} else if (inputConfig.type == 'enum') {
				new EnumPicker(sectionElem, this.simUI.player, { ...inputConfig, inline: true });
			}
		});
	}

	private configureIconSection(sectionElem: HTMLElement, iconPickers: Array<any>, adjustColumns?: boolean) {
		if (!iconPickers.length) {
			sectionElem.classList.add('hide');
		} else if (adjustColumns) {
			if (iconPickers.length <= 4) {
				sectionElem.style.gridTemplateColumns = `repeat(${iconPickers.length}, 1fr)`;
			} else if (iconPickers.length > 4 && iconPickers.length < 8) {
				sectionElem.style.gridTemplateColumns = `repeat(${Math.ceil(iconPickers.length / 2)}, 1fr)`;
			}
		}
	}

	private buildSavedDataPickers() {
		const savedRotationsManager = new SavedDataManager<Player<any>, SavedRotation>(this.rightPanel, this.simUI.player, {
			label: 'Rotation',
			header: { title: 'Saved Rotations' },
			storageKey: this.simUI.getSavedRotationStorageKey(),
			getData: (player: Player<any>) =>
				SavedRotation.create({
					rotation: APLRotation.clone(player.aplRotation),
				}),
			setData: (eventID: EventID, player: Player<any>, newRotation: SavedRotation) => {
				TypedEvent.freezeAllAndDo(() => {
					player.setAplRotation(eventID, newRotation.rotation || APLRotation.create());
				});
			},
			changeEmitters: [this.simUI.player.rotationChangeEmitter, this.simUI.player.talentsChangeEmitter],
			equals: (a: SavedRotation, b: SavedRotation) => {
				// Uncomment this to debug equivalence checks with preset rotations (e.g. the chip doesn't highlight)
				//console.log(`Rot A: ${SavedRotation.toJsonString(a, {prettySpaces: 2})}\n\nRot B: ${SavedRotation.toJsonString(b, {prettySpaces: 2})}`);
				return SavedRotation.equals(a, b);
			},
			toJson: (a: SavedRotation) => SavedRotation.toJson(a),
			fromJson: (obj: any) => SavedRotation.fromJson(obj),
		});

		this.simUI.sim.waitForInit().then(() => {
			savedRotationsManager.loadUserData();
			(this.simUI.individualConfig.presets.rotations || []).forEach(presetRotation => {
				const rotData = presetRotation.rotation;
				// Fill default values so the equality checks always work.
				if (!rotData.rotation) rotData.rotation = APLRotation.create();

				savedRotationsManager.addSavedData({
					name: presetRotation.name,
					tooltip: presetRotation.tooltip,
					isPreset: true,
					data: rotData,
					enableWhen: presetRotation.enableWhen,
					onLoad: presetRotation.onLoad,
				});
			});
		});
	}
}
