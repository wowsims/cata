package database

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/tailscale/hujson"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

// Example db input file: https://nether.wowhead.com/cata/data/gear-planner?dv=100

func ParseWowheadDB(dbContents string) WowheadDatabase {
	var wowheadDB WowheadDatabase

	// Each part looks like 'WH.setPageData("wow.gearPlanner.some.name", {......});'
	parts := strings.Split(dbContents, "WH.setPageData(")

	for _, dbPart := range parts {
		//fmt.Printf("Part len: %d\n", len(dbPart))
		if len(dbPart) < 10 {
			continue
		}
		dbPart = strings.TrimSpace(dbPart)
		dbPart = strings.TrimRight(dbPart, ");")

		if dbPart[0] != '"' {
			continue
		}
		secondQuoteIdx := strings.Index(dbPart[1:], "\"")
		if secondQuoteIdx == -1 {
			continue
		}
		dbName := dbPart[1 : secondQuoteIdx+1]
		//fmt.Printf("DB name: %s\n", dbName)

		commaIdx := strings.Index(dbPart, ",")
		dbContents := dbPart[commaIdx+1:]
		if dbName == "wow.gearPlanner.cata.item" {
			standardized, err := hujson.Standardize([]byte(dbContents)) // Removes invalid JSON, such as trailing commas
			if err != nil {
				log.Fatalf("Failed to standardize json %s\n\n%s\n\n%s", err, dbContents[0:30], dbContents[len(dbContents)-30:])
			}

			err = json.Unmarshal(standardized, &wowheadDB.Items)
			if err != nil {
				log.Fatalf("failed to parse wowhead item db to json %s\n\n%s", err, dbContents[0:30])
			}
		}

		if dbName == "wow.gearPlanner.cata.randomEnchant" {
			standardized, err := hujson.Standardize([]byte(dbContents)) // Removes invalid JSON, such as trailing commas
			if err != nil {
				log.Fatalf("Failed to standardize json %s\n\n%s\n\n%s", err, dbContents[0:30], dbContents[len(dbContents)-30:])
			}

			err = json.Unmarshal(standardized, &wowheadDB.RandomSuffixes)
			if err != nil {
				log.Fatalf("failed to parse wowhead random suffix db to json %s\n\n%s", err, dbContents[0:30])
			}
		}
	}

	fmt.Printf("\n--\nWowhead DB items loaded: %d\n--\n", len(wowheadDB.Items))
	fmt.Printf("\n--\nWowhead DB random suffixes loaded: %d\n--\n", len(wowheadDB.RandomSuffixes))

	return wowheadDB
}

type WowheadDatabase struct {
	Items          map[string]WowheadItem
	RandomSuffixes map[string]WowheadRandomSuffix
}

type WowheadRandomSuffix struct {
	ID    int32                    `json:"id"`
	Name  string                   `json:"name"`
	Stats WowheadRandomSuffixStats `json:"stats"`
}

type WowheadRandomSuffixStats struct {
	Strength          int32 `json:"str"`
	Agility           int32 `json:"agi"`
	Stamina           int32 `json:"sta"`
	Intellect         int32 `json:"int"`
	Spirit            int32 `json:"spi"`
	SpellPower        int32 `json:"spldmg"`
	MP5               int32 `json:"manargn"`
	HitRating         int32 `json:"hitrtng"`
	CritRating        int32 `json:"critstrkrtng"`
	HasteRating       int32 `json:"hastertng"`
	AttackPower       int32 `json:"mleatkpwr"`
	Expertise         int32 `json:"exprtng"`
	Armor             int32 `json:"armor"`
	RangedAttackPower int32 `json:"rgdatkpwr"`
	Block             int32 `json:"blockrtng"`
	Dodge             int32 `json:"dodgertng"`
	Parry             int32 `json:"parryrtng"`
	ArcaneResistance  int32 `json:"arcres"`
	FireResistance    int32 `json:"firres"`
	FrostResistance   int32 `json:"frores"`
	NatureResistance  int32 `json:"natres"`
	ShadowResistance  int32 `json:"shares"`
	Mastery           int32 `json:"mastrtng"`
}

