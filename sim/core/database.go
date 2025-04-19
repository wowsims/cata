package core

import (
	"fmt"
	"math"
	"slices"
	"sync"
	"time"

	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
	"google.golang.org/protobuf/encoding/protojson"
)

var WITH_DB = false

var ItemsByID = map[int32]Item{}
var GemsByID = map[int32]Gem{}
var RandomSuffixesByID = map[int32]RandomSuffix{}
var EnchantsByEffectID = map[int32]Enchant{}
var ReforgeStatsByID = map[int32]ReforgeStat{}
var ConsumablesByID = map[int32]Consumable{}
var SpellEffectsById = map[int32]*proto.SpellEffect{}
var RandomPropPointsByIlvl = map[int32]*proto.QualityAllocations{}
var ArmorValueDatabase = &proto.ArmorValueDatabase{}
var WeaponValueDatabase = &proto.WeaponDamageDatabase{}
var ItemArmorTotalByIlvl = map[int32]*proto.ItemArmorTotal{}
var mutex = &sync.Mutex{}

func addToDatabase(newDB *proto.SimDatabase) {
	// create mutex lock here and lock it
	// defer unlock it
	mutex.Lock()
	defer mutex.Unlock()

	for _, v := range newDB.Items {
		if _, ok := ItemsByID[v.Id]; !ok {
			ItemsByID[v.Id] = ItemFromProto(v)
		}
	}

	for _, v := range newDB.RandomSuffixes {
		if _, ok := RandomSuffixesByID[v.Id]; !ok {
			RandomSuffixesByID[v.Id] = RandomSuffixFromProto(v)
		}
	}

	for _, v := range newDB.Enchants {
		if _, ok := EnchantsByEffectID[v.EffectId]; !ok {
			EnchantsByEffectID[v.EffectId] = EnchantFromProto(v)
		}
	}

	for _, v := range newDB.Gems {
		if _, ok := GemsByID[v.Id]; !ok {
			GemsByID[v.Id] = GemFromProto(v)
		}
	}

	for _, v := range newDB.ReforgeStats {
		if _, ok := ReforgeStatsByID[v.Id]; !ok {
			ReforgeStatsByID[v.Id] = ReforgeStatFromProto(v)
		}
	}
	for _, v := range newDB.Consumables {
		if _, ok := ConsumablesByID[v.Id]; !ok {
			ConsumablesByID[v.Id] = ConsumableFromProto(v)
		}
	}
	for _, v := range newDB.SpellEffects {
		if _, ok := SpellEffectsById[v.Id]; !ok {
			SpellEffectsById[v.Id] = v
		}
	}
	for i, v := range newDB.RandomPropPoints {
		if _, ok := RandomPropPointsByIlvl[i]; !ok {
			RandomPropPointsByIlvl[i] = v
		}
	}
	for _, v := range newDB.ArmorTotalValue {
		if _, ok := ItemArmorTotalByIlvl[v.ItemLevel]; !ok {
			ItemArmorTotalByIlvl[v.ItemLevel] = v
		}
	}
	ArmorValueDatabase = newDB.ArmorDb
	WeaponValueDatabase = newDB.WeaponDamageDb
}

type ReforgeStat struct {
	ID         int32
	FromStat   proto.Stat
	ToStat     proto.Stat
	Multiplier float64
}

// ReforgeStatFromProto converts a protobuf ReforgeStat to a Go ReforgeStat
func ReforgeStatFromProto(protoStat *proto.ReforgeStat) ReforgeStat {
	return ReforgeStat{
		ID:         protoStat.GetId(),
		FromStat:   protoStat.GetFromStat(),
		ToStat:     protoStat.GetToStat(),
		Multiplier: protoStat.GetMultiplier(),
	}
}

// ReforgeStatToProto converts a Go ReforgeStat to a protobuf ReforgeStat
func ReforgeStatToProto(stat ReforgeStat) *proto.ReforgeStat {
	return &proto.ReforgeStat{
		Id:         stat.ID,
		FromStat:   stat.FromStat,
		ToStat:     stat.ToStat,
		Multiplier: stat.Multiplier,
	}
}

