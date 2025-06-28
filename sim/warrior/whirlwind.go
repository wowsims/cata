package warrior

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (warrior *Warrior) RegisterWhirlwindSpell() {
	actionID := core.ActionID{SpellID: 1680}
	maxHits := warrior.Env.TotalTargetCount() // Whirlwind is uncapped in Cata
	results := make([]*core.SpellResult, maxHits)

	var whirlwindOH *core.Spell
	if warrior.AutoAttacks.IsDualWielding && warrior.GetOHWeapon().WeaponType != proto.WeaponType_WeaponTypeStaff &&
		warrior.GetOHWeapon().WeaponType != proto.WeaponType_WeaponTypePolearm {
		whirlwindOH = warrior.RegisterSpell(core.SpellConfig{
			ActionID:       actionID.WithTag(2),
			SpellSchool:    core.SpellSchoolPhysical,
			ProcMask:       core.ProcMaskMeleeOHSpecial, //TODO: needs testing to check if it procs auras, according to pre-cata it didn't
			ClassSpellMask: SpellMaskWhirlwindOh,
			Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagNoOnCastComplete,

			DamageMultiplier: 1.0,
			ThreatMultiplier: 1.25,
			CritMultiplier:   warrior.DefaultMeleeCritMultiplier(),

			BonusCoefficient: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				curTarget := target
				numHits := sim.Environment.ActiveTargetCount()
				for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
					baseDamage := 0.65 * spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
					results[hitIndex] = whirlwindOH.CalcDamage(sim, curTarget, baseDamage, whirlwindOH.OutcomeMeleeWeaponSpecialHitAndCrit)

					curTarget = sim.Environment.NextActiveTargetUnit(curTarget)
				}

				for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
					whirlwindOH.DealDamage(sim, results[hitIndex])
				}
			},
		})
	}

	warrior.Whirlwind = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID.WithTag(1),
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskWhirlwind | SpellMaskSpecialAttack,

		RageCost: core.RageCostOptions{
			Cost: 25,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 10,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(BerserkerStance)
		},

		DamageMultiplier: 1.0,
		ThreatMultiplier: 1.25,
		CritMultiplier:   warrior.DefaultMeleeCritMultiplier(),

		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			curTarget := target
			numHits := sim.Environment.ActiveTargetCount()
			numLandedHits := 0
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				baseDamage := 0.65 * spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
				results[hitIndex] = spell.CalcDamage(sim, curTarget, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
				if results[hitIndex].Landed() {
					numLandedHits++
				}

				curTarget = sim.Environment.NextActiveTargetUnit(curTarget)
			}

			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				spell.DealDamage(sim, results[hitIndex])
			}

			if numLandedHits >= 4 {
				spell.CD.Reduce(time.Second * 6)
			}

			if whirlwindOH != nil {
				whirlwindOH.Cast(sim, target)
			}
		},
	})
}
