package priest

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

const PowerInfusionDuration = time.Second * 20
const PowerInfusionCD = time.Minute * 2

func (priest *Priest) registerPowerInfusionSpell() {
	if !priest.Talents.PowerInfusion {
		return
	}
	actionID := core.ActionID{SpellID: 10060}
	piAura := priest.GetOrRegisterAura(core.Aura{
		Label:    "PowerInfusion-Aura",
		ActionID: actionID,
		Duration: PowerInfusionDuration,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.05,
	}).AttachMultiplyCastSpeed(1.2)

	piAura.NewExclusiveEffect("ManaCost", true, core.ExclusiveEffect{
		Priority: -20,
		OnGain: func(ee *core.ExclusiveEffect, sim *core.Simulation) {
			ee.Aura.Unit.PseudoStats.SpellCostPercentModifier -= 20
		},
		OnExpire: func(ee *core.ExclusiveEffect, sim *core.Simulation) {
			ee.Aura.Unit.PseudoStats.SpellCostPercentModifier += 20
		},
	})

	piSpell := priest.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagHelpful,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: PowerInfusionCD,
			},
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, _ *core.Spell) {
			piAura.Activate(sim)
		},
	})

	priest.AddMajorCooldown(core.MajorCooldown{
		Spell:    piSpell,
		Priority: core.CooldownPriorityBloodlust,
		Type:     core.CooldownTypeMana,
	})
}
