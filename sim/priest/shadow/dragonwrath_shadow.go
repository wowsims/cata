package shadow

import (
	"github.com/wowsims/cata/sim/common/cata"
	"github.com/wowsims/cata/sim/core/proto"
)

func init() {
	cata.AddClassConfig(proto.Spec_SpecShadowPriest, 0.082).
		SupressImpact(2944) // Improved Devouring Plague

}
