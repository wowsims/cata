package encounters

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/encounters/bwd"
	"github.com/wowsims/cata/sim/encounters/icc"
	"github.com/wowsims/cata/sim/encounters/naxxramas"
	"github.com/wowsims/cata/sim/encounters/toc"
	"github.com/wowsims/cata/sim/encounters/ulduar"
)

func init() {
	naxxramas.Register()
	ulduar.Register()
	toc.Register()
	icc.Register()
	bwd.Register()
}

func AddSingleTargetBossEncounter(presetTarget *core.PresetTarget) {
	core.AddPresetTarget(presetTarget)
	core.AddPresetEncounter(presetTarget.Config.Name, []string{
		presetTarget.Path(),
	})
}
