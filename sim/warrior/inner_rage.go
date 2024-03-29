package warrior

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (warrior *Warrior) RegisterInnerRage() {
	actionID := core.ActionID{SpellID: 1134}
	warrior.InnerRageAura = warrior.RegisterAura(core.Aura{
		Label:    "Inner Rage",
		ActionID: actionID,
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.HeroicStrike.CostMultiplier *= 0.5
			warrior.Cleave.CostMultiplier *= 0.5
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.HeroicStrike.CostMultiplier /= 0.5
			warrior.Cleave.CostMultiplier /= 0.5
		},
	})

	ir := warrior.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		Flags:       core.SpellFlagHelpful | core.SpellFlagNoOnCastComplete | core.SpellFlagMCD | core.SpellFlagAPL,
		SpellSchool: core.SpellSchoolPhysical,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 0,
			},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 30,
			},
		},
		ThreatMultiplier: 0.0,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			warrior.InnerRageAura.Activate(sim)
		},
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: ir,
		Type:  core.CooldownTypeDPS,
	})
}
