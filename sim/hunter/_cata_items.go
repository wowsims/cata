package hunter

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

func (hunter *Hunter) newFlamingArrowSpell(spellID int32) core.SpellConfig {
	actionID := core.ActionID{SpellID: spellID} // actually 99058

	return core.SpellConfig{
		ActionID:    actionID.WithTag(3),
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

		DamageMultiplier: 0.8,
		CritMultiplier:   hunter.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.RangedWeaponDamage(sim, spell.RangedAttackPower())
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
				Kind:      core.SpellMod_PowerCost_Pct,
				ClassMask: HunterSpellsTierTwelve,
				IntValue:  -100,
			})
			var burningAdrenaline = hunter.RegisterAura(core.Aura{
				Label:    "Burning Adrenaline",
				Duration: time.Second * 15,
				ActionID: core.ActionID{SpellID: 99060},
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					baMod.Activate()
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if !spell.Matches(HunterSpellsTierTwelve) || spell.ActionID.SpellID == 0 {
						return
					}

					if hunter.HasActiveAura("Ready, Set, Aim...") && spell.Matches(HunterSpellAimedShot) {
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
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			hunter := agent.(HunterAgent).GetHunter()

			// Handled in Cobra and Steady code respectively
			hunter.T13_2pc = setBonusAura
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			hunter := agent.(HunterAgent).GetHunter()
			var chronoHunter = hunter.RegisterAura(core.Aura{ // 105919
				Label:    "Chronohunter",
				Duration: time.Second * 15,
				ActionID: core.ActionID{SpellID: 105919},
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.MultiplyRangedSpeed(sim, 1.3)
					if hunter.Pet != nil {
						hunter.Pet.Unit.MultiplyAttackSpeed(sim, 1.3)
					}
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.MultiplyRangedSpeed(sim, 1/1.3)
					if hunter.Pet != nil {
						hunter.Pet.Unit.MultiplyAttackSpeed(sim, 1/1.3)
					}
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
				ClassMask: HunterSpellSteadyShot,
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
	hunter.RegisterPvPGloveMod(
		[]int32{64991, 64709, 60424, 65544, 70534, 70260, 70441, 72369, 73717, 73583},
		core.SpellModConfig{
			ClassMask: HunterSpellExplosiveTrap | HunterSpellBlackArrow,
			Kind:      core.SpellMod_Cooldown_Flat,
			TimeValue: -time.Second * 2,
		})
}
