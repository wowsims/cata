package warrior

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

// TODO (maybe) https://github.com/magey/wotlk-warrior/issues/23 - Rend is not benefitting from Two-Handed Weapon Specialization
func (warrior *Warrior) RegisterRendSpell() {
	dotTicks := int32(5)

	warrior.Rend = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 47465},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskRend,

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

		ThreatMultiplier: 1,
		CritMultiplier:   warrior.DefaultMeleeCritMultiplier(),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Rend",
				Tag:   "Rend",
			},
			NumberOfTicks:       dotTicks,
			TickLength:          time.Second * 3,
			AffectedByCastSpeed: true,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				avgWeaponDamage := warrior.AutoAttacks.MH().CalculateAverageWeaponDamage(dot.Spell.MeleeAttackPower())
				ap := dot.Spell.MeleeAttackPower() / 14.0
				dot.SnapshotBaseDamage = 529 + (0.25 * 6 * (avgWeaponDamage + ap*warrior.AutoAttacks.MH().SwingSpeed))

				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex], true)
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
