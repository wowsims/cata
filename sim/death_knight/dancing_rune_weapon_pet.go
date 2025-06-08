package death_knight

import (
	"github.com/wowsims/mop/sim/core"
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

type RuneWeaponPet struct {
	core.Pet

	dkOwner *DeathKnight

	// Diseases
	FrostFeverSpell  *core.Spell
	BloodPlagueSpell *core.Spell

	drwDmgSnapshot  float64
	drwPhysSnapshot float64

	RuneWeaponSpells map[core.ActionID]*core.Spell
}

func (runeWeapon *RuneWeaponPet) Initialize() {
	runeWeapon.dkOwner.registerDrwFrostFever()
	runeWeapon.dkOwner.registerDrwBloodPlague()
	runeWeapon.AddCopySpell(OutbreakActionID, runeWeapon.dkOwner.registerDrwOutbreak())
	runeWeapon.AddCopySpell(IcyTouchActionID, runeWeapon.dkOwner.registerDrwIcyTouch())
	// runeWeapon.AddCopySpell(BloodBoilActionID, runeWeapon.dkOwner.registerDrwBloodBoilSpell())
	// runeWeapon.AddCopySpell(DeathCoilActionID, runeWeapon.dkOwner.registerDrwDeathCoilSpell())
	runeWeapon.AddCopySpell(PestilenceActionID, runeWeapon.dkOwner.registerDrwPestilence())
	runeWeapon.AddCopySpell(PlagueStrikeActionID, runeWeapon.dkOwner.registerDrwPlagueStrike())
	// runeWeapon.AddCopySpell(DeathStrikeActionID, runeWeapon.dkOwner.registerDrwDeathStrikeSpell())
}

func (runeWeapon *RuneWeaponPet) DiseasesAreActive(target *core.Unit) bool {
	return runeWeapon.FrostFeverSpell.Dot(target).IsActive() || runeWeapon.BloodPlagueSpell.Dot(target).IsActive()
}

func (runeWeapon *RuneWeaponPet) GetDiseaseMulti(target *core.Unit, base float64, increase float64) float64 {
	count := 0
	if runeWeapon.FrostFeverSpell.Dot(target).IsActive() {
		count++
	}
	if runeWeapon.BloodPlagueSpell.Dot(target).IsActive() {
		count++
	}
	return base + increase*float64(count)
}

func (runeWeapon *RuneWeaponPet) AddCopySpell(actionId core.ActionID, spell *core.Spell) {
	runeWeapon.RuneWeaponSpells[actionId] = spell
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
				hitRating := ownerStats[stats.HitRating]
				expertiseRating := ownerStats[stats.ExpertiseRating]
				combined := (hitRating + expertiseRating) * 0.5

				return stats.Stats{
					stats.AttackPower:         ownerStats[stats.AttackPower],
					stats.CritRating:          ownerStats[stats.CritRating],
					stats.ExpertiseRating:     combined,
					stats.HasteRating:         ownerStats[stats.HasteRating],
					stats.HitRating:           combined,
					stats.PhysicalCritPercent: ownerStats[stats.PhysicalCritPercent],
					stats.SpellCritPercent:    ownerStats[stats.SpellCritPercent],
				}
			},
			EnabledOnStart:                  false,
			IsGuardian:                      true,
			HasDynamicMeleeSpeedInheritance: true,
		}),
		dkOwner: dk,
	}

	runeWeapon.RuneWeaponSpells = map[core.ActionID]*core.Spell{}

	runeWeapon.OnPetEnable = runeWeapon.enable
	runeWeapon.OnPetDisable = runeWeapon.disable

	mhWeapon := dk.WeaponFromMainHand(dk.DefaultCritMultiplier())

	// baseDamage := mhWeapon.AverageDamage() / mhWeapon.SwingSpeed * 3.5
	baseDamage := dk.CalcScalingSpellDmg(3.0)
	mhWeapon.BaseDamageMin = baseDamage
	mhWeapon.BaseDamageMax = baseDamage

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
	runeWeapon.drwDmgSnapshot = runeWeapon.dkOwner.PseudoStats.DamageDealtMultiplier * 0.5
	runeWeapon.PseudoStats.DamageDealtMultiplier *= runeWeapon.drwDmgSnapshot

	runeWeapon.drwPhysSnapshot = runeWeapon.dkOwner.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical]
	runeWeapon.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= runeWeapon.drwPhysSnapshot
}

func (runeWeapon *RuneWeaponPet) disable(sim *core.Simulation) {
	// Clear snapshot damage multipliers
	runeWeapon.PseudoStats.DamageDealtMultiplier /= runeWeapon.drwDmgSnapshot
	runeWeapon.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= runeWeapon.drwPhysSnapshot
	runeWeapon.drwPhysSnapshot = 1
	runeWeapon.drwDmgSnapshot = 1
}
