package fire

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/mage"
)

const (
	DDBC_Pyromaniac int = iota
	DDBC_Total
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

	return fireMage
}

type FireMage struct {
	*mage.Mage

	Combustion   *core.Spell
	Ignite       *core.Spell
	Pyroblast    *core.Spell
	InfernoBlast *core.Spell

	pyromaniacAuras core.AuraArray

	combustionDotEstimate int32
}

func (fireMage *FireMage) GetMage() *mage.Mage {
	return fireMage.Mage
}

func (fireMage *FireMage) Reset(sim *core.Simulation) {
	fireMage.Mage.Reset(sim)
}

func (fireMage *FireMage) Initialize() {
	fireMage.Mage.Initialize()

	fireMage.registerPassives()
	fireMage.registerSpells()
}

func (fireMage *FireMage) registerPassives() {
	fireMage.registerMastery()
	fireMage.registerCriticalMass()
	fireMage.registerPyromaniac()
}

func (fireMage *FireMage) registerSpells() {
	fireMage.registerCombustionSpell()
	fireMage.registerFireballSpell()
	fireMage.registerInfernoBlastSpell()
	fireMage.registerDragonsBreathSpell()
	fireMage.registerPyroblastSpell()
	fireMage.registerScorchSpell()
}
