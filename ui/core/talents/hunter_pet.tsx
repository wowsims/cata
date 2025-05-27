import * as InputHelpers from '../components/input_helpers.js';
import { Player } from '../player.js';
import { HunterOptions_PetType as PetType, PetSpec } from '../proto/hunter.js';
import { ActionId } from '../proto_utils/action_id.js';
import { HunterSpecs } from '../proto_utils/utils.js';
import { EventID, TypedEvent } from '../typed_event.js';

export function makePetTypeInputConfig<SpecType extends HunterSpecs>(): InputHelpers.TypedIconEnumPickerConfig<any, PetType> {
	return InputHelpers.makeClassOptionsEnumIconInput<SpecType, PetType>({
		extraCssClasses: ['pet-type-picker'],
		fieldName: 'petType',
		numColumns: 5,
		values: [
			{ value: PetType.PetNone, actionId: ActionId.fromPetName(''), tooltip: 'No Pet' },
			{ value: PetType.Bat, actionId: ActionId.fromPetName('Bat'), tooltip: 'Bat' },
			{ value: PetType.Bear, actionId: ActionId.fromPetName('Bear'), tooltip: 'Bear' },
			{ value: PetType.BirdOfPrey, actionId: ActionId.fromPetName('Bird of Prey'), tooltip: 'Bird of Prey' },
			{ value: PetType.Boar, actionId: ActionId.fromPetName('Boar'), tooltip: 'Boar' },
			{ value: PetType.CarrionBird, actionId: ActionId.fromPetName('Carrion Bird'), tooltip: 'Carrion Bird' },
			{ value: PetType.Cat, actionId: ActionId.fromPetName('Cat'), tooltip: 'Cat' },
			{ value: PetType.Chimaera, actionId: ActionId.fromPetName('Chimaera'), tooltip: 'Chimaera' },
			{ value: PetType.CoreHound, actionId: ActionId.fromPetName('Core Hound'), tooltip: 'Core Hound' },
			{ value: PetType.Crab, actionId: ActionId.fromPetName('Crab'), tooltip: 'Crab' },
			{ value: PetType.Crocolisk, actionId: ActionId.fromPetName('Crocolisk'), tooltip: 'Crocolisk' },
			{ value: PetType.Devilsaur, actionId: ActionId.fromPetName('Devilsaur'), tooltip: 'Devilsaur' },
			{ value: PetType.Dragonhawk, actionId: ActionId.fromPetName('Dragonhawk'), tooltip: 'Dragonhawk' },
			{ value: PetType.Fox, actionId: ActionId.fromPetName('Fox'), tooltip: 'Fox' },
			{ value: PetType.Gorilla, actionId: ActionId.fromPetName('Gorilla'), tooltip: 'Gorilla' },
			{ value: PetType.Hyena, actionId: ActionId.fromPetName('Hyena'), tooltip: 'Hyena' },
			{ value: PetType.Moth, actionId: ActionId.fromPetName('Moth'), tooltip: 'Moth' },
			{ value: PetType.NetherRay, actionId: ActionId.fromPetName('Nether Ray'), tooltip: 'Nether Ray' },
			{ value: PetType.Raptor, actionId: ActionId.fromPetName('Raptor'), tooltip: 'Raptor' },
			{ value: PetType.Ravager, actionId: ActionId.fromPetName('Ravager'), tooltip: 'Ravager' },
			{ value: PetType.Rhino, actionId: ActionId.fromPetName('Rhino'), tooltip: 'Rhino' },
			{ value: PetType.Scorpid, actionId: ActionId.fromPetName('Scorpid'), tooltip: 'Scorpid' },
			{ value: PetType.Serpent, actionId: ActionId.fromPetName('Serpent'), tooltip: 'Serpent' },
			{ value: PetType.Silithid, actionId: ActionId.fromPetName('Silithid'), tooltip: 'Silithid' },
			{ value: PetType.Spider, actionId: ActionId.fromPetName('Spider'), tooltip: 'Spider' },
			{ value: PetType.SpiritBeast, actionId: ActionId.fromPetName('Spirit Beast'), tooltip: 'Spirit Beast' },
			{ value: PetType.SporeBat, actionId: ActionId.fromPetName('Spore Bat'), tooltip: 'Spore Bat' },
			{ value: PetType.Tallstrider, actionId: ActionId.fromPetName('Tallstrider'), tooltip: 'Tallstrider' },
			{ value: PetType.Turtle, actionId: ActionId.fromPetName('Turtle'), tooltip: 'Turtle' },
			{ value: PetType.WarpStalker, actionId: ActionId.fromPetName('Warp Stalker'), tooltip: 'Warp Stalker' },
			{ value: PetType.Wasp, actionId: ActionId.fromPetName('Wasp'), tooltip: 'Wasp' },
			{ value: PetType.WindSerpent, actionId: ActionId.fromPetName('Wind Serpent'), tooltip: 'Wind Serpent' },
			{ value: PetType.Wolf, actionId: ActionId.fromPetName('Wolf'), tooltip: 'Wolf' },
			{ value: PetType.Worm, actionId: ActionId.fromPetName('Worm'), tooltip: 'Worm' },
			{ value: PetType.ShaleSpider, actionId: ActionId.fromPetName('Shale Spider'), tooltip: 'Shale Spider' },
			{ value: PetType.Goat, actionId: ActionId.fromPetName('Goat'), tooltip: 'Goat' },
			{ value: PetType.Porcupine, actionId: ActionId.fromPetName('Porcupine'), tooltip: 'Porcupine' },
			{ value: PetType.Monkey, actionId: ActionId.fromPetName('Monkey'), tooltip: 'Monkey' },
			{ value: PetType.Basilisk, actionId: ActionId.fromPetName('Basilisk'), tooltip: 'Basilisk' },
			{ value: PetType.Crane, actionId: ActionId.fromPetName('Crane'), tooltip: 'Crane' },
			{ value: PetType.Dog, actionId: ActionId.fromPetName('Dog'), tooltip: 'Dog' },
			{ value: PetType.Beetle, actionId: ActionId.fromPetName('Beetle'), tooltip: 'Beetle' },
			{ value: PetType.Quilen, actionId: ActionId.fromPetName('Quilen'), tooltip: 'Quilen' },
			{ value: PetType.WaterStrider, actionId: ActionId.fromPetName('Water Strider'), tooltip: 'Water Strider' },
		],
	});
}

const defaultSpec = PetSpec.Ferocity;

export class HunterPet<SpecType extends HunterSpecs> {
	readonly player: Player<SpecType>;

	private spec: PetSpec;

	readonly specChangeEmitter: TypedEvent<void>;

	constructor(player: Player<SpecType>) {
		this.player = player;
		this.spec = defaultSpec;
		this.specChangeEmitter = this.player.specOptionsChangeEmitter;
	}

	getSpec(): PetSpec {
		return this.spec;
	}

	setSpec(eventID: EventID, newSpec: PetSpec) {
		if (newSpec == this.spec) return;

		const options = this.player.getClassOptions();
		options.petSpec = newSpec;
		this.spec = newSpec;
		this.player.setClassOptions(eventID, options);
	}
}
