package brewmaster

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/monk"
)

func (bm *BrewmasterMonk) registerAvertHarm() {
	actionID := core.ActionID{SpellID: 115213}
	duration := 6 * time.Second

	bm.AvertHarmAura = core.MakeProcTriggerAura(&bm.Unit, core.ProcTrigger{
		Name:     "Avert Harm" + bm.Label,
		ActionID: actionID,
		Duration: duration,
		Outcome:  core.OutcomeHit,
		Callback: core.CallbackOnSpellHitTaken,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.RelatedSelfBuff != nil && result.Target.CurrentHealthPercent() <= 0.1 {
				spell.RelatedSelfBuff.Deactivate(sim)
			}
		},
	})

	spell := bm.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,
		ClassSpellMask: monk.MonkSpellAvertHarm,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    bm.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return bm.StanceMatches(monk.SturdyOx)
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			bm.AvertHarmAura.Activate(sim)
		},
		RelatedSelfBuff: bm.AvertHarmAura,
	})

	bm.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeSurvival,
	})

}
