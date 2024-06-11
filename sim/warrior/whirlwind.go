package warrior

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (warrior *Warrior) RegisterWhirlwindSpell() {
	actionID := core.ActionID{SpellID: 1680}
	numHits := warrior.Env.GetNumTargets() // Whirlwind is uncapped in Cata
	results := make([]*core.SpellResult, numHits)

	var whirlwindOH *core.Spell
	if warrior.AutoAttacks.IsDualWielding && warrior.GetOHWeapon().WeaponType != proto.WeaponType_WeaponTypeStaff &&
		warrior.GetOHWeapon().WeaponType != proto.WeaponType_WeaponTypePolearm {
		whirlwindOH = warrior.RegisterSpell(core.SpellConfig{
			ActionID:    actionID.WithTag(1),
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty, // whirlwind offhand hits usually don't proc auras
			Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagNoOnCastComplete,

			DamageMultiplier: 1.0 * warrior.DualWieldSpecialization(),
			ThreatMultiplier: 1.25,
			CritMultiplier:   warrior.DefaultMeleeCritMultiplier(),

			BonusCoefficient: 1,
		})
	}

	warrior.Whirlwind = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
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
			numLandedHits := 0
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				baseDamage := 0.65 * spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
				results[hitIndex] = spell.CalcDamage(sim, curTarget, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
				if results[hitIndex].Landed() {
					numLandedHits++
				}

				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}

			curTarget = target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				spell.DealDamage(sim, results[hitIndex])
				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}

			if numLandedHits > 4 {
				spell.CD.Reduce(time.Second * 6)
			}

			if whirlwindOH != nil {
				curTarget = target
				for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
					baseDamage := 0.65 * spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
					results[hitIndex] = whirlwindOH.CalcDamage(sim, curTarget, baseDamage, whirlwindOH.OutcomeMeleeWeaponSpecialHitAndCrit)

					curTarget = sim.Environment.NextTargetUnit(curTarget)
				}

				curTarget = target
				for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
					whirlwindOH.DealDamage(sim, results[hitIndex])
					curTarget = sim.Environment.NextTargetUnit(curTarget)
				}
			}
		},
	})
}
