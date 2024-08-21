package database

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

type CataTooltipManager struct {
	TooltipManager
}

func (wtm *CataTooltipManager) Read() map[int32]CataItemResponse {
	strDB := wtm.TooltipManager.Read()
	return core.MapMap(strDB, func(id int32, tooltip string) (int32, CataItemResponse) {
		// Reformat the tooltip so it looks more like a wowhead tooltip.
		tooltip = strings.Replace(tooltip, fmt.Sprintf("$WowheadPower.registerItem('%d', 0, ", id), "", 1)
		tooltip = strings.TrimSuffix(tooltip, ";")
		tooltip = strings.TrimSuffix(tooltip, ")")
		tooltip = strings.ReplaceAll(tooltip, "\n", "")
		tooltip = strings.ReplaceAll(tooltip, "\t", "")
		tooltip = strings.Replace(tooltip, "name_enus: '", "\"name\": \"", 1)
		tooltip = strings.Replace(tooltip, "quality:", "\"quality\":", 1)
		tooltip = strings.Replace(tooltip, "icon: '", "\"icon\": \"", 1)
		tooltip = strings.Replace(tooltip, "tooltip_enus: '", "\"tooltip\": \"", 1)
		tooltip = strings.ReplaceAll(tooltip, "',", "\",")
		tooltip = strings.ReplaceAll(tooltip, "\\'", "'")
		// replace the '} with "}
		if strings.HasSuffix(tooltip, "'}") {
			tooltip = tooltip[:len(tooltip)-2] + "\"}"
		}

		return id, NewCataItemResponse(id, tooltip)
	})
}

func NewCataItemTooltipManager(filePath string) *CataTooltipManager {
	return &CataTooltipManager{
		TooltipManager{
			FilePath:   filePath,
			UrlPattern: "https://cata.evowow.com/?item=%s&power",
		},
	}
}

type CataItemResponse struct {
	ID      int32
	Name    string `json:"name"`
	Quality int    `json:"quality"`
	Icon    string `json:"icon"`
	Tooltip string `json:"tooltip"`
}

func NewCataItemResponse(id int32, tooltip string) CataItemResponse {
	response := CataItemResponse{}
	err := json.Unmarshal([]byte(tooltip), &response)
	if err != nil {
		fmt.Printf("Failed to decode tooltipBytes: %s\n", tooltip)
		panic(err)
	}
	response.ID = id
	return response
}

func (item CataItemResponse) GetName() string {
	return item.Name
}
func (item CataItemResponse) GetQuality() int {
	return item.Quality
}
func (item CataItemResponse) GetIcon() string {
	return item.Icon
}

func (item CataItemResponse) TooltipWithoutSetBonus() string {
	setIdx := strings.Index(item.Tooltip, "Set : ")
	if setIdx == -1 {
		return item.Tooltip
	} else {
		return item.Tooltip[:setIdx]
	}
}

func (item CataItemResponse) GetTooltipRegexString(pattern *regexp.Regexp, matchIdx int) string {
	return GetRegexStringValue(item.TooltipWithoutSetBonus(), pattern, matchIdx)
}

func (item CataItemResponse) GetTooltipRegexValue(pattern *regexp.Regexp, matchIdx int) int {
	return GetRegexIntValue(item.TooltipWithoutSetBonus(), pattern, matchIdx)
}

func (item CataItemResponse) GetIntValue(pattern *regexp.Regexp) int {
	return item.GetTooltipRegexValue(pattern, 1)
}

// TODO: Cata update regexes
var wotlkdbArmorRegex = regexp.MustCompile("<!--amr-->([0-9]+) Armor")
var wotlkdbAgilityRegex = regexp.MustCompile(`<!--stat3-->\+([0-9]+) Agility`)
var wotlkdbStrengthRegex = regexp.MustCompile(`<!--stat4-->\+([0-9]+) Strength`)
var wotlkdbIntellectRegex = regexp.MustCompile(`<!--stat5-->\+([0-9]+) Intellect`)
var wotlkdbSpiritRegex = regexp.MustCompile(`<!--stat6-->\+([0-9]+) Spirit`)
var wotlkdbStaminaRegex = regexp.MustCompile(`<!--stat7-->\+([0-9]+) Stamina`)
var wotlkdbSpellPowerRegex = regexp.MustCompile("Equip: Increases spell power by ([0-9]+)")
var wotlkdbSpellPowerRegex2 = regexp.MustCompile("Equip: Increases spell power by <!--rtg45-->([0-9]+)")

var wotlkdbHitRegex = regexp.MustCompile("Improves hit rating by <!--rtg31-->([0-9]+)")
var wotlkdbCritRegex = regexp.MustCompile("Improves critical strike rating by <!--rtg32-->([0-9]+)")
var wotlkdbHasteRegex = regexp.MustCompile("Increases your haste rating by <!--rtg36-->([0-9]+)")

