#!/usr/bin/python3

import csv
import os
import sys
import json
from typing import List, Mapping, Set

import requests

if len(sys.argv) < 3:
    raise Exception("Missing arguments, expected db_path, output_file_path")

input_db_path = sys.argv[1]
output_file_path = sys.argv[2]
branch = "wow_classic_beta"

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

download_file(f"https://wago.tools/db2/SpellItemEnchantment/csv?branch={branch}", f"/tmp/SpellItemEnchantment.csv")

class SpellEnchant:
    def __init__(self, id, description):
        self.enchantID = int(id)
        self.description = description
        
def loadEnchants() -> Mapping[int, SpellEnchant]:
    enchants = {}
    with open('/tmp/SpellItemEnchantment.csv', newline='') as csvfile:
        reader = csv.DictReader(csvfile)
        for row in reader:
            enchant = SpellEnchant(
                row['ID'],
                row['Name_lang']
            )
            
            enchants[int(row['ID'])] = enchant
    
    return enchants

def loadDB() -> Set[int]:
    with open(input_db_path) as dbfile:
        dbjson = json.load(dbfile)
        enchants = []
        for enchant in dbjson['enchants']:
            enchants.append(enchant['effectId'])

    return sorted(enchants)


enchantDesc = loadEnchants()
enchantDB = loadDB()

output = {}
for enchant in enchantDB:
    if not enchant in enchantDesc:
        print (f"WARN: {enchant} not found in SpellItemEnchantment.")
        continue
    
    output[enchant] = enchantDesc[enchant].description

with open(output_file_path, 'w') as o:
    json.dump(output, o)