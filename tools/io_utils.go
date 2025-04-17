// Helper functions for reading/writing data.
package tools

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"golang.org/x/exp/constraints"
	protojson "google.golang.org/protobuf/encoding/protojson"
	googleProto "google.golang.org/protobuf/proto"
)

var readWebThreads = flag.Int("readWebThreads", 8, "number of parallel workers to fetch web pages")

func ReadFile(filePath string) string {
	b, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to open %s: %s", filePath, err)
	}
	return string(b)
}
func readFileLinesInternal(filePath string, throwIfMissing bool) []string {
	file, err := os.Open(filePath)
	if err != nil {
		if throwIfMissing {
			log.Fatalf("Failed to open %s: %s", filePath, err)
		} else {
			return nil
		}
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

func ReadMapOrNil(filePath string) map[string]string {
	return readMapInternal(filePath, false)
}
func readMapInternal(filePath string, throwIfMissing bool) map[string]string {
	res := make(map[string]string)
	if lines := readFileLinesInternal(filePath, throwIfMissing); lines != nil {
		for _, line := range lines {
			splitIndex := strings.Index(line, ",")
			keyStr := line[:splitIndex]
			valStr := line[splitIndex+1:]
			res[keyStr] = valStr
		}
	}
	return res
}

func WriteFile(filePath string, content string) {
	err := os.WriteFile(filePath, []byte(content), 0666)
	if err != nil {
		log.Fatalf("Failed to write file %s: %s", filePath, err)
	}
}

func WriteFileLines(filePath string, lines []string) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Failed to open %s for write: %s", filePath, err)
	}

	for _, line := range lines {
		file.WriteString(line)
		file.WriteString("\n")
	}
}

func WriteMapSortByIntKey(filePath string, contents map[string]string) {
	WriteMapCustomSort(filePath, contents, func(a, b string) int {
		intA, err1 := strconv.Atoi(a)
		intB, err2 := strconv.Atoi(b)
		if err1 != nil {
			panic(err1)
		}
		if err2 != nil {
			panic(err2)
		}
		return intA - intB
	})
}
func WriteMapCustomSort(filePath string, contents map[string]string, sortFunc func(a, b string) int) {
	type Elem struct {
		key string
		val string
	}

	elems := make([]Elem, len(contents))
	i := 0
	for k, v := range contents {
		elems[i] = Elem{key: k, val: v}
		i++
	}

	// Sort so the output is stable.
	slices.SortStableFunc(elems, func(a, b Elem) int {
		return sortFunc(a.key, b.key)
	})

	lines := core.MapSlice(elems, func(elem Elem) string {
		return fmt.Sprintf("%s,%s", elem.key, elem.val)
	})

	WriteFileLines(filePath, lines)
}

func WriteProtoArrayToBuffer[T googleProto.Message](arr []T, buffer *bytes.Buffer, name string) {
	buffer.WriteString("\"")
	buffer.WriteString(name)
	buffer.WriteString("\":[\n")

	for i, elem := range arr {
		jsonBytes, err := protojson.MarshalOptions{UseEnumNumbers: true}.Marshal(elem)
		if err != nil {
			log.Printf("[ERROR] Failed to marshal: %s", err.Error())
		}

		// Format using Compact() so we get a stable output (no random diffs for version control).
		json.Compact(buffer, jsonBytes)

		if i != len(arr)-1 {
			buffer.WriteString(",")
		}
		buffer.WriteString("\n")
	}
	buffer.WriteString("]")
}
func WriteProtoMapToBuffer[K constraints.Ordered, V googleProto.Message](
	m map[K]V,
	buffer *bytes.Buffer,
	name string,
) {
	buffer.WriteString(`"`)
	buffer.WriteString(name)
	buffer.WriteString(`":{` + "\n")

	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for i, k := range keys {
		jsonBytes, err := protojson.MarshalOptions{UseEnumNumbers: true}.Marshal(m[k])
		if err != nil {
			log.Printf("[ERROR] Failed to marshal key %v: %v", k, err)
			continue
		}

		buffer.WriteString(`"`)
		buffer.WriteString(fmt.Sprint(k))
		buffer.WriteString(`":`)
		// Format using Compact() so we get a stable output (no random diffs for version control).
		json.Compact(buffer, jsonBytes)

		if i != len(keys)-1 {
			buffer.WriteString(",")
		}
		buffer.WriteString("\n")
	}

	buffer.WriteString("}")
}
func WriteWeaponDamageToBuffer(wd *proto.WeaponDamageDatabase, buffer *bytes.Buffer, name string) {
	buffer.WriteString("\"")
	buffer.WriteString(name)
	buffer.WriteString("\": {\n")

	writeField := func(fieldName string, data []*proto.ItemQualityValue, isLast bool) {
		buffer.WriteString("\"")
		buffer.WriteString(fieldName)
		buffer.WriteString("\": [\n")

		for i, elem := range data {
			jsonBytes, err := protojson.MarshalOptions{UseEnumNumbers: true}.Marshal(elem)
			if err != nil {
				log.Printf("Couldn't to marshal: %s", err.Error())
			}
			json.Compact(buffer, jsonBytes)
			if i != len(data)-1 {
				buffer.WriteString(",")
			}
			buffer.WriteString("\n")
		}
		buffer.WriteString("]")
		if !isLast {
			buffer.WriteString(",\n")
		} else {
			buffer.WriteString("\n")
		}
	}

	writeField("caster_1h", wd.Caster_1H, false)
	writeField("melee_1h", wd.Melee_1H, false)
	writeField("caster_2h", wd.Caster_2H, false)
	writeField("melee_2h", wd.Melee_2H, false)
	writeField("ranged", wd.Ranged, false)
	writeField("thrown", wd.Thrown, false)
	writeField("wand", wd.Wand, true)

	buffer.WriteString("}")
}

