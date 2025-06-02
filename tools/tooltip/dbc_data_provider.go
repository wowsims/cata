package tooltip

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/tools/database/dbc"
)

type EffectMap map[int]dbc.SpellEffect
type DBCTooltipDataProvider struct {
	DBC *dbc.DBC
}

func GetEffectByIndex(effects map[int]dbc.SpellEffect, index int) *dbc.SpellEffect {
	if len(effects) <= index {
		return nil
	}

	// quick check
	effect := effects[index]
	if effect.EffectIndex == index {
		return &effect
	}

	// did not find
	for _, e := range effects {
		if e.EffectIndex == index {
			return &e
		}
	}

	return nil
}

// GetSpellPPM implements TooltipDataProvider.
func (d DBCTooltipDataProvider) GetSpellPPM(spellId int64) float64 {
	spell, ok := d.DBC.Spells[int(spellId)]
	if !ok {
		return 1
	}

	return float64(spell.SpellProcsPerMinute)
}

// GetSpellProcCooldown implements TooltipDataProvider.
func (d DBCTooltipDataProvider) GetSpellProcCooldown(spellId int64) time.Duration {
	spell, ok := d.DBC.Spells[int(spellId)]
	if !ok {
		return 1
	}

	return time.Duration(spell.ProcCategoryRecovery) * time.Millisecond
}

func (d DBCTooltipDataProvider) GetSpellMaxTargets(spellId int64) int64 {
	spell, ok := d.DBC.Spells[int(spellId)]
	if !ok {
		return 1
	}

	return int64(spell.MaxTargets)
}

func (d DBCTooltipDataProvider) GetEffectAmplitude(spellId int64, effectIdx int64) float64 {
	effectEntries, ok := d.DBC.SpellEffects[int(spellId)]
	if !ok {
		return 1
	}

	effect := GetEffectByIndex(effectEntries, int(effectIdx))
	if effect == nil {
		return 0
	}

	return effect.EffectAmplitude
}

func (d DBCTooltipDataProvider) GetEffectChainAmplitude(spellId int64, effectIdx int64) float64 {
	effectEntries, ok := d.DBC.SpellEffects[int(spellId)]
	if !ok {
		return 1
	}

	effect := GetEffectByIndex(effectEntries, int(effectIdx))
	if effect == nil {
		return 0
	}

	return effect.EffectChainAmplitude
}

func (d DBCTooltipDataProvider) GetEffectPointsPerResource(spellId int64, effectIdx int64) float64 {
	effectEntries, ok := d.DBC.SpellEffects[int(spellId)]
	if !ok {
		return 1
	}

	effect := GetEffectByIndex(effectEntries, int(effectIdx))
	if effect == nil {
		return 0
	}

	return effect.EffectPointsPerResource
}

// GetEffectMaxTargets implements TooltipDataProvider.
func (d DBCTooltipDataProvider) GetEffectMaxTargets(spellId int64, effectIdx int64) int64 {
	effectEntries, ok := d.DBC.SpellEffects[int(spellId)]
	if !ok {
		return 1
	}

	effect := GetEffectByIndex(effectEntries, int(effectIdx))
	if effect == nil {
		return 0
	}

	return int64(effect.EffectChainTargets)
}

// GetSpellProcChance implements TooltipDataProvider.
func (d DBCTooltipDataProvider) GetSpellProcChance(spellId int64) float64 {
	spellEntry, ok := d.DBC.Spells[int(spellId)]
	if !ok {
		return 0
	}

	return float64(spellEntry.ProcChance)
}

func (d DBCTooltipDataProvider) GetSpecNum() int64 {
	return 0
}

func (d DBCTooltipDataProvider) GetSpellIcon(spellId int64) string {
	spellEntry, ok := d.DBC.Spells[int(spellId)]
	if !ok {
		return ""
	}

	return spellEntry.IconPath
}

// GetMainHandWeapon implements TooltipDataProvider.
func (d DBCTooltipDataProvider) GetMainHandWeapon() *core.Weapon {
	// Item: 103727 as dummy for now
	return &core.Weapon{
		BaseDamageMin: 10257,
		BaseDamageMax: 19050,
		SwingSpeed:    2.6,
	}
}

// GetOffHandWeapon implements TooltipDataProvider.
func (d DBCTooltipDataProvider) GetOffHandWeapon() *core.Weapon {
	return &core.Weapon{
		BaseDamageMin: 10257,
		BaseDamageMax: 19050,
		SwingSpeed:    2.6,
	}
}

