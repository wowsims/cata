package monk

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

func (monk *Monk) ApplyTalents() {
	if monk.Level >= 15 {
		monk.registerCelerity()
		monk.registerTigersLust()
		monk.registerMomentum()
	}

	if monk.Level >= 30 {
		monk.registerChiWave()
		monk.registerZenSphere()
		monk.registerChiBurst()
	}

	if monk.Level >= 45 {
		monk.registerPowerStrikes()
		monk.registerAscension()
		monk.registerChiBrew()
	}

	if monk.Level >= 75 {
		monk.registerHealingElixirs()
		monk.registerDampenHarm()
		monk.registerDiffuseMagic()
	}

	if monk.Level >= 90 {
		monk.registerRushingJadeWind()
		monk.registerInvokeXuenTheWhiteTiger()
		monk.registerChiTorpedo()
	}
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
func (monk *Monk) registerChiWave() {
	if !monk.Talents.ChiWave {
		return
	}

	avgScaling := monk.CalcScalingSpellDmg(0.4499999881)
	var nextTarget *core.Unit
	tickIndex := 0

	var chiWaveDamageSpell *core.Spell
	var chiWaveHealingSpell *core.Spell
	chiWaveDamageSpell = monk.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 132467},
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagPassiveSpell,
		ClassSpellMask: MonkSpellChiWave,
		MissileSpeed:   8,

		DamageMultiplier: 1.0,
		ThreatMultiplier: 1.0,
		CritMultiplier:   monk.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := avgScaling + spell.MeleeAttackPower()*0.45

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)

				if tickIndex < 7 {
					tickIndex++

					nextTarget = nextTarget.Env.NextTargetUnit(nextTarget)
					chiWaveHealingSpell.Cast(sim, &monk.Unit)
				}
			})
		},
	})

	chiWaveHealingSpell = monk.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 132463},
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellHealing,
		Flags:          core.SpellFlagHelpful | core.SpellFlagPassiveSpell,
		ClassSpellMask: MonkSpellChiWave,
		MissileSpeed:   8,

		DamageMultiplier: 1.0,
		ThreatMultiplier: 1.0,
		CritMultiplier:   monk.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseHealing := avgScaling + spell.MeleeAttackPower()*0.45

			result := spell.CalcHealing(sim, target, baseHealing, spell.OutcomeHealingCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealHealing(sim, result)

				if tickIndex < 7 {
					tickIndex++
					chiWaveDamageSpell.Cast(sim, nextTarget)
				}
			})
		},
	})

	monk.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 115098},
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MonkSpellChiWave,
		MaxRange:       40,

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
	})
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

	numTargets := monk.Env.GetNumTargets()

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
		ActionID:       core.ActionID{SpellID: 125033},
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagPassiveSpell,
		ClassSpellMask: MonkSpellZenSphere,
		MaxRange:       10,

		DamageMultiplier: 1.0,
		ThreatMultiplier: 1.0,
		CritMultiplier:   monk.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := avgDetonateDmgScaling + spell.MeleeAttackPower()*avgDetonateDmgBonusCoefficient

			results := make([]*core.SpellResult, numTargets)

			for idx := int32(0); idx < numTargets; idx++ {
				currentTarget := sim.Environment.GetTargetUnit(idx)
				results[idx] = spell.CalcDamage(sim, currentTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}

			for idx := int32(0); idx < numTargets; idx++ {
				spell.DealDamage(sim, results[idx])
			}
		},
	})

	detonateHealingSpell := monk.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 124101},
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

	var currentTargetIndex int32
	zenSphereDotSpell := monk.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 124098},
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagPassiveSpell,
		ClassSpellMask: MonkSpellZenSphere,
		MaxRange:       10,

		DamageMultiplier: 1.0,
		ThreatMultiplier: 1.0,
		CritMultiplier:   monk.DefaultCritMultiplier(),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Zen Sphere (Damage)" + monk.Label,
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					detonateDamageSpell.Cast(sim, aura.Unit)
				},
			},
			NumberOfTicks:       8,
			TickLength:          time.Second * 2,
			AffectedByCastSpeed: false,
			OnTick: func(sim *core.Simulation, _ *core.Unit, dot *core.Dot) {
				healValue := avgTickScaling + dot.Spell.MeleeAttackPower()*avgTickBonusCoefficient
				target := sim.GetTargetUnit(currentTargetIndex)
				dot.Spell.CalcAndDealPeriodicDamage(sim, target, healValue, dot.OutcomeTickMagicCrit)
				currentTargetIndex = sim.NextTargetUnit(target).UnitIndex
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			currentTargetIndex = target.UnitIndex
			spell.Dot(target).Activate(sim)
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
				Label: "Zen Sphere (Heal)" + monk.Label,
				OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if aura.Unit.CurrentHealthPercent() <= 0.35 {
						aura.Deactivate(sim)
						zenSphereDotSpell.Dot(aura.Unit.CurrentTarget).Deactivate(sim)
					}
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					detonateHealingSpell.Cast(sim, aura.Unit)
				},
			},
			NumberOfTicks:       8,
			TickLength:          time.Second * 2,
			AffectedByCastSpeed: false,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				healValue := avgTickScaling + dot.Spell.MeleeAttackPower()*avgTickBonusCoefficient
				dot.Spell.CalcAndDealPeriodicHealing(sim, target, healValue, dot.OutcomeTickHealingCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			target := spell.Unit

			if target.CurrentHealthPercent() <= 0.35 {
				detonateHealingSpell.Cast(sim, target)
				detonateDamageSpell.Cast(sim, target.CurrentTarget)
				return
			}

			spell.Hot(target).Activate(sim)
			zenSphereDotSpell.Cast(sim, target.CurrentTarget)
		},
	})
}

