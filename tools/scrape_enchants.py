#!/usr/bin/python3

from enum import Enum
import requests
import csv
import os

from typing import Callable, List, Mapping, KeysView

## SET MINIMUM SPELLID HERE
MIN_SPELL_ID = 70165
IGNORE_SPELLS = [
    359641,
    359858,
    359640,
    359847,
    359949,
    359950,
    359639,
    359642,
    359685,
    359895,
]

def download_file(url, file_path):
    if os.path.exists(file_path):
        return

    response = requests.get(url)
    if response.status_code == 200:
        with open(file_path, 'wb') as file:
            file.write(response.content)
        print(f"File downloaded successfully to {file_path}")
    else:
        print(f"Failed to download file from {url}")


branch = "wow_classic_beta"
files = [
    "SpellEffect",
    "SpellItemEnchantment",
    "Spell",
    "SpellName",
    "ItemEffect",
    "SpellEquippedItems",
    "ItemSparse",
    "SkillLineAbility"
]

asset_dir = "/tmp"
for file in files:
    download_file(f"https://wago.tools/db2/{file}/csv?branch={branch}", f"{asset_dir}/{file}.csv")

statMapping: Mapping[str, Callable[[int], str]] = {
    "0": lambda val: f"stats.Mana: {val}",
    "1": lambda val: f"stats.Health: {val}",
    "3": lambda val: f"stats.Agility: {val}",
    "4": lambda val: f"stats.Strength: {val}",
    "5": lambda val: f"stats.Intellect: {val}",
    "6": lambda val: f"stats.Spirit: {val}",
    "7": lambda val: f"stats.Stamina: {val}",
    "13": lambda val: f"stats.Dodge: {val}",
    "14": lambda val: f"stats.Parry: {val}",
    "15": lambda val: f"stats.Block: {val}",
    # ranged hit - but not supported by core, so melee + spell i.g.?
    "17": lambda val: f"stats.MeleeHit: {val}, stats.SpellHit: {val}",
    "19": lambda val: f"stats.MeleeCrit: {val}",
    "20": lambda val: f"stats.MeleeCrit: {val}",

    # ranged haste - not supported by core so melee + spell i.g.?
    "29": lambda val: f"stats.MeleeHaste: {val}, stats.SpellHaste: {val}",

    # hit is the same for spell and melee in cata
    "31": lambda val: f"stats.MeleeHit: {val}, stats.SpellHit: {val}",
    "32": lambda val: f"stats.MeleeCrit: {val}, stats.SpellCrit: {val}",
    "35": lambda val: f"stats.Resilience: {val}",
    "36": lambda val: f"stats.MeleeHaste: {val}, stats.SpellHaste: {val}",
    "37": lambda val: f"stats.Expertise: {val}",
    "38": lambda val: f"stats.AttackPower: {val}, stats.RangedAttackPower: {val}",
    "39": lambda val: f"stats.RangedAttackPower: {val}",
    "43": lambda val: f"stats.MP5: {val}",
    "45": lambda val: f"stats.SpellPower: {val}",
    "47": lambda val: f"stats.SpellPenetration: {val}",
    "49": lambda val: f"stats.Mastery: {val}",
}

armorMapping: Mapping[str, Callable[[int], str]] = {
    "0": lambda val: f"stats.Armor: {val}",
    "2": lambda val: f"stats.FireResistance: {val}",
    "3": lambda val: f"stats.NatureResistance: {val}",
    "4": lambda val: f"stats.FrostResistance: {val}",
    "5": lambda val: f"stats.ShadowResistance: {val}",
    "6": lambda val: f"stats.ArcaneResistance: {val}"
}

class ItemClass(Enum):
    WEAPON = 2
    ARMOR = 4
    MISC = 15

