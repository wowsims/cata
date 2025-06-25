package enhancement

import (
	"math"
	"strconv"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

type SpiritWolf struct {
	core.Pet

	shamanOwner *EnhancementShaman

	SpiritBite *core.Spell
	enabledAt  time.Duration
}

type SpiritWolves struct {
	SpiritWolf1 *SpiritWolf
	SpiritWolf2 *SpiritWolf
}

func (SpiritWolves *SpiritWolves) EnableWithTimeout(sim *core.Simulation) {
	SpiritWolves.SpiritWolf1.enabledAt = sim.CurrentTime
	SpiritWolves.SpiritWolf2.enabledAt = sim.CurrentTime
	SpiritWolves.SpiritWolf1.EnableWithTimeout(sim, SpiritWolves.SpiritWolf1, time.Second*30)
	SpiritWolves.SpiritWolf2.EnableWithTimeout(sim, SpiritWolves.SpiritWolf2, time.Second*30)
}

func (SpiritWolves *SpiritWolves) CancelGCDTimer(sim *core.Simulation) {
	SpiritWolves.SpiritWolf1.CancelGCDTimer(sim)
	SpiritWolves.SpiritWolf2.CancelGCDTimer(sim)
}

var spiritWolfBaseStats = stats.Stats{
	stats.Stamina: 3137,
}

func (enh *EnhancementShaman) NewSpiritWolf(index int) *SpiritWolf {
	spiritWolf := &SpiritWolf{
		Pet: core.NewPet(core.PetConfig{
			Name:                            "Spirit Wolf " + strconv.Itoa(index),
			Owner:                           &enh.Character,
			BaseStats:                       spiritWolfBaseStats,
			StatInheritance:                 enh.makeStatInheritance(),
			EnabledOnStart:                  false,
			IsGuardian:                      true,
			HasDynamicMeleeSpeedInheritance: true,
			HasDynamicCastSpeedInheritance:  true,
		}),
		shamanOwner: enh,
	}
	baseMeleeDamage := enh.CalcScalingSpellDmg(0.5)
	spiritWolf.EnableAutoAttacks(spiritWolf, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  baseMeleeDamage,
			BaseDamageMax:  baseMeleeDamage,
			SwingSpeed:     1.5,
			CritMultiplier: spiritWolf.DefaultCritMultiplier(),
		},
		AutoSwingMelee: true,
	})

	spiritWolf.OnPetEnable = func(sim *core.Simulation) {
	}

	enh.AddPet(spiritWolf)

	return spiritWolf
}

func (enh *EnhancementShaman) makeStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		ownerHitRating := ownerStats[stats.HitRating]
		ownerExpertiseRating := ownerStats[stats.ExpertiseRating]
		ownerSpellCritPercent := ownerStats[stats.SpellCritPercent]
		ownerPhysicalCritPercent := ownerStats[stats.PhysicalCritPercent]
		ownerHasteRating := ownerStats[stats.HasteRating]
		hitExpRating := (ownerHitRating + ownerExpertiseRating) / 2
		critPercent := core.TernaryFloat64(math.Abs(ownerPhysicalCritPercent) > math.Abs(ownerSpellCritPercent), ownerPhysicalCritPercent, ownerSpellCritPercent)

		return stats.Stats{
			stats.Stamina:             ownerStats[stats.Stamina] * 0.3,
			stats.AttackPower:         ownerStats[stats.AttackPower] * 0.5,
			stats.HitRating:           hitExpRating,
			stats.ExpertiseRating:     hitExpRating,
			stats.PhysicalCritPercent: critPercent,
			stats.SpellCritPercent:    critPercent,
			stats.HasteRating:         ownerHasteRating,
		}
	}
}

func (spiritWolf *SpiritWolf) Initialize() {
	spiritWolf.registerSpiritBite()
}

func (spiritWolf *SpiritWolf) ExecuteCustomRotation(sim *core.Simulation) {
	/*
		Spirit Bite on Cd starting 3.3s in
	*/
	target := spiritWolf.CurrentTarget

	if sim.CurrentTime >= spiritWolf.enabledAt+time.Millisecond*3300 {
		spiritWolf.SpiritBite.Cast(sim, target)
	}

	if !spiritWolf.GCD.IsReady(sim) {
		return
	}

	minCd := spiritWolf.SpiritBite.CD.ReadyAt()
	spiritWolf.ExtendGCDUntil(sim, max(minCd, sim.CurrentTime+time.Second))

}

func (spiritWolf *SpiritWolf) Reset(sim *core.Simulation) {
	spiritWolf.Disable(sim)
	if sim.Log != nil {
		spiritWolf.Log(sim, "Base Stats: %s", spiritWolfBaseStats)
		inheritedStats := spiritWolf.shamanOwner.makeStatInheritance()(spiritWolf.shamanOwner.GetStats())
		spiritWolf.Log(sim, "Inherited Stats: %s", inheritedStats)
		spiritWolf.Log(sim, "Total Stats: %s", spiritWolf.GetStats())
	}
}

func (spiritWolf *SpiritWolf) GetPet() *core.Pet {
	return &spiritWolf.Pet
}

func (spiritWolf *SpiritWolf) registerSpiritBite() {
	spiritWolf.SpiritBite = spiritWolf.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 58859},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellDamage,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    spiritWolf.NewTimer(),
				Duration: time.Millisecond * 7300,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   spiritWolf.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damageScaling := spiritWolf.shamanOwner.CalcAndRollDamageRange(sim, 1.05, 0.40000000596)
			baseDamage := damageScaling + spell.MeleeAttackPower()*0.30000001192
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}
