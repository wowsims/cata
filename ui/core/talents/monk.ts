import { MonkMajorGlyph, MonkMinorGlyph, MonkTalents } from '../proto/monk.js';
import { GlyphsConfig } from './glyphs_picker.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import MonkTalentJson from './trees/monk.json';export const monkTalentsConfig: TalentsConfig<MonkTalents> = newTalentsConfig(MonkTalentJson);

export const monkGlyphsConfig: GlyphsConfig = {
	majorGlyphs: {
		[MonkMajorGlyph.GlyphOfRapidRolling]: {
			name: "Glyph of Rapid Rolling",
			description: "For 5s seconds after using Roll or Chi Torpedo, your next Roll or Chi Torpedo will go 30% farther.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_monk_roll.jpg",
		},
		[MonkMajorGlyph.GlyphOfTranscendence]: {
			name: "Glyph of Transcendence",
			description: "Reduces the cooldown of your Transcendence: Transfer spell by 5 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/monk_ability_transcendence.jpg",
		},
		[MonkMajorGlyph.GlyphOfBreathOfFire]: {
			name: "Glyph of Breath of Fire",
			description: "When you use Breath of Fire on targets afflicted with your Dizzying Haze, they become Disoriented for 3s.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_monk_breathoffire.jpg",
		},
		[MonkMajorGlyph.GlyphOfClash]: {
			name: "Glyph of Clash",
			description: "Increases the range of your Clash ability by 10 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_monk_clashingoxcharge.jpg",
		},
		[MonkMajorGlyph.GlyphOfEnduringHealingSphere]: {
			name: "Glyph of Enduring Healing Sphere",
			description: "Increases the duration of your Healing Spheres by 3 minutes.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_monk_healthsphere.jpg",
		},
		[MonkMajorGlyph.GlyphOfGuard]: {
			name: "Glyph of Guard",
			description: "Increases the amount your Guard absorbs by 10%, but your Guard can only absorb magical damage.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_monk_guard.jpg",
		},
		[MonkMajorGlyph.GlyphOfManaTea]: {
			name: "Glyph of Mana Tea",
			description: "Your Mana Tea is instant instead of channeled and consumes two stacks when used, but causes a 10 sec cooldown.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/monk_ability_cherrymanatea.jpg",
		},
		[MonkMajorGlyph.GlyphOfZenMeditation]: {
			name: "Glyph of Zen Meditation",
			description: "You can now channel Zen Meditation while moving.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_monk_zenmeditation.jpg",
		},
		[MonkMajorGlyph.GlyphOfRenewingMists]: {
			name: "Glyph of Renewing Mists",
			description: "Your Renewing Mist travels to the furthest injured target within 40 yards rather than the closest injured target within 20 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_monk_renewingmists.jpg",
		},
		[MonkMajorGlyph.GlyphOfSpinningCraneKick]: {
			name: "Glyph of Spinning Crane Kick",
			description: "You move at full speed while channeling Spinning Crane Kick.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_monk_cranekick_new.jpg",
		},
		[MonkMajorGlyph.GlyphOfSurgingMist]: {
			name: "Glyph of Surging Mist",
			description: "Your Surging Mist no longer requires a target, and instead heals the lowest health target within 40 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_monk_surgingmist.jpg",
		},
		[MonkMajorGlyph.GlyphOfTouchOfDeath]: {
			name: "Glyph of Touch of Death",
			description: "Your Touch of Death no longer has a Chi cost, but the cooldown is increased by 2 minutes.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_monk_touchofdeath.jpg",
		},
		[MonkMajorGlyph.GlyphOfNimbleBrew]: {
			name: "Glyph of Nimble Brew",
			description: "Clearing an effect with Nimble Brew heals you for 10% of your maximum health.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_monk_nimblebrew.jpg",
		},
		[MonkMajorGlyph.GlyphOfAfterlife]: {
			name: "Glyph of Afterlife",
			description: "Increases the chance to summon a Healing Sphere when you kill an enemy while gaining experience or honor to 100%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_priest_finalprayer.jpg",
		},
		[MonkMajorGlyph.GlyphOfFistsOfFury]: {
			name: "Glyph of Fists of Fury",
			description: "When channeling Fists of Fury, your parry chance is increased by 100%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/monk_ability_fistoffury.jpg",
		},
		[MonkMajorGlyph.GlyphOfFortifyingBrew]: {
			name: "Glyph of Fortifying Brew",
			description: "Your Fortifying Brew reduces damage taken by an additional 5%, but increases your health by 10% rather than 20%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_monk_fortifyingale_new.jpg",
		},
		[MonkMajorGlyph.GlyphOfLeerOfTheOx]: {
			name: "Glyph of Leer of the Ox",
			description: "Teaches you the spell Leer of the Ox.\u000D\u000A\u000D\u000A |CFFFFFFFFLeer of the Ox|R\u000D\u000A Your Ox Statue stirs hatred in the target, reducing their movement speed by 50% and causing them to attack your Black Ox Statue for 8s. The statue must be within 40 yards of the target.\u000D\u000A\u000D\u000A Requires Black Ox Statue to be active.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_monk_leeroftheox.jpg",
		},
		[MonkMajorGlyph.GlyphOfLifeCocoon]: {
			name: "Glyph of Life Cocoon",
			description: "Life Cocoon can now be cast while stunned.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_monk_chicocoon.jpg",
		},
		[MonkMajorGlyph.GlyphOfFortuitousSpheres]: {
			name: "Glyph of Fortuitous Spheres",
			description: "Falling below 25% health will automatically summon a healing sphere near you at no cost. This effect cannot occur more often than once every 30 seconds.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_monk_healthsphere.jpg",
		},
		[MonkMajorGlyph.GlyphOfParalysis]: {
			name: "Glyph of Paralysis",
			description: "Your Paralysis ability also removes all damage over time effects from the target.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_monk_paralysis.jpg",
		},
		[MonkMajorGlyph.GlyphOfSparring]: {
			name: "Glyph of Sparring",
			description: "While Sparring, you also have a 5% chance to deflect spells from attackers in front of you, stacking up to 3 times.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_monk_sparring.jpg",
		},
		[MonkMajorGlyph.GlyphOfDetox]: {
			name: "Glyph of Detox",
			description: "Detox heals your target for 5% when it successfully removes a harmful effect.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_dispelmagic.jpg",
		},
		[MonkMajorGlyph.GlyphOfTouchOfKarma]: {
			name: "Glyph of Touch of Karma",
			description: "Your Touch of Karma now has a 25 yard range.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_monk_touchofkarma.jpg",
		},
		[MonkMajorGlyph.GlyphOfTargetedExpulsion]: {
			name: "Glyph of Targeted Expulsion",
			description: "Expel Harm can now be used on other allies, but the healing is reduced by 50% on them.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_monk_expelharm.jpg",
		},
	},
	minorGlyphs: {
		[MonkMinorGlyph.GlyphOfSpinningFireBlossom]: {
			name: "Glyph of Spinning Fire Blossom",
			description: "Your Spinning Fire Blossom requires an enemy target rather than traveling in front of you, but is no longer capable of rooting targets.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_monk_explodingjadeblossom.jpg",
		},
		[MonkMinorGlyph.GlyphOfCracklingTigerLightning]: {
			name: "Glyph of Crackling Tiger Lightning",
			description: "Your Crackling Jade Lightning visual is altered to the color of the White Tiger celestial.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_monk_cracklingjadelightning.jpg",
		},
		[MonkMinorGlyph.GlyphOfFlyingSerpentKick]: {
			name: "Glyph of Flying Serpent Kick",
			description: "Your Flying Serpent Kick automatically ends when you fly into an enemy, triggering the area of effect damage and snare.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_monk_flyingdragonkick.jpg",
		},
		[MonkMinorGlyph.GlyphOfHonor]: {
			name: "Glyph of Honor",
			description: "You honorably bow after each successful Touch of Death.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/pandarenracial_innerpeace.jpg",
		},
		[MonkMinorGlyph.GlyphOfJab]: {
			name: "Glyph of Jab",
			description: "You always will attack with hands and fist with Jab, even with non-fist weapons equipped.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_monk_jab.jpg",
		},
		[MonkMinorGlyph.GlyphOfRisingTigerKick]: {
			name: "Glyph of Rising Tiger Kick",
			description: "Your Rising Sun Kick\'s visual is altered to the color of the White Tiger.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_monk_risingsunkick.jpg",
		},
		[MonkMinorGlyph.GlyphOfSpiritRoll]: {
			name: "Glyph of Spirit Roll",
			description: "You can cast Roll or Chi Torpedo while dead as a spirit.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_guardianspirit.jpg",
		},
		[MonkMinorGlyph.GlyphOfFightingPose]: {
			name: "Glyph of Fighting Pose",
			description: "Your spirit now appears in a fighting pose when using Transcendence.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_monk_dpsstance.jpg",
		},
		[MonkMinorGlyph.GlyphOfWaterRoll]: {
			name: "Glyph of Water Roll",
			description: "You can Roll or Chi Torpedo over water.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_monk_roll.jpg",
		},
		[MonkMinorGlyph.GlyphOfZenFlight]: {
			name: "Glyph of Zen Flight",
			description: "Teaches you the spell Zen Flight. Zen Flight requires a Flight Master\'s License in order to be cast.\u000D\u000A\u000D\u000A |CFFFFFFFFZen Flight|R\u000D\u000A You fly through the air at a quick speed on a meditative cloud.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_monk_zenflight.jpg",
		},
		[MonkMinorGlyph.GlyphOfBlackoutKick]: {
			name: "Glyph of Blackout Kick",
			description: "Your Blackout Kick always deals 20% additional damage over 4s regardless of positioning but you\'re unable to trigger the healing effect.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_monk_blackoutkick.jpg",
		},
	},
};
