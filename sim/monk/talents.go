package monk

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func (monk *Monk) ApplyTalents() {
	// Level 15
	monk.registerCelerity()
	monk.registerTigersLust()
	monk.registerMomentum()

	// Level 30
	monk.registerChiWave()
	monk.registerZenSphere()
	monk.registerChiBurst()

	// Level 45
	monk.registerPowerStrikes()
	monk.registerAscension()
	monk.registerChiBrew()

	// Level 75
	monk.registerHealingElixirs()
	monk.registerDampenHarm()
	monk.registerDiffuseMagic()

	// Level 90
	monk.registerRushingJadeWind()
	monk.registerInvokeXuenTheWhiteTiger()
	monk.registerChiTorpedo()
}

func (monk *Monk) registerCelerity() {
}

func (monk *Monk) registerTigersLust() {
}

func (monk *Monk) registerMomentum() {
}

/*
Tooltip:
You cause a wave of Chi energy to flow through friend and foe, dealing $<damage> Nature damage or $<healing> healing. Bounces up to 7 times to the nearest targets within 25 yards.

When bouncing to allies, Chi Wave will prefer those injured over full health.

$damage=${<avg>+$ap*0.45}
$healing=${<avg>+$ap*0.45}
*/
var chiWaveActionID = core.ActionID{SpellID: 115098}
var chiWaveDamageActionID = core.ActionID{SpellID: 132467}
var chiWaveHealActionID = core.ActionID{SpellID: 132463}
var chiWaveMaxBounces = 7
var chiWaveBonusCoeff = 0.45
var chiWaveScaling = core.CalcScalingSpellAverageEffect(proto.Class_ClassMonk, chiWaveBonusCoeff)

func chiWaveSpellConfig(_ *Monk, isSEFClone bool, overrides core.SpellConfig) core.SpellConfig {
	config := core.SpellConfig{
		ActionID:       chiWaveActionID,
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MonkSpellChiWave,
		MaxRange:       40,

		Cast: overrides.Cast,

		ApplyEffects: overrides.ApplyEffects,
	}

	if isSEFClone {
		config.ActionID = config.ActionID.WithTag(SEFSpellID)
	}

	return config
}
func chiWaveDamageSpellConfig(monk *Monk, isSEFClone bool, overrides core.SpellConfig) core.SpellConfig {
	config := core.SpellConfig{
		ActionID:       chiWaveDamageActionID,
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagPassiveSpell,
		ClassSpellMask: MonkSpellChiWave,
		MaxRange:       40,
		MissileSpeed:   8,

		DamageMultiplier: 1.0,
		ThreatMultiplier: 1.0,
		CritMultiplier:   monk.DefaultCritMultiplier(),

		ApplyEffects: overrides.ApplyEffects,
	}

	if isSEFClone {
		config.ActionID = config.ActionID.WithTag(SEFSpellID)
	}

	return config
}

func chiWaveHealSpellConfig(monk *Monk, isSEFClone bool, overrides core.SpellConfig) core.SpellConfig {
	config := core.SpellConfig{
		ActionID:       chiWaveHealActionID,
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellHealing,
		Flags:          core.SpellFlagHelpful | core.SpellFlagPassiveSpell,
		ClassSpellMask: MonkSpellChiWave,
		MissileSpeed:   8,

		DamageMultiplier: 1.0,
		ThreatMultiplier: 1.0,
		CritMultiplier:   monk.DefaultCritMultiplier(),

		ApplyEffects: overrides.ApplyEffects,
	}

	if isSEFClone {
		config.ActionID = config.ActionID.WithTag(SEFSpellID)
	}

	return config
}

