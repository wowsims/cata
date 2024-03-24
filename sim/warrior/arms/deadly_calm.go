package arms

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (war *ArmsWarrior) RegisterDeadlyCalm() {
	if !war.Talents.DeadlyCalm {
		return
	}

	dcActionID := core.ActionID{SpellID: 85730}
	dcAura := war.RegisterAura(core.Aura{
		Label:    "Deadly Calm",
		ActionID: dcActionID,
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			war.PseudoStats.CostMultiplier -= 1
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			war.PseudoStats.CostMultiplier += 1
		},
	})

	dc := war.RegisterSpell(core.SpellConfig{
		ActionID:    dcActionID,
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagAPL | core.SpellFlagMCD | core.SpellFlagNoOnDamageDealt | core.SpellFlagHelpful,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 0,
			},
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: time.Minute * 2,
			},
		},
		ProcMask: core.ProcMaskEmpty,
		RageCost: core.RageCostOptions{Cost: 0},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return war.InnerRage == nil || !war.InnerRage.IsActive()
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dcAura.Activate(sim)
		},
	})

	war.AddMajorCooldown(core.MajorCooldown{
		Spell: dc,
		Type:  core.CooldownTypeDPS,
	})
}
