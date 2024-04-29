package druid

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (druid *Druid) registerMoonfireSpell() {
	// TODO: Shooting stars proc on periodic damage
	// TODO: Glyph of Moonfire increase to periodic damage
	numTicks := druid.moonfireTicks()
	//hasMoonfireGlyph := druid.HasPrimeGlyph(proto.DruidPrimeGlyph_GlyphOfMoonfire)
	//bonusPeriodicDamageMultiplier := core.TernaryFloat64(hasMoonfireGlyph, 0.2, 0)

	druid.Moonfire = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 8921},
		SpellSchool:    core.SpellSchoolArcane,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: DruidSpellMoonfire,
		Flags:          SpellFlagNaturesGrace | SpellFlagOmenTrigger | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.21,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: 1,

		CritMultiplier:   druid.BalanceCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 0.18,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Moonfire",
			},
			NumberOfTicks:       druid.moonfireTicks(),
			TickLength:          time.Second * 2,
			AffectedByCastSpeed: true,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				// dot.Spell.DamageMultiplier = baseDamageMultiplier + bonusPeriodicDamageMultiplier
				// dot.SnapshotBaseDamage = 200 + 0.13*dot.Spell.SpellPower()
				// attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
				// dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
				// dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable, true)
				// dot.Spell.DamageMultiplier = baseDamageMultiplier - malusInitialDamageMultiplier
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			min, max := core.CalcScalingSpellEffectVarianceMinMax(proto.Class_ClassDruid, 0.221, 0.2)
			baseDamage := sim.Roll(min, max)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if result.Landed() {
				druid.ExtendingMoonfireStacks = 3
				dot := spell.Dot(target)
				dot.NumberOfTicks = numTicks
				dot.Apply(sim)
			}

			spell.DealDamage(sim, result)
		},
	})
}

func (druid *Druid) moonfireTicks() int32 {
	return 6
}
