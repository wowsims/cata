package frost

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func (Mage *FrostMage) registerSummonWaterElementalSpell() {

	Mage.SummonWaterElemental = Mage.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 31687},
		Flags:    core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 3,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 1500 * time.Millisecond,
			},
			CD: core.Cooldown{
				Timer:    Mage.NewTimer(),
				Duration: time.Minute * 1,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			Mage.waterElemental.Enable(sim, Mage.waterElemental)
		},
	})

	Mage.AddMajorCooldown(core.MajorCooldown{
		Spell: Mage.SummonWaterElemental,
		Type:  core.CooldownTypeDPS,
	})
}

type WaterElemental struct {
	core.Pet

	mageOwner *FrostMage

	Waterbolt               *core.Spell
	waterElementalDamageMod *core.SpellMod
}

func (Mage *FrostMage) NewWaterElemental() *WaterElemental {

	hasGlyph := 

	waterElemental := &WaterElemental{
		Pet: core.NewPet(core.PetConfig{
			Name:                           "Water Elemental",
			Owner:                          &Mage.Character,
			BaseStats:                      waterElementalBaseStats,
			StatInheritance:                waterElementalStatInheritance,
			HasDynamicCastSpeedInheritance: true,
			EnabledOnStart:                 true,
			IsGuardian:                     true,
		}),
		mageOwner: Mage,
	}
	waterElemental.EnableManaBar()

	Mage.AddPet(waterElemental)

	return waterElemental
}

func (we *WaterElemental) GetPet() *core.Pet {
	return &we.Pet
}

func (we *WaterElemental) Initialize() {
	we.registerWaterboltSpell()
}

func (we *WaterElemental) Reset(_ *core.Simulation) {
}

func (we *WaterElemental) ExecuteCustomRotation(sim *core.Simulation) {
	spell := we.Waterbolt
	spell.Cast(sim, we.CurrentTarget)
}

var waterElementalBaseStats = stats.Stats{
	// Mana seems to always be at 300k on beta
	stats.Mana: 300000,
}

var waterElementalStatInheritance = func(ownerStats stats.Stats) stats.Stats {
	// Water elemental usually has about half the HP of the caster, with glyph this is bumped by an additional 40%
	waterElementalStaminaRatio = 0.5 
	if Mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfWaterElemental) {
		waterElementalStaminaRatio *= 1.4
	}
	return stats.Stats{
		stats.Stamina:          ownerStats[stats.Stamina] * waterElementalStaminaRatio,
		stats.SpellPower:       ownerStats[stats.SpellPower],
		stats.HasteRating:      ownerStats[stats.HasteRating],
		stats.SpellCritPercent: ownerStats[stats.SpellCritPercent],
		// this (crit) needs to be tested more thoroughly when pet hit is not bugged
	}
}

var waterboltVariance = 0.25   // Per https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=31707 Field: "Variance"
var waterboltScale = 0.5       // Per https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=31707 Field: "Coefficient"
var waterboltCoefficient = 0.5 // Per https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=31707 Field: "BonusCoefficient"

func (we *WaterElemental) registerWaterboltSpell() {

	if Mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfWaterElemental) {
		we.AddStaticMod(core.SpellModConfig{
			Kind:      core.SpellMod_AllowCastWhileMoving,
			ClassMask: MageWaterElementalSpellWaterBolt,
		})
	}

	we.Waterbolt = we.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 31707},
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: MageWaterElementalSpellWaterBolt,

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
		CritMultiplier:   we.mageOwner.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: waterboltCoefficient,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := we.CalcAndRollDamageRange(sim, waterboltScale, waterboltVariance)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
