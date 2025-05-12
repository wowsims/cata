package database

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

type WowheadTooltipManager struct {
	TooltipManager
}

func (wtm *WowheadTooltipManager) Read() map[int32]WowheadItemResponse {
	strDB := wtm.TooltipManager.Read()
	return core.MapMap(strDB, func(id int32, tooltip string) (int32, WowheadItemResponse) {
		return id, NewWowheadItemResponse(id, tooltip)
	})
}

func NewWowheadItemTooltipManager(filePath string) *WowheadTooltipManager {
	return &WowheadTooltipManager{
		TooltipManager{
			FilePath:   filePath,
			UrlPattern: "https://nether.wowhead.com/mop-classic/tooltip/item/%s?lvl=90",
		},
	}
}

func NewWowheadSpellTooltipManager(filePath string) *WowheadTooltipManager {
	return &WowheadTooltipManager{
		TooltipManager{
			FilePath:   filePath,
			UrlPattern: "https://nether.wowhead.com/mop-classic/tooltip/spell/%s",
		},
	}
}

type ItemResponse interface {
	GetName() string
	GetQuality() int
	GetIcon() string
	HasBuff() bool
	TooltipWithoutSetBonus() string
	GetTooltipRegexString(pattern *regexp.Regexp, matchIdx int) string
	GetTooltipRegexValue(pattern *regexp.Regexp, matchIdx int) int
	GetIntValue(pattern *regexp.Regexp) int
	GetStats() stats.Stats
	IsEquippable() bool
	GetItemLevel() int
	GetPhase() int
	GetUnique() bool
	GetItemType() proto.ItemType
	GetArmorType() proto.ArmorType
	GetWeaponType() proto.WeaponType
	GetHandType() proto.HandType
	GetRangedWeaponType() proto.RangedWeaponType
	GetWeaponDamage() (float64, float64)
	GetWeaponSpeed() float64
	GetGemSockets() []proto.GemColor
	GetSocketBonus() stats.Stats
	GetSocketColor() proto.GemColor
	GetGemStats() stats.Stats
	GetItemSetName() string
	IsHeroic() bool
	GetRequiredProfession() proto.Profession
}

type WowheadItemResponse struct {
	ID      int32
	Name    string `json:"name"`
	Quality int    `json:"quality"`
	Icon    string `json:"icon"`
	Tooltip string `json:"tooltip"`
	Buff    string `json:"buff"`
}

func NewWowheadItemResponse(id int32, tooltip string) WowheadItemResponse {
	response := WowheadItemResponse{}
	err := json.Unmarshal([]byte(tooltip), &response)
	if err != nil {
		fmt.Printf("Failed to decode tooltipBytes: %s\n", tooltip)
		panic(err)
	}
	response.ID = id
	return response
}

func (item WowheadItemResponse) GetName() string {
	return item.Name
}
func (item WowheadItemResponse) GetQuality() int {
	return item.Quality
}
func (item WowheadItemResponse) GetIcon() string {
	return item.Icon
}
func (item WowheadItemResponse) HasBuff() bool {
	return item.Buff != ""
}

func GetRegexStringValue(srcStr string, pattern *regexp.Regexp, matchIdx int) string {
	match := pattern.FindStringSubmatch(srcStr)
	if match == nil {
		return ""
	} else {
		return match[matchIdx]
	}
}
func GetRegexIntValue(srcStr string, pattern *regexp.Regexp, matchIdx int) int {
	matchStr := GetRegexStringValue(srcStr, pattern, matchIdx)
	matchStr = strings.Replace(matchStr, ",", "", -1)

	val, err := strconv.Atoi(matchStr)
	if err != nil {
		return 0
	}

	return val
}
func GetBestRegexIntValue(srcStr string, patterns []*regexp.Regexp, matchIdx int) int {
	best := 0
	for _, pattern := range patterns {
		newVal := GetRegexIntValue(srcStr, pattern, matchIdx)
		if newVal > best {
			best = newVal
		}
	}
	return best
}

func (item WowheadItemResponse) TooltipWithoutSetBonus() string {
	setIdx := strings.Index(item.Tooltip, "Set : ")
	if setIdx == -1 {
		return item.Tooltip
	} else {
		return item.Tooltip[:setIdx]
	}
}

func (item WowheadItemResponse) GetTooltipRegexString(pattern *regexp.Regexp, matchIdx int) string {
	return GetRegexStringValue(item.TooltipWithoutSetBonus(), pattern, matchIdx)
}

func (item WowheadItemResponse) GetTooltipRegexValue(pattern *regexp.Regexp, matchIdx int) int {
	return GetRegexIntValue(item.TooltipWithoutSetBonus(), pattern, matchIdx)
}

func (item WowheadItemResponse) GetIntValue(pattern *regexp.Regexp) int {
	return item.GetTooltipRegexValue(pattern, 1)
}

var expansionRegex = regexp.MustCompile(`(tbc|wotlk|cata|mop)`)