func (monk *Monk) registerChiWave() {
	if !monk.Talents.ChiWave {
		return
	}

	var nextTarget *core.Unit
	tickIndex := 0

	var chiWaveHealingSpell *core.Spell
	chiWaveDamageSpell := monk.RegisterSpell(chiWaveDamageSpellConfig(monk, false, core.SpellConfig{
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := chiWaveScaling + spell.MeleeAttackPower()*chiWaveBonusCoeff

			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialNoBlockDodgeParryNoCritNoHitCounter)
			if result.Landed() {
				spell.WaitTravelTime(sim, func(sim *core.Simulation) {
					result = spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicCrit)
					if tickIndex < chiWaveMaxBounces {
						tickIndex++
						nextTarget = nextTarget.Env.NextTargetUnit(nextTarget)
						chiWaveHealingSpell.Cast(sim, &monk.Unit)
					}
				})
			}
		},
	}))

	chiWaveHealingSpell = monk.RegisterSpell(chiWaveHealSpellConfig(monk, false, core.SpellConfig{
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseHealing := chiWaveScaling + spell.MeleeAttackPower()*chiWaveBonusCoeff

			result := spell.CalcHealing(sim, target, baseHealing, spell.OutcomeHealingCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealHealing(sim, result)

				if tickIndex < chiWaveMaxBounces {
					tickIndex++
					chiWaveDamageSpell.Cast(sim, nextTarget)
				}
			})
		},
	}))

	monk.RegisterSpell(chiWaveSpellConfig(monk, false, core.SpellConfig{
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    monk.NewTimer(),
				Duration: time.Second * 15,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			tickIndex = 0
			if monk.IsOpponent(target) {
				nextTarget = target.Env.NextTargetUnit(target)
				chiWaveDamageSpell.Cast(sim, target)
			} else {
				nextTarget = target.CurrentTarget
				chiWaveHealingSpell.Cast(sim, target)
			}
		},
	}))
}

func (pet *StormEarthAndFirePet) registerSEFChiWave() {
	if pet.owner.Spec != proto.Spec_SpecWindwalkerMonk || !pet.owner.Talents.ChiWave {
		return
	}

	var nextTarget *core.Unit
	tickIndex := 0

	var chiWaveHealingSpell *core.Spell
	chiWaveDamageSpell := pet.RegisterSpell(chiWaveDamageSpellConfig(pet.owner, true, core.SpellConfig{
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := chiWaveScaling + spell.MeleeAttackPower()*chiWaveBonusCoeff

			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialNoBlockDodgeParryNoCritNoHitCounter)

			if result.Landed() {
				spell.WaitTravelTime(sim, func(sim *core.Simulation) {
					spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicCrit)
					if tickIndex < chiWaveMaxBounces {
						tickIndex++
						nextTarget = nextTarget.Env.NextTargetUnit(nextTarget)
						chiWaveHealingSpell.Cast(sim, &pet.Unit)
					}
				})
			}
		},
	}))

	chiWaveHealingSpell = pet.RegisterSpell(chiWaveHealSpellConfig(pet.owner, true, core.SpellConfig{
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseHealing := chiWaveScaling + spell.MeleeAttackPower()*chiWaveBonusCoeff

			result := spell.CalcHealing(sim, target, baseHealing, spell.OutcomeHealingCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealHealing(sim, result)

				if tickIndex < chiWaveMaxBounces {
					tickIndex++
					chiWaveDamageSpell.Cast(sim, nextTarget)
				}
			})
		},
	}))

	pet.RegisterSpell(chiWaveSpellConfig(pet.owner, true, core.SpellConfig{
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			tickIndex = 0
			if pet.IsOpponent(target) {
				nextTarget = target.Env.NextTargetUnit(target)
				chiWaveDamageSpell.Cast(sim, target)
			} else {
				nextTarget = target.CurrentTarget
				chiWaveHealingSpell.Cast(sim, target)
			}
		},
	}))
}

/*
Tooltip:
Forms a Zen Sphere above the target, healing the target for $<healingperiodic> and dealing $<damageperiodic> Nature damage to the nearest enemy within 10 yards of the target every 2 sec for 16 sec. Only two Zen Spheres can be summoned at any one time.

If the target of the Zen Sphere reaches 35% or lower health or if the Zen Sphere is dispelled or expires it will detonate, dealing $<damagedetonate> Nature damage and $<healingdetonate> healing to all targets within 10 yards.

$damageperiodic=${(<avg>+$ap*0.09)}
$healingperiodic=${(<avg>+$ap*0.09)}
$damagedetonate=${(<avg>+$ap*0.368)}
$healingdetonate=${(<avg>+$ap*0.234)}
*/

