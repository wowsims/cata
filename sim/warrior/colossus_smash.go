package warrior

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (war *Warrior) registerColossusSmash() {
	if war.Spec == proto.Spec_SpecProtectionWarrior {
		return
	}

	actionID := core.ActionID{SpellID: 86346}
	hasGlyph := war.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfColossusSmash)

	war.ColossusSmashAuras = war.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Colossus Smash",
			ActionID: actionID,
			Duration: time.Millisecond * 6500,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				war.AttackTables[aura.Unit.UnitIndex].IgnoreArmor = true
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				war.AttackTables[aura.Unit.UnitIndex].IgnoreArmor = false
			},
		})
	})

	physVulnerabilityAuras := war.NewEnemyAuraArray(core.PhysVulnerabilityAura)

	war.ColossusSmash = war.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagAPL | core.SpellFlagMeleeMetrics,
		ClassSpellMask: SpellMaskColossusSmash,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   20,
			Refund: 0.8,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: time.Second * 20,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1.75,
		CritMultiplier:   war.DefaultCritMultiplier(),

		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) * 1.77999997139
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			if !result.Landed() {
				spell.IssueRefund(sim)
			} else {
				csAura := war.ColossusSmashAuras.Get(target)
				csAura.Activate(sim)

				physVulnAura := physVulnerabilityAuras.Get(target)
				physVulnAura.Activate(sim)

				if hasGlyph {
					war.TryApplySunderArmorEffect(sim, target)
				}
			}
		},

		RelatedAuraArrays: war.ColossusSmashAuras.ToMap(),
	})
}
