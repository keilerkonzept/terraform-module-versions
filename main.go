package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/sgreben/versions/pkg/semver"

	"github.com/sgreben/flagvar"
)

var (
	appName                 = "terrafile-module-versions"
	version                 = "SNAPSHOT"
	terrafileVersionDefault = "master"
)

var config struct {
	Paths               []string
	ModuleNames         flagvar.StringSet
	PrintVersionAndExit bool
	Quiet               bool
	Updates             bool
}

func init() {
	log.SetFlags(0)
	log.SetPrefix("[" + appName + "] ")
	log.SetOutput(os.Stderr)
	flag.BoolVar(&config.PrintVersionAndExit, "version", config.PrintVersionAndExit, "print version and exit")
	flag.BoolVar(&config.Updates, "updates", config.Updates, "check for updates")
	flag.BoolVar(&config.Updates, "u", config.Updates, "(alias for -updates)")
	flag.BoolVar(&config.Quiet, "quiet", config.Quiet, "suppress log output (stderr)")
	flag.BoolVar(&config.Quiet, "q", config.Quiet, "(alias for -quiet)")
	flag.Var(&config.ModuleNames, "module", "include this module (may be specified repeatedly. by default, all modules are included)")
	flag.Parse()

	if config.PrintVersionAndExit {
		fmt.Println(version)
		os.Exit(0)
	}

	if config.Quiet {
		log.SetOutput(ioutil.Discard)
	}

	config.Paths = flag.Args()
}

func main() {
	var scanner scanner
	for _, path := range config.Paths {
		if err := scanner.ScanFile(path); err != nil {
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
		updates(included)
	default:
		list(included)
	}
}

func list(rs []*moduleReference) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	for _, r := range rs {
		src := r.SourceStruct()
		enc.Encode(outputList{
			Path:              r.Path,
			Name:              r.Name,
			Source:            r.Source,
			Version:           src.InferredVersion(),
			VersionConstraint: r.Version,
			Type:              src.Type(),
		})
	}
}

func updates(rs []*moduleReference) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	out := make(chan outputUpdates, 1)
	var wg sync.WaitGroup
	outputDone := make(chan bool)
	wg.Add(len(rs))
	for _, r := range rs {
		r := r
		go func() {
			defer wg.Done()
			updatesSingle(r, out)
		}()
	}
	go func() {
		for o := range out {
			enc.Encode(o)
		}
		outputDone <- true
	}()
	wg.Wait()
	close(out)
	<-outputDone
}

func updatesSingle(r *moduleReference, out chan outputUpdates) {
	src := r.SourceStruct()
	currentVersion := src.InferredVersion()
	var currentVersionStruct *semver.Version
	versionConstraint := "*"
	if currentVersion != nil {
		currentVersionStruct, _ = semver.Parse(*currentVersion)
	}
	if r.Version != nil {
		versionConstraint = *r.Version
	}
	versionConstraintStruct, err := semver.ParseConstraint(versionConstraint)
	if err != nil {
		log.Printf("%q: %v", versionConstraint, err)
		return
	}
	versions, err := src.Versions()
	if err != nil {
		log.Printf("fetch versions for %q: %v", r.Source, err)
		return
	}
	var versionsStrings []string
	var hasMajorUpdate, hasMinorUpdate, hasPatchUpdate bool
	var latest string
	var oldest *semver.Version
	for _, v := range versions {
		if !versionConstraintStruct.Check(v.VersionStruct) {
			continue
		}
		if oldest == nil {
			oldest = v.VersionStruct
		}
		compareAgainst := currentVersionStruct
		if compareAgainst == nil {
			compareAgainst = oldest
		}
		if !v.VersionStruct.GreaterThan(compareAgainst) {
			continue
		}
		if v.VersionStruct.Major > compareAgainst.Major {
			hasMajorUpdate = true
		}
		if v.VersionStruct.Minor > compareAgainst.Minor {
			hasMinorUpdate = true
		}
		if v.VersionStruct.Patch > compareAgainst.Patch {
			hasPatchUpdate = true
		}
		latest = v.Version
		versionsStrings = append(versionsStrings, v.Version)
	}
	if len(versionsStrings) == 0 {
		return
	}
	out <- outputUpdates{
		Path:              r.Path,
		Name:              r.Name,
		Source:            r.Source,
		Version:           currentVersion,
		VersionConstraint: r.Version,
		Type:              src.Type(),
		Updates:           versionsStrings,
		UpdateLatest:      latest,
		HasMajorUpdate:    hasMajorUpdate,
		HasMinorUpdate:    hasMinorUpdate,
		HasPatchUpdate:    hasPatchUpdate,
	}
}