func WriteArmorValuesToBuffer(ad *proto.ArmorValueDatabase, buffer *bytes.Buffer, name string) {
	buffer.WriteString("\"")
	buffer.WriteString(name)
	buffer.WriteString("\": {\n")

	writeField := func(fieldName string, data []*proto.ItemQualityValue, isLast bool) {
		buffer.WriteString("\"")
		buffer.WriteString(fieldName)
		buffer.WriteString("\": [\n")

		for i, elem := range data {
			jsonBytes, err := protojson.MarshalOptions{UseEnumNumbers: true}.Marshal(elem)
			if err != nil {
				log.Printf("[ERROR] Failed to marshal: %s", err.Error())
			}
			json.Compact(buffer, jsonBytes)
			if i != len(data)-1 {
				buffer.WriteString(",")
			}
			buffer.WriteString("\n")
		}
		buffer.WriteString("]")
		if !isLast {
			buffer.WriteString(",\n")
		} else {
			buffer.WriteString("\n")
		}
	}

	writeField("armor_values", ad.ArmorValues, false)
	writeField("shield_armor_values", ad.ShieldArmorValues, true)

	buffer.WriteString("}")
}

// Fetches web results a single url, and returns the page contents as a string.
func ReadWeb(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	resp.Body.Close()
	return string(body), nil
}
func ReadWebRequired(url string) string {
	body, err := ReadWeb(url)
	if err != nil {
		panic(err)
	}
	return body
}

// Fetches web results from all the given urls, and returns a parallel array of page contents.
func ReadWebMulti(urls []string) []string {
	threads := *readWebThreads
	if threads > len(urls) {
		threads = len(urls)
	}

	type WebResult struct {
		urlIdx int
		body   string
	}
	webResults := make(chan WebResult, 10)
	wg := &sync.WaitGroup{}

	for thread := 0; thread < threads; thread++ {
		startIdx := len(urls) * thread / threads
		endIdx := len(urls) * (thread + 1) / threads
		wg.Add(1)
		go func(min, max int) {
			fmt.Printf("ReadWebMulti Starting worker for URL block %d to %d\n", min, max-1)
			for i := min; i < max; i++ {
				url := urls[i]
				body, err := ReadWeb(url)
				if err != nil {
					fmt.Printf("ReadWebMulti Error fetching %s: %s\n", url, err)
					continue
				}
				webResults <- WebResult{urlIdx: i, body: body}
			}
			wg.Done()
		}(startIdx, endIdx)
	}

	go func() {
		wg.Wait()
		close(webResults)
	}()

	results := make([]string, len(urls))

	totalComplete := 0
	var lastUpdate time.Time
	for res := range webResults {
		totalComplete++

		if time.Since(lastUpdate).Seconds() > 2 {
			lastUpdate = time.Now()
			fmt.Printf("ReadWebMulti %d/%d complete\n", totalComplete, len(urls))
		}

		results[res.urlIdx] = res.body
	}
	fmt.Printf("ReadWebMulti %d/%d complete\n", totalComplete, len(urls))

	return results
}

// Like ReadWebMulti, but uses a lambda function for converting keys --> urls
// and returns a map of keys to web contents.
func ReadWebMultiMap[K comparable](keys []K, keyToUrl func(K) string) map[K]string {
	urls := core.MapSlice(keys, keyToUrl)
	results := ReadWebMulti(urls)

	mapResults := make(map[K]string, len(urls))
	for i := 0; i < len(urls); i++ {
		mapResults[keys[i]] = results[i]
	}
	return mapResults
}
