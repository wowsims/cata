import * as PresetUtils from '../../core/preset_utils';
import { Conjured, Consumes, Faction, Flask, Food, Glyphs, Potions, Profession, Spec, UnitReference } from '../../core/proto/common';
import {
	ArcaneMage_Options as MageOptions,
	ArcaneMage_Rotation as MageRotation,
	MageMajorGlyph,
	MageMinorGlyph,
	MageOptions_ArmorType as ArmorType,
} from '../../core/proto/mage';
import { SavedTalents } from '../../core/proto/ui';
import ArcaneApl from './apls/arcane.apl.json';
import ArcaneAoeApl from './apls/arcane_aoe.apl.json';
import P1ArcaneGear from './gear_sets/p1_arcane.gear.json';
import P2ArcaneGear from './gear_sets/p2_arcane.gear.json';
import P3ArcaneAllianceGear from './gear_sets/p3_arcane_alliance.gear.json';
import P3ArcaneHordeGear from './gear_sets/p3_arcane_horde.gear.json';
import P4ArcaneAllianceGear from './gear_sets/p4_arcane_alliance.gear.json';
import P4ArcaneHordeGear from './gear_sets/p4_arcane_horde.gear.json';
import PreraidArcaneGear from './gear_sets/preraid_arcane.gear.json';
// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
export const ARCANE_PRERAID_PRESET = PresetUtils.makePresetGear('Arcane Preraid Preset', PreraidArcaneGear, { talentTree: 0 });
export const ARCANE_P1_PRESET = PresetUtils.makePresetGear('Arcane P1 Preset', P1ArcaneGear, { talentTree: 0 });
export const ARCANE_P2_PRESET = PresetUtils.makePresetGear('Arcane P2 Preset', P2ArcaneGear, { talentTree: 0 });
export const ARCANE_P3_PRESET_ALLIANCE = PresetUtils.makePresetGear('Arcane P3 Preset [A]', P3ArcaneAllianceGear, { talentTree: 0, faction: Faction.Alliance });
export const ARCANE_P3_PRESET_HORDE = PresetUtils.makePresetGear('Arcane P3 Preset [H]', P3ArcaneHordeGear, { talentTree: 0, faction: Faction.Horde });
export const ARCANE_P4_PRESET_ALLIANCE = PresetUtils.makePresetGear('Arcane P4 Preset [A]', P4ArcaneAllianceGear, { talentTree: 0, faction: Faction.Alliance });
export const ARCANE_P4_PRESET_HORDE = PresetUtils.makePresetGear('Arcane P4 Preset [H]', P4ArcaneHordeGear, { talentTree: 0, faction: Faction.Horde });

export const DefaultSimpleRotation = MageRotation.create({
	only3ArcaneBlastStacksBelowManaPercent: 0.15,
	blastWithoutMissileBarrageAboveManaPercent: 0.2,
	missileBarrageBelowManaPercent: 0,
	useArcaneBarrage: false,
});

export const ROTATION_PRESET_SIMPLE = PresetUtils.makePresetSimpleRotation('Simple Default', Spec.SpecArcaneMage, DefaultSimpleRotation);
export const ARCANE_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Arcane', ArcaneApl, { talentTree: 0 });
export const ARCANE_ROTATION_PRESET_AOE = PresetUtils.makePresetAPLRotation('Arcane AOE', ArcaneAoeApl, { talentTree: 0 });

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/cata/talent-calc and copy the numbers in the url.
export const ArcaneTalents = {
	name: 'Arcane',
	data: SavedTalents.create({
		talentsString: '23000513310033015032310250532-03-023303001',
		glyphs: Glyphs.create({
			major1: MageMajorGlyph.GlyphOfArcaneBlast,
			major2: MageMajorGlyph.GlyphOfArcaneMissiles,
			major3: MageMajorGlyph.GlyphOfMoltenArmor,
			minor1: MageMinorGlyph.GlyphOfSlowFall,
			minor2: MageMinorGlyph.GlyphOfFrostWard,
			minor3: MageMinorGlyph.GlyphOfBlastWave,
		}),
	}),
};

export const DefaultArcaneOptions = MageOptions.create({
	classOptions: {
		armor: ArmorType.MoltenArmor,
	},
	focusMagicPercentUptime: 99,
	focusMagicTarget: UnitReference.create(),
});

export const DefaultArcaneConsumes = Consumes.create({
	defaultPotion: Potions.PotionOfSpeed,
	defaultConjured: Conjured.ConjuredDarkRune,
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFirecrackerSalmon,
});

export const OtherDefaults = {
	distanceFromTarget: 20,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
};
