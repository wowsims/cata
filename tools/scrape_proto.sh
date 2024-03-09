#!/usr/bin/env bash

# WoW Classic classes
classes=("druid" "death-knight" "hunter" "mage" "paladin" "priest" "rogue" "shaman" "warlock" "warrior")

for class in "${classes[@]}"; do
    # Replace "-" with "" for filenames
    class_filename="${class//-/}"
    #python3 ./scrape_talents_proto.py $class ../proto/${class_filename}.proto
    python3 ./scrape_glyphs.py $class ../proto/tmp/${class_filename}.proto.test
    #python3 ./scrape_talents_config.py $class ../ui/core/talents/trees/${class_filename}.json
done
