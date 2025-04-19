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
		Items:            make([]*proto.SimItem, len(db.Items)),
		Enchants:         make([]*proto.SimEnchant, len(db.Enchants)),
		Gems:             make([]*proto.SimGem, len(db.Gems)),
		ReforgeStats:     make([]*proto.ReforgeStat, len(db.ReforgeStats)),
		RandomSuffixes:   make([]*proto.ItemRandomSuffix, len(db.RandomSuffixes)),
		Consumables:      make([]*proto.Consumable, len(db.Consumables)),
		SpellEffects:     make([]*proto.SpellEffect, len(db.SpellEffects)),
		RandomPropPoints: make(map[int32]*proto.QualityAllocations, len(db.RandomPropPoints)),
		ArmorTotalValue:  make(map[int32]*proto.ItemArmorTotal, len(db.ArmorTotalValue)),
		ArmorDb:          &proto.ArmorValueDatabase{},
		WeaponDamageDb:   &proto.WeaponDamageDatabase{},
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
			SetId:            item.SetId,
			Ilvl:             item.Ilvl,
			Quality:          item.Quality,
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

	for i, consumable := range db.Consumables {
		simDB.Consumables[i] = consumable
	}

	for i, effect := range db.SpellEffects {
		simDB.SpellEffects[i] = effect
	}
	for i, v := range db.RandomPropPoints {
		simDB.RandomPropPoints[i] = v
	}
	for _, v := range db.ArmorTotalValue {
		simDB.ArmorTotalValue[v.ItemLevel] = v
	}
	simDB.ArmorDb = db.ArmorDb
	simDB.WeaponDamageDb = db.WeaponDamageDb
	addToDatabase(simDB)
}
