import tippy from 'tippy.js';
import { ref } from 'tsx-vanilla';

import { IndividualSimUI } from '../../individual_sim_ui';
import { PresetBuild } from '../../preset_utils';
import { APLRotation, APLRotation_Type } from '../../proto/apl';
import { Encounter, EquipmentSpec, HealingModel, Spec } from '../../proto/common';
import { SavedTalents } from '../../proto/ui';
import { TypedEvent } from '../../typed_event';
import { Component } from '../component';
import { ContentBlock } from '../content_block';

type PresetConfigurationCategory = 'gear' | 'talents' | 'rotation' | 'encounter' | 'race';

export class PresetConfigurationPicker extends Component {
	readonly simUI: IndividualSimUI<Spec>;
	readonly builds: Array<PresetBuild>;

	constructor(parentElem: HTMLElement, simUI: IndividualSimUI<Spec>, types?: PresetConfigurationCategory[]) {
		super(parentElem, 'preset-configuration-picker-root');
		this.rootElem.classList.add('saved-data-manager-root');

		this.simUI = simUI;
		this.builds = (this.simUI.individualConfig.presets.builds ?? []).filter(build =>
			Object.keys(build).some(category => types?.includes(category as PresetConfigurationCategory) && !!build[category as PresetConfigurationCategory]),
		);

		if (!this.builds.length) {
			this.rootElem.classList.add('hide');
			return;
		}

		const contentBlock = new ContentBlock(this.rootElem, 'saved-data', {
			header: {
				title: 'Preset Configurations',
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
					this.simUI.player.changeEmitter,
					this.simUI.sim.settingsChangeEmitter,
					this.simUI.sim.raid.changeEmitter,
					this.simUI.sim.encounter.changeEmitter,
				]).on(checkActive);
			});
			contentBlock.bodyElement.replaceChildren(container);
		});
	}

	private applyBuild({ gear, rotation, rotationType, talents, epWeights, encounter, race }: PresetBuild) {
		const eventID = TypedEvent.nextEventID();
		TypedEvent.freezeAllAndDo(() => {
			if (gear) this.simUI.player.setGear(eventID, this.simUI.sim.db.lookupEquipmentSpec(gear.gear));
			if (race) this.simUI.player.setRace(eventID, race);
			if (talents) {
				this.simUI.player.setTalentsString(eventID, talents.data.talentsString);
				if (talents.data.glyphs) this.simUI.player.setGlyphs(eventID, talents.data.glyphs);
			}
			if (rotationType) {
				this.simUI.player.aplRotation.type = rotationType;
				this.simUI.player.rotationChangeEmitter.emit(eventID);
			} else  if (rotation?.rotation.rotation) {
				this.simUI.player.setAplRotation(eventID, rotation.rotation.rotation);
			}
			if (epWeights) this.simUI.player.setEpWeights(eventID, epWeights.epWeights);
			if (encounter) {
				if (encounter.encounter) this.simUI.sim.encounter.fromProto(eventID, encounter.encounter);
				if (encounter.healingModel) this.simUI.player.setHealingModel(eventID, encounter.healingModel);
				if (encounter.tanks) this.simUI.sim.raid.setTanks(eventID, encounter.tanks);
				if (encounter.buffs) this.simUI.player.setBuffs(eventID, encounter.buffs);
				if (encounter.debuffs) this.simUI.sim.raid.setDebuffs(eventID, encounter.debuffs);
				if (encounter.raidBuffs) this.simUI.sim.raid.setBuffs(eventID, encounter.raidBuffs);
				if (encounter.consumes) this.simUI.player.setConsumes(eventID, encounter.consumes);
			}
		});
	}

	private isBuildActive({ gear, rotation, rotationType, talents, epWeights, encounter, race }: PresetBuild): boolean {
		const hasGear = gear ? EquipmentSpec.equals(gear.gear, this.simUI.player.getGear().asSpec()) : true;
		const hasRace = typeof race === 'number' ? race === this.simUI.player.getRace() : true;
		const hasTalents = talents
			? SavedTalents.equals(
					talents.data,
					SavedTalents.create({
						talentsString: this.simUI.player.getTalentsString(),
						glyphs: this.simUI.player.getGlyphs(),
					}),
			  )
			: true;
		let hasRotation = true;
		if (rotationType) {
			hasRotation = rotationType === this.simUI.player.getRotationType();
		} else if (rotation) {
			const activeRotation = this.simUI.player.getResolvedAplRotation();
			// Ensure that the auto rotation can be matched with a preset
			if (activeRotation.type === APLRotation_Type.TypeAuto) activeRotation.type = APLRotation_Type.TypeAPL;
			if (rotation.rotation?.rotation?.type === APLRotation_Type.TypeSimple && rotation.rotation.rotation?.simple?.specRotationJson) {
				hasRotation = this.simUI.player.specTypeFunctions.rotationEquals(
					this.simUI.player.specTypeFunctions.rotationFromJson(JSON.parse(rotation.rotation.rotation.simple.specRotationJson)),
					this.simUI.player.getSimpleRotation(),
				);
			} else {
				hasRotation = APLRotation.equals(rotation.rotation.rotation, activeRotation);
			}
		}
		const hasEpWeights = epWeights ? this.simUI.player.getEpWeights().equals(epWeights.epWeights) : true;
		const hasEncounter = encounter?.encounter ? Encounter.equals(encounter.encounter, this.simUI.sim.encounter.toProto()) : true;
		const hasHealingModel = encounter?.healingModel ? HealingModel.equals(encounter.healingModel, this.simUI.player.getHealingModel()) : true;

		return hasGear && hasRace && hasTalents && hasRotation && hasEpWeights && hasEncounter && hasHealingModel;
	}
}
