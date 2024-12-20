package hunter

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

func (hunter *Hunter) newFlamingArrowSpell(spellID int32) core.SpellConfig {
	actionID := core.ActionID{SpellID: spellID} // actually 99058

	return core.SpellConfig{
		ActionID:    actionID.WithTag(3),
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

		DamageMultiplier: 0.8,
		CritMultiplier:   hunter.CritMultiplier(false, false, false),
		ThreatMultiplier: 1,

		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.RangedWeaponDamage(sim, spell.RangedAttackPower(target))
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
		},
	}
}

var ItemSetFlameWakersBattleGear = core.NewItemSet(core.ItemSet{
	Name: "Flamewaker's Battlegear",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			hunter := agent.(HunterAgent).GetHunter()

			var flamingArrowSpellForSteadyShot = hunter.RegisterSpell(hunter.newFlamingArrowSpell(56641))
			var flamingArrowSpellForCobraShot = hunter.RegisterSpell(hunter.newFlamingArrowSpell(77767))

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:           "T12 2-set",
				ClassSpellMask: HunterSpellSteadyShot | HunterSpellCobraShot,
				ProcChance:     0.1,
				Callback:       core.CallbackOnSpellHitDealt,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.Matches(HunterSpellSteadyShot) {
						flamingArrowSpellForSteadyShot.Cast(sim, result.Target)
					} else {
						flamingArrowSpellForCobraShot.Cast(sim, result.Target)
					}
				},
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			hunter := agent.(HunterAgent).GetHunter()
			var baMod = hunter.AddDynamicMod(core.SpellModConfig{
				Kind:       core.SpellMod_PowerCost_Pct,
				ClassMask:  HunterSpellsTierTwelve,
				FloatValue: -1,
			})
			var burningAdrenaline = hunter.RegisterAura(core.Aura{
				Label:    "Burning Adrenaline",
				Duration: time.Second * 15,
				ActionID: core.ActionID{SpellID: 99060},
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					baMod.Activate()
				},
				OnApplyEffects: func(aura *core.Aura, sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
					if spell.ClassSpellMask&HunterSpellsTierTwelve == 0 || spell.ActionID.SpellID == 0 {
						return
					}
					// https://www.bluetracker.gg/wow/topic/eu-en/510644-cataclysm-classic-hotfixes-11-december/
					// Breaks both Arcane Shot and Explosive Shot consuming of Burning Adrenaline
					// Arcane Shot is free if Lock and Load is up
					// if hunter.LockAndLoadAura.IsActive() && (spell.ClassSpellMask == HunterSpellExplosiveShot || spell.SpellID == 53301 || spell.SpellID == 1215485) {
					// 	return
					// }

					if hunter.HasActiveAura("Ready, Set, Aim...") && spell.ClassSpellMask == HunterSpellAimedShot {
						return
					}

					baMod.Deactivate()
					aura.Deactivate(sim)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					baMod.Deactivate()
				},
			})
			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:       "T12 4-set",
				ProcChance: 0.1,
				ProcMask:   core.ProcMaskRangedAuto,
				Callback:   core.CallbackOnSpellHitDealt,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					burningAdrenaline.Activate(sim)
				},
			})
		},
	},
})

var ItemSetWyrmstalkerBattleGear = core.NewItemSet(core.ItemSet{
	Name: "Wyrmstalker Battlegear",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(_ core.Agent, setBonusAura *core.Aura) {
			// Handled in Cobra Shot
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Flat,
				FloatValue: 0.1,
				ClassMask:  HunterSpellSteadyShot,
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			hunter := agent.(HunterAgent).GetHunter()
			var chronoHunter = hunter.RegisterAura(core.Aura{ // 105919
				Label:    "Chronohunter",
				Duration: time.Second * 15,
				ActionID: core.ActionID{SpellID: 105919},
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.MultiplyRangedSpeed(sim, 1.3)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.MultiplyRangedSpeed(sim, 1/1.3)
				},
			})

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:           "T13 4-set",
				Callback:       core.CallbackOnCastComplete,
				ClassSpellMask: HunterSpellArcaneShot,
				ProcChance:     0.4,
				ICD:            time.Second * 110,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					chronoHunter.Activate(sim)
				},
			})
		},
	},
})

func (hunter *Hunter) has2pcT13() bool {
	return hunter.HasActiveSetBonus(ItemSetWyrmstalkerBattleGear.Name, 2)
}

var ItemSetLightningChargedBattleGear = core.NewItemSet(core.ItemSet{
	Name: "Lightning-Charged Battlegear",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(_ core.Agent, setBonusAura *core.Aura) {
			// 5% Crit on SS
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_BonusCrit_Percent,
				ClassMask:  HunterSpellSerpentSting,
				FloatValue: 5,
			})
		},
		4: func(_ core.Agent, setBonusAura *core.Aura) {
			// Cobra & Steady Shot -0.2s cast time
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:      core.SpellMod_CastTime_Flat,
				ClassMask: HunterSpellCobraShot,
				TimeValue: -200 * time.Millisecond,
			})
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:      core.SpellMod_CastTime_Flat,
				ClassMask: HunterSpellCobraShot,
				TimeValue: -200 * time.Millisecond,
			})
		},
	},
})

var ItemSetGladiatorsPursuit = core.NewItemSet(core.ItemSet{
	ID:   920,
	Name: "Gladiator's Pursuit",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(_ core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachStatBuff(stats.Agility, 70)
		},
		4: func(_ core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachStatBuff(stats.Agility, 90)

			// Multiply focus regen 1.05
			focusRegenMultiplier := 1.05
			setBonusAura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.MultiplyFocusRegenSpeed(sim, focusRegenMultiplier)
			})
			setBonusAura.ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.MultiplyFocusRegenSpeed(sim, 1/focusRegenMultiplier)
			})

		},
	},
})

func (hunter *Hunter) addBloodthirstyGloves() {
	spellMod := hunter.AddDynamicMod(core.SpellModConfig{
		ClassMask: HunterSpellExplosiveTrap | HunterSpellBlackArrow,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: -time.Second * 2,
	})

	checkGloves := func() {
		switch hunter.Hands().ID {
		case 64991, 64709, 60424, 65544, 70534, 70260, 70441, 72369, 73717, 73583:
			spellMod.Activate()
			return
		default:
			spellMod.Deactivate()
		}
	}

	if hunter.ItemSwap.IsEnabled() {
		hunter.RegisterItemSwapCallback([]proto.ItemSlot{proto.ItemSlot_ItemSlotHands}, func(_ *core.Simulation, _ proto.ItemSlot) {
			checkGloves()
		})
	} else {
		checkGloves()
	}
}