type Consumable struct {
	Id            int32
	Type          proto.ConsumableType
	Stats         stats.Stats
	BuffsMainStat bool
	Name          string
	BuffDuration  time.Duration
	EffectIds     []int32
}

func ConsumableFromProto(consumable *proto.Consumable) Consumable {
	return Consumable{
		Id:            consumable.Id,
		Type:          consumable.Type,
		Stats:         stats.FromProtoArray(consumable.Stats),
		BuffsMainStat: consumable.BuffsMainStat,
		Name:          consumable.Name,
		BuffDuration:  time.Second * time.Duration(consumable.BuffDuration),
		EffectIds:     consumable.EffectIds,
	}
}

type Item struct {
	ID        int32
	Type      proto.ItemType
	ArmorType proto.ArmorType
	// Weapon Stats
	WeaponType       proto.WeaponType
	HandType         proto.HandType
	RangedWeaponType proto.RangedWeaponType
	WeaponDamageMin  float64
	WeaponDamageMax  float64
	SwingSpeed       float64

	Name    string
	Stats   stats.Stats // Stats applied to wearer
	Quality proto.ItemQuality
	SetName string // Empty string if not part of a set.
	SetID   int32  // 0 if not part of a set.

	GemSockets  []proto.GemColor
	SocketBonus stats.Stats

	// Modified for each instance of the item.
	RandomSuffix RandomSuffix
	Gems         []Gem
	Enchant      Enchant
	Reforging    *ReforgeStat

	//Internal use
	TempEnchant     int32
	Ilvl            int32
	DmgVariance     float64
	ArmorModifier   float64
	QualityModifier float64
	SocketModifier  []float64
	StatAllocation  stats.Stats //key=stat id, value=allocation number
}

func (item *Item) UpgradeItemLevelBy(upgradeLevel int) int {
	if item.Quality == proto.ItemQuality_ItemQualityUncommon {
		return upgradeLevel * 8
	}
	if item.Quality == proto.ItemQuality_ItemQualityRare {
		return upgradeLevel * 4
	}
	if item.Quality == proto.ItemQuality_ItemQualityEpic {
		return upgradeLevel * 4
	}
	return 0
}

func ItemFromProto(pData *proto.SimItem) Item {
	return Item{
		ID:               pData.Id,
		Name:             pData.Name,
		Type:             pData.Type,
		ArmorType:        pData.ArmorType,
		WeaponType:       pData.WeaponType,
		Quality:          pData.Quality,
		HandType:         pData.HandType,
		RangedWeaponType: pData.RangedWeaponType,
		WeaponDamageMin:  pData.WeaponDamageMin,
		WeaponDamageMax:  pData.WeaponDamageMax,
		SwingSpeed:       pData.WeaponSpeed,
		Stats:            stats.FromProtoArray(pData.Stats),
		GemSockets:       pData.GemSockets,
		SocketBonus:      stats.FromProtoArray(pData.SocketBonus),
		SetName:          pData.SetName,
		SetID:            pData.SetId,
		Ilvl:             pData.Ilvl,
		DmgVariance:      pData.DmgVariance,
		ArmorModifier:    pData.ArmorModifier,
		QualityModifier:  pData.QualityModifier,
		SocketModifier:   pData.SocketModifier,
		StatAllocation:   stats.FromProtoArray(pData.StatAllocation),
	}
}

func (item *Item) ToItemSpecProto() *proto.ItemSpec {
	itemSpec := &proto.ItemSpec{
		Id:           item.ID,
		RandomSuffix: item.RandomSuffix.ID,
		Enchant:      item.Enchant.EffectID,
		Gems:         MapSlice(item.Gems, func(gem Gem) int32 { return gem.ID }),
	}

	// Check if Reforging is not nil before accessing ID
	// The idea here is to convert a reforging ID to sim stats
	if item.Reforging != nil {
		itemSpec.Reforging = item.Reforging.ID
	} else {
		itemSpec.Reforging = 0
	}

	return itemSpec
}

