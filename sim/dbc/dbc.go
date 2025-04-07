package dbc

import (
	"encoding/json"
	"log"
	"sync"
)

type DBC struct {
	Items                  map[int]*Item                      // Item ID
	Gems                   map[int]*Gem                       // Item ID
	Enchants               map[int]*Enchant                   // ItemEchantment ID
	ItemStatEffects        map[int]*ItemStatEffect            // ItemID? something anyway
	SpellEffects           map[int]map[int]SpellEffect        // Search by spellID and effect index
	RandomSuffix           map[int]*RandomSuffix              // Item level
	ItemDamageTable        map[string]map[int]ItemDamageTable // By Table name and item level
	RandomPropertiesByIlvl map[int32]RandomPropAllocationMap
	ItemArmorQuality       map[int]ItemArmorQuality
	ItemArmorShield        map[int]ItemArmorShield
	ItemArmorTotal         map[int]ItemArmorTotal
	ArmorLocation          map[int]ArmorLocation
}

func NewDBC() *DBC {
	return &DBC{
		Items:                  make(map[int]*Item),
		Gems:                   make(map[int]*Gem),
		Enchants:               make(map[int]*Enchant),
		ItemStatEffects:        make(map[int]*ItemStatEffect),
		SpellEffects:           make(map[int]map[int]SpellEffect),
		RandomSuffix:           make(map[int]*RandomSuffix),
		ItemDamageTable:        make(map[string]map[int]ItemDamageTable),
		RandomPropertiesByIlvl: make(map[int32]RandomPropAllocationMap),
		ItemArmorQuality:       make(map[int]ItemArmorQuality),
		ItemArmorShield:        make(map[int]ItemArmorShield),
		ItemArmorTotal:         make(map[int]ItemArmorTotal),
		ArmorLocation:          make(map[int]ArmorLocation),
	}
}

var (
	dbcInstance *DBC
	once        sync.Once
)

func GetDBC() *DBC {
	once.Do(func() {
		dbcInstance = NewDBC()

		if err := dbcInstance.loadItems("./assets/db_inputs/dbc/items.json"); err != nil {
			log.Fatalf("Error loading items: %v", err)
		}
		if err := dbcInstance.loadGems("./assets/db_inputs/dbc/gems.json"); err != nil {
			log.Fatalf("Error loading gems: %v", err)
		}
		if err := dbcInstance.loadEnchants("./assets/db_inputs/dbc/enchants.json"); err != nil {
			log.Fatalf("Error loading enchants: %v", err)
		}
		if err := dbcInstance.loadItemStatEffects("./assets/db_inputs/dbc/item_stat_effects.json"); err != nil {
			log.Fatalf("Error loading item stat effects: %v", err)
		}
		if err := dbcInstance.loadSpellEffects("./assets/db_inputs/dbc/spell_effects.json"); err != nil {
			log.Fatalf("Error loading spell effects: %v", err)
		}
		if err := dbcInstance.loadRandomSuffix("./assets/db_inputs/dbc/random_suffix.json"); err != nil {
			log.Fatalf("Error loading random suffixes: %v", err)
		}
		if err := dbcInstance.loadRandomPropertiesByIlvl("./assets/db_inputs/dbc/rand_prop_points.json"); err != nil {
			log.Fatalf("Error loading item damage tables: %v", err)
		}
		if err := dbcInstance.loadItemDamageTables("./assets/db_inputs/dbc/item_damage_tables.json"); err != nil {
			log.Fatalf("Error loading item damage tables: %v", err)
		}
		if err := dbcInstance.LoadItemArmorQuality("./assets/db_inputs/dbc/item_armor_quality.json"); err != nil {
			log.Fatalf("Error loading item damage tables: %v", err)
		}
		if err := dbcInstance.LoadItemArmorTotal("./assets/db_inputs/dbc/item_armor_total.json"); err != nil {
			log.Fatalf("Error loading item damage tables: %v", err)
		}
		if err := dbcInstance.LoadItemArmorShield("./assets/db_inputs/dbc/item_armor_shield.json"); err != nil {
			log.Fatalf("Error loading item damage tables: %v", err)
		}
		if err := dbcInstance.LoadArmorLocation("./assets/db_inputs/dbc/armor_location.json"); err != nil {
			log.Fatalf("Error loading item damage tables: %v", err)
		}
	})
	return dbcInstance
}

