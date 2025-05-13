package warrior

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

const cdDuration = time.Second * 3

func (warrior *Warrior) RegisterHeroicStrikeSpell() {
	warrior.HeroicStrike = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 78},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagAPL | core.SpellFlagMeleeMetrics,
		ClassSpellMask: SpellMaskHeroicStrike | SpellMaskSpecialAttack,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   30,
			Refund: 0.8,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 0,
			},
			CD: core.Cooldown{
				Timer:    warrior.hsCleaveCD,
				Duration: cdDuration,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		FlatThreatBonus:  259,
		CritMultiplier:   warrior.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 8 + (spell.MeleeAttackPower() * 0.6)

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}

			spell.DealDamage(sim, result)
		},
	})
}

func (warrior *Warrior) RegisterCleaveSpell() {
	targets := core.TernaryInt32(warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfCleaving), 3, 2)
	numHits := min(targets, warrior.Env.GetNumTargets())
	results := make([]*core.SpellResult, numHits)

	warrior.Cleave = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 845},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagAPL | core.SpellFlagMeleeMetrics,
		ClassSpellMask: SpellMaskCleave | SpellMaskSpecialAttack,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost: 30,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 0,
			},
			CD: core.Cooldown{
				Timer:    warrior.hsCleaveCD,
				Duration: cdDuration,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		FlatThreatBonus:  225,
		CritMultiplier:   warrior.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			curTarget := target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				baseDamage := 6 + (spell.MeleeAttackPower() * 0.45)
				results[hitIndex] = spell.CalcDamage(sim, curTarget, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}

			curTarget = target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				spell.DealDamage(sim, results[hitIndex])
				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}
		},
	})
}
