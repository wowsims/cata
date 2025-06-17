package paladin

const (
	// Shared abilities
	SpellMaskAvengingWrath int64 = 1 << iota
	SpellMaskCensure
	SpellMaskCrusaderStrike
	SpellMaskDevotionAura
	SpellMaskDivineProtection
	SpellMaskDivineShield
	SpellMaskFlashOfLight
	SpellMaskGuardianOfAncientKings
	SpellMaskHammerOfWrath
	SpellMaskHammerOfTheRighteousMelee
	SpellMaskHammerOfTheRighteousAoe
	SpellMaskHandOfProtection
	SpellMaskJudgment
	SpellMaskLayOnHands
	SpellMaskSealOfInsight
	SpellMaskSealOfRighteousness
	SpellMaskSealOfTruth
	SpellMaskShieldOfTheRighteous
	SpellMaskTemplarsVerdict
	SpellMaskWordOfGlory

	// Retribution abilities
	SpellMaskDivineStorm
	SpellMaskExorcism
	SpellMaskInquisition
	SpellMaskSealOfJustice

	// Protection abilities
	SpellMaskArdentDefender
	SpellMaskAvengersShield
	SpellMaskConsecration
	SpellMaskHolyWrath

	// Holy abilities
	SpellMaskBeaconOfLight
	SpellMaskDaybreak
	SpellMaskDenounce
	SpellMaskDivineFavor
	SpellMaskDivineLight
	SpellMaskDivinePlea
	SpellMaskHolyLight
	SpellMaskHolyRadiance
	SpellMaskHolyShockDamage
	SpellMaskHolyShockHeal
	SpellMaskLightOfDawn

	// Glyphs
	SpellMaskHarshWords
)

const SpellMaskBuilderBase = SpellMaskCrusaderStrike |
	SpellMaskHammerOfTheRighteous

const SpellMaskBuilderRet = SpellMaskBuilderBase |
	SpellMaskJudgment |
	SpellMaskExorcism |
	SpellMaskHammerOfWrath

const SpellMaskBuilderProt = SpellMaskBuilderBase |
	SpellMaskJudgment |
	SpellMaskAvengersShield

const SpellMaskBuilderHoly = SpellMaskBuilderBase |
	// SpellMaskJudgment | only if Selfless Healer is talented
	SpellMaskHolyShock |
	SpellMaskHolyRadiance

const SpellMaskSpender = SpellMaskTemplarsVerdict |
	SpellMaskDivineStorm |
	SpellMaskHarshWords |
	SpellMaskInquisition |
	SpellMaskShieldOfTheRighteous |
	SpellMaskWordOfGlory

const SpellMaskSanctityOfBattleBase = SpellMaskCrusaderStrike |
	SpellMaskJudgment |
	SpellMaskHammerOfWrath |
	SpellMaskHarshWords |
	SpellMaskWordOfGlory

const SpellMaskSanctityOfBattleBaseGcd = SpellMaskCrusaderStrike |
	SpellMaskHammerOfWrath |
	SpellMaskJudgment |
	SpellMaskHarshWords |
	SpellMaskWordOfGlory

const SpellMaskSanctityOfBattleRet = SpellMaskSanctityOfBattleBase |
	// SpellMaskHammerOfTheRighteous | // Will be handled by Crusader Strike, since they share CD
	SpellMaskExorcism

const SpellMaskSanctityOfBattleRetGcd = SpellMaskSanctityOfBattleBaseGcd |
	// SpellMaskHammerOfTheRighteous | // Will be handled by Crusader Strike, since they share CD
	SpellMaskDivineStorm |
	SpellMaskTemplarsVerdict

const SpellMaskSanctityOfBattleProt = SpellMaskSanctityOfBattleBase |
	// SpellMaskHammerOfTheRighteous | // Will be handled by Crusader Strike, since they share CD
	SpellMaskAvengersShield |
	SpellMaskConsecration |
	SpellMaskHolyWrath |
	SpellMaskShieldOfTheRighteous

const SpellMaskSanctityOfBattleProtGcd = SpellMaskSanctityOfBattleBaseGcd

const SpellMaskSanctityOfBattleHoly = SpellMaskSanctityOfBattleBase |
	SpellMaskHolyShock

const SpellMaskSanctityOfBattleHolyGcd = SpellMaskSanctityOfBattleBaseGcd

const SpellMaskHolyShock = SpellMaskHolyShockDamage | SpellMaskHolyShockHeal

const SpellMaskHammerOfTheRighteous = SpellMaskHammerOfTheRighteousMelee | SpellMaskHammerOfTheRighteousAoe

const SpellMaskCanTriggerSealOfJustice = SpellMaskCrusaderStrike |
	SpellMaskHammerOfTheRighteousMelee |
	SpellMaskShieldOfTheRighteous |
	SpellMaskTemplarsVerdict

const SpellMaskCanTriggerSealOfInsight = SpellMaskCanTriggerSealOfJustice

const SpellMaskCanTriggerSealOfRighteousness = SpellMaskCrusaderStrike |
	SpellMaskTemplarsVerdict |
	SpellMaskDivineStorm |
	SpellMaskHammerOfTheRighteousMelee |
	SpellMaskShieldOfTheRighteous

const SpellMaskCanTriggerSealOfTruth = SpellMaskCrusaderStrike |
	SpellMaskTemplarsVerdict |
	SpellMaskJudgment |
	SpellMaskHammerOfTheRighteousMelee |
	SpellMaskShieldOfTheRighteous

const SpellMaskCanTriggerAncientPower = SpellMaskCanTriggerSealOfTruth

const SpellMaskCanTriggerHandOfLight = SpellMaskCrusaderStrike |
	SpellMaskDivineStorm |
	SpellMaskTemplarsVerdict |
	SpellMaskHammerOfTheRighteous |
	SpellMaskHammerOfWrath

const SpellMaskDamageModifiedBySwordOfLight = SpellMaskSealOfTruth |
	SpellMaskSealOfJustice |
	SpellMaskSealOfRighteousness |
	SpellMaskDivineStorm |
	SpellMaskHammerOfWrath |
	SpellMaskJudgment

const SpellMaskSeals = SpellMaskSealOfJustice |
	SpellMaskSealOfInsight |
	SpellMaskSealOfRighteousness |
	SpellMaskSealOfTruth

const SpellMaskModifiedBySealOfInsight = SpellMaskDivineLight |
	SpellMaskFlashOfLight |
	SpellMaskHolyLight |
	SpellMaskHolyRadiance |
	SpellMaskLayOnHands |
	SpellMaskLightOfDawn |
	SpellMaskWordOfGlory

const SpellMaskCausesForbearance = SpellMaskDivineShield |
	SpellMaskHandOfProtection |
	SpellMaskLayOnHands
