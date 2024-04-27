package arms

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/warrior"
)

func (war *ArmsWarrior) RegisterDeadlyCalm() {
	if !war.Talents.DeadlyCalm {
		return
	}

	dcActionID := core.ActionID{SpellID: 85730}

	dcMod := war.AddDynamicMod(core.SpellModConfig{
		ClassMask:  warrior.SpellMaskSpecialAttack,
		Kind:       core.SpellMod_PowerCost_Pct,
		FloatValue: -1,
	})

	war.DeadlyCalmAura = war.RegisterAura(core.Aura{
		Label:    "Deadly Calm",
		ActionID: dcActionID,
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dcMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dcMod.Deactivate()
		},
	})

	dc := war.RegisterSpell(core.SpellConfig{
		ActionID:       dcActionID,
		SpellSchool:    core.SpellSchoolPhysical,
		Flags:          core.SpellFlagAPL | core.SpellFlagMCD | core.SpellFlagNoOnDamageDealt | core.SpellFlagHelpful,
		ClassSpellMask: warrior.SpellMaskDeadlyCalm,

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
			return !war.InnerRageAura.IsActive() && !war.RecklessnessAura.IsActive()
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			war.DeadlyCalmAura.Activate(sim)
		},
	})

	war.AddMajorCooldown(core.MajorCooldown{
		Spell: dc,
		Type:  core.CooldownTypeDPS,
	})
}
