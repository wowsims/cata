package shaman

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (shaman *Shaman) ApplyTalents() {

	//"Hotfix (2013-09-23): Lightning Bolt's damage has been increased by 10%."
	//Additive with shamanism
	shaman.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskLightningBolt | SpellMaskLightningBoltOverload,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: 0.1,
	})

	shaman.ApplyElementalMastery()
	shaman.ApplyAncestralSwiftness()
	shaman.ApplyEchoOfTheElements()
	shaman.ApplyUnleashedFury()
	shaman.ApplyPrimalElementalist()
	shaman.ApplyElementalBlast()

	shaman.ApplyGlyphs()
}

func (shaman *Shaman) ApplyElementalMastery() {
	if !shaman.Talents.ElementalMastery {
		return
	}

	eleMasterActionID := core.ActionID{SpellID: 16166}

	buffAura := shaman.RegisterAura(core.Aura{
		Label:    "Elemental Mastery Buff",
		ActionID: core.ActionID{SpellID: 64701},
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			shaman.MultiplyCastSpeed(1.30)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			shaman.MultiplyCastSpeed(1 / 1.30)
		},
	})

	eleMastSpell := shaman.RegisterSpell(core.SpellConfig{
		ActionID:       eleMasterActionID,
		ClassSpellMask: SpellMaskElementalMastery,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Second * 90,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			buffAura.Activate(sim)
		},
	})

	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell: eleMastSpell,
		Type:  core.CooldownTypeDPS,
	})
}

func (shaman *Shaman) ApplyAncestralSwiftness() {
	if !shaman.Talents.AncestralSwiftness {
		return
	}

	asCdTimer := shaman.NewTimer()
	asCd := time.Second * 90

	affectedSpells := SpellMaskLightningBolt | SpellMaskChainLightning | SpellMaskEarthShock | SpellMaskElementalBlast
	AncestralSwiftnessBuffAura := core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label:    "Ancestral swiftness Buff",
		ActionID: core.ActionID{SpellID: 64701},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			shaman.MultiplyCastSpeed(1.05)
			shaman.MultiplyMeleeSpeed(sim, 1.1)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			shaman.MultiplyCastSpeed(1 / 1.30)
			shaman.MultiplyMeleeSpeed(sim, 1/1.05)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.ClassSpellMask&affectedSpells == 0 {
				return
			}
			asCdTimer.Set(sim.CurrentTime + asCd)
			shaman.UpdateMajorCooldowns()
			aura.Deactivate(sim)
		},
	}))

	asSpell := shaman.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 16188},
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			//TODO goes on cd only when buff is consumed
			CD: core.Cooldown{
				Timer:    asCdTimer,
				Duration: asCd,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			AncestralSwiftnessBuffAura.Activate(sim)
		},
	})

	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell: asSpell,
		Type:  core.CooldownTypeDPS,
	})
}

func (shaman *Shaman) ApplyEchoOfTheElements() {
	if !shaman.Talents.EchoOfTheElements {
		return
	}

	core.MakeProcTriggerAura(&shaman.Unit, core.ProcTrigger{
		Name:       "Echo of The Elements Dummy",
		Callback:   core.CallbackOnSpellHitDealt,
		SpellFlags: SpellFlagShamanSpell,
		Outcome:    core.OutcomeHit,
		ProcChance: 0.05,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			tempFlags := spell.Flags
			spell.Flags &= ^core.SpellFlagNoOnCastComplete
			spell.SkipCastAndApplyEffects(sim, result.Target)
			spell.Flags = tempFlags
		},
	})
}

func (shaman *Shaman) ApplyUnleashedFury() {
	if !shaman.Talents.UnleashedFury {
		return
	}

	flametongueDebuffAura := shaman.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Unleashed Fury FT-" + shaman.Label,
			ActionID: core.ActionID{SpellID: 118470},
			Duration: time.Second * 10,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				core.EnableDamageDoneByCaster(DDBC_UnleashedFury, DDBC_Total, shaman.AttackTables[aura.Unit.Index], func(sim *core.Simulation, spell *core.Spell, attackTable *core.AttackTable) float64 {
					if spell.ClassSpellMask&(SpellMaskLightningBolt|SpellMaskLightningBoltOverload) > 0 {
						return 1.3
					}
					if spell.ClassSpellMask&(SpellMaskLavaBurst|SpellMaskLavaBurstOverload) > 0 {
						return 1.1
					}
					return 1.0
				})
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				core.DisableDamageDoneByCaster(DDBC_UnleashedFury, shaman.AttackTables[aura.Unit.Index])
			},
		})
	})

	windfuryProcAura := core.MakeProcTriggerAura(&shaman.Unit, core.ProcTrigger{
		Name:     "Unleashed Fury WF Proc Aura",
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeWhiteHit,
		Outcome:  core.OutcomeHit,
		Duration: time.Second * 8,
		ExtraCondition: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
			return shaman.SelfBuffs.Shield == proto.ShamanShield_LightningShield
		},
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			shaman.LightningShieldDamage.Cast(sim, result.Target)
		},
	})

	core.MakeProcTriggerAura(&shaman.Unit, core.ProcTrigger{
		Name:           "Unleashed Fury",
		ActionID:       core.ActionID{SpellID: 117012},
		Callback:       core.CallbackOnCastComplete,
		ClassSpellMask: SpellMaskUnleashElements,
		ProcChance:     1.0,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			switch shaman.SelfBuffs.ImbueMH {
			case proto.ShamanImbue_FlametongueWeapon:
				flametongueDebuffAura.Get(result.Target).Activate(sim)
			case proto.ShamanImbue_WindfuryWeapon:
				windfuryProcAura.Activate(sim)
			case proto.ShamanImbue_EarthlivingWeapon:
			case proto.ShamanImbue_FrostbrandWeapon:
			case proto.ShamanImbue_RockbiterWeapon:
			}
			if shaman.SelfBuffs.ImbueOH != proto.ShamanImbue_NoImbue && shaman.SelfBuffs.ImbueOH != shaman.SelfBuffs.ImbueMH {
				switch shaman.SelfBuffs.ImbueOH {
				case proto.ShamanImbue_FlametongueWeapon:
					flametongueDebuffAura.Get(result.Target).Activate(sim)
				case proto.ShamanImbue_WindfuryWeapon:
					windfuryProcAura.Activate(sim)
				case proto.ShamanImbue_EarthlivingWeapon:
				case proto.ShamanImbue_FrostbrandWeapon:
				case proto.ShamanImbue_RockbiterWeapon:
				}
			}
		},
	})
}

func (shaman *Shaman) ApplyPrimalElementalist() {
	if !shaman.Talents.PrimalElementalist {
		return
	}

	//TODO
}

func (shaman *Shaman) ApplyElementalBlast() {
	if !shaman.Talents.ElementalBlast {
		return
	}
	shaman.RegisterElementalBlastSpell()
}
