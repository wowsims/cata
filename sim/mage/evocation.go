package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (mage *Mage) registerEvocation() {

	if mage.Talents.RuneOfPower {
		return
	}

	actionID := core.ActionID{SpellID: 12051}
	manaMetrics := mage.NewManaMetrics(actionID)
	manaPerTick := 0.0

	evocation := mage.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagHelpful | core.SpellFlagChanneled | core.SpellFlagAPL,
		ClassSpellMask: MageSpellEvocation,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Minute * 2,
			},
		},

		Hot: core.DotConfig{
			SelfOnly: true,
			Aura: core.Aura{
				Label: "Evocation",
			},
			NumberOfTicks:        3,
			TickLength:           time.Second * 2,
			AffectedByCastSpeed:  true,
			HasteReducesDuration: true,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				mage.AddMana(sim, manaPerTick, manaMetrics)
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			manaPerTick = mage.MaxMana() * 0.15
			spell.SelfHot().Apply(sim)
			spell.SelfHot().TickOnce(sim)
		},
	})

	invocationCooldownMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  MageSpellEvocation,
		FloatValue: -1.0,
		Kind:       core.SpellMod_Cooldown_Multiplier,
	})

	invocationSpeedUp := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  MageSpellEvocation,
		FloatValue: -1.0,
		Kind:       core.SpellMod_DotTickLength_Flat,
	})

	invocationDamageMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  MageSpellsAllDamaging,
		FloatValue: 0.15,
		Kind:       core.SpellMod_DamageDone_Pct,
	})

	mage.invocationAura = mage.RegisterAura(core.Aura{
		Label:    "Invocation Aura",
		ActionID: actionID,
		Duration: time.Minute,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			invocationDamageMod.Activate()
			aura.Activate(sim)
		},
	})

	if mage.Talents.Invocation {
		invocationCooldownMod.Activate()
		invocationSpeedUp.Activate()
	}

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: evocation,
		Type:  core.CooldownTypeMana,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			if sim.GetRemainingDuration() < 12*time.Second {
				return false
			}

			return character.CurrentManaPercent() < 0.1
		},
	})
}
