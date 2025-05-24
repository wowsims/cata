package brewmaster

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/monk"
)

func (bm *BrewmasterMonk) registerPassives() {
	bm.registerBrewmasterTraining()
	bm.registerElusiveBrew()
	bm.registerGiftOfTheOx()
	bm.registerDesperateMeasures()
}

func (bm *BrewmasterMonk) registerBrewmasterTraining() {
	// Fortifying Brew
	// Also increases your Stagger amount by 20% while active.
	// Fortifying Brew Stagger mod is implemented in stagger.go

	// Tiger Palm
	// Tiger Palm no longer costs Chi, and when you deal damage with Tiger Palm the amount of your next Guard is increased by 15%. Lasts 30 sec.
	// Tiger Palm Chi mod is implemented in tiger_palm.go
	bm.PowerGuardAura = bm.RegisterAura(core.Aura{
		Label:    "Power Guard",
		ActionID: core.ActionID{SpellID: 118636},
		Duration: 30 * time.Second,
	})

	core.MakeProcTriggerAura(&bm.Unit, core.ProcTrigger{
		Name:           "Power Guard Trigger",
		ClassSpellMask: monk.MonkSpellTigerPalm,
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			bm.PowerGuardAura.Activate(sim)
		},
	})

	// Blackout Kick
	// After you Blackout Kick, you gain Shuffle, increasing your parry chance by 20%
	// and your Stagger amount by an additional 20% for 6 sec.
	// Stagger amount is implemented in stagger.go
	bm.ShuffleAura = bm.RegisterAura(core.Aura{
		Label:    "Shuffle",
		ActionID: core.ActionID{SpellID: 115307},
		Duration: 6 * time.Second,
	}).AttachAdditivePseudoStatBuff(&bm.PseudoStats.BaseParryChance, 0.2)

	core.MakeProcTriggerAura(&bm.Unit, core.ProcTrigger{
		Name:           "Shuffle Trigger",
		ClassSpellMask: monk.MonkSpellBlackoutKick,
		Callback:       core.CallbackOnCastComplete,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if bm.ShuffleAura.IsActive() {
				bm.ShuffleAura.UpdateExpires(bm.ShuffleAura.ExpiresAt() + 6*time.Second)
			} else {
				bm.ShuffleAura.Activate(sim)
			}
		},
	})
}

func (bm *BrewmasterMonk) registerElusiveBrew() {
	buffActionID := core.ActionID{SpellID: 115308}
	stackActionID := core.ActionID{SpellID: 128938}

	stackingAura := core.MakePermanent(bm.RegisterAura(core.Aura{
		Label:     "Brewing: Elusive Brew" + bm.Label,
		ActionID:  stackActionID,
		Duration:  30 * time.Second,
		MaxStacks: 15,
	}))

	bm.Monk.RegisterOnNewBrewStacks(func(sim *core.Simulation, stacksToAdd int32) {
		stackingAura.Activate(sim)
		stackingAura.SetStacks(sim, stackingAura.GetStacks()+stacksToAdd)
	})

	bm.ElusiveBrewAura = bm.RegisterAura(core.Aura{
		Label:    "Elusive Brew" + bm.Label,
		ActionID: buffActionID,
		Duration: 0,
	}).AttachAdditivePseudoStatBuff(&bm.PseudoStats.BaseDodgeChance, 0.3)

	core.MakeProcTriggerAura(&bm.Unit, core.ProcTrigger{
		Name:     "Brewing: Elusive Brew Proc",
		ActionID: stackActionID,
		Outcome:  core.OutcomeCrit,
		ProcMask: core.ProcMaskMeleeWhiteHit,
		Callback: core.CallbackOnSpellHitDealt,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			stacks := 0.0
			if bm.HandType == proto.HandType_HandTypeOneHand {
				stacks = 1.5 * bm.MainHand().SwingSpeed / 2.6
			} else {
				stacks = 3 * bm.MainHand().SwingSpeed / 3.6
			}

			if sim.Proc(math.Mod(stacks, 1), "Brewing: Elusive Brew") {
				stacks = math.Ceil(stacks)
			} else {
				stacks = math.Floor(stacks)
			}

			stackingAura.Activate(sim)
			stackingAura.SetStacks(sim, stackingAura.GetStacks()+int32(stacks))
		},
	})

	spell := bm.RegisterSpell(core.SpellConfig{
		ActionID:       buffActionID,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,
		ClassSpellMask: monk.MonkSpellElusiveBrew,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    bm.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return bm.StanceMatches(monk.SturdyOx) && stackingAura.GetStacks() > 0
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			stacks := stackingAura.GetStacks()
			bm.ElusiveBrewAura.Duration = time.Duration(stacks) * time.Second
			bm.ElusiveBrewStacks = stacks

			bm.ElusiveBrewAura.Activate(sim)
			stackingAura.SetStacks(sim, 0)
		},
		RelatedSelfBuff: bm.ElusiveBrewAura,
	})

	bm.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeSurvival,
	})
}

