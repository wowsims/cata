package warlock

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (warlock *Warlock) ApplyDestructionTalents() {
	// Bane
	warlock.AddStaticMod(core.SpellModConfig{
		ClassMask: WarlockSpellShadowBolt | WarlockSpellChaosBolt | WarlockSpellImmolate,
		Kind:      core.SpellMod_CastTime_Flat,
		TimeValue: time.Duration([]int{0, -100, -300, -500}[warlock.Talents.Bane]) * time.Millisecond,
	})

	// Shadow And Flame
	warlock.AddStaticMod(core.SpellModConfig{
		ClassMask:  WarlockSpellShadowBolt | WarlockSpellIncinerate,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: []float64{0.0, 0.04, 0.08, 0.12}[warlock.Talents.ShadowAndFlame],
	})

	// Improved Immolate
	warlock.AddStaticMod(core.SpellModConfig{
		ClassMask:  WarlockSpellImmolate | WarlockSpellImmolateDot | WarlockSpellConflagrate,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: []float64{0.0, 0.1, 0.2}[warlock.Talents.ImprovedImmolate],
	})

	// Emberstorm
	warlock.AddStaticMod(core.SpellModConfig{
		ClassMask: WarlockSpellSoulFire,
		Kind:      core.SpellMod_CastTime_Flat,
		TimeValue: time.Duration(-500*warlock.Talents.Emberstorm) * time.Millisecond,
	})
	warlock.AddStaticMod(core.SpellModConfig{
		ClassMask: WarlockSpellIncinerate,
		Kind:      core.SpellMod_CastTime_Flat,
		TimeValue: time.Duration([]float64{0, -130, -250}[warlock.Talents.Emberstorm]) * time.Millisecond,
	})

	warlock.registerImprovedSearingPain()
	warlock.registerBackdraft()
	warlock.registerShadowBurnSpell()
	warlock.registerBurningEmbers()
	warlock.registerSoulLeech()

	// FireAndBrimstoneDamage mod is in Immolate
	warlock.AddStaticMod(core.SpellModConfig{
		ClassMask:  WarlockSpellConflagrate,
		Kind:       core.SpellMod_BonusCrit_Percent,
		FloatValue: 5.0 * float64(warlock.Talents.FireAndBrimstone),
	})

	warlock.registerEmpoweredImp()

	// TODO: BANE OF HAVOC
}

func (warlock *Warlock) registerImprovedSearingPain() {
	if warlock.Talents.ImprovedSearingPain <= 0 {
		return
	}

	improvedSearingPain := warlock.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusCrit_Percent,
		ClassMask:  WarlockSpellSearingPain,
		FloatValue: 20 * float64(warlock.Talents.ImprovedSearingPain),
	})

	warlock.RegisterResetEffect(func(sim *core.Simulation) {
		improvedSearingPain.Deactivate()
		sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute int32) {
			if isExecute == 25 {
				improvedSearingPain.Activate()
			}
		})
	})
}

func (warlock *Warlock) registerBackdraft() {
	if warlock.Talents.Backdraft <= 0 {
		return
	}

	castTimeMod := warlock.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_CastTime_Pct,
		ClassMask:  WarlockSpellShadowBolt | WarlockSpellIncinerate | WarlockSpellChaosBolt,
		FloatValue: -0.10 * float64(warlock.Talents.Backdraft),
	})

	backdraft := warlock.RegisterAura(core.Aura{
		Label:     "Backdraft",
		ActionID:  core.ActionID{SpellID: 54277},
		Duration:  15 * time.Second,
		MaxStacks: 3,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			castTimeMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			castTimeMod.Deactivate()
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(WarlockSpellShadowBolt|WarlockSpellIncinerate|WarlockSpellChaosBolt) &&
				// DTR procs don't consume backdraft stacks
				!spell.ProcMask.Matches(core.ProcMaskSpellProc) {
				aura.RemoveStack(sim)
			}
		},
	})

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label:    "Backdraft Hidden Aura",
		ActionID: core.ActionID{SpellID: 47260},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(WarlockSpellConflagrate) && result.Landed() {
				backdraft.Activate(sim)
				backdraft.SetStacks(sim, 3)
			}
		},
	}))
}

