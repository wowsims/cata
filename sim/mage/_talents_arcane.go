package mage

import (
	//"github.com/wowsims/mop/sim/core/proto"

	"time"

	"github.com/wowsims/mop/sim/core"
)

func (mage *Mage) ApplyArcaneTalents() {

	// Cooldowns/Special Implementations
	mage.applyArcaneConcentration()
	mage.registerPresenceOfMindCD()
	mage.applyArcanePotency()
	mage.applyFocusMagic()
	mage.registerArcanePowerCD()

	// Netherwind Presence
	if mage.Talents.NetherwindPresence > 0 {
		mage.PseudoStats.CastSpeedMultiplier *= 1 + (0.01 * float64(mage.Talents.NetherwindPresence))
	}

	// Torment the Weak
	if mage.Talents.TormentTheWeak > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			ClassMask:  MageSpellArcaneBarrage | MageSpellArcaneBlast | MageSpellArcaneExplosion | MageSpellArcaneMissilesTick,
			FloatValue: 0.02 * float64(mage.Talents.TormentTheWeak),
			Kind:       core.SpellMod_DamageDone_Flat,
		})
	}

	//Improved Arcane Missiles
	if mage.Talents.ImprovedArcaneMissiles > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			ClassMask: MageSpellArcaneMissilesCast,
			IntValue:  mage.Talents.ImprovedArcaneMissiles,
			Kind:      core.SpellMod_DotNumberOfTicks_Flat,
		})
	}

	// Arcane Flows
	if mage.Talents.ArcaneFlows > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			ClassMask:  MageSpellArcanePower | MageSpellPresenceOfMind,
			FloatValue: []float64{0, 0.88, 0.75}[mage.Talents.ArcaneFlows],
			Kind:       core.SpellMod_Cooldown_Multiplier,
		})

		mage.AddStaticMod(core.SpellModConfig{
			ClassMask: MageSpellEvocation,
			TimeValue: -time.Minute * time.Duration(mage.Talents.ArcaneFlows),
			Kind:      core.SpellMod_Cooldown_Flat,
		})
	}

	// Missile Barrage
	if mage.Talents.MissileBarrage > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			ClassMask: MageSpellArcaneMissilesCast,
			TimeValue: time.Millisecond * time.Duration(-100*mage.Talents.MissileBarrage),
			Kind:      core.SpellMod_DotTickLength_Flat,
		})
	}

	// Improved Arcane Explosion
	if mage.Talents.ImprovedArcaneExplosion > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			ClassMask: MageSpellArcaneExplosion,
			TimeValue: -1 * time.Duration(0.3*float64(mage.Talents.ImprovedArcaneExplosion)),
			Kind:      core.SpellMod_GlobalCooldown_Flat,
		})

		mage.AddStaticMod(core.SpellModConfig{
			ClassMask: MageSpellArcaneExplosion,
			IntValue:  -25 * mage.Talents.ImprovedArcaneExplosion,
			Kind:      core.SpellMod_PowerCost_Pct,
		})
	}

}

func (mage *Mage) applyArcaneConcentration() {
	if mage.Talents.ArcaneConcentration == 0 {
		return
	}

	// The result that caused the proc. Used to check we don't deactivate from the same proc.
	var procCheckAt time.Duration
	var procSpell *core.Spell
	procChance := []float64{0, 0.13, 0.27, 0.4}[mage.Talents.ArcaneConcentration]

	clearCastingMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask: MageSpellsAllDamaging,
		IntValue:  -1000,
		Kind:      core.SpellMod_PowerCost_Pct,
	})

	// The Clearcasting proc
	clearcastingAura := mage.RegisterAura(core.Aura{
		Label:    "Clearcasting",
		ActionID: core.ActionID{SpellID: 12536},
		Duration: time.Second * 15,
		Icd: &core.Cooldown{
			Timer:    mage.NewTimer(),
			Duration: time.Second * 15,
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if mage.arcanePotencyAura != nil {
				mage.arcanePotencyAura.Activate(sim)
			}
			clearCastingMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			clearCastingMod.Deactivate()
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.ClassSpellMask&MageSpellsAllDamaging == 0 {
				return
			}
			if spell.DefaultCast.Cost == 0 {
				return
			}
			if procCheckAt == sim.CurrentTime && procSpell == spell {
				// Means this is another hit from the same cast that procced CC.
				return
			}
			aura.Deactivate(sim)
		},
	})

	// Tracks if Clearcasting should proc
	core.MakePermanent(mage.RegisterAura(core.Aura{
		Label: "Arcane Concentration",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ClassSpellMask&MageSpellsAllDamaging == 0 {
				return
			}
			if !clearcastingAura.Icd.IsReady(sim) {
				return
			}
			if !result.Landed() {
				return
			}
			// Arcane Concentration does 1 roll for aoe spells as long as one target is hit
			if procCheckAt == sim.CurrentTime {
				return
			}

			procCheckAt = sim.CurrentTime
			procSpell = spell

			curProcChance := procChance
			// Arcane Missile ticks can proc CC, just at a low rate of about 1.5% with 5/5 Arcane Concentration
			if spell.ClassSpellMask == MageSpellArcaneMissilesTick {
				curProcChance *= 0.15
			}

			if !sim.Proc(curProcChance, "Arcane Concentration") {
				return
			}

			clearcastingAura.Icd.Use(sim)
			clearcastingAura.Activate(sim)
		},
	}))
}

