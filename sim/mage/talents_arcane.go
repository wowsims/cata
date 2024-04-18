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

	// Arcane Flows is inside each relevant spell

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
	var proccedAt time.Duration
	var proccedSpell *core.Spell

	// Tracks if Clearcasting should proc
	mage.RegisterAura(core.Aura{
		Label:    "Arcane Concentration",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Flags.Matches(SpellFlagMage) || spell == mage.ArcaneMissiles {
				return
			}
			if !result.Landed() {
				return
			}

			procChance := []float64{0, 0.03, 0.06, 0.1}[mage.Talents.ArcaneConcentration]
			// Arcane Missile ticks can proc CC, just at a low rate of about 1.5% with 5/5 Arcane Concentration
			if spell == mage.ArcaneMissilesTickSpell {
				procChance *= 0.15
			}
			if !sim.Proc(procChance, "Arcane Concentration") {
				return
			}
			proccedAt = sim.CurrentTime
			proccedSpell = spell

			//mage.ArcanePotencyAura.Activate(sim)
			mage.ClearcastingAura.Activate(sim)
			mage.ArcaneBlastAura.GetStacks()
		},
	})

	/* 	if mage.Talents.ArcanePotency > 0 {
		mage.ArcanePotencyAura = mage.RegisterAura(core.Aura{
			Label:    "Arcane Potency",
			ActionID: core.ActionID{SpellID: 31572},
			Duration: time.Hour,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.AddStatDynamic(sim, stats.SpellCrit, bonusCrit)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.AddStatDynamic(sim, stats.SpellCrit, -bonusCrit)
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if !spell.Flags.Matches(SpellFlagMage) {
					return
				}
				// Don't spend on the spell that procced it
				if proccedAt == sim.CurrentTime {
					return
				}
				aura.Deactivate(sim)
			},
		})
	} */

	clearCastingMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  MageSpellsAllDamaging,
		FloatValue: -1,
		Kind:       core.SpellMod_PowerCost_Pct,
	})
	// The Clearcasting proc
	mage.ClearcastingAura = mage.RegisterAura(core.Aura{
		Label:    "Clearcasting",
		ActionID: core.ActionID{SpellID: 12536},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			//mage.ArcanePotencyAura.Activate(sim)
			clearCastingMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			clearCastingMod.Deactivate()
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !spell.Flags.Matches(SpellFlagMage) {
				return
			}
			if spell.DefaultCast.Cost == 0 {
				return
			}
			if spell == mage.ArcaneMissiles && mage.ArcaneMissilesProcAura.IsActive() {
				return
			}
			if proccedAt == sim.CurrentTime && proccedSpell == spell {
				// Means this is another hit from the same cast that procced CC.
				return
			}
			aura.Deactivate(sim)
		},
	})
}

func (mage *Mage) registerPresenceOfMindCD() {
	if !mage.Talents.PresenceOfMind {
		return
	}

	actionID := core.ActionID{SpellID: 12043}
	var spellToUse *core.Spell
	mage.Env.RegisterPostFinalizeEffect(func() {
		spellToUse = mage.ArcaneBlast
	})

	spell := mage.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second * time.Duration(120*(1-[]float64{0.0, 0.07, 0.15}[mage.Talents.ArcaneFlows])),
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			if !mage.GCD.IsReady(sim) {
				return false
			}
			if mage.ArcanePowerAura.IsActive() {
				return false
			}

			manaCost := spellToUse.DefaultCast.Cost * mage.PseudoStats.CostMultiplier
			if spellToUse == mage.ArcaneBlast {
				manaCost *= float64(mage.ArcaneBlastAura.GetStacks()) * 1.75
			}
			return mage.CurrentMana() >= manaCost
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if mage.ArcanePotencyAura != nil {
				mage.ArcanePotencyAura.Activate(sim)
			}
			normalCastTime := spellToUse.DefaultCast.CastTime
			spellToUse.DefaultCast.CastTime = 0
			spellToUse.Cast(sim, mage.CurrentTarget)
			spellToUse.DefaultCast.CastTime = normalCastTime
		},
	})
	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}

func (mage *Mage) applyArcanePotency() {
	if mage.Talents.ArcanePotency == 0 {
		return
	}

	arcanePotencyMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  MageSpellsAllDamaging,
		FloatValue: []float64{0.0, 0.07, 0.15}[mage.Talents.ArcanePotency] * core.CritRatingPerCritChance,
		Kind:       core.SpellMod_BonusCrit_Rating,
	})

	mage.ArcanePotencyAura = mage.RegisterAura(core.Aura{
		Label:     "Arcane Potency",
		Duration:  core.NeverExpires,
		MaxStacks: 2,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			//prevent the spell that procced it from spending it
			core.StartDelayedAction(sim, core.DelayedActionOptions{
				DoAt: sim.CurrentTime + time.Millisecond*10,
				OnAction: func(sim *core.Simulation) {
					aura.SetStacks(sim, 2)
					arcanePotencyMod.Activate()
				},
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			arcanePotencyMod.Deactivate()
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// Only remove a stack if it's an applicable spell
			if spell.ClassSpellMask == arcanePotencyMod.ClassMask {
				aura.RemoveStack(sim)
			}
		},
	},
	)
}

func (mage *Mage) registerArcanePowerCD() {
	if !mage.Talents.ArcanePower {
		return
	}
	actionID := core.ActionID{SpellID: 12042}

	var affectedSpells []*core.Spell
	mage.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Flags.Matches(SpellFlagMage) {
			affectedSpells = append(affectedSpells, spell)
		}
	})

	mage.ArcanePowerAura = mage.RegisterAura(core.Aura{
		Label:    "Arcane Power",
		ActionID: actionID,
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.DamageMultiplierAdditive += 0.2
				spell.CostMultiplier += 0.1
			}
			mage.arcanePowerGCDmod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.DamageMultiplierAdditive -= 0.2
				spell.CostMultiplier -= 0.2
			}
			mage.arcanePowerGCDmod.Deactivate()
		},
	})
	core.RegisterPercentDamageModifierEffect(mage.ArcanePowerAura, 1.2)

	spell := mage.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second * time.Duration(120*(1-[]float64{0.0, 0.07, 0.15}[mage.Talents.ArcaneFlows])),
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			mage.ArcanePowerAura.Activate(sim)
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !mage.ArcanePotencyAura.IsActive()
		},
	})
	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}