func (monk *Monk) registerZenSphere() {
	if !monk.Talents.ZenSphere {
		return
	}

	targetDummies := monk.Env.Raid.GetTargetDummies()
	maxTargets := int32(max(1, len(targetDummies)))

	zenSphereAura := monk.RegisterAura(core.Aura{
		Label:     "Zen Sphere" + monk.Label,
		ActionID:  core.ActionID{SpellID: 124081}.WithTag(1),
		Duration:  core.NeverExpires,
		MaxStacks: maxTargets,
	})

	avgTickScaling := monk.CalcScalingSpellDmg(0.1040000021)
	// The 15% extra is from a hotfix not represented in the tooltip.
	avgTickBonusCoefficient := 0.09 * 1.15

	avgDetonateHealScaling := monk.CalcScalingSpellDmg(0.2689999938)
	// The 15% extra is from a hotfix not represented in the tooltip.
	avgDetonateHealBonusCoefficient := 0.234 * 1.15

	avgDetonateDmgScaling := monk.CalcScalingSpellDmg(0.4230000079)
	// The 15% extra is from a hotfix not represented in the tooltip.
	avgDetonateDmgBonusCoefficient := 0.368 * 1.15

	detonateDamageSpell := monk.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 124081}.WithTag(5), // Real Spell ID: 125033
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAoE | core.SpellFlagPassiveSpell,
		ClassSpellMask: MonkSpellZenSphere,
		MaxRange:       10,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   monk.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, target := range sim.Encounter.TargetUnits {
				baseDamage := avgDetonateDmgScaling + spell.MeleeAttackPower()*avgDetonateDmgBonusCoefficient
				result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialNoBlockDodgeParryNoCritNoHitCounter)

				if result.Landed() {
					spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicCrit)
				}
			}
		},
	})

	detonateHealingSpell := monk.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 124081}.WithTag(4), // Real Spell ID: 124101
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellHealing,
		Flags:          core.SpellFlagHelpful | core.SpellFlagPassiveSpell,
		ClassSpellMask: MonkSpellZenSphere,
		MaxRange:       10,

		DamageMultiplier: 1.0,
		ThreatMultiplier: 1.0,
		CritMultiplier:   monk.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseHealing := avgDetonateHealScaling + spell.MeleeAttackPower()*avgDetonateHealBonusCoefficient
			spell.CalcAndDealHealing(sim, target, baseHealing, spell.OutcomeHealingCrit)
		},
	})

	zenSphereDotTick := monk.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 124081}.WithTag(3), // Real Spell ID: 124098
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagPassiveSpell,
		ClassSpellMask: MonkSpellZenSphere,
		MaxRange:       10,

		DamageMultiplier: 1.0,
		ThreatMultiplier: 1.0,
		CritMultiplier:   monk.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			healValue := avgTickScaling + spell.MeleeAttackPower()*avgTickBonusCoefficient
			result := spell.CalcDamage(sim, target, healValue, spell.OutcomeTickMagicHitAndCrit)
			spell.DealPeriodicDamage(sim, result)
		},
	})

	monk.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 124081},
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellHealing,
		Flags:          core.SpellFlagAPL | core.SpellFlagHelpful,
		ClassSpellMask: MonkSpellZenSphere,
		MaxRange:       40,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    monk.NewTimer(),
				Duration: time.Second * 10,
			},
		},

		DamageMultiplier: 1.0,
		ThreatMultiplier: 1.0,
		CritMultiplier:   monk.DefaultCritMultiplier(),

		Hot: core.DotConfig{
			Aura: core.Aura{
				Label:    "Zen Sphere (Heal)" + monk.Label,
				ActionID: core.ActionID{SpellID: 124081}.WithTag(2),
				OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if aura.Unit.CurrentHealthPercent() <= 0.35 {
						aura.Deactivate(sim)
					}
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					detonateHealingSpell.Cast(sim, aura.Unit)
					detonateDamageSpell.Cast(sim, aura.Unit)
					if zenSphereAura.IsActive() {
						zenSphereAura.RemoveStack(sim)
					}
				},
			},
			NumberOfTicks:       8,
			TickLength:          time.Second * 2,
			AffectedByCastSpeed: false,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.SnapshotBaseDamage = avgTickScaling + dot.Spell.MeleeAttackPower()*avgTickBonusCoefficient
				dot.SnapshotAttackerMultiplier = dot.Spell.CasterHealingMultiplier()
				dot.CalcAndDealPeriodicSnapshotHealing(sim, target, dot.OutcomeTickHealingCrit)
				dot.Spell.RelatedDotSpell.Cast(sim, target.CurrentTarget)
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !zenSphereAura.IsActive() || zenSphereAura.GetStacks() < zenSphereAura.MaxStacks
		},

		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			var target *core.Unit

			if len(targetDummies) > 1 {
				for _, dummy := range targetDummies {
					unit := &dummy.Unit
					if !spell.Hot(unit).IsActive() {
						target = unit
						break
					}
				}
			}

			if target == nil {
				return
			}

			if target.CurrentHealthPercent() <= 0.35 {
				detonateHealingSpell.Cast(sim, target)
				detonateDamageSpell.Cast(sim, target.CurrentTarget)
				return
			}

			zenSphereAura.Activate(sim)
			zenSphereAura.AddStack(sim)
			spell.Hot(target).Activate(sim)
		},
		RelatedDotSpell: zenSphereDotTick,
	})
}