/*
Tooltip:
You summon a torrent of Chi energy and hurl it forward, up to 40 yds, dealing $<damage> Nature damage to all enemies, and $<healing> healing to all allies in its path. Chi Burst will always heal the Monk.

While casting Chi Burst, you continue to dodge, parry, and auto-attack.

$damage=${<avg>+$ap*1.21}
$healing=${<avg>+$ap}
*/
func (monk *Monk) registerChiBurst() {
	if !monk.Talents.ChiBurst {
		return
	}

	numTargets := monk.Env.GetNumTargets()
	avgDmgScaling := monk.CalcScalingSpellDmg(1.2100000381)

	chiBurstDamageSpell := monk.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 148135},
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagPassiveSpell,
		ClassSpellMask: MonkSpellChiBurst,
		MissileSpeed:   30,

		DamageMultiplier: 1.0,
		ThreatMultiplier: 1.0,
		CritMultiplier:   monk.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := avgDmgScaling + spell.MeleeAttackPower()*1.21

			results := make([]*core.SpellResult, numTargets)

			for idx := int32(0); idx < numTargets; idx++ {
				currentTarget := sim.Environment.GetTargetUnit(idx)
				results[idx] = spell.CalcDamage(sim, currentTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}

			spell.WaitTravelTime(sim, func(simulation *core.Simulation) {
				for idx := int32(0); idx < numTargets; idx++ {
					spell.DealDamage(sim, results[idx])
				}
			})
		},
	})

	chiBurstHealingSpell := monk.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 130654},
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
	})

	actionID := core.ActionID{SpellID: 123986}
	chiBurstFakeCastAura := monk.RegisterAura(core.Aura{
		Label:    "Chi Burst" + monk.Label,
		ActionID: actionID,
		Duration: time.Second,
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			chiBurstDamageSpell.Cast(sim, monk.CurrentTarget)
			chiBurstHealingSpell.Cast(sim, &monk.Unit)
		},
	})

	monk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
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

func (monk *Monk) registerPowerStrikes() {
	if !monk.Talents.PowerStrikes {
		return
	}

	chiSphereSpellActionID := core.ActionID{SpellID: 121283}
	chiSphereChiMetrics := monk.NewChiMetrics(chiSphereSpellActionID)

	monk.ChiSphereAura = monk.RegisterAura(core.Aura{
		Label:     "Chi Sphere" + monk.Label,
		ActionID:  core.ActionID{SpellID: 121286},
		Duration:  time.Minute * 2,
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
			return !chiSphereUseAura.IsActive() && monk.ChiSphereAura.IsActive() && monk.ChiSphereAura.GetStacks() > 0
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

	if monk.ComboPoints() == monk.MaxComboPoints() {
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
			monk.MultiplyEnergyRegenSpeed(sim, 1.15)
			monk.SetMaxComboPoints(5)

			if monk.HasManaBar() {
				monk.MultiplyStat(stats.Mana, 1.15)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			monk.MultiplyEnergyRegenSpeed(sim, 1.0/1.15)
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

	chiBrewAura := monk.RegisterAura(core.Aura{
		Label:     "Chi Brew" + monk.Label,
		ActionID:  actionID,
		Duration:  core.NeverExpires,
		MaxStacks: 2,

		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			if monk.chiBrewRecharge != nil {
				monk.chiBrewRecharge.Cancel(sim)
			}

			aura.Activate(sim)
			aura.SetStacks(sim, 2)
		},

		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			if !aura.IsActive() {
				return
			}

			if newStacks < oldStacks {
				nextRecharge := &core.PendingAction{
					NextActionAt: sim.CurrentTime + time.Second*45,
					OnAction: func(sim *core.Simulation) {
						aura.Activate(sim)
						aura.AddStack(sim)
					},
				}

				if monk.chiBrewRecharge != nil {
					// If we have an existing stack recharging, set this new one as current when it's done.
					// This way we can always check next recharge time from the APL.
					oldAction := monk.chiBrewRecharge.OnAction
					monk.chiBrewRecharge.OnAction = func(sim *core.Simulation) {
						monk.chiBrewRecharge = nextRecharge
						oldAction(sim)
					}
				} else {
					monk.chiBrewRecharge = nextRecharge
				}

				sim.AddPendingAction(nextRecharge)
			} else if newStacks > oldStacks {
				if newStacks == 2 || monk.chiBrewRecharge.IsConsumed() {
					monk.chiBrewRecharge = nil
				}
			}
		},
	})

	monk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,
		ClassSpellMask: MonkSpellChiBrew,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return chiBrewAura.GetStacks() >= 1
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// TODO: Add 5 Elusive Brew stacks for Brewmasters
			// TODO: Add 2 Mana Tea stacks for Mistweavers

			monk.AddChi(sim, spell, 2, chiMetrics)
			monk.AddBrewStacks(sim, 2)

			chiBrewAura.RemoveStack(sim)
		},
	})
}

