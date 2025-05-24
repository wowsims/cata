package protection

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/paladin"
)

func RegisterProtectionPaladin() {
	core.RegisterAgentFactory(
		proto.Player_ProtectionPaladin{},
		proto.Spec_SpecProtectionPaladin,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewProtectionPaladin(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_ProtectionPaladin)
			if !ok {
				panic("Invalid spec value for Protection Paladin!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewProtectionPaladin(character *core.Character, options *proto.Player) *ProtectionPaladin {
	protOptions := options.GetProtectionPaladin()

	prot := &ProtectionPaladin{
		Paladin: paladin.NewPaladin(character, options.TalentsString, protOptions.Options.ClassOptions),
		Options: protOptions.Options,
	}

	return prot
}

type ProtectionPaladin struct {
	*paladin.Paladin

	Options *proto.ProtectionPaladin_Options
}

func (prot *ProtectionPaladin) GetPaladin() *paladin.Paladin {
	return prot.Paladin
}

func (prot *ProtectionPaladin) Initialize() {
	prot.Paladin.Initialize()
	prot.ActivateRighteousFury()
	prot.registerAvengersShieldSpell()
	prot.registerHolyWrath()
	prot.registerConsecrationSpell()
	prot.registerSpecializationEffects()
}

func (prot *ProtectionPaladin) ApplyTalents() {
	prot.Paladin.ApplyTalents()
	prot.ApplyArmorSpecializationEffect(stats.Stamina, proto.ArmorType_ArmorTypePlate, 86525)
}

func (prot *ProtectionPaladin) Reset(sim *core.Simulation) {
	prot.Paladin.Reset(sim)
	prot.RighteousFuryAura.Activate(sim)
}

func (prot *ProtectionPaladin) registerSpecializationEffects() {
	prot.registerMastery()

	prot.applyGuardedByTheLight()
	prot.applySanctuary()
	prot.applyJudgmentsOfTheWise()
	prot.applyGrandCrusader()
	prot.applyArdentDefender()

	// Vengeance
	prot.RegisterVengeance(84839, nil)

	prot.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  paladin.SpellMaskSealOfTruth | paladin.SpellMaskCensure,
		FloatValue: 0.2,
	})
}

func (prot *ProtectionPaladin) registerMastery() {
	core.MakePermanent(prot.RegisterAura(core.Aura{
		Label:      "Mastery: Divine Bulwark" + prot.Label,
		ActionID:   core.ActionID{SpellID: 76671},
		BuildPhase: core.CharacterBuildPhaseBuffs,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			prot.ShieldOfTheRighteousAdditiveMultiplier = prot.getMasteryPercent()
		},
	})).AttachStatBuff(stats.BlockPercent, prot.getMasteryPercent())

	// Keep it updated when mastery changes
	prot.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMasteryRating float64, newMasteryRating float64) {
		prot.AddStatDynamic(sim, stats.BlockPercent, core.MasteryRatingToMasteryPoints(newMasteryRating-oldMasteryRating))
		prot.ShieldOfTheRighteousAdditiveMultiplier = prot.getMasteryPercent()
	})
}

func (prot *ProtectionPaladin) getMasteryPercent() float64 {
	return (8.0 + prot.GetMasteryPoints()) / 100.0
}

func (prot *ProtectionPaladin) applyGuardedByTheLight() {
	actionID := core.ActionID{SpellID: 53592}
	manaMetrics := prot.NewManaMetrics(actionID)

	oldGetSpellPowerValue := prot.GetSpellPowerValue
	newGetSpellPowerValue := func(spell *core.Spell) float64 {
		return spell.MeleeAttackPower() * 0.5
	}

	core.MakePermanent(prot.RegisterAura(core.Aura{
		Label:      "Guarded by the Light" + prot.Label,
		ActionID:   actionID,
		BuildPhase: core.CharacterBuildPhaseTalents,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:   time.Second * 2,
				Priority: core.ActionPriorityRegen,
				OnAction: func(*core.Simulation) {
					prot.AddMana(sim, 0.15*prot.MaxMana(), manaMetrics)
				},
			})

			prot.GetSpellPowerValue = newGetSpellPowerValue
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			prot.GetSpellPowerValue = oldGetSpellPowerValue
		},
	})).AttachStatDependency(
		prot.NewDynamicMultiplyStat(stats.Stamina, 1.25),
	).AttachMultiplicativePseudoStatBuff(
		&prot.PseudoStats.BaseBlockChance,
		0.1,
	).AttachAdditivePseudoStatBuff(
		&prot.PseudoStats.ReducedCritTakenChance,
		0.06,
	).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Pct,
		ClassMask: paladin.SpellMaskCrusaderStrike,
		IntValue:  -80,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Pct,
		ClassMask: paladin.SpellMaskJudgment,
		IntValue:  -40,
	})
}