/*
Tooltip:
You summon a torrent of Chi energy and hurl it forward, up to 40 yds, dealing $<damage> Nature damage to all enemies, and $<healing> healing to all allies in its path. Chi Burst will always heal the Monk.

While casting Chi Burst, you continue to dodge, parry, and auto-attack.

$damage=${<avg>+$ap*1.21}
$healing=${<avg>+$ap}
*/
var chiBurstActionID = core.ActionID{SpellID: 123986}
var chiBurstDamageActionID = core.ActionID{SpellID: 148135}
var chiBurstHealActionID = core.ActionID{SpellID: 130654}
var chiBurstBonusCoeff = 1.21
var chiBurstScaling = core.CalcScalingSpellAverageEffect(proto.Class_ClassMonk, chiWaveBonusCoeff)

func chiBurstDamageSpellConfig(monk *Monk, isSEFClone bool) core.SpellConfig {

	config := core.SpellConfig{
		ActionID:       chiBurstDamageActionID,
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAoE | core.SpellFlagPassiveSpell,
		ClassSpellMask: MonkSpellChiBurst,
		MissileSpeed:   30,
		MaxRange:       40,

		DamageMultiplier: 1.0,
		ThreatMultiplier: 1.0,
		CritMultiplier:   monk.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {

			spell.WaitTravelTime(sim, func(simulation *core.Simulation) {
				for _, target := range sim.Encounter.TargetUnits {
					baseDamage := chiBurstScaling + spell.MeleeAttackPower()*chiBurstBonusCoeff
					result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialNoBlockDodgeParryNoCritNoHitCounter)

					if result.Landed() {
						spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicCrit)
					}
				}
			})

		},
	}

	if isSEFClone {
		config.ActionID = config.ActionID.WithTag(SEFSpellID)
	}

	return config

}
func chiBurstHealSpellConfig(monk *Monk, isSEFClone bool) core.SpellConfig {
	config := core.SpellConfig{
		ActionID:       chiBurstHealActionID,
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellHealing,
		Flags:          core.SpellFlagHelpful | core.SpellFlagPassiveSpell,
		ClassSpellMask: MonkSpellChiBurst,
		MissileSpeed:   30,

		DamageMultiplier: 1.0,
		ThreatMultiplier: 1.0,
		CritMultiplier:   monk.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseHealing := monk.CalcScalingSpellDmg(1) + spell.MeleeAttackPower()

			result := spell.CalcHealing(sim, target, baseHealing, spell.OutcomeHealingCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealHealing(sim, result)
			})
		},
	}

	if isSEFClone {
		config.ActionID = config.ActionID.WithTag(SEFSpellID)
	}

	return config
}

func (monk *Monk) registerChiBurst() {
	if !monk.Talents.ChiBurst {
		return
	}

	chiBurstDamageSpell := monk.RegisterSpell(chiBurstDamageSpellConfig(monk, false))

	chiBurstHealingSpell := monk.RegisterSpell(chiBurstHealSpellConfig(monk, false))

	chiBurstFakeCastAura := monk.RegisterAura(core.Aura{
		Label:    "Chi Burst" + monk.Label,
		ActionID: chiBurstActionID,
		Duration: time.Second,
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			chiBurstDamageSpell.Cast(sim, monk.CurrentTarget)
			chiBurstHealingSpell.Cast(sim, &monk.Unit)
		},
	})

	monk.RegisterSpell(core.SpellConfig{
		ActionID:       chiBurstActionID,
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MonkSpellChiBurst,
		MaxRange:       40,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    monk.NewTimer(),
				Duration: time.Second * 30,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			chiBurstFakeCastAura.Activate(sim)
		},
	})
}

