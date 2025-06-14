import { DruidMajorGlyph, DruidMinorGlyph, DruidTalents } from '../proto/druid.js';
import { GlyphsConfig } from './glyphs_picker.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import DruidTalentJson from './trees/druid.json';export const druidTalentsConfig: TalentsConfig<DruidTalents> = newTalentsConfig(DruidTalentJson);

export const druidGlyphsConfig: GlyphsConfig = {
	majorGlyphs: {
		[DruidMajorGlyph.GlyphOfFrenziedRegeneration]: {
			name: "Glyph of Frenzied Regeneration",
			description: "For 6s after activating Frenzied Regeneration, healing effects on you are 40% more powerful. However, your Frenzied Regeneration now always costs 50 Rage and no longer converts Rage into health.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_bullrush.jpg",
		},
		[DruidMajorGlyph.GlyphOfMaul]: {
			name: "Glyph of Maul",
			description: "Your Maul ability now hits 1 additional target for 50% damage.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_druid_maul.jpg",
		},
		[DruidMajorGlyph.GlyphOfOmens]: {
			name: "Glyph of Omens",
			description: "While you are not in an Eclipse, the following abilities now grant 10 Solar or Lunar Energy: Entangling Roots, Cyclone, Faerie Fire, Faerie Swarm, Mass Entanglement, Typhoon, Disorienting Roar, Ursol\'s Vortex, and Mighty Bash.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_vehicle_sonicshockwave.jpg",
		},
		[DruidMajorGlyph.GlyphOfShred]: {
			name: "Glyph of Shred",
			description: "While Berserk or Tiger\'s Fury is active, Shred has no positional requirement.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_vampiricaura.jpg",
		},
		[DruidMajorGlyph.GlyphOfProwl]: {
			name: "Glyph of Prowl",
			description: "Reduces the movement penalty of Prowl by 100%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_druid_prowl.jpg",
		},
		[DruidMajorGlyph.GlyphOfPounce]: {
			name: "Glyph of Pounce",
			description: "Increases the range of your Pounce by 8 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_druid_supriseattack.jpg",
		},
		[DruidMajorGlyph.GlyphOfStampede]: {
			name: "Glyph of Stampede",
			description: "You can now cast Stampeding Roar without being in Bear Form or Cat Form.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_druid_stampedingroar_cat.jpg",
		},
		[DruidMajorGlyph.GlyphOfInnervate]: {
			name: "Glyph of Innervate",
			description: "When Innervate is cast on a target other than the caster, both the caster and target will benefit, but at 40% reduced effect.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_lightning.jpg",
		},
		[DruidMajorGlyph.GlyphOfRebirth]: {
			name: "Glyph of Rebirth",
			description: "Players resurrected by Rebirth are returned to life with 100% health.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_reincarnation.jpg",
		},
		[DruidMajorGlyph.GlyphOfRegrowth]: {
			name: "Glyph of Regrowth",
			description: "Increases the critical strike chance of your Regrowth by 40%, but removes the periodic component of the spell.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_resistnature.jpg",
		},
		[DruidMajorGlyph.GlyphOfRejuvenation]: {
			name: "Glyph of Rejuvenation",
			description: "When you have Rejuvenation active on three or more targets, the cast time of your Nourish spell is reduced by 30%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_rejuvenation.jpg",
		},
		[DruidMajorGlyph.GlyphOfHealingTouch]: {
			name: "Glyph of Healing Touch",
			description: "When you cast Healing Touch, the cooldown on your Nature\'s Swiftness is reduced by 3 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_healingtouch.jpg",
		},
		[DruidMajorGlyph.GlyphOfEfflorescence]: {
			name: "Glyph of Efflorescence",
			description: "The Efflorescence effect is now caused by your Wild Mushroom instead of by Swiftmend, and lasts as long as the Wild Mushroom is active. Additionally, increases the healing done by Swiftmend by 20%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_healingtouch.jpg",
		},
		[DruidMajorGlyph.GlyphOfGuidedStars]: {
			name: "Glyph of Guided Stars",
			description: "Your Starfall only hits targets affected by your Moonfire or Sunfire.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_mage_arcanebarrage.jpg",
		},
		[DruidMajorGlyph.GlyphOfHurricane]: {
			name: "Glyph of Hurricane",
			description: "Your Hurricane and Astral Storm abilities now also slow the movement speed of their victims by 50%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_cyclone.jpg",
		},
		[DruidMajorGlyph.GlyphOfSkullBash]: {
			name: "Glyph of Skull Bash",
			description: "Increases the duration of your Skull Bash interrupt by 2 sec, but increases the cooldown by 5 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_misc_bone_taurenskull_01.jpg",
		},
		[DruidMajorGlyph.GlyphOfNaturesGrasp]: {
			name: "Glyph of Nature's Grasp",
			description: "Reduces the cooldown of Nature\'s Grasp by 45 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_natureswrath.jpg",
		},
		[DruidMajorGlyph.GlyphOfSavagery]: {
			name: "Glyph of Savagery",
			description: "Savage Roar can now be used with 0 combo points, resulting in a 12s duration.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_druid_skinteeth.jpg",
		},
		[DruidMajorGlyph.GlyphOfEntanglingRoots]: {
			name: "Glyph of Entangling Roots",
			description: "Reduces the cast time of your Entangling Roots by 0.2 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_stranglevines.jpg",
		},
		[DruidMajorGlyph.GlyphOfBlooming]: {
			name: "Glyph of Blooming",
			description: "Increases the bloom heal of your Lifebloom when it expires by 50%, but its duration is reduced by 5 sec and your Healing Touch, Nourish, and Regrowth abilities no longer refresh the duration.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_protectionformnature.jpg",
		},
		[DruidMajorGlyph.GlyphOfDash]: {
			name: "Glyph of Dash",
			description: "Reduces the cooldown of your Dash ability by 60 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_druid_dash.jpg",
		},
		[DruidMajorGlyph.GlyphOfMasterShapeshifter]: {
			name: "Glyph of Master Shapeshifter",
			description: "Reduces the mana cost of all shapeshifts by 100%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_druid_mastershapeshifter.jpg",
		},
		[DruidMajorGlyph.GlyphOfSurvivalInstincts]: {
			name: "Glyph of Survival Instincts",
			description: "Reduces the cooldown of Survival Instincts by 60 sec, but reduces its duration by 50%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_druid_tigersroar.jpg",
		},
		[DruidMajorGlyph.GlyphOfWildGrowth]: {
			name: "Glyph of Wild Growth",
			description: "Wild Growth can affect 1 additional target, but its cooldown is increased by 2 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_druid_flourish.jpg",
		},
		[DruidMajorGlyph.GlyphOfMightOfUrsoc]: {
			name: "Glyph of Might of Ursoc",
			description: "Increases the health gain from Might of Ursoc by 20%, but increases the cooldown by 2 min.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_druid_mightofursoc.jpg",
		},
		[DruidMajorGlyph.GlyphOfStampedingRoar]: {
			name: "Glyph of Stampeding Roar",
			description: "Increases the radius of Stampeding Roar by 30 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_druid_stamedingroar.jpg",
		},
		[DruidMajorGlyph.GlyphOfCyclone]: {
			name: "Glyph of Cyclone",
			description: "Increases the range of your Cyclone spell by 5 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_earthbind.jpg",
		},
		[DruidMajorGlyph.GlyphOfBarkskin]: {
			name: "Glyph of Barkskin",
			description: "Reduces the chance you\'ll be critically hit by 25% while Barkskin is active.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_stoneclawtotem.jpg",
		},
		[DruidMajorGlyph.GlyphOfFerociousBite]: {
			name: "Glyph of Ferocious Bite",
			description: "Your Ferocious Bite ability heals you for 2% of your maximum health for each 10 Energy used.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_druid_ferociousbite.jpg",
		},
		[DruidMajorGlyph.GlyphOfFaeSilence]: {
			name: "Glyph of Fae Silence",
			description: "Faerie Fire used in Bear Form also silences the target for 3s, but triggers a 15 sec cooldown.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_druid_challangingroar.jpg",
		},
		[DruidMajorGlyph.GlyphOfFaerieFire]: {
			name: "Glyph of Faerie Fire",
			description: "Increases the range of your Faerie Fire by 10 yds.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_faeriefire.jpg",
		},
		[DruidMajorGlyph.GlyphOfCatForm]: {
			name: "Glyph of Cat Form",
			description: "Increases healing done to you by 20% while in Cat Form.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_druid_catform.jpg",
		},
	},
	minorGlyphs: {
		[DruidMinorGlyph.GlyphOfTheStag]: {
			name: "Glyph of the Stag",
			description: "Your Travel Form can now be used as a mount by party members. This glyph is disabled while Glyph of the Cheetah is active.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/trade_archaeology_antleredcloakclasp.jpg",
		},
		[DruidMinorGlyph.GlyphOfTheOrca]: {
			name: "Glyph of the Orca",
			description: "Your Aquatic Form now appears as an Orca.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_druid_aquaticform.jpg",
		},
		[DruidMinorGlyph.GlyphOfAquaticForm]: {
			name: "Glyph of Aquatic Form",
			description: "Increases your swim speed by 50% while in Aquatic Form.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_enchant_essencemagiclarge.jpg",
		},
		[DruidMinorGlyph.GlyphOfGrace]: {
			name: "Glyph of Grace",
			description: "Feline Grace reduces falling damage even while not in Cat Form.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_feather_01.jpg",
		},
		[DruidMinorGlyph.GlyphOfTheChameleon]: {
			name: "Glyph of the Chameleon",
			description: "Each time you shapeshift into Cat Form or Bear Form, your shapeshifted form will have a random hair color.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_druid_mastershapeshifter.jpg",
		},
		[DruidMinorGlyph.GlyphOfCharmWoodlandCreature]: {
			name: "Glyph of Charm Woodland Creature",
			description: "Teaches you the ability Charm Woodland Creature.\u000D\u000A\u000D\u000A Allows the Druid to befriend an ambient creature, which will follow the Druid for 1hr.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_misc_rabbit.jpg",
		},
		[DruidMinorGlyph.GlyphOfStars]: {
			name: "Glyph of Stars",
			description: "Your Moonkin Form now appears as Astral Form, conferring all the same benefits, but appearing as an astrally enhanced version of your normal humanoid form.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/achievement_boss_algalon_01.jpg",
		},
		[DruidMinorGlyph.GlyphOfThePredator]: {
			name: "Glyph of the Predator",
			description: "Your Track Humanoids ability now also tracks beasts.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_tracking.jpg",
		},
		[DruidMinorGlyph.GlyphOfTheTreant]: {
			name: "Glyph of the Treant",
			description: "Teaches you the ability Treant Form.\u000D\u000A\u000D\u000A Shapeshift into Treant Form.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_druid_treeoflife.jpg",
		},
		[DruidMinorGlyph.GlyphOfTheCheetah]: {
			name: "Glyph of the Cheetah",
			description: "Your Travel Form appears as a Cheetah. This glyph will prevent Glyph of the Stag from functioning.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_hunter_catlikereflexes.jpg",
		},
		[DruidMinorGlyph.GlyphOfFocus]: {
			name: "Glyph of Focus",
			description: "Reduces Starfall\'s radius by 50%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_arcane_arcanepotency.jpg",
		},
		[DruidMinorGlyph.GlyphOfTheSproutingMushroom]: {
			name: "Glyph of the Sprouting Mushroom",
			description: "Your Wild Mushroom spell can now be placed on the ground instead of underneath a target.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_druid_wildmushroom.jpg",
		},
		[DruidMinorGlyph.GlyphOfOneWithNature]: {
			name: "Glyph of One with Nature",
			description: "Grants you the ability to teleport to a random natural location.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_druid_manatree.jpg",
		},
	},
};
