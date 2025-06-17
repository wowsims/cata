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

	arcaneMissilesProcAura *core.Aura
	arcanePowerAura        *core.Aura

	arcaneMissiles          *core.Spell
	arcaneMissilesTickSpell *core.Spell
	arcanePower             *core.Spell

	arcaneMissileCritSnapshot float64
}

func NewArcaneMage(character *core.Character, options *proto.Player) *ArcaneMage {
	arcaneOptions := options.GetArcaneMage().Options

	arcane := &ArcaneMage{
		Mage: mage.NewMage(character, options, arcaneOptions.ClassOptions),
	}
	arcane.ArcaneOptions = arcaneOptions

	return arcane
}

func (arcaneMage *ArcaneMage) GetMage() *mage.Mage {
	return arcaneMage.Mage
}

func (arcaneMage *ArcaneMage) Reset(sim *core.Simulation) {
	arcaneMage.Mage.Reset(sim)
}

func (arcane *ArcaneMage) Initialize() {
	arcane.Mage.Initialize()

	arcane.registerPassives()
	arcane.registerSpells()
}

func (arcane *ArcaneMage) registerPassives() {
	arcane.registerMastery()
	arcane.registerArcaneCharges()
}

func (arcane *ArcaneMage) registerSpells() {
	arcane.registerArcaneBarrageSpell()
	arcane.registerArcaneBlastSpell()
	arcane.registerArcaneMissilesSpell()
	arcane.registerArcanePowerCD()
}
