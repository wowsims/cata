package elemental

import (
	"time"

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

	totems := &proto.ShamanTotems{}
	if eleOptions.ClassOptions.Totems != nil {
		totems = eleOptions.ClassOptions.Totems
	}

	inRange := eleOptions.ThunderstormRange == proto.ElementalShaman_Options_TSInRange
	ele := &ElementalShaman{
		Shaman: shaman.NewShaman(character, options.TalentsString, totems, selfBuffs, inRange),
	}

	if mh := ele.GetMHWeapon(); mh != nil {
		// ele.ApplyFlametongueImbueToItem(mh)
		ele.SelfBuffs.ImbueMH = proto.ShamanImbue_FlametongueWeapon
	}

	return ele
}

func (eleShaman *ElementalShaman) Initialize() {
	eleShaman.Shaman.Initialize()

	eleShaman.registerThunderstormSpell()

	// Shamanism
	eleShaman.AddStaticMod(core.SpellModConfig{
		ClassMask: shaman.SpellMaskLavaBurst | shaman.SpellMaskChainLightning | shaman.SpellMaskLightningBolt,
		Kind:      core.SpellMod_CastTime_Flat,
		TimeValue: time.Millisecond * -500,
	})

	eleShaman.AddStaticMod(core.SpellModConfig{
		ClassMask: shaman.SpellMaskLavaBurst | shaman.SpellMaskChainLightning | shaman.SpellMaskLightningBolt |
			shaman.SpellMaskLightningBoltOverload | shaman.SpellMaskChainLightningOverload | shaman.SpellMaskLavaBurstOverload,
		Kind:       core.SpellMod_BonusCoeffecient_Flat,
		FloatValue: 0.36,
	})

	// Elemental Fury
	eleShaman.AddStaticMod(core.SpellModConfig{
		ClassMask: shaman.SpellMaskFire | shaman.SpellMaskNature |
			shaman.SpellMaskFrost | shaman.SpellMaskMagmaTotem | shaman.SpellMaskSearingTotem | shaman.SpellMaskEarthquake,
		Kind:       core.SpellMod_CritMultiplier_Flat,
		FloatValue: 1.0,
	})

	eleShaman.AddStaticMod(core.SpellModConfig{
		ClassMask: shaman.SpellMaskChainLightning,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: time.Second * -3,
	})
}

func (ele *ElementalShaman) ApplyTalents() {
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
