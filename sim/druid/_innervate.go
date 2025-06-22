package druid

import (
	"github.com/wowsims/mop/sim/core"
)

// Returns the time to wait before the next action, or 0 if innervate is on CD
// or disabled.
func (druid *Druid) registerInnervateCD() {
	innervateTarget := druid.GetUnit(druid.SelfBuffs.InnervateTarget)
	if innervateTarget == nil {
		return
	}
	innervateTargetChar := druid.Env.Raid.GetPlayerFromUnit(innervateTarget).GetCharacter()

	actionID := core.ActionID{SpellID: 29166, Tag: druid.Index}
	var innervateSpell *DruidSpell

	innervateCD := core.InnervateCD

	amount := 0.05
	if innervateTarget == &druid.Unit {
		amount = 0.2 + float64(druid.Talents.Dreamstate)*0.15
	}

	var innervateAura = core.InnervateAura(innervateTargetChar, actionID.Tag, amount)
	innervateManaThreshold := core.InnervateManaThreshold(innervateTargetChar)

	innervateSpell = druid.RegisterSpell(Humanoid|Moonkin|Tree, core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagReadinessTrinket,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: innervateCD,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			// If target already has another innervate, don't cast.
			return !innervateTarget.HasActiveAuraWithTag(core.InnervateAuraTag)
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			innervateAura.Activate(sim)
		},
	})

	druid.AddMajorCooldown(core.MajorCooldown{
		Spell: innervateSpell.Spell,
		Type:  core.CooldownTypeMana,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			// Innervate needs to be activated as late as possible to maximize DPS. The issue is that
			// innervate gives so much mana that it can cause Super Mana Potion or Dark Rune usages
			// to be delayed, if they come off CD soon after innervate. This delay is minimized by
			// activating innervate from the smallest amount of mana possible.
			return innervateTarget.CurrentMana() <= innervateManaThreshold
		},
	})
}
