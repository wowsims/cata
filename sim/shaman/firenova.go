package shaman

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (shaman *Shaman) registerFireNovaSpell() {
	shaman.FireNova = shaman.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 1535},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          SpellFlagFocusable | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskFireNova,
		ManaCost: core.ManaCostOptions{
			BaseCost:   0.22,
			Multiplier: 1 - 0.05*float64(shaman.Talents.Convection) - shaman.GetMentalQuicknessBonus(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Second * time.Duration(4),
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   shaman.ElementalFuryCritMultiplier(0),
		ThreatMultiplier: 1,
		BonusCoefficient: 0.164,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				if aoeTarget.GetAura("FlameShock") != nil {
					for _, newTarget := range sim.Encounter.TargetUnits {
						if newTarget != aoeTarget {
							spell.CalcAndDealDamage(sim, newTarget, 789, spell.OutcomeMagicHitAndCrit)
						}
					}
				}
			}
		},
	})
}

func (shaman *Shaman) IsFireNovaCastable(sim *core.Simulation) bool {
	if shaman.FireNova.IsReady(sim) {
		for _, aoeTarget := range sim.Encounter.TargetUnits {
			if aoeTarget.GetAura("FlameShock") != nil {
				return true
			}
		}
	}

	return false
}
