package hunter

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (hunter *Hunter) ApplyBMTalents() {
	if hunter.Talents.ImprovedKillCommand > 0 {
		hunter.Pet.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_BonusCrit_Percent,
			ClassMask:  HunterSpellKillCommand,
			FloatValue: float64(hunter.Talents.ImprovedKillCommand) * 5,
		})
	}
	hunter.applyKillingStreak()
	hunter.applyCobraStrikes()
	hunter.applyInvigoration()
	hunter.applyFocusFireCD()
	hunter.applyFervorCD()
	hunter.applySpiritBond()
}

func (hunter *Hunter) applyCobraStrikes() {
	if hunter.Talents.CobraStrikes == 0 || hunter.Pet == nil {
		return
	}

	actionID := core.ActionID{SpellID: 53260}
	procChance := 0.05 * float64(hunter.Talents.CobraStrikes)

	hunter.Pet.CobraStrikesAura = hunter.Pet.RegisterAura(core.Aura{
		Label:     "Cobra Strikes",
		ActionID:  actionID,
		Duration:  time.Second * 10,
		MaxStacks: 2,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			hunter.Pet.focusDump.BonusCritPercent += 100
			if hunter.Pet.specialAbility != nil {
				hunter.Pet.specialAbility.BonusCritPercent += 100
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			hunter.Pet.focusDump.BonusCritPercent -= 100
			if hunter.Pet.specialAbility != nil {
				hunter.Pet.specialAbility.BonusCritPercent -= 100
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskMeleeSpecial | core.ProcMaskSpellDamage) {
				aura.RemoveStack(sim)
			}
		},
	})

	hunter.RegisterAura(core.Aura{
		Label:    "Cobra Strikes",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell != hunter.ArcaneShot { // Only arcane shot, but also can proc on non crits
				return
			}

			if sim.RandomFloat("Cobra Strikes") < procChance {
				hunter.Pet.CobraStrikesAura.Activate(sim)
				hunter.Pet.CobraStrikesAura.SetStacks(sim, 2)
			}
		},
	})
}

func (hunter *Hunter) applySpiritBond() {
	if hunter.Talents.SpiritBond == 0 || hunter.Pet == nil {
		return
	}

	hunter.PseudoStats.HealingTakenMultiplier *= 1 + 0.05*float64(hunter.Talents.SpiritBond)
	hunter.Pet.PseudoStats.HealingTakenMultiplier *= 1 + 0.05*float64(hunter.Talents.SpiritBond)

	actionID := core.ActionID{SpellID: 20895}
	healthMultiplier := 0.01 * float64(hunter.Talents.SpiritBond)
	healthMetrics := hunter.NewHealthMetrics(actionID)
	petHealthMetrics := hunter.Pet.NewHealthMetrics(actionID)

	hunter.RegisterResetEffect(func(sim *core.Simulation) {
		core.StartPeriodicAction(sim, core.PeriodicActionOptions{
			Period: time.Second * 10,
			OnAction: func(sim *core.Simulation) {
				hunter.GainHealth(sim, hunter.MaxHealth()*healthMultiplier, healthMetrics)
				hunter.Pet.GainHealth(sim, hunter.Pet.MaxHealth()*healthMultiplier, petHealthMetrics)
			},
		})
	})
}

func (hunter *Hunter) applyInvigoration() {
	if hunter.Talents.Invigoration == 0 || hunter.Pet == nil {
		return
	}

	focusMetrics := hunter.NewFocusMetrics(core.ActionID{SpellID: 53253})

	hunter.Pet.RegisterAura(core.Aura{
		Label:    "Invigoration",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMeleeSpecial | core.ProcMaskSpellDamage) {
				return
			}

			if !result.DidCrit() {
				return
			}

			hunter.AddFocus(sim, 3*float64(hunter.Talents.Invigoration), focusMetrics)
		},
	})
}

func (hunter *Hunter) registerBestialWrathCD() {
	if !hunter.Talents.BestialWrath {
		return
	}
	if hunter.Talents.TheBeastWithin {
		hunter.PseudoStats.DamageDealtMultiplier *= 1.1
	}
	bwCostMod := hunter.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Pct,
		ClassMask: HunterSpellsAll,
		IntValue:  -50,
	})
	actionID := core.ActionID{SpellID: 19574}

	bestialWrathPetAura := hunter.Pet.RegisterAura(core.Aura{
		Label:    "Bestial Wrath Pet",
		ActionID: actionID,
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.2
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.2
		},
	})

	bestialWrathAura := hunter.RegisterAura(core.Aura{
		Label:    "Bestial Wrath",
		ActionID: actionID,
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.2
			if hunter.Talents.TheBeastWithin {
				bwCostMod.Activate()
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.2
			if hunter.Talents.TheBeastWithin {
				bwCostMod.Deactivate()
			}
		},
	})
	core.RegisterPercentDamageModifierEffect(bestialWrathAura, 1.2)

	bwSpell := hunter.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: HunterSpellBestialWrath,
		FocusCost: core.FocusCostOptions{
			Cost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 1,
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: hunter.applyLongevity(time.Minute * 2),
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			bestialWrathPetAura.Activate(sim)

			if hunter.Talents.TheBeastWithin {
				bestialWrathAura.Activate(sim)
			}
		},
	})

	hunter.AddMajorCooldown(core.MajorCooldown{
		Spell: bwSpell,
		Type:  core.CooldownTypeDPS,
	})
}

