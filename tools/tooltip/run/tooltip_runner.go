package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/wowsims/mop/tools/database/dbc"
	"github.com/wowsims/mop/tools/tooltip"
)

// Sample program for testing
func main() {
	fmt.Println("Loading DB...")
	dbc := dbc.GetDBC()
	fmt.Printf("Loaded %d spell tooltips.\n", len(dbc.Spells))

	start := time.Now()
	errorCounter := 0
	successCounter := 0
	keys := make([]int, 0, len(dbc.Spells))

	// used for specific Tooltip eval
	const fixedID = 0

	// run spells in order
	for key := range dbc.Spells {
		if fixedID > 0 && key != fixedID {
			continue
		}

		keys = append(keys, key)
	}

	sort.Ints(keys)
	for _, spellId := range keys {
		spell := dbc.Spells[spellId]
		parsedTooltip, error := tooltip.ParseTooltip(spell.Description, tooltip.DBCTooltipDataProvider{DBC: dbc}, int64(spell.ID))
		if error != nil {
			fmt.Printf("[%d]: Failed to render. %s\n", spell.ID, error.Error())
			errorCounter++
			if errorCounter > 20 {
				panic("Too many errors")
			}
		} else {
			successCounter++
			fmt.Printf("[%d]: %s\n", spell.ID, parsedTooltip)
		}
	}

	fmt.Printf("Prased %d Tooltips. In %.1fs Failed: %d", successCounter, float64(time.Since(start))/float64(time.Second), errorCounter)
}
