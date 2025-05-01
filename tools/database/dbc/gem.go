package dbc

import (
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

type Gem struct {
	ItemId       int
	Name         string
	FDID         int
	GemType      GemType
	Effects      []int
	EffectPoints []int
	EffectArgs   []int
	MinItemLevel int
	Quality      ItemQuality
	IsJc         bool
	Flags0       ItemStaticFlags0
}

func (gem *Gem) ToProto() *proto.UIGem {
	uiGem := &proto.UIGem{
		Id:   int32(gem.ItemId),
		Name: gem.Name,
		//Icon:    strings.ToLower(GetIconName(iconsMap, gem.FDID)),
		Quality: gem.Quality.ToProto(),
		Color:   gem.GemType.ToProto(),
		Unique:  gem.Flags0.Has(UNIQUE_EQUIPPABLE),
		Stats:   gem.GetItemEnchantmentStats().ToProtoArray(),
	}
	if gem.IsJc {
		uiGem.RequiredProfession = proto.Profession_Jewelcrafting
	}

	return uiGem
}

func (gem *Gem) GetItemEnchantmentStats() stats.Stats {
	stats := stats.Stats{}
	processEnchantmentEffects(gem.Effects, gem.EffectArgs, gem.EffectPoints, &stats, false)
	return stats
}
