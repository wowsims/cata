//go:build wasm

package core

import (
	"github.com/wowsims/cata/sim/core/proto"
)

// Note: WASM can't do threads with go, so there's no reason to even compile the whole concurrency code. Instead just run sims directly.

func RunConcurrentRaidSimAsync(request *proto.RaidSimRequest, progress chan *proto.ProgressMetrics) {
	go RunSim(request, progress, nil)
}

func RunConcurrentRaidSimSync(request *proto.RaidSimRequest) *proto.RaidSimResult {
	return RunSim(request, nil, nil)
}