var armorRegex = regexp.MustCompile(`<!--amr-->([0-9]+) Armor`)
var agilityRegex = regexp.MustCompile(`<!--stat3-->\+([0-9]+) Agility`)
var strengthRegex = regexp.MustCompile(`<!--stat4-->\+([0-9]+) Strength`)
var intellectRegex = regexp.MustCompile(`<!--stat5-->\+([0-9]+) Intellect`)
var spiritRegex = regexp.MustCompile(`<!--stat6-->\+([0-9]+) Spirit`)
var staminaRegex = regexp.MustCompile(`<!--stat7-->\+([0-9]+) Stamina`)
var spellPowerRegex = regexp.MustCompile(`Increases spell power by ([0-9]{1,3}(,[0-9]{3})*)\.`)
var spellPowerRegex2 = regexp.MustCompile(`Increases spell power by <!--rtg45-->([0-9]{1,3}(,[0-9]{3})*)\.`)
var masteryRegex = regexp.MustCompile(`<!--rtg49-->([0-9]+)\s*Mastery`)

/*
// Not sure these exist anymore?
var arcaneSpellPowerRegex = regexp.MustCompile(`Increases Arcane power by ([0-9]+)\.`)
var fireSpellPowerRegex = regexp.MustCompile(`Increases Fire power by ([0-9]+)\.`)
var frostSpellPowerRegex = regexp.MustCompile(`Increases Frost power by ([0-9]+)\.`)
var holySpellPowerRegex = regexp.MustCompile(`Increases Holy power by ([0-9]+)\.`)
var natureSpellPowerRegex = regexp.MustCompile(`Increases Nature power by ([0-9]+)\.`)
var shadowSpellPowerRegex = regexp.MustCompile(`Increases Shadow power by ([0-9]+)\.`)
*/

var hitRegex = regexp.MustCompile(`Improves hit rating by <!--rtg31-->([0-9]+)\.`)
var critRegex = regexp.MustCompile(`Improves critical strike rating by <!--rtg32-->([0-9]+)\.`)
var hasteRegex = regexp.MustCompile(`Improves haste rating by <!--rtg36-->([0-9]+)\.`)

var spellPenetrationRegex = regexp.MustCompile(`Increases your spell penetration by ([0-9]+)\.`)
var mp5Regex = regexp.MustCompile(`Restores ([0-9]+) mana per 5 sec\.`)
var attackPowerRegex = regexp.MustCompile(`Increases attack power by ([0-9]+)\.`)
var attackPowerRegex2 = regexp.MustCompile(`Increases attack power by <!--rtg38-->([0-9]+)\.`)

var rangedAttackPowerRegex = regexp.MustCompile(`Increases ranged attack power by ([0-9]+)\.`)
var rangedAttackPowerRegex2 = regexp.MustCompile(`Increases ranged attack power by <!--rtg39-->([0-9]+)\.`)

var armorPenetrationRegex = regexp.MustCompile(`Increases armor penetration rating by ([0-9]+)`)
var armorPenetrationRegex2 = regexp.MustCompile(`Increases your armor penetration by <!--rtg44-->([0-9]+)\.`)

var expertiseRegex = regexp.MustCompile(`Increases your expertise rating by <!--rtg37-->([0-9]+)\.`)
var weaponDamageRegex = regexp.MustCompile(`<!--dmg-->([0-9]+) - ([0-9]+)`)
var weaponDamageRegex2 = regexp.MustCompile(`<!--dmg-->([0-9]+) Damage`)
var weaponSpeedRegex = regexp.MustCompile(`<!--spd-->(([0-9]+).([0-9]+))`)

var defenseRegex = regexp.MustCompile(`Increases defense rating by <!--rtg12-->([0-9]+)\.`)
var defenseRegex2 = regexp.MustCompile(`Increases defense rating by ([0-9]+)\.`)
var blockRegex = regexp.MustCompile(`Increases your shield block rating by <!--rtg15-->([0-9]+)\.`)
var blockRegex2 = regexp.MustCompile(`Increases your shield block rating by ([0-9]+)\.`)
var dodgeRegex = regexp.MustCompile(`Increases your dodge rating by <!--rtg13-->([0-9]+)\.`)
var dodgeRegex2 = regexp.MustCompile(`Increases your dodge rating by ([0-9]+)\.`)
var parryRegex = regexp.MustCompile(`Increases your parry rating by <!--rtg14-->([0-9]+)\.`)
var parryRegex2 = regexp.MustCompile(`Increases your parry rating by ([0-9]+)\.`)
var resilienceRegex = regexp.MustCompile(`Improves your resilience rating by <!--rtg35-->([0-9]+)\.`)
var arcaneResistanceRegex = regexp.MustCompile(`\+([0-9]+) Arcane Resistance`)
var fireResistanceRegex = regexp.MustCompile(`\+([0-9]+) Fire Resistance`)
var frostResistanceRegex = regexp.MustCompile(`\+([0-9]+) Frost Resistance`)
var natureResistanceRegex = regexp.MustCompile(`\+([0-9]+) Nature Resistance`)
var shadowResistanceRegex = regexp.MustCompile(`\+([0-9]+) Shadow Resistance`)
var bonusArmorRegex = regexp.MustCompile(`Has ([0-9]+) bonus armor`)
var bonusArmorRegex2 = regexp.MustCompile(`([\d,\.]+) Bonus Armor`)