func (wrs WowheadRandomSuffix) ToProto() *proto.ItemRandomSuffix {
	suffixStats := stats.Stats{
		stats.Strength:          float64(wrs.Stats.Strength),
		stats.Agility:           float64(wrs.Stats.Agility),
		stats.Stamina:           float64(wrs.Stats.Stamina),
		stats.Intellect:         float64(wrs.Stats.Intellect),
		stats.Spirit:            float64(wrs.Stats.Spirit),
		stats.SpellPower:        float64(wrs.Stats.SpellPower),
		stats.MP5:               float64(wrs.Stats.MP5),
		stats.HitRating:         float64(wrs.Stats.HitRating),
		stats.CritRating:        float64(wrs.Stats.CritRating),
		stats.HasteRating:       float64(wrs.Stats.HasteRating),
		stats.AttackPower:       float64(wrs.Stats.AttackPower),
		stats.ExpertiseRating:   float64(wrs.Stats.Expertise),
		stats.Armor:             float64(wrs.Stats.Armor),
		stats.RangedAttackPower: float64(wrs.Stats.RangedAttackPower),
		stats.DodgeRating:       float64(wrs.Stats.Dodge),
		stats.ParryRating:       float64(wrs.Stats.Parry),
		stats.ArcaneResistance:  float64(wrs.Stats.ArcaneResistance),
		stats.FireResistance:    float64(wrs.Stats.FireResistance),
		stats.FrostResistance:   float64(wrs.Stats.FrostResistance),
		stats.NatureResistance:  float64(wrs.Stats.NatureResistance),
		stats.ShadowResistance:  float64(wrs.Stats.ShadowResistance),
		stats.MasteryRating:     float64(wrs.Stats.Mastery),
	}

	return &proto.ItemRandomSuffix{
		Id:    wrs.ID,
		Name:  wrs.Name,
		Stats: suffixStats.ToProtoArray(),
	}
}

