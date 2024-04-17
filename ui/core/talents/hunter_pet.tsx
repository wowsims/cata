// eslint-disable-next-line @typescript-eslint/no-unused-vars
import { element } from 'tsx-vanilla';

import { Component } from '../components/component.js';
import { IconEnumPicker } from '../components/icon_enum_picker';
import * as InputHelpers from '../components/input_helpers.js';
import { SavedDataManager } from '../components/saved_data_manager.js';
import { Player } from '../player.js';
import { Spec } from '../proto/common';
import { HunterOptions_PetType as PetType, HunterPetTalents } from '../proto/hunter.js';
import { ActionId } from '../proto_utils/action_id.js';
import { getTalentTree, getTalentTreePoints, HunterSpecs } from '../proto_utils/utils.js';
import { SimUI } from '../sim_ui.js';
import { EventID, TypedEvent } from '../typed_event.js';
import { protoToTalentString, talentStringToProto } from './factory.js';
import { newTalentsConfig, TalentsConfig, TalentsPicker } from './talents_picker.jsx';
import HunterPetCunningJson from './trees/hunter_cunning.json';
import HunterPetFerocityJson from './trees/hunter_ferocity.json';
import HunterPetTenacityJson from './trees/hunter_tenacity.json';

export function makePetTypeInputConfig<SpecType extends HunterSpecs>(): InputHelpers.TypedIconEnumPickerConfig<Player<SpecType>, PetType> {
	return InputHelpers.makeClassOptionsEnumIconInput<SpecType, PetType>({
		extraCssClasses: ['pet-type-picker'],
		fieldName: 'petType',
		numColumns: 5,
		values: [
			{ value: PetType.PetNone, tooltip: 'No Pet' },
			{ actionId: ActionId.fromPetName('Bat'), tooltip: 'Bat', value: PetType.Bat },
			{ actionId: ActionId.fromPetName('Bear'), tooltip: 'Bear', value: PetType.Bear },
			{ actionId: ActionId.fromPetName('Bird of Prey'), tooltip: 'Bird of Prey', value: PetType.BirdOfPrey },
			{ actionId: ActionId.fromPetName('Boar'), tooltip: 'Boar', value: PetType.Boar },
			{ actionId: ActionId.fromPetName('Carrion Bird'), tooltip: 'Carrion Bird', value: PetType.CarrionBird },
			{ actionId: ActionId.fromPetName('Cat'), tooltip: 'Cat', value: PetType.Cat },
			{ actionId: ActionId.fromPetName('Chimaera'), tooltip: 'Chimaera (Exotic)', value: PetType.Chimaera },
			{ actionId: ActionId.fromPetName('Core Hound'), tooltip: 'Core Hound (Exotic)', value: PetType.CoreHound },
			{ actionId: ActionId.fromPetName('Crab'), tooltip: 'Crab', value: PetType.Crab },
			{ actionId: ActionId.fromPetName('Crocolisk'), tooltip: 'Crocolisk', value: PetType.Crocolisk },
			{ actionId: ActionId.fromPetName('Devilsaur'), tooltip: 'Devilsaur (Exotic)', value: PetType.Devilsaur },
			{ actionId: ActionId.fromPetName('Dragonhawk'), tooltip: 'Dragonhawk', value: PetType.Dragonhawk },
			{ actionId: ActionId.fromPetName('Gorilla'), tooltip: 'Gorilla', value: PetType.Gorilla },
			{ actionId: ActionId.fromPetName('Hyena'), tooltip: 'Hyena', value: PetType.Hyena },
			{ actionId: ActionId.fromPetName('Moth'), tooltip: 'Moth', value: PetType.Moth },
			{ actionId: ActionId.fromPetName('Nether Ray'), tooltip: 'Nether Ray', value: PetType.NetherRay },
			{ actionId: ActionId.fromPetName('Raptor'), tooltip: 'Raptor', value: PetType.Raptor },
			{ actionId: ActionId.fromPetName('Ravager'), tooltip: 'Ravager', value: PetType.Ravager },
			{ actionId: ActionId.fromPetName('Rhino'), tooltip: 'Rhino', value: PetType.Rhino },
			{ actionId: ActionId.fromPetName('Scorpid'), tooltip: 'Scorpid', value: PetType.Scorpid },
			{ actionId: ActionId.fromPetName('Serpent'), tooltip: 'Serpent', value: PetType.Serpent },
			{ actionId: ActionId.fromPetName('Silithid'), tooltip: 'Silithid (Exotic)', value: PetType.Silithid },
			{ actionId: ActionId.fromPetName('Spider'), tooltip: 'Spider', value: PetType.Spider },
			{ actionId: ActionId.fromPetName('Spirit Beast'), tooltip: 'Spirit Beast (Exotic)', value: PetType.SpiritBeast },
			{ actionId: ActionId.fromPetName('Spore Bat'), tooltip: 'Spore Bat', value: PetType.SporeBat },
			{ actionId: ActionId.fromPetName('Tallstrider'), tooltip: 'Tallstrider', value: PetType.Tallstrider },
			{ actionId: ActionId.fromPetName('Turtle'), tooltip: 'Turtle', value: PetType.Turtle },
			{ actionId: ActionId.fromPetName('Warp Stalker'), tooltip: 'Warp Stalker', value: PetType.WarpStalker },
			{ actionId: ActionId.fromPetName('Wasp'), tooltip: 'Wasp', value: PetType.Wasp },
			{ actionId: ActionId.fromPetName('Wind Serpent'), tooltip: 'Wind Serpent', value: PetType.WindSerpent },
			{ actionId: ActionId.fromPetName('Wolf'), tooltip: 'Wolf', value: PetType.Wolf },
			{ actionId: ActionId.fromPetName('Worm'), tooltip: 'Worm (Exotic)', value: PetType.Worm },
		],
	});
}