func (item WowheadItemResponse) GetStats() stats.Stats {
	sp := float64(item.GetIntValue(spellPowerRegex)) + float64(item.GetIntValue(spellPowerRegex2))
	baseAP := float64(item.GetIntValue(attackPowerRegex)) + float64(item.GetIntValue(attackPowerRegex2))
	armor, bonusArmor := item.GetArmorValues()
	return stats.Stats{
		stats.Armor:             float64(armor),
		stats.BonusArmor:        float64(bonusArmor),
		stats.Strength:          float64(item.GetIntValue(strengthRegex)),
		stats.Agility:           float64(item.GetIntValue(agilityRegex)),
		stats.Stamina:           float64(item.GetIntValue(staminaRegex)),
		stats.Intellect:         float64(item.GetIntValue(intellectRegex)),
		stats.Spirit:            float64(item.GetIntValue(spiritRegex)),
		stats.SpellPower:        sp,
		stats.HitRating:         float64(item.GetIntValue(hitRegex)),
		stats.CritRating:        float64(item.GetIntValue(critRegex)),
		stats.HasteRating:       float64(item.GetIntValue(hasteRegex)),
		stats.SpellPenetration:  float64(item.GetIntValue(spellPenetrationRegex)),
		stats.MP5:               float64(item.GetIntValue(mp5Regex)),
		stats.AttackPower:       baseAP,
		stats.RangedAttackPower: baseAP + float64(item.GetIntValue(rangedAttackPowerRegex)) + float64(item.GetIntValue(rangedAttackPowerRegex2)),
		stats.ExpertiseRating:   float64(item.GetIntValue(expertiseRegex)),
		stats.DodgeRating:       float64(item.GetIntValue(dodgeRegex) + item.GetIntValue(dodgeRegex2)),
		stats.ParryRating:       float64(item.GetIntValue(parryRegex) + item.GetIntValue(parryRegex2)),
		stats.ResilienceRating:  float64(item.GetIntValue(resilienceRegex)),
		stats.ArcaneResistance:  float64(item.GetIntValue(arcaneResistanceRegex)),
		stats.FireResistance:    float64(item.GetIntValue(fireResistanceRegex)),
		stats.FrostResistance:   float64(item.GetIntValue(frostResistanceRegex)),
		stats.NatureResistance:  float64(item.GetIntValue(natureResistanceRegex)),
		stats.ShadowResistance:  float64(item.GetIntValue(shadowResistanceRegex)),
		stats.MasteryRating:     float64(item.GetIntValue(masteryRegex)),
	}
}

var patternRegexes = []*regexp.Regexp{
	regexp.MustCompile(`Design:`),
	regexp.MustCompile(`Recipe:`),
	regexp.MustCompile(`Pattern:`),
	regexp.MustCompile(`Plans:`),
	regexp.MustCompile(`Schematic:`),
}

func (item WowheadItemResponse) IsPattern() bool {
	for _, pattern := range patternRegexes {
		if pattern.MatchString(item.Tooltip) {
			return true
		}
	}
	return false
}

var randomEnchantRegex = regexp.MustCompile(`Random enchantment`)

func (item WowheadItemResponse) IsRandomEnchant() bool {
	return randomEnchantRegex.MatchString(item.Tooltip)
}

func (item WowheadItemResponse) IsEquippable() bool {
	return item.GetItemType() != proto.ItemType_ItemTypeUnknown &&
		!item.IsPattern() &&
		item.GetItemLevel() <= 416
}

var itemLevelRegex = regexp.MustCompile(`Item Level <!--ilvl-->([0-9]+)<`)

func (item WowheadItemResponse) GetItemLevel() int {
	return item.GetIntValue(itemLevelRegex)
}

var phaseRegex = regexp.MustCompile(`Phase ([0-9])`)

func (item WowheadItemResponse) GetPhase() int {
	phase := item.GetIntValue(phaseRegex)
	if phase != 0 {
		return phase
	}

	ilvl := item.GetItemLevel()
	if ilvl <= 284 { // TBC items
		return 0
	} else {
		return 1
	}
}

var uniqueRegex = regexp.MustCompile(`Unique`)
var jcGemsRegex = regexp.MustCompile(`Jeweler's Gems`)

func (item WowheadItemResponse) GetUnique() bool {
	return uniqueRegex.MatchString(item.Tooltip) && !jcGemsRegex.MatchString(item.Tooltip)
}

