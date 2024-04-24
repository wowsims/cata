package death_knight

import (
	//"github.com/wowsims/cata/sim/core/proto"

	//"time"

	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

func (dk *DeathKnight) ApplyBloodTalents() {
	// Butchery
	dk.applyButchery()

	// Blade barrier
	if dk.Talents.BladeBarrier > 0 {
		dk.PseudoStats.DamageTakenMultiplier *= 1.0 - 0.02*float64(dk.Talents.BladeBarrier)
	}

	if dk.Talents.ImprovedBloodTap > 0 {
		dk.AddStaticMod(core.SpellModConfig{
			Kind:      core.SpellMod_Cooldown_Flat,
			TimeValue: time.Second * -15 * time.Duration(dk.Talents.ImprovedBloodTap),
			ClassMask: DeathKnightSpellBloodTap,
		})
	}

	// Bladed Armor
	if dk.Talents.BladedArmor > 0 {
		coeff := float64(dk.Talents.BladedArmor) * 2
		dk.AddStatDependency(stats.Armor, stats.AttackPower, coeff/180.0)
	}

	// Scent of Blood
	dk.applyScentOfBlood()

	//Toughness
	if dk.Talents.Toughness > 0 {
		dk.ApplyEquipScaling(stats.Armor, []float64{1.0, 1.03, 1.06, 1.1}[dk.Talents.Toughness])
	}

	// Abomination's Might
	if dk.Talents.AbominationsMight > 0 {
		strengthCoeff := 0.01 * float64(dk.Talents.AbominationsMight)
		dk.MultiplyStat(stats.Strength, 1.0+strengthCoeff)
	}

	// Blood Parasite
	dk.applyBloodworms()

	// Will of the Necropolis
	dk.applyWillOfTheNecropolis()
}

func (dk *DeathKnight) applyWillOfTheNecropolis() {
	if dk.Talents.WillOfTheNecropolis == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 50150}
	wotnAura := dk.RegisterAura(core.Aura{
		Label:    "Will of The Necropolis",
		ActionID: actionID,
		Duration: core.NeverExpires,
	})

	damageMitigation := 1.0 - (0.05 * float64(dk.Talents.WillOfTheNecropolis))
	dk.AddDynamicDamageTakenModifier(func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
		if (dk.CurrentHealth()-result.Damage)/dk.MaxHealth() <= 0.35 {
			result.Damage *= damageMitigation
			if (dk.CurrentHealth()-result.Damage)/dk.MaxHealth() <= 0.35 {
				wotnAura.Activate(sim)
			}
			return
		}
		wotnAura.Deactivate(sim)
	})
}

func (dk *DeathKnight) applyScentOfBlood() {
	if dk.Talents.ScentOfBlood == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 49509}
	procChance := 0.15

	rpMetrics := dk.NewRunicPowerMetrics(actionID)

	procAura := dk.RegisterAura(core.Aura{
		Label:     "Scent of Blood Proc",
		ActionID:  actionID,
		Duration:  core.NeverExpires,
		MaxStacks: dk.Talents.ScentOfBlood,

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}

			dk.AddRunicPower(sim, 10.0, rpMetrics)
			aura.RemoveStack(sim)
		},
	})

	core.MakePermanent(dk.GetOrRegisterAura(core.Aura{
		Label: "Scent of Blood",
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if sim.Proc(procChance, "Scent Of Blood Proc Chance") {
				procAura.Activate(sim)
				procAura.SetStacks(sim, procAura.MaxStacks)
			}
		},
	}))
}

func (dk *DeathKnight) applyButchery() {
	if dk.Talents.Butchery == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 49483}
	amountOfRunicPower := 1.0 * float64(dk.Talents.Butchery)
	rpMetrics := dk.NewRunicPowerMetrics(actionID)

	var pa *core.PendingAction
	core.MakePermanent(dk.RegisterAura(core.Aura{
		ActionID: actionID,
		Label:    "Butchery",
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			pa = core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period: time.Second * 5,
				OnAction: func(sim *core.Simulation) {
					dk.AddRunicPower(sim, amountOfRunicPower, rpMetrics)
				},
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			pa.Cancel(sim)
		},
	}))
}

func (dk *DeathKnight) applyBloodworms() {
	// if dk.Talents.BloodParasite == 0 {
	// 	return
	// }

	// procChance := 0.03 * float64(dk.Talents.BloodParasite)
	// icd := core.Cooldown{
	// 	Timer:    dk.NewTimer(),
	// 	Duration: time.Second * 20,
	// }

	// // For tracking purposes
	// procSpell := dk.RegisterSpell(core.SpellConfig{
	// 	ActionID: core.ActionID{SpellID: 49543},
	// 	ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
	// 		// Summon Bloodworms
	// 		random := int(math.Round(sim.RandomFloat("Bloodworms count")*2.0)) + 2
	// 		for i := 0; i < random; i++ {
	// 			dk.Bloodworm[i].EnableWithTimeout(sim, dk.Bloodworm[i], time.Second*20)
	// 			dk.Bloodworm[i].CancelGCDTimer(sim)
	// 		}
	// 	},
	// })

	// core.MakePermanent(dk.RegisterAura(core.Aura{
	// 	Label: "Bloodworms Proc",
	// 	OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	// 		if !spell.ProcMask.Matches(core.ProcMaskMelee) {
	// 			return
	// 		}

	// 		if !icd.IsReady(sim) {
	// 			return
	// 		}

	// 		if sim.RandomFloat("Bloodworms proc") < procChance {
	// 			icd.Use(sim)
	// 			procSpell.Cast(sim, result.Target)
	// 		}
	// 	},
	// }))
}
