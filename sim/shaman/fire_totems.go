package shaman

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (shaman *Shaman) registerSearingTotemSpell() {
	shaman.SearingTotem = shaman.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 3599},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          SpellFlagTotem | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskSearingTotem,
		ManaCost: core.ManaCostOptions{
			BaseCost:   0.05,
			Multiplier: 1 - 0.15*float64(shaman.Talents.TotemicFocus) - shaman.GetMentalQuicknessBonus(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   shaman.DefaultSpellCritMultiplier(),
		BonusCoefficient: 0.167,
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "SearingTotem",
			},
			// Actual searing totem cast in game is currently 1500 milliseconds with a slight random
			// delay inbetween each cast so using an extra 20 milliseconds to account for the delay
			// subtracting 1 tick so that it doesn't shoot after its actual expiration
			NumberOfTicks: int32(40*(1.0+0.20*float64(shaman.Talents.TotemicFocus))) - 1,
			TickLength:    time.Millisecond * (1500 + 20),
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				baseDamage := shaman.CalcAndRollDamageRange(sim, 0.09600000083, 0.30000001192)
				dot.Spell.CalcAndDealDamage(sim, target, baseDamage, dot.Spell.OutcomeMagicHitAndCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			shaman.MagmaTotem.AOEDot().Deactivate(sim)
			shaman.FireElemental.Disable(sim)
			spell.Dot(sim.GetTargetUnit(0)).Apply(sim)
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
		Flags:          SpellFlagTotem | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskMagmaTotem,
		ManaCost: core.ManaCostOptions{
			BaseCost:   0.18,
			Multiplier: 1 - 0.15*float64(shaman.Talents.TotemicFocus) - shaman.GetMentalQuicknessBonus(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   shaman.DefaultSpellCritMultiplier(),

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "MagmaTotem",
			},
			NumberOfTicks:    int32(30 * (1.0 + 0.20*float64(shaman.Talents.TotemicFocus))),
			TickLength:       time.Second * 2,
			BonusCoefficient: 0.08,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				results := make([]*core.SpellResult, shaman.Env.GetNumTargets())
				baseDamage := shaman.ClassSpellScaling * 0.26699998975
				aoeMult := sim.Encounter.AOECapMultiplier()
				dot.Spell.DamageMultiplier *= aoeMult
				for i, aoeTarget := range sim.Encounter.TargetUnits {
					results[i] = dot.Spell.CalcDamage(sim, aoeTarget, baseDamage, dot.Spell.OutcomeMagicHitAndCrit)
				}
				for i, _ := range sim.Encounter.TargetUnits {
					dot.Spell.DealDamage(sim, results[i])
				}
				dot.Spell.DamageMultiplier /= aoeMult
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
