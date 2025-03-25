package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/wowsims/cata/sim/core/proto"
)

func QueryItems() ([]*proto.UIItem, error) {
	helper, err := NewDBHelper("./tools/database/wowsims.db")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database helper: %v", err)
	}
	defer helper.Close()

	query := `
		SELECT 
			i.ID, 
			s.Display_lang AS Name,
			i.InventoryType,
			s.ItemDelay,
			s.OverallQualityID,
			s.DmgVariance,
			s.MinDamage, 
			s.MaxDamage,
			s.ItemLevel,
			ic.ClassName_lang AS ItemClassName,
			isc.VerboseName_lang AS ItemSubClassName,
			rpp.Epic as RPPEpic,
			rpp.Superior as RPPSuperior,
			rpp.Good as RPPGood,
			s.Field_1_15_3_55112_014 as StatValue,
			s.StatModifier_bonusStat as bonusStat,
			(at.Cloth * al.Clothmodifier) AS clothArmorValue,
			(at.Leather * al.LeatherModifier) AS leatherArmorValue,
			(at.Mail * al.Chainmodifier) AS mailArmorValue,
			(at.Plate * al.Platemodifier) AS plateArmorValue,
			CASE 
				WHEN s.InventoryType = 20 THEN 5 
				ELSE s.InventoryType 
			END AS ArmorLocationID,
			ias.Quality as shieldArmorValues,
			s.StatPercentEditor as StatPercentEditor,
			s.SocketType as SocketTypes,
			s.Socket_match_enchantment_ID as SocketEnchantmentId,
			s.Flags_0 as Flags_0
		FROM Item i
		JOIN ItemSparse s ON i.ID = s.ID
		JOIN ItemClass ic ON i.ClassID = ic.ClassID
		JOIN ItemSubClass isc ON i.ClassID = isc.ClassID AND i.SubClassID = isc.SubClassID
		JOIN RandPropPoints rpp ON s.ItemLevel = rpp.ID
		LEFT JOIN ArmorLocation al ON al.ID = ArmorLocationId
		LEFT JOIN ItemArmorShield ias ON s.ItemLevel = ias.ItemLevel
		JOIN ItemArmorTotal at ON s.ItemLevel = at.ItemLevel
		WHERE s.ID = ?;
	`

	var items []*proto.UIItem

	helper.QueryAndProcess(query, func(rows *sql.Rows) error {
		for rows.Next() {
			item, err := processItemRow(helper, rows)
			if err != nil {
				log.Printf("Error processing item row: %v", err)
				continue
			}
			items = append(items, item)
			jsonItem, err := json.MarshalIndent(item, "", "  ")

			if err != nil {
				log.Printf("Error marshaling item (ID: %d): %v", item.Id, err)
			} else {
				fmt.Println("Marshalled UIItem:")
				fmt.Println(string(jsonItem))
			}
		}
		return rows.Err()
	}, 78489) // Pass the item ID or remove this and WHERE s.ID = ?; if you want all items

	return items, nil
}
