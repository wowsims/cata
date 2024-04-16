package shaman

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (shaman *Shaman) registerLavaLashSpell() {
	damageMultiplier := 2.6
	if shaman.SelfBuffs.ImbueOH == proto.ShamanImbue_FlametongueWeapon {
		damageMultiplier += 0.4
	}

	shaman.LavaLash = shaman.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 78146},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskMeleeOHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskLavaLash,
		ManaCost: core.ManaCostOptions{
			BaseCost: 0.04,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		DamageMultiplier: damageMultiplier,
		CritMultiplier:   shaman.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.OHWeaponDamage(sim, spell.MeleeAttackPower())
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				return
			}
			searingFlames := shaman.SearingFlames.Dot(target)
			if searingFlames.GetStacks() > 0 {
				numberSpread := 0
				maxTargets := 4
				for _, otherTarget := range sim.Encounter.TargetUnits {
					if otherTarget != target {
						shaman.FlameShock.Cast(sim, otherTarget)
						numberSpread++
					}

					if numberSpread >= maxTargets {
						return
					}
				}

				searingFlames.SetStacks(sim, 0)
			}
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return shaman.HasOHWeapon()
		},
	})
}

func (shaman *Shaman) IsLavaLashCastable(sim *core.Simulation) bool {
	return shaman.LavaLash.IsReady(sim)
}