func (pet *StormEarthAndFirePet) registerSEFChiBurst() {
	if pet.owner.Spec != proto.Spec_SpecWindwalkerMonk || !pet.owner.Talents.ChiBurst {
		return
	}

	pet.RegisterSpell(chiBurstDamageSpellConfig(pet.owner, true))
	pet.RegisterSpell(chiBurstHealSpellConfig(pet.owner, true))
}

func (monk *Monk) registerPowerStrikes() {
	if !monk.Talents.PowerStrikes {
		return
	}

	chiSphereSpellActionID := core.ActionID{SpellID: 121283}
	chiSphereChiMetrics := monk.NewChiMetrics(chiSphereSpellActionID)

	hasGlyph := monk.HasMajorGlyph(proto.MonkMajorGlyph_GlyphOfEnduringHealingSphere)
	chiSphereduration := time.Minute*2 + core.TernaryDuration(hasGlyph, time.Minute*3, 0)

	monk.ChiSphereAura = monk.RegisterAura(core.Aura{
		Label:     "Chi Sphere" + monk.Label,
		ActionID:  core.ActionID{SpellID: 121286},
		Duration:  time.Minute * chiSphereduration,
		MaxStacks: 10,
	})

	chiSphereUseAura := monk.RegisterAura(core.Aura{
		Label:    "Chi Sphere (Use)" + monk.Label,
		ActionID: chiSphereSpellActionID,
		Duration: core.NeverExpires,
	})

	monk.RegisterSpell(core.SpellConfig{
		ActionID:       chiSphereSpellActionID,
		SpellSchool:    core.SpellSchoolPhysical,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagAPL,
		ClassSpellMask: MonkSpellChiSphere,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return monk.GetChi() != monk.GetMaxChi() && !chiSphereUseAura.IsActive() && monk.ChiSphereAura.IsActive() && monk.ChiSphereAura.GetStacks() > 0
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			chiSphereUseAura.Activate(sim)

			// Orbs spawn ~4 yards away, simulate movement to grab the sphere.
			moveDuration := core.DurationFromSeconds(4.0 / monk.GetMovementSpeed())
			monk.MoveDuration(moveDuration, sim)
			sim.AddPendingAction(&core.PendingAction{
				NextActionAt: sim.CurrentTime + moveDuration,
				OnAction: func(sim *core.Simulation) {
					monk.ChiSphereAura.RemoveStack(sim)
					chiSphereUseAura.Deactivate(sim)
					monk.AddChi(sim, spell, 1, chiSphereChiMetrics)
				},
			})
		},
	})

	powerStrikesAuraActionID := core.ActionID{SpellID: 129914}
	monk.PowerStrikesChiMetrics = monk.NewChiMetrics(powerStrikesAuraActionID)

	monk.PowerStrikesAura = monk.RegisterAura(core.Aura{
		Label:    "Power Strikes" + monk.Label,
		ActionID: powerStrikesAuraActionID,
		Duration: core.NeverExpires,
	})

	monk.RegisterResetEffect(func(sim *core.Simulation) {
		// Start at a random time
		startAt := sim.RandomFloat("Power Strikes Start") * 20.0
		sim.AddPendingAction(&core.PendingAction{
			NextActionAt: core.DurationFromSeconds(startAt),
			OnAction: func(sim *core.Simulation) {
				core.StartPeriodicAction(sim, core.PeriodicActionOptions{
					Period:          time.Second * 20,
					Priority:        core.ActionPriorityLow,
					TickImmediately: true,
					OnAction: func(sim *core.Simulation) {
						monk.PowerStrikesAura.Activate(sim)
					},
				})
			},
		})
	})
}

func (monk *Monk) TriggerPowerStrikes(sim *core.Simulation) {
	if !monk.PowerStrikesAura.IsActive() {
		return
	}

	if monk.GetChi() == monk.GetMaxChi() {
		monk.ChiSphereAura.Activate(sim)
		monk.ChiSphereAura.AddStack(sim)
	} else {
		monk.AddChi(sim, nil, 1, monk.PowerStrikesChiMetrics)
	}

	monk.PowerStrikesAura.Deactivate(sim)
}

func (monk *Monk) registerAscension() {
	if !monk.Talents.Ascension {
		return
	}

	core.MakePermanent(monk.GetOrRegisterAura(core.Aura{
		Label:    "Ascension" + monk.Label,
		ActionID: core.ActionID{SpellID: 115396},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			monk.ApplyAdditiveEnergyRegenBonus(sim, 0.15)
			monk.SetMaxComboPoints(5)

			if monk.HasManaBar() {
				monk.MultiplyStat(stats.Mana, 1.15)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			monk.ApplyAdditiveEnergyRegenBonus(sim, -0.15)
			monk.SetMaxComboPoints(4)

			if monk.HasManaBar() {
				monk.MultiplyStat(stats.Mana, 1.0/1.15)
			}
		},
	}))
}

