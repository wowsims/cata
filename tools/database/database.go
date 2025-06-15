package database

import (
	"bytes"
	"log"
	"os"
	"slices"

	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/tools"
	"golang.org/x/exp/maps"
	"google.golang.org/protobuf/encoding/protojson"
	googleProto "google.golang.org/protobuf/proto"
)

type EnchantDBKey struct {
	EffectID int32
	ItemID   int32
	SpellID  int32
}

func EnchantToDBKey(enchant *proto.UIEnchant) EnchantDBKey {
	return EnchantDBKey{
		EffectID: enchant.EffectId,
		ItemID:   enchant.ItemId,
		SpellID:  enchant.SpellId,
	}
}

type WowDatabase struct {
	Items          map[int32]*proto.UIItem
	RandomSuffixes map[int32]*proto.ItemRandomSuffix
	Enchants       map[EnchantDBKey]*proto.UIEnchant
	Gems           map[int32]*proto.UIGem

	Zones map[int32]*proto.UIZone
	Npcs  map[int32]*proto.UINPC

	ItemIcons                map[int32]*proto.IconData
	SpellIcons               map[int32]*proto.IconData
	ReforgeStats             map[int32]*proto.ReforgeStat
	ItemEffectRandPropPoints map[int32]*proto.ItemEffectRandPropPoints

	Encounters []*proto.PresetEncounter
	GlyphIDs   []*proto.GlyphID

	Consumables map[int32]*proto.Consumable
	Effects     map[int32]*proto.SpellEffect
}

func NewWowDatabase() *WowDatabase {
	return &WowDatabase{
		Items:          make(map[int32]*proto.UIItem),
		RandomSuffixes: make(map[int32]*proto.ItemRandomSuffix),
		Enchants:       make(map[EnchantDBKey]*proto.UIEnchant),
		Gems:           make(map[int32]*proto.UIGem),
		Zones:          make(map[int32]*proto.UIZone),
		Npcs:           make(map[int32]*proto.UINPC),

		ItemIcons:                make(map[int32]*proto.IconData),
		SpellIcons:               make(map[int32]*proto.IconData),
		ReforgeStats:             make(map[int32]*proto.ReforgeStat),
		ItemEffectRandPropPoints: make(map[int32]*proto.ItemEffectRandPropPoints),

		Consumables: make(map[int32]*proto.Consumable),
		Effects:     make(map[int32]*proto.SpellEffect),
	}
}

func (db *WowDatabase) Clone() *WowDatabase {
	return &WowDatabase{
		Items:          maps.Clone(db.Items),
		RandomSuffixes: maps.Clone(db.RandomSuffixes),
		Enchants:       maps.Clone(db.Enchants),
		Gems:           maps.Clone(db.Gems),
		Zones:          maps.Clone(db.Zones),
		Npcs:           maps.Clone(db.Npcs),

		ItemIcons:                maps.Clone(db.ItemIcons),
		SpellIcons:               maps.Clone(db.SpellIcons),
		ReforgeStats:             maps.Clone(db.ReforgeStats),
		ItemEffectRandPropPoints: maps.Clone(db.ItemEffectRandPropPoints),

		Consumables: maps.Clone(db.Consumables),
		Effects:     maps.Clone(db.Effects),
	}
}

func (db *WowDatabase) MergeItems(arr []*proto.UIItem) {
	for _, item := range arr {
		db.MergeItem(item)
	}
}
func (db *WowDatabase) MergeItem(src *proto.UIItem) {
	if dst, ok := db.Items[src.Id]; ok {
		// googleproto.Merge concatenates lists, but we want replacement, so do them manually.
		if src.Stats != nil {
			dst.Stats = src.Stats
			src.Stats = nil
		}
		if src.SocketBonus != nil {
			dst.SocketBonus = src.SocketBonus
			src.SocketBonus = nil
		}
		googleProto.Merge(dst, src)
	} else {
		db.Items[src.Id] = src
	}
}

func (db *WowDatabase) MergeEnchants(arr []*proto.UIEnchant) {
	for _, enchant := range arr {
		db.MergeEnchant(enchant)
	}
}

func (db *WowDatabase) MergeEnchant(src *proto.UIEnchant) {
	key := EnchantToDBKey(src)
	if dst, ok := db.Enchants[key]; ok {
		// googleproto.Merge concatenates lists, but we want replacement, so do them manually.
		if src.Stats != nil {
			dst.Stats = src.Stats
			src.Stats = nil
		}
		googleProto.Merge(dst, src)
	} else {
		db.Enchants[key] = src
	}
}

