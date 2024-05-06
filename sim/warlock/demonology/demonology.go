package demonology

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/warlock"
)

func RegisterDemonologyWarlock() {
	core.RegisterAgentFactory(
		proto.Player_DemonologyWarlock{},
		proto.Spec_SpecDemonologyWarlock,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewDemonologyWarlock(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_DemonologyWarlock)
			if !ok {
				panic("Invalid spec value for Demonology Warlock!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewDemonologyWarlock(character *core.Character, options *proto.Player) *DemonologyWarlock {
	demoOptions := options.GetDemonologyWarlock().Options

	demonology := &DemonologyWarlock{
		Warlock: warlock.NewWarlock(character, options, demoOptions.ClassOptions),
	}

	return demonology
}

type DemonologyWarlock struct {
	*warlock.Warlock

	MasterDemonologistOwnerMod *core.SpellMod
}

func (demonology DemonologyWarlock) getMasteryBonus() float64 {
	return 0.18 + 0.023*demonology.GetMasteryPoints()
}

func (demonology *DemonologyWarlock) GetWarlock() *warlock.Warlock {
	return demonology.Warlock
}

func (demonology *DemonologyWarlock) Initialize() {
	demonology.Warlock.Initialize()

	demonology.registerHandOfGuldanSpell()
	demonology.registerMetamorphosisSpell()
	demonology.registerSummonFelguardSpell()
}

func (demonology *DemonologyWarlock) ApplyTalents() {
	demonology.Warlock.ApplyTalents()

	//Mastery: Master Demonologist
	//TODO: Is there a better way to apply this mod/activate/deactivate it from all pets?
	demonology.MasterDemonologistOwnerMod = demonology.AddDynamicMod(core.SpellModConfig{
		Kind: core.SpellMod_DamageDone_Pct,
	})

	MasterDemonologistPetMod_Felguard := demonology.Felguard.AddDynamicMod(core.SpellModConfig{
		Kind: core.SpellMod_DamageDone_Pct,
	})

	MasterDemonologistPetMod_Felhunter := demonology.Felhunter.AddDynamicMod(core.SpellModConfig{
		Kind: core.SpellMod_DamageDone_Pct,
	})

	MasterDemonologistPetMod_Imp := demonology.Imp.AddDynamicMod(core.SpellModConfig{
		Kind: core.SpellMod_DamageDone_Pct,
	})

	MasterDemonologistPetMod_Succubus := demonology.Succubus.AddDynamicMod(core.SpellModConfig{
		Kind: core.SpellMod_DamageDone_Pct,
	})

	MasterDemonologistPetMod_Doomguard := demonology.Doomguard.AddDynamicMod(core.SpellModConfig{
		Kind: core.SpellMod_DamageDone_Pct,
	})

	MasterDemonologistPetMod_Infernal := demonology.Infernal.AddDynamicMod(core.SpellModConfig{
		Kind: core.SpellMod_DamageDone_Pct,
	})

	MasterDemonologistPetMod_EbonImp := demonology.EbonImp.AddDynamicMod(core.SpellModConfig{
		Kind: core.SpellMod_DamageDone_Pct,
	})

	demonology.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery float64, newMastery float64) {
		demonology.MasterDemonologistOwnerMod.UpdateFloatValue(demonology.getMasteryBonus())
		MasterDemonologistPetMod_Felguard.UpdateFloatValue(demonology.getMasteryBonus())
		MasterDemonologistPetMod_Felhunter.UpdateFloatValue(demonology.getMasteryBonus())
		MasterDemonologistPetMod_Imp.UpdateFloatValue(demonology.getMasteryBonus())
		MasterDemonologistPetMod_Succubus.UpdateFloatValue(demonology.getMasteryBonus())
		MasterDemonologistPetMod_Doomguard.UpdateFloatValue(demonology.getMasteryBonus())
		MasterDemonologistPetMod_Infernal.UpdateFloatValue(demonology.getMasteryBonus())
		MasterDemonologistPetMod_EbonImp.UpdateFloatValue(demonology.getMasteryBonus())
	})

	core.MakePermanent(demonology.GetOrRegisterAura(core.Aura{
		Label:    "Mastery: Master Demonologist",
		ActionID: core.ActionID{SpellID: 77219},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			demonology.MasterDemonologistOwnerMod.UpdateFloatValue(demonology.getMasteryBonus())
			MasterDemonologistPetMod_Felguard.Activate()
			MasterDemonologistPetMod_Felhunter.Activate()
			MasterDemonologistPetMod_Imp.Activate()
			MasterDemonologistPetMod_Succubus.Activate()
			MasterDemonologistPetMod_Doomguard.Activate()
			MasterDemonologistPetMod_Infernal.Activate()
			MasterDemonologistPetMod_EbonImp.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			MasterDemonologistPetMod_Felguard.Deactivate()
			MasterDemonologistPetMod_Felhunter.Deactivate()
			MasterDemonologistPetMod_Imp.Deactivate()
			MasterDemonologistPetMod_Succubus.Deactivate()
			MasterDemonologistPetMod_Doomguard.Deactivate()
			MasterDemonologistPetMod_Infernal.Deactivate()
			MasterDemonologistPetMod_EbonImp.Deactivate()
		},
	}))

	// Demonic Knowledge
	demonology.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  warlock.WarlockShadowDamage | warlock.WarlockFireDamage,
		FloatValue: 0.15,
	})
}

func (demonology *DemonologyWarlock) Reset(sim *core.Simulation) {
	demonology.Warlock.Reset(sim)
}

func (demonology *DemonologyWarlock) registerSummonFelguardSpell() {
	demonology.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 30146},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: warlock.WarlockSpellSummonFelguard,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.8,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 6,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			demonology.SoulBurnAura.Deactivate(sim)
			demonology.ChangeActivePet(sim, warlock.PetFelguard)
		},
	})
}
