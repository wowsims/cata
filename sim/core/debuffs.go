package core

import (
	"strconv"
	"time"

	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

func applyDebuffEffects(target *Unit, targetIdx int, debuffs *proto.Debuffs, raid *proto.Raid) {
	if debuffs.Judgement && targetIdx == 0 {
		MakePermanent(JudgementOfLightAura(target))
	}

	// +8% Spell DMG
	if debuffs.CurseOfElements && targetIdx == 0 {
		MakePermanent(CurseOfElementsAura(target))
	}

	if debuffs.EbonPlaguebringer {
		MakePermanent(EbonPlaguebringerAura(nil, target, 2, 3))
	}

	if debuffs.EarthAndMoon && targetIdx == 0 {
		MakePermanent(EarthAndMoonAura(target))
	}

	if debuffs.MasterPoisoner && targetIdx == 0 {
		MakePermanent(MasterPoisonerDebuff(target))
	}

	if debuffs.FireBreath && targetIdx == 0 {
		MakePermanent(FireBreathDebuff(target))
	}

	if debuffs.LightningBreath && targetIdx == 0 {
		MakePermanent(LightningBreath(target))
	}

	// +4% Phsyical Damage
	if debuffs.BloodFrenzy && targetIdx < 4 {
		MakePermanent(BloodFrenzyAura(target, 2))
		MakePermanent(TraumaAura(target, 2))
	}

	if debuffs.SavageCombat {
		MakePermanent(SavageCombatAura(target, 2))
	}

	if debuffs.FrostFever || debuffs.BrittleBones {
		MakePermanent(FrostFeverAura(target, TernaryInt32(debuffs.BrittleBones, 2, 0)))
	}

	if debuffs.AcidSpit && targetIdx == 0 {
		MakePermanent(AcidSpitAura(target))
	}

	// Bleed Damage
	// Blood Frenzy @4% Physical Damage
	if debuffs.Mangle && targetIdx == 0 {
		MakePermanent(MangleAura(target))
	}

	if debuffs.Hemorrhage && targetIdx == 0 {
		MakePermanent(HemorrhageAura(target))
	}

	if debuffs.Stampede && targetIdx == 0 {
		MakePermanent(StampedeAura(target))
	}

	// Spell Crit
	if debuffs.CriticalMass && targetIdx == 0 {
		MakePermanent(CriticalMassAura(target))
	}

	if debuffs.ShadowAndFlame && targetIdx == 0 {
		MakePermanent(CriticalMassAura(target))
	}

	if debuffs.ExposeArmor && targetIdx == 0 {
		aura := ExposeArmorAura(target, false)
		ScheduledMajorArmorAura(aura, PeriodicActionOptions{
			Period:   time.Second * 3,
			NumTicks: 1,
			OnAction: func(sim *Simulation) {
				aura.Activate(sim)
			},
		}, raid)
	}

	if debuffs.SunderArmor && targetIdx == 0 {
		aura := SunderArmorAura(target)
		ScheduledMajorArmorAura(aura, PeriodicActionOptions{
			Period:          time.Millisecond * 1500,
			NumTicks:        3,
			TickImmediately: true,
			Priority:        ActionPriorityDOT, // High prio so it comes before actual warrior sunders.
			OnAction: func(sim *Simulation) {
				aura.Activate(sim)
				if aura.IsActive() {
					aura.AddStack(sim)
				}
			},
		}, raid)
	}

	if debuffs.CorrosiveSpit && targetIdx == 0 {
		aura := CorrosiveSpitAura(target)
		ScheduledMajorArmorAura(aura, PeriodicActionOptions{
			Period:          time.Second * 10,
			NumTicks:        3,
			TickImmediately: true,
			Priority:        ActionPriorityDOT,
			OnAction: func(sim *Simulation) {
				aura.Activate(sim)
				if aura.IsActive() {
					aura.AddStack(sim)
				}
			},
		}, raid)
	}

	if debuffs.FaerieFire && targetIdx == 0 {
		aura := FaerieFireAura(target)
		ScheduledMajorArmorAura(aura, PeriodicActionOptions{
			Period:          time.Millisecond * 1500,
			NumTicks:        3,
			TickImmediately: true,
			Priority:        ActionPriorityDOT,
			OnAction: func(sim *Simulation) {
				aura.Activate(sim)
				if aura.IsActive() {
					aura.AddStack(sim)
				}
			},
		}, raid)
	}

	// -10% Physical Damage
	if debuffs.CurseOfWeakness {
		MakePermanent(CurseOfWeaknessAura(target))
	}

	if debuffs.DemoralizingRoar {
		MakePermanent(DemoralizingRoarAura(target))
	}

	if debuffs.DemoralizingShout {
		MakePermanent(DemoralizingShoutAura(target, false))
	}

	if debuffs.DemoralizingScreech {
		MakePermanent(DemoralizingScreechAura(target))
	}

	if debuffs.Vindication {
		MakePermanent(VindicationAura(target))
	}

	if debuffs.ScarletFever {
		MakePermanent(ScarletFeverAura(target, 2, 0))
	}

	// Atk spd reduction
	if debuffs.ThunderClap {
		MakePermanent(ThunderClapAura(target))
	}

	if debuffs.InfectedWounds && targetIdx == 0 {
		MakePermanent(InfectedWoundsAura(target, 2))
	}

	if debuffs.JudgementsOfTheJust && targetIdx == 0 {
		MakePermanent(JudgementsOfTheJustAura(target, 2))
	}

	if debuffs.DustCloud && targetIdx == 0 {
		MakePermanent(DustCloud(target))
	}
}

func ScheduledMajorArmorAura(aura *Aura, options PeriodicActionOptions, raid *proto.Raid) {
	aura.OnReset = func(aura *Aura, sim *Simulation) {
		aura.Duration = NeverExpires
		StartPeriodicAction(sim, options)
	}
}

var JudgementOfLightAuraLabel = "Judgement of Light"

func JudgementOfLightAura(target *Unit) *Aura {
	actionID := ActionID{SpellID: 20271}

	return target.GetOrRegisterAura(Aura{
		Label:    JudgementOfLightAuraLabel,
		ActionID: actionID,
		Duration: time.Second * 20,
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if !spell.ProcMask.Matches(ProcMaskMelee) || !result.Landed() {
				return
			}
		},
	})
}

