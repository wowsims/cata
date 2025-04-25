package warlock

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

func (warlock *Warlock) ApplyDemonologyTalents() {
	// Demonic Embrace
	if warlock.Talents.DemonicEmbrace > 0 {
		warlock.MultiplyStat(stats.Stamina, []float64{1.0, 1.04, 1.07, 1.10}[warlock.Talents.DemonicEmbrace])
	}

	// Dark Arts
	if warlock.Talents.DarkArts > 0 {
		warlock.Imp.AddStaticMod(core.SpellModConfig{
			ClassMask: WarlockSpellImpFireBolt,
			Kind:      core.SpellMod_CastTime_Flat,
			TimeValue: time.Duration(-250*warlock.Talents.DarkArts) * time.Millisecond,
		})

		warlock.Felguard.AddStaticMod(core.SpellModConfig{
			ClassMask:  WarlockSpellFelGuardLegionStrike,
			Kind:       core.SpellMod_DamageDone_Flat,
			FloatValue: .05 * float64(warlock.Talents.DarkArts),
		})

		warlock.Felhunter.AddStaticMod(core.SpellModConfig{
			ClassMask:  WarlockSpellFelHunterShadowBite,
			Kind:       core.SpellMod_DamageDone_Flat,
			FloatValue: .05 * float64(warlock.Talents.DarkArts),
		})
	}

	warlock.registerManaFeed()
	warlock.registerMasterSummoner()
	warlock.registerImpendingDoom()
	warlock.registerMoltenCore()

	// Inferno
	if warlock.Talents.Inferno {
		warlock.AddStaticMod(core.SpellModConfig{
			ClassMask: WarlockSpellImmolateDot,
			Kind:      core.SpellMod_DotNumberOfTicks_Flat,
			IntValue:  2,
		})
	}

	warlock.registerDecimation()
	warlock.registerCremation()

	if warlock.Talents.DemonicPact {
		warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= 1.02
		warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= 1.02
	}
}

func (warlock *Warlock) registerManaFeed() {
	if warlock.Talents.ManaFeed <= 0 && !warlock.Talents.Metamorphosis {
		return
	}

	actionID := core.ActionID{SpellID: 85175}
	manaMetrics := warlock.NewManaMetrics(actionID)
	manaReturn := 0.02 * float64(warlock.Talents.ManaFeed)

	aura := core.Aura{
		Label:    "Mana Feed",
		ActionID: actionID,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidCrit() {
				if spell.Matches(WarlockSpellSuccubusLashOfPain | WarlockSpellImpFireBolt) {
					warlock.AddMana(sim, manaReturn*warlock.MaxMana(), manaMetrics)
				} else if spell.Matches(WarlockSpellFelGuardLegionStrike | WarlockSpellFelHunterShadowBite) {
					// felguard and felhunter gain 4x the mana
					warlock.AddMana(sim, 4*manaReturn*warlock.MaxMana(), manaMetrics)
				}
			}
		},
	}

	for _, pet := range warlock.Pets {
		if !pet.IsGuardian() {
			core.MakePermanent(pet.RegisterAura(aura))
		}
	}
}

func (warlock *Warlock) registerMasterSummoner() {
	if warlock.Talents.MasterSummoner <= 0 {
		return
	}

	warlock.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_CastTime_Flat,
		ClassMask:  WarlockSummonSpells,
		FloatValue: float64(-500*warlock.Talents.MasterSummoner) * float64(time.Millisecond),
	})

	warlock.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Pct,
		ClassMask: WarlockSummonSpells,
		IntValue:  -50 * warlock.Talents.MasterSummoner,
	})
}