var qualityReaders = map[proto.ItemQuality]func(*proto.QualityValues) float64{
	proto.ItemQuality_ItemQualityCommon:    func(v *proto.QualityValues) float64 { return v.Common },
	proto.ItemQuality_ItemQualityUncommon:  func(v *proto.QualityValues) float64 { return v.Uncommon },
	proto.ItemQuality_ItemQualityRare:      func(v *proto.QualityValues) float64 { return v.Rare },
	proto.ItemQuality_ItemQualityEpic:      func(v *proto.QualityValues) float64 { return v.Epic },
	proto.ItemQuality_ItemQualityLegendary: func(v *proto.QualityValues) float64 { return v.Legendary },
	proto.ItemQuality_ItemQualityArtifact:  func(v *proto.QualityValues) float64 { return v.Artifact },
	proto.ItemQuality_ItemQualityHeirloom:  func(v *proto.QualityValues) float64 { return v.Heirloom },
}

func (item *Item) ValueForQuality(vals *proto.QualityValues) float64 {
	if fn, ok := qualityReaders[item.Quality]; ok {
		return fn(vals)
	}
	return 0
}

func (item *Item) WeaponDps(itemLevel int32) float64 {
	var ilvl int32
	if itemLevel > 0 {
		ilvl = itemLevel
	} else {
		ilvl = item.Ilvl
	}
	caster := item.Stats[stats.SpellPower] > 0

	var table []*proto.ItemQualityValue

	switch item.Type {
	case proto.ItemType_ItemTypeWeapon:
		if item.HandType == proto.HandType_HandTypeTwoHand {
			if caster {
				table = WeaponValueDatabase.Caster_2H
			} else {
				table = WeaponValueDatabase.Melee_2H
			}
		} else {
			if caster {
				table = WeaponValueDatabase.Caster_1H
			} else {
				table = WeaponValueDatabase.Melee_1H
			}
		}
	case proto.ItemType_ItemTypeRanged:
		switch item.RangedWeaponType {
		case proto.RangedWeaponType_RangedWeaponTypeBow,
			proto.RangedWeaponType_RangedWeaponTypeCrossbow,
			proto.RangedWeaponType_RangedWeaponTypeGun:
			table = WeaponValueDatabase.Ranged

		case proto.RangedWeaponType_RangedWeaponTypeThrown:
			table = WeaponValueDatabase.Thrown

		case proto.RangedWeaponType_RangedWeaponTypeWand:
			table = WeaponValueDatabase.Wand

		default:
			return 0
		}
	}
	idx := int(ilvl)
	if idx < 0 || idx >= len(table) || table[idx] == nil {
		return 0
	}
	return item.ValueForQuality(table[ilvl].Quality)
}

func (item *Item) DamageMin(itemLevel int32) float64 {
	if itemLevel == 0 {
		itemLevel = item.Ilvl
	}
	total := item.WeaponDps(itemLevel)*(float64(item.SwingSpeed))*(1-item.DmgVariance/2) +
		(item.QualityModifier * (float64(item.SwingSpeed)))
	if total < 0 {
		total = 1
	}
	return math.Floor(total)
}

func (item *Item) DamageMax(itemLevel int32) float64 {
	if itemLevel == 0 {
		itemLevel = item.Ilvl
	}
	total := item.WeaponDps(itemLevel)*(float64(item.SwingSpeed))*(1+item.DmgVariance/2) +
		(item.QualityModifier * (float64(item.SwingSpeed)))
	if total < 0 {
		total = 1
	}
	return math.Floor(total + 0.5)
}
func (item *Item) GetArmorValue(itemLevel int32) int {
	if item.Quality > 5 {
		return 0
	}

	var ilvl int32
	if itemLevel > 0 {
		ilvl = itemLevel
	} else {
		ilvl = item.Ilvl
	}
	if item.WeaponType == proto.WeaponType_WeaponTypeShield {
		return int(math.Floor(item.ValueForQuality(ArmorValueDatabase.ShieldArmorValues[ilvl].Quality)) + 0.5)
	}
	if item.ArmorType == proto.ArmorType_ArmorTypeUnknown {
		return 0
	}
	total_armor := 0.0
	quality := 0.0
	armorTotal := ItemArmorTotalByIlvl[ilvl]
	switch item.ArmorType {
	case proto.ArmorType_ArmorTypeCloth:
		total_armor = armorTotal.Cloth
	case proto.ArmorType_ArmorTypeLeather:
		total_armor = armorTotal.Leather
	case proto.ArmorType_ArmorTypeMail:
		total_armor = armorTotal.Mail
	case proto.ArmorType_ArmorTypePlate:
		total_armor = armorTotal.Plate
	}

	quality = item.ValueForQuality(ArmorValueDatabase.ArmorValues[ilvl].Quality)
	return int(math.Floor(total_armor*quality*item.ArmorModifier + 0.5))
}

