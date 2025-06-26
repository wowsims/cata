package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

const (
	HurricaneBonusCoeff = 0.31
	HurricaneCoeff      = 0.31
)

func (druid *Druid) registerHurricaneSpell() {
	druid.HurricaneTickSpell = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 42231},
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellProc,
		Flags:          core.SpellFlagAoE,
		ClassSpellMask: DruidSpellHurricane,

		CritMultiplier:   druid.DefaultCritMultiplier(),
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: HurricaneBonusCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := druid.CalcScalingSpellDmg(HurricaneCoeff)

			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, damage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})

	druid.Hurricane = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 16914},
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagChanneled | core.SpellFlagAPL,
		ClassSpellMask: DruidSpellHurricane,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 50.3,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "Hurricane (Aura)",
			},
			NumberOfTicks:       10,
			TickLength:          time.Second * 1,
			AffectedByCastSpeed: true,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				druid.HurricaneTickSpell.Cast(sim, target)
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.AOEDot().Apply(sim)
		},
	})
}
