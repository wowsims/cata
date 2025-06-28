package warrior

import (
	"time"

	"github.com/wowsims/mop/sim/common/shared"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (warrior *Warrior) ApplyTalents() {
	// Level 15
	warrior.registerJuggernaut()

	// Level 30
	warrior.registerImpendingVictory()

	// Level 45

	// Level 60
	warrior.registerBladestorm()
	warrior.registerShockwave()
	warrior.registerDragonRoar()

	// Level 75

	// Level 90
	warrior.registerAvatar()
	warrior.registerBloodbath()
	warrior.registerStormBolt()
}

func (war *Warrior) registerJuggernaut() {
	if !war.Talents.Juggernaut {
		return
	}

	war.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskCharge,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: -8 * time.Second,
	})
}

func (war *Warrior) registerImpendingVictory() {
	if !war.Talents.ImpendingVictory {
		return
	}

	actionID := core.ActionID{SpellID: 103840}
	healthMetrics := war.NewHealthMetrics(actionID)

	war.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagAPL | core.SpellFlagMeleeMetrics,
		ClassSpellMask: SpellMaskImpendingVictory,

		RageCost: core.RageCostOptions{
			Cost:   10,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: time.Second * 30,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   war.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			war.VictoryRushAura.Deactivate(sim)

			baseDamage := 56 + spell.MeleeAttackPower()*0.56
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			healthMultiplier := core.TernaryFloat64(war.T15Tank2P != nil && war.T15Tank2P.IsActive(), 0.4, 0.2)

			if result.Landed() {
				war.GainHealth(sim, war.MaxHealth()*healthMultiplier, healthMetrics)
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}

func (war *Warrior) registerDragonRoar() {
	if !war.Talents.DragonRoar {
		return
	}

	actionID := core.ActionID{SpellID: 118000}

	damageMultipliers := []float64{1, 0.75, 0.65, 0.55, 0.50}

	spell := war.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ClassSpellMask: SpellMaskDragonRoar,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagAPL | core.SpellFlagMeleeMetrics | core.SpellFlagIgnoreArmor | core.SpellFlagReadinessTrinket,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: time.Minute * 1,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		CritMultiplier:   war.DefaultCritMultiplier(),
		BonusCritPercent: 100,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damageMultiplier := damageMultipliers[min(war.Env.GetNumTargets()-1, 4)]
			baseDamage := 126 + spell.MeleeAttackPower()*1.39999997616
			spell.DamageMultiplier *= damageMultiplier
			for _, enemyTarget := range sim.Encounter.ActiveTargets {
				spell.CalcAndDealDamage(sim, &enemyTarget.Unit, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
			}
			spell.DamageMultiplier /= damageMultiplier
		},
	})

	war.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}

func (war *Warrior) registerBladestorm() {
	if !war.Talents.Bladestorm {
		return
	}

	actionID := core.ActionID{SpellID: 46924}

	damageMultiplier := 1.2
	if war.Spec == proto.Spec_SpecArmsWarrior {
		damageMultiplier += 0.6
	} else if war.Spec == proto.Spec_SpecProtectionWarrior {
		damageMultiplier *= 1.33
	}

	results := make([]*core.SpellResult, war.Env.GetNumTargets())

	mhSpell := war.RegisterSpell(core.SpellConfig{
		ActionID:       actionID.WithTag(1), // Real Spell ID: 50622
		SpellSchool:    core.SpellSchoolPhysical,
		ClassSpellMask: SpellMaskBladestormMH,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagPassiveSpell,

		DamageMultiplier: damageMultiplier,
		CritMultiplier:   war.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for i, enemyTarget := range sim.Encounter.ActiveTargets {
				baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
				results[i] = spell.CalcDamage(sim, &enemyTarget.Unit, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			}

			war.CastNormalizedSweepingStrikesAttack(results, sim)

			for _, result := range results {
				spell.DealDamage(sim, result)
			}
		},
	})

	ohSpell := war.RegisterSpell(core.SpellConfig{
		ActionID:       actionID.WithTag(2), // Real Spell ID: 95738,
		SpellSchool:    core.SpellSchoolPhysical,
		ClassSpellMask: SpellMaskBladestormOH,
		ProcMask:       core.ProcMaskMeleeOHSpecial,
		Flags:          core.SpellFlagPassiveSpell,

		DamageMultiplier: damageMultiplier,
		CritMultiplier:   war.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			for _, enemyTarget := range sim.Encounter.ActiveTargets {
				baseDamage := spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
				spell.CalcAndDealDamage(sim, &enemyTarget.Unit, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			}
		},
	})

	spell := war.RegisterSpell(core.SpellConfig{
		ActionID:       actionID.WithTag(0),
		SpellSchool:    core.SpellSchoolPhysical,
		ClassSpellMask: SpellMaskBladestorm,
		Flags:          core.SpellFlagChanneled | core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ProcMask:       core.ProcMaskEmpty,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: time.Minute * 1,
			},
		},

		DamageMultiplier: 1.0,
		CritMultiplier:   war.DefaultCritMultiplier(),

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "Bladestorm",
			},
			NumberOfTicks: 6,
			TickLength:    time.Second * 1,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				mhSpell.Cast(sim, target)

				if war.OffHand() != nil && war.OffHand().WeaponType != proto.WeaponType_WeaponTypeUnknown {
					ohSpell.Cast(sim, target)
				}
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dot := spell.AOEDot()
			dot.Apply(sim)
			dot.TickOnce(sim)
		},
	})

	war.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}

