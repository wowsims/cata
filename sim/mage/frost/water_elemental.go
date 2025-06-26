package frost

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/mage"
)

func (mage *FrostMage) registerSummonWaterElementalSpell() {

	mage.SummonWaterElemental = mage.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 31687},
		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 3,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 1500 * time.Millisecond,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Minute * 1,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			mage.waterElemental.Enable(sim, mage.waterElemental)
		},
	})
}

type WaterElemental struct {
	core.Pet

	mageOwner *FrostMage

	Waterbolt *core.Spell
}

func (mage *FrostMage) NewWaterElemental() *WaterElemental {
	waterElementalStatInheritance := func(ownerStats stats.Stats) stats.Stats {
		// Water elemental usually has about half the HP of the caster, with glyph this is bumped by an additional 40%
		return stats.Stats{
			stats.Stamina:          ownerStats[stats.Stamina] * 0.5,
			stats.SpellPower:       ownerStats[stats.SpellPower],
			stats.HasteRating:      ownerStats[stats.HasteRating],
			stats.SpellCritPercent: ownerStats[stats.SpellCritPercent],
			// this (crit) needs to be tested more thoroughly when pet hit is not bugged
		}
	}

	waterElementalBaseStats := stats.Stats{
		// Mana seems to always be at 300k on beta
		stats.Mana: 300000,
	}

	waterElemental := &WaterElemental{
		Pet: core.NewPet(core.PetConfig{
			Name:                           "Water Elemental",
			Owner:                          &mage.Character,
			BaseStats:                      waterElementalBaseStats,
			NonHitExpStatInheritance:       waterElementalStatInheritance,
			HasDynamicCastSpeedInheritance: true,
			EnabledOnStart:                 true,
			IsGuardian:                     true,
		}),
		mageOwner: mage,
	}
	waterElemental.EnableManaBar()
	waterElemental.EnableDynamicStats(waterElementalStatInheritance)

	mage.AddPet(waterElemental)

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

func (we *WaterElemental) registerWaterboltSpell() {
	waterboltVariance := 0.25   // Per https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=31707 Field: "Variance"
	waterboltScale := 0.5       // Per https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=31707 Field: "Coefficient"
	waterboltCoefficient := 0.5 // Per https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=31707 Field: "BonusCoefficient"
	if we.mageOwner.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfWaterElemental) {
		we.AddStaticMod(core.SpellModConfig{
			Kind:      core.SpellMod_AllowCastWhileMoving,
			ClassMask: mage.MageWaterElementalSpellWaterBolt,
		})
	}

	we.Waterbolt = we.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 31707},
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: mage.MageWaterElementalSpellWaterBolt,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2500,
			},
		},

		DamageMultiplier: 1 * 1.2, // 2013-09-23 Ice Lance's damage has been increased by 20%
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
