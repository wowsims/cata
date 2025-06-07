package arcane

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/mage"
)

func RegisterArcaneMage() {
	core.RegisterAgentFactory(
		proto.Player_ArcaneMage{},
		proto.Spec_SpecArcaneMage,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewArcaneMage(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_ArcaneMage)
			if !ok {
				panic("Invalid spec value for Arcane Mage!")
			}
			player.Spec = playerSpec
		},
	)
}

type ArcaneMage struct {
	*mage.Mage

	Options *proto.ArcaneMage_Options

	arcaneChargesAura      *core.Aura
	arcaneMissilesProcAura *core.Aura
	arcanePowerAura        *core.Aura

	arcaneMissilesTickSpell *core.Spell
	arcanePower             *core.Spell

	arcaneMissileCritSnapshot float64
}

func NewArcaneMage(character *core.Character, options *proto.Player) *ArcaneMage {
	arcaneOptions := options.GetArcaneMage().Options

	arcaneMage := &ArcaneMage{
		Mage: mage.NewMage(character, options, arcaneOptions.ClassOptions),
	}
	arcaneMage.ArcaneOptions = arcaneOptions

	return arcaneMage
}

func (arcaneMage *ArcaneMage) GetMage() *mage.Mage {
	return arcaneMage.Mage
}

func (arcaneMage *ArcaneMage) Reset(sim *core.Simulation) {
	arcaneMage.Mage.Reset(sim)
}

func (arcaneMage *ArcaneMage) Initialize() {
	arcaneMage.Mage.Initialize()

	arcaneMage.registerArcaneBarrageSpell()
	arcaneMage.registerArcaneBlastSpell()
	arcaneMage.registerArcaneCharges()
	arcaneMage.registerArcaneMissilesSpell()
	arcaneMage.registerArcanePowerCD()

	arcaneMage.registerArcanePowerCD()
}

func (arcane *ArcaneMage) ApplyTalents() {

	arcane.Mage.ApplyTalents()
	arcane.ApplyMastery()

}

func (arcaneMage *ArcaneMage) GetArcaneMasteryBonus() float64 {
	return (0.16 + 0.02*arcaneMage.GetMasteryPoints())
}

func (arcaneMage *ArcaneMage) ArcaneMasteryValue() float64 {
	return arcaneMage.GetArcaneMasteryBonus() * (arcaneMage.CurrentMana() / arcaneMage.MaxMana())
}
