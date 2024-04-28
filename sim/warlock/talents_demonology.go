package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

func (warlock *Warlock) ApplyDemonologyTalents() {
	// Demonic Embrace
	if warlock.Talents.DemonicEmbrace > 0 {
		warlock.MultiplyStat(stats.Stamina, []float64{1.04, 1.07, 1.10}[warlock.Talents.DemonicEmbrace])
	}

	//Dark Arts
	if warlock.Talents.DarkArts > 0 {
		warlock.Imp.AddStaticMod(core.SpellModConfig{
			ClassMask: WarlockSpellImpFireBolt,
			Kind:      core.SpellMod_CastTime_Flat,
			TimeValue: time.Millisecond * time.Duration(-250*warlock.Talents.DarkArts),
		})

		//TODO: Add/Mult
		warlock.Felguard.AddStaticMod(core.SpellModConfig{
			ClassMask:  WarlockSpellFelGuardLegionStrike,
			Kind:       core.SpellMod_DamageDone_Pct,
			FloatValue: .05 * float64(warlock.Talents.DarkArts),
		})

		//TODO: Add/Mult
		warlock.Felhunter.AddStaticMod(core.SpellModConfig{
			ClassMask:  WarlockSpellFelHunterShadowBite,
			Kind:       core.SpellMod_DamageDone_Pct,
			FloatValue: .05 * float64(warlock.Talents.DarkArts),
		})
	}

	warlock.registerManaFeed()
	warlock.registerImpendingDoom()
	warlock.registerMoltenCore()

	// Inferno
	if warlock.Talents.Inferno {
		warlock.AddStaticMod(core.SpellModConfig{
			ClassMask: WarlockSpellImmolate,
			Kind:      core.SpellMod_DotNumberOfTicks_Flat,
			TimeValue: 2,
		})
	}

	warlock.registerDecimation()
	warlock.registerCremation()

	if warlock.Talents.DemonicPact {
		warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= 1.1
		warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= 1.1
		warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexArcane] *= 1.1
		warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexNature] *= 1.1
		warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexHoly] *= 1.1
	}
}

func (warlock *Warlock) registerManaFeed() {
	if warlock.Talents.ManaFeed <= 0 {
		return
	}

	actionID := core.ActionID{SpellID: 85175}
	manaMetrics := warlock.NewManaMetrics(actionID)
	manaReturn := 0.02 * float64(warlock.Talents.ManaFeed)

	core.MakePermanent(
		warlock.RegisterAura(core.Aura{
			Label:    "Mana Feed",
			ActionID: actionID,
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.DidCrit() && (spell.ClassSpellMask == WarlockSpellImpFireBolt ||
					spell.ClassSpellMask == WarlockSpellFelGuardLegionStrike ||
					spell.ClassSpellMask == WarlockSpellSuccubusLashOfPain ||
					spell.ClassSpellMask == WarlockSpellFelHunterShadowBite) {
					restore := manaReturn * warlock.GetStat(stats.Mana)
					warlock.AddMana(sim, restore, manaMetrics)
				}
			},
		}))
}

func (warlock *Warlock) registerImpendingDoom() {
	if warlock.Talents.ImpendingDoom <= 0 {
		return
	}

	impendingDoomProcChance := 0.05 * float64(warlock.Talents.ImpendingDoom)

	core.MakePermanent(
		warlock.RegisterAura(core.Aura{
			Label:    "Impending Doom",
			ActionID: core.ActionID{SpellID: 85107},
			//TODO: Do they need to hit?
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if (spell == warlock.ShadowBolt || spell == warlock.HandOfGuldan || spell == warlock.SoulFire || spell == warlock.Incinerate) &&
					!warlock.Metamorphosis.CD.IsReady(sim) && sim.Proc(impendingDoomProcChance, "Impending Doom") {
					*warlock.Metamorphosis.CD.Timer = core.Timer(time.Duration(*warlock.Metamorphosis.CD.Timer) - time.Second*15)
					warlock.UpdateMajorCooldowns()
				}
			},
		}))
}

