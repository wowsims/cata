package paladin

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

func (paladin *Paladin) registerDivineProtectionSpell() {
	glyphOfDivineProtection := paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfDivineProtection)

	actionID := core.ActionID{SpellID: 498}
	paladin.DivineProtectionAura = paladin.RegisterAura(core.Aura{
		Label:    "Divine Protection" + paladin.Label,
		ActionID: actionID,
		Duration: time.Second * 10,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if glyphOfDivineProtection {
				paladin.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] *= 0.6
				paladin.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] *= 0.6
				paladin.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] *= 0.6
				paladin.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexHoly] *= 0.6
				paladin.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] *= 0.6
				paladin.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] *= 0.6
			} else {
				paladin.PseudoStats.DamageTakenMultiplier *= 0.8
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if glyphOfDivineProtection {
				paladin.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] /= 0.6
				paladin.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] /= 0.6
				paladin.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] /= 0.6
				paladin.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexHoly] /= 0.6
				paladin.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] /= 0.6
				paladin.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] /= 0.6
			} else {
				paladin.PseudoStats.DamageTakenMultiplier /= 0.8
			}
		},
	})

	paladin.DivineProtection = paladin.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskDivineProtection,

		ManaCost: core.ManaCostOptions{
			BaseCostFraction: 0.03,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Minute * 1,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			paladin.DivineProtectionAura.Activate(sim)
		},
	})

	if paladin.Spec == proto.Spec_SpecProtectionPaladin {
		paladin.AddMajorCooldown(core.MajorCooldown{
			Spell: paladin.DivineProtection,
			Type:  core.CooldownTypeSurvival,
		})
	}
}
