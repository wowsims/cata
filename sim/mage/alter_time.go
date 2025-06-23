package mage

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (mage *Mage) registerAlterTimeCD() {

	auraState := map[string]*core.AuraState{}
	var allAuras []*core.Aura
	actionID := core.ActionID{SpellID: 108978}
	mageSavedMana := 0.0
	mageSavedHitPoints := 0.0
	manaMetrics := mage.NewManaMetrics(actionID.WithTag(1))
	healthMetrics := mage.NewHealthMetrics(actionID.WithTag(1))

	restoreState := func(sim *core.Simulation) {
		if manaDiff := mage.CurrentMana() - mageSavedMana; manaDiff > 0.0 {
			mage.SpendMana(sim, math.Abs(manaDiff), manaMetrics)
		} else {
			mage.AddMana(sim, math.Abs(manaDiff), manaMetrics)
		}

		if healthDiff := mage.CurrentHealth() - mageSavedHitPoints; healthDiff > 0.0 {
			mage.RemoveHealth(sim, math.Abs(healthDiff))
		} else {
			mage.GainHealth(sim, math.Abs(healthDiff), healthMetrics)
		}

		for _, aura := range allAuras {
			state := auraState[aura.Label]
			if state != nil {
				aura.RestoreState(*state, sim)
			} else if aura.IsActive() {
				aura.Deactivate(sim)
			}
		}
	}

	mage.AlterTimeAura = mage.RegisterAura(core.Aura{
		Label:    "Alter Time",
		ActionID: actionID,
		Duration: time.Second * 6,
		OnInit: func(alterTimeAura *core.Aura, sim *core.Simulation) {
			allAuras = core.FilterSlice(mage.GetAuras(), func(aura *core.Aura) bool {
				return aura.Duration != core.NeverExpires && aura != alterTimeAura
			})
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			mageSavedMana = mage.CurrentMana()
			mageSavedHitPoints = mage.CurrentHealth()
			for _, aura := range allAuras {
				if aura.IsActive() {
					state := aura.SaveState(sim)
					auraState[aura.Label] = &state
				} else {
					auraState[aura.Label] = nil
				}
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if aura.StartedAt()+aura.Duration <= sim.CurrentTime {
				restoreState(sim)
			}
		},
	})

	mage.AlterTime = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 108978},
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

		RelatedSelfBuff: mage.AlterTimeAura,
	})

	mage.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 108978}.WithTag(1),
		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return mage.AlterTimeAura.IsActive()
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			mage.AlterTimeAura.Deactivate(sim)
			restoreState(sim)
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: mage.AlterTime,
		Type:  core.CooldownTypeDPS,
	})
}
