import { IndividualSimUI } from '../../individual_sim_ui';
import { Player } from '../../player';
import { EquipmentSpec, UnitStats } from '../../proto/common';
import { SavedGearSet } from '../../proto/ui';
import { Stats } from '../../proto_utils/stats';
import { EventID, TypedEvent } from '../../typed_event';
import GearPicker from '../gear_picker/gear_picker';
import { SavedDataManager } from '../saved_data_manager';
import { SimTab } from '../sim_tab';
import { GemSummary } from './gem_summary';
import { ReforgeSummary } from './reforge_summary';

export class GearTab extends SimTab {
	protected simUI: IndividualSimUI<any>;

	readonly leftPanel: HTMLElement;
	readonly rightPanel: HTMLElement;

	constructor(parentElem: HTMLElement, simUI: IndividualSimUI<any>) {
		super(parentElem, simUI, { identifier: 'gear-tab', title: 'Gear' });
		this.simUI = simUI;

		this.leftPanel = document.createElement('div');
		this.leftPanel.classList.add('gear-tab-left', 'tab-panel-left');

		this.rightPanel = document.createElement('div');
		this.rightPanel.classList.add('gear-tab-right', 'tab-panel-right');

		this.contentContainer.appendChild(this.leftPanel);
		this.contentContainer.appendChild(this.rightPanel);

		this.buildTabContent();
	}

	protected buildTabContent() {
		this.buildGearPickers();
		this.buildSummaryTablesContainer();
		this.buildSavedGearsetPicker();
	}

	private buildSummaryTablesContainer() {
		const container = document.createElement('div');
		container.classList.add('summary-tables-container');
		this.leftPanel.appendChild(container);

		new GemSummary(container, this.simUI, this.simUI.player);
		new ReforgeSummary(container, this.simUI, this.simUI.player);
	}

	private buildGearPickers() {
		new GearPicker(this.leftPanel, this.simUI, this.simUI.player);
	}

	private buildSavedGearsetPicker() {
		const savedGearManager = new SavedDataManager<Player<any>, SavedGearSet>(this.rightPanel, this.simUI.player, {
			header: { title: 'Gear Sets' },
			label: 'Gear Set',
			storageKey: this.simUI.getSavedGearStorageKey(),
			getData: (player: Player<any>) => {
				return SavedGearSet.create({
					gear: player.getGear().asSpec(),
					bonusStatsStats: player.getBonusStats().toProto(),
				});
			},
			setData: (eventID: EventID, player: Player<any>, newSavedGear: SavedGearSet) => {
				TypedEvent.freezeAllAndDo(() => {
					player.setGear(eventID, this.simUI.sim.db.lookupEquipmentSpec(newSavedGear.gear || EquipmentSpec.create()));
					player.setBonusStats(eventID, Stats.fromProto(newSavedGear.bonusStatsStats || UnitStats.create()));
				});
			},
			changeEmitters: [this.simUI.player.changeEmitter],
			equals: (a: SavedGearSet, b: SavedGearSet) => SavedGearSet.equals(a, b),
			toJson: (a: SavedGearSet) => SavedGearSet.toJson(a),
			fromJson: (obj: any) => SavedGearSet.fromJson(obj),
		});

		this.simUI.sim.waitForInit().then(() => {
			savedGearManager.loadUserData();
			this.simUI.individualConfig.presets.gear.forEach(presetGear => {
				savedGearManager.addSavedData({
					name: presetGear.name,
					tooltip: presetGear.tooltip,
					isPreset: true,
					data: SavedGearSet.create({
						// Convert to gear and back so order is always the same.
						gear: this.simUI.sim.db.lookupEquipmentSpec(presetGear.gear).asSpec(),
						bonusStatsStats: new Stats().toProto(),
					}),
					enableWhen: presetGear.enableWhen,
					onLoad: presetGear.onLoad,
				});
			});
		});
	}
}
