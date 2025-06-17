package monk

import (
	"fmt"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (monk *Monk) NewAPLValue(_ *core.APLRotation, config *proto.APLValue) core.APLValue {
	switch config.Value.(type) {
	case *proto.APLValue_MonkCurrentChi:
		return monk.newValueCurrentChi(config.GetMonkCurrentChi(), config.Uuid)
	case *proto.APLValue_MonkMaxChi:
		return monk.newValueMaxChi(config.GetMonkMaxChi(), config.Uuid)
	default:
		return nil
	}
}

type APLValueCurrentChi struct {
	core.DefaultAPLValueImpl
	monk *Monk
}

func (monk *Monk) newValueCurrentChi(_ *proto.APLValueMonkCurrentChi, _ *proto.UUID) core.APLValue {
	return &APLValueCurrentChi{
		monk: monk,
	}
}
func (value *APLValueCurrentChi) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeInt
}
func (value *APLValueCurrentChi) GetInt(_ *core.Simulation) int32 {
	return value.monk.GetChi()
}
func (value *APLValueCurrentChi) String() string {
	return "Current Chi"
}

type APLValueMaxChi struct {
	core.DefaultAPLValueImpl
	maxChi int32
}

func (monk *Monk) newValueMaxChi(_ *proto.APLValueMonkMaxChi, _ *proto.UUID) core.APLValue {
	return &APLValueMaxChi{
		maxChi: monk.GetMaxChi(),
	}
}
func (value *APLValueMaxChi) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeInt
}
func (value *APLValueMaxChi) GetInt(_ *core.Simulation) int32 {
	return value.maxChi
}
func (value *APLValueMaxChi) String() string {
	return fmt.Sprintf("Max Chi(%d)", value.maxChi)
}