var wotlkdbSpellPenetrationRegex = regexp.MustCompile("Increases your spell penetration by ([0-9]+)")
var wotlkdbMp5Regex = regexp.MustCompile("Restores ([0-9]+) mana per 5 sec")
var wotlkdbAttackPowerRegex = regexp.MustCompile(`Increases attack power by ([0-9]+)\.`)
var wotlkdbAttackPowerRegex2 = regexp.MustCompile(`Increases attack power by <!--rtg38-->([0-9]+)\.`)
var wotlkdbRangedAttackPowerRegex = regexp.MustCompile("Increases ranged attack power by ([0-9]+)")
var wotlkdbExpertiseRegex = regexp.MustCompile("Increases expertise rating by <!--rtg37-->([0-9]+)")

var wotlkdbBlockRegex = regexp.MustCompile(`Equip: Increases your shield block rating by <!--rtg15-->([0-9]+)`)
var wotlkdbBlockRegex2 = regexp.MustCompile("Equip: Increases your shield block rating by ([0-9]+)")
var wotlkdbDodgeRegex = regexp.MustCompile("Increases your dodge rating by <!--rtg13-->([0-9]+)")
var wotlkdbDodgeRegex2 = regexp.MustCompile("Increases your dodge rating by ([0-9]+)")
var wotlkdbParryRegex = regexp.MustCompile("Increases your parry rating by <!--rtg14-->([0-9]+)")
var wotlkdbParryRegex2 = regexp.MustCompile("Increases your parry rating by ([0-9]+)")
var wotlkdbResilienceRegex = regexp.MustCompile("Increases your resilience rating by <!--rtg35-->([0-9]+)")
var wotlkdbArcaneResistanceRegex = regexp.MustCompile(`\+([0-9]+) Arcane Resistance`)
var wotlkdbFireResistanceRegex = regexp.MustCompile(`\+([0-9]+) Fire Resistance`)
var wotlkdbFrostResistanceRegex = regexp.MustCompile(`\+([0-9]+) Frost Resistance`)
var wotlkdbNatureResistanceRegex = regexp.MustCompile(`\+([0-9]+) Nature Resistance`)
var wotlkdbShadowResistanceRegex = regexp.MustCompile(`\+([0-9]+) Shadow Resistance`)

func (item CataItemResponse) GetStats() stats.Stats {
	sp := float64(item.GetIntValue(wotlkdbSpellPowerRegex)) + float64(item.GetIntValue(wotlkdbSpellPowerRegex2))
	return stats.Stats{
		stats.Armor:             float64(item.GetIntValue(wotlkdbArmorRegex)),
		stats.Strength:          float64(item.GetIntValue(wotlkdbStrengthRegex)),
		stats.Agility:           float64(item.GetIntValue(wotlkdbAgilityRegex)),
		stats.Stamina:           float64(item.GetIntValue(wotlkdbStaminaRegex)),
		stats.Intellect:         float64(item.GetIntValue(wotlkdbIntellectRegex)),
		stats.Spirit:            float64(item.GetIntValue(wotlkdbSpiritRegex)),
		stats.SpellPower:        sp,
		stats.HitRating:         float64(item.GetIntValue(wotlkdbHitRegex)),
		stats.CritRating:        float64(item.GetIntValue(wotlkdbCritRegex)),
		stats.HasteRating:       float64(item.GetIntValue(wotlkdbHasteRegex)),
		stats.SpellPenetration:  float64(item.GetIntValue(wotlkdbSpellPenetrationRegex)),
		stats.MP5:               float64(item.GetIntValue(wotlkdbMp5Regex)),
		stats.AttackPower:       float64(item.GetIntValue(wotlkdbAttackPowerRegex) + item.GetIntValue(wotlkdbAttackPowerRegex2)),
		stats.RangedAttackPower: float64(item.GetIntValue(wotlkdbAttackPowerRegex) + item.GetIntValue(wotlkdbAttackPowerRegex2) + item.GetIntValue(wotlkdbRangedAttackPowerRegex)),
		stats.ExpertiseRating:   float64(item.GetIntValue(wotlkdbExpertiseRegex)),
		stats.DodgeRating:       float64(item.GetIntValue(wotlkdbDodgeRegex) + item.GetIntValue(wotlkdbDodgeRegex2)),
		stats.ParryRating:       float64(item.GetIntValue(wotlkdbParryRegex) + item.GetIntValue(wotlkdbParryRegex2)),
		stats.ResilienceRating:  float64(item.GetIntValue(wotlkdbResilienceRegex)),
		stats.ArcaneResistance:  float64(item.GetIntValue(wotlkdbArcaneResistanceRegex)),
		stats.FireResistance:    float64(item.GetIntValue(wotlkdbFireResistanceRegex)),
		stats.FrostResistance:   float64(item.GetIntValue(wotlkdbFrostResistanceRegex)),
		stats.NatureResistance:  float64(item.GetIntValue(wotlkdbNatureResistanceRegex)),
		stats.ShadowResistance:  float64(item.GetIntValue(wotlkdbShadowResistanceRegex)),
	}
}

