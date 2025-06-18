package death_knight

import "github.com/wowsims/mop/sim/core"

// T14 DPS
var ItemSetBattlegearOfTheLostCatacomb = core.NewItemSet(core.ItemSet{
	Name: "Battlegear of the Lost Catacomb",
	ID:   1123,
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(_ core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				ClassMask:  DeathKnightSpellFrostStrike | DeathKnightSpellObliterate | DeathKnightSpellScourgeStrike,
				FloatValue: 0.04,
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			dk := agent.(DeathKnightAgent).GetDeathKnight()

			// Handled in sim/core/buffs.go and sim/death_knight/frost/pillar_of_frost.go
			dk.T14Dps4pc = setBonusAura
		},
	},
})
