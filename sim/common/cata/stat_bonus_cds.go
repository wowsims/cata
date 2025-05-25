package cata

import (
	"time"

	"github.com/wowsims/cata/sim/common/shared"
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

func init() {

	// HASTE
	shared.NewHasteActive(67152, 617, time.Second*20, time.Minute*2) // Lady La-La's Singing Shell

	// Shard of Woe
	core.NewSimpleStatOffensiveTrinketEffectWithOtherEffects(60233, stats.Stats{stats.HasteRating: 1935}, time.Second*10, time.Minute*1, func(agent core.Agent, _ proto.ItemLevelState) {
		character := agent.GetCharacter()

		aura := core.MakePermanent(character.RegisterAura(core.Aura{
			Label:    "Spell Cost Reduction",
			ActionID: core.ActionID{ItemID: 91171},
		}).AttachSpellMod(core.SpellModConfig{
			Kind:         core.SpellMod_PowerCost_Flat,
			ResourceType: proto.ResourceType_ResourceTypeMana,
			IntValue:     -205,
		}))

		character.ItemSwap.RegisterProc(60233, aura)
	})
	// CRIT
	shared.NewSimpleStatActive(66879) // Bottled Lightning
	shared.NewSimpleStatActive(61448) // Oremantle's Favor
	shared.NewSimpleStatActive(70144) // Ricket's Magnetic Fireball

	// STRENGTH
	shared.NewSimpleStatActive(59685) // Kvaldir Battle Standard - Alliance
	shared.NewSimpleStatActive(59689) // Kvaldir Battle Standard - Horde
	shared.NewSimpleStatActive(55251) // Might of the Ocean
	shared.NewSimpleStatActive(56285) // Might of the Ocean (Heroic)
	shared.NewSimpleStatActive(55814) // Magnetite Mirror
	shared.NewSimpleStatActive(56345) // Magnetite Mirror (Heroic)
	shared.NewSimpleStatActive(66994) // Soul's Anguish
	shared.NewSimpleStatActive(52351) // Figurine - King of Boars
	shared.NewSimpleStatActive(64689) // Bloodthirsty Gladiator's Badge of Victory
	shared.NewSimpleStatActive(62464) // Impatience of Youth (Alliance)
	shared.NewSimpleStatActive(62469) // Impatience of Youth (Horde)
	shared.NewSimpleStatActive(61034) // Vicious Gladiator's Badge of Victory - 365
	shared.NewSimpleStatActive(70519) // Vicious Gladiator's Badge of Victory - 371
	shared.NewSimpleStatActive(70400) // Ruthless Gladiator's Badge of Victory - 384
	shared.NewSimpleStatActive(72450) // Ruthless Gladiator's Badge of Victory - 390
	shared.NewSimpleStatActive(73496) // Cataclysmic Gladiator's Badge of Victory
	shared.NewSimpleStatActive(69002) // Essence of the Eternal Flame
	shared.NewSimpleStatActive(77116) // Rotting Skull 397 - Valor Points

	// AGILITY
	shared.NewSimpleStatActive(63840) // Juju of Nimbleness
	shared.NewSimpleStatActive(63843) // Blood-Soaked Ale Mug
	shared.NewSimpleStatActive(64687) // Bloodthirsty Gladiator's Badge of Conquest
	shared.NewSimpleStatActive(52199) // Figurine - Demon Panther
	shared.NewSimpleStatActive(62468) // Unsolvable Riddle (Alliance)
	shared.NewSimpleStatActive(62463) // Unsolvable Riddle (Horde)
	shared.NewSimpleStatActive(68709) // Unsolvable Riddle (No Faction)
	shared.NewSimpleStatActive(61033) // Vicious Gladiator's Badge of Conquest - 365
	shared.NewSimpleStatActive(70517) // Vicious Gladiator's Badge of Conquest - 371
	shared.NewSimpleStatActive(70399) // Ruthless Gladiator's Badge of Conquest - 384
	shared.NewSimpleStatActive(72304) // Ruthless Gladiator's Badge of Conquest - 390
	shared.NewSimpleStatActive(73648) // Cataclysmic Gladiator's Badge of Conquest
	shared.NewSimpleStatActive(69001) // Ancient Petrified Seed
	shared.NewSimpleStatActive(77113) // Kiroptyric Sigil 397 - Valor Points

	// SPIRIT
	shared.NewSimpleStatActive(67101) // Unquenchable Flame
	shared.NewSimpleStatActive(52354) // Figurine - Dream Owl
	shared.NewSimpleStatActive(58184) // Core of Ripeness

	// DODGE
	shared.NewSimpleStatActive(67037)  // Binding Promise
	shared.NewSimpleStatActive(52352)  // Figurine - Earthen Guardian
	shared.NewSimpleStatActive(59515)  // Vial of Stolen Memories
	shared.NewSimpleStatActive(65109)  // Vial of Stolen Memories (Heroic)
	shared.NewSimpleStatActive(70143)  // Moonwell Phial
	shared.NewSimpleStatActive(232015) // Brawler's Trophy
	shared.NewSimpleStatActive(77117)  // Fire of the Deep 397 - Valor Points

	// SPELLPOWER
	shared.NewSimpleStatActive(61429) // Insignia of the Earthen Lord
	shared.NewSimpleStatActive(55256) // Sea Star
	shared.NewSimpleStatActive(56290) // Sea Star (Heroic)
	shared.NewSimpleStatActive(52353) // Figurine - Jeweled Serpent
	shared.NewSimpleStatActive(64688) // Bloodthirsty Gladiator's Badge of Dominance
	shared.NewSimpleStatActive(58183) // Soul Casket
	shared.NewSimpleStatActive(61035) // Vicious Gladiator's Badge of Dominance - 365
	shared.NewSimpleStatActive(70518) // Vicious Gladiator's Badge of Dominance - 371
	shared.NewSimpleStatActive(70401) // Ruthless Gladiator's Badge of Dominance - 384
	shared.NewSimpleStatActive(72448) // Ruthless Gladiator's Badge of Dominance - 390
	shared.NewSimpleStatActive(73498) // Cataclysmic Gladiator's Badge of Dominance
	shared.NewSimpleStatActive(77114) // Bottled Wishes 397 - Valor Points
	shared.NewSimpleStatActive(77115) // Reflection of the Light 397 - Valor Points

	// HEALTH
	shared.NewSimpleStatActive(61433) // Insignia of Diplomacy
	shared.NewSimpleStatActive(55845) // Heart of Thunder
	shared.NewSimpleStatActive(56370) // Heart of Thunder
	shared.NewSimpleStatActive(64740) // Bloodthirsty Gladiator's Emblem of Cruelty
	shared.NewSimpleStatActive(64741) // Bloodthirsty Gladiator's Emblem of Meditation
	shared.NewSimpleStatActive(64742) // Bloodthirsty Gladiator's Emblem of Tenacity
	shared.NewSimpleStatActive(62048) // Darkmoon Card: Earthquake
	shared.NewSimpleStatActive(61026) // Vicious Gladiator's Emblem of Cruelty
	shared.NewSimpleStatActive(61028) // Vicious Gladiator's Emblem of Alacrity
	shared.NewSimpleStatActive(61029) // Vicious Gladiator's Emblem of Prowess
	shared.NewSimpleStatActive(61032) // Vicious Gladiator's Emblem of Tenacity
	shared.NewSimpleStatActive(61030) // Vicious Gladiator's Emblem of Proficiency
	shared.NewSimpleStatActive(61027) // Vicious Gladiator's Emblem of Accuracy

	// INT
	shared.NewSimpleStatActive(67118) // Electrospark Heartstarter
	shared.NewSimpleStatActive(68998) // Rune of Zeth
	shared.NewSimpleStatActive(69000) // Fiery Quintessence

	// MASTERY
	shared.NewSimpleStatActive(63745) // Za'brox's Lucky Tooth - Alliance
	shared.NewSimpleStatActive(63742) // Za'brox's Lucky Tooth - Horde
	shared.NewSimpleStatActive(56115) // Skardyn's Grace
	shared.NewSimpleStatActive(56440) // Skardyn's Grace (Heroic)
	shared.NewSimpleStatActive(56132) // Mark of Khardros
	shared.NewSimpleStatActive(56458) // Mark of Khardros (Heroic)
	shared.NewSimpleStatActive(70142) // Moonwell Chalice

	// PARRY
	shared.NewSimpleStatActive(55881) // Impetuous Query
	shared.NewSimpleStatActive(56406) // Impetuous Query (Heroic)

	// RESISTANCE
	shared.NewSimpleStatActive(62466) // Mirror of Broken Images (Alliance)
	shared.NewSimpleStatActive(62471) // Mirror of Broken Images (Horde)

}
