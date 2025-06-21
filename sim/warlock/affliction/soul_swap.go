package affliction

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (affliction *AfflictionWarlock) registerSoulSwap() {
	var inhaleTarget *core.Unit
	var debuffState map[int32]core.DotState
	dotRefs := []**core.Spell{&affliction.Corruption, &affliction.Agony, &affliction.Seed, &affliction.UnstableAffliction}

	inhaleBuff := affliction.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 86211},
		Label:    "Soul Swap",
		Duration: time.Second * 3,
	})

	// Exhale
	affliction.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 86213},
		Flags:       core.SpellFlagAPL,
		ProcMask:    core.ProcMaskEmpty,
		SpellSchool: core.SpellSchoolShadow,

		ThreatMultiplier: 1,
		CritMultiplier:   affliction.DefaultCritMultiplier(),
		DamageMultiplier: 1,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return inhaleBuff.IsActive() && target != inhaleTarget
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// restore states
			for _, spellRef := range dotRefs {
				dot := (*spellRef).Dot(target)
				state, ok := debuffState[dot.ActionID.SpellID]
				if !ok {
					// not stored, was not active
					continue
				}

				(*spellRef).Proc(sim, target)
				dot.RestoreState(state, sim)
			}
		},
	})

	// used to not allocate a result for every check
	expectedDamage := &core.SpellResult{}

	// we dont use seed in the expected calculations as it's not applied by exhale
	expectedDotRefs := []**core.Spell{&affliction.Corruption, &affliction.Agony, &affliction.UnstableAffliction}

	// Inhale
	affliction.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 86121},
		Flags:       core.SpellFlagAPL,
		ProcMask:    core.ProcMaskEmpty,
		SpellSchool: core.SpellSchoolShadow,

		ThreatMultiplier: 1,
		CritMultiplier:   affliction.DefaultCritMultiplier(),
		DamageMultiplier: 1,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return (anyDoTActive(dotRefs, target) || affliction.SoulBurnAura.IsActive()) && !inhaleBuff.IsActive()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if affliction.SoulBurnAura.IsActive() {
				affliction.Agony.Proc(sim, target)
				affliction.Corruption.Proc(sim, target)
				affliction.UnstableAffliction.Proc(sim, target)
				affliction.SoulBurnAura.Deactivate(sim)
				return
			}

			inhaleTarget = target
			debuffState = map[int32]core.DotState{}

			// store states
			for _, spellRef := range dotRefs {
				dot := (*spellRef).Dot(target)
				if dot.IsActive() {
					debuffState[dot.ActionID.SpellID] = dot.SaveState(sim)
				}
			}

			inhaleBuff.Activate(sim)
		},

		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			expectedDamage.Damage = 0
			if useSnapshot {
				for _, spellRef := range expectedDotRefs {
					dot := (*spellRef).Dot(target)
					expectedDamage.Damage += dot.Spell.ExpectedTickDamageFromCurrentSnapshot(sim, target)
				}

				return expectedDamage
			}

			for _, spellRef := range expectedDotRefs {
				dot := (*spellRef).Dot(target)
				expectedDamage.Damage += dot.Spell.ExpectedTickDamage(sim, target)
			}

			return expectedDamage
		},
	})
}

func anyDoTActive(dots []**core.Spell, target *core.Unit) bool {
	for _, spellRef := range dots {
		if (*spellRef).Dot(target).IsActive() {
			return true
		}
	}

	return false
}