func CurseOfElementsAura(target *Unit) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Curse of Elements",
		ActionID: ActionID{SpellID: 1490},
		Duration: time.Minute * 5,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatsDynamic(sim, stats.Stats{stats.ArcaneResistance: -215, stats.FireResistance: -215, stats.FrostResistance: -215, stats.ShadowResistance: -215, stats.NatureResistance: -215})
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatsDynamic(sim, stats.Stats{stats.ArcaneResistance: 215, stats.FireResistance: 215, stats.FrostResistance: 215, stats.ShadowResistance: 215, stats.NatureResistance: 215})
		},
	})
	spellDamageEffect(aura, 1.08)
	return aura
}

func EarthAndMoonAura(target *Unit) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Earth And Moon",
		ActionID: ActionID{SpellID: 60433},
		Duration: time.Second * 12,
	})
	spellDamageEffect(aura, 1.08)
	return aura
}

func MasterPoisonerDebuff(target *Unit) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Master Poisoner",
		ActionID: ActionID{SpellID: 58410},
		Duration: time.Second * 15,
	})
	spellDamageEffect(aura, 1.08)
	return aura
}

func FireBreathDebuff(target *Unit) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Fire Breath",
		ActionID: ActionID{SpellID: 34889},
		Duration: time.Second * 45,
	})
	spellDamageEffect(aura, 1.08)
	return aura
}

func LightningBreath(target *Unit) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Lightning Breath",
		ActionID: ActionID{SpellID: 24844},
		Duration: time.Second * 45,
	})
	spellDamageEffect(aura, 1.08)
	return aura
}

