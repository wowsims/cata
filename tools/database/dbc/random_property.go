package dbc

import "github.com/wowsims/cata/sim/core/proto"

type RandomPropAllocation struct {
	Epic0     int32 `json:"Epic_0"`
	Epic1     int32 `json:"Epic_1"`
	Epic2     int32 `json:"Epic_2"`
	Epic3     int32 `json:"Epic_3"`
	Epic4     int32 `json:"Epic_4"`
	Superior0 int32 `json:"Superior_0"`
	Superior1 int32 `json:"Superior_1"`
	Superior2 int32 `json:"Superior_2"`
	Superior3 int32 `json:"Superior_3"`
	Superior4 int32 `json:"Superior_4"`
	Good0     int32 `json:"Good_0"`
	Good1     int32 `json:"Good_1"`
	Good2     int32 `json:"Good_2"`
	Good3     int32 `json:"Good_3"`
	Good4     int32 `json:"Good_4"`
}

type RandomPropAllocationMap map[proto.ItemQuality][5]int32

type RandomPropAllocationsByIlvl map[int]RandomPropAllocationMap

func (r RandomPropAllocation) ToProto() *proto.RandomPropAllocation {
	return &proto.RandomPropAllocation{
		Allocations: &proto.QualityAllocations{
			Good:     []int32{r.Good0, r.Good1, r.Good2, r.Good3, r.Good4},
			Superior: []int32{r.Superior0, r.Superior1, r.Superior2, r.Superior3, r.Superior4},
			Epic:     []int32{r.Epic0, r.Epic1, r.Epic2, r.Epic3, r.Epic4},
		},
	}
}
