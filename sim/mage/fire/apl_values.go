package fire

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (mage *FireMage) NewAPLValue(rot *core.APLRotation, config *proto.APLValue) core.APLValue {
	switch config.Value.(type) {
	case *proto.APLValue_MageCurrentCombustionDotEstimate:
		return mage.newValueCurrentCombustionDotEstimate(config.GetMageCurrentCombustionDotEstimate(), config.Uuid)
	default:
		return nil
	}
}

type APLValueMageCurrentCombustionDotEstimate struct {
	core.DefaultAPLValueImpl
	mage                  *FireMage
	combustionDotEstimate int32
}

func (mage *FireMage) newValueCurrentCombustionDotEstimate(_ *proto.APLValueMageCurrentCombustionDotEstimate, _ *proto.UUID) core.APLValue {
	if mage.Spec != proto.Spec_SpecFireMage {
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

	if value.mage.combustionDotEstimate != value.combustionDotEstimate {
		value.combustionDotEstimate = value.mage.combustionDotEstimate
		if sim.Log != nil {
			value.mage.Log(sim, "Combustion Dot Estimate: %d", value.combustionDotEstimate)
		}
	}

	return value.combustionDotEstimate
}

func (value *APLValueMageCurrentCombustionDotEstimate) String() string {
	return "Combustion Dot Estimated Value"
}
