package assassination

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/rogue"
)

var MutilateSpellID int32 = 1329

func (sinRogue *AssassinationRogue) newMutilateHitSpell(isMH bool) *core.Spell {
	actionID := core.ActionID{SpellID: MutilateSpellID, Tag: 1}
	procMask := core.ProcMaskMeleeMHSpecial
	if !isMH {
		actionID = core.ActionID{SpellID: MutilateSpellID, Tag: 2}
		procMask = core.ProcMaskMeleeOHSpecial
	}
	mutBaseDamage := sinRogue.GetBaseDamageFromCoefficient(0.25)
	mutWeaponPercent := 2.8

	return sinRogue.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       procMask,
		Flags:          core.SpellFlagMeleeMetrics | rogue.SpellFlagSealFate,
		ClassSpellMask: rogue.RogueSpellMutilateHit,

		DamageMultiplier:         mutWeaponPercent,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           sinRogue.CritMultiplier(true),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			var baseDamage float64
			if isMH {
				baseDamage = mutBaseDamage + spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			} else {
				baseDamage = mutBaseDamage + spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			}

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialBlockAndCrit)
		},
	})
}

func (sinRogue *AssassinationRogue) registerMutilateSpell() {
	sinRogue.MutilateMH = sinRogue.newMutilateHitSpell(true)
	sinRogue.MutilateOH = sinRogue.newMutilateHitSpell(false)

	sinRogue.Mutilate = sinRogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: MutilateSpellID, Tag: 0},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskEmpty, // Mutilate (Cast) no longer appears to proc anything
		Flags:          core.SpellFlagMeleeMetrics | rogue.SpellFlagBuilder | core.SpellFlagAPL,
		ClassSpellMask: rogue.RogueSpellMutilate,

		EnergyCost: core.EnergyCostOptions{
			Cost:   55,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:    time.Second,
				GCDMin: time.Millisecond * 700,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return sinRogue.HasDagger(core.MainHand) && sinRogue.HasDagger(core.OffHand)
		},

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			sinRogue.BreakStealth(sim)
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				sinRogue.AddComboPointsOrAnticipation(sim, 2, spell.ComboPointMetrics())
				sinRogue.MutilateOH.Cast(sim, target)
				sinRogue.MutilateMH.Cast(sim, target)
			} else {
				spell.IssueRefund(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}
