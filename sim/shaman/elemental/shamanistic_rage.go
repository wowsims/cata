package elemental

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/shaman"
)

func (ele *ElementalShaman) registerShamanisticRageSpell() {

	actionID := core.ActionID{SpellID: 30823}
	srAura := ele.RegisterAura(core.Aura{
		Label:    "Shamanistic Rage",
		ActionID: actionID,
		Duration: time.Second * 15,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct,
		FloatValue: -2,
	})

	spell := ele.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: shaman.SpellMaskShamanisticRage,
		Flags:          core.SpellFlagReadinessTrinket,
		Cast: core.CastConfig{
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    ele.NewTimer(),
				Duration: time.Minute * 1,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			srAura.Activate(sim)
		},
		RelatedSelfBuff: srAura,
	})

	ele.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeMana,
		ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
			return ele.CurrentManaPercent() < 0.05
		},
	})
}