func (bm *BrewmasterMonk) registerGiftOfTheOx() {
	giftOfTheOxPassiveActionID := core.ActionID{SpellID: 124502}
	giftOfTheOxHealActionID := core.ActionID{SpellID: 124507}

	hasGlyph := bm.HasMajorGlyph(proto.MonkMajorGlyph_GlyphOfEnduringHealingSphere)
	sphereDuration := time.Minute*1 + core.TernaryDuration(hasGlyph, time.Minute*3, 0)

	giftOfTheOxStackingAura := bm.RegisterAura(core.Aura{
		Label:     "Gift Of The Ox" + bm.Label,
		ActionID:  giftOfTheOxPassiveActionID,
		Duration:  sphereDuration,
		MaxStacks: math.MaxInt32,
	})

	bm.RegisterSpell(core.SpellConfig{
		ActionID:       giftOfTheOxHealActionID,
		ClassSpellMask: monk.MonkSpellGiftOfTheOx,
		SpellSchool:    core.SpellSchoolNature,
		Flags:          core.SpellFlagPassiveSpell | core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,
		ProcMask:       core.ProcMaskSpellHealing,

		DamageMultiplier: 1,
		CritMultiplier:   1,
		ThreatMultiplier: 1,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return giftOfTheOxStackingAura.IsActive() && giftOfTheOxStackingAura.GetStacks() > 0
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			heal := bm.CalcScalingSpellDmg(4.5) + spell.MeleeAttackPower()*0.2508
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

			// Source:
			// https://www.wowhead.com/blue-tracker/topic/beta-class-balance-analysis-pt-ii-6397900436#2
			// https://www.wowhead.com/blue-tracker/topic/beta-class-balance-analysis-5889309137#59115992048
			// https://web.archive.org/web/20130801205930/http://elitistjerks.com/f99/t131791-like_water_brewmasters_resource_8_1_13_a/#Gift_of_the_Ox
			if spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) || spell.Matches(monk.MonkSpellTigerStrikes) {
				weapon := core.Ternary(spell.IsMH(), bm.MainHand(), bm.OffHand())
				procChance = core.Ternary(bm.HandType == proto.HandType_HandTypeOneHand, 0.051852, 0.06) * weapon.SwingSpeed
			} else if spell.ProcMask.Matches(core.ProcMaskMeleeSpecial) {
				procChance = 0.10
			}

			if sim.Proc(procChance, "Gift of The Ox") {
				giftOfTheOxStackingAura.Activate(sim)
				giftOfTheOxStackingAura.AddStack(sim)
			}
		},
	})

}

func (bm *BrewmasterMonk) registerDesperateMeasures() {
	actionID := core.ActionID{SpellID: 126060}

	aura := bm.RegisterAura(core.Aura{
		Label:    "Desperate Measures" + bm.Label,
		ActionID: actionID,
		Duration: core.NeverExpires,
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask:  monk.MonkSpellExpelHarm,
		Kind:       core.SpellMod_Cooldown_Multiplier,
		FloatValue: -1,
	})

	core.MakeProcTriggerAura(&bm.Unit, core.ProcTrigger{
		Name:     "Desperate Measures Health Monitor" + bm.Label,
		ActionID: actionID,
		Duration: 0,
		Outcome:  core.OutcomeHit,
		Callback: core.CallbackOnSpellHitTaken,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Target.CurrentHealthPercent() <= 0.35 {
				aura.Activate(sim)
			} else {
				aura.Deactivate(sim)
			}
		},
	})
}
