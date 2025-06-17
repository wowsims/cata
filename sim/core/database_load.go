//go:build with_db

// Only include this file in the build when we specify the 'with_db' tag.
// Without the tag, the database will start out completely empty.
package core

import (
	"github.com/wowsims/mop/assets/database"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	db := database.Load()
	WITH_DB = true

	simDB := &proto.SimDatabase{
		Items:                    make([]*proto.SimItem, len(db.Items)),
		Enchants:                 make([]*proto.SimEnchant, len(db.Enchants)),
		Gems:                     make([]*proto.SimGem, len(db.Gems)),
		ReforgeStats:             make([]*proto.ReforgeStat, len(db.ReforgeStats)),
		ItemEffectRandPropPoints: make([]*proto.ItemEffectRandPropPoints, len(db.ItemEffectRandPropPoints)),
		RandomSuffixes:           make([]*proto.ItemRandomSuffix, len(db.RandomSuffixes)),
		Consumables:              make([]*proto.Consumable, len(db.Consumables)),
		SpellEffects:             make([]*proto.SpellEffect, len(db.SpellEffects)),
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
			GemSockets:       item.GemSockets,
			SocketBonus:      item.SocketBonus,
			WeaponSpeed:      item.WeaponSpeed,
			SetName:          item.SetName,
			SetId:            item.SetId,
			ScalingOptions:   item.ScalingOptions,
			ItemEffect:       item.ItemEffect,
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
			EffectId:      enchant.EffectId,
			Stats:         enchant.Stats,
			EnchantEffect: enchant.EnchantEffect,
			Name:          enchant.Name,
			Type:          enchant.Type,
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

	for i, itemEffectRpp := range db.ItemEffectRandPropPoints {
		simDB.ItemEffectRandPropPoints[i] = &proto.ItemEffectRandPropPoints{
			Ilvl:           itemEffectRpp.Ilvl,
			RandPropPoints: itemEffectRpp.RandPropPoints,
		}
	}

	for i, consumable := range db.Consumables {
		simDB.Consumables[i] = consumable
	}

	for i, effect := range db.SpellEffects {
		simDB.SpellEffects[i] = effect
	}

	addToDatabase(simDB)
}
