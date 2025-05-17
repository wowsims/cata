package bwd

import (
	"fmt"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/encounters/default_ai"
)

func createMagmawPreset(bossPrefix string, raidSize int, isHeroic bool,
	npcId int32, health float64, minBaseDamage float64,
	addNpcId int32, addHealth float64, addMinBaseDamage float64) {

	targetName := fmt.Sprintf("Magmaw %d", raidSize)
	targetNameAdd := fmt.Sprintf("Blazing Construct %d", raidSize)
	if isHeroic {
		targetName = targetName + " H"
		targetNameAdd = targetNameAdd + " H"
	}
	core.AddPresetTarget(&core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: &proto.Target{
			Id:        npcId,
			Name:      targetName,
			Level:     88,
			MobType:   proto.MobType_MobTypeBeast,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      health,
				stats.Armor:       11977,
				stats.AttackPower: 0,
			}.ToProtoArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       2.5,
			MinBaseDamage:    minBaseDamage,
			DamageSpread:     0.4,
			SuppressDodge:    false,
			ParryHaste:       false,
			DualWield:        false,
			DualWieldPenalty: false,
			TargetInputs: []*proto.TargetInput{
				// TODO: Figure out how to make Size and Heroic
				// pickable for a preset. Right now we have to
				// make up to 4 presets per boss which sucks...
				// {
				// 	Label:     "Raid Size",
				// 	Tooltip:   "The size of the Raid",
				// 	InputType: proto.InputType_Enum,
				// 	EnumValue: 1,
				// 	EnumOptions: []string{
				// 		"10", "25",
				// 	},
				// },
				// {
				// 	Label:     "Heroic",
				// 	Tooltip:   "Is the encounter in Heroic Mode",
				// 	InputType: proto.InputType_Bool,
				// 	BoolValue: true,
				// },
				{
					Label:       "Impale Reaction Time",
					Tooltip:     "How long will the Raid take to Impale Head in Seconds. (After the initial 10s)",
					InputType:   proto.InputType_Number,
					NumberValue: 5.0,
				},
			},
		},
		AI: func() core.TargetAI {
			return makeMagmawAI(raidSize, isHeroic)
		},
	})

	if isHeroic {
		core.AddPresetTarget(&core.PresetTarget{
			PathPrefix: bossPrefix,
			Config: &proto.Target{
				Id:        addNpcId,
				Name:      targetNameAdd,
				Level:     87,
				MobType:   proto.MobType_MobTypeBeast,
				TankIndex: 1,

				Stats: stats.Stats{
					stats.Health:      addHealth,
					stats.Armor:       11977,
					stats.AttackPower: 0,
				}.ToProtoArray(),

				SpellSchool:   proto.SpellSchool_SpellSchoolPhysical,
				SwingSpeed:    2.0,
				MinBaseDamage: addMinBaseDamage,
				DamageSpread:  0.5,
				TargetInputs:  []*proto.TargetInput{},
			},
			AI: default_ai.NewDefaultAI([]default_ai.TargetAbility{
				{
					InitialCD:   time.Second * 5,
					ChanceToUse: 0,
					MakeSpell: func(target *core.Target) *core.Spell {
						// Fiery Slash Next melee Spell
						nextMeleeSpell := target.GetOrRegisterSpell(core.SpellConfig{
							ActionID:    core.ActionID{SpellID: 92144},
							SpellSchool: core.SpellSchoolFire,
							ProcMask:    core.ProcMaskSpellDamage,

							Cast: core.CastConfig{
								CD: core.Cooldown{
									Timer:    target.NewTimer(),
									Duration: time.Second * 7,
								},
							},

							DamageMultiplier: 0.75,

							ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
								spell.CalcAndDealDamage(sim, target, spell.Unit.AutoAttacks.MH().EnemyWeaponDamage(sim, spell.MeleeAttackPower(), 0.5), spell.OutcomeEnemyMeleeWhite)
							},
						})

						target.AutoAttacks.SetReplaceMHSwing(func(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
							if nextMeleeSpell.CanCast(sim, target.CurrentTarget) && sim.Proc(0.75, "Fiery Slash Cast") {
								return nextMeleeSpell
							}
							return mhSwingSpell
						})

						target.AutoAttacks.MHConfig().ActionID.Tag = 49416

						return nextMeleeSpell
					},
				},
			}),
		})
		core.AddPresetEncounter(targetName, []string{
			bossPrefix + "/" + targetName,
			bossPrefix + "/" + targetNameAdd,
		})
	} else {
		core.AddPresetEncounter(targetName, []string{
			bossPrefix + "/" + targetName,
		})
	}

}

