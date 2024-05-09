import { Player } from '../../player';
import { BattleElixir, Class, Conjured, Consumes, Explosive, Flask, Food, GuardianElixir, Potions, Profession, Spec, Stat, TinkerHands } from '../../proto/common';
import { ActionId } from '../../proto_utils/action_id';
import { EventID, TypedEvent } from '../../typed_event';
import { IconEnumValueConfig } from '../icon_enum_picker';
import { makeBooleanConsumeInput } from '../icon_inputs';
import * as InputHelpers from '../input_helpers';
import { ActionInputConfig, ItemStatOption } from './stat_options';

export interface ConsumableInputConfig<T> extends ActionInputConfig<T> {
	value: T;
}

export interface ConsumableStatOption<T> extends ItemStatOption<T> {
	config: ConsumableInputConfig<T>;
}

export interface ConsumeInputFactoryArgs<T extends number> {
	consumesFieldName: keyof Consumes;
	// Additional callback if logic besides syncing consumes is required
	onSet?: (eventactionId: EventID, player: Player<any>, newValue: T) => void;
	showWhen?: (player: Player<any>) => boolean;
}

function makeConsumeInputFactory<T extends number, SpecType extends Spec>(
	args: ConsumeInputFactoryArgs<T>,
): (options: ConsumableStatOption<T>[], tooltip?: string) => InputHelpers.TypedIconEnumPickerConfig<Player<SpecType>, T> {
	return (options: ConsumableStatOption<T>[], tooltip?: string) => {
		return {
			type: 'iconEnum',
			tooltip: tooltip,
			numColumns: options.length > 5 ? 2 : 1,
			values: [{ value: 0 } as unknown as IconEnumValueConfig<Player<SpecType>, T>].concat(
				options.map(option => {
					const rtn = {
						actionId: option.config.actionId,
						value: option.config.value,
						showWhen: (player: Player<SpecType>) =>
							(!option.config.showWhen || option.config.showWhen(player)) &&
							(option.config.faction || player.getFaction()) == player.getFaction(),
					} as IconEnumValueConfig<Player<SpecType>, T>;

					return rtn;
				}),
			),
			equals: (a: T, b: T) => a == b,
			zeroValue: 0 as T,
			changedEvent: (player: Player<any>) => TypedEvent.onAny([player.consumesChangeEmitter, player.gearChangeEmitter, player.professionChangeEmitter]),
			showWhen: (player: Player<any>) => !args.showWhen || args.showWhen(player),
			getValue: (player: Player<any>) => player.getConsumes()[args.consumesFieldName] as T,
			setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
				const newConsumes = player.getConsumes();

				if (newConsumes[args.consumesFieldName] === newValue) {
					return;
				}

				(newConsumes[args.consumesFieldName] as number) = newValue;
				TypedEvent.freezeAllAndDo(() => {
					player.setConsumes(eventID, newConsumes);
					if (args.onSet) {
						args.onSet(eventID, player, newValue as T);
					}
				});
			},
		};
	};
}

///////////////////////////////////////////////////////////////////////////
//                                 CONJURED
///////////////////////////////////////////////////////////////////////////

export const ConjuredDarkRune = {
	actionId: ActionId.fromItemId(12662),
	value: Conjured.ConjuredDarkRune,
};
export const ConjuredFlameCap = {
	actionId: ActionId.fromItemId(22788),
	value: Conjured.ConjuredFlameCap,
};
export const ConjuredHealthstone = {
	actionId: ActionId.fromItemId(22105),
	value: Conjured.ConjuredHealthstone,
};
export const ConjuredRogueThistleTea = {
	actionId: ActionId.fromItemId(7676),
	value: Conjured.ConjuredRogueThistleTea,
	showWhen: <SpecType extends Spec>(player: Player<SpecType>) => player.getClass() == Class.ClassRogue,
};

export const CONJURED_CONFIG = [
	{ config: ConjuredRogueThistleTea, stats: [] },
	{ config: ConjuredHealthstone, stats: [Stat.StatStamina] },
	{ config: ConjuredDarkRune, stats: [Stat.StatIntellect] },
	{ config: ConjuredFlameCap, stats: [] },
] as ConsumableStatOption<Conjured>[];

export const makeConjuredInput = makeConsumeInputFactory({ consumesFieldName: 'defaultConjured' });

