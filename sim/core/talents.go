package core

import (
	"encoding/json"
	"log"
	"slices"
	"strconv"
	"strings"
)

type TalentTree struct {
	name     string
	spellIds []int32
	isKnown  []bool
}

type Talents struct {
	trees []TalentTree
}

type TalentConfig struct {
	FieldName string `json:"fieldName"`
	// Spell ID for each rank of this talent.
	// Omitted ranks will be inferred by incrementing from the last provided rank.
	SpellIds  []int32 `json:"spellIds"`
	MaxPoints int32   `json:"maxPoints"`
}

type TalentTreeConfig struct {
	Name          string         `json:"name"`
	BackgroundUrl string         `json:"backgroundUrl"`
	Talents       []TalentConfig `json:"talents"`
}

func (talents Talents) IsKnown(spellId int32) bool {
	for _, tree := range talents.trees {
		index := slices.Index(tree.spellIds, spellId)
		if index > -1 {
			return tree.isKnown[index]
		}
	}
	return false
}

func (character *Character) FillTalentsData(talentsJson string, talentsString string) {
	talentsData := Talents{
		trees: make([]TalentTree, 0),
	}

	var talents []TalentTreeConfig

	err := json.Unmarshal([]byte(talentsJson), &talents)
	if err != nil {
		log.Fatalf("failed to parse talent to json %s", err)
	}

	knownTalentTrees := strings.Split(talentsString, "-")

	for treeIdx, tree := range talents {
		talentsData.trees = append(talentsData.trees, TalentTree{
			name:     tree.Name,
			spellIds: make([]int32, 0),
			isKnown:  make([]bool, 0),
		})

		knownLenght := len(knownTalentTrees[treeIdx])
		knownValues := strings.Split(knownTalentTrees[treeIdx], "")

		for talentIdx, talent := range tree.Talents {
			talentsData.trees[treeIdx].spellIds = append(talentsData.trees[treeIdx].spellIds, talent.SpellIds...)

			knownValue := 0
			if knownLenght > talentIdx {
				var _ error
				knownValue, _ = strconv.Atoi(knownValues[talentIdx])
			}
			for knownIdx := 0; knownIdx < int(talent.MaxPoints); knownIdx++ {
				talentsData.trees[treeIdx].isKnown = append(talentsData.trees[treeIdx].isKnown, knownValue > knownIdx)
			}

			// Infer omitted spell IDs.
			if len(talent.SpellIds) < int(talent.MaxPoints) {
				curSpellId := talent.SpellIds[len(talent.SpellIds)-1]
				for i := len(talent.SpellIds); i < int(talent.MaxPoints); i++ {
					curSpellId++
					talentsData.trees[treeIdx].spellIds = append(talentsData.trees[treeIdx].spellIds, curSpellId)
				}
			}
		}
	}

	character.talents = &talentsData
}
