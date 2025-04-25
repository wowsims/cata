package warrior

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (warrior *Warrior) RegisterSunderArmor() *core.Spell {
	warrior.SunderArmorAuras = warrior.NewEnemyAuraArray(core.SunderArmorAura)

	hasGlyph := warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfSunderArmor)
	numTargets := warrior.Env.GetNumTargets()
	config := core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 7386},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskSunderArmor | SpellMaskSpecialAttack,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   15,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.CanApplySunderAura(target)
		},

		ThreatMultiplier: 1,
		FlatThreatBonus:  360,

		RelatedAuraArrays: warrior.SunderArmorAuras.ToMap(),
	}

	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
		result.Threat = spell.ThreatFromDamage(result.Outcome, 0.05*spell.MeleeAttackPower())

		if result.Landed() {
			warrior.TryApplySunderArmorEffect(sim, target)
			// https://www.wowhead.com/mop-classic/item=43427/glyph-of-sunder-armor - also applies to devastate in cata
			if hasGlyph && numTargets > 1 {
				nextTarget := warrior.Env.NextTargetUnit(target)
				warrior.TryApplySunderArmorEffect(sim, nextTarget)
			}
		} else {
			spell.IssueRefund(sim)
		}

		spell.DealOutcome(sim, result)

	}
	return warrior.RegisterSpell(config)
}

func (warrior *Warrior) CanApplySunderAura(target *core.Unit) bool {
	return warrior.SunderArmorAuras.Get(target).IsActive() || !warrior.SunderArmorAuras.Get(target).ExclusiveEffects[0].Category.AnyActive()
}

func (warrior *Warrior) TryApplySunderArmorEffect(sim *core.Simulation, target *core.Unit) {
	if warrior.CanApplySunderAura(target) {
		aura := warrior.SunderArmorAuras.Get(target)
		aura.Activate(sim)
		if aura.IsActive() {
			aura.AddStack(sim)
		}
	}
}
