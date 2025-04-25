package priest

import (
	"github.com/wowsims/mop/sim/core"
)

func (priest *Priest) registerPowerInfusionSpell() {
	if !priest.Talents.PowerInfusion {
		return
	}

	actionID := core.ActionID{SpellID: 10060, Tag: priest.Index}

	powerInfusionTarget := priest.GetUnit(priest.SelfBuffs.PowerInfusionTarget)
	powerInfusionAuras := priest.NewAllyAuraArray(func(unit *core.Unit) *core.Aura {
		if unit.Type == core.PetUnit {
			return nil
		}
		return core.PowerInfusionAura(unit, actionID.Tag)
	})

	piSpell := priest.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagHelpful,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 16,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: core.PowerInfusionCD,
			},
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, _ *core.Spell) {
			if powerInfusionTarget != nil {
				powerInfusionAuras.Get(powerInfusionTarget).Activate(sim)
			} else {
				powerInfusionAuras.Get(target).Activate(sim)
			}
		},
	})

	priest.AddMajorCooldown(core.MajorCooldown{
		Spell:    piSpell,
		Priority: core.CooldownPriorityBloodlust,
		Type:     core.CooldownTypeMana,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			// How can we determine the target will be able to continue casting
			// 	for the next 15s at 20% reduced mana cost? Arbitrary value until then.
			//if powerInfusionTarget.CurrentMana() < 3000 {
			//	return false
			//}
			return powerInfusionTarget != nil && !powerInfusionTarget.HasActiveAuraWithTag(core.BloodlustAuraTag)
		},
	})
}
