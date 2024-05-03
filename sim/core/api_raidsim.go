//go:build wasm

package core

import (
	"github.com/wowsims/cata/sim/core/proto"
)

func runConcurrentSim(request *proto.RaidSimRequest, progress chan *proto.ProgressMetrics) {
	go RunSim(request, progress)
}
