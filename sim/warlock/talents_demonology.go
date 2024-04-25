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
		// pet.AddStaticMod(core.SpellModConfig{
		// 	ClassMask: WarlockSpellImpFireBolt,
		// 	Kind:      core.SpellMod_CastTime_Flat,
		// 	TimeValue: time.Millisecond * -250 * time.Duration(warlock.Talents.DarkArts),
		// })

		// //TODO: Add/Mult
		// pet.AddStaticMod(core.SpellModConfig{
		// 	ClassMask:  WarlockSpellFelGuardLegionStrike,
		// 	Kind:       core.SpellMod_DamageDone_Pct,
		// 	FloatValue: .05 * float64(warlock.Talents.DarkArts),
		// })

		// pet.AddStaticMod(core.SpellModConfig{
		// 	ClassMask:  WarlockSpellFelHunterShadowBite,
		// 	Kind:       core.SpellMod_DamageDone_Pct,
		// 	FloatValue: .05 * float64(warlock.Talents.DarkArts),
		// })
	}

	//TODO: Mana Feed

	warlock.registerImpendingDoom()
	//TODO: Bane Of Doom Mod

	warlock.registerMoltenCore()

	//TODO: Ancient Grimoire

	// Inferno
	if warlock.Talents.Inferno {
		warlock.AddStaticMod(core.SpellModConfig{
			ClassMask: WarlockSpellImmolate,
			Kind:      core.SpellMod_DotNumberOfTicks_Flat,
			TimeValue: 2,
		})
	}

	warlock.registerDecimation()

	//TODO: Cremation

	if warlock.Talents.DemonicPact {
		warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= 1.1
		warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= 1.1
		warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexArcane] *= 1.1
		warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexNature] *= 1.1
		warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexHoly] *= 1.1
	}

}

func (warlock *Warlock) registerImpendingDoom() {
	if warlock.Talents.ImpendingDoom <= 0 {
		return
	}

	impendingDoomProcChance := 0.05 * float64(warlock.Talents.ImpendingDoom)

	warlock.NightfallProcAura = warlock.RegisterAura(core.Aura{
		Label:    "Impending Doom",
		ActionID: core.ActionID{SpellID: 85107},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		//TODO: Do they need to hit?
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if (spell == warlock.ShadowBolt || spell == warlock.HandOfGuldan || spell == warlock.SoulFire || spell == warlock.Incinerate) &&
				!warlock.Metamorphosis.CD.IsReady(sim) && sim.Proc(impendingDoomProcChance, "Impending Doom") {
				*warlock.Metamorphosis.CD.Timer = core.Timer(time.Duration(*warlock.Metamorphosis.CD.Timer) - time.Second*15)
				warlock.UpdateMajorCooldowns()
			}
		},
	})
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
			//TODO: Is this still valid? Not in either wotlk or cata tooltip.
			warlock.Incinerate.DefaultCast.GCD = time.Duration(float64(warlock.Incinerate.DefaultCast.GCD) * (1 - castReduction))
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.Incinerate.DamageMultiplier /= moltenCoreDamageBonus
			warlock.Incinerate.CastTimeMultiplier += castReduction
			//TODO: Is this still valid? Not in either wotlk or cata tooltip.
			warlock.Incinerate.DefaultCast.GCD = time.Duration(float64(warlock.Incinerate.DefaultCast.GCD) / (1 - castReduction))
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell == warlock.Incinerate {
				aura.RemoveStack(sim)
			}
		},
	})

	warlock.RegisterAura(core.Aura{
		Label: "Molten Core Hidden Aura",
		// ActionID: core.ActionID{SpellID: 47247},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		//TODO: Can this occur on the initial Immolate damage?
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell == warlock.Immolate {
				if sim.Proc(0.02*float64(warlock.Talents.MoltenCore), "Molten Core") {
					warlock.MoltenCoreAura.Activate(sim)
					warlock.MoltenCoreAura.SetStacks(sim, 3)
				}
			}
		},
	})
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
			//TODO: Is this still valid? Not in either wotlk or cata tooltip.
			warlock.SoulFire.DefaultCast.GCD = time.Duration(float64(warlock.SoulFire.DefaultCast.GCD) * (1 - decimationMod))
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.SoulFire.CastTimeMultiplier += decimationMod
			//TODO: Is this still valid? Not in either wotlk or cata tooltip.
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
