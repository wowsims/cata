package brewmaster

import (
	"fmt"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (monk *BrewmasterMonk) NewAPLValue(_ *core.APLRotation, config *proto.APLValue) core.APLValue {
	switch config.Value.(type) {
	case *proto.APLValue_BrewmasterMonkCurrentStaggerPercent:
		return monk.newValueCurrentStaggerPercent(config.GetBrewmasterMonkCurrentStaggerPercent(), config.Uuid)
	default:
		return nil
	}
}

type APLValueBrewmasterMonkCurrentStaggerPercent struct {
	core.DefaultAPLValueImpl
	monk *BrewmasterMonk
	aura *core.Aura
}

func (monk *BrewmasterMonk) newValueCurrentStaggerPercent(config *proto.APLValueBrewmasterMonkCurrentStaggerPercent, uuid *proto.UUID) core.APLValue {
	return &APLValueBrewmasterMonkCurrentStaggerPercent{
		monk: monk,
		aura: monk.Stagger.SelfHot().Aura,
	}
}
func (value *APLValueBrewmasterMonkCurrentStaggerPercent) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueBrewmasterMonkCurrentStaggerPercent) GetFloat(sim *core.Simulation) float64 {
	return float64(value.aura.GetStacks()) / value.monk.MaxHealth()
}
func (value *APLValueBrewmasterMonkCurrentStaggerPercent) String() string {
	return fmt.Sprintf("Current Stagger %%")
}
