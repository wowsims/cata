package dbc

import (
	"slices"
	"sort"

	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

type Enchant struct {
	EffectId           int
	Name               string
	SpellId            int
	ItemId             int
	ProfessionId       int
	Effects            []int
	EffectPoints       []int
	EffectArgs         []int
	IsWeaponEnchant    bool
	InventoryType      InventoryTypeFlag
	SubClassMask       int
	ClassMask          int
	FDID               int
	Quality            ItemQuality
	RequiredProfession int
	EffectName         string
}

func (enchant *Enchant) HasEnchantEffect() bool {
	for idx, effect := range enchant.Effects {
		if effect == ITEM_ENCHANTMENT_COMBAT_SPELL {
			return true
		}

		// We apply a buff here, check if it's a trigger
		if effect == ITEM_ENCHANTMENT_EQUIP_SPELL {
			spellId := enchant.EffectArgs[idx]
			spellEffects := dbcInstance.SpellEffects[spellId]
			for _, spellEffect := range spellEffects {
				if spellEffect.EffectAura == A_PROC_TRIGGER_SPELL ||
					spellEffect.EffectAura == A_PROC_TRIGGER_SPELL_WITH_VALUE {
					return true
				}
			}
		}
	}

	return false
}

func (enchant *Enchant) ToProto() *proto.UIEnchant {
	uiEnchant := &proto.UIEnchant{
		Name:               enchant.Name,
		ItemId:             int32(enchant.ItemId),
		SpellId:            int32(enchant.SpellId),
		EffectId:           int32(enchant.EffectId),
		ClassAllowlist:     GetClassesFromClassMask(enchant.ClassMask),
		ExtraTypes:         []proto.ItemType{},
		Stats:              stats.Stats{}.ToProtoArray(),
		Quality:            enchant.Quality.ToProto(),
		RequiredProfession: GetProfession(enchant.RequiredProfession),
	}

	if enchant.HasEnchantEffect() {
		eff := ItemEffect{TriggerType: 2, SpellID: enchant.SpellId}
		parsedEffect, hasStats := eff.ToProto(0, 0)
		if hasStats {
			uiEnchant.EnchantEffect = parsedEffect
		}
		if uiEnchant.EnchantEffect.GetOnUse() == nil && uiEnchant.EnchantEffect.GetProc() == nil {
			uiEnchant.EnchantEffect = nil
		}
	}

	if enchant.FDID == 0 {
		uiEnchant.Icon = "trade_engraving"
	}

	if enchant.IsWeaponEnchant {
		// Process weapon enchants.
		uiEnchant.Type = proto.ItemType_ItemTypeWeapon
		if enchant.SubClassMask == ITEM_SUBCLASS_BIT_WEAPON_STAFF {
			// Staff only.
			uiEnchant.EnchantType = proto.EnchantType_EnchantTypeStaff
		}
		if enchant.SubClassMask == rangedMask {
			uiEnchant.Type = proto.ItemType_ItemTypeRanged
		}
		if enchant.SubClassMask == twoHandMask {
			// Two-handed weapon.
			uiEnchant.EnchantType = proto.EnchantType_EnchantTypeTwoHand
		}
	} else {
		// Process non-weapon enchants.
		if enchant.SubClassMask == OffHandValue {
			uiEnchant.EnchantType = proto.EnchantType_EnchantTypeOffHand
			uiEnchant.Type = proto.ItemType_ItemTypeWeapon
		}
		if enchant.SubClassMask == ITEM_SUBCLASS_BIT_ARMOR_SHIELD || enchant.SubClassMask == 64 { // idk where the 64 comes from but shield spikes are this
			uiEnchant.EnchantType = proto.EnchantType_EnchantTypeShield
			uiEnchant.Type = proto.ItemType_ItemTypeWeapon
		}
		// Sort flags for consistent generation
		var flags []int
		for flag := range MapInventoryTypeToEnchantMetaType {
			flags = append(flags, int(flag))
		}

		sort.Ints(flags)

		for _, f := range flags {
			flag := InventoryTypeFlag(f)
			m := MapInventoryTypeToEnchantMetaType[flag]
			if enchant.InventoryType.Has(flag) {
				if uiEnchant.Type != proto.ItemType_ItemTypeUnknown {
					uiEnchant.ExtraTypes = append(uiEnchant.ExtraTypes, m.ItemType)
				} else {
					uiEnchant.Type = m.ItemType
				}
			}
		}
		slices.Sort(uiEnchant.ExtraTypes)
	}
	stats := stats.Stats{}
	processEnchantmentEffects(enchant.Effects, enchant.EffectArgs, enchant.EffectPoints, &stats, true)
	uiEnchant.Stats = stats.ToProtoArray()
	return uiEnchant
}