var itemTypePatterns = map[proto.ItemType]*regexp.Regexp{
	proto.ItemType_ItemTypeHead:     regexp.MustCompile(`<td>Head</td>`),
	proto.ItemType_ItemTypeNeck:     regexp.MustCompile(`<td>Neck</td>`),
	proto.ItemType_ItemTypeShoulder: regexp.MustCompile(`<td>Shoulder</td>`),
	proto.ItemType_ItemTypeBack:     regexp.MustCompile(`<td>Back</td>`),
	proto.ItemType_ItemTypeChest:    regexp.MustCompile(`<td>Chest</td>`),
	proto.ItemType_ItemTypeWrist:    regexp.MustCompile(`<td>Wrist</td>`),
	proto.ItemType_ItemTypeHands:    regexp.MustCompile(`<td>Hands</td>`),
	proto.ItemType_ItemTypeWaist:    regexp.MustCompile(`<td>Waist</td>`),
	proto.ItemType_ItemTypeLegs:     regexp.MustCompile(`<td>Legs</td>`),
	proto.ItemType_ItemTypeFeet:     regexp.MustCompile(`<td>Feet</td>`),
	proto.ItemType_ItemTypeFinger:   regexp.MustCompile(`<td>Finger</td>`),
	proto.ItemType_ItemTypeTrinket:  regexp.MustCompile(`<td>Trinket</td>`),
	proto.ItemType_ItemTypeWeapon:   regexp.MustCompile(`<td>((Main Hand)|(Two-Hand)|(One-Hand)|(Off Hand)|(Held In Off-hand)|(Held In Off-Hand))</td>`),
	proto.ItemType_ItemTypeRanged:   regexp.MustCompile(`<td>(Ranged|Thrown|Relic)</td>`),
}

func (item WowheadItemResponse) GetItemType() proto.ItemType {
	for itemType, pattern := range itemTypePatterns {
		if pattern.MatchString(item.Tooltip) {
			return itemType
		}
	}
	return proto.ItemType_ItemTypeUnknown
}

func (item WowheadItemResponse) IsScalableArmorSlot() bool {
	// Special case shields as Base Armor
	if item.GetWeaponType() == proto.WeaponType_WeaponTypeShield {
		return true
	}

	itemType := item.GetItemType()
	switch itemType {
	case
		proto.ItemType_ItemTypeNeck,
		proto.ItemType_ItemTypeFinger,
		proto.ItemType_ItemTypeTrinket,
		proto.ItemType_ItemTypeWeapon:
		return false
	}
	return true
}

func (item WowheadItemResponse) GetArmorValues() (int, int) {
	armorValue := item.GetIntValue(armorRegex)
	bonusArmorValue1 := item.GetIntValue(bonusArmorRegex)
	bonusArmorValue2 := item.GetIntValue(bonusArmorRegex2)
	bonusArmorValue := bonusArmorValue1 + bonusArmorValue2

	if item.IsScalableArmorSlot() {
		armorValue -= bonusArmorValue1
	} else if bonusArmorValue2 == 0 {
		bonusArmorValue = armorValue
		armorValue = 0
	}

	return armorValue, bonusArmorValue
}

var armorTypePatterns = map[proto.ArmorType]*regexp.Regexp{
	proto.ArmorType_ArmorTypeCloth:   regexp.MustCompile(`<span class="q1">(?:<!--asc1-->)?Cloth</span>`),
	proto.ArmorType_ArmorTypeLeather: regexp.MustCompile(`<span class="q1">(?:<!--asc2-->)?Leather</span>`),
	proto.ArmorType_ArmorTypeMail:    regexp.MustCompile(`<span class="q1">(?:<!--asc3-->)?Mail</span>`),
	proto.ArmorType_ArmorTypePlate:   regexp.MustCompile(`<span class="q1">(?:<!--asc4-->)?Plate</span>`),
}

func (item WowheadItemResponse) GetArmorType() proto.ArmorType {
	for armorType, pattern := range armorTypePatterns {
		if pattern.MatchString(item.Tooltip) {
			return armorType
		}
	}
	return proto.ArmorType_ArmorTypeUnknown
}

var weaponTypePatterns = map[proto.WeaponType]*regexp.Regexp{
	proto.WeaponType_WeaponTypeAxe:     regexp.MustCompile(`<span class="q1">Axe</span>`),
	proto.WeaponType_WeaponTypeDagger:  regexp.MustCompile(`<span class="q1">Dagger</span>`),
	proto.WeaponType_WeaponTypeFist:    regexp.MustCompile(`<span class="q1">Fist Weapon</span>`),
	proto.WeaponType_WeaponTypeMace:    regexp.MustCompile(`<span class="q1">Mace</span>`),
	proto.WeaponType_WeaponTypeOffHand: regexp.MustCompile(`<td>Held In Off-hand</td>`),
	proto.WeaponType_WeaponTypePolearm: regexp.MustCompile(`<span class="q1">Polearm</span>`),
	proto.WeaponType_WeaponTypeShield:  regexp.MustCompile(`<span class="q1">Shield</span>`),
	proto.WeaponType_WeaponTypeStaff:   regexp.MustCompile(`<span class="q1">Staff</span>`),
	proto.WeaponType_WeaponTypeSword:   regexp.MustCompile(`<span class="q1">Sword</span>`),
}

func (item WowheadItemResponse) GetWeaponType() proto.WeaponType {
	for weaponType, pattern := range weaponTypePatterns {
		if pattern.MatchString(item.Tooltip) {
			return weaponType
		}
	}
	return proto.WeaponType_WeaponTypeUnknown
}

var handTypePatterns = map[proto.HandType]*regexp.Regexp{
	proto.HandType_HandTypeMainHand: regexp.MustCompile(`<td>Main Hand</td>`),
	proto.HandType_HandTypeOneHand:  regexp.MustCompile(`<td>One-Hand</td>`),
	proto.HandType_HandTypeOffHand:  regexp.MustCompile(`<td>((Off Hand)|(Held In Off-hand)|(Held In Off-Hand))</td>`),
	proto.HandType_HandTypeTwoHand:  regexp.MustCompile(`<td>Two-Hand</td>`),
}