enum PetCategory {
	Cunning,
	Ferocity,
	Tenacity,
}

const petCategories: Record<PetType, PetCategory> = {
	[PetType.PetNone]: PetCategory.Ferocity,
	[PetType.Bat]: PetCategory.Cunning,
	[PetType.Bear]: PetCategory.Tenacity,
	[PetType.BirdOfPrey]: PetCategory.Cunning,
	[PetType.Boar]: PetCategory.Tenacity,
	[PetType.CarrionBird]: PetCategory.Ferocity,
	[PetType.Cat]: PetCategory.Ferocity,
	[PetType.Chimaera]: PetCategory.Cunning,
	[PetType.CoreHound]: PetCategory.Ferocity,
	[PetType.Crab]: PetCategory.Tenacity,
	[PetType.Crocolisk]: PetCategory.Tenacity,
	[PetType.Devilsaur]: PetCategory.Ferocity,
	[PetType.Dragonhawk]: PetCategory.Cunning,
	[PetType.Gorilla]: PetCategory.Tenacity,
	[PetType.Hyena]: PetCategory.Ferocity,
	[PetType.Moth]: PetCategory.Ferocity,
	[PetType.NetherRay]: PetCategory.Cunning,
	[PetType.Raptor]: PetCategory.Ferocity,
	[PetType.Ravager]: PetCategory.Cunning,
	[PetType.Rhino]: PetCategory.Tenacity,
	[PetType.Scorpid]: PetCategory.Tenacity,
	[PetType.Serpent]: PetCategory.Cunning,
	[PetType.Silithid]: PetCategory.Cunning,
	[PetType.Spider]: PetCategory.Cunning,
	[PetType.SpiritBeast]: PetCategory.Ferocity,
	[PetType.SporeBat]: PetCategory.Cunning,
	[PetType.Tallstrider]: PetCategory.Ferocity,
	[PetType.Turtle]: PetCategory.Tenacity,
	[PetType.WarpStalker]: PetCategory.Tenacity,
	[PetType.Wasp]: PetCategory.Ferocity,
	[PetType.WindSerpent]: PetCategory.Cunning,
	[PetType.Wolf]: PetCategory.Ferocity,
	[PetType.Worm]: PetCategory.Tenacity,
};