func EbonPlaguebringerAura(caster *Character, target *Unit, epidemicPoints int32, ebonPlaguebringerPoints int32) *Aura {
	// On application, Ebon Plaguebringer trigger extra 'ghost' procs.
	var ghostSpell *Spell
	label := "External"
	if caster != nil {
		label = caster.Label
		ghostSpell = caster.RegisterSpell(SpellConfig{
			ActionID:    ActionID{SpellID: 52789},
			SpellSchool: SpellSchoolShadow,
			ProcMask:    ProcMaskSpellDamage,
			Flags:       SpellFlagNoLogs | SpellFlagNoMetrics | SpellFlagNoOnCastComplete | SpellFlagIgnoreModifiers,

			DamageMultiplier: 1,
			ThreatMultiplier: 0,

			ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
				// Just deal 0 damage as the "Harmful Spell" is implemented on spell damage
				spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
			},
		})
	}

	aura := target.GetOrRegisterAura(Aura{
		Label:    "EbonPlaguebringer" + label, // Support multiple DKs having their EP up
		Tag:      "EbonPlaguebringer",
		ActionID: ActionID{SpellID: 65142},
		Duration: time.Second * (21 + []time.Duration{0, 4, 8, 12}[epidemicPoints]),
		OnGain: func(aura *Aura, sim *Simulation) {
			if ghostSpell != nil {
				ghostSpell.Cast(sim, aura.Unit)
			}
		},
	})

	if ebonPlaguebringerPoints > 0 {
		spellDamageEffect(aura, 1.08)
	}
	return aura
}

func spellDamageEffect(aura *Aura, multiplier float64) *ExclusiveEffect {
	return aura.NewExclusiveEffect("SpellDamageTaken%", false, ExclusiveEffect{
		Priority: multiplier,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] *= multiplier
			ee.Aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] *= multiplier
			ee.Aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] *= multiplier
			ee.Aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] *= multiplier
			ee.Aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] *= multiplier
			ee.Aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexHoly] *= multiplier
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] /= multiplier
			ee.Aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] /= multiplier
			ee.Aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] /= multiplier
			ee.Aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] /= multiplier
			ee.Aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] /= multiplier
			ee.Aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexHoly] /= multiplier
		},
	})
}

func BloodFrenzyAura(target *Unit, points int32) *Aura {
	return bloodFrenzySavageCombatAura(target, "Blood Frenzy", ActionID{SpellID: 29859}, points)
}
func SavageCombatAura(target *Unit, points int32) *Aura {
	return bloodFrenzySavageCombatAura(target, "Savage Combat", ActionID{SpellID: 58413}, points)
}
func AcidSpitAura(target *Unit) *Aura {
	return bloodFrenzySavageCombatAura(target, "Acid Spit", ActionID{SpellID: 55754}, 2)
}

func bloodFrenzySavageCombatAura(target *Unit, label string, id ActionID, points int32) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    label + "-" + strconv.Itoa(int(points)),
		ActionID: id,
		// No fixed duration, lasts as long as the bleed that activates it.
		Duration: NeverExpires,
	})

	multiplier := 1 + 0.02*float64(points)
	PhysDamageTakenEffect(aura, multiplier)
	return aura
}

func MangleAura(target *Unit) *Aura {
	return bleedDamageAura(target, Aura{
		Label:    "Mangle",
		ActionID: ActionID{SpellID: 33876},
		Duration: time.Minute,
	}, 1.3)
}

func HemorrhageAura(target *Unit) *Aura {
	return bleedDamageAura(target, Aura{
		Label:    "Hemorrhage",
		ActionID: ActionID{SpellID: 16511},
		Duration: time.Minute,
	}, 1.3)
}

func StampedeAura(target *Unit) *Aura {
	return bleedDamageAura(target, Aura{
		Label:    "Stampede",
		ActionID: ActionID{SpellID: 57386},
		Duration: time.Second * 30,
	}, 1.3)
}

func TraumaAura(target *Unit, points int32) *Aura {
	return bleedDamageAura(target, Aura{
		Label:    "Trauma",
		ActionID: ActionID{SpellID: TernaryInt32(points == 1, 46856, 46857)},
		Duration: 1 * time.Minute,
	}, TernaryFloat64(points == 1, 1.15, 1.3))
}

// Bleed Damage Multiplier category
const BleedEffectCategory = "BleedDamage"

