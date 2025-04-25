package rogue

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (rogue *Rogue) registerExposeArmorSpell() {
	rogue.ExposeArmorAuras = rogue.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.ExposeArmorAura(target, rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfExposeArmor))
	})
	durationBonus := core.TernaryDuration(rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfExposeArmor), time.Second*12, 0)
	rogue.exposeArmorDurations = [6]time.Duration{
		0,
		time.Second*6 + durationBonus,
		time.Second*12 + durationBonus,
		time.Second*18 + durationBonus,
		time.Second*24 + durationBonus,
		time.Second*30 + durationBonus,
	}

	rogue.ExposeArmor = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 8647},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | SpellFlagFinisher | core.SpellFlagAPL,
		MetricSplits:   6,
		ClassSpellMask: RogueSpellExposeArmor,

		EnergyCost: core.EnergyCostOptions{
			Cost:          25.0,
			Refund:        0.8,
			RefundMetrics: rogue.EnergyRefundMetrics,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				spell.SetMetricsSplit(spell.Unit.ComboPoints())
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return rogue.ComboPoints() > 0
		},

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				spentPoints := rogue.ComboPoints()
				debuffAura := rogue.ExposeArmorAuras.Get(target)
				debuffAura.Duration = rogue.exposeArmorDurations[spentPoints]
				debuffAura.Activate(sim)
				rogue.ApplyFinisher(sim, spell)
				if rogue.Talents.ImprovedExposeArmor > 0 {
					procChance := 0.5 * float64(rogue.Talents.ImprovedExposeArmor)
					if sim.Proc(procChance, "Improved Expose Armor") {
						rogue.AddComboPoints(sim, spentPoints, spell.ComboPointMetrics())
					}
				}
			} else {
				spell.IssueRefund(sim)
			}
			spell.DealOutcome(sim, result)
		},

		RelatedAuraArrays: rogue.ExposeArmorAuras.ToMap(),
	})
}