func (warlock *Warlock) registerMoltenCore() {
	if warlock.Talents.MoltenCore <= 0 {
		return
	}

	castReduction := 0.06 * float64(warlock.Talents.MoltenCore)
	moltenCoreDamageBonus := 1 + 0.06*float64(warlock.Talents.MoltenCore)

	warlock.MoltenCoreAura = warlock.RegisterAura(core.Aura{
		Label:     "Molten Core Proc Aura",
		ActionID:  core.ActionID{SpellID: 71165},
		Duration:  time.Second * 15,
		MaxStacks: 3,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warlock.Incinerate.DamageMultiplier *= moltenCoreDamageBonus
			warlock.Incinerate.CastTimeMultiplier -= castReduction
			warlock.Incinerate.DefaultCast.GCD = time.Duration(float64(warlock.Incinerate.DefaultCast.GCD) * (1 - castReduction))
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.Incinerate.DamageMultiplier /= moltenCoreDamageBonus
			warlock.Incinerate.CastTimeMultiplier += castReduction
			warlock.Incinerate.DefaultCast.GCD = time.Duration(float64(warlock.Incinerate.DefaultCast.GCD) / (1 - castReduction))
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell == warlock.Incinerate {
				aura.RemoveStack(sim)
			}
		},
	})

	core.MakePermanent(
		warlock.RegisterAura(core.Aura{
			Label: "Molten Core Hidden Aura",
			//TODO: Can this occur on the initial Immolate damage?
			OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell == warlock.Immolate {
					if sim.Proc(0.02*float64(warlock.Talents.MoltenCore), "Molten Core") {
						warlock.MoltenCoreAura.Activate(sim)
						warlock.MoltenCoreAura.SetStacks(sim, 3)
					}
				}
			},
		}))
}

func (warlock *Warlock) registerDecimation() {
	if warlock.Talents.Decimation <= 0 {
		return
	}

	decimationMod := 0.2 * float64(warlock.Talents.Decimation)
	warlock.DecimationAura = warlock.RegisterAura(core.Aura{
		Label:    "Decimation Proc Aura",
		ActionID: core.ActionID{SpellID: 63167},
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warlock.SoulFire.CastTimeMultiplier -= decimationMod
			warlock.SoulFire.DefaultCast.GCD = time.Duration(float64(warlock.SoulFire.DefaultCast.GCD) * (1 - decimationMod))
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.SoulFire.CastTimeMultiplier += decimationMod
			warlock.SoulFire.DefaultCast.GCD = time.Duration(float64(warlock.SoulFire.DefaultCast.GCD) / (1 - decimationMod))
		},
	})

	decimation := warlock.RegisterAura(core.Aura{
		Label:    "Decimation Talent Hidden Aura",
		Duration: core.NeverExpires,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && (spell == warlock.ShadowBolt || spell == warlock.Incinerate || spell == warlock.SoulFire) {
				warlock.DecimationAura.Activate(sim)
			}
		},
	})

	warlock.RegisterResetEffect(func(sim *core.Simulation) {
		sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute int32) {
			if isExecute == 25 {
				decimation.Activate(sim)
			}
		})
	})
}

func (warlock *Warlock) registerCremation() {
	if warlock.Talents.Cremation <= 0 {
		return
	}

	procChance := []float64{0.5, 1.0}[warlock.Talents.Cremation]

	cremationAura := warlock.RegisterAura(core.Aura{
		Label:    "Cremation Talent",
		Duration: core.NeverExpires,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell == warlock.HandOfGuldan {
				if warlock.ImmolateDot.Dot(aura.Unit).IsActive() && sim.Proc(procChance, "Cremation") {
					//TODO: Should this Rollover or Apply like other dots in cata?
					warlock.ImmolateDot.Dot(aura.Unit).Rollover(sim)
				}
			}
		},
	})

	warlock.RegisterResetEffect(func(sim *core.Simulation) {
		sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute int32) {
			if isExecute == 25 {
				cremationAura.Activate(sim)
			}
		})
	})
}
