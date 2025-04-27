package brewmaster

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/monk"
)

/*
Tooltip:
You Guard against future attacks, absorbing [(Attack power * 1.971) + 16202] damage for 30 sec.

Any heals you apply to yourself while Guarding are increased by 30%.

-- Glyph of Guard --
Increases the amount your Guard absorbs by 10%, but your Guard can only absorb magical damage.
-- Glyph of Guard --
*/
func (bm *BrewmasterMonk) registerGuard() {
	actionID := core.ActionID{SpellID: 115295}

	spellSchool := core.SpellSchoolPhysical | core.SpellSchoolArcane | core.SpellSchoolFire | core.SpellSchoolFrost | core.SpellSchoolHoly | core.SpellSchoolNature | core.SpellSchoolShadow

	if bm.HasMajorGlyph(proto.MonkMajorGlyph_GlyphOfGuard) {
		spellSchool ^= core.SpellSchoolPhysical
	}

	aura := bm.NewDamageAbsorptionAuraForSchool(
		"Guard Absorb",
		actionID.WithTag(1),
		30*time.Second,
		spellSchool,
		func(_ *core.Unit) float64 {
			return bm.GetStat(stats.AttackPower)*1.971 + bm.CalcScalingSpellDmg(13)
		},
	)

	aura.Aura.AttachMultiplicativePseudoStatBuff(&bm.PseudoStats.HealingTakenMultiplier, 1.3)

	bm.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          monk.SpellFlagSpender | core.SpellFlagAPL,
		ClassSpellMask: monk.MonkSpellGuard,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    bm.NewTimer(),
				Duration: time.Second * 30,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return bm.ComboPoints() >= 2
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			aura.Activate(sim)
		},
	})
}