///////////////////////////////////////////////////////////////////////////
//                                 EXPLOSIVES
///////////////////////////////////////////////////////////////////////////

export const ExplosiveSaroniteBomb = {
	actionId: ActionId.fromItemId(41119),
	value: Explosive.ExplosiveSaroniteBomb,
};
export const ExplosiveCobaltFragBomb = {
	actionId: ActionId.fromItemId(40771),
	value: Explosive.ExplosiveCobaltFragBomb,
};

export const EXPLOSIVES_CONFIG = [
	{ config: ExplosiveSaroniteBomb, stats: [] },
	{ config: ExplosiveCobaltFragBomb, stats: [] },
] as ConsumableStatOption<Explosive>[];

export const makeExplosivesInput = makeConsumeInputFactory({
	consumesFieldName: 'fillerExplosive',
	showWhen: (player: Player<any>) => player.hasProfession(Profession.Engineering),
});

export const ThermalSapper = makeBooleanConsumeInput({
	actionId: ActionId.fromItemId(42641),
	fieldName: 'thermalSapper',
	showWhen: (player: Player<any>) => player.hasProfession(Profession.Engineering),
});
export const ExplosiveDecoy = makeBooleanConsumeInput({
	actionId: ActionId.fromItemId(40536),
	fieldName: 'explosiveDecoy',
	showWhen: (player: Player<any>) => player.hasProfession(Profession.Engineering),
});

///////////////////////////////////////////////////////////////////////////
//                                 Tinkers
///////////////////////////////////////////////////////////////////////////

export const TinkerHandsSynapseSprings = {
	actionId: ActionId.fromSpellId(82174),
	value: TinkerHands.TinkerHandsSynapseSprings,
};
export const TinkerHandsQuickflipDeflectionPlates = {
	actionId: ActionId.fromSpellId(82176),
	value: TinkerHands.TinkerHandsQuickflipDeflectionPlates,
};
export const TinkerHandsTazikShocker = {
	actionId: ActionId.fromSpellId(82179),
	value: TinkerHands.TinkerHandsTazikShocker,
};
export const TinkerHandsSpinalHealingInjector = {
	actionId: ActionId.fromSpellId(82184),
	value: TinkerHands.TinkerHandsSpinalHealingInjector,
};
export const TinkerHandsZ50ManaGulper = {
	actionId: ActionId.fromSpellId(82186),
	value: TinkerHands.TinkerHandsZ50ManaGulper,
};

export const TINKERS_HANDS_CONFIG = [
	{ config: TinkerHandsSynapseSprings, stats: [] },
	{ config: TinkerHandsQuickflipDeflectionPlates, stats: [] },
	{ config: TinkerHandsTazikShocker, stats: [] },
	{ config: TinkerHandsSpinalHealingInjector, stats: [] },
	{ config: TinkerHandsZ50ManaGulper, stats: [] },
] as ConsumableStatOption<TinkerHands>[];

export const makeTinkerHandsInput = makeConsumeInputFactory({
	consumesFieldName: 'tinkerHands',
	showWhen: (player: Player<any>) => player.hasProfession(Profession.Engineering),
});

///////////////////////////////////////////////////////////////////////////
//                                 FLASKS + ELIXIRS
///////////////////////////////////////////////////////////////////////////

// Flasks
export const FlaskOfTitanicStrength = {
	actionId: ActionId.fromItemId(58088), // Use the correct item ID
	value: Flask.FlaskOfTitanicStrength,
};

export const FlaskOfTheWinds = {
	actionId: ActionId.fromItemId(58087), // Use the correct item ID
	value: Flask.FlaskOfTheWinds,
};

export const FlaskOfSteelskin = {
	actionId: ActionId.fromItemId(58085), // Use the correct item ID
	value: Flask.FlaskOfSteelskin,
};

export const FlaskOfFlowingWater = {
	actionId: ActionId.fromItemId(67438), // Use the correct item ID
	value: Flask.FlaskOfFlowingWater,
};

export const FlaskOfTheDraconicMind = {
	actionId: ActionId.fromItemId(58086), // Use the correct item ID
	value: Flask.FlaskOfTheDraconicMind,
};

