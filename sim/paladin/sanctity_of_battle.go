package paladin

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (paladin *Paladin) registerSanctityOfBattle() {
	var classMask int64
	if paladin.Spec == proto.Spec_SpecProtectionPaladin {
		classMask = SpellMaskSanctityOfBattleProt
	} else if paladin.Spec == proto.Spec_SpecHolyPaladin {
		classMask = SpellMaskSanctityOfBattleHoly
	} else {
		classMask = SpellMaskSanctityOfBattleRet
	}

	cooldownMod := paladin.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Multiplier,
		ClassMask: classMask,
	})

	updateFloatValue := func(castSpeed float64) {
		cooldownMod.UpdateFloatValue(castSpeed)
	}

	paladin.AddOnCastSpeedChanged(func(_ float64, castSpeed float64) {
		updateFloatValue(castSpeed)
	})

	core.MakePermanent(paladin.GetOrRegisterAura(core.Aura{
		Label:    "Sanctity of Battle",
		ActionID: core.ActionID{SpellID: 25956},

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			updateFloatValue(paladin.CastSpeed)
			cooldownMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			cooldownMod.Deactivate()
		},
	}))
}