func (warlock *Warlock) registerImpendingDoom() {
	if warlock.Talents.ImpendingDoom <= 0 {
		return
	}

	impendingDoomProcChance := 0.05 * float64(warlock.Talents.ImpendingDoom)

	if !warlock.Talents.Metamorphosis {
		return
	}

	core.MakePermanent(
		warlock.RegisterAura(core.Aura{
			Label:    "Impending Doom",
			ActionID: core.ActionID{SpellID: 85107},
			//TODO: Do they need to hit?
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if spell.Matches(WarlockSpellShadowBolt|WarlockSpellHandOfGuldan|WarlockSpellSoulFire|WarlockSpellIncinerate) && !spell.ProcMask.Matches(core.ProcMaskSpellProc|core.ProcMaskSpellDamageProc) {
					if !warlock.Metamorphosis.CD.IsReady(sim) && sim.Proc(impendingDoomProcChance, "Impending Doom") {
						warlock.Metamorphosis.CD.Reduce(15 * time.Second)
						warlock.UpdateMajorCooldowns()
					}
				}
			},
		}))
}

func (warlock *Warlock) registerMoltenCore() {
	if warlock.Talents.MoltenCore <= 0 {
		return
	}

	damageMod := warlock.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		ClassMask:  WarlockSpellIncinerate,
		FloatValue: 0.06 * float64(warlock.Talents.MoltenCore),
	})

	castTimeMod := warlock.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_CastTime_Pct,
		ClassMask:  WarlockSpellIncinerate,
		FloatValue: -0.1 * float64(warlock.Talents.MoltenCore),
	})

	moltenCoreAura := warlock.RegisterAura(core.Aura{
		Label:     "Molten Core Proc Aura",
		ActionID:  core.ActionID{SpellID: 71165},
		Duration:  15 * time.Second,
		MaxStacks: 3,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Activate()
			castTimeMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Deactivate()
			castTimeMod.Deactivate()
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// if the incinerate cast started BEFORE we got the molten core buff, the incinerate benefits from it but
			// does not consume a stack. Detect this and only remove a stack if that's not the case
			if spell.Matches(WarlockSpellIncinerate) && sim.CurrentTime-spell.CurCast.CastTime > aura.StartedAt() {
				aura.RemoveStack(sim)
			}
		},
	})

	procChance := 0.02 * float64(warlock.Talents.MoltenCore)
	onHit := func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
		if spell.Matches(WarlockSpellImmolate|WarlockSpellImmolateDot) && !spell.ProcMask.Matches(core.ProcMaskSpellProc|core.ProcMaskSpellDamageProc) && sim.Proc(procChance, "Molten Core") {
			moltenCoreAura.Activate(sim)
			moltenCoreAura.SetStacks(sim, 3)
		}
	}

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label:                 "Molten Core Hidden Aura",
		OnSpellHitDealt:       onHit,
		OnPeriodicDamageDealt: onHit,
	}))
}

func (warlock *Warlock) registerDecimation() {
	if warlock.Talents.Decimation <= 0 {
		return
	}

	decimationMod := warlock.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_CastTime_Pct,
		ClassMask:  WarlockSpellSoulFire,
		FloatValue: -0.2 * float64(warlock.Talents.Decimation),
	})

	decimationAura := warlock.RegisterAura(core.Aura{
		Label:    "Decimation Proc Aura",
		ActionID: core.ActionID{SpellID: 63167},
		Duration: 10 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			decimationMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			decimationMod.Deactivate()
		},
	})

	decimation := warlock.RegisterAura(core.Aura{
		Label:    "Decimation Talent Hidden Aura",
		Duration: core.NeverExpires,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && spell.Matches(WarlockSpellShadowBolt|WarlockSpellIncinerate|WarlockSpellSoulFire) {
				decimationAura.Activate(sim)
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

	procChance := []float64{0.0, 0.5, 1.0}[warlock.Talents.Cremation]

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label: "Cremation Talent",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(WarlockSpellHandOfGuldan) {
				if warlock.Immolate.Dot(result.Target).IsActive() && sim.Proc(procChance, "Cremation") {
					warlock.Immolate.Dot(result.Target).Apply(sim)
				}
			}
		},
	}))
}
