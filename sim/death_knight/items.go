package death_knight

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

// TODO: T12
// TODO: T13

var ItemSetMagmaPlatedBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Magma Plated Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the critical strike chance of your Death Coil and Frost Strike abilities by 5%.
			agent.GetCharacter().AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_BonusCrit_Rating,
				ClassMask:  DeathKnightSpellDeathCoil | DeathKnightSpellFrostStrike,
				FloatValue: 5 * core.CritRatingPerCritChance,
			})
		},
		4: func(agent core.Agent) {
			// Each time you gain a Death Rune, you also gain 1% increased attack power for 30 sec. Stacks up to 3 times.
			// Also activated whenever KM procs
			character := agent.GetCharacter()

			apDep := make([]*stats.StatDependency, 3)
			for i := 1; i <= 3; i++ {
				apDep[i-1] = character.NewDynamicMultiplyStat(stats.AttackPower, 1.0+float64(i)*0.01)
			}

			aura := character.GetOrRegisterAura(core.Aura{
				Label:     "Death Eater",
				ActionID:  core.ActionID{SpellID: 90507},
				Duration:  time.Second * 30,
				MaxStacks: 3,
				OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
					if oldStacks > 0 {
						character.DisableDynamicStatDep(sim, apDep[oldStacks-1])
					}
					if newStacks > 0 {
						character.EnableDynamicStatDep(sim, apDep[newStacks-1])
					}
				},
			})

			core.MakeProcTriggerAura(&agent.GetCharacter().Unit, core.ProcTrigger{
				Name:           "Magma Plated Battlegear",
				Callback:       core.CallbackOnCastComplete,
				ClassSpellMask: DeathKnightSpellConvertToDeathRune | DeathKnightSpellKillingMachine,
				ICD:            time.Millisecond * 10, // Batch together double rune converts
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					aura.Activate(sim)
					aura.AddStack(sim)
				},
			})
		},
	},
})

var ItemSetMagmaPlatedBattlearmor = core.NewItemSet(core.ItemSet{
	Name: "Magma Plated Battlearmor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the damage done by your Death Strike ability by 5%.
			agent.GetCharacter().AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Flat,
				ClassMask:  DeathKnightSpellDeathStrike,
				FloatValue: 0.05,
			})
		},
		4: func(agent core.Agent) {
			// Increases the duration of your Icebound Fortitude ability by 50%.
			// Implemented in icebound_fortitude.go
		},
	},
})

func init() {
}
