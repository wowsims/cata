package sim

import (
	_ "github.com/wowsims/cata/sim/common"
	"github.com/wowsims/cata/sim/death_knight/blood"
	frostDeathKnight "github.com/wowsims/cata/sim/death_knight/frost"
	"github.com/wowsims/cata/sim/death_knight/unholy"
	"github.com/wowsims/cata/sim/druid/balance"
	"github.com/wowsims/cata/sim/druid/feral"
	"github.com/wowsims/cata/sim/druid/guardian"
	restoDruid "github.com/wowsims/cata/sim/druid/restoration"
	_ "github.com/wowsims/cata/sim/encounters"
	"github.com/wowsims/cata/sim/hunter/beast_mastery"
	"github.com/wowsims/cata/sim/hunter/marksmanship"
	"github.com/wowsims/cata/sim/hunter/survival"
	"github.com/wowsims/cata/sim/mage/arcane"
	"github.com/wowsims/cata/sim/mage/fire"
	frostMage "github.com/wowsims/cata/sim/mage/frost"
	holyPaladin "github.com/wowsims/cata/sim/paladin/holy"
	protPaladin "github.com/wowsims/cata/sim/paladin/protection"
	"github.com/wowsims/cata/sim/paladin/retribution"
	"github.com/wowsims/cata/sim/priest/discipline"
	holyPriest "github.com/wowsims/cata/sim/priest/holy"
	"github.com/wowsims/cata/sim/priest/shadow"
	"github.com/wowsims/cata/sim/rogue/assassination"
	"github.com/wowsims/cata/sim/rogue/combat"
	"github.com/wowsims/cata/sim/rogue/subtlety"
	"github.com/wowsims/cata/sim/shaman/elemental"
	"github.com/wowsims/cata/sim/shaman/enhancement"
	restoShaman "github.com/wowsims/cata/sim/shaman/restoration"
	"github.com/wowsims/cata/sim/warlock/affliction"
	"github.com/wowsims/cata/sim/warlock/demonology"
	"github.com/wowsims/cata/sim/warlock/destruction"
	"github.com/wowsims/cata/sim/warrior/arms"
	"github.com/wowsims/cata/sim/warrior/fury"
	protWarrior "github.com/wowsims/cata/sim/warrior/protection"
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
}
