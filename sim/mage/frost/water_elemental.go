package frost

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

// The numbers in this file are VERY rough approximations based on logs.

func (Mage *FrostMage) registerSummonWaterElementalSpell() {

	Mage.SummonWaterElemental = Mage.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 31687},
		Flags:    core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 16,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    Mage.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			Mage.waterElemental.Enable(sim, Mage.waterElemental)
		},
	})

	Mage.AddMajorCooldown(core.MajorCooldown{
		Spell:    Mage.SummonWaterElemental,
		Priority: core.CooldownPriorityDrums + 1000, // Always prefer to cast before drums or lust so the ele gets their benefits.
		Type:     core.CooldownTypeDPS,
	})
}

type WaterElemental struct {
	core.Pet

	mageOwner *FrostMage

	// Water Ele almost never just stands still and spams like we want, it sometimes
	// does its own thing. This controls how much it does that.
	disobeyChance float64

	Waterbolt *core.Spell
}

func (Mage *FrostMage) NewWaterElemental(disobeyChance float64) *WaterElemental {
	waterElemental := &WaterElemental{
		Pet: core.NewPet(core.PetConfig{
			Name:            "Water Elemental",
			Owner:           &Mage.Character,
			BaseStats:       waterElementalBaseStats,
			StatInheritance: waterElementalStatInheritance,
			EnabledOnStart:  true,
			IsGuardian:      true,
		}),
		mageOwner:     Mage,
		disobeyChance: disobeyChance,
	}
	waterElemental.EnableManaBarWithModifier(0.333)

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

	if sim.Proc(we.disobeyChance, "Disobey") {
		// Water ele has decided not to cooperate, so just wait for the cast time
		// instead of casting.
		we.WaitUntil(sim, sim.CurrentTime+spell.DefaultCast.CastTime)
		return
	}

	spell.Cast(sim, we.CurrentTarget)
}

// These numbers are just rough guesses based on looking at some logs.
var waterElementalBaseStats = stats.Stats{
	// TODO update. taken at level 80 on beta
	stats.Mana:      16123,
	stats.Intellect: 369, //unsure on beta
}

var waterElementalStatInheritance = func(ownerStats stats.Stats) stats.Stats {
	// These numbers are just rough guesses based on looking at some logs.
	return stats.Stats{
		// TODO Pet stats scale dynamically in combat
		stats.Stamina:    ownerStats[stats.Stamina] * 0.2,
		stats.Intellect:  ownerStats[stats.Intellect] * 0.3,
		stats.SpellPower: ownerStats[stats.SpellPower] * 0.333,

		// TODO test crit chance. It does crit, so figure out how often and if it scales
		/* Results: owner 5% crit, Waterbolt 13% crit
		owner 18% crit, waterbolt 18% crit
		*/
		stats.HasteRating:      ownerStats[stats.HasteRating],
		stats.SpellCritPercent: ownerStats[stats.SpellCritPercent],
	}
}

func (we *WaterElemental) registerWaterboltSpell() {
	we.Waterbolt = we.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 31707},
		SpellSchool: core.SpellSchoolFrost,
		ProcMask:    core.ProcMaskSpellDamage,

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
		BonusCoefficient: 0.833,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := .405 * we.mageOwner.ClassSpellScaling
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}
