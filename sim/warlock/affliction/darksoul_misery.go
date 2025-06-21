package affliction

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

func (affliction *AfflictionWarlock) registerDarkSoulMisery() {
	buff := affliction.RegisterAura(core.Aura{
		Label:    "Dark Soul: Misery",
		ActionID: core.ActionID{SpellID: 113860},
		Duration: time.Second * 20,
	}).AttachMultiplyCastSpeed(1.3)

	spell := affliction.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 113860},
		DamageMultiplier: 1,
		ProcMask:         core.ProcMaskEmpty,
		SpellSchool:      core.SpellSchoolShadow,
		ClassSpellMask:   warlock.WarlockSpellDarkSoulMisery,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{NonEmpty: true},
			CD: core.Cooldown{
				Timer:    affliction.NewTimer(),
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
		RelatedSelfBuff: buff,
	})

	affliction.AddMajorCooldown(core.MajorCooldown{
		Spell:    spell,
		Priority: core.CooldownPriorityDefault,
		Type:     core.CooldownTypeDPS,
	})
}
