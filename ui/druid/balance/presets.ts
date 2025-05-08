import * as PresetUtils from '../../core/preset_utils.js';
import { ConsumesSpec, Debuffs, Glyphs, IndividualBuffs, PartyBuffs, Profession, RaidBuffs, Stat, UnitReference } from '../../core/proto/common.js';
import { BalanceDruid_Options as BalanceDruidOptions, DruidMajorGlyph } from '../../core/proto/druid.js';
import { SavedTalents } from '../../core/proto/ui.js';
import { Stats } from '../../core/proto_utils/stats';
import T11Apl from './apls/t11.apl.json';
import T12Apl from './apls/t12.apl.json';
import T13Apl from './apls/t13.apl.json';
import PreraidGear from './gear_sets/preraid.gear.json';
import T11Gear from './gear_sets/t11.gear.json';
import T12Gear from './gear_sets/t12.gear.json';
import T13Gear from './gear_sets/t13.gear.json';
import T13ItemSwapGear from './gear_sets/t13_item_swap.gear.json';

export const PreraidPresetGear = PresetUtils.makePresetGear('Pre-raid', PreraidGear);
export const T11PresetGear = PresetUtils.makePresetGear('T11', T11Gear);
export const T12PresetGear = PresetUtils.makePresetGear('T12', T12Gear);
export const T13PresetGear = PresetUtils.makePresetGear('T13', T13Gear);

export const T13PresetItemSwapGear = PresetUtils.makePresetItemSwapGear('T13 - Item Swap', T13ItemSwapGear);

export const T11PresetRotation = PresetUtils.makePresetAPLRotation('T11 4P', T11Apl);
export const T12PresetRotation = PresetUtils.makePresetAPLRotation('T12', T12Apl);
export const T13PresetRotation = PresetUtils.makePresetAPLRotation('T13', T13Apl);

export const StandardEPWeights = PresetUtils.makePresetEpWeights(
	'Standard',
	Stats.fromMap({
		[Stat.StatIntellect]: 1.3,
		[Stat.StatSpirit]: 1.27,
		[Stat.StatSpellPower]: 1,
		[Stat.StatHitRating]: 1.27,
		[Stat.StatCritRating]: 0.41,
		[Stat.StatHasteRating]: 0.8,
		[Stat.StatMasteryRating]: 0.56,
	}),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/mop-classic/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '',
		glyphs: Glyphs.create({
			major2: DruidMajorGlyph.GlyphOfRebirth,
		}),
	}),
};

export const DefaultOptions = BalanceDruidOptions.create({
	classOptions: {
		innervateTarget: UnitReference.create(),
	},
});

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 58086, // Flask of the Draconic Mind
	foodId: 62290, // Seafood Magnifique Feast
	potId: 58091, // Volcanic Potion
	prepotId: 58091, // Volcanic Potion
});
export const DefaultRaidBuffs = RaidBuffs.create({});

export const DefaultIndividualBuffs = IndividualBuffs.create({});

export const DefaultPartyBuffs = PartyBuffs.create({});

export const DefaultDebuffs = Debuffs.create({
	// bloodFrenzy: true,
	// sunderArmor: true,
	// ebonPlaguebringer: true,
	// mangle: true,
	// criticalMass: true,
	// demoralizingShout: true,
	// frostFever: true,
});

export const OtherDefaults = {
	distanceFromTarget: 20,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
};

export const PresetBuildPreraid = PresetUtils.makePresetBuild('Balance Pre-raid', {
	gear: PreraidPresetGear,
	talents: StandardTalents,
	rotation: T13PresetRotation,
	epWeights: StandardEPWeights,
});

export const PresetBuildT11 = PresetUtils.makePresetBuild('Balance T11', {
	gear: T11PresetGear,
	talents: StandardTalents,
	rotation: T11PresetRotation,
	epWeights: StandardEPWeights,
});

export const PresetBuildT12 = PresetUtils.makePresetBuild('Balance T12', {
	gear: T12PresetGear,
	talents: StandardTalents,
	rotation: T12PresetRotation,
	epWeights: StandardEPWeights,
});

export const PresetBuildT13 = PresetUtils.makePresetBuild('Balance T13', {
	gear: T13PresetGear,
	talents: StandardTalents,
	rotation: T13PresetRotation,
	epWeights: StandardEPWeights,
});
