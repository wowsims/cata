package mage

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (mage *Mage) NewAPLValue(rot *core.APLRotation, config *proto.APLValue) core.APLValue {
	switch config.Value.(type) {
	case *proto.APLValue_MageCurrentCombustionDotEstimate:
		return mage.newValueCurrentCombustionDotEstimate(config.GetMageCurrentCombustionDotEstimate())
	default:
		return nil
	}
}

type APLValueMageCurrentCombustionDotEstimate struct {
	core.DefaultAPLValueImpl
	mage *Mage
}

func (mage *Mage) newValueCurrentCombustionDotEstimate(_ *proto.APLValueMageCurrentCombustionDotEstimate) core.APLValue {
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
	mage := value.mage

	if !mage.Combustion.IsReady(sim) {
		return int32(0)
	}

	dot := mage.Combustion.RelatedDotSpell.Dot(mage.CurrentTarget)
	combustionDotDamage := 0.0
	tickCount := int(dot.ExpectedTickCount())

	for i := 0; i < tickCount; i++ {
		damage := mage.Combustion.RelatedDotSpell.ExpectedTickDamage(sim, mage.CurrentTarget)
		combustionDotDamage += damage
	}

	return int32(combustionDotDamage)
}

func (value *APLValueMageCurrentCombustionDotEstimate) String() string {
	return "Combustion Dot Estimated Value"
}
