package mage

import (
	//"github.com/wowsims/cata/sim/core/proto"

	"time"

	"github.com/wowsims/cata/sim/core"
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
		mage.AddStaticMod(core.SpellModConfig{
			ClassMask:  MageSpellArcaneBarrage,
			FloatValue: -0.01 * float64(mage.Talents.NetherwindPresence) * core.HasteRatingPerHastePercent,
			Kind:       core.SpellMod_CastTime_Pct,
		})
	}

	// Torment the Weak
	if mage.Talents.TormentTheWeak > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			ClassMask:  MageSpellArcaneBarrage | MageSpellArcaneBlast | MageSpellArcaneExplosion, //| MageSpellArcaneMissiles,
			FloatValue: 0.02 * float64(mage.Talents.TormentTheWeak),
			Kind:       core.SpellMod_DamageDone_Flat,
		})
	}

	//Improved Arcane Missiles
	if mage.Talents.ImprovedArcaneMissiles > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			ClassMask: MageSpellArcaneMissilesCast,
			IntValue:  int64(mage.Talents.ImprovedArcaneMissiles),
			Kind:      core.SpellMod_DotNumberOfTicks_Flat,
		})
	}

	// Arcane Flows
	if mage.Talents.ArcaneFlows > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			ClassMask:  MageSpellArcanePower | MageSpellPresenceOfMind,
			FloatValue: -[]float64{0, 0.12, 0.25}[mage.Talents.ArcaneFlows],
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
			ClassMask:  MageSpellArcaneExplosion,
			FloatValue: -0.25 * float64(mage.Talents.ImprovedArcaneExplosion),
			Kind:       core.SpellMod_PowerCost_Pct,
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
	procChance := []float64{0, 0.03, 0.06, 0.1}[mage.Talents.ArcaneConcentration]

	clearCastingMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  MageSpellsAllDamaging,
		FloatValue: -1,
		Kind:       core.SpellMod_PowerCost_Pct,
	})

	// The Clearcasting proc
	clearcastingAura := mage.RegisterAura(core.Aura{
		Label:    "Clearcasting",
		ActionID: core.ActionID{SpellID: 12536},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if mage.ArcanePotencyAura != nil {
				mage.ArcanePotencyAura.Activate(sim)
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

			clearcastingAura.Activate(sim)
		},
	}))
}

func (mage *Mage) registerPresenceOfMindCD() {
	if !mage.Talents.PresenceOfMind {
		return
	}

	presenceOfMindMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  MageSpellsAll ^ MageSpellInstantCast,
		FloatValue: -1,
		Kind:       core.SpellMod_CastTime_Pct,
	})

	var pomSpell *core.Spell

	presenceOfMindAura := mage.RegisterAura(core.Aura{
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
			if spell.ClassSpellMask&(MageSpellsAll^MageSpellInstantCast) == 0 {
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
			return mage.GCD.IsReady(sim)
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if mage.ArcanePotencyAura != nil {
				mage.ArcanePotencyAura.Activate(sim)
				mage.ArcanePotencyAura.SetStacks(sim, 2)
			}
			presenceOfMindAura.Activate(sim)
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
		FloatValue: []float64{0.0, 7, 15}[mage.Talents.ArcanePotency] * core.CritRatingPerCritChance,
		Kind:       core.SpellMod_BonusCrit_Rating,
	})

	var procTime time.Duration
	mage.ArcanePotencyAura = mage.RegisterAura(core.Aura{
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
		ClassMask:  MageSpellsAll,
		FloatValue: 0.1,
		Kind:       core.SpellMod_PowerCost_Pct,
	})

	arcanePowerDmgMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  MageSpellsAll,
		FloatValue: 0.2,
		Kind:       core.SpellMod_DamageDone_Pct,
	})

	arcanePowerAura := mage.RegisterAura(core.Aura{
		Label:    "Arcane Power Aura",
		ActionID: actionID,
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if mage.arcanePowerGCDmod != nil {
				mage.arcanePowerGCDmod.Activate()
			}
			arcanePowerCostMod.Activate()
			arcanePowerDmgMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if mage.arcanePowerGCDmod != nil {
				mage.arcanePowerGCDmod.Deactivate()
			}
			arcanePowerCostMod.Deactivate()
			arcanePowerDmgMod.Deactivate()
		},
	})

	arcanePower := mage.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second * 120,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			arcanePowerAura.Activate(sim)
		},
		/*ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !mage.ArcanePotencyAura.IsActive()
		}, */
	})
	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: arcanePower,
		Type:  core.CooldownTypeDPS,
	})
}
