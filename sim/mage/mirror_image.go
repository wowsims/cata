package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

/* Cata dummy testing

damage did NOT change dynamically (equipped staff midway and spells did same damage)
if a frost bolt is mid-air when mirror images expire, frostbolt does not land
*/

func (mage *Mage) registerMirrorImageCD() {
	spell := mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 55342},
		SpellSchool:    core.SpellSchoolArcane,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellMirrorImage,
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 10,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			mage.mirrorImage.EnableWithTimeout(sim, mage.mirrorImage, time.Second*30)
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}

type MirrorImage struct {
	core.Pet

	mageOwner *Mage

	Frostbolt   *core.Spell
	Fireblast   *core.Spell
	Fireball    *core.Spell
	ArcaneBlast *core.Spell

	hasGlyph bool
	Spec     proto.Spec
}

func (mage *Mage) NewMirrorImage() *MirrorImage {
	hasGlyph := mage.HasMinorGlyph(proto.MageMinorGlyph_GlyphOfMirrorImage)

	mirrorImage := &MirrorImage{
		Pet: core.NewPet(core.PetConfig{
			Name:            "Mirror Image",
			Owner:           &mage.Character,
			BaseStats:       mirrorImageBaseStats,
			StatInheritance: createMirrorImageInheritance(),
			EnabledOnStart:  false,
			IsGuardian:      true,
		}),
		mageOwner: mage,
		hasGlyph:  hasGlyph,
		Spec:      mage.Spec,
	}

	mirrorImage.EnableManaBar()

	mage.AddPet(mirrorImage)

	return mirrorImage
}

func (mi *MirrorImage) GetPet() *core.Pet {
	return &mi.Pet
}

func (mi *MirrorImage) Initialize() {
	mi.registerFireblastSpell()
	mi.registerFrostboltSpell()
	mi.registerArcaneBlastSpell()
	mi.registerFireballSpell()
}

func (mi *MirrorImage) Reset(_ *core.Simulation) {
}

func (mi *MirrorImage) ExecuteCustomRotation(sim *core.Simulation) {
	var spell *core.Spell
	//TODO implement glyph, where mirror images cast your spec's main filler

	// Arcane Spec & Glyphed
	if mi.Spec == 10 && mi.hasGlyph {
		spell = mi.ArcaneBlast
	} else if mi.Spec == 11 && mi.hasGlyph {
		// Fire Spec & Glyphed
		spell = mi.Fireball
	} else {
		// Frost spec or no glyph
		spell = mi.Frostbolt
		if mi.Fireblast.CD.IsReady(sim) {
			spell = mi.Fireblast
		}
	}

	if success := spell.Cast(sim, mi.CurrentTarget); !success {
		mi.Disable(sim)
	}
}

var mirrorImageBaseStats = stats.Stats{
	stats.Mana: 27020, // Confirmed via ingame bars at 80

	// seems to be about 8% baseline in wotlk
	stats.SpellCritPercent: 8,
}

var createMirrorImageInheritance = func() func(stats.Stats) stats.Stats {
	return func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			stats.SpellHitPercent: ownerStats[stats.SpellHitPercent],
			stats.SpellPower:      ownerStats[stats.SpellPower] * 0.33,
		}
	}
}

func (mi *MirrorImage) registerFrostboltSpell() {
	mi.Frostbolt = mi.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 59638}, // Confirmed via logs
		SpellSchool:  core.SpellSchoolFrost,
		ProcMask:     core.ProcMaskSpellDamage,
		MissileSpeed: 24,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 3,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   mi.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			//3x damage for 3 mirror images
			baseDamage := (220 + 0.25*spell.SpellPower()) * 3
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}

// *******************************************************
// If Fire spec with glyph, will chain cast Fireball
// *******************************************************
func (mi *MirrorImage) registerFireballSpell() {
	mi.Fireball = mi.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 88082}, // confirmed via logs
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
			//3x damage for 3 mirror images
			baseDamage := (317 + 0.338*spell.SpellPower()) * 3
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}

// *******************************************************
// If Arcane spec with glyph, will chain cast Arcane Blast
// *******************************************************
func (mi *MirrorImage) registerArcaneBlastSpell() {
	mi.ArcaneBlast = mi.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 88084}, //Confirmed via logs
		SpellSchool: core.SpellSchoolArcane,
		ProcMask:    core.ProcMaskSpellDamage,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 7,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 3,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   mi.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			//3x damage for 3 mirror images
			baseDamage := (257.76 + 0.275*spell.SpellPower()) * 3 //unsure how to get MI scaling, just used mage's # but can't call it certain
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}

func (mi *MirrorImage) registerFireblastSpell() {

	mi.Fireblast = mi.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 59637},
		SpellSchool: core.SpellSchoolFrost,
		ProcMask:    core.ProcMaskSpellDamage,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDMin,
			},
			CD: core.Cooldown{
				Timer:    mi.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   mi.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			//3x damage for 3 mirror images
			baseDamage := (88 + 0.15*spell.SpellPower()) * 3
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}