func (item *Item) GetStats(itemLevel int32) *stats.Stats {
	statsArr := &stats.Stats{}
	for stat := range item.StatAllocation {
		statsArr[stat] = item.GetScaledStat(stats.Stat(stat), itemLevel)
		if stat == int(stats.AttackPower) {
			statsArr[stats.RangedAttackPower] = item.GetScaledStat(stats.Stat(stat), itemLevel)
		}
	}
	return statsArr
}

func (item *Item) GetScaledStat(stat stats.Stat, itemLevel int32) float64 {
	slotType := item.Type
	budget := 0.0

	if slotType == proto.ItemType_ItemTypeUnknown {
		return 0.0
	}

	budget = float64(item.RandPropPoints(itemLevel))
	if item.StatAllocation[stat] > 0 && budget > 0 {
		rawValue := math.Round(item.StatAllocation[stat] * budget * 0.0001)
		return rawValue - item.SocketModifier[stat]
	} else {
		return math.Floor(item.Stats[stat] * item.ApproximateScaleCoeff(item.Ilvl, itemLevel))
	}
}

func (item *Item) ApproximateScaleCoeff(currIlvl int32, newIlvl int32) float64 {
	if currIlvl == 0 || newIlvl == 0 {
		return 1.0
	}

	diff := float64(currIlvl - newIlvl)
	return 1.0 / math.Pow(1.15, diff/15.0)
}

type RandomSuffix struct {
	ID    int32
	Name  string
	Stats stats.Stats
}

func RandomSuffixFromProto(pData *proto.ItemRandomSuffix) RandomSuffix {
	return RandomSuffix{
		ID:    pData.Id,
		Name:  pData.Name,
		Stats: stats.FromProtoArray(pData.Stats),
	}
}

type Enchant struct {
	EffectID int32 // Used by UI to apply effect to tooltip
	Stats    stats.Stats
}

func EnchantFromProto(pData *proto.SimEnchant) Enchant {
	return Enchant{
		EffectID: pData.EffectId,
		Stats:    stats.FromProtoArray(pData.Stats),
	}
}

type Gem struct {
	ID    int32
	Name  string
	Stats stats.Stats
	Color proto.GemColor
}

func GemFromProto(pData *proto.SimGem) Gem {
	return Gem{
		ID:    pData.Id,
		Name:  pData.Name,
		Stats: stats.FromProtoArray(pData.Stats),
		Color: pData.Color,
	}
}

type ItemSpec struct {
	ID            int32
	RandomSuffix  int32
	Enchant       int32
	Gems          []int32
	Reforging     int32
	UpgradeLevel  proto.UpgradeLevel
	ChallengeMode bool
}

type Equipment [proto.ItemSlot_ItemSlotRanged + 1]Item

func (equipment *Equipment) MainHand() *Item {
	return &equipment[proto.ItemSlot_ItemSlotMainHand]
}

func (equipment *Equipment) OffHand() *Item {
	return &equipment[proto.ItemSlot_ItemSlotOffHand]
}

func (equipment *Equipment) Ranged() *Item {
	return &equipment[proto.ItemSlot_ItemSlotRanged]
}

func (equipment *Equipment) Head() *Item {
	return &equipment[proto.ItemSlot_ItemSlotHead]
}

func (equipment *Equipment) Neck() *Item {
	return &equipment[proto.ItemSlot_ItemSlotNeck]
}

func (equipment *Equipment) Shoulder() *Item {
	return &equipment[proto.ItemSlot_ItemSlotShoulder]
}

func (equipment *Equipment) Back() *Item {
	return &equipment[proto.ItemSlot_ItemSlotBack]
}

func (equipment *Equipment) Chest() *Item {
	return &equipment[proto.ItemSlot_ItemSlotChest]
}

func (equipment *Equipment) Wrist() *Item {
	return &equipment[proto.ItemSlot_ItemSlotWrist]
}

