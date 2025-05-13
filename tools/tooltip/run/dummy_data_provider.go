package main

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

type DummyTooltipDataProvider struct{}

// GetEffectAmplitude implements tooltip.TooltipDataProvider.
func (d DummyTooltipDataProvider) GetEffectAmplitude(spellId int64, effectIdx int64) float64 {
	return 2000
}

// GetEffectChainAmplitude implements tooltip.TooltipDataProvider.
func (d DummyTooltipDataProvider) GetEffectChainAmplitude(spellId int64, effectidx int64) float64 {
	return 3
}

// GetEffectMaxTargets implements tooltip.TooltipDataProvider.
func (d DummyTooltipDataProvider) GetEffectMaxTargets(spellId int64, effectIdx int64) int64 {
	return 3
}

// GetEffectPointsPerResource implements tooltip.TooltipDataProvider.
func (d DummyTooltipDataProvider) GetEffectPointsPerResource(spellId int64, effectIdx int64) float64 {
	return 20
}

// GetSpellMaxTargets implements tooltip.TooltipDataProvider.
func (d DummyTooltipDataProvider) GetSpellMaxTargets(spellId int64) int64 {
	return 5
}

// GetSpellPPM implements tooltip.TooltipDataProvider.
func (d DummyTooltipDataProvider) GetSpellPPM(spellId int64) float64 {
	return 0.85
}

// GetSpellProcCooldown implements tooltip.TooltipDataProvider.
func (d DummyTooltipDataProvider) GetSpellProcCooldown(spellId int64) time.Duration {
	return time.Second * 45
}

// GetSpellStacks implements tooltip.TooltipDataProvider.
func (d DummyTooltipDataProvider) GetSpellStacks(spellId int64) int64 {
	return 20
}

func (d DummyTooltipDataProvider) GetPlayerLevel() float64 {
	return 90
}

func (d DummyTooltipDataProvider) GetSpellProcChance(spellId int64) float64 {
	return 0.5
}

func (d DummyTooltipDataProvider) GetEffectBaseValue(spellId int64, effectIndex int64) float64 {
	return 100
}

// GetEffectRadius implements TooltipDataProvider.
func (d DummyTooltipDataProvider) GetEffectRadius(spellid int64, effectIdx int64) float64 {
	return 0
}

// GetStacks implements TooltipDataProvider.
func (d DummyTooltipDataProvider) GetStacks(spellId int64) int64 {
	return 5
}

// GetAttackPower implements TooltipDataProvider.
func (d DummyTooltipDataProvider) GetAttackPower() float64 {
	return 1
}

// GetEffectScaledValue implements TooltipDataProvider.
func (d DummyTooltipDataProvider) GetEffectScaledValue(spellId int64, effectIdx int64) float64 {
	// Dummy level 90 priest
	return 1045.69
}

func (d DummyTooltipDataProvider) GetSpellIcon(spellId int64) string {
	return "interface\\icons\\ability_deathknight_brittlebones.blp"
}

// GetDescriptionVariableString implements TooltipDataProvider.
func (d DummyTooltipDataProvider) GetDescriptionVariableString(spellId int64) string {
	return "$stnc=$?a103985[${1.2}][${1.0}] $dwm1=$?a108561[${1}][${0.898882275}] $dwm=$?a115697[${1}][${$<dwm1>}] $bm=$?s120267[${1}][${1}] $offm1=$?a108561[${0}][${1}] $offm=$?a115697[${0}][${$<offm1>}] $apc=$?s120267[${$AP/14}][${$AP/14}] $offlow=$?!s124146[${$mwb/2/$mws}][${$owb/2/$ows}] $offhigh=$?!s124146[${$MWB/2/$mws}][${$OWB/2/$ows}] $low=${$<stnc>*($<bm>*$<dwm>*(($mwb)/($MWS)+$<offm>*$<offlow>)+$<apc>-1)} $test=${$<stnc>*($<bm>*$<dwm>*(($MWB)/($MWS)+$<offm>*$<offhigh>)+$<apc>+1)}"
}

// GetSpellDuration implements TooltipDataProvider.
func (d DummyTooltipDataProvider) GetSpellDuration(spellId int64) time.Duration {
	return time.Second * 12
}

// GetSpellPower implements TooltipDataProvider.
func (d DummyTooltipDataProvider) GetSpellPower() float64 {
	return 1
}

// GetSpellRange implements TooltipDataProvider.
func (d DummyTooltipDataProvider) GetSpellRange(spellId int64) float64 {
	return 30
}

// HasAura implements TooltipDataProvider.
func (d DummyTooltipDataProvider) HasAura(auraId int64) bool {
	return true
}

// HasPassive implements TooltipDataProvider.
func (d DummyTooltipDataProvider) HasPassive(auraId int64) bool {
	return true
}

// IsMaleGender implements TooltipDataProvider.
func (d DummyTooltipDataProvider) IsMaleGender() bool {
	return true
}

// KnowsSpell implements TooltipDataProvider.
func (d DummyTooltipDataProvider) KnowsSpell(spellId int64) bool {
	return true
}

func (d DummyTooltipDataProvider) GetSpellDescription(spellId int64) string {
	return "This dummy spell has a chance to allow ${15-$max($PL-70,0)/2}% of your mana regeneration to continue while casting for $38346d."
}

func (d DummyTooltipDataProvider) GetSpellName(spellId int64) string {
	return "DummyRefSpell"
}

func (d DummyTooltipDataProvider) GetEffectPeriod(spellId int64, effectIdx int64) time.Duration {
	return time.Second * 2
}

func (d DummyTooltipDataProvider) GetSpecNum() int64 {
	return 1
}

func (d DummyTooltipDataProvider) GetMainHandWeapon() *core.Weapon {
	return &core.Weapon{
		BaseDamageMin: 100,
		BaseDamageMax: 200,
		SwingSpeed:    2.6,
	}
}

func (d DummyTooltipDataProvider) GetOffHandWeapon() *core.Weapon {
	return &core.Weapon{
		BaseDamageMin: 100,
		BaseDamageMax: 200,
		SwingSpeed:    2.6,
	}
}

func (d DummyTooltipDataProvider) GetEffectEnchantValue(enchantId int64, effectIdx int64) float64 {
	return 10
}
