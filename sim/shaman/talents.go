package shaman

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (shaman *Shaman) ApplyTalents() {

	//"Hotfix (2013-09-23): Lightning Bolt's damage has been increased by 10%."
	shaman.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskLightningBolt | SpellMaskLightningBoltOverload,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.1,
	})
	//"Hotfix (2013-09-23): Flametongue Weapon's Flametongue Attack effect now deals 50% more damage."
	//"Hotfix (2013-09-23): Windfury Weapon's Windfury Attack effect now deals 50% more damage."
	shaman.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskFlametongueWeapon | SpellMaskWindfuryWeapon,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.5,
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
		Label:    "Elemental Mastery",
		ActionID: eleMasterActionID,
		Duration: time.Second * 20,
	}).AttachMultiplyCastSpeed(1.3)

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

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label:      "Ancestral Swiftness Passive",
		BuildPhase: core.CharacterBuildPhaseTalents,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			shaman.MultiplyMeleeSpeed(sim, 1.1)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			shaman.MultiplyMeleeSpeed(sim, 1/1.1)
		},
	}).AttachMultiplyCastSpeed(1.05))

	asCdTimer := shaman.NewTimer()
	asCd := time.Second * 90

	affectedSpells := SpellMaskLightningBolt | SpellMaskChainLightning | SpellMaskElementalBlast
	shaman.AncestralSwiftnessInstantAura = shaman.RegisterAura(core.Aura{
		Label:    "Ancestral swiftness",
		ActionID: core.ActionID{SpellID: 16188},
		Duration: core.NeverExpires,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !spell.Matches(affectedSpells) {
				return
			}
			//If both AS and MW 5 stacks buff are active, only MW gets consumed.
			//As i don't know which OnCastComplete is going to be executed first, check here if MW has not just been consumed/is active
			if shaman.Spec == proto.Spec_SpecEnhancementShaman && (shaman.MaelstromWeaponAura.TimeInactive(sim) == 0 && (!shaman.MaelstromWeaponAura.IsActive() || shaman.MaelstromWeaponAura.GetStacks() == 5)) {
				return
			}
			asCdTimer.Set(sim.CurrentTime + asCd)
			shaman.UpdateMajorCooldowns()
			aura.Deactivate(sim)
		},
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask:  affectedSpells,
		Kind:       core.SpellMod_CastTime_Pct,
		FloatValue: -100,
	})

	asSpell := shaman.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 16188},
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    asCdTimer,
				Duration: asCd,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			shaman.AncestralSwiftnessInstantAura.Activate(sim)
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

	var copySpells = map[*core.Spell]*core.Spell{}
	var alreadyProcced = map[*core.Spell]bool{}
	var lastTimestamp time.Duration

	core.MakePermanent(shaman.GetOrRegisterAura(core.Aura{
		Label: "Echo of The Elements Dummy",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || spell.Flags.Matches(SpellFlagIsEcho) || !spell.Flags.Matches(SpellFlagShamanSpell) {
				return
			}
			if sim.CurrentTime == lastTimestamp && alreadyProcced[spell] {
				return
			} else if sim.CurrentTime != lastTimestamp {
				lastTimestamp = sim.CurrentTime
				alreadyProcced = map[*core.Spell]bool{}
			}
			procChance := core.TernaryFloat64(shaman.Spec == proto.Spec_SpecElementalShaman, 0.06, 0.3)
			if spell.Matches(SpellMaskElementalBlast | SpellMaskElementalBlastOverload) {
				procChance = 0.06
			}
			if !sim.Proc(procChance, "Echo of The Elements") {
				return
			}
			alreadyProcced[spell] = true
			if copySpells[spell] == nil {
				copySpells[spell] = spell.Unit.RegisterSpell(core.SpellConfig{
					ActionID:                 core.ActionID{SpellID: spell.SpellID, Tag: core.TernaryInt32(spell.Tag == CastTagLightningOverload, 8, 7)},
					SpellSchool:              spell.SpellSchool,
					ProcMask:                 core.ProcMaskSpellProc,
					ApplyEffects:             spell.ApplyEffects,
					ManaCost:                 core.ManaCostOptions{},
					CritMultiplier:           spell.Unit.Env.Raid.GetPlayerFromUnit(spell.Unit).GetCharacter().DefaultCritMultiplier(),
					BonusCritPercent:         spell.BonusCritPercent,
					DamageMultiplier:         core.TernaryFloat64(spell.Tag == CastTagLightningOverload, 0.75, 1),
					DamageMultiplierAdditive: 1,
					MissileSpeed:             spell.MissileSpeed,
					ClassSpellMask:           spell.ClassSpellMask,
					BonusCoefficient:         spell.BonusCoefficient,
					Flags:                    spell.Flags & ^core.SpellFlagAPL | core.SpellFlagNoOnCastComplete | SpellFlagIsEcho,
					RelatedDotSpell:          spell.RelatedDotSpell,
				})
			}
			copySpell := copySpells[spell]
			copySpell.SpellMetrics[result.Target.UnitIndex].Casts--
			copySpell.Cast(sim, result.Target)
		},
	}))
}

func (shaman *Shaman) ApplyUnleashedFury() {
	if !shaman.Talents.UnleashedFury {
		return
	}

	unleashedFuryDDBCHandler := func(sim *core.Simulation, spell *core.Spell, attackTable *core.AttackTable) float64 {
		if spell.Matches(SpellMaskLightningBolt | SpellMaskLightningBoltOverload) {
			return 1.3
		}
		if spell.Matches(SpellMaskLavaBurst | SpellMaskLavaBurstOverload) {
			return 1.1
		}
		return 1.0
	}

	flametongueDebuffAura := shaman.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Unleashed Fury FT-" + shaman.Label,
			ActionID: core.ActionID{SpellID: 118470},
			Duration: time.Second * 10,
		}).AttachDDBC(DDBC_UnleashedFury, DDBC_Total, &shaman.AttackTables, unleashedFuryDDBCHandler)
	})

	windfuryProcAura := core.MakeProcTriggerAura(&shaman.Unit, core.ProcTrigger{
		Name:            "Unleashed Fury WF Proc Aura",
		MetricsActionID: core.ActionID{SpellID: 118472},
		Callback:        core.CallbackOnSpellHitDealt,
		ProcMask:        core.ProcMaskMeleeWhiteHit,
		Outcome:         core.OutcomeLanded,
		Duration:        time.Second * 8,
		ProcChance:      0.45,
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
		Callback:       core.CallbackOnApplyEffects,
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
	//In the corresponding pet files
}

func (shaman *Shaman) ApplyElementalBlast() {
	if !shaman.Talents.ElementalBlast {
		return
	}
	shaman.registerElementalBlastSpell()
}
