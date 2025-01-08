package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (warlock *Warlock) registerDemonSoul() {

	demonSoulImp := warlock.RegisterAura(core.Aura{
		Label:    "Demon Soul: Imp",
		ActionID: core.ActionID{SpellID: 79459},
		Duration: 20 * time.Second,
	})

	demonSoulImp.AttachSpellMod(core.SpellModConfig{
		Kind: core.SpellMod_BonusCrit_Percent,
		ClassMask: WarlockSpellShadowBolt | WarlockSpellIncinerate | WarlockSpellSoulFire | WarlockSpellChaosBolt |
			WarlockSpellImmolate | WarlockSpellImmolateDot,
		FloatValue: 30,
	})

	demonSoulFelhunter := warlock.RegisterAura(core.Aura{
		Label:    "Demon Soul: Felhunter",
		ActionID: core.ActionID{SpellID: 79460},
		Duration: 20 * time.Second,
	})

	demonSoulFelhunter.AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		ClassMask:  WarlockPeriodicShadowDamage,
		FloatValue: 0.2,
	})

	felguardHasteMulti := 1.15

	demonSoulFelguard := warlock.RegisterAura(core.Aura{
		Label:    "Demon Soul: Felguard",
		ActionID: core.ActionID{SpellID: 79462},
		Duration: 20 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warlock.MultiplyCastSpeed(felguardHasteMulti)
			warlock.MultiplyAttackSpeed(sim, felguardHasteMulti)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.MultiplyCastSpeed(1 / felguardHasteMulti)
			warlock.MultiplyAttackSpeed(sim, 1/felguardHasteMulti)
		},
	})

	demonSoulFelguard.AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		School:     core.SpellSchoolShadow | core.SpellSchoolFire,
		FloatValue: 0.1,
	})

	demonSoulSuccubus := warlock.RegisterAura(core.Aura{
		Label:    "Demon Soul: Succubus",
		ActionID: core.ActionID{SpellID: 79463},
		Duration: 20 * time.Second,
	})

	demonSoulSuccubus.AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  WarlockSpellShadowBolt,
		FloatValue: 0.1,
	})

	demonSoul := warlock.RegisterSpell(core.SpellConfig{
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
				Duration: 2 * time.Minute,
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

	warlock.AddMajorCooldown(core.MajorCooldown{
		Spell: demonSoul,
		Type:  core.CooldownTypeDPS,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return true
		},
	})
}
