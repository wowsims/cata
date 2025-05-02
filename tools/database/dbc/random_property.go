package dbc

import "github.com/wowsims/mop/sim/core/proto"

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
