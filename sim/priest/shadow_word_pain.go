package priest

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (priest *Priest) registerShadowWordPainSpell() {
	priest.ShadowWordPain = priest.RegisterSpell(PriestSpellShadowWordPain, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 589},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.22,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "ShadowWordPain",
			},

			NumberOfTicks:       6,
			TickLength:          time.Second * 3,
			AffectedByCastSpeed: true,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = 194.709 + 0.161*dot.Spell.SpellPower()
				dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				if priest.Talents.Shadowform {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
				} else {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {

				// why? Was here before
				spell.SpellMetrics[target.UnitIndex].Hits--
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},

		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			if useSnapshot {
				dot := spell.Dot(target)
				return dot.CalcSnapshotDamage(sim, target, dot.OutcomeExpectedMagicSnapshotCrit)
			} else {
				baseDamage := 194.709 + 0.161*spell.SpellPower()
				return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicCrit)
			}
		},
	})
}
