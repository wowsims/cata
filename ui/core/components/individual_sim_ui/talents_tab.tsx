import { IndividualSimUI } from '../../individual_sim_ui';
import { Player } from '../../player';
import { Class, Glyphs, Spec } from '../../proto/common';
import { SavedTalents } from '../../proto/ui';
import { classTalentsConfig } from '../../talents/factory';
import { TalentsPicker } from '../../talents/talents_picker';
import { EventID, TypedEvent } from '../../typed_event';
import { PetSpecPicker } from '../pickers/pet_spec_picker';
import { SavedDataManager } from '../saved_data_manager';
import { SimTab } from '../sim_tab';
import { PresetConfigurationPicker } from './preset_configuration_picker';

export class TalentsTab<SpecType extends Spec> extends SimTab {
	protected simUI: IndividualSimUI<any>;

	readonly leftPanel: HTMLElement;
	readonly rightPanel: HTMLElement;

	constructor(parentElem: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parentElem, simUI, { identifier: 'talents-tab', title: 'Talents' });
		this.simUI = simUI;

		this.leftPanel = (<div className="talents-tab-left tab-panel-left" />) as HTMLElement;
		this.rightPanel = (<div className="talents-tab-right tab-panel-right within-raid-sim-hide" />) as HTMLElement;

		this.contentContainer.appendChild(this.leftPanel);
		this.contentContainer.appendChild(this.rightPanel);

		this.buildTabContent();
	}

	protected buildTabContent() {
		this.buildTalentsPicker(this.leftPanel);

		this.buildPresetConfigurationPicker();
		this.buildSavedTalentsPicker();

		this.buildHunterPetPicker(this.leftPanel);
	}
	private buildHunterPetPicker(parentElem: HTMLElement) {
		if (this.simUI.player.isClass(Class.ClassHunter)) {
			new PetSpecPicker(parentElem, this.simUI.player);
		}
	}
	private buildTalentsPicker(parentElem: HTMLElement) {
		new TalentsPicker(parentElem, this.simUI.player, {
			playerClass: this.simUI.player.getClass(),
			playerSpec: this.simUI.player.getSpec(),
			tree: classTalentsConfig[this.simUI.player.getClass()]!,
			changedEvent: (player: Player<any>) => player.talentsChangeEmitter,
			getValue: (player: Player<any>) => player.getTalentsString(),
			setValue: (eventID: EventID, player: Player<any>, newValue: string) => {
				player.setTalentsString(eventID, newValue);
			},
		});
	}

	private buildPresetConfigurationPicker() {
		new PresetConfigurationPicker(this.rightPanel, this.simUI, ['talents']);
	}

	private buildSavedTalentsPicker() {
		const savedTalentsManager = new SavedDataManager<Player<any>, SavedTalents>(this.rightPanel, this.simUI.player, {
			label: 'Talents',
			header: { title: 'Saved Talents' },
			storageKey: this.simUI.getSavedTalentsStorageKey(),
			getData: (player: Player<any>) =>
				SavedTalents.create({
					talentsString: player.getTalentsString(),
					glyphs: player.getGlyphs(),
				}),
			setData: (eventID: EventID, player: Player<any>, newTalents: SavedTalents) => {
				TypedEvent.freezeAllAndDo(() => {
					player.setTalentsString(eventID, newTalents.talentsString);
					player.setGlyphs(eventID, newTalents.glyphs || Glyphs.create());
				});
			},
			changeEmitters: [this.simUI.player.talentsChangeEmitter, this.simUI.player.glyphsChangeEmitter],
			equals: (a: SavedTalents, b: SavedTalents) => SavedTalents.equals(a, b),
			toJson: (a: SavedTalents) => SavedTalents.toJson(a),
			fromJson: (obj: any) => SavedTalents.fromJson(obj),
		});

		this.simUI.sim.waitForInit().then(() => {
			savedTalentsManager.loadUserData();
			this.simUI.individualConfig.presets.talents.forEach(config => {
				config.isPreset = true;
				savedTalentsManager.addSavedData({
					name: config.name,
					isPreset: true,
					data: config.data,
					onLoad: config.onLoad,
				});
			});
		});
	}
}
