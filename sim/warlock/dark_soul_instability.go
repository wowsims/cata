package warlock

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

func (warlock *Warlock) registerDarkSoulInstability() {
	buff := warlock.NewTemporaryStatsAura(
		"Dark Soul: Instability",
		core.ActionID{SpellID: 113858},
		stats.Stats{stats.CritRating: 30 * core.CritRatingPerCritPercent},
		time.Second*20,
	)

	spell := warlock.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 113858},
		DamageMultiplier: 1,
		ProcMask:         core.ProcMaskEmpty,
		SpellSchool:      core.SpellSchoolShadow,
		ClassSpellMask:   WarlockSpellDarkSoulInsanity,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{NonEmpty: true},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Minute * 2,
			},
		},
		RechargeTime: time.Minute * 2,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			buff.Activate(sim)
		},
	})

	warlock.AddMajorCooldown(core.MajorCooldown{
		Spell:    spell,
		BuffAura: buff,
		Priority: core.CooldownPriorityDefault,
		Type:     core.CooldownTypeDPS,
	})
}
