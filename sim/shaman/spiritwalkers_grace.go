package shaman

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (shaman *Shaman) spiritwalkersGraceBaseDuration() time.Duration {
	return 15 * time.Second
}

func (shaman *Shaman) registerSpiritwalkersGraceSpell() {
	actionID := core.ActionID{SpellID: 79206}

	spiritwalkersGraceAura := shaman.RegisterAura(core.Aura{
		Label:    "Spiritwalker's Grace" + shaman.Label,
		ActionID: actionID,
		Duration: 15 * time.Second,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_AllowCastWhileMoving,
		ClassMask: SpellMaskNone,
	})

	shaman.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolNature,
		Flags:          core.SpellFlagAPL | core.SpellFlagReadinessTrinket,
		ClassSpellMask: SpellMaskSpiritwalkersGrace,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 14.1,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Second * 120,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.RelatedSelfBuff.Activate(sim)
		},
		RelatedSelfBuff: spiritwalkersGraceAura,
	})
}