// burning embers triggers before the spell lands for imp, when it lands for soul fire
// not affected by improved soul fire
func (warlock *Warlock) registerBurningEmbers() {
	if warlock.Talents.BurningEmbers <= 0 {
		return
	}

	ticks := int32(7)
	spMult := (0.7 * float64(warlock.Talents.BurningEmbers)) / float64(ticks)
	baseDmg := warlock.CalcScalingSpellDmg(0.07349999994*float64(warlock.Talents.BurningEmbers)) / float64(ticks)

	warlock.BurningEmbers = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 85421},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage, // TODO: even the imp hits can proc some trinkets, though not most
		Flags:          core.SpellFlagIgnoreModifiers | core.SpellFlagNoSpellMods | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,
		ClassSpellMask: WarlockSpellBurningEmbers,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura:          core.Aura{Label: "Burning Embers"},
			NumberOfTicks: ticks,
			TickLength:    1 * time.Second,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				// damage is capped to a limit that depends on SP; the spell will not "remember" the damage from
				// previous hits even if our SP increases so it's safe to do this here
				dot.SnapshotAttackerMultiplier = 1
				dot.SnapshotBaseDamage = min(baseDmg+spMult*dot.Spell.SpellPower(), dot.SnapshotBaseDamage)
			},

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},
	})

	// we don't do imp's firebolt portion of this here, because it applies _before_ the spell actually hits, and
	// unfortunately the OnCastComplete hook does not provide the result of the cast
	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label:    "Burning Embers Hidden Aura",
		ActionID: core.ActionID{SpellID: 85112},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(WarlockSpellSoulFire) && result.Landed() && !spell.ProcMask.Matches(core.ProcMaskSpellProc|core.ProcMaskSpellDamageProc) {
				dot := warlock.BurningEmbers.Dot(result.Target)
				if !dot.IsActive() {
					dot.SnapshotBaseDamage = 0.0 // ensure we don't use old dot data
				}
				dot.SnapshotBaseDamage += result.Damage * 0.25 * float64(warlock.Talents.BurningEmbers)
				dot.Apply(sim)
			}
		},
	}))
}

func (warlock *Warlock) registerSoulLeech() {
	if warlock.Talents.SoulLeech <= 0 {
		return
	}

	actionID := core.ActionID{SpellID: 30295}
	restore := 0.02 * float64(warlock.Talents.SoulLeech)
	manaMetrics := warlock.NewManaMetrics(actionID)

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label:    "Soul Leech Hidden Aura",
		ActionID: actionID,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(WarlockSpellShadowBurn | WarlockSpellSoulFire | WarlockSpellChaosBolt) {
				warlock.AddMana(sim, restore*warlock.MaxMana(), manaMetrics)
				// also restores health but probably NA
			}
		},
	}))
}

func (warlock *Warlock) registerEmpoweredImp() {
	if warlock.Talents.EmpoweredImp <= 0 || warlock.Options.Summon != proto.WarlockOptions_Imp {
		return
	}

	procChance := 0.02 * float64(warlock.Talents.EmpoweredImp)

	castTimeMod := warlock.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_CastTime_Pct,
		ClassMask:  WarlockSpellSoulFire,
		FloatValue: -1,
	})

	empoweredImpAura := warlock.RegisterAura(core.Aura{
		Label:    "Empowered Imp",
		ActionID: core.ActionID{SpellID: 47221},
		Duration: 8 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			castTimeMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			castTimeMod.Deactivate()
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// if the soul fire cast started BEFORE we got the empowered imp buff, it does not get consumed
			if spell.Matches(WarlockSpellSoulFire) && sim.CurrentTime-spell.CurCast.CastTime > aura.StartedAt() {
				aura.Deactivate(sim)
			}
		},
	})

	core.MakePermanent(warlock.Imp.RegisterAura(core.Aura{
		Label: "Empowered Imp Hidden Aura",

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(WarlockSpellImpFireBolt) && sim.Proc(procChance, "Empowered Imp") {
				empoweredImpAura.Activate(sim)
			}
		},
	}))
}