func (item CataItemResponse) IsPattern() bool {
	for _, pattern := range patternRegexes {
		if pattern.MatchString(item.Tooltip) {
			return true
		}
	}
	return false
}

func (item CataItemResponse) IsRandomEnchant() bool {
	return randomEnchantRegex.MatchString(item.Tooltip)
}

func (item CataItemResponse) IsEquippable() bool {
	return item.GetItemType() != proto.ItemType_ItemTypeUnknown &&
		!item.IsPattern() &&
		!item.IsRandomEnchant()
}

var wotlkItemLevelRegex = regexp.MustCompile("Item Level ([0-9]+)<")

func (item CataItemResponse) GetItemLevel() int {
	return item.GetIntValue(wotlkItemLevelRegex)
}

// WOTLK DB has no phase info
func (item CataItemResponse) GetPhase() int {

	ilvl := item.GetItemLevel()
	if ilvl < 200 || ilvl == 200 || ilvl == 213 || ilvl == 226 {
		return 1
	} else if ilvl == 219 || ilvl == 226 || ilvl == 239 {
		return 2
	} else if ilvl == 232 || ilvl == 245 || ilvl == 258 {
		return 3
	} else if ilvl == 251 || ilvl == 258 || ilvl == 259 || ilvl == 264 || ilvl == 268 || ilvl == 270 || ilvl == 271 || ilvl == 272 {
		return 4
	} else if ilvl == 277 || ilvl == 284 {
		return 5
	}

	// default to 1
	return 1
}

func (item CataItemResponse) GetUnique() bool {
	return uniqueRegex.MatchString(item.Tooltip) && !jcGemsRegex.MatchString(item.Tooltip)
}

func (item CataItemResponse) GetItemType() proto.ItemType {
	for itemType, pattern := range itemTypePatterns {
		if pattern.MatchString(item.Tooltip) {
			return itemType
		}
	}
	return proto.ItemType_ItemTypeUnknown
}

var wotlkArmorTypePatterns = map[proto.ArmorType]*regexp.Regexp{
	proto.ArmorType_ArmorTypeCloth:   regexp.MustCompile("<th><!--asc1-->Cloth</th>"),
	proto.ArmorType_ArmorTypeLeather: regexp.MustCompile("<th><!--asc2-->Leather</th>"),
	proto.ArmorType_ArmorTypeMail:    regexp.MustCompile("<th><!--asc3-->Mail</th>"),
	proto.ArmorType_ArmorTypePlate:   regexp.MustCompile("<th><!--asc4-->Plate</th>"),
}

func (item CataItemResponse) GetArmorType() proto.ArmorType {
	for armorType, pattern := range wotlkArmorTypePatterns {
		if pattern.MatchString(item.Tooltip) {
			return armorType
		}
	}
	return proto.ArmorType_ArmorTypeUnknown
}

var wotlkWeaponTypePatterns = map[proto.WeaponType]*regexp.Regexp{
	proto.WeaponType_WeaponTypeAxe:     regexp.MustCompile("<th>Axe</th>"),
	proto.WeaponType_WeaponTypeDagger:  regexp.MustCompile("<th>Dagger</th>"),
	proto.WeaponType_WeaponTypeFist:    regexp.MustCompile("<th>Fist Weapon</th>"),
	proto.WeaponType_WeaponTypeMace:    regexp.MustCompile("<th>Mace</th>"),
	proto.WeaponType_WeaponTypeOffHand: regexp.MustCompile("<td>Held In Off-Hand</td>"),
	proto.WeaponType_WeaponTypePolearm: regexp.MustCompile("<th>Polearm</th>"),
	proto.WeaponType_WeaponTypeShield:  regexp.MustCompile("<th><!--asc6-->Shield</th>"),
	proto.WeaponType_WeaponTypeStaff:   regexp.MustCompile("<th>Staff</th>"),
	proto.WeaponType_WeaponTypeSword:   regexp.MustCompile("<th>Sword</th>"),
}

func (item CataItemResponse) GetWeaponType() proto.WeaponType {
	for weaponType, pattern := range wotlkWeaponTypePatterns {
		if pattern.MatchString(item.Tooltip) {
			return weaponType
		}
	}
	return proto.WeaponType_WeaponTypeUnknown
}

