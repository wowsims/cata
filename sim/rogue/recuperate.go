package rogue

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (rogue *Rogue) registerRecuperate() {
	actionID := core.ActionID{SpellID: 73651}

	energeticRecoveryAction := core.ActionID{SpellID: 79152}
	energeticRecoveryMetrics := rogue.NewEnergyMetrics(energeticRecoveryAction)

	rogue.Recuperate = rogue.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		Flags:        SpellFlagFinisher | core.SpellFlagAPL,
		MetricSplits: 6,

		EnergyCost: core.EnergyCostOptions{
			Cost: 30,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				spell.SetMetricsSplit(spell.Unit.ComboPoints())
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return rogue.ComboPoints() > 0
		},
		Hot: core.DotConfig{
			Aura: core.Aura{
				Label: "Recuperate",
			},
			NumberOfTicks:       0, // Decided at cast time
			TickLength:          time.Second * 3,
			AffectedByCastSpeed: false,
			BonusCoefficient:    1,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				rogue.RecuperateAura = dot.Aura
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				// Maybe implement Recup heal?
				healValue := rogue.MaxHealth() * (.03 + .005*float64(rogue.Talents.ImprovedRecuperate))
				dot.Spell.CalcAndDealPeriodicHealing(sim, target, healValue, dot.OutcomeTick)

				if rogue.Talents.EnergeticRecovery > 0 {
					energyRegen := float64(rogue.Talents.EnergeticRecovery) * 4.0
					// Trigger Energetic Recovery after small delay to prevent aura refresh loops
					// https://i.gyazo.com/dc845a371102294abfb207c6fd586bfa.png
					core.StartDelayedAction(sim, core.DelayedActionOptions{
						DoAt:     sim.CurrentTime + 1,
						Priority: core.ActionPriorityDOT,
						OnAction: func(s *core.Simulation) {
							rogue.AddEnergy(sim, energyRegen, energeticRecoveryMetrics)
						},
					})
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			aura := spell.Hot(spell.Unit)
			aura.Duration = time.Duration(rogue.ComboPoints()) * time.Second * 6
			aura.NumberOfTicks = rogue.ComboPoints() * 2
			aura.Activate(sim)
			rogue.ApplyFinisher(sim, spell)
		},
	})
}
