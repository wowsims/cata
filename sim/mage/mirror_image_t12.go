package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

/*
	From Wowhead comments, this miror image is distinctly different from the base images.
	It has less health and hits harder.
	1. Casts fireballs that hit for ~6.5 - 7.5k each. The MI will cast 4 of these before de-spawning. Lasts for ~15 seconds.
		- Another comment suggests 6490 to 7543
	2. Has only ~3.5k HP but is immune to many forms of PVE AoE.
	3. Has an ICD of approximately 45 seconds, you will never get more than 1 at a time.
	4. Does not appear to be affected by the base stats of the mage aside from hit (if you are not hit capped the MI can miss). I tested this by removing several pieces of gear and examining the MI's DPS.
	5. Does not follow your target if you switch targets.
	6. Is not affected by buffs or CDs.
*/

type T12MirrorImage struct {
	core.Pet

	mageOwner *Mage

	Fireball *core.Spell
}

func (mage *Mage) NewT12MirrorImage() *T12MirrorImage {
	mirrorImage := &T12MirrorImage{
		Pet:       core.NewPet("Mirror Image T12 2pc", &mage.Character, t12MirrorImageBaseStats, createT12MirrorImageInheritance(), false, true),
		mageOwner: mage,
	}

	mirrorImage.EnableManaBar()

	mage.AddPet(mirrorImage)

	return mirrorImage
}

func (mi *T12MirrorImage) GetPet() *core.Pet {
	return &mi.Pet
}

func (mi *T12MirrorImage) Initialize() {
	mi.registerFireballSpell()
}

func (mi *T12MirrorImage) Reset(_ *core.Simulation) {
}

func (mi *T12MirrorImage) ExecuteCustomRotation(sim *core.Simulation) {
	if success := mi.Fireball.Cast(sim, mi.CurrentTarget); !success {
		mi.Disable(sim)
	}
}

var t12MirrorImageBaseStats = stats.Stats{
	stats.Mana: 27020, // Confirmed via ingame bars at 80

	// seems to be about 8% baseline in wotlk
	stats.SpellCritPercent: 8,
}

func createT12MirrorImageInheritance() func(stats.Stats) stats.Stats {
	return func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			// According to the simcraft implementation, the T12 Mirror Images snapshot the owner's crit
			stats.SpellCritPercent: ownerStats[stats.SpellCritPercent],
			stats.SpellHitPercent:  ownerStats[stats.SpellHitPercent],
		}
	}
}

func (mi *T12MirrorImage) registerFireballSpell() {
	mi.Fireball = mi.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 99062},
		SpellSchool:  core.SpellSchoolFire,
		ProcMask:     core.ProcMaskSpellDamage,
		MissileSpeed: 24,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.01,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2500,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   mi.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(6490, 7543)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