class WeaponType(Enum):
    AXE_1H = 0
    AXE_2H = 1  
    BOW = 2
    GUN = 3
    MACE_1H = 4
    MACE_2H = 5  
    POLEARM = 6
    SWORD_1H = 7
    SWORD_2H = 8  
    STAFF = 10
    FIST_WEAPON = 13
    DAGGER = 15
    CROSSBOW = 18
    WAND = 19
    FISHING_POLE = 20

    def mask(weapon) -> int:
        return 1 << weapon.value
    
    @staticmethod
    def mask_oneHanded() -> int:
        return WeaponType.mask(WeaponType.AXE_1H)\
            | WeaponType.mask(WeaponType.MACE_1H)\
            | WeaponType.mask(WeaponType.SWORD_1H)\
            | WeaponType.mask(WeaponType.DAGGER)\
            | WeaponType.mask(WeaponType.FIST_WEAPON)

    @staticmethod
    def mask_twoHanded() -> int:
        return WeaponType.mask(WeaponType.AXE_2H)\
            | WeaponType.mask(WeaponType.MACE_2H)\
            | WeaponType.mask(WeaponType.SWORD_2H)\
            | WeaponType.mask(WeaponType.POLEARM)
    
    @staticmethod
    def mask_projectile() -> int:
        return WeaponType.mask(WeaponType.CROSSBOW)\
            | WeaponType.mask(WeaponType.BOW)\
            | WeaponType.mask(WeaponType.GUN)
    
    @staticmethod
    def mask_sharp() -> int:
        return WeaponType.mask(WeaponType.AXE_2H)\
            | WeaponType.mask(WeaponType.SWORD_2H)\
            | WeaponType.mask(WeaponType.SWORD_1H)\
            | WeaponType.mask(WeaponType.AXE_1H)\
            | WeaponType.mask(WeaponType.POLEARM)\
            | WeaponType.mask(WeaponType.DAGGER)
    
    @staticmethod
    def mask_blunt() -> int:
        return WeaponType.mask(WeaponType.MACE_2H)\
            | WeaponType.mask(WeaponType.MACE_1H)
    

class ItemSubclassArmor(Enum):
    AMOR_CLOTH = 1
    ARMOR_LEATHER = 2
    ARMOR_MAIL = 3
    ARMOR_PLATE = 4
    SHIELD = 6

    @classmethod
    def isArmor(cls, mask: int) -> bool:
        return mask & ((1 << cls.AMOR_CLOTH) | (1 << cls.ARMOR_LEATHER) | (1 << cls.ARMOR_MAIL) | (1 << cls.ARMOR_PLATE)) > 0
    
    @classmethod
    def isShield(cls, mask: int) -> bool:
        return mask & (1 << cls.SHIELD) > 0
    
class InventoryType(Enum):
    NO_EQUIP = 0
    HEAD = 1
    NECK = 2
    SHOULDERS = 3
    CHEST = 5
    WAIST = 6
    LEGS = 7
    FEET = 8
    WRISTS = 9
    HANDS = 10
    FINGER= 11
    TRINKET = 12
    SHIELD = 14
    RANGED = 15
    CLOAK = 16
    WEAPON = 17
    OFFHAND = 23



class SpellEffectRecord:
    def __init__(self, spellID, enchantID, itemRef):
        self.spellID: int = int(spellID)
        self.enchantID: int = int(enchantID)

        # item ref set for tournament realm scroll enchants
        self.itemRef: int = int(itemRef)

class SpellNameRecord:
    def __init__(self, spellID, spellName):
        self.spellID: int = int(spellID)
        self.name: str = spellName

class SpellDescription:
    def __init__(self, spellID, description):
        self.spellID: int = int(spellID)
        self.description = description

class EnchantEffect:
    def __init__(self, type, points, misc):
        self.type = type
        self.points = points
        self.misc = misc

class SparseItemInfo:
    def __init__(self, id, quality, allowedClassMask):
        self.id: int = int(id)
        self.quality: int = int(quality)
        self.allowedClassMask: int = int(allowedClassMask)

