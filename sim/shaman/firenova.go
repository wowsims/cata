package shaman

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

// TODO: Numbers updated for cata. Need to check spell coefficient
func (shaman *Shaman) registerFireNovaSpell() {
	shaman.FireNova = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 1535},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagFocusable | core.SpellFlagAPL,

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

		DamageMultiplier: 1 + float64(shaman.Talents.CallOfFlame)*0.1,
		CritMultiplier:   shaman.ElementalFuryCritMultiplier(0),
		ThreatMultiplier: 0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := 789 + 0.2142*spell.SpellPower()
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				if aoeTarget.GetAura("FlameShock") != nil {
					for _, newTarget := range sim.Encounter.TargetUnits {
						if newTarget != aoeTarget {
							// TODO: Uncomment this
							//baseDamage *= sim.Encounter.AOECapMultiplier()
							spell.CalcAndDealDamage(sim, newTarget, damage, spell.OutcomeMagicHitAndCrit)
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
