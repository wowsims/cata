package shaman

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

var StormstrikeActionID = core.ActionID{SpellID: 17364}

// TODO: Confirm how this affects lightning shield
func (shaman *Shaman) StormstrikeDebuffAura(target *core.Unit) *core.Aura {
	return target.GetOrRegisterAura(core.Aura{
		Label:    "Stormstrike-" + shaman.Label,
		ActionID: StormstrikeActionID,
		Duration: time.Second * 15,
	})
}

func (shaman *Shaman) calcDamageStormstrikeCritChance(sim *core.Simulation, target *core.Unit, baseDamage float64, spell *core.Spell) *core.SpellResult {
	var result *core.SpellResult
	if target.HasActiveAura("Stormstrike-" + shaman.Label) {
		critPercentBonus := core.TernaryFloat64(shaman.HasPrimeGlyph(proto.ShamanPrimeGlyph_GlyphOfStormstrike), 35, 25)
		spell.BonusCritPercent += critPercentBonus
		result = spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		spell.BonusCritPercent -= critPercentBonus
	} else {
		result = spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
	}
	return result
}

func (shaman *Shaman) newStormstrikeHitSpell(isMH bool) *core.Spell {
	var procMask core.ProcMask
	var actionTag int32
	if isMH {
		procMask = core.ProcMaskMeleeMHSpecial
		actionTag = 1
	} else {
		procMask = core.ProcMaskMeleeOHSpecial
		actionTag = 2
	}

	return shaman.RegisterSpell(core.SpellConfig{
		ActionID:       StormstrikeActionID.WithTag(actionTag),
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       procMask,
		Flags:          core.SpellFlagMeleeMetrics,
		ClassSpellMask: SpellMaskStormstrikeDamage,

		ThreatMultiplier: 1,
		DamageMultiplier: 2.25,
		CritMultiplier:   shaman.DefaultCritMultiplier(),

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

func (shaman *Shaman) registerStormstrikeSpell() {
	mhHit := shaman.newStormstrikeHitSpell(true)
	ohHit := shaman.newStormstrikeHitSpell(false)

	ssDebuffAuras := shaman.NewEnemyAuraArray(shaman.StormstrikeDebuffAura)

	shaman.Stormstrike = shaman.RegisterSpell(core.SpellConfig{
		ActionID:       StormstrikeActionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskStormstrikeCast,
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Second * 8,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHitNoHitCounter)
			if result.Landed() {
				ssDebuffAura := ssDebuffAuras.Get(target)
				ssDebuffAura.Activate(sim)

				if shaman.HasMHWeapon() {
					mhHit.Cast(sim, target)
				}

				if shaman.AutoAttacks.IsDualWielding && shaman.HasOHWeapon() {
					ohHit.Cast(sim, target)
				}
			}
			spell.DealOutcome(sim, result)
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return shaman.HasMHWeapon() || shaman.HasOHWeapon()
		},
	})
}
