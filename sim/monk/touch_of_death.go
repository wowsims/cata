package monk

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (monk *Monk) registerTouchOfDeath() {
	hasGlyph := monk.HasMajorGlyph(proto.MonkMajorGlyph_GlyphOfTouchOfDeath)
	actionID := core.ActionID{SpellID: 115080}
	chiMetrics := monk.NewChiMetrics(actionID)
	cooldown := time.Second*90 + core.TernaryDuration(hasGlyph, 2*time.Minute, 0)

	monk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagCannotBeDodged | core.SpellFlagIgnoreArmor | core.SpellFlagIgnoreModifiers | SpellFlagBuilder | core.SpellFlagAPL,
		ClassSpellMask: MonkSpellTouchOfDeath,
		MaxRange:       core.MaxMeleeRange,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    monk.NewTimer(),
				Duration: cooldown,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   monk.DefaultCritMultiplier(),

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return (hasGlyph || monk.GetChi() >= 3) && sim.GetRemainingDuration() <= time.Second*1
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.RollWithLabel(0, spell.Unit.MaxHealth(), "Touch Of Death Damage")
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParryNoCrit)

			if result.Landed() {
				spell.DealDamage(sim, result)
				if !hasGlyph {
					monk.SpendChi(sim, 3, chiMetrics)
				}
			}
		},
	})
}
