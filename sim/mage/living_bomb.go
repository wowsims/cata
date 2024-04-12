package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (mage *Mage) registerLivingBombSpell() {
	// Cata version has a cap of 3 active dots at once
	// Research this implementation
	livingBombExplosionSpell := mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 44461},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | HotStreakSpells,

		DamageMultiplierAdditive: 1 +
			.01*float64(mage.Talents.FirePower) +
			.05*float64(mage.Talents.CriticalMass) +
			core.TernaryFloat64(mage.HasPrimeGlyph(proto.MagePrimeGlyph_GlyphOfLivingBomb), .03, 0),
		CritMultiplier:   mage.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.5*mage.ScalingBaseDamage + 0.515*spell.SpellPower()
			baseDamage *= sim.Encounter.AOECapMultiplier()
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})

	onTick := func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
		dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
	}

	mage.LivingBomb = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 44457},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.17,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		DamageMultiplier: mage.GetFireMasteryBonusMultiplier(),
		DamageMultiplierAdditive: 1 +
			.01*float64(mage.Talents.FirePower) +
			.05*float64(mage.Talents.CriticalMass),
		CritMultiplier:   mage.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "LivingBomb",
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					livingBombExplosionSpell.Cast(sim, aura.Unit)
				},
			},

			NumberOfTicks:       4, // TODO haste breakpoint makes a new tick based on haste to stay between 10.5s and 13.5s 🥴
			TickLength:          time.Second * 3,
			AffectedByCastSpeed: true,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = 345 + 0.2*dot.Spell.SpellPower()
				dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: onTick,
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.SpellMetrics[target.UnitIndex].Hits--
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}
