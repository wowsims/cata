package protection

import (
	"fmt"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (prot *ProtectionPaladin) NewAPLValue(_ *core.APLRotation, config *proto.APLValue) core.APLValue {
	switch config.Value.(type) {
	case *proto.APLValue_ProtectionPaladinDamageTakenLastGlobal:
		return prot.newValueDamageTakenLastGlobal(config.GetProtectionPaladinDamageTakenLastGlobal(), config.Uuid)
	default:
		return nil
	}
}

type APLValueProtectionPaladinDamageTakenLastGlobal struct {
	core.DefaultAPLValueImpl
	prot *ProtectionPaladin
}

func (prot *ProtectionPaladin) newValueDamageTakenLastGlobal(_ *proto.APLValueProtectionPaladinDamageTakenLastGlobal, _ *proto.UUID) core.APLValue {
	return &APLValueProtectionPaladinDamageTakenLastGlobal{
		prot: prot,
	}
}
func (value *APLValueProtectionPaladinDamageTakenLastGlobal) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueProtectionPaladinDamageTakenLastGlobal) GetFloat(_ *core.Simulation) float64 {
	return value.prot.DamageTakenLastGlobal
}
func (value *APLValueProtectionPaladinDamageTakenLastGlobal) String() string {
	return fmt.Sprintf("Damage Taken Last Global(%f)", value.prot.DamageTakenLastGlobal)
}
