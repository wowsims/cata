package shaman

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

type FireElemental struct {
	core.Pet

	FireBlast *core.Spell
	FireNova  *core.Spell
	Immolate  *core.Spell
	Empower   *core.Spell

	shamanOwner *Shaman

	fireBlastAutocast bool
	fireNovaAutocast  bool
	immolateAutocast  bool
	empowerAutocast   bool
}

var FireElementalSpellPowerScaling = 0.36

func (shaman *Shaman) NewFireElemental(isGuardian bool) *FireElemental {
	fireElemental := &FireElemental{
		Pet: core.NewPet(core.PetConfig{
			Name:                            core.Ternary(isGuardian, "Greater Fire Elemental", "Primal Fire Elemental"),
			Owner:                           &shaman.Character,
			BaseStats:                       shaman.fireElementalBaseStats(isGuardian),
			StatInheritance:                 shaman.fireElementalStatInheritance(isGuardian),
			EnabledOnStart:                  false,
			IsGuardian:                      isGuardian,
			HasDynamicCastSpeedInheritance:  true,
			HasDynamicMeleeSpeedInheritance: true,
		}),
		shamanOwner:       shaman,
		fireBlastAutocast: shaman.FeleAutocast.AutocastFireblast || isGuardian,
		fireNovaAutocast:  shaman.FeleAutocast.AutocastFirenova || isGuardian,
		immolateAutocast:  shaman.FeleAutocast.AutocastImmolate,
		empowerAutocast:   shaman.FeleAutocast.AutocastEmpower,
	}
	scalingDamage := shaman.CalcScalingSpellDmg(1.0)
	baseMeleeDamage := core.TernaryFloat64(isGuardian, scalingDamage, scalingDamage*1.8)
	fireElemental.EnableManaBar()
	fireElemental.EnableAutoAttacks(fireElemental, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  baseMeleeDamage,
			BaseDamageMax:  baseMeleeDamage,
			SwingSpeed:     1.4,
			CritMultiplier: fireElemental.DefaultCritMultiplier(),
			SpellSchool:    core.SpellSchoolFire,
		},
		AutoSwingMelee: true,
	})
	fireElemental.AutoAttacks.MHConfig().ProcMask |= core.ProcMaskSpellDamage
	fireElemental.AutoAttacks.MHConfig().ClassSpellMask |= SpellMaskFireElementalMelee

	fireElemental.OnPetEnable = fireElemental.enable(isGuardian)
	fireElemental.OnPetDisable = fireElemental.disable

	shaman.AddPet(fireElemental)

	return fireElemental
}

func (fireElemental *FireElemental) enable(isGuardian bool) func(*core.Simulation) {
	return func(sim *core.Simulation) {
		if fireElemental.empowerAutocast {
			if fireElemental.Empower.Cast(sim, &fireElemental.shamanOwner.Unit) {
				fireElemental.AutoAttacks.StopMeleeUntil(sim, fireElemental.Empower.Hot(&fireElemental.shamanOwner.Unit).ExpiresAt(), false)
			}
		}
	}
}

func (fireElemental *FireElemental) disable(sim *core.Simulation) {
	fireElemental.Empower.Hot(&fireElemental.shamanOwner.Unit).Deactivate(sim)
}

func (fireElemental *FireElemental) GetPet() *core.Pet {
	return &fireElemental.Pet
}

func (fireElemental *FireElemental) Initialize() {

	fireElemental.registerFireBlast()
	fireElemental.registerFireNova()
	fireElemental.registerImmolate()
	fireElemental.registerEmpower()
}

func (fireElemental *FireElemental) Reset(_ *core.Simulation) {

}

func (fireElemental *FireElemental) ExecuteCustomRotation(sim *core.Simulation) {
	/*
		Fire Blast on CD, Fire nova on CD when 2+ targets, Immolate on CD if not up on a target
	*/
	target := fireElemental.CurrentTarget

	if fireElemental.immolateAutocast {
		for _, target := range sim.Encounter.TargetUnits {
			if fireElemental.Immolate.Dot(target).RemainingDuration(sim) < fireElemental.Immolate.Dot(target).TickPeriod() && fireElemental.TryCast(sim, target, fireElemental.Immolate) {
				break
			}
		}
	}
	if fireElemental.fireNovaAutocast && len(sim.Encounter.TargetUnits) > 2 {
		fireElemental.TryCast(sim, target, fireElemental.FireNova)
	}
	if fireElemental.fireBlastAutocast {
		fireElemental.FireBlast.Cast(sim, target)
	}

	if !fireElemental.GCD.IsReady(sim) {
		return
	}

	minCd := min(fireElemental.FireBlast.CD.ReadyAt(), fireElemental.FireNova.CD.ReadyAt(), fireElemental.Immolate.CD.ReadyAt())
	fireElemental.ExtendGCDUntil(sim, max(minCd, sim.CurrentTime+time.Second))

}

func (fireElemental *FireElemental) TryCast(sim *core.Simulation, target *core.Unit, spell *core.Spell) bool {
	if !spell.Cast(sim, target) {
		return false
	}
	// all spell casts reset the elemental's swing timer
	fireElemental.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime+spell.CurCast.CastTime, false)
	return true
}

func (shaman *Shaman) fireElementalBaseStats(isGuardian bool) stats.Stats {
	return stats.Stats{
		stats.Mana:    9916,
		stats.Stamina: core.TernaryFloat64(isGuardian, 7843, 7843*1.2),
	}
}

func (shaman *Shaman) fireElementalStatInheritance(isGuardian bool) core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		ownerHitRating := ownerStats[stats.HitRating]
		ownerExpertiseRating := ownerStats[stats.ExpertiseRating]
		ownerSpellCritPercent := ownerStats[stats.SpellCritPercent]
		ownerPhysicalCritPercent := ownerStats[stats.PhysicalCritPercent]
		ownerHasteRating := ownerStats[stats.HasteRating]
		hitExpRating := (ownerHitRating + ownerExpertiseRating) / 2
		critPercent := core.TernaryFloat64(math.Abs(ownerPhysicalCritPercent) > math.Abs(ownerSpellCritPercent), ownerPhysicalCritPercent, ownerSpellCritPercent)

		power := core.TernaryFloat64(shaman.Spec == proto.Spec_SpecEnhancementShaman, ownerStats[stats.AttackPower]*0.65, ownerStats[stats.SpellPower])

		return stats.Stats{
			stats.Stamina:    ownerStats[stats.Stamina] * core.TernaryFloat64(isGuardian, 0.75, 0.75*1.2),
			stats.SpellPower: power * core.TernaryFloat64(isGuardian, FireElementalSpellPowerScaling, FireElementalSpellPowerScaling*1.8),

			stats.HitRating:           hitExpRating,
			stats.ExpertiseRating:     hitExpRating,
			stats.SpellCritPercent:    critPercent,
			stats.PhysicalCritPercent: critPercent,
			stats.HasteRating:         ownerHasteRating,
		}
	}
}
