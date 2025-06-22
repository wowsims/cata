package balance

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/druid"
)

const (
	EnergyGainPerTick           = 25.0
	EnergyGainPerTickDuringSotF = 100.0
)

func (moonkin *BalanceDruid) registerAstralCommunionSpell() {
	actionID := core.ActionID{SpellID: 127663}

	eclipseEnergyGain := EnergyGainPerTick

	solarMetric := moonkin.NewSolarEnergyMetrics(actionID)
	lunarMetric := moonkin.NewLunarEnergyMetrics(actionID)

	moonkin.AstralCommunion = moonkin.RegisterSpell(druid.Humanoid|druid.Moonkin, core.SpellConfig{
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
			NumberOfTicks:       4,
			TickLength:          time.Second * 1,
			AffectedByCastSpeed: false,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				if moonkin.CanGainEnergy(SolarAndLunarEnergy) {
					moonkin.AddEclipseEnergy(eclipseEnergyGain, LunarEnergy, sim, lunarMetric, dot.Spell)
				} else {
					moonkin.AddEclipseEnergy(eclipseEnergyGain, SolarEnergy, sim, solarMetric, dot.Spell)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.SelfHot().Apply(sim)

			if moonkin.AstralInsight.IsActive() {
				eclipseEnergyGain = EnergyGainPerTickDuringSotF

				spell.SelfHot().TickOnce(sim)
				spell.SelfHot().Deactivate(sim)

				eclipseEnergyGain = EnergyGainPerTick
				moonkin.AstralInsight.Deactivate(sim)
			}
		},
	})

	moonkin.AddEclipseCallback(func(_ Eclipse, gained bool, sim *core.Simulation) {
		if gained && moonkin.AstralCommunion.SelfHot().IsActive() {
			moonkin.AstralCommunion.SelfHot().Deactivate(sim)
		}
	})
}
