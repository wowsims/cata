package shadow

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/priest"
)

const swdScale = 2.392
const swdCoeff = 2.157

func (shadow *ShadowPriest) registerShadowWordDeathSpell() {
	actionId := core.ActionID{SpellID: 32379}
	swdAura := shadow.RegisterAura(core.Aura{
		Label:    "Shadow Word: Death",
		ActionID: actionId.WithTag(1),
		Duration: 9 * time.Second,
	})

	shadow.RegisterSpell(core.SpellConfig{
		ActionID:       actionId,
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: priest.PriestSpellShadowWordDeath,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 2.6,
		},

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           shadow.DefaultCritMultiplier(),
		ThreatMultiplier:         1,
		BonusCoefficient:         swdCoeff,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shadow.NewTimer(),
				Duration: time.Second * 8,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, shadow.CalcScalingSpellDmg(swdScale), spell.OutcomeMagicHitAndCrit)
			if swdAura.IsActive() {
				swdAura.Deactivate(sim)
				return
			}

			shadow.ShadowOrbs.Gain(1, actionId, sim)
			swdAura.Activate(sim)
			spell.CD.Reset()
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return sim.IsExecutePhase20()
		},
	})
}
