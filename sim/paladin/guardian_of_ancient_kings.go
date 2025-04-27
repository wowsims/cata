package paladin

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func (paladin *Paladin) goakBaseDuration() time.Duration {
	if paladin.Spec == proto.Spec_SpecProtectionPaladin {
		return time.Second * 12
	} else {
		return time.Second * 30
	}
}

func (paladin *Paladin) registerGuardianOfAncientKings() {
	duration := paladin.goakBaseDuration()

	var spell *core.Spell
	switch paladin.Spec {
	case proto.Spec_SpecHolyPaladin:
		spell = paladin.registerHolyGuardian(duration)
	case proto.Spec_SpecProtectionPaladin:
		spell = paladin.registerProtectionGuardian(duration)
	default:
	case proto.Spec_SpecRetributionPaladin:
		spell = paladin.registerRetributionGuardian(duration)
	}

	paladin.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}

func (paladin *Paladin) registerHolyGuardian(duration time.Duration) *core.Spell {
	actionID := core.ActionID{SpellID: 86150}

	paladin.GoakAura = paladin.RegisterAura(core.Aura{
		Label:    "Guardian of Ancient Kings" + paladin.Label,
		ActionID: actionID,
		Duration: duration,

		OnHealDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// TODO: Trigger guardians heal if `spell` a single-target heal
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			paladin.AncientGuardian.Pet.Disable(sim)
		},
	})

	return paladin.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskGuardianOfAncientKings,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Minute * 5,
			},
		},

		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			paladin.GoakAura.Activate(sim)
			paladin.AncientGuardian.Enable(sim, paladin.AncientGuardian)
			paladin.AncientGuardian.CancelGCDTimer(sim)
		},
	})
}

func (paladin *Paladin) registerProtectionGuardian(duration time.Duration) *core.Spell {
	actionID := core.ActionID{SpellID: 86150}

	paladin.GoakAura = paladin.RegisterAura(core.Aura{
		Label:    "Guardian of Ancient Kings" + paladin.Label,
		ActionID: actionID,
		Duration: duration,

		// TODO: Perhaps refactor this to also be a pet with a channeled cast?
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			paladin.PseudoStats.DamageTakenMultiplier *= 0.5
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			paladin.PseudoStats.DamageTakenMultiplier /= 0.5
		},
	})

	return paladin.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskGuardianOfAncientKings,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Minute * 5,
			},
		},

		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			paladin.GoakAura.Activate(sim)
		},
	})
}

func (paladin *Paladin) registerRetributionGuardian(duration time.Duration) *core.Spell {
	var strDepByStackCount = map[int32]*stats.StatDependency{}

	for i := 1; i <= 20; i++ {
		strDepByStackCount[int32(i)] = paladin.NewDynamicMultiplyStat(stats.Strength, 1.0+0.01*float64(i))
	}

	paladin.AncientPowerAura = paladin.RegisterAura(core.Aura{
		Label:     "Ancient Power" + paladin.Label,
		ActionID:  core.ActionID{SpellID: 86700},
		Duration:  duration + time.Second*1,
		MaxStacks: 20,

		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			if oldStacks > 0 {
				paladin.DisableDynamicStatDep(sim, strDepByStackCount[oldStacks])
			}

			if newStacks > 0 {
				paladin.EnableDynamicStatDep(sim, strDepByStackCount[newStacks])
			}
		},
	})

	ancientFuryMinDamage, ancientFuryMaxDamage :=
		core.CalcScalingSpellEffectVarianceMinMax(proto.Class_ClassPaladin, 0.23659999669, 0.30000001192)
	numTargets := paladin.Env.GetNumTargets()
	results := make([]*core.SpellResult, numTargets)

	ancientFury := paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 86704},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: SpellMaskAncientFury,
		Flags:          core.SpellFlagPassiveSpell,

		MaxRange: 10,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.RollWithLabel(ancientFuryMinDamage, ancientFuryMaxDamage, "Ancient Fury"+paladin.Label) +
				0.06100000069*spell.SpellPower()

			// Deals X Holy damage per application of Ancient Power,
			// divided evenly among all targets within 10 yards.
			baseDamage *= float64(paladin.AncientPowerAura.GetStacks())
			baseDamage /= float64(numTargets)

			for idx := int32(0); idx < numTargets; idx++ {
				currentTarget := sim.Environment.GetTargetUnit(idx)
				results[idx] = spell.CalcDamage(sim, currentTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}

			for idx := int32(0); idx < numTargets; idx++ {
				spell.DealDamage(sim, results[idx])
			}
		},
	})

	actionID := core.ActionID{SpellID: 86150}

	paladin.GoakAura = paladin.RegisterAura(core.Aura{
		Label:    "Guardian of Ancient Kings" + paladin.Label,
		ActionID: actionID,
		Duration: duration,

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || spell.ClassSpellMask&SpellMaskCanTriggerAncientPower == 0 {
				return
			}

			paladin.AncientPowerAura.AddStack(sim)
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			paladin.AncientPowerAura.Activate(sim)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			paladin.AncientGuardian.Pet.Disable(sim)
			ancientFury.Cast(sim, paladin.CurrentTarget)
			paladin.AncientPowerAura.Deactivate(sim)
		},
	})

	return paladin.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL,
		ProcMask:       core.ProcMaskEmpty,
		ClassSpellMask: SpellMaskGuardianOfAncientKings,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Minute * 5,
			},
		},

		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			paladin.GoakAura.Activate(sim)
			paladin.AncientGuardian.Enable(sim, paladin.AncientGuardian)
			paladin.AncientGuardian.CancelGCDTimer(sim)
		},
	})
}
