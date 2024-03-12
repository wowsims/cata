#!/usr/bin/env bash

# WoW Classic classes
classes=("druid" "death-knight" "hunter" "mage" "paladin" "priest" "rogue" "shaman" "warlock" "warrior")

for class in "${classes[@]}"; do
    # Replace "-" with "" for filenames
    #python3 ./scrape_talents_proto.py $class ../proto/${class}.proto
    python3 ./scrape_glyphs.py $class ../proto/tmp/${class}.proto.test
    #python3 ./scrape_talents_config.py $class ../ui/core/talents/trees/${class}.json
done