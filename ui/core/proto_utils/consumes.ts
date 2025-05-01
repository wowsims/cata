import { BattleElixir, Conjured, ConsumableType, Consumes, ConsumesSpec, Flask, Food, GuardianElixir, Potions, TinkerHands } from '../proto/common';
import { findInputItemForEnum } from '../utils';
import { Database } from './database';

export function convertConsumesToSpec(consumes: Consumes | undefined, db: Database, existingSpec?: ConsumesSpec): ConsumesSpec {
	const spec = existingSpec ?? ConsumesSpec.create();
	if (!consumes) return spec;

	const enumMappings: Array<{
		consumesKey: keyof Consumes;
		specKey: keyof ConsumesSpec;
		enumType: any;
		unknownValue: any;
		cType: ConsumableType;
	}> = [
		{
			consumesKey: 'prepopPotion',
			specKey: 'prepotId',
			enumType: Potions,
			unknownValue: Potions.UnknownPotion,
			cType: ConsumableType.ConsumableTypePotion,
		},
		{ consumesKey: 'defaultPotion', specKey: 'potId', enumType: Potions, unknownValue: Potions.UnknownPotion, cType: ConsumableType.ConsumableTypePotion },
		{ consumesKey: 'flask', specKey: 'flaskId', enumType: Flask, unknownValue: Flask.FlaskUnknown, cType: ConsumableType.ConsumableTypeFlask },
		{ consumesKey: 'food', specKey: 'foodId', enumType: Food, unknownValue: Food.FoodUnknown, cType: ConsumableType.ConsumableTypeFood },
		{
			consumesKey: 'guardianElixir',
			specKey: 'guardianElixirId',
			enumType: GuardianElixir,
			unknownValue: GuardianElixir.GuardianElixirUnknown,
			cType: ConsumableType.ConsumableTypeGuardianElixir,
		},
		{
			consumesKey: 'battleElixir',
			specKey: 'battleElixirId',
			enumType: BattleElixir,
			unknownValue: BattleElixir.BattleElixirUnknown,
			cType: ConsumableType.ConsumableTypeBattleElixir,
		},
	];

	for (const { consumesKey, specKey, enumType, unknownValue, cType } of enumMappings) {
		const v = consumes[consumesKey] as any;
		if (v !== unknownValue && (spec[specKey] as number) === 0) {
			spec[specKey] = findInputItemForEnum(enumType, v, db.getConsumablesByType(cType))?.id ?? 0;
		}
	}

	if (consumes.highpoweredBoltGun && spec.explosiveId === 0) spec.explosiveId = 82207;
	if (consumes.explosiveBigDaddy && spec.explosiveId === 0) spec.explosiveId = 89637;

	const conjuredMap: Record<Conjured, number> = {
		[Conjured.ConjuredDarkRune]: 20520,
		[Conjured.ConjuredHealthstone]: 5512,
		[Conjured.ConjuredRogueThistleTea]: 7676,
		[Conjured.ConjuredUnknown]: 0,
	};
	if (consumes.defaultConjured !== Conjured.ConjuredUnknown && spec.conjuredId === 0) {
		spec.conjuredId = conjuredMap[consumes.defaultConjured] ?? 0;
	}

	const tinkerMap: Record<TinkerHands, number> = {
		[TinkerHands.TinkerHandsSynapseSprings]: 82174,
		[TinkerHands.TinkerHandsTazikShocker]: 82180,
		[TinkerHands.TinkerHandsQuickflipDeflectionPlates]: 82177,
		[TinkerHands.TinkerHandsSpinalHealingInjector]: 82184,
		[TinkerHands.TinkerHandsZ50ManaGulper]: 82186,
		[TinkerHands.TinkerHandsNone]: 0,
	};
	if (consumes.tinkerHands !== TinkerHands.TinkerHandsNone && spec.tinkerId === 0) {
		spec.tinkerId = tinkerMap[consumes.tinkerHands] ?? 0;
	}

	return spec;
}
