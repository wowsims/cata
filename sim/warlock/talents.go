package warlock

import (
	"fmt"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (warlock *Warlock) registerHarvestLife() {
	if !warlock.Talents.HarvestLife {
		return
	}

	warlock.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.5,
		ClassMask:  WarlockSpellDrainLife,
	})
}

func (warlock *Warlock) registerArchimondesDarkness() {
	if !warlock.Talents.ArchimondesDarkness {
		return
	}

	warlock.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_ModCharges_Flat,
		IntValue:  2,
		ClassMask: WarlockDarkSoulSpell,
	})

	warlock.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: -time.Second * 100,
		ClassMask: WarlockDarkSoulSpell,
	})
}

func (warlock *Warlock) registerKilJaedensCunning() {
	if !warlock.Talents.KiljaedensCunning {
		return
	}

	warlock.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_AllowCastWhileMoving,
		ClassMask: WarlockSpellIncinerate | WarlockSpellShadowBolt | WarlockSpellMaleficGrasp,
	})
}

func (warlock *Warlock) registerMannarothsFury() {
	if !warlock.Talents.MannorothsFury {
		return
	}

	buff := warlock.RegisterAura(core.Aura{
		Label:    "Mannaroth's Fury",
		ActionID: core.ActionID{SpellID: 108508},
		Duration: time.Second * 10,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  WarlockSpellRainOfFire | WarlockSpellSeedOfCorruptionExposion | WarlockSpellSeedOfCorruption | WarlockSpellHellfire | WarlockSpellImmolationAura,
		FloatValue: 1,
	})

	warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 108508},
		Flags:       core.SpellFlagAPL,
		ProcMask:    core.ProcMaskEmpty,
		SpellSchool: core.SpellSchoolShadow,

		ThreatMultiplier: 1,
		DamageMultiplier: 1,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},

			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Minute,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			buff.Activate(sim)
		},
	})
}

func (warlock *Warlock) registerGrimoireOfSupremacy() {
	if !warlock.Talents.GrimoireOfSupremacy {
		return
	}

	// It's honestly just a repaint and a 20% damage mod slapped on top..
	fireBolt := warlock.Imp.GetSpell(petActionFireBolt)

	// It's now Felbolt!
	fireBolt.ActionID = core.ActionID{SpellID: 115746}
	fireBolt.DamageMultiplier *= 1.2
	updateName(&warlock.Imp.Pet, "Fel Imp")

	// Spell stays the same
	warlock.Voidwalker.PseudoStats.DamageDealtMultiplier *= 1.2
	updateName(&warlock.Voidwalker.Pet, "Voidlord")

	// Now Tongue Lash
	shadowBite := warlock.Felhunter.GetSpell(petActionShadowBite)
	shadowBite.ActionID = core.ActionID{SpellID: 115778}
	warlock.Felhunter.PseudoStats.DamageDealtMultiplier *= 1.2
	updateName(&warlock.Felhunter.Pet, "Observer")

	// Succubus get's larger makeover
	// Now dualwield with 1.5x less base damage
	weaponConfig := ScaledAutoAttackConfig(3)
	weaponConfig.MainHand.BaseDamageMax /= 1.5
	weaponConfig.MainHand.BaseDamageMin /= 1.5
	weaponConfig.OffHand = weaponConfig.MainHand

	warlock.Succubus.EnableAutoAttacks(warlock.Succubus, *weaponConfig)
	warlock.Succubus.ChangeStatInheritance(warlock.SimplePetStatInheritanceWithScale(1 + 1.0/9.0))
	lashOfPain := warlock.Succubus.GetSpell(petActionLashOfPain)
	lashOfPain.ActionID = core.ActionID{SpellID: 115748}
	warlock.Succubus.PseudoStats.DamageDealtMultiplier *= 1.2
	warlock.Succubus.PseudoStats.DisableDWMissPenalty = true
	updateName(&warlock.Succubus.Pet, "Shivarra")

	updateName(&warlock.Infernal.Pet, "Abyssal")
	warlock.Infernal.PseudoStats.DamageDealtMultiplier *= 1.2

	updateName(&warlock.Doomguard.Pet, "Terrorguard")
	warlock.Doomguard.PseudoStats.DamageDealtMultiplier *= 1.2

	warlock.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct,
		FloatValue: -0.5,
		ClassMask:  WarlockSpellSummonInfernal | WarlockSpellSummonDoomguard,
	})
}

