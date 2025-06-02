import { ShamanMajorGlyph, ShamanMinorGlyph, ShamanTalents } from '../proto/shaman.js';
import { GlyphsConfig } from './glyphs_picker.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import ShamanTalentJson from './trees/shaman.json';export const shamanTalentsConfig: TalentsConfig<ShamanTalents> = newTalentsConfig(ShamanTalentJson);

export const shamanGlyphsConfig: GlyphsConfig = {
	majorGlyphs: {
		[ShamanMajorGlyph.GlyphOfUnstableEarth]: {
			name: "Glyph of Unstable Earth",
			description: "Causes your Earthquake spell to also reduce the movement speed of affected targets by 40% for 3s.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shaman_earthquake.jpg",
		},
		[ShamanMajorGlyph.GlyphOfChainLightning]: {
			name: "Glyph of Chain Lightning",
			description: "Your Chain Lightning spell now strikes 2 additional targets, but deals 10% less damage.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_chainlightning.jpg",
		},
		[ShamanMajorGlyph.GlyphOfSpiritWalk]: {
			name: "Glyph of Spirit Walk",
			description: "Reduces the cooldown of your Spirit Walk ability by 25%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_tracking.jpg",
		},
		[ShamanMajorGlyph.GlyphOfCapacitorTotem]: {
			name: "Glyph of Capacitor Totem",
			description: "Reduces the charging time of your Capacitor Totem by 2 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_brilliance.jpg",
		},
		[ShamanMajorGlyph.GlyphOfPurge]: {
			name: "Glyph of Purge",
			description: "Your Purge dispels 1 additional Magic effect but has a 6 sec cooldown.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_purge.jpg",
		},
		[ShamanMajorGlyph.GlyphOfFireElementalTotem]: {
			name: "Glyph of Fire Elemental Totem",
			description: "Reduces the cooldown and duration of your Fire Elemental Totem by 50%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_fire_elemental_totem.jpg",
		},
		[ShamanMajorGlyph.GlyphOfFireNova]: {
			name: "Glyph of Fire Nova",
			description: "Increases the radius of your Fire Nova spell by 5 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shaman_firenova.jpg",
		},
		[ShamanMajorGlyph.GlyphOfFlameShock]: {
			name: "Glyph of Flame Shock",
			description: "When your Flame Shock deals damage, it heals you for 30% of the damage dealt.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_fire_flameshock.jpg",
		},
		[ShamanMajorGlyph.GlyphOfWindShear]: {
			name: "Glyph of Wind Shear",
			description: "Increases the school lockout duration of Wind Shear by 1 sec, but also increases the cooldown by 3 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_cyclone.jpg",
		},
		[ShamanMajorGlyph.GlyphOfHealingStreamTotem]: {
			name: "Glyph of Healing Stream Totem",
			description: "When your Healing Stream Totem heals an ally, it also reduces their Fire, Frost, and Nature damage taken by 10% for 6s.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_spear_04.jpg",
		},
		[ShamanMajorGlyph.GlyphOfHealingWave]: {
			name: "Glyph of Healing Wave",
			description: "Your Healing Wave also heals you for 20% of the healing effect when you heal someone else.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_healingwavegreater.jpg",
		},
		[ShamanMajorGlyph.GlyphOfTotemicRecall]: {
			name: "Glyph of Totemic Recall",
			description: "Causes your Totemic Recall ability to return the full mana cost of any recalled totems.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shaman_totemrecall.jpg",
		},
		[ShamanMajorGlyph.GlyphOfTelluricCurrents]: {
			name: "Glyph of Telluric Currents",
			description: "Causes your Lightning Bolt to restore 2% of your mana when it strikes an enemy.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_lightning_lightningbolt01.jpg",
		},
		[ShamanMajorGlyph.GlyphOfGroundingTotem]: {
			name: "Glyph of Grounding Totem",
			description: "Instead of absorbing a spell, your Grounding Totem reflects the next harmful spell back at its caster, but the cooldown of your Grounding Totem is increased by 20 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_groundingtotem.jpg",
		},
		[ShamanMajorGlyph.GlyphOfSpiritwalkersGrace]: {
			name: "Glyph of Spiritwalker's Grace",
			description: "Increases the duration of your Spiritwalker\'s Grace by 5 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shaman_spiritwalkersgrace.jpg",
		},
		[ShamanMajorGlyph.GlyphOfWaterShield]: {
			name: "Glyph of Water Shield",
			description: "Increases the mana generated reactively by your Water Shield when you are attacked by 50%, but reduces the passive mana generation by 15%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_shaman_watershield.jpg",
		},
		[ShamanMajorGlyph.GlyphOfCleansingWaters]: {
			name: "Glyph of Cleansing Waters",
			description: "When you dispel a harmful Magic or Curse effect from an ally, you also heal the target for 5% of your maximum health.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_regeneration_02.jpg",
		},
		[ShamanMajorGlyph.GlyphOfFrostShock]: {
			name: "Glyph of Frost Shock",
			description: "Decreases the cooldown incurred by your Frost Shock by 2 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_frost_frostshock.jpg",
		},
		[ShamanMajorGlyph.GlyphOfChaining]: {
			name: "Glyph of Chaining",
			description: "Increases the jump distance of your Chain Heal spell by 100%, but gives the spell a 2 sec cooldown.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_druid_empoweredtouch.jpg",
		},
		[ShamanMajorGlyph.GlyphOfHealingStorm]: {
			name: "Glyph of Healing Storm",
			description: "Each application of Maelstrom Weapon also increases your direct healing done by 20%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shaman_maelstromweapon.jpg",
		},
		[ShamanMajorGlyph.GlyphOfGhostWolf]: {
			name: "Glyph of Ghost Wolf",
			description: "While in Ghost Wolf form, you are less hindered by effects that would reduce movement speed.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_spiritwolf.jpg",
		},
		[ShamanMajorGlyph.GlyphOfThunder]: {
			name: "Glyph of Thunder",
			description: "Reduces the cooldown on Thunderstorm by 10 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shaman_thunderstorm.jpg",
		},
		[ShamanMajorGlyph.GlyphOfFeralSpirit]: {
			name: "Glyph of Feral Spirit",
			description: "Increases the healing done by your Feral Spirits\' Spirit Hunt by 40%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shaman_feralspirit.jpg",
		},
		[ShamanMajorGlyph.GlyphOfRiptide]: {
			name: "Glyph of Riptide",
			description: "Removes the cooldown of Riptide, but reduces the initial direct healing by 75%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_riptide.jpg",
		},
		[ShamanMajorGlyph.GlyphOfShamanisticRage]: {
			name: "Glyph of Shamanistic Rage",
			description: "Activating your Shamanistic Rage ability also cleanses you of all dispellable harmful Magic effects.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_shamanrage.jpg",
		},
		[ShamanMajorGlyph.GlyphOfHex]: {
			name: "Glyph of Hex",
			description: "Reduces the cooldown of your Hex spell by 10 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shaman_hex.jpg",
		},
		[ShamanMajorGlyph.GlyphOfTotemicVigor]: {
			name: "Glyph of Totemic Vigor",
			description: "Increases the health of your totems by 5% of your maximum health.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_relics_totemofrebirth.jpg",
		},
		[ShamanMajorGlyph.GlyphOfLightningShield]: {
			name: "Glyph of Lightning Shield",
			description: "When your Lightning Shield is triggered, you take 10% less damage for 6s.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_lightningshield.jpg",
		},
		[ShamanMajorGlyph.GlyphOfPurging]: {
			name: "Glyph of Purging",
			description: "Successfully Purging a target now grants a stack of Maelstrom Weapon.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shaman_improvedreincarnation.jpg",
		},
		[ShamanMajorGlyph.GlyphOfEternalEarth]: {
			name: "Glyph of Eternal Earth",
			description: "Your Lightning Bolt has a chance to add a charge to your currently active Earth Shield. This cannot cause Earth Shield to exceed 9 charges.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shaman_improvedearthshield.jpg",
		},
	},
	minorGlyphs: {
		[ShamanMinorGlyph.GlyphOfTheLakestrider]: {
			name: "Glyph of the Lakestrider",
			description: "You automatically gain Water Walking while you are in Ghost Wolf form.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_frost_windwalkon.jpg",
		},
		[ShamanMinorGlyph.GlyphOfLavaLash]: {
			name: "Glyph of Lava Lash",
			description: "Your Lava Lash ability no longer spreads Flame Shock to nearby targets.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_shaman_lavalash.jpg",
		},
		[ShamanMinorGlyph.GlyphOfAstralRecall]: {
			name: "Glyph of Astral Recall",
			description: "Reduces the cooldown of your Astral Recall spell by 5 minutes.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_astralrecal.jpg",
		},
		[ShamanMinorGlyph.GlyphOfFarSight]: {
			name: "Glyph of Far Sight",
			description: "Your Far Sight spell may be used indoors.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_farsight.jpg",
		},
		[ShamanMinorGlyph.GlyphOfTheSpectralWolf]: {
			name: "Glyph of the Spectral Wolf",
			description: "Alters the appearance of your Ghost Wolf transformation, causing it to resemble a large, spectral wolf.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_mount_whitedirewolf.jpg",
		},
		[ShamanMinorGlyph.GlyphOfTotemicEncirclement]: {
			name: "Glyph of Totemic Encirclement",
			description: "When you cast a totem spell, you also place unempowered totems for any elements that are not currently active. These totems have 5 health and produce no other effects.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_shaman_totemrelocation.jpg",
		},
		[ShamanMinorGlyph.GlyphOfThunderstorm]: {
			name: "Glyph of Thunderstorm",
			description: "Removes the knockback effect from your Thunderstorm spell.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shaman_thunderstorm.jpg",
		},
		[ShamanMinorGlyph.GlyphOfDeluge]: {
			name: "Glyph of Deluge",
			description: "Your Chain Heal now has a watery appearance.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_healingwavegreater.jpg",
		},
		[ShamanMinorGlyph.GlyphOfSpiritRaptors]: {
			name: "Glyph of Spirit Raptors",
			description: "Your Spirit Wolves are replaced with Spirit Raptors.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/trade_archaeology_dinosaurskeleton.jpg",
		},
		[ShamanMinorGlyph.GlyphOfLingeringAncestors]: {
			name: "Glyph of Lingering Ancestors",
			description: "Resurrecting someone with Ancestral Spirit causes a ghostly ancestor to follow them around for a short time.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shaman_ancestralawakening.jpg",
		},
		[ShamanMinorGlyph.GlyphOfSpiritWolf]: {
			name: "Glyph of Spirit Wolf",
			description: "Ghost Wolf can be now be used while you are a ghost.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_spiritwolf.jpg",
		},
		[ShamanMinorGlyph.GlyphOfFlamingSerpent]: {
			name: "Glyph of Flaming Serpent",
			description: "Your Searing Totem now resembles Vol\'jin\'s Serpent Ward.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_guardianward.jpg",
		},
		[ShamanMinorGlyph.GlyphOfTheCompy]: {
			name: "Glyph of the Compy",
			description: "Your Hex now transforms enemies into a Compy.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_hunter_pet_raptor.jpg",
		},
		[ShamanMinorGlyph.GlyphOfElementalFamiliars]: {
			name: "Glyph of Elemental Familiars",
			description: "Summons a random Fire, Water, or Nature familiar. Familiars of different types have a tendency to fight each other.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_pet_pandarenelemental.jpg",
		},
		[ShamanMinorGlyph.GlyphOfAstralFixation]: {
			name: "Glyph of Astral Fixation",
			description: "Astral Recall now takes you to your capital\'s Earthshrine.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_astralrecalgroup.jpg",
		},
		[ShamanMinorGlyph.GlyphOfRainOfFrogs]: {
			name: "Glyph of Rain of Frogs",
			description: "You summon a rain storm of frogs at your targeted location.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_pet_toad_blue.jpg",
		},
	},
};
