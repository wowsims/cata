package death_knight

import (
	"math"

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

	if sourceSpell.RelatedDotSpell != nil {
		CopySpellMultipliers(sourceSpell.RelatedDotSpell, targetSpell.RelatedDotSpell, target)
	}
}

type RuneWeaponPet struct {
	core.Pet

	dkOwner *DeathKnight

	// Diseases
	FrostFeverSpell  *core.Spell
	BloodPlagueSpell *core.Spell

	drwDmgSnapshot       float64
	drwSchoolDmgSnapshot [stats.SchoolLen]float64

	StrikeWeapon       *core.Weapon
	StrikeWeaponDamage float64

	RuneWeaponSpells map[core.ActionID]*core.Spell
}

func (runeWeapon *RuneWeaponPet) Initialize() {
	runeWeapon.dkOwner.registerDrwFrostFever()
	runeWeapon.dkOwner.registerDrwBloodPlague()
	runeWeapon.AddCopySpell(BloodBoilActionID, runeWeapon.dkOwner.registerDrwBloodBoil())
	runeWeapon.AddCopySpell(DeathCoilActionID, runeWeapon.dkOwner.registerDrwDeathCoil())
	runeWeapon.AddCopySpell(DeathStrikeActionID, runeWeapon.dkOwner.registerDrwDeathStrike())
	runeWeapon.AddCopySpell(IcyTouchActionID, runeWeapon.dkOwner.registerDrwIcyTouch())
	runeWeapon.AddCopySpell(OutbreakActionID, runeWeapon.dkOwner.registerDrwOutbreak())
	runeWeapon.AddCopySpell(PestilenceActionID, runeWeapon.dkOwner.registerDrwPestilence())
	runeWeapon.AddCopySpell(PlagueStrikeActionID, runeWeapon.dkOwner.registerDrwPlagueStrike())
	runeWeapon.AddCopySpell(SoulReaperActionID.WithTag(1), runeWeapon.dkOwner.registerDrwSoulReaper())

	runeWeapon.registerFirstHitDamageAura()
}

// The absolute first white hit doesn't get the -50% modifier.
// Trying to abuse it by macroing e.g. DS and DRW together doesn't seem to work.
func (runeWeapon *RuneWeaponPet) registerFirstHitDamageAura() {
	var damageMultiplier float64
	var firstHitAura *core.Aura
	firstHitAura = runeWeapon.RegisterAura(core.Aura{
		Label:    "First Hit Penalty Bypass" + runeWeapon.Label,
		Duration: core.SpellBatchWindow,

		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
	}).AttachProcTrigger(core.ProcTrigger{
		Callback: core.CallbackOnApplyEffects,
		ProcMask: core.ProcMaskWhiteHit,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			damageMultiplier = spell.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical]
			spell.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= 0.5
		},
	}).AttachProcTrigger(core.ProcTrigger{
		Callback: core.CallbackOnCastComplete,
		ProcMask: core.ProcMaskWhiteHit,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			spell.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] = damageMultiplier
			firstHitAura.Deactivate(sim)
		},
	})
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

	baseDamage := dk.CalcScalingSpellDmg(3.0)
	runeWeapon.EnableAutoAttacks(runeWeapon, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:        baseDamage,
			BaseDamageMax:        baseDamage,
			SwingSpeed:           3.5,
			NormalizedSwingSpeed: 3.5,
			CritMultiplier:       dk.DefaultCritMultiplier(),
			AttackPowerPerDPS:    core.DefaultAttackPowerPerDPS,
			MaxRange:             core.MaxMeleeRange,
		},
		AutoSwingMelee: true,
	})

	// Special weapon used for some strikes like DS and HS
	strikeWeaponBaseDamage := math.Floor(dk.CalcScalingSpellDmg(1.6))
	runeWeapon.StrikeWeapon = &core.Weapon{
		BaseDamageMin:        strikeWeaponBaseDamage,
		BaseDamageMax:        strikeWeaponBaseDamage,
		SwingSpeed:           3.5,
		NormalizedSwingSpeed: 3.5,
		CritMultiplier:       dk.DefaultCritMultiplier(),
		AttackPowerPerDPS:    core.DefaultAttackPowerPerDPS,
		MaxRange:             core.MaxMeleeRange,
	}
	runeWeapon.StrikeWeaponDamage = math.Floor(runeWeapon.StrikeWeapon.DPS()) * 3.5

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
	runeWeapon.drwDmgSnapshot = runeWeapon.dkOwner.PseudoStats.DamageDealtMultiplier
	runeWeapon.drwSchoolDmgSnapshot = runeWeapon.dkOwner.PseudoStats.SchoolDamageDealtMultiplier

	runeWeapon.PseudoStats.DamageDealtMultiplier *= runeWeapon.drwDmgSnapshot
	for i := range stats.SchoolLen {
		runeWeapon.PseudoStats.SchoolDamageDealtMultiplier[i] *= runeWeapon.drwSchoolDmgSnapshot[i]
		if i == stats.SchoolIndexPhysical {
			runeWeapon.PseudoStats.SchoolDamageDealtMultiplier[i] *= 0.5
		}
	}
}

func (runeWeapon *RuneWeaponPet) disable(sim *core.Simulation) {
	// Clear snapshot damage multipliers
	runeWeapon.PseudoStats.DamageDealtMultiplier /= runeWeapon.drwDmgSnapshot
	for i := range stats.SchoolLen {
		if i == stats.SchoolIndexPhysical {
			runeWeapon.PseudoStats.SchoolDamageDealtMultiplier[i] /= 0.5
		}
		runeWeapon.PseudoStats.SchoolDamageDealtMultiplier[i] /= runeWeapon.drwSchoolDmgSnapshot[i]
	}
	runeWeapon.drwSchoolDmgSnapshot = stats.NewSchoolFloatArray()
	runeWeapon.drwDmgSnapshot = 1
}
