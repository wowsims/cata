package balance

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/druid"
)

func (moonkin *BalanceDruid) registerAstralCommunionSpell() {
	actionID := core.ActionID{SpellID: 127663}

	channelTickLength := time.Second * 1
	numberOfTicks := 4
	eclipseEnergyGain := 25.0

	solarMetric := moonkin.NewSolarEnergyMetrics(actionID)
	lunarMetric := moonkin.NewLunarEnergyMetrics(actionID)

	moonkin.RegisterSpell(druid.Humanoid|druid.Moonkin, core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolArcane,
		Flags:          core.SpellFlagHelpful | core.SpellFlagChanneled | core.SpellFlagAPL,
		ClassSpellMask: druid.DruidSpellAstralCommunion,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{GCD: core.GCDDefault},
		},

		Hot: core.DotConfig{
			SelfOnly:            true,
			Aura:                core.Aura{Label: "Astral Communion"},
			NumberOfTicks:       int32(numberOfTicks),
			TickLength:          channelTickLength,
			AffectedByCastSpeed: false,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {

				if moonkin.CanGainEnergy(SolarAndLunarEnergy) {
					moonkin.AddEclipseEnergy(eclipseEnergyGain, LunarEnergy, sim, lunarMetric, dot.Spell)
				} else {
					moonkin.AddEclipseEnergy(eclipseEnergyGain, SolarEnergy, sim, solarMetric, dot.Spell)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.SelfHot().Apply(sim)
		},
	})
}
