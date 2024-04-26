package druid

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (druid *Druid) registerWrathSpell() {
	spellCoeff := 0.879
	wrathlMetric := druid.NewLunarEnergyMetrics(core.ActionID{SpellID: 5176})

	druid.Wrath = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 5176},
		SpellSchool:  core.SpellSchoolNature,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        SpellFlagOmenTrigger | core.SpellFlagAPL,
		MissileSpeed: 20,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.09,
			Multiplier: 1 - 0.03*float64(druid.Talents.Moonglow),
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond*2500 - time.Millisecond*100*time.Duration(druid.Talents.StarlightWrath),
			},
		},

		BonusCritRating: 2 * float64(druid.Talents.NaturesMajesty) * core.CritRatingPerCritChance,

		DamageMultiplier: 1 + core.TernaryFloat64(druid.HasPrimeGlyph(proto.DruidPrimeGlyph_GlyphOfWrath), 0.1, 0),

		CritMultiplier: druid.BalanceCritMultiplier(),

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(675, 761) + spellCoeff*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if result.Landed() {
				druid.AddEclipseEnergy(13+1.0/3.0, LunarEnergy, sim, wrathlMetric)

				spell.WaitTravelTime(sim, func(sim *core.Simulation) {
					spell.DealDamage(sim, result)
				})
			}

		},
	})
}
