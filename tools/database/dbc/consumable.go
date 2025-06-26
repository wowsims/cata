package dbc

import (
	"slices"

	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

type Consumable struct {
	Id                 int             // Item ID
	Name               string          // Item name
	ItemLevel          int             // Item level
	RequiredLevel      int             // Required level to use
	ClassId            int             // Item class ID (should be 0 for consumables)
	SubClassId         ConsumableClass // Item subclass ID
	IconFileDataID     int             // Icon file data ID
	SpellCategoryID    int             // Spell category ID
	SpellCategoryFlags int             // Spell category flags
	ItemEffects        []int           // Item effect IDs
	ElixirType         int
	Duration           int // In milliseconds
	CooldownDuration   int // In milliseconds
}

func (c *Consumable) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"Id":                 c.Id,
		"Name":               c.Name,
		"ItemLevel":          c.ItemLevel,
		"RequiredLevel":      c.RequiredLevel,
		"ClassId":            c.ClassId,
		"SubClassId":         c.SubClassId,
		"IconFileDataID":     c.IconFileDataID,
		"SpellCategoryID":    c.SpellCategoryID,
		"SpellCategoryFlags": c.SpellCategoryFlags,
		"ItemEffects":        c.ItemEffects,
	}
}

// ToProto converts the Consumable to a proto representation.
func (c *Consumable) ToProto() *proto.Consumable {
	return &proto.Consumable{
		Id:               int32(c.Id),
		Type:             c.GetConsumableType(),
		Stats:            c.GetStatModifiers().ToProtoArray(),
		Name:             c.Name,
		BuffsMainStat:    false, // Todo: Should be food currently, might be more in MoP, figure out how to tell
		BuffDuration:     int32(c.Duration / 1000),
		CooldownDuration: int32(c.CooldownDuration / 1000),
		EffectIds:        c.GetNonStatEffectIds(),
	}
}
func (c *Consumable) GetConsumableType() proto.ConsumableType {
	if c.SubClassId == ELIXIR {
		switch c.ElixirType {
		case 1:
			return proto.ConsumableType_ConsumableTypeGuardianElixir
		case 2:
			return proto.ConsumableType_ConsumableTypeBattleElixir
		}
	}
	if val, ok := consumableClassToProto[c.SubClassId]; ok {
		return val
	}
	return proto.ConsumableType_ConsumableTypeUnknown
}

func (s ConsumableClass) ToProto() proto.ConsumableType {
	if val, ok := consumableClassToProto[s]; ok {
		return val
	}
	return proto.ConsumableType_ConsumableTypeUnknown
}

func (consumable *Consumable) GetNonStatEffectIds() []int32 {
	var effectIds []int32

	statAuraTypes := map[SpellEffectType]bool{
		E_HEAL:     true,
		E_ENERGIZE: true,
	}
	slices.Sort(consumable.ItemEffects)
	for _, effectID := range consumable.ItemEffects {
		effect := GetItemEffect(effectID)
		if effect.ID != 0 {
			if spellEffects, ok := dbcInstance.SpellEffects[effect.SpellID]; ok {
				for _, spellEffect := range spellEffects {
					if statAuraTypes[spellEffect.EffectType] {
						effectIds = append(effectIds, int32(spellEffect.ID))
					}
				}
			}
		}
	}
	slices.Sort(effectIds)
	return effectIds
}
func (consumable *Consumable) GetStatModifiers() *stats.Stats {
	stats := &stats.Stats{}
	for _, effectID := range consumable.ItemEffects {
		effect := GetItemEffect(effectID)
		if effect.ID != 0 {
			if spellEffects, ok := dbcInstance.SpellEffects[effect.SpellID]; ok {
				for _, spellEffect := range spellEffects {
					stat := spellEffect.ParseStatEffect(spellEffect.Coefficient != 0, 0)
					stats.AddInplace(stat)
				}
			}
		}
	}
	return stats
}
