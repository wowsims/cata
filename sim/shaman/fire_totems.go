package shaman

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

// TODO: Updated for cata. Not sure if the coefficient is changed. Need to check.
//
//	Also not sure if the number of ticks needs to be updated when the talent (20%/40% duration)
func (shaman *Shaman) registerSearingTotemSpell() {
	shaman.SearingTotem = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 3599},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       SpellFlagTotem | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.05,
			Multiplier: 1 - 0.15*float64(shaman.Talents.TotemicFocus) - shaman.GetMentalQuicknessBonus(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
		},

		DamageMultiplier: 1 + float64(shaman.Talents.CallOfFlame)*0.1,
		CritMultiplier:   shaman.ElementalFuryCritMultiplier(0),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "SearingTotem",
			},
			// These are the real tick values, but searing totem doesn't start its next
			// cast until the previous missile hits the target. We don't have an option
			// for target distance yet so just pretend the tick rate is lower.
			// https://wotlk.wowhead.com/spell=25530/attack
			//NumberOfTicks:        30,
			//TickLength:           time.Second * 2.2,
			NumberOfTicks: 24,
			TickLength:    time.Second * 60 / 24,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				baseDamage := 90 + 0.167*dot.Spell.SpellPower()
				result := dot.Spell.CalcAndDealDamage(sim, target, baseDamage, dot.Spell.OutcomeMagicHitAndCrit)

				// TODO: Just getting this ready. No idea how this works... the tooltip says this dot stacks up to five times. Similar to deadly poison?
				if shaman.Talents.SearingFlames > 0 && result.Landed() {
					if shaman.Talents.SearingFlames == 3 || sim.RandomFloat("Searing Flames") < 0.33*float64(shaman.Talents.SearingFlames) {
						searingFlamesDot := shaman.SearingFlamesDot.Dot(target)
						searingFlamesDot.SnapshotBaseDamage = result.Damage / float64(dot.NumberOfTicks)
						searingFlamesDot.SnapshotAttackerMultiplier = 1
						searingFlamesDot.Spell.Cast(sim, target)
					}
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			shaman.MagmaTotem.AOEDot().Cancel(sim)
			shaman.FireElemental.Disable(sim)
			spell.Dot(sim.GetTargetUnit(0)).Apply(sim)

			bonusDuration := 1.0 + 0.20*float64(shaman.Talents.TotemicFocus)
			// +1 needed because of rounding issues with totem tick time.
			shaman.TotemExpirations[FireTotem] = sim.CurrentTime + time.Second*60*time.Duration(bonusDuration) + 1
		},
	})
}

// TODO: Updated for cata. Not sure if the coefficient is changed. Need to check.
//
//	Also not sure if the number of ticks needs to be updated when the talent (20%/40% duration)
func (shaman *Shaman) registerMagmaTotemSpell() {
	shaman.MagmaTotem = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 8190},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       SpellFlagTotem | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.18,
			Multiplier: 1 - 0.15*float64(shaman.Talents.TotemicFocus) - shaman.GetMentalQuicknessBonus(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
		},

		DamageMultiplier: 1 + float64(shaman.Talents.CallOfFlame)*0.1,
		CritMultiplier:   shaman.ElementalFuryCritMultiplier(0),

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "MagmaTotem",
			},
			NumberOfTicks: 10,
			TickLength:    time.Second * 2,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				baseDamage := 268 + 0.1*dot.Spell.SpellPower()
				baseDamage *= sim.Encounter.AOECapMultiplier()
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					dot.Spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, dot.Spell.OutcomeMagicHitAndCrit)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			shaman.SearingTotem.Dot(shaman.CurrentTarget).Cancel(sim)
			shaman.FireElemental.Disable(sim)
			spell.AOEDot().Apply(sim)
			bonusDuration := 1.0 + 0.20*float64(shaman.Talents.TotemicFocus)
			// +1 needed because of rounding issues with totem tick time.
			shaman.TotemExpirations[FireTotem] = sim.CurrentTime + time.Second*60*time.Duration(bonusDuration) + 1
		},
	})
}
