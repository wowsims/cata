package sim

import (
	"github.com/wowsims/mop/sim/common"
	"github.com/wowsims/mop/sim/death_knight/blood"
	frostDeathKnight "github.com/wowsims/mop/sim/death_knight/frost"
	"github.com/wowsims/mop/sim/death_knight/unholy"
	"github.com/wowsims/mop/sim/druid/balance"
	"github.com/wowsims/mop/sim/druid/feral"
	"github.com/wowsims/mop/sim/druid/guardian"
	restoDruid "github.com/wowsims/mop/sim/druid/restoration"
	_ "github.com/wowsims/mop/sim/encounters"
	"github.com/wowsims/mop/sim/hunter/beast_mastery"
	"github.com/wowsims/mop/sim/hunter/marksmanship"
	"github.com/wowsims/mop/sim/hunter/survival"
	"github.com/wowsims/mop/sim/mage/arcane"
	"github.com/wowsims/mop/sim/mage/fire"
	frostMage "github.com/wowsims/mop/sim/mage/frost"
	"github.com/wowsims/mop/sim/monk/brewmaster"
	"github.com/wowsims/mop/sim/monk/mistweaver"
	"github.com/wowsims/mop/sim/monk/windwalker"
	holyPaladin "github.com/wowsims/mop/sim/paladin/holy"
	protPaladin "github.com/wowsims/mop/sim/paladin/protection"
	"github.com/wowsims/mop/sim/paladin/retribution"
	"github.com/wowsims/mop/sim/priest/discipline"
	holyPriest "github.com/wowsims/mop/sim/priest/holy"
	"github.com/wowsims/mop/sim/priest/shadow"
	"github.com/wowsims/mop/sim/rogue/assassination"
	"github.com/wowsims/mop/sim/rogue/combat"
	"github.com/wowsims/mop/sim/rogue/subtlety"
	"github.com/wowsims/mop/sim/shaman/elemental"
	"github.com/wowsims/mop/sim/shaman/enhancement"
	restoShaman "github.com/wowsims/mop/sim/shaman/restoration"
	"github.com/wowsims/mop/sim/warlock/affliction"
	"github.com/wowsims/mop/sim/warlock/demonology"
	"github.com/wowsims/mop/sim/warlock/destruction"
	"github.com/wowsims/mop/sim/warrior/arms"
	"github.com/wowsims/mop/sim/warrior/fury"
	protWarrior "github.com/wowsims/mop/sim/warrior/protection"
)

var registered = false

func RegisterAll() {
	if registered {
		return
	}
	registered = true

	blood.RegisterBloodDeathKnight()
	frostDeathKnight.RegisterFrostDeathKnight()
	unholy.RegisterUnholyDeathKnight()

	balance.RegisterBalanceDruid()
	feral.RegisterFeralDruid()
	guardian.RegisterGuardianDruid()
	restoDruid.RegisterRestorationDruid()

	beast_mastery.RegisterBeastMasteryHunter()
	marksmanship.RegisterMarksmanshipHunter()
	survival.RegisterSurvivalHunter()

	arcane.RegisterArcaneMage()
	fire.RegisterFireMage()
	frostMage.RegisterFrostMage()

	holyPaladin.RegisterHolyPaladin()
	protPaladin.RegisterProtectionPaladin()
	retribution.RegisterRetributionPaladin()

	discipline.RegisterDisciplinePriest()
	holyPriest.RegisterHolyPriest()
	shadow.RegisterShadowPriest()

	assassination.RegisterAssassinationRogue()
	combat.RegisterCombatRogue()
	subtlety.RegisterSubtletyRogue()

	elemental.RegisterElementalShaman()
	enhancement.RegisterEnhancementShaman()
	restoShaman.RegisterRestorationShaman()

	affliction.RegisterAfflictionWarlock()
	demonology.RegisterDemonologyWarlock()
	destruction.RegisterDestructionWarlock()

	arms.RegisterArmsWarrior()
	fury.RegisterFuryWarrior()
	protWarrior.RegisterProtectionWarrior()

	brewmaster.RegisterBrewmasterMonk()
	mistweaver.RegisterMistweaverMonk()
	windwalker.RegisterWindwalkerMonk()

	common.RegisterAllEffects()
}
