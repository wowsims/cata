package dbc

import (
	"compress/gzip"
	"io"
	"os"

	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

func SocketCost(level int) float64 {
	cost := 0.0

	if level >= 19 {
		cost++
	}
	if level >= 31 {
		cost++
	}
	if level >= 43 {
		cost++
	}
	if level >= 55 {
		cost++
	}
	if level >= 89 {
		cost++
	}
	if level >= 100 {
		cost++
	}
	if level >= 178 {
		cost += 2
	}
	return cost
}

func GetProfession(id int) proto.Profession {
	if profession, ok := MapProfessionIdToProfession[id]; ok {
		return profession
	}
	return 0
}

func GetClassesFromClassMask(mask int) []proto.Class {
	var result []proto.Class
	for _, class := range classes {
		// Calculate the bit flag using 1 << (ID - 1)
		if mask&(1<<(class.ID-1)) != 0 {
			result = append(result, class.protoClass)
		}
	}
	return result
}

type DbcClass struct {
	protoClass proto.Class
	ID         int
}

var classes = []DbcClass{
	{proto.Class_ClassWarrior, 1},
	{proto.Class_ClassPaladin, 2},
	{proto.Class_ClassHunter, 3},
	{proto.Class_ClassRogue, 4},
	{proto.Class_ClassPriest, 5},
	{proto.Class_ClassDeathKnight, 6},
	{proto.Class_ClassShaman, 7},
	{proto.Class_ClassMage, 8},
	{proto.Class_ClassWarlock, 9},
	{proto.Class_ClassDruid, 11},
}

func readGzipFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, DataLoadError{
			Source:   filename,
			DataType: "gzip file",
			Reason:   err.Error(),
		}
	}
	defer f.Close()

	gzReader, err := gzip.NewReader(f)
	if err != nil {
		return nil, DataLoadError{
			Source:   filename,
			DataType: "gzip",
			Reason:   err.Error(),
		}
	}
	defer gzReader.Close()

	data, err := io.ReadAll(gzReader)
	if err != nil {
		return nil, DataLoadError{
			Source:   filename,
			DataType: "decompression",
			Reason:   err.Error(),
		}
	}

	return data, nil
}

func processEnchantmentEffects(
	effects []int,
	effectArgs []int,
	effectPoints []int,
	outStats *stats.Stats,
	addRanged bool,
) {
	for i, effect := range effects {
		switch effect {
		case ITEM_ENCHANTMENT_RESISTANCE:
			stat, _ := MapResistanceToStat(effectArgs[i])
			outStats[stat] = float64(effectPoints[i])
		case ITEM_ENCHANTMENT_STAT:
			stat, _ := MapBonusStatIndexToStat(effectArgs[i])
			outStats[stat] = float64(effectPoints[i])
			// If the bonus stat is attack power, copy it to ranged attack power
			if addRanged && stat == proto.Stat_StatAttackPower {
				outStats[proto.Stat_StatRangedAttackPower] = float64(effectPoints[i])
			}
		case ITEM_ENCHANTMENT_EQUIP_SPELL:
			spellEffects := dbcInstance.SpellEffects[effectArgs[i]]
			for _, spellEffect := range spellEffects {
				if spellEffect.EffectMiscValues[0] == -1 &&
					spellEffect.EffectType == E_APPLY_AURA &&
					spellEffect.EffectAura == A_MOD_STAT {
					// Apply bonus to all stats
					outStats[proto.Stat_StatAgility] += float64(spellEffect.EffectBasePoints)
					outStats[proto.Stat_StatIntellect] += float64(spellEffect.EffectBasePoints)
					outStats[proto.Stat_StatSpirit] += float64(spellEffect.EffectBasePoints)
					outStats[proto.Stat_StatStamina] += float64(spellEffect.EffectBasePoints)
					outStats[proto.Stat_StatStrength] += float64(spellEffect.EffectBasePoints)
				}
			}
		case ITEM_ENCHANTMENT_COMBAT_SPELL:
			// Not processed (chance on hit, ignore for now)
		case ITEM_ENCHANTMENT_USE_SPELL:
			// Not processed
		}
	}
}
