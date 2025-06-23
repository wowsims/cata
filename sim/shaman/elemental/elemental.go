package elemental

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/shaman"
)

func RegisterElementalShaman() {
	core.RegisterAgentFactory(
		proto.Player_ElementalShaman{},
		proto.Spec_SpecElementalShaman,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewElementalShaman(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_ElementalShaman)
			if !ok {
				panic("Invalid spec value for Elemental Shaman!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewElementalShaman(character *core.Character, options *proto.Player) *ElementalShaman {
	eleOptions := options.GetElementalShaman().Options

	selfBuffs := shaman.SelfBuffs{
		Shield: eleOptions.ClassOptions.Shield,
	}

	inRange := eleOptions.ThunderstormRange == proto.ElementalShaman_Options_TSInRange
	ele := &ElementalShaman{
		Shaman: shaman.NewShaman(character, options.TalentsString, selfBuffs, inRange, eleOptions.ClassOptions.FeleAutocast),
	}

	if mh := ele.GetMHWeapon(); mh != nil {
		ele.ApplyFlametongueImbueToItem(mh)
		ele.SelfBuffs.ImbueMH = proto.ShamanImbue_FlametongueWeapon
	}

	return ele
}

func (eleShaman *ElementalShaman) Initialize() {
	eleShaman.Shaman.Initialize()

	eleShaman.registerThunderstormSpell()
	eleShaman.registerLavaBurstSpell()
	eleShaman.registerEarthquakeSpell()
	eleShaman.registerLavaBeamSpell()
	eleShaman.registerShamanisticRageSpell()
}

func (ele *ElementalShaman) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	raidBuffs.ElementalOath = true
	ele.Shaman.AddRaidBuffs(raidBuffs)
}

func (ele *ElementalShaman) ApplyTalents() {
	ele.ApplyElementalTalents()
	ele.Shaman.ApplyTalents()
	ele.ApplyArmorSpecializationEffect(stats.Intellect, proto.ArmorType_ArmorTypeMail, 86529)
}

type ElementalShaman struct {
	*shaman.Shaman
}

func (eleShaman *ElementalShaman) GetShaman() *shaman.Shaman {
	return eleShaman.Shaman
}

func (eleShaman *ElementalShaman) Reset(sim *core.Simulation) {
	eleShaman.Shaman.Reset(sim)
}
