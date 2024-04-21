package shaman

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

// TODO: Do the number of ticks needs to be updated when the talent changes? (20%/40% duration)
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
		BonusCoefficient: 0.2,
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
				dot.Spell.CalcAndDealDamage(sim, target, 90, dot.Spell.OutcomeMagicHitAndCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			shaman.MagmaTotem.AOEDot().Cancel(sim)
			shaman.FireElemental.Disable(sim)
			spell.Dot(sim.GetTargetUnit(0)).Apply(sim)
			duration := 60 * (1.0 + 0.20*float64(shaman.Talents.TotemicFocus))
			shaman.TotemExpirations[FireTotem] = sim.CurrentTime + time.Duration(duration)*time.Second
		},
	})
}

// TODO: Do the number of ticks needs to be updated when the talent changes? (20%/40% duration)
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
				baseDamage := 268 * sim.Encounter.AOECapMultiplier()
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					dot.Spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, dot.Spell.OutcomeMagicHitAndCrit)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			shaman.SearingTotem.Dot(shaman.CurrentTarget).Cancel(sim)
			shaman.FireElemental.Disable(sim)
			spell.AOEDot().Apply(sim)

			duration := 60 * (1.0 + 0.20*float64(shaman.Talents.TotemicFocus))
			shaman.TotemExpirations[FireTotem] = sim.CurrentTime + time.Duration(duration)*time.Second
		},
	})
}
