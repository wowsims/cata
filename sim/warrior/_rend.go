package warrior

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (warrior *Warrior) RegisterRendSpell() {
	baseTickDamage := warrior.ClassSpellScaling * 0.0939999968

	warrior.Rend = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 772},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagAPL | core.SpellFlagMeleeMetrics | SpellFlagBleed | core.SpellFlagPassiveSpell,
		ClassSpellMask: SpellMaskRend,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   10,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(BattleStance | DefensiveStance)
		},

		DamageMultiplier: 1.0,
		ThreatMultiplier: 1,
		CritMultiplier:   warrior.DefaultMeleeCritMultiplier(),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:    "Rend",
				ActionID: core.ActionID{SpellID: 94009},
				Tag:      "Rend",
			},
			NumberOfTicks: 5,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				weaponMH := warrior.AutoAttacks.MH()
				avgMHDamage := weaponMH.CalculateAverageWeaponDamage(dot.Spell.MeleeAttackPower())

				dot.SnapshotPhysical(target, baseTickDamage+0.25*avgMHDamage)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				dot := spell.Dot(target)

				// Rend ticks once on application
				dot.Apply(sim)
				dot.TickOnce(sim)
			} else {
				spell.IssueRefund(sim)
			}

			spell.DealOutcome(sim, result)
		},
	})
}
