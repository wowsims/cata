package paladin

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

func (paladin *Paladin) applyProtectionTalents() {
	paladin.applySealsOfThePure()
	paladin.applyToughness()
	paladin.applyHallowedGround()
	paladin.applySanctuary()
	paladin.applyWrathOfTheLightbringer()
	paladin.applyHammerOfTheRighteous()
	paladin.applyReckoning()
	paladin.applyShieldOfTheRighteous()
	paladin.applyShieldOfTheTemplar()
}

func (paladin *Paladin) applySealsOfThePure() {
	if paladin.Talents.SealsOfThePure == 0 {
		return
	}

	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskSealOfRighteousness | SpellMaskSealOfTruth | SpellMaskSealOfJustice,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.06 * float64(paladin.Talents.SealsOfThePure),
	})
}

func (paladin *Paladin) applyToughness() {
	if paladin.Talents.Toughness == 0 {
		return
	}

	paladin.ApplyEquipScaling(stats.Armor, []float64{1.0, 1.03, 1.06, 1.1}[paladin.Talents.Toughness])
}

func (paladin *Paladin) applyHallowedGround() {
	if paladin.Talents.HallowedGround == 0 {
		return
	}

	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskConsecration,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.2 * float64(paladin.Talents.HallowedGround),
	})

	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskConsecration,
		Kind:       core.SpellMod_PowerCost_Pct,
		FloatValue: -(0.4 * float64(paladin.Talents.HallowedGround)),
	})
}

func (paladin *Paladin) applySanctuary() {
	if paladin.Talents.Sanctuary == 0 {
		return
	}

	paladin.PseudoStats.ReducedCritTakenChance += 0.02 * float64(paladin.Talents.Sanctuary)
	paladin.PseudoStats.DamageTakenMultiplier *= 1.0 - []float64{0, 0.03, 0.07, 0.1}[paladin.Talents.Sanctuary]

	manaReturnActionID := core.ActionID{SpellID: []int32{0, 57319, 84626, 84627}[paladin.Talents.Sanctuary]}
	manaMetrics := paladin.NewManaMetrics(manaReturnActionID)
	manaReturnPct := 0.01 * float64(paladin.Talents.Sanctuary)

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:     "Sanctuary",
		Callback: core.CallbackOnSpellHitTaken,
		ProcMask: core.ProcMaskMelee,
		Outcome:  core.OutcomeBlock | core.OutcomeDodge,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			paladin.AddMana(sim, manaReturnPct*paladin.MaxMana(), manaMetrics)
		},
	})
}

func (paladin *Paladin) applyHammerOfTheRighteous() {
	if !paladin.Talents.HammerOfTheRighteous {
		return
	}

	aoeMinDamage, aoeMaxDamage :=
		core.CalcScalingSpellEffectVarianceMinMax(proto.Class_ClassPaladin, 0.708, 0.4)

	numTargets := paladin.Env.GetNumTargets()
	actionId := core.ActionID{SpellID: 53595}
	hpMetrics := paladin.NewHolyPowerMetrics(actionId)

	hammerOfTheRighteousAoe := paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 88263},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagMeleeMetrics,
		ClassSpellMask: SpellMaskHammerOfTheRighteousAoe,

		MaxRange: 8,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(aoeMinDamage, aoeMaxDamage) +
				0.18*spell.MeleeAttackPower()
			results := make([]*core.SpellResult, numTargets)

			for idx := int32(0); idx < numTargets; idx++ {
				currentTarget := sim.Environment.GetTargetUnit(idx)
				results[idx] = spell.CalcDamage(sim, currentTarget, baseDamage, spell.OutcomeMagicCrit)
			}

			for idx := int32(0); idx < numTargets; idx++ {
				spell.DealDamage(sim, results[idx])
			}
		},
	})

	paladin.HammerOfTheRighteous = paladin.RegisterSpell(core.SpellConfig{
		ActionID:       actionId,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskHammerOfTheRighteousMelee,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.10,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    paladin.sharedBuilderTimer,
				Duration: paladin.sharedBuilderBaseCD,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return paladin.MainHand().HandType != proto.HandType_HandTypeTwoHand
		},

		DamageMultiplier: 0.3,
		CritMultiplier:   paladin.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				hammerOfTheRighteousAoe.Cast(sim, target)
				paladin.GainHolyPower(sim, 1, hpMetrics)
			}
		},
	})
}

func (paladin *Paladin) applyWrathOfTheLightbringer() {
	if paladin.Talents.WrathOfTheLightbringer == 0 {
		return
	}

	dmgIncrease := 0.5 * float64(paladin.Talents.WrathOfTheLightbringer)

	// For some reason, only Crusader Strike and JoT are additive, while the rest are multiplicative.
	// Dunno if this is actually correct but that's how simc does it.
	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskCrusaderStrike | SpellMaskJudgementOfTruth,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: dmgIncrease,
	})

	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskJudgementOfJustice | SpellMaskJudgementOfInsight | SpellMaskJudgementOfRighteousness,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: dmgIncrease,
	})

	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskHammerOfWrath | SpellMaskHolyWrath,
		Kind:       core.SpellMod_BonusCrit_Rating,
		FloatValue: 15 * float64(paladin.Talents.WrathOfTheLightbringer) * core.CritRatingPerCritChance,
	})
}

func (paladin *Paladin) applyReckoning() {
	if paladin.Talents.Reckoning == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 20178}
	procChance := 0.1 * float64(paladin.Talents.Reckoning)

	var reckoningSpell *core.Spell

	procAura := paladin.RegisterAura(core.Aura{
		Label:     "Reckoning Proc",
		ActionID:  actionID,
		Duration:  time.Second * 8,
		MaxStacks: 4,

		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			config := *paladin.AutoAttacks.MHConfig()
			config.ActionID = actionID
			reckoningSpell = paladin.GetOrRegisterSpell(config)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell == paladin.AutoAttacks.MHAuto() {
				reckoningSpell.Cast(sim, result.Target)
			}
		},
	})

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:       "Reckoning",
		ProcMask:   core.ProcMaskMelee,
		ProcChance: procChance,
		Outcome:    core.OutcomeBlock,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			procAura.Activate(sim)
			procAura.SetStacks(sim, 4)
		},
	})
}

func (paladin *Paladin) applyShieldOfTheRighteous() {
	if !paladin.Talents.ShieldOfTheRighteous {
		return
	}

	actionId := core.ActionID{SpellID: 53600}
	hpMetrics := paladin.NewHolyPowerMetrics(actionId)

	shieldDmg := core.CalcScalingSpellAverageEffect(proto.Class_ClassPaladin, 0.593)

	paladin.ShieldOfTheRighteous = paladin.RegisterSpell(core.SpellConfig{
		ActionID:       actionId,
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskShieldOfTheRighteous,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return paladin.GetHolyPowerValue() > 0
		},

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {

			baseDamage := []float64{0, 1, 3, 6}[paladin.GetHolyPowerValue()] *
				(shieldDmg + 0.1*spell.MeleeAttackPower())

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				paladin.SpendHolyPower(sim, hpMetrics)
			}
		},
	})
}

func (paladin *Paladin) applyShieldOfTheTemplar() {
	if paladin.Talents.ShieldOfTheTemplar == 0 {
		return
	}

	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskGuardianOfAncientKings,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: -(time.Second * time.Duration(40*paladin.Talents.ShieldOfTheTemplar)),
	})

	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskAvengingWrath,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: -(time.Second * time.Duration(20*paladin.Talents.ShieldOfTheTemplar)),
	})
}
