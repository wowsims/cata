package enhancement

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/shaman"
)

var StormstrikeActionID = core.ActionID{SpellID: 17364}

// TODO: Confirm how this affects lightning shield
func (enh *EnhancementShaman) StormstrikeDebuffAura(target *core.Unit) *core.Aura {
	return target.GetOrRegisterAura(core.Aura{
		Label:    "Stormstrike-" + enh.Label,
		ActionID: StormstrikeActionID,
		Duration: time.Second * 15,
	})
}

func (enh *EnhancementShaman) newStormstrikeHitSpell(isMH bool) *core.Spell {
	var procMask core.ProcMask
	var actionTag int32
	if isMH {
		procMask = core.ProcMaskMeleeMHSpecial
		actionTag = 1
	} else {
		procMask = core.ProcMaskMeleeOHSpecial
		actionTag = 2
	}

	return enh.RegisterSpell(core.SpellConfig{
		ActionID:       StormstrikeActionID.WithTag(actionTag),
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       procMask,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,
		ClassSpellMask: shaman.SpellMaskStormstrikeDamage,

		ThreatMultiplier: 1,
		DamageMultiplier: 2.25,
		CritMultiplier:   enh.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			var baseDamage float64
			if isMH {
				baseDamage = spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
			} else {
				baseDamage = spell.Unit.OHWeaponDamage(sim, spell.MeleeAttackPower())
			}
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialBlockAndCrit)
			spell.SpellMetrics[target.UnitIndex].Casts--
		},
	})
}

func (enh *EnhancementShaman) registerStormstrikeSpell() {
	mhHit := enh.newStormstrikeHitSpell(true)
	ohHit := enh.newStormstrikeHitSpell(false)

	ssDebuffAuras := enh.NewEnemyAuraArray(enh.StormstrikeDebuffAura)

	enh.Stormstrike = enh.RegisterSpell(core.SpellConfig{
		ActionID:       StormstrikeActionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: shaman.SpellMaskStormstrikeCast,
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 9.4,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    enh.NewTimer(),
				Duration: time.Second * 8,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHitNoHitCounter)
			if result.Landed() {
				ssDebuffAura := ssDebuffAuras.Get(target)
				ssDebuffAura.Activate(sim)

				if enh.HasMHWeapon() {
					mhHit.Cast(sim, target)
				}

				if enh.AutoAttacks.IsDualWielding && enh.HasOHWeapon() {
					ohHit.Cast(sim, target)
				}
			}
			spell.DealOutcome(sim, result)
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return enh.HasMHWeapon() || enh.HasOHWeapon()
		},
	})
}
