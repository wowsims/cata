package rogue

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (rogue *Rogue) registerExposeArmorSpell() {
	rogue.ExposeArmorAuras = rogue.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.WeakenedArmorAura(target)
	})

	cpMetric := rogue.NewComboPointMetrics(core.ActionID{SpellID: 8647})
	hasGlyph := rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfExposeArmor)

	rogue.ExposeArmor = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 8647},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | SpellFlagBuilder | core.SpellFlagAPL,
		ClassSpellMask: RogueSpellExposeArmor,

		EnergyCost: core.EnergyCostOptions{
			Cost:   25.0,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
				// Omitting the GCDMin - does not appear affected by either Shadow Blades or Adrenaline Rush
			},
			IgnoreHaste: true,
		},

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				debuffAura := rogue.ExposeArmorAuras.Get(target)
				debuffAura.Activate(sim)
				if hasGlyph {
					// just set the stacks to 3
					debuffAura.SetStacks(sim, 3)
				} else {
					debuffAura.AddStack(sim)
				}

				rogue.AddComboPointsOrAnticipation(sim, 1, cpMetric)
			} else {
				spell.IssueRefund(sim)
			}
			spell.DealOutcome(sim, result)
		},

		RelatedAuraArrays: rogue.ExposeArmorAuras.ToMap(),
	})
}
