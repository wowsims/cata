package warlock

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

const hellFireScale = 0.20999999344
const hellFireCoeff = 0.20999999344

func (warlock *Warlock) RegisterHellfire(callback WarlockSpellCastedCallback) *core.Spell {
	var baseDamage = warlock.CalcScalingSpellDmg(hellFireScale)
	results := make([]core.SpellResult, len(warlock.Env.Encounter.TargetUnits))

	hellfireActionID := core.ActionID{SpellID: 1949}
	manaMetric := warlock.NewManaMetrics(hellfireActionID)
	warlock.Hellfire = warlock.RegisterSpell(core.SpellConfig{
		ActionID:         hellfireActionID,
		SpellSchool:      core.SpellSchoolFire,
		Flags:            core.SpellFlagAoE | core.SpellFlagChanneled | core.SpellFlagAPL,
		ProcMask:         core.ProcMaskSpellDamage,
		ClassSpellMask:   WarlockSpellHellfire,
		ThreatMultiplier: 1,
		DamageMultiplier: 1,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		ManaCost: core.ManaCostOptions{BaseCostPercent: 2},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Hellfire",
			},

			IsAOE:                true,
			TickLength:           time.Second,
			NumberOfTicks:        14,
			HasteReducesDuration: true,
			AffectedByCastSpeed:  true,
			BonusCoefficient:     hellFireCoeff,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				for idx, unit := range sim.Encounter.TargetUnits {
					results[idx] = *dot.Spell.CalcAndDealPeriodicDamage(sim, unit, baseDamage, dot.Spell.OutcomeMagicHit)
				}

				warlock.SpendMana(sim, warlock.MaxMana()*0.02, manaMetric)
				if callback != nil {
					callback(results, dot.Spell, sim)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.AOEDot().Apply(sim)
		},
	})

	return warlock.Hellfire
}
