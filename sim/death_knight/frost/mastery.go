package frost

import "github.com/wowsims/mop/sim/core"

// Increases all Frost damage done by (16 + (<Mastery Rating>/600)*2)%.
func (fdk *FrostDeathKnight) registerMastery() {
	masteryMod := fdk.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		School:     core.SpellSchoolFrost,
		FloatValue: fdk.getMasteryPercent(),
	})

	fdk.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery float64, newMastery float64) {
		masteryMod.UpdateFloatValue(fdk.getMasteryPercent())
	})

	core.MakePermanent(fdk.RegisterAura(core.Aura{
		Label:    "Frozen Heart" + fdk.Label,
		ActionID: core.ActionID{SpellID: 77514},

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			masteryMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			masteryMod.Deactivate()
		},
	}))
}

func (fdk *FrostDeathKnight) getMasteryPercent() float64 {
	return 0.16 + 0.02*fdk.GetMasteryPoints()
}