func (item WowheadItemResponse) GetHandType() proto.HandType {
	for handType, pattern := range handTypePatterns {
		if pattern.MatchString(item.Tooltip) {
			return handType
		}
	}
	return proto.HandType_HandTypeUnknown
}

var rangedWeaponTypePatterns = map[proto.RangedWeaponType]*regexp.Regexp{
	proto.RangedWeaponType_RangedWeaponTypeBow:      regexp.MustCompile(`<span class="q1">Bow</span>`),
	proto.RangedWeaponType_RangedWeaponTypeCrossbow: regexp.MustCompile(`<span class="q1">Crossbow</span>`),
	proto.RangedWeaponType_RangedWeaponTypeGun:      regexp.MustCompile(`<span class="q1">Gun</span>`),
	proto.RangedWeaponType_RangedWeaponTypeRelic:    regexp.MustCompile(`<td>Relic</td>`),
	proto.RangedWeaponType_RangedWeaponTypeThrown:   regexp.MustCompile(`<span class="q1">Thrown</span>`),
	proto.RangedWeaponType_RangedWeaponTypeWand:     regexp.MustCompile(`<span class="q1">Wand</span>`),
}

func (item WowheadItemResponse) GetRangedWeaponType() proto.RangedWeaponType {
	for rangedWeaponType, pattern := range rangedWeaponTypePatterns {
		if pattern.MatchString(item.Tooltip) {
			return rangedWeaponType
		}
	}
	return proto.RangedWeaponType_RangedWeaponTypeUnknown
}

// Returns min/max of weapon damage
func (item WowheadItemResponse) GetWeaponDamage() (float64, float64) {
	noCommas := strings.ReplaceAll(item.Tooltip, ",", "")
	if matches := weaponDamageRegex.FindStringSubmatch(noCommas); len(matches) > 0 {
		min, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			log.Fatalf("Failed to parse weapon damage: %s", err)
		}
		max, err := strconv.ParseFloat(matches[2], 64)
		if err != nil {
			log.Fatalf("Failed to parse weapon damage: %s", err)
		}
		if min > max {
			log.Fatalf("Invalid weapon damage for item %s: min = %0.1f, max = %0.1f", item.Name, min, max)
		}
		return min, max
	} else if matches := weaponDamageRegex2.FindStringSubmatch(noCommas); len(matches) > 0 {
		val, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			log.Fatalf("Failed to parse weapon damage: %s", err)
		}
		return val, val
	}
	return 0, 0
}

func (item WowheadItemResponse) GetWeaponSpeed() float64 {
	if matches := weaponSpeedRegex.FindStringSubmatch(item.Tooltip); len(matches) > 0 {
		speed, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			log.Fatalf("Failed to parse weapon damage: %s", err)
		}
		return speed
	}
	return 0
}

var gemColorsRegex = regexp.MustCompile("(Meta|Yellow|Blue|Red|Cogwheel|Prismatic) Socket")

func (item WowheadItemResponse) GetGemSockets() []proto.GemColor {
	matches := gemColorsRegex.FindAllStringSubmatch(item.Tooltip, -1)
	if matches == nil {
		return []proto.GemColor{}
	}

	numSockets := len(matches)
	gemColors := make([]proto.GemColor, numSockets)
	for socketIdx, match := range matches {
		gemColorName := "GemColor" + match[1]
		gemColors[socketIdx] = proto.GemColor(proto.GemColor_value[gemColorName])
	}
	return gemColors
}

var socketBonusRegex = regexp.MustCompile(`<span class="q0">Socket Bonus: (.*?)</span>`)
var strengthSocketBonusRegexes = []*regexp.Regexp{regexp.MustCompile(`\+([0-9]+) Strength`)}
var agilitySocketBonusRegexes = []*regexp.Regexp{regexp.MustCompile(`\+([0-9]+) Agility`)}
var staminaSocketBonusRegexes = []*regexp.Regexp{regexp.MustCompile(`\+([0-9]+) Stamina`)}
var intellectSocketBonusRegexes = []*regexp.Regexp{regexp.MustCompile(`\+([0-9]+) Intellect`)}
var spiritSocketBonusRegexes = []*regexp.Regexp{regexp.MustCompile(`\+([0-9]+) Spirit`)}
var spellPowerSocketBonusRegexes = []*regexp.Regexp{regexp.MustCompile(`\+([0-9]+) Spell Power`)}
var spellHitSocketBonusRegexes = []*regexp.Regexp{regexp.MustCompile(`\+([0-9]+) Hit Rating`)}
var spellCritSocketBonusRegexes = []*regexp.Regexp{regexp.MustCompile(`\+([0-9]+) Critical Strike Rating`)}
var hasteSocketBonusRegexes = []*regexp.Regexp{regexp.MustCompile(`\+([0-9]+) Haste Rating`)}
var mp5SocketBonusRegexes = []*regexp.Regexp{
	regexp.MustCompile(`([0-9]+) Mana per 5 sec`),
	regexp.MustCompile(`([0-9]+) mana per 5 sec`),
}
var attackPowerSocketBonusRegexes = []*regexp.Regexp{regexp.MustCompile(`\+([0-9]+) Attack Power`)}
var armorPenSocketBonusRegexes = []*regexp.Regexp{regexp.MustCompile(`\+([0-9]+) Armor Penetration Rating`)}
var expertiseSocketBonusRegexes = []*regexp.Regexp{regexp.MustCompile(`\+([0-9]+) Expertise Rating`)}
var blockSocketBonusRegexes = []*regexp.Regexp{regexp.MustCompile(`\+([0-9]+) Block Rating`)}
var dodgeSocketBonusRegexes = []*regexp.Regexp{regexp.MustCompile(`\+([0-9]+) Dodge Rating`)}
var parrySocketBonusRegexes = []*regexp.Regexp{regexp.MustCompile(`\+([0-9]+) Parry Rating`)}
var resilienceSocketBonusRegexes = []*regexp.Regexp{regexp.MustCompile(`\+([0-9]+) Resilience Rating`)}
var masterySocketBonusRegexes = []*regexp.Regexp{regexp.MustCompile(`\+([0-9]+) Mastery Rating`)}