func (monk *Monk) registerHealingElixirs() {
}

func (monk *Monk) registerDampenHarm() {
}

func (monk *Monk) registerDiffuseMagic() {
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
func (monk *Monk) registerRushingJadeWind() {
	if !monk.Talents.RushingJadeWind {
		return
	}

	actionID := core.ActionID{SpellID: 116847}
	debuffActionID := core.ActionID{SpellID: 148187}
	chiMetrics := monk.NewChiMetrics(actionID)
	numTargets := monk.Env.GetNumTargets()

	rushingJadeWindTickSpell := monk.RegisterSpell(core.SpellConfig{
		ActionID:       debuffActionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell,
		ClassSpellMask: MonkSpellRushingJadeWind,
		MaxRange:       8,

		DamageMultiplier: 1.4, // 1.59 * (1.4 / 1.59),
		ThreatMultiplier: 1,
		CritMultiplier:   monk.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			results := make([]*core.SpellResult, numTargets)

			for idx := int32(0); idx < numTargets; idx++ {
				currentTarget := sim.Environment.GetTargetUnit(idx)
				baseDamage := monk.CalculateMonkStrikeDamage(sim, spell)
				results[idx] = spell.CalcDamage(sim, currentTarget, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			}

			for idx := int32(0); idx < numTargets; idx++ {
				spell.DealDamage(sim, results[idx])
			}
		},
	})

	baseCooldown := time.Second * 6

	rushingJadeWindBuff := monk.RegisterAura(core.Aura{
		Label:    "Rushing Jade Wind" + monk.Label,
		ActionID: actionID,
		Duration: baseCooldown,
	})

	isWiseSerpent := monk.StanceMatches(WiseSerpent)
	var rushingJadeWindSpell *core.Spell
	rushingJadeWindSpell = monk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          SpellFlagBuilder | core.SpellFlagAPL,
		ClassSpellMask: MonkSpellRushingJadeWind,

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
			IsAOE: true,
			Aura: core.Aura{
				Label:    "Rushing Jade Wind" + monk.Label,
				ActionID: debuffActionID,
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					rushingJadeWindSpell.CD.Duration = monk.ApplyCastSpeed(baseCooldown)
					rushingJadeWindBuff.Duration = rushingJadeWindSpell.CD.Duration
				},
			},
			NumberOfTicks:        8,
			TickLength:           time.Millisecond * 750,
			AffectedByCastSpeed:  true,
			HasteReducesDuration: true,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				rushingJadeWindTickSpell.Cast(sim, target)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dot := spell.AOEDot()
			dot.Apply(sim)
			dot.TickOnce(sim)

			remainingDuration := dot.RemainingDuration(sim)
			rushingJadeWindSpell.CD.Duration = remainingDuration
			rushingJadeWindBuff.Duration = remainingDuration
			rushingJadeWindBuff.Activate(sim)

			if numTargets >= 3 {
				monk.AddChi(sim, spell, 1, chiMetrics)
			}
		},
	})

	monk.AddOnCastSpeedChanged(func(_ float64, _ float64) {
		rushingJadeWindSpell.CD.Duration = monk.ApplyCastSpeed(baseCooldown)
		rushingJadeWindBuff.Duration = rushingJadeWindSpell.CD.Duration
	})
}

func (monk *Monk) registerInvokeXuenTheWhiteTiger() {
}

func (monk *Monk) registerChiTorpedo() {
}