func (d *DBC) loadRandomPropertiesByIlvl(filename string) error {
	data, err := readGzipFile(filename)
	if err != nil {
		return err
	}
	var properties RandomPropAllocationsByIlvl
	if err = json.Unmarshal(data, &properties); err != nil {
		return err
	}
	d.RandomPropertiesByIlvl = properties
	return nil
}

func (d *DBC) loadItems(filename string) error {
	data, err := readGzipFile(filename)
	if err != nil {
		return err
	}
	var items []Item
	if err = json.Unmarshal(data, &items); err != nil {
		return err
	}
	for i := range items {
		item := &items[i]
		d.Items[item.Id] = item
	}
	return nil
}

func (d *DBC) loadGems(filename string) error {
	data, err := readGzipFile(filename)
	if err != nil {
		return err
	}
	var gems []Gem
	if err = json.Unmarshal(data, &gems); err != nil {
		return err
	}
	for i := range gems {
		gem := &gems[i]
		d.Gems[gem.ItemId] = gem
	}
	return nil
}

func (d *DBC) loadEnchants(filename string) error {
	data, err := readGzipFile(filename)
	if err != nil {
		return err
	}
	var enchants []Enchant
	if err = json.Unmarshal(data, &enchants); err != nil {
		return err
	}
	for i := range enchants {
		enchant := &enchants[i]
		d.Enchants[enchant.EffectId] = enchant
	}
	return nil
}

func (d *DBC) loadItemStatEffects(filename string) error {
	data, err := readGzipFile(filename)
	if err != nil {
		return err
	}
	var effects []ItemStatEffect
	if err = json.Unmarshal(data, &effects); err != nil {
		return err
	}
	for i := range effects {
		effect := &effects[i]
		d.ItemStatEffects[effect.ID] = effect
	}
	return nil
}

func (d *DBC) loadSpellEffects(filename string) error {
	data, err := readGzipFile(filename)
	if err != nil {
		return err
	}
	var effects map[int]map[int]SpellEffect
	if err = json.Unmarshal(data, &effects); err != nil {
		return err
	}
	d.SpellEffects = effects
	return nil
}

func (d *DBC) loadRandomSuffix(filename string) error {
	data, err := readGzipFile(filename)
	if err != nil {
		return err
	}
	var suffixes []RandomSuffix
	if err = json.Unmarshal(data, &suffixes); err != nil {
		return err
	}
	for i := range suffixes {
		suffix := &suffixes[i]
		d.RandomSuffix[suffix.ID] = suffix
	}
	return nil
}

func (d *DBC) LoadItemArmorQuality(filename string) error {
	data, err := readGzipFile(filename)
	if err != nil {
		return err
	}
	var tables map[int]ItemArmorQuality
	if err = json.Unmarshal(data, &tables); err != nil {
		return err
	}
	d.ItemArmorQuality = tables
	return nil
}
func (d *DBC) LoadArmorLocation(filename string) error {
	data, err := readGzipFile(filename)
	if err != nil {
		return err
	}
	var tables map[int]ArmorLocation
	if err = json.Unmarshal(data, &tables); err != nil {
		return err
	}
	d.ArmorLocation = tables
	return nil
}
func (d *DBC) LoadItemArmorShield(filename string) error {
	data, err := readGzipFile(filename)
	if err != nil {
		return err
	}
	var tables map[int]ItemArmorShield
	if err = json.Unmarshal(data, &tables); err != nil {
		return err
	}
	d.ItemArmorShield = tables
	return nil
}

func (d *DBC) LoadItemArmorTotal(filename string) error {
	data, err := readGzipFile(filename)
	if err != nil {
		return err
	}
	var tables map[int]ItemArmorTotal
	if err = json.Unmarshal(data, &tables); err != nil {
		return err
	}
	d.ItemArmorTotal = tables
	return nil
}

func (d *DBC) loadItemDamageTables(filename string) error {
	data, err := readGzipFile(filename)
	if err != nil {
		return err
	}
	var tables map[string]map[int]ItemDamageTable
	if err = json.Unmarshal(data, &tables); err != nil {
		return err
	}
	d.ItemDamageTable = tables
	return nil
}