skillLineMap: Mapping[int, str] = {
    164: "proto.Profession_Blacksmithing",
    165: "proto.Profession_Leatherworking",
    197: "proto.Profession_Tailoring",
    202: "proto.Profession_Engineering",
    333: "proto.Profession_Enchanting",
    755: "proto.Profession_Jewelcrafting",
    773: "proto.Profession_Inscription",
    776: "runeforging"
}

class ClassMask(Enum):
    WARRIOR = 0x1
    PALADIN = 0x2
    HUNTER = 0x4
    ROGUE = 0x8
    PRIEST = 0x10
    DEATH_KNIGHT = 0x20
    SHAMAN = 0x40
    MAGE = 0x80
    WARLOCK = 0x100
    DRUID = 0x200

class SkillLineAbility:
    def __init__(self, spellId, skillLine):
        self.spellID: int = spellId
        self.skillLine: int = int(skillLine)

class SpellEnchant:
    def __init__(self, enchantID, effect_0, points_0, misc_0, effect_1, points_1, misc_1, effect_2, points_2, misc_2, tradeSkill):
        self.enchantID: int = int(enchantID)
        self.effects : List[EnchantEffect] = []
        if effect_0 != "0":
            self.effects.append(EnchantEffect(effect_0, points_0, misc_0))

        if effect_1  != "0":
            self.effects.append(EnchantEffect(effect_1, points_1, misc_1))

        if effect_2  != "0":
            self.effects.append(EnchantEffect(effect_2, points_2, misc_2))

        self.tradeSkill: int = int(tradeSkill)

    def getStatString(self) -> str:
        result = ""
        for effect in self.effects:
            lookup = None
            if  effect.type == "4":
                lookup = armorMapping
            elif effect.type == "5":
                lookup = statMapping
            else:
                continue

            strStat = "Unk"
            if not effect.misc in lookup:
                print(F"WARN: {effect.misc}. Unknown stat.")
            else:
                strStat = lookup[effect.misc]

            if len(result) > 0:
                result += ", "
            
            result += f"{strStat(effect.points)}"
        return result
    
slotProtoMap: Mapping[InventoryType, str] = {
    InventoryType.HEAD: "proto.ItemType_ItemTypeHead",
    InventoryType.SHOULDERS: "proto.ItemType_ItemTypeShoulder",
    InventoryType.CLOAK: "proto.ItemType_ItemTypeBack",
    InventoryType.CHEST: "proto.ItemType_ItemTypeChest",
    InventoryType.WRISTS: "proto.ItemType_ItemTypeWrist",
    InventoryType.HANDS: "proto.ItemType_ItemTypeHands",
    InventoryType.WAIST: "proto.ItemType_ItemTypeWaist",
    InventoryType.LEGS: "proto.ItemType_ItemTypeLegs",
    InventoryType.FEET: "proto.ItemType_ItemTypeFeet",
    InventoryType.FINGER: "proto.ItemType_ItemTypeFinger",
    InventoryType.WEAPON: "proto.ItemType_ItemTypeWeapon",
    InventoryType.RANGED: "proto.ItemType_ItemTypeRanged",
}

class SpellItemRequirement:
    def __init__(self, spellID, itemClass, itemSubClass, inventoryType):
        self.spellID: int = int(spellID)
        self.itemClass : ItemClass = ItemClass(int(itemClass))
        self.itemSubClass : int = int(itemSubClass)
        self.inventoryTypeMask : int = int(inventoryType)
    
    def getItemTypeStr(self) -> str:
        inv = self.getInventoryType()
        return slotProtoMap[inv]

    def getInventoryType(self) -> InventoryType:
        if self.itemClass == ItemClass.WEAPON:
            if WeaponType.mask_projectile() & self.itemSubClass > 0:
                return InventoryType.RANGED
            return InventoryType.WEAPON
        
        if self.itemClass == ItemClass.ARMOR:
            return self._getArmorSlot()

        print("WARN: Unknown item class: " + str(self.itemClass))

    def _getArmorSlot(self) -> InventoryType:
        if self.itemSubClass & (1 << ItemSubclassArmor.SHIELD.value) > 0 or self.inventoryTypeMask & (1 << InventoryType.SHIELD.value) > 0:
            return InventoryType.WEAPON
        
        for slot in InventoryType:
            if self.inventoryTypeMask & (1 <<  slot.value) > 0:
                return slot
            
        print("WARN: Unknown armor slot(s): " + str(self.inventoryTypeMask))


