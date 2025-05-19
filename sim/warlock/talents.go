package warlock

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (warlock *Warlock) registerArchimondesDarkness() {
	if !warlock.Talents.ArchimondesDarkness {
		return
	}

	warlock.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_ModCharges_Flat,
		IntValue:  2,
		ClassMask: WarlockSpellDarkSoulInsanity,
	})

	warlock.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: -time.Second * 100,
		ClassMask: WarlockSpellDarkSoulInsanity,
	})
}

func (warlock *Warlock) registerKilJaedensCunning() {
	if !warlock.Talents.KiljaedensCunning {
		return
	}

	warlock.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_AllowCastWhileMoving,
		ClassMask: WarlockSpellIncinerate | WarlockSpellShadowBolt | WarlockSpellMaleficGrasp,
	})
}

func (warlock *Warlock) registerMannarothsFury() {
	if !warlock.Talents.MannorothsFury {
		return
	}

	buff := warlock.RegisterAura(core.Aura{
		Label:    "Mannaroth's Fury",
		ActionID: core.ActionID{SpellID: 108508},
		Duration: time.Second * 10,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  WarlockSpellRainOfFire | WarlockSpellSeedOfCorruptionExposion | WarlockSpellSeedOfCorruption | WarlockSpellImmolationAura | WarlockSpellHellfire,
		FloatValue: 1,
	})

	warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 108508},
		Flags:       core.SpellFlagAPL,
		ProcMask:    core.ProcMaskEmpty,
		SpellSchool: core.SpellSchoolShadow,

		ThreatMultiplier: 1,
		DamageMultiplier: 1,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},

			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Minute,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			buff.Activate(sim)
		},
	})
}