func (equipment *Equipment) Hands() *Item {
	return &equipment[proto.ItemSlot_ItemSlotHands]
}

func (equipment *Equipment) Waist() *Item {
	return &equipment[proto.ItemSlot_ItemSlotWaist]
}

func (equipment *Equipment) Legs() *Item {
	return &equipment[proto.ItemSlot_ItemSlotLegs]
}

func (equipment *Equipment) Feet() *Item {
	return &equipment[proto.ItemSlot_ItemSlotFeet]
}

func (equipment *Equipment) Trinket1() *Item {
	return &equipment[proto.ItemSlot_ItemSlotTrinket1]
}

func (equipment *Equipment) Trinket2() *Item {
	return &equipment[proto.ItemSlot_ItemSlotTrinket2]
}

func (equipment *Equipment) Finger1() *Item {
	return &equipment[proto.ItemSlot_ItemSlotFinger1]
}

func (equipment *Equipment) Finger2() *Item {
	return &equipment[proto.ItemSlot_ItemSlotFinger2]
}

func (equipment *Equipment) EquipItem(item Item) {
	if item.Type == proto.ItemType_ItemTypeFinger {
		if equipment.Finger1().ID == 0 {
			*equipment.Finger1() = item
		} else {
			*equipment.Finger2() = item
		}
	} else if item.Type == proto.ItemType_ItemTypeTrinket {
		if equipment.Trinket1().ID == 0 {
			*equipment.Trinket1() = item
		} else {
			*equipment.Trinket2() = item
		}
	} else if item.Type == proto.ItemType_ItemTypeWeapon {
		if item.WeaponType == proto.WeaponType_WeaponTypeShield && equipment.MainHand().HandType != proto.HandType_HandTypeTwoHand {
			*equipment.OffHand() = item
		} else if item.HandType == proto.HandType_HandTypeMainHand || item.HandType == proto.HandType_HandTypeUnknown {
			*equipment.MainHand() = item
		} else if item.HandType == proto.HandType_HandTypeOffHand {
			*equipment.OffHand() = item
		} else if item.HandType == proto.HandType_HandTypeOneHand || item.HandType == proto.HandType_HandTypeTwoHand {
			if equipment.MainHand().ID == 0 {
				*equipment.MainHand() = item
			} else if equipment.OffHand().ID == 0 {
				*equipment.OffHand() = item
			}
		}
	} else {
		equipment[ItemTypeToSlot(item.Type)] = item
	}
}

func (equipment *Equipment) containsEnchantInSlot(effectID int32, slot proto.ItemSlot) bool {
	return (equipment[slot].Enchant.EffectID == effectID) || (equipment[slot].TempEnchant == effectID)
}

func (equipment *Equipment) containsEnchantInSlots(effectID int32, possibleSlots []proto.ItemSlot) bool {
	return slices.ContainsFunc(possibleSlots, func(slot proto.ItemSlot) bool {
		return equipment.containsEnchantInSlot(effectID, slot)
	})
}

func (equipment *Equipment) containsItemInSlots(itemID int32, possibleSlots []proto.ItemSlot) bool {
	return slices.ContainsFunc(possibleSlots, func(slot proto.ItemSlot) bool {
		return equipment[slot].ID == itemID
	})
}

func (equipment *Equipment) ToEquipmentSpecProto() *proto.EquipmentSpec {
	return &proto.EquipmentSpec{
		Items: MapSlice(equipment[:], func(item Item) *proto.ItemSpec {
			return item.ToItemSpecProto()
		}),
	}
}

// Structs used for looking up items/gems/enchants
type EquipmentSpec [proto.ItemSlot_ItemSlotRanged + 1]ItemSpec

func ProtoToEquipmentSpec(es *proto.EquipmentSpec) EquipmentSpec {
	var coreEquip EquipmentSpec
	for i, item := range es.Items {
		coreEquip[i] = ItemSpec{
			ID:            item.Id,
			RandomSuffix:  item.RandomSuffix,
			Enchant:       item.Enchant,
			Gems:          item.Gems,
			Reforging:     item.Reforging,
			UpgradeLevel:  item.UpgradeLevel,
			ChallengeMode: item.ChallengeMode,
		}
	}
	return coreEquip
}