class EnchantItemInfo:
    def __init__(self, spellID, itemID):
        self.spellID :int = int(spellID)
        self.itemID :int = int(itemID)

def loadEffects() -> Mapping[int, SpellEffectRecord]:
    effects = {}
    with open(f'{asset_dir}/SpellEffect.csv', newline='') as csvfile:
        reader = csv.DictReader(csvfile)
        for row in reader:
            if row['Effect'] != "53":
                continue

            spellID = int(row['SpellID'])
            if spellID < MIN_SPELL_ID:
                continue
            if spellID in IGNORE_SPELLS:
                continue

            effect = SpellEffectRecord(spellID, row['EffectMiscValue_0'], row['EffectItemType'])
            effects[effect.spellID] = effect

    return effects

def loadEnchants() -> Mapping[int, SpellEnchant]:
    enchants = {}
    with open(f'{asset_dir}/SpellItemEnchantment.csv', newline='') as csvfile:
        reader = csv.DictReader(csvfile)
        for row in reader:
            enchant = SpellEnchant(
                row['ID'],
                row['Effect_0'],
                row['EffectPointsMin_0'],
                row['EffectArg_0'],
                row['Effect_1'],
                row['EffectPointsMin_1'],
                row['EffectArg_1'],
                row['Effect_2'],
                row['EffectPointsMin_2'],
                row['EffectArg_2'],
                row['RequiredSkillID'])
            enchants[enchant.enchantID] = enchant

    return enchants

def loadSpellNames(knownSpells: KeysView[int]) -> Mapping[int, SpellNameRecord]:
    names = {}
    with open(f'{asset_dir}/SpellName.csv', newline='') as csvfile:
        reader = csv.DictReader(csvfile)
        for row in reader:
            id = int(row['ID'])
            if not (id in knownSpells):
                continue

            names[id] = SpellNameRecord(id, row['Name_lang'])
    return names

def loadSpellDescriptions(knownSpells: KeysView[int]) -> Mapping[int, SpellDescription]:
    names = {}
    with open(f'{asset_dir}/Spell.csv', newline='') as csvfile:
        reader = csv.DictReader(csvfile)
        for row in reader:
            id = int(row['ID'])
            if not (id in knownSpells):
                continue

            names[id] = SpellDescription(id, row['Description_lang'])
    return names

def loadSlotRequirements(knownSpells: KeysView[int]) -> Mapping[int, SpellItemRequirement]:
    names = {}
    with open(f'{asset_dir}/SpellEquippedItems.csv', newline='') as csvfile:
        reader = csv.DictReader(csvfile)
        for row in reader:
            id = int(row['SpellID'])
            if not (id in knownSpells):
                continue

            names[id] = SpellItemRequirement(id, row['EquippedItemClass'], row['EquippedItemSubclass'], row["EquippedItemInvTypes"])
    return names

def loadItemInfos(knownSpells: KeysView[int]) -> Mapping[int, EnchantItemInfo]:
    itemInfo = {}
    with open(f'{asset_dir}/ItemEffect.csv', newline='') as csvfile:
        reader = csv.DictReader(csvfile)
        for row in reader:
            id = int(row['SpellID'])
            if not (id in knownSpells):
                continue

            itemInfo[id] = EnchantItemInfo(id, row['ParentItemID'])
    return itemInfo

