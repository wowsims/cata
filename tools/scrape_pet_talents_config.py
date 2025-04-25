#!/usr/bin/python

# This tool generates the talents config code, e.g. in ui/core/talents/shaman.ts.

import json
import sys

from typing import List

from selenium import webdriver
from selenium.webdriver.chrome.service import Service as ChromeService
from webdriver_manager.chrome import ChromeDriverManager
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC

chrome_options = Options()
chrome_options.add_argument('--headless')
chrome_options.add_argument('--no-sandbox')
chrome_options.add_argument('--disable-dev-shm-usage')
driver = webdriver.Chrome(service=ChromeService(ChromeDriverManager().install()),options=chrome_options)

if len(sys.argv) < 3:
	raise Exception("Missing arguments, expected className and outputFilePath")
className = sys.argv[1]
outputFilePath = sys.argv[2]

def _between(s, start, end):
	return s[(i := s.find(start) + len(start)): i + s[i:].find(end)]


driver.implicitly_wait(2)


def _get_spell_id_from_link(link):
    parts = link.split("/")
    for part in parts:
        if "spell=" in part:
            return int(part.split("=")[-1])
    raise ValueError("No spell ID found in the link")



def get_other_spell_ranks(spell_name: str, ignore_int: int) -> List[int]:
	overrides = {
		"T N T": "t.n.t",
	}

	formatted_spell_name = overrides.get(spell_name, spell_name.replace(' ', '+'))
	driver.get(f"https://www.wowhead.com/mop-classic/spells/pet-abilities/{className}/name:{formatted_spell_name}")
	driver.implicitly_wait(2)

	spell_ids = []
	table = driver.find_element(By.CLASS_NAME, "listview-mode-default")
	rows = table.find_elements(By.CLASS_NAME, "listview-row")
	print(f"Found {len(rows)} for {spell_name}")
	for row in rows:
		questionable_elements = row.find_elements(By.XPATH, ".//*[contains(@style, 'inv_misc_questionmark.jpg')]")

		if questionable_elements:
            # If any questionable elements found, skip this row
			print("Skipping row")
			continue
		a_element = row.find_element(By.CLASS_NAME, "listview-cleartext")
		href = a_element.get_attribute("href")
		spell_id = _get_spell_id_from_link(href)
		if spell_id != ignore_int:
			spell_ids.append(spell_id)

	return spell_ids

def rowcol(v):
	return v["location"]["rowIdx"] + v["location"]["colIdx"]/10


to_export = []

driver.get('https://wowhead.com/mop-classic/pet-talent-calc/')
try:
    wait = WebDriverWait(driver, 10)
    accept_button = wait.until(EC.element_to_be_clickable((By.ID, "onetrust-accept-btn-handler")))
    accept_button.click()
    print("Clicked the 'I Accept' button.")
except Exception as e:
    print("No 'I Accept' button to click or error clicking it:", e)
driver.implicitly_wait(2)
trees = driver.find_elements(By.CLASS_NAME, "ctc-tree")
for tree in trees:
	_working_talents = {}

	talents = tree.find_elements(By.CLASS_NAME, "ctc-tree-talent")
	print("found %d talents\n".format(len(talents)))
	for talent in talents:
		row, col = int(talent.get_attribute("data-row")), int(talent.get_attribute("data-col"))
		max_points = int(talent.get_attribute("data-max-points"))
		link = talent.find_element(By.XPATH, "./div/a").get_attribute("href")
		name = "".join(word if i == 0 else word.title() for i, word in enumerate(link.split("/")[-1].split("-")))
		fancyName = " ".join(word.title() for i, word in enumerate(link.split("/")[-1].split("-")))
		print(link)
		print(_get_spell_id_from_link(link))
		_working_talents[(row, col)] = {
			"fieldName": name,
			"fancyName": fancyName,
			"location": {
				"rowIdx": row,
				"colIdx": col,
			},
			"spellIds": [_get_spell_id_from_link(link)],
			"maxPoints": max_points,
		}

	arrows = tree.find_elements(By.CLASS_NAME, "ctc-tree-talent-arrow")
	for arrow in arrows:
		prereq_row, prereq_col = int(arrow.get_attribute("data-row")), int(arrow.get_attribute("data-col"))
		length = 0
		dsstr = arrow.get_attribute("data-size")
		if dsstr:
			length = int(dsstr)

		direction = arrow.get_attribute("class").split()[-1].split("-")[-1]
		offset_row, offset_col = {"left": (0, -1), "right": (0, 1), "down": (1, 0)}[direction]

		end_row = prereq_row + offset_row * length
		end_col = prereq_col + offset_col * length

		_working_talents[(end_row, end_col)]["prereqLocation"] = {
			"rowIdx": prereq_row,
			"colIdx": prereq_col,
		}

	title = tree.find_element(By.XPATH, ".//b").text
	background = tree.find_element(By.CLASS_NAME, "ctc-tree-talents-background").get_attribute("style")
	values = list(_working_talents.values())
	values.sort(key=rowcol)
	to_export.append({
		"name": title,
		"backgroundUrl": _between(background, '"', '"'),
		"talents": values,
	})

for subtree in to_export:
	for talent in subtree["talents"]:
		if talent["maxPoints"] > 1:
			talent["spellIds"] += get_other_spell_ranks(talent["fancyName"], talent["spellIds"][0])

json_data = json.dumps(to_export, indent=2)
with open(outputFilePath, "w") as outfile:
	outfile.write(json_data)
