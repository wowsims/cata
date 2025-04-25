package arms

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warrior"
)

func (war *ArmsWarrior) RegisterDeadlyCalm() {
	if !war.Talents.DeadlyCalm {
		return
	}

	dcActionID := core.ActionID{SpellID: 85730}

	dcMod := war.AddDynamicMod(core.SpellModConfig{
		ClassMask: warrior.SpellMaskSpecialAttack,
		Kind:      core.SpellMod_PowerCost_Pct,
		IntValue:  -100,
	})

	dcAura := war.RegisterAura(core.Aura{
		Label:    "Deadly Calm",
		ActionID: dcActionID,
		Duration: time.Second * 10,
		Tag:      warrior.InnerRageExclusionTag,
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
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: time.Minute * 2,
			},
			SharedCD: core.Cooldown{
				Timer:    war.RecklessnessDeadlyCalmLock(),
				Duration: 10 * time.Second,
			},
		},
		ProcMask: core.ProcMaskEmpty,
		RageCost: core.RageCostOptions{Cost: 0},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !war.HasActiveAuraWithTag(warrior.InnerRageExclusionTag)
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