export const FlaskOfTheFrostWyrm = {
	actionId: ActionId.fromItemId(46376),
	value: Flask.FlaskOfTheFrostWyrm,
};
export const FlaskOfEndlessRage = {
	actionId: ActionId.fromItemId(46377),
	value: Flask.FlaskOfEndlessRage,
};
export const FlaskOfPureMojo = {
	actionId: ActionId.fromItemId(46378),
	value: Flask.FlaskOfPureMojo,
};
export const FlaskOfStoneblood = {
	actionId: ActionId.fromItemId(46379),
	value: Flask.FlaskOfStoneblood,
};
export const LesserFlaskOfToughness = {
	actionId: ActionId.fromItemId(40079),
	value: Flask.LesserFlaskOfToughness,
};
export const LesserFlaskOfResistance = {
	actionId: ActionId.fromItemId(44939),
	value: Flask.LesserFlaskOfResistance,
};

export const FLASKS_CONFIG = [
	{ config: FlaskOfTheDraconicMind, stats: [Stat.StatSpellPower] },
	{ config: FlaskOfTitanicStrength, stats: [Stat.StatStrength] },
	{ config: FlaskOfTheWinds, stats: [Stat.StatAgility] },
	{ config: FlaskOfSteelskin, stats: [Stat.StatStamina] },
	{ config: FlaskOfFlowingWater, stats: [Stat.StatSpirit] },
	{ config: FlaskOfTheFrostWyrm, stats: [Stat.StatSpellPower] },
	{ config: FlaskOfEndlessRage, stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower] },
	{ config: FlaskOfPureMojo, stats: [Stat.StatMP5] },
	{ config: FlaskOfStoneblood, stats: [Stat.StatStamina] },
	{ config: LesserFlaskOfToughness, stats: [Stat.StatResilience] },
	{
		config: LesserFlaskOfResistance,
		stats: [Stat.StatArcaneResistance, Stat.StatFireResistance, Stat.StatFrostResistance, Stat.StatNatureResistance, Stat.StatShadowResistance],
	},
] as ConsumableStatOption<Flask>[];

export const makeFlasksInput = makeConsumeInputFactory({
	consumesFieldName: 'flask',
	onSet: (eventID: EventID, player: Player<any>, newValue: Flask) => {
		if (newValue) {
			const newConsumes = player.getConsumes();
			newConsumes.battleElixir = BattleElixir.BattleElixirUnknown;
			newConsumes.guardianElixir = GuardianElixir.GuardianElixirUnknown;
			player.setConsumes(eventID, newConsumes);
		}
	},
});

// Battle Elixirs
export const ElixirOfTheMaster = {
	actionId: ActionId.fromItemId(58148),
	value: BattleElixir.ElixirOfTheMaster,
};

export const ElixirOfMightySpeed = {
	actionId: ActionId.fromItemId(58144),
	value: BattleElixir.ElixirOfMightySpeed,
};

export const ElixirOfImpossibleAccuracy = {
	actionId: ActionId.fromItemId(58094),
	value: BattleElixir.ElixirOfImpossibleAccuracy,
};

export const ElixirOfTheCobra = {
	actionId: ActionId.fromItemId(58092),
	value: BattleElixir.ElixirOfTheCobra,
};

export const ElixirOfTheNaga = {
	actionId: ActionId.fromItemId(58089),
	value: BattleElixir.ElixirOfTheNaga,
};

export const GhostElixir = {
	actionId: ActionId.fromItemId(58084),
	value: BattleElixir.GhostElixir,
};
export const ElixirOfAccuracy = {
	actionId: ActionId.fromItemId(44325),
	value: BattleElixir.ElixirOfAccuracy,
};
export const ElixirOfArmorPiercing = {
	actionId: ActionId.fromItemId(44330),
	value: BattleElixir.ElixirOfArmorPiercing,
};
export const ElixirOfDeadlyStrikes = {
	actionId: ActionId.fromItemId(44327),
	value: BattleElixir.ElixirOfDeadlyStrikes,
};
export const ElixirOfExpertise = {
	actionId: ActionId.fromItemId(44329),
	value: BattleElixir.ElixirOfExpertise,
};
export const ElixirOfLightningSpeed = {
	actionId: ActionId.fromItemId(44331),
	value: BattleElixir.ElixirOfLightningSpeed,
};
export const ElixirOfMightyAgility = {
	actionId: ActionId.fromItemId(39666),
	value: BattleElixir.ElixirOfMightyAgility,
};
export const ElixirOfMightyStrength = {
	actionId: ActionId.fromItemId(40073),
	value: BattleElixir.ElixirOfMightyStrength,
};
export const GurusElixir = {
	actionId: ActionId.fromItemId(40076),
	value: BattleElixir.GurusElixir,
};
export const SpellpowerElixir = {
	actionId: ActionId.fromItemId(40070),
	value: BattleElixir.SpellpowerElixir,
};
export const WrathElixir = {
	actionId: ActionId.fromItemId(40068),
	value: BattleElixir.WrathElixir,
};
export const ElixirOfDemonslaying = {
	actionId: ActionId.fromItemId(9224),
	value: BattleElixir.ElixirOfDemonslaying,
};

