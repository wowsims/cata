package death_knight

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func CopySpellMultipliers(sourceSpell *core.Spell, targetSpell *core.Spell, target *core.Unit) {
	targetSpell.DamageMultiplier = sourceSpell.DamageMultiplier
	targetSpell.DamageMultiplierAdditive = sourceSpell.DamageMultiplierAdditive
	targetSpell.BonusCritPercent = sourceSpell.BonusCritPercent
	targetSpell.BonusHitPercent = sourceSpell.BonusHitPercent
	targetSpell.CritMultiplier = sourceSpell.CritMultiplier
	targetSpell.ThreatMultiplier = sourceSpell.ThreatMultiplier

	if sourceSpell.Dot(target) != nil {
		sourceDot := sourceSpell.Dot(target)
		targetDot := targetSpell.Dot(target)

		targetDot.BaseTickCount = sourceDot.BaseTickCount
		targetDot.BaseTickLength = sourceDot.BaseTickLength
	}
}

func (dk *DeathKnight) registerDancingRuneWeaponSpell() {
	if !dk.Talents.DancingRuneWeapon {
		return
	}

	duration := time.Second * 12

	hasGlyph := dk.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfDancingRuneWeapon)

	t124PAura := dk.RegisterAura(core.Aura{
		Label:    "Flaming Rune Weapon",
		ActionID: core.ActionID{SpellID: 101162},
		Duration: time.Second * 12,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dk.PseudoStats.BaseParryChance += 0.15
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.PseudoStats.BaseParryChance -= 0.15
		},
	})

	dancingRuneWeaponAura := dk.RegisterAura(core.Aura{
		Label:    "Dancing Rune Weapon",
		ActionID: core.ActionID{SpellID: 81256},
		Duration: duration,
		// Casts copy
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			copySpell := dk.RuneWeapon.runeWeaponSpells[spell.ActionID]
			if copySpell == nil {
				return
			}

			CopySpellMultipliers(spell, copySpell, dk.CurrentTarget)

			copySpell.Cast(sim, dk.CurrentTarget)
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dk.PseudoStats.BaseParryChance += 0.20
			if hasGlyph {
				dk.PseudoStats.ThreatMultiplier *= 1.5
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.PseudoStats.BaseParryChance -= 0.20
			if hasGlyph {
				dk.PseudoStats.ThreatMultiplier /= 1.5
			}
			if dk.T12Tank4pc.IsActive() {
				t124PAura.Activate(sim)
			}
		},
	})

	spell := dk.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 49028},
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: DeathKnightSpellDancingRuneWeapon,

		RuneCost: core.RuneCostOptions{
			RunicPowerCost: 60,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Second * 90,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			dk.RuneWeapon.EnableWithTimeout(sim, dk.RuneWeapon, duration)
			dk.RuneWeapon.CancelGCDTimer(sim)
			dancingRuneWeaponAura.Activate(sim)
		},
	})

	dk.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}

type RuneWeaponPet struct {
	core.Pet

	dkOwner *DeathKnight

	// Diseases
	FrostFeverSpell  *core.Spell
	BloodPlagueSpell *core.Spell

	drwDmgSnapshot  float64
	drwPhysSnapshot float64

	runeWeaponSpells map[core.ActionID]*core.Spell
}

func (runeWeapon *RuneWeaponPet) Initialize() {
	runeWeapon.dkOwner.registerDrwFrostFever()
	runeWeapon.dkOwner.registerDrwBloodPlague()
	runeWeapon.AddCopySpell(OutbreakActionID, runeWeapon.dkOwner.registerDrwOutbreakSpell())
	runeWeapon.AddCopySpell(IcyTouchActionID, runeWeapon.dkOwner.registerDrwIcyTouchSpell())
	runeWeapon.AddCopySpell(BloodBoilActionID, runeWeapon.dkOwner.registerDrwBloodBoilSpell())
	runeWeapon.AddCopySpell(DeathCoilActionID, runeWeapon.dkOwner.registerDrwDeathCoilSpell())
	runeWeapon.AddCopySpell(PlagueStrikeActionID.WithTag(1), runeWeapon.dkOwner.registerDrwPlagueStrikeSpell())
	runeWeapon.AddCopySpell(DeathStrikeActionID.WithTag(1), runeWeapon.dkOwner.registerDrwDeathStrikeSpell())
	runeWeapon.AddCopySpell(RuneStrikeActionID.WithTag(1), runeWeapon.dkOwner.registerDrwRuneStrikeSpell())
	runeWeapon.AddCopySpell(FesteringStrikeActionID.WithTag(1), runeWeapon.dkOwner.registerDrwFesteringStrikeSpell())
	runeWeapon.AddCopySpell(BloodStrikeActionID.WithTag(1), runeWeapon.dkOwner.registerDrwBloodStrikeSpell())
}

