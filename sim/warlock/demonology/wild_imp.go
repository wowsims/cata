package demonology

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/warlock"
)

// wild imps will cast 10 casts then despawn
// they fight like any other guardian imp
// we can potentially spawn a lot of imps due to Doom being able to proc them so.. fingers crossed >.<

type WildImpPet struct {
	core.Pet

	Fireball *core.Spell
}

// registers the wild imp spell and handlers
// count The number of imps that shoudl be registered. It will be the upper limit the sim can spawn simultaniously
func (demonology *DemonologyWarlock) registerWildImp(count int) {
	demonology.WildImps = make([]*WildImpPet, count)
	for idx := 0; idx < count; idx++ {
		demonology.WildImps[idx] = demonology.buildWildImp(count)
		demonology.AddPet(demonology.WildImps[idx])
	}

	// register passiv
	demonology.registerWildImpPassive()
}

func (demonology *DemonologyWarlock) buildWildImp(counter int) *WildImpPet {
	pet := &WildImpPet{
		Pet: core.NewPet(core.PetConfig{
			Name:                            "Wild Imp",
			Owner:                           &demonology.Character,
			BaseStats:                       stats.Stats{stats.Health: 48312.8, stats.Armor: 19680},
			StatInheritance:                 demonology.SimplePetStatInheritanceWithScale(0),
			EnabledOnStart:                  false,
			IsGuardian:                      true,
			HasDynamicMeleeSpeedInheritance: true,
			HasDynamicCastSpeedInheritance:  true,
			HasResourceRegenInheritance:     true,
		}),
	}

	// set pet class for proper scaling values
	pet.Class = pet.Owner.Class
	pet.EnableEnergyBar(core.EnergyBarOptions{
		MaxEnergy:  10,
		HasNoRegen: true,
	})

	pet.registerFireboltSpell()
	return pet
}

func (pet *WildImpPet) GetPet() *core.Pet {
	return &pet.Pet
}

func (pet *WildImpPet) Reset(sim *core.Simulation) {
}

func (pet *WildImpPet) ExecuteCustomRotation(sim *core.Simulation) {
	spell := pet.Fireball
	if spell.CanCast(sim, pet.CurrentTarget) {
		spell.Cast(sim, pet.CurrentTarget)
		return
	}

	if pet.CurrentEnergy() == 0 {
		if sim.Log != nil {
			pet.Log(sim, "Wild Imp despawned.")
		}

		sim.AddPendingAction(&core.PendingAction{
			NextActionAt: sim.CurrentTime,
			Priority:     core.ActionPriorityAuto,
			OnAction: func(sim *core.Simulation) {
				pet.Disable(sim)
			},
		})

		return
	}

	var offset = time.Duration(0)
	if pet.Hardcast.Expires > sim.CurrentTime {
		offset = pet.Hardcast.Expires - sim.CurrentTime
	}

	pet.WaitUntil(sim, sim.CurrentTime+offset+time.Millisecond*100)
}

// Hotfixes already included
const felFireBoltScale = 0.242 * 1.43 // 2025.06.13 Changes to Beta - Wild Imp Damage increased by 43%
const felFireBoltVariance = 0.05
const felFireBoltCoeff = 0.242 * 1.43

func (pet *WildImpPet) registerFireboltSpell() {
	pet.Fireball = pet.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 104318},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: warlock.WarlockSpellImpFireBolt,
		MissileSpeed:   16,

		EnergyCost: core.EnergyCostOptions{
			Cost: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      time.Second * 1,
				CastTime: time.Second * 2,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   2,
		ThreatMultiplier: 1,
		BonusCoefficient: felFireBoltCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			pet.Owner.Unit.GetSecondaryResourceBar().Gain(sim, 5, spell.ActionID)
			result := spell.CalcDamage(sim, target, pet.CalcAndRollDamageRange(sim, felFireBoltScale, felFireBoltVariance), spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}

func (warlock *DemonologyWarlock) SpawnImp(sim *core.Simulation) {
	for _, pet := range warlock.WildImps {
		if pet.IsActive() {
			continue
		}

		pet.Enable(sim, pet)
		return
	}

	panic("TOO MANY IMPS!")
}

func (demonology *DemonologyWarlock) registerWildImpPassive() {
	var trigger *core.Aura
	trigger = core.MakeProcTriggerAura(&demonology.Unit, core.ProcTrigger{
		MetricsActionID: core.ActionID{SpellID: 114925},
		Name:            "Demonic Calling",
		Callback:        core.CallbackOnCastComplete,
		ClassSpellMask:  warlock.WarlockSpellShadowBolt | warlock.WarlockSpellSoulFire | warlock.WarlockSpellTouchOfChaos,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			demonology.SpawnImp(sim)
			trigger.Deactivate(sim)
		},
	})

	getCD := func() time.Duration {
		return time.Duration(
			core.TernaryFloat64(
				demonology.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfImpSwarm), 24, 20)/
				demonology.TotalSpellHasteMultiplier()) * time.Second
	}

	var triggerAction *core.PendingAction
	var controllerImpSpawn func(sim *core.Simulation)
	controllerImpSpawn = func(sim *core.Simulation) {
		if demonology.ImpSwarm == nil || demonology.ImpSwarm.CD.IsReady(sim) {
			trigger.Activate(sim)
		}

		triggerAction = &core.PendingAction{
			NextActionAt: sim.CurrentTime + getCD(),
			Priority:     core.ActionPriorityAuto,
			OnAction:     controllerImpSpawn,
		}

		sim.AddPendingAction(triggerAction)
	}

	core.MakePermanent(demonology.RegisterAura(core.Aura{
		Label: "Wild Imp - Controller",
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			cd := time.Duration(sim.Roll(float64(time.Second), float64(getCD())))

			// initially do random timer to simulate real world scenario more appropiate
			triggerAction = &core.PendingAction{
				NextActionAt: sim.CurrentTime + cd,
				Priority:     core.ActionPriorityAuto,
				OnAction:     controllerImpSpawn,
			}

			sim.AddPendingAction(triggerAction)
		},
	}))

	core.MakeProcTriggerAura(&demonology.Unit, core.ProcTrigger{
		Name:           "Wild Imp - Doom Monitor",
		ClassSpellMask: warlock.WarlockSpellDoom,
		Outcome:        core.OutcomeCrit,
		Callback:       core.CallbackOnPeriodicDamageDealt,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			demonology.SpawnImp(sim)
		},
	})
}
