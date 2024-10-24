package druid

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (druid *Druid) registerWildMushrooms() {

	wildMushroomsStackAura := druid.GetOrRegisterAura(core.Aura{
		Label:     "Wild Mushroom Stacks",
		ActionID:  core.ActionID{SpellID: 88747},
		Duration:  core.NeverExpires,
		MaxStacks: 3,
	})

	druid.WildMushrooms = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID: core.ActionID{SpellID: 88747},
		Flags:    core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.11,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDMin,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			wildMushroomsStackAura.Activate(sim)
			wildMushroomsStackAura.AddStack(sim)
		},
	})

	wildMushroomsDamage := druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 78777},
		SpellSchool:      core.SpellSchoolNature,
		Flags:            core.SpellFlagPassiveSpell,
		ProcMask:         core.ProcMaskSpellDamage,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		ClassSpellMask:   DruidSpellWildMushroomDetonate,
		CritMultiplier:   druid.DefaultSpellCritMultiplier(),
		BonusCoefficient: 0.6032,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			min, max := core.CalcScalingSpellEffectVarianceMinMax(proto.Class_ClassDruid, 0.9464, 0.19)
			baseDamage := sim.Roll(min, max)
			baseDamage *= sim.Encounter.AOECapMultiplier()

			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})

	druid.WildMushroomsDetonate = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 88751},
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: DruidSpellWildMushroomDetonate,
		Flags:          core.SpellFlagAPL | SpellFlagOmenTrigger | core.SpellFlagPassiveSpell,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second * 10,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   druid.DefaultSpellCritMultiplier(),
		BonusCoefficient: 0.6032,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {

			for i := wildMushroomsStackAura.GetStacks(); i > 0; i-- {
				wildMushroomsDamage.Cast(sim, target)
				if !spell.ProcMask.Matches(core.ProcMaskSpellProc) {
					wildMushroomsStackAura.RemoveStack(sim)
				}
			}
		},
	})
}
