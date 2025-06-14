package enhancement

import (
	"slices"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/shaman"
)

func (enh *EnhancementShaman) getSearingFlamesMultiplier() float64 {
	return enh.SearingFlamesMultiplier + core.TernaryFloat64(enh.T12Enh2pc.IsActive(), 0.05, 0)
}

func (enh *EnhancementShaman) registerLavaLashSpell() {
	damageMultiplier := 3.0
	if enh.SelfBuffs.ImbueOH == proto.ShamanImbue_FlametongueWeapon {
		damageMultiplier *= 1.4
	}

	enh.LavaLash = enh.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 60103},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskMeleeOHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: shaman.SpellMaskLavaLash,
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 4,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    enh.NewTimer(),
				Duration: time.Second * 10,
			},
		},
		DamageMultiplier: damageMultiplier,
		CritMultiplier:   enh.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			searingFlamesBonus := 1.0

			baseDamage *= searingFlamesBonus
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if !result.Landed() {
				return
			}

			if !enh.HasMinorGlyph(proto.ShamanMinorGlyph_GlyphOfLavaLash) {
				flameShockDot := enh.FlameShock.Dot(target)

				if flameShockDot != nil && flameShockDot.IsActive() {
					numberSpread := 0
					maxTargets := min(4, len(sim.Encounter.TargetUnits))
					sortedTargets := make([]*core.Unit, len(sim.Encounter.TargetUnits))
					copy(sortedTargets, sim.Encounter.TargetUnits)
					slices.SortFunc(sortedTargets, func(a *core.Unit, b *core.Unit) int {
						aDot := enh.FlameShock.Dot(a)
						if aDot == nil || !aDot.IsActive() {
							return -1
						}
						bDot := enh.FlameShock.Dot(b)
						if bDot == nil || !bDot.IsActive() {
							return 1
						}
						return int(aDot.RemainingDuration(sim) - bDot.RemainingDuration(sim))
					})

					for _, otherTarget := range sortedTargets {
						if otherTarget == target {
							return
						}

						enh.FlameShock.RelatedDotSpell.Dot(otherTarget).CopyDotAndApply(sim, flameShockDot)
						numberSpread++

						if numberSpread >= maxTargets {
							return
						}
					}
				}
			}
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return enh.HasOHWeapon()
		},
	})
}

func (enh *EnhancementShaman) IsLavaLashCastable(sim *core.Simulation) bool {
	return enh.LavaLash.IsReady(sim)
}
