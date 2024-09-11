package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

// T11
var ItemSetMaleficRaiment = core.NewItemSet(core.ItemSet{
	Name: "Shadowflame Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			agent.(WarlockAgent).GetWarlock().AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_CastTime_Pct,
				ClassMask:  WarlockSpellChaosBolt | WarlockSpellHandOfGuldan | WarlockSpellHaunt,
				FloatValue: -0.1,
			})
		},
		4: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()

			dmgMod := warlock.AddDynamicMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Flat,
				ClassMask:  WarlockSpellFelFlame,
				FloatValue: 3.0,
			})

			aura := warlock.RegisterAura(core.Aura{
				Label:     "Fel Spark",
				ActionID:  core.ActionID{SpellID: 89937},
				Duration:  15 * time.Second,
				MaxStacks: 2,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					dmgMod.Activate()
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					dmgMod.Deactivate()
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.Matches(WarlockSpellFelFlame) && result.Landed() {
						aura.RemoveStack(sim)
					}
				},
			})

			core.MakePermanent(warlock.RegisterAura(core.Aura{
				Label:           "Item - Warlock T11 4P Bonus",
				ActionID:        core.ActionID{SpellID: 89935},
				ActionIDForProc: aura.ActionID,
				OnPeriodicDamageDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					if spell.Matches(WarlockSpellImmolateDot|WarlockSpellUnstableAffliction) &&
						sim.Proc(0.02, "Warlock 4pT11") {
						aura.Activate(sim)
						aura.SetStacks(sim, 2)
					}
				},
			}))
		},
	},
})

// T12
type FieryImpPet struct {
	core.Pet

	FlameBlast *core.Spell
}

func (warlock *Warlock) NewFieryImp() *FieryImpPet {
	baseStats := stats.Stats{stats.SpellCritPercent: 0} // rough guess, seems to only get crit from debuffs?

	statInheritance := func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			// unclear what exactly the scaling is here, but seems to not miss(?)
			stats.HitRating: ownerStats[stats.SpellHitPercent] * core.SpellHitRatingPerHitPercent,
		}
	}

	imp := &FieryImpPet{Pet: core.NewPet("Fiery Imp", &warlock.Character, baseStats, statInheritance, false, true)}

	warlock.AddPet(imp)
	imp.registerFlameBlast(warlock)

	return imp
}

func (imp *FieryImpPet) GetPet() *core.Pet { return &imp.Pet }

func (imp *FieryImpPet) Initialize() {}

func (imp *FieryImpPet) Reset(_ *core.Simulation) {}

func (imp *FieryImpPet) ExecuteCustomRotation(sim *core.Simulation) {
	if imp.FlameBlast.CanCast(sim, imp.CurrentTarget) {
		imp.FlameBlast.Cast(sim, imp.CurrentTarget)
		delay := time.Duration(sim.RollWithLabel(150.0, 750.0, "Imp Cast Delay")) * time.Millisecond
		imp.WaitUntil(sim, imp.NextGCDAt()+delay)
		return
	}
}

func (pet *FieryImpPet) registerFlameBlast(warlock *Warlock) {
	pet.FlameBlast = pet.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 99226},
		SpellSchool:  core.SpellSchoolFire,
		ProcMask:     core.ProcMaskSpellDamage,
		MissileSpeed: 16,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 1500 * time.Millisecond,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           2,
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDmg := sim.Roll(2750+1, 2750+1514)
			result := spell.CalcDamage(sim, target, baseDmg, spell.OutcomeMagicHitAndCrit)

			// TODO: old wowhead comments seem to suggest that this spell does trigger burning embers, but does it do it
			// the same way that fire bolt does?
			if warlock.Talents.BurningEmbers > 0 && result.Landed() {
				dot := warlock.BurningEmbers.Dot(result.Target)
				dot.SnapshotBaseDamage += result.Damage * 0.25 * float64(warlock.Talents.BurningEmbers)
				dot.Apply(sim)
			}

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}

var ItemSetBalespidersBurningVestments = core.NewItemSet(core.ItemSet{
	Name: "Balespider's Burning Vestments",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()

			core.MakePermanent(warlock.RegisterAura(core.Aura{
				Label:    "Item - Warlock T12 2P Bonus",
				ActionID: core.ActionID{SpellID: 99220},
				Icd: &core.Cooldown{
					Timer:    warlock.NewTimer(),
					Duration: 45 * time.Second,
				},
				OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					if aura.Icd.IsReady(sim) && sim.Proc(0.05, "Warlock 2pT12") {
						warlock.FieryImp.EnableWithTimeout(sim, warlock.FieryImp, 15*time.Second)
						aura.Icd.Use(sim)
					}
				},
			}))
		},
		4: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()

			dmgMod := warlock.AddDynamicMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				School:     core.SpellSchoolShadow | core.SpellSchoolFire,
				FloatValue: 0.20,
			})

			aura := warlock.RegisterAura(core.Aura{
				Label:    "Apocalypse",
				ActionID: core.ActionID{SpellID: 99232},
				Duration: 8 * time.Second,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					dmgMod.Activate()
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					dmgMod.Deactivate()
				},
			})

			core.MakePermanent(warlock.RegisterAura(core.Aura{
				Label:           "Item - Warlock T12 4P Bonus",
				ActionID:        core.ActionID{SpellID: 99229},
				ActionIDForProc: aura.ActionID,
				OnCastComplete: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell.Matches(WarlockSpellShadowBolt|WarlockSpellIncinerate|WarlockSpellSoulFire|WarlockSpellDrainSoul) &&
						sim.Proc(0.05, "Warlock 4pT12") {
						aura.Activate(sim)
					}
				},
			}))
		},
	},
})

var ItemSetGladiatorsFelshroud = core.NewItemSet(core.ItemSet{
	ID:   910,
	Name: "Gladiator's Felshroud",

	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			agent.(WarlockAgent).GetWarlock().AddStats(stats.Stats{
				stats.Intellect: 70,
			})
		},
		4: func(agent core.Agent) {
			lock := agent.(WarlockAgent).GetWarlock()
			lock.AddStats(stats.Stats{
				stats.Intellect: 90,
			})

			// TODO: enable if we ever implement death coil
			// lock.AddStaticMod(core.SpellModConfig{
			// 	Kind:       core.SpellMod_Cooldown_Flat,
			// 	ClassMask:  WarlockSpellDeathCoil,
			// 	FloatValue: -30 * time.Second,
			// })
		},
	},
})
