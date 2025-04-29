package brewmaster

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/monk"
)

func (bm *BrewmasterMonk) registerElusiveBrew() {
	buffActionID := core.ActionID{SpellID: 115308}
	stackActionID := core.ActionID{SpellID: 128938}

	stackingAura := core.MakePermanent(bm.RegisterAura(core.Aura{
		Label:     "Brewing: Elusive Brew" + bm.Label,
		ActionID:  stackActionID,
		Duration:  core.NeverExpires,
		MaxStacks: 15,
	}))

	buffAura := bm.RegisterAura(core.Aura{
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
			buffAura.Duration = time.Duration(stackingAura.GetStacks()) * time.Second
			buffAura.Activate(sim)
			stackingAura.SetStacks(sim, 0)
		},
	})

	bm.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeSurvival,
	})
}
