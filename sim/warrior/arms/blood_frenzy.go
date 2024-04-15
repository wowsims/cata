package arms

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/warrior"
)

func (war *ArmsWarrior) RegisterBloodFrenzy() {
	if war.Talents.BloodFrenzy == 0 {
		return
	}

	bfAuras := war.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.BloodFrenzyAura(target, war.Talents.BloodFrenzy)
	})

	// Trauma is also applied by the Blood Frenzy talent in Cata
	traumaAuras := war.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.TraumaAura(target, int(war.Talents.BloodFrenzy))
	})

	bfRageProc := core.ActionID{SpellID: 92576}
	procChance := 0.05 * float64(war.Talents.BloodFrenzy)
	bfRageMetrics := war.NewRageMetrics(bfRageProc)
	core.MakePermanent(war.RegisterAura(core.Aura{
		Label: "Blood Frenzy Buff Trigger",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			if spell.ActionID.IsOtherAction(proto.OtherAction_OtherActionAttack) {
				if sim.RandomFloat("Blood Frenzy Rage Proc") < procChance {
					war.AddRage(sim, 20, bfRageMetrics)
				}
			}

			if spell.Flags.Matches(warrior.SpellFlagBleed) {
				dot := spell.Dot(result.Target)

				// Apply Blood Frenzy, it lasts as long as the dot is on the target
				bf := bfAuras.Get(result.Target)
				bf.Duration = dot.TickLength * time.Duration(dot.NumberOfTicks)
				bf.Activate(sim)

				// Apply Trauma, has fixed duration regardless of bleeds
				trauma := traumaAuras.Get(result.Target)
				trauma.Activate(sim)
			}
		},
	}))
}