func (hunter *Hunter) applyFervorCD() {
	if !hunter.Talents.Fervor {
		return
	}

	actionID := core.ActionID{SpellID: 82726}
	focusMetrics := hunter.NewFocusMetrics(actionID)
	fervorSpell := hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		FocusCost: core.FocusCostOptions{
			Cost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second * 1,
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Minute * 2,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			hunter.AddFocus(sim, 50, focusMetrics)
			hunter.Pet.AddFocus(sim, 50, focusMetrics)
		},
	})

	hunter.AddMajorCooldown(core.MajorCooldown{
		Spell: fervorSpell,
		Type:  core.CooldownTypeDPS,
	})
}
func (hunter *Hunter) applyFocusFireCD() {
	if !hunter.Talents.FocusFire || hunter.Pet == nil {
		return
	}

	actionID := core.ActionID{SpellID: 82692}
	petFocusMetrics := hunter.Pet.NewFocusMetrics(actionID)
	focusFireAura := hunter.RegisterAura(core.Aura{
		Label:    "Focus Fire",
		ActionID: actionID,
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			hunter.Pet.FrenzyStacksSnapshot = float64(hunter.Pet.FrenzyAura.GetStacks())
			if hunter.Pet.FrenzyStacksSnapshot >= 1 {
				hunter.Pet.FrenzyAura.Deactivate(sim)
				hunter.Pet.AddFocus(sim, 4, petFocusMetrics)
				aura.Unit.MultiplyRangedSpeed(sim, 1+(float64(hunter.Pet.FrenzyStacksSnapshot)*0.03))
				if sim.Log != nil {
					hunter.Pet.Log(sim, "Consumed %0f stacks of Frenzy for Focus Fire.", hunter.Pet.FrenzyStacksSnapshot)
				}
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if hunter.Pet.FrenzyStacksSnapshot > 0 {
				aura.Unit.MultiplyRangedSpeed(sim, 1/(1+(float64(hunter.Pet.FrenzyStacksSnapshot)*0.03)))
			}
		},
	})

	focusFireSpell := hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		FocusCost: core.FocusCostOptions{
			Cost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 1,
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 15,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if focusFireAura.IsActive() {
				focusFireAura.Deactivate(sim) // Want to apply new one
			}
			focusFireAura.Activate(sim)
			//focusFireAura.OnGain(focusFireAura, sim)
		},
	})

	hunter.AddMajorCooldown(core.MajorCooldown{
		Spell: focusFireSpell,
		Type:  core.CooldownTypeDPS,
	})

}
func (hunter *Hunter) applyLongevity(dur time.Duration) time.Duration {
	return time.Duration(float64(dur) * (1.0 - 0.1*float64(hunter.Talents.Longevity)))
}
func (hunter *Hunter) applyFrenzy() {
	if hunter.Talents.Frenzy == 0 {
		return
	}
	actionID := core.ActionID{SpellID: 19622}
	hunter.Pet.FrenzyAura = hunter.Pet.RegisterAura(core.Aura{
		Label:     "Frenzy",
		Duration:  time.Second * 10,
		ActionID:  actionID,
		MaxStacks: 5,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyMeleeSpeed(sim, 1+(float64(hunter.Talents.Frenzy)*0.02))
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyMeleeSpeed(sim, 1/(1+float64(hunter.Talents.Frenzy)*0.02))
		},
	})

	hunter.Pet.RegisterAura(core.Aura{
		Label:    "FrenzyHandler",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMeleeSpecial | core.ProcMaskSpellDamage) {
				return
			}
			if hunter.Pet.FrenzyAura.IsActive() {
				if hunter.Pet.FrenzyAura.GetStacks() != 5 {
					hunter.Pet.FrenzyAura.AddStack(sim)
					hunter.Pet.FrenzyAura.Refresh(sim)
				}
			} else {
				hunter.Pet.FrenzyAura.Activate(sim)
				hunter.Pet.FrenzyAura.SetStacks(sim, 1)
			}
		},
	})
}
func (hunter *Hunter) applyKillingStreak() {
	if hunter.Talents.KillingStreak == 0 {
		return
	}
	damageMod := hunter.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  HunterSpellKillCommand,
		FloatValue: float64(hunter.Talents.KillingStreak) * 0.1,
	})
	costMod := hunter.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Flat,
		ClassMask: HunterSpellKillCommand,
		IntValue:  -hunter.Talents.KillingStreak * 5,
	})
	if hunter.Pet != nil {
		hunter.KillingStreakAura = hunter.Pet.RegisterAura(core.Aura{
			Label:    "Killing Streak",
			ActionID: core.ActionID{SpellID: 82748},
			Duration: core.NeverExpires,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				damageMod.Activate()
				costMod.Activate()
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				damageMod.Deactivate()
				costMod.Deactivate()
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell == hunter.Pet.KillCommand {
					aura.Deactivate(sim)
				}
			},
		})
		hunter.KillingStreakCounterAura = hunter.Pet.RegisterAura(core.Aura{
			Label:     "Killing Streak (KC Crit)",
			Duration:  core.NeverExpires,
			MaxStacks: 2,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell == hunter.Pet.KillCommand {
					if aura.GetStacks() == 2 && result.DidCrit() {
						hunter.KillingStreakAura.Activate(sim)
						aura.SetStacks(sim, 1)
						return
					}
					if result.DidCrit() {
						aura.AddStack(sim)
					} else {
						aura.SetStacks(sim, 1)
					}
				}
			},
		})
	}
}
