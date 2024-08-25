import tippy from 'tippy.js';
import { ref } from 'tsx-vanilla';

import { IndividualSimUI } from '../../individual_sim_ui';
import { PresetBuild } from '../../preset_utils';
import { APLRotation, APLRotation_Type } from '../../proto/apl';
import { Encounter, EquipmentSpec, HealingModel, Spec } from '../../proto/common';
import { TypedEvent } from '../../typed_event';
import { Component } from '../component';
import { ContentBlock } from '../content_block';

type PresetConfigurationCategory = 'gear' | 'talents' | 'rotation' | 'encounter';

export class PresetConfigurationPicker extends Component {
	readonly simUI: IndividualSimUI<Spec>;
	readonly builds: Array<PresetBuild>;

	constructor(parentElem: HTMLElement, simUI: IndividualSimUI<Spec>, type?: PresetConfigurationCategory) {
		super(parentElem, 'preset-configuration-picker-root');
		this.rootElem.classList.add('saved-data-manager-root');

		this.simUI = simUI;
		this.builds = (this.simUI.individualConfig.presets.builds ?? []).filter(build =>
			Object.keys(build).some(category => category === type && !!build[category]),
		);

		if (!this.builds.length) {
			this.rootElem.classList.add('hide');
			return;
		}

		const contentBlock = new ContentBlock(this.rootElem, 'saved-data', {
			header: {
				title: 'Preset configurations',
				tooltip: 'Preset configurations can apply an optimal combination of gear, talents, rotation and encounter settings.',
			},
		});

		const buildsContainerRef = ref<HTMLDivElement>();

		const container = (
			<div className="saved-data-container">
				<div className="saved-data-presets" ref={buildsContainerRef}></div>
			</div>
		);

		this.simUI.sim.waitForInit().then(() => {
			this.builds.forEach(build => {
				const dataElemRef = ref<HTMLButtonElement>();
				buildsContainerRef.value!.appendChild(
					<button className="saved-data-set-chip badge rounded-pill" ref={dataElemRef}>
						<span className="saved-data-set-name" attributes={{ role: 'button' }} onclick={() => this.applyBuild(build)}>
							{build.name}
						</span>
					</button>,
				);

				tippy(dataElemRef.value!, {
					content: (
						<>
							<p className="mb-1">This preset affects the following settings:</p>
							<ul className="mb-0 text-capitalize">
								{Object.keys(build)
									.filter(c => c !== 'name')
									.map(category => (build[category as PresetConfigurationCategory] ? <li>{category}</li> : undefined))}
							</ul>
						</>
					),
				});

				const checkActive = () => dataElemRef.value!.classList[this.isBuildActive(build) ? 'add' : 'remove']('active');

				checkActive();
				TypedEvent.onAny([
					this.simUI.player.gearChangeEmitter,
					this.simUI.player.talentsChangeEmitter,
					this.simUI.player.rotationChangeEmitter,
					this.simUI.player.epWeightsChangeEmitter,
					this.simUI.player.healingModelChangeEmitter,
					this.simUI.sim.settingsChangeEmitter,
					this.simUI.sim.encounter.changeEmitter,
				]).on(checkActive);
			});
			contentBlock.bodyElement.replaceChildren(container);
		});
	}

	private applyBuild({ gear, rotation, talents, epWeights, encounter }: PresetBuild) {
		const eventID = TypedEvent.nextEventID();
		TypedEvent.freezeAllAndDo(() => {
			if (gear) this.simUI.player.setGear(eventID, this.simUI.sim.db.lookupEquipmentSpec(gear.gear));
			if (talents) this.simUI.player.setTalentsString(eventID, talents.data.talentsString);
			if (rotation?.rotation.rotation) {
				const type = rotation.rotation.rotation?.type;
				switch (type) {
					case APLRotation_Type.TypeAPL:
						this.simUI.player.setAplRotation(eventID, rotation.rotation.rotation);
						break;
					case APLRotation_Type.TypeSimple:
						this.simUI.player.setSimpleRotation(eventID, rotation.rotation.rotation);
						break;
					default:
						break;
				}
			}
			if (epWeights) this.simUI.player.setEpWeights(eventID, epWeights.epWeights);
			if (encounter) {
				if (encounter.encounter) this.simUI.sim.encounter.fromProto(eventID, encounter.encounter);
				if (encounter.healingModel) this.simUI.player.setHealingModel(eventID, encounter.healingModel);
			}
		});
	}

	private isBuildActive({ gear, rotation, talents, epWeights, encounter }: PresetBuild): boolean {
		const hasGear = gear ? EquipmentSpec.equals(gear.gear, this.simUI.player.getGear().asSpec()) : true;
		const hasTalents = talents ? talents.data.talentsString == this.simUI.player.getTalentsString() : true;
		const hasRotation = rotation ? APLRotation.equals(rotation.rotation.rotation, this.simUI.player.getResolvedAplRotation()) : true;
		const hasEpWeights = epWeights ? this.simUI.player.getEpWeights().equals(epWeights.epWeights) : true;
		const hasEncounter = encounter?.encounter ? Encounter.equals(encounter.encounter, this.simUI.sim.encounter.toProto()) : true;
		const hasHealingModel = encounter?.healingModel ? HealingModel.equals(encounter.healingModel, this.simUI.player.getHealingModel()) : true;

		return hasGear && hasTalents && hasRotation && hasEpWeights && hasEncounter && hasHealingModel;
	}
}