export const BATTLE_ELIXIRS_CONFIG = [
	{ config: ElixirOfTheMaster, stats: [Stat.StatMastery] },
	{ config: ElixirOfMightySpeed, stats: [Stat.StatMeleeHaste, Stat.StatSpellHaste] },
	{ config: ElixirOfImpossibleAccuracy, stats: [Stat.StatMeleeHit, Stat.StatSpellHit] },
	{ config: ElixirOfTheCobra, stats: [Stat.StatMeleeCrit, Stat.StatSpellCrit] },
	{ config: ElixirOfTheNaga, stats: [Stat.StatExpertise] },
	{ config: GhostElixir, stats: [Stat.StatSpirit] },
	{ config: ElixirOfAccuracy, stats: [Stat.StatMeleeHit, Stat.StatSpellHit] },
	{ config: ElixirOfDeadlyStrikes, stats: [Stat.StatMeleeCrit, Stat.StatSpellCrit] },
	{ config: ElixirOfExpertise, stats: [Stat.StatExpertise] },
	{ config: ElixirOfLightningSpeed, stats: [Stat.StatMeleeHaste, Stat.StatSpellHaste] },
	{ config: ElixirOfMightyAgility, stats: [Stat.StatAgility] },
	{ config: ElixirOfMightyStrength, stats: [Stat.StatStrength] },
	{
		config: GurusElixir,
		stats: [Stat.StatStamina, Stat.StatAgility, Stat.StatStrength, Stat.StatSpirit, Stat.StatIntellect],
	},
	{ config: SpellpowerElixir, stats: [Stat.StatSpellPower] },
	{ config: WrathElixir, stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower] },
] as ConsumableStatOption<BattleElixir>[];

export const makeBattleElixirsInput = makeConsumeInputFactory({
	consumesFieldName: 'battleElixir',
	onSet: (eventID: EventID, player: Player<any>, newValue: BattleElixir) => {
		if (newValue) {
			const newConsumes = player.getConsumes();
			newConsumes.flask = Flask.FlaskUnknown;
			player.setConsumes(eventID, newConsumes);
		}
	},
});

// Guardian Elixirs
export const ElixirOfDeepEarth = {
	actionId: ActionId.fromItemId(58093),
	value: GuardianElixir.ElixirOfDeepEarth,
};
export const PrismaticElixir = {
	actionId: ActionId.fromItemId(58143),
	value: GuardianElixir.PrismaticElixir,
};
export const ElixirOfMightyDefense = {
	actionId: ActionId.fromItemId(44328),
	value: GuardianElixir.ElixirOfMightyDefense,
};
export const ElixirOfMightyFortitude = {
	actionId: ActionId.fromItemId(40078),
	value: GuardianElixir.ElixirOfMightyFortitude,
};
export const ElixirOfMightyMageblood = {
	actionId: ActionId.fromItemId(40109),
	value: GuardianElixir.ElixirOfMightyMageblood,
};
export const ElixirOfMightyThoughts = {
	actionId: ActionId.fromItemId(44332),
	value: GuardianElixir.ElixirOfMightyThoughts,
};
export const ElixirOfProtection = {
	actionId: ActionId.fromItemId(40097),
	value: GuardianElixir.ElixirOfProtection,
};
export const ElixirOfSpirit = {
	actionId: ActionId.fromItemId(40072),
	value: GuardianElixir.ElixirOfSpirit,
};

