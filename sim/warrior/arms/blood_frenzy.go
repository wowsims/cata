package arms

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/warrior"
)

func (war *ArmsWarrior) applyBloodFrenzy() {
	if war.Talents.BloodFrenzy == 0 {
		return
	}

	bfAuras := war.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.BloodFrenzyAura(target, war.Talents.BloodFrenzy)
	})

	// Trauma is also applied by the Blood Frenzy talent in Cata
	traumaAuras := war.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.TraumaAura(target, war.Talents.BloodFrenzy)
	})

	bfRageProc := core.ActionID{SpellID: 92576}
	procChance := 0.05 * float64(war.Talents.BloodFrenzy)
	bfRageMetrics := war.NewRageMetrics(bfRageProc)
	core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:       "Blood Frenzy Rage Trigger",
		Callback:   core.CallbackOnSpellHitDealt,
		Outcome:    core.OutcomeLanded,
		ProcMask:   core.ProcMaskMeleeMHAuto,
		ProcChance: procChance,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			war.AddRage(sim, 20, bfRageMetrics)
		},
	})
	core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:       "Blood Frenzy/Trauma Debuff Trigger",
		Callback:   core.CallbackOnSpellHitDealt,
		SpellFlags: warrior.SpellFlagBleed,
		Outcome:    core.OutcomeLanded,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			dot := spell.Dot(result.Target)

			// Apply Blood Frenzy, it lasts as long as the dot is on the target
			bf := bfAuras.Get(result.Target)
			bf.Duration = dot.TickLength * time.Duration(dot.NumberOfTicks)
			bf.Activate(sim)

			// Apply Trauma, has fixed duration regardless of bleeds
			trauma := traumaAuras.Get(result.Target)
			trauma.Activate(sim)
		},
	})
}
