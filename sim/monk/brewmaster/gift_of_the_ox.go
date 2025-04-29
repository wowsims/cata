package brewmaster

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/monk"
)

func (bm *BrewmasterMonk) registerGiftOfTheOx() {
	giftOfTheOxPassiveActionID := core.ActionID{SpellID: 124502}
	giftOfTheOxHealActionID := core.ActionID{SpellID: 124507}

	giftOfTheOxStackingAura := bm.RegisterAura(core.Aura{
		Label:     "Gift Of The Ox" + bm.Label,
		ActionID:  giftOfTheOxPassiveActionID,
		Duration:  time.Minute * 1,
		MaxStacks: 3,
	})

	bm.RegisterSpell(core.SpellConfig{
		ActionID:       giftOfTheOxHealActionID,
		ClassSpellMask: monk.MonkSpellChiSphere,
		SpellSchool:    core.SpellSchoolNature,
		Flags:          core.SpellFlagPassiveSpell | core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,
		ProcMask:       core.ProcMaskSpellHealing,

		DamageMultiplier: 1,
		CritMultiplier:   1,

		BonusCoefficient: 0.5025,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return giftOfTheOxStackingAura.IsActive() && giftOfTheOxStackingAura.GetStacks() > 0
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			heal := bm.CalcScalingSpellDmg(4.5) + spell.MeleeAttackPower()*spell.BonusCoefficient
			spell.CalcAndDealHealing(sim, spell.Unit, heal, spell.OutcomeHealing)
			giftOfTheOxStackingAura.RemoveStack(sim)
		},
		RelatedSelfBuff: giftOfTheOxStackingAura,
	})

	core.MakeProcTriggerAura(&bm.Unit, core.ProcTrigger{
		Name:     "Gift of The Ox Proc",
		ActionID: giftOfTheOxPassiveActionID,
		ProcMask: core.ProcMaskMelee,
		Callback: core.CallbackOnSpellHitDealt,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			procChance := 0.0
			if spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				if bm.HandType == proto.HandType_HandTypeOneHand {
					procChance = 0.03 * bm.MainHand().SwingSpeed
				} else {
					procChance = 0.06 * bm.MainHand().SwingSpeed
				}
			} else if spell.ProcMask.Matches(core.ProcMaskMeleeSpecial) {
				if bm.HandType == proto.HandType_HandTypeOneHand {
					procChance = bm.MainHand().SwingSpeed / 50
				} else {
					procChance = bm.MainHand().SwingSpeed / 25
				}
			}

			if sim.Proc(procChance, "Gift of The Ox") {
				giftOfTheOxStackingAura.Activate(sim)
				giftOfTheOxStackingAura.AddStack(sim)
			}
		},
	})

}