func (monk *Monk) registerChiBrew() {
	if !monk.Talents.ChiBrew {
		return
	}

	actionID := core.ActionID{SpellID: 115399}
	chiMetrics := monk.NewChiMetrics(actionID)
	monk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,
		ClassSpellMask: MonkSpellChiBrew,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},

		Charges:      2,
		RechargeTime: time.Second * 45,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// TODO: Add 2 Mana Tea stacks for Mistweavers

			if monk.Spec == proto.Spec_SpecBrewmasterMonk {
				monk.AddBrewStacks(sim, 5)
			} else if monk.Spec == proto.Spec_SpecWindwalkerMonk {
				monk.AddBrewStacks(sim, 2)
			}

			monk.AddChi(sim, spell, 2, chiMetrics)
		},
	})
}

func (monk *Monk) registerHealingElixirs() {
	if !monk.Talents.HealingElixirs {
		return
	}
}

func (monk *Monk) registerDampenHarm() {
	if !monk.Talents.DampenHarm {
		return
	}

	actionId := core.ActionID{SpellID: 122278}

	monk.DampenHarmAura = monk.RegisterAura(core.Aura{
		Label:     "Dampen Harm" + monk.Label,
		ActionID:  actionId.WithTag(1),
		MaxStacks: 3,
		Duration:  time.Second * 45,
	})

	// Dampen Harms Damage Reduction for BRM is implemented in stagger.go
	if monk.Spec != proto.Spec_SpecBrewmasterMonk {
		monk.AddDynamicDamageTakenModifier(func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult, isPeriodic bool) {
			if !monk.DampenHarmAura.IsActive() || !result.Landed() || result.Damage < result.Target.MaxHealth()*0.2 {
				return
			}

			monk.DampenHarmAura.RemoveStack(sim)
			result.Damage /= 2
		})
	}

	spell := monk.RegisterSpell(core.SpellConfig{
		ActionID:       actionId,
		SpellSchool:    core.SpellSchoolPhysical,
		ClassSpellMask: MonkSpellDampenHarm,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			CD: core.Cooldown{
				Timer:    monk.NewTimer(),
				Duration: 90 * time.Second,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			monk.DampenHarmAura.Activate(sim)
			monk.DampenHarmAura.SetStacks(sim, 3)
		},
	})

	monk.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeSurvival,
		ShouldActivate: func(_ *core.Simulation, _ *core.Character) bool {
			return monk.Spec == proto.Spec_SpecBrewmasterMonk && monk.CurrentHealthPercent() < 0.5
		},
	})
}

func (monk *Monk) registerDiffuseMagic() {
	if !monk.Talents.DiffuseMagic {
		return
	}
}

/*
Tooltip:
You summon a whirling tornado around you which deals ${1.59*(1.4/1.59)*$<low>} to ${1.59*(1.4/1.59)*$<high>} damage to all nearby enemies

-- Teachings of the Monastery --
and $117640m1 healing to nearby allies
-- Teachings of the Monastery --

every 0.75 sec, within 8 yards. Generates 1 Chi, if it hits at least 3 targets. Lasts 6 sec.

Replaces Spinning Crane Kick.
*/
var rushingJadeWindActionID = core.ActionID{SpellID: 116847}
var rushingJadeWindDebuffActionID = core.ActionID{SpellID: 148187}

func rushingJadeWindTickSpellConfig(monk *Monk, isSEFClone bool) core.SpellConfig {
	config := core.SpellConfig{
		ActionID:       rushingJadeWindDebuffActionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagAoE | core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell,
		ClassSpellMask: MonkSpellRushingJadeWind,
		MaxRange:       8,

		DamageMultiplier: 1.4, // 1.59 * (1.4 / 1.59),
		ThreatMultiplier: 1,
		CritMultiplier:   monk.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, target := range sim.Encounter.TargetUnits {
				baseDamage := monk.CalculateMonkStrikeDamage(sim, spell)
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			}
		},
	}

	if isSEFClone {
		config.ActionID = config.ActionID.WithTag(SEFSpellID)
	}

	return config

}