func (item WowheadItemResponse) GetSocketBonus() stats.Stats {
	match := socketBonusRegex.FindStringSubmatch(item.Tooltip)
	if match == nil {
		return stats.Stats{}
	}

	bonusStr := match[1]
	//fmt.Printf("\n%s\n", bonusStr)

	stats := stats.Stats{
		stats.Strength:          float64(GetBestRegexIntValue(bonusStr, strengthSocketBonusRegexes, 1)),
		stats.Agility:           float64(GetBestRegexIntValue(bonusStr, agilitySocketBonusRegexes, 1)),
		stats.Stamina:           float64(GetBestRegexIntValue(bonusStr, staminaSocketBonusRegexes, 1)),
		stats.Intellect:         float64(GetBestRegexIntValue(bonusStr, intellectSocketBonusRegexes, 1)),
		stats.Spirit:            float64(GetBestRegexIntValue(bonusStr, spiritSocketBonusRegexes, 1)),
		stats.HasteRating:       float64(GetBestRegexIntValue(bonusStr, hasteSocketBonusRegexes, 1)),
		stats.SpellPower:        float64(GetBestRegexIntValue(bonusStr, spellPowerSocketBonusRegexes, 1)),
		stats.HitRating:         float64(GetBestRegexIntValue(bonusStr, spellHitSocketBonusRegexes, 1)),
		stats.CritRating:        float64(GetBestRegexIntValue(bonusStr, spellCritSocketBonusRegexes, 1)),
		stats.MP5:               float64(GetBestRegexIntValue(bonusStr, mp5SocketBonusRegexes, 1)),
		stats.AttackPower:       float64(GetBestRegexIntValue(bonusStr, attackPowerSocketBonusRegexes, 1)),
		stats.RangedAttackPower: float64(GetBestRegexIntValue(bonusStr, attackPowerSocketBonusRegexes, 1)),
		stats.ExpertiseRating:   float64(GetBestRegexIntValue(bonusStr, expertiseSocketBonusRegexes, 1)),
		stats.DodgeRating:       float64(GetBestRegexIntValue(bonusStr, dodgeSocketBonusRegexes, 1)),
		stats.ParryRating:       float64(GetBestRegexIntValue(bonusStr, parrySocketBonusRegexes, 1)),
		stats.ResilienceRating:  float64(GetBestRegexIntValue(bonusStr, resilienceSocketBonusRegexes, 1)),
		stats.MasteryRating:     float64(GetBestRegexIntValue(bonusStr, masterySocketBonusRegexes, 1)),
	}

	return stats
}

var gemSocketColorPatterns = map[proto.GemColor]*regexp.Regexp{
	proto.GemColor_GemColorMeta:      regexp.MustCompile(`Only fits in a meta gem slot\.`),
	proto.GemColor_GemColorBlue:      regexp.MustCompile(`Matches a Blue ([Ss])ocket\.`),
	proto.GemColor_GemColorRed:       regexp.MustCompile(`Matches a Red [Ss]ocket\.`),
	proto.GemColor_GemColorYellow:    regexp.MustCompile(`Matches a Yellow [Ss]ocket\.`),
	proto.GemColor_GemColorOrange:    regexp.MustCompile(`Matches a ((Yellow)|(Red)) or ((Yellow)|(Red)) [Ss]ocket\.`),
	proto.GemColor_GemColorPurple:    regexp.MustCompile(`Matches a ((Blue)|(Red)) or ((Blue)|(Red)) [Ss]ocket\.`),
	proto.GemColor_GemColorGreen:     regexp.MustCompile(`Matches a ((Yellow)|(Blue)) or ((Yellow)|(Blue)) [Ss]ocket\.`),
	proto.GemColor_GemColorPrismatic: regexp.MustCompile(`(Matches any [Ss]ocket)|(Matches a Red, Yellow or Blue [Ss]ocket)`),
	proto.GemColor_GemColorCogwheel:  regexp.MustCompile(`Only fits in a Cogwheel socket.`),
}

