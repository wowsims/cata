package death_knight

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

// Incites a friendly party or raid member into a killing frenzy for 30 sec, increasing the target's melee and ranged haste by 20%, but causing them to lose health equal to 2% of their maximum health every 3 sec.
func (dk *DeathKnight) registerUnholyFrenzy() {
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