export const GUARDIAN_ELIXIRS_CONFIG = [
	{ config: ElixirOfDeepEarth, stats: [Stat.StatArmor] },
	{
		config: PrismaticElixir,
		stats: [Stat.StatArcaneResistance, Stat.StatFireResistance, Stat.StatFrostResistance, Stat.StatNatureResistance, Stat.StatShadowResistance],
	},
	{ config: ElixirOfMightyDefense, stats: [Stat.StatDefense] },
	{ config: ElixirOfMightyFortitude, stats: [Stat.StatStamina] },
	{ config: ElixirOfMightyMageblood, stats: [Stat.StatMP5] },
	{ config: ElixirOfMightyThoughts, stats: [Stat.StatIntellect] },
	{ config: ElixirOfProtection, stats: [Stat.StatArmor] },
	{ config: ElixirOfSpirit, stats: [Stat.StatSpirit] },
] as ConsumableStatOption<GuardianElixir>[];

export const makeGuardianElixirsInput = makeConsumeInputFactory({
	consumesFieldName: 'guardianElixir',
	onSet: (eventID: EventID, player: Player<any>, newValue: GuardianElixir) => {
		if (newValue) {
			const newConsumes = player.getConsumes();
			newConsumes.flask = Flask.FlaskUnknown;
			player.setConsumes(eventID, newConsumes);
		}
	},
});

///////////////////////////////////////////////////////////////////////////
//                                 FOOD
///////////////////////////////////////////////////////////////////////////

export const FoodFishFeast = { actionId: ActionId.fromItemId(43015), value: Food.FoodFishFeast };
export const FoodGreatFeast = { actionId: ActionId.fromItemId(34753), value: Food.FoodGreatFeast };
export const FoodBlackenedDragonfin = {
	actionId: ActionId.fromItemId(42999),
	value: Food.FoodBlackenedDragonfin,
};
export const FoodHeartyRhino = {
	actionId: ActionId.fromItemId(42995),
	value: Food.FoodHeartyRhino,
};
export const FoodMegaMammothMeal = {
	actionId: ActionId.fromItemId(34754),
	value: Food.FoodMegaMammothMeal,
};
export const FoodSpicedWormBurger = {
	actionId: ActionId.fromItemId(34756),
	value: Food.FoodSpicedWormBurger,
};
export const FoodRhinoliciousWormsteak = {
	actionId: ActionId.fromItemId(42994),
	value: Food.FoodRhinoliciousWormsteak,
};
export const FoodImperialMantaSteak = {
	actionId: ActionId.fromItemId(34769),
	value: Food.FoodImperialMantaSteak,
};
export const FoodSnapperExtreme = {
	actionId: ActionId.fromItemId(42996),
	value: Food.FoodSnapperExtreme,
};
export const FoodMightyRhinoDogs = {
	actionId: ActionId.fromItemId(34758),
	value: Food.FoodMightyRhinoDogs,
};
export const FoodFirecrackerSalmon = {
	actionId: ActionId.fromItemId(34767),
	value: Food.FoodFirecrackerSalmon,
};
export const FoodCuttlesteak = {
	actionId: ActionId.fromItemId(42998),
	value: Food.FoodCuttlesteak,
};
export const FoodDragonfinFilet = {
	actionId: ActionId.fromItemId(43000),
	value: Food.FoodDragonfinFilet,
};

