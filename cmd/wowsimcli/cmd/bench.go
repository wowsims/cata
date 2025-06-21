package cmd

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"google.golang.org/protobuf/encoding/protojson"
	pb "google.golang.org/protobuf/proto"
)

var (
	benchCmd = &cobra.Command{
		Use:   "bench",
		Short: "sweep stats in increments and record results",
		Run:   benchMain,
	}

	stats []string
	step  int
	steps int
)

func init() {
	benchCmd.Flags().StringVar(&infile, "infile", "input.json", "location of base input file (RaidSimRequest in protojson format)")
	benchCmd.Flags().StringVar(&outfile, "outfile", "bench_results.csv", "path to write CSV results")
	benchCmd.Flags().StringSliceVar(&stats, "stats", []string{}, "list of stat names to sweep (e.g. Strength, Agility)")
	benchCmd.Flags().IntVar(&step, "step", 160, "increment amount per step")
	benchCmd.Flags().IntVar(&steps, "steps", 5, "number of steps up and down")
	benchCmd.Flags().BoolVar(&verbose, "verbose", false, "show detailed progress")

	benchCmd.MarkFlagRequired("stats")
}

// name->proto.Stat mapping
var statMap = map[string]proto.Stat{
	"Strength":        proto.Stat_StatStrength,
	"Agility":         proto.Stat_StatAgility,
	"Stamina":         proto.Stat_StatStamina,
	"Intellect":       proto.Stat_StatIntellect,
	"Spirit":          proto.Stat_StatSpirit,
	"HitRating":       proto.Stat_StatHitRating,
	"CritRating":      proto.Stat_StatCritRating,
	"HasteRating":     proto.Stat_StatHasteRating,
	"ExpertiseRating": proto.Stat_StatExpertiseRating,
	"DodgeRating":     proto.Stat_StatDodgeRating,
	"ParryRating":     proto.Stat_StatParryRating,
	"MasteryRating":   proto.Stat_StatMasteryRating,
	// add others as needed
}

func benchMain(cmd *cobra.Command, args []string) {
	// load base request
	data, err := os.ReadFile(infile)
	if err != nil {
		log.Fatalf("failed to load input json file %q: %v", infile, err)
	}
	base := &proto.RaidSimRequest{}
	if err := (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(data, base); err != nil {
		log.Fatalf("failed to parse input json: %v", err)
	}
	// prepare CSV output
	f, err := os.Create(outfile)
	if err != nil {
		log.Fatalf("could not create outfile: %v", err)
	}
	defer f.Close()
	writer := csv.NewWriter(f)
	defer writer.Flush()

	// header: Stat,Delta,Metric
	writer.Write([]string{"Stat", "Delta", "AverageDPS"})

	// for each stat, sweep
	for _, name := range stats {
		statKey, ok := statMap[name]
		if !ok {
			log.Fatalf("unknown stat name: %s", name)
		}

		for i := -steps; i <= steps; i++ {
			delta := float64(i * step)
			// clone base
			runReq := pb.Clone(base).(*proto.RaidSimRequest)
			// adjust first party player
			runReq.Raid.Parties[0].Players[0].BonusStats.Stats[statKey] += delta

			if verbose {
				fmt.Printf("Running %s %+d...\n", name, delta)
			}

			// run sim
			reporter := make(chan *proto.ProgressMetrics, 1)
			core.RunRaidSimConcurrentAsync(runReq, reporter, fmt.Sprintf("bench-%s-%d", strings.ToLower(name), i))
			var final *proto.RaidSimResult
			for m := range reporter {
				if m.FinalRaidResult != nil {
					final = m.FinalRaidResult
					break
				}
			}

			// record a metric, e.g., average DPS
			met := "0"
			if final.EncounterMetrics != nil {
				met = fmt.Sprintf("%.2f", final.RaidMetrics.Dps.Avg)
			}
			writer.Write([]string{name, strconv.Itoa(int(delta)), met})
		}
	}

	if verbose {
		fmt.Printf("Wrote bench results to %s\n", outfile)
	}
}
