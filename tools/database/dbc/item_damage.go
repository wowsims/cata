package dbc

type ItemDamageTable struct {
	ItemLevel int
	Quality   []float64
}

func (item *Item) WeaponDps(itemLevel int) float64 {
	quality := item.OverallQuality
	if item.OverallQuality > 6 {
		quality = 4 // Heirlooms = epic
	}

	ilvl := 0
	if itemLevel > 0 {
		ilvl = itemLevel
	} else {
		ilvl = item.ItemLevel
	}

	switch item.InventoryType {
	case INVTYPE_WEAPON, INVTYPE_WEAPONMAINHAND, INVTYPE_WEAPONOFFHAND:
		{
			if item.Flags1.Has(CASTER_WEAPON) {
				return dbcInstance.ItemDamageTable["ItemDamageOneHandCaster"][ilvl].Quality[quality]
			} else {
				return dbcInstance.ItemDamageTable["ItemDamageOneHand"][ilvl].Quality[quality]
			}
		}
	case INVTYPE_2HWEAPON:
		if item.Flags1.Has(CASTER_WEAPON) {
			return dbcInstance.ItemDamageTable["ItemDamageTwoHandCaster"][ilvl].Quality[quality]
		} else {
			return dbcInstance.ItemDamageTable["ItemDamageTwoHand"][ilvl].Quality[quality]
		}
	case INVTYPE_RANGED, INVTYPE_THROWN, INVTYPE_RANGEDRIGHT:
		switch item.ItemSubClass {
		case ITEM_SUBCLASS_WEAPON_BOW, ITEM_SUBCLASS_WEAPON_GUN, ITEM_SUBCLASS_WEAPON_CROSSBOW:
			return dbcInstance.ItemDamageTable["ItemDamageRanged"][ilvl].Quality[quality]
		case ITEM_SUBCLASS_WEAPON_THROWN:
			return dbcInstance.ItemDamageTable["ItemDamageThrown"][ilvl].Quality[quality]
		case ITEM_SUBCLASS_WEAPON_WAND:
			return dbcInstance.ItemDamageTable["ItemDamageWand"][ilvl].Quality[quality]
		}
	}
	return 0
}
