package frost

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/mage"
)

func (frost *FrostMage) registerFrozenOrbSpell() {

	frost.frozenOrb = frost.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 84714},
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: mage.MageSpellFrozenOrb,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 6,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    frost.NewTimer(),
				Duration: time.Minute,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			frost.frozenOrb.EnableWithTimeout(sim, frost.frozenOrb, time.Second*10)
		},
	})

	frost.AddMajorCooldown(core.MajorCooldown{
		Spell: frost.frozenOrb,
		Type:  core.CooldownTypeDPS,
	})
}

type FrozenOrb struct {
	core.Pet

	mageOwner *FrostMage

	FrozenOrbTick *core.Spell

	TickCount int64
}

func (frost *FrostMage) NewFrozenOrb() *FrozenOrb {
	frozenOrb := &FrozenOrb{
		Pet: core.NewPet(core.PetConfig{
			Name:            "Frozen Orb",
			Owner:           &frost.Character,
			BaseStats:       frozenOrbBaseStats,
			StatInheritance: createFrozenOrbInheritance(),
			EnabledOnStart:  false,
			IsGuardian:      true,
		}),
		mageOwner: frost,
		TickCount: 0,
	}

	frozenOrb.Pet.OnPetEnable = frozenOrb.enable

	frost.AddPet(frozenOrb)

	return frozenOrb
}

func (frozenOrb *FrozenOrb) enable(sim *core.Simulation) {

	frozenOrb.EnableDynamicStats(func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			stats.SpellPower: ownerStats[stats.SpellPower],
		}
	})
}

func (frozenOrb *FrozenOrb) GetPet() *core.Pet {
	return &frozenOrb.Pet
}

func (frozenOrb *FrozenOrb) Initialize() {
	frozenOrb.registerFrozenOrbTickSpell()
}

func (frozenOrb *FrozenOrb) Reset(_ *core.Simulation) {
	frozenOrb.TickCount = 0
}

func (frozenOrb *FrozenOrb) ExecuteCustomRotation(sim *core.Simulation) {
	spell := frozenOrb.FrozenOrbTick
	if success := spell.Cast(sim, frozenOrb.CurrentTarget); !success {
		frozenOrb.Disable(sim)
	}
}

var frozenOrbBaseStats = stats.Stats{}

var createFrozenOrbInheritance = func() func(stats.Stats) stats.Stats {
	return func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			stats.SpellHitPercent:  ownerStats[stats.SpellHitPercent],
			stats.SpellCritPercent: ownerStats[stats.SpellCritPercent],
			stats.SpellPower:       ownerStats[stats.SpellPower],
		}
	}
}

// Values found at https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=84721
var frozenOrbCoefficient = 0.65
var frozenOrbScaling = 0.51
var frozenOrbVariance = 0.25

func (frozenOrb *FrozenOrb) registerFrozenOrbTickSpell() {
	frozenOrb.FrozenOrbTick = frozenOrb.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 84721},
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: mage.MageSpellFrozenOrb,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   frozenOrb.mageOwner.DefaultCritMultiplier(),
		BonusCoefficient: frozenOrbCoefficient,
		ThreatMultiplier: 1,
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return frozenOrb.TickCount < 10
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {

			damage := frozenOrb.mageOwner.CalcAndRollDamageRange(sim, frozenOrbScaling, frozenOrbVariance)
			spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)
			frozenOrb.TickCount += 1
		},
	})
}