const categoryOrder = [PetCategory.Cunning, PetCategory.Ferocity, PetCategory.Tenacity];
const categoryClasses = ['cunning', 'ferocity', 'tenacity'];

export function getPetTalentsConfig(petType: PetType): TalentsConfig<HunterPetTalents> {
	const petCategory = petCategories[petType];
	const categoryIdx = categoryOrder.indexOf(petCategory);
	return petTalentsConfig[categoryIdx];
}

export const cunningDefault: HunterPetTalents = HunterPetTalents.create({
	serpentSwiftness: 2,
	dive: true,
	owlsFocus: 2,
	spikedCollar: 3,
	cullingTheHerd: 3,
	feedingFrenzy: 2,
	roarOfRecovery: true,
	wolverineBite: true,
	wildHunt: 2,
});
export const ferocityDefault: HunterPetTalents = HunterPetTalents.create({
	serpentSwiftness: 2,
	dive: true,
	spikedCollar: 3,
	bloodthirsty: 1,
	cullingTheHerd: 3,
	spidersBite: 3,
	rabid: true,
	callOfTheWild: true,
	sharkAttack: 2,
});
export const tenacityDefault: HunterPetTalents = HunterPetTalents.create({
	serpentSwiftness: 2,
	charge: true,
	spikedCollar: 3,
	boarsSpeed: true,
	cullingTheHerd: 3,
	thunderstomp: true,
	graceOfTheMantis: 2,
	roarOfSacrifice: true,
	intervene: true,
	wildHunt: 2,
});
const defaultTalents = [cunningDefault, ferocityDefault, tenacityDefault];

export const cunningBMDefault: HunterPetTalents = HunterPetTalents.create({
	serpentSwiftness: 2,
	dive: true,
	owlsFocus: 2,
	spikedCollar: 3,
	cullingTheHerd: 3,
	feedingFrenzy: 2,
	roarOfRecovery: true,
	wolverineBite: true,
	wildHunt: 2,
});
export const ferocityBMDefault: HunterPetTalents = HunterPetTalents.create({
	serpentSwiftness: 2,
	dive: true,
	spikedCollar: 3,
	bloodthirsty: 1,
	cullingTheHerd: 3,
	spidersBite: 3,
	rabid: true,
	callOfTheWild: true,
	sharkAttack: 2,
});
export const tenacityBMDefault: HunterPetTalents = HunterPetTalents.create({
	serpentSwiftness: 2,
	charge: true,
	spikedCollar: 3,
	boarsSpeed: true,
	cullingTheHerd: 3,
	thunderstomp: true,
	graceOfTheMantis: 2,
	roarOfSacrifice: true,
	intervene: true,
	wildHunt: 2,
});
const defaultBMTalents = [cunningBMDefault, ferocityBMDefault, tenacityBMDefault];

export const cunningPetTalentsConfig: TalentsConfig<HunterPetTalents> = newTalentsConfig(HunterPetCunningJson);
export const ferocityPetTalentsConfig: TalentsConfig<HunterPetTalents> = newTalentsConfig(HunterPetFerocityJson);
export const tenacityPetTalentsConfig: TalentsConfig<HunterPetTalents> = newTalentsConfig(HunterPetTenacityJson);

export const petTalentsConfig = [cunningPetTalentsConfig, ferocityPetTalentsConfig, tenacityPetTalentsConfig];

export class HunterPet<SpecType extends HunterSpecs> {
	readonly player: Player<SpecType>;

	private talents: HunterPetTalents;
	private talentsConfig: TalentsConfig<HunterPetTalents>;
	private talentsString: string;

	readonly talentsChangeEmitter: TypedEvent<void>;

	constructor(player: Player<SpecType>) {
		this.player = player;
		this.talents = player.getClassOptions().petTalents ?? HunterPetTalents.create();
		this.talentsConfig = getPetTalentsConfig(player.getClassOptions().petType);
		this.talentsString = protoToTalentString(this.talents, this.talentsConfig);
		this.talentsChangeEmitter = this.player.specOptionsChangeEmitter;
	}

	getTalents(): HunterPetTalents {
		return this.talents;
	}

	getTalentsString(): string {
		return protoToTalentString(this.talents, this.talentsConfig);
	}