func NewItem(itemSpec ItemSpec) Item {
	item := Item{}
	if foundItem, ok := ItemsByID[itemSpec.ID]; ok {
		item = foundItem
	} else {
		panic(fmt.Sprintf("No item with id: %d", itemSpec.ID))
	}
	if itemSpec.UpgradeLevel != proto.UpgradeLevel_UPGRADE_LEVEL_NONE {

		scaledIlvl := item.Ilvl + int32(item.UpgradeItemLevelBy(int(itemSpec.UpgradeLevel)))

		item.Stats = *item.GetStats(scaledIlvl)
		if item.WeaponType != proto.WeaponType_WeaponTypeUnknown {
			item.WeaponDamageMin = item.DamageMin(scaledIlvl)
			item.WeaponDamageMax = item.DamageMax(scaledIlvl)
		}
		item.Ilvl = scaledIlvl
	}

	if itemSpec.RandomSuffix != 0 {
		if randomSuffix, ok := RandomSuffixesByID[itemSpec.RandomSuffix]; ok {
			item.RandomSuffix = randomSuffix
		} else {
			panic(fmt.Sprintf("No random suffix with id: %d", itemSpec.RandomSuffix))
		}
	}

	if itemSpec.Enchant != 0 {
		if enchant, ok := EnchantsByEffectID[itemSpec.Enchant]; ok {
			item.Enchant = enchant
		}
		// else {
		// 	panic(fmt.Sprintf("No enchant with id: %d", itemSpec.Enchant))
		// }
	}

	if itemSpec.Reforging > 112 { // There is no id below 113
		reforge := ReforgeStatsByID[itemSpec.Reforging]

		if validateReforging(&item, reforge) {
			item.Reforging = &reforge
		} else {
			panic(fmt.Sprintf("When validating reforging for item %d, the reforging could not be validated. %d", itemSpec.ID, itemSpec.Reforging))
		}
	}

	if len(itemSpec.Gems) > 0 {
		// Need to do this to account for possible extra gem sockets.
		numGems := len(item.GemSockets)
		if len(itemSpec.Gems) > numGems {
			numGems = len(itemSpec.Gems)
		}

		item.Gems = make([]Gem, numGems)
		for gemIdx, gemID := range itemSpec.Gems {
			if gem, ok := GemsByID[gemID]; ok {
				item.Gems[gemIdx] = gem
			} else {
				if gemID != 0 {
					panic(fmt.Sprintf("When parsing item %d, socket %d had gem with id: %d\nThis gem is not in the database.", itemSpec.ID, gemIdx, gemID))
				}
			}
		}
	}
	return item
}

func validateReforging(item *Item, reforging ReforgeStat) bool {
	// Validate that the item can reforge these to stats
	reforgeableStats := stats.Stats{}
	if item.RandomSuffix.ID != 0 {
		fmt.Println("item.RandPropPoints(item.ItemLevel)", item.RandPropPoints(item.Ilvl), item.Ilvl)
		reforgeableStats = reforgeableStats.Add(item.RandomSuffix.Stats.Multiply(float64(item.RandPropPoints(item.Ilvl)) / 10000.).Floor())
	} else {
		reforgeableStats = reforgeableStats.Add(item.Stats)
	}

	return (reforgeableStats[reforging.FromStat] > 0) && (reforgeableStats[reforging.ToStat] == 0)
}

func NewEquipmentSet(equipSpec EquipmentSpec) Equipment {
	equipment := Equipment{}
	for _, itemSpec := range equipSpec {
		if itemSpec.ID != 0 {
			equipment.EquipItem(NewItem(itemSpec))
		}
	}
	return equipment
}

func ProtoToEquipment(es *proto.EquipmentSpec) Equipment {
	return NewEquipmentSet(ProtoToEquipmentSpec(es))
}

// Like ItemSpec, but uses names for reference instead of ID.
type ItemStringSpec struct {
	Name    string
	Enchant string
	Gems    []string
}

func EquipmentSpecFromJsonString(jsonString string) *proto.EquipmentSpec {
	es := &proto.EquipmentSpec{}

	data := []byte(jsonString)
	if err := protojson.Unmarshal(data, es); err != nil {
		panic(err)
	}
	return es
}

