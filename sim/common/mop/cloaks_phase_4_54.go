package mop

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/common/shared"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

type cloakConfig struct {
	label  string
	itemID int32
}

func init() {
	// Xing-Ho, Breath of Yu'lon
	// Your damaging spells have a chance to grant 100% critical strike chance for 4 sec.  (Approximately 0.62 procs per minute)
	core.NewItemEffect(102246, func(agent core.Agent, state proto.ItemLevelState) {
		character := agent.GetCharacter()
		label := "Xing-Ho, Breath of Yu'lon"

		aura := core.MakePermanent(character.RegisterAura(core.Aura{
			Label: label,
		}))

		shared.RegisterIgniteEffect(&character.Unit, shared.IgniteConfig{
			ActionID:           core.ActionID{SpellID: 146198},
			SpellSchool:        core.SpellSchoolNature,
			DisableCastMetrics: true,
			DotAuraLabel:       "Essence of Yu'lon",
			TickLength:         1 * time.Second,
			NumberOfTicks:      4,
			NumberAOEOfTargets: 5,
			ParentAura:         aura,

			ProcTrigger: core.ProcTrigger{
				Name:    label,
				Harmful: true,
				DPM: character.NewRPPMProcManager(102246, false, core.ProcMaskSpecial|core.ProcMaskProc, core.RPPMConfig{
					PPM: 2.61100006104,
				}.WithHasteMod().
					WithSpecMod(0.25, proto.Spec_SpecArcaneMage).
					WithSpecMod(0.20000000298, proto.Spec_SpecFireMage).
					WithSpecMod(0.20000000298, proto.Spec_SpecFrostMage).
					WithSpecMod(-0.75, proto.Spec_SpecProtectionPaladin).
					WithSpecMod(-0.75, proto.Spec_SpecProtectionWarrior).
					WithSpecMod(0.10000000149, proto.Spec_SpecBalanceDruid).
					WithSpecMod(-0.75, proto.Spec_SpecGuardianDruid).
					WithSpecMod(-0.75, proto.Spec_SpecBloodDeathKnight).
					WithSpecMod(0, proto.Spec_SpecShadowPriest).
					WithSpecMod(0.05000000075, proto.Spec_SpecElementalShaman).
					WithSpecMod(0.10000000149, proto.Spec_SpecAfflictionWarlock).
					WithSpecMod(0.25, proto.Spec_SpecDemonologyWarlock).
					WithSpecMod(0.15000000596, proto.Spec_SpecDestructionWarlock).
					WithSpecMod(-0.75, proto.Spec_SpecBrewmasterMonk),
				),
				Callback: core.CallbackOnSpellHitDealt,
			},

			DamageCalculator: func(spell *core.Spell, result *core.SpellResult) float64 {
				return spell.SpellPower() * 2
			},
		})

		eligibleSlots := character.ItemSwap.EligibleSlotsForItem(102246)
		character.ItemSwap.RegisterProcWithSlots(102246, aura, eligibleSlots)
	})

	newXuenCloakEffect := func(config cloakConfig) {
		core.NewItemEffect(config.itemID, func(agent core.Agent, state proto.ItemLevelState) {
			character := agent.GetCharacter()
			numHits := min(5, character.Env.GetNumTargets())

			flurrySpell := character.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 146194},
				SpellSchool: core.SpellSchoolPhysical,
				ProcMask:    core.ProcMaskEmpty,
				Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

				DamageMultiplier: 1,
				CritMultiplier:   character.DefaultCritMultiplier(),
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					baseDamage := max(spell.MeleeAttackPower(), spell.RangedAttackPower()) * 0.2

					core.StartPeriodicAction(sim, core.PeriodicActionOptions{
						NumTicks: 10,
						Period:   300 * time.Millisecond,
						OnAction: func(sim *core.Simulation) {
							targetsHit := int32(0)
							for range numHits {
								if target == nil || targetsHit >= numHits {
									break
								}
								target = sim.Environment.NextTargetUnit(target)
								spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeAlwaysHit)
								targetsHit++
							}
						},
					})
				},
			})

			procTrigger := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:     config.label,
				ActionID: core.ActionID{SpellID: 146195},
				Harmful:  true,
				ICD:      time.Second * 3,
				DPM: character.NewRPPMProcManager(config.itemID,
					false,
					core.ProcMaskMeleeOrMeleeProc|core.ProcMaskRangedOrRangedProc,
					core.RPPMConfig{
						PPM: 1.74000000954,
					}.WithHasteMod().
						WithSpecMod(-0.40000000596, proto.Spec_SpecProtectionPaladin).
						WithSpecMod(0.44999998808, proto.Spec_SpecRetributionPaladin).
						WithSpecMod(0.34999999404, proto.Spec_SpecArmsWarrior).
						WithSpecMod(0.05000000075, proto.Spec_SpecFuryWarrior).
						WithSpecMod(-0.40000000596, proto.Spec_SpecProtectionWarrior).
						WithSpecMod(0.30000001192, proto.Spec_SpecFeralDruid).
						WithSpecMod(-0.40000000596, proto.Spec_SpecGuardianDruid).
						WithSpecMod(-0.40000000596, proto.Spec_SpecBloodDeathKnight).
						WithSpecMod(0.5, proto.Spec_SpecFrostDeathKnight).
						WithSpecMod(0.05000000075, proto.Spec_SpecUnholyDeathKnight).
						WithSpecMod(0, proto.Spec_SpecBeastMasteryHunter).
						WithSpecMod(0.20000000298, proto.Spec_SpecMarksmanshipHunter).
						WithSpecMod(0.15000000596, proto.Spec_SpecSurvivalHunter).
						WithSpecMod(0.55000001192, proto.Spec_SpecAssassinationRogue).
						WithSpecMod(0.15000000596, proto.Spec_SpecCombatRogue).
						WithSpecMod(0, proto.Spec_SpecSubtletyRogue).
						WithSpecMod(0.55000001192, proto.Spec_SpecEnhancementShaman).
						WithSpecMod(-0.40000000596, proto.Spec_SpecBrewmasterMonk).
						WithSpecMod(0.20000000298, proto.Spec_SpecWindwalkerMonk),
				),
				Callback: core.CallbackOnSpellHitDealt,
				Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					flurrySpell.Cast(sim, spell.Unit.CurrentTarget)
				},
			})

			eligibleSlots := character.ItemSwap.EligibleSlotsForItem(config.itemID)
			character.ItemSwap.RegisterProcWithSlots(config.itemID, procTrigger, eligibleSlots)
		})
	}

	// Fen-Yu, Fury of Xuen
	// Your damaging attacks have a chance to trigger a Flurry of Xuen, causing you to deal 1 damage
	// to up to 5 enemies in front of you, every 0.3 sec for 3 sec. (Approximately [1.74 + Haste] procs per minute)
	newXuenCloakEffect(cloakConfig{
		label:  "Fen-Yu, Fury of Xuen",
		itemID: 102248,
	})

	// Gong-Lu, Strength of Xuen
	// Your damaging attacks have a chance to trigger a Flurry of Xuen, causing you to deal 1 damage
	// to up to 5 enemies in front of you, every 0.3 sec for 3 sec. (Approximately [1.74 + Haste] procs per minute)
	newXuenCloakEffect(cloakConfig{
		label:  "Gong-Lu, Strength of Xuen",
		itemID: 102249,
	})

	newNiuzaoCloakEffect := func(config cloakConfig) {
		core.NewItemEffect(config.itemID, func(agent core.Agent, state proto.ItemLevelState) {
			character := agent.GetCharacter()

			dummyAura := core.MakePermanent(character.RegisterAura(core.Aura{
				Label:    "Endurance of Niuzao - Dummy",
				Duration: core.NeverExpires,
			}))

			shieldAura := character.RegisterAura(core.Aura{
				Label:     "Endurance of Niuzao",
				ActionID:  core.ActionID{SpellID: 146193},
				Duration:  core.NeverExpires,
				MaxStacks: math.MaxInt32,
			})

			icd := core.Cooldown{
				Timer:    character.NewTimer(),
				Duration: time.Minute * 2,
			}

			character.AddDynamicDamageTakenModifier(func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult, isPeriodic bool) {
				if character.CurrentHealth()-result.Damage <= 0 && dummyAura.IsActive() && icd.IsReady(sim) {
					shieldAura.Activate(sim)
					absorbedDamage := result.Damage
					result.Damage = 0
					shieldAura.SetStacks(sim, int32(absorbedDamage))
					shieldAura.Deactivate(sim)
					icd.Use(sim)
				}
			})

			eligibleSlots := character.ItemSwap.EligibleSlotsForItem(config.itemID)
			character.ItemSwap.RegisterProcWithSlots(config.itemID, dummyAura, eligibleSlots)
		})
	}

	// Qian-Le, Courage of Niuzao
	// The Endurance of Niuzao fully absorbs the damage of one attack that would normally kill you. This effect has a 2 min cooldown. Does not function for non-Tank-specialized characters.
	newNiuzaoCloakEffect(cloakConfig{
		label:  "Qian-Le, Courage of Niuzao",
		itemID: 102245,
	})

	// Qian-Ying, Fortitude of Niuzao
	// The Endurance of Niuzao fully absorbs the damage of one attack that would normally kill you. This effect has a 2 min cooldown. Does not function for non-Tank-specialized characters.
	newNiuzaoCloakEffect(cloakConfig{
		label:  "Qian-Ying, Fortitude of Niuzao",
		itemID: 102250,
	})
}
