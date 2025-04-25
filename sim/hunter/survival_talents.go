package hunter

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

func (hunter *Hunter) ApplySurvivalTalents() {
	if hunter.Talents.Pathing > 0 {
		bonus := 0.01 * float64(hunter.Talents.Pathing)
		core.MakePermanent(hunter.RegisterAura(core.Aura{
			BuildPhase: core.CharacterBuildPhaseBase,
			Label:      "Pathing",
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				hunter.PseudoStats.RangedSpeedMultiplier /= 1 + bonus
			},
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				hunter.PseudoStats.RangedSpeedMultiplier *= 1 + bonus
			},
		}))
	}

	if hunter.Talents.HunterVsWild > 0 {
		bonus := 0.05 * float64(hunter.Talents.HunterVsWild)
		hunter.MultiplyStat(stats.Stamina, 1+bonus)
	}

	if hunter.Talents.HuntingParty {
		agiBonus := 0.02
		hunter.MultiplyStat(stats.Agility, 1.0+agiBonus)
	}

	if hunter.Talents.Resourcefulness > 0 {
		hunter.AddStaticMod(core.SpellModConfig{
			Kind:      core.SpellMod_Cooldown_Flat,
			ClassMask: HunterSpellBlackArrow | HunterSpellExplosiveTrap,
			TimeValue: -(time.Second * 2 * time.Duration(hunter.Talents.Resourcefulness)),
		})
	}
	if hunter.Talents.TrapMastery > 0 {
		hunter.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Flat,
			ClassMask:  HunterSpellBlackArrow | HunterSpellExplosiveTrap,
			FloatValue: .10 * float64(hunter.Talents.TrapMastery),
		})
	}
	if hunter.Talents.Toxicology > 0 {
		hunter.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_CritMultiplier_Flat,
			ClassMask:  HunterSpellBlackArrow | HunterSpellSerpentSting,
			FloatValue: float64(hunter.Talents.Toxicology) * 0.5,
		})
	}
	if hunter.Talents.ImprovedSerpentSting > 0 {
		hunter.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_BonusCrit_Percent,
			ClassMask:  HunterSpellSerpentSting,
			FloatValue: float64(hunter.Talents.ImprovedSerpentSting) * 5,
		})
	}

	hunter.applyTNT()
	hunter.applySniperTraining()
	hunter.applyThrillOfTheHunt()
}

func (hunter *Hunter) applyThrillOfTheHunt() {
	if hunter.Talents.ThrillOfTheHunt == 0 {
		return
	}

	procChance := float64(hunter.Talents.ThrillOfTheHunt) * 0.05
	focusMetrics := hunter.NewFocusMetrics(core.ActionID{SpellID: 34499})

	hunter.RegisterAura(core.Aura{
		Label:    "Thrill of the Hunt",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// mask 256
			if spell == hunter.ArcaneShot || spell.ClassSpellMask == HunterSpellExplosiveShot || spell == hunter.BlackArrow {
				if sim.Proc(procChance, "ThrillOfTheHunt") {
					hunter.AddFocus(sim, spell.DefaultCast.Cost*0.4, focusMetrics)
				}
			}
		},
	})
}

func (hunter *Hunter) applySniperTraining() {
	if hunter.Talents.SniperTraining == 0 {
		return
	}

	uptime := hunter.SurvivalOptions.SniperTrainingUptime
	if uptime <= 0 {
		return
	}
	uptime = min(1, uptime)

	dmgMod := hunter.AddDynamicMod(core.SpellModConfig{
		ClassMask:  HunterSpellCobraShot | HunterSpellSteadyShot,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: .02 * float64(hunter.Talents.SniperTraining),
	})

	stAura := hunter.RegisterAura(core.Aura{
		Label:    "Sniper Training",
		ActionID: core.ActionID{SpellID: 53304},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dmgMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dmgMod.Deactivate()
		},
	})

	core.ApplyFixedUptimeAura(stAura, uptime, time.Second*15, 1)

	hunter.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusCrit_Percent,
		ClassMask:  HunterSpellKillShot,
		FloatValue: 5 * float64(hunter.Talents.SniperTraining),
	})
}

// Todo: Should we support precasting freezing/ice trap?
func (hunter *Hunter) applyTNT() {
	if hunter.Talents.TNT == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 56343}
	procChance := []float64{0, 0.10, 0.20}[hunter.Talents.TNT]

	icd := core.Cooldown{
		Timer:    hunter.NewTimer(),
		Duration: time.Second * 10,
	}

	hunter.LockAndLoadAura = hunter.RegisterAura(core.Aura{
		Icd:       &icd,
		Label:     "Lock and Load Proc",
		ActionID:  actionID,
		Duration:  time.Second * 12,
		MaxStacks: 2,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if hunter.ExplosiveShot != nil {
				hunter.ExplosiveShot.Cost.PercentModifier -= 100
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if hunter.ExplosiveShot != nil {
				hunter.ExplosiveShot.Cost.PercentModifier += 100
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell == hunter.ExplosiveShot {
				hunter.ExplosiveShot.CD.Reset()
				// Weird check but..
				if !aura.Unit.HasActiveAura("Burning Adrenaline") {
					aura.RemoveStack(sim)
				}
			}
		},
		// OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
		// 	if spell == hunter.ExplosiveShot {

		// 		hunter.ExplosiveShot.CD.Reset()

		// 		aura.RemoveStack(sim)
		// 	}
		// },
	})

	hunter.RegisterAura(core.Aura{
		Label:    "TNT Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell != hunter.BlackArrow && spell != hunter.ExplosiveTrap {
				return
			}

			if !icd.IsReady(sim) {
				return
			}

			if sim.RandomFloat("Lock and Load") < procChance {
				icd.Use(sim)
				hunter.LockAndLoadAura.Activate(sim)
				hunter.LockAndLoadAura.SetStacks(sim, 2)
				if hunter.ExplosiveShot != nil {
					hunter.ExplosiveShot.CD.Reset()
				}
			}
		},
	})
}
