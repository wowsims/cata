package rogue

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (rogue *Rogue) applyPoisons() {
	rogue.applyDeadlyPoison()
	rogue.applyWoundPoison()
}

func (rogue *Rogue) registerPoisonAuras() {
	rogue.MasterPoisonerDebuffAuras = rogue.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		aura := core.MasterPoisonerDebuff(target)
		aura.Duration = core.NeverExpires
		return aura
	})
}

func (rogue *Rogue) GetLethalPoisonProcChance() float64 {
	return 0.3 + core.TernaryFloat64(rogue.Spec == proto.Spec_SpecAssassinationRogue, 0.2, 0)
}

func (rogue *Rogue) UpdateLethalPoisonPPH(bonusChance float64) {
	pph := rogue.GetLethalPoisonProcChance() + bonusChance
	rogue.deadlyPoisonPPHM = rogue.NewFixedProcChanceManager(pph, core.ProcMaskMelee)
	rogue.woundPoisonPPHM = rogue.NewFixedProcChanceManager(pph, core.ProcMaskMelee)
}

func (rogue *Rogue) registerDeadlyPoisonSpell() {
	dot_baseDamage := rogue.GetBaseDamageFromCoefficient(0.60000002384)
	dot_apScaling := 0.21299999952

	hit_baseDamage := rogue.GetBaseDamageFromCoefficient(0.31299999356)
	hit_apScaling := 0.10899999738
	hit_variance := 0.28000000119 * hit_baseDamage
	hit_minimum := hit_baseDamage - hit_variance/2

	// Register the hit as a distinct spell for Results UI separation
	dpHit := rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 2818, Tag: 2},
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellDamageProc,
		ClassSpellMask: RogueSpellDeadlyPoison,
		Flags:          core.SpellFlagPassiveSpell,

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           rogue.CritMultiplier(false),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := hit_minimum +
				sim.RandomFloat("Deadly Poison Hit")*hit_variance +
				hit_apScaling*spell.MeleeAttackPower()
			// DoT spell already checked if we hit, just send the damage
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
		},
	})

	rogue.DeadlyPoison = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 2818, Tag: 1},
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellDamageProc,
		ClassSpellMask: RogueSpellDeadlyPoison,
		Flags:          core.SpellFlagPassiveSpell,

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           rogue.CritMultiplier(false),
		ThreatMultiplier:         1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     "Deadly Poison",
				MaxStacks: 1,
				Duration:  time.Second * 12,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					rogue.MasterPoisonerDebuffAuras.Get(aura.Unit).Activate(sim)

				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					rogue.MasterPoisonerDebuffAuras.Get(aura.Unit).Deactivate(sim)

				},
			},
			NumberOfTicks: 4,
			TickLength:    time.Second * 3,

			OnSnapshot: func(_ *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = (dot_baseDamage + dot_apScaling*dot.Spell.MeleeAttackPower())
				attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
				dot.SnapshotCritChance = dot.Spell.PhysicalCritChance(attackTable)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable, true)
			},

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialNoBlockDodgeParry)
			if !result.Landed() {
				return
			}

			dot := spell.Dot(target)
			if dot.IsActive() {
				dpHit.Cast(sim, result.Target)
				dot.Refresh(sim)
				dot.TakeSnapshot(sim, false)
			} else {
				dot.Apply(sim)
				dot.Refresh(sim)
				dot.TakeSnapshot(sim, false)
			}
		},
	})
}

func (rogue *Rogue) registerWoundPoisonSpell() {
	woundPoisonDebuffAura := core.Aura{
		Label:    "Wound Poison",
		ActionID: core.ActionID{SpellID: 8680},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.MasterPoisonerDebuffAuras.Get(aura.Unit).Activate(sim)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.MasterPoisonerDebuffAuras.Get(aura.Unit).Deactivate(sim)
		},
	}

	rogue.WoundPoisonDebuffAuras = rogue.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.RegisterAura(woundPoisonDebuffAura)
	})

	wpBaseDamage := rogue.GetBaseDamageFromCoefficient(0.41699999571)
	apScaling := 0.11999999732
	variance := 0.28000000119 * wpBaseDamage
	minBaseDamage := wpBaseDamage - variance/2

	rogue.WoundPoison = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 8680},
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellDamageProc,
		ClassSpellMask: RogueSpellWoundPoison,

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           rogue.DefaultCritMultiplier(),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := minBaseDamage +
				sim.RandomFloat("Wound Poison")*variance +
				spell.MeleeAttackPower()*apScaling

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				rogue.WoundPoisonDebuffAuras.Get(target).Activate(sim)
			}
		},
	})
}

func (rogue *Rogue) applyDeadlyPoison() {
	if rogue.Options.LethalPoison == proto.RogueOptions_DeadlyPoison {
		rogue.UpdateLethalPoisonPPH(0)

		core.MakeProcTriggerAura(&rogue.Unit, core.ProcTrigger{
			Name:     "Deadly Poison",
			Outcome:  core.OutcomeLanded,
			Callback: core.CallbackOnSpellHitDealt,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procMask := core.Ternary(spell.SpellID == 86392, core.ProcMaskMeleeMH, spell.ProcMask)
				if rogue.deadlyPoisonPPHM.Proc(sim, procMask, "Deadly Poison") {
					rogue.DeadlyPoison.Cast(sim, result.Target)
				}
			},
		})
	}
}

func (rogue *Rogue) applyWoundPoison() {
	if rogue.Options.LethalPoison == proto.RogueOptions_WoundPoison {
		rogue.UpdateLethalPoisonPPH(0)

		core.MakeProcTriggerAura(&rogue.Unit, core.ProcTrigger{
			Name:     "Wound Poison",
			Outcome:  core.OutcomeLanded,
			Callback: core.CallbackOnSpellHitDealt,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procMask := core.Ternary(spell.SpellID == 86392, core.ProcMaskMeleeMH, spell.ProcMask)
				if rogue.woundPoisonPPHM.Proc(sim, procMask, "Wound Poison") {
					rogue.WoundPoison.Cast(sim, result.Target)
				}
			},
		})
	}
}