func (item CataItemResponse) GetHandType() proto.HandType {
	for handType, pattern := range handTypePatterns {
		if pattern.MatchString(item.Tooltip) {
			return handType
		}
	}
	return proto.HandType_HandTypeUnknown
}

var wotlkRangedWeaponTypePatterns = map[proto.RangedWeaponType]*regexp.Regexp{
	proto.RangedWeaponType_RangedWeaponTypeBow:      regexp.MustCompile("<th>Bow</th>"),
	proto.RangedWeaponType_RangedWeaponTypeCrossbow: regexp.MustCompile("<th>Crossbow</th>"),
	proto.RangedWeaponType_RangedWeaponTypeGun:      regexp.MustCompile("<th>Gun</th>"),
	proto.RangedWeaponType_RangedWeaponTypeRelic:    regexp.MustCompile("<td>Relic</td>"),
	proto.RangedWeaponType_RangedWeaponTypeThrown:   regexp.MustCompile("<th>Thrown</th>"),
	proto.RangedWeaponType_RangedWeaponTypeWand:     regexp.MustCompile("<th>Wand</th>"),
}

func (item CataItemResponse) GetRangedWeaponType() proto.RangedWeaponType {
	for rangedWeaponType, pattern := range wotlkRangedWeaponTypePatterns {
		if pattern.MatchString(item.Tooltip) {
			return rangedWeaponType
		}
	}
	return proto.RangedWeaponType_RangedWeaponTypeUnknown
}

// Returns min/max of weapon damage
func (item CataItemResponse) GetWeaponDamage() (float64, float64) {
	if matches := weaponDamageRegex.FindStringSubmatch(item.Tooltip); len(matches) > 0 {
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
	} else if matches := weaponDamageRegex2.FindStringSubmatch(item.Tooltip); len(matches) > 0 {
		val, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			log.Fatalf("Failed to parse weapon damage: %s", err)
		}
		return val, val
	}
	return 0, 0
}

func (item CataItemResponse) GetWeaponSpeed() float64 {
	if matches := weaponSpeedRegex.FindStringSubmatch(item.Tooltip); len(matches) > 0 {
		speed, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			log.Fatalf("Failed to parse weapon damage: %s", err)
		}
		return speed
	}
	return 0
}

func (item CataItemResponse) GetGemSockets() []proto.GemColor {
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

func (item CataItemResponse) GetSocketBonus() stats.Stats {
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
	}

	return stats
}

func (item CataItemResponse) GetSocketColor() proto.GemColor {
	for socketColor, pattern := range gemSocketColorPatterns {
		if pattern.MatchString(item.Tooltip) {
			return socketColor
		}
	}
	// fmt.Printf("Could not find socket color for gem %s\n", item.Name)
	return proto.GemColor_GemColorUnknown
}

func (item CataItemResponse) GetGemStats() stats.Stats {
	stats := stats.Stats{
		stats.Strength:  float64(GetBestRegexIntValue(item.Tooltip, strengthGemStatRegexes, 1)),
		stats.Agility:   float64(GetBestRegexIntValue(item.Tooltip, agilityGemStatRegexes, 1)),
		stats.Stamina:   float64(GetBestRegexIntValue(item.Tooltip, staminaGemStatRegexes, 1)),
		stats.Intellect: float64(GetBestRegexIntValue(item.Tooltip, intellectGemStatRegexes, 1)),
		stats.Spirit:    float64(GetBestRegexIntValue(item.Tooltip, spiritGemStatRegexes, 1)),

		stats.HitRating:   float64(GetBestRegexIntValue(item.Tooltip, hitGemStatRegexes, 1)),
		stats.CritRating:  float64(GetBestRegexIntValue(item.Tooltip, critGemStatRegexes, 1)),
		stats.HasteRating: float64(GetBestRegexIntValue(item.Tooltip, hasteGemStatRegexes, 1)),

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

// TODO: Cata check regex
var wotlkItemSetNameRegex = regexp.MustCompile("<a href=\\\"\\?itemset=([0-9]+)\\\" class=\\\"q\\\">([^<]+)<")

func (item CataItemResponse) GetItemSetName() string {
	return item.GetTooltipRegexString(wotlkItemSetNameRegex, 2)
}

func (item CataItemResponse) IsHeroic() bool {
	return strings.Contains(item.Tooltip, "<span class=\"q2\">Heroic</span>")
}

func (item CataItemResponse) GetRequiredProfession() proto.Profession {
	if jcGemsRegex.MatchString(item.Tooltip) {
		return proto.Profession_Jewelcrafting
	}

	return proto.Profession_ProfessionUnknown
}
