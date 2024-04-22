package warrior

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

// TODO (maybe) https://github.com/magey/wotlk-warrior/issues/23 - Rend is not benefitting from Two-Handed Weapon Specialization
func (warrior *Warrior) RegisterRendSpell() {
	dotTicks := int32(5)

	warrior.Rend = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 772},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagAPL | core.SpellFlagMeleeMetrics | SpellFlagBleed,
		ClassSpellMask: SpellMaskRend | SpellMaskSpecialAttack,

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
			NumberOfTicks: dotTicks,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				weaponMH := warrior.AutoAttacks.MH()
				avgMHDamage := weaponMH.CalculateAverageWeaponDamage(dot.Spell.MeleeAttackPower())

				ap := dot.Spell.MeleeAttackPower() / 14.0
				dot.Snapshot(target, (529+(0.25*6*(avgMHDamage+ap*weaponMH.SwingSpeed)))/float64(dot.NumberOfTicks))
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickPhysicalCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				dot := spell.Dot(target)

				// Rend ticks once on application, including on refreshes
				dot.Apply(sim)
				dot.TickOnce(sim)
			} else {
				spell.IssueRefund(sim)
			}

			spell.DealOutcome(sim, result)
		},
	})
}
