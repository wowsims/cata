package death_knight

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

func (dk *DeathKnight) registerSummonGargoyleSpell() {
	if !dk.Talents.SummonGargoyle {
		return
	}

	trackingAura := dk.RegisterAura(core.Aura{
		Label:    "Summon Gargoyle",
		ActionID: core.ActionID{SpellID: 49206},
		Duration: time.Second * 30,
	})

	spell := dk.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 49206},
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: DeathKnightSpellSummonGargoyle,

		RuneCost: core.RuneCostOptions{
			RunicPowerCost: 60,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			trackingAura.Activate(sim)

			dk.Gargoyle.expireTime = sim.CurrentTime + time.Second*30
			dk.Gargoyle.EnableWithTimeout(sim, dk.Gargoyle, time.Second*30)
			// Start casting after a 2.5s delay to simulate the summon animation
			dk.Gargoyle.SetGCDTimer(sim, sim.CurrentTime+time.Millisecond*2500)
		},
	})

	dk.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}

type GargoylePet struct {
	core.Pet

	expireTime time.Duration
	dkOwner    *DeathKnight

	GargoyleStrike *core.Spell
}

func (dk *DeathKnight) NewGargoyle() *GargoylePet {
	gargoyle := &GargoylePet{
		Pet: core.NewPet(core.PetConfig{
			Name:  "Gargoyle",
			Owner: &dk.Character,
			BaseStats: stats.Stats{
				stats.Stamina:         1000,
				stats.SpellHitPercent: -float64(dk.Talents.NervesOfColdSteel) * HitCapRatio,
			},
			StatInheritance: func(ownerStats stats.Stats) stats.Stats {
				return stats.Stats{
					stats.SpellPower:       ownerStats[stats.AttackPower] * 0.7,
					stats.SpellHitPercent:  ownerStats[stats.PhysicalHitPercent] * HitCapRatio,
					stats.HasteRating:      ownerStats[stats.HasteRating],
					stats.SpellCritPercent: ownerStats[stats.SpellCritPercent],
				}
			},
			EnabledOnStart: false,
			IsGuardian:     true,
		}),
		dkOwner: dk,
	}

	gargoyle.OnPetEnable = func(sim *core.Simulation) {
		gargoyle.PseudoStats.CastSpeedMultiplier = 1 // guardians are not affected by raid buffs
		gargoyle.MultiplyCastSpeed(dk.PseudoStats.MeleeSpeedMultiplier)

		// No longer updates dynamically
		// gargoyle.EnableDynamicMeleeSpeed(func(amount float64) {
		// 	gargoyle.MultiplyCastSpeed(amount)
		// })

		// gargoyle.EnableDynamicStats(func(ownerStats stats.Stats) stats.Stats {
		// 	return stats.Stats{
		// 		stats.SpellHaste: ownerStats[stats.MeleeHaste],
		// 		stats.SpellHit:   ownerStats[stats.MeleeHit] * PetSpellHitScale,
		// 		stats.SpellCrit:  ownerStats[stats.SpellCrit],
		// 	}
		// })
	}

	dk.AddPet(gargoyle)

	return gargoyle
}

func (garg *GargoylePet) GetPet() *core.Pet {
	return &garg.Pet
}

func (garg *GargoylePet) Initialize() {
	garg.registerGargoyleStrikeSpell()
}

func (garg *GargoylePet) Reset(_ *core.Simulation) {
}

func (garg *GargoylePet) ExecuteCustomRotation(sim *core.Simulation) {
	if garg.GargoyleStrike.CanCast(sim, garg.CurrentTarget) {
		gargCastTime := garg.ApplyCastSpeedForSpell(garg.GargoyleStrike.DefaultCast.CastTime, garg.GargoyleStrike)
		if sim.CurrentTime+gargCastTime > garg.expireTime {
			// If the cast wont finish before expiration time just dont cast
			return
		}

		garg.GargoyleStrike.Cast(sim, garg.CurrentTarget)
	}
}

func (garg *GargoylePet) registerGargoyleStrikeSpell() {
	garg.GargoyleStrike = garg.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 51963},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellDamage,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2000,
			},
			IgnoreHaste: true,
			// Custom modify cast to not lower GCD
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.CastTime = spell.Unit.ApplyCastSpeedForSpell(spell.DefaultCast.CastTime, spell)
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   2,
		ThreatMultiplier: 1,

		BonusCoefficient: 0.453,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 291.0
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.DealDamage(sim, result)
		},
	})
}
