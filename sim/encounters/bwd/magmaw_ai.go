package bwd

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

func addMagmaw(bossPrefix string) {
	core.AddPresetTarget(&core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: &proto.Target{
			Id:        41570,
			Name:      "Magmaw",
			Level:     88,
			MobType:   proto.MobType_MobTypeBeast,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      120_000_000,
				stats.Armor:       11977,
				stats.AttackPower: 650,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       2.5,
			MinBaseDamage:    209609,
			DamageSpread:     0.4,
			SuppressDodge:    false,
			ParryHaste:       false,
			DualWield:        false,
			DualWieldPenalty: false,
			TargetInputs: []*proto.TargetInput{
				{
					Label:     "Raid Size",
					Tooltip:   "The size of the Raid",
					InputType: proto.InputType_Enum,
					EnumValue: 1,
					EnumOptions: []string{
						"10", "25",
					},
				},
				{
					Label:     "Heroic",
					Tooltip:   "Is the encounter in Heroic Mode",
					InputType: proto.InputType_Bool,
					BoolValue: true,
				},
				{
					Label:       "Impale Reaction Time",
					Tooltip:     "How long will the Raid take to Impale Head in Seconds. (After the initial 10s)",
					InputType:   proto.InputType_Number,
					NumberValue: 5.0,
				},
			},
		},
		AI: func() core.TargetAI {
			return &MagmawAI{}
		},
	})
	core.AddPresetEncounter("Magmaw", []string{
		bossPrefix + "/Magmaw",
	})
}

type MagmawAI struct {
	Target *core.Target

	canAct bool

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

	if target.Env.Raid.Size() <= 1 {
		// Individual Sims - use the input configuration
		ai.raidSize = []int{10, 25}[config.TargetInputs[0].EnumValue]
	} else {
		// Raid sim - Set from number of players
		ai.raidSize = 10
		if target.Env.Raid.Size() > 10 {
			ai.raidSize = 25
		}
	}

	ai.isHeroic = config.TargetInputs[1].BoolValue
	ai.impaleDelay = config.TargetInputs[2].NumberValue

	ai.registerSpells()
}

func (ai *MagmawAI) Reset(sim *core.Simulation) {
	ai.canAct = true
}

const BossGCD = time.Millisecond * 1620

func (ai *MagmawAI) ExecuteCustomRotation(sim *core.Simulation) {
	if !ai.canAct {
		ai.Target.WaitUntil(sim, sim.CurrentTime+BossGCD)
		return
	}

	// Mangle
	if ai.mangle.CanCast(sim, ai.Target.CurrentTarget) {
		ai.mangle.Cast(sim, ai.Target.CurrentTarget)
		return
	}

	// Lava Spew
	if ai.lavaSpew.CanCast(sim, ai.Target.CurrentTarget) && sim.Proc(0.7, "Lava Spew Cast Roll") {
		ai.lavaSpew.Cast(sim, ai.Target.CurrentTarget)
		return
	}

	// Magma Spit
	if ai.magmaSpit.CanCast(sim, ai.Target.CurrentTarget) && sim.Proc(0.6, "Magma Spit Cast Roll") {
		ai.magmaSpit.Cast(sim, ai.Target.CurrentTarget)
		return
	}

	ai.Target.WaitUntil(sim, sim.CurrentTime+BossGCD)
}

func (ai *MagmawAI) registerSpells() {
	// 0 - 10N, 1 - 25N, 2 - 10H, 3 - 25H
	scalingIndex := core.TernaryInt(ai.raidSize == 10, core.TernaryInt(ai.isHeroic, 2, 0), core.TernaryInt(ai.isHeroic, 3, 1))
	isIndividualSim := ai.Target.Env.Raid.Size() == 1

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
		110463,
		132556,
		132556,
		154648,
	}[scalingIndex]

	magmaSpitVariance := []float64{
		17914,
		21496,
		21496,
		25080,
	}[scalingIndex]

	magmaSpitDamageRoll := func(sim *core.Simulation) float64 {
		return magmaSpitBase + magmaSpitVariance*sim.RandomFloat("Magma Spit Damage")
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
				if numTargets >= sim.Environment.GetNumTargets() {
					for _, aoeTarget := range sim.Raid.AllPlayerUnits {
						spell.CalcAndDealDamage(sim, aoeTarget, magmaSpitDamageRoll(sim), spell.OutcomeAlwaysHit)
					}
				} else {
					validTargets := make([]int32, 0)
					for idx, _ := range sim.Raid.AllPlayerUnits {
						validTargets = append(validTargets, int32(idx))
					}
					hitTargets := make([]int32, 0)
					for idx := int32(0); idx < numTargets; idx++ {
						targetRoll := int(sim.RandomFloat("Magma Spit Target Roll") * float64(len(validTargets)))
						rolledTarget := validTargets[targetRoll]
						hitTargets = append(hitTargets, int32(rolledTarget))

						remove := func(s []int32, i int32) []int32 {
							s[i] = s[len(s)-1]
							return s[:len(s)-1]
						}
						validTargets = remove(validTargets, rolledTarget)
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
					// Activate Expose
					ai.pointOfVulnerability.Activate(sim)
				},
			},

			TickLength:    time.Second * 2,
			NumberOfTicks: 5 + int32(ai.impaleDelay/2.0), // Simulate Mangle Duration as 10s + Input Delay

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				baseDamage := mangleTick // + mangleVariance*sim.RandomFloat("Magmaw Mangle Tick")
				dot.Spell.CalcAndDealPeriodicDamage(sim, target, baseDamage, dot.Spell.OutcomeAlwaysHit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) * 1.5
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeAlwaysHit)
			spell.Dot(target).Apply(sim)
		},
	})

	ai.Target.RegisterResetEffect(func(sim *core.Simulation) {
		ai.mangle.CD.Use(sim)
		ai.magmaSpit.CD.Set(time.Second * 5)
		ai.lavaSpew.CD.Set(time.Second * 20)
	})
}