func addMagmaw(bossPrefix string) {
	// size, heroic, boss hp, boss min damage, add hp, add min damage
	createMagmawPreset(bossPrefix, 10, false, 41570, 26_798_304, 110000, 0, 0, 0)
	createMagmawPreset(bossPrefix, 25, false, 41571, 81_082_048, 150000, 0, 0, 0)
	createMagmawPreset(bossPrefix, 10, true, 41572, 39_200_000, 150000, 49416, 1_410_000, 44000)
	createMagmawPreset(bossPrefix, 25, true, 41573, 120_016_403, 210000, 49417, 4_500_000, 80000)
}

func makeMagmawAI(raidSize int, isHeroic bool) core.TargetAI {
	return &MagmawAI{
		raidSize: raidSize,
		isHeroic: isHeroic,
	}
}

type MagmawAI struct {
	Target *core.Target

	canAct             bool
	individualTankSwap bool

	lastMangleTarget *core.Unit

	raidSize    int
	isHeroic    bool
	impaleDelay float64

	mangle    *core.Spell
	magmaSpit *core.Spell
	lavaSpew  *core.Spell

	pointOfVulnerability *core.Aura
	swelteringArmor      core.AuraArray
}

func (ai *MagmawAI) Initialize(target *core.Target, config *proto.Target) {
	ai.Target = target

	// if target.Env.Raid.Size() <= 1 {
	// 	// Individual Sims - use the input configuration
	// 	ai.raidSize = []int{10, 25}[config.TargetInputs[0].EnumValue]
	// } else {
	// 	// Raid sim - Set from number of players
	// 	ai.raidSize = 10
	// 	if target.Env.Raid.Size() > 10 {
	// 		ai.raidSize = 25
	// 	}
	// }

	// ai.isHeroic = config.TargetInputs[1].BoolValue
	// ai.impaleDelay = config.TargetInputs[2].NumberValue

	ai.Target.AutoAttacks.MHConfig().ActionID.Tag = 41570

	ai.impaleDelay = config.TargetInputs[0].NumberValue
	ai.registerSpells()
}

func (ai *MagmawAI) Reset(sim *core.Simulation) {
	ai.canAct = true
	ai.individualTankSwap = false
}

const BossGCD = time.Millisecond * 1620

func (ai *MagmawAI) ExecuteCustomRotation(sim *core.Simulation) {
	if !ai.canAct {
		ai.Target.WaitUntil(sim, sim.CurrentTime+BossGCD)
		return
	}

	target := ai.Target.CurrentTarget
	if target == nil {
		// For individual non tank sims we still want abilities to work
		target = &ai.Target.Env.Raid.Parties[0].Players[0].GetCharacter().Unit
		ai.individualTankSwap = true
	}

	// Mangle
	if ai.mangle.CanCast(sim, target) {
		ai.mangle.Cast(sim, target)
		return
	}

	// Lava Spew
	if ai.lavaSpew.CanCast(sim, target) && sim.Proc(0.7, "Lava Spew Cast Roll") {
		ai.lavaSpew.Cast(sim, target)
		return
	}

	// Magma Spit
	if ai.magmaSpit.CanCast(sim, target) && sim.Proc(0.6, "Magma Spit Cast Roll") {
		ai.magmaSpit.Cast(sim, target)
		ai.Target.ExtendGCDUntil(sim, sim.CurrentTime+BossGCD)
		return
	}

	ai.Target.WaitUntil(sim, sim.CurrentTime+BossGCD)
}

