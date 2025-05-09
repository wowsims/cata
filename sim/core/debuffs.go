package core

import (
	"time"

	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

// applyRaidDebuffEffects applies all raid-level debuffs based on the provided Debuffs proto.
func applyDebuffEffects(target *Unit, targetIdx int, debuffs *proto.Debuffs, raid *proto.Raid) {
	// –10% Physical damage dealt for 30s
	if debuffs.WeakenedBlows {
		MakePermanent(WeakenedBlowsAura(target))
	}

	// +4% Physical damage taken for 30s
	if debuffs.PhysicalVulnerability {
		MakePermanent(PhysVulnerabilityAura(target))
	}

	// –4% Armor for 30s, stacks 3 times
	if debuffs.WeakenedArmor {
		aura := MakePermanent(WeakenedArmorAura(target))

		aura.ApplyOnGain(func(aura *Aura, sim *Simulation) {
			aura.SetStacks(sim, 3)
		})
	}

	// Spell‐damage‐taken sources
	if debuffs.FireBreath {
		MakePermanent(FireBreathDebuff(target))
	}
	if debuffs.LightningBreath {
		MakePermanent(LightningBreathDebuff(target))
	}
	if debuffs.MasterPoisoner {
		MakePermanent(MasterPoisonerDebuff(target))
	}
	if debuffs.CurseOfElements {
		MakePermanent(CurseOfElementsAura(target))
	}

	// Casting‐speed‐reduction sources
	if debuffs.NecroticStrike {
		MakePermanent(NecroticStrikeAura(target))
	}
	if debuffs.LavaBreath {
		MakePermanent(LavaBreathAura(target))
	}
	if debuffs.SporeCloud {
		MakePermanent(SporeCloud(target))
	}
	if debuffs.Slow {
		MakePermanent(SlowAura(target))
	}
	if debuffs.MindNumbingPoison {
		MakePermanent(MindNumbingPoisonAura(target))
	}
	if debuffs.CurseOfEnfeeblement {
		MakePermanent(CurseOfEnfeeblement(target))
	}
}

// –10% Physical damage dealt
func WeakenedBlowsAura(target *Unit) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Weakened Blows",
		ActionID: ActionID{SpellID: 115798},
		Duration: time.Second * 30,
	})
	PhysDamageReductionEffect(aura, 0.1)
	return aura
}

// +4% Physical damage taken
func PhysVulnerabilityAura(target *Unit) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Physical Vulnerability",
		ActionID: ActionID{SpellID: 81326},
		Duration: time.Second * 30,
	})
	PhysDamageTakenEffect(aura, 1.04)
	return aura
}

// –4% Armor stacks 3
func WeakenedArmorAura(target *Unit) *Aura {
	var effect *ExclusiveEffect
	aura := target.GetOrRegisterAura(Aura{
		Label:     "Weakened Armor",
		ActionID:  ActionID{SpellID: 113746},
		Duration:  time.Second * 30,
		MaxStacks: 3,
		OnStacksChange: func(_ *Aura, sim *Simulation, oldStacks int32, newStacks int32) {
			effect.SetPriority(sim, 0.04*float64(newStacks))
		},
	})
	effect = registerMajorArpEffect(aura, 0)
	return aura
}

func MortalWoundsAura(target *Unit) *Aura {
	return majorHealingReductionAura(target, "Mortal Wounds", 115804, 0.25)
}

// Spell‐damage‐taken sources
func FireBreathDebuff(target *Unit) *Aura {
	return spellDamageEffectAura(Aura{Label: "Fire Breath", ActionID: ActionID{SpellID: 34889}, Duration: time.Second * 45}, target, 1.05)
}
func LightningBreathDebuff(target *Unit) *Aura {
	return spellDamageEffectAura(Aura{Label: "Lightning Breath", ActionID: ActionID{SpellID: 24844}, Duration: time.Second * 45}, target, 1.05)
}
func MasterPoisonerDebuff(target *Unit) *Aura {
	return spellDamageEffectAura(Aura{Label: "Master Poisoner", ActionID: ActionID{SpellID: 58410}, Duration: time.Second * 15}, target, 1.05)
}
func CurseOfElementsAura(target *Unit) *Aura {
	return spellDamageEffectAura(Aura{Label: "Curse of Elements", ActionID: ActionID{SpellID: 1490}, Duration: time.Minute * 5}, target, 1.05)
}

func majorHealingReductionAura(target *Unit, label string, spellID int32, multiplier float64) *Aura {
	aura := target.GetOrRegisterAura(Aura{Label: label, ActionID: ActionID{SpellID: spellID}, Duration: time.Second * 30})
	aura.NewExclusiveEffect("HealingReduction", false, ExclusiveEffect{
		Priority: multiplier,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.HealingTakenMultiplier *= multiplier
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.HealingTakenMultiplier /= multiplier
		},
	})
	return aura
}

// Casting‐speed‐reduction sources
func NecroticStrikeAura(target *Unit) *Aura {
	return castSpeedReductionAura(target, "Necrotic Strike", 73975, 1.5)
}
func LavaBreathAura(target *Unit) *Aura {
	return castSpeedReductionAura(target, "Lava Breath", 58604, 1.5)
}
func SporeCloud(target *Unit) *Aura { return castSpeedReductionAura(target, "Spore Cloud", 50274, 1.5) }
func MindNumbingPoisonAura(target *Unit) *Aura {
	return castSpeedReductionAura(target, "Mind-numbing Poison", 5761, 1.5)
}
func CurseOfEnfeeblement(target *Unit) *Aura {
	return castSpeedReductionAura(target, "Curse of Enfeeblement", 109466, 1.5)
}
func SlowAura(target *Unit) *Aura {
	return castSpeedReductionAura(target, "Slow", 31589, 1.5)
}
func castSpeedReductionAura(target *Unit, label string, spellID int32, multiplier float64) *Aura {
	aura := target.GetOrRegisterAura(Aura{Label: label, ActionID: ActionID{SpellID: spellID}, Duration: time.Second * 30})
	aura.NewExclusiveEffect("CastSpdReduction", false, ExclusiveEffect{
		Priority: multiplier,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.MultiplyCastSpeed(1 / multiplier)
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.MultiplyCastSpeed(multiplier)
		},
	})
	return aura
}

const SpellDamageEffectAuraTag = "SpellDamageAuraTag"

func spellDamageEffectAura(auraConfig Aura, target *Unit, multiplier float64) *Aura {
	auraConfig.Tag = SpellDamageEffectAuraTag
	aura := target.GetOrRegisterAura(auraConfig)
	aura.NewExclusiveEffect("SpellDamageTaken%", true, ExclusiveEffect{
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
	return aura
}

var majorArmorReductionEffectCategory = "MajorArmorReduction"

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

var ShatteringThrowAuraTag = "ShatteringThrow"
var ShatteringThrowDuration = time.Second * 10

func ShatteringThrowAura(target *Unit, actionTag int32) *Aura {
	armorReduction := 0.2

	return target.GetOrRegisterAura(Aura{
		Label:    "Shattering Throw",
		Tag:      ShatteringThrowAuraTag,
		ActionID: ActionID{SpellID: 64382, Tag: actionTag},
		Duration: ShatteringThrowDuration,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.ArmorMultiplier *= (1.0 - armorReduction)
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.ArmorMultiplier /= (1.0 - armorReduction)
		},
	})
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
