package monk

import (
	"fmt"
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (monk *Monk) registerHealingSphere() {
	healingSphereActionID := core.ActionID{SpellID: 115460}
	healingSphereHealActionID := core.ActionID{SpellID: 115464}

	stacksAura := monk.RegisterAura(core.Aura{
		Label:     "Healing Sphere Stacks" + monk.Label,
		ActionID:  healingSphereActionID.WithTag(1),
		Duration:  core.NeverExpires,
		MaxStacks: 3,
	})

	healingSpheres := make([]*core.Aura, stacksAura.MaxStacks)

	for i := range healingSpheres {
		healingSpheres[i] = monk.RegisterAura(core.Aura{
			Label:    fmt.Sprintf("Healing Sphere #%v %v", i, monk.Label),
			ActionID: healingSphereActionID,
			Duration: time.Minute * 1,
		})
	}

	addHealingSphere := func(sim *core.Simulation) {
		for _, healingSphere := range healingSpheres {
			if !healingSphere.IsActive() {
				healingSphere.Activate(sim)
				break
			}
		}
	}

	removeHealingSphere := func(sim *core.Simulation) {
		for _, healingSphere := range healingSpheres {
			if healingSphere.IsActive() {
				stacksAura.RemoveStack(sim)
				healingSphere.Deactivate(sim)
				break
			}
		}
	}

	// Healing Sphere - Heal
	monk.RegisterSpell(core.SpellConfig{
		ActionID:    healingSphereHealActionID,
		SpellSchool: core.SpellSchoolNature,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagHelpful | core.SpellFlagAPL,
		ProcMask:    core.ProcMaskSpellHealing,

		DamageMultiplier: 1,

		BonusCoefficient: 0.5025,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return stacksAura.GetStacks() > 0
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			heal := monk.CalcScalingSpellDmg(9.122) + spell.MeleeAttackPower()*spell.BonusCoefficient
			spell.CalcAndDealHealing(sim, spell.Unit, heal, spell.OutcomeHealing)
			removeHealingSphere(sim)
		},
	})

	// Healing Sphere - Use
	monk.RegisterSpell(core.SpellConfig{
		ActionID:       healingSphereActionID,
		ClassSpellMask: MonkSpellHealingSphere,
		SpellSchool:    core.SpellSchoolNature,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,
		ProcMask:       core.ProcMaskSpellHealing,
		MaxRange:       40,

		DamageMultiplier: 1,
		CritMultiplier:   1,

		EnergyCost: core.EnergyCostOptions{
			Cost: core.TernaryInt32(monk.StanceMatches(WiseSerpent), 0, 40),
		},
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: core.TernaryInt32(monk.StanceMatches(WiseSerpent), 2, 0),
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 500 * time.Microsecond,
			},
			CD: core.Cooldown{
				Timer:    monk.NewTimer(),
				Duration: 500 * time.Millisecond,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return monk.StanceMatches(FierceTiger|SturdyOx|WiseSerpent) && stacksAura.GetStacks() <= stacksAura.MaxStacks
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			stacksAura.Activate(sim)
			stacksAura.AddStack(sim)

			// Delay the Healing Sphere so it can't be used immediately
			startAt := sim.RandomFloat("Healing Sphere Drop") * 200.0
			sim.AddPendingAction(&core.PendingAction{
				NextActionAt: time.Millisecond * time.Duration(startAt),
				OnAction: func(sim *core.Simulation) {
					addHealingSphere(sim)
				},
			})
		},
	})
}
