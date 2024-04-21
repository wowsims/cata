package shaman

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

var StormstrikeActionID = core.ActionID{SpellID: 17364}

// TODO: Confirm how this affects lightning shield
func (shaman *Shaman) StormstrikeDebuffAura(target *core.Unit) *core.Aura {
	return target.GetOrRegisterAura(core.Aura{
		Label:    "Stormstrike-" + shaman.Label,
		ActionID: StormstrikeActionID,
		Duration: time.Second * 15,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Unit != &shaman.Unit {
				return
			}

			if spell == shaman.LightningBolt || spell == shaman.ChainLightning || spell == shaman.LightningShield || spell == shaman.EarthShock {
				spell.CritMultiplier += core.TernaryFloat64(shaman.HasPrimeGlyph(proto.ShamanPrimeGlyph_GlyphOfStormstrike), 0.35, 0.25)
			}
		},
	})
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
		ActionID:    StormstrikeActionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagIncludeTargetBonusDamage,

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
