package dbc

import (
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

type RandomSuffix struct {
	ID            int
	Name          string
	AllocationPct []int // AllocationPct_0-4
	EffectArgs    []int // EffectArg_0-4
	Effects       []int // Effect_0-4
}

func (raw RandomSuffix) ToProto() *proto.ItemRandomSuffix {
	suffix := &proto.ItemRandomSuffix{
		Name:  raw.Name,
		Id:    int32(raw.ID),
		Stats: stats.Stats{}.ToProtoArray(),
	}
	for i, effect := range raw.Effects {
		var stat proto.Stat
		var matchFound bool

		if effect == 5 {
			stat, matchFound = MapBonusStatIndexToStat(raw.EffectArgs[i])
		} else if effect == 4 {
			stat, matchFound = MapResistanceToStat(raw.EffectArgs[i])
		}

		if !matchFound {
			continue
		}

		amount := raw.AllocationPct[i]
		suffix.Stats[stat] = float64(amount)
	}
	return suffix
}
