import * as InputHelpers from '../core/components/input_helpers.js';
import { Player } from '../core/player.js';
import { WarlockOptions_Armor as Armor, WarlockOptions_Summon as Summon, WarlockOptions_WeaponImbue as WeaponImbue } from '../core/proto/warlock.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { WarlockSpecs } from '../core/proto_utils/utils';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const ArmorInput = <SpecType extends WarlockSpecs>() =>
	InputHelpers.makeClassOptionsEnumIconInput<SpecType, Armor>({
		fieldName: 'armor',
		values: [
			{ value: Armor.NoArmor, tooltip: 'No Armor' },
			{ actionId: ActionId.fromSpellId(47893), value: Armor.FelArmor },
			{ actionId: ActionId.fromSpellId(47889), value: Armor.DemonArmor },
		],
	});

export const WeaponImbueInput = <SpecType extends WarlockSpecs>() =>
	InputHelpers.makeClassOptionsEnumIconInput<SpecType, WeaponImbue>({
		fieldName: 'weaponImbue',
		values: [
			{ value: WeaponImbue.NoWeaponImbue, tooltip: 'No Weapon Stone' },
			{ actionId: ActionId.fromItemId(41174), value: WeaponImbue.GrandFirestone },
			{ actionId: ActionId.fromItemId(41196), value: WeaponImbue.GrandSpellstone },
		],
	});

export const PetInput = <SpecType extends WarlockSpecs>() =>
	InputHelpers.makeClassOptionsEnumIconInput<SpecType, Summon>({
		fieldName: 'summon',
		values: [
			{ value: Summon.NoSummon, tooltip: 'No Pet' },
			{ actionId: ActionId.fromSpellId(688), value: Summon.Imp },
			{ actionId: ActionId.fromSpellId(712), value: Summon.Succubus },
			{ actionId: ActionId.fromSpellId(691), value: Summon.Felhunter },
			{
				actionId: ActionId.fromSpellId(30146),
				value: Summon.Felguard,
				showWhen: (player: Player<SpecType>) => player.getTalents().summonFelguard,
			},
		],
		changeEmitter: (player: Player<SpecType>) => player.changeEmitter,
	});

export const DetonateSeed = <SpecType extends WarlockSpecs>() =>
	InputHelpers.makeClassOptionsBooleanInput<SpecType>({
		fieldName: 'detonateSeed',
		label: 'Detonate Seed on Cast',
		labelTooltip: 'Simulates raid doing damage to targets such that seed detonates immediately on cast.',
	});
