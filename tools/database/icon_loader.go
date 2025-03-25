package database

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func LoadArtTexturePaths(filePath string) (map[int]string, error) {
	paths := make(map[int]string)

	// Open the file only once.
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Regex to capture the key and its corresponding path.
	// It matches lines like:
	//    [130646]="Interface/AbilitiesFrame/UI-AbilityPanel-BotLeft",
	// and captures "130646" and "Interface/AbilitiesFrame/UI-AbilityPanel-BotLeft".
	re := regexp.MustCompile(`\[\s*(\d+)\s*\]\s*=\s*"(.*?)"`)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) == 3 {
			key, _ := strconv.Atoi(matches[1])
			value := matches[2]
			paths[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return paths, nil
}

func GetIconName(artPaths map[int]string, fdid int) string {
	path, ok := artPaths[fdid]
	if !ok {
		return ""
	}

	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}