export const FoodBlackenedBasilisk = {
	actionId: ActionId.fromItemId(27657),
	value: Food.FoodBlackenedBasilisk,
};
export const FoodGrilledMudfish = {
	actionId: ActionId.fromItemId(27664),
	value: Food.FoodGrilledMudfish,
};
export const FoodRavagerDog = { actionId: ActionId.fromItemId(27655), value: Food.FoodRavagerDog };
export const FoodRoastedClefthoof = {
	actionId: ActionId.fromItemId(27658),
	value: Food.FoodRoastedClefthoof,
};
export const FoodSpicyHotTalbuk = {
	actionId: ActionId.fromItemId(33872),
	value: Food.FoodSpicyHotTalbuk,
};
export const FoodSkullfishSoup = {
	actionId: ActionId.fromItemId(33825),
	value: Food.FoodSkullfishSoup,
};
export const FoodFishermansFeast = {
	actionId: ActionId.fromItemId(33052),
	value: Food.FoodFishermansFeast,
};
export const FoodSeafoodMagnifiqueFeast = {
	actionId: ActionId.fromItemId(62290),
	value: Food.FoodSeafoodFeast,
};
export const FoodFortuneCookie = {
	actionId: ActionId.fromItemId(62649),
	value: Food.FoodFortuneCookie,
};
export const FoodSeveredSagefishHead = {
	actionId: ActionId.fromItemId(62671),
	value: Food.FoodSeveredSagefish,
};
export const FoodBeerBastedCrocolisk = {
	actionId: ActionId.fromItemId(62670),
	value: Food.FoodBeerBasedCrocolisk,
};
export const FoodBakedRockfish = {
	actionId: ActionId.fromItemId(62661),
	value: Food.FoodBakedRockfish,
};
export const FoodBasiliskLiverdog = {
	actionId: ActionId.fromItemId(62665),
	value: Food.FoodBasiliskLiverdog,
};
export const FoodBlackbellySushi = {
	actionId: ActionId.fromItemId(62668),
	value: Food.FoodBlackbellySushi,
};
export const FoodSkeweredEll = {
	actionId: ActionId.fromItemId(62669),
	value: Food.FoodSkeweredEel,
};
export const FoodCrocoliskAuGratin = {
	actionId: ActionId.fromItemId(62664),
	value: Food.FoodCrocoliskAuGratin,
};
export const FoodDeliciousSagefishTail = {
	actionId: ActionId.fromItemId(62666),
	value: Food.FoodDeliciousSagefishTail,
};
export const FoodMushroomSauceMudfish = {
	actionId: ActionId.fromItemId(62667),
	value: Food.FoodMushroomSauceMudfish,
};
export const FoodGrilledDragon = {
	actionId: ActionId.fromItemId(62662),
	value: Food.FoodGrilledDragon,
};
export const FoodLavascaleMinestrone = {
	actionId: ActionId.fromItemId(62663),
	value: Food.FoodLavascaleMinestrone,
};

export const FOOD_CONFIG = [
	{
		config: FoodSeafoodMagnifiqueFeast,
		stats: [Stat.StatStamina, Stat.StatStrength, Stat.StatIntellect, Stat.StatAgility],
	},
	{ config: FoodFortuneCookie, stats: [Stat.StatAgility, Stat.StatStamina, Stat.StatStrength, Stat.StatAgility] },
	{ config: FoodSeveredSagefishHead, stats: [Stat.StatIntellect] },
	{ config: FoodBeerBastedCrocolisk, stats: [Stat.StatStrength] },
	{ config: FoodSkeweredEll, stats: [Stat.StatAgility] },
	{ config: FoodDeliciousSagefishTail, stats: [Stat.StatSpirit] },
	{ config: FoodBakedRockfish, stats: [Stat.StatMeleeCrit, Stat.StatSpellCrit] },
	{ config: FoodBasiliskLiverdog, stats: [Stat.StatMeleeHaste, Stat.StatSpellHaste] },
	{ config: FoodLavascaleMinestrone, stats: [Stat.StatMastery] },
	{ config: FoodGrilledDragon, stats: [Stat.StatMeleeHit, Stat.StatSpellHit] },
	{ config: FoodCrocoliskAuGratin, stats: [Stat.StatExpertise] },
	{ config: FoodMushroomSauceMudfish, stats: [Stat.StatDodge] },
	{ config: FoodBlackbellySushi, stats: [Stat.StatParry] },
	{
		config: FoodFishFeast,
		stats: [Stat.StatStamina, Stat.StatAttackPower, Stat.StatRangedAttackPower, Stat.StatSpellPower],
	},
	{
		config: FoodGreatFeast,
		stats: [Stat.StatStamina, Stat.StatAttackPower, Stat.StatRangedAttackPower, Stat.StatSpellPower],
	},
	{ config: FoodBlackenedDragonfin, stats: [Stat.StatAgility] },
	{ config: FoodDragonfinFilet, stats: [Stat.StatStrength] },
	{ config: FoodCuttlesteak, stats: [Stat.StatSpirit] },
	{ config: FoodMegaMammothMeal, stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower] },
	{ config: FoodHeartyRhino, stats: [Stat.StatMeleeCrit, Stat.StatSpellCrit] },
	{ config: FoodRhinoliciousWormsteak, stats: [Stat.StatExpertise] },
	{ config: FoodFirecrackerSalmon, stats: [Stat.StatSpellPower] },
	{ config: FoodSnapperExtreme, stats: [Stat.StatMeleeHit, Stat.StatSpellHit] },
	{ config: FoodSpicedWormBurger, stats: [Stat.StatMeleeCrit, Stat.StatSpellCrit] },
	{ config: FoodImperialMantaSteak, stats: [Stat.StatMeleeHaste, Stat.StatSpellHaste] },
	{ config: FoodMightyRhinoDogs, stats: [Stat.StatMP5] },
] as ConsumableStatOption<Food>[];

