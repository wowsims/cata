package frost

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/mage"
)

func (frost *FrostMage) registerFrozenOrbSpell() {

	frozenOrb := frost.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 84714},
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: mage.MageSpellFrozenOrb,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 10,
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
			// Frozen Orb gives a stack of FoF upon reaching an enemy. It does it a good bit before actually hitting testing on beta, so instant I think is fine.
			frost.FingersOfFrostAura.Activate(sim)
			frost.FingersOfFrostAura.AddStack(sim)
			frost.frozenOrb.EnableWithTimeout(sim, frost.frozenOrb, time.Second*10)
		},
	})

	frost.AddMajorCooldown(core.MajorCooldown{
		Spell: frozenOrb,
		Type:  core.CooldownTypeDPS,
	})
}

type FrozenOrb struct {
	core.Pet

	mageOwner *FrostMage

	FrozenOrbTick *core.Spell

	TickCount int64
}

func createFrozenOrbInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			stats.SpellCritPercent: ownerStats[stats.SpellCritPercent],
			stats.SpellPower:       ownerStats[stats.SpellPower],
		}
	}
}

func (frost *FrostMage) NewFrozenOrb() *FrozenOrb {

	frozenOrbBaseStats := stats.Stats{}
	frozenOrb := &FrozenOrb{
		Pet: core.NewPet(core.PetConfig{
			Name:                     "Frozen Orb",
			Owner:                    &frost.Character,
			BaseStats:                frozenOrbBaseStats,
			NonHitExpStatInheritance: createFrozenOrbInheritance(),
			EnabledOnStart:           false,
			IsGuardian:               false,
		}),
		mageOwner: frost,
		TickCount: 0,
	}

	frozenOrb.Pet.OnPetEnable = frozenOrb.enable

	frost.AddPet(frozenOrb)

	return frozenOrb
}

func (frozenOrb *FrozenOrb) enable(sim *core.Simulation) {
	frozenOrb.TickCount = 0
	frozenOrb.EnableDynamicStats(createFrozenOrbInheritance())
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
	if frozenOrb.FrozenOrbTick.CanCast(sim, frozenOrb.CurrentTarget) {
		frozenOrb.FrozenOrbTick.Cast(sim, frozenOrb.CurrentTarget)
	}
}

func (frozenOrb *FrozenOrb) registerFrozenOrbTickSpell() {
	// Values found at https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=84721
	frozenOrbCoefficient := 0.65
	frozenOrbScaling := 0.51
	frozenOrbVariance := 0.25
	frozenOrb.FrozenOrbTick = frozenOrb.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 84721},
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: mage.MageSpellFrozenOrbTick,
		Flags:          core.SpellFlagAoE,

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
			anyLanded := false
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				result := spell.CalcAndDealDamage(sim, aoeTarget, damage, spell.OutcomeMagicHitAndCrit)
				if !anyLanded && result.Landed() {
					anyLanded = true
				}
			}
			if anyLanded && sim.Proc(0.15, "FingersOfFrostProc") {
				frozenOrb.mageOwner.FingersOfFrostAura.Activate(sim)
				frozenOrb.mageOwner.FingersOfFrostAura.AddStack(sim)
			}
			frozenOrb.TickCount += 1

			if frozenOrb.TickCount >= 10 {
				frozenOrb.Disable(sim)
			}
		},
	})
}