func (d DBCTooltipDataProvider) GetPlayerLevel() float64 {
	return core.CharacterLevel
}

func (d DBCTooltipDataProvider) GetSpellDescription(spellId int64) string {
	spellEntry, ok := d.DBC.Spells[int(spellId)]
	if !ok {
		return ""
	}

	return spellEntry.Description
}

func (d DBCTooltipDataProvider) GetSpellName(spellId int64) string {
	spellEntry, ok := d.DBC.Spells[int(spellId)]
	if !ok {
		return ""
	}

	return spellEntry.NameLang
}

// GetAttackPower implements TooltipDataProvider.
func (d DBCTooltipDataProvider) GetAttackPower() float64 {
	return 1
}

func (d DBCTooltipDataProvider) ShouldUseBaseScaling(spellId int64) bool {
	spellEntry, ok := d.DBC.Spells[int(spellId)]
	if !ok {
		return false
	}

	// need proper scaling entry
	return spellEntry.SpellClassSet > 0
}
func (d DBCTooltipDataProvider) GetClass(spellId int64) proto.Class {
	spellEntry, ok := d.DBC.Spells[int(spellId)]
	if !ok {
		return proto.Class_ClassUnknown
	}

	switch spellEntry.SpellClassSet {
	case 53:
		return proto.Class_ClassMonk
	case 15:
		return proto.Class_ClassDeathKnight
	case 11:
		return proto.Class_ClassShaman
	case 10:
		return proto.Class_ClassPaladin
	case 9:
		return proto.Class_ClassHunter
	case 8:
		return proto.Class_ClassRogue
	case 7:
		return proto.Class_ClassDruid
	case 6:
		return proto.Class_ClassPriest
	case 5:
		return proto.Class_ClassWarlock
	case 4:
		return proto.Class_ClassWarrior
	case 3:
		return proto.Class_ClassMage
	default:
		return proto.Class_ClassUnknown
	}
}

// GetEffectBaseDamage implements TooltipDataProvider.
func (d DBCTooltipDataProvider) GetEffectScaledValue(spellId int64, effectIdx int64) float64 {
	effectEntries, ok := d.DBC.SpellEffects[int(spellId)]
	class := d.GetClass(spellId)

	if !ok {
		return 1
	}

	// some spells are just fucked..
	if int(effectIdx) >= len(effectEntries) {
		effectIdx = int64(len(effectEntries) - 1)
	}

	effect := GetEffectByIndex(effectEntries, int(effectIdx))
	if effect == nil {
		return 1
	}

	baseDamage := 0.0

	// using class scaling
	if effect.Coefficient > 0 && d.ShouldUseBaseScaling(spellId) {
		baseValue := 0.0

		// for now use generic unk13 scaling for level 90
		if class == proto.Class_ClassUnknown {
			baseValue = 1710.000000
		} else {
			baseValue = core.ClassBaseScaling[class]
		}

		baseDamage += baseValue * effect.Coefficient
	} else {
		baseDamage += float64(effect.EffectBasePoints)
		spell := d.DBC.Spells[int(spellId)]
		if spell.MaxScalingLevel > 0 {
			baseDamage += effect.EffectRealPointsPerLevel * math.Min(float64(spell.MaxScalingLevel), core.CharacterLevel)
		}
	}

	shouldScale := false
	switch effect.EffectType {
	case dbc.E_SCHOOL_DAMAGE:
		shouldScale = true

	case dbc.E_APPLY_AURA:
		fallthrough
	case dbc.E_APPLY_AREA_AURA_ENEMY:
		fallthrough
	case dbc.E_APPLY_AREA_AURA_FRIEND:
		fallthrough
	case dbc.E_APPLY_AREA_AURA_PARTY:
		fallthrough
	case dbc.E_APPLY_AREA_AURA_OWNER:
		fallthrough
	case dbc.E_APPLY_AREA_AURA_RAID:
		fallthrough
	case dbc.E_APPLY_AREA_AURA_PARTY_NONRANDOM:
		fallthrough
	case dbc.E_APPLY_AREA_AURA_PET:
		fallthrough
	case dbc.E_APPLY_AURA_ON_PET:
		switch effect.EffectAura {
		case dbc.A_PERIODIC_DAMAGE:
			fallthrough
		case dbc.A_PERIODIC_HEAL:
			shouldScale = true
		}
	}

	if !shouldScale {
		return baseDamage
	}

	if effect.BonusCoefficientFromAP > 0 {
		baseDamage += d.GetAttackPower() * effect.BonusCoefficientFromAP
	}

	if effect.EffectBonusCoefficient > 0 {
		baseDamage += d.GetSpellPower() * effect.EffectBonusCoefficient
	}

	return baseDamage
}

