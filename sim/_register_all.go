package sim

import (
	_ "github.com/wowsims/cata/sim/common"
	dpsDeathKnight "github.com/wowsims/cata/sim/deathknight/dps"
	tankDeathKnight "github.com/wowsims/cata/sim/deathknight/tank"
	"github.com/wowsims/cata/sim/druid/balance"
	"github.com/wowsims/cata/sim/druid/feral"
	restoDruid "github.com/wowsims/cata/sim/druid/restoration"
	feralTank "github.com/wowsims/cata/sim/druid/tank"
	_ "github.com/wowsims/cata/sim/encounters"
	"github.com/wowsims/cata/sim/hunter"
	"github.com/wowsims/cata/sim/mage"
	holyPaladin "github.com/wowsims/cata/sim/paladin/holy"
	protectionPaladin "github.com/wowsims/cata/sim/paladin/protection"
	"github.com/wowsims/cata/sim/paladin/retribution"
	healingPriest "github.com/wowsims/cata/sim/priest/healing"
	"github.com/wowsims/cata/sim/priest/shadow"
	"github.com/wowsims/cata/sim/priest/smite"
	"github.com/wowsims/cata/sim/rogue"
	"github.com/wowsims/cata/sim/shaman/elemental"
	"github.com/wowsims/cata/sim/shaman/enhancement"
	restoShaman "github.com/wowsims/cata/sim/shaman/restoration"
	"github.com/wowsims/cata/sim/warlock"
	dpsWarrior "github.com/wowsims/cata/sim/warrior/dps"
	protectionWarrior "github.com/wowsims/cata/sim/warrior/protection"
)

var registered = false

func RegisterAll() {
	if registered {
		return
	}
	registered = true

	balance.RegisterBalanceDruid()
	feral.RegisterFeralDruid()
	feralTank.RegisterFeralTankDruid()
	restoDruid.RegisterRestorationDruid()
	elemental.RegisterElementalShaman()
	enhancement.RegisterEnhancementShaman()
	restoShaman.RegisterRestorationShaman()
	hunter.RegisterHunter()
	mage.RegisterMage()
	healingPriest.RegisterHealingPriest()
	shadow.RegisterShadowPriest()
	smite.RegisterSmitePriest()
	rogue.RegisterRogue()
	dpsWarrior.RegisterDpsWarrior()
	protectionWarrior.RegisterProtectionWarrior()
	holyPaladin.RegisterHolyPaladin()
	protectionPaladin.RegisterProtectionPaladin()
	retribution.RegisterRetributionPaladin()
	warlock.RegisterWarlock()
	dpsDeathKnight.RegisterDpsDeathknight()
	tankDeathKnight.RegisterTankDeathknight()
}
