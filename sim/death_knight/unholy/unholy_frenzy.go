package unholy

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

// Incites a friendly party or raid member into a killing frenzy for 30 sec, increasing the target's melee and ranged haste by 20%, but causing them to lose health equal to 2% of their maximum health every 3 sec.
func (uhdk *UnholyDeathKnight) registerUnholyFrenzy() {
	actionID := core.ActionID{SpellID: 49016, Tag: uhdk.Index}

	unholyFrenzyAuras := uhdk.NewAllyAuraArray(func(u *core.Unit) *core.Aura {
		if u.Type == core.PetUnit {
			return nil
		}

		return core.UnholyFrenzyAura(u, actionID.Tag, func() bool {
			return uhdk.T14Dps4pc.IsActive()
		})
	})
	unholyFrenzyTarget := uhdk.GetUnit(uhdk.Inputs.UnholyFrenzyTarget)
	if unholyFrenzyTarget == nil {
		unholyFrenzyTarget = &uhdk.Unit
	}

	if unholyFrenzyTarget == nil {
		return
	}

	unholyFrenzy := uhdk.Character.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL | core.SpellFlagHelpful | core.SpellFlagReadinessTrinket,
		ClassSpellMask: death_knight.DeathKnightSpellUnholyFrenzy,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    uhdk.NewTimer(),
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

	uhdk.AddMajorCooldown(core.MajorCooldown{
		Spell:    unholyFrenzy,
		Priority: core.CooldownPriorityBloodlust,
		Type:     core.CooldownTypeDPS,
	})
}
