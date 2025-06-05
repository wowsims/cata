package retribution

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/paladin"
)

/*
Increases the damage you deal with two-handed melee weapons by 30%.

Your spell power is now equal to 50% of your attack power, and you no longer benefit from other sources of spell power.

Grants 6% of your maximum mana every 2 sec.

Increases the healing done by Word of Glory by 60% and Flash of Light by 100%.
*/
func (ret *RetributionPaladin) registerSwordOfLight() {
	actionID := core.ActionID{SpellID: 53503}
	swordOfLightHpActionID := core.ActionID{SpellID: 141459}
	ret.CanTriggerHolyAvengerHpGain(swordOfLightHpActionID)

	oldGetSpellPowerValue := ret.GetSpellPowerValue
	newGetSpellPowerValue := func(spell *core.Spell) float64 {
		return math.Floor(spell.MeleeAttackPower() * 0.5)
	}

	core.MakePermanent(ret.RegisterAura(core.Aura{
		Label:      "Sword of Light" + ret.Label,
		ActionID:   actionID,
		BuildPhase: core.CharacterBuildPhaseTalents,

		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			// Not in tooltip: Hammer of Wrath is usable during Avenging Wrath
			oldExtraCastCondition := ret.HammerOfWrath.ExtraCastCondition
			ret.HammerOfWrath.ExtraCastCondition = func(sim *core.Simulation, target *core.Unit) bool {
				return (oldExtraCastCondition != nil && oldExtraCastCondition(sim, target)) ||
					ret.AvengingWrathAura.IsActive()
			}
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			ret.GetSpellPowerValue = newGetSpellPowerValue
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			ret.GetSpellPowerValue = oldGetSpellPowerValue
		},
	})).AttachProcTrigger(core.ProcTrigger{
		// Not in tooltip: Hammer of Wrath generates a charge of Holy Power
		ClassSpellMask: paladin.SpellMaskHammerOfWrath,
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			ret.HolyPower.Gain(sim, 1, swordOfLightHpActionID)
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
		// Not in tooltip: Crusader Strike costs 40% less mana
		Kind:      core.SpellMod_PowerCost_Pct,
		ClassMask: paladin.SpellMaskCrusaderStrike,
		IntValue:  -80,
	}).AttachSpellMod(core.SpellModConfig{
		// Not in tooltip: Judgment costs 40% less mana
		Kind:      core.SpellMod_PowerCost_Pct,
		ClassMask: paladin.SpellMaskJudgment,
		IntValue:  -40,
	}).AttachSpellMod(core.SpellModConfig{
		// Not in tooltip: Cooldown of Avenging Wrath is reduced by 1 minute
		Kind:      core.SpellMod_Cooldown_Flat,
		ClassMask: paladin.SpellMaskAvengingWrath,
		TimeValue: time.Minute * -1,
	})

	manaMetrics := ret.NewManaMetrics(actionID)
	core.MakePermanent(ret.RegisterAura(core.Aura{
		Label: "Sword of Light Mana Regen" + ret.Label,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:   time.Second * 2,
				Priority: core.ActionPriorityRegen,
				OnAction: func(*core.Simulation) {
					ret.AddMana(sim, 0.06*ret.MaxMana(), manaMetrics)
				},
			})
		},
	}))

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
