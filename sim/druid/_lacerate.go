package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (druid *Druid) registerLacerateSpell() {
	tickDamage := 0.0700000003 * druid.ClassSpellScaling     // ~69
	initialDamage := 3.65700006485 * druid.ClassSpellScaling // ~3608

	initialDamageMul := 1.0
	// Set bonuses can scale up the ticks relative to the initial hit
	getTickDamageMultiplier := func() float64 { return core.TernaryFloat64(druid.T11Feral2pBonus.IsActive(), 1.1, 1) }

	druid.Lacerate = druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 33745},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagIgnoreResists,

		RageCost: core.RageCostOptions{
			Cost:   15,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		BonusCritPercent: core.TernaryFloat64(druid.HasPrimeGlyph(proto.DruidPrimeGlyph_GlyphOfLacerate), 5, 0),
		DamageMultiplier: initialDamageMul,
		CritMultiplier:   druid.DefaultCritMultiplier(),
		ThreatMultiplier: 1, // Changed in Cata
		MaxRange:         core.MaxMeleeRange,
		FlatThreatBonus:  0, // Removed in Cata

		Dot: core.DotConfig{
			Aura: druid.applyRendAndTear(core.Aura{
				Label:     "Lacerate",
				MaxStacks: 3,
				Duration:  time.Second * 15,
			}),
			NumberOfTicks: 5,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = tickDamage + 0.00369*dot.Spell.MeleeAttackPower()
				dot.SnapshotBaseDamage *= float64(dot.Aura.GetStacks())

				attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
				dot.Spell.DamageMultiplier = getTickDamageMultiplier()
				dot.SnapshotCritChance = dot.Spell.PhysicalCritChance(attackTable)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable, true)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)

				if (druid.BerserkProcAura != nil) && sim.Proc(0.5, "Berserk") {
					druid.BerserkProcAura.Activate(sim)
					druid.WaitUntil(sim, sim.CurrentTime+druid.ReactionTime)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := initialDamage + 0.0552*spell.MeleeAttackPower()
			if druid.BleedCategories.Get(target).AnyActive() {
				baseDamage *= 1.3
			}

			spell.DamageMultiplier = initialDamageMul
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				dot := spell.Dot(target)
				if dot.IsActive() {
					dot.Refresh(sim)
					dot.AddStack(sim)
					dot.TakeSnapshot(sim, false)
				} else {
					dot.Apply(sim)
					dot.SetStacks(sim, 1)
					dot.TakeSnapshot(sim, false)
				}
			} else {
				spell.IssueRefund(sim)
			}

			spell.DealDamage(sim, result)
		},
	})
}
