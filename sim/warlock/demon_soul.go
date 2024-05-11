package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (warlock *Warlock) registerDemonSoulSpell() {

	impMod := warlock.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusCrit_Rating,
		ClassMask:  WarlockSpellShadowBolt | WarlockSpellIncinerate | WarlockSpellSoulFire | WarlockSpellChaosBolt,
		FloatValue: 30 * core.CritRatingPerCritChance,
	})

	demonSoulImp := warlock.RegisterAura(core.Aura{
		Label:    "Demon Soul: Imp",
		ActionID: core.ActionID{SpellID: 79459},
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			impMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			impMod.Deactivate()
		},
	})

	felhunterMod := warlock.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  WarlockPeriodicShadowDamage,
		FloatValue: 0.2,
	})

	demonSoulFelhunter := warlock.RegisterAura(core.Aura{
		Label:    "Demon Soul: Felhunter",
		ActionID: core.ActionID{SpellID: 79460},
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			felhunterMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			felhunterMod.Deactivate()
		},
	})

	felguardHasteMod := warlock.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_CastTime_Pct,
		ClassMask:  WarlockFireDamage | WarlockShadowDamage,
		FloatValue: -0.15,
	})

	felguardDamageMod := warlock.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  WarlockFireDamage | WarlockShadowDamage,
		FloatValue: 0.1,
	})

	demonSoulFelguard := warlock.RegisterAura(core.Aura{
		Label:    "Demon Soul: Felguard",
		ActionID: core.ActionID{SpellID: 79462},
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			felguardHasteMod.Activate()
			felguardDamageMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			felguardHasteMod.Deactivate()
			felguardDamageMod.Deactivate()
		},
	})

	succubusDamageMod := warlock.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  WarlockSpellShadowBolt,
		FloatValue: 0.1,
	})

	demonSoulSuccubus := warlock.RegisterAura(core.Aura{
		Label:    "Demon Soul: Succubus",
		ActionID: core.ActionID{SpellID: 79463},
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			succubusDamageMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			succubusDamageMod.Deactivate()
		},
	})

	warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 77801},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellDemonSoul,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.15,
			Multiplier: 1,
		},

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Minute * 2,
			},
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if warlock.Felguard.IsActive() {
				demonSoulFelguard.Activate(sim)
			} else if warlock.Felhunter.IsActive() {
				demonSoulFelhunter.Activate(sim)
			} else if warlock.Imp.IsActive() {
				demonSoulImp.Activate(sim)
			} else if warlock.Succubus.IsActive() {
				demonSoulSuccubus.Activate(sim)
			}
		},
	})
}
