import * as PresetUtils from '../../core/preset_utils';
import {
	Consumes,
	Debuffs,
	Faction,
	Flask,
	Food,
	Glyphs,
	IndividualBuffs,
	PetFood,
	Potions,
	Profession,
	RaidBuffs,
	TristateEffect,
} from '../../core/proto/common';
import { SavedTalents } from '../../core/proto/ui';
import {
	AfflictionWarlock_Options as WarlockOptions,
	WarlockMajorGlyph as MajorGlyph,
	WarlockMinorGlyph as MinorGlyph,
	WarlockOptions_Armor as Armor,
	WarlockOptions_Summon as Summon,
	WarlockOptions_WeaponImbue as WeaponImbue,
} from '../../core/proto/warlock';
import AfflictionApl from './apls/affliction.apl.json';
import P1AfflictionGear from './gear_sets/p1_affliction.gear.json';
import P2AfflictionGear from './gear_sets/p2_affliction.gear.json';
import P3AfflictionAllianceGear from './gear_sets/p3_affliction_alliance.gear.json';
import P3AfflictionHordeGear from './gear_sets/p3_affliction_horde.gear.json';
import P4AfflictionGear from './gear_sets/p4_affliction.gear.json';
import PreraidAfflictionGear from './gear_sets/preraid_affliction.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const BIS_TOOLTIP = "This gear preset is inspired from Zephan's Affliction guide: https://www.warcrafttavern.com/wotlk/guides/pve-affliction-warlock/";

export const PRERAID_AFFLICTION_PRESET = PresetUtils.makePresetGear('Preraid Affliction', PreraidAfflictionGear, { tooltip: BIS_TOOLTIP });
export const P1_AFFLICTION_PRESET = PresetUtils.makePresetGear('P1 Affliction', P1AfflictionGear, { tooltip: BIS_TOOLTIP });
export const P2_AFFLICTION_PRESET = PresetUtils.makePresetGear('P2 Affliction', P2AfflictionGear, { tooltip: BIS_TOOLTIP });
export const P3_AFFLICTION_ALLIANCE_PRESET = PresetUtils.makePresetGear('P3 Affliction [A]', P3AfflictionAllianceGear, {
	tooltip: BIS_TOOLTIP,
	faction: Faction.Alliance,
});
export const P3_AFFLICTION_HORDE_PRESET = PresetUtils.makePresetGear('P3 Affliction [H]', P3AfflictionHordeGear, {
	tooltip: BIS_TOOLTIP,
	faction: Faction.Horde,
});
export const P4_AFFLICTION_PRESET = PresetUtils.makePresetGear('P4 Affliction', P4AfflictionGear, { tooltip: BIS_TOOLTIP });

export const APL_Affliction_Default = PresetUtils.makePresetAPLRotation('Affliction', AfflictionApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/cata/talent-calc and copy the numbers in the url.

export const AfflictionTalents = {
	name: 'Affliction',
	data: SavedTalents.create({
		// talentsString: '2350002030023510253500331151--550000051',
		// glyphs: Glyphs.create({
		// 	major1: MajorGlyph.GlyphOfQuickDecay,
		// 	major2: MajorGlyph.GlyphOfLifeTap,
		// 	major3: MajorGlyph.GlyphOfHaunt,
		// 	minor1: MinorGlyph.GlyphOfSouls,
		// 	minor2: MinorGlyph.GlyphOfDrainSoul,
		// 	minor3: MinorGlyph.GlyphOfSubjugateDemon,
		// }),
	}),
};

export const AfflictionOptions = WarlockOptions.create({
	classOptions: {
		armor: Armor.FelArmor,
		summon: Summon.Felhunter,
		weaponImbue: WeaponImbue.GrandSpellstone,
		detonateSeed: true,
	},
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFishFeast,
	petFood: PetFood.PetFoodSpicedMammothTreats,
	defaultPotion: Potions.VolcanicPotion,
	prepopPotion: Potions.VolcanicPotion,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	bloodlust: true,
	markOfTheWild: true,
	icyTalons: true,
	moonkinForm: true,
	leaderOfThePack: true,
	powerWordFortitude: true,
	strengthOfEarthTotem: true,
	trueshotAura: true,
	wrathOfAirTotem: true,
	demonicPact: true,
	blessingOfKings: true,
	blessingOfMight: true,
	communion: true,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
});

export const DestroIndividualBuffs = IndividualBuffs.create({
});

export const DefaultDebuffs = Debuffs.create({
	bloodFrenzy: true,
	sunderArmor: true,
	ebonPlaguebringer: true,
	mangle: true,
	criticalMass: true,
	demoralizingShout: true,
	frostFever: true,
	judgement: true,
});

export const DestroDebuffs = Debuffs.create({
	bloodFrenzy: true,
	sunderArmor: true,
	ebonPlaguebringer: true,
	mangle: true,
	criticalMass: true,
	demoralizingShout: true,
	frostFever: true,
	judgement: true,
});

export const OtherDefaults = {
	distanceFromTarget: 25,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
	channelClipDelay: 150,
};
