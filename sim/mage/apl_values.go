package mage

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (mage *Mage) NewAPLValue(rot *core.APLRotation, config *proto.APLValue) core.APLValue {
	switch config.Value.(type) {
	case *proto.APLValue_MageCurrentCombustionDotEstimate:
		return mage.newValueCurrentCombustionDotEstimate(config.GetMageCurrentCombustionDotEstimate(), config.Uuid)
	default:
		return nil
	}
}

type APLValueMageCurrentCombustionDotEstimate struct {
	core.DefaultAPLValueImpl
	mage *Mage
}

func (mage *Mage) newValueCurrentCombustionDotEstimate(_ *proto.APLValueMageCurrentCombustionDotEstimate, _ *proto.UUID) core.APLValue {
	if !mage.Talents.Combustion {
		return nil
	}

	return &APLValueMageCurrentCombustionDotEstimate{
		mage: mage,
	}
}

func (value *APLValueMageCurrentCombustionDotEstimate) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeInt
}

func (value *APLValueMageCurrentCombustionDotEstimate) GetInt(sim *core.Simulation) int32 {
	return value.mage.combustionDotEstimate
}

func (value *APLValueMageCurrentCombustionDotEstimate) String() string {
	return "Combustion Dot Estimated Value"
}
