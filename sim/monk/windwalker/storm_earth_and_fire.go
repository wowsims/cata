package monk

import (
	"time"

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

func (ww *WindwalkerMonk) registerDancingRuneWeaponSpell() {
	duration := time.Second * 12

	sefAura := ww.RegisterAura(core.Aura{
		Label:    "Storm, Earth, and Fire",
		ActionID: core.ActionID{SpellID: 137639},
		Duration: core.NeverExpires,
		// Casts copy
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			copySpell := ww.RuneWeapon.runeWeaponSpells[spell.ActionID]
			if copySpell == nil {
				return
			}

			CopySpellMultipliers(spell, copySpell, ww.CurrentTarget)

			copySpell.Cast(sim, ww.CurrentTarget)
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {

		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {

		},
	})

	spell := ww.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 137639},
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MonkSpellStormEarthAndFire,

		EnergyCost: core.EnergyCostOptions{
			Cost: 10,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			ww.RuneWeapon.EnableWithTimeout(sim, ww.RuneWeapon, duration)
			ww.RuneWeapon.CancelGCDTimer(sim)
			dancingRuneWeaponAura.Activate(sim)
		},
	})

	ww.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}

type StormEarthAndFirePet struct {
	core.Pet

	dkOwner *DeathKnight

	// Diseases
	FrostFeverSpell  *core.Spell
	BloodPlagueSpell *core.Spell

	drwDmgSnapshot  float64
	drwPhysSnapshot float64

	runeWeaponSpells map[core.ActionID]*core.Spell
}

func (runeWeapon *StormEarthAndFirePet) Initialize() {
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

func (runeWeapon *StormEarthAndFirePet) AddCopySpell(actionId core.ActionID, spell *core.Spell) {
	runeWeapon.runeWeaponSpells[actionId] = spell
}

func (dk *WindwalkerMonk) NewRuneWeapon() *StormEarthAndFirePet {
	runeWeapon := &StormEarthAndFirePet{
		Pet: core.NewPet("Rune Weapon", &ww.Character, stats.Stats{
			stats.Stamina: 100,
		}, func(ownerStats stats.Stats) stats.Stats {
			return stats.Stats{
				stats.AttackPower: ownerStats[stats.AttackPower],
				stats.HasteRating: ownerStats[stats.HasteRating],

				stats.PhysicalHitPercent: ownerStats[stats.PhysicalHitPercent],
				stats.SpellHitPercent:    ownerStats[stats.PhysicalHitPercent] * HitCapRatio,

				stats.ExpertiseRating: ownerStats[stats.PhysicalHitPercent] * PetExpertiseRatingScale,

				stats.PhysicalCritPercent: ownerStats[stats.PhysicalCritPercent],
				stats.SpellCritPercent:    ownerStats[stats.SpellCritPercent],
			}
		}, false, true),
		dkOwner: dk,
	}

	runeWeapon.runeWeaponSpells = map[core.ActionID]*core.Spell{}

	runeWeapon.OnPetEnable = runeWeapon.enable
	runeWeapon.OnPetDisable = runeWeapon.disable

	mhWeapon := ww.WeaponFromMainHand(ww.DefaultCritMultiplier())

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

	ww.AddPet(runeWeapon)

	return runeWeapon
}

func (runeWeapon *StormEarthAndFirePet) GetPet() *core.Pet {
	return &runeWeapon.Pet
}

func (runeWeapon *StormEarthAndFirePet) Reset(_ *core.Simulation) {
}

func (runeWeapon *StormEarthAndFirePet) ExecuteCustomRotation(_ *core.Simulation) {
}

func (runeWeapon *StormEarthAndFirePet) enable(sim *core.Simulation) {
	// Snapshot extra % speed modifiers from dk owner
	runeWeapon.PseudoStats.MeleeSpeedMultiplier = 1
	runeWeapon.MultiplyMeleeSpeed(sim, runeWeapon.dkOwner.PseudoStats.MeleeSpeedMultiplier)

	runeWeapon.drwDmgSnapshot = runeWeapon.dkOwner.PseudoStats.DamageDealtMultiplier * 0.5
	runeWeapon.PseudoStats.DamageDealtMultiplier *= runeWeapon.drwDmgSnapshot

	runeWeapon.drwPhysSnapshot = runeWeapon.dkOwner.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical]
	runeWeapon.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= runeWeapon.drwPhysSnapshot

}

func (runeWeapon *StormEarthAndFirePet) disable(sim *core.Simulation) {
	// Clear snapshot speed
	runeWeapon.PseudoStats.MeleeSpeedMultiplier = 1
	runeWeapon.MultiplyMeleeSpeed(sim, 1)

	// Clear snapshot damage multipliers
	runeWeapon.PseudoStats.DamageDealtMultiplier /= runeWeapon.drwDmgSnapshot
	runeWeapon.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= runeWeapon.drwPhysSnapshot
	runeWeapon.drwPhysSnapshot = 1
	runeWeapon.drwDmgSnapshot = 1
}
