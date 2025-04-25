package mage

import (
	//"github.com/wowsims/mop/sim/core/proto"

	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

func (mage *Mage) ApplyFrostTalents() {

	// Cooldowns/Special Implementations
	mage.registerColdSnapCD()
	mage.registerIcyVeinsCD()
	mage.applyFingersOfFrost()
	// mage.applyImprovedFreeze()
	mage.applyBrainFreeze()

	//Early Frost
	if mage.Talents.EarlyFrost > 0 {
		earlyFrostMod := mage.AddDynamicMod(core.SpellModConfig{
			ClassMask: MageSpellFrostbolt,
			TimeValue: []time.Duration{0, time.Millisecond * -300, time.Millisecond * -600}[mage.Talents.EarlyFrost],
			Kind:      core.SpellMod_CastTime_Flat,
		})

		mage.RegisterAura(core.Aura{
			Label:    "Early Frost",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
				earlyFrostMod.Activate()
			},
			Icd: &core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second * 15,
			},

			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if aura.Icd.IsReady(sim) {
					aura.Icd.Use(sim)
					earlyFrostMod.Activate()
				}
				if spell.ClassSpellMask == MageSpellFrostfireBolt && earlyFrostMod.IsActive {
					core.StartDelayedAction(sim, core.DelayedActionOptions{
						DoAt: sim.CurrentTime + 10*time.Millisecond,
						OnAction: func(sim *core.Simulation) {
							earlyFrostMod.Deactivate()
						},
					})
				}
			},
		})
	}

	// Piercing Ice
	if mage.Talents.PiercingIce > 0 {
		mage.AddStat(stats.SpellCritPercent, 1*float64(mage.Talents.PiercingIce))
	}

	// Ice Floes inside spells

	//Piercing Chill

	//Ice Shards inside blizzard

}

func (mage *Mage) registerIcyVeinsCD() {
	if !mage.Talents.IcyVeins {
		return
	}

	icyVeinsMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  MageSpellsAll,
		FloatValue: -0.2,
		Kind:       core.SpellMod_CastTime_Pct,
	})

	actionID := core.ActionID{SpellID: 12472}
	icyVeinsAura := mage.RegisterAura(core.Aura{
		Label:    "Icy Veins",
		ActionID: actionID,
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			icyVeinsMod.Activate()
		},
		OnExpire: func(_ *core.Aura, sim *core.Simulation) {
			icyVeinsMod.Deactivate()
			if mage.t13ProcAura != nil {
				mage.t13ProcAura.Deactivate(sim)
			}
		},
	})

	mage.IcyVeins = mage.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: MageSpellIcyVeins,
		Flags:          core.SpellFlagNoOnCastComplete,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 3,
		},

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second * time.Duration(180*[]float64{1, .93, .86, .80}[mage.Talents.IceFloes]),
			},
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			// Need to check for icy veins already active in case Cold Snap is used right after.
			return !icyVeinsAura.IsActive() && !mage.frostfireOrb.IsEnabled()
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			icyVeinsAura.Activate(sim)
			if mage.T13_4pc.IsActive() {
				spell.CD.Reduce(time.Second * time.Duration(15*mage.t13ProcAura.GetStacks()))
			}
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: mage.IcyVeins,
		Type:  core.CooldownTypeDPS,
	})
}

func (mage *Mage) applyFingersOfFrost() {
	if mage.Talents.FingersOfFrost == 0 {
		return
	}

	//Talent gives 7/14/20 percent chance to proc FoF on spell hit
	procChance := []float64{0, 0.07, 0.14, 0.20}[mage.Talents.FingersOfFrost]

	fingersOfFrostIceLanceDamageMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  MageSpellIceLance,
		FloatValue: 0.25,
		Kind:       core.SpellMod_DamageDone_Pct,
	})

	fingersOfFrostFrozenCritMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  MageSpellIceLance | MageSpellDeepFreeze,
		FloatValue: mage.GetStat(stats.SpellCritPercent) * 2, // TODO: Shouldn't this be evaluated dynamically when the proc happens?
		Kind:       core.SpellMod_BonusCrit_Percent,
	})

	mage.FingersOfFrostAura = mage.RegisterAura(core.Aura{
		Label:     "Fingers of Frost Proc",
		ActionID:  core.ActionID{SpellID: 44545},
		Duration:  time.Second * 15,
		MaxStacks: 2,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			fingersOfFrostFrozenCritMod.Activate()
			fingersOfFrostIceLanceDamageMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			fingersOfFrostFrozenCritMod.Deactivate()
			fingersOfFrostIceLanceDamageMod.Deactivate()
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.ClassSpellMask&(MageSpellIceLance|MageSpellDeepFreeze) > 0 {
				aura.RemoveStack(sim)
			}
		},
	})

	mage.RegisterAura(core.Aura{
		Label:    "Fingers of Frost Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if mage.hasChillEffect(spell) && sim.Proc(procChance, "FingersOfFrostProc") {
				mage.FingersOfFrostAura.Activate(sim)
				mage.FingersOfFrostAura.AddStack(sim)
			}
		},
	})
}

func (mage *Mage) applyBrainFreeze() {
	if mage.Talents.BrainFreeze == 0 {
		return
	}

	mage.brainFreezeProcChance = .05 * float64(mage.Talents.BrainFreeze)

	brainFreezeCostMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask: MageSpellBrainFreeze,
		IntValue:  -100,
		Kind:      core.SpellMod_PowerCost_Pct,
	})

	brainFreezeCastMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  MageSpellBrainFreeze,
		FloatValue: -1,
		Kind:       core.SpellMod_CastTime_Pct,
	})

	brainFreezeAura := mage.GetOrRegisterAura(core.Aura{
		Label:    "Brain Freeze Proc",
		ActionID: core.ActionID{SpellID: 57761},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			brainFreezeCostMod.Activate()
			brainFreezeCastMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			brainFreezeCostMod.Deactivate()
			brainFreezeCastMod.Deactivate()
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.ClassSpellMask&(MageSpellFrostfireBolt|MageSpellFireball) > 0 {
				aura.Deactivate(sim)
			}
		},
	})

	mage.RegisterAura(core.Aura{
		Label:    "Brain Freeze Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if mage.hasChillEffect(spell) && sim.Proc(mage.brainFreezeProcChance, "Brain Freeze") {
				brainFreezeAura.Activate(sim)
			}
		},
	})
}

func (mage *Mage) registerColdSnapCD() {
	if !mage.Talents.ColdSnap {
		return
	}

	actionID := core.ActionID{SpellID: 11958}
	spell := mage.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Minute * 8,
			},
		},
		/* 		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			// Don't use if there are no cooldowns to reset.
			return (mage.IcyVeins != nil && !mage.IcyVeins.IsReady(sim)) ||
				(mage.SummonWaterElemental != nil && !mage.SummonWaterElemental.IsReady(sim))
		}, */
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if mage.IcyVeins != nil {
				mage.IcyVeins.CD.Reset()
			}
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			// Ideally wait for both water ele and icy veins so we can reset both.
			if mage.IcyVeins != nil && mage.IcyVeins.IsReady(sim) {
				return false
			}

			if mage.SummonWaterElemental != nil && mage.SummonWaterElemental.IsReady(sim) {
				return false
			}
			// You want to reset orb
			if mage.FrostfireOrb.IsReady(sim) {
				return false
			}
			return true
		},
	})
}