func (ai *MagmawAI) registerSpells() {
	// 0 - 10N, 1 - 25N, 2 - 10H, 3 - 25H
	scalingIndex := core.TernaryInt(ai.raidSize == 10, core.TernaryInt(ai.isHeroic, 2, 0), core.TernaryInt(ai.isHeroic, 3, 1))
	isIndividualSim := ai.Target.Env.Raid.Size() == 1
	tankUnit := &ai.Target.Env.Raid.Parties[0].Players[0].GetCharacter().Unit

	// Exposed Aura
	ai.pointOfVulnerability = ai.Target.GetOrRegisterAura(core.Aura{
		Label:    "Point of Vulnerability",
		ActionID: core.ActionID{SpellID: 79010},
		Duration: time.Second * 30,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageTakenMultiplier *= 2
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageTakenMultiplier /= 2

			if sim.CurrentTime >= sim.Duration {
				return
			}

			// TODO: Move this to an APL Action
			if !isIndividualSim && ai.Target.Env.GetNumTargets() > 1 {
				addTarget := ai.Target.Env.NextTargetUnit(&ai.Target.Unit)
				if addTarget.CurrentTarget != nil {
					// Swap Tanks
					addTank := addTarget.CurrentTarget
					bossTank := ai.Target.CurrentTarget

					addTank.CurrentTarget = &ai.Target.Unit
					bossTank.CurrentTarget = addTarget

					addTarget.CurrentTarget = bossTank
					ai.Target.CurrentTarget = addTank
				}
			} else if isIndividualSim {
				// Individual sim fake tank swaps
				if tankUnit.Metrics.IsTanking() {
					if !ai.individualTankSwap {
						// Remove boss target
						ai.individualTankSwap = true
						ai.Target.CurrentTarget = nil

						// Set add target
						if ai.Target.Env.GetNumTargets() > 1 {
							addTarget := ai.Target.Env.NextTargetUnit(&ai.Target.Unit)
							tankUnit.CurrentTarget = addTarget

							addTarget.CurrentTarget = tankUnit
							addTarget.AutoAttacks.EnableAutoSwing(sim)
						}
					} else {
						ai.individualTankSwap = false

						// Set boss target
						ai.Target.CurrentTarget = tankUnit
						tankUnit.CurrentTarget = &ai.Target.Unit

						// Remove add target
						if ai.Target.Env.GetNumTargets() > 1 {
							addTarget := ai.Target.Env.NextTargetUnit(&ai.Target.Unit)
							addTarget.AutoAttacks.CancelAutoSwing(sim)
							addTarget.CurrentTarget = nil
						}
					}
				}
			}

			ai.canAct = true
			ai.Target.AutoAttacks.EnableAutoSwing(sim)
		},
	})

	// Mangle Debuff Aura
	ai.swelteringArmor = ai.Target.NewAllyAuraArray(func(unit *core.Unit) *core.Aura {
		if unit.Type == core.PetUnit {
			return nil
		}
		return unit.GetOrRegisterAura(core.Aura{
			Label:    "Sweltering Armor",
			ActionID: core.ActionID{SpellID: 78199},
			Duration: time.Second * 90,

			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.PseudoStats.ArmorMultiplier *= 0.5
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.PseudoStats.ArmorMultiplier /= 0.5
			},
		})
	})

	lavaSpewBase := []float64{
		14799,
		14799,
		20811,
		24049,
	}[scalingIndex]

	lavaSpewVariance := []float64{
		2401,
		2401,
		3376,
		3901,
	}[scalingIndex]

	lavaSpewDamageRoll := func(sim *core.Simulation) float64 {
		return lavaSpewBase + lavaSpewVariance*sim.RandomFloat("Lava Spew Damage")
	}

	ai.lavaSpew = ai.Target.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 77690},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,

		DamageMultiplier: 1,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    ai.Target.NewTimer(),
				Duration: time.Second * 30,
			},
			IgnoreHaste: true,
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 2,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				spell.Unit.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime+cast.CastTime, false)
			},
		},

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label:    "Lava Spew",
				ActionID: core.ActionID{SpellID: 77690},
			},

			TickLength:    time.Second * 2,
			NumberOfTicks: 3,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				for _, aoeTarget := range sim.Raid.AllPlayerUnits {
					dot.Spell.CalcAndDealDamage(sim, aoeTarget, lavaSpewDamageRoll(sim), dot.Spell.OutcomeAlwaysHit)
				}

				// This tick delays melees by up to 300ms after it lands
				meleeMinAt := sim.CurrentTime + time.Millisecond*300
				nextMeleeAt := dot.Spell.Unit.AutoAttacks.NextAttackAt()
				if nextMeleeAt < meleeMinAt {
					dot.Spell.Unit.AutoAttacks.DelayMeleeBy(sim, meleeMinAt-nextMeleeAt)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.AOEDot().Apply(sim)
		},
	})

	magmaSpitBase := []float64{
		30624,
		30624,
		34999,
		39374,
	}[scalingIndex]

	magmaSpitVariance := []float64{
		8751,
		8751,
		10001,
		11251,
	}[scalingIndex]

	magmaSpitDamageRoll := func(sim *core.Simulation) float64 {
		return magmaSpitBase + magmaSpitVariance*sim.RandomFloat("Magma Spit Damage")
	}

	removeAt := func(s []int32, i int32) []int32 {
		s[i] = s[len(s)-1]
		return s[:len(s)-1]
	}
	numTargets := core.TernaryInt32(ai.raidSize == 10, 3, 8)
	ai.magmaSpit = ai.Target.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 78359},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,

		DamageMultiplier: 1,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    ai.Target.NewTimer(),
				Duration: time.Second * 7,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if isIndividualSim {
				chanceToBeHit := float64(numTargets) / float64(ai.raidSize)
				if sim.Proc(float64(chanceToBeHit), "Magma Spit Hit") {
					spell.CalcAndDealDamage(sim, target, magmaSpitDamageRoll(sim), spell.OutcomeAlwaysHit)
				}
			} else {
				if int(numTargets) >= len(sim.Raid.AllPlayerUnits) {
					for _, aoeTarget := range sim.Raid.AllPlayerUnits {
						spell.CalcAndDealDamage(sim, aoeTarget, magmaSpitDamageRoll(sim), spell.OutcomeAlwaysHit)
					}
				} else {
					validTargets := make([]int32, 0)
					for idx := range sim.Raid.AllPlayerUnits {
						validTargets = append(validTargets, int32(idx))
					}
					hitTargets := make([]int32, 0)
					for idx := int32(0); idx < numTargets; idx++ {
						targetRoll := int32(sim.RandomFloat("Magma Spit Target Roll") * float64(len(validTargets)))
						hitTargets = append(hitTargets, validTargets[targetRoll])
						validTargets = removeAt(validTargets, targetRoll)
					}

					for idx := int32(0); idx < numTargets; idx++ {
						spell.CalcAndDealDamage(sim, sim.Raid.AllPlayerUnits[hitTargets[idx]], magmaSpitDamageRoll(sim), spell.OutcomeAlwaysHit)
					}
				}
			}
		},
	})

	mangleTick := []float64{
		110463,
		132556,
		132556,
		154648,
	}[scalingIndex]

	// Variance listed in DB2 but not observed in logs
	// mangleVariance := []float64{
	// 	17914,
	// 	21496,
	// 	21496,
	// 	25080,
	// }[scalingIndex]

	ai.mangle = ai.Target.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 89773},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagApplyArmorReduction,

		DamageMultiplier: 1,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    ai.Target.NewTimer(),
				Duration: time.Second * 90,
			},
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:    "Magmaw Mangle",
				ActionID: core.ActionID{SpellID: 89773},
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					ai.Target.AutoAttacks.CancelAutoSwing(sim)
					ai.canAct = false
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					if sim.CurrentTime >= sim.Duration {
						return
					}

					// Activate Expose
					ai.pointOfVulnerability.Activate(sim)

					if !isIndividualSim || (!ai.individualTankSwap && tankUnit.Metrics.IsTanking()) {
						ai.swelteringArmor.Get(ai.lastMangleTarget).Activate(sim)
					}
				},
			},

			TickLength:    time.Second * 2,
			NumberOfTicks: 5 + int32(ai.impaleDelay/2.0), // Simulate Mangle Duration as 10s + Input Delay

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				doDamage := !isIndividualSim || ai.Target.Env.Raid.Parties[0].Players[0].GetCharacter().Unit.Metrics.IsTanking()
				if doDamage {
					if isIndividualSim && ai.individualTankSwap {
						return
					}
					baseDamage := mangleTick // + mangleVariance*sim.RandomFloat("Magmaw Mangle Tick")
					dot.Spell.CalcAndDealPeriodicDamage(sim, target, baseDamage, dot.Spell.OutcomeAlwaysHit)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			ai.lastMangleTarget = target
			doDamage := !isIndividualSim || ai.Target.Env.Raid.Parties[0].Players[0].GetCharacter().Unit.Metrics.IsTanking()
			if doDamage && (!isIndividualSim || !ai.individualTankSwap) {
				baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) * 1.5
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeAlwaysHit)
			}
			spell.Dot(target).Apply(sim)
		},
	})

	ai.Target.RegisterResetEffect(func(sim *core.Simulation) {
		ai.mangle.CD.Use(sim)
		ai.magmaSpit.CD.Set(time.Second * 5)
		ai.lavaSpew.CD.Set(time.Second * 20)
	})
}
