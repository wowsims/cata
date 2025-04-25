package subtlety

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/rogue"
)

func (subRogue *SubtletyRogue) registerSanguinaryVein() {
	if subRogue.Talents.SanguinaryVein == 0 {
		return
	}

	svBonus := 1 + 0.08*float64(subRogue.Talents.SanguinaryVein)
	hasGlyph := subRogue.HasPrimeGlyph(proto.RoguePrimeGlyph_GlyphOfHemorrhage)
	isApplied := false

	svDebuffArray := subRogue.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Sanguinary Vein Debuff",
			Duration: core.NeverExpires,
			// Action ID Suppressed to not fill debuff log
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				if !isApplied {
					isApplied = true
					subRogue.AttackTables[aura.Unit.UnitIndex].DamageTakenMultiplier *= svBonus
				}
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				if !aura.Unit.HasAuraWithTag(rogue.RogueBleedTag) {
					subRogue.AttackTables[aura.Unit.UnitIndex].DamageTakenMultiplier /= svBonus
					isApplied = false
				}
			},
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				if isApplied && !subRogue.Options.AssumeBleedActive {
					isApplied = false
					subRogue.AttackTables[aura.Unit.UnitIndex].DamageTakenMultiplier /= svBonus
				}
			},
		})
	})

	subRogue.Env.RegisterPreFinalizeEffect(func() {
		if subRogue.Rupture != nil {
			subRogue.Rupture.RelatedAuraArrays = subRogue.Rupture.RelatedAuraArrays.Append(svDebuffArray)
		}
		if subRogue.Hemorrhage != nil && hasGlyph {
			subRogue.Hemorrhage.RelatedAuraArrays = subRogue.Hemorrhage.RelatedAuraArrays.Append(svDebuffArray)
		}
	})

	subRogue.RegisterPrepullAction(-1, func(sim *core.Simulation) {
		if subRogue.Options.AssumeBleedActive {
			for _, target := range subRogue.Env.Encounter.TargetUnits {
				aura := svDebuffArray.Get(target)
				aura.Duration = core.NeverExpires
				aura.Activate(sim)
			}
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

			if spell == subRogue.Rupture {
				aura := svDebuffArray.Get(result.Target)
				dot := spell.Dot(result.Target)
				aura.Duration = dot.BaseTickLength * time.Duration(dot.BaseTickCount)
				aura.Activate(sim)
			} else if spell == subRogue.Hemorrhage && hasGlyph {
				aura := svDebuffArray.Get(result.Target)
				aura.Duration = 24 * time.Second
				aura.Activate(sim)
			}
		},
	})
}
