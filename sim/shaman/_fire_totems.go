package shaman

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
)

func searingTickCount(shaman *Shaman, offset float64) int32 {
	return int32(math.Ceil(40*(1.0+0.20*float64(shaman.Talents.TotemicFocus)+offset))) - 1
}

func (shaman *Shaman) registerSearingTotemSpell() {
	results := make([]*core.SpellResult, shaman.Env.GetNumTargets())

	shaman.SearingTotem = shaman.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 3599},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          SpellFlagTotem | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskSearingTotem,
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 5,
			PercentModifier: 100 - (15 * shaman.Talents.TotemicFocus) - shaman.GetMentalQuicknessBonus(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   shaman.DefaultCritMultiplier(),
		BonusCoefficient: 0.167,
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "SearingTotem",
			},
			// Actual searing totem cast in game is currently 1500 milliseconds with a slight random
			// delay inbetween each cast so using an extra 20 milliseconds to account for the delay
			// subtracting 1 tick so that it doesn't shoot after its actual expiration
			NumberOfTicks: searingTickCount(shaman, 0),
			TickLength:    time.Millisecond * (1500 + 20),
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				baseDamage := shaman.CalcAndRollDamageRange(sim, 0.09600000083, 0.30000001192)
				dot.Spell.CalcAndDealDamage(sim, target, baseDamage, dot.Spell.OutcomeMagicHitAndCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			shaman.MagmaTotem.AOEDot().Deactivate(sim)
			shaman.FireElemental.Disable(sim)
			if sim.CurrentTime < 0 {
				dropTime := sim.CurrentTime
				pa := &core.PendingAction{
					NextActionAt: 0,
					Priority:     core.ActionPriorityGCD,

					OnAction: func(sim *core.Simulation) {
						spell.Dot(sim.GetTargetUnit(0)).BaseTickCount = searingTickCount(shaman, dropTime.Minutes())
						spell.Dot(sim.GetTargetUnit(0)).Apply(sim)
					},
				}
				sim.AddPendingAction(pa)
			} else {
				spell.Dot(sim.GetTargetUnit(0)).BaseTickCount = searingTickCount(shaman, 0)
				spell.Dot(sim.GetTargetUnit(0)).Apply(sim)
			}
			duration := 60 * (1.0 + 0.20*float64(shaman.Talents.TotemicFocus))
			shaman.TotemExpirations[FireTotem] = sim.CurrentTime + time.Duration(duration)*time.Second
		},
	})
}

func (shaman *Shaman) registerMagmaTotemSpell() {
	shaman.MagmaTotem = shaman.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 8190},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAoE | SpellFlagTotem | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskMagmaTotem,
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 18,
			PercentModifier: 100 - (15 * shaman.Talents.TotemicFocus) - shaman.GetMentalQuicknessBonus(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   shaman.DefaultCritMultiplier(),

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "MagmaTotem",
			},
			NumberOfTicks:    int32(30 * (1.0 + 0.20*float64(shaman.Talents.TotemicFocus))),
			TickLength:       time.Second * 2,
			BonusCoefficient: 0.06700000167,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				baseDamage := shaman.CalcScalingSpellDmg(0.26699998975)
				for i, aoeTarget := range sim.Encounter.TargetUnits {
					results[i] = dot.Spell.CalcPeriodicDamage(sim, aoeTarget, baseDamage, dot.Spell.OutcomeMagicHitAndCrit)
				}
				for i := range sim.Encounter.TargetUnits {
					dot.Spell.DealPeriodicDamage(sim, results[i])
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			shaman.SearingTotem.Dot(shaman.CurrentTarget).Deactivate(sim)
			shaman.FireElemental.Disable(sim)
			spell.AOEDot().Apply(sim)

			duration := 60 * (1.0 + 0.20*float64(shaman.Talents.TotemicFocus))
			shaman.TotemExpirations[FireTotem] = sim.CurrentTime + time.Duration(duration)*time.Second
		},
	})
}
