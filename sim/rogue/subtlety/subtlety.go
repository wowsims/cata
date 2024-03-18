package subtlety

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/rogue"
)

func RegisterSubtletyRogue() {
	core.RegisterAgentFactory(
		proto.Player_SubtletyRogue{},
		proto.Spec_SpecSubtletyRogue,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewSubtletyRogue(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_SubtletyRogue)
			if !ok {
				panic("Invalid spec value for Subtlety Rogue!")
			}
			player.Spec = playerSpec
		},
	)
}

func (subRogue *SubtletyRogue) Initialize() {
	subRogue.Rogue.Initialize()

	subRogue.registerHemorrhageSpell()
	subRogue.registerSanguinaryVein()
}

func NewSubtletyRogue(character *core.Character, options *proto.Player) *SubtletyRogue {
	subOptions := options.GetSubtletyRogue().Options

	subRogue := &SubtletyRogue{
		Rogue: rogue.NewRogue(character, subOptions.ClassOptions, options.TalentsString),
	}
	subRogue.SubtletyOptions = subOptions

	return subRogue
}

type SubtletyRogue struct {
	*rogue.Rogue
}

func (subRogue *SubtletyRogue) GetRogue() *rogue.Rogue {
	return subRogue.Rogue
}

func (subRogue *SubtletyRogue) Reset(sim *core.Simulation) {
	subRogue.Rogue.Reset(sim)
}

// func (rogue *Rogue) applyInitiative() {
// 	if rogue.Talents.Initiative == 0 {
// 		return
// 	}

// 	procChance := 0.5*float65(rogue.Talents.Initiative)
// 	cpMetrics := rogue.NewComboPointMetrics(core.ActionID{SpellID: 13980})

// 	rogue.RegisterAura(core.Aura{
// 		Label:    "Initiative",
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			if spell == rogue.Garrote || spell == rogue.Ambush {
// 				if result.Landed() {
// 					if sim.Proc(procChance, "Initiative") {
// 						rogue.AddComboPoints(sim, 1, cpMetrics)
// 					}
// 				}
// 			}
// 		},
// 	})
// }

// func (rogue *Rogue) registerHonorAmongThieves() {
// 	// When anyone in your group critically hits with a damage or healing spell or ability,
// 	// you have a [33%/66%/100%] chance to gain a combo point on your current target.
// 	// This effect cannot occur more than once per second.
// 	if rogue.Talents.HonorAmongThieves == 0 {
// 		return
// 	}

// 	procChance := []float64{0, 0.33, 0.66, 1}[rogue.Talents.HonorAmongThieves]
// 	comboMetrics := rogue.NewComboPointMetrics(core.ActionID{SpellID: 51701})
// 	honorAmongThievesID := core.ActionID{SpellID: 51701}

// 	icd := core.Cooldown{
// 		Timer:    rogue.NewTimer(),
// 		Duration: time.Second,
// 	}

// 	maybeProc := func(sim *core.Simulation) {
// 		if icd.IsReady(sim) && sim.Proc(procChance, "honor of thieves") {
// 			rogue.AddComboPoints(sim, 1, comboMetrics)
// 			icd.Use(sim)
// 		}
// 	}

// 	rogue.HonorAmongThieves = rogue.RegisterAura(core.Aura{
// 		Label:    "Honor Among Thieves",
// 		ActionID: honorAmongThievesID,
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnGain: func(_ *core.Aura, sim *core.Simulation) {
// 			// In an ideal party, you'd probably get up to 6 ability crits/s (Rate = 600).
// 			//  Survival Hunters, Enhancement Shamans, and Assassination Rogues are particularly good.
// 			if rogue.SubtletyOptions.HonorAmongThievesCritRate <= 0 {
// 				return
// 			}

// 			if rogue.SubtletyOptions.HonorAmongThievesCritRate > 2000 {
// 				rogue.SubtletyOptions.HonorAmongThievesCritRate = 2000 // limited, so performance doesn't suffer
// 			}

// 			rateToDuration := float64(time.Second) * 100 / float64(rogue.SubtletyOptions.HonorAmongThievesCritRate)

// 			pa := &core.PendingAction{}
// 			pa.OnAction = func(sim *core.Simulation) {
// 				maybeProc(sim)
// 				pa.NextActionAt = sim.CurrentTime + time.Duration(sim.RandomExpFloat("next party crit")*rateToDuration)
// 				sim.AddPendingAction(pa)
// 			}
// 			pa.NextActionAt = sim.CurrentTime + time.Duration(sim.RandomExpFloat("next party crit")*rateToDuration)
// 			sim.AddPendingAction(pa)
// 		},
// 		OnSpellHitDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			if result.DidCrit() && !spell.ProcMask.Matches(core.ProcMaskMeleeMHAuto|core.ProcMaskMeleeOHAuto|core.ProcMaskRangedAuto) {
// 				maybeProc(sim)
// 			}
// 		},
// 		OnPeriodicDamageDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			if result.DidCrit() {
// 				maybeProc(sim)
// 			}
// 		},
// 	})
// }

func (subRogue *SubtletyRogue) registerSanguinaryVein() {
	if subRogue.Talents.SanguinaryVein == 0 {
		return
	}

	svBonus := 1 + 0.08*float64(subRogue.Talents.SanguinaryVein)
	svSpellID := core.TernaryInt32(subRogue.Talents.SanguinaryVein == 1, 79146, 79147)

	svDebuffArray := subRogue.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Sanguinary Vein Debuff",
			ActionID: core.ActionID{SpellID: svSpellID},
			Duration: core.NeverExpires,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				subRogue.AttackTables[aura.Unit.UnitIndex].DamageTakenMultiplier *= svBonus
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				subRogue.AttackTables[aura.Unit.UnitIndex].DamageTakenMultiplier /= svBonus
			},
		})
	})

	subRogue.Env.RegisterPreFinalizeEffect(func() {
		if subRogue.Rupture != nil {
			subRogue.Rupture.RelatedAuras = append(subRogue.Rupture.RelatedAuras, svDebuffArray)
		}
		if subRogue.Hemorrhage != nil && subRogue.HasPrimeGlyph(proto.RoguePrimeGlyph_GlyphOfHemorrhage) {
			subRogue.Hemorrhage.RelatedAuras = append(subRogue.Hemorrhage.RelatedAuras, svDebuffArray)
		}
	})

	subRogue.RegisterAura(core.Aura{
		Label:    "Sanguinary Vein Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			if spell == subRogue.Rupture || (spell == subRogue.Hemorrhage && subRogue.HasPrimeGlyph(proto.RoguePrimeGlyph_GlyphOfHemorrhage)) {
				aura := svDebuffArray.Get(result.Target)
				dot := spell.Dot(result.Target)
				aura.Duration = dot.TickLength * time.Duration(dot.NumberOfTicks)
				aura.Activate(sim)
			}
		},
	})
}
