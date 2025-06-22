package blood

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/death_knight"
)

/*
-- Glyph of Vampiric Blood --

Increases the amount of health the Death Knight receives from healing spells and effects by 40% for 10 sec

-- else --

Temporarily grants the Death Knight 15% of maximum health and increases the amount of health received from healing spells and effects by 25% for 10 sec.
After the effect expires, the health is lost.

----------
*/
func (bdk *BloodDeathKnight) registerVampiricBlood() {
	actionID := core.ActionID{SpellID: 55233}
	healthMetrics := bdk.NewHealthMetrics(actionID)

	hasGlyph := bdk.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfVampiricBlood)
	healBonus := core.TernaryFloat64(hasGlyph, 1.40, 1.25)

	vampiricBloodAura := bdk.RegisterAura(core.Aura{
		Label:    "Vampiric Blood" + bdk.Label,
		ActionID: actionID,
		Duration: time.Second * 10,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if !hasGlyph {
				bdk.VampiricBloodBonusHealth = bdk.MaxHealth() * 0.15
				bdk.UpdateMaxHealth(sim, bdk.VampiricBloodBonusHealth, healthMetrics)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if !hasGlyph {
				bdk.UpdateMaxHealth(sim, -bdk.VampiricBloodBonusHealth, healthMetrics)
			}
		},
	}).AttachMultiplicativePseudoStatBuff(&bdk.PseudoStats.HealingTakenMultiplier, healBonus)

	spell := bdk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		Flags:          core.SpellFlagAPL | core.SpellFlagHelpful | core.SpellFlagReadinessTrinket,
		ClassSpellMask: death_knight.DeathKnightSpellVampiricBlood,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    bdk.NewTimer(),
				Duration: time.Minute,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.RelatedSelfBuff.Activate(sim)
		},

		RelatedSelfBuff: vampiricBloodAura,
	})

	bdk.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeSurvival,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return bdk.CurrentHealthPercent() < 0.4
		},
	})
}
