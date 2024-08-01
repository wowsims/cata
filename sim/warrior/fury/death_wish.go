package fury

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/warrior"
)

func (war *FuryWarrior) RegisterDeathWish() {
	hasGlyph := war.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfDeathWish)
	var bonusSnapshot float64
	dwAura := war.RegisterAura(core.Aura{
		Label:    "Death Wish",
		ActionID: core.ActionID{SpellID: 12292},
		Tag:      warrior.EnrageTag,
		Duration: 30 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if sim.CurrentTime < 0 && war.Options.PrepullMastery > 0 {
				prepullMultiplier := war.GetMasteryBonusMultiplierFromMasteryRating(float64(war.Options.PrepullMastery))
				bonusSnapshot = 1.0 + (0.2 * prepullMultiplier)
			} else {
				bonusSnapshot = 1.0 + (0.2 * war.EnrageEffectMultiplier)
			}

			war.PseudoStats.DamageDealtMultiplier *= bonusSnapshot
			if !hasGlyph {
				war.PseudoStats.DamageTakenMultiplier *= bonusSnapshot
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			war.PseudoStats.DamageDealtMultiplier /= bonusSnapshot
			if !hasGlyph {
				war.PseudoStats.DamageTakenMultiplier /= bonusSnapshot
			}
		},
	})

	dwSpell := war.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 12292},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL | core.SpellFlagMCD,
		ClassSpellMask: warrior.SpellMaskDeathWish,

		RageCost: core.RageCostOptions{
			Cost: 10,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 0,
			},

			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: 3 * time.Minute,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dwAura.Activate(sim)
		},
	})

	core.RegisterPercentDamageModifierEffect(dwAura, 1.2)

	war.AddMajorCooldown(core.MajorCooldown{
		Spell: dwSpell,
		Type:  core.CooldownTypeDPS,
	})
}
