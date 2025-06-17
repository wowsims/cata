package destruction

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/warlock"
)

func (destruction *DestructionWarlock) registerDarkSoulInstability() {
	buff := destruction.NewTemporaryStatsAura(
		"Dark Soul: Instability",
		core.ActionID{SpellID: 113858},
		stats.Stats{stats.CritRating: 30 * core.CritRatingPerCritPercent},
		time.Second*20,
	)

	spell := destruction.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 113858},
		DamageMultiplier: 1,
		ProcMask:         core.ProcMaskEmpty,
		SpellSchool:      core.SpellSchoolShadow,
		ClassSpellMask:   warlock.WarlockSpellDarkSoulInsanity,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{NonEmpty: true},
			CD: core.Cooldown{
				Timer:    destruction.NewTimer(),
				Duration: time.Minute * 2,
			},
		},
		RechargeTime: time.Minute * 2,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			buff.Activate(sim)
		},
		RelatedSelfBuff: buff.Aura,
	})
	destruction.AddMajorCooldown(core.MajorCooldown{
		Spell:    spell,
		BuffAura: buff,
		Priority: core.CooldownPriorityDefault,
		Type:     core.CooldownTypeDPS,
	})
}
