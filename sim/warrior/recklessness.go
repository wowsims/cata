package warrior

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (warrior *Warrior) RegisterRecklessnessCD() {
	actionID := core.ActionID{SpellID: 1719}

	critMod := warrior.AddDynamicMod(core.SpellModConfig{
		ClassMask:  SpellMaskSpecialAttack,
		Kind:       core.SpellMod_BonusCrit_Rating,
		FloatValue: 50 * core.CritRatingPerCritChance,
	})

	warrior.RecklessnessAura = warrior.RegisterAura(core.Aura{
		Label:    "Recklessness",
		ActionID: actionID,
		Duration: time.Second * 12,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.DamageTakenMultiplier *= 1.2
			critMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.DamageTakenMultiplier /= 1.2
			critMod.Deactivate()
		},
	})

	reckSpell := warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL | core.SpellFlagMCD,
		ClassSpellMask: SpellMaskRecklessness,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Minute * 5,
			},
			IgnoreHaste: true,
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			if warrior.DeadlyCalmAura != nil {
				return !warrior.DeadlyCalmAura.IsActive()
			}
			return true
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			warrior.RecklessnessAura.Activate(sim)
		},
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: reckSpell,
		Type:  core.CooldownTypeDPS,
	})
}
