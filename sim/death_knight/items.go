package death_knight

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

// T14 DPS
var ItemSetBattlegearOfTheLostCatacomb = core.NewItemSet(core.ItemSet{
	Name: "Battlegear of the Lost Catacomb",
	ID:   1123,
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(_ core.Agent, setBonusAura *core.Aura) {
			// Your Obliterate, Frost Strike, and Scourge Strike deal 4% increased damage.
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				ClassMask:  DeathKnightSpellFrostStrike | DeathKnightSpellObliterate | DeathKnightSpellScourgeStrike,
				FloatValue: 0.04,
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			// Your Pillar of Frost ability grants 5% additional Strength, and your Unholy Frenzy ability grants 10% additional haste.
			dk := agent.(DeathKnightAgent).GetDeathKnight()

			// Handled in sim/core/buffs.go and sim/death_knight/frost/pillar_of_frost.go
			dk.T14Dps4pc = setBonusAura
		},
	},
})

// T14 Tank
var PlateOfTheLostCatacomb = core.NewItemSet(core.ItemSet{
	Name: "Plate of the Lost Catacomb",
	ID:   1124,
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(_ core.Agent, setBonusAura *core.Aura) {
			// Reduces the cooldown of your Vampiric Blood ability by 20 sec.
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:      core.SpellMod_Cooldown_Flat,
				ClassMask: DeathKnightSpellVampiricBlood,
				TimeValue: time.Second * -20,
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			// Increases the healing received from your Death Strike by 10%.
			dk := agent.(DeathKnightAgent).GetDeathKnight()

			setBonusAura.AttachMultiplicativePseudoStatBuff(
				&dk.deathStrikeHealingMultiplier, 1.1,
			)
		},
	},
})