func ItemSwapFromJsonString(jsonString string) *proto.ItemSwap {
	is := &proto.ItemSwap{}

	data := []byte(jsonString)
	if err := protojson.Unmarshal(data, is); err != nil {
		panic(err)
	}
	return is
}

func (equipment *Equipment) Stats() stats.Stats {
	equipStats := stats.Stats{}

	for _, item := range equipment {
		equipStats = equipStats.Add(ItemEquipmentStats(item))
	}
	return equipStats
}

func (e *Item) RandPropPoints(itemLevel int32) int32 {
	ilvl := e.Ilvl
	if itemLevel != e.Ilvl {
		ilvl = itemLevel
	}
	allocs := RandomPropPointsByIlvl[ilvl]
	fmt.Println("ALLOCS", allocs)
	bucket := pickQualityArray(allocs, e.Quality)
	fmt.Println("BUCKET", bucket)
	if len(bucket) == 0 {
		return 0
	}
	idx := calculateRandomPropSlotIndex(e)
	fmt.Println("INDEX", idx)
	return bucket[idx]
}

func ItemEquipmentStats(item Item) stats.Stats {
	equipStats := stats.Stats{}

	if item.ID == 0 {
		return equipStats
	}

	equipStats = equipStats.Add(item.Stats)

	// Random suffix stats can be Reforged, so apply those prior to any Reforges
	rawSuffixStats := item.RandomSuffix.Stats

	equipStats = equipStats.Add(rawSuffixStats.Multiply(float64(item.RandPropPoints(item.Ilvl)) / 10000.).Floor())

	// Apply reforging
	if item.Reforging != nil {
		itemStats := item.Stats.Add(rawSuffixStats.Multiply(float64(item.RandPropPoints(item.Ilvl)) / 10000.).Floor())
		reforgingChanges := stats.Stats{}
		fromStat := item.Reforging.FromStat

		if itemStats[fromStat] > 0 {
			reduction := math.Floor(itemStats[fromStat] * item.Reforging.Multiplier)
			reforgingChanges[fromStat] = -reduction
			reforgingChanges[item.Reforging.ToStat] = reduction
		}

		equipStats = equipStats.Add(reforgingChanges)
	}

	equipStats = equipStats.Add(item.Enchant.Stats)

	for _, gem := range item.Gems {
		equipStats = equipStats.Add(gem.Stats)
	}
	// Check socket bonus
	if len(item.GemSockets) > 0 && len(item.Gems) >= len(item.GemSockets) {
		allMatch := true
		for gemIndex, socketColor := range item.GemSockets {
			if !ColorIntersects(socketColor, item.Gems[gemIndex].Color) {
				allMatch = false
				break
			}
		}

		if allMatch {
			equipStats = equipStats.Add(item.SocketBonus)
		}
	}

	return equipStats
}

func GetItemByID(id int32) *Item {
	if item, ok := ItemsByID[id]; ok {
		return &item
	}
	return nil
}

func ItemTypeToSlot(it proto.ItemType) proto.ItemSlot {
	switch it {
	case proto.ItemType_ItemTypeHead:
		return proto.ItemSlot_ItemSlotHead
	case proto.ItemType_ItemTypeNeck:
		return proto.ItemSlot_ItemSlotNeck
	case proto.ItemType_ItemTypeShoulder:
		return proto.ItemSlot_ItemSlotShoulder
	case proto.ItemType_ItemTypeBack:
		return proto.ItemSlot_ItemSlotBack
	case proto.ItemType_ItemTypeChest:
		return proto.ItemSlot_ItemSlotChest
	case proto.ItemType_ItemTypeWrist:
		return proto.ItemSlot_ItemSlotWrist
	case proto.ItemType_ItemTypeHands:
		return proto.ItemSlot_ItemSlotHands
	case proto.ItemType_ItemTypeWaist:
		return proto.ItemSlot_ItemSlotWaist
	case proto.ItemType_ItemTypeLegs:
		return proto.ItemSlot_ItemSlotLegs
	case proto.ItemType_ItemTypeFeet:
		return proto.ItemSlot_ItemSlotFeet
	case proto.ItemType_ItemTypeFinger:
		return proto.ItemSlot_ItemSlotFinger1
	case proto.ItemType_ItemTypeTrinket:
		return proto.ItemSlot_ItemSlotTrinket1
	case proto.ItemType_ItemTypeWeapon:
		return proto.ItemSlot_ItemSlotMainHand
	case proto.ItemType_ItemTypeRanged:
		return proto.ItemSlot_ItemSlotRanged
	}

	return 255
}

