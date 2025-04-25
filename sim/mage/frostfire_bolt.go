package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (mage *Mage) registerFrostfireBoltSpell() {

	hasGlyph := mage.HasPrimeGlyph(proto.MagePrimeGlyph_GlyphOfFrostfire)

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 44614},
		SpellSchool:    core.SpellSchoolFire | core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellFrostfireBolt,
		MissileSpeed:   28,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 9,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2500,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   mage.DefaultSpellCritMultiplier(),
		BonusCoefficient: 0.977,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     "FrostfireBolt",
				MaxStacks: 3,
				Duration:  time.Second * 12,
			},
			NumberOfTicks:       4,
			TickLength:          time.Second * 3,
			AffectedByCastSpeed: true,
			BonusCoefficient:    0.00733,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.Snapshot(target, 0.00712*mage.ClassSpellScaling)
				dot.SnapshotBaseDamage = dot.SnapshotBaseDamage * float64(dot.Aura.GetStacks())
			},

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.949 * mage.ClassSpellScaling
			// Not sure if double dipping exists in Cata. Removed for now.
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				if result.Landed() {
					dot := spell.Dot(target)
					if hasGlyph && dot != nil {
						if dot.IsActive() {
							dot.Refresh(sim)
							dot.AddStack(sim)
						} else {
							dot.Apply(sim)
						}
						dot.TakeSnapshot(sim, true)
					}
					spell.DealDamage(sim, result)
				}
			})
		},
	})
}