func bleedDamageEffect(aura *Aura, multiplier float64) *Aura {
	aura.NewExclusiveEffect(BleedEffectCategory, true, ExclusiveEffect{
		Priority: multiplier,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.PeriodicPhysicalDamageTakenMultiplier *= multiplier
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.PeriodicPhysicalDamageTakenMultiplier /= multiplier
		},
	})

	return aura
}
func bleedDamageAura(target *Unit, config Aura, multiplier float64) *Aura {
	aura := target.GetOrRegisterAura(config)
	return bleedDamageEffect(aura, multiplier)
}

const SpellCritEffectCategory = "spellcritdebuff"

func CriticalMassAura(target *Unit) *Aura {
	return majorSpellCritDebuffAura(target, "Critical Mass", ActionID{SpellID: 22959}, 5)
}

func ShadowAndFlameAura(target *Unit) *Aura {
	return majorSpellCritDebuffAura(target, "Shadow and Flame", ActionID{SpellID: 17800}, 5)
}

func majorSpellCritDebuffAura(target *Unit, label string, actionID ActionID, percent float64) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    label,
		ActionID: actionID,
		Duration: time.Second * 30,
	})

	bonusSpellCrit := percent * CritRatingPerCritChance
	aura.NewExclusiveEffect(SpellCritEffectCategory, true, ExclusiveEffect{
		Priority: percent,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.BonusSpellCritRatingTaken += bonusSpellCrit
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.BonusSpellCritRatingTaken -= bonusSpellCrit
		},
	})
	return aura
}

var majorArmorReductionEffectCategory = "MajorArmorReduction"

func SunderArmorAura(target *Unit) *Aura {
	var effect *ExclusiveEffect
	aura := target.GetOrRegisterAura(Aura{
		Label:     "Sunder Armor",
		ActionID:  ActionID{SpellID: 58567},
		Duration:  time.Second * 30,
		MaxStacks: 3,
		OnStacksChange: func(_ *Aura, sim *Simulation, oldStacks int32, newStacks int32) {
			effect.SetPriority(sim, 0.04*float64(newStacks))
		},
	})
	effect = registerMajorArpEffect(aura, 0)
	return aura
}

func FaerieFireAura(target *Unit) *Aura {
	var effect *ExclusiveEffect
	aura := target.GetOrRegisterAura(Aura{
		Label:     "Faerie Fire",
		ActionID:  ActionID{SpellID: 770},
		Duration:  time.Minute * 5,
		MaxStacks: 3,
		OnStacksChange: func(aura *Aura, sim *Simulation, oldStacks int32, newStacks int32) {
			effect.SetPriority(sim, 0.04*float64(newStacks))
		},
	})
	effect = registerMajorArpEffect(aura, 0)
	return aura
}

func CorrosiveSpitAura(target *Unit) *Aura {
	var effect *ExclusiveEffect
	aura := target.GetOrRegisterAura(Aura{
		Label:     "Corrosive Spit",
		ActionID:  ActionID{SpellID: 55754},
		Duration:  time.Second * 30,
		MaxStacks: 3,
		OnStacksChange: func(aura *Aura, sim *Simulation, oldStacks int32, newStacks int32) {
			effect.SetPriority(sim, 0.04*float64(newStacks))
		},
	})
	effect = registerMajorArpEffect(aura, 0)
	return aura
}

func registerMajorArpEffect(aura *Aura, initialArp float64) *ExclusiveEffect {
	return aura.NewExclusiveEffect(majorArmorReductionEffectCategory, true, ExclusiveEffect{
		Priority: initialArp,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.ArmorMultiplier *= 1 - ee.Priority
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.ArmorMultiplier /= 1 - ee.Priority
		},
	})
}

func ExposeArmorAura(target *Unit, hasGlyph bool) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "ExposeArmor",
		ActionID: ActionID{SpellID: 8647},
		Duration: time.Second * TernaryDuration(hasGlyph, 42, 30),
	})

	registerMajorArpEffect(aura, 0.12)
	return aura
}

var ShatteringThrowAuraTag = "ShatteringThrow"
var ShatteringThrowDuration = time.Second * 10

func ShatteringThrowAura(target *Unit) *Aura {
	armorReduction := 0.2

	return target.GetOrRegisterAura(Aura{
		Label:    "Shattering Throw",
		Tag:      ShatteringThrowAuraTag,
		ActionID: ActionID{SpellID: 64382},
		Duration: ShatteringThrowDuration,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.ArmorMultiplier *= (1.0 - armorReduction)
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.ArmorMultiplier /= (1.0 - armorReduction)
		},
	})
}