func (mage *Mage) registerPresenceOfMindCD() {
	if !mage.Talents.PresenceOfMind {
		return
	}

	presenceOfMindMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  MageSpellsAll ^ MageSpellInstantCast ^ MageSpellEvocation,
		FloatValue: -1,
		Kind:       core.SpellMod_CastTime_Pct,
	})

	var pomSpell *core.Spell

	mage.presenceOfMindAura = mage.RegisterAura(core.Aura{
		Label:    "Presence of Mind",
		ActionID: core.ActionID{SpellID: 12043},
		Duration: time.Hour,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			presenceOfMindMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			presenceOfMindMod.Deactivate()
			pomSpell.CD.Use(sim)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.ClassSpellMask&(MageSpellsAll^MageSpellInstantCast^MageSpellEvocation) == 0 {
				return
			}
			if spell.DefaultCast.CastTime == 0 {
				return
			}
			aura.Deactivate(sim)
		},
	})

	pomSpell = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 12043},
		Flags:          core.SpellFlagNoOnCastComplete,
		ClassSpellMask: MageSpellPresenceOfMind,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second * 120,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return mage.GCD.IsReady(sim) && !mage.arcanePowerAura.IsActive()
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if mage.arcanePotencyAura != nil {
				mage.arcanePotencyAura.Activate(sim)
				mage.arcanePotencyAura.SetStacks(sim, 2)
			}
			mage.presenceOfMindAura.Activate(sim)
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: pomSpell,
		Type:  core.CooldownTypeDPS,
	})
}

func (mage *Mage) applyArcanePotency() {
	if mage.Talents.ArcanePotency == 0 {
		return
	}

	arcanePotencyMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  MageSpellsAllDamaging,
		FloatValue: []float64{0.0, 7, 15}[mage.Talents.ArcanePotency],
		Kind:       core.SpellMod_BonusCrit_Percent,
	})

	var procTime time.Duration
	mage.arcanePotencyAura = mage.RegisterAura(core.Aura{
		Label:     "Arcane Potency",
		ActionID:  core.ActionID{SpellID: 57531},
		Duration:  time.Hour,
		MaxStacks: 2,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			procTime = sim.CurrentTime
			arcanePotencyMod.Activate()
			aura.SetStacks(sim, 2)
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			arcanePotencyMod.Deactivate()
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// Only remove a stack if it's an applicable spell
			if sim.CurrentTime == procTime {
				return
			}

			if spell.ClassSpellMask&((MageSpellsAllDamaging^MageSpellArcaneMissilesTick)|MageSpellArcaneMissilesCast) == 0 {
				return
			}

			aura.RemoveStack(sim)
		},
	})

}

func (mage *Mage) registerArcanePowerCD() {
	if !mage.Talents.ArcanePower {
		return
	}

	actionID := core.ActionID{SpellID: 12042}
	arcanePowerCostMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask: MageSpellsAllDamaging,
		IntValue:  10,
		Kind:      core.SpellMod_PowerCost_Pct,
	})

	arcanePowerDmgMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  MageSpellsAllDamaging,
		FloatValue: 0.2,
		Kind:       core.SpellMod_DamageDone_Pct,
	})

	mage.arcanePowerAura = mage.RegisterAura(core.Aura{
		Label:    "Arcane Power Aura",
		ActionID: actionID,
		Duration: time.Second * 15,
	})

	mage.arcanePowerAura.NewExclusiveEffect("ManaCost", true, core.ExclusiveEffect{
		Priority: 10,
		OnGain: func(ee *core.ExclusiveEffect, sim *core.Simulation) {
			if mage.arcanePowerGCDmod != nil {
				mage.arcanePowerGCDmod.Activate()
			}

			arcanePowerCostMod.UpdateIntValue(core.TernaryInt32(mage.T12_4pc.IsActive(), -10, 20))
			arcanePowerCostMod.Activate()

			arcanePowerDmgMod.Activate()
		},
		OnExpire: func(ee *core.ExclusiveEffect, sim *core.Simulation) {
			if mage.arcanePowerGCDmod != nil {
				mage.arcanePowerGCDmod.Deactivate()
			}

			arcanePowerCostMod.Deactivate()
			arcanePowerDmgMod.Deactivate()
			if mage.t13ProcAura != nil {
				mage.t13ProcAura.Deactivate(sim)
			}
		},
	})

	arcanePower := mage.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagNoOnCastComplete,
		ClassSpellMask: MageSpellArcanePower,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second * 120,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !mage.presenceOfMindAura.IsActive()
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			mage.arcanePowerAura.Activate(sim)
			if mage.T13_4pc.IsActive() {
				// We need to manually set the CD to the correct value because of Arcane Flows talent being applied after the CD
				spell.CD.Set(sim.CurrentTime + time.Duration(float64(spell.CD.Duration-time.Second*time.Duration(7*mage.t13ProcAura.GetStacks()))*spell.CdMultiplier))
			}
		},
	})
	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: arcanePower,
		Type:  core.CooldownTypeDPS,
	})
}