func updateName(pet *core.Pet, name string) {
	pet.Name = name
	pet.Label = fmt.Sprintf("%s - %s", pet.Owner.Label, name)
}

func (warlock *Warlock) registerGrimoireOfService() {
	if !warlock.Talents.GrimoireOfService {
		return
	}

	// build all pets as they're additional summons
	imp := warlock.registerImpWithName("Grimoire: Imp", false, true)
	felHunter := warlock.registerFelHunterWithName("Grimoire: Felhunter", false, true)
	voidWalker := warlock.registerVoidWalkerWithName("Grimoire: Voidwalker", false, true)
	succubus := warlock.registerSuccubusWithName("Grimoire: Succubus", false, true)

	warlock.serviceTimer = warlock.NewTimer()

	warlock.BuildAndRegisterSummonSpell(111859, imp)
	warlock.BuildAndRegisterSummonSpell(111895, voidWalker)
	warlock.BuildAndRegisterSummonSpell(111896, succubus)
	warlock.BuildAndRegisterSummonSpell(111897, felHunter)
}

func (warlock *Warlock) BuildAndRegisterSummonSpell(id int32, pet *WarlockPet) {
	for _, spell := range pet.AutoCastAbilities {
		spell.Flags &= ^core.SpellFlagAPL
	}

	warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: id},
		SpellSchool: core.SpellSchoolShadow,
		Flags:       core.SpellFlagAPL,
		ProcMask:    core.ProcMaskEmpty,

		ThreatMultiplier: 1,
		DamageMultiplier: 1,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    warlock.serviceTimer,
				Duration: time.Minute * 2,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			pet.EnableWithTimeout(sim, pet, time.Second*20)
		},
	})
}

func (warlock *Warlock) registerGrimoireOfSacrifice() {
	if !warlock.Talents.GrimoireOfSacrifice {
		return
	}

	buff := warlock.RegisterAura(core.Aura{
		Label:    "Grimioire of Sacrifice",
		ActionID: core.ActionID{SpellID: 108503},
		Duration: time.Hour,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Matches(WarlockSpellChaosBolt) || !result.Landed() {
				return
			}

			warlock.ApplyDotWithPandemic(spell.Dot(result.Target), sim)
		},
	})

	switch warlock.Spec {
	case proto.Spec_SpecDemonologyWarlock:
		buff.AttachSpellMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Pct,
			FloatValue: 0.25,
			ClassMask:  WarlockSpellShadowBolt | WarlockSpellSoulBurn | WarlockSpellHandOfGuldan | WarlockSpellChaosWave | WarlockSpellTouchOfChaos | WarlockSpellDemonicSlash | WarlockSpellVoidray | WarlockSpellSoulFire,
		})
	case proto.Spec_SpecAfflictionWarlock:
		buff.AttachSpellMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Pct,
			FloatValue: 0.2,
			ClassMask:  WarlockSpellDrainSoul | WarlockSpellMaleficGrasp | WarlockSpellFelFlame,
		})
	case proto.Spec_SpecDestructionWarlock:
		buff.AttachSpellMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Pct,
			FloatValue: 0.15,
			ClassMask:  WarlockSpellConflagrate | WarlockSpellShadowBurn | WarlockSpellFelFlame | WarlockSpellIncinerate | WarlockSpellDrainLife,
		})
	}

	applyPetHook := func(pet *WarlockPet) {
		oldEnable := pet.OnPetEnable
		pet.OnPetEnable = func(sim *core.Simulation) {
			if oldEnable != nil {
				oldEnable(sim)
			}

			if buff.IsActive() {
				buff.Deactivate(sim)
			}
		}
	}

	warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 108503},
		SpellSchool: core.SpellSchoolFire,
		Flags:       core.SpellFlagAPL,
		ProcMask:    core.ProcMaskEmpty,

		ThreatMultiplier: 1,
		DamageMultiplier: 1,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},

			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Second * 30,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if warlock.ActivePet != nil {
				warlock.ActivePet.Disable(sim)
			}

			buff.Activate(sim)
		},
	})

	for _, pet := range warlock.Pets {
		pet.DisableOnStart()
	}

	applyPetHook(warlock.Imp)
	applyPetHook(warlock.Succubus)
	applyPetHook(warlock.Felhunter)
	applyPetHook(warlock.Voidwalker)
}