func (item WowheadItemResponse) GetSocketColor() proto.GemColor {
	for socketColor, pattern := range gemSocketColorPatterns {
		if pattern.MatchString(item.Tooltip) {
			return socketColor
		}
	}
	// fmt.Printf("Could not find socket color for gem %s\n", item.Name)
	return proto.GemColor_GemColorUnknown
}
func (item WowheadItemResponse) IsGem() bool {
	return item.GetSocketColor() != proto.GemColor_GemColorUnknown &&
		!strings.Contains(item.GetName(), "Design:")
}
func (item WowheadItemResponse) ToItemProto() *proto.UIItem {
	weaponDamageMin, weaponDamageMax := item.GetWeaponDamage()
	return &proto.UIItem{
		Id:   item.ID,
		Name: item.GetName(),
		Icon: item.GetIcon(),

		Type:             item.GetItemType(),
		ArmorType:        item.GetArmorType(),
		WeaponType:       item.GetWeaponType(),
		HandType:         item.GetHandType(),
		RangedWeaponType: item.GetRangedWeaponType(),

		Stats:       item.GetStats().ToProtoArray(),
		GemSockets:  item.GetGemSockets(),
		SocketBonus: item.GetSocketBonus().ToProtoArray(),

		WeaponDamageMin: weaponDamageMin,
		WeaponDamageMax: weaponDamageMax,
		WeaponSpeed:     item.GetWeaponSpeed(),

		Ilvl:    int32(item.GetItemLevel()),
		Phase:   int32(item.GetPhase()),
		Quality: proto.ItemQuality(item.GetQuality()),
		Unique:  item.GetUnique(),
		Heroic:  item.IsHeroic(),

		RequiredProfession: item.GetRequiredProfession(),
		SetName:            item.GetItemSetName(),
	}
}
func (item WowheadItemResponse) ToGemProto() *proto.UIGem {
	return &proto.UIGem{
		Id:    item.ID,
		Name:  item.GetName(),
		Icon:  item.GetIcon(),
		Color: item.GetSocketColor(),

		Stats: item.GetGemStats().ToProtoArray(),

		Phase:              int32(item.GetPhase()),
		Quality:            proto.ItemQuality(item.GetQuality()),
		Unique:             item.GetUnique(),
		RequiredProfession: item.GetRequiredProfession(),
	}
}

var strengthGemStatRegexes = []*regexp.Regexp{regexp.MustCompile(`\+([0-9]+) Strength`), regexp.MustCompile(`\+([0-9]+) (to )?All Stats`)}
var agilityGemStatRegexes = []*regexp.Regexp{regexp.MustCompile(`\+([0-9]+) Agility`), regexp.MustCompile(`\+([0-9]+) (to )?All Stats`)}
var staminaGemStatRegexes = []*regexp.Regexp{regexp.MustCompile(`\+([0-9]+) Stamina`), regexp.MustCompile(`\+([0-9]+) (to )?All Stats`)}
var intellectGemStatRegexes = []*regexp.Regexp{regexp.MustCompile(`\+([0-9]+) Intellect`), regexp.MustCompile(`\+([0-9]+) (to )?All Stats`)}
var spiritGemStatRegexes = []*regexp.Regexp{regexp.MustCompile(`\+([0-9]+) Spirit`), regexp.MustCompile(`\+([0-9]+) (to )?All Stats`)}
var spellPowerGemStatRegexes = []*regexp.Regexp{regexp.MustCompile(`\+([0-9]+) Spell Power`)}
var hitGemStatRegexes = []*regexp.Regexp{regexp.MustCompile(`\+([0-9]+) Hit Rating`)}
var critGemStatRegexes = []*regexp.Regexp{
	regexp.MustCompile(`\+([0-9]+) Crit Rating`),
	regexp.MustCompile(`\+([0-9]+) Critical Strike Rating`),
	regexp.MustCompile(`\+([0-9]+) Critical`),
}
var hasteGemStatRegexes = []*regexp.Regexp{
	regexp.MustCompile(`\+([0-9]+) Haste Rating`),
}
var armorPenetrationGemStatRegexes = []*regexp.Regexp{regexp.MustCompile(`\+([0-9]+) Armor Penetration`)}
var spellPenetrationGemStatRegexes = []*regexp.Regexp{regexp.MustCompile(`\+([0-9]+) Spell Penetration`)}
var mp5GemStatRegexes = []*regexp.Regexp{
	regexp.MustCompile(`([0-9]+) Mana per 5 sec`),
	regexp.MustCompile(`([0-9]+) mana per 5 sec`),
	regexp.MustCompile(`([0-9]+) Mana every 5 seconds`),
}
var attackPowerGemStatRegexes = []*regexp.Regexp{regexp.MustCompile(`\+([0-9]+) Attack Power`)}
var expertiseGemStatRegexes = []*regexp.Regexp{regexp.MustCompile(`\+([0-9]+) Expertise Rating`)}
var dodgeGemStatRegexes = []*regexp.Regexp{regexp.MustCompile(`\+([0-9]+) Dodge Rating`)}
var parryGemStatRegexes = []*regexp.Regexp{regexp.MustCompile(`\+([0-9]+) Parry Rating`)}
var resilienceGemStatRegexes = []*regexp.Regexp{regexp.MustCompile(`\+([0-9]+) Resilience Rating`)}
var allResistGemStatRegexes = []*regexp.Regexp{regexp.MustCompile(`\+([0-9]+) Resist All`)}
var masteryGemStatRegexes = []*regexp.Regexp{regexp.MustCompile(`\+([0-9]+) Mastery Rating`)} //https://www.wowhead.com/mop-classic/spell=101748/zen-elven-peridot

