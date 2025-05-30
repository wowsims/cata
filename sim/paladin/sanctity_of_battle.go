package paladin

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

/*
Melee haste effects lower the cooldown and global cooldown of your

-- Holy Insight --
Holy Shock,
-- /Holy Insight --

Judgment, Crusader Strike,

-- Guarded by the Light --
Hammer of the Righteous, Consecration, Holy Wrath, Avenger's Shield, Shield of the Righteous
-- /Guarded of the Light --

-- Sword of Light --
Hammer of the Righteous, Exorcism
-- /Sword of Light --

and Hammer of Wrath.
*/
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

	updateFloatValue := func(attackSpeed float64) {
		cooldownMod.UpdateFloatValue(1 / attackSpeed)
	}

	paladin.AddOnMeleeAttackSpeedChanged(func(_ float64, attackSpeed float64) {
		updateFloatValue(attackSpeed)
	})

	core.MakePermanent(paladin.GetOrRegisterAura(core.Aura{
		Label:    "Sanctity of Battle",
		ActionID: core.ActionID{SpellID: 25956},

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			updateFloatValue(paladin.SwingSpeed())
			cooldownMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			cooldownMod.Deactivate()
		},
	}))
}
