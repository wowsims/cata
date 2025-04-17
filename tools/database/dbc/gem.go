package dbc

import (
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
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

// Todo: All dbc types should have a ToProto
// Though we will need a data layer
// func (gem Gem) ToProto() *proto.UIGem {

// }
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

	gem.SetGemSpellEffects()

	return uiGem
}

func (gem *Gem) SetGemSpellEffects() {
	// For example spell pen
	// or all stats
	// etc.
	//Todo: Unfinished
	// Could append actual effects on the gem that is applied on sim start???
	// Same for items and enchants hmmmm
}

func (gem *Gem) GetItemEnchantmentStats() stats.Stats {
	stats := stats.Stats{}
	processEnchantmentEffects(gem.Effects, gem.EffectArgs, gem.EffectPoints, &stats, false)
	return stats
}
