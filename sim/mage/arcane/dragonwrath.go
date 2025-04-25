package arcane

import (
	cata "github.com/wowsims/mop/sim/common/cata"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	cata.CreateDTRClassConfig(proto.Spec_SpecArcaneMage, 1.0/12)
}