type WowheadItem struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`

	Quality int32 `json:"quality"`
	Ilvl    int32 `json:"itemLevel"`
	Phase   int32 `json:"contentPhase"`

	RaceMask  uint32 `json:"raceMask"`
	ClassMask uint16 `json:"classMask"`

	Stats               WowheadItemStats `json:"stats"`
	RandomSuffixOptions []int32          `json:"randomEnchants"`

	SourceTypes   []int32             `json:"source"` // 1 = Crafted, 2 = Dropped by, 3 = sold by zone vendor? barely used, 4 = Quest, 5 = Sold by
	SourceDetails []WowheadItemSource `json:"sourcemore"`
}
type WowheadItemStats struct {
	Armor int32 `json:"armor"`
}
type WowheadItemSource struct {
	C        int32  `json:"c"`
	S        int32  `json:"s"`
	Name     string `json:"n"`    // Name of crafting spell
	Icon     string `json:"icon"` // Icon corresponding to the named entity
	EntityID int32  `json:"ti"`   // Crafting Spell ID / NPC ID / ?? / Quest ID
	ZoneID   int32  `json:"z"`    // Only for drop / sold by sources
}

func (wi WowheadItem) ToProto() *proto.UIItem {
	var sources []*proto.UIItemSource
	for i, details := range wi.SourceDetails {
		switch wi.SourceTypes[i] {
		case 1: // Crafted
			profession, ok := wowheadProfessionIds[details.S]
			if ok {
				sources = append(sources, &proto.UIItemSource{
					Source: &proto.UIItemSource_Crafted{
						Crafted: &proto.CraftedSource{
							SpellId:    details.EntityID,
							Profession: profession,
						},
					},
				})
			}
		case 2: // Dropped by
			// Do nothing, we'll get this from AtlasLoot.
		case 3: // Sold by zone vendor? barely used
		case 4: // Quest
			sources = append(sources, &proto.UIItemSource{
				Source: &proto.UIItemSource_Quest{
					Quest: &proto.QuestSource{
						Id:   details.EntityID,
						Name: details.Name,
					},
				},
			})
		case 5: // Sold by
			sources = append(sources, &proto.UIItemSource{
				Source: &proto.UIItemSource_SoldBy{
					SoldBy: &proto.SoldBySource{
						NpcId:   details.EntityID,
						NpcName: details.Name,
						ZoneId:  details.ZoneID,
					},
				},
			})
		}
	}

	return &proto.UIItem{
		Id:                  wi.ID,
		Name:                wi.Name,
		Icon:                wi.Icon,
		Ilvl:                wi.Ilvl,
		Phase:               wi.Phase,
		FactionRestriction:  wi.getFactionRstriction(),
		ClassAllowlist:      wi.getClassRestriction(),
		Sources:             sources,
		RandomSuffixOptions: wi.RandomSuffixOptions,
	}
}

var wowheadProfessionIds = map[int32]proto.Profession{
	//"FirstAid": proto.Profession_FirstAid,
	//"Cooking":        proto.Profession_Cooking,
	164: proto.Profession_Blacksmithing,
	165: proto.Profession_Leatherworking,
	171: proto.Profession_Alchemy,
	186: proto.Profession_Mining,
	197: proto.Profession_Tailoring,
	202: proto.Profession_Engineering,
	333: proto.Profession_Enchanting,
	755: proto.Profession_Jewelcrafting,
	773: proto.Profession_Inscription,
	794: proto.Profession_Archeology,
}

func (wi WowheadItem) getFactionRstriction() proto.UIItem_FactionRestriction {
	if wi.RaceMask == 2098253 {
		return proto.UIItem_FACTION_RESTRICTION_ALLIANCE_ONLY
	} else if wi.RaceMask == 946 {
		return proto.UIItem_FACTION_RESTRICTION_HORDE_ONLY
	} else {
		return proto.UIItem_FACTION_RESTRICTION_UNSPECIFIED
	}
}

type ClassMask uint16

const (
	ClassMaskWarrior     ClassMask = 1 << iota
	ClassMaskPaladin               // 2
	ClassMaskHunter                // 4
	ClassMaskRogue                 // 8
	ClassMaskPriest                // 16
	ClassMaskDeathKnight           // 32
	ClassMaskShaman                // 64
	ClassMaskMage                  // 128
	ClassMaskWarlock               // 256
	ClassMaskUnknown               // 512 seemingly unused?
	ClassMaskDruid                 // 1024
)

func (wi WowheadItem) getClassRestriction() []proto.Class {
	classAllowlist := []proto.Class{}
	if wi.ClassMask&uint16(ClassMaskWarrior) != 0 {
		classAllowlist = append(classAllowlist, proto.Class_ClassWarrior)
	}
	if wi.ClassMask&uint16(ClassMaskPaladin) != 0 {
		classAllowlist = append(classAllowlist, proto.Class_ClassPaladin)
	}
	if wi.ClassMask&uint16(ClassMaskHunter) != 0 {
		classAllowlist = append(classAllowlist, proto.Class_ClassHunter)
	}
	if wi.ClassMask&uint16(ClassMaskRogue) != 0 {
		classAllowlist = append(classAllowlist, proto.Class_ClassRogue)
	}
	if wi.ClassMask&uint16(ClassMaskPriest) != 0 {
		classAllowlist = append(classAllowlist, proto.Class_ClassPriest)
	}
	if wi.ClassMask&uint16(ClassMaskDruid) != 0 {
		classAllowlist = append(classAllowlist, proto.Class_ClassDruid)
	}
	if wi.ClassMask&uint16(ClassMaskShaman) != 0 {
		classAllowlist = append(classAllowlist, proto.Class_ClassShaman)
	}
	if wi.ClassMask&uint16(ClassMaskMage) != 0 {
		classAllowlist = append(classAllowlist, proto.Class_ClassMage)
	}
	if wi.ClassMask&uint16(ClassMaskWarlock) != 0 {
		classAllowlist = append(classAllowlist, proto.Class_ClassWarlock)
	}
	if wi.ClassMask&uint16(ClassMaskDeathKnight) != 0 {
		classAllowlist = append(classAllowlist, proto.Class_ClassDeathKnight)
	}

	return classAllowlist
}
