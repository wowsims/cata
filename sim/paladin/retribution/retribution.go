package retribution

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/paladin"
)

func RegisterRetributionPaladin() {
	core.RegisterAgentFactory(
		proto.Player_RetributionPaladin{},
		proto.Spec_SpecRetributionPaladin,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewRetributionPaladin(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_RetributionPaladin)
			if !ok {
				panic("Invalid spec value for Retribution Paladin!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewRetributionPaladin(character *core.Character, options *proto.Player) *RetributionPaladin {
	retOptions := options.GetRetributionPaladin()

	ret := &RetributionPaladin{
		Paladin: paladin.NewPaladin(character, options.TalentsString, retOptions.Options.ClassOptions),
	}
	ret.StartingHolyPower = retOptions.Options.StartingHolyPower

	return ret
}

type RetributionPaladin struct {
	*paladin.Paladin

	HoLDamage float64
}

func (ret *RetributionPaladin) GetPaladin() *paladin.Paladin {
	return ret.Paladin
}

func (ret *RetributionPaladin) Initialize() {
	ret.Paladin.Initialize()
	ret.registerSpecializationEffects()
	ret.registerSealOfJustice()
	ret.registerInquisition()
	ret.registerExorcism()
	ret.registerDivineStorm()
}

func (ret *RetributionPaladin) ApplyTalents() {
	ret.Paladin.ApplyTalents()
	ret.ApplyArmorSpecializationEffect(stats.Strength, proto.ArmorType_ArmorTypePlate, 86525)
}

func (ret *RetributionPaladin) Reset(sim *core.Simulation) {
	ret.Paladin.Reset(sim)
}

func (ret *RetributionPaladin) registerSpecializationEffects() {
	ret.registerMastery()

	ret.applyJudgmentsOfTheBold()
	ret.applyArtOfWar()
	ret.applySwordOfLight()
}

func (ret *RetributionPaladin) registerMastery() {
	handOfLight := ret.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 96172},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIgnoreModifiers | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

		DamageMultiplier: 1.0,
		ThreatMultiplier: 1.0,
		CritMultiplier:   0.0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := ret.HoLDamage
			if target.HasActiveAuraWithTag(core.SpellDamageEffectAuraTag) {
				baseDamage *= 1.05
			}
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeAlwaysHit)
		},
	})

	core.MakeProcTriggerAura(&ret.Unit, core.ProcTrigger{
		Name:           "Mastery: Hand of Light" + ret.Label,
		ActionID:       core.ActionID{SpellID: 76672},
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: paladin.SpellMaskCanTriggerHandOfLight,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			ret.HoLDamage = ret.getMasteryPercent() * result.Damage
			handOfLight.Cast(sim, result.Target)
		},
	})
}

func (ret *RetributionPaladin) getMasteryPercent() float64 {
	return ((8.0 + ret.GetMasteryPoints()) * 1.85000002384) / 100.0
}

func (ret *RetributionPaladin) applyJudgmentsOfTheBold() {
	actionID := core.ActionID{SpellID: 111528}
	ret.CanTriggerHolyAvengerHpGain(actionID)
	auraArray := ret.NewEnemyAuraArray(core.PhysVulnerabilityAura)
	core.MakeProcTriggerAura(&ret.Unit, core.ProcTrigger{
		Name:           "Judgments of the Bold" + ret.Label,
		ActionID:       core.ActionID{SpellID: 111529},
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: paladin.SpellMaskJudgment,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			ret.HolyPower.Gain(1, actionID, sim)

			auraArray.Get(result.Target).Activate(sim)
		},
	})
}

func (ret *RetributionPaladin) applyArtOfWar() {
	artOfWarAura := ret.RegisterAura(core.Aura{
		Label:    "The Art Of War" + ret.Label,
		ActionID: core.ActionID{SpellID: 59578},
		Duration: time.Second * 6,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			ret.Exorcism.CD.Reset()
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(paladin.SpellMaskExorcism) {
				aura.Deactivate(sim)
			}
		},
	})

	core.MakeProcTriggerAura(&ret.Unit, core.ProcTrigger{
		Name:       "Art of War" + ret.Label,
		ActionID:   core.ActionID{SpellID: 87138},
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeWhiteHit,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.20,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			artOfWarAura.Activate(sim)
		},
	})
}

func (ret *RetributionPaladin) applySwordOfLight() {
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
		Name:           "Hammer of Wrath Holy Power Gain" + ret.Label,
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
