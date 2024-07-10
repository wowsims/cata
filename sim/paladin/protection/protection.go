package protection

import (
	"math"
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
	"github.com/wowsims/cata/sim/paladin"
)

func RegisterProtectionPaladin() {
	core.RegisterAgentFactory(
		proto.Player_ProtectionPaladin{},
		proto.Spec_SpecProtectionPaladin,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewProtectionPaladin(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_ProtectionPaladin)
			if !ok {
				panic("Invalid spec value for Protection Paladin!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewProtectionPaladin(character *core.Character, options *proto.Player) *ProtectionPaladin {
	protOptions := options.GetProtectionPaladin()

	prot := &ProtectionPaladin{
		Paladin:   paladin.NewPaladin(character, options.TalentsString, protOptions.Options.ClassOptions),
		Options:   protOptions.Options,
		vengeance: &core.VengeanceTracker{},
	}

	healingModel := options.HealingModel
	if healingModel != nil {
		if healingModel.InspirationUptime > 0.0 {
			core.ApplyInspiration(&prot.Unit, healingModel.InspirationUptime)
		}
	}

	return prot
}

type ProtectionPaladin struct {
	*paladin.Paladin

	Options *proto.ProtectionPaladin_Options

	vengeance *core.VengeanceTracker
}

func (prot *ProtectionPaladin) GetPaladin() *paladin.Paladin {
	return prot.Paladin
}

func (prot *ProtectionPaladin) Initialize() {
	prot.Paladin.Initialize()
	prot.ActivateRighteousFury()
	prot.registerAvengersShieldSpell()
	prot.RegisterSpecializationEffects()
}

func (prot *ProtectionPaladin) ApplyTalents() {
	prot.Paladin.ApplyTalents()
	prot.ApplyArmorSpecializationEffect(stats.Stamina, proto.ArmorType_ArmorTypePlate)
}

func (prot *ProtectionPaladin) Reset(sim *core.Simulation) {
	prot.Paladin.Reset(sim)
	prot.RighteousFuryAura.Activate(sim)
}

func (prot *ProtectionPaladin) RegisterSpecializationEffects() {
	// Divine Bulwark
	prot.RegisterMastery()

	// Touched by the Light
	prot.AddStatDependency(stats.Strength, stats.SpellPower, 0.6)
	prot.AddStat(stats.SpellHit, core.SpellHitRatingPerHitChance*8)
	prot.MultiplyStat(stats.Stamina, 1.15)
	core.MakePermanent(prot.GetOrRegisterAura(core.Aura{
		Label:    "Touched by the Light",
		ActionID: core.ActionID{SpellID: 53592},
	}))

	// Judgements of the Wise
	prot.ApplyJudgementsOfTheWise()

	// Vengeance
	core.ApplyVengeanceEffect(&prot.Character, prot.vengeance, 84839)
}

func (prot *ProtectionPaladin) RegisterMastery() {
	// Divine Bulwark
	masteryBlockChance := 18.0 + prot.GetMasteryPoints()*2.25
	prot.AddStat(stats.Block, masteryBlockChance*core.BlockRatingPerBlockChance)

	// Keep it updated when mastery changes
	prot.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery, newMastery float64) {
		oldBlockRating := (2.25 * core.MasteryRatingToMasteryPoints(oldMastery)) * core.BlockRatingPerBlockChance
		newBlockRating := (2.25 * core.MasteryRatingToMasteryPoints(newMastery)) * core.BlockRatingPerBlockChance

		prot.AddStatDynamic(sim, stats.Block, -oldBlockRating+newBlockRating)
	})
}

func (prot *ProtectionPaladin) ApplyJudgementsOfTheWise() {
	actionID := core.ActionID{SpellID: 31878}
	manaMetrics := prot.NewManaMetrics(actionID)

	// It's 30% of base mana over 10 seconds, with haste adding ticks.
	manaPerTick := math.Round(0.030 * prot.BaseMana)

	jotw := prot.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagHelpful | core.SpellFlagNoMetrics | core.SpellFlagNoLogs,

		Hot: core.DotConfig{
			SelfOnly: true,
			Aura: core.Aura{
				Label: "Judgements of the Wise",
			},
			NumberOfTicks:        10,
			TickLength:           time.Second * 1,
			AffectedByCastSpeed:  true,
			HasteReducesDuration: false,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				prot.AddMana(sim, manaPerTick, manaMetrics)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.SelfHot().Apply(sim)
		},
	})

	core.MakeProcTriggerAura(&prot.Unit, core.ProcTrigger{
		Name:           "Judgements of the Wise Trigger",
		ActionID:       actionID,
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: paladin.SpellMaskJudgement,
		ProcChance:     1.0,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			jotw.Cast(sim, &prot.Unit)
		},
	})
}
