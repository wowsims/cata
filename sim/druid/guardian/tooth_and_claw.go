package guardian

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/druid"
)

func (bear *GuardianDruid) registerToothAndClawPassive() {
	// First register the stackable debuff on enemy units.
	debuffConfig := core.Aura{
		Label:     "Tooth and Claw",
		ActionID:  core.ActionID{SpellID: 135597},
		Duration:  time.Second * 15,
		MaxStacks: math.MaxInt32,

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskWhiteHit) && result.Landed() {
				aura.Deactivate(sim)
			}
		},
	}

	bear.ToothAndClawDebuffs = bear.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(debuffConfig)
	})

	// Then register the absorb effect on friendly units.
	absorbHandler := func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult, _ bool) {
		if !spell.ProcMask.Matches(core.ProcMaskWhiteHit) || (result.Damage <= 0) {
			return
		}

		debuff := bear.ToothAndClawDebuffs.Get(spell.Unit)

		if !debuff.IsActive() {
			return
		}

		absorbedDamage := min(float64(debuff.GetStacks()), result.Damage)
		result.Damage -= absorbedDamage

		if sim.Log != nil {
			result.Target.Log(sim, "Tooth and Claw absorbed %.1f damage from incoming auto-attack.", absorbedDamage)
		}
	}

	for _, unit := range bear.Env.AllUnits {
		if unit.Type != core.EnemyUnit {
			unit.AddDynamicDamageTakenModifier(absorbHandler)
		}
	}

	// Next, register the personal buff that empowers Maul.
	bear.ToothAndClawBuff = bear.RegisterAura(core.Aura{
		Label:    "Tooth and Claw",
		ActionID: core.ActionID{SpellID: 135286},
		Duration: time.Second * 10,

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if bear.Maul.IsEqual(spell) && result.Landed() {
				debuff := bear.ToothAndClawDebuffs.Get(result.Target)
				debuff.Activate(sim)
				addedAbsorbAmount := max((bear.GetStat(stats.AttackPower)-2*bear.GetStat(stats.Agility))*2.2, bear.GetStat(stats.Stamina)*2.5) * 0.4
				debuff.SetStacks(sim, debuff.GetStacks()+int32(math.Round(addedAbsorbAmount)))
				aura.Deactivate(sim)
			}
		},
	})

	// Finally, register the trigger for the personal buff.
	core.MakeProcTriggerAura(&bear.Unit, core.ProcTrigger{
		Name:       "Tooth and Claw Trigger",
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskWhiteHit,
		Outcome:    core.OutcomeLanded,
		Harmful:    true,
		ProcChance: 0.4,

		Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			if bear.InForm(druid.Bear) {
				bear.ToothAndClawBuff.Activate(sim)
			}
		},
	})
}
