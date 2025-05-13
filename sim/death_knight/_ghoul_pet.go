package death_knight

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

type GhoulPet struct {
	core.Pet

	dkOwner *DeathKnight

	DarkTransformationAura *core.Aura
	Claw                   *core.Spell

	uptimePercent float64
}

func (dk *DeathKnight) NewArmyGhoulPet(_ int) *GhoulPet {
	ghoulPet := &GhoulPet{
		Pet:     core.NewPet("Army of the Dead", &dk.Character, dk.ghoulBaseStats(), dk.ghoulStatInheritance(), false, true),
		dkOwner: dk,
	}

	ghoulPet.PseudoStats.DamageTakenMultiplier *= 0.1

	dk.SetupGhoul(ghoulPet, 14/0.0055)

	// command doesn't apply to army ghoul
	if dk.Race == proto.Race_RaceOrc {
		ghoulPet.PseudoStats.DamageDealtMultiplier /= 1.05
	}

	return ghoulPet
}

func (dk *DeathKnight) NewGhoulPet(permanent bool) *GhoulPet {
	ghoulPet := &GhoulPet{
		Pet:     core.NewPet("Ghoul", &dk.Character, dk.ghoulBaseStats(), dk.ghoulStatInheritance(), permanent, !permanent),
		dkOwner: dk,
	}

	dk.SetupGhoul(ghoulPet, 14)

	return ghoulPet
}

func (dk *DeathKnight) SetupGhoul(ghoulPet *GhoulPet, apScaling float64) {

	ghoulPet.EnableAutoAttacks(ghoulPet, core.AutoAttackOptions{
		MainHand: core.Weapon{
			// Base 240 DPS with observed around 300 range
			BaseDamageMin:     (240 - 75) * 2,
			BaseDamageMax:     (240 + 75) * 2,
			SwingSpeed:        2,
			CritMultiplier:    2,
			AttackPowerPerDPS: apScaling,
		},
		AutoSwingMelee: true,
	})

	ghoulPet.AddStatDependency(stats.Strength, stats.AttackPower, 2)

	ghoulPet.Pet.OnPetEnable = ghoulPet.enable

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
	if !ghoulPet.IsGuardian() {
		ghoulPet.uptimePercent = min(1, max(0, ghoulPet.dkOwner.Inputs.PetUptime))
	} else {
		ghoulPet.uptimePercent = 1.0
	}
}

func (ghoulPet *GhoulPet) ExecuteCustomRotation(sim *core.Simulation) {
	if ghoulPet.uptimePercent < 1.0 { // Apply uptime for permanent pet ghoul
		if sim.GetRemainingDurationPercent() < 1.0-ghoulPet.uptimePercent { // once fight is % completed, disable pet.
			ghoulPet.Pet.Disable(sim)
			return
		}
	}

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

func (ghoulPet *GhoulPet) enable(sim *core.Simulation) {
	if ghoulPet.IsGuardian() {
		ghoulPet.PseudoStats.MeleeSpeedMultiplier = 1 // guardians are not affected by raid buffs
		ghoulPet.MultiplyMeleeSpeed(sim, ghoulPet.dkOwner.PseudoStats.MeleeSpeedMultiplier)
		return
	}

	ghoulPet.MultiplyMeleeSpeed(sim, ghoulPet.dkOwner.PseudoStats.MeleeSpeedMultiplier)

	ghoulPet.EnableDynamicMeleeSpeed(func(amount float64) {
		ghoulPet.MultiplyMeleeSpeed(sim, amount)
	})
}

func (dk *DeathKnight) ghoulBaseStats() stats.Stats {
	nocsHitPercent := -float64(dk.Talents.NervesOfColdSteel)
	return stats.Stats{
		stats.Stamina:             388,
		stats.Agility:             3343 - 10, // We remove 10 to not mess with crit conversion
		stats.Strength:            476,
		stats.AttackPower:         -20,
		stats.PhysicalCritPercent: 5,

		// Remove bonus hit that would be transfered from the DKs Nerves of Cold Steel
		stats.PhysicalHitPercent: nocsHitPercent,
		stats.ExpertiseRating:    nocsHitPercent * PetExpertiseRatingScale,
	}
}

func (dk *DeathKnight) ghoulStatInheritance() core.PetStatInheritance {
	glyphBonusSta := core.TernaryFloat64(dk.HasPrimeGlyph(proto.DeathKnightPrimeGlyph_GlyphOfRaiseDead), 1.4, 1.0)
	glyphBonusStr := core.TernaryFloat64(dk.HasPrimeGlyph(proto.DeathKnightPrimeGlyph_GlyphOfRaiseDead), 0.4254, .0)

	return func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			stats.Stamina:  ownerStats[stats.Stamina] * (0.904 * glyphBonusSta),
			stats.Strength: ownerStats[stats.Strength] * (1.01 + glyphBonusStr),

			stats.PhysicalHitPercent: ownerStats[stats.PhysicalHitPercent],
			stats.ExpertiseRating:    ownerStats[stats.PhysicalHitPercent] * PetExpertiseRatingScale,

			stats.HasteRating:         ownerStats[stats.HasteRating],
			stats.PhysicalCritPercent: ownerStats[stats.PhysicalCritPercent],
		}
	}
}

func (ghoulPet *GhoulPet) registerClaw() *core.Spell {
	return ghoulPet.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 47468},
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
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1.25,
		CritMultiplier:   2,
		ThreatMultiplier: 1,

		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			results := make([]*core.SpellResult, min(int32(2), min(ghoulPet.Env.GetNumTargets(), core.TernaryInt32(ghoulPet.DarkTransformationAura.IsActive(), 2, 1))))

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
	if ghoulPet.IsGuardian() {
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