func (db *WowDatabase) MergeGems(arr []*proto.UIGem) {
	for _, gem := range arr {
		db.MergeGem(gem)
	}
}
func (db *WowDatabase) MergeGem(src *proto.UIGem) {
	if dst, ok := db.Gems[src.Id]; ok {
		// googleproto.Merge concatenates lists, but we want replacement, so do them manually.
		if src.Stats != nil {
			dst.Stats = src.Stats
			src.Stats = nil
		}
		googleProto.Merge(dst, src)
	} else {
		db.Gems[src.Id] = src
	}
}

func (db *WowDatabase) MergeZones(arr []*proto.UIZone) {
	for _, zone := range arr {
		db.MergeZone(zone)
	}
}
func (db *WowDatabase) MergeZone(src *proto.UIZone) {
	if dst, ok := db.Zones[src.Id]; ok {
		googleProto.Merge(dst, src)
	} else {
		db.Zones[src.Id] = src
	}
}

func (db *WowDatabase) MergeNpcs(arr []*proto.UINPC) {
	for _, npc := range arr {
		db.MergeNpc(npc)
	}
}

func (db *WowDatabase) MergeNpc(src *proto.UINPC) {
	if dst, ok := db.Npcs[src.Id]; ok {
		googleProto.Merge(dst, src)
	} else {
		db.Npcs[src.Id] = src
	}
}

func (db *WowDatabase) MergeConsumable(src *proto.Consumable) {
	if dst, ok := db.Consumables[src.Id]; ok {
		if src.Stats != nil {
			dst.Stats = src.Stats
			src.Stats = nil
		}
		googleProto.Merge(dst, src)
	} else {
		db.Consumables[src.Id] = src
	}
}

func (db *WowDatabase) MergeEffect(src *proto.SpellEffect) {
	if dst, ok := db.Effects[src.Id]; ok {
		googleProto.Merge(dst, src)
	} else {
		db.Effects[src.Id] = src
	}
}

func (db *WowDatabase) AddItemIcon(id int32, icon string, name string) {
	if id == 0 {
		return
	}
	db.ItemIcons[id] = &proto.IconData{Id: id, Name: name, Icon: icon}

}

func (db *WowDatabase) AddSpellIcon(id int32, icon string, name string) {
	if id == 0 {
		return
	}
	db.SpellIcons[id] = &proto.IconData{Id: id, Name: name, Icon: icon}

}

type idKeyed interface {
	GetId() int32
}

func mapToSlice[T idKeyed](m map[int32]T) []T {
	vs := make([]T, 0, len(m))
	for _, v := range m {
		vs = append(vs, v)
	}
	slices.SortFunc(vs, func(a, b T) int {
		return int(a.GetId() - b.GetId())
	})
	return vs
}

type ilvlKeyed interface {
	GetIlvl() int32
}

func mapToSliceByIlvl[T ilvlKeyed](m map[int32]T) []T {
	vs := make([]T, 0, len(m))
	for _, v := range m {
		vs = append(vs, v)
	}
	slices.SortFunc(vs, func(a, b T) int {
		return int(a.GetIlvl() - b.GetIlvl())
	})
	return vs
}

func (db *WowDatabase) ToUIProto() *proto.UIDatabase {
	enchants := make([]*proto.UIEnchant, 0, len(db.Enchants))
	for _, v := range db.Enchants {
		enchants = append(enchants, v)
	}
	slices.SortFunc(enchants, func(v1, v2 *proto.UIEnchant) int {
		if v1.EffectId != v2.EffectId {
			return int(v1.EffectId - v2.EffectId)
		}
		return int(v1.Type - v2.Type)
	})

	return &proto.UIDatabase{
		Items:                    mapToSlice(db.Items),
		RandomSuffixes:           mapToSlice(db.RandomSuffixes),
		Enchants:                 enchants,
		Gems:                     mapToSlice(db.Gems),
		Encounters:               db.Encounters,
		Zones:                    mapToSlice(db.Zones),
		Npcs:                     mapToSlice(db.Npcs),
		ItemIcons:                mapToSlice(db.ItemIcons),
		SpellIcons:               mapToSlice(db.SpellIcons),
		GlyphIds:                 db.GlyphIDs,
		ReforgeStats:             mapToSlice(db.ReforgeStats),
		ItemEffectRandPropPoints: mapToSliceByIlvl(db.ItemEffectRandPropPoints),
		Consumables:              mapToSlice(db.Consumables),
		SpellEffects:             mapToSlice(db.Effects),
	}
}