def loadItemSparse(relevantItems: KeysView[int]) -> Mapping[int, SparseItemInfo]:
    itemInfo = {}
    with open(f'{asset_dir}/ItemSparse.csv', newline='') as csvfile:
        reader = csv.DictReader(csvfile)
        for row in reader:
            id = int(row['ID'])
            if not (id in relevantItems):
                continue

            itemInfo[id] = SparseItemInfo(id, row['OverallQualityID'], row['AllowableClass'])
    return itemInfo

def loadSkillLines(knownSpells: List[int]) -> Mapping[int, SkillLineAbility]:
    skillLineInfo = {}
    with open(f'{asset_dir}/SkillLineAbility.csv', newline='') as csvfile:
        reader = csv.DictReader(csvfile)
        for row in reader:
            id = int(row['Spell'])
            if not (id in knownSpells):
                continue

            skillLineInfo[id] = SkillLineAbility(id, row['SkillLine'])
    return skillLineInfo

qualityMap: Mapping[int, str] = {
    0: "ItemQuality_ItemQualityJunk",
    1: "proto.ItemQuality_ItemQualityCommon",
    2: "proto.ItemQuality_ItemQualityUncommon",
    3: "proto.ItemQuality_ItemQualityRare",
    4: "proto.ItemQuality_ItemQualityEpic",
    5: "proto.ItemQuality_ItemQualityLegendary",
    7: "proto.ItemQuality_ItemQualityRare",  # Account bound
}

classMap: Mapping[ClassMask, str] = {
    ClassMask.WARRIOR: "proto.Class_ClassWarrior",
    ClassMask.DEATH_KNIGHT: "proto.Class_ClassDeathKnight",
    ClassMask.DRUID: "proto.Class_ClassDruid",
    ClassMask.HUNTER: "proto.Class_ClassHunter",
    ClassMask.MAGE: "proto.Class_ClassMage",
    ClassMask.PALADIN: "proto.Class_ClassPaladin",
    ClassMask.PRIEST: "proto.Class_ClassPriest",
    ClassMask.ROGUE: "proto.Class_ClassRogue",
    ClassMask.SHAMAN: "proto.Class_ClassShaman",
    ClassMask.WARLOCK: "proto.Class_ClassWarlock",
}

spellEffects = loadEffects()
enchants = loadEnchants()
spellNames = loadSpellNames(spellEffects.keys())
spellDescriptions = loadSpellDescriptions(spellEffects.keys())
slotRequirements = loadSlotRequirements(spellEffects.keys())
itemInfo = loadItemInfos(spellEffects.keys())
itemSparse = loadItemSparse(list(map(lambda info: info.itemID, itemInfo.values())))
skillLines = loadSkillLines(spellEffects.keys())

def getQuality(spellId: str) -> str: 
    if spellId in itemInfo:
        itemId = itemInfo[spellId].itemID
        return qualityMap[itemSparse[itemId].quality]
    
    return qualityMap[1]

def getItemIDPart(effect: SpellEffectRecord) -> str:
    if effect.itemRef > 0 or not effect.spellID in itemInfo:
        return ""
    
    return f", ItemId: {itemInfo[effect.spellID].itemID}"

def getRequiredSkillLine(effect: SpellEffectRecord) -> int:
    enchant = enchants[effect.enchantID]
    if enchant.tradeSkill > 0:
        return enchant.tradeSkill

    if effect.spellID in skillLines:
        skillLine = skillLines[effect.spellID]

        # not enchanting
        if skillLine.skillLine != 333:
            return skillLines[effect.spellID].skillLine
    
    return 0
    
def getProfessionPart(effect: SpellEffectRecord) -> str:
    tradeSkill = getRequiredSkillLine(effect)
    if tradeSkill > 0:
        return ", RequiredProfession: " + skillLineMap[tradeSkill]
        
    return ""

