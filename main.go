package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/sgreben/versions/pkg/semver"

	"github.com/sgreben/flagvar"
)

var (
	appName                 = "terrafile-module-versions"
	version                 = "2-SNAPSHOT"
	terrafileVersionDefault = "master"
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
	flag.BoolVar(&config.UpdatesFoundNonzeroExit, "updates-found-nonzero-exit", config.UpdatesFoundNonzeroExit, "exit with a nonzero code when modules with upates are found (implies -updates)")
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
	var scanner scanner
	if len(config.Paths) == 0 {
		config.Paths, _ = filepath.Glob("*.tf")
	}
	for _, path := range config.Paths {
		if err := scanner.ScanDir(path); err != nil {
			log.Fatal(err)
		}
	}
	var included []*moduleReference
	moduleNamesFilter := config.ModuleNames.Value
	moduleNamesFilterEmpty := len(moduleNamesFilter) == 0
	for _, r := range scanner.Results {
		include := moduleNamesFilterEmpty || moduleNamesFilter[r.Name]
		if !include {
			continue
		}
		included = append(included, r)
	}

	switch {
	case config.Updates:
		if config.Pretty {
			updatesPretty(included)
			return
		}
		updatesJSON(included)
	default:
		if config.Pretty {
			listPretty(included)
			return
		}
		listJSON(included)
	}
}

func listPretty(rs []*moduleReference) {
	var out []outputList
	for _, r := range rs {
		src := r.SourceStruct()
		version := ""
		constraint := ""
		if v := src.InferredVersion(); v != nil {
			version = *v
		}
		if c := src.Version; c != nil {
			constraint = *src.Version
		}
		out = append(out, outputList{
			Path:              r.Path,
			Name:              r.Name,
			Source:            r.Source,
			Version:           version,
			VersionConstraint: constraint,
			Type:              src.Type(),
		})
	}
	listPrettyPrint(os.Stdout, out)
}

func listJSON(rs []*moduleReference) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	for _, r := range rs {
		src := r.SourceStruct()
		version := ""
		constraint := ""
		if v := src.InferredVersion(); v != nil {
			version = *v
		}
		if c := src.Version; c != nil {
			constraint = *src.Version
		}
		enc.Encode(outputList{
			Path:              r.Path,
			Name:              r.Name,
			Source:            r.Source,
			Version:           version,
			VersionConstraint: constraint,
			Type:              src.Type(),
		})
	}
}

func updatesJSON(rs []*moduleReference) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	out := make(chan outputUpdates, len(rs))
	var wg sync.WaitGroup
	outputDone := make(chan bool)
	wg.Add(len(rs))
	for _, r := range rs {
		r := r
		go func(r *moduleReference) {
			defer wg.Done()
			if err := updates(r, out); err != nil {
				log.Printf("%v: %v", r, err)
			}
		}(r)
	}
	var matchingUpdatesFound bool
	go func() {
		for o := range out {
			enc.Encode(o)
			if o.MatchingUpdate {
				matchingUpdatesFound = true
			}
		}
		outputDone <- true
	}()
	wg.Wait()
	close(out)
	<-outputDone
	if config.UpdatesFoundNonzeroExit && matchingUpdatesFound {
		os.Exit(1)
	}
}

func updatesPretty(rs []*moduleReference) {
	out := make(chan outputUpdates, len(rs))
	var wg sync.WaitGroup
	wg.Add(len(rs))
	for _, r := range rs {
		r := r
		go func(r *moduleReference) {
			defer wg.Done()
			if err := updates(r, out); err != nil {
				log.Printf("%v: %v", r, err)
			}
		}(r)
	}
	wg.Wait()
	close(out)
	var output []outputUpdates
	var matchingUpdatesFound bool
	for o := range out {
		output = append(output, o)
		if o.MatchingUpdate {
			matchingUpdatesFound = true
		}
	}
	updatePrettyPrint(os.Stdout, output)
	if config.UpdatesFoundNonzeroExit && matchingUpdatesFound {
		os.Exit(1)
	}
}

func updates(r *moduleReference, out chan outputUpdates) error {
	src := r.SourceStruct()
	currentVersionString := src.InferredVersion()
	var currentVersion *semver.Version
	var versionConstraintString string
	if currentVersionString != nil {
		var err error
		currentVersion, err = semver.Parse(*currentVersionString)
		if err != nil {
			return fmt.Errorf("parse version %q: %v", *currentVersionString, err)
		}
	}
	haveConstraint := r.Version != nil && *r.Version != ""
	haveCurrentVersion := currentVersion != nil
	switch {
	case !haveConstraint && !haveCurrentVersion:
		versionConstraintString = "*"
	case !haveConstraint && haveCurrentVersion:
		versionConstraintString = fmt.Sprintf(">=%s", currentVersion.String())
	case haveConstraint && !haveCurrentVersion:
		versionConstraintString = *r.Version
	case haveConstraint && haveCurrentVersion:
		versionConstraintString = fmt.Sprintf("%s,>=%s", *r.Version, currentVersion.String())
	}
	versionConstraint, err := semver.ParseConstraint(versionConstraintString)
	if err != nil {
		return fmt.Errorf("parse version constraint %q: %v", versionConstraintString, err)
	}
	versions, err := src.Versions()
	if err != nil {
		return fmt.Errorf("fetch versions: %v", err)
	}
	var versionsCollection semver.Collection
	for _, v := range versions {
		versionsCollection = append(versionsCollection, v.Version)
	}
	var matchingUpdate bool
	var constraintUpdate bool
	latest := versionConstraint.LatestMatching(versionsCollection)
	var latestString string
	if latest != nil {
		latestString = latest.Original
		if currentVersion != nil {
			matchingUpdate = latest.GreaterThan(currentVersion)
		}
		oldest := versionConstraint.OldestMatching(versionsCollection)
		if latest.GreaterThan(oldest) {
			constraintUpdate = true
		}
	}
	var latestOverallString string
	var nonMatchingUpdate bool
	if len(versionsCollection) > 0 {
		sort.Sort(versionsCollection)
		latestOverall := versionsCollection[len(versionsCollection)-1]
		latestOverallString = latestOverall.Original
		if !versionConstraint.Check(latestOverall) {
			nonMatchingUpdate = true
		}
	}
	version := ""
	constraint := ""
	if v := src.InferredVersion(); v != nil {
		version = *v
	}
	if c := src.Version; c != nil {
		constraint = *src.Version
	}
	out <- outputUpdates{
		Path:                    r.Path,
		Name:                    r.Name,
		Source:                  r.Source,
		Version:                 version,
		VersionConstraint:       constraint,
		ConstraintUpdate:        constraintUpdate,
		Latest:                  latestString,
		MatchingUpdate:          matchingUpdate,
		LatestWithoutConstraint: latestOverallString,
		NonMatchingUpdate:       nonMatchingUpdate,
	}
	return nil
}
