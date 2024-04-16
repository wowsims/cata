package death_knight

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
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
			dk.Gargoyle.EnableWithTimeout(sim, dk.Gargoyle, time.Second*30)
			dk.Gargoyle.CancelGCDTimer(sim)

			trackingAura.Activate(sim)

			// Start casting after a 2.5s delay to simulate the summon animation
			pa := core.PendingAction{
				NextActionAt: sim.CurrentTime + time.Millisecond*2500,
				Priority:     core.ActionPriorityAuto,
				OnAction: func(s *core.Simulation) {
					dk.Gargoyle.GargoyleStrike.Cast(sim, dk.CurrentTarget)
				},
			}
			sim.AddPendingAction(&pa)
		},
	})

	dk.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}

type GargoylePet struct {
	core.Pet

	dkOwner *DeathKnight

	GargoyleStrike *core.Spell
}

func (dk *DeathKnight) NewGargoyle() *GargoylePet {
	gargoyle := &GargoylePet{
		Pet: core.NewPet("Gargoyle", &dk.Character, stats.Stats{
			stats.Stamina: 1000,
		}, func(ownerStats stats.Stats) stats.Stats {
			return stats.Stats{
				stats.SpellPower: ownerStats[stats.AttackPower],
				stats.SpellHit:   ownerStats[stats.MeleeHit] * PetSpellHitScale,
				stats.SpellHaste: ownerStats[stats.MeleeHaste],
				stats.SpellCrit:  ownerStats[stats.SpellCrit],
			}
		}, false, true),
		dkOwner: dk,
	}

	gargoyle.OnPetEnable = func(sim *core.Simulation) {
		gargoyle.PseudoStats.CastSpeedMultiplier = 1 // guardians are not affected by raid buffs
		gargoyle.MultiplyCastSpeed(dk.PseudoStats.MeleeSpeedMultiplier)

		gargoyle.EnableDynamicMeleeSpeed(func(amount float64) {
			gargoyle.MultiplyCastSpeed(amount)
		})

		gargoyle.EnableDynamicStats(func(ownerStats stats.Stats) stats.Stats {
			return stats.Stats{
				stats.SpellHaste: ownerStats[stats.MeleeHaste],
				stats.SpellHit:   ownerStats[stats.MeleeHit] * PetSpellHitScale,
				stats.SpellCrit:  ownerStats[stats.SpellCrit],
			}
		})
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

func (garg *GargoylePet) ExecuteCustomRotation(_ *core.Simulation) {
}

func (garg *GargoylePet) registerGargoyleStrikeSpell() {
	garg.GargoyleStrike = garg.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 51963},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellDamage,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				CastTime: time.Millisecond * 2000,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   1.5,
		ThreatMultiplier: 1,

		BonusCoefficient: 0.317,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 291.0
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.DealDamage(sim, result)

			garg.GargoyleStrike.Cast(sim, garg.CurrentTarget)
		},
	})
}