func (prot *ProtectionPaladin) applySanctuary() {
	core.MakePermanent(prot.RegisterAura(core.Aura{
		Label:      "Sanctuary" + prot.Label,
		ActionID:   core.ActionID{SpellID: 105805},
		BuildPhase: core.CharacterBuildPhaseBuffs,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			prot.ApplyDynamicEquipScaling(sim, stats.Armor, 1.1)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			prot.RemoveDynamicEquipScaling(sim, stats.Armor, 1.1)
		},
	})).AttachAdditivePseudoStatBuff(
		&prot.PseudoStats.BaseDodgeChance,
		0.02,
	).AttachMultiplicativePseudoStatBuff(
		&prot.PseudoStats.DamageTakenMultiplier,
		0.85,
	)
}

func (prot *ProtectionPaladin) applyJudgmentsOfTheWise() {
	jotwHpActionID := core.ActionID{SpellID: 105427}
	prot.CanTriggerHolyAvengerHpGain(jotwHpActionID)
	core.MakeProcTriggerAura(&prot.Unit, core.ProcTrigger{
		Name:           "Judgments of the Wise" + prot.Label,
		ActionID:       core.ActionID{SpellID: 105424},
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: paladin.SpellMaskJudgment,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			prot.HolyPower.Gain(1, jotwHpActionID, sim)
		},
	})
}

func (prot *ProtectionPaladin) applyGrandCrusader() {
	hpActionID := core.ActionID{SpellID: 98057}
	prot.CanTriggerHolyAvengerHpGain(hpActionID)

	var grandCrusaderAura *core.Aura
	grandCrusaderAura = prot.RegisterAura(core.Aura{
		Label:    "Grand Crusader" + prot.Label,
		ActionID: core.ActionID{SpellID: 85416},
		Duration: time.Second * 6,
	}).AttachProcTrigger(core.ProcTrigger{
		Name:           "Grand Crusader Consume Trigger" + prot.Label,
		Callback:       core.CallbackOnCastComplete,
		ClassSpellMask: paladin.SpellMaskAvengersShield,
		Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
			prot.HolyPower.Gain(1, hpActionID, sim)
			grandCrusaderAura.Deactivate(sim)
		},
	})

	core.MakeProcTriggerAura(&prot.Unit, core.ProcTrigger{
		Name:       "Grand Crusader Trigger" + prot.Label,
		ActionID:   core.ActionID{SpellID: 85043},
		Callback:   core.CallbackOnSpellHitTaken,
		Outcome:    core.OutcomeDodge | core.OutcomeParry,
		ProcChance: 0.3,
		ICD:        time.Second,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			prot.AvengersShield.CD.Reset()
			grandCrusaderAura.Activate(sim)
		},
	})
}

func (prot *ProtectionPaladin) applyArdentDefender() {
	actionID := core.ActionID{SpellID: 31850}

	adAura := prot.RegisterAura(core.Aura{
		Label:    "Ardent Defender" + prot.Label,
		ActionID: actionID,
		Duration: time.Second * 10,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			prot.PseudoStats.DamageTakenMultiplier *= 0.8
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			prot.PseudoStats.DamageTakenMultiplier /= 0.8
		},
	})

	ardentDefender := prot.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		Flags:       core.SpellFlagAPL,
		SpellSchool: core.SpellSchoolHoly,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    prot.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			adAura.Activate(sim)
		},
	})

	adHealAmount := 0.0

	// Spell to heal you when AD has procced; fire this before fatal damage so that a Death is not detected
	adHeal := prot.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 66235},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagHelpful,

		CritMultiplier:   1,
		ThreatMultiplier: 0,
		DamageMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealHealing(sim, &prot.Unit, adHealAmount, spell.OutcomeHealing)
		},
	})

	// >= 15% hp, hit gets reduced so we end up at 15% without heal
	// < 15% hp, hit gets reduced to 0 and we heal the remaining health up to 15%
	prot.AddDynamicDamageTakenModifier(func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult, isPeriodic bool) {
		if adAura.IsActive() && result.Damage >= prot.CurrentHealth() {
			maxHealth := prot.MaxHealth()
			currentHealth := prot.CurrentHealth()
			incomingDamage := result.Damage

			if currentHealth/maxHealth >= 0.15 {
				// Incoming attack gets reduced so we end up at 15% hp
				// TODO: Overkill counted as absorb but not as healing in logs
				result.Damage = currentHealth - maxHealth*0.15
				if sim.Log != nil {
					prot.Log(sim, "Ardent Defender absorbed %.1f damage", incomingDamage-result.Damage)
				}
			} else {
				// Incoming attack gets reduced to 0
				// Heal up to 15% hp
				// TODO: Overkill counted as absorb but not as healing in logs
				result.Damage = 0
				adHealAmount = maxHealth*0.15 - currentHealth
				adHeal.Cast(sim, &prot.Unit)
				if sim.Log != nil {
					prot.Log(sim, "Ardent Defender absorbed %.1f damage and healed for %.1f", incomingDamage, adHealAmount)
				}
			}

			adAura.Deactivate(sim)
		}
	})

	prot.AddMajorCooldown(core.MajorCooldown{
		Spell: ardentDefender,
		Type:  core.CooldownTypeSurvival,
	})
}
