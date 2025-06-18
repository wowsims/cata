package subtlety

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/rogue"
)

func (subRogue *SubtletyRogue) registerSanguinaryVein() {
	svBonus := 1.35
	hasHemoGlyph := subRogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfHemorraghingVeins)

	svDebuffArray := subRogue.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Sanguinary Vein Debuff",
			Duration: core.NeverExpires,
			// Action ID Suppressed to not fill debuff log
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				core.EnableDamageDoneByCaster(DDBC_SanguinaryVein, DDBC_Total, subRogue.AttackTables[aura.Unit.UnitIndex], func(sim *core.Simulation, spell *core.Spell, attackTable *core.AttackTable) float64 {
					if spell.Matches(rogue.RogueSpellsAll) || spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
						return svBonus
					}
					return 1.0
				})
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				core.DisableDamageDoneByCaster(DDBC_SanguinaryVein, subRogue.AttackTables[aura.Unit.UnitIndex])
			},
		})
	})

	subRogue.Env.RegisterPreFinalizeEffect(func() {
		if subRogue.Rupture != nil {
			subRogue.Rupture.RelatedAuraArrays = subRogue.Rupture.RelatedAuraArrays.Append(svDebuffArray)
		}
		if subRogue.Garrote != nil {
			subRogue.Garrote.RelatedAuraArrays = subRogue.Garrote.RelatedAuraArrays.Append(svDebuffArray)
		}
		if subRogue.Hemorrhage != nil && hasHemoGlyph {
			subRogue.Hemorrhage.RelatedAuraArrays = subRogue.Hemorrhage.RelatedAuraArrays.Append(svDebuffArray)
		}
		if subRogue.CrimsonTempest != nil {
			subRogue.CrimsonTempestDoT.RelatedAuraArrays = subRogue.CrimsonTempestDoT.RelatedAuraArrays.Append(svDebuffArray)
		}
	})

	subRogue.RegisterAura(core.Aura{
		Label:    "Sanguinary Vein Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			if spell == subRogue.Rupture || spell == subRogue.Garrote || spell == subRogue.CrimsonTempestDoT {
				aura := svDebuffArray.Get(result.Target)
				dot := spell.Dot(result.Target)
				aura.Duration = dot.BaseTickLength * time.Duration(dot.BaseTickCount)
				aura.Activate(sim)
			} else if spell == subRogue.Hemorrhage && hasHemoGlyph {
				aura := svDebuffArray.Get(result.Target)
				aura.Duration = 24 * time.Second
				aura.Activate(sim)
			}
		},
	})
}
