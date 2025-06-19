package paladin

import (
	"time"

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
	var cooldownMask int64
	var gcdMask int64
	if paladin.Spec == proto.Spec_SpecProtectionPaladin {
		cooldownMask = SpellMaskSanctityOfBattleProt
		gcdMask = SpellMaskSanctityOfBattleProtGcd
	} else if paladin.Spec == proto.Spec_SpecHolyPaladin {
		cooldownMask = SpellMaskSanctityOfBattleHoly
		gcdMask = SpellMaskSanctityOfBattleHolyGcd
	} else {
		cooldownMask = SpellMaskSanctityOfBattleRet
		gcdMask = SpellMaskSanctityOfBattleRetGcd
	}

	cooldownMod := paladin.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Multiplier,
		ClassMask: cooldownMask,
	})

	gcdMod := paladin.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_GlobalCooldown_Flat,
		ClassMask: gcdMask,
	})

	updateFloatValue := func(meleeHaste float64) {
		multiplier := 1 / meleeHaste
		cooldownMod.UpdateFloatValue(multiplier)
		gcdMod.UpdateTimeValue(-(core.DurationFromSeconds(min(0.5, 1.5-1.5*multiplier)).Round(time.Millisecond)))
	}

	paladin.AddOnMeleeAndRangedHasteChanged(func(_ float64, meleeHaste float64) {
		updateFloatValue(meleeHaste)
	})

	core.MakePermanent(paladin.GetOrRegisterAura(core.Aura{
		Label:    "Sanctity of Battle",
		ActionID: core.ActionID{SpellID: 25956},

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			updateFloatValue(paladin.TotalRealHasteMultiplier())
			cooldownMod.Activate()
			gcdMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			gcdMod.Activate()
			cooldownMod.Deactivate()
		},
	}))
}
