package paladin

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func (paladin *Paladin) goakBaseDuration() time.Duration {
	switch paladin.Spec {
	case proto.Spec_SpecHolyPaladin:
		return time.Second * 15
	case proto.Spec_SpecProtectionPaladin:
		return time.Second * 12
	default:
		return time.Second * 30
	}
}

func (paladin *Paladin) registerGuardianOfAncientKings() {
	duration := paladin.goakBaseDuration()

	switch paladin.Spec {
	case proto.Spec_SpecHolyPaladin:
		paladin.registerHolyGuardian(duration)
	case proto.Spec_SpecProtectionPaladin:
		paladin.registerProtectionGuardian(duration)
	default:
	case proto.Spec_SpecRetributionPaladin:
		paladin.registerRetributionGuardian(duration)
	}
}

/*
Summons a Guardian of Ancient Kings to help you heal for 15 sec.

The Guardian of Ancient Kings will heal the targets of your heals for an additional 100% of the amount healed and grants you 10% haste for its duration.
*/
func (paladin *Paladin) registerHolyGuardian(duration time.Duration) {
	actionID := core.ActionID{SpellID: 86669}

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
	}).AttachMultiplyCastSpeed(1.1)

	spell := paladin.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskGuardianOfAncientKings,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			paladin.GoakAura.Activate(sim)
			paladin.AncientGuardian.Enable(sim, paladin.AncientGuardian)
			paladin.AncientGuardian.CancelGCDTimer(sim)
		},
	})

	paladin.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeSurvival,
	})
}

/*
Summons a Guardian of Ancient Kings to protect you for 12 sec.

The Guardian of Ancient Kings reduces damage taken by 50%.
*/
func (paladin *Paladin) registerProtectionGuardian(duration time.Duration) {
	actionID := core.ActionID{SpellID: 86659}

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

	spell := paladin.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL | core.SpellFlagHelpful,
		ClassSpellMask: SpellMaskGuardianOfAncientKings,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			paladin.GoakAura.Activate(sim)
		},
	})

	paladin.AddDefensiveCooldownAura(paladin.GoakAura)
	paladin.AddMajorCooldown(core.MajorCooldown{
		Spell:    spell,
		Type:     core.CooldownTypeSurvival,
		Priority: core.CooldownPriorityLow + 20,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return !paladin.AnyActiveDefensiveCooldown()
		},
	})
}

/*
Summons a Guardian of Ancient Kings to help you deal damage for 30 sec.

The Guardian of Ancient Kings will attack your current enemy.
Both your attacks and the attacks of the Guardian will infuse you with Ancient Power that is unleashed as Ancient Fury when the Guardian departs.

Ancient Power
Strength increased by 1%.
When your Guardian of Ancient Kings departs, you release Ancient Fury, causing (<229-311> + 0.107 * <SP>) Holy damage, split among all enemies within 10 yards.

Ancient Fury
Unleash the fury of ancient kings, causing (<229-311> + 0.107 * <SP>) Holy damage per application of Ancient Power, divided evenly among all targets within 10 yards.
*/
func (paladin *Paladin) registerRetributionGuardian(duration time.Duration) {
	var strDepByStackCount = map[int32]*stats.StatDependency{}

	for i := 1; i <= 12; i++ {
		strDepByStackCount[int32(i)] = paladin.NewDynamicMultiplyStat(stats.Strength, 1.0+0.01*float64(i))
	}

	paladin.AncientPowerAura = paladin.RegisterAura(core.Aura{
		Label:     "Ancient Power" + paladin.Label,
		ActionID:  core.ActionID{SpellID: 86700},
		Duration:  duration + time.Second*1,
		MaxStacks: 12,

		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			if oldStacks > 0 {
				paladin.DisableDynamicStatDep(sim, strDepByStackCount[oldStacks])
			}

			if newStacks > 0 {
				paladin.EnableDynamicStatDep(sim, strDepByStackCount[newStacks])
			}
		},
	})

	numTargets := paladin.Env.GetNumTargets()

	ancientFury := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 86704},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagPassiveSpell,

		MaxRange: 10,

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := paladin.CalcAndRollDamageRange(sim, 0.23659999669, 0.30000001192) +
				0.10700000077*spell.SpellPower()

			// Deals X Holy damage per application of Ancient Power,
			// divided evenly among all targets within 10 yards.
			baseDamage *= float64(paladin.AncientPowerAura.GetStacks())
			baseDamage /= float64(numTargets)

			results := make([]*core.SpellResult, numTargets)
			for idx := range numTargets {
				currentTarget := sim.Environment.GetTargetUnit(idx)
				results[idx] = spell.CalcDamage(sim, currentTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}

			for idx := range numTargets {
				spell.DealDamage(sim, results[idx])
			}
		},
	})

	actionID := core.ActionID{SpellID: 86698}

	paladin.GoakAura = paladin.RegisterAura(core.Aura{
		Label:    "Guardian of Ancient Kings" + paladin.Label,
		ActionID: actionID,
		Duration: duration,

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() ||
				(!spell.Matches(SpellMaskCanTriggerAncientPower) && !spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit)) {
				return
			}

			paladin.AncientPowerAura.AddStack(sim)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			paladin.AncientGuardian.Pet.Disable(sim)
			ancientFury.Cast(sim, paladin.CurrentTarget)
		},
	}).AttachDependentAura(paladin.AncientPowerAura)

	spell := paladin.RegisterSpell(core.SpellConfig{
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
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			spell.RelatedSelfBuff.Activate(sim)
			paladin.AncientGuardian.Enable(sim, paladin.AncientGuardian)
			paladin.AncientGuardian.CancelGCDTimer(sim)
		},

		RelatedSelfBuff: paladin.GoakAura,
	})

	paladin.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}
