#!/usr/bin/python

# This tool generates the talents proto code, e.g. 'ShamanTalents' found in proto/shaman.proto.

import json
import sys

from selenium import webdriver
from selenium.webdriver.chrome.service import Service as ChromeService
from webdriver_manager.chrome import ChromeDriverManager
from selenium.webdriver.common.by import By
from selenium.webdriver.chrome.options import Options
chrome_options = Options()
chrome_options.add_argument('--headless')
chrome_options.add_argument('--no-sandbox')
chrome_options.add_argument('--disable-dev-shm-usage')
driver = webdriver.Chrome(service=ChromeService(ChromeDriverManager().install()),options=chrome_options)
if len(sys.argv) < 3:
	raise Exception("Missing arguments, expected class_name and output_file_path")
class_name = sys.argv[1]
output_file_path = sys.argv[2]

driver.implicitly_wait(2)

driver.get('https://wowhead.com/mop-classic/talent-calc/' + class_name)
trees = driver.find_elements(By.CLASS_NAME, "ctc-tree")

with open(output_file_path, "w") as outfile:
	# Convert "death-knight" to DeathKnight
	pretty_class_name = "".join(word.title() for i, word in enumerate(class_name.split("-")))
	outfile.write("message {}Talents {{\n".format(pretty_class_name))

	talent_idx = 1
	print(trees)
	for tree_idx, tree in enumerate(trees):
		title = tree.find_element(By.XPATH, ".//b").text
		outfile.write("\t// {}\n".format(title))

		tree_talents_data = []
		talents = tree.find_elements(By.CLASS_NAME, "ctc-tree-talent")
		for talent in talents:
			max_points = int(talent.get_attribute("data-max-points"))
			field_type = "bool" if max_points == 1 else "int32"

			link = talent.find_element(By.XPATH, "./div/a").get_attribute("href")
			name = "_".join(word for i, word in enumerate(link.split("/")[-1].split("-")))

			print("Talent: " + name)
			tree_talents_data.append({
				"row": int(talent.get_attribute("data-row")),
				"col": int(talent.get_attribute("data-col")),
				"name": name,
				"field_type": field_type,
			})

		tree_talents_data.sort(key=lambda data: data["row"] * 4 + data["col"])
		for data in tree_talents_data:
			outfile.write("\t{} {} = {};\n".format(data["field_type"], data["name"], talent_idx))
			talent_idx += 1

		if tree_idx != len(trees) - 1:
			outfile.write("\n")

	outfile.write("}")
