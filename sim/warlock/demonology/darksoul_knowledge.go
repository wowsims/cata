package demonology

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/warlock"
)

func (demonology *DemonologyWarlock) registerDarksoulKnowledge() {
	buff := demonology.NewTemporaryStatsAura(
		"Dark Soul: Knowledge",
		core.ActionID{SpellID: 113858},
		stats.Stats{stats.MasteryRating: 30 * core.MasteryRatingPerMasteryPoint},
		time.Second*20,
	)

	spell := demonology.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 113861},
		DamageMultiplier: 1,
		ProcMask:         core.ProcMaskEmpty,
		SpellSchool:      core.SpellSchoolShadow,
		ClassSpellMask:   warlock.WarlockSpellDarkSoulKnowledge,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{NonEmpty: true},
			CD: core.Cooldown{
				Timer:    demonology.NewTimer(),
				Duration: time.Minute * 2,
			},
		},
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 5,
		},
		RechargeTime: time.Minute * 2,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			buff.Activate(sim)
		},
		RelatedSelfBuff: buff.Aura,
	})

	demonology.AddMajorCooldown(core.MajorCooldown{
		Spell:    spell,
		BuffAura: buff,
		Priority: core.CooldownPriorityDefault,
		Type:     core.CooldownTypeDPS,
	})
}
