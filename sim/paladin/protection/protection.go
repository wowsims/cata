package protection

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/paladin"
)

func RegisterProtectionPaladin() {
	core.RegisterAgentFactory(
		proto.Player_ProtectionPaladin{},
		proto.Spec_SpecProtectionPaladin,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewProtectionPaladin(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_ProtectionPaladin)
			if !ok {
				panic("Invalid spec value for Protection Paladin!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewProtectionPaladin(character *core.Character, options *proto.Player) *ProtectionPaladin {
	protOptions := options.GetProtectionPaladin()

	prot := &ProtectionPaladin{
		Paladin: paladin.NewPaladin(character, options.TalentsString, protOptions.Options.ClassOptions),
		Options: protOptions.Options,
	}

	return prot
}

type ProtectionPaladin struct {
	*paladin.Paladin

	Options *proto.ProtectionPaladin_Options
}

func (prot *ProtectionPaladin) GetPaladin() *paladin.Paladin {
	return prot.Paladin
}

func (prot *ProtectionPaladin) Initialize() {
	prot.Paladin.Initialize()

	prot.registerMastery()

	prot.registerArdentDefender()
	prot.registerAvengersShieldSpell()
	prot.registerConsecrationSpell()
	prot.registerGrandCrusader()
	prot.registerGuardedByTheLight()
	prot.registerHolyWrath()
	prot.registerJudgmentsOfTheWise()
	prot.registerRighteousFury()
	prot.registerSanctuary()

	// Vengeance
	prot.RegisterVengeance(84839, nil)

	prot.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  paladin.SpellMaskSealOfTruth | paladin.SpellMaskCensure,
		FloatValue: 0.2,
	})
}

func (prot *ProtectionPaladin) ApplyTalents() {
	prot.Paladin.ApplyTalents()
	prot.ApplyArmorSpecializationEffect(stats.Stamina, proto.ArmorType_ArmorTypePlate, 86525)
}

func (prot *ProtectionPaladin) Reset(sim *core.Simulation) {
	prot.Paladin.Reset(sim)
}