func (runeWeapon *RuneWeaponPet) AddCopySpell(actionId core.ActionID, spell *core.Spell) {
	runeWeapon.runeWeaponSpells[actionId] = spell
}

func (dk *DeathKnight) NewRuneWeapon() *RuneWeaponPet {
	runeWeapon := &RuneWeaponPet{
		Pet: core.NewPet(core.PetConfig{
			Name:  "Rune Weapon",
			Owner: &dk.Character,
			BaseStats: stats.Stats{
				stats.Stamina: 100,
			},
			StatInheritance: func(ownerStats stats.Stats) stats.Stats {
				return stats.Stats{
					stats.AttackPower: ownerStats[stats.AttackPower],
					stats.HasteRating: ownerStats[stats.HasteRating],

					stats.PhysicalHitPercent: ownerStats[stats.PhysicalHitPercent],
					stats.SpellHitPercent:    ownerStats[stats.PhysicalHitPercent] * HitCapRatio,

					stats.ExpertiseRating: ownerStats[stats.PhysicalHitPercent] * PetExpertiseRatingScale,

					stats.PhysicalCritPercent: ownerStats[stats.PhysicalCritPercent],
					stats.SpellCritPercent:    ownerStats[stats.SpellCritPercent],
				}
			},
			EnabledOnStart: false,
			IsGuardian:     true,
		}),
		dkOwner: dk,
	}

	runeWeapon.runeWeaponSpells = map[core.ActionID]*core.Spell{}

	runeWeapon.OnPetEnable = runeWeapon.enable
	runeWeapon.OnPetDisable = runeWeapon.disable

	mhWeapon := dk.WeaponFromMainHand(dk.DefaultCritMultiplier())

	baseDamage := mhWeapon.AverageDamage() / mhWeapon.SwingSpeed * 3.5
	mhWeapon.BaseDamageMin = baseDamage - 150
	mhWeapon.BaseDamageMax = baseDamage + 150

	mhWeapon.SwingSpeed = 3.5
	mhWeapon.NormalizedSwingSpeed = 3.3

	runeWeapon.EnableAutoAttacks(runeWeapon, core.AutoAttackOptions{
		MainHand:       mhWeapon,
		AutoSwingMelee: true,
	})

	runeWeapon.PseudoStats.DamageTakenMultiplier = 0

	dk.AddPet(runeWeapon)

	return runeWeapon
}

func (runeWeapon *RuneWeaponPet) GetPet() *core.Pet {
	return &runeWeapon.Pet
}

func (runeWeapon *RuneWeaponPet) Reset(_ *core.Simulation) {
}

func (runeWeapon *RuneWeaponPet) ExecuteCustomRotation(_ *core.Simulation) {
}

func (runeWeapon *RuneWeaponPet) enable(sim *core.Simulation) {
	// Snapshot extra % speed modifiers from dk owner
	runeWeapon.PseudoStats.MeleeSpeedMultiplier = 1
	runeWeapon.MultiplyMeleeSpeed(sim, runeWeapon.dkOwner.PseudoStats.MeleeSpeedMultiplier)

	runeWeapon.drwDmgSnapshot = runeWeapon.dkOwner.PseudoStats.DamageDealtMultiplier * 0.5
	runeWeapon.PseudoStats.DamageDealtMultiplier *= runeWeapon.drwDmgSnapshot

	runeWeapon.drwPhysSnapshot = runeWeapon.dkOwner.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical]
	runeWeapon.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= runeWeapon.drwPhysSnapshot

}

func (runeWeapon *RuneWeaponPet) disable(sim *core.Simulation) {
	// Clear snapshot speed
	runeWeapon.PseudoStats.MeleeSpeedMultiplier = 1
	runeWeapon.MultiplyMeleeSpeed(sim, 1)

	// Clear snapshot damage multipliers
	runeWeapon.PseudoStats.DamageDealtMultiplier /= runeWeapon.drwDmgSnapshot
	runeWeapon.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= runeWeapon.drwPhysSnapshot
	runeWeapon.drwPhysSnapshot = 1
	runeWeapon.drwDmgSnapshot = 1
}
