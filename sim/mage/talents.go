package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

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
			if !spell.Matches(MageSpellsAll ^ MageSpellInstantCast ^ MageSpellEvocation) {
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
				Duration: time.Second * 90,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return mage.GCD.IsReady(sim)
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			mage.presenceOfMindAura.Activate(sim)
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: pomSpell,
		Type:  core.CooldownTypeDPS,
	})
}

func (mage *Mage) registerIceFloesCD() {
	if !mage.Talents.IceFloes {
		return
	}

	iceFloesMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask: MageSpellsAll ^ MageSpellInstantCast ^ MageSpellEvocation,
		Kind:      core.SpellMod_AllowCastWhileMoving,
	})

	var iceFloesSpell *core.Spell

	mage.iceFloesfAura = mage.RegisterAura(core.Aura{
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
			if !spell.Matches(MageSpellsAll ^ MageSpellInstantCast ^ MageSpellEvocation) {
				return
			}
			if spell.DefaultCast.CastTime == 0 {
				return
			}
			aura.Deactivate(sim)
		},
	})

	iceFloesSpell = mage.RegisterSpell(core.SpellConfig{
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
			mage.presenceOfMindAura.Activate(sim)
		},
	})

}

func (mage *Mage) registerRuneOfPower() {
	if !mage.Talents.RuneOfPower {
		return
	}

	mage.runeOfPowerAura = mage.RegisterAura(core.Aura{
		Label:    "Rune of Power",
		ActionID: core.ActionID{SpellID: 116011},
		Duration: time.Minute,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: .15,
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
			mage.runeOfPowerAura.Activate(sim)
		},
	})
}
