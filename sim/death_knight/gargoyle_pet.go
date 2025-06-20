package death_knight

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

type GargoylePet struct {
	core.Pet

	expireTime time.Duration
	dkOwner    *DeathKnight

	GargoyleStrike *core.Spell
}

func (dk *DeathKnight) NewGargoyle() *GargoylePet {
	gargoyle := &GargoylePet{
		Pet: core.NewPet(core.PetConfig{
			Name:      "Gargoyle",
			Owner:     &dk.Character,
			BaseStats: stats.Stats{},
			StatInheritance: func(ownerStats stats.Stats) stats.Stats {
				hitRating := ownerStats[stats.HitRating]
				expertiseRating := ownerStats[stats.ExpertiseRating]
				combined := (hitRating + expertiseRating) * 0.5

				return stats.Stats{
					stats.Armor:            ownerStats[stats.Armor],
					stats.CritRating:       ownerStats[stats.CritRating],
					stats.ExpertiseRating:  combined,
					stats.HasteRating:      ownerStats[stats.HasteRating],
					stats.Health:           ownerStats[stats.Health],
					stats.HitRating:        combined,
					stats.SpellCritPercent: ownerStats[stats.SpellCritPercent],
					stats.SpellPower:       ownerStats[stats.AttackPower] * 0.7,
					stats.Stamina:          ownerStats[stats.Stamina],
				}
			},
			EnabledOnStart:                 false,
			IsGuardian:                     true,
			HasDynamicCastSpeedInheritance: true,
		}),
		dkOwner: dk,
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

func (garg *GargoylePet) SetExpireTime(expireTime time.Duration) {
	garg.expireTime = expireTime
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
		SpellSchool: core.SpellSchoolPlague,
		ProcMask:    core.ProcMaskSpellDamage,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDMin,
				CastTime: time.Millisecond * 2000,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   garg.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		BonusCoefficient: 0.8259999752,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := garg.dkOwner.CalcAndRollDamageRange(sim, 0.5, 0.5)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}