const HuntersMarkAuraTag = "HuntersMark"

func HuntersMarkAura(target *Unit) *Aura {
	bonus := 1772.0 // 443.000000 * 4 @ VoraciousGhost - Hunters Mark and Hawk uses the Unknown class in the SpellScaling
	//Todo: Validate calculation

	aura := target.GetOrRegisterAura(Aura{
		Label:    "HuntersMark",
		Tag:      HuntersMarkAuraTag,
		ActionID: ActionID{SpellID: 1130},
		Duration: NeverExpires,
	})

	aura.NewExclusiveEffect("HuntersMark", true, ExclusiveEffect{
		Priority: bonus,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.BonusRangedAttackPowerTaken += bonus
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.BonusRangedAttackPowerTaken -= bonus
		},
	})

	return aura
}

func CurseOfWeaknessAura(target *Unit) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Curse of Weakness",
		ActionID: ActionID{SpellID: 702},
		Duration: time.Minute * 2,
	})
	PhysDamageReductionEffect(aura, 0.1)
	return aura
}

func DemoralizingRoarAura(target *Unit) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "DemoralizingRoar",
		ActionID: ActionID{SpellID: 99},
		Duration: time.Second * 30,
	})
	PhysDamageReductionEffect(aura, 0.1)
	return aura
}

func DemoralizingShoutAura(target *Unit, glyph bool) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "DemoralizingShout",
		ActionID: ActionID{SpellID: 1160},
		Duration: time.Second*30 + TernaryDuration(glyph, time.Second*15, 0),
	})
	PhysDamageReductionEffect(aura, 0.1)
	return aura
}

func VindicationAura(target *Unit) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Vindication",
		ActionID: ActionID{SpellID: 26017},
		Duration: time.Second * 30,
	})
	PhysDamageReductionEffect(aura, 0.1)
	return aura
}

func DemoralizingScreechAura(target *Unit) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "DemoralizingScreech",
		ActionID: ActionID{SpellID: 24423},
		Duration: time.Second * 10,
	})
	PhysDamageReductionEffect(aura, 0.1)
	return aura
}

func ScarletFeverAura(target *Unit, points int32, epidemic int32) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Scarlet Fever",
		ActionID: ActionID{SpellID: 81130},
		Duration: time.Second * time.Duration(21+epidemic*4),
	})
	PhysDamageReductionEffect(aura, 0.05*float64(points))
	return aura
}

func PhysDamageTakenEffect(aura *Aura, multiplier float64) *ExclusiveEffect {
	return aura.NewExclusiveEffect("PhysicalDmg", false, ExclusiveEffect{
		Priority: multiplier,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexPhysical] *= multiplier
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexPhysical] /= multiplier
		},
	})
}

func PhysDamageReductionEffect(aura *Aura, dmgReduction float64) *ExclusiveEffect {
	reductionMult := 1.0 - dmgReduction
	return aura.NewExclusiveEffect("PhysDamageReduction", false, ExclusiveEffect{
		Priority: dmgReduction,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= reductionMult
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= reductionMult
		},
	})
}

func apReductionEffect(aura *Aura, apReduction float64) *ExclusiveEffect {
	statReduction := stats.Stats{stats.AttackPower: -apReduction}
	return aura.NewExclusiveEffect("APReduction", false, ExclusiveEffect{
		Priority: apReduction,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.AddStatsDynamic(sim, statReduction)
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.AddStatsDynamic(sim, statReduction.Invert())
		},
	})
}

func ThunderClapAura(target *Unit) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "ThunderClap",
		ActionID: ActionID{SpellID: 6343},
		Duration: time.Second * 30,
	})
	AtkSpeedReductionEffect(aura, 1.2)
	return aura
}

func InfectedWoundsAura(target *Unit, points int32) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "InfectedWounds-" + strconv.Itoa(int(points)),
		ActionID: ActionID{SpellID: 48485},
		Duration: time.Second * 12,
	})
	AtkSpeedReductionEffect(aura, 1+0.1*float64(points))
	return aura
}

