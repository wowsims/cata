package mage

import (
	"github.com/wowsims/mop/sim/core"
)

// T14
var ItemSetRegaliaOfTheBurningScroll = core.NewItemSet(core.ItemSet{
	Name: "Regalia of the Burning Scroll",
	Bonuses: map[int32]core.ApplySetBonus{
		// Increases the damage done by your Arcane Missiles spell by 7%, increases the damage done by your Pyroblast spell by 8%, and increases the damage done by your Ice Lance spell by 12%.
		2: func(_ core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				ClassMask:  MageSpellIceLance,
				FloatValue: 0.12,
			}).AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				ClassMask:  MageSpellArcaneMissilesTick,
				FloatValue: 0.07,
			}).AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				ClassMask:  MageSpellPyroblast | MageSpellPyroblastDot,
				FloatValue: 0.08,
			})
			setBonusAura.ExposeToAPL(123097)
		},
		// Increases the damage bonus of Arcane Power by an additional 10%, reduces the cooldown of Icy Veins by 50%, and reduces the cooldown of Combustion by 20%.
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			mage := agent.(MageAgent).GetMage()
			mage.T14_4pc = setBonusAura

			setBonusAura.AttachSpellMod(core.SpellModConfig{
				FloatValue: 0.5,
				Kind:       core.SpellMod_Cooldown_Multiplier,
				ClassMask:  MageSpellIcyVeins,
			}).AttachSpellMod(core.SpellModConfig{
				FloatValue: 1 - 0.2,
				Kind:       core.SpellMod_Cooldown_Multiplier,
				ClassMask:  MageSpellCombustion,
			})
			setBonusAura.ExposeToAPL(123101)
		},
	},
})