// See getEligibleItemSlots in proto_utils/utils.ts.
var itemTypeToSlotsMap = map[proto.ItemType][]proto.ItemSlot{
	proto.ItemType_ItemTypeHead:     {proto.ItemSlot_ItemSlotHead},
	proto.ItemType_ItemTypeNeck:     {proto.ItemSlot_ItemSlotNeck},
	proto.ItemType_ItemTypeShoulder: {proto.ItemSlot_ItemSlotShoulder},
	proto.ItemType_ItemTypeBack:     {proto.ItemSlot_ItemSlotBack},
	proto.ItemType_ItemTypeChest:    {proto.ItemSlot_ItemSlotChest},
	proto.ItemType_ItemTypeWrist:    {proto.ItemSlot_ItemSlotWrist},
	proto.ItemType_ItemTypeHands:    {proto.ItemSlot_ItemSlotHands},
	proto.ItemType_ItemTypeWaist:    {proto.ItemSlot_ItemSlotWaist},
	proto.ItemType_ItemTypeLegs:     {proto.ItemSlot_ItemSlotLegs},
	proto.ItemType_ItemTypeFeet:     {proto.ItemSlot_ItemSlotFeet},
	proto.ItemType_ItemTypeFinger:   {proto.ItemSlot_ItemSlotFinger1, proto.ItemSlot_ItemSlotFinger2},
	proto.ItemType_ItemTypeTrinket:  {proto.ItemSlot_ItemSlotTrinket1, proto.ItemSlot_ItemSlotTrinket2},
	proto.ItemType_ItemTypeRanged:   {proto.ItemSlot_ItemSlotRanged},
	// ItemType_ItemTypeWeapon is excluded intentionally - the slot cannot be decided based on type alone for weapons.
}

func eligibleSlotsForItem(item *Item, isFuryWarrior bool) []proto.ItemSlot {
	if item == nil {
		return nil
	}

	if slots, ok := itemTypeToSlotsMap[item.Type]; ok {
		return slots
	}

	if item.Type == proto.ItemType_ItemTypeWeapon {
		if isFuryWarrior {
			return []proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand, proto.ItemSlot_ItemSlotOffHand}
		}

		switch item.HandType {
		case proto.HandType_HandTypeTwoHand, proto.HandType_HandTypeMainHand:
			return []proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand}
		case proto.HandType_HandTypeOffHand:
			return []proto.ItemSlot{proto.ItemSlot_ItemSlotOffHand}
		case proto.HandType_HandTypeOneHand:
			return []proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand, proto.ItemSlot_ItemSlotOffHand}
		}
	}

	return nil
}

func ColorIntersects(g proto.GemColor, o proto.GemColor) bool {
	if g == o {
		return true
	}
	if g == proto.GemColor_GemColorPrismatic || o == proto.GemColor_GemColorPrismatic {
		return true
	}
	if g == proto.GemColor_GemColorMeta {
		return o == proto.GemColor_GemColorMeta
	}
	if g == proto.GemColor_GemColorRed {
		return o == proto.GemColor_GemColorOrange || o == proto.GemColor_GemColorPurple
	}
	if g == proto.GemColor_GemColorBlue {
		return o == proto.GemColor_GemColorGreen || o == proto.GemColor_GemColorPurple
	}
	if g == proto.GemColor_GemColorYellow {
		return o == proto.GemColor_GemColorGreen || o == proto.GemColor_GemColorOrange
	}
	if g == proto.GemColor_GemColorOrange {
		return o == proto.GemColor_GemColorYellow || o == proto.GemColor_GemColorRed
	}
	if g == proto.GemColor_GemColorGreen {
		return o == proto.GemColor_GemColorYellow || o == proto.GemColor_GemColorBlue
	}
	if g == proto.GemColor_GemColorPurple {
		return o == proto.GemColor_GemColorBlue || o == proto.GemColor_GemColorRed
	}

	return false // dunno what else could be.
}