// Note: Paladin code might apply this as part of their judgement auras instead
// of using another separate aura.
func JudgementsOfTheJustAura(target *Unit, points int32) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "JudgementsOfTheJust-" + strconv.Itoa(int(points)),
		ActionID: ActionID{SpellID: 53696},
		Duration: time.Second * 30,
	})
	AtkSpeedReductionEffect(aura, 1.0+0.1*float64(points))
	return aura
}

func DustCloud(target *Unit) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Dust Cloud",
		ActionID: ActionID{SpellID: 50285},
		Duration: time.Second * 30,
	})
	AtkSpeedReductionEffect(aura, 1.2)
	return aura
}

func FrostFeverAura(target *Unit, britleBones int32) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "FrostFeverDebuff",
		ActionID: ActionID{SpellID: 55095},
		Duration: NeverExpires,
	})
	AtkSpeedReductionEffect(aura, 1.2)
	if britleBones > 0 {
		PhysDamageTakenEffect(aura, 1+0.02*float64(britleBones))
	}
	return aura
}

func AtkSpeedReductionEffect(aura *Aura, speedMultiplier float64) *ExclusiveEffect {
	return aura.NewExclusiveEffect("AtkSpdReduction", false, ExclusiveEffect{
		Priority: speedMultiplier,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.MultiplyAttackSpeed(sim, 1/speedMultiplier)
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.MultiplyAttackSpeed(sim, speedMultiplier)
		},
	})
}

func MarkOfBloodAura(target *Unit) *Aura {
	actionId := ActionID{SpellID: 49005}

	var healthMetrics *ResourceMetrics
	aura := target.GetOrRegisterAura(Aura{
		Label:     "MarkOfBlood",
		ActionID:  actionId,
		Duration:  20 * time.Second,
		MaxStacks: 20,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.SetStacks(sim, aura.MaxStacks)

			target := aura.Unit.CurrentTarget

			if healthMetrics == nil && target != nil {
				healthMetrics = target.NewHealthMetrics(actionId)
			}
		},
		OnSpellHitDealt: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			target := aura.Unit.CurrentTarget

			if target != nil && result.Landed() {
				// Vampiric Blood bonus max health is ignored in MoB calculation (maybe other Max health effects as well?)
				targetHealth := target.MaxHealth()
				if target.HasActiveAura("Vampiric Blood") {
					targetHealth /= 1.15
				}
				// Current testing shows 5% healing instead of 4% as stated in the tooltip
				target.GainHealth(sim, targetHealth*0.05*target.PseudoStats.HealingTakenMultiplier, healthMetrics)
				aura.RemoveStack(sim)

				if aura.GetStacks() == 0 {
					aura.Deactivate(sim)
				}
			}
		},
	})
	return aura
}

func InsectSwarmAura(target *Unit) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "InsectSwarmMiss",
		ActionID: ActionID{SpellID: 27013},
		Duration: time.Second * 12,
	})
	increasedMissEffect(aura, 0.03)
	return aura
}

func increasedMissEffect(aura *Aura, increasedMissChance float64) *ExclusiveEffect {
	return aura.NewExclusiveEffect("IncreasedMiss", false, ExclusiveEffect{
		Priority: increasedMissChance,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.IncreasedMissChance += increasedMissChance
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.IncreasedMissChance -= increasedMissChance
		},
	})
}

func minorCritDebuffAura(target *Unit, label string, actionID ActionID, duration time.Duration, critBonus float64) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    label,
		ActionID: actionID,
		Duration: duration,
	})
	critBonusEffect(aura, critBonus)
	return aura
}

func critBonusEffect(aura *Aura, critBonus float64) *ExclusiveEffect {
	return aura.NewExclusiveEffect("CritBonus", false, ExclusiveEffect{
		Priority: critBonus,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.BonusCritRatingTaken += critBonus
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.BonusCritRatingTaken -= critBonus
		},
	})
}

func CrystalYieldAura(target *Unit) *Aura {
	return target.GetOrRegisterAura(Aura{
		Label:    "Crystal Yield",
		ActionID: ActionID{SpellID: 15235},
		Duration: 2 * time.Minute,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.stats[stats.Armor] -= 200
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.stats[stats.Armor] += 200
		},
	})
}