func (item WowheadItemResponse) GetGemStats() stats.Stats {
	stats := stats.Stats{
		stats.Strength:  float64(GetBestRegexIntValue(item.Tooltip, strengthGemStatRegexes, 1)),
		stats.Agility:   float64(GetBestRegexIntValue(item.Tooltip, agilityGemStatRegexes, 1)),
		stats.Stamina:   float64(GetBestRegexIntValue(item.Tooltip, staminaGemStatRegexes, 1)),
		stats.Intellect: float64(GetBestRegexIntValue(item.Tooltip, intellectGemStatRegexes, 1)),
		stats.Spirit:    float64(GetBestRegexIntValue(item.Tooltip, spiritGemStatRegexes, 1)),

		stats.HitRating:     float64(GetBestRegexIntValue(item.Tooltip, hitGemStatRegexes, 1)),
		stats.CritRating:    float64(GetBestRegexIntValue(item.Tooltip, critGemStatRegexes, 1)),
		stats.HasteRating:   float64(GetBestRegexIntValue(item.Tooltip, hasteGemStatRegexes, 1)),
		stats.MasteryRating: float64(GetBestRegexIntValue(item.Tooltip, masteryGemStatRegexes, 1)),

		stats.SpellPower:        float64(GetBestRegexIntValue(item.Tooltip, spellPowerGemStatRegexes, 1)),
		stats.AttackPower:       float64(GetBestRegexIntValue(item.Tooltip, attackPowerGemStatRegexes, 1)),
		stats.RangedAttackPower: float64(GetBestRegexIntValue(item.Tooltip, attackPowerGemStatRegexes, 1)),
		stats.SpellPenetration:  float64(GetBestRegexIntValue(item.Tooltip, spellPenetrationGemStatRegexes, 1)),
		stats.MP5:               float64(GetBestRegexIntValue(item.Tooltip, mp5GemStatRegexes, 1)),
		stats.ExpertiseRating:   float64(GetBestRegexIntValue(item.Tooltip, expertiseGemStatRegexes, 1)),
		stats.DodgeRating:       float64(GetBestRegexIntValue(item.Tooltip, dodgeGemStatRegexes, 1)),
		stats.ParryRating:       float64(GetBestRegexIntValue(item.Tooltip, parryGemStatRegexes, 1)),
		stats.ResilienceRating:  float64(GetBestRegexIntValue(item.Tooltip, resilienceGemStatRegexes, 1)),
		stats.ArcaneResistance:  float64(GetBestRegexIntValue(item.Tooltip, allResistGemStatRegexes, 1)),
		stats.FireResistance:    float64(GetBestRegexIntValue(item.Tooltip, allResistGemStatRegexes, 1)),
		stats.FrostResistance:   float64(GetBestRegexIntValue(item.Tooltip, allResistGemStatRegexes, 1)),
		stats.NatureResistance:  float64(GetBestRegexIntValue(item.Tooltip, allResistGemStatRegexes, 1)),
		stats.ShadowResistance:  float64(GetBestRegexIntValue(item.Tooltip, allResistGemStatRegexes, 1)),
	}

	return stats
}

var itemSetNameRegex = regexp.MustCompile(fmt.Sprintf(`<a href="\/%s\/item-set=-?([0-9]+)\/(.*)" class="q">([^<]+)<`, expansionRegex))

func (item WowheadItemResponse) GetItemSetName() string {
	original := item.GetTooltipRegexString(itemSetNameRegex, 4)

	// Strip out the 10/25 man prefixes from set names
	withoutTier := strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(original, "Heroes' "), "Valorous "), "Conqueror's "), "Triumphant "), "Sanctified ")
	if original != withoutTier { // if we found a tier prefix, return now.
		return withoutTier
	}

	// Now strip out the season prefix from any pvp set names
	withoutPvp := strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(original, "Savage Glad", "Glad", 1), "Hateful Glad", "Glad", 1), "Deadly Glad", "Glad", 1), "Furious Glad", "Glad", 1), "Relentless Glad", "Glad", 1), "Wrathful Glad", "Glad", 1)
	return withoutPvp
}

func (item WowheadItemResponse) IsHeroic() bool {
	return strings.Contains(item.Tooltip, "<span class=\"q2\">Heroic</span>")
}

func (item WowheadItemResponse) GetRequiredProfession() proto.Profession {
	if jcGemsRegex.MatchString(item.Tooltip) {
		return proto.Profession_Jewelcrafting
	}

	return proto.Profession_ProfessionUnknown
}
