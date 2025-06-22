package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (mage *Mage) registerAlterTimeCD() {

	auraState := map[int32]core.AuraState{}
	var allAuras []*core.Aura
	actionID := core.ActionID{SpellID: 108978}
	mageCurrentMana := 0.0
	mageCurrentHitpoints := 0.0
	manaMetrics := mage.NewManaMetrics(actionID)
	healthMetrics := mage.NewHealthMetrics(actionID)

	mage.AlterTimeAura = mage.RegisterAura(core.Aura{
		Label:    "Alter Time",
		ActionID: actionID,
		Duration: time.Second * 6,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			allAuras = core.FilterSlice(mage.GetAuras(), func(aura *core.Aura) bool {
				return aura.Duration == core.NeverExpires
			})
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			mageCurrentMana = mage.CurrentMana()
			mageCurrentHitpoints = mage.CurrentHealth()
			for _, aura := range allAuras {
				auraState[aura.ActionID.SpellID] = aura.SaveState(sim)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			mage.SetCurrentMana(sim, mageCurrentMana, manaMetrics)
			mage.SetCurrentHealth(sim, mageCurrentHitpoints, healthMetrics)
			for _, aura := range allAuras {
				state := auraState[aura.ActionID.SpellID]
				if state != aura.SaveState(sim) {
					aura.RestoreState(state, sim)
				}
			}
		},
	})

	mage.AlterTime = mage.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: MageSpellAlterTime,
		Flags:          core.SpellFlagNoOnCastComplete,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 1,
		},

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Minute * 3,
			},
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			mage.AlterTimeAura.Activate(sim)
			mage.WaitUntil(sim, sim.CurrentTime+mage.ReactionTime)
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: mage.AlterTime,
		Type:  core.CooldownTypeDPS,
	})
}
