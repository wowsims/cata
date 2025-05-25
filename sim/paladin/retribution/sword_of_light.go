package retribution

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/paladin"
)

func (ret *RetributionPaladin) registerSwordOfLight() {
	actionID := core.ActionID{SpellID: 53503}
	manaMetrics := ret.NewManaMetrics(actionID)
	swordOfLightHpActionID := core.ActionID{SpellID: 141459}
	ret.CanTriggerHolyAvengerHpGain(swordOfLightHpActionID)

	oldGetSpellPowerValue := ret.GetSpellPowerValue
	newGetSpellPowerValue := func(spell *core.Spell) float64 {
		return spell.MeleeAttackPower() * 0.5
	}

	core.MakePermanent(ret.RegisterAura(core.Aura{
		Label:      "Sword of Light" + ret.Label,
		ActionID:   actionID,
		BuildPhase: core.CharacterBuildPhaseTalents,

		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			oldExtraCastCondition := ret.HammerOfWrath.ExtraCastCondition
			ret.HammerOfWrath.ExtraCastCondition = func(sim *core.Simulation, target *core.Unit) bool {
				return (oldExtraCastCondition != nil && oldExtraCastCondition(sim, target)) ||
					ret.AvengingWrathAura.IsActive()
			}
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:   time.Second * 2,
				Priority: core.ActionPriorityRegen,
				OnAction: func(*core.Simulation) {
					ret.AddMana(sim, 0.06*ret.MaxMana(), manaMetrics)
				},
			})

			ret.GetSpellPowerValue = newGetSpellPowerValue
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			ret.GetSpellPowerValue = oldGetSpellPowerValue
		},
	})).AttachProcTrigger(core.ProcTrigger{
		ClassSpellMask: paladin.SpellMaskHammerOfWrath,
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			ret.HolyPower.Gain(1, swordOfLightHpActionID, sim)
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  paladin.SpellMaskWordOfGlory,
		FloatValue: 0.6,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  paladin.SpellMaskFlashOfLight,
		FloatValue: 1.0,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Pct,
		ClassMask: paladin.SpellMaskCrusaderStrike,
		IntValue:  -80,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Pct,
		ClassMask: paladin.SpellMaskJudgment,
		IntValue:  -40,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		ClassMask: paladin.SpellMaskAvengingWrath,
		TimeValue: time.Minute * -1,
	})

	holyTwoHandDamageMod := ret.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  paladin.SpellMaskDamageModifiedBySwordOfLight,
		FloatValue: 0.3,
	})

	checkWeaponType := func() {
		mhWeapon := ret.GetMHWeapon()
		if mhWeapon != nil && mhWeapon.HandType == proto.HandType_HandTypeTwoHand {
			if !holyTwoHandDamageMod.IsActive {
				ret.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1.3
			}
			holyTwoHandDamageMod.Activate()
		} else if holyTwoHandDamageMod.IsActive {
			ret.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= 1.3
			holyTwoHandDamageMod.Deactivate()
		}
	}

	checkWeaponType()

	ret.RegisterItemSwapCallback([]proto.ItemSlot{proto.ItemSlot_ItemSlotHands}, func(_ *core.Simulation, _ proto.ItemSlot) {
		checkWeaponType()
	})
}
