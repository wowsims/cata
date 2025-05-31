package cata

import (
	"time"

	"github.com/wowsims/mop/sim/common/shared"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func init() {
	// HASTE
	shared.NewHasteActive(67152, 617, time.Second*20, time.Minute*2) // Lady La-La's Singing Shell

	// Shard of Woe
	core.NewSimpleStatOffensiveTrinketEffectWithOtherEffects(60233, stats.Stats{stats.HasteRating: 1935}, time.Second*10, time.Minute*1, func(agent core.Agent) {
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
	shared.NewCritActive(66879, 512, time.Second*20, time.Minute*2)  // Bottled Lightning
	shared.NewCritActive(61448, 970, time.Second*20, time.Minute*2)  // Oremantle's Favor
	shared.NewCritActive(70144, 1700, time.Second*20, time.Minute*2) // Ricket's Magnetic Fireball

	// STRENGTH
	shared.NewStrengthActive(59685, 830, time.Second*10, time.Minute)     // Kvaldir Battle Standard - Alliance
	shared.NewStrengthActive(59689, 830, time.Second*10, time.Minute)     // Kvaldir Battle Standard - Horde
	shared.NewStrengthActive(55251, 765, time.Second*15, time.Second*90)  // Might of the Ocean
	shared.NewStrengthActive(56285, 1425, time.Second*15, time.Second*90) // Might of the Ocean (Heroic)
	shared.NewStrengthActive(55814, 1075, time.Second*15, time.Second*90) // Magnetite Mirror
	shared.NewStrengthActive(56345, 1425, time.Second*15, time.Second*90) // Magnetite Mirror (Heroic)
	shared.NewStrengthActive(66994, 765, time.Second*15, time.Second*90)  // Soul's Anguish
	shared.NewStrengthActive(52351, 1425, time.Second*20, time.Minute*2)  // Figurine - King of Boars
	shared.NewStrengthActive(64689, 1520, time.Second*20, time.Minute*2)  // Bloodthirsty Gladiator's Badge of Victory
	shared.NewStrengthActive(62464, 1605, time.Second*20, time.Minute*2)  // Impatience of Youth (Alliance)
	shared.NewStrengthActive(62469, 1605, time.Second*20, time.Minute*2)  // Impatience of Youth (Horde)
	shared.NewStrengthActive(61034, 1605, time.Second*20, time.Minute*2)  // Vicious Gladiator's Badge of Victory - 365
	shared.NewStrengthActive(70519, 1794, time.Second*20, time.Minute*2)  // Vicious Gladiator's Badge of Victory - 371
	shared.NewStrengthActive(70400, 2029, time.Second*20, time.Minute*2)  // Ruthless Gladiator's Badge of Victory - 384
	shared.NewStrengthActive(72450, 2144, time.Second*20, time.Minute*2)  // Ruthless Gladiator's Badge of Victory - 390
	shared.NewStrengthActive(73496, 2419, time.Second*20, time.Minute*2)  // Cataclysmic Gladiator's Badge of Victory
	shared.NewStrengthActive(69002, 1277, time.Second*15, time.Minute)    // Essence of the Eternal Flame
	shared.NewStrengthActive(77116, 2290, time.Second*15, time.Second*90) // Rotting Skull 397 - Valor Points

	// AGILITY
	shared.NewAgilityActive(63840, 1095, time.Second*15, time.Second*90) // Juju of Nimbleness
	shared.NewAgilityActive(63843, 1095, time.Second*15, time.Second*90) // Blood-Soaked Ale Mug
	shared.NewAgilityActive(64687, 1520, time.Second*20, time.Second*90) // Bloodthirsty Gladiator's Badge of Conquest
	shared.NewAgilityActive(52199, 1425, time.Second*20, time.Minute*2)  // Figurine - Demon Panther
	shared.NewAgilityActive(62468, 1605, time.Second*20, time.Minute*2)  // Unsolvable Riddle (Alliance)
	shared.NewAgilityActive(62463, 1605, time.Second*20, time.Minute*2)  // Unsolvable Riddle (Horde)
	shared.NewAgilityActive(68709, 1605, time.Second*20, time.Minute*2)  // Unsolvable Riddle (No Faction)
	shared.NewAgilityActive(61033, 1605, time.Second*20, time.Minute*2)  // Vicious Gladiator's Badge of Conquest - 365
	shared.NewAgilityActive(70517, 1794, time.Second*20, time.Minute*2)  // Vicious Gladiator's Badge of Conquest - 371
	shared.NewAgilityActive(70399, 2029, time.Second*20, time.Minute*2)  // Ruthless Gladiator's Badge of Conquest - 384
	shared.NewAgilityActive(72304, 2144, time.Second*20, time.Minute*2)  // Ruthless Gladiator's Badge of Conquest - 390
	shared.NewAgilityActive(73648, 2419, time.Second*20, time.Minute*2)  // Cataclysmic Gladiator's Badge of Conquest
	shared.NewAgilityActive(69001, 1277, time.Second*15, time.Minute)    // Ancient Petrified Seed
	shared.NewAgilityActive(77113, 2290, time.Second*15, time.Second*90) // Kiroptyric Sigil 397 - Valor Points

	// SPIRIT
	shared.NewSpiritActive(67101, 555, time.Second*20, time.Minute*2)  // Unquenchable Flame
	shared.NewSpiritActive(52354, 1425, time.Second*20, time.Minute*2) // Figurine - Dream Owl
	shared.NewSpiritActive(58184, 1926, time.Second*20, time.Minute*2) // Core of Ripeness

	// DODGE
	shared.NewDodgeActive(67037, 512, time.Second*20, time.Minute*2)   // Binding Promise
	shared.NewDodgeActive(52352, 1425, time.Second*20, time.Minute*2)  // Figurine - Earthen Guardian
	shared.NewDodgeActive(59515, 1605, time.Second*20, time.Minute*2)  // Vial of Stolen Memories
	shared.NewDodgeActive(65109, 1812, time.Second*20, time.Minute*2)  // Vial of Stolen Memories (Heroic)
	shared.NewDodgeActive(70143, 1700, time.Second*20, time.Minute*2)  // Moonwell Phial
	shared.NewDodgeActive(232015, 1520, time.Second*20, time.Minute*2) // Brawler's Trophy
	shared.NewDodgeActive(77117, 2290, time.Second*15, time.Second*90) // Fire of the Deep 397 - Valor Points

	// SpellPower
	shared.NewSpellPowerActive(61429, 970, time.Second*15, time.Second*90)  // Insignia of the Earthen Lord
	shared.NewSpellPowerActive(55256, 765, time.Second*20, time.Minute*2)   // Sea Star
	shared.NewSpellPowerActive(56290, 1425, time.Second*20, time.Minute*2)  // Sea Star (Heroic)
	shared.NewSpellPowerActive(52353, 1425, time.Second*20, time.Minute*2)  // Figurine - Jeweled Serpent
	shared.NewSpellPowerActive(64688, 1520, time.Second*20, time.Minute*2)  // Bloodthirsty Gladiator's Badge of Dominance
	shared.NewSpellPowerActive(58183, 1926, time.Second*20, time.Minute*2)  // Soul Casket
	shared.NewSpellPowerActive(61035, 1605, time.Second*20, time.Minute*2)  // Vicious Gladiator's Badge of Dominance - 365
	shared.NewSpellPowerActive(70518, 1794, time.Second*20, time.Minute*2)  // Vicious Gladiator's Badge of Dominance - 371
	shared.NewSpellPowerActive(70401, 2029, time.Second*20, time.Minute*2)  // Ruthless Gladiator's Badge of Dominance - 384
	shared.NewSpellPowerActive(72448, 2144, time.Second*20, time.Minute*2)  // Ruthless Gladiator's Badge of Dominance - 390
	shared.NewSpellPowerActive(73498, 2419, time.Second*20, time.Minute*2)  // Cataclysmic Gladiator's Badge of Dominance
	shared.NewSpellPowerActive(77114, 2290, time.Second*15, time.Second*90) // Bottled Wishes 397 - Valor Points
	shared.NewSpellPowerActive(77115, 2290, time.Second*15, time.Second*90) // Reflection of the Light 397 - Valor Points

	// HEALTH
	shared.NewHealthActive(61433, 6985, time.Second*15, time.Minute*3)  // Insignia of Diplomacy
	shared.NewHealthActive(55845, 7740, time.Second*15, time.Minute*3)  // Heart of Thunder
	shared.NewHealthActive(56370, 10260, time.Second*15, time.Minute*3) // Heart of Thunder
	shared.NewHealthActive(64740, 15315, time.Second*15, time.Minute*2) // Bloodthirsty Gladiator's Emblem of Cruelty
	shared.NewHealthActive(64741, 15315, time.Second*15, time.Minute*2) // Bloodthirsty Gladiator's Emblem of Meditation
	shared.NewHealthActive(64742, 15315, time.Second*15, time.Minute*2) // Bloodthirsty Gladiator's Emblem of Tenacity
	shared.NewHealthActive(62048, 15500, time.Second*15, time.Minute*2) // Darkmoon Card: Earthquake
	shared.NewHealthActive(61026, 16196, time.Second*15, time.Minute*2) // Vicious Gladiator's Emblem of Cruelty
	shared.NewHealthActive(61028, 16196, time.Second*15, time.Minute*2) // Vicious Gladiator's Emblem of Alacrity
	shared.NewHealthActive(61029, 16196, time.Second*15, time.Minute*2) // Vicious Gladiator's Emblem of Prowess
	shared.NewHealthActive(61032, 16196, time.Second*15, time.Minute*2) // Vicious Gladiator's Emblem of Tenacity
	shared.NewHealthActive(61030, 16196, time.Second*15, time.Minute*2) // Vicious Gladiator's Emblem of Proficiency
	shared.NewHealthActive(61027, 16196, time.Second*15, time.Minute*2) // Vicious Gladiator's Emblem of Accuracy

	// INT
	shared.NewIntActive(67118, 567, time.Second*20, time.Minute*2)   // Electrospark Heartstarter
	shared.NewIntActive(68998, 1277, time.Second*15, time.Minute)    // Rune of Zeth
	shared.NewIntActive(69000, 1149, time.Second*25, time.Second*90) // Fiery Quintessence

	// MASTERY
	shared.NewMasteryActive(63745, 1095, time.Second*15, time.Second*90) // Za'brox's Lucky Tooth - Alliance
	shared.NewMasteryActive(63742, 1095, time.Second*15, time.Second*90) // Za'brox's Lucky Tooth - Horde
	shared.NewMasteryActive(56115, 1260, time.Second*20, time.Minute*2)  // Skardyn's Grace
	shared.NewMasteryActive(56440, 1425, time.Second*20, time.Minute*2)  // Skardyn's Grace (Heroic)
	shared.NewMasteryActive(56132, 1260, time.Second*15, time.Second*90) // Mark of Khardros
	shared.NewMasteryActive(56458, 1425, time.Second*15, time.Second*90) // Mark of Khardros (Heroic)
	shared.NewMasteryActive(70142, 1700, time.Second*20, time.Minute*2)  // Moonwell Chalice

	// PARRY
	shared.NewParryActive(55881, 1260, time.Second*10, time.Minute) // Impetuous Query
	shared.NewParryActive(56406, 1425, time.Second*10, time.Minute) // Impetuous Query (Heroic)
}