func (war *Warrior) registerShockwave() {
	if !war.Talents.Shockwave {
		return
	}

	war.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 46968},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		ClassSpellMask: SpellMaskShockwave,
		Flags:          core.SpellFlagAoE | core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagReadinessTrinket,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: 40 * time.Second,
			},
		},

		DamageMultiplier: 0.75,
		CritMultiplier:   war.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			numLandedHits := 0
			baseDamage := spell.MeleeAttackPower()
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				result := spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

				if result.Landed() {
					numLandedHits++
				}
			}
			if numLandedHits >= 3 {
				spell.CD.Reduce(time.Second * 20)
			}
		},
	})
}

func (war *Warrior) registerAvatar() {
	if !war.Talents.Avatar {
		return
	}

	actionId := core.ActionID{SpellID: 107574}
	avatarAura := war.RegisterAura(core.Aura{
		Label:    "Avatar",
		ActionID: actionId,
		Duration: 24 * time.Second,
	}).AttachMultiplicativePseudoStatBuff(&war.Unit.PseudoStats.DamageDealtMultiplier, 1.2)

	avatar := war.RegisterSpell(core.SpellConfig{
		ActionID:       actionId,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL | core.SpellFlagReadinessTrinket,
		ClassSpellMask: SpellMaskAvatar,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: 3 * time.Minute,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			avatarAura.Activate(sim)
		},
	})

	war.AddMajorCooldown(core.MajorCooldown{
		Spell: avatar,
		Type:  core.CooldownTypeDPS,
	})
}

func (war *Warrior) registerBloodbath() {
	if !war.Talents.Bloodbath {
		return
	}

	spellActionID := core.ActionID{SpellID: 12292}
	dotActionID := core.ActionID{SpellID: 113344}

	aura := war.RegisterAura(core.Aura{
		Label:    "Bloodbath",
		ActionID: spellActionID,
		Duration: 12 * time.Second,
	})

	shared.RegisterIgniteEffect(&war.Unit, shared.IgniteConfig{
		ActionID:       dotActionID,
		ClassSpellMask: SpellMaskBloodbathDot,
		DotAuraLabel:   "Bloodbath Dot",
		DotAuraTag:     "BloodbathDot",
		TickLength:     1 * time.Second,
		NumberOfTicks:  6,
		ParentAura:     aura,

		ProcTrigger: core.ProcTrigger{
			Name:     "Bloodbath - Trigger",
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: core.ProcMaskMeleeSpecial,
			Outcome:  core.OutcomeLanded,
		},

		DamageCalculator: func(result *core.SpellResult) float64 {
			return result.Damage * 0.3
		},
	})

	spell := war.RegisterSpell(core.SpellConfig{
		ActionID:       spellActionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ClassSpellMask: SpellMaskBloodbath,
		Flags:          core.SpellFlagAPL | core.SpellFlagReadinessTrinket,
		ProcMask:       core.ProcMaskEmpty,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: time.Minute * 1,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   war.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			aura.Activate(sim)
		},
	})

	war.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}

func (war *Warrior) registerStormBolt() {
	if !war.Talents.StormBolt {
		return
	}

	actionID := core.ActionID{SpellID: 107570}

	stormBoltOH := war.RegisterSpell(core.SpellConfig{
		ActionID:       actionID.WithTag(2),
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeOHSpecial,
		ClassSpellMask: SpellMaskStormBoltOH,
		MaxRange:       30,

		DamageMultiplier: 5,
		ThreatMultiplier: 1,
		CritMultiplier:   war.DefaultCritMultiplier(),

		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)
		},
	})

	war.RegisterSpell(core.SpellConfig{
		ActionID:       actionID.WithTag(1),
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagReadinessTrinket,
		ClassSpellMask: SpellMaskStormBolt,
		MaxRange:       30,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: 30 * time.Second,
			},
		},

		DamageMultiplier: 5,
		ThreatMultiplier: 1,
		CritMultiplier:   war.DefaultCritMultiplier(),

		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)

			if result.Landed() && war.Spec == proto.Spec_SpecFuryWarrior && war.OffHand() != nil && war.OffHand().WeaponType != proto.WeaponType_WeaponTypeUnknown {
				stormBoltOH.Cast(sim, target)
			}
		},
	})
}
