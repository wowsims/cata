package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
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
		Pet: core.NewPet(core.PetConfig{
			Name:            "Mirror Image T12 2pc",
			Owner:           &mage.Character,
			BaseStats:       t12MirrorImageBaseStats,
			StatInheritance: createT12MirrorImageInheritance(),
			EnabledOnStart:  false,
			IsGuardian:      true,
		}),
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
	if mi.Fireball.CanCast(sim, mi.CurrentTarget) {
		mi.Fireball.Cast(sim, mi.CurrentTarget)
		minDelay := 680.0
		maxDelay := 790.0
		delayRange := maxDelay - minDelay
		// ~680-790ms delay between casts resulting in ~735 ms average
		// Research:https://docs.google.com/spreadsheets/d/e/2PACX-1vTvD34UWX5Eb9dIGmH7EPRQyuLdJDOpNR7_8cmZlWRZb1W7RlRE-y7ffSnvM55o_GZ5dPusxAW1STH3/pubchart?oid=96701738&format=image
		randomDelay := time.Duration(minDelay+delayRange*sim.RandomFloat("T12 Mirror Image Cast Delay")) * time.Millisecond
		mi.WaitUntil(sim, mi.NextGCDAt()+randomDelay)
		return
	}
}

var t12MirrorImageBaseStats = stats.Stats{
	stats.Mana:             27020, // Confirmed via ingame bars at 80
	stats.SpellCritPercent: 0,
}

func createT12MirrorImageInheritance() func(stats.Stats) stats.Stats {
	return func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			stats.SpellHitPercent: ownerStats[stats.SpellHitPercent],
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
			BaseCostPercent: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2500,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   mi.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.RollWithLabel(6490, 7543, "T12 Mirror Image Fireball")
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
