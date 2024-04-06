package priest

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (priest *Priest) registerVampiricTouchSpell() {
	if !priest.Talents.VampiricTouch {
		return
	}

	replSrc := priest.Env.Raid.NewReplenishmentSource(core.ActionID{SpellID: 34914})

	priest.VampiricTouch = priest.RegisterSpell(PriestSpellVampiricTouch, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 34914},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.17,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
		},
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "VampiricTouch",
				OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if result.Landed() && spell == priest.MindBlast.Spell {
						priest.Env.Raid.ProcReplenishment(sim, replSrc)
					}
				},
			},
			NumberOfTicks:       5,
			TickLength:          time.Second * 3,
			AffectedByCastSpeed: true,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = priest.ScalingBaseDamage*0.101 + 0.352*dot.Spell.SpellPower()
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
				spell.Dot(target).Apply(sim)
			}

			spell.DealOutcome(sim, result)
		},
		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			if useSnapshot {
				dot := spell.Dot(target)
				return dot.CalcSnapshotDamage(sim, target, dot.OutcomeExpectedMagicSnapshotCrit)
			} else {
				baseDamage := priest.ScalingBaseDamage*0.101 + 0.352*spell.SpellPower()
				return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicCrit)
			}
		},
	})
}