def getClassRestrictionPart(effect: SpellEffectRecord) -> str:
    tradeSkill = getRequiredSkillLine(effect)
    if tradeSkill == 776: # DK
        return ", ClassAllowlist: []proto.Class{proto.Class_DeathKnight}"
    
    if not effect.spellID in itemInfo:
        return ""
    
    allowed = itemSparse[itemInfo[effect.spellID].itemID].allowedClassMask
    if allowed == 0 or allowed == -1:
        return ""
    
    classes = list(filter(lambda t: (allowed & t.value) > 0, ClassMask))
    if len(classes) == 0:
        return ""
    
    classProto = map(lambda cl: classMap[cl], classes)
    return ", ClassAllowlist: []proto.Class{" + ", ".join(classProto) + "}"
    

def getExtraSlotsPart(effect: SpellEffectRecord) -> str:
    if slotRequirements[effect.spellID].itemClass == ItemClass.WEAPON:
        return ""    
    slots = list(filter(lambda t: (slotRequirements[effect.spellID].inventoryTypeMask & (1 << t.value)) > 0 and not t == type, InventoryType))
    slots = list(filter(lambda t: t != InventoryType.OFFHAND, slots))
    if len(slots) < 2:
        return ""
    slots = slots[1:]
    protoSlots = map(lambda inv: slotProtoMap[inv], slots)
    return ", ExtraTypes: []proto.ItemType{" + ', '.join(protoSlots) + "}"

def getEnchantTypePart(effect: SpellEffectRecord) -> str:
    equipInfo = slotRequirements[effect.spellID]
    enchantType = ""

    if equipInfo.itemClass == ItemClass.ARMOR and equipInfo.inventoryTypeMask & (1 << InventoryType.OFFHAND.value) > 0:
        enchantType = "proto.EnchantType_EnchantTypeOffHand"
    elif equipInfo.itemClass == ItemClass.ARMOR and equipInfo.itemSubClass & (1 << ItemSubclassArmor.SHIELD.value) > 0:
        enchantType = "proto.EnchantType_EnchantTypeShield"
    if equipInfo.itemClass == ItemClass.WEAPON and equipInfo.itemSubClass == (1 << WeaponType.STAFF.value):
        enchantType = "proto.EnchantType_EnchantTypeStaff"
    elif equipInfo.itemSubClass & WeaponType.mask_twoHanded() == WeaponType.mask_twoHanded() and \
        equipInfo.itemSubClass & WeaponType.mask_oneHanded() == 0:
        enchantType = "proto.EnchantType_EnchantTypeTwoHand"

    if len(enchantType) > 0:
        return ", EnchantType: " + enchantType
    
    return ""

def shouldSkip(effect: SpellEffectRecord) -> bool:
    if effect.spellID not in slotRequirements:
        return True
    
    if effect.spellID not in spellNames:
        return True
    
    if slotRequirements[effect.spellID].itemClass == ItemClass.MISC:
        return True
    
    name = spellNames[effect.spellID].name
    if name.startswith("QA") or name.find("test") > 0 or name.find("Test") > 0:
        return True
    
    return False

# order by slot
spellMap: Mapping[InventoryType, List[str]] = {}
for id in sorted(spellEffects.keys()):
    effect = spellEffects[id]
    if shouldSkip(effect):
        continue


    it = slotRequirements[effect.spellID].getInventoryType()
    if not it in spellMap:
        spellMap[it] = []
    
    spellMap[it].append(f"{{EffectId: {effect.enchantID}{getItemIDPart(effect)}, SpellId: {id}, Name: \"{spellNames[id].name}\", Quality: {getQuality(id)}, Stats: stats.Stats{{{enchants[effect.enchantID].getStatString()}}}.ToFloatArray(), Type: {slotRequirements[id].getItemTypeStr()}{getEnchantTypePart(effect)}{getExtraSlotsPart(effect)}{getClassRestrictionPart(effect)}{getProfessionPart(effect)}}},\n")


with open("/tmp/enchants.tmp", 'w') as file:
    for type in spellMap.keys():
        file.write(f"// {type.name}\n")
        file.writelines(spellMap[type])
        file.write("\n\n")