export const makeFoodInput = makeConsumeInputFactory({ consumesFieldName: 'food' });

///////////////////////////////////////////////////////////////////////////
//                                 PET
///////////////////////////////////////////////////////////////////////////

export const PetScrollOfAgilityV = makeBooleanConsumeInput({
	actionId: ActionId.fromItemId(27498),
	fieldName: 'petScrollOfAgility',
	value: 5,
});
export const PetScrollOfStrengthV = makeBooleanConsumeInput({
	actionId: ActionId.fromItemId(27503),
	fieldName: 'petScrollOfStrength',
	value: 5,
});

///////////////////////////////////////////////////////////////////////////
//                                 POTIONS
///////////////////////////////////////////////////////////////////////////
export const GolembloodPotion = {
	actionId: ActionId.fromItemId(58146),
	value: Potions.GolembloodPotion,
};

export const PotionOfTheTolvir = {
	actionId: ActionId.fromItemId(58145),
	value: Potions.PotionOfTheTolvir,
};

export const PotionOfConcentration = {
	actionId: ActionId.fromItemId(57194),
	value: Potions.PotionOfConcentration,
};

export const VolcanicPotion = {
	actionId: ActionId.fromItemId(58091),
	value: Potions.VolcanicPotion,
};

export const EarthenPotion = {
	actionId: ActionId.fromItemId(58090),
	value: Potions.EarthenPotion,
};

export const MightyRejuvenationPotion = {
	actionId: ActionId.fromItemId(57193),
	value: Potions.MightyRejuvenationPotion,
};

export const MythicalHealingPotion = {
	actionId: ActionId.fromItemId(57191),
	value: Potions.MythicalHealingPotion,
};

export const MythicalManaPotion = {
	actionId: ActionId.fromItemId(57192),
	value: Potions.MythicalManaPotion,
};
export const PotionOfSpeed = { actionId: ActionId.fromItemId(40211), value: Potions.PotionOfSpeed };
export const HastePotion = { actionId: ActionId.fromItemId(22838), value: Potions.HastePotion };
export const MightyRagePotion = {
	actionId: ActionId.fromItemId(13442),
	value: Potions.MightyRagePotion,
};

export const POTIONS_CONFIG = [
	{ config: GolembloodPotion, stats: [Stat.StatStrength] },
	{ config: PotionOfTheTolvir, stats: [Stat.StatAgility] },
	{ config: PotionOfConcentration, stats: [Stat.StatMana] },
	{ config: VolcanicPotion, stats: [Stat.StatIntellect] },
	{ config: EarthenPotion, stats: [Stat.StatArmor] },
	{ config: MightyRejuvenationPotion, stats: [Stat.StatIntellect, Stat.StatHealth] },
	{ config: MythicalHealingPotion, stats: [Stat.StatHealth] },
	{ config: MythicalManaPotion, stats: [Stat.StatIntellect] },
	{ config: PotionOfSpeed, stats: [Stat.StatMeleeHaste, Stat.StatSpellHaste] },
] as ConsumableStatOption<Potions>[];

export const PRE_POTIONS_CONFIG = [
	{ config: GolembloodPotion, stats: [Stat.StatStrength] },
	{ config: PotionOfTheTolvir, stats: [Stat.StatAgility] },
	{ config: VolcanicPotion, stats: [Stat.StatIntellect] },
	{ config: EarthenPotion, stats: [Stat.StatArmor] },
	{ config: PotionOfSpeed, stats: [Stat.StatMeleeHaste, Stat.StatSpellHaste] },
] as ConsumableStatOption<Potions>[];

export const makePotionsInput = makeConsumeInputFactory({ consumesFieldName: 'defaultPotion' });
export const makePrepopPotionsInput = makeConsumeInputFactory({
	consumesFieldName: 'prepopPotion',
});
