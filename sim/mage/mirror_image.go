package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

/* MOP dummy testing

damage DOES change dynamically (equipped staff midway and spells did more damage, on the next cast)
*/

func (mage *Mage) registerMirrorImageCD() {
	mage.SummonMirrorImages = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 55342},
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellMirrorImage,
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 2,
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
			for _, mirrorImage := range mage.mirrorImages {
				mirrorImage.EnableWithTimeout(sim, mirrorImage, time.Second*30)
			}
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: mage.SummonMirrorImages,
		Type:  core.CooldownTypeDPS,
	})
}

type MirrorImage struct {
	core.Pet

	mageOwner *Mage

	mainSpell   *core.Spell // Spell that mirror images actually use.
	Frostbolt   *core.Spell
	Fireblast   *core.Spell
	Fireball    *core.Spell
	ArcaneBlast *core.Spell

	arcaneChargesAura *core.Aura

	hasGlyph bool
}

func (mage *Mage) NewMirrorImage() *MirrorImage {
	hasGlyph := mage.HasMinorGlyph(proto.MageMinorGlyph_GlyphOfMirrorImage)

	mirrorImageStatInheritance := func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			stats.Stamina:          ownerStats[stats.Stamina],
			stats.SpellPower:       ownerStats[stats.SpellPower] * 0.05,
			stats.HasteRating:      ownerStats[stats.HasteRating],
			stats.SpellCritPercent: ownerStats[stats.SpellCritPercent],
		}
	}

	mirrorImageBaseStats := stats.Stats{
		stats.Mana: 27020, // Confirmed via ingame bars at 80
	}

	mirrorImage := &MirrorImage{
		Pet: core.NewPet(core.PetConfig{
			Name:                     "Mirror Image",
			Owner:                    &mage.Character,
			BaseStats:                mirrorImageBaseStats,
			NonHitExpStatInheritance: mirrorImageStatInheritance,
			EnabledOnStart:           false,
			IsGuardian:               true,
		}),
		mageOwner: mage,
		hasGlyph:  hasGlyph,
	}

	mirrorImage.EnableManaBar()
	mirrorImage.EnableDynamicStats(mirrorImageStatInheritance)

	mage.AddPet(mirrorImage)

	return mirrorImage
}

func (mi *MirrorImage) GetPet() *core.Pet {
	return &mi.Pet
}

func (mi *MirrorImage) Initialize() {
	mi.registerFrostboltSpell()
	mi.registerArcaneBlastSpell()
	mi.registerFireballSpell()

	mi.mainSpell = mi.Frostbolt
	if mi.hasGlyph {
		if mi.mageOwner.Spec == proto.Spec_SpecArcaneMage {
			mi.mainSpell = mi.ArcaneBlast
		} else if mi.mageOwner.Spec == proto.Spec_SpecFireMage {
			mi.mainSpell = mi.Fireball
		}
	}

}

func (mi *MirrorImage) Reset(_ *core.Simulation) {
}

func (mi *MirrorImage) ExecuteCustomRotation(sim *core.Simulation) {
	mi.mainSpell.Cast(sim, mi.CurrentTarget)
}

func (mi *MirrorImage) registerFrostboltSpell() {

	frostBoltCoefficient := 1.65
	frostBoltScaling := 1.65
	frostBoltVariance := 0.1

	mi.Frostbolt = mi.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 59638},
		SpellSchool:  core.SpellSchoolFrost,
		ProcMask:     core.ProcMaskSpellDamage,
		MissileSpeed: 24,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: .1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 2,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   mi.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: frostBoltCoefficient,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := mi.Owner.CalcAndRollDamageRange(sim, frostBoltScaling, frostBoltVariance)
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

	fireBallCoefficient := 1.8
	fireBallScaling := 1.8
	fireBallVariance := 0.2

	mi.Fireball = mi.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 88082}, // confirmed via logs
		SpellSchool:  core.SpellSchoolFire,
		ProcMask:     core.ProcMaskSpellDamage,
		MissileSpeed: 24,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: .1,
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
		BonusCoefficient: fireBallCoefficient,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := mi.Owner.CalcAndRollDamageRange(sim, fireBallScaling, fireBallVariance)
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

	arcaneBlastCoefficient := .9
	arcaneBlastScaling := .9
	arcaneBlastVariance := 0.15

	mi.ArcaneBlast = mi.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 88084}, //Confirmed via logs
		SpellSchool:    core.SpellSchoolArcane,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: MageMirrorImageSpellArcaneBlast,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: .1,
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
		BonusCoefficient: arcaneBlastCoefficient,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := mi.Owner.CalcAndRollDamageRange(sim, arcaneBlastScaling, arcaneBlastVariance)
			result := spell.CalcAndDealDamage(sim, mi.CurrentTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				mi.arcaneChargesAura.Activate(sim)
				mi.arcaneChargesAura.AddStack(sim)
			}
		},
	})

	abDamageMod := mi.AddDynamicMod(core.SpellModConfig{
		ClassMask:  MageMirrorImageSpellArcaneBlast,
		FloatValue: .5 * mi.mageOwner.T15_4PC_ArcaneChargeEffect,
		Kind:       core.SpellMod_DamageDone_Flat,
	})
	abCostMod := mi.AddDynamicMod(core.SpellModConfig{
		ClassMask:  MageMirrorImageSpellArcaneBlast,
		FloatValue: 1.5 * mi.mageOwner.T15_4PC_ArcaneChargeEffect,
		Kind:       core.SpellMod_PowerCost_Pct,
	})

	mi.arcaneChargesAura = mi.GetOrRegisterAura(core.Aura{
		Label:     "Mirror Images: Arcane Charges Aura",
		ActionID:  core.ActionID{SpellID: 36032}, //idk if it gets its own
		Duration:  time.Second * 10,
		MaxStacks: 4,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			abDamageMod.Activate()
			abCostMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			abDamageMod.Deactivate()
			abCostMod.Deactivate()
		},
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			stacks := float64(newStacks)
			abDamageMod.UpdateFloatValue(0.5 * stacks * mi.mageOwner.T15_4PC_ArcaneChargeEffect)
			abCostMod.UpdateFloatValue(1.5 * stacks * mi.mageOwner.T15_4PC_ArcaneChargeEffect)
		},
	})
}
