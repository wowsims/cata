package frost

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/mage"
)

func RegisterFrostMage() {
	core.RegisterAgentFactory(
		proto.Player_FrostMage{},
		proto.Spec_SpecFrostMage,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewFrostMage(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_FrostMage)
			if !ok {
				panic("Invalid spec value for Frost Mage!")
			}
			player.Spec = playerSpec
		},
	)
}

type FrostMage struct {
	*mage.Mage

	waterElemental *WaterElemental
}

func NewFrostMage(character *core.Character, options *proto.Player) *FrostMage {
	frostOptions := options.GetFrostMage().Options

	frostMage := &FrostMage{
		Mage: mage.NewMage(character, options, frostOptions.ClassOptions),
	}
	frostMage.waterElemental = frostMage.NewWaterElemental(0.20)

	return frostMage
}

func (frostMage *FrostMage) GetMage() *mage.Mage {
	return frostMage.Mage
}

func (frostMage *FrostMage) Reset(sim *core.Simulation) {
	frostMage.Mage.Reset(sim)
}

func (frostMage *FrostMage) Initialize() {
	frostMage.Mage.Initialize()

	frostMage.registerSummonWaterElementalSpell()
}

func (frostMage *FrostMage) ApplyTalents() {
	frostMage.Mage.ApplyTalents()

	// Frost  Specialization Bonus
	frostMage.Mage.AddStaticMod(core.SpellModConfig{
		School:     core.SpellSchoolFrost,
		FloatValue: 0.25,
		Kind:       core.SpellMod_DamageDone_Pct,
	})

	frostMage.waterElemental.AddStaticMod(core.SpellModConfig{
		School:     core.SpellSchoolFrost,
		FloatValue: 0.25,
		Kind:       core.SpellMod_DamageDone_Pct,
	})

	frostMage.Mage.AddStaticMod(core.SpellModConfig{
		ClassMask:  mage.MageSpellFrostbolt,
		FloatValue: 0.15,
		Kind:       core.SpellMod_DamageDone_Flat,
	})

	// Frost Mastery Bonus

	frostMasteryMod := frostMage.Mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  mage.MageSpellIceLance | mage.MageSpellDeepFreeze,
		FloatValue: frostMage.GetMasteryBonus(),
		Kind:       core.SpellMod_DamageDone_Pct,
	})

	core.MakePermanent(frostMage.Mage.GetOrRegisterAura(core.Aura{
		Label:    "Frostburn",
		ActionID: core.ActionID{SpellID: 76595},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			frostMasteryMod.UpdateFloatValue(frostMage.GetMasteryBonus())
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			frostMasteryMod.Deactivate()
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if frostMage.Mage.FingersOfFrostAura.IsActive() {
				frostMasteryMod.Activate()
			}
		},
	}))

	frostMage.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery, newMastery float64) {
		frostMasteryMod.UpdateFloatValue(frostMage.GetMasteryBonus())
	})
}

func (frostMage *FrostMage) GetMasteryBonus() float64 {
	return (.05 + 0.025*frostMage.Mage.GetMasteryPoints())
}
