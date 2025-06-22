package death_knight

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

type GhoulPet struct {
	core.Pet

	dkOwner     *DeathKnight
	clawSpellID int32
	summonDelay bool

	DarkTransformationAura *core.Aura
	ShadowInfusionAura     *core.Aura
	Claw                   *core.Spell
}

func (dk *DeathKnight) NewArmyGhoulPet() *GhoulPet {
	return dk.newGhoulPetInternal("Army of the Dead", false, 0.5)
}

func (dk *DeathKnight) NewGhoulPet(permanent bool) *GhoulPet {
	return dk.newGhoulPetInternal("Ghoul", permanent, 0.8)
}

func (dk *DeathKnight) NewFallenZandalariPet() *GhoulPet {
	troll := dk.newGhoulPetInternal("Fallen Zandalari", false, 0.8)
	troll.summonDelay = false

	// Fallen Zandalari use their own spell called Strike, which does 150% damage
	troll.clawSpellID = 138537
	troll.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  GhoulSpellClaw,
		FloatValue: 0.2,
	})

	// Command doesn't apply to Fallen Zandalari
	if dk.Race == proto.Race_RaceOrc {
		troll.PseudoStats.DamageDealtMultiplier /= 1.02
	}
	return troll
}

func (dk *DeathKnight) newGhoulPetInternal(name string, permanent bool, scalingCoef float64) *GhoulPet {
	ghoulPet := &GhoulPet{
		Pet: core.NewPet(core.PetConfig{
			Name:                            name,
			Owner:                           &dk.Character,
			BaseStats:                       stats.Stats{stats.AttackPower: -20},
			StatInheritance:                 dk.ghoulStatInheritance(scalingCoef),
			EnabledOnStart:                  permanent,
			IsGuardian:                      !permanent,
			HasDynamicMeleeSpeedInheritance: true,
		}),
		dkOwner:     dk,
		clawSpellID: 91776,
		summonDelay: true,
	}

	ghoulPet.PseudoStats.DamageTakenMultiplier *= 0.1

	dk.SetupGhoul(ghoulPet, scalingCoef)

	return ghoulPet
}

func (dk *DeathKnight) SetupGhoul(ghoulPet *GhoulPet, scalingCoef float64) {
	baseDamage := dk.CalcScalingSpellDmg(scalingCoef)
	ghoulPet.EnableAutoAttacks(ghoulPet, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:     baseDamage,
			BaseDamageMax:     baseDamage,
			SwingSpeed:        2,
			CritMultiplier:    dk.DefaultCritMultiplier(),
			AttackPowerPerDPS: core.DefaultAttackPowerPerDPS,
		},
		AutoSwingMelee: true,
	})

	ghoulPet.Unit.EnableFocusBar(100, 10.0, false, nil)

	dk.AddPet(ghoulPet)
}

func (ghoulPet *GhoulPet) GetPet() *core.Pet {
	return &ghoulPet.Pet
}

func (ghoulPet *GhoulPet) Initialize() {
	ghoulPet.Claw = ghoulPet.registerClaw()
}

func (ghoulPet *GhoulPet) Reset(_ *core.Simulation) {
}

func (ghoulPet *GhoulPet) ExecuteCustomRotation(sim *core.Simulation) {
	if !ghoulPet.GCD.IsReady(sim) {
		return
	}

	if ghoulPet.CurrentFocus() < ghoulPet.Claw.DefaultCast.Cost {
		ghoulPet.ExtendGCDUntil(sim, sim.CurrentTime+core.DurationFromSeconds((ghoulPet.Claw.DefaultCast.Cost-ghoulPet.CurrentFocus())/ghoulPet.FocusRegenPerSecond()))
		return
	}

	if ghoulPet.Claw.CanCast(sim, ghoulPet.CurrentTarget) {
		ghoulPet.Claw.Cast(sim, ghoulPet.CurrentTarget)
	}
}

func (dk *DeathKnight) ghoulStatInheritance(apCoef float64) core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		hitRating := ownerStats[stats.HitRating]
		expertiseRating := ownerStats[stats.ExpertiseRating]
		combined := (hitRating + expertiseRating) * 0.5

		return stats.Stats{
			stats.Armor:               ownerStats[stats.Armor],
			stats.AttackPower:         ownerStats[stats.AttackPower] * apCoef,
			stats.CritRating:          ownerStats[stats.CritRating],
			stats.ExpertiseRating:     combined,
			stats.HasteRating:         ownerStats[stats.HasteRating],
			stats.Health:              ownerStats[stats.Health],
			stats.HitRating:           combined,
			stats.PhysicalCritPercent: ownerStats[stats.PhysicalCritPercent],
			stats.Stamina:             ownerStats[stats.Stamina],
		}
	}
}

func (ghoulPet *GhoulPet) registerClaw() *core.Spell {
	return ghoulPet.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: ghoulPet.clawSpellID},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics,
		ClassSpellMask: GhoulSpellClaw,

		FocusCost: core.FocusCostOptions{
			Cost:   40,
			Refund: 0.8,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDMin,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1.25,
		CritMultiplier:   ghoulPet.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			results := make([]*core.SpellResult, min(int32(3), min(ghoulPet.Env.GetNumTargets(), core.TernaryInt32(ghoulPet.DarkTransformationAura.IsActive(), 3, 1))))

			for idx := range results {
				baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
				results[idx] = spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
				target = sim.Environment.NextTargetUnit(target)
			}

			for idx, result := range results {
				if idx == 0 && !result.Landed() {
					spell.IssueRefund(sim)
				}
				spell.DealDamage(sim, result)
			}
		},
	})
}

func (ghoulPet *GhoulPet) Enable(sim *core.Simulation, petAgent core.PetAgent) {
	if ghoulPet.IsGuardian() && ghoulPet.summonDelay {
		// The ghoul takes around 4.5s - 5s to from summon to first hit, depending on your distance from the target.
		randomDelay := core.DurationFromSeconds(sim.RollWithLabel(4.5, 5, "Raise Dead Delay"))
		ghoulPet.Pet.EnableWithStartAttackDelay(sim, petAgent, randomDelay)
	} else {
		ghoulPet.Pet.Enable(sim, petAgent)
	}
}

const (
	GhoulSpellNone int64 = 0
	GhoulSpellClaw int64 = 1 << iota

	GhoulSpellLast
	GhoulSpellsAll = DeathKnightSpellLast<<1 - 1
)
