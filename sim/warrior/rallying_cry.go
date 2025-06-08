package warrior

import (
	"github.com/wowsims/mop/sim/core"
)

func (war *Warrior) registerRallyingCry() {
	war.RallyingCryAura = core.RallyingCryAura(&war.Character, war.Index)

	spell := war.RegisterSpell(core.SpellConfig{
		ActionID:       core.RallyingCryActionID,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskRallyingCry,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: core.RallyingCryCD,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return (war.LastStandAura != nil && !war.LastStandAura.IsActive()) || war.LastStandAura == nil
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			war.RallyingCryAura.Activate(sim)
		},
	})

	war.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeSurvival,
	})
}
