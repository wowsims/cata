package paladin

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

func applyT11Prot4pcBonus(duration time.Duration) time.Duration {
	return time.Millisecond * time.Duration(float64(duration.Milliseconds())*1.5)
}

func (paladin *Paladin) registerGuardianOfAncientKings() {
	hasT11Prot4pc := paladin.HasSetBonus(ItemSetReinforcedSapphiriumBattlearmor, 4)

	var duration time.Duration
	if paladin.Spec == proto.Spec_SpecProtectionPaladin {
		duration = time.Second * 12
	} else {
		duration = time.Second * 30
	}

	if hasT11Prot4pc {
		duration = applyT11Prot4pcBonus(duration)
	}

	var spell *core.Spell
	switch paladin.Spec {
	case proto.Spec_SpecHolyPaladin:
		spell = paladin.registerHolyGuardian(duration)
	case proto.Spec_SpecProtectionPaladin:
		spell = paladin.registerProtectionGuardian(duration)
	default:
	case proto.Spec_SpecRetributionPaladin:
		spell = paladin.registerRetributionGuardian(duration, paladin.SnapshotGuardian && !hasT11Prot4pc)
	}

	paladin.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}

func (paladin *Paladin) registerHolyGuardian(duration time.Duration) *core.Spell {
	actionID := core.ActionID{SpellID: 86150}

	goakAura := paladin.RegisterAura(core.Aura{
		Label:    "Guardian of Ancient Kings",
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
			goakAura.Activate(sim)
			paladin.AncientGuardian.Enable(sim, paladin.AncientGuardian)
			paladin.AncientGuardian.CancelGCDTimer(sim)
		},
	})
}

func (paladin *Paladin) registerProtectionGuardian(duration time.Duration) *core.Spell {
	actionID := core.ActionID{SpellID: 86150}

	goakAura := paladin.RegisterAura(core.Aura{
		Label:    "Guardian of Ancient Kings",
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
			goakAura.Activate(sim)
		},
	})
}

func (paladin *Paladin) registerRetributionGuardian(duration time.Duration, snapshotGuardian bool) *core.Spell {
	var strDepByStackCount = map[int32]*stats.StatDependency{}

	for i := 1; i <= 20; i++ {
		strDepByStackCount[int32(i)] = paladin.NewDynamicMultiplyStat(stats.Strength, 1.0+0.01*float64(i))
	}

	ancientPowerDuration := duration + time.Second*1
	ancientPower := paladin.RegisterAura(core.Aura{
		Label:     "Ancient Power",
		ActionID:  core.ActionID{SpellID: 86700},
		Duration:  ancientPowerDuration,
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
		core.CalcScalingSpellEffectVarianceMinMax(proto.Class_ClassPaladin, 0.2366, 0.3)
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
		CritMultiplier:   paladin.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.RollWithLabel(ancientFuryMinDamage, ancientFuryMaxDamage, "Ancient Fury") +
				0.061*spell.SpellPower()

			// Deals X Holy damage per application of Ancient Power,
			// divided evenly among all targets within 10 yards.
			baseDamage *= float64(ancientPower.GetStacks())
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

	goakAura := paladin.RegisterAura(core.Aura{
		Label:    "Guardian of Ancient Kings",
		ActionID: actionID,
		Duration: duration,

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || spell.ClassSpellMask&SpellMaskCanTriggerAncientPower == 0 {
				return
			}

			ancientPower.AddStack(sim)
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			ancientPower.Activate(sim)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			paladin.AncientGuardian.Pet.Disable(sim)
			ancientFury.Cast(sim, paladin.CurrentTarget)
			ancientPower.Deactivate(sim)

			// Deactivate T11 Prot 4pc bonus if configured and activated during prepull
			if snapshotGuardian && (aura.Duration != duration || ancientPower.Duration != ancientPowerDuration) {
				aura.Duration = duration
				ancientPower.Duration = ancientPowerDuration
			}
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
			// Activate T11 Prot 4pc bonus if configured and activated during prepull
			if sim.CurrentTime < 0 && snapshotGuardian {
				goakAura.Duration = applyT11Prot4pcBonus(duration)
				ancientPower.Duration = applyT11Prot4pcBonus(ancientPowerDuration)
			}

			goakAura.Activate(sim)
			paladin.AncientGuardian.Enable(sim, paladin.AncientGuardian)
			paladin.AncientGuardian.CancelGCDTimer(sim)
		},
	})
}