func sliceToMap[T idKeyed](vs []T) map[int32]T {
	m := make(map[int32]T, len(vs))
	for _, v := range vs {
		m[v.GetId()] = v
	}
	return m
}
func iLvlKeyedSliceToMap[T ilvlKeyed](vs []T) map[int32]T {
	m := make(map[int32]T, len(vs))
	for _, v := range vs {
		m[v.GetIlvl()] = v
	}
	return m
}
func ReadDatabaseFromJson(jsonStr string) *WowDatabase {
	dbProto := &proto.UIDatabase{}
	if err := protojson.Unmarshal([]byte(jsonStr), dbProto); err != nil {
		panic(err)
	}

	enchants := make(map[EnchantDBKey]*proto.UIEnchant, len(dbProto.Enchants))
	for _, v := range dbProto.Enchants {
		enchants[EnchantToDBKey(v)] = v
	}

	return &WowDatabase{
		Items:                    sliceToMap(dbProto.Items),
		RandomSuffixes:           sliceToMap(dbProto.RandomSuffixes),
		Enchants:                 enchants,
		Gems:                     sliceToMap(dbProto.Gems),
		Zones:                    sliceToMap(dbProto.Zones),
		Npcs:                     sliceToMap(dbProto.Npcs),
		ItemIcons:                sliceToMap(dbProto.ItemIcons),
		SpellIcons:               sliceToMap(dbProto.SpellIcons),
		ReforgeStats:             sliceToMap(dbProto.ReforgeStats),
		ItemEffectRandPropPoints: iLvlKeyedSliceToMap(dbProto.ItemEffectRandPropPoints),
		Consumables:              sliceToMap(dbProto.Consumables),
		Effects:                  sliceToMap(dbProto.SpellEffects),
	}
}

func (db *WowDatabase) WriteBinaryAndJson(binFilePath, jsonFilePath string) {
	db.WriteBinary(binFilePath)
	db.WriteJson(jsonFilePath)
}

func (db *WowDatabase) WriteBinary(binFilePath string) {
	uidb := db.ToUIProto()

	// Write database as a binary file.
	protoBytes, err := googleProto.Marshal(uidb)
	if err != nil {
		log.Fatalf("[ERROR] Failed to marshal db: %s", err.Error())
	}
	os.WriteFile(binFilePath, protoBytes, 0666)
}

func (db *WowDatabase) WriteJson(jsonFilePath string) {
	// Also write in JSON format, so we can manually inspect the contents.
	// Write it out line-by-line, so we can have 1 line / item, making it more human-readable.
	uidb := db.ToUIProto()

	buffer := new(bytes.Buffer)
	buffer.WriteString("{\n")

	tools.WriteProtoArrayToBuffer(uidb.Items, buffer, "items")
	buffer.WriteString(",\n")
	tools.WriteProtoArrayToBuffer(uidb.RandomSuffixes, buffer, "randomSuffixes")
	buffer.WriteString(",\n")
	tools.WriteProtoArrayToBuffer(uidb.Enchants, buffer, "enchants")
	buffer.WriteString(",\n")
	tools.WriteProtoArrayToBuffer(uidb.Gems, buffer, "gems")
	buffer.WriteString(",\n")
	tools.WriteProtoArrayToBuffer(uidb.Zones, buffer, "zones")
	buffer.WriteString(",\n")
	tools.WriteProtoArrayToBuffer(uidb.Npcs, buffer, "npcs")
	buffer.WriteString(",\n")
	tools.WriteProtoArrayToBuffer(uidb.ReforgeStats, buffer, "reforgeStats")
	buffer.WriteString(",\n")
	tools.WriteProtoArrayToBuffer(uidb.ItemEffectRandPropPoints, buffer, "itemEffectRandPropPoints")
	buffer.WriteString(",\n")
	tools.WriteProtoArrayToBuffer(uidb.ItemIcons, buffer, "itemIcons")
	buffer.WriteString(",\n")
	tools.WriteProtoArrayToBuffer(uidb.SpellIcons, buffer, "spellIcons")
	buffer.WriteString(",\n")
	tools.WriteProtoArrayToBuffer(uidb.Encounters, buffer, "encounters")
	buffer.WriteString(",\n")
	tools.WriteProtoArrayToBuffer(uidb.GlyphIds, buffer, "glyphIds")
	buffer.WriteString(",\n")
	tools.WriteProtoArrayToBuffer(uidb.Consumables, buffer, "consumables")
	buffer.WriteString(",\n")
	tools.WriteProtoArrayToBuffer(uidb.SpellEffects, buffer, "spellEffects")
	buffer.WriteString("\n")
	buffer.WriteString("}")
	os.WriteFile(jsonFilePath, buffer.Bytes(), 0666)
}