func rushingJadeWindSpellConfig(monk *Monk, isSEFClone bool, overrides core.SpellConfig) core.SpellConfig {
	config := core.SpellConfig{
		ActionID:       rushingJadeWindActionID,
		Flags:          SpellFlagBuilder | core.SpellFlagAPL,
		ClassSpellMask: MonkSpellRushingJadeWind,

		EnergyCost: overrides.EnergyCost,
		ManaCost:   overrides.ManaCost,
		Cast:       overrides.Cast,

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label:    "Rushing Jade Wind" + monk.Label,
				ActionID: rushingJadeWindDebuffActionID,
				OnInit:   overrides.Dot.Aura.OnInit,
			},
			NumberOfTicks:        8,
			TickLength:           time.Millisecond * 750,
			AffectedByCastSpeed:  true,
			HasteReducesDuration: true,

			OnTick: overrides.Dot.OnTick,
		},

		ApplyEffects: overrides.ApplyEffects,
	}

	if isSEFClone {
		config.ActionID = config.ActionID.WithTag(SEFSpellID)
		config.Flags &= ^(core.SpellFlagAPL | SpellFlagBuilder)
	}

	return config
}

func (monk *Monk) registerRushingJadeWind() {
	if !monk.Talents.RushingJadeWind {
		return
	}

	chiMetrics := monk.NewChiMetrics(rushingJadeWindActionID)
	numTargets := monk.Env.GetNumTargets()
	baseCooldown := time.Second * 6

	rushingJadeWindTickSpell := monk.RegisterSpell(rushingJadeWindTickSpellConfig(monk, false))

	rushingJadeWindBuff := monk.RegisterAura(core.Aura{
		Label:    "Rushing Jade Wind" + monk.Label,
		ActionID: rushingJadeWindActionID,
		Duration: baseCooldown,
	})

	isWiseSerpent := monk.StanceMatches(WiseSerpent)
	monk.RegisterSpell(rushingJadeWindSpellConfig(monk, false, core.SpellConfig{
		EnergyCost: core.EnergyCostOptions{
			Cost: core.TernaryInt32(isWiseSerpent, 0, 40),
		},
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: core.TernaryFloat64(isWiseSerpent, 7.15, 0),
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    monk.NewTimer(),
				Duration: baseCooldown,
			},
		},

		Dot: core.DotConfig{
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				rushingJadeWindTickSpell.Cast(sim, target)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dot := spell.AOEDot()
			dot.Apply(sim)
			dot.TickOnce(sim)

			remainingDuration := dot.RemainingDuration(sim)
			spell.CD.Set(sim.CurrentTime + remainingDuration)
			rushingJadeWindBuff.Duration = remainingDuration
			rushingJadeWindBuff.Activate(sim)

			if numTargets >= 3 {
				monk.AddChi(sim, spell, 1, chiMetrics)
			}
		},
	}))
}

func (pet *StormEarthAndFirePet) registerSEFRushingJadeWind() {
	if pet.owner.Spec != proto.Spec_SpecWindwalkerMonk || !pet.owner.Talents.RushingJadeWind {
		return
	}

	rushingJadeWindTickSpell := pet.RegisterSpell(rushingJadeWindTickSpellConfig(pet.owner, true))

	pet.RegisterSpell(rushingJadeWindSpellConfig(pet.owner, true, core.SpellConfig{
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			IgnoreHaste: true,
		},

		Dot: core.DotConfig{
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				rushingJadeWindTickSpell.Cast(sim, target)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dot := spell.AOEDot()
			dot.Apply(sim)
			dot.TickOnce(sim)
		},
	}))
}

func (monk *Monk) registerInvokeXuenTheWhiteTiger() {
	if !monk.Talents.InvokeXuenTheWhiteTiger {
		return
	}

	actionID := core.ActionID{SpellID: 123904}

	// For timeline only
	monk.XuenAura = monk.RegisterAura(core.Aura{
		ActionID: actionID,
		Label:    "Xuen, the White Tiger",
		Duration: time.Second * 45.0,
	})

	spell := monk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MonkSpellInvokeXuenTheWhiteTiger,

		MaxRange: 40,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    monk.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			monk.XuenPet.EnableWithTimeout(sim, monk.XuenPet, time.Second*45.0)
			monk.XuenAura.Activate(sim)
		},
	})

	monk.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})

}

func (monk *Monk) registerChiTorpedo() {
}
