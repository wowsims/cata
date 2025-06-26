package fire

import (
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
	extraTargets := core.Ternary(hasGlyph, 4, 3)
	numTargets := len(fire.Env.Encounter.TargetUnits) - 1

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

		DamageMultiplier: 1,
		CritMultiplier:   fire.DefaultCritMultiplier(),
		BonusCoefficient: infernoBlastCoefficient,
		ThreatMultiplier: 1,
		BonusCritPercent: 100,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			mainTarget := target
			debuffState := map[int32]core.DotState{}
			dotRefs := []**core.Spell{&fire.Pyroblast.RelatedDotSpell, &fire.Combustion.RelatedDotSpell, &fire.Ignite}
			var igniteRef *core.Dot

			for _, spellRef := range dotRefs {
				dot := (*spellRef).Dot(mainTarget)
				if dot != nil && dot.IsActive() {
					debuffState[dot.ActionID.SpellID] = dot.SaveState(sim)
					if dot.Spell.Matches(mage.MageSpellIgnite) {
						igniteRef = dot
					}
				}
			}

			for range min(extraTargets, numTargets) {
				aoeTarget := fire.Env.NextTargetUnit(target)
				for _, spellRef := range dotRefs {
					dot := (*spellRef).Dot(aoeTarget)
					state, ok := debuffState[dot.ActionID.SpellID]
					if !ok {
						// not stored, was not active
						continue
					}
					if dot.Spell.Matches(mage.MageSpellIgnite) {
						currentTargetHasIgnite := dot.IsActive() // Storing this here so we can do the common steps without overwriting how it was for the checks.
						newDamage := igniteRef.OutstandingDmg()
						currentDamage := dot.OutstandingDmg()
						totalDamage := currentDamage + newDamage

						// Current Target has Ignite so add the Ignite bank from the main target
						// and then follow the default renew logic
						if currentTargetHasIgnite {
							newTickCount := dot.BaseTickCount + 1
							damagePerTick := totalDamage / float64(newTickCount)
							dot.SnapshotBaseDamage = damagePerTick
							dot.Apply(sim)
						} else {
							// Current Target does not have so Ignite bank from the main target
							// and then start a 2 tick (4s) ignite if the Ignite being spread is less than 1 second, otherwise 3 (6s)
							newTickCount := dot.BaseTickCount + core.TernaryInt32(igniteRef.RemainingDuration(sim) < time.Second, 0, 1)
							damagePerTick := totalDamage / float64(newTickCount)
							dot.SnapshotBaseDamage = damagePerTick
							dot.BaseTickCount = newTickCount
							dot.Apply(sim)
						}
						dot.Aura.SetStacks(sim, int32(dot.SnapshotBaseDamage))
						continue
					}
					(*spellRef).Proc(sim, aoeTarget)
					dot.RestoreState(state, sim)
				}
			}

			baseDamage := fire.CalcAndRollDamageRange(sim, infernoBlastScaling, infernoBlastVariance)
			result := spell.CalcAndDealDamage(sim, mainTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			fire.HandleHeatingUp(sim, spell, result)
		},
	})
}