	setTalentsString(eventID: EventID, newTalentsString: string) {
		if (newTalentsString == this.talentsString) return;

		const options = this.player.getClassOptions();
		options.petTalents = talentStringToProto(HunterPetTalents.create(), newTalentsString, this.talentsConfig);

		this.talents = options.petTalents;
		this.talentsString = newTalentsString;
		this.player.setClassOptions(eventID, options);
	}

	getTalentTree(): number {
		return getTalentTree(this.getTalentsString());
	}

	getTalentTreePoints(): Array<number> {
		return getTalentTreePoints(this.getTalentsString());
	}
}

export class HunterPetTalentsPicker<SpecType extends HunterSpecs> extends Component {
	private readonly simUI: SimUI;
	private readonly player: Player<SpecType>;
	private curCategory: PetCategory | null;
	private curTalents: HunterPetTalents;

	// Not saved to storage, just holds last-used values for this session.
	private savedSets: Array<HunterPetTalents>;

	constructor(parent: HTMLElement, simUI: SimUI, player: Player<SpecType>) {
		super(parent, 'hunter-pet-talents-picker');
		this.simUI = simUI;
		this.player = player;

		this.curCategory = this.getCategoryFromPlayer();
		this.curTalents = this.getPetTalentsFromPlayer();
		this.savedSets = defaultTalents.slice();
		this.savedSets[this.curCategory] = this.curTalents;

		this.rootElem.classList.add(categoryClasses[this.curCategory]);

		const talentsContainer = <div className="pet-talents-container" />;
		this.rootElem.appendChild(talentsContainer);

		simUI.sim.waitForInit().then(() => {
			const pet = new HunterPet(player);
			categoryOrder.map((_category, i) => {
				const pickerContainer = document.createElement('div');
				pickerContainer.classList.add('hunter-pet-talents-' + categoryClasses[i]);
				talentsContainer.appendChild(pickerContainer);

				const talentsConfig = petTalentsConfig[i];
				const picker = new TalentsPicker(pickerContainer, pet, {
					klass: player.getClass(),
					trees: talentsConfig,
					changedEvent: (pet: HunterPet<SpecType>) => pet.player.specOptionsChangeEmitter,
					getValue: (pet: HunterPet<SpecType>) => pet.getTalentsString(),
					setValue: (eventID: EventID, pet: HunterPet<SpecType>, newValue: string) => {
						pet.setTalentsString(eventID, newValue);
						this.savedSets[i] = pet.getTalents();
						this.curTalents = pet.getTalents();
					},
					pointsPerRow: 3,
				});

				return picker;
			});
		});

		new IconEnumPicker(this.rootElem, this.player, makePetTypeInputConfig());

		player.specOptionsChangeEmitter.on(() => {
			const petCategory = this.getCategoryFromPlayer();
			const categoryIdx = categoryOrder.indexOf(petCategory);
			console.log('should change pet talents', petCategory);

			if (petCategory != this.curCategory) {
				this.curCategory = petCategory;
				this.rootElem.classList.remove(...categoryClasses);
				this.rootElem.classList.add(categoryClasses[categoryIdx]);

				const curTalents = this.getPetTalentsFromPlayer();
				console.log(curTalents);
				if (!HunterPetTalents.equals(curTalents, this.curTalents)) {
					// If the current talents have also changed, this was probably a load so we shouldn't switch sets.
					this.curTalents = curTalents;
					this.savedSets[this.curCategory] = this.curTalents;
				} else {
					// Revert to the talents from last time the user was editing this category.
					const options = this.player.getClassOptions();
					options.petTalents = this.savedSets[this.curCategory];
					this.player.setClassOptions(TypedEvent.nextEventID(), options);
					this.curTalents = options.petTalents;
				}
			}
		});
	}

	getPetTalentsFromPlayer(): HunterPetTalents {
		return this.player.getClassOptions().petTalents || HunterPetTalents.create();
	}

	getCategoryFromPlayer(): PetCategory {
		const petType = this.player.getClassOptions().petType;
		return petCategories[petType];
	}
}
