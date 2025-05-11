import { WarlockMajorGlyph, WarlockMinorGlyph, WarlockTalents } from '../proto/warlock.js';
import { GlyphsConfig } from './glyphs_picker.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import WarlockTalentJson from './trees/warlock.json';export const warlockTalentsConfig: TalentsConfig<WarlockTalents> = newTalentsConfig(WarlockTalentJson);

export const warlockGlyphsConfig: GlyphsConfig = {
	majorGlyphs: {
		[WarlockMajorGlyph.GlyphOfConflagrate]: {
			name: "Glyph of Conflagrate",
			description: "Conflagrate no longer requires Immolate to snare the target.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_fire_fireball.jpg",
		},
		[WarlockMajorGlyph.GlyphOfSiphonLife]: {
			name: "Glyph of Siphon Life",
			description: "Your Immolate spell will heal you for 0.5% of your maximum health when dealing periodic damage.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_requiem.jpg",
		},
		[WarlockMajorGlyph.GlyphOfFear]: {
			name: "Glyph of Fear",
			description: "Your Fear causes the target to tremble in place instead of fleeing in fear.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_possession.jpg",
		},
		[WarlockMajorGlyph.GlyphOfDemonTraining]: {
			name: "Glyph of Demon Training",
			description: "Improves your demon\'s special abilities:\u000D\u000A\u000D\u000A Your Fel Imp \'s Firebolt cast time is reduced by 50% and fires in bursts of three.\u000D\u000A\u000D\u000A Increases your Voidlord \'s total armor by 10%.\u000D\u000A\u000D\u000A Your Shivarra \'s Mesmerizes ability also removes all damage over time effects from the target.\u000D\u000A\u000D\u000A When your Observer uses Clone Magic, you will also be healed for that amount.\u000D\u000A\u000D\u000A Increases your Wrathguard \'s total health by 20%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_summonfelhunter.jpg",
		},
		[WarlockMajorGlyph.GlyphOfHealthstone]: {
			name: "Glyph of Healthstone",
			description: "You receive 100% more healing from using a healthstone, but the health is restored over 10 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_stone_04.jpg",
		},
		[WarlockMajorGlyph.GlyphOfCurseOfTheElements]: {
			name: "Glyph of Curse of the Elements",
			description: "Curse of the Elements hits 2 additional nearby targets.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/warlock_curse_shadow.jpg",
		},
		[WarlockMajorGlyph.GlyphOfImpSwarm]: {
			name: "Glyph of Imp Swarm",
			description: "Teaches you the ability Imp Swarm.\u000D\u000A Requires Demonology.\u000D\u000A\u000D\u000A |Tinterface/icons/ability_warlock_impoweredimp.blp:24|t Imp Swarm\u000D\u000A Summons 4 Wild Imps from the Twisting Nether to attack the target.\u000D\u000A\u000D\u000A The Wild Imps passive effect is disabled while Imp Swarm is on cooldown. Imp Swarm\'s cooldown is reduced by spell haste. Also increases Wild Imp\'s cooldown by 4 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warlock_impoweredimp.jpg",
		},
		[WarlockMajorGlyph.GlyphOfHavoc]: {
			name: "Glyph of Havoc",
			description: "Havoc gains 3 additional charges, but the cooldown is increased by 35 seconds.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warlock_baneofhavoc.jpg",
		},
		[WarlockMajorGlyph.GlyphOfSoulstone]: {
			name: "Glyph of Soulstone",
			description: "Players resurrected by Soulstone are returned to life with 100% health.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_misc_orb_04.jpg",
		},
		[WarlockMajorGlyph.GlyphOfUnstableAffliction]: {
			name: "Glyph of Unstable Affliction",
			description: "Reduces the cast time of Unstable Affliction by 25%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_unstableaffliction_3.jpg",
		},
		[WarlockMajorGlyph.GlyphOfSoulConsumption]: {
			name: "Glyph of Soul Consumption",
			description: "Your Drain Soul restores 20% of your total health after you kill a target that yields experience or honor. You restore 20% of your total health after you kill a target in Demon Form that yields experience or honor. You restore 20% of your total health after you kill a target with Chaos Bolt that yields experience or honor.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warlock_soulsiphon.jpg",
		},
		[WarlockMajorGlyph.GlyphOfCurseOfExhaustion]: {
			name: "Glyph of Curse of Exhaustion",
			description: "Curse of Exhaustion now reduces the targets movement speed by 30%, lasts half as long and has a 10 second cooldown.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_grimward.jpg",
		},
		[WarlockMajorGlyph.GlyphOfDrainLife]: {
			name: "Glyph of Drain Life",
			description: "Increases the healing of your Drain Life by 30%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_lifedrain02.jpg",
		},
		[WarlockMajorGlyph.GlyphOfDemonHunting]: {
			name: "Glyph of Demon Hunting",
			description: "Requires Demonology.\u000D\u000A\u000D\u000A Teaches you the ability Dark Apotheosis.\u000D\u000A\u000D\u000A Dark Apotheosis\u000D\u000A You imbue yourself with demonic energies, reducing physical damage taken by 10.00%, reduces magic damage taken by 15%, and allows the use of various demonic abilities.\u000D\u000A\u000D\u000A In addition, Soulshatter becomes Provocation which taunts your target, Twilight Ward becomes Fury Ward which will absorb all schools of damage, Shadow Bolt becomes Demonic Slash, and Fear becomes Sleep.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_metamorphosis.jpg",
		},
		[WarlockMajorGlyph.GlyphOfEmberTap]: {
			name: "Glyph of Ember Tap",
			description: "Ember Tap heals you for an additional 5% of your health, but the health is restored over 10 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_ember.jpg",
		},
		[WarlockMajorGlyph.GlyphOfDemonicCircle]: {
			name: "Glyph of Demonic Circle",
			description: "Reduces the cooldown on Demonic Circle by 4 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_demoniccircleteleport.jpg",
		},
		[WarlockMajorGlyph.GlyphOfUnendingResolve]: {
			name: "Glyph of Unending Resolve",
			description: "The damage reduction of Unending Resolve is reduced by 20%, but the cooldown is reduced by 60 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_demonictactics.jpg",
		},
		[WarlockMajorGlyph.GlyphOfLifeTap]: {
			name: "Glyph of Life Tap",
			description: "Your Life Tap no longer consumes health, but instead absorbs 0 healing received. This effect stacks.\u000D\u000A\u000D\u000A The absorb lasts 30 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_burningspirit.jpg",
		},
		[WarlockMajorGlyph.GlyphOfEternalResolve]: {
			name: "Glyph of Eternal Resolve",
			description: "Unending Resolve can no longer be activated, but passively provides 10% damage reduction from all sources.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_demonictactics.jpg",
		},
		[WarlockMajorGlyph.GlyphOfSupernova]: {
			name: "Glyph of Supernova",
			description: "When you are killed, all enemies within 8 yards take damage equal to 10% of your maximum health per Burning Ember held.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_fire_ragnaros_supernova.jpg",
		},
	},
	minorGlyphs: {
		[WarlockMinorGlyph.GlyphOfHandOfGuldan]: {
			name: "Glyph of Hand of Gul'dan",
			description: "Your Hand of Gul\'dan can now be targeted at a location.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warlock_handofguldan.jpg",
		},
		[WarlockMinorGlyph.GlyphOfVerdantSpheres]: {
			name: "Glyph of Verdant Spheres",
			description: "Your Soul Shards are transformed into Verdant Spheres.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_misc_enchantedpearlb.jpg",
		},
		[WarlockMinorGlyph.GlyphOfNightmares]: {
			name: "Glyph of Nightmares",
			description: "Your Felsteed and Dreadsteed can cross water while running and leave a trail of flames.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_mount_nightmarehorse.jpg",
		},
		[WarlockMinorGlyph.GlyphOfFelguard]: {
			name: "Glyph of Felguard",
			description: "Your Felguard will equip a random two-handed axe, sword or polearm from your backpack.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_summonfelguard.jpg",
		},
		[WarlockMinorGlyph.GlyphOfHealthFunnel]: {
			name: "Glyph of Health Funnel",
			description: "Your Health Funnel instantly restores 15% of your demon\'s health, but has a 10 sec. cooldown.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_lifedrain.jpg",
		},
		[WarlockMinorGlyph.GlyphOfSubtlety]: {
			name: "Glyph of Subtlety",
			description: "Your Soul Shards no longer display while out of combat.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_rogue_masterofsubtlety.jpg",
		},
		[WarlockMinorGlyph.GlyphOfShadowBolt]: {
			name: "Glyph of Shadow Bolt",
			description: "Splits your Shadow Bolt into three smaller attacks.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_shadowbolt.jpg",
		},
		[WarlockMinorGlyph.GlyphOfCarrionSwarm]: {
			name: "Glyph of Carrion Swarm",
			description: "Your Carrion Swarm no longer knocks targets back.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warlock_demonicpower.jpg",
		},
		[WarlockMinorGlyph.GlyphOfFallingMeteor]: {
			name: "Glyph of Falling Meteor",
			description: "If you use Demonic Leap while falling, you slam into the ground rapidly and will not die from falling damage.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_fire_meteorstorm.jpg",
		},
		[WarlockMinorGlyph.GlyphOfUnendingBreath]: {
			name: "Glyph of Unending Breath",
			description: "Increases the swim speed of targets affected by your Unending Breath spell by 20%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_demonbreath.jpg",
		},
		[WarlockMinorGlyph.GlyphOfEyeOfKilrogg]: {
			name: "Glyph of Eye of Kilrogg",
			description: "Your Eye of Kilrogg is no longer stealthed and can now place your Demonic Circle. The casting Warlock must be within line of sight of the Eye of Kilrogg to place the Demonic Circle.\u000D\u000A\u000D\u000A In addition, the movement speed of your Eye of Kilrogg is increased by 50% and allows it to fly in areas where flying mounts are enabled.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_evileye.jpg",
		},
		[WarlockMinorGlyph.GlyphOfSubjugateDemon]: {
			name: "Glyph of Subjugate Demon",
			description: "Reduces the cast time of your Subjugate Demon spell by 50%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_enslavedemon.jpg",
		},
		[WarlockMinorGlyph.GlyphOfSoulwell]: {
			name: "Glyph of Soulwell",
			description: "Your soulwell glows with an eerie light.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_shadesofdarkness.jpg",
		},
		[WarlockMinorGlyph.GlyphOfCrimsonBanish]: {
			name: "Glyph of Crimson Banish",
			description: "Your Banish spell is now red.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_cripple.jpg",
		},
		[WarlockMinorGlyph.GlyphOfGatewayAttunement]: {
			name: "Glyph of Gateway Attunement",
			description: "Demonic Gateways will automatically activate when you step near them.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_warlock_demonicportal_green.jpg",
		},
	},
};
