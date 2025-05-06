import * as PresetUtils from '../../core/preset_utils';
import { ConsumesSpec, Debuffs, Glyphs, Profession, RaidBuffs, Stat, UnitReference } from '../../core/proto/common';
import { ArcaneMage_Options as MageOptions, MageMajorGlyph as MajorGlyph } from '../../core/proto/mage';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import ArcaneApl from './apls/arcane.apl.json';
import P1ArcaneBisGear from './gear_sets/p1.gear.json';
import P3ArcaneBisGear from './gear_sets/p3.gear.json';
import P4ArcaneBisGear from './gear_sets/p4.gear.json';
import P3ArcanePrebisGear from './gear_sets/prebis.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
export const ARCANE_P1_PRESET = PresetUtils.makePresetGear('P1 Preset', P1ArcaneBisGear);
export const ARCANE_P3_PREBIS_PRESET = PresetUtils.makePresetGear('Pre-raid ', P3ArcanePrebisGear);
export const ARCANE_P3_PRESET = PresetUtils.makePresetGear('P3 Preset', P3ArcaneBisGear);
export const ARCANE_P4_PRESET = PresetUtils.makePresetGear('P4', P4ArcaneBisGear);
/* export const DefaultSimpleRotation = MageRotation.create({
	only3ArcaneBlastStacksBelowManaPercent: 0.15,
	blastWithoutMissileBarrageAboveManaPercent: 0.2,
	missileBarrageBelowManaPercent: 0,
	useArcaneBarrage: false,
}); */

//export const ROTATION_PRESET_SIMPLE = PresetUtils.makePresetSimpleRotation('Simple Default', Spec.SpecArcaneMage, DefaultSimpleRotation);
export const ARCANE_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Arcane', ArcaneApl);
//export const ARCANE_ROTATION_PRESET_AOE = PresetUtils.makePresetAPLRotation('Arcane AOE', ArcaneAoeApl);

// Preset options for EP weights
export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Default',
	Stats.fromMap({
		[Stat.StatIntellect]: 1.8,
		[Stat.StatSpellPower]: 1,
		[Stat.StatHitRating]: 1.52,
		[Stat.StatCritRating]: 0.65,
		[Stat.StatHasteRating]: 0.7,
		[Stat.StatMasteryRating]: 0.67,
	}),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/mop-classic/talent-calc and copy the numbers in the url.
export const ArcaneTalents = {
	name: 'Default',
	data: SavedTalents.create({
		talentsString: '',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfEvocation,
			major2: MajorGlyph.GlyphOfArcanePower,
			major3: MajorGlyph.GlyphOfManaGem,
		}),
	}),
};

export const DefaultArcaneOptions = MageOptions.create({
	classOptions: {},
	focusMagicPercentUptime: 90,
	focusMagicTarget: UnitReference.create(),
});
export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 58086, // Flask of the Draconic Mind
	foodId: 62290, // Seafood Magnifique Feast
	potId: 58091, // Volcanic Potion
	prepotId: 58091, // Volcanic Potion
	tinkerId: 82174, // Synapse Springs
});
export const DefaultRaidBuffs = RaidBuffs.create({});

export const DefaultDebuffs = Debuffs.create({
	// ebonPlaguebringer: true,
	// shadowAndFlame: true,
});

export const OtherDefaults = {
	distanceFromTarget: 20,
	profession1: Profession.Engineering,
	profession2: Profession.Blacksmithing,
};
