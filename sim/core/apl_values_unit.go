package core

import (
	"github.com/wowsims/cata/sim/core/proto"
)

type APLValueUnitIsMoving struct {
	DefaultAPLValueImpl
	unit UnitReference
}

func (rot *APLRotation) newValueCharacterIsMoving(config *proto.APLValueUnitIsMoving, uuid *proto.UUID) APLValue {
	unit := rot.GetSourceUnit(config.SourceUnit)
	if unit.Get() == nil {
		return nil
	}
	return &APLValueUnitIsMoving{
		unit: unit,
	}
}
func (value *APLValueUnitIsMoving) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueUnitIsMoving) GetBool(sim *Simulation) bool {
	return value.unit.Get().Moving
}
func (value *APLValueUnitIsMoving) String() string {
	return "Is Moving"
}
