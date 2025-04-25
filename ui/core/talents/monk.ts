import { MonkMajorGlyph, MonkMinorGlyph, MonkTalents } from '../proto/monk.js';
import { GlyphsConfig } from './glyphs_picker.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import MonkTalentJson from './trees/monk.json';

export const monkTalentsConfig: TalentsConfig<MonkTalents> = newTalentsConfig(MonkTalentJson);

export const monkGlyphsConfig: GlyphsConfig = {
	primeGlyphs: {},
	majorGlyphs: {
		[MonkMajorGlyph.MonkMajorGlyphAfterlife]: {
			name: 'Glyph of Afterlife',
			description: 'Increases the chance to summon a Healing Sphere when you kill an enemy while gaining experience or honor by 100%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_priest_finalprayer.jpg',
		},
		[MonkMajorGlyph.MonkMajorGlyphNimbleBrew]: {
			name: 'Glyph of Nimble Brew',
			description: "Clearing an effect with Nimble Brew heals you for 10% of your maximum health.",
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_monk_nimblebrew.jpg',
		},
		[MonkMajorGlyph.MonkMajorGlyphEnduringHealingSphere]: {
			name: 'Glyph of Enduring Healing Sphere',
			description: 'Increases the duration of your Healing Spheres by 3 minutes.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_monk_healthsphere.jpg',
		},
		[MonkMajorGlyph.MonkMajorGlyphFistsOfFury]: {
			name: "Glyph of Fists of Fury",
			description: "When channeling Fists of Fury, your parry chance is increased by 100%.",
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/monk_ability_fistoffury.jpg',
		},
		[MonkMajorGlyph.MonkMajorGlyphFortifyingBrew]: {
			name: 'Glyph of Fortifying Brew',
			description: 'Your Fortifying Brew reduces damage taken by an additional 5%, but increases your health by 10% rather than 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_monk_fortifyingale_new.jpg',
		},
		[MonkMajorGlyph.MonkMajorGlyphFortuitousSpheres]: {
			name: 'Glyph of Fortuitous Spheres',
			description: 'Falling below 25% health will automatically summon a healing sphere near you at no cost. This effect cannot occur more often than once every 30 seconds.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_monk_healthsphere.jpg',
		},
		[MonkMajorGlyph.MonkMajorGlyphSparring]: {
			name: 'Glyph of Sparring',
			description: 'While Sparring, you also have a 5% chance to deflect spells from attackers in front of you, stacking up to 3 times.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_monk_sparring.jpg',
		},
		[MonkMajorGlyph.MonkMajorGlyphSpinningCraneKick]: {
			name: 'Glyph of Spinning Crane Kick',
			description: 'You move at full speed while channeling Spinning Crane Kick.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_monk_cranekick_new.jpg',
		},
		[MonkMajorGlyph.MonkMajorGlyphDetox]: {
			name: 'Glyph of Detox',
			description: 'Detox heals your target for 5% when it successfully removes a harmful effect.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_dispelmagic.jpg',
		},
		[MonkMajorGlyph.MonkMajorGlyphTouchOfDeath]: {
			name: 'Glyph of Touch of Death',
			description: 'Your Touch of Death no longer has a Chi cost, but the cooldown is increased by 2 minutes.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_monk_touchofdeath.jpg',
		},
		[MonkMajorGlyph.MonkMajorGlyphTouchOfKarma]: {
			name: 'Glyph of Touch of Karma',
			description: 'Your Touch of Karma now has a 20 yard range.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_monk_touchofkarma.jpg',
		},
		[MonkMajorGlyph.MonkMajorGlyphTranscendence]: {
			name: 'Glyph of Transcendence',
			description: 'Increases the range of your Transcendence: Transfer spell by 10 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/monk_ability_transcendence.jpg',
		},
		[MonkMajorGlyph.MonkMajorGlyphZenMeditation]: {
			name: 'Glyph of Zen Meditation',
			description: 'You can now channel Zen Meditation while moving.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_monk_zenmeditation.jpg',
		},
	},
	minorGlyphs: {
		[MonkMinorGlyph.MonkMinorGlyphBlackoutKick]: {
			name: 'Glyph of Blackout Kick',
			description: "Your Blackout Kick always deals 20% additional damage over 4 sec regardless of positioning but you're unable to trigger the healing effect.",
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_monk_blackoutkick.jpg',
		},
		[MonkMinorGlyph.MonkMinorGlyphCracklingTigerLightning]: {
			name: 'Glyph of Crackling Tiger Lightning',
			description: 'Your Crackling Jade Lightning visual is altered to the color of the White Tiger celestial.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_monk_cracklingjadelightning.jpg',
		},
		[MonkMinorGlyph.MonkMinorGlyphFightingPose]: {
			name: 'Glyph of Fighting Pose',
			description: 'Your spirit now appears in a fighting pose when using Transcendence.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_monk_dpsstance.jpg',
		},
		[MonkMinorGlyph.MonkMinorGlyphFlyingSerpentKick]: {
			name: 'Glyph of Flying Serpent Kick',
			description: 'Your Flying Serpent Kick automatically ends when you fly into an enemy, triggering the area of effect damage and snare.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_monk_flyingdragonkick.jpg',
		},
		[MonkMinorGlyph.MonkMinorGlyphHonor]: {
			name: 'Glyph of Honor',
			description: 'You honorably bow after each successful Touch of Death.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/pandarenracial_innerpeace.jpg',
		},
		[MonkMinorGlyph.MonkMinorGlyphJab]: {
			name: 'Glyph of Jab',
			description: 'You always will attack with hands and fist with Jab, even with non-fist weapons equipped.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_monk_jab.jpg',
		},
		[MonkMinorGlyph.MonkMinorGlyphRisingTigerKick]: {
			name: 'Glyph of Rising Tiger Kick',
			description: "Your Rising Sun Kick's visual is altered to the color of the White Tiger.",
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_monk_risingsunkick.jpg',
		},
		[MonkMinorGlyph.MonkMinorGlyphSpinningFireBlossom]: {
			name: 'Glyph of Spinning Fire Blossom',
			description: 'Your Spinning Fire Blossom requires an enemy target rather than traveling in front of you, but is no longer capable of rooting targets.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_monk_explodingjadeblossom.jpg',
		},
		[MonkMinorGlyph.MonkMinorGlyphSpiritRoll]: {
			name: 'Glyph of Spirit Roll',
			description: 'You can cast Roll or Chi Torpedo while dead as a spirit.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_guardianspirit.jpg',
		},
		[MonkMinorGlyph.MonkMinorGlyphWaterRoll]: {
			name: 'Glyph of Water Roll',
			description: 'You can Roll or Chi Torpedo over water',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_monk_roll.jpg',
		},
		[MonkMinorGlyph.MonkMinorGlyphZenFlight]: {
			name: 'Glyph of Zen Flight',
			description: "Teaches you the spell Zen Flight. Zen Flight requires a Flight Master's License in order to be cast.",
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_monk_zenflight.jpg',
		},
	},
};
