package fire

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/mage"
)

func RegisterFireMage() {
	core.RegisterAgentFactory(
		proto.Player_FireMage{},
		proto.Spec_SpecFireMage,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewFireMage(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_FireMage)
			if !ok {
				panic("Invalid spec value for Fire Mage!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewFireMage(character *core.Character, options *proto.Player) *FireMage {
	fireOptions := options.GetFireMage().Options

	fireMage := &FireMage{
		Mage: mage.NewMage(character, options, fireOptions.ClassOptions),
	}
	fireMage.FireOptions = fireOptions

	fireMage.ApplyFireSpecInnate()
	return fireMage
}

type FireMage struct {
	*mage.Mage
}

func (fireMage *FireMage) GetMage() *mage.Mage {
	return fireMage.Mage
}

func (fireMage *FireMage) Reset(sim *core.Simulation) {
	fireMage.Mage.Reset(sim)
}

func (fireMage *FireMage) Initialize() {
	fireMage.Mage.Initialize()

	fireMage.registerPyroblastSpell()
}

func (fireMage *FireMage) GetMasteryBonus() float64 {
	return (0.22 + 0.028*fireMage.GetMasteryPoints())
}

func (fireMage *FireMage) ApplyFireSpecInnate() {

	// Fire  Specialization Bonus
	fireMage.AddStaticMod(core.SpellModConfig{
		School:     core.SpellSchoolFire,
		FloatValue: 0.25,
		Kind:       core.SpellMod_DamageDone_Pct,
	})

	fireMastery := fireMage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  mage.MageSpellFireMastery, // Ignite is done inside
		FloatValue: fireMage.GetMasteryBonus(),
		Kind:       core.SpellMod_DamageDone_Pct,
	})

	fireMage.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery, newMastery float64) {
		fireMastery.UpdateFloatValue(fireMage.GetMasteryBonus())
	})

	core.MakePermanent(fireMage.GetOrRegisterAura(core.Aura{
		Label:    "Flashburn",
		ActionID: core.ActionID{SpellID: 76595},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			fireMastery.UpdateFloatValue(fireMage.GetMasteryBonus())
			fireMastery.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			fireMastery.Deactivate()
		},
	}))
}
