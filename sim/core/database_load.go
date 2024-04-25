// Only include this file in the build when we specify the 'with_db' tag.
// Without the tag, the database will start out completely empty.
//go:build with_db

package core

import (
	"github.com/wowsims/cata/assets/database"
	"github.com/wowsims/cata/sim/core/proto"
)

func init() {
	db := database.Load()
	WITH_DB = true

	simDB := &proto.SimDatabase{
		Items:          make([]*proto.SimItem, len(db.Items)),
		Enchants:       make([]*proto.SimEnchant, len(db.Enchants)),
		Gems:           make([]*proto.SimGem, len(db.Gems)),
		ReforgeStats:   make([]*proto.ReforgeStat, len(db.ReforgeStats)),
		RandomSuffixes: make([]*proto.ItemRandomSuffix, len(db.RandomSuffixes)),
	}

	for i, item := range db.Items {
		simDB.Items[i] = &proto.SimItem{
			Id:               item.Id,
			Name:             item.Name,
			Type:             item.Type,
			ArmorType:        item.ArmorType,
			WeaponType:       item.WeaponType,
			HandType:         item.HandType,
			RangedWeaponType: item.RangedWeaponType,
			Stats:            item.Stats,
			GemSockets:       item.GemSockets,
			SocketBonus:      item.SocketBonus,
			WeaponDamageMin:  item.WeaponDamageMin,
			WeaponDamageMax:  item.WeaponDamageMax,
			WeaponSpeed:      item.WeaponSpeed,
			SetName:          item.SetName,
			RandPropPoints:   item.RandPropPoints,
		}
	}

	for i, suffix := range db.RandomSuffixes {
		simDB.RandomSuffixes[i] = &proto.ItemRandomSuffix{
			Id:    suffix.Id,
			Name:  suffix.Name,
			Stats: suffix.Stats,
		}
	}

	for i, enchant := range db.Enchants {
		simDB.Enchants[i] = &proto.SimEnchant{
			EffectId: enchant.EffectId,
			Stats:    enchant.Stats,
		}
	}

	for i, gem := range db.Gems {
		simDB.Gems[i] = &proto.SimGem{
			Id:    gem.Id,
			Name:  gem.Name,
			Color: gem.Color,
			Stats: gem.Stats,
		}
	}

	for i, reforgeStat := range db.ReforgeStats {
		simDB.ReforgeStats[i] = &proto.ReforgeStat{
			Id:         reforgeStat.Id,
			FromStat:   reforgeStat.FromStat,
			ToStat:     reforgeStat.ToStat,
			Multiplier: reforgeStat.Multiplier,
		}
	}

	addToDatabase(simDB)
}
