package death_knight

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (dk *DeathKnight) registerUnholyFrenzySpell() {
	if !dk.Talents.UnholyFrenzy {
		return
	}

	actionID := core.ActionID{SpellID: 49016, Tag: dk.Index}

	unholyFrenzyAuras := dk.NewAllyAuraArray(func(u *core.Unit) *core.Aura {
		if u.Type == core.PetUnit {
			return nil
		}
		return core.UnholyFrenzyAura(u, actionID.Tag)
	})
	unholyFrenzyTarget := dk.GetUnit(dk.Inputs.UnholyFrenzyTarget)
	if unholyFrenzyTarget == nil {
		unholyFrenzyTarget = &dk.Unit
	}

	if unholyFrenzyTarget == nil {
		return
	}

	unholyFrenzy := dk.Character.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL | core.SpellFlagHelpful,
		ClassSpellMask: DeathKnightSpellUnholyFrenzy,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, _ *core.Spell) {
			if unholyFrenzyTarget != nil {
				unholyFrenzyAuras.Get(unholyFrenzyTarget).Activate(sim)
			} else if target.Type == core.PlayerUnit {
				unholyFrenzyAuras.Get(target).Activate(sim)
			}
		},
	})

	dk.AddMajorCooldown(core.MajorCooldown{
		Spell:    unholyFrenzy,
		Priority: core.CooldownPriorityBloodlust,
		Type:     core.CooldownTypeDPS,
	})
}
