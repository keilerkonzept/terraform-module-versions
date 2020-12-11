package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/keilerkonzept/terraform-module-versions/pkg/modulecall"
	"github.com/keilerkonzept/terraform-module-versions/pkg/output"
	"github.com/keilerkonzept/terraform-module-versions/pkg/registry"
	"github.com/keilerkonzept/terraform-module-versions/pkg/scan"
	"github.com/keilerkonzept/terraform-module-versions/pkg/update"

	"github.com/sgreben/flagvar"
)

var (
	appName       = "terrafile-module-versions"
	version       = "2-SNAPSHOT"
	updatesClient = update.Client{
		Registry: registry.Client{
			HTTP: http.DefaultClient,
		},
	}
)

var config struct {
	Paths                   []string
	ModuleNames             flagvar.StringSet
	PrintVersionAndExit     bool
	Quiet                   bool
	Updates                 bool
	UpdatesFoundNonzeroExit bool
	Pretty                  bool
}

func init() {
	log.SetFlags(0)
	log.SetPrefix("[" + appName + "] ")
	log.SetOutput(os.Stderr)
	flag.BoolVar(&config.PrintVersionAndExit, "version", config.PrintVersionAndExit, "print version and exit")
	flag.BoolVar(&config.Updates, "updates", config.Updates, "check for updates")
	flag.BoolVar(&config.Updates, "update", config.Updates, "(alias for -updates)")
	flag.BoolVar(&config.Updates, "u", config.Updates, "(alias for -updates)")
	flag.BoolVar(&config.Quiet, "quiet", config.Quiet, "suppress log output (stderr)")
	flag.BoolVar(&config.Quiet, "q", config.Quiet, "(alias for -quiet)")
	flag.BoolVar(&config.Pretty, "pretty", config.Pretty, "human-readable output")
	flag.BoolVar(&config.Pretty, "p", config.Pretty, "(alias for -pretty)")
	flag.BoolVar(&config.UpdatesFoundNonzeroExit, "e", config.UpdatesFoundNonzeroExit, "(alias for -updates-found-nonzero-exit, implies -updates)")
	flag.BoolVar(&config.UpdatesFoundNonzeroExit, "updates-found-nonzero-exit", config.UpdatesFoundNonzeroExit, "exit with a nonzero code when modules with updates are found (implies -updates)")
	flag.Var(&config.ModuleNames, "module", "include this module (may be specified repeatedly. by default, all modules are included)")
	flag.Parse()

	if config.PrintVersionAndExit {
		fmt.Println(version)
		os.Exit(0)
	}

	if config.Quiet {
		log.SetOutput(ioutil.Discard)
	}

	if config.UpdatesFoundNonzeroExit {
		config.Updates = true
	}

	config.Paths = flag.Args()
}

func main() {
	if len(config.Paths) == 0 {
		config.Paths, _ = filepath.Glob("*.tf")
	}
	scanResults, err := scan.Scan(config.Paths)
	if err != nil {
		log.Fatal(err)
	}
	moduleNamesFilter := config.ModuleNames.Value
	moduleNamesFilterEmpty := len(moduleNamesFilter) == 0
	scanResultsFiltered := make([]scan.Result, 0, len(scanResults))
	for _, r := range scanResults {
		include := moduleNamesFilterEmpty || moduleNamesFilter[r.ModuleCall.Name]
		if !include {
			continue
		}
		scanResultsFiltered = append(scanResultsFiltered, r)
	}
	scanResults = scanResultsFiltered

	if !config.Updates {
		list(scanResults)
		return
	}
	updates(scanResults)
}

func list(scanResults []scan.Result) {
	var out output.Modules
	for _, m := range scanResults {
		parsed, err := modulecall.Parse(m.ModuleCall)
		if err != nil {
			log.Printf("error: %v", err)
			out = append(out, output.Module{
				Path:              m.Path,
				Name:              m.ModuleCall.Name,
				Source:            m.ModuleCall.Source,
				VersionConstraint: m.ModuleCall.Version,
			})
			continue
		}
		out = append(out, output.Module{
			Path:              m.Path,
			Name:              m.ModuleCall.Name,
			Source:            m.ModuleCall.Source,
			VersionConstraint: parsed.ConstraintsString,
			Version:           parsed.VersionString,
			Type:              parsed.Source.Type(),
		})
	}
	switch {
	case config.Pretty:
		out.WriteTable(os.Stdout)
	default:
		_ = out.WriteJSONL(os.Stdout)
	}
}

func updates(scanResults []scan.Result) {
	var (
		out                     output.Updates
		foundMatchingUpdates    bool
		foundNonMatchingUpdates bool
	)
	for _, m := range scanResults {
		parsed, err := modulecall.Parse(m.ModuleCall)
		if err != nil {
			log.Printf("error: %v", err)
			continue
		}
		update, err := updatesClient.Update(*parsed.Source, parsed.Version, parsed.Constraints)
		if err != nil {
			log.Printf("error: %v", err)
			continue
		}
		updateOutput := output.Update{
			Path:              m.Path,
			Name:              m.ModuleCall.Name,
			VersionConstraint: parsed.ConstraintsString,
			Version:           parsed.VersionString,
			LatestMatching:    update.LatestMatchingVersion,
			MatchingUpdate:    update.LatestMatchingUpdate != "",
			LatestOverall:     update.LatestOverallVersion,
			NonMatchingUpdate: update.LatestOverallUpdate != "" && update.LatestOverallUpdate != update.LatestMatchingVersion,
		}
		if updateOutput.MatchingUpdate {
			foundMatchingUpdates = true
		}
		if updateOutput.NonMatchingUpdate {
			foundNonMatchingUpdates = true
		}
		out = append(out, updateOutput)
	}
	switch {
	case config.Pretty:
		out.WriteTable(os.Stdout)
	default:
		_ = out.WriteJSONL(os.Stdout)
	}
	if config.UpdatesFoundNonzeroExit {
		if foundMatchingUpdates || foundNonMatchingUpdates { // not distinguishing between these for now
			os.Exit(1)
		}
	}
}
