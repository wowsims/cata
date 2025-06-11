import * as Mechanics from '../../core/constants/mechanics';
import * as PresetUtils from '../../core/preset_utils.js';
import { ConsumesSpec, Debuffs, Glyphs, Profession, PseudoStat, RaidBuffs, Stat } from '../../core/proto/common.js';
import {
	AirTotem,
	CallTotem,
	EarthTotem,
	EnhancementShaman_Options as EnhancementShamanOptions,
	FireTotem,
	ShamanImbue,
	ShamanMajorGlyph,
	ShamanMinorGlyph,
	ShamanShield,
	ShamanSyncType,
	ShamanTotems,
	TotemSet,
	WaterTotem,
} from '../../core/proto/shaman.js';
import { SavedTalents } from '../../core/proto/ui.js';
import { Stats } from '../../core/proto_utils/stats';
import DefaultApl from './apls/default.apl.json';
import P1NonOrcGear from './gear_sets/p1.non-orc.gear.json';
import P1OrcGear from './gear_sets/p1.orc.gear.json';
import P3NonOrcGear from './gear_sets/p3.non-orc.gear.json';
import P3OrcGear from './gear_sets/p3.orc.gear.json';
import P4OrcGear from './gear_sets/p4.orc.gear.json';
import ItemSwapP4 from './gear_sets/p4_item_swap.gear.json';
import PreraidGear from './gear_sets/preraid.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_PRESET = PresetUtils.makePresetGear('Pre-raid', PreraidGear);

export const P1_ORC_PRESET = PresetUtils.makePresetGear('P1 - Orc', P1OrcGear);
export const P1_NON_ORC_PRESET = PresetUtils.makePresetGear('P1 - Non-Orc', P1NonOrcGear);
export const P3_ORC_PRESET = PresetUtils.makePresetGear('P3 - Orc', P3OrcGear);
export const P3_NON_ORC_PRESET = PresetUtils.makePresetGear('P3 - Non-Orc', P3NonOrcGear);
export const P4_ORC_PRESET = PresetUtils.makePresetGear('P4', P4OrcGear);

export const P4_ITEM_SWAP = PresetUtils.makePresetItemSwapGear('P4 - Item Swap', ItemSwapP4);

export const ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Default', DefaultApl);

// Preset options for EP weights
export const P3_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Default',
	Stats.fromMap(
		{
			[Stat.StatIntellect]: 0.07,
			[Stat.StatAgility]: 2.47,
			[Stat.StatSpellPower]: 0,
			[Stat.StatHitRating]: 0.89 + 0.6,
			[Stat.StatCritRating]: 0.26 + 0.58,
			[Stat.StatHasteRating]: 0.22 + 0.44,
			[Stat.StatAttackPower]: 1.0,
			[Stat.StatExpertiseRating]: 1.3,
			[Stat.StatMasteryRating]: 1.21,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 3.05,
			[PseudoStat.PseudoStatOffHandDps]: 2.56,
			[PseudoStat.PseudoStatSpellHitPercent]: 0.89 * Mechanics.SPELL_HIT_RATING_PER_HIT_PERCENT,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 0.6 * Mechanics.PHYSICAL_HIT_RATING_PER_HIT_PERCENT,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/mop-classic/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '',
		glyphs: Glyphs.create({
			major1: ShamanMajorGlyph.GlyphOfLightningShield,
			major2: ShamanMajorGlyph.GlyphOfHealingStreamTotem,
			major3: ShamanMajorGlyph.GlyphOfFireNova,
		}),
	}),
};

export const DefaultOptions = EnhancementShamanOptions.create({
	classOptions: {
		shield: ShamanShield.LightningShield,
		call: CallTotem.Elements,
		totems: ShamanTotems.create({
			elements: TotemSet.create({
				earth: EarthTotem.StrengthOfEarthTotem,
				air: AirTotem.WrathOfAirTotem,
				fire: FireTotem.SearingTotem,
				water: WaterTotem.ManaSpringTotem,
			}),
			ancestors: TotemSet.create({
				earth: EarthTotem.StrengthOfEarthTotem,
				air: AirTotem.WrathOfAirTotem,
				fire: FireTotem.SearingTotem,
				water: WaterTotem.ManaSpringTotem,
			}),
			spirits: TotemSet.create({
				earth: EarthTotem.StrengthOfEarthTotem,
				air: AirTotem.WrathOfAirTotem,
				fire: FireTotem.SearingTotem,
				water: WaterTotem.ManaSpringTotem,
			}),
			earth: EarthTotem.StrengthOfEarthTotem,
			air: AirTotem.WrathOfAirTotem,
			fire: FireTotem.SearingTotem,
			water: WaterTotem.ManaSpringTotem,
		}),
		imbueMh: ShamanImbue.WindfuryWeapon,
	},
	imbueOh: ShamanImbue.FlametongueWeapon,
	syncType: ShamanSyncType.Auto,
});

export const OtherDefaults = {
	distanceFromTarget: 5,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
};

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 58087, // Flask of the Winds
	foodId: 62290, // Seafood Magnifique Feast
	potId: 58145, // Potion of the Tol'vir
	prepotId: 58145, // Potion of the Tol'vir
});

export const DefaultRaidBuffs = RaidBuffs.create({});

export const DefaultDebuffs = Debuffs.create({
	// bloodFrenzy: true,
	// faerieFire: true,
	// ebonPlaguebringer: true,
	// mangle: true,
	// criticalMass: true,
	// demoralizingShout: true,
	// frostFever: true,
});
