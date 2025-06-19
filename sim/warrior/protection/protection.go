package protection

import (
	"math"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/warrior"
)

func RegisterProtectionWarrior() {
	core.RegisterAgentFactory(
		proto.Player_ProtectionWarrior{},
		proto.Spec_SpecProtectionWarrior,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewProtectionWarrior(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_ProtectionWarrior)
			if !ok {
				panic("Invalid spec value for Protection Warrior!")
			}
			player.Spec = playerSpec
		},
	)
}

type ProtectionWarrior struct {
	*warrior.Warrior

	Options *proto.ProtectionWarrior_Options
}

func NewProtectionWarrior(character *core.Character, options *proto.Player) *ProtectionWarrior {
	protOptions := options.GetProtectionWarrior().Options

	war := &ProtectionWarrior{
		Warrior: warrior.NewWarrior(character, protOptions.ClassOptions, options.TalentsString, warrior.WarriorInputs{}),
		Options: protOptions,
	}

	return war
}

func (war *ProtectionWarrior) CalculateMasteryBlockChance() float64 {
	return math.Floor(0.5*(8.0+war.GetMasteryPoints())) / 100.0
}

func (war *ProtectionWarrior) CalculateMasteryCriticalBlockChance() float64 {
	return math.Floor(2.2*(8.0+war.GetMasteryPoints())) / 100.0
}

func (war *ProtectionWarrior) GetWarrior() *warrior.Warrior {
	return war.Warrior
}

func (war *ProtectionWarrior) Initialize() {
	war.Warrior.Initialize()
	war.registerPassives()

	war.registerDevastate()
	war.registerRevenge()
	war.registerShieldSlam()
	war.registerShieldBlock()
	war.registerShieldBarrier()
	war.registerDemoralizingShout()
	war.registerLastStand()
}

func (war *ProtectionWarrior) registerPassives() {
	war.ApplyArmorSpecializationEffect(stats.Stamina, proto.ArmorType_ArmorTypePlate, 86526)

	// Critical block
	war.registerMastery()

	war.registerUnwaveringSentinel()
	war.registerBastionOfDefense()
	war.registerSwordAndBoard()
	war.registerRiposte()

	// Vengeance
	war.RegisterVengeance(93098, war.DefensiveStanceAura)
}

func (war *ProtectionWarrior) registerMastery() {

	dummyCriticalBlockSpell := war.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 76857}, // Doesn't seem like there's an actual spell ID for the block itself, so use the mastery ID
		Flags:    core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,
	})

	war.Blockhandler = func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
		if !spell.SpellSchool.Matches(core.SpellSchoolPhysical) {
			return
		}

		if result.Outcome.Matches(core.OutcomeBlock) && !result.Outcome.Matches(core.OutcomeMiss) && !result.Outcome.Matches(core.OutcomeParry) && !result.Outcome.Matches(core.OutcomeDodge) {
			procChance := war.GetCriticalBlockChance()
			if sim.Proc(procChance, "Critical Block Roll") {
				result.Damage = result.Damage * (1 - war.BlockDamageReduction()*2)
				dummyCriticalBlockSpell.Cast(sim, spell.Unit)
				return
			}
			result.Damage = result.Damage * (1 - war.BlockDamageReduction())
		}
	}

	// Crit block mastery also applies an equal amount to regular block
	// set initial block % from both Masteries
	war.CriticalBlockChance[0] = war.CalculateMasteryCriticalBlockChance()
	war.CriticalBlockChance[1] = war.CalculateMasteryBlockChance()
	war.AddStat(stats.BlockPercent, (war.CriticalBlockChance[0]+war.CriticalBlockChance[1])*100.0)

	// and keep it updated when mastery changes
	war.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMasteryRating float64, newMasteryRating float64) {
		war.CriticalBlockChance[0] = war.CalculateMasteryCriticalBlockChance()
		war.CriticalBlockChance[1] = war.CalculateMasteryBlockChance()
		masteryBlockStat := 0.5 * core.MasteryRatingToMasteryPoints(newMasteryRating-oldMasteryRating)
		masteryCriticalBlockStat := 2.2 * core.MasteryRatingToMasteryPoints(newMasteryRating-oldMasteryRating)
		war.AddStatDynamic(sim, stats.BlockPercent, masteryCriticalBlockStat+masteryBlockStat)
	})
}

func (war *ProtectionWarrior) Reset(sim *core.Simulation) {
	war.Warrior.Reset(sim)
}
