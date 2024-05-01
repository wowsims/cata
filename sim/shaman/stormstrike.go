package shaman

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

var StormstrikeActionID = core.ActionID{SpellID: 17364}

func Tier12StormstrikeBonus(sim *core.Simulation, spell *core.Spell, attackTable *core.AttackTable) float64 {
	if spell.ClassSpellMask&(SpellMaskFireNova|SpellMaskFlameShock|SpellMaskLavaBurst|SpellMaskUnleashFlame|SpellMaskFlametongueWeapon) > 0 {
		return 1.06
	}
	return 1.0
}

// TODO: Confirm how this affects lightning shield
func (shaman *Shaman) StormstrikeDebuffAura(target *core.Unit) *core.Aura {
	hasT12P4 := false // todo
	return target.GetOrRegisterAura(core.Aura{
		Label:    "Stormstrike-" + shaman.Label,
		ActionID: StormstrikeActionID,
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if hasT12P4 {
				core.EnableDamageDoneByCaster(DDBC_T12P2, DDBC_Total, shaman.AttackTables[aura.Unit.UnitIndex], Tier12StormstrikeBonus)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if hasT12P4 {
				core.DisableDamageDoneByCaster(DDBC_T12P2, shaman.AttackTables[aura.Unit.UnitIndex])
			}
		},
	})
}

func (shaman *Shaman) calcDamageStormstrikeCritChance(sim *core.Simulation, target *core.Unit, baseDamage float64, spell *core.Spell) *core.SpellResult {
	var result *core.SpellResult
	if target.HasActiveAura("Stormstrike-" + shaman.Label) {
		critBonus := core.CritRatingPerCritChance * core.TernaryFloat64(shaman.HasPrimeGlyph(proto.ShamanPrimeGlyph_GlyphOfStormstrike), 35, 25)
		spell.BonusCritRating += critBonus
		result = spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		spell.BonusCritRating -= critBonus
	} else {
		result = spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
	}
	return result
}

func (shaman *Shaman) newStormstrikeHitSpell(isMH bool) func(*core.Simulation, *core.Unit, *core.Spell) {
	var procMask core.ProcMask
	if isMH {
		procMask = core.ProcMaskMeleeMHSpecial
	} else {
		procMask = core.ProcMaskMeleeOHSpecial
	}

	return func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		var baseDamage float64
		spell.ProcMask = procMask
		if isMH {
			baseDamage = spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
		} else {
			baseDamage = spell.Unit.OHWeaponDamage(sim, spell.MeleeAttackPower())
		}

		spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
	}
}

func (shaman *Shaman) registerStormstrikeSpell() {
	mhHit := shaman.newStormstrikeHitSpell(true)
	ohHit := shaman.newStormstrikeHitSpell(false)

	ssDebuffAuras := shaman.NewEnemyAuraArray(shaman.StormstrikeDebuffAura)

	shaman.Stormstrike = shaman.RegisterSpell(core.SpellConfig{
		ActionID:       StormstrikeActionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagIncludeTargetBonusDamage,
		ClassSpellMask: SpellMaskStormstrike,
		ManaCost: core.ManaCostOptions{
			BaseCost: 0.08,
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

		ThreatMultiplier: 1,
		DamageMultiplier: 2.25,
		CritMultiplier:   shaman.DefaultMeleeCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				ssDebuffAura := ssDebuffAuras.Get(target)
				ssDebuffAura.Activate(sim)

				if shaman.HasMHWeapon() {
					mhHit(sim, target, spell)
				}

				if shaman.AutoAttacks.IsDualWielding && shaman.HasOHWeapon() {
					ohHit(sim, target, spell)
				}

				shaman.Stormstrike.SpellMetrics[target.UnitIndex].Hits--
			}
			spell.DealOutcome(sim, result)
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return shaman.HasMHWeapon() || shaman.HasOHWeapon()
		},
	})
}
