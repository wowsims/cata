package fire

import (
	"fmt"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/mage"
)

func (fire *FireMage) registerInfernoBlastSpell() {

	infernoBlastVariance := 0.17   // Per https://wago.tools/db2/SpellEffect?build=5.5.0.61217&filter%5BSpellID%5D=108853 Field: "Variance"
	infernoBlastScaling := .60     // Per https://wago.tools/db2/SpellEffect?build=5.5.0.61217&filter%5BSpellID%5D=108853 Field: "Coefficient"
	infernoBlastCoefficient := .60 // Per https://wago.tools/db2/SpellEffect?build=5.5.0.61217&filter%5BSpellID%5D=108853 Field: "BonusCoefficient"

	hasGlyph := fire.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfInfernoBlast)
	extraTargets := core.Ternary(hasGlyph, 4, 5)

	fire.InfernoBlast = fire.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 108853},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: mage.MageSpellInfernoBlast,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 2,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    fire.NewTimer(),
				Duration: time.Second * 8,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           fire.DefaultCritMultiplier(),
		BonusCoefficient:         infernoBlastCoefficient,
		ThreatMultiplier:         1,
		BonusCritPercent:         100,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			debuffState := map[int32]core.DotState{}
			dotRefs := []**core.Spell{&fire.Pyroblast.RelatedDotSpell, &fire.Combustion.RelatedDotSpell, &fire.Ignite}

			for _, spellRef := range dotRefs {
				fmt.Println(*spellRef)
				dot := (*spellRef).Dot(target)
				if dot.IsActive() {
					fmt.Println(dot.SaveState(sim))
					debuffState[dot.ActionID.SpellID] = dot.SaveState(sim)
				}
			}

			currTarget := target
			for range min(extraTargets, len(sim.Encounter.TargetUnits)-1) {
				currTarget = fire.Env.NextTargetUnit(currTarget)
				for _, spellRef := range dotRefs {
					dot := (*spellRef).Dot(currTarget)
					state, ok := debuffState[dot.ActionID.SpellID]
					if !ok {
						// not stored, was not active
						continue
					}

					(*spellRef).Proc(sim, currTarget)
					dot.RestoreState(state, sim)
				}
			}

			baseDamage := fire.CalcAndRollDamageRange(sim, infernoBlastScaling, infernoBlastVariance)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}
