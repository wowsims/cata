package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (mage *Mage) ApplyTalents() {

	// Level 15
	mage.registerPresenceOfMind()
	mage.registerIceFloes()

	// Level 30

	// Level 45

	// Level 75
	mage.registerNetherTempest()
	mage.registerLivingBomb()
	mage.registerFrostBomb()

	// Level 90
	mage.registerRuneOfPower()
	mage.registerInvocation()

}

func (mage *Mage) registerPresenceOfMind() {
	if !mage.Talents.PresenceOfMind {
		return
	}

	presenceOfMindMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  MageSpellsAll ^ (MageSpellInstantCast | MageSpellBlizzard | MageSpellEvocation),
		FloatValue: -1,
		Kind:       core.SpellMod_CastTime_Pct,
	})

	pomSpell := mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 12043},
		Flags:          core.SpellFlagNoOnCastComplete,
		ClassSpellMask: MageSpellPresenceOfMind,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second * 90,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return mage.GCD.IsReady(sim)
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			mage.PresenceOfMindAura.Activate(sim)
		},
	})

	mage.PresenceOfMindAura = mage.RegisterAura(core.Aura{
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
			if !spell.Matches(MageSpellsAll ^ (MageSpellInstantCast | MageSpellEvocation)) {
				return
			}
			if spell.DefaultCast.CastTime == 0 {
				return
			}
			aura.Deactivate(sim)
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: pomSpell,
		Type:  core.CooldownTypeDPS,
	})
}

func (mage *Mage) registerIceFloes() {
	if !mage.Talents.IceFloes {
		return
	}

	iceFloesMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask: MageSpellsAll ^ (MageSpellInstantCast | MageSpellEvocation),
		Kind:      core.SpellMod_AllowCastWhileMoving,
	})

	iceFloesSpell := mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 108839},
		Flags:          core.SpellFlagNoOnCastComplete, //Need to investigate this
		ClassSpellMask: MageSpellIceFloes,
		Charges:        3,
		RechargeTime:   time.Second * 20,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return mage.GCD.IsReady(sim)
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			mage.PresenceOfMindAura.Activate(sim)
		},
	})

	mage.IceFloesAura = mage.RegisterAura(core.Aura{
		Label:    "Ice Floes",
		ActionID: core.ActionID{SpellID: 108839},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			iceFloesMod.Activate()
			iceFloesSpell.CD.Use(sim)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			iceFloesMod.Deactivate()
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !spell.Matches(MageSpellsAll ^ (MageSpellInstantCast | MageSpellEvocation)) {
				return
			}
			if spell.DefaultCast.CastTime == 0 {
				return
			}
			aura.Deactivate(sim)
		},
	})

}

func (mage *Mage) registerInvocation() {
	if !mage.Talents.Invocation {
		return
	}

	mage.AddStaticMod(core.SpellModConfig{
		ClassMask:  MageSpellEvocation,
		FloatValue: -1,
		Kind:       core.SpellMod_Cooldown_Multiplier,
	})

	mage.AddStaticMod(core.SpellModConfig{
		ClassMask: MageSpellEvocation,
		TimeValue: time.Second * -1.0,
		Kind:      core.SpellMod_DotTickLength_Flat,
	})

	mage.InvocationAura = mage.RegisterAura(core.Aura{
		Label:    "Invocation Aura",
		ActionID: core.ActionID{SpellID: 116257},
		Duration: time.Minute,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			mage.MultiplyManaRegenSpeed(sim, 0.5)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			mage.MultiplyManaRegenSpeed(sim, 1/0.5)
		},
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask:  MageSpellsAllDamaging,
		FloatValue: 0.15,
		Kind:       core.SpellMod_DamageDone_Pct,
	})

}

func (mage *Mage) registerRuneOfPower() {
	if !mage.Talents.RuneOfPower {
		return
	}

	mage.RuneOfPowerAura = mage.RegisterAura(core.Aura{
		Label:    "Rune of Power",
		ActionID: core.ActionID{SpellID: 116011},
		Duration: time.Minute,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			mage.MultiplyManaRegenSpeed(sim, 1.75)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			mage.MultiplyManaRegenSpeed(sim, 1/1.75)
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.15,
		ClassMask:  MageSpellsAllDamaging,
	})

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 116011},
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MagespellRuneOfPower,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return mage.GCD.IsReady(sim)
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			mage.RuneOfPowerAura.Activate(sim)
		},
	})
}
