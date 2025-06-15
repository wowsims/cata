package dbc

import (
	"bufio"
	_ "embed"
	"strconv"
	"strings"

	"github.com/wowsims/mop/sim/core/proto"
)

//go:embed GameTables/SpellScaling.txt
var spellScalingFile string

type SpellScaling struct {
	Level  int
	Values map[proto.Class]float64
}

func (dbc *DBC) LoadSpellScaling() error {
	scanner := bufio.NewScanner(strings.NewReader(spellScalingFile))

	scanner.Scan() // Skip first line

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) < 14 {
			continue // consider handling or logging this situation
		}

		level, err := strconv.Atoi(parts[0])
		if err != nil {
			continue // consider handling or logging this situation
		}

		scaling := SpellScaling{
			Level: level,
			Values: map[proto.Class]float64{
				proto.Class_ClassWarrior:     parseScalingValue(parts[1]),
				proto.Class_ClassPaladin:     parseScalingValue(parts[2]),
				proto.Class_ClassHunter:      parseScalingValue(parts[3]),
				proto.Class_ClassRogue:       parseScalingValue(parts[4]),
				proto.Class_ClassPriest:      parseScalingValue(parts[5]),
				proto.Class_ClassDeathKnight: parseScalingValue(parts[6]),
				proto.Class_ClassShaman:      parseScalingValue(parts[7]),
				proto.Class_ClassMage:        parseScalingValue(parts[8]),
				proto.Class_ClassWarlock:     parseScalingValue(parts[9]),
				proto.Class_ClassMonk:        parseScalingValue(parts[10]),
				proto.Class_ClassDruid:       parseScalingValue(parts[11]),
				proto.Class_ClassExtra1:      parseScalingValue(parts[12]),
				proto.Class_ClassExtra2:      parseScalingValue(parts[13]),
				proto.Class_ClassExtra3:      parseScalingValue(parts[14]),
				proto.Class_ClassExtra4:      parseScalingValue(parts[15]),
				proto.Class_ClassExtra5:      parseScalingValue(parts[16]),
				proto.Class_ClassExtra6:      parseScalingValue(parts[17]),
			},
		}
		dbc.SpellScalings[level] = scaling
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (dbc *DBC) SpellScaling(class proto.Class, level int) float64 {
	if scaling, ok := dbc.SpellScalings[level]; ok {
		if value, ok := scaling.Values[class]; ok {
			return value
		}
	}
	return 0.0 // return a default or error value if not found
}

func parseScalingValue(value string) float64 {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0.0 // consider how to handle or log this error properly
	}
	return v
}