// GetDescriptionVariableString implements TooltipDataProvider.
func (d DBCTooltipDataProvider) GetDescriptionVariableString(spellId int64) string {
	spellEntry, ok := d.DBC.Spells[int(spellId)]
	if !ok {
		return ""
	}

	return spellEntry.Variables
}

// GetEffectBaseValue implements TooltipDataProvider.
func (d DBCTooltipDataProvider) GetEffectBaseValue(spellId int64, effectIdx int64) float64 {
	effectEntries, ok := d.DBC.SpellEffects[int(spellId)]
	if !ok {
		return 0
	}

	effect := GetEffectByIndex(effectEntries, int(effectIdx))
	if effect == nil {
		return 0
	}

	return float64(effect.EffectBasePoints)
}

// GetEffectPeriod implements TooltipDataProvider.
func (d DBCTooltipDataProvider) GetEffectPeriod(spellId int64, effectIdx int64) time.Duration {
	effectEntries, ok := d.DBC.SpellEffects[int(spellId)]
	if !ok {
		return 0
	}

	effect := GetEffectByIndex(effectEntries, int(effectIdx))
	if effect == nil {
		return 0
	}

	return time.Duration(effect.EffectAuraPeriod) * time.Millisecond
}

// GetEffectRadius implements TooltipDataProvider.
func (d DBCTooltipDataProvider) GetEffectRadius(spellId int64, effectIdx int64) float64 {
	effectEntries, ok := d.DBC.SpellEffects[int(spellId)]
	if !ok {
		return 0
	}

	effect := GetEffectByIndex(effectEntries, int(effectIdx))
	if effect == nil {
		return 0
	}

	return effect.GetRadiusMax()
}

// GetSpellDuration implements TooltipDataProvider.
func (d DBCTooltipDataProvider) GetSpellDuration(spellId int64) time.Duration {
	spell, ok := d.DBC.Spells[int(spellId)]
	if !ok {
		return 0
	}

	if spell.Duration < 0 {
		return 0
	}

	return time.Duration(spell.Duration) * time.Millisecond
}

// GetSpellPower implements TooltipDataProvider.
func (d DBCTooltipDataProvider) GetSpellPower() float64 {
	return 15000
}

// GetSpellRange implements TooltipDataProvider.
func (d DBCTooltipDataProvider) GetSpellRange(spellId int64) float64 {
	spell, ok := d.DBC.Spells[int(spellId)]
	if !ok {
		return 0
	}

	return float64(spell.MaxRange)
}

// GetStacks implements TooltipDataProvider.
func (d DBCTooltipDataProvider) GetSpellStacks(spellId int64) int64 {
	spell, ok := d.DBC.Spells[int(spellId)]
	if !ok {
		return 0
	}

	if spell.ProcCharges > 0 {
		return int64(spell.ProcCharges)
	}

	if spell.MaxCumulativeStacks > 0 {
		return int64(spell.MaxCumulativeStacks)
	}

	return 0
}

// HasAura implements TooltipDataProvider.
func (d DBCTooltipDataProvider) HasAura(auraId int64) bool {
	return true
}

// HasPassive implements TooltipDataProvider.
func (d DBCTooltipDataProvider) HasPassive(auraId int64) bool {
	return true
}

// IsMaleGender implements TooltipDataProvider.
func (d DBCTooltipDataProvider) IsMaleGender() bool {
	return true
}

// KnowsSpell implements TooltipDataProvider.
func (d DBCTooltipDataProvider) KnowsSpell(spellId int64) bool {
	return true
}

func (d DBCTooltipDataProvider) GetEffectEnchantValue(enchantId int64, effectIdx int64) float64 {
	enchantInfo, ok := d.DBC.Enchants[int(enchantId)]
	if !ok {
		return 0
	}

	return float64(enchantInfo.EffectPoints[effectIdx])
}
