package warlock

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (warlock *Warlock) registerDrainLife() {
	warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 689},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagChanneled | core.SpellFlagHauntSE | core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellDrainLife,

		ManaCost: core.ManaCostOptions{BaseCostPercent: 12},
		Cast:     core.CastConfig{DefaultCast: core.Cast{GCD: core.GCDDefault}},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           warlock.DefaultCritMultiplier(),
		ThreatMultiplier:         1,
		BonusCoefficient:         0.53899997473,

		Dot: core.DotConfig{
			Aura:                 core.Aura{Label: "Drain Life"},
			NumberOfTicks:        3,
			TickLength:           1 * time.Second,
			AffectedByCastSpeed:  true,
			HasteReducesDuration: true,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				baseDmg := warlock.CalcScalingSpellDmg(0.08550000191) + 0.12899999321*dot.Spell.SpellPower()
				dot.SnapshotBaseDamage = baseDmg * warlock.calcSoulSiphonMult(target)
				dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex], true)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			warlock.SoulBurnAura.Deactivate(sim)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
				spell.DealOutcome(sim, result)
			}
		},
	})
}
